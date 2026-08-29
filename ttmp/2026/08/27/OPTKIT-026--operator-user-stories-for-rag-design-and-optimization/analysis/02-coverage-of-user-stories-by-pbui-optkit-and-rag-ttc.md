---
Title: Coverage of the OPTKIT-026 User Stories by pbui, optkit, and rag-ttc
Ticket: OPTKIT-026
Status: active
Topics:
    - design
    - ui
    - rag-ttc
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/27/OPTKIT-026--operator-user-stories-for-rag-design-and-optimization/analysis/01-user-stories-the-person-who-makes-the-rag-system-good.md
      Note: The needs-first story set this report measures against
    - Path: /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/apps/workbench/web/src/apps/index.ts
      Note: |-
        The 22 implemented workbench tiles — concrete evidence of what is built
        The 22 implemented tiles — evidence of what is built
    - Path: /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/judgeinstrument/instrument.go
      Note: The LLM judge exists standalone but is not wired into campaigns
    - Path: /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/optkitcampaign/campaign.go
      Note: |-
        scoreObservation — still one hardcoded coverage instrument; the judge seam (OPTKIT-025 G3) is not wired
        scoreObservation still one hardcoded instrument; judge seam not wired
    - Path: /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/workbenchhost/documents.go
      Note: |-
        Document format validators — corpus-notes/questionset-draft are absent, proving Theme 2 is unbuilt
        Document format validators — corpus-notes/questionset-draft absent, proving Theme 2 unbuilt
ExternalSources: []
Summary: A per-story audit of the original S1-S48 operator user stories against the code that actually exists in pbui, optkit, and rag-ttc; S49-S50 were added after this audit and are intentionally not included in its totals.
LastUpdated: 2026-08-28T21:45:00-04:00
WhatFor: The evidence-backed gap map that drives the next ticket slicing; the acceptance companion to the user-stories document.
WhenToUse: Before opening a new implementation ticket; before claiming a story is "done"; to check whether a need already has a home.
---




# Coverage of the OPTKIT-026 User Stories by pbui, optkit, and rag-ttc

This report measures the original forty-eight stories, S1–S48, in
[`01-user-stories-the-person-who-makes-the-rag-system-good.md`](./01-user-stories-the-person-who-makes-the-rag-system-good.md)
against the code that existed when the audit was performed. **Scope note:**
S49 (end-to-end grounded-answer benchmark) and S50 (one durable
build-to-benchmark workflow) were added on 2026-08-28 after testing the story
set against a concrete operator goal. They are intentionally not included in
the 9 covered / 11 partial / 23 designed / 5 not-covered totals below; a
future coverage refresh should audit all fifty stories. This report
distinguishes four states:

- **Covered (built)** — implemented and usable from the workbench or CLI now.
- **Partial (built-ish)** — a real mechanism exists, but the user-facing
  story is incomplete (often the recording/projection side is built and the
  surface is not, or a CLI path exists and the UI does not).
- **Designed (not built)** — a design document exists in an open ticket but
  no code implements it.
- **Not covered** — neither design nor implementation exists.

Every claim is anchored to a file, a tile, an endpoint, or a ticket. The
method was: read the three existing intern guides (OPTKIT-025/027/029) and
the OPTKIT-026 design docs, then verify the implementation state in the
code — the implemented tiles in `apps/workbench/web/src/apps/`, the document
format validators in `workbenchhost/documents.go`, the executor/instrument
seam in `optkitcampaign/campaign.go`, and the API packages `exploreapi`,
`experimentsapi`, `workbenchapi`, `specialistapi`.

## 1. The repositories and what each one is

The three repositories are not peers; they are layers, and which layer a
story lives in determines what "covered" means for it.

