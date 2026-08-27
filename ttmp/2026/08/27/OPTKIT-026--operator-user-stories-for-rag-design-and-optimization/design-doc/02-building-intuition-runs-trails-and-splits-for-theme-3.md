---
Title: 'Building Intuition: Runs, Trails, and Splits for Theme 3'
Ticket: OPTKIT-026
Status: active
Topics:
    - design
    - ui
    - rag-ttc
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://rag-ttc/pkg/ttc/search/service.go
      Note: RetrievalStage and StageCandidate — the recorded per-candidate data this design surfaces
    - Path: repo://rag-ttc/pkg/ttc/specialistapi/types.go
      Note: StageSummary, the lossy projection this design un-coarsens
    - Path: repo://rag-ttc/apps/workbench/web/src/apps/AutopsyApp.tsx
      Note: The stage tile extended with candidate rows
    - Path: repo://rag-ttc/apps/workbench/web/src/apps/ChunkApp.tsx
      Note: The chunk tile extended with matched-representation marking
    - Path: repo://rag-ttc/cmd/rag-ttc/cmds/indexes/build.go
      Note: selectChunker — the chunkers the split preview offers
ExternalSources: []
Summary: A design for stories S9-S13 built on one unification — a retrieval run has two provenances, scratch and recorded — yielding two new tiles, three extensions to existing ones, and a projection fix that surfaces per-candidate data the pipeline already records.
WhatFor: The implementation contract for the exploration surfaces of the workbench.
WhenToUse: Before implementing ad-hoc query, rank explanation, or chunking comparison; as the acceptance vocabulary for stories S9-S13.
---

# Building Intuition: Runs, Trails, and Splits for Theme 3

## 1. Scope and the discipline applied

This document covers **S9** (ask a question right now), **S10** (understand
why a result ranked where it did), **S11** (see how a document splits under
different chunking strategies), **S12** (see what the retriever actually
matched on), and **S13** (two configurations side by side on one question).

The same two constraints as the previous document apply: smallest surface
that makes each story true, and reuse over invention. Two findings from
reading the code shaped everything below.

**The first finding is that the data already exists.** `ttcsearch.
RetrievalStage` records, per stage, a `Candidates []StageCandidate` where
each candidate carries `{ChunkID, Rank, Channel, RepresentationID, Score,
RerankerScore, Contributions}`, and `rag.Contribution` is `{Channel, Rank,
Weight, Value}`. That is a complete answer to "why did this rank here",
recorded on every episode since OPTKIT-022. It never reaches the browser
because `specialistapi.StageSummary` drops the `Candidates` field. **S10 is
mostly a projection fix, not a feature.**

**The second finding is that a retrieval run has two provenances and one
shape.** An episode recorded in a campaign and an ad-hoc query typed into a
box both produce a `ttcsearch.SearchOutput`. If the vocabulary models a
`run` with a provenance rather than modelling "episode" and "ad-hoc query"
separately, then one set of tiles serves both, and every explanation surface
built for exploration works unchanged on recorded results. That unification
is why this theme adds only two tiles.

## 2. What the design adds, in total

| Kind | Count | Names |
|---|---|---|
| Presentation types | 3 | `run`, `hit`, `chunkPreview` |
| New tiles | 2 | `trail`, `split` |
| Extended tiles | 3 | `ask` (lanes), `autopsy` (candidate rows), `chunk` (matched representation) |
| Workspaces | 3 | Probe, Bake-off, Chunking lab |
| Document formats | 0 | scratch runs are client state; pointers reuse `focus` |
| Read projections | 2 | ad-hoc search, chunk-split preview |
| Projection changes | 1 | `StageSummary` carries candidates |

## 3. Semantic objects

