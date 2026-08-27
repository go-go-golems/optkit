---
Title: Investigation Diary
Ticket: OPTKIT-003
Status: active
Topics:
    - optkit
    - ui
    - query-plane
    - visualization
    - scientific-workflow
    - architecture
    - local-development
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://internal/web/server_test.go
      Note: GET-only, static, campaign, and SSE integration evidence
    - Path: repo://internal/web/static/index.html
      Note: Semantic explorer document shell
    - Path: repo://internal/web/static/styles.css
      Note: Pure-white retro-monochrome visual system
    - Path: repo://store/sqlite/query.go
      Note: Stable read-only campaign catalog
    - Path: repo://ttmp/2026/08/24/OPTKIT-003--read-only-scientific-campaign-explorer/scripts/01-capture-numbergame-evidence.sh
      Note: Reproducible scientific evidence capture
    - Path: repo://ttmp/2026/08/24/OPTKIT-003--read-only-scientific-campaign-explorer/scripts/02-validate-explorer.sh
      Note: Complete explorer validation workflow
ExternalSources: []
Summary: Chronological evidence, implementation decisions, failures, validation, and review instructions for the read-only Optkit scientific campaign explorer.
LastUpdated: 2026-08-25T00:40:00-04:00
WhatFor: Explain how the explorer architecture and implementation were derived, tested, corrected, and delivered.
WhenToUse: Read before reviewing, extending, debugging, or resuming OPTKIT-003.
---


# Diary

## Goal

Design, explain, implement, validate, and deliver a read-only Optkit scientific campaign explorer that helps a new engineer and scientist navigate CLI-authored campaign evidence without adding browser mutation paths.

## Step 1: Create the explorer research ticket and map the evidence substrate

I created OPTKIT-003 as a separate ticket because the query plane and scientific explorer have a coherent scope distinct from OPTKIT-002's broad framework/RAG implementation program. The investigation started from the actual imported Optkit packages and a newly generated Numbergame campaign rather than treating the earlier scenario prose as an API specification.

The resulting primary guide explains journal authority, immutable artifacts, episode trajectories, observations and epochs, trials and paired estimates, budgets, query interfaces, projections, GET APIs, SSE, search, artifact safety, navigation, visualizations, testing, and phased implementation. It repeatedly distinguishes current code, implemented ticket work, and remaining roadmap.

### Prompt Context

**User prompt (verbatim):**

> Create a new docmgr ticket, and Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Create an independent, exhaustive ticket that teaches an intern how Optkit evidence works and specifies the read-only scientific explorer from storage through UI, then validate and deliver the bundle to reMarkable.

**Inferred user intent:** Produce a durable onboarding and implementation reference that makes the explorer understandable without relying on conversation history.

### What I did

- Added vocabulary for `ui`, `query-plane`, `visualization`, and `scientific-workflow`.
- Created `OPTKIT-003 — Read-Only Scientific Campaign Explorer`.
- Added a primary design/implementation guide and strict investigation diary.
- Added tasks for repository inventory, scientist workflow, backend/API/SSE/search design, visualization design, implementation guidance, validation, implementation, and reMarkable delivery.
- Read campaign journal/controller/state, SQLite journal/schema, artifact store/ref, episode event/writer, trial/estimate, observations/epochs, scheduler, budgets, local profile, CLI, architecture tests, and the complete Numbergame execution/measurement/analysis path.
- Added `scripts/01-capture-numbergame-evidence.sh`.
- Generated `sources/01-numbergame-evidence.md` from a fresh temporary campaign and deleted the temporary store after capture.
- Wrote the primary guide with evidence-backed current state, diagrams, APIs, pseudocode, UI workflow, visual specifications, phases, tests, decisions, risks, and file references.

### Why

- The explorer must visualize authoritative semantics rather than private SQLite implementation details.
- A fresh Numbergame capture provides concrete counts and payloads for every proposed page.
- A new intern needs the distinction between control events and episode trajectory events before implementing a timeline.
- Scientific correctness requires intervention integrity, missing-pair handling, epoch compatibility, and evidence provenance to appear in the UI design.

