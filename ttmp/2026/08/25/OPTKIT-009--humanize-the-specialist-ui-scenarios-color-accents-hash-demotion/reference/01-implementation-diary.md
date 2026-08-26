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

## Step 3: Answer "what changed and what came back" — per-layer widgets over data the store already had

The user, looking at the live comparison, asked the fair question: what actually changed in retrieval, how do I see results, what ranked them — and was the base data even sufficient? Investigation showed the store already held everything: arm config values in the campaign spec, and full `SearchOutput` (ranked results with scores, per-channel fusion contributions, chunk text, effective limit and its source) inside the output artifact — blocked only by a 4 KiB preview cap and unexposed fields. No data recreation was needed.

### Prompt Context

**User prompt (verbatim):** "ok, is it implemented? like i'm looking at http://127.0.0.1:5197/campaigns/campaign%3Afa569e901e85694649ac99615de83c78/compare/limit-1/limit-2 and legit I don't know what changed in the retrieval, and how to look at the results or so, what was used for ranking. Did you already go through all these things?" — followed mid-turn by "or is the base data not enough and you should recreate it ?" and "alsso I guess retrieval/fusion/reranking/etc... all will need different types of widgets for displaying (themselves or the diff, etc...)"

**Assistant interpretation:** Verify honestly what is and isn't implemented; determine whether the recorded data suffices or campaigns need re-recording; and architect per-layer display widgets (values, diffs, outputs) rather than one generic view.

**Inferred user intent:** Standing on the comparison screen, be able to answer: what did I change (in values), what came back (ranked results), and what ranked it (scores/channels) — without reading raw JSON.

**Commit (code):** 4ddd59f03 — "Specialist: show what changed and what came back, per layer"

### What I did

- Backend: `ArmSummary.Config` now carries the arm's executable `RetrievalConfig` (values were already persisted in `CampaignSpec.Arms`); raised `MaximumArtifactPreviewSize` 4096 → 16384 so retrieval outputs (~5 KB) and trajectories (~9 KB) preview inline.
- Frontend: introduced `src/layerwidgets/` — a registry mapping layer names to value-diff widgets and artifact schemas to output widgets, with generic fallbacks. Retrieval registers both: `RetrievalDiffView` renders "candidates kept per retriever: 2 (was 1)" inside the Configuration diff panel; `RetrievalOutputView` renders the Episode output as the question asked, the limit applied and its source, and a ranked-evidence table (rank, chunk, score, per-channel contributions like `bm25 #2 · 0.0161 ×1` / `vector #1 · 0.0164 ×1`, and the chunk text itself). Empty result lists say "for a negative case, that can be exactly right."
- `ArtifactBlock` renders a registered widget when the preview's schema matches, keeping raw JSON behind a disclosure.
- 4 new widget tests (36 total); validated live — no campaign re-run required since only the projector changed. Screenshots 13 (value diff) and 14 (ranked evidence) archived.

### Why

- The bottleneck was projection, not recording: re-running campaigns would have changed nothing. Establishing that before touching data answered the "recreate it?" question with evidence.
- One generic table can't explain retrieval limits and fusion weights and reranker cutoffs; a registry gives each layer its own vocabulary while unregistered layers degrade gracefully.

### What worked

- Chunk text turned out to already be in the output artifact (`SearchResult.Text`) — so "what does this chunk say" is partially answered without chunk-lab, at least for retrieved chunks.
- Fusion contributions (`rag.Contribution{channel, rank, weight, value}`) directly answer "what was used for ranking".

### What didn't work

- First commit attempt failed the pre-commit gofmt check on `types.go` (hand-aligned struct tags); `gofmt -w` fixed it. The earlier chained `git commit | tail` also masked the failure — the commit had not landed until retried.

### What was tricky to build

- Preview typing: previews arrive as `unknown` JSON, so `RetrievalOutputPreview` treats every field as optional and the widget renders honestly labeled gaps ("unscored", "no channel breakdown recorded") instead of assuming shape.

### What warrants a second pair of eyes

- The 16 KiB preview cap now exposes full trajectories inline; confirm that's acceptable under the sensitivity policy (both artifacts are `internal`, same as before — only size policy changed).
- `RetrievalDiffView` reads arm-level config; when future campaigns vary multiple layers, per-layer value surfacing (beyond retrieval) needs the API to expose resolved layer values, not just the retrieval config.

### What should be done in the future

- Register widgets for fusion, reranking, chunking as their value schemas surface; case-level side-by-side pipeline diff remains the next big view.

### Code review instructions