```ts
/** One retrieval, whatever produced it. */
export interface RunRef {
  runId: string;
  provenance: "recorded" | "scratch";
  query: string;
  armId: string;
  /** Present only when provenance is "recorded". */
  campaignId?: string;
  episodeId?: string;
  resultCount?: number;
}

/** One chunk as one run saw it — the subject of "why did this rank here". */
export interface HitRef {
  runId: string;
  chunkId: string;
  rank: number;
  score?: number;
  channel?: string;            // which retriever surfaced it
  representationId?: string;   // what it matched on
  title?: string;
}

/** A hypothetical chunk from a chunking preview. NOT a chunk: it has no
 *  identity in any index, so it deliberately gets a near-empty menu. */
export interface ChunkPreviewRef {
  documentId: string;
  settingId: string;           // "A" | "B" | …
  index: number;
  byteStart: number;
  byteEnd: number;
  runes: number;
  headingPath?: string;
  orphaned?: boolean;          // lost its heading context
}
```

Type-graph placement:

```ts
{ id: "run",          parents: ["inspectable"] },
{ id: "hit",          parents: ["inspectable"] },
{ id: "chunkPreview", parents: ["inspectable"] },
```

`hit` is deliberately not `watchable`: a hit is a fact about one run, not an
object you track over time. The chunk behind it is watchable, and the hit's
menu offers a route to it.

One fact joins `WorkbenchFacts`:

```ts
  /** The run whose lens the chunk tile renders under ("matched on…"), and
   *  the default subject for trail verbs. Null outside an exploration
   *  workspace, which is why the chunk tile's matched-representation row
   *  is conditional rather than always-on. */
  activeRunId: string | null;
```

## 4. Verbs, by object

```ts
export type RunVerb =
  | { kind: "ask.run";        query: string; lanes: string[]; limit: number }  // COMMAND
  | { kind: "ask.addLane";    armId: string }                                  // LOCAL
  | { kind: "ask.removeLane"; armId: string }                                  // LOCAL
  | { kind: "ask.clear" }                                                      // LOCAL
  | { kind: "open.trail";     runId: string; chunkId: string }
  | { kind: "open.autopsy.run"; runId: string }        // autopsy on a scratch run
  | { kind: "run.pin";        runId: string }          // keep a bounded summary
  | { kind: "run.activate";   runId: string };         // LOCAL — sets activeRunId

export type SplitVerb =
  | { kind: "open.split";     corpusId: string; documentId: string }
  | { kind: "split.addSetting";    documentId: string;
                                   chunker: string; maxRunes: number;
                                   overlapRunes: number }
  | { kind: "split.removeSetting"; documentId: string; settingId: string };

export type RepresentationVerb =
  | { kind: "representation.flag";   chunkId: string; representationId: string;
                                     reason: string }
  | { kind: "representation.unflag"; chunkId: string; representationId: string };
```

Menu contributions:

| Object | Label | Verb | Notes |
|---|---|---|---|
| `run` | Open stage view (primary) | `open.autopsy.run` | the same autopsy tile, either provenance |
| | Use as lens | `run.activate` | makes the chunk tile show "matched on…" |
| | Pin this run | `run.pin` | scratch only; writes a bounded summary |
| `hit` | Why did this rank here? (primary) | `open.trail` | |
| | Open the chunk | `open.chunk` | existing verb, reused |
| | Mark as … target | `target.attach` | **Theme 2's contribution, unchanged** |
| | Add the chunk to the watchlist | inherited via chunk | |
| `chunkPreview` | Inspect | inherited | *the whole menu* — see §9 |
| `representation` | Flag as misleading… | `representation.flag` | writes to `ragttc.corpus-notes/v1` |
| | Remove flag | `representation.unflag` | |
| `document` | Compare chunking… | `open.split` | Theme 2's type, new contribution |
| `case` *(existing)* | Ask this against… | `ask.run` | a recorded case becomes a scratch query |

The reuse to notice: **`target.attach` from Theme 2 needs no change to work
on a hit.** The hit's menu resolves the chunk behind it and attaches that.
Authoring ground truth from an ad-hoc query — story S6 — is therefore
already implemented once this theme's `hit` type exists.

## 5. Tiles

### `trail` — one result's journey (doc-bound: focus)

The centrepiece, and the transpose of the autopsy: where the autopsy shows
*stages × all candidates*, the trail shows *one candidate × all stages*.

