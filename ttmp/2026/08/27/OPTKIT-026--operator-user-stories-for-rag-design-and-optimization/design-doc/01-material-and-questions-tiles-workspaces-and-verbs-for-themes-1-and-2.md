---
Title: 'Material and Questions: Tiles, Workspaces, and Verbs for Themes 1 and 2'
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
    - Path: repo://rag-ttc/apps/workbench/web/src/pbui/types.ts
      Note: The presentation vocabulary this design extends with five types
    - Path: repo://rag-ttc/apps/workbench/web/src/pbui/actions.ts
      Note: The action registry — every verb below is a contribution here
    - Path: repo://rag-ttc/apps/workbench/web/src/pbui/verbs.ts
      Note: The product verb union the new families join
    - Path: repo://rag-ttc/pkg/ttc/workbenchhost/documents.go
      Note: Document format validators — two new formats land here
    - Path: repo://rag-ttc/pkg/ttc/experimentworkbench/workbench_service.go
      Note: The compile/preview/seal service the questionset seal mirrors
    - Path: repo://rag-ttc/cmd/rag-ttc/cmds/indexes/settings.go
      Note: The build inputs corpus exclusions must eventually reach
ExternalSources: []
Summary: A minimal design covering user-story themes 1 and 2 — five new presentation types, five new tiles, three workspace arrangements, two document formats, and the backend projections behind them — built by reusing the proposal-draft authoring pattern rather than inventing a second one.
WhatFor: The implementation contract for the corpus-knowledge and question-authoring slice of the workbench.
WhenToUse: Before implementing any corpus or question-set surface; as the acceptance vocabulary for stories S1–S8.
---

# Material and Questions: Tiles, Workspaces, and Verbs for Themes 1 and 2

## 1. Scope and the discipline applied

This document covers stories **S1–S3** (knowing my material) and **S4–S8**
(defining what good means) from the user-story set. Nothing else. Two
constraints shaped every decision below:

- **Do not over-engineer.** Each story is served by the smallest surface that
  makes it true. Where a story wants something elaborate, this design takes
  the honest cheap version and says so in §9.
- **Leverage what exists.** The workbench already has an action kernel, a
  presentation vocabulary, a document-sync layer, an authoring-document
  pattern with server-side compilation and durable sealing, an approval
  flow, and an agent seat. A new surface that does not reuse those is a
  design failure, not a feature.

The central claim of this design is one sentence: **a question set is a
proposal draft.** It is a document that a human and an agent co-edit through
verbs, whose durable form is produced by an explicit seal, whose sealed
result is content-addressed and referenced by everything measured against
it. That is exactly the shape OPTKIT-023 built and OPTKIT-024 gave an agent
seat to. Recognising the isomorphism is what makes this slice small.

## 2. What the design adds, in total

| Kind | Count | Names |
|---|---|---|
| Presentation types | 5 | `corpus`, `document`, `questionSet`, `question`, `target` |
| Tiles | 5 | `corpus`, `document`, `questions`, `question`, `ask` |
| Workspaces | 3 | Material, Questions, Curation |
| Document formats | 2 | `ragttc.questionset-draft/v1`, `ragttc.corpus-notes/v1` |
| Verb families | 2 | `question.*` / `target.*` (draft family), `corpus.*` (local + command) |
| Read projections | 4 | corpus survey, document list, document detail, question-set versions |
| Command-API operations | 2 | `questionset.seal`, `coverage.sweep` |

Everything else is reuse: the inspector, watchlist and trace tiles are
unchanged; the `chunk` type gains two contributions and nothing more; the
approval flow, the document sync, the verb router, the agent seat, and the
trace all work without modification.

## 3. Semantic objects

New presentation values, following the existing rule from `pbui/types.ts`:
*a value carries what its menu needs to decide — ids plus the few fields
labels and rules read, never full projections.*

