---
Title: Diary
Ticket: OPTKIT-029
Status: active
Topics:
    - backend
    - design
    - rag-ttc
    - ui
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://optkit/ttmp/2026/08/27/OPTKIT-029--reading-results-the-experiments-projection-the-results-tile-and-findings/scripts/01-per-query-win-loss.py
      Note: Step 1 — the throwaway script written down as the projection's target
    - Path: ws://rag-ttc/apps/workbench/web/src/apps/FindingApp.tsx
      Note: Step 6 — the tile that re-reads its own comparison
    - Path: ws://rag-ttc/apps/workbench/web/src/pbui/actions.ts
      Note: Steps 5 and 6 — the results contributions and the citable node
    - Path: ws://rag-ttc/pkg/ttc/evalartifact/evalartifact.go
      Note: Steps 2 and 6 — the shared schema and the hit_rate asymmetry
    - Path: ws://rag-ttc/pkg/ttc/experimentsapi/project.go
      Note: Step 2 — legs, distributions, the per-query comparison
    - Path: ws://rag-ttc/pkg/ttc/experimentsapi/experimentsapi_test.go
      Note: Runs against the real OPTKIT-027 artifact and reproduces all three handoff numbers
ExternalSources: []
Summary: Seven steps from "the 18/11/115 came from a script nobody kept" to a projection, two CLI commands, three tiles, and a claim that can be contradicted by its own evidence.
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: Reviewing OPTKIT-029, and knowing which decisions were load-bearing.
WhenToUse: Before changing anything under pkg/ttc/experimentsapi or the Results workspace.
---


# Diary

## Goal

OPTKIT-027 ended with a real result — LLM summaries do not pay for themselves —
and no way to read it. The per-query counts that made it legible were computed
with a throwaway script. This ticket makes reading a measurement a first-class
part of the product: a projection over the artifact, a CLI over the same
functions, tiles that never show a mean alone, and a finding that can be
wrong.

## Step 1: Ground Truth Before Code (optkit 6f6a26e)

The handoff's headline is a specific triple — 18 better, 11 worse, 115
unchanged, on rrf — computed by a script that no longer exists. Writing the Go
projection first and then checking it against a number in a paragraph would
have meant checking it against a memory. So the first commit is the script,
rebuilt and checked in, with an `--expect` flag that fails loudly.

It reproduced all three of the handoff's claims exactly: bm25 recall@10
+0.0197, vector nDCG@10 −0.0063 (the handoff rounds to −0.0064), and rrf
18/11/115 over 144 queries. That agreement is what let me trust my reading of
the artifact schema for everything after.

### Prompt Context

**User prompt (verbatim):** "The handoff is at ttmp/2026/08/27/OPTKIT-027…/reference/02-handoff.md: […] The next ticket is **OPTKIT-029**, scoped in the last section of the diary: an `experiments` read projection over the evaluation artifact (the per-query win/loss counts that made today's result legible were computed with a throwaway script — that is a missing projection), a `results` tile showing means *and* distribution with strategy as an axis, and a `finding` document where the agent's reading of a result becomes a first-class, disputable object with evidence chips rather than prose. […] commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill). Print out a brutalist work slip with the plan / different phases for the ticket. then before stsarting a phase, plrint a split about the phase, and print one when the phase is done."

**Assistant interpretation:** Pick up OPTKIT-029 from the OPTKIT-027 handoff
and build the three pieces it names, committing as I go, keeping a diary in
the skill's format, and printing thermal work slips for the plan and for each
phase boundary.

**Inferred user intent:** A working results surface, plus the paper trail that
lets the next person continue without re-deriving anything — and a physical
record of the phases as they happen.

**Commit (code):** optkit 6f6a26e

### What I did
- Created the docmgr ticket.
- Wrote `scripts/01-per-query-win-loss.py`: reads the artifact, computes
  per-query deltas between two legs under one metric, prints the counts and
  the moved queries, and exits non-zero if `--expect better/worse/unchanged`
  does not match.
- Ran it against `.cache/rag-ttc/eval-raw-vs-summary.json` for all three
  strategies and cross-checked the aggregate deltas.

### Why
- A projection needs an oracle. Reconstructing the oracle first means the Go
  code has something to be wrong against.

### What worked
- `--expect 18/11/115` passed on the first run, which meant my reading of
  `report.per_query` was right.