```text
┌ ⠿ TRAIL · ptr-51ac ────────────────────────────────────────────┐
│ ▌chunk c17… "Limelight · Mature size"                          │
│ IN ▌run scratch-4f21 · ▌arm baseline                           │
│ QUERY  blue ice vs limelight — which for shade?                │
│ FINAL  rank 2 of 12 · score 0.0490                             │
│ PATH                                                           │
│  lexical.raw         rank 7    bm25 4.11   on ▌raw             │
│  lexical.collapsed   rank 5    bm25 4.11   kept, best of 3     │
│  lexical.policy      kept                  no rule matched     │
│  vector.knn          rank 2    cos 0.81    on ▌summary ⚑       │
│  fusion.rrf          rank 2    0.0490                          │
│     ▌lexical  rank 5   weight 1.0  →  0.0167                   │
│     ▌vector   rank 2   weight 1.0  →  0.0323                   │
│  final               rank 2    returned                        │
│ BEATEN BY · 1                                                  │
│  1 ▌chunk 8f2…   lex 1 · vec 3  →  0.0492   [why?]             │
│ DROPPED AT  —                                                  │
└────────────────────────────────────────────────────────────────┘
widgets  hit identity · run identity · per-stage path rows ·
         contribution rows · competitor rows
data     the run's stages (recorded: the un-coarsened projection;
         scratch: the search response) — no new endpoint
```

Serves **S10** whole. Two honesty rules are load-bearing here: a stage that
recorded only membership shows *kept* with no score rather than a fabricated
one, and a candidate that vanished between stages shows *dropped at
<stage>* without inventing a reason — the recorded data does not carry a
per-candidate drop reason today (§11).

### `split` — chunking side by side (doc-bound: focus)

```text
┌ ⠿ SPLIT · ptr-3e70 ────────────────────────────────────────────┐
│ ▌doc Blue Ice Hydrangea Care · 2,148 runes                     │
│ SETTINGS  A markdown 512/64    B markdown-heading 512/64       │
│           [+ setting]                    pure CPU · no spend   │
│ A · 6 chunks   median 402 runes   ⚠ 2 orphaned from heading    │
│ B · 4 chunks   median 537 runes     0 orphaned                 │
│ ┌─ A ─────────────────────────┐┌─ B ──────────────────────────┐│
│ │ ▌1 [0,438)     # Blue Ice…  ││ ▌1 [0,712)    # Blue Ice Care││
│ │ ▌2 [374,810)   …morning sun ││    (heading retained)        ││
│ │ ▌3 [746,1180)  ## Winter…   ││ ▌2 [712,1408) ## Winter care ││
│ │ ▌4 [1116,1554) …mulch to a… ││ ▌3 [1408,1902)## Feeding     ││
│ │    ⚠ heading lost           ││ ▌4 [1902,2148)## Problems    ││
│ │ ▌5 …  ▌6 …                  ││                              ││
│ └─────────────────────────────┘└──────────────────────────────┘│
│ SOURCE · boundaries marked inline                              │
│  # Blue Ice Hydrangea Care ⟦A1 B1⟧                             │
│  Blue Ice grows best in morning sun ⟦A2⟧ with afternoon shade… │
│  ## Winter care ⟦A3 B2⟧                                        │
└────────────────────────────────────────────────────────────────┘
widgets  document identity · setting chips · per-setting stats ·
         two boundary lists · source text with inline markers
data     POST /corpora/{id}/documents/{docId}/chunk  (pure CPU)
```

Serves **S11** whole, and it is the cheapest valuable surface in the entire
program: chunking is deterministic local computation, so this tile costs
nothing to run and needs neither a bundle nor the live executor.

### `ask` — extended with lanes (singleton)

Specified in the Theme 1–2 document as a single-lane scratchpad; Theme 3
adds lanes and the diff, which is all S13 requires.

```text
┌ ⠿ ASK ─────────────────────────────────────────────────────────┐
│ [blue ice vs limelight — which for shade?               ][ask] │
│ LANES  A ▌arm baseline    B ▌arm rrf-30    [+ lane]   limit 20 │
│ TARGETING ▌q-comparison ◀      2 queries · 624 tokens · 88 ms  │
│  rank   A · baseline               B · rrf-30                  │
│   1   ▌chunk 8f2… Blue Ice·Light   ▌chunk 8f2…           =     │
│   2   ▌chunk c17… Limelight·Size   ▌chunk a04…           ▲ new │
│   3   ▌chunk a04… Colour Guide     ▌chunk c17…           ▼ -1  │
│   4   ▌chunk d91… Pruning          —                     ✖ out │
│  DIFF   2 moved · 1 new · 1 dropped · overlap 3 of 4           │
│ [promote query into the set]   [pin]                  [clear]  │
└────────────────────────────────────────────────────────────────┘
Every row is a `hit`: its menu carries "why did this rank here",
"open the chunk", and Theme 2's three target-marking actions.
```

