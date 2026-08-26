---
Title: Implementation Diary
Ticket: OPTKIT-009
Status: active
Topics:
    - implementation
    - rag-ttc
    - ui
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://rag-ttc/apps/specialist/web/src/components/IdChip.tsx
      Note: IdTray disclosure demoting identifiers (commit 1e77d16cd)
    - Path: ws://rag-ttc/apps/specialist/web/src/components/graphs.tsx
      Note: 'Tufte data graphics: slope graph, chunk grid, stage flow, meter bars (commit baf6ce93f)'
    - Path: ws://rag-ttc/apps/specialist/web/src/screens/ComparisonScreen.tsx
      Note: Verdict strip, show-identifiers toggles, |delta|-sorted cases (commit 1e77d16cd)
    - Path: ws://rag-ttc/apps/specialist/web/src/styles/system.css
      Note: Accent palette semantics, flat sections, square buttons (commit 1e77d16cd)
    - Path: ws://rag-ttc/assets/configs/experiments/optkit-rag/semantic-limit-v1.yaml
      Note: Manifest rewritten with authored prose for experiment, arms, cases (commit f5e5a7e29)
    - Path: ws://rag-ttc/pkg/ttc/specialistapi/projector.go
      Note: Surfaces descriptions and case query text in cockpit/comparison (commit f5e5a7e29)
ExternalSources: []
Summary: Diary of the humanize pass on the specialist UI — scenarios document, hyperslop accent palette, flat sections, hash demotion, verdict strip.
LastUpdated: 2026-08-25T22:05:00-04:00
WhatFor: Record the UX-direction change and how it was implemented and validated.
WhenToUse: Read before continuing specialist UI product work (side-by-side diff, previews, playground).
---



# Diary

## Goal

Capture the shift of the specialist UI from journal-shaped tables toward a tool a RAG engineer does work in: a scenarios/questions design doc, plus a restyle (flat sections, hyperslop accent colors on white, square buttons) and hash demotion into identifier trays.

## Step 1: Scenarios doc, accent restyle, hash demotion, verdict-first comparison

Wrote the persona-driven scenarios and questions document, then implemented the design consequences that need no backend change: removed the fake window chrome (framed pinstripe panels, offset shadows) in favor of flat rule-under-title sections, adopted the hyperslop.systems accent palette on a white ground, moved identifiers out of the reading line everywhere, added a verdict strip and |Δ|-sorted cases to the comparison screen, and gave the cockpit a plain-language status sentence. All 28 tests still pass; validated live and archived new screenshots.

### Prompt Context

**User prompt (verbatim):** "Ok, make a list of questions I would like to know as I am working on this RAG system, where I want to also see and maybe even try actual work to get a feel, and not just stare at tables. Also, remove the brutalist dropshadow, and you can use color accents (especially font colors), look at the colors on https://hyperslop.systems/ for example. (but stay white background), so that we can create a friendlier UI. showing hashes doesn't really help, those are atmost details that can be viewed in a detail box/sidebar/whatever." — followed mid-turn by: "go through a list of scenarios where you are the user and want to use a nice UI to do your work." and "also no rounded corners on the buttons and no need for that fake window chrome around the sections, but otherwise this is the look i was going for."

**Assistant interpretation:** Produce a scenarios-and-questions design doc from the user's seat (an engineer tuning a RAG system who wants to see and try real work), and restyle the UI: white background, hyperslop accent colors especially in text, no drop shadows, no rounded buttons, no window-chrome framing, hashes demoted to detail disclosures.

**Inferred user intent:** Evolve the contract-faithful v1 into a friendly instrument that answers working questions at a glance, with identifiers available but out of the way; establish a product roadmap grounded in real usage scenarios.

**Commit (code):** 1e77d16cd — "Specialist UI: friendlier flat look with accent colors, hashes demoted"

### What I did

