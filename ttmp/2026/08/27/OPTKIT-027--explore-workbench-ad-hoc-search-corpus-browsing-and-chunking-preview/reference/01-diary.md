---
Title: Diary
Ticket: OPTKIT-027
Status: active
Topics:
    - design
    - ui
    - rag-ttc
    - backend
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Chronological record of the explore workbench build — search command, HTTP endpoint, ask/trail tiles, corpus projections, and the chunking split.
LastUpdated: 2026-08-27T15:10:00-04:00
WhatFor: ""
WhenToUse: ""
---

# Diary

## Goal

Chronological record of turning a built index bundle into something you can
search, browse, and reason about — from the CLI seam through the HTTP
projections to the workbench tiles.

## Step 1: The Search Command (rag-ttc 9e5cc8cb2)

The first real work of the whole program that is not a design document. A
glazed command that opens an index bundle, runs one query through the live
hybrid pipeline, and prints either the returned results or every recorded
stage candidate. It is deliberately the smallest possible consumer of the
retrieval composition, because everything above it — the HTTP endpoint, the
ask tile, and eventually the campaign live executor — assembles the same
handle.

The important discovery of this step was not in the code I wrote but in what
running it revealed: the pipeline has **twelve** stages, not the six the
design documents had been describing, and the corpus and its evaluation set
were already committed in the repository all along.

### Prompt Context

**User prompt (verbatim):** "Ok, where would you start, mostly about just
like, finally doing some real work? Ideally I would like to just chunk and
index some stuff, see the pipelines going, and running some searches. I
don't really need the optimization framework part, I finally want to just
get things churning a bit." … then "Ok, create a new docmgr ticket, Create
a detailed analysis / design / implementation guide that is for a new
intern… Do slice 1-3 in that ticket."

**Assistant interpretation:** Stop designing; build the shortest path to a
searchable corpus. Ticket it, write the intern guide, then implement.

**Inferred user intent:** Something churning that they can point at and use,
without the optimization machinery in the way.

**Commit (code):** rag-ttc 9e5cc8cb2 — "search: run one query against a
built bundle and show every stage (OPTKIT-027 P1)"

### What I did
- `cmd/rag-ttc/cmds/search/command.go`: a glazed command over
  `ragsearch.Open` → `Handle.NewSessionRegistry` → `SearchTool.RunRoute`,
  with `--stages` switching from result rows to per-candidate stage rows.
- Registered it in `cmd/rag-ttc/main.go` as a top-level `rag-ttc search`.
- `command_test.go`: input refusal order, rune-counting truncation, and
  channel naming — the three pure behaviours worth pinning without a bundle.
- Built the real 200-document bundle and ran actual queries against it.
- Corrected the intern guide's stage list from six invented names to the
  twelve real ones, with a note that they were verified by running the tool.

### Why
- Task 9tnl. The command is the composition seam: proving it works proves
  the bundle opens, the embedder resolves, the tool config parses, and the
  hybrid path runs — with no server and no browser in the way.

### What worked
- The retrieval is genuinely good. "where is my order?" returns, in order:
  order-editing, shipping-notification, shipping-method, order-status, and
  order-edit FAQs. Real content, sensible ranking, 200 documents.
- `cache_hits 1979 / provider_work_calls 0` on the build: every embedding
  was already in `.cache/rag-ttc`, so the first real bundle cost nothing.
- Stage mode immediately paid for itself by contradicting the design docs.

### What didn't work
- `--output-root .cache/…` failed with `invalid uri authority: .cache`. The
  content store builds a SQLite URI from the path, and a leading dot is read
  as a URI authority. Absolute paths work. Worth a friendlier error upstream.
- `ragsearch.Open` requires the bundle path be **relative to the repository
  root** (`safePath`), so the absolute path that the build command *emits*
  is rejected by the search command that consumes it. Both are defensible in
  isolation; together they are a papercut, and the guide now shows the
  working invocation.