### What worked

- The capture reproduced a completed 81-event campaign with eight episodes, 24 observations, eight budget reservations/usage commits, a +0.71041675 paired delta, and an eligible decision.
- Existing package boundaries naturally support a query plane: journal/artifacts are authoritative while `projection.RebuildOverview` is already pure and rebuildable.
- Numbergame is sufficiently deep to drive lineage, matrix, timeline, measurement, budget, estimate, and decision views.

### What didn't work

- N/A during the research/capture step. Source reading, demo execution, inspection, and verification succeeded.

### What I learned

- `campaign.Journal.Read` is cursor-based but not page-bounded; a scalable query adapter must not mistake this for database-side pagination.
- The artifact store validates bytes but does not decide browser disclosure. Sensitivity and reachability belong to query policy.
- External CLI/workers can write while the server runs, so polling authoritative journal heads is more correct than an in-process-only event bus.
- Generic control events need schema-owned projectors to produce rich domain views without coupling campaign core to Numbergame or RAG.

### What was tricky to build

- The guide had to be implementation-specific without promising missing packages as current reality. Every major section therefore names existing symbols and labels future interfaces explicitly.
- “Read-only” still requires complex behavior: reconnectable live follow, replay, search, safe previews, deep links, and derived scientific visualizations.
- The guide had to preserve both a generic architecture and a complete Numbergame story; the schema registry/plugin boundary keeps those concerns separate.

### What warrants a second pair of eyes

- Review the proposed separation between generic query models and domain projectors.
- Review artifact reachability/sensitivity defaults before RAG data is exposed.
- Review whether event page limits must be implemented inside SQLite before any large campaign is served.
- Review if multiple store roots belong in v1; the guide assumes one configured local root.

### What should be done in the future

- Complete the implemented v0 gaps listed in the guide: entity routes, nested trajectories, server search, replay slider, comparison, and scalable pagination.
- Keep the guide synchronized with implementation commits rather than leaving “proposed” claims stale.

### Code review instructions

- Start with the primary guide Sections 1-5, 7-12, 15-17, and 23-27.
- Reproduce evidence using `scripts/01-capture-numbergame-evidence.sh`.
- Compare `sources/01-numbergame-evidence.md` event counts with `examples/numbergame/demo.go`.
- Confirm every authoritative claim maps to the primary file references in Section 31.

### Technical details

- Ticket: `ttmp/2026/08/24/OPTKIT-003--read-only-scientific-campaign-explorer`.
- Primary guide size after implementation updates: approximately 50 KiB.
- Captured event counts: 81 control events, 8 completed episodes, 24 observations.
- Capture source revision: `5db1815dabbf1150d59a0b7554e023df17cf860a`.

## Step 2: Implement the dependency-free white monochrome explorer

The scope changed from documentation-only to a working first explorer slice. I implemented a read-only query package, stable campaign listing, GET JSON endpoints, cursor-based SSE, static embedding, a CLI serve command, and a dependency-free HTML/CSS/JavaScript interface. The interface renders actual Numbergame campaign data rather than a hard-coded mock.

The visual design is a flat white scientific document with retro Macintosh restraint: black/gray structure, modern system fonts, monospace evidence, and color used only for link/status/diagnostic text. There is no simulated window chrome, menu bar, paper texture, beige surface, framework, Node dependency, or frontend build step.

### Prompt Context

**User prompt (verbatim):**

> implement as simple HTML (retro monochrome macos look, but no window chrome / menu bar, use modern fonts, color accents as font color) and JS.
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Replace the proposed React/Vite frontend with plain semantic HTML, CSS, and JavaScript and implement a working read-only explorer matching the requested visual language.

**Inferred user intent:** Keep the explorer lightweight, understandable, distributable in one Go binary, and visually disciplined rather than introducing a frontend application stack.

**Commit (implementation):** `1b09c3c6ccae76dff66501fd1401066ef6353c6a` — "feat(ui): add read-only campaign explorer"

