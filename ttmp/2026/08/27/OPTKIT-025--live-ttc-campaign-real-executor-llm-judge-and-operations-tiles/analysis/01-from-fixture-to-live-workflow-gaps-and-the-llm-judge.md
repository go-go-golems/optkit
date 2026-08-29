---
Title: 'From Fixture to Live: Workflow, Gaps, and the LLM Judge'
Ticket: OPTKIT-025
Status: active
Topics:
    - rag-ttc
    - design
    - backend
    - ui
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://judgekit/GLOSSARY.md
      Note: The measurement-theory chain the judge instrument is built on
    - Path: ws://rag-ttc/pkg/ttc/optkitcampaign/campaign.go
      Note: scoreObservation (the hardcoded instrument the seam replaces) and budget limits
ExternalSources: []
Summary: The operator's phased workflow for the first real TTC campaign, the exact gap list between today's fixture machinery and a live run, the LLM-as-judge design (judgekit-backed second instrument, budgeted and calibrated from day one), and mockups of the three new operations tiles.
LastUpdated: 2026-08-27T13:05:00-04:00
WhatFor: The implementation contract for OPTKIT-025..027 and the operator's playbook for the first live campaign.
WhenToUse: Before implementing the live executor or judge; during the first real campaign as the workflow checklist.
---


# From Fixture to Live: Workflow, Gaps, and the LLM Judge

Companion to this ticket's intern guide (which explains what exists). This
document says what is MISSING for a real campaign over the real TTC
knowledge base, how the LLM judge is wired in **from the start**, what the
human sees and steers, and how the work slices into tickets.

## 1. The gap list, exactly

| # | Gap | Where it lands | Size |
|---|---|---|---|
| G1 | Live executor + `rag.ttc-live/v1` preparation wrapping the ragsearch composition (index bundle + embedding provider + tool config) | `pkg/ttc/optkitcampaign` (new `live.go`) + a bundle-opening composition helper reusable from `internal/customer/ragsearch` | M |
| G2 | Catalog ↔ `SearchConfig` mapping completed (every catalog variable binds a real field; every field worth optimizing has a variable) | `pkg/ttc/optimization` config/catalog | S |
| G3 | Multi-instrument scoring: `scoreObservation` (one hardcoded instrument) becomes an instrument list; episodes yield one observation PER instrument | `pkg/ttc/optkitcampaign/campaign.go` | M |
| G4 | **LLM judge instrument** (judgekit-backed; §3) | new `pkg/ttc/judgeinstrument` wiring judgekit | M–L |
| G5 | Budget resources beyond retrieval: `embeddings.tokens`, `judge.calls`, `judge.tokens` — reserved before, committed after, from executor/judge usage reports | `optkitcampaign` limits + manifest schema | S |
| G6 | Real case suite with targets (curated queries + expected chunks/docs + judge rubric refs) | `assets/configs/experiments/optkit-rag/ttc-live-*.yaml` | human+agent work, not code |
| G7 | Operations read projections: `GET /campaigns/{c}/work` (queue/leases), budget detail; substrate/provenance endpoint | `pkg/ttc/specialistapi` | S–M |
| G8 | Three tiles: `runner`, `budget`, `substrate` (§5) + Go catalog entries | `apps/workbench/web` + `workbenchhost/catalog.go` | M |
| G9 | Estimates per construct (today `calculateEstimates` assumes one score stream; with two instruments, per-construct means/deltas — the comparison API already speaks "construct") | `optkitcampaign` + `specialistapi` | S |
| G10 | `trial.run` wiring (candidates → arms via the API; today CLI-only `resume`) | `experimentworkbench` + command API | M (can trail) |
| G11 | Judge calibration surface: gold labels, agreement stats, judge-vs-deterministic disagreement (§3.4) | judgekit `calibration` + judge tile extension | M (v1 minimal) |

Deliberately NOT in scope: autonomous trial execution, agent-authored
layouts beyond `view.open`, sandbox/widget vocabulary, the model-backed
chat conversation (release-gated task 75ds in OPTKIT-024 — the scripted
agent seat works today).

## 2. The live executor (G1) — sketch

