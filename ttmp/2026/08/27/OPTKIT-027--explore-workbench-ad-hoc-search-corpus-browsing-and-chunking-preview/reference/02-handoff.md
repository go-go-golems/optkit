---
Title: Handoff
Ticket: OPTKIT-027
Status: complete
Topics:
    - design
    - ui
    - rag-ttc
    - backend
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/ttmp/2026/08/27/OPTKIT-027--explore-workbench-ad-hoc-search-corpus-browsing-and-chunking-preview/design-doc/01-intern-guide-the-explore-workbench-end-to-end.md
      Note: Read this first — §10 has the exact commands that work
    - Path: repo://optkit/ttmp/2026/08/27/OPTKIT-026--operator-user-stories-for-rag-design-and-optimization/analysis/01-user-stories-the-person-who-makes-the-rag-system-good.md
      Note: The 48 stories the whole programme is derived from
    - Path: repo://rag-ttc/datasets/ttc/README.md
      Note: The corpus and the 148-query evaluation set, already in the repo
ExternalSources: []
Summary: Two-paragraph orientation for the colleague picking this up — what exists, what to run first, what is true about the corpus, and what the next ticket is.
WhatFor: Getting a new person productive on the explore workbench in an hour.
WhenToUse: Read before anything else on day one.
---

# Handoff

Start by running the system rather than reading about it: §10 of this
ticket's [intern guide](../design-doc/01-intern-guide-the-explore-workbench-end-to-end.md)
has the exact command sequence that works, including the two path quirks
that each cost an hour to discover (`index build` needs absolute
`--output-root`; `search` needs a `--bundle` relative to
`--repository-root`, so the absolute path the build *prints* is rejected by
the command that consumes it). Everything you need is already committed:
`datasets/ttc/` holds a 200-, 2,000- and 3,149-document corpus plus an
evaluation set of **148 real queries with 243 graded judgments**, so you do
not have to author ground truth to measure anything. Two bundles are built
under `.cache/rag-ttc/indexes/`: `rk-65790ee2…` (raw) and `rk-7e257c3a…`
(raw + LLM summaries). Serve both at once — `--index-bundle` is repeatable
— and the workbench gives you three workspaces: **Evidence** (campaign
tiles, empty without a campaign, which is fine), **Explore** (ask a
question, "+ lane" to compare two bundles, click a result for its trail
through all twelve pipeline stages), and **Material** (corpus survey,
document detail, chunking split). The code to read, in order:
`internal/customer/ragsearch/ragsearch.go` (the composition every consumer
reuses), `pkg/ttc/exploreapi/` (the projections), then
`apps/workbench/web/src/apps/TrailApp.tsx` and the three contributions in
`src/pbui/actions.ts` — that last file is where every menu item in the
product is declared, and it is the fastest way to understand how the UI
works.

Three findings are load-bearing and you should not re-derive them. **The
chunker is windowing, not structuring**: 1,780 of 1,979 chunks end at
exactly the 1200-rune limit rather than at a document boundary (p25 = p50 =
p75 = p95 = 1200), because TTC content is largely flat prose with almost no
markdown headings — so `--chunk-runes` is the lever, and `markdown-heading`
will not behave differently here. **Summaries did not pay off**: 1,979 LLM
calls bought +0.0197 recall@10 on bm25, −0.0064 nDCG@10 on vector, and at
the RRF fusion that actually serves users a net of 18 better / 11 worse /
115 unchanged across 144 queries — noise (`.cache/rag-ttc/eval-raw-vs-summary.json`,
reproduced by §10 step 5). Note the shape of that mistake, because it is the
main methodological lesson here: the two-lane ask comparison made summaries
look like a clear win on one query, and the full evaluation said otherwise
— **the lane comparison is a hypothesis generator, `index evaluate` is the
thing that answers**. The next ticket is **OPTKIT-029**, scoped in the last
section of the diary: an `experiments` read projection over the evaluation
artifact (the per-query win/loss counts that made today's result legible
were computed with a throwaway script — that is a missing projection), a
`results` tile showing means *and* distribution with strategy as an axis,
and a `finding` document where the agent's reading of a result becomes a
first-class, disputable object with evidence chips rather than prose. The
open fork you should not settle alone: whether `index evaluate` becomes a
proper campaign (OPTKIT-025 G3, the instrument seam) or stays a standalone
runner — 029 is deliberately scoped to *reading* results so that decision
stays open. Also of interest: **OPTKIT-028 dropped in priority today**,
because its value is watching hour-long builds and the run that justified it
is the full-corpus summary build the measurement just argued against.
