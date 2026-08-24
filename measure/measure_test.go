package measure

import (
	"testing"
	"time"
)

func TestEpochIsolationChangesIdentity(t *testing.T) {
	v1, err := NewEpoch(EpochDefinition{Construct: "accuracy", Instrument: "fake-judge", Protocol: "judge/v1", Implementation: "fake/v1"})
	if err != nil {
		t.Fatal(err)
	}
	v2, err := NewEpoch(EpochDefinition{Construct: "accuracy", Instrument: "fake-judge", Protocol: "judge/v2", Implementation: "fake/v1"})
	if err != nil {
		t.Fatal(err)
	}
	if v1.ID == v2.ID {
		t.Fatal("protocol change did not alter epoch ID")
	}
	obs1, err := NewObservation(Draft{Role: RoleMeasurement, Subject: SubjectRef{Kind: "episode", ID: "episode:test"}, Construct: "accuracy", Instrument: "fake-judge", Protocol: "judge/v1", Epoch: v1, Status: StatusMeasured, Value: NumberValue("0.9")}, time.Unix(10, 0))
	if err != nil {
		t.Fatal(err)
	}
	obs2, err := NewObservation(Draft{Role: RoleMeasurement, Subject: SubjectRef{Kind: "episode", ID: "episode:test"}, Construct: "accuracy", Instrument: "fake-judge", Protocol: "judge/v2", Epoch: v2, Status: StatusMeasured, Value: NumberValue("0.9")}, time.Unix(10, 0))
	if err != nil {
		t.Fatal(err)
	}
	if obs1.ID == obs2.ID {
		t.Fatal("epoch change did not alter observation identity")
	}
}
