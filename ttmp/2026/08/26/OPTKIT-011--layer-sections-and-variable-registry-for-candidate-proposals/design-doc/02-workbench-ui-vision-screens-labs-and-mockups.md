---
Title: "Workbench UI Vision: Screens, Labs, and Mockups"
Ticket: OPTKIT-011
Status: active
Topics:
    - design
    - rag-ttc
    - ui
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://rag-ttc/apps/specialist/web/src/screens/LabScreen.tsx
      Note: The implemented numbergame Experiment Lab this vision extends
    - Path: repo://rag-ttc/apps/specialist/web/src/layerwidgets/index.tsx
      Note: The per-layer widget registry the stage tabs build on
ExternalSources: []
Summary: The motivation document for the section/variable/candidate work — the screens it exists to serve, with real screenshots of what is built and detailed mockups of what comes next.
WhatFor: Show why the registry, sections, and structured candidates exist, by showing exactly what a user does with them.
WhenToUse: Read alongside the intern guide (doc 01); this document is the intent, that document is the mechanism.
---

# Workbench UI Vision: Screens, Labs, and Mockups

The intern guide (document 01 in this ticket) specifies mechanism: sections, variables, registries, patch-style candidates. This document is the reason those mechanisms exist. Every abstraction in the guide was reverse-engineered from a screen a person needs, and the screens came from one reframing that occurred midway through OPTKIT-009: the user of this system is not an auditor of experiment records — it is an engineer making a chatbot's retrieval better, whose working loop is *read failures → form a hypothesis → change one thing → re-run → see what got fixed and what broke*. Screens exist to serve that loop. Where a mechanism decision seems arbitrary in document 01, this document is where its justification lives.

The document has three parts: what is already built (with screenshots from the running system), what is designed but not built (with detailed mockups), and the traceability table connecting every UI element to the mechanism it requires.

## Part 1 — What is built

### The specialist screens: from journal-shaped to question-shaped

The read-only specialist UI (OPTKIT-008) began contract-faithful and journal-shaped: correct tables of identifiers. Three passes (OPTKIT-009) reshaped it around questions. The comparison screen now leads with a verdict sentence, explains the configuration change in values, sorts cases by movement, and plots every paired case with the mean as one line among them:

![Comparison screen: verdict strip, authored arm prose, value-level config diff, invalidation plan, paired-case slope graph](/home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/various/screenshots/12-comparison-prose-graphs.png)

Design rules visible above, all load-bearing for what follows: color encodes meaning (purple = the change you made; yellow = recomputation it forced; green = healthy; red = failure; uncolored dither = missing data), every colored mark carries its word, identifiers are folded into disclosures, and every panel opens with one or two plain sentences saying what it is. Those sentences are exactly what the `Short`/`Long` documentation fields on sections and variables generalize: today the prose is hand-written per panel; the registry makes it authored once, next to the code, and rendered everywhere.

### The pipeline as evidence, not metadata

Three recording changes made the pipeline screen answer debugging questions. Each stage shows its scored candidates; a chunk catalog carries the content of everything any stage touched, including chunks the policy filter removed; and each hit names the representation text it matched on:

![Pipeline screen: stage glossary, per-candidate scores, count-flow and chunk-presence graphics](/home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/various/screenshots/15-stage-data.png)

![Chunk catalog: every touched chunk readable in full, including filtered ones](/home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/various/screenshots/16-chunk-content.png)