```ts
export interface CorpusRef {
  corpusId: string;        // "ttc-site-2026-08"
  digest: string;          // content digest of the document set
  documentCount?: number;
}

export interface DocumentRef {
  corpusId: string;
  documentId: string;
  title?: string;
  kind?: string;           // guide | product | faq | post
  runes?: number;
  flag?: "note" | "exclude" | null;   // what the menu label needs
}

export interface QuestionSetRef {
  setDocId: string;        // the DRAFT document being edited
  setId: string;           // "ttc-live-v1"
  questionCount?: number;
  dirty?: boolean;         // unsealed edits exist
}

export interface QuestionRef {
  setDocId: string;
  questionId: string;
  text?: string;
  gold?: boolean;
  targetCount?: number;
  coverage?: "covered" | "partial" | "uncovered" | "unknown";
}

export interface TargetRef {
  setDocId: string;
  questionId: string;
  targetType: "document" | "chunk";
  targetId: string;
  grade: "essential" | "helpful" | "non-answer";
}
```

Type-graph placement in `workbenchTypeDefinitions` (`pbui/actions.ts`):

```ts
{ id: "corpus",      parents: ["inspectable"] },
{ id: "document",    parents: ["watchable"]  },   // watch a document you are tracking
{ id: "questionSet", parents: ["inspectable"] },
{ id: "question",    parents: ["watchable"]  },   // watch a question you are fighting with
{ id: "target",      parents: ["inspectable"] },
```

Two new facts join `WorkbenchFacts`, mirroring `activeDraftDocId` and
`draftEvidenceKeys` exactly:

```ts
  /** The question "mark as target" verbs land on. */
  activeQuestionId: { setDocId: string; questionId: string } | null;
  /** Stable keys already attached as targets on the ACTIVE question — so
   *  "Remove target" appears contextually wherever that object is shown. */
  activeTargetKeys: ReadonlySet<string>;
  /** Documents carrying an operator flag, by stable key. */
  flaggedKeys: ReadonlySet<string>;
```

The consequence is the design's best property: because rules read facts and
not tables, **"Mark as essential target" appears on a chunk in the ask tile,
on a chunk in the autopsy tile, on a document in the corpus tile, and on a
mention in chat — from one contribution.**

## 4. Verbs, by object

New verb kinds join the existing union in `pbui/verbs.ts`. Draft-family
verbs mutate a document and never touch the store; command-family verbs call
the API; danger verbs are gated by the existing approval flow.

```ts
export type CorpusVerb =
  | { kind: "open.corpus";      corpusId: string }
  | { kind: "open.document";    corpusId: string; documentId: string }
  | { kind: "corpus.flag";      corpusId: string; documentId: string;
                                flag: "note" | "exclude"; reason: string }
  | { kind: "corpus.unflag";    corpusId: string; documentId: string }
  | { kind: "corpus.exportExclusions"; corpusId: string };        // COMMAND

export type QuestionVerb =
  | { kind: "open.questions";   setDocId: string }
  | { kind: "open.question";    setDocId: string; questionId: string }
  | { kind: "question.activate"; setDocId: string; questionId: string }   // LOCAL
  | { kind: "question.add";     setDocId: string; text: string }
  | { kind: "question.setText"; setDocId: string; questionId: string; text: string }
  | { kind: "question.setRubric"; setDocId: string; questionId: string; text: string }
  | { kind: "question.setCategory"; setDocId: string; questionId: string; category: string }
  | { kind: "question.remove";  setDocId: string; questionId: string }
  | { kind: "question.markGold"; setDocId: string; questionId: string;
                                 score?: number; note?: string }
  | { kind: "question.unmarkGold"; setDocId: string; questionId: string }
  | { kind: "question.ask";     setDocId: string; questionId: string }    // LOCAL
  | { kind: "target.attach";    setDocId: string; questionId: string;
                                ref: Ref; grade: "essential" | "helpful" | "non-answer" }
  | { kind: "target.setGrade";  setDocId: string; questionId: string;
                                targetKey: string; grade: string }
  | { kind: "target.remove";    setDocId: string; questionId: string; targetKey: string }
  | { kind: "questionset.seal"; setDocId: string; idempotencyKey: string;
                                approvalId?: string }                     // DANGER
  | { kind: "coverage.sweep";   setDocId: string; armId: string;
                                approvalId?: string };                    // DANGER

export type AskVerb =
  | { kind: "ask.run";   query: string; armId: string; limit: number }    // COMMAND
  | { kind: "ask.clear" }                                                 // LOCAL
  | { kind: "ask.promote"; setDocId: string; query: string };             // draft
```

