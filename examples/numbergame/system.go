package numbergame

import (
	"context"
	"fmt"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/episode"
	"github.com/go-go-golems/optkit/record"
)

type Invocation struct {
	Case Case  `json:"case"`
	Seed int64 `json:"seed"`
}

type Executable struct {
	Config Config
	Clock  episode.Clock
}

func (e Executable) Run(ctx context.Context, invocation Invocation, sink episode.Sink) (episode.RunResult, error) {
	if err := e.Config.Validate(); err != nil {
		return episode.RunResult{}, err
	}
	clock := e.Clock
	if clock == nil {
		clock = episode.RealClock{}
	}
	started := clock.Now().UTC()
	root := record.SpanID("span:numbergame-root")
	if _, err := sink.Emit(ctx, episode.Emission{
		Kind: "input.received", Schema: "schema:numbergame.case/v1", Span: root,
		Payload: invocation.Case, Sensitivity: artifact.SensitivityInternal,
	}); err != nil {
		return episode.RunResult{}, err
	}
	calculation := record.SpanID("span:numbergame-calculation")
	if _, err := sink.Emit(ctx, episode.Emission{
		Kind: "numbergame.configured", Schema: "schema:numbergame.config/v1", Span: calculation, Parent: &root,
		Payload: e.Config, Tags: map[string]string{"role": "effective-config"}, Sensitivity: artifact.SensitivityInternal,
	}); err != nil {
		return episode.RunResult{}, err
	}
	noise := sampleNoise(invocation.Seed, e.Config.Noise)
	if _, err := sink.Emit(ctx, episode.Emission{
		Kind: "numbergame.noise.sampled", Schema: "schema:numbergame.noise-sample/v1", Span: calculation, Parent: &root,
		Payload: map[string]any{"seed": invocation.Seed, "mode": e.Config.Noise, "value": noise}, Sensitivity: artifact.SensitivityInternal,
	}); err != nil {
		return episode.RunResult{}, err
	}
	value := e.Config.Multiplier*invocation.Case.Input + noise
	absoluteError := value - invocation.Case.Target
	if absoluteError < 0 {
		absoluteError = -absoluteError
	}
	output := Output{
		Input: invocation.Case.Input, Multiplier: e.Config.Multiplier,
		NoiseMode: e.Config.Noise, NoiseValue: noise, Value: value,
		Target: invocation.Case.Target, AbsoluteError: absoluteError,
	}
	outputRef, err := sink.AttachJSON(ctx, "schema:numbergame.output/v1", artifact.SensitivityInternal, output)
	if err != nil {
		return episode.RunResult{}, err
	}
	if _, err := sink.Emit(ctx, episode.Emission{
		Kind: "output.produced", Schema: "schema:numbergame.output-event/v1", Span: root,
		Payload: map[string]any{"output": outputRef}, Sensitivity: artifact.SensitivityInternal,
	}); err != nil {
		return episode.RunResult{}, err
	}
	if _, err := sink.Emit(ctx, episode.Emission{
		Kind: "episode.completed", Schema: "schema:numbergame.completed/v1", Span: root,
		Payload: map[string]any{"status": "completed"}, Sensitivity: artifact.SensitivityInternal,
	}); err != nil {
		return episode.RunResult{}, err
	}
	finished := clock.Now().UTC()
	if finished.Before(started) {
		return episode.RunResult{}, fmt.Errorf("clock moved backwards")
	}
	return episode.RunResult{
		Status: episode.StatusCompleted, Output: &outputRef,
		Usage: []episode.ResourceUsage{
			{Resource: "episodes", Units: 1},
			{Resource: "numbergame.operations", Units: 1},
		},
		StartedAt: started, FinishedAt: finished,
	}, nil
}

func sampleNoise(seed int64, mode Noise) int {
	if mode == NoiseNone {
		return 0
	}
	x := uint64(seed) + 0x9e3779b97f4a7c15
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	x *= 0x2545f4914f6cdd1d
	if mode == NoiseSmall {
		return int(x%3) - 1
	}
	return int(x%7) - 3
}
