package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/go-go-golems/optkit/record"
)

var (
	ErrLeaseNotFound    = errors.New("work lease not found")
	ErrTerminalConflict = errors.New("conflicting terminal work result")
	ErrWorkConflict     = errors.New("work semantic key conflict")
	ErrWorkNotFound     = errors.New("work item not found")
)

type Queue interface {
	Enqueue(context.Context, []WorkItem) error
	Lease(context.Context, string, LeaseRequest) ([]Lease, error)
	Heartbeat(context.Context, record.LeaseID, string, time.Duration) error
	Complete(context.Context, record.LeaseID, WorkResult) error
	Fail(context.Context, record.LeaseID, WorkFailure) error
	ReclaimExpired(context.Context, time.Time) (int64, error)
	Get(context.Context, record.WorkID) (WorkRecord, error)
}