Menu contributions, per type — this is the table a reviewer should read
against `pbui/actions.ts`:

| Object | Label | Verb | Notes |
|---|---|---|---|
| `corpus` | Open survey (primary) | `open.corpus` | |
| | Export exclusions to build inputs | `corpus.exportExclusions` | command; writes a build input file |
| | Run coverage sweep… | `coverage.sweep` | **danger** — spends embeddings |
| `document` | Open document (primary) | `open.document` | |
| | Mark as essential target | `target.attach` grade=essential | unavailable without an active question |
| | Mark as helpful target | `target.attach` grade=helpful | |
| | Mark as non-answer | `target.attach` grade=non-answer | |
| | Remove target | `target.remove` | shown only when in `activeTargetKeys` |
| | Flag as note… / Flag as excluded… | `corpus.flag` | label flips to "Remove flag" when flagged |
| | Add to watchlist | inherited | |
| `questionSet` | Open questions (primary) | `open.questions` | |
| | Add question… | `question.add` | |
| | Seal version | `questionset.seal` | **danger** — durable |
| `question` | Open question (primary) | `open.question` | |
| | Make active | `question.activate` | the target-attach destination |
| | Ask this question | `question.ask` | prefills the ask tile |
| | Mark gold / Unmark gold | `question.markGold` | |
| | Remove question | `question.remove` | |
| | Add to watchlist | inherited | |
| `target` | Open the document | `open.document` | |
| | Set grade → essential/helpful/non-answer | `target.setGrade` | |
| | Remove | `target.remove` | |
| `chunk` *(existing)* | Mark as … target | `target.attach` | **the reuse win**: same three contributions, defined once |
| | Remove target | `target.remove` | |

Availability follows the existing kernel discipline: the target verbs are
`unavailable` with the reason *"no active question"* rather than hidden, so
an agent planning against the vocabulary is told why, and a human sees a
greyed row that explains itself.

## 5. Tiles

Five new tiles. Each is small, each reads one projection or one document,
and each is composable into more than one workspace.

### `corpus` — the survey (doc-bound: focus)

```text
┌ ⠿ CORPUS · ptr-a41f ───────────────────────────────────────────┐
│ ▌corpus ttc-site-2026-08 · sha256:9ac2f1…      1,257 documents │
│ KINDS    guide 412 · product 623 · faq 141 · post 81           │
│ LENGTH   ▁▃▇▅▂▁   median 1.4k · p95 9.2k · 37 under 200 runes ⚠│
│ FLAGS    ⚑18 excluded · ⚑6 outdated note · ⚑2 restricted       │
│ NEAR-DUPLICATE   9 clusters · 214 documents                    │
│  ▌cluster 03 · 41 docs · "Care & handling" boilerplate         │
│  ▌cluster 07 · 22 docs · product blurb repeated verbatim       │
│ DOCUMENTS · showing 40 of 1,257                                │
│  [kind ▾] [length ▾] [flagged ▾] [used-as-target ▾]            │
│  ▌doc Blue Ice Hydrangea Care    guide    2.1k                 │
│  ▌doc Emerald Green Arborvitae   guide    1.7k                 │
│  ▌doc Shipping FAQ               faq      0.4k  ⚑ excluded     │
│  ▌doc Blue Ice (product page)    product  0.2k  ⚠ thin         │
└────────────────────────────────────────────────────────────────┘
widgets  identity tray · kind counters · length histogram ·
         flag counters · cluster chips · filtered document list
data     GET /corpora/{id}/survey · GET /corpora/{id}/documents
rows     every document and cluster is a live presentation
```

Serves **S1** whole, **S3** through the flag counters and filters, and gives
**S2** its landing place once the coverage sweep exists.

### `document` — one piece of material (doc-bound: focus)

