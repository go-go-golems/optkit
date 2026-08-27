---
Title: Diary
Ticket: OPTKIT-022
Status: active
Topics:
    - design
    - implementation
    - ui
    - rag-ttc
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-08-26T23:18:00.562906497-04:00
WhatFor: ""
WhenToUse: ""
---

# Diary

## Goal

<!-- What is the purpose of this reference document? -->

## Context

<!-- Provide background context needed to use this reference -->

## Quick Reference

<!-- Provide copy/paste-ready content, API contracts, or quick-look tables -->

## Usage Examples

<!-- Show how to use this reference in practice -->

## Related

<!-- Link to related documents or resources -->

## Goal

Chronological record of the ragttc workbench implementation: the guide
retargeting to pbui 0.8.0 and the build of the product package and evidence
tiles.

## Step 1: Retarget the 022/023/024 Guides to pbui 0.8.0

The greenfield ruling deleted pbui's legacy surfaces (pbui repo
PBUI-ACTIONS-3, commit 6efeaeb → version 0.8.0), so the three
implementation guides were rewritten before any code: this product must be
born kernel-native, and its guides must never name a deleted mechanism.

### Prompt Context

**User prompt (verbatim):** (see OPTKIT-021 diary step 6 — same directive:
"perfect, update the design docs to have the simplification that not having
to worry about lgacy affords us." followed by the five-step execution plan)

**Assistant interpretation:** Rewrite the 022/023/024 guides against pbui
0.8.0 per the approved simplification plan.

**Inferred user intent:** Implementation (which starts immediately after)
follows the guides literally; stale guides would rebuild deleted APIs.

### What I did
- 022 §3: menu mechanics rewritten to kernel resolution (four-state
  availability, no bound verb on disabled rows, ambiguity as data, fresh
  revalidation, primary invocation for click-to-open).
- 022 §4 skeleton: added `pbui/actions.ts` (graph + rules + snapshotFor +
  `workbenchTileContributions()`) and `translators.ts`; runtime signature
  updated; descriptors marked representation-only.
- 022 §5: the case example split into a three-line descriptor and four
  kernel rules (autopsy as `primary: true`; evidence-attach unavailable with
  reason until 023).
- 022 §10/§11: "descriptor tests assert verbs" replaced by GOLDENS-FIRST
  menu row specs (`registry.resolve()` against literal snapshots — the
  pbui-workbench standalone-row-spec pattern); build order step 2 updated.
- 023: variable primary rule wording; translator edge names
  (`ragttc.delta-to-case`, `ragttc.layer-to-section`); SealBar renders the
  resolved `proposal.seal` action (capability + sealability in one rule) and
  clicks go through `performAction` — the hand-rolled stale guard now
  protects only the preview cache.
- 024: vocabulary export declared GENERATED from the registry/graph/
  translators/descriptors (PBUI-ACTIONS-3 Phase B generator, golden-pinned);
  approval implemented as a one-shot capability grant re-resolved through
  the kernel.

### Why
- The guides are the implementation contract; OPTKIT-022 starts from them.

### What worked
- No ADR decision changed — every edit was vocabulary-level substitution,
  which is itself evidence the kernel matches the design's intent.

### What didn't work
- N/A

### What I learned
- The goldens-first inversion is the natural greenfield form of the
  PBUI-ACTIONS-2 equivalence fence: with no old implementation to record
  from, the golden IS the spec rather than the fence.

### What was tricky to build
- The SealBar/kernel split in 023: deciding that fresh revalidation replaces
  the seal-path stale guard (but NOT the compile-cache guard) required
  tracing which staleness each guard actually protects against.

### What warrants a second pair of eyes
- 022 §5's rule ids (`ragttc.case.*`) and action ids (`case.*`) — they
  become frozen wire names under ADR L once implemented.

### What should be done in the future
- Implementation begins at build-order step 1 (scaffold) in the next step.

### Code review instructions
- Diff the three guides at this commit; grep them for `disabledBecause`,
  `conversions`, `activate` — all gone except historical notes.

### Technical details
- pbui 0.8.0 API: `createPbui({registry, defaultEnvironment, actions,
  snapshotFor, translators})`; `metadata.primary`; four-state availability.

## Step 2: The Product, Live — and the Document Host (rag-ttc 2f9bc865a, 935d0c295, 71ed43e29)

The workbench product exists and has been driven end to end against a live
fixture campaign: cockpit → accept-mode arm comparison → worst delta →
challenger autopsy (12 stages, artifact digests) → chunk content, every hop
through kernel-resolved menus and idempotent pointer documents. The Go
document host now serves layouts and pointer documents beside the command
API.

### Prompt Context

**User prompt (verbatim):** (see Step 1 — the five-step execution plan;
mid-implementation the user reported "this looks trash, there's a css issue
or so")

**Assistant interpretation:** Build OPTKIT-022 per the retargeted guide;
fix the visual defect immediately.

**Inferred user intent:** A working, honest product on pbui 0.8.0.

**Commit (code):** rag-ttc 2f9bc865a (product), 935d0c295 (css + chunk
context), 71ed43e29 (workbenchhost)

### What I did
- Scaffolded `rag-ttc/apps/workbench/web` (link: deps on local pbui 0.8.0
  packages + plot; specialist dev-proxy pattern; react dedupe).
- `src/pbui/`: 14 representation-only descriptors; `actions.ts` with the
  type graph (abstract `inspectable`/`watchable`), per-type rules, the
  inherited Inspect + watch-toggle (label function flips on the watched
  fact), the shared tile fragment, and `createSnapshotFor` with a
  fact-composed revision; five translator edges; runtime wiring.
- Verb sink: family dispatch (workbench/navigate/local), pointer-document
  write-or-find (`ensurePointerDocument` + protobuf documentPut), ADR K
  trace records, accept-bridge slot so `compare.with` can arm accept mode
  from outside React.
- Nine tiles over the live specialist projections; failures ranks
  worst-first client-side until the server projection lands.
- Menu goldens written FIRST as the row spec (8 tests) — they passed on the
  first resolve, which is the kernel doing what the guide promised.
- Live smoke via Playwright against a real fixture campaign served on a
  fresh port; three screenshots delivered.
- `pkg/ttc/workbenchhost`: catalog (9 apps, binding rules), strict
  ragttc.focus/comparison validators, revision-guarded file store with SSE
  bumps, HTTP routes at /api/rag/workbench-docs/v1, serve.go mount behind
  `--workbench-docs-store`. Boundary respected: no domain imports.

### Why
- OPTKIT-022's task list; the backend chain and pbui 0.8.0 were both ready.

### What worked
- The goldens-first inversion: writing the row spec before the registry
  existed produced zero iteration on menu behavior.
- The kernel's inherited rules earned their keep immediately: Inspect on
  fourteen types and the watch toggle on five are two declarations.

### What didn't work
- `pnpm dev` initially rendered unstyled ("this looks trash"): only
  `presentation-parts.css` was imported; the workbench split-tree needs
  `pbui/styles.css` + `pbui-workbench/styles.css`. Fixed in 935d0c295.
- First `go get github.com/hyperslop-systems/pbui@latest` failed (private
  repo); pinning datalab's pseudo-version v0.0.0-20260730225710-6f20852567e1
  resolved from the module cache.
- The protojson Mutation oneof flattens (`{"documentPut":…}`, not
  `{"body":{"documentPut":…}}`) — first HTTP test failed with
  `unknown field "body"`.
- Port 8091 was already occupied by an unrelated server, producing
  confusing 404s against a healthy /health — the recorded 8090-conflict
  lesson, one port up. Fresh high port fixed it.

### What I learned
- Walking the live UI found a contract gap no test had: a chunk chip's verb
  carried no campaign/episode, so its pointer document could not name the
  catalog holding the content. The fix is the datalab value rule (the
  presenting component supplies context), and ADR I's focus format gained
  optional `episode`/`chunk` fields — recorded as an amendment in the 021
  guide.

### What was tricky to build
- `useWorkbenchStore` + parsed documents: a selector returning a freshly
  parsed object each call would loop `useSyncExternalStore`; the selector
  returns the stable payload and parsing memoizes on it.
- Accept mode lives in Provider state but the sink lives outside React;
  the accept-bridge slot (the conversationFacts pattern) late-binds
  `context.accept` without a module cycle.

### What warrants a second pair of eyes
- The serve.go change touches the OPTKIT-018 composition (additive: one
  flag, one wrapping mux) — the backend owner should glance at it.
- `workbenchhost.Store` uses protojson files per workbench; if multiple
  server replicas ever share a root, the revision guard is per-process.

### What should be done in the future
- Plot adapter (`renderInteractive`) for the compare slope graph (task
  bubz); frontend sync to the document host replacing localStorage; the
  accessibility sweep (task h3q0); the server-side worst-first projection.

### Code review instructions
- Start: `apps/workbench/web/src/pbui/actions.ts` against the golden spec
  in `src/test/menu-goldens.test.ts`; then `src/sink.ts`; then
  `pkg/ttc/workbenchhost/` with its test.
- Validate: `pnpm typecheck && pnpm test && pnpm build` in the app;
  `GOWORK=off go test ./pkg/ttc/workbenchhost/` in rag-ttc.

### Technical details
- Live evidence campaign: campaign:bb66b363386e1b36edf7aed238b17f9d
  (2 arms, 3 cases, paired Δ 0.1667, journal seq 47 verified).

## Step 3: The Slope Graph — Plot Marks Become Kernel Objects (rag-ttc b1c139589)

With 023 sealed, the 022 tails came due. The first was the compare tile's
missing centerpiece: the paired slope graph as a plot document whose marks
are live presentations. The adapter boundary from §7 now exists (`src/plot/`
is the only importer of plot types), and every point on the graph carries
the same delta menu as the per-case row — right-clicking a dot on a chart
opens the challenger autopsy, because the kernel does not know the dot is a
dot.

### Prompt Context

**User prompt (verbatim):** "make sure to read the ticket design doc and
diaries, and commit at appropriate intervals and keep a detailed diary as
you work (using the diary format from the skill) Print out a brutalist work
slip with the plan / different phases for the ticket. then before stsarting
a phase, plrint a split about the phase, and print one when the phase is
done. (in case you forgot)"

**Assistant interpretation:** Continue the OPTKIT-022 tails (slope graph,
doc-host sync, a11y, worst-first) with the full working discipline: design
doc + diaries first, work slips per phase, diary steps, commits at
intervals.

**Inferred user intent:** The tails should get the same rigor as the main
ticket phases — reviewable history, printed plan, no silent scope.

**Commit (code):** rag-ttc b1c139589 — "workbench-ui: paired slope graph in
the compare tile via the plot adapter"

### What I did
- `src/plot/adapters.ts`: the sole plot boundary. `slopeProjection()`
  projects `PairedCaseResult[]` into a `hyperslop.plot` document — a
  per-case line layer (group=case, color=trend up/down/flat), a bold mean
  line (stat summary mean, §6.6 "mean among its cases"), and an interactive
  point layer — then calls `renderPlot`. Returns the scene plus `plotted`
  and `missing[]` (case ids excluded because a side is unmeasured).
- `src/plot/SlopeGraph.tsx`: renders `PlotHost` with `renderInteractive`
  wrapping every point hit in `<Presentation svg reference={delta…}>`; the
  missing list is named beneath the graphic.
- CompareApp: new "slope" section; case rows now sorted by |Δ| (missing
  last) per §6.6.
- styles: slope tone tokens referencing the existing tone palette; plot
  stylesheet imported in main.tsx (it inherits the pbui tokens on its own).
- Goldens in `src/test/slope-plot.test.ts` (5 tests): measured pairs
  plotted, missing named, 6 interactive points carrying case ids, mean path
  spans both sides, all-missing input yields no plot rather than zeros.
- Live verification: resumed store-live's campaign (6 episodes ran,
  completed, paired Δ 0.1667), opened a live comparison via arm menu →
  "Compare with…" → accept, confirmed 6 `<g data-pbui="presentation">`
  points and the delta menu (Open comparison / Open challenger autopsy /
  Inspect / Add to watchlist) opening from a plot mark.

### Why
- Task bubz's last open piece; §7 makes the adapter boundary a structural
  rule ("nothing else imports plot types"), and the slope graph is the
  proof that plot marks and kernel objects compose.

### What worked
- Goldens-first again: the scene-walk test found the real mark count and
  layer ids on the first run after one grammar fix.
- The plot stylesheet inherits pbui tokens by falling back through
  `var(--pbui-…)`, so the graph matched the workbench theme with zero
  theme-bridging CSS.

### What didn't work
- `mapping: { color: { kind: "constant", value: "mean" } }` — the grammar
  refuses constant color mappings: `channel.type.invalid: Channel color
  requires a field mapping.` Fixed by adding a constant-valued `series`
  column to every row and mapping the mean layer's color to that field.
- Playwright could not click the first slope point: two cases share score
  1.0 at both sides, so their circles coincide and the top one intercepts
  the pointer. Real overlap, not a bug — clicked the unique q-comparison
  point instead.

### What I learned
- Line geoms never carry `interaction` in the scene (only symbols, bars,
  errorbars, boxplots do), so the interactive surface of a slope graph is
  its endpoints — which is also the honest surface: a click means "this
  case on this side".
- `hit.values` is the raw projected row, so the adapter can smuggle the
  case id through the scene without any side table.

### What was tricky to build
- The mean line needed a color identity to be distinguishable, but the
  grammar only maps fields. Symptom: null scene with one error diagnostic.
  Cause: constant refs are legal for most channels but not color. Solution:
  a `series: "mean"` column on every row; the summary stat groups by x and
  the constant column collapses to one line. The categorical color scale
  then carries both fields' values (up/down/flat from trend, mean from
  series) in one domain.

### What warrants a second pair of eyes
- Overlapping equal-score points: ties render as coincident circles and
  only the topmost is clickable. Fine for 3 fixture cases; a jitter or
  dodge position may be wanted for real densities.
- The slope viewport is fixed at 480×250 inside an overflow-x container; a
  resize-observer viewport would track the tile instead.

### What should be done in the future
- Churn triage list (fixed / broke / still-broken) from §6.6 still absent —
  it needs a per-case verdict threshold the read API does not expose yet.

### Code review instructions
- Start at `apps/workbench/web/src/plot/adapters.ts` (the document and
  projection), then `SlopeGraph.tsx` (renderInteractive), then the goldens.
- Validate: `pnpm test` and `pnpm typecheck` in `apps/workbench/web`.

### Technical details
- store-live campaign campaign:3f84c0554dcfe54df51919e994352c1b is now
  COMPLETED (49 events, 6 episodes, paired Δ 0.1667) — it refuses new
  candidates from here; future seal demos need a fresh store.