```text
# pkg/ttc/optkitcampaign/live.go (new)
const TTCLivePreparation = "rag.ttc-live/v1"

type LiveExecutor struct {
    tool     *ttcsearch.SearchTool     # built ONCE per process from:
    #   indexbundle.Open(bundlePath)  → Lexical, Vector, Content, Manifest
    #   raggeppetto.NewBundle(Roles{Embedding:true})  → query embedder
    #   toolconfig.Load(configPath)
    usage    chan ResourceUsage        # per-episode: queries, results,
}                                      #   embeddings.tokens

func (e LiveExecutor) Execute(ctx, config PipelineConfig, case) (SearchOutput, error):
    if config.Retrieval.Preparation != TTCLivePreparation: refuse
    cfg := e.tool.ConfigWith(          # map the optimization knobs onto
        RRFK:            config.Fusion.RRFK,          # the REAL SearchConfig
        FinalResultLimit: config.Retrieval.FinalResultLimit,
        BM25TopK / VectorTopK / Route / LexicalRepresentation: …)
    out := e.tool.RunRoute(ctx, {Query: case.Query, Limit: …}, config.Route)
    report usage {retrieval.queries:1, retrieval.results:len,
                  embeddings.tokens: measured}
    return out
```

Key decisions:

- The executor is **process-local state** (open indexes, credentials) —
  exactly what the `Executor` doc comment reserves. The campaign store
  records the bundle's manifest digest at `CampaignCreated` (substrate
  provenance, G7/G8's substrate tile reads it back).
- Config mapping is TOTAL or refused: a `PipelineConfig` field the live
  tool cannot honor is an error at `Prepare`, not a silent default —
  the runtime-honest catalog rule, enforced at the seam.

## 3. LLM as judge, from day one (G3+G4)

### 3.1 Why a second instrument, not a replacement

Deterministic target coverage is cheap, reproducible, and blind: it cannot
tell "retrieved the right chunk" from "answered the question". The judge
measures a second construct on the same episodes. Two constructs, two
instruments, one epoch discipline — the measure model (intern guide §6.3)
was built for exactly this, and judgekit supplies the missing middle:

```text
construct   retrieval.answer-support            (does retrieved evidence support answering?)
instrument  ttc.judge.retrieval-support/v1      (judgekit spec+protocol, model-pinned)
protocol    complete-block/v1                   (same block design as coverage)
epoch       frozen (construct, instrument, protocol, implementation,
            model id, prompt digest)            → scores never silently mixed
```

### 3.2 The instrument seam (G3)

```text
# optkitcampaign: replace the hardcoded scoreObservation with:
type Instrument interface {
    Definition() measure.EpochDefinition        # names + versions
    Score(ctx, spec EpisodeSpec, case RetrievalCase,
          output SearchOutput) (measure.Observation, Usage, error)
}
RunOptions.Instruments []Instrument             # default: [Coverage{}]

per episode:
  for inst in instruments:
      obs, usage = inst.Score(...)              # judge failures → status
      record ObservationRecorded                #   "failed", NEVER a zero
      commit usage                              # judge.calls / judge.tokens
```

### 3.3 The judge instrument (G4)

`pkg/ttc/judgeinstrument` (new), built on judgekit's chain
(spec → protocol → instance+evidence → assessment):

```text
JudgeInstrument.Score(spec, case, output):
  evidence  = top-K retrieved chunks (text from the content store)
  instance  = judgekit eval.NewInstance(
                input:     case.Query,
                candidate: rendered evidence set,
                reference: case's rubric / expected-answer notes,
                facts:     required claims per case)
  assessment = run via geppetto engine (model pinned in the epoch)
  score     = assessment → [0,1] per the spec's dimension weights
  artifact  = FULL assessment (claims, rationale, per-dimension scores)
              stored content-addressed; the observation's Evidence ref
  usage     = {judge.calls: 1, judge.tokens: measured}
```

Non-negotiables, inherited from the honesty rules:

- The judge's **rationale is an artifact**, not a log line — the judge
  tile renders it, the journal can reproduce it.
- A judge API failure is observation status `failed` — coverage still
  records; the campaign completes with a partial second construct rather
  than fabricating.
- The model + prompt digest live IN the epoch definition: change either
  and comparisons across the change are structurally impossible.