```text
┌ ⠿ DOCUMENT · ptr-c2b9 ─────────────────────────────────────────┐
│ ▌doc Blue Ice Hydrangea Care                                   │
│  source   https://thetreecenter.com/blue-ice-hydrangea         │
│  kind guide · 2,148 runes · updated 2026-06-14                 │
│  ▌corpus ttc-site-2026-08                                      │
│ FLAGS    none                    [note…]  [exclude…]           │
│ TEXT                                                           │
│  # Blue Ice Hydrangea Care                                     │
│  Blue Ice grows best in morning sun with afternoon shade in    │
│  zones 7 and warmer…                                           │
│  ## Winter care                                                │
│  Mulch to a depth of three inches after the first hard freeze… │
│ USED AS TARGET BY · 3 questions                                │
│  ▌q-overwinter essential · ▌q-mulch helpful · ▌q-shade helpful │
└────────────────────────────────────────────────────────────────┘
widgets  identity tray · flag bar · full text · back-reference chips
data     GET /corpora/{id}/documents/{docId}  + the set draft document
```

Serves **S1** step 3, **S3** steps 1–3, and — through *used as target by* —
gives **S7** its blast-radius answer for free: editing this document's role
shows which questions depend on it.

### `questions` — the set (doc-bound: the question-set draft)

```text
┌ ⠿ QUESTIONS · set-draft-7f21 ──────────────────────────────────┐
│ ▌set ttc-live-v1 · sealed v3 · 4 unsealed edits                │
│ 38 questions · ★11 gold · ⚠4 without targets                   │
│ COVERAGE  ✔29 answerable · ◐5 partial · ✖4 uncovered           │
│ [all | gold | no-targets | uncovered]            [+ question]  │
│  ★ ▌q-overwinter    "how do I overwinter a blue…"  3 ✔         │
│    ▌q-comparison    "blue ice vs limelight — sha…" 2 ✔    ◀    │
│    ▌q-inventory     "do you have 5ft arborvitae…"  0 ✖ ⚠       │
│    ▌q-mulch-depth   "how deep should mulch be…"    1 ◐         │
│  ★ ▌q-shade-list    "which hydrangeas take shade…" 4 ✔         │
│ ── seal bar ──────────────────────────────────────────────────  │
│ Sealing writes evaluation set v4 — durable, content-addressed,  │
│ and referenced by every run measured afterwards.                │
│                        [SEAL VERSION — durable]  [discard v4]   │
└────────────────────────────────────────────────────────────────┘
◀ marks the ACTIVE question — where "mark as target" lands
widgets  set identity · counters · coverage row · filter chips ·
         question list · seal bar (the OPTKIT-023 SealBar, reused)
```

Serves **S4** and **S8**, and is the home of **S7**'s versioning.

### `question` — one question and its truth (doc-bound: focus)

```text
┌ ⠿ QUESTION · ptr-9d10 ─────────────────────────────────────────┐
│ ▌q-comparison            ACTIVE ◀      [★ gold]  [ask this]    │
│ TEXT     [blue ice vs limelight hydrangea — which for shade?]  │
│ CATEGORY comparison            WEIGHT  high                    │
│ RUBRIC                                                         │
│  [Evidence must let someone say which of the two tolerates     │
│   shade, and how they differ in mature size.]                  │
│ TARGETS · 2                                                    │
│  ▌doc Blue Ice Hydrangea Care     essential  [grade ▾]  [×]    │
│  ▌doc Limelight Hydrangea         essential  [grade ▾]  [×]    │
│ NON-ANSWERS · 1                                                │
│  ▌doc Hydrangea Colour Guide      looks relevant, is not       │
│ MY LABEL (gold)   0.8   "both retrieved; shade claim is thin"  │
│ COVERAGE  ✔ answerable · checked 2026-08-27 against ▌arm base  │
└────────────────────────────────────────────────────────────────┘
widgets  question identity · prose fields (the IntentApp field
         pattern, reused) · target rows · non-answer rows ·
         gold label · coverage line
```

Serves **S5** whole and **S8**'s labelling.

