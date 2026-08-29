---
Title: 'Running the Work: Jobs, Phases, Units, and Conversations for Theme 9'
Ticket: OPTKIT-026
Status: active
Topics:
    - design
    - ui
    - rag-ttc
    - backend
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://flowkit/flow/report.go
      Note: StepReport — per-phase counters, retries by class, spend and meters this design projects
    - Path: ws://flowkit/flow/classify.go
      Note: ErrorClass (transient/data/fatal) — the failure grouping primitive that already exists
    - Path: ws://ragopt/pkg/runstore/types.go
      Note: Manifest, Status, Summary — the build run record made observable while running
    - Path: ws://optkit/episode/types.go
      Note: Event.Span and Event.Parent — the span tree that makes turns representable today
    - Path: ws://optkit/measure/observation.go
      Note: SubjectRef{Kind,ID} is generic — turn-level observations need no schema change
    - Path: repo://rag-ttc/pkg/ttc/optkitcampaign/campaign.go
      Note: work_items, leases and budget — the second work engine this design projects from
ExternalSources: []
Summary: A design for stories S37-S48 that unifies heterogeneous long-running work at the projection layer rather than rewriting either execution engine, and shows that multi-turn evaluation needs a new executor and scenario cases but no change to the measurement chain.
WhatFor: The implementation contract for the operational and multi-turn evaluation surfaces of the workbench.
WhenToUse: Before implementing job progress, control, or conversation evaluation; as the acceptance vocabulary for stories S37-S48.
---

# Running the Work: Jobs, Phases, Units, and Conversations for Theme 9

## 1. Scope and the two findings that shaped it

This document covers **S37–S48**: seeing and steering long-running work
(builds, embedding runs, judge batches, campaigns, sweeps) and evaluating
multi-turn conversations.

The same discipline applies as in the previous two documents — smallest
surface per story, reuse over invention. Reading the code produced two
findings that changed the design substantially.

**Finding one: there are two work engines, and the honest move is to unify
them in the projection, not in the runtime.**

| | Engine A — campaigns | Engine B — builds and experiments |
|---|---|---|
| Where | `pkg/ttc/optkitcampaign` | `flowkit/flow` + `ragopt/pkg/runstore` |
| Unit of work | episode (`work_items` rows with leases) | flow item within a named step |
| Progress | journal events, queryable live | `flow.Report`, **produced at the end** |
| Failure model | `failed_terminal` + attempts | `ErrorClass{transient, data, fatal}`, quarantine records |
| Budget | reserve/commit per episode | `flow.Resource{Ceiling, Budget, UnitUSD}` per step |
| Identity | campaign journal | `runstore.Manifest` with inputs copied and hashed |

Both are good. Rewriting either to serve a progress view would be
destructive and slow. Both already produce almost exactly the data the
stories want — `flow.StepReport` carries `{Items, Hits, Misses, WorkCalls,
Retries, Quarantined, Skipped, RetriesByClass, Spend, Meters}`, which is
S38 and S40 nearly in full. What is missing is that engine B's report is
**terminal**: you get it when the run finishes, as files on disk. So the
design defines one read-side vocabulary that both engines project into, and
makes exactly one behavioural change — engine B publishes its in-flight
report periodically instead of only at the end.

**Finding two: multi-turn evaluation does not need the measurement chain
extended.** The story document's closing observation said it did; reading
the code shows a more encouraging picture. `episode.Event` already carries
`Span record.SpanID` and `Parent *record.SpanID` — a trajectory is a span
tree, so a conversation's turns are already representable as spans with
retrieval events nested beneath them. And `measure.SubjectRef` is
`{Kind string, ID string}` — generic, not typed to episodes — so a
turn-level observation is `SubjectRef{Kind: "turn", ID: "<episode>/<span>"}`
with no schema change at all. What multi-turn actually needs is a
**conversation executor**, **scenario cases**, and **projections** — the
same three things single-query evaluation needed, not a new measurement
theory.

## 2. What the design adds, in total

