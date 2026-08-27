---
Title: 'Intern Guide: Reading Results End to End'
Ticket: OPTKIT-029
Status: active
Topics:
    - backend
    - design
    - rag-ttc
    - ui
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://rag-ttc/pkg/ttc/evalartifact/evalartifact.go
      Note: The artifact schema, written by index evaluate and read by everything here
    - Path: ws://rag-ttc/pkg/ttc/experimentsapi/project.go
      Note: The projection — legs, distributions, and the per-query comparison
    - Path: ws://rag-ttc/cmd/rag-ttc/cmds/indexes/results.go
      Note: The CLI seam; check any number a tile shows without a browser
    - Path: ws://rag-ttc/apps/workbench/web/src/apps/ResultsApp.tsx
      Note: Strategy as an axis; the mean never appears alone
    - Path: ws://rag-ttc/apps/workbench/web/src/apps/FindingApp.tsx
      Note: The claim re-reads its own comparison, which is what makes it falsifiable
ExternalSources: []
Summary: How the evaluation artifact becomes something you can read — the projection, the two CLI commands, the three tiles — and the four decisions that shaped it.
WhatFor: Getting a colleague productive on the results surface, and telling them what NOT to re-derive.
WhenToUse: Before touching anything under pkg/ttc/experimentsapi or the Results workspace.
---

# Intern Guide: Reading Results End to End

## 0. Read this first

`rag-ttc index evaluate` has always been able to measure. Until this ticket
nothing could **read** what it measured except a person squinting at a JSON
file, and the one number that decided OPTKIT-027 — "18 better, 11 worse, 115
unchanged" — was computed with a throwaway script that no longer exists.

This ticket turns that into a projection, a CLI, and three tiles. It reads;
it does not measure. Nothing here runs an evaluation, spends a token, or
writes a file, and that restraint is deliberate: whether `index evaluate`
should become a proper campaign (OPTKIT-025 G3) or stay a standalone runner
is an open question, and every line of this surface works either way.

Start by running §7. It takes two minutes and no provider.

## 1. The one idea: a LEG

An evaluation artifact holds bundles; each bundle was measured under several
retrieval strategies. **A leg is one (bundle, strategy) pair**, and it is the
only thing that has metrics.

Everything on this surface is organised around that, because the alternative
is the mistake OPTKIT-027 nearly shipped. Adding LLM summaries to the corpus:

| strategy | raw | raw+summary | verdict |
|---|---|---|---|
| bm25 | 0.7378 | 0.7531 | helped |
| vector | 0.8599 | 0.8536 | hurt |
| rrf (what serves users) | 0.8623 | 0.8635 | noise |

Three answers to one question. A surface that let you compare "the raw bundle"
against "the summary bundle" without naming the strategy would let any of the
three be quoted as *the* result, and the first one anybody looks at is
usually the one that agrees with them.

A leg's wire id is `<bundle_id>:<strategy>`.

## 2. The second idea: the mean is not the result

The rrf row above moves by **+0.0012**. That is a rounding error, and it is
also the number an operator would write in a report.

Here is the same change read per query:

```
18 better · 11 worse · 115 unchanged      mean moved +0.0012
29 of 144 questions moved (20%)
worst: ttc-y-003  1.0000 → 0.6131  (−0.3869)
```

Twenty percent of the evaluation set moved, four questions lost a
first-rank hit outright, and the count of perfectly-answered questions went
**down** from 91 to 90 — while the mean went up. The mean and the
distribution disagree, and which one you read decides whether you ship.

So: **no surface in this ticket ever shows a mean alone.** Every leg carries
its distribution (ten fixed bins over [0,1], the perfect and zero counts,
p25/p50/p75/p95); every comparison leads with counts and prints the mean
difference beside them. That is the entire design brief.

Two details that follow from it:

- **Bins are fixed, not fitted.** Ten equal intervals over [0,1] for every
  leg. Fitting bins per leg would rescale the picture and make identical
  distributions look different, which is the only thing a side-by-side
  histogram is for.
- **`perfect` and `zero` are reported separately.** Retrieval metrics are
  bimodal: a mean of 0.86 is usually "most questions perfect, a handful
  catastrophic", and which of the two it is changes what you do next.

## 3. The third idea: a finding can be wrong

