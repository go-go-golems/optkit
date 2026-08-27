---
Title: Scientific Campaign Explorer Architecture and Implementation Guide
Ticket: OPTKIT-003
Status: active
Topics:
    - optkit
    - ui
    - query-plane
    - visualization
    - scientific-workflow
    - architecture
    - local-development
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://campaign/event.go
      Note: Control event identity, sequence, provenance, and payload references
    - Path: repo://campaign/journal.go
      Note: Authoritative campaign read/replay contract
    - Path: repo://episode/writer.go
      Note: Verified episode trajectory evidence model
    - Path: repo://examples/numbergame/demo.go
      Note: Complete scientist workflow and UI fixture source
    - Path: repo://internal/web/server.go
      Note: Implemented GET API, SSE, security, and static embed
    - Path: repo://internal/web/static/app.js
      Note: Implemented dependency-free campaign navigation and visualization
    - Path: repo://query/service.go
      Note: Implemented read-only projections, limits, and payload policy
ExternalSources: []
Summary: Intern-oriented architecture and implementation guide for a read-only Optkit campaign explorer driven by CLI-authored journal facts, rebuildable projections, search, deep links, SSE replay, and scientific visualizations.
LastUpdated: 2026-08-25T00:30:00-04:00
WhatFor: Explain the existing Optkit evidence model and specify the backend, API, frontend, visualization, safety, testing, and delivery work needed to build the read-only scientific campaign explorer.
WhenToUse: Before implementing or reviewing the Optkit query server, embedded vanilla HTML/CSS/JavaScript application, search index, SSE stream, campaign projections, or scientific visualizations.
---


# Scientific Campaign Explorer Architecture and Implementation Guide

## 1. Executive summary

Optkit is a local-first framework for running attributable optimization experiments. Its authoritative record is not a mutable “campaign row.” It is an append-only campaign journal whose events point to immutable content-addressed artifacts. Episodes carry their own sealed trajectories. Measurements are immutable observations under explicit measurement epochs. Estimates and decisions are also artifacts recorded by journal events. A projection rebuilds useful views from those facts.

The requested UI should **not** create campaigns, edit parameters, start workers, pause campaigns, or record decisions. Those writes belong to CLI and application-command workflows, commonly operated by LLM agents. The browser is a read-only scientific query plane. Its job is to help a scientist answer:

- What question was tested?
- What changed between baseline and candidate?
- Which cases and arms ran?
- Did the intended intervention actually execute?
- What happened inside each episode?
- Which evidence supports each observation, estimate, and decision?
- Were pairs missing, budgets exceeded, or measurement epochs mixed?
- How does this campaign compare with another?
- Can I give an agent a stable deep link to the evidence I reviewed?

The first implementation should use the Numbergame campaign already in the repository. Numbergame is small enough to understand completely but exercises the real substrate: typed snapshots and patches, candidate lineage, a complete-block trial, eight durable episodes, two finite budgets, artifact-backed trajectories, 24 observations, a paired estimate, a decision, 81 control events, process restarts, and journal verification.

The recommended architecture is:

```text
LLM agent / scientist CLI
          |
          | commands and workers (writes)
          v
 campaign journal + queue + budgets    filesystem CAS
          |                                  |
          +-------------- facts ------------+
                             |
                      read-only query service
                    projections + search index
                             |
                  GET JSON + replayable GET SSE
                             |
              embedded HTML/CSS/JavaScript explorer
       navigation + search + graphs + tables + charts
```

The browser never reads SQLite or the artifact filesystem directly. It never sends a mutation command. The backend uses standard-library `net/http` with Go 1.22+ `http.ServeMux`, serves `/api/v1/...`, and embeds plain HTML, CSS, and JavaScript directly into the Optkit binary. There is no Node, pnpm, Vite, React, or frontend build step.

## 2. Non-negotiable product boundary

### 2.1 Writes remain outside the UI

Campaign creation and mutation remain in CLI/application code:

```text
create campaign
materialize baseline
propose candidate
compile trial
start/pause/resume/stop
run workers
record decision
```

The UI does not reproduce those validations in forms. It observes committed facts after they appear in the journal or artifact store.

### 2.2 The UI may still be interactive

Read-only does not mean passive. The user can:

- search, filter, sort, group, and compare;
- pan and zoom lineage and trajectory graphs;
- replay the journal to any sequence;
- follow a running campaign live;
- open evidence drawers;
- switch chart/table representations;
- copy IDs, digests, deep links, and CLI inspection commands;
- download permitted exports;
- save view state in the URL.

These interactions change browser presentation, not campaign facts.

### 2.3 The handoff loop

```text
agent authors through CLI
        |
        v
Optkit records immutable facts
        |
        v
scientist explores deep-linked evidence
        |
        v
scientist gives evidence-linked direction
        |
        v
agent creates a follow-up campaign through CLI
```

A later CLI improvement should print the canonical explorer URL after every write-oriented workflow.

## 3. Existing system: the mental model

An intern should learn four layers before writing UI code.

### 3.1 Semantic records

A semantic record answers “what is this thing?” Examples include:

- baseline and challenger snapshots;
- patches and candidates;
- dataset and trial manifests;
- episode work and results;
- observations and epochs;
- estimates and decisions.

Many IDs are derived from canonical content. Equal semantics therefore produce equal content identities.

### 3.2 Campaign control journal

The campaign journal answers “what facts were accepted, and in what order?” `campaign.Journal` at `campaign/journal.go:31` exposes:

```go
type Journal interface {
    Append(ctx, campaignID, expectedVersion, events) (AppendResult, error)
    Read(ctx, campaignID, afterSequence) ([]ControlEvent, error)
    Head(ctx, campaignID) (Head, error)
    Verify(ctx, campaignID) error
}
```

A `campaign.ControlEvent` at `campaign/event.go:55` contains:

```text
campaign + sequence + kind + schema + subject
occurred/recorded time + actor
command/causation/correlation
previous digest + payload artifact ref + tags + event digest
```

The sequence is the authoritative replay cursor. The event payload is not inlined; it is an immutable artifact reference.

### 3.3 Episode trajectory