- Checking bm25 and vector too surfaced the shape of the real finding: 8/1/135
  and 7/14/123 — summaries helped one channel and hurt another.

### What I learned
- The handoff said 029 was "scoped in the last section of the diary". It was
  not; the diary ends at Step 6 (typography) and the scope lives only in the
  handoff paragraph itself. Not a problem — the paragraph is specific — but
  worth recording, because I spent a few minutes looking for a section that
  does not exist.

### What was tricky to build
- Nothing yet. The one care taken: `metric_value` returns `None` rather than
  `0.0` for a missing cutoff, which mattered a great deal two steps later.

### What warrants a second pair of eyes
- N/A — the script is a check, not a dependency.

### What should be done in the future
- N/A; superseded by Step 2, which the script now guards.

### Code review instructions
- `scripts/01-per-query-win-loss.py`.
- Validate: `python3 scripts/01-per-query-win-loss.py <artifact> --expect 18/11/115`.

## Step 2: The Projection, and One Asymmetry That Nearly Ate 137 Queries (rag-ttc b9e99b314)

The artifact had a writer and no reader: `evaluationArtifact` was a private
struct in `cmd/rag-ttc/cmds/indexes/evaluate.go`. I moved the schema into
`pkg/ttc/evalartifact` and made both the command and the new projection depend
on it, so a schema change now breaks the reader at compile time rather than at
render time.

The projection itself (`pkg/ttc/experimentsapi`) is organised around the LEG —
one (bundle, strategy) pair — because that is the unit that has metrics, and
because comparing bundles without holding the strategy fixed is exactly how
OPTKIT-027 nearly shipped the wrong conclusion.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Build the read projection the handoff asked for.

**Inferred user intent:** The win/loss counts stop being a script and become
something the product knows how to compute.

**Commit (code):** rag-ttc b9e99b314

### What I did
- `pkg/ttc/evalartifact`: the artifact types, `Load` with schema refusal, and
  the metric vocabulary (`MetricID`/`ParseMetricID`/`MetricIDs`/`Aggregate`/
  `PerQuery`). Cutoffs are DISCOVERED from the report, not hardcoded, so a run
  at k=20 cannot silently report k=10 numbers.
- `pkg/ttc/experimentsapi`: `Summarize`, `Legs`, `Compare`, and an HTTP server
  with a six-code refusal taxonomy.
- `pkg/ttc/apienvelope`: pulled the `{error:{code,message}}` envelope out of
  `exploreapi` so both surfaces share one definition.
- 13 tests against the real OPTKIT-027 artifact, committed gzipped (37 KB
  against 760 KB).

### Why
- The handoff's "that is a missing projection" is precise: the counts existed,
  nothing computed them. Everything else in the ticket reads this.

### What worked
- The tests reproduce all three strategies' counts exactly.
- `TestCountsPartitionTheEvaluationSet` caught a real bug while I wrote it: my
  first `Compare` built the query universe from the metric tables rather than
  from the reports, so a question ungraded on both legs would have vanished
  from the counts entirely without the totals ever failing to add up
  *visibly*. Fixed by taking the universe from `report.PerQuery`.

### What I learned — the asymmetry
`hit_rate_at` **omits** a cutoff whose per-query value is zero. Upstream
writes `HitRateAt[cutoff] = 1` only on a hit, and its own averaging reads the
missing key back as `0` via Go's zero-value map read. Every other metric
writes its zeros.

I found this by checking, not by reading: I computed, for every metric on
every leg, the mean over per-query values two ways and compared both against
the report's own aggregate.

```
hit_rate_at@1: 137 of 144 per-query records have no "1" key
hit_rate_at@5:  43        "
hit_rate_at@10: 38        "
every other metric: 0 missing keys anywhere
```

So the honest rule is *not* the obvious one. "Missing is never zero" is a
program-wide invariant here, and applying it uniformly would have reported
hit_rate@1 as 7 graded questions out of 144. The rule that is actually true:

- `hit_rate_at`: absent cutoff means **false**.
- everything else: absent cutoff means **unmeasured**, and stays absent.

`evalartifact.PerQuery` encodes exactly that, with the reasoning in the
docstring, and `TestPerQueryMeansMatchTheReport` pins it — every metric, every
leg, agreeing with the artifact's own aggregate to 1e-9.

### What was tricky to build
- The above. The symptom would have been silent: a distribution that looked
  plausible, a mean that was wrong by a factor of twenty, and no error
  anywhere. The approach that found it was refusing to trust either reading
  and checking both against a third number the artifact already carried.