**Commit (white surfaces):** `d90ae43061ff8ab61d03ad19c410eaf9a0d4cecc` — "style(ui): use pure white explorer surfaces"

**Commit (no-CGO test gate):** `368ca452ff34a2665b3d7c33080249948379d537` — "test(ui): gate SQLite explorer tests on cgo"

### What I did

- Added `store/sqlite/query.go` with stable read-only campaign listing.
- Added `query/service.go` with read-only interfaces and campaign/event models.
- Added campaign overview, budget, integrity, system/dataset/trial identity, event-kind inventory, and bounded JSON event payload projection.
- Added event response limits: default 200 and hard maximum 500.
- Denied confidential/restricted event payload previews; allowed bounded valid JSON for public/internal payloads.
- Added `internal/web/server.go` with standard-library Go 1.22+ `http.ServeMux`.
- Added GET health, campaign list, campaign detail, events, and SSE routes.
- Added CSP, nosniff, frame denial, and referrer policy.
- Added direct `go:embed` for `index.html`, `styles.css`, and `app.js`.
- Added `optkit serve --store PATH --listen 127.0.0.1:8080`.
- Added the vanilla browser UI with campaign list, hash navigation, global search, stage rail, metrics, configuration lineage, trial matrix, paired deltas, budgets, event-kind inventory, 81-event timeline, evidence drawer, and EventSource follow.
- Added query and HTTP integration tests over real Numbergame stores.
- Added `scripts/02-validate-explorer.sh` and `sources/02-explorer-validation.txt`.
- Updated the primary guide from React/Vite assumptions to the actual dependency-free stack and marked implemented versus remaining capabilities.

### Why

- Plain assets exactly meet the user requirement and eliminate an unnecessary Node toolchain.
- Direct embed makes development/release topology one Go process and one binary.
- The first UI must operate against journal facts to validate the architecture, not only display a sketch.
- Hash routing provides deep links without an SPA fallback or server-side route duplication.
- Native `<progress>` elements provide data-driven bars under a strict no-inline-style CSP.

### What worked

- Full `go test ./... -count=1` passes.
- Focused race tests for `query`, `internal/web`, and `cmd/optkit` pass.
- Full repository `make ci-check`, `make race`, and `make lint` pass through the validation script.
- GET API/static integration tests pass and POST `/api/v1/campaigns` returns 405.
- SSE replay after sequence 80 returns event 81.
- Chromium rendered one real completed campaign with all six stages complete, 8/8 matrix cells, +0.710417 paired delta, budgets, decision, inventory, and timeline.
- Global search reduced the timeline to `1 / 81 EVENTS` for `EstimateRecorded`.
- Evidence drawer opened the exact estimate JSON.
- A fresh Chromium tab reported zero console errors or warnings.

### What didn't work

- The first browser load failed JavaScript parsing because a template string used `\0`/`\00` as a map-key separator. Chromium reported:

  `Octal escape sequences are not allowed in template strings.`

  I replaced the separator with `::` and restarted the embedded server.
- The first campaign render produced 14 CSP errors because dynamic chart and budget geometry used inline `style` attributes while CSP allowed only `style-src 'self'`. I replaced those elements with native `<progress>` controls and external pseudo-element styling. A new browser tab then reported zero errors.
- The first lint run found an exhaustive sensitivity switch and new tests using `testing.T.Context`, which requires Go 1.24 while `go.mod` declares Go 1.23. I made public/internal and unknown sensitivity cases explicit and used `context.Background()` in the new tests.
- The first complete validation script failed in `CGO_ENABLED=0 go test ./...` because query/web integration tests create SQLite stores. I added `//go:build cgo` to those integration tests. Normal no-CGO package compilation remains covered, while SQLite integration runs under CGO.
- The first visual request used off-white “paper” tokens. The user clarified:

  > white background, no paper

  I changed every surface to pure white, removed the paper token names, and kept only monochrome structure.

### What I learned

