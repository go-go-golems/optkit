package numbergame

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/episode"
	"github.com/go-go-golems/optkit/measure"
)

type Assessment struct {
	Instrument string `json:"instrument"`
	Protocol   string `json:"protocol"`
	Score      string `json:"score"`
	Label      string `json:"label"`
	Rationale  string `json:"rationale"`
}

type Measurements struct {
	Intervention  measure.Observation
	AbsoluteError measure.Observation
	Judge         measure.Observation
	Assessment    artifact.Ref
	Score         float64
}

func Measure(ctx context.Context, store artifact.Store, result episode.Result, expected Config, repeat int, createdAt time.Time) (Measurements, error) {
	if result.Output == nil {
		return Measurements{}, fmt.Errorf("numbergame result has no output")
	}
	var output Output
	if err := artifact.DecodeJSON(ctx, store, *result.Output, &output); err != nil {
		return Measurements{}, err
	}
	trajectory, err := episode.LoadTrajectory(ctx, store, result.Trajectory)
	if err != nil {
		return Measurements{}, err
	}
	var observed Config
	foundConfig := false
	for _, event := range trajectory.Events {
		if event.Kind != "numbergame.configured" {
			continue
		}
		if err := artifact.DecodeJSON(ctx, store, event.Payload, &observed); err != nil {
			return Measurements{}, err
		}
		foundConfig = true
		break
	}
	exercised := foundConfig && observed == expected
	interventionEpoch, err := measure.NewEpoch(measure.EpochDefinition{
		Construct: "intervention.numbergame.config", Instrument: "numbergame.config-probe",
		Protocol: "numbergame.config-probe/v1", Implementation: "numbergame/v1",
	})
	if err != nil {
		return Measurements{}, err
	}
	interventionStatus := measure.StatusViolated
	if exercised {
		interventionStatus = measure.StatusSatisfied
	}
	intervention, err := measure.NewObservation(measure.Draft{
		Role: measure.RoleIntervention, Subject: measure.SubjectRef{Kind: "episode", ID: string(trajectory.Episode)},
		Construct: interventionEpoch.Definition.Construct, Instrument: interventionEpoch.Definition.Instrument,
		Protocol: interventionEpoch.Definition.Protocol, Epoch: interventionEpoch, Status: interventionStatus,
		Value: measure.BoolValue(exercised), Evidence: []artifact.Ref{result.Trajectory},
		Diagnostics: map[string]any{"expected": expected, "observed": observed, "found": foundConfig}, Repeat: repeat,
	}, createdAt)
	if err != nil {
		return Measurements{}, err
	}

	errorEpoch, err := measure.NewEpoch(measure.EpochDefinition{
		Construct: "target.absolute_error", Instrument: "numbergame.absolute-error",
		Protocol: "numbergame.absolute-error/v1", Implementation: "numbergame/v1",
	})
	if err != nil {
		return Measurements{}, err
	}
	absoluteError, err := measure.NewObservation(measure.Draft{
		Role: measure.RoleMeasurement, Subject: measure.SubjectRef{Kind: "episode", ID: string(trajectory.Episode)},
		Construct: errorEpoch.Definition.Construct, Instrument: errorEpoch.Definition.Instrument,
		Protocol: errorEpoch.Definition.Protocol, Epoch: errorEpoch, Status: measure.StatusMeasured,
		Value: measure.NumberValue(strconv.Itoa(output.AbsoluteError)), Evidence: []artifact.Ref{result.Trajectory, *result.Output},
		Diagnostics: map[string]any{"target": output.Target, "value": output.Value}, Repeat: repeat,
	}, createdAt)
	if err != nil {
		return Measurements{}, err
	}

	score := 1.0 / (1.0 + float64(output.AbsoluteError))
	label := "needs_improvement"
	if output.AbsoluteError == 0 {
		label = "exact"
	}
	assessment := Assessment{
		Instrument: "numbergame.fake-judge", Protocol: "numbergame.fake-judge/v1",
		Score: strconv.FormatFloat(score, 'f', 6, 64), Label: label,
		Rationale: fmt.Sprintf("Deterministic fake judge: absolute error is %d.", output.AbsoluteError),
	}
	assessmentRef, err := artifact.PutCanonical(ctx, store, "schema:numbergame.fake-judge-assessment/v1", artifact.SensitivityInternal, assessment)
	if err != nil {
		return Measurements{}, err
	}
	judgeEpoch, err := measure.NewEpoch(measure.EpochDefinition{
		Construct: "target.accuracy", Instrument: "numbergame.fake-judge",
		Protocol: "numbergame.fake-judge/v1", Implementation: "deterministic-fake/v1",
	})
	if err != nil {
		return Measurements{}, err
	}
	judge, err := measure.NewObservation(measure.Draft{
		Role: measure.RoleMeasurement, Subject: measure.SubjectRef{Kind: "episode", ID: string(trajectory.Episode)},
		Construct: judgeEpoch.Definition.Construct, Instrument: judgeEpoch.Definition.Instrument,
		Protocol: judgeEpoch.Definition.Protocol, Epoch: judgeEpoch, Status: measure.StatusMeasured,
		Value:       measure.NumberValue(strconv.FormatFloat(score, 'f', 6, 64)),
		Evidence:    []artifact.Ref{result.Trajectory, *result.Output, assessmentRef},
		Diagnostics: map[string]any{"label": label}, Repeat: repeat,
	}, createdAt)
	if err != nil {
		return Measurements{}, err
	}
	return Measurements{Intervention: intervention, AbsoluteError: absoluteError, Judge: judge, Assessment: assessmentRef, Score: score}, nil
}