### 3.4 Calibration from the start (G11, minimal v1)

judgekit's `calibration` package (gold sets, confusion, recall) is wired,
not aspirational: the case suite marks ~10 cases as **gold** with human
labels for the judge's construct. After each campaign:

- agreement stats (judge vs gold) land as a campaign artifact;
- the judge tile shows per-case judge score beside coverage score, and
  flags disagreement (judge high / coverage low and vice versa) — those
  rows are exactly where the human should look first;
- a drifting judge (agreement drop across epochs) is a red banner, not a
  silent trend.

### 3.5 What changes downstream (G9)

- Estimates become per-construct: cockpit shows two means per arm;
  comparisons return two metric rows (the API shape already carries
  `construct` — the fixture just only ever filled one).
- The failures projection gains `?construct=` (default: coverage).
- The propose loop is unchanged — hypothesis text now typically names
  which construct it expects to move (`expected_improvement.metric`
  already exists for exactly this).

## 4. Budgets for real money (G5)

```yaml
# manifest addition
budget:
  retrieval.queries:   {limit: 600}
  retrieval.results:   {limit: 60000}
  embeddings.tokens:   {limit: 2_000_000}
  judge.calls:         {limit: 1200}     # episodes × instruments × repeats
  judge.tokens:        {limit: 6_000_000}
```

Reservation stays per-episode (claim the worst case, commit the actual);
`budgetQuantities` already folds arbitrary resources. The budget tile (§5)
makes burn visible while the campaign runs; `trial.run` and `proposal.seal`
stay danger verbs whose approval prose states the budget consequence.

## 5. The three operations tiles (G7+G8)

### runner — live episode queue (singleton, campaign-scoped via focus)

```text
┌ ⠿ RUNNER · ptr-… ──────────────────────────────────────┐
│ QUEUE  6 queued · 2 active · 41 done · 1 failed        │
│ THROUGHPUT  ▂▄▅▇▆▅  ~4.2 episodes/min                  │
│ ACTIVE                                                 │
│ ▌ep 9bf95… q-comparison × rrf-30   lease 00:41  try 1  │
│ ▌ep 0587a… q-hybrid     × rrf-30   lease 00:12  try 1  │
│ FAILED (terminal)                                      │
│ ▌ep d6e15… q-policy     × limit-1  embed timeout  ×3   │ ← chip menu:
│ QUEUED · next 4 of 6                                   │   autopsy,
│ ▌q-inventory × rrf-30    ▌q-care × rrf-30  …           │   provenance
└────────────────────────────────────────────────────────┘
widgets: counters row · throughput sparkline · three chip lists
data:    GET /campaigns/{c}/work  (work_items + lease ages + attempts),
         poll or SSE; every episode row is a live presentation
```

### budget — per-resource burn (doc-bound: focus)

```text
┌ ⠿ BUDGET · ptr-… ──────────────────────────────────────┐
│ retrieval.queries   ████████░░░░░░  312/600   reserved 8│
│ retrieval.results   ███░░░░░░░░░░░  9.1k/60k            │
│ embeddings.tokens   █████████░░░░░  1.31M/2M  ⚠ 65%     │
│ judge.calls         ████░░░░░░░░░░  388/1200            │
│ judge.tokens        ███░░░░░░░░░░░  1.9M/6M             │
│ BURN  embeddings ~41k tokens/min → exhausts in ~17 min  │
│ (bar: used █ / reserved ▒ / remaining ░; ⚠ over 60%)   │
└────────────────────────────────────────────────────────┘
widgets: one bar row per resource · burn-rate + projection line
data:    budget snapshot (already in cockpit) + reservation detail
```

### substrate — what we are searching over (doc-bound: focus)

```text
┌ ⠿ SUBSTRATE · ptr-… ───────────────────────────────────┐
│ INDEX BUNDLE  ▌bundle sha256:9ac2f1…   built 2026-08-25│
│  chunks         14,382    documents  1,257             │
│  representations 3  (raw, summary, qa)                 │
│  embedding model ▌text-embedding-…  dim 1536           │
│  lexical         bm25 · tantivy-…                      │
│ PINNED BY  ▌campaign:…  at CampaignCreated             │
│ DRIFT  live bundle == campaign bundle ✓                │ ← red when the
└────────────────────────────────────────────────────────┘   serve process
widgets: identity tray · counts · drift check                bundle differs
data:    indexbundle.Manifest via a small provenance endpoint
```

