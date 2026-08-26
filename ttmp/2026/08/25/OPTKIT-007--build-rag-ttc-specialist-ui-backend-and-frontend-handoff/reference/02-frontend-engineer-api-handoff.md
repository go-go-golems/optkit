---
Title: Frontend Engineer API Handoff
Ticket: OPTKIT-007
Status: active
Topics:
    - implementation
    - rag-ttc
    - ui
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://rag-ttc/pkg/ttc/specialistapi/http.go
      Note: Frontend HTTP contract
    - Path: repo://rag-ttc/pkg/ttc/specialistapi/pipeline.go
      Note: Pipeline and provenance semantics
    - Path: repo://rag-ttc/pkg/ttc/specialistapi/projector.go
      Note: Campaign and comparison semantics
ExternalSources: []
Summary: Routes, TypeScript sketches, navigation flow, fixtures, and integration constraints for the specialist UI.
LastUpdated: 2026-08-26T00:45:00Z
WhatFor: Let a frontend engineer build the first specialist workflow without reading journals or Go internals.
WhenToUse: Use while implementing campaign, comparison, pipeline, provenance, and config-diff screens.
---


# Frontend Engineer API Handoff

## Start the backend

Build RAG-TTC and create the deterministic campaign:

```bash
cd rag-ttc
GOWORK=off go build -o /tmp/rag-ttc ./cmd/rag-ttc

/tmp/rag-ttc experiment optkit-rag campaign run \
  --manifest assets/configs/experiments/optkit-rag/semantic-limit-v1.yaml \
  --store /tmp/rag-specialist \
  --format json
```

Copy the returned campaign ID, then start the API:

```bash
/tmp/rag-ttc experiment optkit-rag campaign serve \
  --store /tmp/rag-specialist \
  --listen 127.0.0.1:8090
```

Health check:

```bash
curl http://127.0.0.1:8090/api/rag/v1/health
```

The server is read-only. It intentionally exposes no CORS wildcard; use a same-origin development proxy.

## Navigation workflow

```text
CampaignCockpit
  └─ choose baseline + treatment
      └─ ComparisonView
          └─ choose paired case and one episode
              └─ PipelineView
                  └─ open ProvenanceView
```

A comparison case contains both baseline and treatment episode IDs. The UI should offer both pipeline links rather than assuming treatment is always the interesting arm.

## Routes

```text
GET /api/rag/v1/campaigns/{campaign}/cockpit
GET /api/rag/v1/campaigns/{campaign}/comparisons/{baseline}/{treatment}
GET /api/rag/v1/campaigns/{campaign}/cases?baseline={arm}&treatment={arm}&limit=50&after={cursor}
GET /api/rag/v1/campaigns/{campaign}/episodes/{episode}/pipeline
GET /api/rag/v1/campaigns/{campaign}/provenance/episode/{episode}
```

Treat path IDs and cursors as opaque strings. Use `encodeURIComponent` for every path segment and query value.

## TypeScript sketches

These sketches show the stable v1 shape. Generate or refine local frontend types from the archived JSON fixtures rather than guessing optionality.

```ts
type ID = string;
type Digest = string;

type Diagnostic = {
  code: string;
  severity: "info" | "warning" | "error" | string;
  message: string;
  subject?: string;
};

type ArmSummary = {
  id: string;
  snapshot: ID;
  graph_id?: ID;
  mean?: number;
  sample_size: number;
};

type CampaignCockpit = {
  api_version: "rag-ttc.specialist-api/v1";
  schema: "rag-ttc.campaign-cockpit/v1";
  campaign: ID;
  trial: ID;
  status: string;
  manifest_id?: ID;
  integrity: { journal_verified: boolean; error?: string };
  arms: ArmSummary[];
  cases: Array<{ id: string; mode: string }>;
  estimates: Array<{
    baseline: string;
    treatment: string;
    value?: number;
    pairs: number;
  }>;
  budget: {
    campaign: ID;
    violated: boolean;
    resources: Array<{
      resource: string;
      limit: number;
      reserved: number;
      committed: number;
      available: number;
    }>;
  };
  episodes: {
    total: number;
    queued: number;
    active: number;
    completed: number;
    failed_terminal: number;
  };
  diagnostics: Diagnostic[];
  through_seq: number;
};

type ObservationValue = {
  status: "measured" | "failed" | "unknown" | "inapplicable" | string;
  value?: number;
  epoch?: ID;
};

type PairedCaseResult = {
  case_id: string;
  repeat: number;
  baseline_episode?: ID;
  treatment_episode?: ID;
  baseline: ObservationValue;
  treatment: ObservationValue;
  delta?: number;
  pipeline_available: boolean;
  provenance_available: boolean;
};

type CasePage = {
  api_version: string;
  schema: "rag-ttc.campaign-case-page/v1";
  campaign: ID;
  baseline: string;
  treatment: string;
  cases: PairedCaseResult[];
  next_cursor?: string;
  has_more: boolean;
  through_seq: number;
};
```