- Second, smaller: `Compare`'s counts must partition the evaluation set, which
  means the ungraded case has to be counted rather than skipped, and the
  moved-list must be built from a universe that includes queries neither leg
  graded.

### What warrants a second pair of eyes
- `evalartifact.PerQuery`'s asymmetry. It is correct for this upstream
  version, and it is the kind of thing that silently becomes wrong if ragkit
  ever starts writing hit_rate's zeros. The test would catch that, which is
  the point — but a reviewer should know the test is load-bearing rather than
  decorative.
- The fixed [0,1] bins in `buckets()`. They assume every metric is in [0,1],
  which is true of all five today. A metric outside that range would clamp
  into the end bins rather than fail.

### What should be done in the future
- If a metric outside [0,1] is ever added, `buckets` needs a range, and every
  histogram in the tile needs to stop being comparable across legs — or the
  metric needs excluding from the distribution view.

### Code review instructions
- Start at `pkg/ttc/evalartifact/evalartifact.go`'s `PerQuery` docstring, then
  `project.go`'s `Compare`.
- Validate: `go test ./pkg/ttc/experimentsapi/...` — it runs against the real
  artifact and reproduces the handoff numbers.

### Technical details

```go
// The rule, in full:
value, ok := table[cutoff]
if !ok && name == "hit_rate_at" {
    return 0, true   // absent means false, for this metric only
}
return value, ok     // absent means unmeasured, everywhere else
```

## Step 3: The CLI Seam (rag-ttc 09df2d277)

`rag-ttc search` proved the retrieval composition without a browser in
OPTKIT-027; `index results` and `index compare` do the same for this
projection. They live in the `index` group beside the command that wrote the
artifact, so the sequence reads as one thing: `index evaluate --artifact X`,
then `index results --artifact X`.

They call the same functions the tiles call. That is the whole value: a number
that looks wrong in a tile can be checked in a terminal, with no server, no
token, and no browser.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Give the projection a terminal seam before
wiring any UI to it.

**Commit (code):** rag-ttc 09df2d277

### What I did
- `cmd/rag-ttc/cmds/indexes/results.go`: two glazed commands, no provider
  section (reading a measurement resolves nothing and spends nothing).
- `--summary` on compare emits ONE row — the counts beside the mean
  difference — mirroring `search --stages`.

### Why
- Row output gives table/json/yaml for free, and a per-query comparison is
  genuinely tabular.

### What worked
- The first run printed the ticket's finding without my having to look for it:

```
strategy  representations  mean     perfect  zero
rrf       raw              0.8623   91       2
rrf       raw+summary      0.8635   90       2
```

  The mean went up and the count of perfect answers went down. That one pair
  of rows is a better argument than the paragraph I would have written.

### What didn't work
- `--output table` is not a flag; glazed spells it `--format`. Minor, but it
  cost a round trip:
  `Error: unknown flag: --output`.

### What I learned
- Putting the read commands in the `index` group rather than inventing a
  `results` top-level group made the whole surface discoverable from where the
  operator already is. There is a `cmds/experiments/` package already, holding
  the optkit-rag experiment app — reusing that name for this would have
  collided conceptually with the API's meaning of "experiment".

### What was tricky to build
- `--all-queries` has to reconstruct the unchanged questions locally, because
  the projection deliberately does not ship 115 rows that all say "nothing".
  The reconstruction is exact (delta zero by definition) but it is a second
  place that knows the shape, so it is commented as such.

### What warrants a second pair of eyes
- `chosenMetrics`: an unknown `--metric` errors with the artifact's actual
  metric list attached. Check the error text is genuinely useful.

### What should be done in the future
- N/A.

### Code review instructions
- `cmd/rag-ttc/cmds/indexes/results.go`.
- Validate: §7 of the intern guide, steps 1–3.

## Step 4: Serving It (rag-ttc 2a84264e5)

One repeatable flag, three routes, and a deliberate independence: the
experiments routes need no index bundle, no embedder, and no tool
configuration.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Commit (code):** rag-ttc 2a84264e5

### What I did
- `--evaluation-artifact` on `campaign serve`, repeatable, first is the
  default; artifacts are loaded and validated at startup.
- Mounted `experimentsapi.Server` beside `exploreapi.Server` on the same mux.

### Why
- A bad path failing at startup is diagnosable; a server that starts and then
  refuses one route is not.