| Kind | Count | Names |
|---|---|---|
| Presentation types | 5 | `job`, `phase`, `unit`, `failureGroup`, `turn` |
| New tiles | 4 | `jobs`, `job`, `unit`, `conversation` |
| Kind panels (inside `job`) | 4 | build, campaign, judge, conversations |
| Workspaces | 3 | Operations, Judging live, Conversations |
| Read projections | 5 | work list, job detail, units, failure groups, conversation |
| Live stream | 1 | work SSE |
| Command operations | 6 | start, pause, stop, resume, retry group, set concurrency |
| Executor | 1 | `rag.ttc-conversation/v1` |

Reused unchanged: Theme 3's `trail` and `chunk` tiles serve a turn's
retrieval, so S44 step 4 — *"open that turn's retrieval with the same
explanation surfaces a single query gets"* — is satisfied by tiles that
already exist. A scenario is a **question with follow-ups**, so Theme 2's
entire authoring machinery (rubric, targets, gold labels, sealing) covers
scenario authoring with one added field.

## 3. Semantic objects

```ts
/** One intention: build this index, run this campaign, judge these
 *  episodes, run this conversation suite. */
export interface JobRef {
  jobId: string;
  kind: "build" | "campaign" | "judge" | "conversations" | "sweep";
  label: string;
  state: "queued" | "running" | "paused" | "stopped" | "done" | "failed";
  /** What it waits on; a queued job with unmet dependencies cannot start. */
  waitsOn?: string[];
  supervised?: boolean;      // can this UI control it, or only watch it?
}

/** An ordered stage within a job. Engine B: a flow step. Engine A: a
 *  journal-derived stage (plan / execute / score). */
export interface PhaseRef {
  jobId: string;
  phase: string;
  state: "pending" | "running" | "done" | "failed";
  done?: number;
  total?: number;
  ratePerMin?: number;
}

/** One item of work: a chunk summarization, an episode, a judge call, a
 *  conversation. */
export interface UnitRef {
  jobId: string;
  phase: string;
  unitId: string;
  state: "queued" | "running" | "done" | "failed";
  attempt?: number;
  subjectLabel?: string;     // "chunk-8f2…", "q-comparison × rrf-30"
}

/** Units grouped by why they failed — the diagnosis, not the list. */
export interface FailureGroupRef {
  jobId: string;
  cause: string;             // "rate_limit", "context_too_long", …
  class: "transient" | "data" | "fatal";
  count: number;
  growing?: boolean;
}

/** One turn of a conversation. Backed by a trajectory span. */
export interface TurnRef {
  campaignId?: string;
  runId: string;             // the conversation run (Theme 3's run id)
  turnId: string;            // the span id
  index: number;
  text?: string;
  score?: number;
  status?: "measured" | "failed" | "unknown";
  firstDegradation?: boolean;
}
```

Type-graph placement:

```ts
{ id: "job",          parents: ["watchable"]  },   // watch a long build
{ id: "phase",        parents: ["inspectable"] },
{ id: "unit",         parents: ["inspectable"] },
{ id: "failureGroup", parents: ["inspectable"] },
{ id: "turn",         parents: ["watchable"]  },   // watch a turn that keeps failing
```

Two facts join `WorkbenchFacts`:

```ts
  /** Jobs this UI supervises and may therefore control. A job started from
   *  a terminal is visible but not controllable — its control verbs are
   *  unavailable with that reason, never rendered as dead buttons. */
  supervisedJobIds: ReadonlySet<string>;
  /** The bearer principal holds the job-control grant. */
  canControlJobs: boolean;
```

## 4. Verbs, by object