Serves **S9** and **S13**. The diff is computed client-side from two runs;
no endpoint knows about comparison.

### `autopsy` — extended with candidate rows (existing tile)

```text
┌ ⠿ AUTOPSY · ptr-8c40 ──────────────────────────────────────────┐
│ ▌episode 9bf95… q-comparison × rrf-30 · completed              │
│ STAGES · 6                                                     │
│ ▌lexical.raw         in 0   out 40                             │
│ ▌lexical.collapsed   in 40  out 22                             │
│ ▌lexical.policy      in 22  out 19    3 dropped                │
│ ▌vector.knn          in 0   out 40                             │
│ ▌fusion.rrf          in 59  out 51    ▼ expanded               │
│     rank  chunk              lex   vec    fused                │
│      1  ▌hit chunk 8f2…       1     3    0.0492   [why?]       │
│      2  ▌hit chunk c17…       5     2    0.0490   [why?]       │
│      3  ▌hit chunk a04…      12     4    0.0401   [why?]       │
│ ▌final               in 51  out 10                             │
│      (membership only — this stage recorded no scores)         │
└────────────────────────────────────────────────────────────────┘
```

One tile, both provenances: bound to a focus document naming an episode it
fetches the projection; naming a scratch run it reads client state.

### `chunk` — extended with the matched representation (existing tile)

```text
┌ ⠿ CHUNK · ptr-c2b9 ────────────────────────────────────────────┐
│ ▌chunk c17…   from ▌doc Limelight Hydrangea                    │
│ IN ▌run scratch-4f21   matched on ▌summary via vector.knn r2   │
│ TEXT · raw                                                     │
│  "Limelight reaches 6 to 8 feet at maturity and tolerates…"    │
│ REPRESENTATIONS · 3                                            │
│  ▌raw                        matched by lexical (rank 7)       │
│  ▌summary        ◀ MATCHED   "A large panicle hydrangea for…"  │
│     ⚑ flagged: claims 'full shade'; source says 'partial'      │
│  ▌question                   "How big does Limelight get?"     │
└────────────────────────────────────────────────────────────────┘
```

Serves **S12**. The "IN ▌run" row appears only when `activeRunId` is set —
outside an exploration workspace the tile is exactly what it is today.

## 6. Workspaces

### Probe — "what does this do, and why?" (S9, S10, S12)

```text
┌──────────────────────────────┬────────────────────┐
│ ASK                          │ TRAIL              │
│ query → hits                 │ why this ranked    │
├──────────────────────────────┼────────────────────┤
│ CHUNK                        │ INSPECTOR          │
│ text · representations ·     ├────────────────────┤
│ what matched                 │ TRACE              │
└──────────────────────────────┴────────────────────┘
```

The loop this workspace exists for: ask, see something wrong, open its
trail, discover it matched on a bad summary, open the chunk, flag the
summary. Four clicks from confusion to a recorded finding.

### Bake-off — "does this change do anything?" (S13)

```text
┌────────────────────────────────────────┬──────────────┐
│ ASK   two lanes · diff marks           │ TRAIL · A    │
│                                        ├──────────────┤
│                                        │ TRAIL · B    │
├────────────────────────────────────────┴──────────────┤
│ AUTOPSY   stage counts for the lane in focus          │
└───────────────────────────────────────────────────────┘
```

Two trail tiles on the same chunk under different lanes is the fastest way
to see *which stage* the configuration change actually moved — the answer
is usually one row differing between two otherwise identical paths.

### Chunking lab — "how should this be cut?" (S11)