## 6. The operator workflow, phase by phase

“You” = the human. “Agent” = the coding agent (this assistant or the
workbench agent seat). Every phase ends at something you can SEE.

### Phase 0 — substrate (once per corpus change)

1. Agent: corpus ingest + representation build + embedding run from the
   knowledge DB (`rag-ttc inspect corpus …`, representations commands);
   produces the index bundle + manifest.
2. You: read the substrate summary (chunk/doc counts, model, digest).
   *Steering point: is this the corpus you mean to optimize over?*
3. The bundle digest goes into the manifest; the substrate tile will
   verify serve-vs-campaign bundle identity forever after.

### Phase 1 — cases and rubrics (the highest-leverage human hours)

1. You: pick 20–50 real queries (customer logs, support tickets, your own
   hard cases). For each: expected target chunks/docs (coverage ground
   truth) and, for the judge, a short rubric note ("a good retrieval
   supports answering X and Y").
2. Agent: drafts the case YAML; marks ~10 as gold with your labels.
3. You: red-pen the drafts. *This defines what "better" means; nothing
   downstream can fix a bad case suite.*

### Phase 2 — baseline campaign, eyes open

```bash
campaign dry-run --manifest ttc-live-baseline.yaml --store S
   # → episodes count, per-resource budget claims. Read this table.
campaign run --manifest … --store S --stop-after after_lease
campaign serve --store S … --agent-token …
pnpm dev   # open the workbench: cockpit + runner + budget + substrate
```

*Steering points: the dry-run budget table (before any spend); the
substrate drift check; budgets in the manifest are YOUR ceilings.*

### Phase 3 — execution, watched

`campaign resume` works the queue. You watch: runner (jobs in flight,
failures with reasons), budget (burn + projection), failures tile filling
worst-first per construct, judge tile showing rationale beside coverage,
disagreement flags. Anything odd → autopsy the episode's real stages,
open the chunks. *Steering: pause is ctrl-C (leases expire, resume
continues); raising a budget is a manifest decision, not a click.*

### Phase 4 — the propose loop (unchanged, now on real data)

Catalog → draft (you or the agent seat) → compile + preview probes on
real cases → human adjusts values in the shared document → seal (approval
flow if agent-proposed; proposer provenance lands in the journal) → trial
of the candidate → compare tile: per-construct deltas + slope. Decide.

### Phase 5 — iterate and audit

`campaign verify` after anything important; the journal is the record of
what was tried, by whom, under which approval, measured by which epoch.
Judge calibration report per campaign; recalibrate the rubric before
trusting a judge-led delta.

## 7. Ticket slicing

- **OPTKIT-025** (this ticket): G1–G5 + G9 — live executor, catalog
  mapping, instrument seam, judge instrument with calibration-minimal,
  budget resources, per-construct estimates. Exit: a 5-case live
  campaign on the real bundle completes with BOTH observations per
  episode and verifies.
- **OPTKIT-026**: G7+G8 — work/budget/substrate projections + the three
  tiles. Exit: the phase-3 watching experience exists.
- **OPTKIT-027**: G6 + first full campaign + G11 full calibration
  surface (+ G10 trial wiring when the backend colleague lands it).
  Exit: the phase-by-phase workflow above executed end to end on the
  real knowledge base, decisions recorded.

## 8. Open questions (flag, do not decide silently)

- Judge model residency: which provider/model for the judge epoch, and is
  its spend metered on the same budget ledger as embeddings (proposed:
  yes, separate resources)?
- Case provenance: do curated cases eventually become documents in the
  workbench (a `cases` tile with the draft pattern) or stay manifest
  YAML? Proposed: YAML through OPTKIT-027, revisit after.
- Sensitivity: chunk text reaching the judge is an
  `artifact.read.restricted`-class flow when the corpus has restricted
  sources — the judge instrument must respect the same policy filters the
  pipeline applies (it scores what was RETRIEVED, which already passed
  policy — but rubrics must not leak restricted references).