- A strict CSP meaningfully shapes even a simple frontend: data-driven geometry should avoid inline styles or use a reviewed nonce strategy.
- Browser console state can accumulate errors from prior navigations; validation used a fresh tab to prove the corrected build had zero new messages.
- Hash routing is sufficient for the current two-page explorer and lets exact `/` remain a simple document route.
- Real event payload previews make a useful domain-specific UI possible before a full server-side schema projector registry exists.
- The v0 event API applies response bounds after reading the remaining journal; this is correct for Numbergame but must move paging into SQLite for RAG-scale campaigns.

### What was tricky to build

- The service had to remain read-only at the type level. `query.MetadataStore` includes only list/read/head/verify/snapshot methods, even though the concrete SQLite object also implements mutation interfaces.
- A single generic event API had to support both rich Numbergame derivation and unknown future schemas. It returns complete metadata and optional bounded JSON rather than claiming a universal typed payload.
- The trial matrix is derived in JavaScript by joining `TrialPlanned` cases/arms with `EpisodeCompleted` payloads. This validates the information model but should move to a server projector before scale/other domains.
- Color accents had to remain textual. Structural borders/progress marks are monochrome; green/blue/red/amber/violet communicate through text and remain nonessential.
- SSE must observe external writers. It polls journal pages by sequence rather than relying on an in-memory server event bus.

### What warrants a second pair of eyes

- Review the `query.Service.populatePayload` local sensitivity policy before non-Numbergame data is served.
- Review the fact that event payload previews return internal JSON to a loopback-local browser.
- Review SSE slow-client behavior; current direct writes rely on HTTP connection backpressure and do not yet impose an explicit queued-byte ceiling.
- Review event database pagination before large RAG campaigns.
- Review JavaScript derivation of lineage/trial/analysis for missing/failed partial campaigns.
- Review the 20 KiB single `app.js`; split it only when continued feature growth harms readability.
- Review visual contrast and keyboard behavior with an automated accessibility audit.

### What should be done in the future

- Add entity routes and nested episode/result/trajectory/artifact traversal.
- Add a replay slider with `at=<sequence>` projection cutoff.
- Add server-side structured search and projectors.
- Add measurement QA grouped by case/arm/repeat and epoch.
- Add large-campaign pagination/virtualization and cross-campaign comparison.
- Add CLI output for the canonical explorer URL after campaign workflows.

### Code review instructions

- Start at `query/service.go` to verify authority and preview boundaries.
- Review `store/sqlite/query.go` for read-only catalog behavior.
- Review `internal/web/server.go` for GET-only routes, SSE cursor handling, CSP, and exact static routing.
- Review `internal/web/static/index.html`, then CSS tokens/layout, then `app.js` derivations.
- Run `scripts/02-validate-explorer.sh` and require `EXPLORER_VALIDATION=PASS` and `POST_MUTATION_STATUS=405`.
- Run the browser against a fresh demo and verify search, campaign navigation, estimate drawer, and a clean console.

### Technical details

- API version: `optkit.query/v1`.
- Default/hard event limits: 200/500.
- Payload preview ceiling: 256 KiB.
- Default bind: `127.0.0.1:8080`.
- Browser routes: `#/campaigns`, `#/campaigns/:id`.
- Static files: directly embedded; no build step.
- Final validation marker: `EXPLORER_VALIDATION=PASS`.

## Step 3: Enforce pure white surfaces and deliver the bundle

The final pass removed the remaining paper metaphor from CSS names and values, reran the complete validation workflow, validated ticket health, and uploaded the architecture package to reMarkable. The implementation and documentation now consistently describe a pure-white, flat, dependency-free explorer.

The uploaded bundle contains the ticket index, 1,575-line intern guide, diary through implementation, and 709-line generated Numbergame evidence capture. The local ticket remains authoritative for this delivery receipt and the full validation transcript.

### Prompt Context

**User prompt (verbatim):**

> white background, no paper
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Remove off-white/paper styling completely while retaining the retro monochrome structure, modern typography, and text-only color accents.