```ts
export type JobVerb =
  | { kind: "open.jobs" }
  | { kind: "open.job";        jobId: string }
  | { kind: "open.unit";       jobId: string; unitId: string }
  | { kind: "job.start";       jobId: string; approvalId?: string }  // DANGER
  | { kind: "job.pause";       jobId: string }                       // COMMAND
  | { kind: "job.stop";        jobId: string }                       // COMMAND
  | { kind: "job.resume";      jobId: string; approvalId?: string }  // DANGER
  | { kind: "job.retryGroup";  jobId: string; cause: string; approvalId?: string }
  | { kind: "job.setConcurrency"; jobId: string;
                                  workers?: number; batchSize?: number }
  | { kind: "failures.acknowledge"; jobId: string; cause: string };  // LOCAL

export type ConversationVerb =
  | { kind: "open.conversation"; runId: string; campaignId?: string }
  | { kind: "open.turnTrail";    runId: string; turnId: string }
  | { kind: "conversation.rerun"; runId: string; armId: string;
                                  approvalId?: string };             // DANGER
```

| Object | Label | Verb | Notes |
|---|---|---|---|
| `job` | Open (primary) | `open.job` | |
| | Pause / Stop | `job.pause` / `job.stop` | unavailable when not supervised, with the reason |
| | Resume | `job.resume` | **danger** — resumes spending |
| | Start now | `job.start` | **danger**; refused when dependencies unmet (S48) |
| | Add to watchlist | inherited | a long build worth checking on |
| `phase` | Show its units | `open.job` + phase filter | |
| | Inspect | inherited | counters, rates, spend, projection |
| `failureGroup` | Retry this group | `job.retryGroup` | **danger** — spends again |
| | Acknowledge | `failures.acknowledge` | expected failures stop shouting |
| | Open one member | `open.unit` | |
| `unit` | Open in full (primary) | `open.unit` | prompt, output, error, cost |
| | Open its configuration | `inspect` | |
| `turn` | Open this turn's retrieval (primary) | `open.turnTrail` | **Theme 3's trail tile, unchanged** |
| | Ask this turn's question | `ask.run` | Theme 3's verb, reused |
| | Mark as … target | `target.attach` | Theme 2's verbs, reused |
| `run` *(Theme 3)* | Open as conversation | `open.conversation` | when the run has turns |
| | Re-run against… | `conversation.rerun` | **danger** — spends |

The reuse to notice: a `turn` needs **no new explanation surface**. Its menu
routes into Theme 3's trail and Theme 2's target marking, because a turn's
retrieval is a retrieval like any other.

## 5. Tiles

### `jobs` — everything at once (singleton)

```text
┌ ⠿ JOBS ────────────────────────────────────────────────────────┐
│ RUNNING · 3                              spend today  $4.18    │
│ ▌job build ttc-site-2026-08         embed      8,412/14,382    │
│    ▇▇▇▇▇▇░░░░  58% · 41/min · ~2h11m · $3.02 of $12.00         │
│ ▌job judge campaign:51717a…         judging      388/1,200     │
│    ▇▇▇░░░░░░░  32% · 22/min · ~37m · $1.16 of $6.00            │
│ ▌job campaign rrf-sweep             executing     41/60  ✖1    │
│    ▇▇▇▇▇▇▇░░░  68% · 4.2/min · ~4m                             │
│ QUEUED · 2                                                     │
│ ▌job campaign ttc-live-baseline     waits on ▌job build   ⏸    │
│ ▌job conversations ttc-scenarios    waits on ▌job campaign ⏸   │
│ FINISHED · 24h                                                 │
│ ▌job build ttc-site-2026-08-a   ✔ 3h02m · $9.41                │
│ ▌job judge campaign:3f84c0…     ■ stopped by me · 210 done     │
└────────────────────────────────────────────────────────────────┘
widgets  three grouped lists · per-job progress bar · rate · ETA ·
         spend against ceiling · dependency chips
data     GET /api/rag/v1/work  +  SSE /work/stream
```

Serves **S37**, the queue half of **S48**, and is where **S46** lands in the
morning: the finished list with outcomes and costs is the overnight answer.

### `job` — one job, in detail (doc-bound: focus)

The generic frame is the same for every kind — identity, inputs, phases,
failures, throughput, control — and a **kind panel** renders what is
specific. That registry is the same idea as the specialist app's
schema-keyed output widgets (`apps/specialist/web/src/layerwidgets/`), and
it is what keeps this from becoming four bespoke tiles.