- Created OPTKIT-009 and wrote `design-doc/01-specialist-work-scenarios-and-questions.md`: seven scenarios (morning check, did-my-change-help, case debugging, getting-a-feel, judge skepticism, writing up, planning), fourteen ranked questions, and an answerability map (now / new-view / backend-ask / deferred).
- Pulled the real palette from hyperslop.systems (`green #2db878, purple #805bd7, red #ef4038, yellow #f2ad00`) and defined darkened text-safe variants plus washes in `system.css`. Semantic mapping: green = healthy/measured/reuse, purple = the operator's direct change + links, yellow = transitive recompute/warnings, red = failure/violation. Missing data keeps the uncolored dither swatch (absence of color = absence of data), now as a tidy 8px swatch inside a dashed chip.
- Removed the System-1 window chrome per follow-up feedback: panels are now flat sections (small-caps title over a 2px rule), no offset shadows, square buttons (`border-radius` gone), pinstripe token deleted. IBM Plex Mono added to the mono stack.
- Demoted hashes: new `IdTray` disclosure component; cockpit campaign IDs and arm snapshots, comparison arm cards and graph IDs go into trays; config-diff and invalidation tables hide their identity/digest columns behind a "Show identifiers" toggle; pipeline artifact metadata folds into an "Artifact record" disclosure while the preview/no-preview reason stays on the reading line. Provenance stays the full-detail screen.
- Verdict-first comparison: a purple-barred sentence ("On retrieval.target-coverage, limit-2 is ahead of limit-1: 0.833 → 1.000 (+0.167 over 3 pairs)"), a "You changed retrieval; the other 11 layers are identical" summary above the diff, and paired cases sorted by |Δ| with signed deltas colored green/red. Cockpit opens with "This campaign is completed: 6 of 6 episodes completed; the resource budget held."

### Why

- The v1 screens answered "what is recorded"; the scenarios doc reorders the UI around "what am I trying to learn". Hashes are appendix material (scenario F) — needed for citation, hostile to reading.
- Color went onto meaning, not decoration, so the monochrome missing-data vocabulary survives: dither still marks the absence of data, now the only uncolored material.

### What worked

- All 28 tests stayed green through the restyle with zero test edits — the tests assert semantics (order, words, links), not presentation, and `<details>` content remains queryable in the DOM.
- The mid-turn feedback loop ("no rounded corners, no fake window chrome... otherwise this is the look") landed cleanly because panels/badges are centralized in `system.css` and `Panel.tsx`.

### What didn't work

- WebFetch of hyperslop.systems returned only text (markdown conversion strips CSS); recovered by curling the raw HTML and grepping hex values/custom properties out of the inline `<style>`.
- First draft of the missing-value chip (dither background clipped to the border box with a white content plate) rendered ragged at small sizes — replaced with a dashed gray chip carrying an 8px dither swatch.

### What I learned

- Keeping component tests on ARIA roles and visible words (not classes/styles) makes an aggressive restyle nearly free.
- A verdict sentence composed purely from existing response fields (`metrics[].baseline/treatment/delta/pairs`) delivers most of the "stop making me read tables" value before any backend work.

### What was tricky to build

- Encoding severity and change-kind twice: color alone fails accessibility and monochrome-print, so every colored badge still carries its word, and dither remains the missing-data material. The mapping had to stay consistent across cockpit badges, diff badges, plan badges, stage headers, and diagnostics — it lives in CSS token names (`badge-direct`, `badge-upstream`, `badge-unchanged`, `badge-good`, `badge-error`) rather than per-screen styles.
- Sorting paired cases without breaking pagination semantics: sorting happens within the fetched page only (`CaseRows` sorts a copy), so the opaque cursor contract is untouched; the caption states the sort so the reader doesn't assume server order.

### What warrants a second pair of eyes

- Verdict phrasing assumes higher metric = better ("ahead of"); if a construct is ever loss-like, the sentence will mislead. The projector may eventually need a direction flag per construct.
- Contrast: darkened accent inks were chosen for AA at small sizes (`#1d7f53`, `#6747b5`, `#cf2b24`, `#93690a` on white); worth a pass with a contrast checker on the washes (`badge-upstream` yellow-ink on yellow-wash).

### What should be done in the future

- Next per the scenarios doc: side-by-side pipeline diff view (frontend-only, both routes exist); backend ticket for config-value and case-input previews; later chunk-lab, playground, judge calibration as contracts land.