- **`optkit`** is the domain-neutral control plane: canonical records,
  content-addressed artifacts, configuration space (lenses, variables,
  snapshots, patches, candidates), episode trajectories, measurement
  epochs and observations, complete-block trials, the append-only hash-chained
  campaign journal, the durable SQLite work queue with leases, finite
  budgets, and a read-only query plane (`optkit/query`, `optkit/internal/web`).
  It is **fully built** (see its README and `docs/01-…md`). Optkit is the
  *mechanism* behind the honesty stories (S20, S25, S31) and the budget
  story (S15, S32), even where the workbench surface for those stories is
  not. Optkit never imports ragkit, ragopt, judgekit, or rag-ttc
  (`internal/boundary/boundary_test.go` enforces this).

- **`rag-ttc`** is the TTC product: the retrieval pipeline (`pkg/ttc/search`),
  the optimization model (`pkg/ttc/optimization`), the campaign machinery
  (`pkg/ttc/optkitcampaign`), the LLM judge (`pkg/ttc/judgeinstrument`), the
  three HTTP APIs (`specialistapi`, `workbenchapi`, `workbenchhost`) plus the
  two read surfaces added by OPTKIT-027/029 (`exploreapi`, `experimentsapi`),
  and the workbench frontend (`apps/workbench/web`). Most user stories live
  or will live here.

- **`pbui`** is the presentation/UI framework the workbench is built *on*:
  the `pbui-workbench` package (tile, split-pane, surface, launcher, store,
  `defineApp`), the `pbuichat` package (chat server, agent vocabulary,
  mentions, tools, trace), and the Go `workbench`/`chatui`/`authkit`
  packages. Pbui is the *framework* layer for the agent seat (Theme 8) and
  for every tile every other theme renders through. A story is "covered by
  pbui" when the framework capability the story needs already exists, even
  if the product has not yet assembled a tile from it.

## 2. Summary by theme

| Theme | Stories | Covered | Partial | Designed | Not covered |
|---|---|---|---|---|---|
| 1 — Knowing my material | S1–S3 | 1 | 0 | 2 | 0 |
| 2 — Defining what good means | S4–S8 | 0 | 0 | 5 | 0 |
| 3 — Building intuition | S9–S13 | 4 | 1 | 0 | 0 |
| 4 — Running experiments | S14–S20 | 0 | 5 | 0 | 2 |
| 5 — Judging quality | S21–S25 | 1 | 0 | 4 | 0 |
| 6 — Deciding | S26–S29 | 2 | 0 | 0 | 2 |
| 7 — Living with the system | S30–S33 | 0 | 3 | 0 | 1 |
| 8 — Working with an assistant | S34–S36 | 1 | 2 | 0 | 0 |
| 9 — Running the work | S37–S48 | 0 | 0 | 12 | 0 |
| **Total** | **48** | **9** | **11** | **23** | **5** |

Read across: nine stories are usable today, eleven have a real mechanism but
no complete surface, twenty-three are designed in an open ticket but not
built, and five have neither design nor implementation. The five unowned
gaps are listed in §6 and are the input to the next slicing.

## 3. Per-story audit

Evidence abbreviations: **027** = OPTKIT-027 (built), **029** = OPTKIT-029
(built), **025** = OPTKIT-025 (active, designed/started), **026** =
OPTKIT-026 (active, design only), **028** = OPTKIT-028 (active, designed),
**base** = the foundation built by OPTKIT-001…024.

### Theme 1 — Knowing my material

| Story | Status | Evidence | Gap / next ticket |
|---|---|---|---|
| S1 — see what is in the corpus | **Covered** | 027: `CorpusApp.tsx`, `DocumentApp.tsx`; `GET /api/rag/v1/corpus/survey`, `/corpus/documents`, `/corpus/documents/{id}` over `indexbundle.Inspect`+`Measure`. Counts by kind, length histogram, percentiles, signals, uncounted. | Step 4 (mark exclusions persistently) is S3, not built. |
| S2 — can the corpus answer the questions | **Designed** | 026 design-doc/01 §4/§8: `coverage.sweep` verb + projection. | No implementation ticket. Gated on a live arm (025 G1). |
| S3 — record junk/duplicate/missing | **Designed** | 026 design-doc/01 §3/§7: `ragttc.corpus-notes/v1`, `corpus.flag`/`unflag` verbs, flag counters. | Format absent from `workbenchhost/documents.go` (only focus/comparison/proposal-draft/watchlist/conversation/explore-focus/results-focus/finding/widget exist). No ticket. |