A campaign event says that an episode completed. The episode trajectory explains what happened inside the execution. `episode.Event` at `episode/types.go:14` contains an episode-local sequence, event kind, schema, span, optional parent span, payload artifact, tags, and sensitivity marker.

`episode.Writer.Seal` stores one `schema:optkit.trajectory/v1` artifact. `episode.LoadTrajectory` at `episode/writer.go:155` verifies episode ownership, sequence order, and every event payload before returning it.

### 3.4 Rebuildable projections

A projection is a disposable read model. `projection.RebuildOverview` at `projection/overview.go:23` folds all events into:

```go
type Overview struct {
    Campaign       record.CampaignID
    Version        uint64
    Status         campaign.Status
    Events         uint64
    Episodes       int
    Queued         int
    Active         int
    Completed      int
    FailedTerminal int
    LastDigest     record.Digest
}
```

The UI needs many more projections, but they should follow this rule:

> If a projection is deleted, it can be rebuilt from journal facts and verified artifacts.

## 4. What exists today

### 4.1 Local composition

`local.Profile` at `local/profile.go:13` composes:

```text
<root>/optkit.db     SQLite control metadata, queue, and budgets
<root>/artifacts/    filesystem content-addressed store
```

SQLite uses WAL mode. The artifact store verifies digest and size and stores immutable files.

### 4.2 Current CLI query surface

`cmd/optkit/main.go` currently provides:

```text
optkit demo
optkit campaign inspect
optkit campaign verify
optkit artifact verify
```

`campaign inspect` at `cmd/optkit/main.go:100` verifies the journal, rebuilds the overview, reads the budget snapshot, counts event kinds, and returns a tail. This is a useful first projection but not enough for navigation.

### 4.3 Current Numbergame evidence

The reproducible capture in `sources/01-numbergame-evidence.md` records:

```text
campaign events:       81
episodes:               8
observations:          24
budget reservations:   8
usage commits:          8
terminal failures:      0
paired delta:           +0.71041675
decision:               eligible
```

Event counts are:

```text
CampaignCreated          1
PlanCompiled             1
CampaignStarted          1
CandidateProposed        1
SnapshotMaterialized     1
TrialPlanned             1
BudgetReserved           8
EpisodeScheduled         8
EpisodeLeaseGranted      8
EpisodeAttemptStarted    8
UsageCommitted           8
EpisodeCompleted         8
ObservationRecorded     24
EstimateRecorded         1
DecisionRecorded         1
CampaignCompleted        1
```

### 4.4 Implemented v0 query capabilities

This ticket now implements:

- stable SQLite campaign listing through `store/sqlite/query.go`;
- generic campaign overview, budget, integrity, and event-kind projections;
- event response limits and sequence cursors (the v0 adapter still reads the remaining journal before slicing);
- bounded public/internal JSON payload previews on visible control events;
- GET-only JSON endpoints for campaign list, campaign detail, and events;
- replayable SSE using sequence cursors and cross-process journal polling;
- directly embedded white monochrome HTML/CSS/JavaScript;
- campaign search/filter, deep-link navigation, stage rail, metrics, lineage, trial matrix, paired deltas, budget custody, fact inventory, event timeline, and JSON evidence drawer;
- `optkit serve --store ... --listen ...`;
- API, read-only method, SSE, static asset, race, lint, and Chromium smoke validation.

### 4.5 Remaining query capabilities

The next slices still need:

- entity lookup by subject ID;
- nested result/trajectory/artifact traversal;
- server-side schema projector registry;
- measurement QA grouped by episode/case/arm;
- server-side structured search or derived FTS index;
- replay slider at an exact historical sequence;
- cross-campaign comparison;
- bounded database-side event pagination rather than read-then-slice;
- artifact reachability catalog beyond event payloads;
- large-campaign virtualization and accessibility audit.

Do not expose private SQLite tables directly to fill these gaps.

## 5. Numbergame scientific narrative

### 5.1 Research question

```text
Does multiplier 3 improve target accuracy over multiplier 2
on four development cases without violating the campaign budget?
```

### 5.2 Configuration lineage

```text
baseline snapshot                 challenger snapshot
multiplier: 2                     multiplier: 3
noise: none                       noise: none
       |                                 ^
       +--- patch multiplier 2 -> 3 -----+
                         |
                  candidate hypothesis
```

The variable domain is `1..10`; noise is one of `none`, `small`, or `large`. See `examples/numbergame/model.go`.

### 5.3 Trial

```text
                    baseline       challenger
case-01: 1 -> 3       episode         episode
case-02: 2 -> 6       episode         episode
case-03: 4 -> 12      episode         episode
case-04: 7 -> 21      episode         episode
```

`experiment.NewCompleteBlockTrial` stores two arms, one immutable dataset, one repeat, and `numbergame.execution/v1`. `experiment.Expand` at `experiment/trial.go:115` derives deterministic episode identities and seeds.

### 5.4 Execution

For every episode, `examples/numbergame/system.go` emits:

```text
input.received
  numbergame.configured
  numbergame.noise.sampled
output.produced
episode.completed
```

The output records input, multiplier, noise, result, target, and absolute error.

### 5.5 Measurement

`examples/numbergame/measure.go` emits three observations:

1. intervention integrity: did the observed config equal the assigned config?
2. absolute error: `abs(output - target)`;
3. fake judge: `1 / (1 + absolute_error)`.

Every observation names construct, instrument, protocol, epoch, status, value, evidence, diagnostics, and repeat. The general type is at `measure/observation.go:62`; epoch identity is defined in `measure/epoch.go`.

### 5.6 Analysis

`experiment.PairedMean` at `experiment/estimate.go:26` pairs baseline and challenger by case/repeat. Missing or invalid pairs make the estimate invalid rather than silently shrinking the denominator.

```text
case-01: 0.500 -> 1.000  delta +0.500
case-02: 0.333 -> 1.000  delta +0.667
case-03: 0.200 -> 1.000  delta +0.800
case-04: 0.125 -> 1.000  delta +0.875
mean delta:                     +0.71041675
```

### 5.7 Decision

The deterministic policy requires:

- no budget violation;
- every intervention exercised;
- positive paired accuracy delta.

The resulting status is `eligible`. `RunDemo` then restarts the process and rebuilds the terminal overview from journal/artifacts, proving that the result is not retained only in memory.

