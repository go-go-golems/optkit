---
Title: 'Architect Brief: Modular Optimization Workbench and Candidate Authoring'
Ticket: OPTKIT-011
Status: active
Topics:
    - design
    - optkit
    - rag-ttc
    - ui
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/Downloads/Architect Brief_ Modular Optimization Workbench + Candidate Authoring.md
      Note: Original architect brief imported verbatim into the ticket design-doc collection
ExternalSources: []
Summary: Architect review of OPTKIT-011, including architecture decisions, proposal compiler boundaries, phased tickets, and recommended PR cadence.
LastUpdated: 2026-08-26T14:11:31.522772698-04:00
WhatFor: ""
WhenToUse: ""
---


# Architect Brief: Modular Optimization Workbench + Candidate Authoring

## Mission

Turn the OPTKIT-011 design into an implementation architecture and a sequence of independently shippable tickets.

The end state is not merely a registry-backed form builder. It is a reusable **optimization workbench** where an engineer can:

**inspect a failure → form a hypothesis → mutate a meaningful optimization coordinate → understand recomputation/cost → preview where honest → seal a candidate → execute a trial → compare outcomes against the original hypothesis and risks.**

That workflow is the product requirement behind the two design documents.

The architecture must make this workflow reusable across retrieval, fusion, representations/prompts, reranking, policy, numbergame, and future optimization surfaces without either:

1. building every experiment UI from scratch, or
2. creating an over-generalized server-driven UI DSL.

---

## 1. Ground truth to preserve

### Optkit owns experiment mechanics, not RAG semantics

Keep the existing repository boundary:

- `optkit/` owns generic variables, lenses, domains, patches, snapshots, candidates, journals, trials, measurements, etc.
- `rag-ttc/` owns RAG layer configuration, which variables exist, pipeline graph composition, previews and specialist/workbench behavior.
- `ragkit/` stays algorithm-focused unless a later architectural decision deliberately promotes the layered-config abstraction.

The existing design states the placement rule well: variable definitions stay with the configuration they mutate; generic variable machinery lives below its consumers.

### Preserve the existing mutation algebra

`optkit/space` already has the right foundation:

`Variable[C,V] → Lens + Domain + Codec → PatchBuilder → child Snapshot`.

Patches already record the exact variable assignments and produce content-addressed child snapshots.

Do not replace this with a second ad-hoc mutation mechanism for the UI.

### Preserve the graph/invalidation model

`rag-ttc/pkg/ttc/optimization` already answers a critical question:

> If this configuration changes, which layers change directly, which become invalid through dependencies, and which can be reused?

That is the source of truth behind the recomputation strip/bill in the UI.

### Persist experiment meaning

Once an experiment is created, the store is canonical. The manifest is authoring input and must not be needed later to explain the experiment. The read side projects recorded facts and does not reconstruct/recompute missing historical data.

Candidate hypothesis, expected improvement, risks, motivating evidence, patch and child configuration therefore belong in durable experiment state.

### One vocabulary for playing and proposing

The strongest UI invariant is:

> A knob in a lab and a mutation in a real candidate refer to the same registered optimization variable.

The fusion `k` slider is `fusion.rrf_k`; a retrieval cut is its registered variable; a prompt editor operates on the registered prompt asset. Playing with a configuration and sealing a candidate are two modes over the same mutation model.

---

# 2. Architecture decisions to settle first

Before feature implementation, produce short ADRs for the following. Do not allow these decisions to emerge accidentally across several PRs.

## ADR A — What is the patchable “whole system” configuration?

Current code has an important asymmetry:

- optkit snapshots currently contain `RetrievalConfig`;
- the optimization graph describes all twelve pipeline layers;
- only retrieval has a real operator-editable semantic configuration;
- the remaining graph layers are largely identities.

The design document explicitly identifies this limitation.

### Recommended direction

Introduce a patchable RAG optimization configuration representing the complete pipeline, conceptually:

```go
type PipelineConfig struct {
    Corpus          ...
    Chunking        ...
    Representations ...
    Embeddings      ...
    Indexes         ...
    Retrieval       RetrievalConfig
    Fusion          FusionConfig
    Reranking       ...
    Evidence        ...
    Context         ...
    Answer          ...
    Judge           ...
}
```

Frozen layers may initially contain opaque/versioned configuration values.

The important relationship becomes:

```text
PipelineConfig
      │
      ├── semantic values operators mutate
      │
      └── derive
           ↓
    optimization.Graph
```

