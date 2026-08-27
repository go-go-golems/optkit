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