```text
┌ ⠿ JOB · ptr-71bd ──────────────────────────────────────────────┐
│ ▌job build ttc-site-2026-08 · running · started 09:14 (1h47m)  │
│ INPUTS  corpus sha256:9ac2f1… · chunker markdown 512/64        │
│         summarizer ▌gpt-… · embed ▌text-embedding-…            │
│ PHASES                                                         │
│  ✔ ▌chunk          14,382/14,382    3m11s               free   │
│  ✔ ▌represent      14,382/14,382   52m04s   $2.71 · 2.1M tok   │
│  ▶ ▌embed           8,412/14,382   41/min   $3.02 · 1.3M tok   │
│      ▇▇▇▇▇▇░░░░ 58%   ~2h11m left    ⚠ projects $5.16 of $6.00 │
│  ○ ▌index               0/1        projected 4m                │
│ FAILURES · 61 in 3 groups                                      │
│  ▌group rate_limit        transient  54  ▲ growing   [retry]   │
│  ▌group context_too_long  data        6    stopped   [retry]   │
│  ▌group provider_5xx      transient   1    stopped             │
│ THROUGHPUT  workers 8 · batch 64     [workers −/+] [batch −/+] │
│ CONTROL  [pause] [stop]              supervised by this UI     │
└────────────────────────────────────────────────────────────────┘
widgets  job identity · input chips · phase rows with bar/rate/
         projection/spend · failure-group rows · throughput
         controls · control bar · kind panel
data     GET /work/{jobId} + /work/{jobId}/failures + SSE
```

Serves **S38**, **S39**, **S40**, **S47**, and the projection warning is
S15's ceiling promise honoured at runtime rather than only at planning time.

### the judge kind panel — reading verdicts as they arrive

```text
┌ ⠿ JOB · ptr-9c02 ──────────────────────────────────────────────┐
│ ▌job judge campaign:51717a… · running · 388/1,200 · $1.16      │
│ EPOCH  ttc.judge.retrieval-support/v1 · ▌model … · prompt d3f1…│
│ SCORES FORMING       0 ▁▂▅▇▅▂▁ 1        median 0.62  n=388     │
│ VERDICTS · newest first                                        │
│  ▌unit q-overwinter × rrf-30            0.90                   │
│     "evidence covers mulch depth and timing; the zone claim    │
│      is implied rather than stated"                            │
│  ▌unit q-comparison × rrf-30            0.20                   │
│     "retrieved size and colour, nothing about shade tolerance" │
│  ▌unit q-inventory × rrf-30             0.10                   │
│     "no stock data in the corpus — question may be out of      │
│      scope"                          ⚠ 3rd verdict like this   │
│ CONTROL  [stop]   on restart: [keep judged ▾ | discard]        │
└────────────────────────────────────────────────────────────────┘
```

Serves **S42** whole. The score histogram forming live is what makes a judge
stuck at one value obvious within a dozen calls, and the explicit
keep-or-discard choice on restart is the "my choice, stated explicitly"
the story asks for.

### `unit` — one item of work, in full (doc-bound: focus)

```text
┌ ⠿ UNIT · ptr-33f8 ─────────────────────────────────────────────┐
│ ▌unit represent chunk-8f2… · attempt 2 of 6 · 1.9s · 412 tok   │
│ IN  ▌job build ttc-site-2026-08 · phase ▌represent             │
│ INPUT · prompt summary/v3 (digest d3f1…)                       │
│  "Summarize the following passage for retrieval. Keep plant    │
│   names and hardiness zones verbatim…"                         │
│  ───                                                           │
│  "Blue Ice grows best in morning sun with afternoon shade…"    │
│ OUTPUT                                                         │
│  "A blue-flowered panicle hydrangea for partial shade in…"     │
│ ATTEMPTS                                                       │
│  1  ✖ rate_limit (transient) after 0.4s → retried              │
│  2  ✔ 1.9s                                                     │
└────────────────────────────────────────────────────────────────┘
```