- The first dry run resolved the default profile (`gpt-5.6-luna-codex`) and
  failed on a missing embedding key. `--profile ttc-live-openai` plus BOTH
  registries — the pinocchio default and the repo's own `profiles.yaml` — is
  the incantation. The binary's baked default registry path points at a
  different workspace entirely.

### What I learned
- **The stage list is twelve, and symmetric.** `lexical.{raw,collapsed,
  policy_filtered}`, `vector.{raw,collapsed,policy_filtered}`,
  `retrieval.{fused,policy_recheck,reranked}`, `evidence.{hydrated,admitted,
  returned}`. Policy is applied twice — per channel and again after fusion —
  so "why did this disappear?" has two distinct answers and the stage name
  is what distinguishes them. The trail tile design should show the two
  channels side by side rather than as one list.
- **The corpus was already here.** `datasets/ttc/` holds 200-, 2,000- and
  3,149-document corpora plus an evaluation set with 148 queries and 243
  graded judgments. Two earlier analyses of mine assumed ingest work that
  does not exist. Checking the repository beat reasoning about it.

### What was tricky to build
- Choosing between `ragsearch.Open` and assembling `indexbundle.Open` +
  `ttcsearch.NewService` directly. The former lives under
  `internal/customer/` and is nominally the customer app's composition root,
  which reads oddly for a generic command. But it carries verified-document
  loading, source-catalog construction, role policy, and route compilation —
  roughly sixty lines of careful sequencing whose ordering matters (bleve
  takes an exclusive lock, so verification must happen first). Duplicating
  that to avoid an awkward import would have been the worse trade.

### What warrants a second pair of eyes
- The `internal/customer/ragsearch` import from a top-level command. If this
  composition is now shared by the CLI, the HTTP endpoint, and later the
  campaign executor, it probably wants to move out of `customer/`.
- One session registry per query. Correct — the evidence ledger is
  turn-scoped and must not carry citations across searches — but it does
  construct a tool per call, and a reviewer should confirm that is as cheap
  as it looks over already-open stores.

### What should be done in the future
- P2: the HTTP endpoint over the same composition, opened once at startup.
- The relative-versus-absolute path asymmetry between `index build` and
  `search` deserves fixing in one direction or the other.

### Code review instructions
- Start at `cmd/rag-ttc/cmds/search/command.go`; compare `run()` against
  `internal/customer/ragsearch/ragsearch.go:46` to see what is reused.
- Validate: `go test ./cmd/rag-ttc/cmds/search/`, then run the command
  against a built bundle with and without `--stages`.

### Technical details

```bash
# build (cache-warm: no provider calls)
./rag-ttc index build --corpus datasets/ttc/corpus.json \
  --output-root "$PWD/.cache/rag-ttc/indexes" --cache-directory "$PWD/.cache/rag-ttc" \
  --profile-registries "$HOME/.config/pinocchio/profiles.yaml,$PWD/profiles.yaml" \
  --profile ttc-live-openai --embedding-budget 2500 --max-estimated-usd 1.00 \
  --allow-unpriced-provider
# → rk-65790ee26943caa0ba6b2ec361a0a874 · 200 docs · 1979 chunks · 0 provider calls

# search (bundle path RELATIVE to --repository-root)
./rag-ttc search --bundle ".cache/rag-ttc/indexes/rk-65790ee26943caa0ba6b2ec361a0a874" \
  --tool-config assets/configs/tool-qa/production-v1.yaml \
  --repository-root "$PWD" --scratch-directory ".cache/rag-ttc/scratch" \
  --profile-registries "$HOME/.config/pinocchio/profiles.yaml,$PWD/profiles.yaml" \
  --profile ttc-live-openai --query "where is my order?" --limit 5 [--stages]
```

Observed stage row counts for that query: `lexical.*` 20 each,
`vector.*` 20 each, `retrieval.*` 25 each, `evidence.hydrated` 25,
`evidence.admitted` 1, `evidence.returned` 5.

