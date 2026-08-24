package episode

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type EventKind string

type Event struct {
	ID        record.EventID    `json:"id"`
	Episode   record.EpisodeID  `json:"episode"`
	Seq       uint64            `json:"seq"`
	Time      time.Time         `json:"time"`
	Kind      EventKind         `json:"kind"`
	Schema    record.SchemaID   `json:"schema"`
	Span      record.SpanID     `json:"span"`
	Parent    *record.SpanID    `json:"parent,omitempty"`
	Payload   artifact.Ref      `json:"payload"`
	Tags      map[string]string `json:"tags,omitempty"`
	Sensitive bool              `json:"sensitive"`
}

type Emission struct {
	Kind        EventKind
	Schema      record.SchemaID
	Span        record.SpanID
	Parent      *record.SpanID
	Payload     any
	Tags        map[string]string
	Sensitivity artifact.Sensitivity
}

type Sink interface {
	Emit(context.Context, Emission) (Event, error)
	AttachJSON(context.Context, record.SchemaID, artifact.Sensitivity, any) (artifact.Ref, error)
}

type Executable[I any] interface {
	Run(context.Context, I, Sink) (RunResult, error)
}

type Status string

const (
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

type ResourceUsage struct {
	Resource string `json:"resource"`
	Units    int64  `json:"units"`
}

type Failure struct {
	Class       string          `json:"class"`
	Scope       string          `json:"scope"`
	Retryable   bool            `json:"retryable"`
	Code        string          `json:"code"`
	Message     string          `json:"message"`
	Evidence    []artifact.Ref  `json:"evidence,omitempty"`
	Diagnostics json.RawMessage `json:"diagnostics,omitempty"`
}

type RunResult struct {
	Status     Status          `json:"status"`
	Output     *artifact.Ref   `json:"output,omitempty"`
	Usage      []ResourceUsage `json:"usage,omitempty"`
	Failure    *Failure        `json:"failure,omitempty"`
	StartedAt  time.Time       `json:"started_at"`
	FinishedAt time.Time       `json:"finished_at"`
}

type Result struct {
	RunResult
	Trajectory artifact.Ref `json:"trajectory"`
}

func (r RunResult) WithTrajectory(trajectory artifact.Ref) Result {
	return Result{RunResult: r, Trajectory: trajectory}
}