## 6. Scientist navigation workflow

### 6.1 Campaign discovery

The scientist begins at `/campaigns`, not a creation form.

They can search/filter by:

```text
system:numbergame
status:completed
actor:demo
kind:EpisodeFailed
construct:target.accuracy
schema:optkit.observation/v1
subject:case-03
id:campaign:...
digest:sha256:...
```

Results are typed and deep-linked: campaign, event, episode, observation, estimate, decision, or artifact.

### 6.2 Campaign overview

`/campaigns/{campaignID}` answers the high-level questions:

- campaign status/version/digest;
- system and research question;
- baseline and candidate;
- trial/case/arm counts;
- episode state counts;
- budget utilization;
- estimate and decision;
- journal verification state;
- recent events.

### 6.3 Provenance navigation

`/campaigns/{id}/lineage` visualizes:

```text
baseline -> patch -> challenger -> candidate -> trial
                                          |
                             case x arm episodes
                                          |
                               observations/epochs
                                          |
                                  estimate -> decision
```

Nodes open a side drawer and have stable URLs.

### 6.4 Trial matrix

`/campaigns/{id}/trial` renders case × arm × repeat cells. Cell state distinguishes planned, queued, leased, running, completed, failed, missing observation, and invalid intervention.

### 6.5 Timeline and replay

`/campaigns/{id}/timeline?at=57` reconstructs the view at exact journal sequence 57. A slider changes only the projection cutoff.

### 6.6 Episode inspection

`/campaigns/{id}/episodes/{episodeID}` shows assignment, lifecycle, attempts, result/failure, usage, observations, and trajectory. The trajectory is a span tree with ordered events and payload drawers.

### 6.7 Measurement QA

The scientist checks intervention observations before outcome charts. A violated intervention is not a low score; it is an invalid treatment cell.

### 6.8 Analysis and decision

Paired charts show each case, missingness, compatible epoch, and aggregate result. The decision page links backward through estimate → observations → episode → trajectory → payload.

### 6.9 Cross-campaign comparison

`/compare?campaign=A&campaign=B` warns on incompatible system, dataset, protocol, or epochs before displaying comparative charts.

## 7. Proposed backend package architecture

```text
query/
  service.go          implemented read-only interfaces, models, paging,
                      journal projections, safe payload previews
  registry.go         future schema payload decoders/projectors
  search.go           future structured search AST and result types
  artifact.go         future nested artifact reachability/preview logic
  numbergame.go       future server-side Numbergame projection plugin

projection/
  overview.go         existing generic projection
  lineage.go          future generic relationship graph
  episodes.go         future lifecycle/measurement indexes

store/sqlite/
  query.go            implemented stable campaign listing
  search_index.go     optional derived FTS tables

internal/web/
  server.go           implemented net/http GET API, SSE, security, embed
  server_test.go      API/static/read-only/SSE tests
  static/index.html   semantic page shell
  static/styles.css   white retro-monochrome visual system
  static/app.js       hash navigation, fetch, search, projections, drawers

cmd/optkit/
  main.go             implemented `optkit serve` composition root
```

The existing architecture rule remains: projections and query semantics do not depend on concrete SQLite. Composition roots may wire `store/sqlite` implementations.

## 8. Read-only repository APIs

Do not widen `campaign.Journal` with UI-specific methods. Keep command-side interfaces small and add query-side capabilities.

```go
package query

type CampaignFilter struct {
    Status   []campaign.Status
    System   []record.SystemID
    Actor    []record.ActorRef
    UpdatedAfter *time.Time
    Text     string
}

type PageRequest struct {
    Cursor string
    Limit  int
}

type CampaignCatalog interface {
    ListCampaigns(context.Context, CampaignFilter, PageRequest) (CampaignPage, error)
}

type EventReader interface {
    ReadEventPage(
        context.Context,
        record.CampaignID,
        uint64, // after sequence
        int,    // limit
    ) (EventPage, error)
    Head(context.Context, record.CampaignID) (campaign.Head, error)
}

type BudgetReader interface {
    Snapshot(context.Context, record.CampaignID) (budget.Snapshot, error)
}
```

The SQLite adapter can query efficient indexes. An in-memory adapter can support tests. Neither interface permits append, reserve, lease, complete, or fail.

### 8.1 Why pagination is required

Current `Journal.Read(ctx, id, from)` returns every event after `from`. That is acceptable for 81 events but unsafe for long RAG campaigns. `ReadEventPage` must enforce a server ceiling, for example 500 events.

### 8.2 Campaign listing

`campaign_heads` currently contains IDs, version, and digest but not status/system/title. A v1 catalog may:

1. list heads;
2. fold each small journal;
3. lazily decode the `CampaignCreated` payload.

A scalable version should maintain a rebuildable `query_campaigns` projection table. That table is disposable and must never become command authority.

## 9. Projection registry and schema decoding

Generic control events do not know the structure of Numbergame payloads. The query plane needs explicit schema-owned decoders.

```go
type PayloadDecoder interface {
    Schema() record.SchemaID
    Decode(context.Context, artifact.Store, artifact.Ref) (any, error)
}

type EventProjector interface {
    Schemas() []record.SchemaID
    Apply(context.Context, campaign.ControlEvent, DecodedPayload, *CampaignView) error
}

type Registry struct {
    decoders   map[record.SchemaID]PayloadDecoder
    projectors []EventProjector
}
```

Rules:

- unknown schemas remain visible as generic event/artifact metadata;
- decoder failure is shown as a projection diagnostic, not hidden;
- domain plugins live outside generic journal code;
- schema strings are versioned;
- no reflection-based “guess the payload” decoder.

The Numbergame plugin understands:

```text
schema:numbergame.campaign-spec/v1
schema:numbergame.candidate-proposal/v1
schema:numbergame.trial-plan/v1
schema:numbergame.episode-completion/v1
schema:numbergame.decision/v1
schema:optkit.snapshot/v1
schema:optkit.observation/v1
schema:optkit.estimate/v1
schema:optkit.trajectory/v1
```

## 10. Read models

All API models should carry `api_version`, campaign version, and generation time. They should be display-oriented while retaining IDs needed for navigation.

### 10.1 Campaign summary