### Code review instructions

- Start at `src/styles/system.css` (token semantics), then `ComparisonScreen.tsx` (Verdict, IdColsToggle, CaseRows sort), `IdChip.tsx` (IdTray), `Artifact.tsx` (record disclosure vs visible preview reason).
- Validate: `pnpm typecheck && pnpm test && pnpm build`; visually via tmux session `specialist` (API :8091, UI :5197).
- Evidence: `various/screenshots/07-comparison-restyled.png`, `08-cockpit-restyled.png`, `09-pipeline-restyled.png`.

### Technical details

- Palette tokens: accents `--green #2db878 / --purple #805bd7 / --red #ef4038 / --yellow #f2ad00`; text inks `--green-ink #1d7f53 / --purple-ink #6747b5 / --red-ink #cf2b24 / --yellow-ink #93690a`; washes `#e8f7ef / #f1ecfb / #fdecea / #fdf4dc`; grays `--muted #66686e`, `--line-soft #d5d5d0`.
- The scenarios/questions document is the product roadmap: `design-doc/01-specialist-work-scenarios-and-questions.md`.

## Step 2: Data graphics, failure-first reframing, and prose that lives in the data files

Rewrote the scenarios document from the seat of its real user — someone improving a chatbot's retrieval, whose work runs on failures, judge traces, and examples rather than averages — and implemented the two halves that need no new producers: Tufte-principled data graphics, and authored prose carried in the experiment data files themselves, through the campaign spec, out the API, onto every screen.

### Prompt Context

**User prompt (verbatim):** "add sparklines and graphs and such (think about design principles like tufte, in terms of data representation, not style necessarily). Also, for your work, don't base it on what is already there, think as a human who is basically working to create better IR for a chatbot. We want to much more look at judge traces, verdicts, examples of failures, etc... rephrase." — followed mid-turn by: "all this technical scaffolding is to serve human needs, not code for the sake of code", "lso let's use more human readable prose to accompagny our data files and campaign configs and such. They should not just be terse keywords, we have the opportunity to write things out so that someone browsing the side understands what the things are and why.", and "are you expanding the data files and such with human prose, it shouldn't just be general information."

**Assistant interpretation:** Three connected asks: (1) add data graphics guided by Tufte's representation principles; (2) rephrase the product framing around the human doing IR-improvement work — failure examples, judge verdicts and traces at the center; (3) put real authored prose into the experiment data files themselves (manifest, arms, cases) — not generic UI copy — and surface it wherever the data appears.

**Inferred user intent:** A tool a person can browse and understand without insider knowledge: graphics that answer working questions at a glance, and experiments that explain themselves in their author's words.

**Commit (code):** f5e5a7e29 — "Experiment data carries authored prose for the person reading results"; baf6ce93f — "Specialist UI: data graphics and prose that serve the reading human"

### What I did

