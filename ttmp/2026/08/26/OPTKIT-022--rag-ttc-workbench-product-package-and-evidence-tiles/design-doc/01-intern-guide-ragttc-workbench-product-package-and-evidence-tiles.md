---
Title: 'Intern Guide: The RAG-TTC Workbench Product Package and Evidence Tiles'
Ticket: OPTKIT-022
Status: active
Topics:
    - design
    - implementation
    - ui
    - rag-ttc
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://pbui/packages/pbui-workbench/src/apps.ts
      Note: defineApp/createAppRegistry — the app contract every tile implements
    - Path: repo://pbui/packages/pbui-chat/demo/src/apps/createDemoApps.ts
      Note: The modern compact product shape this scaffold copies
    - Path: repo://pbui/packages/datalab-ui/src/components/organisms/ChartPanel/ChartPanel.tsx
      Note: The renderInteractive pattern that turns plot marks into presentations
    - Path: repo://rag-ttc/pkg/ttc/specialistapi/http.go
      Note: The six GET routes the evidence tiles consume
    - Path: repo://rag-ttc/apps/specialist/web/src/components/graphs.tsx
      Note: Slope graph, chunk grid, stage flow — organisms migrating into tiles
ExternalSources: []
Summary: Build the PBUI product package for the RAG optimization workbench and prove it with the read-only evidence tiles (campaigns, failures, judge, autopsy, chunk, compare, inspector, trace, watch) over the existing specialist API, plus the Go workbench document host.
WhatFor: Establish the product scaffold, presentation registry, verb sink, and workbench host so the authoring workspace in OPTKIT-023 only adds tiles and plugins, never architecture.
WhenToUse: Read when implementing the product package or any evidence tile; assumes the accepted OPTKIT-021 contracts.
---

# Intern Guide: The RAG-TTC Workbench Product Package and Evidence Tiles

## 1. What this ticket solves, and for whom

The user of this product is a person improving a chatbot's retrieval. Their
first questions in every session are: which questions is the bot failing, what
did the judge say, where in the pipeline did the right evidence die, and what
does that chunk actually say. Today the specialist SPA answers these with
routed screens. This ticket answers them with **tiles** — independently
placeable, linkable views over the same read-only projections — and in doing so
builds the entire product scaffold (presentation registry, verb sink, workbench
shell, Go document host) that OPTKIT-023's authoring workspace will extend.

The scope split with its siblings:

- **OPTKIT-021** decided the contracts (packaging, document formats,
  vocabulary, verbs). This ticket implements them; it does not reopen them.
- **This ticket** ships the product package plus the *read-only* tiles. It
  needs only `specialistapi` (plus one small additive projection, section 9)
  and can proceed while OPTKIT-015–018 are still under construction.
- **OPTKIT-023** adds the authoring tiles (`catalog`, `proposal`,
  `invalidation`, `preview`, `intent`), which require the OPTKIT-018 command
  API.

**Exit criterion.** A user opens the workbench, sees the triage workspace,
reads the failure gallery worst-first, right-clicks a case to open its autopsy
beside it, clicks a chunk in the presence grid to read its text, opens a
comparison with the slope graph, and inspects any object's identities in the
inspector — with every object presented, every action a serializable verb, and
the layout persisted as a workbench document validated by the Go host.

## 2. Prerequisites and prior contracts

Read first, in this order:

1. `pbui/docs/playbooks/building-a-new-hyperslop-systems-app-on-pbui.md` — the
   canonical product-building playbook; its section 11 checklist encodes the
   expensive mistakes.
2. OPTKIT-021 intern guide — the accepted ADRs. This guide restates only what
   it needs; OPTKIT-021 is authoritative for vocabulary and formats.
3. `pbui/packages/pbui-chat/demo/src/` — the four tutorial apps
   (`inventory`, `sku`, `metals`, `notes`), each demonstrating one mechanism.

Current backend surface this ticket consumes (all GET, all existing):

```text
/api/rag/v1  (rag-ttc/pkg/ttc/specialistapi/http.go)
  campaign cockpit        campaigns, arms, budgets, journal integrity
  comparison              paired means, per-case rows, verdicts
  paired case             per-case detail for both arms
  pipeline                stage records, chunk catalog, layer identities
  provenance              episodes, artifacts, digests
```