```go
type CampaignSummary struct {
    APIVersion string              `json:"api_version"`
    ID         record.CampaignID   `json:"id"`
    Version    uint64              `json:"version"`
    Status     campaign.Status     `json:"status"`
    System     record.SystemID     `json:"system,omitempty"`
    Title      string              `json:"title,omitempty"`
    Question   string              `json:"question,omitempty"`
    Actor      record.ActorRef     `json:"actor,omitempty"`
    UpdatedAt  time.Time           `json:"updated_at"`
    Counts     CampaignCounts      `json:"counts"`
    Decision   *DecisionSummary    `json:"decision,omitempty"`
    Integrity  IntegritySummary    `json:"integrity"`
}
```

### 10.2 Lineage graph

```go
type Graph struct {
    Nodes []GraphNode `json:"nodes"`
    Edges []GraphEdge `json:"edges"`
}

type GraphNode struct {
    ID       string         `json:"id"`
    Kind     string         `json:"kind"`
    Label    string         `json:"label"`
    Status   string         `json:"status,omitempty"`
    Href     string         `json:"href"`
    Metadata map[string]any `json:"metadata,omitempty"`
}

type GraphEdge struct {
    Source string `json:"source"`
    Target string `json:"target"`
    Kind   string `json:"kind"`
}
```

### 10.3 Episode matrix

```go
type EpisodeCell struct {
    Episode      record.EpisodeID `json:"episode"`
    CaseID       string           `json:"case_id"`
    ArmID        string           `json:"arm_id"`
    Repeat       int              `json:"repeat"`
    Status       string           `json:"status"`
    Attempts     int              `json:"attempts"`
    Intervention string           `json:"intervention_status,omitempty"`
    Measurement  string           `json:"measurement_status,omitempty"`
    Href         string           `json:"href"`
}
```

### 10.4 Timeline item

Return control-event metadata separately from decoded payload preview:

```go
type TimelineItem struct {
    Seq         uint64             `json:"seq"`
    ID          record.EventID     `json:"id"`
    Kind        campaign.EventKind `json:"kind"`
    Subject     string             `json:"subject,omitempty"`
    Schema      record.SchemaID    `json:"schema"`
    Actor       record.ActorRef    `json:"actor"`
    OccurredAt  time.Time          `json:"occurred_at"`
    RecordedAt  time.Time          `json:"recorded_at"`
    Payload     ArtifactLink       `json:"payload"`
    Causation   string             `json:"causation,omitempty"`
    Correlation string             `json:"correlation,omitempty"`
    Tags        map[string]string  `json:"tags,omitempty"`
}
```

### 10.5 Scientific analysis

The analysis model must not force every domain into a Numbergame-shaped score. Use typed series:

```go
type EstimateView struct {
    EstimateID  string       `json:"estimate_id"`
    Construct   string       `json:"construct"`
    Baseline    string       `json:"baseline"`
    Treatment   string       `json:"treatment"`
    Value       Decimal      `json:"value"`
    SampleSize  int          `json:"sample_size"`
    Missing     int          `json:"missing"`
    Epochs      []EpochLink  `json:"epochs"`
    Pairs       []PairedPoint `json:"pairs"`
    Validity    ValidityView `json:"validity"`
}
```

Keep decimal values as strings in contracts when exact semantics matter; convert to numbers only in chart adapters after validation.

## 11. HTTP API

Use `http.NewServeMux()` and Go 1.22+ method/path patterns. Never use chi, gin, echo, or another router.

### 11.1 Route table

```text
GET /api/v1/health
GET /api/v1/capabilities

GET /api/v1/campaigns
GET /api/v1/campaigns/{campaignID}
GET /api/v1/campaigns/{campaignID}/lineage
GET /api/v1/campaigns/{campaignID}/events
GET /api/v1/campaigns/{campaignID}/stream
GET /api/v1/campaigns/{campaignID}/trial
GET /api/v1/campaigns/{campaignID}/episodes
GET /api/v1/campaigns/{campaignID}/measurements
GET /api/v1/campaigns/{campaignID}/analysis
GET /api/v1/campaigns/{campaignID}/decision
GET /api/v1/campaigns/{campaignID}/budgets

GET /api/v1/episodes/{episodeID}
GET /api/v1/episodes/{episodeID}/trajectory
GET /api/v1/observations/{observationID}
GET /api/v1/artifacts/{digest}
GET /api/v1/artifacts/{digest}/preview
GET /api/v1/search
```

There are deliberately no POST, PUT, PATCH, or DELETE routes.

### 11.2 Handler registration

```go
func NewServer(service *query.Service, public fs.FS) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /api/v1/campaigns", handleListCampaigns(service))
    mux.HandleFunc("GET /api/v1/campaigns/{campaignID}", handleCampaign(service))
    mux.HandleFunc("GET /api/v1/campaigns/{campaignID}/events", handleEvents(service))
    mux.HandleFunc("GET /api/v1/campaigns/{campaignID}/stream", handleStream(service))
    mux.HandleFunc("GET /api/v1/episodes/{episodeID}", handleEpisode(service))
    mux.HandleFunc("GET /api/v1/artifacts/{digest}/preview", handlePreview(service))
    mux.HandleFunc("GET /api/v1/search", handleSearch(service))

    web.RegisterSPA(mux, public, web.SPAOptions{APIPrefix: "/api"})
    return securityHeaders(mux)
}
```

The SPA fallback is registered after API routes.

### 11.3 Envelope and errors

```json
{
  "api_version": "optkit.query/v1",
  "campaign_version": 81,
  "generated_at": "2026-08-25T00:00:00Z",
  "data": {}
}
```

Errors use stable codes:

```json
{
  "error": {
    "code": "campaign_not_found",
    "message": "campaign ... was not found",
    "request_id": "..."
  }
}
```

Do not return Go type names, filesystem paths, SQL, or stack traces.

### 11.4 Caching

Use the campaign version and last digest for ETags:

```text
ETag: "campaign:...:81:sha256:..."
```

Immutable artifact previews can use the content digest as a permanent ETag.

## 12. SSE replay and live follow

### 12.1 Why polling the journal is correct initially

CLI agents and workers may run in different processes from the query server. An in-memory publish/subscribe bus would miss external writes. SQLite journal head polling works across processes and keeps the journal authoritative.