- `git show 4ddd59f03`: start at `apps/specialist/web/src/layerwidgets/`, then `specialistapi/types.go` (Config field, preview cap) and `projector.go` (armConfig).
- Validate: backend `GOWORK=off go test ./pkg/ttc/specialistapi/`; frontend `pnpm test`; live at the comparison URL — Configuration diff shows the values table, any episode's output shows ranked evidence.

## Step 4: Record the data, not just the metadata — per-stage candidates with scores

The user pressed on the pipeline screen: stages showed counts and chunk IDs, but not what a stage actually did — its query, its results, its scores. Investigation delivered a split verdict on data sufficiency: the query and final ranked output were recorded but unsurfaced, while **per-stage scores were genuinely never recorded** — every stage constructor had `rag.Hit.Score`, `FusedHit.Contributions`, and `Evidence.RetrievalScore/RerankerScore` in hand at runtime and kept only the chunk IDs (the "candidate artifact" is a digest of the ID list, not stored content). So this step changed the recording itself and rebuilt the store.

### Prompt Context

**User prompt (verbatim):** "ok but i can't see the actual data, ,right? it's just metadata? like how do I understand what #10 retrieval.fused or so did, what query it ran, what results it got, etc... ?" — followed mid-turn by: "we could potentially use protobuf to define the data schemas of the artifacts and thus generate frontend / backend types. Or generate the frontend types to parse them from the current go definition, that way we can expose them via the API in a generic manner and the frontend can register custom widgets and all?"

**Assistant interpretation:** Make each pipeline stage show its actual contents (query, candidates, scores), determine whether that required a recording change (it did), and assess schema-codegen options (protobuf vs Go-to-TS generation) for keeping widget types in sync.

**Inferred user intent:** Stand on a stage row and understand what it did without insider knowledge or raw JSON; and stop hand-maintaining parallel type definitions as the schema surface grows.

**Commit (code):** 20e8d4621 — "Record what each retrieval stage actually held, and say what stages mean"

### What I did

- `ttcsearch.RetrievalStage` gains `Candidates []StageCandidate` (chunk, rank, channel, score, reranker score, fusion contributions), filled by `hitStage`/`fusedStage`/`evidenceStage` from data the runtime already had; `idStage` stays membership-only rather than fabricating scores. Stage trajectory payloads are the stage struct verbatim, so trajectories and previews enriched automatically.
- `PipelineView` gains `query` and `case_description` (looked up from the persisted campaign spec), shown at the top of the pipeline screen.
- Frontend: `STAGE_EXPLANATIONS` glossary — one plain sentence per stage name ("Both channels merged into one ranking — each candidate's score combines its keyword and semantic evidence") shown on every stage row; `StageRecordView` registered for `schema:rag-ttc.retrieval-stage/v1` renders the candidate table inline; membership-only stages say so honestly.
- Rebuilt, recreated the store (`campaign:0ac671bd…`), verified via API and screenshot 15: lexical.raw shows chunk-b at 10.0 / chunk-a at 9.0; retrieval.fused shows RRF scores with the bm25+vector breakdown; the chunk-c drop at vector.policy_filtered is visible with its pre-drop score.

### Why

- "Metadata about data" cannot answer a debugging question; the scores at each stage are the mechanism of every retrieval failure. Recording them costs a few hundred bytes per stage and makes the sealed trajectory a complete story.

### What worked

- Because stage payloads are the struct itself, one type change flowed through trajectory, preview, and widget with no projector edits.
- Episode IDs stayed stable across the store rebuild (candidates are outputs, not part of the episode spec), so old deep links to episodes still resolve.

### What didn't work

- N/A this step — build, full suite, and lint passed first try; the one caveat is operational: stores recorded before this change lack `candidates` and show the membership-only fallback.

### What was tricky to build

- Honesty at the edges: `idStage`-built stages (e.g. evidence.admitted in some paths) have no scores at runtime; the widget renders "membership only; per-candidate scores were not captured for this run" instead of hiding the difference.

### What warrants a second pair of eyes

- Stage payload size grew (~300 B → ~1 KB for scored stages); trajectory sizes scale with stages × candidates × contributions. Fine at limit 2; worth a look before limit-50 campaigns.
- The schema ID `rag-ttc.retrieval-stage/v1` was not bumped for the additive field; if consumers treat stage records as closed schemas, this should have been v2 — a contract call for the backend owner.

### What should be done in the future

- Schema/type generation (user proposal): recommended path is Go-structs-as-source-of-truth with generated TS (e.g. tygo) for read-side types, wired into `go generate` — protobuf is already in-repo (buf, chat server) but migrating *artifact* schemas to proto changes the JSON encoding that semantic digests hash, i.e. a store-identity migration; reserve proto for new schema families (judge traces) where cross-language matters from day one. The frontend widget registry keyed by schema ID is already the extension point either way.