The last piece is a **finding**: somebody's reading of a result, as an object
rather than as prose in a report. It carries a claim, a status anybody can
move (`proposed` → `accepted` / `disputed` / `withdrawn`), and evidence.

The mechanism that makes it worth having is a negative one. **A finding
stores no numbers.** Evidence is reference keys; the host validator refuses
any body key it does not recognise, so `mean_delta` cannot be written into
one. When the tile renders it re-reads the comparison the finding names and
prints *today's* counts under the claim.

That is what "disputable" means here. After a rebuild, the sentence
"summaries are noise at the fusion" sits directly above numbers that either
still say so or no longer do. A finding that carried its own supporting
numbers would agree with itself forever, and a claim that cannot be
contradicted is a slogan.

Two rules enforce the same idea in smaller ways:

- **A dispute states its reason** — refused in the parser, in the verb, and
  in the Go validator. A status that says "someone disagreed" and cannot say
  why is a shrug wearing a label.
- **Disputing is not a menu row.** It carries prose, and a verb that needs
  prose is authored in a tile that can ask for it. A "Dispute" menu item that
  filed an empty reason would be exactly the shrug the validator refuses.

## 4. What exists

### 4.1 Go

```text
pkg/ttc/evalartifact/         the artifact schema — ONE definition, two sides
    Artifact / Bundle / Strategy   what `index evaluate` writes
    Load                           parse + refuse a foreign schema_version
    MetricIDs / Aggregate / PerQuery   the metric vocabulary

pkg/ttc/experimentsapi/       the projection + HTTP
    project.go   Summarize, Legs, Compare — pure, no clock, no I/O
    experimentsapi.go   Server, Load, Register, the refusal taxonomy

pkg/ttc/apienvelope/          the one {error:{code,message}} envelope
```

`evalartifact` exists because the schema now has a writer AND a reader. A
schema with a private writer and a reader-side copy drifts silently: the
reader keeps parsing, it just stops meaning the same thing.

### 4.2 CLI

```text
rag-ttc index results --artifact PATH
    [--metric ndcg_at@10]... [--all-metrics] [--strategy rrf]...
    one row per leg and metric: mean, graded, perfect, zero, p25..p95, latency

rag-ttc index compare --artifact PATH --baseline LEG --challenger LEG
    [--metric M] [--epsilon E] [--summary] [--all-queries]
    --summary  → ONE row: the counts beside the mean difference
    otherwise  → one row per question that MOVED, worst first
```

Both live in the `index` group, beside the command that wrote the artifact.
Neither resolves a provider or spends anything. **They call the same functions
the tiles do**, so any number that looks wrong in a tile can be checked in a
terminal without a server, a browser, or a token. Use this.

### 4.3 HTTP

| Route | Method | Spends |
|---|---|---|
| `/api/rag/v1/experiments` | GET | nothing |
| `/api/rag/v1/experiments/{experiment}` | GET | nothing |
| `/api/rag/v1/experiments/{experiment}/compare` | GET | nothing |

`compare` takes `?baseline=&challenger=&metric=&epsilon=`. It is a GET
because it is a read: deterministic, free, and cacheable in principle.

Refusal codes, each naming a different problem:

```text
no_experiment_loaded   the server was started without --evaluation-artifact
experiment_not_found   a name was given and no artifact has it
leg_not_found          the experiment has no such (bundle, strategy)
metric_not_found       no leg records that metric
invalid_metric         the metric id does not parse
same_leg               baseline and challenger are the same leg
```

`same_leg` is worth explaining. Comparing a leg with itself returns 144
unchanged questions — a result shape that looks exactly like a measured tie.
It is refused because a plausible-looking nothing is worse than an error:
somebody will quote it.

An experiment is named by its **file name** — what the operator typed — so a
URL and a shell history line say the same thing. Two artifacts resolving to
one name is refused at startup.

### 4.4 Serve

```text
--evaluation-artifact PATH   REPEATABLE. The first is the default.
```

Independent of `--index-bundle`: reading what an evaluation measured needs no
open index, no embedder, and no tool configuration. A bad path fails at
startup rather than at first request.

### 4.5 The workbench

A fourth workspace, **Results**, with three tiles:

```text
┌ ⠿ RESULTS ─────────────────────────┐ ┌ ⠿ COMPARISON · ptr-… ───────────┐
│ [eval-raw-vs-summary ▾][ndcg@10 ▾] │ │ NDCG_AT@10 · 144 QUESTIONS      │
│ questions  144 evaluated · 4 skip  │ │ 18 better  11 worse  115 same   │
│ BM25 · NDCG_AT@10                  │ │ the mean moved +0.0012 —        │
│  ▌raw         0.7378 ▁▁▂▃█         │ │   29 of 144 moved (20%)         │
│    76 perfect · 16 zero · p25 .613 │ │ baseline    rrf · raw · 91 perf │
│  ▌raw+summary 0.7531 ▁▁▂▃█         │ │ challenger  rrf · r+s  · 90 perf│
│    77 perfect · 14 zero · p25 .613 │ │ [write a finding about this]    │
│ VECTOR · NDCG_AT@10   …            │ │ MOVED · 29, WORST FIRST         │
│ RRF · NDCG_AT@10      …            │ │  ▌ttc-y-003 −0.3869 1.00→0.613  │
└────────────────────────────────────┘ │    unit:75df… ⟶ unit:75df…      │
                                       └─────────────────────────────────┘
┌ ⠿ FINDING · fnd-… ─────────────────────────────────────────────────────┐
│ FINDING · PROPOSED   ▌eval-raw-vs-summary                              │
│ [Summaries are noise at the fusion: the rrf mean rose 0.0012 while    ]│
│ [the perfect-query count fell from 91 to 90…                         ]│
│ [accept] [why you disagree              ] [dispute]                    │
│ WHAT IT READS, NOW                                                     │
│   counts      18 better · 11 worse · 115 unchanged                     │
│   mean moved  0.00120                                                  │
│   read again just now — not the numbers this claim was written         │
│   against, but the ones it has to survive                              │
│ EVIDENCE · 1                                                           │
│   ▌ttc-y-003: cited   [queryOutcome:eval-raw-vs-summary/…]             │
└────────────────────────────────────────────────────────────────────────┘
```

Three new presentation types (`experiment`, `leg`, `queryOutcome`) plus
`finding`. Two new navigation verbs and eight finding verbs. `compare.with`
— the accept-mode "pick the other one" gesture that already anchored arm
comparisons — now anchors legs too.

The type graph gained an abstract **`citable`** node, with `watchable`
beneath it: everything worth tracking over time is worth citing, which
correctly leaves stages, layers and variables out. They are configuration,
not evidence.

## 5. Invariants (in addition to OPTKIT-027 §9)

- **A mean never appears without its distribution.** Anywhere. If you add a
  view that shows one, you have removed the reason this ticket exists.
- **Fixed bins.** Ten over [0,1], the same for every leg, or two histograms
  stop being comparable.
- **A finding stores no numbers.** The Go validator enforces it by refusing
  unknown body keys; keep it that way.
- **A dispute has a reason.** Three layers check this; do not remove one
  because the other two exist.
- **The counts partition the evaluation set.** better + worse + unchanged +
  ungraded equals what the run evaluated. A test pins it.
- **Ungraded is not zero, and not unchanged.** A question one leg graded and
  the other did not is a fact about the run, counted and never averaged.

## 6. The one asymmetry you must know about

`hit_rate_at` **omits** a cutoff whose per-query value is zero. Upstream
(`ragkit/rag/evaluation`) writes `HitRateAt[cutoff] = 1` only on a hit, and
its own averaging reads the missing key back as 0. Every other metric writes
its zeros.

So for that metric — and only that metric — an absent cutoff means *false*,
not *unknown*. Reading it as ungraded removes exactly the misses from the
denominator and leaves only the hits behind, so **hit_rate@1 reads 1.0000 on
every leg** against true values between 0.7778 and 0.8889 — a metric that can
never report anything but perfection. Reading another metric's absence as zero
makes the opposite error and invents observations nobody made.

`evalartifact.PerQuery` encodes exactly that rule, and
`TestPerQueryMeansMatchTheReport` checks every metric on every leg against
the artifact's own aggregate to float precision. If you touch that function,
that test is the thing that tells you.

## 7. Running it end to end

