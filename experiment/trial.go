package experiment

import (
	"fmt"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type Case struct {
	ID       string            `json:"id"`
	Input    artifact.Ref      `json:"input"`
	Groups   []string          `json:"groups,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type DatasetManifest struct {
	ID     string        `json:"id"`
	Role   string        `json:"role"`
	Cases  []Case        `json:"cases"`
	Digest record.Digest `json:"digest"`
}

func cloneCase(value Case) Case {
	cloned := value
	if value.Input.Schema != nil {
		schema := *value.Input.Schema
		cloned.Input.Schema = &schema
	}
	cloned.Groups = append([]string(nil), value.Groups...)
	if value.Metadata != nil {
		cloned.Metadata = make(map[string]string, len(value.Metadata))
		for key, item := range value.Metadata {
			cloned.Metadata[key] = item
		}
	}
	return cloned
}

func cloneCases(values []Case) []Case {
	cloned := make([]Case, len(values))
	for index, value := range values {
		cloned[index] = cloneCase(value)
	}
	return cloned
}

func cloneDatasetManifest(value DatasetManifest) DatasetManifest {
	value.Cases = cloneCases(value.Cases)
	return value
}

func NewDatasetManifest(role string, cases []Case) (DatasetManifest, error) {
	if role == "" || len(cases) == 0 {
		return DatasetManifest{}, fmt.Errorf("dataset role and cases are required")
	}
	seen := make(map[string]struct{}, len(cases))
	for _, c := range cases {
		if c.ID == "" {
			return DatasetManifest{}, fmt.Errorf("case ID is required")
		}
		if _, ok := seen[c.ID]; ok {
			return DatasetManifest{}, fmt.Errorf("duplicate case %s", c.ID)
		}
		seen[c.ID] = struct{}{}
	}
	immutableCases := cloneCases(cases)
	identity := struct {
		Role  string `json:"role"`
		Cases []Case `json:"cases"`
	}{Role: role, Cases: immutableCases}
	digest, _, err := record.SemanticDigest("schema:optkit.dataset/v1", identity)
	if err != nil {
		return DatasetManifest{}, err
	}
	rawID, err := record.ContentID("dataset", digest)
	if err != nil {
		return DatasetManifest{}, err
	}
	return DatasetManifest{ID: rawID, Role: role, Cases: immutableCases, Digest: digest}, nil
}

type Arm struct {
	ID       string            `json:"id"`
	Snapshot record.SnapshotID `json:"snapshot"`
}

type TrialPlan struct {
	ID       record.TrialID  `json:"id"`
	Arms     []Arm           `json:"arms"`
	Dataset  DatasetManifest `json:"dataset"`
	Repeats  int             `json:"repeats"`
	Protocol string          `json:"protocol"`
}

func NewCompleteBlockTrial(arms []Arm, dataset DatasetManifest, repeats int, protocol string) (TrialPlan, error) {
	if len(arms) < 2 {
		return TrialPlan{}, fmt.Errorf("complete-block trial requires at least two arms")
	}
	if repeats < 1 {
		return TrialPlan{}, fmt.Errorf("repeats must be positive")
	}
	if protocol == "" {
		return TrialPlan{}, fmt.Errorf("execution protocol is required")
	}
	seen := make(map[string]struct{}, len(arms))
	for _, arm := range arms {
		if arm.ID == "" || arm.Snapshot == "" {
			return TrialPlan{}, fmt.Errorf("arm ID and snapshot are required")
		}
		if _, ok := seen[arm.ID]; ok {
			return TrialPlan{}, fmt.Errorf("duplicate arm %s", arm.ID)
		}
		seen[arm.ID] = struct{}{}
	}
	identity := struct {
		Design   string        `json:"design"`
		Arms     []Arm         `json:"arms"`
		Dataset  record.Digest `json:"dataset"`
		Repeats  int           `json:"repeats"`
		Protocol string        `json:"protocol"`
	}{Design: "complete_block", Arms: arms, Dataset: dataset.Digest, Repeats: repeats, Protocol: protocol}
	digest, _, err := record.SemanticDigest("schema:optkit.trial/v1", identity)
	if err != nil {
		return TrialPlan{}, err
	}
	rawID, err := record.ContentID("trial", digest)
	if err != nil {
		return TrialPlan{}, err
	}
	return TrialPlan{ID: record.TrialID(rawID), Arms: append([]Arm(nil), arms...), Dataset: cloneDatasetManifest(dataset), Repeats: repeats, Protocol: protocol}, nil
}

type EpisodeSpec struct {
	ID          record.EpisodeID `json:"id"`
	SemanticKey record.Digest    `json:"semantic_key"`
	Trial       record.TrialID   `json:"trial"`
	Arm         Arm              `json:"arm"`
	Case        Case             `json:"case"`
	Repeat      int              `json:"repeat"`
	Seed        int64            `json:"seed"`
	Protocol    string           `json:"protocol"`
}

func Expand(plan TrialPlan) ([]EpisodeSpec, error) {
	out := make([]EpisodeSpec, 0, len(plan.Arms)*len(plan.Dataset.Cases)*plan.Repeats)
	for _, c := range plan.Dataset.Cases {
		for repeat := 0; repeat < plan.Repeats; repeat++ {
			for _, arm := range plan.Arms {
				identity := struct {
					Trial    record.TrialID    `json:"trial"`
					Snapshot record.SnapshotID `json:"snapshot"`
					Case     string            `json:"case"`
					Input    record.Digest     `json:"input"`
					Repeat   int               `json:"repeat"`
					Seed     int64             `json:"seed"`
					Protocol string            `json:"protocol"`
				}{
					Trial: plan.ID, Snapshot: arm.Snapshot, Case: c.ID, Input: c.Input.Digest,
					Repeat: repeat, Seed: deterministicSeed(plan.ID, c.ID, repeat), Protocol: plan.Protocol,
				}
				digest, _, err := record.SemanticDigest("schema:optkit.episode-spec/v1", identity)
				if err != nil {
					return nil, err
				}
				rawID, err := record.ContentID("episode", digest)
				if err != nil {
					return nil, err
				}
				out = append(out, EpisodeSpec{
					ID: record.EpisodeID(rawID), SemanticKey: digest, Trial: plan.ID,
					Arm: arm, Case: c, Repeat: repeat, Seed: identity.Seed, Protocol: plan.Protocol,
				})
			}
		}
	}
	return out, nil
}

func deterministicSeed(trial record.TrialID, caseID string, repeat int) int64 {
	digest := record.DigestParts([]byte(trial), []byte{0}, []byte(caseID), []byte{0}, []byte(fmt.Sprintf("%d", repeat)))
	hexPart, _ := digest.Hex()
	var seed int64
	for i := 0; i < 16 && i < len(hexPart); i++ {
		seed = seed*16 + int64(hexDigit(hexPart[i]))
	}
	return seed & 0x7fffffffffffffff
}

func hexDigit(b byte) int {
	if b >= '0' && b <= '9' {
		return int(b - '0')
	}
	return int(b-'a') + 10
}