### Code review instructions

- `git show 20e8d4621`: `pkg/ttc/search/service.go` (StageCandidate + constructors), `specialistapi/pipeline.go` (query lookup), `apps/.../layerwidgets/retrieval.tsx` (glossary + StageRecordView).
- Validate: `GOWORK=off go test ./pkg/ttc/search/ ./pkg/ttc/specialistapi/`; frontend `pnpm test`; live — any stage row now opens with its explanation and candidate table. Evidence: `various/screenshots/15-stage-data.png`.

## Step 5: Chunk content in the record — judge relevance, not counts

The user's sharpest test yet: "nothing where I can actually look at the results and say 'well those were bad because XYZ, it retrieved a rose when it should have been a thuja'". Correct — stage tables had scores for chunk IDs, but content existed only for final survivors, so a chunk filtered mid-pipeline could never be read. Fixed at the recording layer: the search tool now resolves every chunk any stage touched through the content store at record time and stores a chunk catalog in the run's output.

### Prompt Context

**User prompt (verbatim):** "ok but retrieval.fused can't show me actual data from the runtime query, it just shows me the meta parameters, right? nothing where I can actually look at the results and say \"well those were bad because XYZ, it retrieved a rose when it should have been a thuja\" or so"

**Assistant interpretation:** Make chunk content readable at every point a chunk ID appears — including candidates that were later filtered out — so a human can judge relevance directly.

**Inferred user intent:** The autopsy workflow: read what was retrieved, in full, at each stage, and decide whether the pipeline's choices were right.

**Commit (code):** 4ff578c1c — "Record chunk content with every run so candidates can be judged, not just counted"

### What I did

- `search.Service.DescribeChunks` resolves IDs via `content.Store.Chunks` plus the source catalog (title/URL); explicitly observability-only — store failure shortens the catalog, never fails retrieval.
- `SearchOutput.ChunkCatalog` records every distinct chunk ID across all stages (first-appearance order) with document, title, and full text — including chunks the policy filter later removed.
- `PipelineView.ChunkCatalog` decodes it from the output artifact (warning diagnostic on decode failure, never fatal).
- UI: "Chunks this run touched" panel renders each chunk readable in full; `ChunkCatalogContext` lets `StageRecordView` add a "What it says" column and title to every candidate row. Stores without catalogs show "content not in catalog" chips.
- Store recreated (`campaign:36ea7379…`); screenshot 16: chunk-c reads "The inventory table stores quantity and warehouse identifiers" — a reviewer can now see it was correctly filtered as restricted inventory data.

### Why

- Relevance is judged by reading, and the moment that matters is often a mid-pipeline drop. Content resolved at record time also keeps the sealed-record property: the catalog shows the corpus as it was during the run, not as it is later.

### What worked

- Single implementation point (`RunRoute`) covers all executors; build, full suite, lint green first try.

### What didn't work

- N/A — though note the catalog silently deduplicates by store response; unknown IDs are skipped by design.

### What was tricky to build

- Failure semantics: the catalog must never break retrieval, so `DescribeChunks` swallows store errors and returns what it can, and the projector degrades to a warning diagnostic. The UI states absence instead of hiding the column.

### What warrants a second pair of eyes

- Catalog size scales with corpus chunk length × touched chunks; at limit-50 with long chunks the output artifact could pass the 16 KiB preview cap, making the catalog invisible again (it's decoded, not previewed — actually unaffected; only the raw preview JSON disclosure is capped). Verify the decode path stays exempt from the preview cap as intended.
- Sensitivity: the catalog embeds text of policy-FILTERED chunks in an `internal` artifact — deliberate for specialist review, but confirm this matches the sensitivity policy's intent (a restricted chunk's text now lives in the run output).

### What should be done in the future

- Ground-truth overlay: cases carry required/forbidden targets — marking each catalog chunk as "expected", "forbidden", or "incidental" in the stage tables would complete the rose-vs-thuja judgment at a glance.

### Code review instructions

- `git show 4ff578c1c`: `pkg/ttc/search/search.go` (catalog collection), `service.go` (DescribeChunks), `specialistapi/pipeline.go` (decode), `apps/.../PipelineScreen.tsx` + `layerwidgets/`.
- Validate live: pipeline of any episode in campaign `campaign:36ea7379a2bc647ff16e6b0192c0a8f6` — read chunk-c's text in "Chunks this run touched". Evidence: `various/screenshots/16-chunk-content.png`.