### Theme 2 — Defining what good means

| Story | Status | Evidence | Gap / next ticket |
|---|---|---|---|
| S4 — assemble a question set | **Designed** | 026 design-doc/01: `ragttc.questionset-draft/v1`, `questions`/`question` tiles, `question.*` verbs, seal. | No `QuestionsetDraft` format, no `QuestionsApp`/`QuestionApp` tile. **No implementation ticket exists.** |
| S5 — say what a right answer looks like | **Designed** | 026 design-doc/01 §5: `question` tile (rubric, targets, non-answers). | Not built. Same gap as S4. |
| S6 — author ground truth from evidence | **Designed** | 026 design-doc/01 §5 (`ask` tile) + design-doc/02 §4 (`hit` carries `target.attach`). The `ask` tile exists (027) but `target.attach` does not. | Not built. Needs S4 + the `target.attach` verb. |
| S7 — change ground truth without corrupting history | **Designed** | 026 design-doc/01 §4/§7: sealed versions, `questionset.seal`, versioned additions. | Not built. No sealed evaluation-set artifact today (only the static `datasets/ttc/evaluation.json`). |
| S8 — keep a small set I trust absolutely | **Designed** | 026 design-doc/01: `question.markGold`, gold labels, gold callouts. | Not built. Gold flags exist as a design only. |

> Theme 2 is the single largest unowned gap. The needs document itself
> flags ground-truth authoring as "the hinge" and "the most human-hours-
> expensive stories"; today none of S4–S8 has an implementation ticket, and
> the only evaluation set in the repository is a committed JSON file. The
> 026 design-doc/01 is a complete implementation contract waiting for a
> ticket (see §6, recommended OPTKIT-030).

### Theme 3 — Building intuition

| Story | Status | Evidence | Gap / next ticket |
|---|---|---|---|
| S9 — ask a question right now | **Covered** | 027: `AskApp.tsx`; `POST /api/rag/v1/search` (metered, not an episode). | — |
| S10 — understand why a result ranked where it did | **Covered** | 027: `TrailApp.tsx` reads `SearchOutput.Stages[].Candidates` (`StageCandidate{ChunkID,Rank,Channel,RepresentationID,Score,Contributions}`) directly from the search response. | Per-candidate drop reason not recorded (026-02 §11 open question). |
| S11 — see how a document splits | **Covered** | 027: `SplitApp.tsx`; `POST /corpus/documents/{id}/chunk` (pure CPU, `chunking.Apply`). | — |
| S12 — see what the retriever matched on | **Partial** | 027: `RepresentationID` surfaces in the trail ("on ▌raw/summary"); `ChunkApp.tsx` shows representations. The `representation.flag` (misleading) verb is a 026-01/02 design, not built, so "counts against the approach" does not. | `representation.flag` (026 design). |
| S13 — two configurations side by side | **Covered** | 027: `AskApp` lanes (two `--index-bundle`s), client-side diff (moved/new/dropped). | The recorded-arm `compare` tile (base) is a different, measured comparison; ad-hoc vs measured are deliberately separate (026-02 §9). |

### Theme 4 — Running experiments