### What worked
- Verified by starting a server with an artifact and **nothing else** — no
  bundle, no profile — and curling all three routes. 18/11/115 came back over
  HTTP, and every refusal named what existed instead.

### What was tricky to build
- Nothing structural. `composeAPIHandler` gained a parameter, which broke
  three call sites in two test files — the compiler found all of them, which
  is the argument for threading it as a parameter rather than reaching for a
  package-level variable.

### What warrants a second pair of eyes
- Route ordering on the mux: the experiments patterns are more specific than
  the `/api/rag/v1/` prefix handler, so Go's mux routes them correctly. The
  same reasoning the explore routes already rely on, but worth a second look.

### What should be done in the future
- N/A.

### Code review instructions
- `cmd/rag-ttc/cmds/experiments/optkitrag/serve.go`.
- Validate: start with only `--evaluation-artifact` and curl `/api/rag/v1/experiments`.

## Step 5: The Tiles, and a Vocabulary That Caught Me (rag-ttc f0496a1dc)

Two tiles. `results` groups legs BY STRATEGY and shows every group at once;
`comparison` leads with counts and prints the mean difference beside them.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Build the results tile the handoff asked for —
means AND distribution, strategy as an axis.

**Commit (code):** rag-ttc f0496a1dc

### What I did
- `src/api/experiments.ts` (RTK Query), three presentation types with
  descriptors, two navigation verbs, menu contributions, `ResultsApp`,
  `ComparisonApp`, a new `Results` workspace, and CSS for the leg row, the
  sparkline, and the verdict counts.
- Extended `compare.with` — the existing accept-mode "pick the other one"
  gesture — to anchor on a leg as well as an arm.
- Put the METRIC in the pointer document, not in tile state.

### Why
- The metric decision is the one worth defending. The same two legs under
  recall@10 are 8/1/135 rather than 18/11/115. A comparison read under a
  different metric is a *different comparison*, and a finding that cites one
  must be able to reopen exactly it. Keeping the metric in component state
  would have made the pointer ambiguous the moment anyone shared it.

### What worked
- The type graph did the work I hoped it would: `leg` and `queryOutcome` are
  watchable, `experiment` is not, and the menus fell out of that with no
  per-type special cases.
- The first render of the results tile made the ticket's argument visible
  without any prose: three strategy groups, disagreeing.

### What didn't work — and this is the good failure
The chat vocabulary has an exhaustiveness fence — a TypeScript assertion that
every `ProductVerb` kind appears in the zod schema. Adding two verbs failed
the build immediately:

```
src/chat/verbs.ts(154,26): error TS2344: Type '"open.experiment"' is not
assignable to type '"open.cockpit" | … | "trial.run"'
```

and `vocabulary.ts` failed separately for the three new presentation types:

```
Type '{ campaign: string; … }' is missing the following properties from type
'Record<keyof Values, string>': experiment, leg, queryOutcome
```

This is the "one vocabulary" invariant enforcing itself. I could not add a
menu item an agent would not know about, which is exactly the property
OPTKIT-024 was built for. Every new verb and type now has a wire doc, and the
golden artifact was regenerated with `pnpm vocab`.

### What I learned
- The invariant is real, not aspirational. Two independent compile errors,
  both pointing at the same omission, before any test ran.

### What was tricky to build
- Layout: a singleton tile is unique per DOCUMENT, not per workspace
  (OPTKIT-027 learned this the hard way). `results` is this workspace's own
  singleton and `comparison` is doc-bound, so nothing collides.
- Distribution bars: my first instinct was to scale each leg's histogram to
  its own range, which looks better and lies. Fixed bins over [0,1] look
  worse — most bars are near-invisible because retrieval metrics pile up at
  1.0 — and are the only version where two legs are comparable. Kept the ugly
  one, and wrote down why in the CSS.

### What warrants a second pair of eyes
- `groupByStrategy` preserves the artifact's order (bm25, vector, rrf) rather
  than sorting. That is the pipeline's order; sorting alphabetically would put
  the fusion that serves users in the middle of its own inputs.
- The sparkline is ten `<span>`s with percentage heights. It is legible but
  cramped; a reviewer may reasonably want it taller.

### What should be done in the future
- The results tile has no empty-metric state for a leg measured under a metric
  another leg lacks; it prints "not measured under X" and skips. That is
  correct and untested against a real heterogeneous artifact, because none
  exists yet.