![Match lineage: vector.raw's candidates with the representation text each hit matched on](/home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/various/screenshots/17-vector-raw-lineage.png)

The lesson these three passes taught, stated once and reused throughout: **the read side never recomputes, so if a screen needs a value, the recording layer must produce it.** Each pass was triggered by a user question the screen could not answer ("what did the stage hold?", "what does that chunk say?", "what led vector.raw here?"), and each required deciding whether the store already had the data (surface it) or never recorded it (change the recording and rebuild).

### The Experiment Lab: experiencing the loop

The strongest reframing came from the request to not just visualize but *experience* experiments. The numbergame lab (`/lab/numbergame`, OPTKIT-010) runs a complete optkit campaign in the browser with exact parity to the Go mechanics: propose a challenger with a written hypothesis, seal the trial plan, execute the eight episodes click by click, watch the journal grow, compute the paired estimate, evaluate the decision gates — then compare against the sealed store, which displays the recorded `space.Candidate`'s hypothesis and regression risks:

![The Experiment Lab: propose, seal, run, estimate, decide, compare against the sealed record](/home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/various/screenshots/18-experiment-lab.png)

Two lab facts matter for the design ahead. First, parity is enforced by tests against the exported store bundle — which caught that estimates are computed from six-decimal *recorded* scores, not exact arithmetic. Second, the moment the sandbox reproduces the sealed record exactly ("reproduced" badge) is the emotional payoff of the whole deterministic-store architecture, and every future lab should build to an equivalent moment.

## Part 2 — What is designed

### The Retrieval Lab

The same do-then-compare pattern applied to rag-ttc. It is feasible with exact parity because the fixture's channel hits are data and everything downstream (collapse, policy filter, RRF, cut, coverage) is deterministic computation. Five stations; station 3 holds per-stage tabs, because each stage is a different kind of object and deserves its own visualization.

```text
┌──────────────────────────────────────────────────────────────────────┐
│ RETRIEVAL LAB          case: q-comparison ▾        arm: sandbox      │
│ "Compare the Blue Ice product with its growing guide."               │
│                                                                      │
│ 1 corpus · 2 proposal · 3 pipeline · 4 trial · 5 verdict             │
│           ────────────                                               │
│ ┌ 3 PIPELINE ──────────────────────────────────────────────────────┐ │
│ │ [channels][collapse][policy][fusion][cut][score]                 │ │
│ │  ▔▔▔▔▔▔▔▔                                                        │ │
│ │  (active stage tab renders here)                                 │ │
│ └──────────────────────────────────────────────────────────────────┘ │
│ candidate: limit 1→2 · sealed   ░░ journal: 14 events ▸              │
└──────────────────────────────────────────────────────────────────────┘
```

The case and its question stay pinned; every tab is always about *this* question. The footer keeps the sealed candidate and the growing journal visible so the campaign framing survives deep dives into a stage.

**Tab: channels — you versus the retrievers.** Before the reveal, the user marks which chunks they would return; then the channels show their hits with the representation each matched, and an agreement line makes the interesting disagreement jump out (vector wants chunk-c because its representation genuinely smells relevant to an inventory question).

```text
│ Your picks first: which chunks answer the question?                  │
│   [x] chunk-a Guide  [x] chunk-b Product  [ ] chunk-c Inventory      │
│                                               [ Reveal channels ]   │
│ BM25 (keyword)                 │ VECTOR (semantic)                   │
│ #1 chunk-a ·9.0                │ #1 chunk-c ·0.90                    │
│    rep-a: "Blue Ice evergreen…"│    rep-c: "inventory table          │
│ #2 chunk-b ·9.0                │           quantity warehouse"       │
│    rep-b: "Blue Ice full sun…" │ #2 chunk-b ·0.80                    │
│ agreement with you: BM25 2/2 ✓   VECTOR 1/2 (it also wants chunk-c) │
```

**Tab: policy — authorization as a toggle.** Flip the audience's roles and watch chunk-c live or die; the negative case's PASS/FAIL flips in the same instant. The removed chunk stays readable (struck, not hidden), and the caption states the ordering rule: filtered before ranking, so nothing forbidden can resurface through fusion.

```text
│ Who is asking?  [x] guide  [x] product  [ ] internal  ← try toggling │
│ chunk-a  role guide    → kept                                        │
│ chunk-b  role product  → kept                                        │
│ chunk-c  role internal → ✂ REMOVED (before ranking)                  │
│          "The inventory table stores quantity and warehouse ids."    │
│ negative case q-policy-negative: forbidden {chunk-c} absent → PASS   │
```

**Tab: fusion — the formula, live.** Stacked contribution bars with the actual recorded numbers, and k as a slider. The caption states RRF's non-obvious property (consensus beats single-channel excellence), and the slider lets you find the k where the ranking flips.

```text
│ RRF: each channel contributes 1/(k + rank).      k = [60]  ◄────►    │
│ chunk-b  ▓▓▓▓▓▓▓▓▓░░░░  0.0325 = bm25 #2 → 1/62  +  vector #1 → 1/61 │
│ chunk-a  ▓▓▓▓░░░░░░░░░  0.0164 = bm25 #1 → 1/61                      │
```

**Tab: cut — the experiment's knob.** A draggable cut line over the ranked list, with the baseline's cut drawn as a ghost, and the coverage consequence recomputed live. This is the limit-1-vs-limit-2 campaign compressed into one interactive moment.

```text
│ Keep how many? limit = [2] ◄────►                                    │
│   #1 chunk-b  0.0325   kept                                          │
│  ───────────── limit-1 would cut here ─────────────                  │
│   #2 chunk-a  0.0164   kept (limit 2)                                │
│ q-comparison (needs BOTH):  limit 1 → 0.500    limit 2 → 1.000       │
```

**Tab: score — the metric as an auditable checklist.** Each required group, its targets, whether final evidence satisfied it — with a deep link to the same measurement in the sealed store, closing the loop between lab and production screens.

### The candidate proposal view

This is where the intern guide's mechanisms become visible furniture. The proposal is the coinvault `candidate.yaml` shape rendered as a form, with the layer strip showing the recompute bill of the mutation:

```text
┌ 2 PROPOSAL ──────────────────────────────────────────────────────────┐
│ proposer  human · you          strategy  manual-coordinate/v1        │
│ parent: limit-1 ──[ mutation: retrieval.limit 1 → 2 ]──▶ challenger  │
│         corpus·chunk·repr·embed·index·[RETR]·fus·rank·evid·ctx·ans   │
│         └──────── reuse ────────┘ └──── recompute (7 layers) ───┘    │
│ HYPOTHESIS (yours)                                                   │
│ │ Keeping 2 candidates per retriever will let multi-source           │
│ │ questions keep both documents in evidence.                         │
│ EXPECTED IMPROVEMENT  metric retrieval.target-coverage               │
│                       groups [multi-source] ▸ predicts: q-comparison │
│ REGRESSION RISKS                                                     │
│  ⚠ more admitted text can distract answer generation   (unchecked)   │
│ MOTIVATING EVIDENCE   q-comparison scored 0.500 in campaign 08864c…  │
│                 [ Seal the proposal → plan 6 episodes ]              │
└──────────────────────────────────────────────────────────────────────┘
```

Risk cards return on the verdict station annotated with whatever evidence the trial produced, so declared risks are never written once and forgotten.

### The mutation editor — where you say what changes

The question that motivated the section/registry design directly: "where do I specify what the candidate actually changes — say the summarization prompt, or the RRF factors?" The answer is a layer-scoped variable editor driven entirely by the registry. Two contrasting cases:

```text
┌ MUTATION ── cheap, downstream-only ──────────────────────────────────┐
│ layer: [fusion ▾]     variable: [rrf_k ▾]                            │
│   rrf_k    60  →  [20]  ◄──────────►      domain 1…200               │
│ cost: corpus…retr reuse (6) │ FUS rank evid ctx ans judge recompute  │
│ preview on q-comparison: k=60 → b,a · k=20 → b,a  (order unchanged)  │
└──────────────────────────────────────────────────────────────────────┘

┌ MUTATION ── expensive, upstream ─────────────────────────────────────┐
│ layer: [representations ▾]  variable: [summary_prompt ▾]  (asset)    │
│ parent sha256:81ffc2…            your edit                           │
│ │ Summarize this chunk in one │ │ List the product names, sizes,   │ │
│ │ sentence for search.        │ │ and care facts as keywords.      │ │
│ cost: corpus chunk reuse (2) │ REPR embed index … recompute (10)  ⚠  │
│ preview: [ Generate rep for chunk-c only ▸ ]  (bounded spot check)   │
└──────────────────────────────────────────────────────────────────────┘
```

Three deliberate properties: the **recompute bill appears at edit time** (where you mutate in the stack *is* the cost of your experiment); **asset variables get a side-by-side diff editor** because a prompt change is a text change and should be reviewed like one; and **previews are honest by kind** — scalar changes recompute instantly client-side, asset changes offer a bounded spot-check rather than a fake instant answer.

Every control in these mockups renders documentation from the registry: the section header shows the layer's `Short` ("Merges the keyword and semantic channels into one ranking") with `Long` behind a disclosure, and each variable shows its `Short` inline — a user who has never heard of RRF reads what k does before touching it. This is the requirement that produced the Short/Long fields in document 01.

### The unifying rule

The lab's knobs and the proposal's mutation space are **the same registry**. The k slider on the fusion tab *is* `fusion.rrf_k`; the cut line *is* `retrieval.limit`; the role toggles are the policy layer's variables. "Seal the proposal" means: freeze your current knob deltas into a patch. Playing and proposing share one vocabulary — which is what makes the lab honest preparation for real experiments rather than a disconnected toy.

## Part 3 — Traceability: UI element → mechanism

| UI element (this doc) | Mechanism (doc 01) |
| --- | --- |
| Section headers with orientation prose; per-control help | `Section.Short/Long`, `VariableDescriptor.Short/Long` (Step 1, 3) |
| Sliders with correct ranges; choice menus | `Domain` on descriptors; registry JSON (Step 1, 3, 5) |
| k slider actually changing fusion | real `FusionConfig{RRFK}` plumbed to the executor (Step 2) |
| Mutation line "retrieval.limit 1 → 2"; layer strip cost | patch records + `optimization.Plan` / `ApplyMutations` (Step 4) |
| Hypothesis / expected improvement / risks / motivating cases on proposal and verdict screens | structured `space.Candidate` fields; manifest `candidates:` block; spec persistence (Step 1, 4) |
| Proposal editor and lab knobs sharing one source of truth | `GET /registry` (Step 5); UI binding (Step 6) |
| Prompt diff editor with spot-check preview | asset variables + representation regeneration producer (Step 7) |

Read column left, and the system is a set of screens a person works in; read column right, and it is the implementation plan. They are the same project.