rather than maintaining executable config and manually authored graph descriptions as partially independent structures.

Then optkit naturally patches `Snapshot[PipelineConfig]`, and RAG-TTC derives the new graph and invalidation plan.

If this recommendation is rejected, the ADR must specify the alternative aggregate model used by `ApplyMutations` and prove that layer configs and graph identities cannot drift.

---

## ADR B — Serializable catalog versus executable bindings

The proposed `Section` / `VariableDescriptor` model is correct as the **serializable face** of optimization variables. Typed variables must remain the actual mutation mechanism. This descriptor/closure split is already explicitly intended by the design.

But the architect must fill in the missing bridge between:

```text
HTTP/YAML:
{variable: "fusion.rrf_k", value: 20}
```

and:

```go
space.Set(builder, FusionRRFK, 20)
```

Define two related concepts:

```text
Catalog
    serializable metadata
    sections
    descriptors
    documentation
    domains

Bindings[C]
    VariableID → executable typed mutation adapter
```

The binding layer must decode an untyped submitted value through the variable's real codec/domain and apply the corresponding typed `space.Variable[C,V]`.

There should be one declaration/construction site per variable so descriptor, domain, codec and implementation cannot silently diverge.

---

## ADR C — Domain/value schema

Fix domain serialization before treating it as an API contract.

Current `ChoiceDomain.Descriptor()` exposes labels but loses underlying choice values. A browser cannot safely reconstruct:

```text
"Small seeded noise" → "small"
```

from that descriptor.

Define a proper discriminated value specification covering at least:

- integer ranges;
- float ranges;
- boolean;
- choices as `{value,label}`;
- strings;
- asset references.

Keep stable machine IDs separate from human labels.

Recommended convention:

```text
id:    "fusion.rrf_k"
key:   "rrf_k"
label: "RRF rank constant"
```

Avoid overloading `Name`, `ID`, and variable key.

---

## ADR D — Candidate/prose identity

Adding structured candidate fields changes content identities.

Explicitly decide which of these participate in candidate identity:

- hypothesis;
- expected improvement;
- regression risks;
- proposer;
- motivating evidence;
- strategy;
- timestamps;
- documentation.

The existing design explicitly flags prose identity as unresolved; do not resolve it silently in implementation.

---

## ADR E — Catalog provenance

Recommended addition beyond the current design:

Give the optimization catalog a content/version identity and persist the relevant catalog ID or descriptor snapshot with sealed experiments.

Otherwise an old candidate may later be rendered using the *current* documentation/domain for a variable whose meaning changed since the experiment ran.

Use:

```text
current authoring → current catalog
sealed experiment → sealed/catalog version used when authored
```

This is consistent with optkit's content-addressed provenance model.

---

## ADR F — Read API versus write/workbench API

The current `specialistapi` is deliberately read-only.

Preserve that property unless there is a strong reason to merge responsibilities.

Recommended split:

```text
specialistapi
    durable historical projections
    comparisons
    pipelines
    provenance

workbench/application service
    compile mutation draft
    preview
    seal candidate
    plan/run experiment

HTTP adapter
    exposes workbench commands when needed
```

The command/application layer should exist independently of HTTP so manifests, CLI and browser authoring can all share it.

---

# 3. Core abstraction to build around

The central new application operation should be a **proposal compiler**.

Conceptually:

```text
CompileProposal(
    parent configuration,
    requested mutations
)
        ↓
decode
validate
normalize
apply typed variables
derive child configuration
derive child graph
compute Diff
compute InvalidationPlan
determine preview capabilities
        ↓
CandidateDraft
```

`CompileProposal` MUST NOT create durable experiment records.

Moving a slider through:

```text
60 → 50 → 40 → 30 → 20
```

must not create five candidates or five journal histories.

A second operation:

```text
SealProposal(
    CandidateDraft,
    hypothesis,
    expected improvement,
    risks,
    proposer,
    evidence
)
```

performs durable actions:

```text
materialize assets
materialize child snapshot
build patch
build Candidate
persist into campaign/spec
emit CandidateProposed
```

This draft/seal distinction should be load-bearing throughout both APIs and React.

---

# 4. UI architecture

Do not make the Go catalog a declarative UI language.

The backend catalog describes **meaning and legality**.

React describes **experience**.

Use two registries deliberately.

## Semantic catalog — backend

Examples:

```text
retrieval.limit
fusion.rrf_k
representations.summary_prompt
```

Each provides semantic information:

```text
kind
domain
label
Short/Long documentation
section
schema
sensitivity / asset information
cost hint
```