- Rewrote `design-doc/01-specialist-work-scenarios-and-questions.md` in the user's voice: the work loop (read failures → hypothesize → change → re-run → churn), scenarios like "Show me what's broken", "Says who?" (judge traces), the failure autopsy, and a re-ranked question list with an answerability map.
- Built `src/components/graphs.tsx`, each graphic titled by its question: **SlopeGraph** (which cases moved; the mean is one bold line among its cases so an average cannot hide outliers), **ChunkGrid** (chunk × stage presence matrix — where did the evidence go), **StageFlow** (candidate counts across stages, narrowing segments in red), **MeterBar** (budget headroom, word-sized). Missing values are excluded and named beneath the graphic, never plotted as zero.
- Backend: added `description` fields to `experimentworkbench.Manifest`/`ManifestArm` and `optkitcampaign.RetrievalCase`/`Arm`/`RunOptions`/`CampaignSpec`; wired `service.Run` and `CampaignArms()` to carry them; surfaced them in `specialistapi` (`CampaignCockpit.Description`, `ArmSummary.Description`, `CaseSummary.Query`+`Description` — the case's actual question text was already persisted, just never surfaced).
- Rewrote `assets/configs/experiments/optkit-rag/semantic-limit-v1.yaml` with real authored prose: what the experiment asks, what each arm trades, what each case exercises and what passing means (e.g. the negative case: "retrieving nothing relevant is the correct behavior").
- Added `lede` support to Panel and wrote plain-language explanations for every section on cockpit, comparison, and pipeline; cases table now shows "The question asked" and "What it checks".
- Rebuilt, re-ran the campaign (new campaign `campaign:fa569e901e…`), verified descriptions flow through the live API, screenshots archived (10–12).

### Why

- The v1 contract shows what is recorded; the person needs to know what it *means*. Descriptions authored at experiment-writing time are the only place the "why" exists — the tool's job is to carry them to the reader, not to substitute generic copy.
- Graphics follow Tufte's representation ethic, not his look: data-ink, direct labels, small multiples, comparisons within eyespan — and the project's own rule, no fabricated zeros.

### What worked

- Full rag-ttc test suite + golangci-lint pass via the pre-commit hook on the backend commit; frontend 32/32 tests green.
- `CaseSummary.Query` was free: the campaign spec already persisted the question text. Surfacing it answers "what was actually asked" with zero new storage.
- The chunk grid makes the fixture legible instantly: chunk-b present in 12/12 stages, chunk-c alive only in the two vector stages before dying at the policy filter.

### What didn't work

- WebFetch strips CSS, so palette extraction needed raw curl+grep (recorded in Step 1; same session).
- No real failures beyond that this round; the main risk was semantic: adding `Description` to `RetrievalCase` changes case artifact digests and episode semantic keys, so old stores' campaigns and new ones are distinct experiment identities. That is correct behavior (the description is part of reviewed intent), but it means re-running campaigns after enriching a manifest.

### What I learned

- The strict manifest decoder (`KnownFields(true)`) means data-file enrichment always requires schema changes first — descriptions cannot be smuggled in.
- Writing the case descriptions was itself analysis: articulating "failing q-hybrid signalss a ranking problem, not a budget problem" is exactly the knowledge the UI should carry to the next reader.

### What was tricky to build

- Slope graph label collisions: with 3-case fixtures two left-side labels share y=1.000; solved with a tiny one-pass overlap resolver that pushes stacked labels apart in sorted order.
- Deciding where descriptions live: per-arm prose must survive from manifest → campaign spec → projector without the manifest being consulted at resume time (spec is canonical), hence Description on `optkitcampaign.Arm` persisted in the spec rather than looked up from the manifest.

### What warrants a second pair of eyes

- The description fields participate in `record.SemanticDigest` of the manifest (identity change on prose edit). If prose-only edits should NOT create a new experiment identity, descriptions would need to be excluded from the digest — a deliberate contract decision for the backend owner.
- ChunkGrid presents stage columns left-to-right as if sequential custody, but early stages are parallel retrieval legs (lexical vs vector). The grid stays truthful (pure presence), but a future version might group columns by leg.

### What should be done in the future

- The remaining scenarios-doc backlog: worst-first failure ordering with verdicts (needs judge/answer stages in trajectories), chunk text previews (chunk-lab), read-only playground, per-case fixed/broke/still-broken churn between journal points.

### Code review instructions

- Backend: `git show f5e5a7e29` — start at `optkitcampaign/system.go` (RetrievalCase), then `campaign.go` (Arm/CampaignSpec), `experimentworkbench/manifest.go`, `specialistapi/types.go`+`projector.go`; validate `GOWORK=off go test ./pkg/ttc/...`.
- Frontend: `git show baf6ce93f` — start at `components/graphs.tsx`, then screen wiring; validate `pnpm typecheck && pnpm test`.
- Live: rebuild, `campaign run` (fresh store — old stores lack descriptions), serve on 8091, open the cockpit; evidence in `various/screenshots/11-cockpit-prose.png` and `12-comparison-prose-graphs.png`.

### Technical details

- New API fields (all optional, additive): `cockpit.description`, `arms[].description`, `cases[].query`, `cases[].description` (cockpit and comparison arm summaries share `ArmSummary`).
- Graphics palette reuses the semantic accent tokens; graphs encode with both color and position/weight, and every SVG carries a full-sentence `aria-label` restating its data.