### `ask` — the query scratchpad (singleton)

```text
┌ ⠿ ASK ─────────────────────────────────────────────────────────┐
│ [blue ice vs limelight hydrangea — which for shade?      ][ask]│
│ AGAINST ▌arm baseline · limit 20 (deliberately generous)       │
│ TARGETING ▌q-comparison ◀        marks land on this question   │
│ 12 results · 1 query · 312 embed tokens · 41 ms                │
│  1 ▌chunk 8f2… Blue Ice · Light needs      0.81 [ess][help][✖] │
│  2 ▌chunk c17… Limelight · Mature size     0.74 [ess][help][✖] │
│  3 ▌chunk a04… Hydrangea Colour Guide      0.66 [ess][help][✖] │
│  4 ▌chunk d91… Pruning panicle hydrangeas  0.51 [ess][help][✖] │
│ [promote this query into the set]                    [clear]   │
└────────────────────────────────────────────────────────────────┘
Nothing here is recorded as an experiment. The mark buttons are the
same target.attach contributions the object menus carry.
```

Serves **S6** — and it is the reason S6 is achievable at all. It belongs to
Theme 3 as much as Theme 2; it is specified here because ground-truth
authoring without it is data entry against document identifiers nobody
knows.

## 6. Workspaces

Workspaces already exist in the layout document (`layout(..., {id, name})`
plus `selectWorkspace` in `workbench.ts`); today the product ships one,
named *Evidence*. This design adds three, and every one of them is built
from the five new tiles plus the three existing ambient singletons.

### Material — "what am I working with?" (S1, S3)

```text
┌────────────────────────────────────┬────────────────────┐
│ CORPUS                             │ DOCUMENT           │
│ survey · clusters · filtered list  │ text · flags ·     │
│                                    │ used-as-target-by  │
│                                    ├────────────────────┤
│                                    │ INSPECTOR          │
│                                    ├────────────────────┤
│                                    │ TRACE              │
└────────────────────────────────────┴────────────────────┘
```

Click a document in the survey; it opens in the detail tile through the
`focus` pointer document — the same navigation mechanism the failures and
autopsy tiles already use. Flags applied here appear on that document
everywhere else, forever.

### Questions — "what does good mean?" (S4, S5, S7, S8)

```text
┌───────────────┬──────────────────────────┬───────────────┐
│ QUESTIONS     │ QUESTION                 │ INSPECTOR     │
│ set · filters │ text · rubric · targets  ├───────────────┤
│ · seal bar    │ · gold label             │ WATCHLIST     │
├───────────────┴──────────────────────────┼───────────────┤
│ ASK   query → results → mark as target   │ TRACE         │
└──────────────────────────────────────────┴───────────────┘
```

This is the workspace where most human hours go. The trace tile matters
here for a reason that is not obvious: it is where the human sees what the
**agent** did to the set — every `question.add` and `target.attach` with its
actor — which is S36 arriving for free.

### Curation — "find the targets retrieval is missing" (S6 step 4)

```text
┌──────────────────────────┬───────────────────────────┐
│ ASK                      │ QUESTION  (active)        │
│ several phrasings, one   │ targets accumulating      │
│ question                 │                           │
├──────────────────────────┼───────────────────────────┤
│ CORPUS   filtered browse │ DOCUMENT                  │
│ for what search missed   │ read it, then mark it     │
└──────────────────────────┴───────────────────────────┘
```

The point of this arrangement is that both routes to a target — *it came
back from a query* and *I found it by browsing* — end in the same verb on
the same object type, so the resulting judgment is identical regardless of
how the human got there.

## 7. Documents and state

### `ragttc.questionset-draft/v1` — the authoring document

Body shape (validated in `workbenchhost/documents.go` alongside the existing
`validateProposalDraft`, using the same opaque-id and prose-field helpers):