The specialist SPA's fixture-replay test setup
(`rag-ttc/apps/specialist/web/src/test/`, MSW over archived API fixtures)
migrates with the organisms; the same fixtures drive tile tests.

## 3. The pbui product model, precisely

Interns new to PBUI: this section is the working model. File references are
into `/home/manuel/workspaces/2026-08-24/use-optkit/pbui/`.

### 3.1 Apps, views, placements

An **app** is a registered kind of tile:

```ts
// packages/pbui-workbench/src/apps.ts
export interface AppDescriptor {
  id: string;                 // stable machine id, mirrored in the Go catalog
  title: string;
  tone: string;               // CSS custom-property reference
  singleton: boolean;         // at most one logical view
  duplicable?: boolean;       // defaults to !singleton
  docBound?: boolean;         // reads view.documents
  bindings?: string[];        // required binding keys
  titleFor?(view: AppView): string;
  group?: string; blurb?: string;
  available?(context: AppAvailability): boolean;
  Component: ComponentType<AppProps>;   // AppProps = { placementId, view }
}
```

A **view** (`AppView`) is a logical instance: which app, its title, its named
document bindings. A **placement** is a leaf of a workspace's binary split
tree showing a view. One view placed twice is a **linked view**: both
placements receive the same `view` object and therefore stay in lockstep. A
tile component holds no state of its own; it is a pure function of
`view.documents` plus product stores (RTK Query caches here).

Registration is an explicit list:

```ts
export const workbench = createWorkbench({
  apps: [...evidenceApps(), ...ambientApps()],
  initial: defaultLayout(),
  onMutate: () => {/* persist via the workbench host */},
  onRejected: (m, err) => log(`layout change refused: ${err.code}`),
});
```

### 3.2 Presentations and verbs

Every domain value a tile renders is wrapped:

```tsx
<Presentation reference={{ type: "case", value: { campaignId, caseId, query } }}
              doc={`<case> ${query}`}>
  <CaseRow ... />
</Presentation>
```