## Step 2: The Search Endpoint, and a Boundary That Caught Me (rag-ttc e96f687bc)

The same composition behind HTTP, opened once at startup because the lexical
index takes an exclusive lock. The response embeds the retrieval output
verbatim rather than projecting it — the per-candidate scores, channels,
matched representation ids and fusion contributions are exactly what an
explanation surface needs, and every projection in this program that has
coarsened them has had to be un-coarsened later.

The interesting part of this step was a test I did not write. The product
boundary test rejected my first design outright, and it was right to.

### Prompt Context

**User prompt (verbatim):** (see Step 1 — "Do slice 1-3 in that ticket")

**Assistant interpretation:** P2 — the HTTP search endpoint and the serve
flags that open a bundle.

**Commit (code):** rag-ttc e96f687bc

### What I did
- `pkg/ttc/exploreapi`: `Server{Retriever, Facts}` with `GET
  /api/rag/v1/explore` (availability plus what is loaded) and `POST
  /api/rag/v1/search`. `SearchResponse` embeds `ttcsearch.SearchOutput`.
- `serve` gained `--index-bundle`, `--tool-config`, `--repository-root`,
  `--scratch-directory`, a profile section, and construction through
  `BuildCobraCommandWithGeppettoMiddlewares` so profile resolution works.
- Nine tests: availability both ways, five typed refusals, and one that
  asserts a membership-only stage carries no score **on its own bytes**.
- Live: served the 200-document bundle and queried it over HTTP.

### Why
- Task 9ond. The ask tile needs a route; the route needs the handle open
  once; and the trail tile needs the candidates the campaign projection drops.

### What worked
- Embedding `SearchOutput` rather than projecting it. Twelve stages arrive
  intact, `retrieval.fused` candidates carry
  `contributions:[{channel,rank,weight,value}]`, and `evidence.admitted`
  arrives with no `candidates` key at all — membership, not zeros.
- The availability route earns its place immediately: without it, a server
  started with no bundle would answer searches with an empty list, which
  reads as "nothing matched" rather than "nothing was configured".

### What didn't work
- **`TestProductBoundaries` failed my design**: `pkg/ttc/exploreapi/handle.go
  imports internal/customer/ragsearch: product boundary forbids this
  dependency`. I had put the handle adapter beside the server. The fix was to
  invert the dependency — the server declares the narrow `Retriever`
  interface it needs, and the adapter moved to
  `cmd/rag-ttc/cmds/experiments/optkitrag/explore.go`, at the composition
  site where importing the customer package is legitimate. The resulting
  design is better than the one I wrote, and the boundary test is why.
  This is exactly the risk I flagged for review in Step 1; the repository
  answered it before a human had to.
- Three test call sites of `composeAPIHandler` needed the new parameter,
  found one at a time because the compiler stops at the first.

### What I learned
- Retrieval on a 200-document corpus is honest about its own limits: "how do
  I plant a magnolia?" returns *How to Plant Hydrangeas* twice and a spring
  blooming guide. There is no magnolia planting content in this slice, and
  the pipeline surfaces the nearest planting material rather than nothing.
  That is a corpus finding, not a retrieval bug — and precisely the kind of
  thing the coverage story (S2) exists to make visible.
- `evidence.returned` (out 3) comes BEFORE `evidence.admitted` (out 3, no
  candidates) in the recorded order. Worth knowing before the trail tile
  renders stages in array order and implies a causality that is not there.

### What was tricky to build
- `serve` was built with the plain cobra builder, which carries no profile
  middleware, so `ResolveCLIEngineSettings` had nothing to read. Switching to
  `BuildCobraCommandWithGeppettoMiddlewares` and making `newServeCommand`
  fallible was the smallest correct change; the alternative — resolving
  settings outside the command — would have put provider wiring in main.

### What warrants a second pair of eyes
- `serve` now constructs through the geppetto middleware chain. Existing
  flags are unchanged and its tests pass, but the flag surface grew by the
  whole profile section, which a reviewer should look at once.