```bash
# 1. read an artifact from a terminal — no server, no provider, no spend
./rag-ttc index results --artifact .cache/rag-ttc/eval-raw-vs-summary.json \
  --output-fields leg,strategy,representations,metric,mean,perfect,zero,p25,p95

# 2. the number that decides things
./rag-ttc index compare --artifact .cache/rag-ttc/eval-raw-vs-summary.json \
  --summary \
  --baseline  rk-65790ee26943caa0ba6b2ec361a0a874:rrf \
  --challenger rk-7e257c3a8503d23edf2331659fc13025:rrf \
  --output-fields metric,better,worse,unchanged,ungraded,mean_delta
#   ndcg_at@10 | 18 | 11 | 115 | 0 | 0.0012003097

# 3. the tail, worst first
./rag-ttc index compare --artifact .cache/rag-ttc/eval-raw-vs-summary.json \
  --baseline rk-65790ee2…:rrf --challenger rk-7e257c3a…:rrf \
  --output-fields query_id,outcome,baseline,challenger,delta

# 4. serve it. NOTHING ELSE IS REQUIRED — no bundle, no tool config, no profile.
./rag-ttc experiment optkit-rag campaign serve \
  --store "$PWD/.cache/rag-ttc/serve-store" --listen 127.0.0.1:8544 \
  --workbench-token results-token --workbench-actor actor:operator \
  --workbench-docs-store "$PWD/.cache/rag-ttc/serve-docs" \
  --evaluation-artifact "$PWD/.cache/rag-ttc/eval-raw-vs-summary.json"

cd apps/workbench/web && SPECIALIST_API=http://127.0.0.1:8544 pnpm dev
# Results workspace → right-click a leg → "Compare with…" → click the other
# leg → "write a finding about this" → type a claim → right-click a moved
# question → "Cite as evidence".
```

Add `--index-bundle` and its friends (OPTKIT-027 §10) if you also want the
Explore and Material workspaces in the same server; the two surfaces do not
depend on each other.

## 8. Tests worth knowing

```text
pkg/ttc/experimentsapi/experimentsapi_test.go
    runs against the REAL OPTKIT-027 artifact, committed gzipped (37 KB).
    Reproduces 18/11/115 on rrf, 8/1/135 on bm25, 7/14/123 on vector.
    TestMeanHidesTheDistribution asserts the premise of the whole ticket.

pkg/ttc/workbenchhost/workbenchhost_test.go
    a disputed finding with no reason is refused; a finding carrying a
    number is refused by key name.

apps/workbench/web/src/test/menu-goldens.test.ts
    the exact menu rows for leg, queryOutcome, experiment.
apps/workbench/web/src/test/findings.test.ts
    dispute-needs-a-reason, half-a-comparison, and evidence keys that
    round-trip into openable chips.
```

The fixture is the real artifact rather than a hand-built one on purpose: the
claim this projection has to support is a specific pair of numbers from a
specific run, and a synthetic fixture only proves the code agrees with itself.

## 9. What this ticket did NOT do

- **No campaign.** `index evaluate` is still a standalone runner. Whether it
  becomes a campaign (OPTKIT-025 G3) is deliberately still open, and this
  surface would not notice either way — a campaign-backed runner publishing
  the same artifact serves these routes unchanged.
- **No writing.** Nothing here starts an evaluation. The button that would do
  that belongs to whoever settles the fork above.
- **No significance testing.** `epsilon` is float noise, not a threshold. The
  projection reports what moved; deciding whether it *matters* is the
  reader's job, and it is what a finding is for.
- **No cross-experiment comparison.** Two artifacts can be served at once and
  read separately; comparing legs across them is not supported, because two
  runs against different evaluation sets are not comparable and the surface
  should not imply they are.

## 10. Glossary

| Term | Meaning here |
|---|---|
| **experiment** | one `index evaluate` artifact, named by its file name |
| **leg** | one (bundle, strategy) pair — the unit that has metrics |
| **metric id** | `mrr` or `<name>@<cutoff>`, e.g. `ndcg_at@10` |
| **graded / ungraded** | whether a question carries a value for a metric |
| **moved** | a question whose metric changed by more than epsilon |
| **finding** | a claim about a result, with a status and evidence keys |
| **citable** | a type a finding may point at (not configuration) |