Serves **S41**. Restricted content follows the existing sensitivity policy:
a prompt the principal may not read renders as the same metadata-only
fallback the tiles already use, never as absence.

### `conversation` — the multi-turn result (doc-bound: focus)

```text
┌ ⠿ CONVERSATION · ptr-4a71 ─────────────────────────────────────┐
│ ▌run conv-8ce2… · scenario ▌q-shade-followup · ▌arm rrf-30     │
│ VERDICT  conversation 0.35   ✖ scenario rubric not met         │
│ TURNS · 4                                                      │
│  1 ▌turn "which hydrangeas take shade?"        0.90 ✔  6 hits  │
│  2 ▌turn "how big does that one get?"          0.85 ✔  4 hits  │
│  3 ▌turn "and in zone 5?"                      0.20 ✖  5 hits ◀│
│      first degradation · retrieval returned zone-7 guidance    │
│  4 ▌turn "so can I plant it in November?"      0.15 ✖  3 hits  │
│      carries turn 3's zone error into the answer               │
│ [open turn 3's retrieval]   [re-run against…]           [pin]  │
└────────────────────────────────────────────────────────────────┘
widgets  run identity · conversation verdict · turn rows with
         per-turn verdict and hit count · first-degradation mark
data     GET /campaigns/{c}/conversations/{runId}
```

Serves **S43** and **S44**. The first-degradation mark is computed, not
judged: it is the first turn whose observation status or score falls below
the conversation's threshold, and it is labelled as a heuristic pointer
rather than a claim about causation. **S45** is the re-run verb plus the
same tile opened twice, side by side.

## 6. Workspaces

### Operations — "is anything wrong right now?" (S37–S41, S46, S47)

```text
┌──────────────────────────┬──────────────────────────────┐
│ JOBS                     │ JOB     phases · failures ·  │
│ everything, all kinds    │         throughput · control │
│                          ├──────────────────────────────┤
│                          │ UNIT    one item in full     │
├──────────────────────────┴──────────────────────────────┤
│ TRACE    every control action, with who took it         │
└─────────────────────────────────────────────────────────┘
```

The trace tile earns its place here: pausing a build and retrying a failure
group are actions with consequences, and they belong in the same record as
every other verb, attributed to whoever performed them.

### Judging live — "is this judge doing what I asked?" (S42)

```text
┌──────────────────────────────┬────────────────────────┐
│ JOB · judge panel            │ QUESTION               │
│ verdicts streaming in        │ the rubric being       │
│ score distribution forming   │ judged against         │
├──────────────────────────────┼────────────────────────┤
│ UNIT   one judge call, full  │ JOBS   spend so far    │
└──────────────────────────────┴────────────────────────┘
```

Reading a verdict beside the rubric it was supposed to apply is the whole
point of this arrangement, and the question tile is Theme 2's, unchanged.

### Conversations — "where did it go wrong?" (S43–S45)

```text
┌──────────────────────────┬─────────────────────────────┐
│ CONVERSATION   turns     │ TRAIL   turn 3's retrieval  │
│ verdicts · degradation   │        (Theme 3's tile)     │
├──────────────────────────┼─────────────────────────────┤
│ JOBS   the suite's job   │ CHUNK   what it found       │
└──────────────────────────┴─────────────────────────────┘
```

Two of these four tiles come from Theme 3 without modification.

## 7. State

Nothing here is a workbench document. Job state belongs to the engines and
is read through projections; the tiles bind to `focus` pointer documents
naming a `jobId`, `unitId`, or `runId`, exactly as the evidence tiles bind
to campaigns and episodes today.

Two consequences worth stating:

- **The live stream is the freshness mechanism**, not polling. The workbench
  already consumes SSE for document revisions; the work stream is a second
  subscription with the same lifecycle, and a dropped connection degrades to
  periodic refetch with the staleness visible in the tile rather than a
  frozen number pretending to be live.
