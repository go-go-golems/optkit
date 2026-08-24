//go:build cgo

package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/artifact/memory"
	"github.com/go-go-golems/optkit/budget"
	"github.com/go-go-golems/optkit/campaign"
	"github.com/go-go-golems/optkit/projection"
	"github.com/go-go-golems/optkit/record"
	"github.com/go-go-golems/optkit/scheduler"
)

type mutableClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *mutableClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *mutableClock) Set(value time.Time) {
	c.mu.Lock()
	c.now = value.UTC()
	c.mu.Unlock()
}

func TestJournalIdempotencyVersioningAndReopen(t *testing.T) {
	ctx := context.Background()
	clock := &mutableClock{now: time.Unix(1000, 0).UTC()}
	path := filepath.Join(t.TempDir(), "optkit.db")
	store, err := OpenWithClock(path, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	artifacts := memory.New()
	campaignID := record.CampaignID("campaign:journal-test")
	payload, err := artifact.PutCanonical(ctx, artifacts, "schema:test.payload/v1", artifact.SensitivityInternal, map[string]string{"step": "create"})
	if err != nil {
		t.Fatal(err)
	}
	commandID := record.CommandID("command:create")
	input := campaign.NewEvent{Kind: campaign.CampaignCreated, Schema: "schema:test.payload/v1", Actor: "actor:test", Command: &commandID, Payload: payload}
	first, err := store.Append(ctx, campaignID, 0, []campaign.NewEvent{input})
	if err != nil {
		t.Fatal(err)
	}
	if first.Version != 1 || first.Duplicate {
		t.Fatalf("unexpected first append: %+v", first)
	}

	duplicate, err := store.Append(ctx, campaignID, 0, []campaign.NewEvent{input})
	if err != nil {
		t.Fatal(err)
	}
	if !duplicate.Duplicate || duplicate.Version != 1 || duplicate.Events[0].ID != first.Events[0].ID {
		t.Fatalf("unexpected duplicate result: %+v", duplicate)
	}

	planPayload, err := artifact.PutCanonical(ctx, artifacts, "schema:test.plan/v1", artifact.SensitivityInternal, map[string]string{"plan": "v1"})
	if err != nil {
		t.Fatal(err)
	}
	planCommand := record.CommandID("command:plan")
	_, err = store.Append(ctx, campaignID, 0, []campaign.NewEvent{{Kind: campaign.PlanCompiled, Schema: "schema:test.plan/v1", Actor: "actor:test", Command: &planCommand, Payload: planPayload}})
	if !errors.Is(err, campaign.ErrVersionConflict) {
		t.Fatalf("stale append error = %v", err)
	}
	second, err := store.Append(ctx, campaignID, 1, []campaign.NewEvent{{Kind: campaign.PlanCompiled, Schema: "schema:test.plan/v1", Actor: "actor:test", Command: &planCommand, Payload: planPayload}})
	if err != nil {
		t.Fatal(err)
	}
	if second.Version != 2 {
		t.Fatalf("version = %d", second.Version)
	}
	if err := store.Verify(ctx, campaignID); err != nil {
		t.Fatal(err)
	}

	events, err := store.Read(ctx, campaignID, 0)
	if err != nil {
		t.Fatal(err)
	}
	overview, err := projection.RebuildOverview(campaignID, events)
	if err != nil {
		t.Fatal(err)
	}
	if overview.Status != campaign.StatusReady || overview.Version != 2 {
		t.Fatalf("unexpected overview: %+v", overview)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := OpenWithClock(path, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if err := reopened.Verify(ctx, campaignID); err != nil {
		t.Fatal(err)
	}
	head, err := reopened.Head(ctx, campaignID)
	if err != nil {
		t.Fatal(err)
	}
	if head.Version != 2 || head.LastDigest != second.Events[0].Digest {
		t.Fatalf("unexpected reopened head: %+v", head)
	}
}

func TestControllerReturnsPriorCommandAfterStateAdvanced(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "controller.db")
	clock := &mutableClock{now: time.Unix(2000, 0).UTC()}
	store, err := OpenWithClock(path, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	artifacts := memory.New()
	controller := campaign.Controller{Journal: store, Artifacts: artifacts}
	campaignID := record.CampaignID("campaign:controller-test")
	create := campaign.Command{ID: "command:create", Campaign: campaignID, ExpectedVersion: 0, Kind: campaign.CreateCampaign, Actor: "actor:test", Schema: "schema:test.command/v1", Payload: map[string]string{"name": "test"}}
	first, err := controller.Handle(ctx, create)
	if err != nil {
		t.Fatal(err)
	}
	compile := campaign.Command{ID: "command:compile", Campaign: campaignID, ExpectedVersion: 1, Kind: campaign.CompilePlan, Actor: "actor:test", Schema: "schema:test.command/v1", Payload: map[string]string{"plan": "v1"}}
	if _, err := controller.Handle(ctx, compile); err != nil {
		t.Fatal(err)
	}
	prior, err := controller.Handle(ctx, create)
	if err != nil {
		t.Fatal(err)
	}
	if !prior.Duplicate || prior.Events[0].ID != first.Events[0].ID {
		t.Fatalf("controller did not return prior command: %+v", prior)
	}
}

func TestQueueLeaseRecoveryAndTerminalIdempotency(t *testing.T) {
	ctx := context.Background()
	baseTime := time.Unix(3000, 0).UTC()
	clock := &mutableClock{now: baseTime}
	path := filepath.Join(t.TempDir(), "queue.db")
	store, err := OpenWithClock(path, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	artifacts := memory.New()
	campaignID := record.CampaignID("campaign:queue-test")
	createPayload, err := artifact.PutCanonical(ctx, artifacts, "schema:test.campaign/v1", artifact.SensitivityInternal, map[string]string{"name": "queue"})
	if err != nil {
		t.Fatal(err)
	}
	cmd := record.CommandID("command:create")
	if _, err := store.Append(ctx, campaignID, 0, []campaign.NewEvent{{Kind: campaign.CampaignCreated, Schema: "schema:test.campaign/v1", Actor: "actor:test", Command: &cmd, Payload: createPayload}}); err != nil {
		t.Fatal(err)
	}

	workPayload, err := artifact.PutCanonical(ctx, artifacts, "schema:test.work/v1", artifact.SensitivityInternal, map[string]int{"value": 7})
	if err != nil {
		t.Fatal(err)
	}
	item, err := scheduler.NewWorkItem(campaignID, "episode", record.SumBytes([]byte("episode-key")), workPayload, 10, baseTime, time.Minute, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Enqueue(ctx, []scheduler.WorkItem{item}); err != nil {
		t.Fatal(err)
	}
	leases, err := store.Lease(ctx, "worker-1", scheduler.LeaseRequest{Limit: 1, Now: baseTime})
	if err != nil {
		t.Fatal(err)
	}
	if len(leases) != 1 || leases[0].Attempt != 1 {
		t.Fatalf("unexpected leases: %+v", leases)
	}
	firstLease := leases[0]
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	clock.Set(baseTime.Add(2 * time.Minute))
	reopened, err := OpenWithClock(path, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	reclaimed, err := reopened.ReclaimExpired(ctx, clock.Now())
	if err != nil {
		t.Fatal(err)
	}
	if reclaimed != 1 {
		t.Fatalf("reclaimed = %d", reclaimed)
	}
	leases, err = reopened.Lease(ctx, "worker-2", scheduler.LeaseRequest{Limit: 1, Now: clock.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if len(leases) != 1 || leases[0].Attempt != 2 || leases[0].ID == firstLease.ID {
		t.Fatalf("unexpected recovered lease: %+v", leases)
	}
	resultRef, err := artifact.PutCanonical(ctx, artifacts, "schema:test.result/v1", artifact.SensitivityInternal, map[string]int{"output": 21})
	if err != nil {
		t.Fatal(err)
	}
	result := scheduler.WorkResult{Artifact: resultRef}
	if err := reopened.Complete(ctx, leases[0].ID, result); err != nil {
		t.Fatal(err)
	}
	if err := reopened.Complete(ctx, leases[0].ID, result); err != nil {
		t.Fatalf("idempotent completion failed: %v", err)
	}
	otherRef, err := artifact.PutCanonical(ctx, artifacts, "schema:test.result/v1", artifact.SensitivityInternal, map[string]int{"output": 22})
	if err != nil {
		t.Fatal(err)
	}
	if err := reopened.Complete(ctx, leases[0].ID, scheduler.WorkResult{Artifact: otherRef}); !errors.Is(err, scheduler.ErrTerminalConflict) {
		t.Fatalf("conflicting completion error = %v", err)
	}
	recordValue, err := reopened.Get(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if recordValue.Status != scheduler.WorkCompleted || recordValue.Attempt != 2 || recordValue.Result == nil || recordValue.Result.Artifact.Digest != resultRef.Digest {
		t.Fatalf("unexpected work record: %+v", recordValue)
	}
}

func TestQueueRetryBackoff(t *testing.T) {
	ctx := context.Background()
	baseTime := time.Unix(4000, 0).UTC()
	clock := &mutableClock{now: baseTime}
	store, err := OpenWithClock(filepath.Join(t.TempDir(), "retry.db"), clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	artifacts := memory.New()
	campaignID := record.CampaignID("campaign:retry-test")
	payload, err := artifact.PutCanonical(ctx, artifacts, "schema:test/v1", artifact.SensitivityInternal, map[string]int{"x": 1})
	if err != nil {
		t.Fatal(err)
	}
	cmd := record.CommandID("command:create")
	if _, err := store.Append(ctx, campaignID, 0, []campaign.NewEvent{{Kind: campaign.CampaignCreated, Schema: "schema:test/v1", Actor: "actor:test", Command: &cmd, Payload: payload}}); err != nil {
		t.Fatal(err)
	}
	item, err := scheduler.NewWorkItem(campaignID, "instrument", record.SumBytes([]byte("instrument-key")), payload, 0, baseTime, time.Minute, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Enqueue(ctx, []scheduler.WorkItem{item}); err != nil {
		t.Fatal(err)
	}
	leases, err := store.Lease(ctx, "worker", scheduler.LeaseRequest{Limit: 1, Now: baseTime})
	if err != nil || len(leases) != 1 {
		t.Fatalf("lease = %+v, err=%v", leases, err)
	}
	if err := store.Fail(ctx, leases[0].ID, scheduler.WorkFailure{Code: "rate_limited", Message: "retry", Retryable: true, Backoff: 30 * time.Second}); err != nil {
		t.Fatal(err)
	}
	leases, err = store.Lease(ctx, "worker", scheduler.LeaseRequest{Limit: 1, Now: baseTime.Add(10 * time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	if len(leases) != 0 {
		t.Fatalf("work ignored backoff: %+v", leases)
	}
	clock.Set(baseTime.Add(31 * time.Second))
	leases, err = store.Lease(ctx, "worker", scheduler.LeaseRequest{Limit: 1, Now: clock.Now()})
	if err != nil || len(leases) != 1 || leases[0].Attempt != 2 {
		t.Fatalf("retry lease = %+v, err=%v", leases, err)
	}
	if err := store.Fail(ctx, leases[0].ID, scheduler.WorkFailure{Code: "invalid", Message: "terminal"}); err != nil {
		t.Fatal(err)
	}
	recordValue, err := store.Get(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if recordValue.Status != scheduler.WorkFailed || recordValue.Failure == nil || recordValue.Failure.Code != "invalid" {
		t.Fatalf("unexpected failed work: %+v", recordValue)
	}
}

func TestJournalTamperingIsDetected(t *testing.T) {
	ctx := context.Background()
	store, err := OpenWithClock(filepath.Join(t.TempDir(), "tamper.db"), func() time.Time { return time.Unix(5000, 0).UTC() })
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	artifacts := memory.New()
	payload, err := artifact.PutCanonical(ctx, artifacts, "schema:test/v1", artifact.SensitivityInternal, map[string]int{"x": 1})
	if err != nil {
		t.Fatal(err)
	}
	commandID := record.CommandID("command:create")
	campaignID := record.CampaignID("campaign:tamper-test")
	if _, err := store.Append(ctx, campaignID, 0, []campaign.NewEvent{{Kind: campaign.CampaignCreated, Schema: "schema:test/v1", Actor: "actor:test", Command: &commandID, Payload: payload, Tags: map[string]string{"original": "yes"}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(ctx, `UPDATE campaign_events SET tags_json = ? WHERE campaign_id = ? AND seq = 1`, []byte(`{"tampered":"yes"}`), string(campaignID)); err != nil {
		t.Fatal(err)
	}
	if err := store.Verify(ctx, campaignID); !errors.Is(err, campaign.ErrJournalCorrupt) {
		t.Fatalf("tamper verification error = %v", err)
	}
}

func TestQueueRejectsSemanticKeyWithDifferentPayload(t *testing.T) {
	ctx := context.Background()
	baseTime := time.Unix(6000, 0).UTC()
	store, err := OpenWithClock(filepath.Join(t.TempDir(), "conflict.db"), func() time.Time { return baseTime })
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	artifacts := memory.New()
	campaignID := record.CampaignID("campaign:work-conflict")
	campaignPayload, err := artifact.PutCanonical(ctx, artifacts, "schema:test/v1", artifact.SensitivityInternal, map[string]int{"x": 0})
	if err != nil {
		t.Fatal(err)
	}
	commandID := record.CommandID("command:create")
	if _, err := store.Append(ctx, campaignID, 0, []campaign.NewEvent{{Kind: campaign.CampaignCreated, Schema: "schema:test/v1", Actor: "actor:test", Command: &commandID, Payload: campaignPayload}}); err != nil {
		t.Fatal(err)
	}
	firstPayload, err := artifact.PutCanonical(ctx, artifacts, "schema:test.work/v1", artifact.SensitivityInternal, map[string]int{"x": 1})
	if err != nil {
		t.Fatal(err)
	}
	secondPayload, err := artifact.PutCanonical(ctx, artifacts, "schema:test.work/v1", artifact.SensitivityInternal, map[string]int{"x": 2})
	if err != nil {
		t.Fatal(err)
	}
	key := record.SumBytes([]byte("same-semantic-key"))
	first, err := scheduler.NewWorkItem(campaignID, "episode", key, firstPayload, 0, baseTime, time.Minute, nil)
	if err != nil {
		t.Fatal(err)
	}
	second, err := scheduler.NewWorkItem(campaignID, "episode", key, secondPayload, 0, baseTime, time.Minute, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Enqueue(ctx, []scheduler.WorkItem{first}); err != nil {
		t.Fatal(err)
	}
	if err := store.Enqueue(ctx, []scheduler.WorkItem{second}); !errors.Is(err, scheduler.ErrWorkConflict) {
		t.Fatalf("semantic conflict error = %v", err)
	}
}

func TestBudgetReservationCommitReleaseAndReopen(t *testing.T) {
	ctx := context.Background()
	clock := &mutableClock{now: time.Unix(7000, 0).UTC()}
	path := filepath.Join(t.TempDir(), "budget.db")
	store, err := OpenWithClock(path, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	campaignID := createBudgetCampaign(t, ctx, store, "campaign:budget-test")
	limits := []budget.Limit{{Resource: "operations", Units: 10}, {Resource: "usd_micros", Units: 1000}}
	if err := store.Define(ctx, campaignID, limits); err != nil {
		t.Fatal(err)
	}
	if err := store.Define(ctx, campaignID, limits); err != nil {
		t.Fatalf("idempotent define failed: %v", err)
	}
	if err := store.Define(ctx, campaignID, []budget.Limit{{Resource: "operations", Units: 11}}); !errors.Is(err, budget.ErrLimitConflict) {
		t.Fatalf("limit conflict error = %v", err)
	}

	first, err := store.Reserve(ctx, budget.ReserveRequest{Campaign: campaignID, Work: "work:budget-one", Claims: []budget.Quantity{{Resource: "operations", Units: 6}}})
	if err != nil {
		t.Fatal(err)
	}
	duplicate, err := store.Reserve(ctx, budget.ReserveRequest{Campaign: campaignID, Work: "work:budget-one", Claims: []budget.Quantity{{Resource: "operations", Units: 6}}})
	if err != nil || duplicate.ID != first.ID {
		t.Fatalf("idempotent reserve = %+v, err=%v", duplicate, err)
	}
	if _, err := store.Reserve(ctx, budget.ReserveRequest{Campaign: campaignID, Work: "work:budget-one", Claims: []budget.Quantity{{Resource: "operations", Units: 5}}}); !errors.Is(err, budget.ErrReservationConflict) {
		t.Fatalf("reservation conflict error = %v", err)
	}
	second, err := store.Reserve(ctx, budget.ReserveRequest{Campaign: campaignID, Work: "work:budget-two", Claims: []budget.Quantity{{Resource: "operations", Units: 4}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Reserve(ctx, budget.ReserveRequest{Campaign: campaignID, Work: "work:budget-three", Claims: []budget.Quantity{{Resource: "operations", Units: 1}}}); !errors.Is(err, budget.ErrInsufficient) {
		t.Fatalf("insufficient budget error = %v", err)
	}

	committed, err := store.Commit(ctx, first.ID, []budget.Quantity{{Resource: "operations", Units: 5}})
	if err != nil {
		t.Fatal(err)
	}
	if committed.Status != budget.StatusCommitted || committed.Overage {
		t.Fatalf("unexpected committed reservation: %+v", committed)
	}
	if _, err := store.Commit(ctx, first.ID, []budget.Quantity{{Resource: "operations", Units: 5}}); err != nil {
		t.Fatalf("idempotent commit failed: %v", err)
	}
	if _, err := store.Commit(ctx, first.ID, []budget.Quantity{{Resource: "operations", Units: 4}}); !errors.Is(err, budget.ErrReservationConflict) {
		t.Fatalf("commit conflict error = %v", err)
	}
	if _, err := store.Release(ctx, first.ID); !errors.Is(err, budget.ErrReservationTerminal) {
		t.Fatalf("release committed error = %v", err)
	}
	if _, err := store.Release(ctx, second.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Release(ctx, second.ID); err != nil {
		t.Fatalf("idempotent release failed: %v", err)
	}

	snapshot, err := store.Snapshot(ctx, campaignID)
	if err != nil {
		t.Fatal(err)
	}
	operations := findBudgetResource(t, snapshot, "operations")
	if operations.Limit != 10 || operations.Reserved != 0 || operations.Committed != 5 || operations.Available != 5 || snapshot.Violated {
		t.Fatalf("unexpected budget snapshot: %+v", snapshot)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := OpenWithClock(path, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	replayed, err := reopened.Snapshot(ctx, campaignID)
	if err != nil {
		t.Fatal(err)
	}
	if got := findBudgetResource(t, replayed, "operations"); got != operations {
		t.Fatalf("reopened snapshot = %+v, want %+v", got, operations)
	}
}

func TestBudgetConcurrentReservationConservesCap(t *testing.T) {
	ctx := context.Background()
	store, err := OpenWithClock(filepath.Join(t.TempDir(), "budget-race.db"), func() time.Time { return time.Unix(8000, 0).UTC() })
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	campaignID := createBudgetCampaign(t, ctx, store, "campaign:budget-race")
	if err := store.Define(ctx, campaignID, []budget.Limit{{Resource: "operations", Units: 10}}); err != nil {
		t.Fatal(err)
	}

	start := make(chan struct{})
	results := make(chan error, 2)
	for _, workID := range []record.WorkID{"work:race-one", "work:race-two"} {
		workID := workID
		go func() {
			<-start
			_, err := store.Reserve(ctx, budget.ReserveRequest{Campaign: campaignID, Work: workID, Claims: []budget.Quantity{{Resource: "operations", Units: 6}}})
			results <- err
		}()
	}
	close(start)
	successes := 0
	insufficient := 0
	for range 2 {
		err := <-results
		switch {
		case err == nil:
			successes++
		case errors.Is(err, budget.ErrInsufficient):
			insufficient++
		default:
			t.Fatalf("unexpected reserve error: %v", err)
		}
	}
	if successes != 1 || insufficient != 1 {
		t.Fatalf("successes=%d insufficient=%d", successes, insufficient)
	}
	snapshot, err := store.Snapshot(ctx, campaignID)
	if err != nil {
		t.Fatal(err)
	}
	operations := findBudgetResource(t, snapshot, "operations")
	if operations.Reserved != 6 || operations.Available != 4 || snapshot.Violated {
		t.Fatalf("budget conservation failed: %+v", snapshot)
	}
}

func TestBudgetCommitsUnavoidableOverageAsVisibleCustody(t *testing.T) {
	ctx := context.Background()
	store, err := OpenWithClock(filepath.Join(t.TempDir(), "budget-overage.db"), func() time.Time { return time.Unix(9000, 0).UTC() })
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	campaignID := createBudgetCampaign(t, ctx, store, "campaign:budget-overage")
	if err := store.Define(ctx, campaignID, []budget.Limit{{Resource: "operations", Units: 10}}); err != nil {
		t.Fatal(err)
	}
	reservation, err := store.Reserve(ctx, budget.ReserveRequest{Campaign: campaignID, Work: "work:overage", Claims: []budget.Quantity{{Resource: "operations", Units: 5}}})
	if err != nil {
		t.Fatal(err)
	}
	committed, err := store.Commit(ctx, reservation.ID, []budget.Quantity{{Resource: "operations", Units: 11}})
	if err != nil {
		t.Fatal(err)
	}
	if !committed.Overage || committed.Status != budget.StatusCommitted {
		t.Fatalf("overage was not retained: %+v", committed)
	}
	snapshot, err := store.Snapshot(ctx, campaignID)
	if err != nil {
		t.Fatal(err)
	}
	operations := findBudgetResource(t, snapshot, "operations")
	if !snapshot.Violated || operations.Committed != 11 || operations.Available != -1 {
		t.Fatalf("unexpected overage snapshot: %+v", snapshot)
	}
}

func createBudgetCampaign(t *testing.T, ctx context.Context, store *Store, raw string) record.CampaignID {
	t.Helper()
	artifacts := memory.New()
	payload, err := artifact.PutCanonical(ctx, artifacts, "schema:test.budget-campaign/v1", artifact.SensitivityInternal, map[string]string{"name": raw})
	if err != nil {
		t.Fatal(err)
	}
	campaignID := record.CampaignID(raw)
	commandID := record.CommandID("command:" + raw[len("campaign:"):])
	if _, err := store.Append(ctx, campaignID, 0, []campaign.NewEvent{{Kind: campaign.CampaignCreated, Schema: "schema:test.budget-campaign/v1", Actor: "actor:test", Command: &commandID, Payload: payload}}); err != nil {
		t.Fatal(err)
	}
	return campaignID
}

func findBudgetResource(t *testing.T, snapshot budget.Snapshot, resource budget.Resource) budget.ResourceSnapshot {
	t.Helper()
	for _, value := range snapshot.Resources {
		if value.Resource == resource {
			return value
		}
	}
	t.Fatalf("resource %s not found in %+v", resource, snapshot)
	return budget.ResourceSnapshot{}
}