### 12.2 Protocol

The client connects:

```http
GET /api/v1/campaigns/{id}/stream?after=37
Accept: text/event-stream
Last-Event-ID: 37
```

The server emits:

```text
id: 38
event: campaign-event
data: {"seq":38,"kind":"EpisodeCompleted",...}

```

Algorithm:

```text
cursor = max(query.after, Last-Event-ID)
loop:
    page = ReadEventPage(campaign, cursor, serverPageLimit)
    for event in page:
        emit event with id=event.seq
        flush
        cursor = event.seq

    head = Head(campaign)
    if campaign terminal and cursor == head.version:
        emit campaign-caught-up
        optionally close

    wait poll interval or request cancellation
    emit comment heartbeat periodically
```

### 12.3 Correctness requirements

- events are emitted strictly in sequence order;
- reconnection may repeat an event but cannot skip committed events;
- the frontend reducer is idempotent by `(campaign, seq)`;
- slow clients have bounded buffers and are disconnected rather than consuming unbounded memory;
- request cancellation stops polling immediately;
- no SSE event contains raw restricted payloads;
- replay from old cursors uses pagination.

## 13. Search architecture

### 13.1 Search is derived

Search results are navigation aids, not authoritative facts. Each result links to a journal entity or artifact whose identity can be verified.

### 13.2 Search fields

Index:

- campaign ID, title, question, system, status;
- event ID/kind/schema/subject/actor/tags;
- snapshot, patch, candidate, trial, episode, observation, epoch, estimate, decision IDs;
- known safe textual summaries from decoded payloads;
- artifact digest/schema/media type/sensitivity;
- case and arm IDs;
- construct, instrument, protocol, status.

Never index confidential/restricted payload text by default.

### 13.3 Query language

Start with a small parser rather than passing raw FTS syntax:

```text
status:completed system:numbergame multiplier
kind:EpisodeFailed
construct:target.accuracy status:violated
case:case-03 arm:challenger
schema:optkit.observation/v1
actor:demo
```

Pseudocode:

```text
tokens = lex(query)
for token:
    if token is known_field:value:
        add structured predicate
    else:
        add escaped text term
reject unknown fields and excessive token counts
execute bounded derived-index query
return typed results with href and explanation
```

### 13.4 Initial implementation choices

- **Slice 1:** structured filters and in-memory text matching over small rebuilt summaries.
- **Slice 2:** rebuildable SQLite FTS5 tables for larger local stores.

Do not block the first UI on FTS.

## 14. Artifact safety

### 14.1 The store does not enforce viewer authorization

`artifact.Store.Open` accepts a complete `artifact.Ref`; it does not decide whether a browser may see it. `artifact.Ref` at `artifact/ref.go:14` carries sensitivity: public, internal, confidential, or restricted.

The query plane must add policy.

### 14.2 Preview policy

Recommended local default:

| Sensitivity | Metadata | JSON preview | Raw download |
|---|---|---|---|
| public | yes | yes | optional |
| internal | yes | yes, bounded | optional |
| confidential | yes | denied by default | denied |
| restricted | minimal | denied | denied |

### 14.3 Safe preview algorithm

```text
resolve digest to an artifact ref already reachable from visible evidence
validate media type and schema
apply sensitivity policy
reject size > preview ceiling before allocation
open through artifact.Store
verify digest/size while reading bounded bytes
decode only registered JSON schema or return escaped text metadata
redact configured fields
return preview with truncation and policy metadata
```

Never accept an arbitrary filesystem path. Never render artifact HTML directly. Never inject raw JSON as HTML.

### 14.4 Reachability

An artifact should be previewable only if it is reachable from a visible campaign/event/entity or an explicit safe catalog entry. Possessing a digest is not sufficient authorization.

## 15. Frontend architecture

The implemented frontend is intentionally dependency-free: one semantic HTML document, one CSS file, and one modern browser JavaScript script. There is no React, TypeScript compiler, package manager, framework router, or client state library.

```text
internal/web/static/
  index.html    document landmarks, search, app root, evidence drawer
  styles.css    tokens, responsive layout, tables, charts, accessibility
  app.js        fetch client, hash router, derived Numbergame views, SSE
```

`app.js` keeps one small in-memory state object:

```text
campaign list
selected campaign summary
selected campaign events
current search text
active EventSource
```

Authoritative state is always refetched from `/api/v1`. JavaScript derives the current Numbergame lineage, trial matrix, estimate chart, decision, budget bars, event inventory, and timeline from returned journal payloads.

### 15.1 Visual language

The accepted style is a flat, white, retro-monochrome Macintosh-inspired scientific document:

- pure white page and surfaces; no paper/beige simulation;
- black and gray rules, text, progress patterns, and table structure;
- no window chrome, menu bar, title-bar imitation, drop shadows, gradients, or rounded card stack;
- modern system sans stack (`Inter`, `ui-sans-serif`, `-apple-system`, `Segoe UI`);
- modern system monospace stack for IDs, schemas, counts, and values;
- blue, green, red, amber, and violet appear only as text accents for links/status/diagnostics;
- information remains complete in monochrome and does not rely on color;
- large editorial headings, compact evidence typography, generous white space.

### 15.2 Route map

The current dependency-free router uses URL hashes so the embedded server only needs to serve `/`:

```text
#/campaigns
#/campaigns/:campaignId
```

Future hashes can add entity routes without server fallback complexity:

```text
#/episodes/:episodeId
#/observations/:observationId
#/artifacts/:digest
#/compare?campaign=A&campaign=B
```

Selection and filters belong in the URL when they must be shareable. The initial global search is transient and filters the visible campaign/event JSON.

### 15.3 State and rendering rules

```text
fetch(): server read models
hash URL: selected campaign
plain object: current response cache and search text
DOM: disposable rendered projection
never client state: authoritative campaign lifecycle or copied journal authority
```

Dynamic values are escaped before HTML insertion. JSON payloads are assigned with `textContent`. CSP prohibits inline scripts and inline styles. Native `<progress>` elements provide data-driven bars without style attributes.

## 16. Navigation design

### 16.1 Persistent shell