## Workbench plugin registry — frontend

Grow the existing:

```ts
layerDiffWidgets
outputWidgets
```

into a workbench extension registry.

Conceptually:

```ts
registerWorkbenchPlugin({
  section: "fusion",

  variableEditors: {
    "fusion.rrf_k": RrfEditor,
  },

  layerInspector: FusionInspector,

  outputRenderers: {...},

  preview: FusionPreview,
});
```

Every extension point must have a generic fallback.

Therefore a newly registered ordinary integer should automatically work without frontend code, while an important domain coordinate can provide an expert visualization.

**Generic semantics; specialized experience.**

The vision explicitly requires distinct stage visualizations while still sharing one optimization vocabulary.

---

# 5. Ticket / phase plan

## Phase 0 — Architecture closure

### Ticket 0.1 — Optimization Workbench ADRs

Produce ADRs A–F above.

Also audit and settle two current semantic inconsistencies:

**`retrieval.limit`:** current runtime passes it as the search request's final/effective result limit. BM25/vector top-K values are separate route settings. Current React copy calling it “candidates kept per retriever” is therefore misleading.

**`RRFConstant`:** runtime uses a positive `float64`; the design sketch currently describes an integer range. Decide intentionally whether RRF becomes integer-valued or the variable system gains a float range.

**Exit criteria**

- ADRs reviewed.
- Stable terminology for section/variable IDs.
- Configuration aggregation decision made.
- Candidate identity policy written down.
- No production behavior changed yet.

---

## Phase 1 — Optkit semantic-variable foundation

### Ticket 1.1 — Sections, catalog and richer value/domain descriptors

Primary files:

```text
optkit/space/section.go
optkit/space/variable.go
optkit/space/domain.go
```

Implement:

- ordered sections;
- Short/Long docs;
- stable variable IDs;
- complete choice values + labels;
- float ranges if ADR requires;
- asset kind/schema support;
- deterministic catalog serialization;
- catalog validation.

Migrate numbergame to declare its existing coordinates through one section as the proof case.

### Ticket 1.2 — Executable variable bindings

Add the generic mechanism for converting serialized requested values back into typed variables.

Must prove:

```text
descriptor exists ⇔ executable binding exists
```

for every registered variable.

Tests should cover bad type, out-of-domain value, unknown variable, duplicate assignment and canonical decoding.

### Ticket 1.3 — Structured Candidate intent

Extend `space.Candidate` with the approved structured fields:

- proposer;
- expected improvement;
- risks;
- motivating evidence;
- strategy/targets migration as decided.

Update identity tests explicitly.

**Phase exit criteria**

A generic test program can enumerate a catalog, submit a serialized mutation by ID, apply it through the real typed optkit variable, produce a patch and inspect the resulting candidate.

---

## Phase 2 — Make RAG configuration honestly mutable

### Ticket 2.1 — Pipeline configuration composition

Implement the ADR-A configuration model.

If using the recommended aggregate, introduce `PipelineConfig` and derive the optimization graph from it.

Add lens composition/lifting support if needed so a variable defined beside `FusionConfig` or `RetrievalConfig` can operate against a whole `PipelineConfig`.

The graph must cease being a second independently editable account of semantic values.

### Ticket 2.2 — Real fusion configuration

Introduce a real:

```go
FusionConfig
```

with `RRFK` using the type/domain decided in Phase 0.

Plumb the value into `pkg/ttc/search` instead of relying on the fixed route value.

Update fixture identities and tests.

State the expected identity migration explicitly: new campaigns will receive different fusion identities; historical recorded stores remain historical facts.

### Ticket 2.3 — First RAG catalog

Register at minimum:

```text
retrieval.<final-limit-name>
fusion.rrf_k
```

with reviewed Short/Long documentation and executable bindings.

Remove hard-coded UI semantics such as `FIELD_LABELS` where catalog data is now authoritative.

**Phase exit criteria**

Changing `fusion.rrf_k` through a registered variable changes the actual search calculation and produces the expected graph `Diff` and `Plan`.

---

## Phase 3 — Proposal compiler

### Ticket 3.1 — CompileProposal application service

Implement one application service shared by manifest, CLI and future browser authoring.

Input:

```text
parent
mutations[]
```

Output should include enough data for UI and sealing:

```text
normalized mutations
before values
after values
parent config/snapshot
child config preview
before graph
after graph
config diff
invalidation plan
validation diagnostics
preview capability information
```

No durable writes.

Tests should make deterministic ordering and no-op behavior explicit.

