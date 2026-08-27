---
Title: Projector and Specialist UI Handoff
Ticket: OPTKIT-005
Status: active
Topics:
    - optkit
    - rag-ttc
    - implementation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/query/service.go
      Note: Existing bounded read-only query baseline
    - Path: repo://rag-ttc/pkg/ttc/judgeinstrument/instrument.go
      Note: Measurement provenance source
    - Path: repo://rag-ttc/pkg/ttc/optimization/graph.go
      Note: Stable graph identities consumed by future config projectors
    - Path: repo://rag-ttc/pkg/ttc/optimization/invalidation.go
      Note: Stable diff and invalidation contracts
    - Path: repo://rag-ttc/pkg/ttc/search/service.go
      Note: Frozen pipeline stage vocabulary
ExternalSources: []
Summary: Stable Phase 0-2 contracts, projector boundaries, and first specialist UI workflow.
LastUpdated: 2026-08-25T18:50:00-04:00
WhatFor: Guide the next projector and read-only specialist UI implementation without reinterpreting raw events in the browser.
WhenToUse: Read before implementing OPTKIT-004 Phase 6 projectors or the first Phase 7 browser workflow.
---


# Projector and Specialist UI Handoff

## Purpose

OPTKIT-005 completed the contracts that a specialist comparison UI needs before rendering scientific claims. This handoff identifies the stable inputs, the first server-side projections, the read-only navigation path, and the screens that remain blocked on later campaign producers.

The browser must not decode raw trajectory schemas, infer invalidation, or make promotion decisions. Product-owned projectors read verified journals and artifacts, decode known versions, and return bounded views.

## Stable contracts available now

### Retrieval and product evidence

- shared retrieval fixture: `rag.semantic-fixture/v1`;
- shared fixture SHA-256: `2fa045999a8a89039e00dd60b3fec2bc17b732d557eb00746e207620a5fbdc7f`;
- optimization fixture: `rag-ttc.optimization-semantic-fixture/v2`;
- optimization fixture SHA-256: `1bb7d696b7ecc8971de03a8f38e8cacff75eb0b4a86fa5a43063831ae50e6a2d`;
- strict context, answer, judge-lineage, hard-constraint, and primary-metric records;
- exported retrieval stage names and `rag-ttc.retrieval-stage/v1` payload schema;
- strict sealed-answer decoding and immutable historical measurement.

### Layered configuration

The canonical layer order is:

```text
corpus
chunking
representations
embeddings
indexes
retrieval
fusion
reranking
evidence
context
answer
judge
```

Each layer exposes:

- layer name;
- typed value schema;
- local semantic identity;
- direct dependency identities;
- resolved digest including upstream meaning.

`optimization.Diff` reports direct local changes. `optimization.Plan` reports `reuse` or `recompute` per layer with `direct_change`, `upstream_change`, or `unchanged` reasons.

### Measurement

The existing Judgekit adapter exposes:

- sealed answer and deterministic contract artifacts;
- recomputed Judgekit instance artifacts;
- protocol and report artifacts;
- prompt/model/cache/token/duration provenance in Judgekit reports;
- epoch-bound Optkit observations;
- distinct measured, failed, unknown/missing, and inapplicable states;
- repeat and cache-bypass attribution.

## Required first projector package

Implement the next work in RAG-TTC:

```text
pkg/ttc/projector/
  campaign.go
  comparison.go
  pipeline.go
  provenance.go
  config.go
```

Do not add RAG stage names or product schemas to Optkit core.

## Projection contracts

### `CampaignCockpit`

The cockpit answers whether a campaign is trustworthy and where investigation should begin.

```go
type CampaignCockpit struct {
    Schema       string
    Campaign     CampaignSummary
    Integrity    IntegritySummary
    Arms         []ArmSummary
    Cases        []CaseSummary
    Estimates    []EstimateSummary
    Constraints  []ConstraintSummary
    Budgets      []BudgetSummary
    Diagnostics  []Diagnostic
    ThroughSeq   uint64
}
```

Required derivation:

- verify the journal before claiming integrity;
- fold lifecycle and episode state from control events;
- load estimates and observations through known schemas;
- preserve missing and failed states;
- bound every collection and paginate cases when needed;
- expose the sequence through which the projection was built.

### `ComparisonView`

The comparison view answers what changed and whether the observed difference is supported by complete pairs.

```go
type ComparisonView struct {
    Schema       string
    Campaign     string
    Baseline     ArmSummary
    Challenger   ArmSummary
    ConfigDiff   optimization.ConfigDiff
    Plan         optimization.InvalidationPlan
    Metrics      []MetricDelta
    Constraints  []ConstraintComparison
    Cases        []PairedCaseResult
    ThroughSeq   uint64
}
```

Required behavior:

- show direct changes separately from transitive recomputation;
- never present failed or missing observations as zero;
- keep hard constraints separate from scalar metrics;
- identify measurement epochs on every metric population;
- expose baseline and challenger graph IDs.

### `PipelineView`

The pipeline view answers where one case first diverged or lost required evidence.

```go
type PipelineView struct {
    Schema       string
    Campaign     string
    Episode      EpisodeSummary
    Stages       []StageSummary
    Artifacts    []ArtifactEdge
    Failures     []Diagnostic
    ThroughSeq   uint64
}
```

The initial stage rail uses the frozen order:

```text
lexical.raw
lexical.collapsed
lexical.policy_filtered
vector.raw
vector.collapsed
vector.policy_filtered
retrieval.fused
retrieval.augmented
retrieval.policy_recheck
retrieval.reranked
evidence.hydrated
evidence.returned
evidence.admitted
context.constructed
answer.sealed
judge.measured
```

Only stages present in an episode are rendered. Retrieval-only episodes therefore stop at their final evidence stage rather than fabricating context or judge rows.

### `ProvenanceView`

The provenance view answers which identities and artifacts support a displayed fact.

```go
type ProvenanceView struct {
    Schema        string
    Candidate     CandidateSummary
    Graph         optimization.Graph
    Snapshot      SnapshotSummary
    Episode       EpisodeSummary
    Trajectory    ArtifactSummary
    Epochs        []EpochSummary
    Observations  []ObservationSummary
    ArtifactEdges []ArtifactEdge
}
```

Payload text remains behind sensitivity and size policy. The view exposes artifact refs and bounded previews, never filesystem paths.

## First browser workflow

Implement one complete read-only navigation path before adding more screens:

```text
campaign cockpit
  -> baseline/challenger comparison
    -> paired case result
      -> pipeline microscope
        -> provenance and config diff
```

Every route must be deep-linkable. The browser may copy IDs and CLI commands but may not submit campaign, candidate, decision, or promotion mutations.

## Query route proposal

Use GET-only routes under a product-owned version:

```text
GET /api/rag/v1/campaigns/{campaign}/cockpit
GET /api/rag/v1/campaigns/{campaign}/comparisons/{baseline}/{challenger}
GET /api/rag/v1/campaigns/{campaign}/cases?after=...&limit=...
GET /api/rag/v1/campaigns/{campaign}/episodes/{episode}/pipeline
GET /api/rag/v1/campaigns/{campaign}/provenance/{kind}/{id}
```

Requirements:

- standard-library `http.ServeMux`;
- bounded limits with server defaults and maxima;
- sequence or opaque database cursor, never browser array offsets;
- strict path and ID validation;
- no POST, PUT, PATCH, or DELETE route;
- payload previews follow the existing 256 KiB ceiling and sensitivity policy;
- unknown schemas produce typed diagnostics rather than guessed decoding.

## Deterministic UI fixture

The Phase 0 fixture provides identity and lineage contracts but not a complete journal-backed UI campaign. The projector ticket should create a fresh durable fixture containing:

- two arms and at least three paired cases;
- one positive and one negative case delta;
- one unauthorized candidate removed before fusion;
- one reranker degradation;
- request-limit clamping;
- admitted evidence;
- one valid answer contract;
- two judge epochs over one sealed answer;
- one missing or failed observation;
- a config diff and invalidation plan;
- all artifact refs stored in CAS and reachable from verified facts.

The fixture must run without network providers.

## Available versus deferred screens

| Screen | Data available | Decision |
|---|---|---|
| Campaign cockpit | Durable campaign, budget, episodes, observations, estimates | Build next. |
| Comparison matrix | Two fixed retrieval arms, paired estimates, config graph | Build next. |
| Pipeline microscope | Frozen retrieval stages and sealed trajectories | Build next. |
| Provenance inspector | Artifacts, snapshots, trajectories, epochs, observations | Build next. |
| Config diff/invalidation | Phase 2 graph, diff, and plan | Build next. |
| Chunk laboratory | No bundle-build campaign or chunk lineage producer | Defer until Phase 3. |
| Rich rank waterfall | Current fixture is small and fixed-arm | Add after Phase 4 search spaces. |
| Answer/evidence studio | Sealed answer proof exists; full answer campaigns do not | Add after Phase 5. |
| Judge calibration dashboard | Attribution exists; gold calibration projectors do not | Add with calibration campaign work. |
| Pareto frontier | No multi-objective dominance records | Defer until Phase 8. |
| Promotion approval | Manifest is proposal-only; no review/gate authority | Defer until Phase 8 and keep browser read-only. |

## Security and integrity guardrails

- Verify journal and artifact custody before displaying “verified.”
- Enforce source policy before fusion and external reranking in producers; projectors only report that evidence.
- Never expose confidential answer, evidence, prompt, or report content without an explicit authorized preview policy.
- Treat runtime and model identities as research fingerprints, not signatures.
- Keep mutation in CLI/application workflows and agents.
- Rebuild projectors from authoritative facts and test deterministic output.

## Acceptance gate for the first UI ticket

The first projector/UI ticket is complete when:

1. all four projections rebuild deterministically from one verified fixture campaign;
2. pagination remains bounded;
3. unknown and sensitive artifacts produce visible diagnostics;
4. the cockpit-to-provenance path works through deep links;
5. config changes and invalidation reasons match `optimization.Plan`;
6. missing values remain distinct from zero;
7. the browser emits no mutation request;
8. keyboard navigation, semantic HTML, CSP, and console checks pass;
9. no frontend build step is introduced.

## Review order

1. `rag-ttc/pkg/ttc/optimization/contracts.go`
2. `rag-ttc/pkg/ttc/optimization/graph.go`
3. `rag-ttc/pkg/ttc/optimization/invalidation.go`
4. `rag-ttc/pkg/ttc/search/service.go`
5. `rag-ttc/pkg/ttc/optkitcampaign/system.go`
6. `rag-ttc/pkg/ttc/judgeinstrument/record.go`
7. `rag-ttc/pkg/ttc/judgeinstrument/instrument.go`
8. `optkit/query/service.go`
9. `optkit/internal/web/server.go`
