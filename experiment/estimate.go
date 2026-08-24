package experiment

import (
	"fmt"
	"math"
	"sort"
)

type NumericObservation struct {
	CaseID string
	ArmID  string
	Repeat int
	Value  float64
	Valid  bool
}

type Estimate struct {
	Baseline   string    `json:"baseline"`
	Treatment  string    `json:"treatment"`
	Value      float64   `json:"value"`
	SampleSize int       `json:"sample_size"`
	Missing    int       `json:"missing"`
	Deltas     []float64 `json:"deltas"`
}

func PairedMean(observations []NumericObservation, baseline, treatment string) (Estimate, error) {
	type key struct {
		caseID string
		repeat int
	}
	pairs := make(map[key]map[string]NumericObservation)
	for _, obs := range observations {
		k := key{caseID: obs.CaseID, repeat: obs.Repeat}
		if pairs[k] == nil {
			pairs[k] = make(map[string]NumericObservation)
		}
		if _, exists := pairs[k][obs.ArmID]; exists {
			return Estimate{}, fmt.Errorf("duplicate observation for case %s repeat %d arm %s", obs.CaseID, obs.Repeat, obs.ArmID)
		}
		pairs[k][obs.ArmID] = obs
	}
	keys := make([]key, 0, len(pairs))
	for k := range pairs {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].caseID == keys[j].caseID {
			return keys[i].repeat < keys[j].repeat
		}
		return keys[i].caseID < keys[j].caseID
	})

	estimate := Estimate{Baseline: baseline, Treatment: treatment}
	for _, k := range keys {
		base, baseOK := pairs[k][baseline]
		treat, treatOK := pairs[k][treatment]
		if !baseOK || !treatOK || !base.Valid || !treat.Valid || math.IsNaN(base.Value) || math.IsNaN(treat.Value) {
			estimate.Missing++
			continue
		}
		estimate.Deltas = append(estimate.Deltas, treat.Value-base.Value)
	}
	if estimate.Missing > 0 {
		return estimate, fmt.Errorf("paired estimate invalid: %d missing or invalid pair(s)", estimate.Missing)
	}
	if len(estimate.Deltas) == 0 {
		return estimate, fmt.Errorf("paired estimate has no complete pairs")
	}
	for _, delta := range estimate.Deltas {
		estimate.Value += delta
	}
	estimate.SampleSize = len(estimate.Deltas)
	estimate.Value /= float64(estimate.SampleSize)
	return estimate, nil
}