### Code review instructions
- `src/apps/ResultsApp.tsx` (the two design decisions are in its docstring),
  then `src/pbui/actions.ts`'s `resultsContributions`.
- Validate: `pnpm vitest run`, then §7 of the intern guide.

## Step 6: A Finding That Can Be Wrong (rag-ttc 134fde59f)

The third piece: a claim about a result, as an object with a status somebody
else can move and evidence as chips.

The design turns on one negative decision. **A finding stores no numbers.**
Evidence is reference keys; the Go validator refuses any body key it does not
recognise, so `mean_delta` cannot be written into a finding at all. The tile
re-reads the comparison the finding names on every render and prints today's
counts under the claim.

That is what makes it disputable rather than published. After a rebuild the
sentence sits above numbers that either still support it or no longer do. A
finding carrying its own supporting numbers would agree with itself forever,
and "disputable" would be decoration.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Build the finding document — the agent's reading
of a result as a first-class, disputable object with evidence chips rather
than prose.

**Commit (code):** rag-ttc 134fde59f

### What I did
- `ragttc.finding/v1`: format, parser, Go validator, eight verbs, a tile, a
  `finding` presentation type, and the citing flow.
- Added an abstract **`citable`** node to the type graph with `watchable`
  beneath it.
- Extended `refFromKey` to reconstruct leg, outcome and experiment keys.

### Why
- Citing needed a home in the type graph. Putting `finding.cite` on
  `inspectable` put a greyed row on all thirty types including layers and
  variables — configuration, which is not evidence. `citable` says which
  types a claim may point at, and `watchable` sitting beneath it encodes
  "everything worth tracking over time is worth citing".

### What worked
- The dispute-needs-a-reason rule is enforced in three places (parser, verb,
  Go validator) and each refusal is specific.
- The derived-state rule turned out to have real teeth here: the validator's
  unknown-key refusal is what makes "stores no numbers" mechanical rather
  than a convention.

### What didn't work — twice, both found by running it
1. A finding created from the experiment chip carries no baseline/challenger,
   so the live re-read — the entire disputability mechanism — was
   **unreachable from the UI**. Every finding I could create through the
   product was one that could not contradict itself. Fixed by adding "write a
   finding about this" to the comparison tile, seeding baseline, challenger
   and metric.
2. Evidence chips rendered the raw key. For a cited outcome that is
   `queryOutcome:eval-raw-vs-summary/rk-65790ee26943caa0ba6b2ec361a0a874:rrf->rk-7e257c3a8503d23edf2331659fc13025:rrf/ndcg_at@10/ttc-y-003`
   — a hundred characters of digest where a label belongs. Now the descriptor
   label with the key in the tray, matching every other identifier here.

Both are in Step 7's commit. Neither would have been found by a test; both
were obvious within ten seconds of using it.

### What I learned
- Deciding what is NOT a danger verb mattered. `finding.accept` and
  `finding.dispute` are reversible judgements. Marking them danger — which I
  did first, by reflex, because they are capitalised in the docs — would have
  put an approval dialog in front of the one gesture this ticket exists to
  make cheap. Only `finding.discard` is danger, because deleting somebody's
  recorded reading is the irreversible act.
- Disputing cannot be a menu row. It carries prose, and the codebase's own
  precedent (`intent.setHypothesis` has no menu row) says a verb needing prose
  is authored in a tile that can ask for it. A "Dispute" item that filed an
  empty reason would be exactly the shrug the validator refuses.

### What was tricky to build
- The `citable` node. My first version put cite on `inspectable`, which added
  a greyed row to every menu in the product and broke six golden tests at
  once. The failure was loud, which was useful: the goldens ARE the
  specification here, so six of them changing meant six menus I had not
  thought about. The fix was not to suppress the row but to say which types
  it belongs on — and the type graph is the mechanism for that, already
  present and already used by `watchable`.
- `refFromKey` for a query outcome: the key has four slash-separated parts and
  none of the four may itself contain a slash. It reconstructs the coordinates
  and deliberately does NOT invent the score — a chip knows which question it
  names, not what that question scores today. There is a test asserting the
  delta comes back `undefined`.

### What warrants a second pair of eyes
- `activeFindingDocId` is a second "active document" pointer beside
  `activeDraftDocId`, and the last-focused-tile-wins hazard is the same one
  drafts already have. With two finding tiles open, "Cite as evidence" lands
  on whichever was focused most recently. Documented in the tile; a reviewer
  may think it deserves stronger signalling.
