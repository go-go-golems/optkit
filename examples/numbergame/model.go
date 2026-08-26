package numbergame

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-go-golems/optkit/record"
	"github.com/go-go-golems/optkit/space"
)

type Noise string

const (
	NoiseNone  Noise = "none"
	NoiseSmall Noise = "small"
	NoiseLarge Noise = "large"
)

type Config struct {
	Multiplier int   `json:"multiplier"`
	Noise      Noise `json:"noise"`
}

func (c Config) Validate() error {
	if c.Multiplier < 1 || c.Multiplier > 10 {
		return fmt.Errorf("multiplier %d is outside [1,10]", c.Multiplier)
	}
	switch c.Noise {
	case NoiseNone, NoiseSmall, NoiseLarge:
		return nil
	default:
		return fmt.Errorf("unknown noise mode %q", c.Noise)
	}
}

type Case struct {
	Input  int `json:"input"`
	Target int `json:"target"`
}

type Output struct {
	Input         int   `json:"input"`
	Multiplier    int   `json:"multiplier"`
	NoiseMode     Noise `json:"noise_mode"`
	NoiseValue    int   `json:"noise_value"`
	Value         int   `json:"value"`
	Target        int   `json:"target"`
	AbsoluteError int   `json:"absolute_error"`
}

func ConfigCodec() space.JSONCodec[Config] {
	return space.NewJSONCodec[Config]("schema:numbergame.config/v1")
}

func mustVariable[C, V any](variable space.Variable[C, V], err error) space.Variable[C, V] {
	if err != nil {
		panic(err)
	}
	return variable
}

func MultiplierVariable() space.Variable[Config, int] {
	codec := space.NewJSONCodec[int]("schema:numbergame.multiplier/v1")
	domain := space.IntRange(1, 10)
	defaultValue := 2
	return mustVariable(space.NewVariable(
		space.VariableMetadata{
			ID: "math.multiplier", Key: "multiplier", Label: "Multiplier",
			Short:    "Integer multiplied by every input.",
			Long:     "Controls the deterministic arithmetic multiplier applied before optional seeded noise.",
			CostHint: space.CostLow, BindingVersion: "numbergame.multiplier/v1",
			Probes: []string{"numbergame.config.exercised/v1"},
		},
		space.Lens[Config, int]{
			Get: func(c Config) int { return c.Multiplier },
			Put: func(c Config, value int) (Config, error) {
				c.Multiplier = value
				return c, c.Validate()
			},
		},
		domain, codec, &defaultValue,
	))
}

func NoiseVariable() space.Variable[Config, Noise] {
	codec := space.NewJSONCodec[Noise]("schema:numbergame.noise/v1")
	domain := space.Choices(map[Noise]string{
		NoiseNone: "No noise", NoiseSmall: "Small seeded noise", NoiseLarge: "Large seeded noise",
	})
	defaultValue := NoiseNone
	return mustVariable(space.NewVariable(
		space.VariableMetadata{
			ID: "noise.mode", Key: "mode", Label: "Noise mode",
			Short:    "Seeded pseudo-random perturbation applied after multiplication.",
			Long:     "Selects the deterministic seeded-noise profile used to perturb multiplied values.",
			CostHint: space.CostLow, BindingVersion: "numbergame.noise/v1",
			Probes: []string{"numbergame.config.exercised/v1"},
		},
		space.Lens[Config, Noise]{
			Get: func(c Config) Noise { return c.Noise },
			Put: func(c Config, value Noise) (Config, error) {
				c.Noise = value
				return c, c.Validate()
			},
		},
		domain, codec, &defaultValue,
	))
}

// SequenceClock supplies deterministic, monotonically increasing times for
// integration tests and replayable examples.
type SequenceClock struct {
	mu      sync.Mutex
	current time.Time
	step    time.Duration
}

func NewSequenceClock(start time.Time, step time.Duration) *SequenceClock {
	if step <= 0 {
		step = time.Millisecond
	}
	return &SequenceClock{current: start.UTC(), step: step}
}

func (c *SequenceClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	value := c.current
	c.current = c.current.Add(c.step)
	return value
}

var SystemID = record.SystemID("system:numbergame/v1")