`ComparisonView` embeds complete `config_diff` and `invalidation_plan` objects. Render their layer arrays in canonical upstream-to-downstream order; do not sort alphabetically.

```ts
type LayerDiff = {
  layer: string;
  before_identity: ID;
  after_identity: ID;
  changed: boolean;
};

type PlanStep = {
  layer: string;
  action: "reuse" | "recompute";
  reason: "unchanged" | "direct_change" | "upstream_change";
  before_resolved: Digest;
  after_resolved: Digest;
};
```

Use these visual distinctions:

- `direct_change`: operator changed this layer;
- `upstream_change`: local config stayed constant but an input changed;
- `unchanged`: existing materialization is semantically reusable if custody exists.

A plan does not prove that a reusable artifact is physically present.

## Pipeline rendering

`PipelineView.stages` contains only stages observed in the sealed trajectory. Preserve server order. Do not create empty context, answer, or judge stages for a retrieval-only episode.

Each stage includes counts, chunk IDs, an artifact ref, and an optional preview:

```ts
type ArtifactSummary = {
  ref: ArtifactRef;
  preview?: unknown;
  preview_available: boolean;
  preview_reason?: "sensitivity_policy" | "size_limit" | "read_failed" | "not_json" | string;
};
```

When `preview_available` is false, display the reason and artifact metadata. Do not retry by constructing a filesystem URL.

## Missing values

Never coerce absent values to zero.

```ts
function renderObservation(value: ObservationValue): string {
  if (value.status !== "measured" || value.value === undefined) {
    return value.status;
  }
  return value.value.toFixed(3);
}
```

Likewise, `delta`, arm means, and estimate values are optional. Zero is meaningful when present.

## Pagination

Use the server cursor verbatim:

```ts
async function nextCases(page: CasePage): Promise<CasePage | null> {
  if (!page.has_more || !page.next_cursor) return null;
  const url = new URL(`/api/rag/v1/campaigns/${encodeURIComponent(page.campaign)}/cases`, location.origin);
  url.searchParams.set("baseline", page.baseline);
  url.searchParams.set("treatment", page.treatment);
  url.searchParams.set("limit", "50");
  url.searchParams.set("after", page.next_cursor);
  return fetchJSON(url);
}
```

Do not calculate the cursor from array length.

## Fetch and error handling

```ts
type APIError = { error: { code: string; message: string } };

async function fetchJSON<T>(input: RequestInfo | URL): Promise<T> {
  const response = await fetch(input, { headers: { Accept: "application/json" } });
  if (!response.ok) {
    const body = (await response.json()) as APIError;
    throw new Error(`${body.error.code}: ${body.error.message}`);
  }
  return (await response.json()) as T;
}
```

Recommended handling:

- 400: invalid deep link or stale cursor; show a recoverable navigation error;
- 404: campaign, arm, or episode unavailable; return to the parent screen;
- 405: frontend attempted a forbidden mutation; treat as an implementation defect;
- 500: projection or custody failure; show the error code and do not fabricate partial scientific claims.

## Caching and refresh

Successful projector responses include an `ETag` derived from campaign ID and `through_seq`, but currently also send `Cache-Control: no-store`. Use `through_seq` to detect whether two views came from the same journal point. Poll the cockpit if live progress is needed; there is no specialist SSE route in v1.

## Archived fixtures

The ticket contains real server responses:

```text
sources/api-fixtures/01-cockpit.json
sources/api-fixtures/02-comparison.json
sources/api-fixtures/03-cases-page-1.json
sources/api-fixtures/04-cases-page-2.json
sources/api-fixtures/05-pipeline.json
sources/api-fixtures/06-provenance.json
```

Use them for Storybook/MSW fixtures and component tests. IDs are opaque and should not be parsed.

## Screen acceptance checklist

### Cockpit

- integrity is visible;
- incomplete/failed episode counts are not hidden;
- budget violation is visually distinct;
- arm means and paired estimates preserve missing values;
- baseline/treatment selection is explicit.

### Comparison

- direct and transitive changes are visually distinct;
- metric pair count is shown;
- each case links to both available episode pipelines;
- missing and failed statuses are not numeric.

### Pipeline

- stages remain in server order;
- count changes and chunk IDs are inspectable;
- unavailable previews show reasons;
- diagnostics are visible.

### Provenance

- manifest, graph, snapshot, episode, result, trajectory, and output IDs are copyable;
- artifact sensitivity and size are visible;
- no local paths are assumed;
- the UI does not imply cryptographic signature semantics.

## Explicitly deferred

Do not build unsupported screens from placeholder frontend state:

- chunk laboratory;
- answer/evidence studio;
- judge calibration dashboard;
- Pareto frontier;
- promotion approval or deployment controls.

Those require later durable producers and projector contracts.
