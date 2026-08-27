---
Title: Handoff
Ticket: OPTKIT-029
Status: complete
Topics:
    - backend
    - design
    - rag-ttc
    - ui
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://optkit/ttmp/2026/08/27/OPTKIT-029--reading-results-the-experiments-projection-the-results-tile-and-findings/design-doc/01-intern-guide-reading-results-end-to-end.md
      Note: Read this first — §7 is the exact command sequence, and it needs no provider
    - Path: ws://optkit/ttmp/2026/08/27/OPTKIT-029--reading-results-the-experiments-projection-the-results-tile-and-findings/reference/01-diary.md
      Note: Seven steps with the failures; the last section scopes OPTKIT-030
    - Path: ws://optkit/ttmp/2026/08/27/OPTKIT-027--explore-workbench-ad-hoc-search-corpus-browsing-and-chunking-preview/reference/02-handoff.md
      Note: The ticket this continues — the explore surface and the summaries result
ExternalSources: []
Summary: Two-paragraph orientation — what to run first, the two things that are true about the data, and the fork that is still open.
WhatFor: Getting the next person productive on the results surface in half an hour.
WhenToUse: Read before anything else on day one.
---

# Handoff

Run it before reading about it, and note that unlike OPTKIT-027 you need
**nothing** to start: `rag-ttc index results --artifact
.cache/rag-ttc/eval-raw-vs-summary.json` reads a measurement off disk with no
bundle, no embedder, no profile and no spend, and `index compare --summary
--baseline <bundle>:rrf --challenger <bundle>:rrf` prints the number this
whole ticket exists to make sayable. §7 of the [intern
guide](../design-doc/01-intern-guide-reading-results-end-to-end.md) has the
full sequence including serving the workbench. The organising idea is the
**leg** — one (bundle, strategy) pair, wire id `<bundle_id>:<strategy>` —
because that is the only thing that has metrics, and because comparing bundles
without holding the strategy fixed is precisely how OPTKIT-027 nearly shipped
the wrong conclusion: summaries helped bm25 (+0.015 nDCG@10), hurt vector
(−0.006), and did nothing at the RRF fusion that actually serves users. The
code to read, in order: `pkg/ttc/evalartifact/evalartifact.go` (the schema,
with one asymmetry documented at length), `pkg/ttc/experimentsapi/project.go`
(the projection — `Summarize`, `Legs`, `Compare`, all pure), then
`apps/workbench/web/src/apps/ResultsApp.tsx` and `FindingApp.tsx`. The CLI and
the tiles call the **same functions**, so any number that looks wrong in a
tile can be checked in a terminal without a server.

Two facts are load-bearing and you should not re-derive them. **The mean is
not the result**: on the committed artifact the rrf change moves nDCG@10 by
+0.0012 — a rounding error — while 29 of 144 questions move, four of them down
by more than a third of a point, and the count of perfectly-answered questions
falls from 91 to 90 *as the mean rises*. That is why nothing on this surface
ever shows a mean alone, why histogram bins are fixed over [0,1] rather than
fitted per leg, and why the comparison tile leads with counts. **And
`hit_rate_at` omits its zeros**: upstream writes `HitRateAt[cutoff] = 1` only
on a hit and averages the absence back as 0, so for that metric — and only
that metric — a missing cutoff means *false*, not *unknown*. Applying the
program-wide "missing is never zero" rule uniformly drops exactly the misses
from the denominator and makes hit_rate@1 read 1.0000 on every leg, against
true values of 0.7778 to 0.8889. `evalartifact.PerQuery` encodes the real
rule and `TestPerQueryMeansMatchTheReport` pins it against the artifact's own
aggregates; treat that test as load-bearing rather than decorative. The fork
you should still not settle alone is the one OPTKIT-027 left open: **whether
`index evaluate` becomes a proper campaign (OPTKIT-025 G3) or stays a
standalone runner.** This ticket was scoped to reading so it stays open — a
campaign-backed runner publishing the same artifact would serve these routes
unchanged — but the next piece is probably *running* an evaluation from the
workbench, and that piece cannot avoid the question. Three smaller follow-ups,
in the order they will start to hurt, are in the last section of the diary: a
findings list (one finding is reachable, six are not), the agent path
(`finding.create` and `finding.cite` are in the agent vocabulary and nothing
exercises them), and cross-experiment reading (deliberately unsupported —
two runs against different evaluation sets are not comparable — but "did last
week's rebuild help?" is a real question with no surface).