- One session registry per HTTP request. Correct for evidence isolation;
  confirm the allocation cost is as small as it looks under load.

### What should be done in the future
- P3: the ask and trail tiles over this response.
- The explore routes are unauthenticated, matching the read tier. When the
  corpus contains restricted sources this needs the same sensitivity
  treatment the campaign projections have.

### Code review instructions
- `pkg/ttc/exploreapi/exploreapi.go` then
  `cmd/rag-ttc/cmds/experiments/optkitrag/explore.go` — the interface and its
  one implementation, deliberately apart.
- Validate: `go test ./pkg/ttc/exploreapi/ ./cmd/rag-ttc/...`, then serve a
  bundle and `curl -s localhost:PORT/api/rag/v1/explore`.

### Technical details

```bash
./rag-ttc experiment optkit-rag campaign serve \
  --store "$PWD/.cache/rag-ttc/serve-store" --listen 127.0.0.1:8531 \
  --workbench-token explore-token --workbench-actor actor:operator \
  --workbench-docs-store "$PWD/.cache/rag-ttc/serve-docs" \
  --index-bundle ".cache/rag-ttc/indexes/rk-65790ee26943caa0ba6b2ec361a0a874" \
  --tool-config assets/configs/tool-qa/production-v1.yaml \
  --repository-root "$PWD" --scratch-directory ".cache/rag-ttc/scratch" \
  --profile-registries "$HOME/.config/pinocchio/profiles.yaml,$PWD/profiles.yaml" \
  --profile ttc-live-openai
```

## Step 3: The Explore Workspace (rag-ttc 49131106a)

Two tiles and a second workspace. You type a question, you get real results
from the real bundle, and clicking one shows its path through all twelve
stages with the fusion contributions that decided its rank. The whole slice
touches no campaign, no journal, and no measurement.

Three bugs surfaced only in the browser, and each was a fact about the
system rather than a typo. The most useful was the last one.

### Prompt Context

**User prompt (verbatim):** (see Step 1 — "Do slice 1-3 in that ticket"),
then "go ahead."

**Assistant interpretation:** P3 — the ask and trail tiles over the P2
endpoint, plus the workspace that holds them.

**Commit (code):** rag-ttc 49131106a

### What I did
- `run` and `hit` presentation types, descriptors, type-graph nodes, and
  three action contributions; `ask.run` / `ask.clear` / `open.trail` verbs.
- `api/explore.ts` (RTK Query), `store/runs.ts` (bounded scratch-run slice),
  `AskApp`, `TrailApp`, explore styles.
- `ragttc.explore-focus/v1`: a pointer with no campaign, validated in Go,
  because exploration has none and the focus format requires one.
- Go catalog entries for `ask` and `trail`; the workbench document became
  two workspaces, Evidence and Explore.
- 8 new tests (55 total green), vocabulary regenerated, live-verified in the
  browser end to end.

### Why
- Task pgsn. The endpoint existed; this is the surface that makes it usable
  without curl.

### What worked
- The kernel did what it promised. `target.attach`-style reuse was not even
  needed: declaring `hit` in the type graph gave it Inspect for free, and
  the disabled "Open the chunk" row renders its own reason —
  *"a scratch run has no episode chunk catalog"* — visibly, in the menu.
- The trail immediately taught something true. For "when should I prune a
  panicle hydrangea?", the chunk containing *"Panicle Hydrangea (but not
  mop-heads)"* was **bm25 rank 15** (0.0887) but **vector rank 5** (0.5397),
  and at fusion the vector channel contributed more (0.01538 vs 0.01333),
  carrying it to rank 4. It was found by meaning, not by words — which is
  exactly the question the tile exists to answer.

### What didn't work
Three browser-only failures, in the order they appeared:

1. **`id_mismatch` on every sync push.** `layout(spec, {id})` had been
   setting the DOCUMENT id; moving to `workspaces([...])` without a second
   argument let the library mint a fresh one, so the sync pushed
   `wb-8ad156f6…` to the path `ragttc-workbench` and was refused. Fixed by
   passing `{ id: "ragttc-workbench" }` explicitly.
2. **`unknown_application`.** The running server predated the catalog
   change. Not a code bug — a reminder that the Go catalog and the React app
   registry are two halves of one contract and both have to be rebuilt.
3. **`duplicate_singleton`.** I put `trace` in both workspaces. **A
   singleton is unique per DOCUMENT, not per workspace**, so the whole
   document was unrepresentable. The Explore workspace is now `ask` (a
   singleton used nowhere else) beside `trail` (doc-bound).

Diagnosing (3) took longer than it should have because the sync logged
`create refused (422)` with no body while the push path logged the full
reason. That asymmetry is now fixed — the create path logs the validator's
message, which named the duplicate instantly.

### What I learned
- Singleton scope is document-wide. Any workspace beyond the first can only
  contain doc-bound tiles plus singletons no other workspace claims. This
  constrains every future workspace and belongs in the workspace design
  guidance, not in one diary entry.
- The `trail` binding must be `Required: false`. A required binding makes a
  tile unrepresentable in a default layout, because the layout is created
  before any pointer exists. The tile already had an honest empty state;
  the catalog just had to agree that the state is legal.

### What was tricky to build
- Deciding where scratch runs live. A workbench document would have made
  them agent-visible and reload-proof, but a `SearchOutput` carries the full
  text of every chunk any stage touched, and the document host has no
  sensitivity model for content. They are a bounded in-memory slice (cap 8,
  evicting from both the order and the map) and the trace records that the
  search happened. Losing them on reload is the correct behaviour.

### What warrants a second pair of eyes
- `verbFamily` classifies `ask.run` as a command so the sink reports its
  outcome like an API call, while `open.trail` is navigation by prefix. The
  classification drives trace semantics, so it is worth one reviewer's look.
- The explore routes are unauthenticated (read tier). Fine for this corpus;
  wrong the moment restricted sources are indexed.

### What should be done in the future
- P4 corpus/document tiles, P5 the chunking split.
- The stage order shows `evidence.returned` before `evidence.admitted`,
  which reads as a causality that is not there. Worth confirming with the
  retrieval owner whether the recorded order is meaningful.

### Code review instructions
- Start at `src/apps/TrailApp.tsx` (the transpose of the autopsy), then
  `src/store/runs.ts` and the three contributions in `src/pbui/actions.ts`.
- Validate: `pnpm test && pnpm exec tsc --noEmit`, then serve a bundle, open
  the Explore workspace, ask something, and click a result.

## Step 4: The Material Workspace — Corpus, Document, Split (rag-ttc fec8316c2)

P4 and P5 together, because they share a pointer, a workspace, and most of a
backend. Three projections over `indexbundle.Inspect`/`Measure`, one pure-CPU
chunking preview, three tiles, and a third workspace.

The tiles were built to answer questions about the corpus. Within about two
minutes of pointing them at the real 200-document bundle they answered two I
had not asked, and both are findings about the RAG system rather than about
the code.

### Prompt Context

**User prompt (verbatim):** "Do P4 + P5"

**Assistant interpretation:** Corpus survey and document detail (P4), then
the chunking split (P5), backend and tiles.

**Commit (code):** rag-ttc fec8316c2

### What I did
- `exploreapi/corpus.go`: survey, document list (sort/filter/page), document
  detail with chunks and source text; `exploreapi/split.go`: the chunking
  preview over `chunking.Apply`.
- `serve` inspects the bundle once at startup and loads the source corpus
  separately, warning rather than failing when it is missing.
- `corpus`, `document`, `split` tiles; `corpus`/`document`/`chunkPreview`
  presentation types; `open.corpus` / `open.document` / `open.split` verbs;
  Go catalog entries; the Material workspace.