**Inferred user intent:** Avoid nostalgic skeuomorphism and keep the explorer visually clean, white, technical, and contemporary.

### What I did

- Replaced off-white surface values with pure `#ffffff`.
- Removed `paper` token names from the CSS implementation.
- Kept all structural rules/progress graphics black or gray and all non-monochrome accents textual.
- Added CGO build constraints to integration tests that require the local SQLite profile.
- Reran `scripts/02-validate-explorer.sh` successfully.
- Ran `docmgr doctor --ticket OPTKIT-003 --stale-after 30` successfully.
- Committed the complete guide, diary, scripts, generated evidence, validation report, relations, tasks, changelog, and vocabulary.
- Ran a successful reMarkable bundle dry-run.
- Uploaded `OPTKIT-003 Scientific Campaign Explorer Guide.pdf` to `/ai/2026/08/25/OPTKIT-003`.

### Why

- Visual token names should match the actual design to prevent future contributors from reintroducing beige/paper simulation.
- Integration tests should state their SQLite/CGO dependency rather than breaking the declared no-CGO package gate.
- The dry-run validates document inclusion and destination before cloud mutation.

### What worked

- The complete validation transcript ends with `EXPLORER_VALIDATION=PASS`.
- POST mutation check returns 405.
- Full tests, no-CGO tests, race tests, vet, build, and lint pass.
- Browser smoke validation showed the expected Numbergame views with zero fresh console errors.
- Docmgr doctor reported all checks passed.
- The real upload returned:

  `OK: uploaded OPTKIT-003 Scientific Campaign Explorer Guide.pdf -> /ai/2026/08/25/OPTKIT-003`

### What didn't work

- The first complete validation run failed under `CGO_ENABLED=0` because the newly added query and web integration tests opened SQLite-backed local profiles. The code itself compiled without CGO; the tests required persistence. I added `//go:build cgo` to those integration files and reran the complete validation successfully.

### What I learned

- Build portability and integration-test capability are separate contracts. The no-CGO gate should compile/query packages while SQLite behavioral tests run only where the system binding exists.
- Direct embedding keeps the final visual correction extremely small: one CSS commit and server restart.
- A 6,715-word guide plus generated evidence remains navigable as a single ToC-equipped reMarkable bundle.

### What was tricky to build

- The final bundle must be uploaded before the diary can contain the upload receipt. The reMarkable copy therefore contains the diary through Step 2; this local Step 3 records delivery. Overwriting the new PDF solely to include its own receipt would require force and risk future annotations, so it was not justified.
- The repository date path remains `2026/08/24`, while actual final delivery occurred after midnight on `2026/08/25`; the remote directory uses the actual delivery date.

### What warrants a second pair of eyes

- Open the PDF on-device and confirm large code/API sections and diagrams remain readable.
- Review the final CSS on high-contrast and reduced-motion settings.
- Review whether confidential/internal preview policy should become configurable before serving real RAG campaigns.

### What should be done in the future

- Continue with entity routes, trajectory navigation, replay slider, structured search, and large-campaign pagination as separate focused tickets.
- Upload a new versioned PDF rather than forcing replacement after reviewers annotate this document.

### Code review instructions

- Review commits `1b09c3c`, `d90ae43`, and `368ca45` before the documentation commit.
- Run `scripts/02-validate-explorer.sh` and inspect `sources/02-explorer-validation.txt`.
- Run `docmgr doctor --ticket OPTKIT-003 --stale-after 30`.
- Launch with `optkit serve` and verify the campaign list, Numbergame dashboard, global event search, and evidence drawer.

### Technical details

- Delivery directory: `/ai/2026/08/25/OPTKIT-003`.
- PDF: `OPTKIT-003 Scientific Campaign Explorer Guide.pdf`.
- Bundle inputs: index, primary guide, investigation diary, Numbergame evidence capture.
- Documentation commit before delivery: `4a14d23b99d445e46ea6f8e0dbb7d39f7d50aabf`.
