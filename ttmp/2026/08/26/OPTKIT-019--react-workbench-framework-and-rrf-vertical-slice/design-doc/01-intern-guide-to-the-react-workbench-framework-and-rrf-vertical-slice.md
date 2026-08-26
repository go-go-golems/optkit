---
Title: Intern Guide to the React Workbench Framework and RRF Vertical Slice
Ticket: OPTKIT-019
Status: active
Topics:
    - architecture
    - design
    - implementation
    - optkit
    - rag-ttc
    - ui
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/design-doc/04-backend-first-optimization-workbench-program-roadmap.md
      Note: Parent program goals dependencies exclusions and exit gates
    - Path: repo://rag-ttc/apps/specialist/web/src/api/specialistApi.ts
      Note: |-
        RTK Query read client extended with workbench endpoints
        RTK Query transport boundary
    - Path: repo://rag-ttc/apps/specialist/web/src/layerwidgets/index.tsx
      Note: |-
        Existing specialized renderer registries generalized by this ticket
        Existing renderer registry precedent
    - Path: repo://rag-ttc/apps/specialist/web/src/screens/ComparisonScreen.tsx
      Note: |-
        Existing historical comparison extended with candidate intent
        Historical verdict/evidence screen
    - Path: repo://rag-ttc/apps/specialist/web/src/screens/LabScreen.tsx
      Note: |-
        Numbergame reference experience but not the reusable shell
        Numbergame reference experience
    - Path: repo://rag-ttc/apps/specialist/web/src/test/lab-sim.test.ts
      Note: Recorded Go parity testing precedent
ExternalSources: []
Summary: Frontend architecture for generic catalog-driven editors, specialized workbench plugins, a reusable shell, and an end-to-end RRF candidate workflow.
LastUpdated: 2026-08-26T14:20:29.258849198-04:00
WhatFor: Teach a React contributor how to build a reusable optimization experience without duplicating backend legality, invalidation, or persistence semantics.
WhenToUse: Implement only after OPTKIT-018 stabilizes catalog, projection, and command API contracts.
---


# Intern Guide to the React Workbench Framework and RRF Vertical Slice

## 1. Executive summary

The React workbench is not a generic JSON form and not a collection of one-off labs. It is a reusable experiment workflow with two extension levels:

- generic controls render ordinary catalog value kinds without new frontend code;
- specialized plugins add expert editing, inspection, output, and preview experiences for important coordinates.

The backend catalog defines meaning and legality. React defines interaction and visualization. The backend compiler returns normalized mutations, graph diff, invalidation, diagnostics, and preview capabilities. React must not infer those facts from layer order or reproduce Optkit patch logic.

The first complete proof is `fusion.rrf_k`: choose a case, edit the rank constant, compile a draft, see actual invalidation, inspect RRF contribution arithmetic, write candidate intent, seal/run, and compare recorded outcomes to the hypothesis and risks.

## 2. Existing frontend architecture

### 2.1 API/store

`api/specialistApi.ts` creates an RTK Query API rooted at `/api/rag/v1` and defines cockpit/comparison/cases/pipeline/provenance reads. `store.ts` installs its reducer/middleware.

Workbench commands may use a second RTK Query API rooted at `/api/rag/workbench/v1` or inject endpoints into one API. Keep cache tags and mutation semantics explicit.

### 2.2 Renderer registries

`layerwidgets/index.tsx:19-29` has:

```ts
layerDiffWidgets: Record<string, ComponentType<LayerDiffProps>>
outputWidgets: Record<string, ComponentType<OutputWidgetProps>>
```

This proves keyed specialization. It currently covers display, not editing or previews.

### 2.3 Hardcoded retrieval copy

`layerwidgets/retrieval.tsx:113-149` declares `FIELD_LABELS` for preparation, route, and limit. Catalog descriptors should replace semantic labels/docs, while a specialized widget may still arrange fields intelligently.

### 2.4 Numbergame lab

`LabScreen.tsx` implements an explorable campaign with local simulation and sealed-record comparison. The OPTKIT-010 diary records an important parity failure: estimates use six-decimal recorded measurements, not theoretical arithmetic. Treat numbergame as a reference experience/parity fixture, not as the reusable workbench architecture.

## 3. Product flow

```text
case → proposal → pipeline/preview → trial → verdict
```

Persistent context:

- selected case/question;
- parent arm/snapshot;
- active mutations;
- direct/transitive recomputation;
- candidate intent;
- seal/run state.

Deep inspection should not make the user forget which question and change they are evaluating.

## 4. State ownership

### Server state (RTK Query)

- current catalog;
- parent/case data;
- compile responses;
- preview responses;
- seal/run responses;
- historical comparison.

### Local draft state

- currently edited raw values before debounced compile;
- selected section/variable;
- hypothesis, expected metric/groups, risks;
- UI station/tab/disclosures;
- transient validation focus.

Do not copy historical API data into Redux slices merely to edit it. Keep draft intent separate and explicitly reset/initialize from selected parent.

## 5. Frontend types

Generate or hand-maintain types according to repository policy, but pin them to backend schemas:

```ts
type ValueSpec =
  | { kind: "int"; integer_range: { minimum: number; maximum: number } }
  | { kind: "float"; float_range: { minimum: number; maximum: number } }
  | { kind: "bool" }
  | { kind: "string"; string?: { pattern?: string } }
  | { kind: "choice"; choices: Array<{ value: unknown; label: string }> }
  | { kind: "artifact_ref"; artifact_schema: string };

type VariableDescriptor = {
  id: string;
  key: string;
  label: string;
  short: string;
  long: string;
  value_schema: string;
  value: ValueSpec;
  default?: unknown;
  sensitive: boolean;
  cost_hint?: string;
  probes?: string[];
};
```

Keep `unknown` at transport edges and narrow by discriminant. Do not use broad `any` for catalog values.

## 6. Generic editor registry

```ts
type EditorProps = {
  descriptor: VariableDescriptor;
  before: unknown;
  value: unknown;
  disabled: boolean;
  diagnostics: Diagnostic[];
  onChange(value: unknown): void;
};

type GenericEditor = ComponentType<EditorProps>;

const genericEditors: Record<ValueSpec["kind"], GenericEditor> = {
  int: IntegerRangeEditor,
  float: FloatRangeEditor,
  bool: BooleanEditor,
  string: StringEditor,
  choice: ChoiceEditor,
  artifact_ref: ArtifactRefEditor,
};
```

Generic controls render:

- label and Short docs always;
- Long docs in an accessible disclosure;
- before/current values;
- domain/choices;
- backend diagnostics;
- missing/default distinctions.

Unknown kinds fail visibly with an unsupported-editor state; never silently coerce.

## 7. Specialized WorkbenchRegistry

```ts
type WorkbenchPlugin = {
  section: string;
  variableEditors?: Record<string, ComponentType<EditorProps>>;
  layerInspector?: ComponentType<LayerInspectorProps>;
  outputRenderers?: Record<string, ComponentType<OutputProps>>;
  previewRenderers?: Record<string, ComponentType<PreviewProps>>;
};

type WorkbenchRegistry = {
  register(plugin: WorkbenchPlugin): void;
  editorFor(descriptor: VariableDescriptor): ComponentType<EditorProps>;
  inspectorFor(layer: string): ComponentType<LayerInspectorProps> | undefined;
  previewFor(schema: string): ComponentType<PreviewProps> | undefined;
};
```

Lookup policy:

```text
specialized variable editor if registered
else generic editor by value kind
else visible unsupported fallback
```

Registration rejects duplicate extension keys in development/tests. Plugins do not redefine value domains or mutation IDs.

## 8. WorkbenchShell

Recommended modules:

```text
src/workbench/
  registry.ts
  types.ts
  WorkbenchShell.tsx
  WorkbenchProvider.tsx
  components/
    CaseHeader.tsx
    PipelineRail.tsx
    MutationEditor.tsx
    InvalidationStrip.tsx
    MutationSummary.tsx
    CandidateIntent.tsx
    EvidencePane.tsx
    RiskEvidence.tsx
    SealBar.tsx
  plugins/fusion/
    index.ts
    RrfEditor.tsx
    FusionInspector.tsx
    RrfPreview.tsx
```

Context shape:

```ts
type WorkbenchContextValue = {
  case: CaseSummary;
  parent: ArmSummary;
  catalog: Catalog;
  rawMutations: DraftMutation[];
  draft?: CandidateDraft;
  intent: CandidateIntentDraft;
  preview?: PreviewResult;
  sealState: "editing" | "compiling" | "sealable" | "sealing" | "sealed";
};
```

Use context for cohesive workflow state, not every transient input keystroke if it causes broad rerenders. Memoize selectors/components where stage data is large.

## 9. Compile interaction

```text
user changes editor
  → update raw local mutation
  → debounce compile request
  → cancel/ignore stale request using RTK Query request identity
  → render normalized value, diagnostics, diff, plan
  → enable seal only when latest response matches current raw draft and Sealable=true
```

Never seal a stale compile response. Compare request/draft digest and current mutation state.

Compile failures do not clear the last successful plan without explanation. Show “checking” or current-error states so the invalidation strip is not mistaken for the latest edit.

## 10. Invalidation visualization

Render `InvalidationPlan.Steps` in server order. Use semantic labels:

- `direct_change`: changed coordinate/layer;
- `upstream_change`: recomputed because an input changed;
- `unchanged`: reusable.

Color is supplemental; every mark includes its word/icon/text. Do not infer plan from canonical layer position.

## 11. RRF plugin

### Specialized editor

The RRF editor may combine number input and slider, but backend descriptor bounds are authoritative. Slider step is presentation, not legality. Always allow exact numeric entry.

Show formula and documentation:

```text
contribution = channel weight / (k + rank)
```

### Fusion inspector

Inputs are actual channel candidates and preview payload. Render one row per fused chunk:

```text
chunk-b  total 0.0325
  bm25 rank 2  → 1/(60+2) = 0.016129
  vector rank 1 → 1/(60+1) = 0.016393
```

Do not calculate from only final ordering if ranks/weights are absent. The preview/backend must provide auditable operands.