- **Control is authorization, not UI state.** `supervisedJobIds` and
  `canControlJobs` decide availability; a job this process did not start and
  a principal without the grant both produce an unavailable action with a
  stated reason. Dead buttons are the failure mode to avoid — an agent
  planning against the vocabulary needs the reason as much as the human.

## 8. Backend changes

### 8.1 The work projection (new)

One vocabulary, two adapters. Neither engine changes shape.

| Endpoint | Returns |
|---|---|
| `GET /api/rag/v1/work` | every job: kind, state, current phase, progress, rate, spend against ceiling, dependencies |
| `GET /api/rag/v1/work/{jobId}` | inputs, ordered phases with counters and projections, throughput settings |
| `GET /api/rag/v1/work/{jobId}/units?phase=&status=` | paged units |
| `GET /api/rag/v1/work/{jobId}/units/{unitId}` | full inputs (prompt), output or error, attempts, cost — sensitivity-gated |
| `GET /api/rag/v1/work/{jobId}/failures` | groups by cause and class, with counts and growth |
| `GET /api/rag/v1/work/stream` | SSE: phase advanced, counters changed, failure group grew, job state changed |

**Engine A adapter** (campaigns): a campaign is a job whose phases are plan,
execute, and score; units are episodes from `work_items`; failure groups
come from terminal failures with their error class; spend comes from the
budget ledger. Everything needed is already in the journal — this adapter is
pure projection and can be built today.

**Engine B adapter** (builds and flow-based experiments): phases are flow
steps, and `flow.StepReport` supplies `Items`, `WorkCalls`, `Retries`,
`RetriesByClass`, `Quarantined`, `Skipped`, `Spend` and `Meters` — which is
almost the entire phase row and the whole failure-group table.

### 8.2 The one behavioural change: make engine B observable while running

`flow.Report` is produced at the end. The minimal change that fixes this
without touching flowkit: **the run directory publishes progress
periodically.** A `progress.json` in the `runstore` run directory, written
every few seconds with the in-flight `flow.Report` plus the current step,
and the server watches run directories under a configured root.

- It works for every existing build and experiment command with a small
  change at the call site, not in the engine.
- It survives the process dying: the last snapshot is on disk with a
  timestamp, so a stalled job is visibly stale rather than silently frozen.
- The better long-term fix is a progress callback on `flow.Policy`, which
  belongs upstream in flowkit. Propose it there; do not block on it.

### 8.3 Control operations (new)

Through the command API with a new `job.control` authorizer action.
Resume, start, retry-group and conversation re-run are **danger verbs**:
they spend money, so the agent may request and a human approves, using the
approval flow already built.

Honest limits, which the UI must state rather than hide:

- Engine A pauses and resumes naturally — leases expire and `campaign
  resume` continues; this is proven behaviour.
- Engine B needs a stop signal. v1: a sentinel the runner polls between
  items, honoured at item boundaries so nothing is interrupted mid-call.
- A job **started from a terminal** is visible but not controllable. Its
  control actions are unavailable with that reason. This is the difference
  between an honest tool and one with buttons that do nothing.

### 8.4 Multi-turn evaluation

Three additions, none of them to the measurement chain:

1. **A conversation executor** — `rag.ttc-conversation/v1`, sibling of
   OPTKIT-025's G1 live executor. It runs a scenario's turns against the
   assistant, each turn performing its own retrieval and tool calls, and
   emits trajectory events with **one span per turn**, retrieval events
   nested beneath their turn's span. `episode.Event` already carries `Span`
   and `Parent`; this is a convention over an existing structure.
2. **Scenario cases** — a scenario is a question with follow-ups. Theme 2's
   `ragttc.questionset-draft/v1` gains `turns: [{text, rubric?}]` on a
   question, and a scenario-level rubric is the question's existing rubric
   field. Authoring, gold labels, targets, versioning and sealing all work
   unchanged.
3. **Turn-level observations** — the instrument seam (OPTKIT-025 G3) scores
   per turn with `SubjectRef{Kind: "turn", ID: "<episode>/<span>"}` and once
   per conversation with `Kind: "episode"`. `SubjectRef` is already generic;
   nothing in `measure` changes.