### Ticket 3.2 — SealProposal

Turn a validated draft into:

- stored mutation values/assets;
- patch;
- child snapshot;
- structured candidate;
- journal event;
- persisted campaign/spec state.

The implementation should use normal optkit mechanics rather than reproducing patch semantics.

**Phase exit criteria**

The same mutation can be compiled many times without durable effects and sealed once into deterministic records.

---

## Phase 4 — Manifest candidate authoring

### Ticket 4.1 — `candidates:` manifest block

Implement the patch-style form from the design:

```yaml
candidates:
- id: rrf-k-20
  parent: baseline
  mutations:
  - section: fusion
    variable: rrf_k
    to: 20
  hypothesis: ...
  expected_improvement: ...
  regression_risks: ...
```

The design's intended composition is parent → mutation → changed layer identity → downstream cascade → persisted patch/candidate.

Use the proposal compiler. Do not write a manifest-specific mutation implementation.

Retain current full-arm form for baseline/backward compatibility.

Strict YAML decoding must continue to work.

### Ticket 4.2 — Campaign persistence and compatibility

Extend `CampaignSpec` so every fact required by future screens survives without the manifest:

- candidates;
- relevant patch records;
- catalog/version identity;
- structured intent.

Old stores must continue to open and render the current arm-description fallback.

**Phase exit criteria**

A campaign authored using `candidates:` can be deleted from its original manifest location and still be completely explainable from the sealed store.

---

## Phase 5 — Read projections + workbench command surface

### Ticket 5.1 — Catalog/read projections

Expose the current catalog for authoring.

Extend comparison projections so candidate-authored treatments can expose:

```text
hypothesis
expected improvement
risks
motivation
patch/mutation summary
```

The current design explicitly expects the comparison screen to put candidate intent beside the resulting evidence.

### Ticket 5.2 — Workbench command API

Wrap `CompileProposal`, preview operations and `SealProposal` behind the chosen command/API boundary.

Do not put business logic in HTTP handlers.

Define authorization/sensitivity behavior before asset editing is enabled.

**Phase exit criteria**

A browser can request catalog metadata, compile a draft and receive its actual invalidation plan without creating a journal event.

---

## Phase 6 — React workbench framework

### Ticket 6.1 — Generic WorkbenchRegistry

Evolve:

```text
src/layerwidgets/index.tsx
```

into a general registry supporting:

- variable editor overrides;
- layer inspectors;
- artifact/output renderers;
- preview renderers;
- generic fallbacks.

Generic editors should cover the ordinary value kinds without custom components.

### Ticket 6.2 — WorkbenchShell

Extract reusable UI structure rather than extending the existing `LabScreen.tsx` indefinitely.

Shared context should include approximately:

```text
case
parent
catalog
draft mutations
before/after graph
diff
plan
candidate intent
preview/evidence
seal state
```

Create reusable pieces such as:

```text
CaseHeader
PipelineRail
MutationEditor
InvalidationStrip
MutationSummary
CandidateIntent
EvidencePane
RiskEvidence
SealBar
```

Keep case/question, active mutation and recomputation consequence visible while drilling into stages.

The existing numbergame Lab should be treated as a reference experience and parity fixture, not as the reusable architecture.

---

## Phase 7 — First complete vertical slice: Fusion RRF

### Ticket 7.1 — Fusion workbench plugin

Make `fusion.rrf_k` the proving variable.

It should exercise:

- catalog discovery;
- generic mutation machinery;
- specialized editor;
- live graph invalidation;
- deterministic preview;
- candidate intent;
- sealing;
- persisted comparison.

The fusion inspector should render actual RRF contribution arithmetic rather than merely showing a slider.

This variable is deliberately valuable as the first slice because it tests a real newly mutable layer rather than extending the already-special retrieval layer.

### Ticket 7.2 — Retrieval Lab integration

Build the first real workbench flow around an existing fixture case.

Use the product shape already described in the UI vision:

```text
case → proposal → pipeline → trial → verdict
```

Reuse current recorded stage candidates/chunk catalog/lineage instead of recomputing historical evidence.

**Phase exit criteria**

A user can open a case, change RRF `k`, immediately understand what is structurally invalidated, preview deterministic fusion behavior, write a hypothesis/risk, seal the candidate, run it and see the resulting comparison against that declared intent.

---

## Phase 8 — Asset-variable proof

### Ticket 8.1 — Representation prompt asset variable

Implement something like:

```text
representations.summary_prompt
```

as a real artifact-valued variable.