| Story | Status | Evidence | Gap / next ticket |
|---|---|---|---|
| S14 — try variants, learn which is best | **Partial** | base: propose loop built (`CatalogApp`,`ProposalApp`,`IntentApp`,`InvalidationApp`,`PreviewApp`; `proposals:compile`/`seal`, `previews:run`, `trial.run` verbs). 029: reading results built. But a real multi-variant campaign over the real bundle needs the live executor (025 G1); today only `SemanticFixtureExecutor` runs. | 025 G1 (live executor). |
| S15 — know cost before it runs | **Partial** | base: `campaign dry-run` prints per-resource budget claims (CLI); optkit `budget` reserves/commits. The workbench cost-estimate UI and runtime ceiling-pause are 028 designs. | 028 (runtime ceiling); 026-03 (planning estimate). |
| S16 — watch a run and intervene | **Designed** | 028: `jobs`/`job` tiles, `GET /work`, `/work/stream` (SSE). | Not built — no `JobsApp`/`JobApp` tiles, no `/work` routes. |
| S17 — resume or repair a broken run | **Partial** | base: `campaign resume` (engine A, leases expire, proven CLI behavior); work queue supports retries/expiry. The UI retry-failed-group and partial-marking are 028 designs. | 028. |
| S18 — change the substrate as a dimension | **Partial** | base: optkit `space` config graph + `experimentworkbench` invalidation planner (reuse vs recompute per layer) built. Real substrate rebuild as a campaign dimension needs the live executor + bundle provenance. | 025 G1/G7; 026-02 §11 (bundle records chunker settings). |
| S19 — sweep a range without hand-building | **Not covered** | No sweep generator; catalog variables are individual. Not designed in any 026 doc. | New design + ticket. |
| S20 — only compare comparable things | **Partial** | 029: legs (`<bundle>:<strategy>`) enforce `same_leg` refusal and make comparison honest at the (bundle,strategy) level; `same_leg` is refused because "a plausible-looking nothing is worse than an error". Full provenance comparability (corpus/judge/question-set identity in the manifest) is a 026-01 design (`evaluation_set` in manifest). | 026-01 (manifest `evaluation_set`); 025 (per-construct epochs). |

### Theme 5 — Judging quality

| Story | Status | Evidence | Gap / next ticket |
|---|---|---|---|
| S21 — measure whether evidence supports an answer | **Designed** | 025 G4 + 026-03 judge panel: `ttc.judge.retrieval-support/v1` second construct. `pkg/ttc/judgeinstrument/` exists (instrument.go, prompts.go, record.go) and is used by the standalone `answerquality` command, but is **not wired into campaigns** (`optkitcampaign/campaign.go:766 scoreObservation` still hardcodes one coverage instrument). | 025 G3 (instrument seam) + G4 (wire judge). |
| S22 — read the judge's reasoning | **Designed** | 025 G4 (rationale stored as a content-addressed artifact); 026-03 judge panel (streaming verdicts, score distribution forming). | Not built — no judge panel tile. |
| S23 — overrule the judge, have it count | **Designed** | 025 G11 (calibration: gold labels, agreement stats); 026-01 gold labels. | Not built. |
| S24 — know whether the judge can be trusted | **Designed** | 025 G11 (drift detection, re-judge under new version). | Not built. |
| S25 — tell "failed" from "scored badly" | **Covered (mechanism)** | optkit `measure.Observation` status (`measured`/`failed`/`unknown`/`inapplicable`); "missing is never zero" enforced in projections (`FailuresApp` shows `MISSING`, never `0`). The judge-failure-as-`failed` specifically is a 025 rule, but the discipline is built. | 025 (judge failure handling). |

### Theme 6 — Deciding

| Story | Status | Evidence | Gap / next ticket |
|---|---|---|---|
| S26 — see what got worse, not just better | **Covered** | 029: `ComparisonApp.tsx` shows "worst first" moved questions, per-query deltas, perfect/zero counts beside the mean. | — |
| S27 — know whether a difference is real | **Not covered** | 029 §9 explicitly defers: "epsilon is float noise, not a threshold … deciding whether it matters is the reader's job." No significance design. | New design + ticket (deliberate deferral). |
| S28 — record the decision and its justification | **Covered** | 029: `FindingApp.tsx` — a finding (claim, status `proposed`→`accepted`/`disputed`/`withdrawn`, evidence keys). A finding stores no numbers; it re-reads its comparison, which is what makes it falsifiable. | — |
| S29 — adopt a configuration, know what it changes | **Not covered** | `trial.run` materializes a candidate as a measured arm, but there is no "adopt into live / diff against running / rebuild+deploy consequences / confirmed revert" workflow. Not designed. | New design + ticket. |