```json
{
  "schema_version": 1,
  "set_id": "ttc-live-v1",
  "corpus_id": "ttc-site-2026-08",
  "sealed_version": "sha256:…",
  "questions": [
    {
      "id": "q-comparison",
      "text": "blue ice vs limelight hydrangea — which for shade?",
      "category": "comparison",
      "weight": "high",
      "rubric": "Evidence must let someone say which tolerates shade…",
      "gold": { "score": 0.8, "note": "both retrieved; shade claim thin" },
      "targets": [
        { "type": "document", "id": "doc-blue-ice", "grade": "essential" },
        { "type": "document", "id": "doc-limelight", "grade": "essential" },
        { "type": "document", "id": "doc-colour-guide", "grade": "non-answer" }
      ],
      "coverage": { "status": "covered", "checked": "2026-08-27", "arm": "baseline" }
    }
  ]
}
```

Three rules inherited from the proposal-draft design, and worth restating
because they are what keep this cheap:

- **Derived state is never stored.** Coverage counters, per-question target
  counts, and the "without targets" warning are computed at render. Only
  the coverage *sweep result* is stored, because it is a measurement, not a
  derivation.
- **Documents hold bounded opaque ids.** A target holds a document id, not a
  document; the tile re-materializes the presentation at read time.
- **Sealing is the only durable act.** Everything above is a working
  document; nothing measured references it until a version is sealed.

### `ragttc.corpus-notes/v1` — operator annotations

```json
{
  "schema_version": 1,
  "corpus_id": "ttc-site-2026-08",
  "notes": [
    { "document_id": "doc-shipping-faq", "flag": "exclude",
      "reason": "navigation boilerplate, no product content",
      "at": "2026-08-27T10:11:00Z" }
  ]
}
```

Deliberately a workbench document rather than a backend annotation service:
it needs no new storage, syncs and merges through the existing host, and is
visible to the agent. Flags become build inputs only through the explicit
`corpus.exportExclusions` verb (§8), so an annotation never silently changes
what gets indexed.

## 8. Backend changes

### Read projections (new)

The existing specialist API is campaign-scoped (`/api/rag/v1/campaigns/…`).
Corpora and question sets are not campaign-scoped, so they land as sibling
groups on the same server, following the same projection discipline
(recorded facts only, missing is never zero, a typed schema per view).

| Endpoint | Returns | Serves |
|---|---|---|
| `GET /api/rag/v1/corpora` | known corpora with digests and counts | S1 |
| `GET /api/rag/v1/corpora/{id}/survey` | counts by kind and source, length histogram buckets, thin-document count, near-duplicate clusters | S1 |
| `GET /api/rag/v1/corpora/{id}/documents?kind=&flagged=&q=&limit=` | paged document list | S1, S6 |
| `GET /api/rag/v1/corpora/{id}/documents/{docId}` | metadata plus full text | S1, S3, S6 |
| `GET /api/rag/v1/questionsets/{id}/versions` | sealed versions with digests and question counts | S7 |

The survey is computed once at startup from the corpus the server was
pointed at, and cached: `serve` gains `--corpus <path-or-bundle>`. Near-
duplicate clustering is a shingle-hash pass over document text — a hundred
lines, no new dependency, and honest about being approximate.

### Command-API operations (new)

Both go through the existing bearer-authenticated command API and the
existing `Authorizer` action model, which means the agent principal gets
them or does not get them by grant rather than by convention.

| Operation | Action | Semantics |
|---|---|---|
| `questionset.seal` | `questionset.seal` (new `Action`) | Validates the draft (every question has text; targets reference existing documents), writes the evaluation set as a **store-level artifact** (the artifact store is already store-scoped, not campaign-scoped — see `local.Open(store).Artifacts`), and returns its digest. Idempotent on `IdempotencyKey`, exactly like `proposals:seal`. |
| `coverage.sweep` | `preview.run` (reused) | Runs one generous retrieval per question against a named arm and records per-question `covered / partial / uncovered`. Budgeted: reserves before, commits after, on the same ledger as everything else. |

`questionset.seal` is a **danger verb**: an agent may compose a set all day
and may not make one durable. That is the OPTKIT-024 approval flow with no
new mechanism — the verb carries `approvalId`, the router parks it, the seal
bar renders the request.

### Manifest and build wiring