```text
┌──────────────────────────────┬────────────────────┐
│ SPLIT   A vs B, one document │ DOCUMENT           │
│                              │ source · flags     │
├──────────────────────────────┼────────────────────┤
│ CORPUS   pick the next doc   │ INSPECTOR          │
└──────────────────────────────┴────────────────────┘
```

Deliberately built from two Theme 1 tiles plus one new one. The workflow is
to walk a handful of documents you know well — a long guide, a thin product
page, a FAQ — and see how each setting treats each shape.

## 7. State: where scratch runs live

Scratch runs are **client state**, not documents: a small store slice keyed
by `runId`, holding the `SearchOutput` as returned. Doc-bound tiles bind to
a `focus` pointer document naming `{runId, chunkId}` and resolve the run
from that slice — the same pattern the chunk tile already uses to resolve a
chunk through an episode's catalog.

The reasoning, stated so a future reader can overturn it deliberately:

- A `SearchOutput` carries a full chunk catalogue with text. Syncing that
  through the document host would push restricted content into a document
  store that has no sensitivity model for it.
- Scratch runs are cheap to reproduce and expected to be discarded. A
  reload losing them is the correct behaviour, not a defect.
- What must survive is recorded elsewhere already: the trace records that
  the ask happened, with its query and lanes and cost.

`run.pin` is the escape hatch: it writes a **bounded** summary — query,
lanes, top-N hit ids and scores, no text — into a comparison document, for
the case where a scratch result deserves to be referred to later.

## 8. Backend changes

### The projection fix (no new endpoint, no G1 dependency)

`StageSummary` (`pkg/ttc/specialistapi/types.go:186`) gains `candidates`:

```go
type StageCandidateSummary struct {
    ChunkID          string             `json:"chunk_id"`
    Rank             int                `json:"rank,omitempty"`
    Channel          string             `json:"channel,omitempty"`
    RepresentationID string             `json:"representation_id,omitempty"`
    Score            *float64           `json:"score,omitempty"`
    RerankerScore    *float64           `json:"reranker_score,omitempty"`
    Contributions    []rag.Contribution `json:"contributions,omitempty"`
}
```

populated in `pipeline.go:44` from the already-recorded
`ttcsearch.RetrievalStage.Candidates`. Notes for the reviewer:

- **No new sensitivity surface.** Scores, ranks and channel names are not
  restricted content; chunk *text* already flows through `chunk_catalog`
  under the existing policy and is untouched here.
- **Pointer scores stay pointers.** `Score *float64` preserves the
  distinction between "scored zero" and "no score recorded" — the same
  discipline as the Go type.
- This alone makes **S10 work on every recorded episode today**, with no
  live executor.

### Ad-hoc search (new, gated on G1)

```text
POST /api/rag/v1/search
  body   {query, arm | inline_config, limit}
  auth   bearer; new Authorizer action `search.run`
  →      full SearchOutput: stages (with candidates), chunk_catalog,
         results, and the identity of the bundle+config that answered
```

Three properties this endpoint must have:

- **Not an episode.** Nothing is journaled, nothing is scored, no campaign
  is touched. The response's identity block is what keeps a scratch result
  from ever being mistaken for a measured one.
- **Metered, not gated.** Usage commits to a `scratch` budget line and the
  tile shows cumulative spend; there is a server-side ceiling. There is
  deliberately **no per-query approval** — an approval dialog per question
  would destroy the loop this whole theme exists to enable.
- **Same executor as campaigns.** It calls the G1 live executor, so what
  you explore is what gets measured.

### Chunk-split preview (new, no dependencies at all)

```text
POST /api/rag/v1/corpora/{id}/documents/{docId}/chunk
  body   {settings: [{id, chunker, max_runes, overlap_runes,
                      min_section_runes}]}
  →      per setting: [{index, byte_start, byte_end, runes,
                        heading_path, orphaned}]
```

Implemented over `chunking.Apply` with the chunkers `selectChunker`
(`cmd/rag-ttc/cmds/indexes/build.go:57`) already offers — `markdown` and
`markdown-heading`. Pure CPU, no budget, no bundle, no model. It previews
**chunking only**: representations and embeddings are not previewed, because
those cost money and belong to the substrate ticket.

### Representation notes (extend Theme 2's format)