- The finding is a workbench DOCUMENT, so it syncs to the document host and
  persists per workbench. It is not campaign-scoped and not in the campaign
  store. That is right for now (the artifact it reads is not either) and is
  the thing that would change if `index evaluate` becomes a campaign.

### What should be done in the future
- A findings *list* — right now a finding is reachable only from the tile that
  created it or by reopening its document. Once there are more than three,
  that stops being enough.
- An agent-authored finding is possible today (`finding.create` and
  `finding.cite` are in the agent's vocabulary; accept and dispute are not).
  Nothing exercises that path yet.

### Code review instructions
- `apps/workbench/web/src/apps/FindingApp.tsx` docstring, then
  `pkg/ttc/workbenchhost/documents.go`'s `validateFinding`.
- Validate: `go test ./pkg/ttc/workbenchhost/...` and
  `pnpm vitest run src/test/findings.test.ts`.

### Technical details

```text
ragttc.finding/v1 body — every key, and nothing else:
  experiment  required   the artifact this reads
  claim       required   may be empty (unfinished ≠ malformed)
  status      required   proposed | accepted | disputed | withdrawn
  baseline    optional   \ all-or-nothing with challenger
  challenger  optional   /
  metric      optional   which metric the comparison is read under
  evidence    optional   reference keys, never values
  dispute     optional   REQUIRED when status is disputed
  author      optional
```

## Step 7: Driving It (rag-ttc afc540d17)

Served both index bundles and the evaluation artifact together, opened the
Results workspace in a real browser, and used it.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Prove the whole thing works against the real
artifact rather than against tests.

**Commit (code):** rag-ttc afc540d17

### What I did
- Right-clicked the rrf `raw` leg → "Compare with…" → clicked the rrf
  `raw+summary` leg. The comparison opened with 18/11/115, mean +0.0012, "29
  of 144 questions moved (20%)", and the tail worst-first with its rankings.
- "write a finding about this" → typed a claim → right-clicked `ttc-y-003` →
  "Cite as evidence" → typed a reason → disputed it.
- Verified the stored document: claim and one evidence key, no numeric field.
- Fixed the two problems in Step 6's "what didn't work", plus a cosmetic one:
  finding ids were `finding-…` and the tile title shows the first eight
  characters of a document id, so every finding tile read
  `finding · finding-`. Now `fnd-`.

### What worked
- The tiles agree with the CLI exactly. Same functions, so they should — but
  "should" is not "did".
- The menus matched the goldens verbatim, including "Cite as evidence — no
  finding is open" greying itself with its reason before a finding existed,
  and becoming available the moment one did.

### What I learned
- Every problem in this step was invisible to the test suite and obvious in
  the product. The re-read section being unreachable is the sharpest example:
  every test passed, the feature worked, and no path through the UI could
  produce a finding that used it.

### What was tricky to build
- Nothing; this step was diagnosis. Worth noting the method: I drove the
  product rather than the API, which is why I found a reachability bug rather
  than a correctness bug.

### What warrants a second pair of eyes
- Screenshots of the three states are in the ticket's `various/`.

### What should be done in the future
- Covered in Step 6.

### Code review instructions
- Follow §7 of the intern guide; it is exactly what I did.

## What OPTKIT-030 should pick up

The fork the handoff told me not to settle alone is still open, and this
ticket was scoped to keep it that way: **does `index evaluate` become a proper
campaign (OPTKIT-025 G3), or stay a standalone runner?** Nothing here depends
on the answer — a campaign-backed runner publishing the same artifact serves
these routes unchanged — but the next piece probably does, because the next
piece is running an evaluation from the workbench, and that needs an
instrument seam or it needs a button.

Three smaller things, in order of how soon they will hurt:

1. **A findings list.** A finding is reachable only from the tile that created
   it. That is fine for one and useless for six.
2. **The agent path.** `finding.create` and `finding.cite` are in the agent
   vocabulary and nothing exercises them. The handoff's phrase was "the
   agent's reading of a result"; today a human can write one and an agent
   could, but has not.
3. **Cross-experiment reading.** Two artifacts can be served and read
   separately. Comparing legs across them is deliberately unsupported — two
   runs against different evaluation sets are not comparable — but the
   question "did last week's rebuild help?" is a real one and currently has no
   surface.