- A campaign manifest gains `evaluation_set: {id, version_digest}`, so every
  result carries the definition of "correct" it was measured against. This
  is what makes **S20** (only compare comparable things) mechanical later.
- `corpus.exportExclusions` writes an exclusion list next to the corpus,
  consumed by `rag-ttc indexes build` on the next build. The bridge from
  annotation to build input is one deliberate operator act, visible in the
  trace.

### Frontend registration

Standard, and small: five entries in `createWorkbenchApps()`, five
descriptors in `pbui/registry.ts`, five type-graph nodes and the
contributions of §4 in `pbui/actions.ts`, three workspaces in
`defaultLayout()`, two format validators in `workbenchhost/documents.go`,
and the corresponding app catalog entries in `workbenchhost/catalog.go`.
The vocabulary export (`pnpm vocab`) regenerates and the golden test pins
the new verbs — which is also how the agent learns they exist.

## 9. What this design deliberately does not build

Each of these was considered and cut. They are recorded so that adding them
later is a decision rather than a drift.

- **No annotation service.** Flags are a synced document. If flags ever need
  cross-workspace queries or per-user attribution beyond the trace, that is
  a real backend feature and should be argued for then.
- **No duplicate-cluster resolution workflow.** The survey lists clusters
  and lets you open documents. Merging, canonical-choosing, and bulk
  exclusion are not built; the filter plus the flag verb covers the actual
  need, which is *stop indexing this boilerplate*.
- **No question-import wizard.** Bulk import is a CLI command that writes a
  draft document. The UI grows questions one at a time and by promotion
  from the ask tile.
- **No chunk-boundary rendering in the document tile yet.** It belongs to
  S11 (Theme 3) and needs the substrate work; the document tile shows text
  and stops.
- **No coverage automation.** The sweep is a verb the human runs and pays
  for, not a background job. Automatic re-checking is a Theme 7 concern.
- **No separate gold-labelling mode.** A gold label is a field on a
  question, set from the question tile while looking at it. A dedicated
  labelling queue is exactly the kind of surface that gets built and never
  used.
- **No new list-tile abstraction.** The temptation is a generic
  filterable-list tile parameterised by source. Tiles in this product are
  semantic; a generic one would make every menu conditional and would be
  the first thing to fight the action kernel.

## 10. Dependencies and honest sequencing

- The `ask` tile and `coverage.sweep` both require **real retrieval over the
  real bundle** — OPTKIT-025's G1 live executor. Before G1 they can run only
  against fixtures, which is enough to build and test the surface but not
  enough to author real ground truth.
- Everything else — corpus survey, document detail, flags, the question-set
  draft, the question tile, sealing — has **no dependency on G1** and can be
  built immediately. That is a genuinely useful property: the highest-
  leverage human work (S4, S5, S8) is unblocked today.
- Suggested order: (1) question-set draft document + `questions` and
  `question` tiles + seal, because it reuses the most and blocks the most;
  (2) corpus projections + `corpus` and `document` tiles + flags; (3) the
  `ask` tile, when G1 lands; (4) `coverage.sweep`, last, because it is the
  only piece that spends money.

## 11. Open questions

- **Where do sealed evaluation sets live for discovery?** They are store-
  level artifacts, which gives content-addressing and verification for free,
  but there is no index of "which sets exist". A small metadata record per
  seal is probably right; confirm against how candidates are indexed today.
- **Do documents need their own stable identity across corpus rebuilds?**
  Targets reference document ids. If a re-ingest changes ids, every target
  breaks. Recommendation: document id derives from source URI, not from
  ingest order — needs confirming against the ingest that produces
  `corpus.json`.
- **Should `question` be watchable?** It is proposed as watchable so a
  question you are fighting with can ride along in the watchlist. That adds
  a sixth watchable type; harmless, but it is a vocabulary decision.
- **Chunk-level targets.** The design allows `targetType: "chunk"` but the
  recommendation stands that early ground truth is document-level, because
  chunk ids do not survive a chunking change — which is the variable most
  worth sweeping. The type permits both; the guidance says start coarse.