`ragttc.corpus-notes/v1` gains representation-scoped entries keyed by
`{chunk_id, representation_id}`, so a misleading summary is flagged where it
is seen and counted against the approach that produced it.

### Frontend registration

Two entries in `createWorkbenchApps()`, three descriptors and three
type-graph nodes, the contributions of §4, three workspaces in
`defaultLayout()`, plus the extensions to `AutopsyApp` and `ChunkApp`. The
vocabulary export regenerates; the golden test pins the new verbs.

## 9. What this design deliberately does not build

- **No new comparison tile.** S13 is two lanes in the ask tile. The existing
  `compare` tile stays what it is: recorded arms over a campaign's cases.
  A live ad-hoc comparison and a measured comparison are different claims
  and should not share a surface.
- **No menu on chunk previews beyond Inspect.** A previewed chunk exists
  in no index; offering "add as target" or "ask about this" would be
  offering operations that cannot mean anything. The near-empty menu is the
  design, not an omission.
- **No embedding or representation preview.** Previewing what a summary
  *would* say costs a model call per chunk. It belongs to the substrate
  ticket, with the invalidation plan and the cost estimate in front of it.
- **No saved query library.** `run.pin` writes one bounded summary; a
  managed collection of favourite queries is what the question set already
  is, and S6's promote path is how a query graduates into it.
- **No score visualisations beyond the numbers.** Distributions, recall
  curves, and rank-position histograms are Theme 4 (comparing arms across a
  whole question set), not Theme 3 (understanding one result).
- **No automatic explanation prose.** The trail shows recorded facts. A
  generated "this ranked low because…" sentence would be a model's opinion
  rendered as a system fact, which is exactly the confusion the trace
  discipline exists to prevent — the agent can say it in chat, where it is
  visibly the agent speaking.

## 10. Dependencies and honest sequencing

The useful property of this theme is that most of it is **not** gated on the
live executor:

| Piece | Needs G1? | Notes |
|---|---|---|
| `StageSummary` candidates | no | works on recorded fixture episodes today |
| `trail` tile (recorded runs) | no | reads the fixed projection |
| `autopsy` candidate rows | no | same |
| `split` tile + chunk preview endpoint | no | pure CPU; needs only document text |
| `chunk` matched-representation row | no | recorded runs carry `representation_id` |
| `ask` tile, lanes, diff | **yes** | needs real retrieval to be worth anything |
| `trail` on scratch runs | **yes** | same source |

Suggested order: (1) the projection fix plus the trail tile and autopsy
candidate rows — this is a self-contained ticket that makes recorded
episodes explicable and ships before anything else; (2) the split tile and
its endpoint, which is cheap, standalone, and immediately useful on the real
corpus; (3) the ask tile and lanes when G1 lands; (4) the chunk tile's
matched-representation row, which is small and rides along.

## 11. Open questions

- **Per-candidate drop reasons.** `StageCandidate` records no reason for a
  candidate disappearing at a filtering stage; `ErrorClass` is stage-level.
  The trail can therefore say *dropped at lexical.policy* but not *by which
  rule*. Adding a `DropReason` to the recorded candidate is a small,
  additive change to the retrieval service and would complete S10 — worth
  raising with the retrieval owner rather than deciding here.
- **Scratch budget ceiling.** What is the right daily ceiling for the
  exploration line, and should it be per-principal? Proposed: a generous
  ceiling visible in the tile, since the failure mode to avoid is a
  reluctant operator, not an expensive one.
- **Should the agent see scratch runs?** Proposed no for v1: the agent can
  read the trace and issue its own `ask.run`. Making runs agent-visible
  means putting a bounded projection somewhere shared, which is `run.pin`
  generalised — worth doing only if the agent seat actually needs it.
- **Does the split preview belong to a bundle?** Today it chunks raw
  document text with given settings, independent of any built bundle. Once
  substrate variation becomes a campaign dimension, the preview should be
  able to say *this is exactly what bundle X did* — which requires the
  bundle to record its chunker settings in a form the preview can replay.
  `indexbundle.ChunkerIdentity` already carries name and maximum runes, so
  this looks reachable; confirm the remaining parameters are recorded.