```text
OPTKIT       [global search________________]       VERIFIED
─────────────────────────────────────────────────────────────
breadcrumbs

campaign title / scientific question
stage rail
metrics
lineage
trial matrix
analysis
budgets
fact inventory
timeline

                                      evidence drawer ->
```

This is a flat web document header, not simulated Macintosh window chrome or a menu bar.

### 16.2 Breadcrumbs

```text
Campaigns
  / Numbergame 468f
    / Episodes
      / challenger · case-02
        / trajectory event 3
```

Every entity link shows kind, short ID, and full ID on focus/hover. Copy actions copy the full canonical ID.

### 16.3 Context drawer

Opening evidence should not destroy the scientist's place in a matrix or chart. Use a route-aware side drawer. A “full page” link provides a stable standalone URL.

### 16.4 Keyboard navigation

At minimum:

```text
/       focus search
j/k     next/previous row or event
enter   open selected entity
esc     close drawer
[ / ]   previous/next replay sequence or pair
c       copy selected entity URL
```

Do not hide required actions behind pointer-only gestures.

## 17. Visualization specifications

### 17.1 Campaign stage rail

Show stages as derived status, not editable workflow steps:

```text
Design ✓  Compile ✓  Execute ✓  Measure ✓  Analyze ✓  Decide ✓
```

Each stage links to the event range and relevant projection.

### 17.2 Lineage graph

Nodes: snapshot, patch, candidate, dataset, trial, episode group, observation group, estimate, decision. Use directed edges with a table fallback. Avoid rendering every episode as a graph node by default; expand groups on demand.

### 17.3 Trial matrix

Rows are cases; columns are arms and repeats. Encode status with icon + text + color. Virtualize large matrices and preserve sticky headers.

### 17.4 Timeline

Virtualized event list synchronized with a replay slider. Group by lifecycle/design/scheduling/execution/measurement/analysis. The raw sequence is always visible.

### 17.5 Trajectory span tree

Use parent span IDs to render a tree or waterfall. Preserve event order independent of span grouping. Provide a plain ordered-list fallback.

### 17.6 Measurement QA table

Columns:

```text
case, arm, repeat, construct, status, value,
instrument, protocol, epoch, evidence count
```

Intervention violations and failed measurements get explicit symbols.

### 17.7 Paired-delta chart

Render baseline and treatment points connected per case/repeat. Never render only the aggregate. Show missing pairs as gaps and list the reason.

### 17.8 Budget chart

For each resource, show limit, reserved, committed, available, and violation. Do not show “100% used” as an error when exact planned consumption is expected.

### 17.9 Decision evidence chain

A vertical graph connects decision → estimate → observations → episodes → trajectories/artifacts. This is the highest-value provenance visualization.

## 18. Accessibility and scientific legibility

- Color is never the only carrier of status.
- Every graph has a table/list alternative.
- SVG nodes and edges have accessible names.
- Charts expose underlying values in a table.
- Focus order follows scientific reading order.
- IDs and decimal values use tabular/monospace formatting.
- Times show UTC and optionally local time.
- Truncation is explicit; full IDs are copyable.
- Empty, missing, failed, violated, and inapplicable are distinct.
- Reduced-motion mode disables animated graph transitions and live autoscroll.

## 19. Performance and scale

### 19.1 Backend limits

Suggested defaults:

```text
campaign page:          50
control event page:    200 (hard max 500)
search results:         50 (hard max 200)
artifact preview:      256 KiB
trajectory event page: 500 or full bounded artifact
SSE client queue:      bounded
```

### 19.2 Incremental projection

Cache campaign views by `(campaign ID, head version)`. To update from version N to M:

```text
load cached projection at N
read events after N through M
apply projectors
store cache at M
```

Always support full rebuild for verification.

### 19.3 Frontend

- paginate or virtualize event lists and large matrices before RAG-scale campaigns;
- lazy-load trajectory and nested artifact payloads;
- aggregate graph nodes until expanded;
- refetch the selected campaign after batched SSE facts;
- keep raw artifact blobs out of the long-lived JavaScript state object.

## 20. Failure behavior

The query plane must make uncertainty visible.

| Failure | UI behavior |
|---|---|
| journal verification fails | prominent integrity error; do not claim trusted projection |
| unknown payload schema | show generic event and artifact metadata |
| payload decode fails | diagnostic attached to entity; other navigation remains available |
| artifact missing/corrupt | evidence link marked unavailable/corrupt |
| incompatible epochs | block aggregate comparison and explain |
| incomplete pair | show missing cell and invalid estimate |
| SSE disconnect | retain last sequence, reconnect with cursor |
| search index stale | show indexed-through version and direct entity navigation |
| campaign still running | show partial projection and live cursor |
| confidential preview | metadata-only policy result, not generic 500 |

## 21. Security and privacy

This is a trusted local-development application, not an adversarial custody service, but preserve meaningful boundaries:

- bind to loopback by default;
- require explicit configuration to bind externally;
- GET-only application routes;
- reject unexpected methods through ServeMux;
- set CSP, `X-Content-Type-Options`, and frame policy;
- escape all text;
- no raw HTML artifact rendering;
- no arbitrary path access;
- sensitivity-aware preview;
- bounded request/query sizes;
- context cancellation and timeouts;
- no credentials in API output;
- no chain-of-thought display requirement;
- no mutation “debug endpoint.”

No signatures, keys, or remote custody system are required for this local explorer.

## 22. Development and production topology

### 22.1 Development

```text
GOWORK=off go run ./cmd/optkit demo --store ./tmp/demo --reset
GOWORK=off go run ./cmd/optkit serve --store ./tmp/demo --listen 127.0.0.1:8080
open http://127.0.0.1:8080/
```

Editing `index.html`, `styles.css`, or `app.js` requires restarting `go run` because the files are compiled into the binary by `go:embed`.

### 22.2 Production

```text
GOWORK=off go build -o ./dist/optkit ./cmd/optkit
./dist/optkit serve --store ./data --listen 127.0.0.1:8080
```

The normal binary directly contains HTML, CSS, and JavaScript. There is no frontend generation, package installation, build tag, or runtime static directory. API and static routes are registered before the exact `/` document route.

## 23. Implementation phases

### Phase A — Query interfaces and Numbergame projections