### Theme 7 — Living with the system

| Story | Status | Evidence | Gap / next ticket |
|---|---|---|---|
| S30 — re-check after the corpus changes | **Not covered** | 026-01 §9 explicitly: "No coverage automation … Automatic re-checking is a Theme 7 concern." No design, no ticket. | New design + ticket. |
| S31 — reproduce an old result | **Partial** | base: content addressing + `campaign verify` (journal hash chain + payload re-hash) + manifest digests built. The one-click "reproduce" surface is not. | Surface ticket. |
| S32 — know what it costs to run | **Partial** | base: optkit `budget` ledger + cockpit budget snapshot built. Per-experiment cost attribution and per-query production cost are 025 G5 / 028 spend views, not built. | 025 G5; 028. |
| S33 — hand the story to someone else | **Partial** | 029 findings + base trace (every verb attributed) + journal provide the raw material. A unified chronological "story / export" view is not built. | Surface ticket. |

### Theme 8 — Working with an assistant

| Story | Status | Evidence | Gap / next ticket |
|---|---|---|---|
| S34 — delegate the tedious parts | **Partial** | base (OPTKIT-024): agent seat built — agent token holds `catalog.read`/`proposal.compile`/`preview.run`, is 403 on `seal`; shared verb vocabulary; chat router. Concrete delegation workflows (draft 80 questions, propose a sweep, triage 40 failures) are not built. | 025/026 authoring tools; S4–S8. |
| S35 — approve anything that costs or becomes permanent | **Covered** | base (OPTKIT-023/024): danger verbs (`proposal.seal`, `trial.run`) park as approvals; `SealBar`; per-principal grants; approval carries `approvalId` into the verb and the trace. | — |
| S36 — see and undo what the assistant did | **Partial** | base: `TraceApp.tsx` records every verb with `actor` and `outcome` ("a verb that touched nothing never reads as performed"). "Undo reversible ones" and "label permanent vs reversible" are not built. | Undo design + ticket. |

### Theme 9 — Running the work

| Story | Status | Evidence | Gap / next ticket |
|---|---|---|---|
| S37 — see everything running in one place | **Designed** | 028: `jobs` tile, `GET /work`. | Not built. |
| S38 — watch a build by phase | **Designed** | 028: `job` tile, `progress.json` snapshot, phases from `flow.StepReport`. | Not built. |
| S39 — start/pause/stop/resume | **Designed** | 028 §8; engine A resume is proven CLI, the UI control bar is design. | Not built. |
| S40 — failures grouped by cause, live | **Designed** | 028: `failureGroup` from `flowkit ErrorClass` + `RetriesByClass`. | Not built. |
| S41 — inspect one unit in flight | **Designed** | 026-03 `unit` tile; 028. | Not built. |
| S42 — watch judgments arrive, stop early | **Designed** | 026-03 judge panel (streaming verdicts, score histogram, keep/discard on restart). | Not built; gated on 025 G4. |
| S43 — evaluate a whole conversation | **Designed** | 026-03 §8.4: `rag.ttc-conversation/v1` executor; scenario = question with follow-ups; turns as spans. | Not built; gated on 025 G1. |
| S44 — find the turn it went wrong | **Designed** | 026-03: `conversation` tile, first-degradation mark, reuses Theme 3 `trail`. | Not built. |
| S45 — re-run one conversation after a change | **Designed** | 026-03: `conversation.rerun` danger verb, side-by-side turns. | Not built. |
| S46 — leave it overnight, know what happened | **Designed** | 028: finished list with outcomes/costs; overnight summary as projection. | Not built. |
| S47 — steer throughput at runtime | **Designed** | 028 §8: `job.setConcurrency` (workers/batch). | Not built. |
| S48 — dependencies; refused when not ready | **Designed** | 026-03 + 028: `waitsOn`, `job.start` refused when deps unmet. | Not built. |