- 20 Go tests, 59 frontend tests, all green; vocabulary regenerated.

### What worked
- Wrapping `indexbundle` was as cheap as the guide predicted. The CLI's
  `inspect corpus` commands and these three routes are now two thin skins
  over the same two library calls.
- The histogram made a corpus fact unmissable at a glance (below).

### What I learned — two findings about the corpus, not the code

**1. The chunker is windowing, not structuring.** The survey reports
`at-limit 1780 chunks / 161 docs` out of 1979 chunks, and p25 = p50 = p75 =
p95 = **1200**. Over 90% of chunks end because the size limit ended them, not
because the document did. The histogram is one bar: 1780 chunks in the
1190–1200 bucket. The TTC content is largely flat prose without markdown
headings, so `markdown` at 1200 runes degenerates into a fixed window. That
is worth knowing before anyone tunes fusion.

**2. A large fraction of some documents is directory noise.** Opening
"Arkansas Trees For Sale" (16 chunks, 16,456 runes) shows chunks #7 through
#14 are nursery names, addresses and phone numbers, and #14 is a bare list of
about a hundred phone numbers with no surrounding text. Half that document is
being chunked, embedded and indexed as retrievable evidence. This is the
"junk in my corpus" story arriving as observation rather than as a hypothesis.

### What didn't work
- **My own heading heuristic was wrong on real data.** `documentHasHeading`
  tested `strings.Contains(text, "#")`, and TTC product copy is full of
  "#7 Gallon" and "sku: #708559" — so the first live run marked **every**
  chunk of a 127-chunk document as orphaned. A heuristic that fires on every
  document tells a person nothing. `firstHeading` now requires one to six
  hashes followed by whitespace, and there is a test built from the real
  strings that fooled it.
- **A shadowed JSON field, caught across the language boundary.**
  `DocumentResponse` embedded `DocumentSummary` (which has a `chunks` COUNT)
  and also declared a `chunks` ARRAY. Go silently lets the outer field win,
  so the count vanished from the wire. The TypeScript mirror refused to
  compile — `Type 'ChunkView[]' is not assignable to type 'number'` — which
  is the only reason I noticed. The response now nests the summary, and both
  numbers survive: the count comes from the manifest and the array from the
  bundle, so a disagreement between them is exactly the corruption worth
  seeing.
- One of my own test cases was wrong (six hashes plus a space IS a valid h6).
  Fixed the test, not the code.
- The Material workspace silently failed to appear because an earlier
  explanatory comment had broken my edit anchor and the replacement never
  applied. Verifying in the browser is what caught it; the typechecker
  could not.

### What was tricky to build
- Keeping "document runes" and "sum of chunk runes" apart everywhere. They
  differ by exactly the overlap, and merging them would misreport every
  windowed document — which, per finding 1, is almost all of them here. They
  are separate fields with separate names in the Go type, the TS type, and
  the tile's key/value rows.

### What warrants a second pair of eyes
- The split preview loads the whole source corpus into memory at startup
  (19MB for the full corpus). Fine at this scale; a streaming read or an
  on-demand seek would be better before anyone points it at something large.
- The corpus routes are unauthenticated like the rest of the read tier, but
  they serve full document text. That is the first place restricted content
  would leak if this corpus ever had any.

### What should be done in the future
- P6: the smoke pass and the guide's command block updated with what
  actually works.
- The at-limit finding deserves a real experiment: `markdown-heading` and
  smaller windows against the 148-query evaluation set. That is exactly what
  the optimization machinery is for, and now there is a reason to run it.

### Code review instructions
- `pkg/ttc/exploreapi/corpus.go` then `split.go`; the tiles in
  `src/apps/CorpusApp.tsx`, `DocumentApp.tsx`, `SplitApp.tsx`.
- Validate: `go test ./pkg/ttc/exploreapi/`, `pnpm test`, then open the
  Material workspace and choose "Compare chunking…" on any document.