Right-click opens the object menu, resolved by the ACTION KERNEL
(pbui 0.8.0): the product declares rules in `src/pbui/actions.ts`
(`createActionRegistry` over a type graph), each with four-state
availability, presentation metadata, and a `bind` returning the serializable
verb. Unavailable rows render greyed with their reason and carry no bound
verb; a tie renders as a non-executable diagnostic row; clicking any row
re-resolves against a fresh snapshot before delegating. Left-click resolves
the `primary` invocation — the unique available rule marked
`metadata.primary` performs (a case chip's primary is `open.autopsy`);
otherwise the menu opens. The mouse-doc line names what is under the pointer
and what L/R will do, derived from the same resolution.

Verbs are serializable payloads routed by one `onPerform` sink (OPTKIT-021
ADR K). In this ticket the sink implements the **navigation** family
(pointer-document + `view.open`) and the workbench layout family; the draft
and command families arrive in OPTKIT-023, but the sink's family-dispatch
skeleton and the trace record shape are built now.

### 3.3 Accept mode

`pbui.accept({types, prompt, filter})` requests an object of a type; every
matching presentation in every tile lights up until one is clicked or Esc
aborts. This ticket ships two accept flows (compare-arm selection and
watchlist-add); the eight translator edges from OPTKIT-021 §7.4 are
registered in full (with `AcceptChooser` mounted) so OPTKIT-023's flows work
without registry changes.

## 4. Package scaffold

```text
rag-ttc/apps/workbench/web/
  package.json          deps: @hyperslop-systems/{pbui,pbui-workbench,
                        workbench-protocol,plot}, react, @reduxjs/toolkit
  vite.config.ts        same-origin /api proxy; SPECIALIST_API override
                        (copy the specialist SPA's proxy block verbatim)
  src/
    pbui/
      types.ts          PresentationValues (OPTKIT-021 §7 tables) + TONES
      verbs.ts          Verb union + describeVerb prose forms
      registry.ts       createPresentationRegistry — representation only
      actions.ts        the kernel: type graph, contributions per type,
                        snapshotFor (facts + capabilities), workbench tile
                        fragment via workbenchTileContributions()
      translators.ts    the eight OPTKIT-021 §7.4 edges
      runtime.tsx       createPbui({registry, defaultEnvironment, actions,
                        snapshotFor, translators})
      descriptors/      campaign.ts arm.ts case.ts verdict.ts layer.ts
                        stage.ts chunk.ts representation.ts artifact.ts
                        trial.ts episode.ts delta.ts decision.ts
                        journalEvent.ts   (label/describe/tone ONLY;
                        authoring types arrive in 023)
    apps/
      CampaignsApp/  FailuresApp/  JudgeApp/  AutopsyApp/  ChunkApp/
      CompareApp/    InspectorApp/ TraceApp/  WatchApp/
    components/
      atoms/          chips: CaseChip, VerdictBadge, LayerChip, IdChip, …
      molecules/      IdTray, ObservationLine, DiagnosticsList, …
      organisms/      SlopeGraphPanel, ChunkGridPanel, StageFlowPanel,
                      FailureGallery, JudgeTracePanel, ComparisonPanel, …
    api/
      specialist.ts   RTK Query — every /api/rag/v1 call, one file
      workbenchDocs.ts RTK Query — /api/rag/workbench-docs/v1 (host, §8)
    documents/
      formats.ts      ragttc.comparison/v1, ragttc.focus/v1 TS types + guards
    plot/
      adapters.ts     the sole boundary into @hyperslop-systems/plot (§7)
    workbench.ts      createWorkbench({...})
    App.tsx           Provider + Surface + MouseDocLine + AcceptBanner
  test/
    apps.contract.test.ts   parse-based descriptor invariants (§10)
    layers.test.ts          import-boundary enforcement (§10)
```

Layer rules (enforced by test, copied from datalab-ui's discipline):
`pbui/` may not import `components/` or `apps/`; descriptors are pure;
`apps/` are thin containers over `components/organisms/`; all network calls
live in `api/`.

## 5. Descriptors, rules, and the verb sink, by example

One descriptor, complete — representation only (the shape every type
follows):

```ts
// src/pbui/descriptors/case.ts
import type { PresentationDescriptor } from "@hyperslop-systems/pbui";
import type { CaseRef, WorkbenchEnvironment } from "../types";

export const caseDescriptor: PresentationDescriptor<CaseRef, WorkbenchEnvironment> = {
  tone: "var(--wb-tone-neutral)",
  label: (v) => v.query || v.caseId,
  describe: (v) => ({ presentationType: "case", ...v }),
};
```

And the case type's menu, declared as kernel rules in `src/pbui/actions.ts`
(the shape every type follows; `define = defineActions<Values, WorkbenchFacts, Verb>()`):

```ts
const caseRules = [
  define.exact("case", {
    id: "ragttc.case.autopsy", action: "case.open-autopsy",
    scopes: ["workbench"],
    metadata: { label: "Open autopsy", order: 0, primary: true },
    bind: ({ subject }) => ({ kind: "open.autopsy",
      campaignId: subject.value.campaignId, caseId: subject.value.caseId, armId: null }),
  }),
  define.exact("case", {
    id: "ragttc.case.judge", action: "case.open-judge",
    scopes: ["workbench"],
    metadata: { label: "Read judge verdicts", order: 1 },
    bind: ({ subject }) => ({ kind: "open.judge",
      campaignId: subject.value.campaignId, caseId: subject.value.caseId, armId: null }),
  }),
  define.exact("case", {
    id: "ragttc.case.watch", action: "case.watch",
    scopes: ["workbench"],
    metadata: { label: "Add to watchlist", order: 2 },
    bind: ({ subject }) => ({ kind: "watch.add",
      ref: { type: "case", value: subject.value } }),
  }),
  define.exact("case", {
    id: "ragttc.case.evidence", action: "case.attach-evidence",
    scopes: ["workbench"],
    test: ({ snapshot }) =>
      snapshot.product.activeDraftDocId
        ? available()
        : unavailable("no proposal draft is open"),
    metadata: { label: "Attach as evidence", order: 3 },
    bind: ({ subject, snapshot }) => ({ kind: "evidence.attach",
      docId: snapshot.product.activeDraftDocId ?? "",
      ref: { type: "case", value: subject.value } }),
  }),
];
```

Note the last rule: it is *declared now* and greyed with its reason until
OPTKIT-023 provides an active draft — and because the kernel binds only
available winners, the greyed row carries no verb at all. Declaring
unavailable actions with reasons is the pbui norm — the menu teaches the
workflow. `primary: true` on autopsy makes a bare left click on any case
chip open its autopsy through fresh revalidation.

The sink skeleton:

```ts
// src/workbench.ts (excerpt)
function onPerform(verb: Verb): boolean {
  if (isWorkbenchVerb(verb)) return performWorkbenchVerb(handlers, verb);
  switch (family(verb)) {
    case "navigate": return performNavigate(verb);   // §6.1
    case "draft":    return refuse(verb, "authoring arrives with OPTKIT-023");
    case "command":  return refuse(verb, "authoring arrives with OPTKIT-023");
  }
}
// Every call appends {seq, actor: "human", verb, target, outcome} to the trace
// store; refuse() records outcome: "refused" with the reason. Never report
// performed for a verb that touched nothing.
```

`performNavigate` implements the pointer-document idiom:

```text
performNavigate(open.autopsy {campaignId, caseId, armId}):
  docId = deterministicId("focus", campaignId, caseId, armId)
  if workbench document lacks docId:
      mutate documentPut{id: docId, format: "ragttc.focus/v1", body: {...}}
  perform view.open {appId: "autopsy", documents: {focus: docId}}
  # view.open reuses an existing view with identical bindings → no tile spam
```

Deterministic pointer-document ids make "open the autopsy for this case"
idempotent across sessions.

## 6. The evidence tiles

Every app follows the two-halves pattern: a thin container in `apps/`, pixels
in an organism. The organisms are, wherever possible, migrations of specialist
SPA components (ADR G migration rule: move, never fork). Per tile — binding,
data, presentations emitted, and the specialist source it draws from:

### 6.1 `campaigns` — singleton, launcher group "evidence"

The cockpit: campaign list, arms with prose descriptions, budget meters,
journal-integrity status. Emits `campaign`, `arm`, `decision`,
`journalEvent` presentations. Source: `CockpitScreen` +
`CampaignEntry` organisms. Arm rows offer `open.compare` (accepts a second
arm via accept mode — the first accept flow to build).

### 6.2 `failures` — docBound to `ragttc.focus/v1` (campaign-only focus)

The front door. Worst-verdict-first gallery: question text, judge label and
words, score (missing rendered as missing), repeat-failure badge. Emits
`case` and `verdict`. Requires the worst-first projection (section 9); until
that lands, the tile builds against a fixture and renders the
"projection unavailable" diagnostic state honestly.

### 6.3 `judge` — docBound to `ragttc.focus/v1` (case+arm)

The judge trace for one (case, arm): what the judge was shown, its verdict
text, its score at recorded precision. Emits `verdict`, `episode`,
`artifact`. Source: episode/observation rendering from the paired-case
screen.

### 6.4 `autopsy` — docBound to `ragttc.focus/v1` (case+arm)

The pipeline autopsy: stage-flow narrowing line, chunk-presence grid,
per-stage ranked candidates with scores and channels, layer identity strip.
Emits `stage`, `chunk` (via stage-candidate → chunk conversion), `layer`,
`representation`. Source: `PipelineScreen`, `graphs.tsx` (ChunkGrid,
StageFlow), `layerwidgets/retrieval.tsx` (StageRecordView, glossary).

### 6.5 `chunk` — docBound (chunk pointer), duplicable

Full chunk text, document/title/URL, representation lineage (kind, model,
prompt digest). Emits `representation`, `artifact`. Source: chunk catalog
rendering from the pipeline screen.

### 6.6 `compare` — docBound to `ragttc.comparison/v1`

Paired slope graph (mean among its cases), verdict summary, per-case rows
sorted by |Δ|, churn triage list (fixed / broke / still-broken), and — once
OPTKIT-018 lands — the `CandidateSummary` intent panel beside the evidence.
Emits `delta`, `case`, `arm`, `decision`. Source: `ComparisonScreen` and
`SlopeGraph`. The intent panel renders "not a candidate-authored arm" for
old campaigns (absence, never fabrication).

### 6.7 `inspector` — singleton, ambient

Renders `describe()` of the most recently inspected object, with the ID tray:
digests, journal seqs, artifact refs. This tile is where identities live —
the design language's "hashes one disclosure away" becomes "hashes one tile
away". Every other tile's `inspect` verb targets it.

### 6.8 `trace` — singleton, ambient

The performed-verb trace (seq, actor, verb prose via `describeVerb`, outcome)
and, below it, the campaign journal event tail. Emits `journalEvent`.

### 6.9 `watch` — singleton, ambient

A watchlist of accepted references (`watch.add` verb; accepts any type).
Persisted in the workbench document as a small product document
(`ragttc.watchlist/v1` — additive format, registered in the host validator).

## 7. Plot integration

Copy the datalab boundary discipline: one adapter file
(`src/plot/adapters.ts`) projects product data into
`PlotDocument`/`PlotSchema`/`PlotData`; nothing else imports plot types. The
slope graph, stage flow, and chunk grid become plot documents rendered by
`PlotHost` with `renderInteractive`:

```tsx
const renderInteractive = (hit: PlotHit, element: ReactElement) => {
  if (hit.kind !== "mark") return element;
  const ref = referenceForDatum(hit);   // e.g. {type:"delta", value:{caseId,...}}
  return ref
    ? <Presentation svg reference={ref} doc={`<${ref.type}> ${labelFor(ref)}`}>{element}</Presentation>
    : element;
};
```

Every slope-graph mark is then a `delta` object with the full menu (open case,
open both pipelines). Constraints inherited from plot: click-only (hover,
tooltips, brushing are plot Phase 6 and out of scope); missing values are
excluded and named beneath the graphic, never plotted as zero — the existing
specialist graphs already obey this and the adapters must preserve it.

Note: the existing specialist `graphs.tsx` components are hand-rolled SVG. The
migration to plot documents may be done per-graphic; a hand-rolled organism
wrapped in presentations is acceptable for v1 where the plot grammar lacks a
needed form, but new graphics start as plot documents.

## 8. The Go workbench host

Implement OPTKIT-021 ADR H:

```go
// rag-ttc/pkg/ttc/workbenchhost/catalog.go
func DefaultCatalog() Catalog {
    c := Catalog{}
    for _, id := range []string{"campaigns", "inspector", "trace", "watch"} {
        c[id] = workbench.ApplicationDescriptor{ID: id, Singleton: true,
            DocumentBindings: map[string]workbench.BindingRule{}}
    }
    for _, id := range []string{"failures", "judge", "autopsy", "chunk", "compare"} {
        c[id] = workbench.ApplicationDescriptor{ID: id,
            DocumentBindings: map[string]workbench.BindingRule{
                bindingKeyFor(id): {Required: true}}}
    }
    return c
}
```

plus `documents.go` (strict validators for `ragttc.focus/v1`,
`ragttc.comparison/v1`, `ragttc.watchlist/v1`; `ragttc.proposal-draft/v1`
registered here in OPTKIT-023) and `http.go` mounting the workbench document
routes under `/api/rag/workbench-docs/v1/` using the pbui Go module's
`Validate`/`Apply`. The datalab reference:
`datalab/pkg/server/handlers_workbenches.go` (snapshot PUT with `If-Match`
revision, `POST .../mutate` with `MutationBatch`, SSE revision stream) —
follow its shapes, do not invent new ones. Storage: SQLite table in the
existing rag-ttc store directory, schema copied from datalab's workbench
store.

**Boundary reminder (tested):** `workbenchhost` imports the pbui module and
`documents/` formats only. A test asserts it does not import `specialistapi`,
`experimentworkbench`, or `optimization`.

The frontend syncs with the same optimistic-apply + snapshot-PUT pattern as
`datalab-ui/src/appkit/useRemoteWorkbench.ts`; local-storage persistence is
the dev fallback when the host is absent.

## 9. The one backend ask: worst-first verdicts

The `failures` tile needs a projection ordered by verdict severity with
per-case failure counts across epochs. This is an **additive read
projection** on `specialistapi` (recorded facts only — verdicts and scores are
already in episode observations; the projector orders and counts, it never
recomputes). It is the amendment OPTKIT-021 §11 assigns to OPTKIT-018;
coordinate with the backend colleague so the projection ships with the other
OPTKIT-018 read additions. Until then the tile runs on fixtures.

Shape (settled here so both sides build against it):

```text
GET /api/rag/v1/campaigns/{campaign}/failures?arm={armId}&limit=N
→ { cases: [ { case_id, query, description,
               worst: {arm_id, score?, label, missing?},
               failure_count, episode_count,
               verdict_excerpt } ],
    ordering: "worst-first",
    diagnostics: [...] }
```

Missing scores stay missing; a case with no measured episodes appears with
`missing: true` and sorts after measured failures, never as zero.

## 10. Testing strategy

- **Contract tests (parse-based).** Following
  `datalab-ui/test/apps.test.ts`: parse each `defineApp` literal and assert
  the OPTKIT-021 invariants (singleton apps have no bindings; docBound apps
  declare their binding key; ids appear in the Go catalog — checked against a
  generated JSON export of `DefaultCatalog()` committed beside the test).
- **Layer tests.** Import-boundary assertions for the §4 rules.
- **Menu golden tests (the row spec).** Because this product is greenfield,
  the goldens are written FIRST as the specification: for representative
  references of every type, assert the exact resolved rows — ids, labels,
  order, danger, unavailable reasons, and bound verbs (absent on disabled
  rows) — via `registry.resolve()` against literal snapshots. This is the
  standalone-row-spec pattern `pbui-workbench/src/actions.test.ts` ended at;
  there is no old implementation to record from.
- **Descriptor tests.** Pure `label`/`describe`/`tone` calls with literal
  values and environments.
- **Fixture replay.** MSW over the archived specialist API fixtures (reuse the
  specialist SPA's setup); tiles tested in loading, empty, diagnostic, error,
  and populated states.
- **Applier parity.** No local mutation code exists to test — the product uses
  the protocol applier; the shared fixtures in
  `pbui/packages/workbench-protocol/fixtures/mutations/` already pin Go/TS
  parity.
- **Accessibility.** Presentations are keyboard-activatable by contract
  (pbui provides Enter/Space and Shift+F10); tests cover tab order through a
  tile and the accept-mode flow via keyboard.

## 11. Implementation order

1. Scaffold the package (§4), `App.tsx` with an empty app list, dev proxy;
   commit.
2. `pbui/types.ts`, `verbs.ts`, read-side descriptors (representation
   only), registry, `actions.ts` (type graph + rules + snapshotFor +
   `workbenchTileContributions()`), `translators.ts`, runtime; menu golden
   tests written as the row spec; commit.
3. Sink skeleton with navigation family + trace store; `inspector` and
   `trace` tiles (they make everything else debuggable); commit.
4. `campaigns` tile + RTK Query `specialist.ts`; first workspace; commit.
5. Go `workbenchhost` (catalog, validators, routes, storage) + frontend
   `workbenchDocs.ts` sync; layout persists server-side; commit.
6. `autopsy`, `chunk`, `judge` tiles (organism migrations); pointer-document
   navigation verbs end-to-end; commit.
7. `compare` tile with the plot adapter + `renderInteractive` presentations;
   accept-mode arm selection; commit.
8. `failures` tile on fixtures; wire to the live projection when OPTKIT-018
   lands; commit.
9. `watch` tile + `ragttc.watchlist/v1`; contract/layer/a11y test sweep;
   default workspaces (Triage, Autopsy) in `workbench.ts`; commit.

Each step lands green: `pnpm test` in the product, `go test ./...` in rag-ttc
when the host changes, plus the docmgr diary/changelog updates per the
program's delivery discipline.

## 12. Exclusions

- No authoring tiles, no command-API calls, no draft documents (OPTKIT-023).
- No prompt/asset editing (OPTKIT-020 via OPTKIT-023's plugin registry).
- No agent integration (OPTKIT-024); the trace record shape is built
  agent-ready but stays local.
- No deletion of specialist SPA screens until their replacement tile reaches
  parity (ADR G migration rule).
- No plot Phase 6 work (hover/brush); click-only interaction.

## 13. Glossary delta

Terms beyond the OPTKIT-021 glossary: **organism** — a presentational
component owning pixels, no data fetching; **container app** — the thin
`apps/` component mapping `view.documents` + store state to organism props;
**pointer document** — a small `DocumentPayload` whose only job is to give a
doc-bound tile its binding identity; **ambient tile** — a singleton app
(`inspector`, `trace`, `watch`) that reacts to objects and verbs from
anywhere rather than to a bound document.