## 4. What the open tickets address

Five tickets are open and together own thirty-eight of the forty-eight
stories (built or designed). The map below is the actionable view: which
ticket, if completed, turns which "designed" into "covered" and which
"partial" into "covered".

### OPTKIT-025 — Live TTC campaign: real executor, LLM judge, operations tiles (active)

The foundational backend slice. Its gap list (analysis §1) is G1–G5, G9, G11.

- **G1 live executor** (`rag.ttc-live/v1`) — turns S14 from partial to
  covered and unblocks S18 (substrate dimension) and every Theme 9 piece that
  needs real retrieval (S42–S45).
- **G3 instrument seam** + **G4 LLM judge** wired into campaigns — turns
  S21, S22 from designed to covered and supplies the judge-failure path for
  S25; `pkg/ttc/judgeinstrument/` already exists standalone.
- **G5 budget resources** (`embeddings.tokens`, `judge.calls`, `judge.tokens`)
  — advances S15, S32.
- **G9 per-construct estimates** — makes the comparison API carry two
  constructs, which S26/S20 rely on.
- **G11 calibration-minimal** — turns S23, S24 from designed to covered.

Status in code: `judgeinstrument/` exists and is exercised by the
`answerquality` command, but `optkitcampaign/campaign.go` still hardcodes the
single coverage instrument and the only `Executor` implementation is
`SemanticFixtureExecutor`. So 025 is started, not landed.

### OPTKIT-026 — Operator user stories (active, design only)

The needs-first document plus three design docs that are implementation
contracts:

- **design-doc/01** — Themes 1–2 (S1–S8): corpus/question tiles, verb tables,
  `questionset-draft`/`corpus-notes` formats, `questionset.seal`,
  `coverage.sweep`. This is the contract for the **unowned Theme 2 work**.
- **design-doc/02** — Theme 3 (S9–S13): `run`/`hit`/`chunkPreview` types,
  `trail`/`split` tiles, `StageSummary` un-coarsening. Largely *already built*
  by 027; the remaining designed piece is the `representation.flag` verb (S12).
- **design-doc/03** — Theme 9 (S37–S48): `job`/`phase`/`unit`/`failureGroup`/
  `turn` types, `jobs`/`job`/`unit`/`conversation` tiles, the work projection
  over both engines, the conversation executor. This is the design 028
  implements.

OPTKIT-026 itself ships no code; it is the acceptance vocabulary and the
design source for 028 and for the missing Theme 2 ticket.

### OPTKIT-027 — Explore workbench (complete, built)

Built and shipped: ad-hoc search, corpus browsing, chunking preview. Owns
S1, S9, S10, S11, S13 (covered) and the core of S12 (partial). It is the
load-bearing exploration loop the needs document says "a large fraction of
the other stories are unreachable without."

### OPTKIT-028 — Watching the work (active, designed)

The implementation ticket for Theme 9. Owns S16, S37–S48 (all designed) and
the runtime-ceiling half of S15. Its sequencing (028 §10) is P1 progress
snapshot → P2 work projection → P3 tiles → P4 engine A adapter → P5 control.
Nothing is built yet (no `JobsApp`/`JobApp`/`UnitApp`/`ConversationApp`,
no `/work` routes, no `progress.json`).

### OPTKIT-029 — Reading results (complete, built)

Built and shipped: the experiments projection, the `results`/`comparison`/
`finding` tiles, the `index results`/`index compare` CLI. Owns S26, S28
(covered), advances S20 (partial), and is the reading half of S14. Its
deliberate non-goals — no significance (S27), no writing/running — are the
source of two of the unowned gaps.

### The foundation: OPTKIT-001…024 (done)

These built the substrate the partial stories stand on: optkit itself
(001–005, 010–014), the propose loop and candidate sealing (016–018), the
React workbench and RRF vertical slice (019, 022, 023), the pbui adoption and
agent seat (021, 024), and the evidence tiles (022). They are why S14, S31,
S32, S35, S36 are "partial" rather than "not covered."