The editor should provide:

- old/new text diff;
- artifact identities;
- recomputation bill;
- sensitivity handling.

### Ticket 8.2 — Bounded server preview

Implement a real bounded preview:

> regenerate the representation for one selected chunk

Do not pretend an expensive upstream mutation can be instant.

The design explicitly requires previews to be honest by mutation kind.

**Phase exit criteria**

The generic workbench handles both a cheap scalar mutation and an expensive asset mutation without changing its core workflow.

At this point the abstraction has been proven.

---

# 6. Testing strategy

Every phase should land green independently.

### Optkit

Test:

- lens laws;
- descriptor/binding agreement;
- domain validation;
- catalog deterministic serialization;
- patch canonical ordering;
- candidate identity behavior;
- asset refs;
- backward compatibility.

### RAG-TTC

Test:

- semantic config → graph identity;
- mutation → direct layer change;
- downstream invalidation;
- unchanged upstream reuse;
- RRF runtime parity;
- manifest strict decoding;
- candidate persistence;
- historical-store loading.

### Browser

Keep the existing principle that browser-side deterministic computations are pinned against recorded/exported Go results.

The current numbergame Lab already demonstrated why this matters: recorded precision, not theoretical arithmetic, is the experiment's truth.

Add equivalent golden parity tests for client-side RRF/cut previews.

Do not duplicate expensive backend algorithms in TypeScript merely to create instant-looking previews.

---

# 7. Definition of architectural success

The project is successful when adding an ordinary new optimization variable generally requires:

```text
1. Define the real configuration field.
2. Define one typed optkit variable + descriptor/binding beside it.
3. Register it in the RAG catalog.
4. Optionally add a specialized React plugin if generic rendering is inadequate.
```

It should **not** require modifying:

- the proposal compiler;
- patch semantics;
- the workbench shell;
- manifest mutation parsing;
- comparison semantics;
- generic React controls.

A specialized variable should be able to improve its own editing/inspection experience without forking the overall experiment workflow.

---

# 8. Things not to do

Do not:

- build a JSON server-driven layout/component DSL;
- duplicate optkit patch logic in React or manifest parsing;
- let React determine recomputation from layer order heuristics;
- create durable candidates while users are merely moving controls;
- make current manifest files necessary to inspect old experiments;
- have the read API recompute historical facts;
- hide missing data by substituting defaults or zero;
- describe variables differently in Go, manifest docs and React copy;
- port expensive backend behavior into JavaScript just to create fake instant previews;
- start with prompt assets before the scalar mutation architecture works;
- broaden this work into gate-policy unification or moving optimization into `ragkit` unless an ADR shows it is required for the present architecture.

---

# 9. Recommended PR cadence

Prefer small PRs that prove an invariant rather than vertical mega-PRs.

A good progression is:

```text
PR 1   ADRs + terminology corrections
PR 2   optkit catalog/domain types
PR 3   executable bindings + numbergame proof
PR 4   structured Candidate
PR 5   PipelineConfig / graph derivation
PR 6   real FusionConfig
PR 7   RAG catalog: retrieval + fusion
PR 8   CompileProposal
PR 9   SealProposal
PR 10  manifest candidates
PR 11  campaign/read projections
PR 12  workbench command API
PR 13  React WorkbenchRegistry + generic editor
PR 14  WorkbenchShell
PR 15  RRF vertical slice
PR 16  Retrieval Lab
PR 17+ asset-variable work
```

Some can be combined if the resulting PR still has one reviewable architectural claim.

---

# Architect's immediate deliverable

Before implementation begins, return with:

1. ADRs A–F.
2. Concrete proposed Go types for `Catalog`, variable binding and whole-pipeline configuration.
3. The `CompileProposal` and `SealProposal` service signatures.
4. Proposed durable `CampaignSpec`/Candidate shape.
5. Proposed frontend `WorkbenchRegistry` interfaces.
6. Migration story for existing manifests/stores.
7. A dependency graph of the tickets above.
8. Any disagreement with this brief, stated explicitly before coding.

Once those interfaces are credible, implement the smallest end-to-end proof:

**numbergame registry → RAG `retrieval.limit` → real `fusion.rrf_k`.**

Do not optimize for the largest number of controls. Optimize for proving that one variable can travel correctly through:

```text
definition
  → catalog
  → editor
  → draft compiler
  → typed mutation
  → child configuration
  → graph diff
  → invalidation plan
  → candidate
  → sealed experiment
  → comparison evidence
```

That chain is the architecture.