Plus one projection: `GET /campaigns/{c}/conversations/{runId}` returning
turns in order with their text, verdicts, hit counts, and the retrieval run
behind each — which is what lets Theme 3's trail tile open a turn.

### 8.5 Frontend registration

Four tiles in `createWorkbenchApps()`, five descriptors and type-graph
nodes, the contributions of §4, the kind-panel registry, three workspaces,
and the SSE subscription beside the existing document one. The vocabulary
export regenerates and pins the control verbs, which is how the agent learns
it can ask to pause a build.

## 9. What this design deliberately does not build

- **No third execution engine.** The temptation is a unified job runner. It
  would be a rewrite of two working systems to serve a read surface, and the
  projection gets the same result for a fraction of the cost.
- **No alerting or notifications.** No email, no push, no webhooks. S46 is
  served by a summary that is there when you come back — which is a
  projection, not a delivery system. Notifications can be added later
  against the same stream if the need survives contact with reality.
- **No scheduling.** Jobs start because a person or a dependency starts
  them. Cron-like recurrence is a Theme 7 concern (re-evaluate after a
  corpus change) and should be designed there, with the comparison
  guarantees in view.
- **No log viewer.** The unit is the unit of inspection. A log pane invites
  reading noise instead of the record, and the record is better.
- **No per-kind tiles.** Four kind panels inside one job tile, not four
  job tiles. The frame — phases, failures, throughput, control — is genuinely
  common; only the detail differs.
- **No automatic remediation.** Retry policy and error classification belong
  to flowkit and are already there. The UI shows the classification and lets
  a human retry a group; it does not invent its own recovery logic on top.
- **No conversation authoring UI beyond follow-up turns.** Branching
  scenarios, personas, and simulated users are a much larger design. A
  scenario is an ordered list of turns.

## 10. Dependencies and honest sequencing

| Piece | Needs | Notes |
|---|---|---|
| Work projection, engine A adapter | nothing | campaigns are journaled today; `jobs` and `job` tiles work immediately |
| `unit` tile for episodes | nothing | reads existing episode records |
| Control for engine A | nothing | pause/resume is proven CLI behaviour |
| Engine B progress snapshots | small change at build call sites | unlocks the story the user actually asked about — watching a big embedding run |
| Engine B control | sentinel + runner polling | v1 stops at item boundaries |
| Judge kind panel | OPTKIT-025 G4 (the judge instrument) | the panel is small; the judge is the work |
| Conversation executor and tiles | G1 live executor first | a conversation is many retrievals |

Suggested order: (1) the work projection with the campaign adapter plus the
`jobs`, `job` and `unit` tiles — self-contained, works on today's data, and
makes campaigns watchable; (2) engine B progress snapshots and its adapter,
which is what makes a corpus build visible; (3) control, starting with
engine A where it is already safe; (4) the judge panel when the judge lands;
(5) conversations, after G1.

## 11. Open questions

- **Where do jobs live between processes?** The projection needs to know
  about runs the serving process did not start. Watching a run-directory
  root is proposed; a small job registry table in the store is the
  alternative. The choice affects whether a build started on another machine
  can appear at all.
- **Should a campaign and its judge batch be one job or two?** They are one
  intention and two engines. Proposed: two jobs with an explicit dependency,
  because that is what S48 asks the system to enforce and it keeps each
  job's phases honest.
- **Turn-level epochs.** If a turn observation and a conversation
  observation use the same instrument, are they the same epoch with
  different subject kinds, or different epochs? Proposed: same epoch,
  different subject kind — but this touches comparability (S20) and deserves
  the measurement owner's opinion before it is built.
- **What is a conversation's ground truth?** A per-turn rubric plus an
  overall one is proposed, but whether turn rubrics are worth authoring by
  hand — as opposed to only judging the conversation and using turn scores
  diagnostically — is an open question that the first real scenario suite
  will answer better than this document can.