Files likely:

```text
query/model.go
query/repository.go
query/service.go
query/registry.go
query/numbergame.go
projection/lineage.go
projection/episodes.go
store/sqlite/query.go
```

Tasks:

1. Add campaign catalog and paged event interfaces.
2. Implement SQLite adapters without exposing SQL upward.
3. Build campaign summary, lineage, episode matrix, measurement, analysis, and decision views.
4. Add schema registry and Numbergame decoders.
5. Add in-memory tests and replay tests over the captured campaign.

Gate: all Numbergame pages can be produced as Go values from a reopened local store.

### Phase B — GET API and safe artifacts

Tasks:

1. Create `server.Server` with `http.ServeMux`.
2. Add envelope, error, pagination, ETag, and request-ID helpers.
3. Add all core GET endpoints.
4. Add artifact reachability/sensitivity/size policy.
5. Add handler tests using `httptest`.
6. Add architecture test that no handler imports internal SQLite APIs directly.

Gate: API snapshot tests match checked-in fixtures and no mutation method succeeds.

### Phase C — SSE replay/live follow

Tasks:

1. Implement cursor parsing and Last-Event-ID precedence.
2. Replay pages in sequence.
3. Poll journal head for external-process writes.
4. Add heartbeat, flush, cancellation, slow-client bounds.
5. Test reconnect, duplicate tolerance, terminal catch-up, and concurrent writers.

Gate: a test client disconnects at sequence N and receives N+1 without gaps after reconnect.

### Phase D — Search

Tasks:

1. Implement query lexer/parser and structured predicates.
2. Build typed search results over projections.
3. Add sensitivity-safe text extraction.
4. Add bounded in-memory v1 implementation.
5. Add optional rebuildable FTS implementation after profiling.

Gate: every search result deep-links to a resolvable verified entity.

### Phase E — Vanilla HTML/CSS/JavaScript shell and navigation

Implemented v0 tasks:

1. Add semantic `index.html` with header, global search, app root, error template, and evidence drawer.
2. Add white retro-monochrome responsive CSS with modern system fonts and text-only color accents.
3. Add hash routes for campaign list and campaign overview.
4. Add `fetch` API client, campaign/event state, global filtering, drawer interaction, and EventSource client.
5. Escape dynamic HTML and render JSON through `textContent` under a strict CSP.
6. Add loading/empty/error/integrity states.

Remaining tasks:

1. Add entity-level hash routes for episodes, observations, and artifacts.
2. Add j/k row navigation and replay-sequence URL state.
3. Split `app.js` only if its growth begins to obscure feature boundaries; do not introduce a framework by default.

Gate: all Numbergame entities are reachable by keyboard and stable hash URL.

### Phase F — Scientific visualizations

Tasks:

1. Stage rail and overview facts.
2. Lineage graph with table fallback.
3. Trial matrix.
4. Timeline and replay slider.
5. Trajectory span tree.
6. Measurement QA and paired-delta charts.
7. Budget visualization.
8. Decision evidence chain.
9. Cross-campaign comparison.

Gate: a scientist can answer the core questions in Section 1 without opening raw JSON.

### Phase G — Embed, CI, and hardening

Implemented v0 tasks:

1. Embed the three static files directly with `go:embed`.
2. Serve exact `/`, `/static/`, GET API, and SSE from one standard-library handler.
3. Add CSP, nosniff, frame denial, and referrer policy.
4. Add Go API/static/read-only/SSE tests.
5. Validate the UI in Chromium against a real Numbergame store with zero new console errors.
6. Run full Go tests, focused race tests, and golangci-lint.

Remaining tasks:

1. Add automated JavaScript DOM tests only when behavior exceeds the current browser smoke coverage.
2. Add accessibility audit and large-fixture performance test.
3. Add a release-binary smoke job.

Gate: a clean checkout builds one normal binary that serves API, SSE, HTML, CSS, and JavaScript.

## 24. Testing strategy

### 24.1 Projection laws

- rebuilding from zero equals incremental application;
- replay through sequence N never includes facts after N;
- event ordering is strict;
- unknown schemas do not remove events;
- failed decode remains visible;
- projection cache key includes campaign version;
- overview counts agree with episode matrix counts.

### 24.2 API tests

- every route is GET-only;
- malformed IDs and cursors return stable 400 errors;
- unknown entities return 404;
- integrity failure returns an explicit state;
- pagination has no duplicates/gaps;
- ETag changes with campaign version;
- immutable artifact ETag is stable;
- preview policy applies sensitivity and size bounds;
- API never leaks filesystem paths.

### 24.3 SSE tests

- replay from 0;
- replay from middle;
- Last-Event-ID reconnect;
- query/header cursor precedence;
- new external write appears;
- cancellation terminates polling;
- heartbeat does not advance cursor;
- duplicate client delivery is idempotent;
- slow client cannot cause unbounded memory.

### 24.4 Search tests

- known field parsing;
- escaping and unknown-field rejection;
- result type/href correctness;
- sensitivity-safe indexing;
- indexed-through version visibility;
- result limits;
- no raw FTS injection.

### 24.5 Frontend tests

- route/deep-link rendering;
- keyboard navigation;
- graph table fallback;
- matrix status semantics;
- replay sequence reflected in URL;
- SSE reconnect reducer idempotence;
- charts match source table values;
- artifact denial state;
- no mutation API calls in the API client;
- accessibility checks with axe or equivalent.

### 24.6 End-to-end review script

```text
run Numbergame through CLI
capture campaign ID
start read-only server
open campaign deep link
verify overview = completed / 8 episodes / 81 events
navigate baseline -> patch -> challenger
open case-02 challenger episode and trajectory
open intervention observation
open paired estimate
open decision evidence chain
replay to pre-decision sequence
reconnect SSE at known cursor
search for target.accuracy
verify no mutation route exists
```

## 25. Intern implementation order

An intern should not begin with graph styling. Use this order:

1. Read `campaign/journal.go`, `campaign/event.go`, and `projection/overview.go`.
2. Run Numbergame and inspect the evidence capture.
3. Read `examples/numbergame/demo.go` from `RunDemo` through restart/analysis.
4. Implement pure Go read models with table-driven tests.
5. Add repository interfaces and SQLite adapters.
6. Add GET handlers and snapshot fixtures.
7. Add SSE replay.
8. Build the plain HTML/CSS/JavaScript shell against a real Numbergame API.
9. Add one visualization at a time with a table fallback.
10. Keep static assets directly embedded; add tooling only when an observed need justifies it.

## 26. Code review guide

Review in this order:

1. **Authority boundary:** Can any UI/API path write journal, queue, budget, or artifact state?
2. **Projection correctness:** Can every displayed claim be traced to event/artifact evidence?
3. **Missingness:** Are absent, failed, violated, and unknown distinct?
4. **Pagination/replay:** Can events be skipped or reordered?
5. **Artifact safety:** Can a digest/path bypass reachability or sensitivity policy?
6. **URL stability:** Is every important entity deep-linkable?
7. **Accessibility:** Is every visualization understandable without color or pointer?
8. **Performance:** Are pages, previews, queues, and frontend collections bounded?
9. **Storage independence:** Do browser/query semantics avoid concrete table coupling?
10. **Scientific semantics:** Are epoch compatibility, paired missingness, and intervention integrity visible?

## 27. Decisions

### Decision: read-only UI

- **Context:** CLI workflows operated by LLM agents own campaign creation and mutation.
- **Decision:** expose navigation/search/visualization only; no browser mutation endpoints.
- **Consequence:** CLI must return IDs/URLs; query projections must be rich.
- **Status:** accepted

### Decision: journal sequence is replay cursor

- **Context:** campaign events are ordered and hash chained.
- **Decision:** use `ControlEvent.Seq` for pagination, SSE IDs, and replay cutoffs.
- **Consequence:** projections support exact historical views.
- **Status:** accepted

### Decision: poll journal head for live updates

- **Context:** writers may be external CLI/worker processes.
- **Decision:** begin with bounded SQLite head polling, not in-memory-only pub/sub.
- **Consequence:** slight local latency in exchange for cross-process correctness.
- **Status:** accepted

### Decision: schema registry for domain payloads

- **Context:** generic events point to domain-specific artifacts.
- **Decision:** register explicit versioned decoders/projectors; unknown schemas degrade generically.
- **Consequence:** Numbergame and RAG-TTC can add plugins without coupling campaign core.
- **Status:** accepted

### Decision: standard-library HTTP and embedded vanilla web assets

- **Context:** Optkit is a Go CLI, the explorer is read-only, and the desired UI is simple HTML and JavaScript.
- **Decision:** use Go 1.22+ `http.ServeMux`, plain semantic HTML, plain CSS, plain JavaScript, hash navigation, Fetch, EventSource, and direct `go:embed`.
- **Consequence:** there is no Node dependency, frontend framework, generated asset pipeline, or second development process; restarting Go recompiles assets.
- **Status:** accepted

### Decision: artifact reachability and sensitivity policy

- **Context:** CAS access itself does not authorize browser disclosure.
- **Decision:** preview only bounded, reachable artifacts under sensitivity policy.
- **Consequence:** query server maintains a derived artifact reference catalog.
- **Status:** accepted

## 28. Risks and mitigations

### Risk: projection becomes a second authority

Mitigation: rebuild tests, cache by journal version, and never accept UI writes.

### Risk: Numbergame-specific assumptions leak into generic API

Mitigation: generic models plus schema-owned projector plugins.

### Risk: large campaigns overwhelm event/graph rendering

Mitigation: pagination, virtualization, grouping, lazy expansion, bounded previews.

### Risk: live stream misses external writes

Mitigation: poll authoritative journal head and replay by sequence.

### Risk: scientists trust an aggregate despite invalid interventions

Mitigation: intervention QA gate and explicit validity panel before charts.

### Risk: search leaks sensitive text

Mitigation: sensitivity-aware extraction and metadata-only indexing.

### Risk: browser code duplicates scientific computation

Mitigation: backend/journal artifacts return authoritative estimates and raw paired points; JavaScript only adapts them for display.

### Risk: deep links break after UI refactor

Mitigation: stable entity routes independent of component structure and snapshot route tests.

## 29. Open questions

1. Should the first server support multiple local store roots or exactly one configured root?
2. Should internal artifact JSON previews be enabled by default or require a flag?
3. What is the largest expected RAG-TTC campaign event count in local use?
4. Is cross-campaign comparison required in the first frontend milestone?
5. Should search v1 include FTS, or are structured filters plus ID/schema matching sufficient?
6. Which graph/chart library best meets accessibility and bundle-size constraints?
7. Should terminal SSE connections close after `campaign-caught-up` or remain open for late facts?
8. How should a CLI discover/configure the explorer base URL for handoff links?
9. Which RAG trajectory payloads require redaction before the UI phase reaches P6?
10. Should artifact export be omitted entirely from v1 and added after preview policy review?

None blocks pure Numbergame projection work.

## 30. Definition of done

The explorer milestone is complete when:

- a CLI-created Numbergame campaign appears without manual database configuration;
- campaigns can be listed, searched, and deep-linked;
- overview, lineage, trial matrix, timeline, episodes, trajectories, measurements, analysis, decision, budgets, and safe artifacts are navigable;
- replay works at an exact sequence;
- SSE reconnects without gaps;
- all visualizations have accessible table/list alternatives;
- the API contains no mutation routes;
- confidential/restricted artifacts cannot be previewed by default;
- checked-in fixtures support frontend development without a live backend;
- frontend and backend tests pass;
- one embedded Optkit binary serves SPA and query API;
- an agent can print a campaign URL and a scientist can trace the decision to evidence.

## 31. Primary file references

### Campaign and storage

- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/campaign/journal.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/campaign/event.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/campaign/reducer.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/store/sqlite/journal.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/local/profile.go`

### Evidence and science

- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/artifact/ref.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/episode/types.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/episode/writer.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/measure/observation.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/measure/epoch.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/experiment/trial.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/experiment/estimate.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/budget/types.go`

### Numbergame and current query surface

- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/examples/numbergame/model.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/examples/numbergame/system.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/examples/numbergame/measure.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/examples/numbergame/demo.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/projection/overview.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/cmd/optkit/main.go`
- `../sources/01-numbergame-evidence.md`
