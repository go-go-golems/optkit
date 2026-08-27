package episode

import (
	"context"
	"testing"
	"time"

	"github.com/go-go-golems/optkit/artifact/memory"
	"github.com/go-go-golems/optkit/record"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

func TestTrajectorySealAndLoad(t *testing.T) {
	ctx := context.Background()
	store := memory.New()
	writer, err := NewWriter("episode:test", store, fixedClock{now: time.Unix(100, 0).UTC()})
	if err != nil {
		t.Fatal(err)
	}
	root := record.SpanID("span:root")
	if _, err := writer.Emit(ctx, Emission{Kind: "input.received", Schema: "schema:test.input/v1", Span: root, Payload: map[string]any{"n": 3}}); err != nil {
		t.Fatal(err)
	}
	child := record.SpanID("span:calc")
	if _, err := writer.Emit(ctx, Emission{Kind: "output.produced", Schema: "schema:test.output/v1", Span: child, Parent: &root, Payload: map[string]any{"n": 9}}); err != nil {
		t.Fatal(err)
	}
	ref, err := writer.Seal(ctx)
	if err != nil {
		t.Fatal(err)
	}
	trajectory, err := LoadTrajectory(ctx, store, ref)
	if err != nil {
		t.Fatal(err)
	}
	if len(trajectory.Events) != 2 || trajectory.Events[1].Seq != 2 {
		t.Fatalf("unexpected trajectory: %+v", trajectory)
	}
	if _, err := writer.Emit(ctx, Emission{Kind: "late", Schema: "schema:test.output/v1", Span: child, Payload: map[string]any{}}); err == nil {
		t.Fatal("emit after seal unexpectedly succeeded")
	}
}

func TestUnknownParentRejected(t *testing.T) {
	writer, err := NewWriter("episode:test", memory.New(), fixedClock{now: time.Unix(100, 0).UTC()})
	if err != nil {
		t.Fatal(err)
	}
	parent := record.SpanID("span:missing")
	_, err = writer.Emit(context.Background(), Emission{Kind: "x", Schema: "schema:test/v1", Span: "span:child", Parent: &parent, Payload: map[string]any{}})
	if err == nil {
		t.Fatal("unknown parent unexpectedly accepted")
	}
}