### Deterministic preview

If client-side preview remains, pin it to exported Go results and backend schema. Prefer rendering backend preview payload for the first vertical slice to avoid duplicated semantics. A client helper may support smooth interaction but cannot be the value sealed or persisted.

## 12. Candidate intent and sealing

Candidate form fields:

- hypothesis (required);
- expected metric and groups;
- regression risks as ordered editable list;
- motivating cases/evidence;
- proposer displayed from authenticated actor;
- strategy selected/derived from supported options.

On seal:

```text
assert latest draft Sealable
submit parent + mutations + draft digest + intent + idempotency key
lock inputs while request active
on success freeze sealed proposal and navigate/show run action
on conflict reload parent/current state with explicit explanation
```

Do not optimistically fabricate candidate IDs.

## 13. Verdict and risk evidence

Comparison screen renders candidate hypothesis near verdict, expected improvement near measured metric, and risks beside observed regressions/case movements. The system may not automatically prove a risk occurred; distinguish declared risk from observed evidence.

Full-arm/historical campaigns without candidates retain existing descriptions and do not show empty candidate cards.

## 14. Routing and deep links

Suggested route:

```text
/workbench/campaigns/:campaign/cases/:case/proposal?parent=:arm
```

After sealing/running, existing comparison routes remain canonical historical destinations. Preserve campaign/case/parent in URL so reload and keyboard navigation work.

## 15. Accessibility

- labels associated with every input;
- slider always paired with numeric input and keyboard behavior;
- disclosures use buttons and `aria-expanded`;
- diagnostics linked through `aria-describedby`;
- station navigation uses links/tabs with correct semantics;
- color not sole status channel;
- focus moves to first error or sealed confirmation appropriately;
- no horizontal-only drag requirement for cut/slider controls.

## 16. Implementation phases

1. Add API types/endpoints and fixtures from OPTKIT-018.
2. Implement generic value editors and tests.
3. Implement WorkbenchRegistry with fallback/duplicate behavior.
4. Build shell/context and shared components.
5. Integrate debounced compile and stale-response protection.
6. Build RRF editor/inspector/preview plugin.
7. Add intent/seal/run/compare flow.
8. Validate responsive layout, accessibility, parity, and live server behavior.

## 17. Test strategy

### Generic editor contract tests

Render every value kind, bounds/choices/docs, missing/default values, diagnostics, disabled state, and unsupported fallback.

### Registry tests

- specialized override selected;
- generic fallback selected;
- duplicate registration rejected;
- unknown kind visible;
- plugins cannot alter descriptor/domain.

### Workflow tests with MSW

- loading/empty/error catalog;
- compile diagnostics;
- stale response ignored;
- seal disabled until current draft sealable;
- seal success/conflict/retry;
- candidate comparison present/absent;
- deep-link reload.

### RRF parity

Pin contribution operands/totals/order to Go preview/exported campaign data. Compare recorded precision rather than theoretical precision, following OPTKIT-010.

### Commands

```bash
cd rag-ttc/apps/specialist/web
pnpm typecheck
pnpm test
pnpm build
```

Then use a live server in tmux and inspect with Playwright at desktop and narrow viewport. Capture screenshots in this ticket's `various/screenshots/`.

## 18. Risks and review focus

- A “generic” component DSL can leak into catalog metadata. Keep layout in React.
- Local RRF math can drift from Go; parity tests and backend preview are required.
- Debounced compile races can seal stale state.
- Large pipeline/candidate payloads can cause broad context rerenders.
- Current hardcoded retrieval labels must migrate without losing specialized presentation.
- Missing values must never become zero/default silently.
- Numbergame code should not be refactored into the shell by copying its local simulation assumptions.

## 19. Out of scope

- artifact prompt editing/preview (OPTKIT-020);
- server-driven layout DSL;
- gate policy/promotion screens;
- duplicating backend algorithms for expensive previews;
- making every current specialist screen editable.

## 20. Exit criteria

- ordinary registered values render without plugin code;
- specialized RRF plugin improves experience without redefining semantics;
- shell supports case/proposal/pipeline/trial/verdict flow;
- backend draft/plan controls seal eligibility;
- complete RRF candidate can be sealed/run/compared;
- parity, accessibility, deep-link, typecheck, tests, build, live inspection, diary, doctor, and upload pass.

## 21. File reference map

- `apps/specialist/web/src/layerwidgets/index.tsx:1-30` — existing display registries.
- `apps/specialist/web/src/layerwidgets/retrieval.tsx:113-149` — hardcoded semantic copy to migrate.
- `apps/specialist/web/src/screens/LabScreen.tsx:48+` — numbergame reference workflow.
- `apps/specialist/web/src/screens/ComparisonScreen.tsx:47-430` — verdict/config/plan/case view.
- `apps/specialist/web/src/api/specialistApi.ts:51-94` — RTK Query read API.
- `apps/specialist/web/src/api/types.ts:178-203` — comparison transport types.
- `apps/specialist/web/src/test/lab-sim.test.ts` — Go parity precedent.
- `apps/specialist/web/src/test/server.ts:23-43` — MSW fixture server.