## 5. Coverage by repository (what each layer provides today)

- **pbui** provides the framework for: the tile system and presentation kernel
  (all themes render through it), the chat server + agent vocabulary + trace
  (Theme 8: S35 covered, S34/S36 partial), and the workbench document model.
  It does not itself implement any operator story end-to-end; it is the layer
  the product assembles stories from.

- **optkit** provides the mechanism for: content addressing + verification
  (S31 partial), the budget ledger (S15/S32 partial), the measurement epoch +
  observation status discipline (S25 covered, the honesty backbone for S20),
  the config-graph invalidation planner (S18 partial), and the durable work
  queue with leases (S17 partial). It is fully built and is the reason the
  "honesty stories are mostly about surfacing guarantees" (needs doc §"What
  this set implies") are already half-true.

- **rag-ttc** provides the product surfaces: the explore workbench (027:
  Theme 1/3 covered), the results surface (029: Theme 6 covered, Theme 4
  reading), the propose loop and agent seat (base: Theme 8 partial), the
  retrieval pipeline and judge (Theme 5 designed), and — designed only — the
  job view and conversation evaluation (Theme 9). It is also where the
  unowned Theme 2 work must land.

## 6. The unowned gaps (no open ticket)

These five stories have neither an implementation nor a design home. They
are the next slicing's input.

1. **Theme 2 — S4–S8 (ground-truth / question-set authoring).** The highest-
   leverage unowned work. 026 design-doc/01 is a complete contract; it needs
   an implementation ticket. **Recommended: OPTKIT-030 "Authoring the
   question set and corpus notes."** Unblocks S6, S7, S8 and feeds S21/S23
   (the judge needs gold labels).
2. **S19 — sweep a range.** No design. Needs a sweep generator over catalog
   variables and a curve result view.
3. **S27 — is a difference real.** Deliberately deferred by 029. Needs a
   significance/spread design that does not produce "false precision."
4. **S29 — adopt a configuration.** No design. Needs a deploy/diff/revert
   workflow and a "what this invalidates" computation (the invalidation
   planner already exists; the adopt surface does not).
5. **S30 — re-check after corpus changes.** Acknowledged as a Theme 7 concern
   in 026-01 §9 but not designed. Needs a scheduled re-evaluation tied to
   corpus-digest change with the S20 comparability guarantees in view.

Three further surfaces are "partial" with no specific ticket: the one-click
**reproduce** surface (S31), the per-experiment **cost** surface (S32), and
the chronological **story/export** view (S33). Their mechanisms exist; the
surfaces do not, and each is small enough to fold into 028 or a Theme 7
ticket rather than standing alone.

## 7. What this implies for the next slicing

- **Unblock Theme 2 first.** It is designed, unowned, and everything in
  Theme 5 that uses gold labels (S23, S24) and everything in Theme 4 that
  needs a real evaluation set depends on it. 026-01 is the contract; a
  single ticket ships S4–S8.
- **Land 025 before 028's judge panel.** The 028/026-03 judge panel (S42)
  and conversation evaluation (S43–S45) are gated on the live executor and
  the judge-in-campaign (025 G1/G3/G4). `judgeinstrument/` exists standalone;
  the wiring is the critical path.
- **028 is the right shape but has not started.** Theme 9 is twelve designed
  stories with zero code. Its P1–P3 (progress snapshot + work projection +
  build tiles) is self-contained and works on today's data; it should not
  wait for 025.
- **Keep the deliberate non-goals.** 029's refusal of significance (S27) and
  of writing/running is correct discipline; those become their own small
  tickets with their own designs, not creep inside 029.
- **The five unowned gaps are the next tickets' names.** OPTKIT-030
  (question set), and small tickets for S19, S27, S29, S30 — each with the
  needs document as the acceptance vocabulary, exactly as 026 intended.
