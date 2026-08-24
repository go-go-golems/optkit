package episode

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now().UTC() }

type Writer struct {
	mu      sync.Mutex
	episode record.EpisodeID
	store   artifact.Store
	clock   Clock
	events  []Event
	spans   map[record.SpanID]struct{}
	sealed  bool
}

func NewWriter(episodeID record.EpisodeID, store artifact.Store, clock Clock) (*Writer, error) {
	if err := record.ValidateID("episode", string(episodeID)); err != nil {
		return nil, err
	}
	if store == nil {
		return nil, fmt.Errorf("episode writer requires an artifact store")
	}
	if clock == nil {
		clock = RealClock{}
	}
	return &Writer{episode: episodeID, store: store, clock: clock, spans: make(map[record.SpanID]struct{})}, nil
}

func (w *Writer) AttachJSON(ctx context.Context, schema record.SchemaID, sensitivity artifact.Sensitivity, value any) (artifact.Ref, error) {
	return artifact.PutCanonical(ctx, w.store, schema, sensitivity, value)
}

func (w *Writer) Emit(ctx context.Context, emission Emission) (Event, error) {
	if err := ctx.Err(); err != nil {
		return Event{}, err
	}
	if emission.Kind == "" {
		return Event{}, fmt.Errorf("event kind is required")
	}
	if err := emission.Schema.Validate(); err != nil {
		return Event{}, err
	}
	if err := record.ValidateID("span", string(emission.Span)); err != nil {
		return Event{}, err
	}
	payload, err := w.AttachJSON(ctx, emission.Schema, emission.Sensitivity, emission.Payload)
	if err != nil {
		return Event{}, err
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	if w.sealed {
		return Event{}, fmt.Errorf("trajectory is already sealed")
	}
	if emission.Parent != nil {
		if _, ok := w.spans[*emission.Parent]; !ok {
			return Event{}, fmt.Errorf("parent span %s has not been observed", *emission.Parent)
		}
	}
	w.spans[emission.Span] = struct{}{}
	seq := uint64(len(w.events) + 1)
	eventIdentity := struct {
		Episode record.EpisodeID `json:"episode"`
		Seq     uint64           `json:"seq"`
		Time    time.Time        `json:"time"`
		Kind    EventKind        `json:"kind"`
		Schema  record.SchemaID  `json:"schema"`
		Span    record.SpanID    `json:"span"`
		Parent  *record.SpanID   `json:"parent,omitempty"`
		Payload record.Digest    `json:"payload"`
		Tags    [][2]string      `json:"tags,omitempty"`
	}{
		Episode: w.episode,
		Seq:     seq,
		Time:    w.clock.Now().UTC(),
		Kind:    emission.Kind,
		Schema:  emission.Schema,
		Span:    emission.Span,
		Parent:  emission.Parent,
		Payload: payload.Digest,
		Tags:    sortedTags(emission.Tags),
	}
	digest, _, err := record.SemanticDigest("schema:optkit.trajectory-event-identity/v1", eventIdentity)
	if err != nil {
		return Event{}, err
	}
	rawID, err := record.ContentID("trajectory-event", digest)
	if err != nil {
		return Event{}, err
	}
	event := Event{
		ID:        record.EventID(rawID),
		Episode:   w.episode,
		Seq:       seq,
		Time:      eventIdentity.Time,
		Kind:      emission.Kind,
		Schema:    emission.Schema,
		Span:      emission.Span,
		Parent:    emission.Parent,
		Payload:   payload,
		Tags:      cloneTags(emission.Tags),
		Sensitive: emission.Sensitivity == artifact.SensitivityConfidential || emission.Sensitivity == artifact.SensitivityRestricted,
	}
	w.events = append(w.events, event)
	return event, nil
}

type Trajectory struct {
	Episode record.EpisodeID `json:"episode"`
	Events  []Event          `json:"events"`
}

func (w *Writer) Seal(ctx context.Context) (artifact.Ref, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.sealed {
		return artifact.Ref{}, fmt.Errorf("trajectory is already sealed")
	}
	if len(w.events) == 0 {
		return artifact.Ref{}, fmt.Errorf("cannot seal an empty trajectory")
	}
	for i, event := range w.events {
		if event.Seq != uint64(i+1) {
			return artifact.Ref{}, fmt.Errorf("event sequence at index %d is %d", i, event.Seq)
		}
	}
	ref, err := artifact.PutCanonical(ctx, w.store, "schema:optkit.trajectory/v1", artifact.SensitivityInternal, Trajectory{
		Episode: w.episode,
		Events:  append([]Event(nil), w.events...),
	})
	if err != nil {
		return artifact.Ref{}, err
	}
	w.sealed = true
	return ref, nil
}

func LoadTrajectory(ctx context.Context, store artifact.Store, ref artifact.Ref) (Trajectory, error) {
	var trajectory Trajectory
	if err := artifact.DecodeJSON(ctx, store, ref, &trajectory); err != nil {
		return Trajectory{}, err
	}
	for i, event := range trajectory.Events {
		if event.Episode != trajectory.Episode {
			return Trajectory{}, fmt.Errorf("event %s belongs to episode %s, trajectory is %s", event.ID, event.Episode, trajectory.Episode)
		}
		if event.Seq != uint64(i+1) {
			return Trajectory{}, fmt.Errorf("trajectory sequence mismatch at index %d", i)
		}
		if err := store.Verify(ctx, event.Payload); err != nil {
			return Trajectory{}, fmt.Errorf("verify event %s payload: %w", event.ID, err)
		}
	}
	return trajectory, nil
}

func sortedTags(tags map[string]string) [][2]string {
	keys := make([]string, 0, len(tags))
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([][2]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, [2]string{key, tags[key]})
	}
	return out
}

func cloneTags(tags map[string]string) map[string]string {
	if len(tags) == 0 {
		return nil
	}
	out := make(map[string]string, len(tags))
	for key, value := range tags {
		out[key] = value
	}
	return out
}
