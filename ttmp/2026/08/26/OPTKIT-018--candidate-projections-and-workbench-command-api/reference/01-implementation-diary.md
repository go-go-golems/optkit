---
Title: Implementation Diary
Ticket: OPTKIT-018
Status: complete
Topics:
    - architecture
    - design
    - implementation
    - optkit
    - rag-ttc
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/design-doc/04-backend-first-optimization-workbench-program-roadmap.md
      Note: Overall program context that keeps the ticket aligned
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/design-doc/01-intern-guide-to-candidate-read-projections-and-the-workbench-command-api.md
      Note: Primary design deliverable whose research and delivery this diary records
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/various/live-workbench-api-smoke.log
      Note: Exact live composed server workflow evidence
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/various/remarkable-implementation-dry-run.log
      Note: Completed workbench API implementation bundle selection evidence
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/various/remarkable-implementation-upload.log
      Note: Successful workbench API implementation upload receipt
    - Path: repo://rag-ttc/pkg/ttc/experimentworkbench/workbench_service.go
      Note: Primary command application implementation
    - Path: repo://rag-ttc/pkg/ttc/workbenchapi/server.go
      Note: Strict command HTTP implementation
ExternalSources: []
Summary: Chronological research, design, implementation, validation, and delivery record for OPTKIT-018.
LastUpdated: 2026-08-26T14:20:28.245641732-04:00
WhatFor: ""
WhenToUse: ""
---





# Diary

## Goal

Design historical candidate projections and a separate authorized workbench command API. This diary records research, design decisions, validation, commits, delivery evidence, failures, and continuation instructions.

## Step 1: Investigate the current system and write the intern design guide

This step established the ticket before implementation starts. I read the parent OPTKIT-011 roadmap and architect brief, traced the current source paths that constrain this ticket, and wrote a self-contained intern guide that explains both current behavior and the proposed implementation sequence.

The design preserves the backend-first dependency order. It records concrete APIs, pseudocode, diagrams, file/line references, tests, risks, exclusions, and an exit gate so later work can proceed without reconstructing the program context.

### Prompt Context

**User prompt (verbatim):**

> Ok, create all these tickets, and for each ticket, Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.
>
> commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)
>
> Youc an run experiments and such, store all scripts you write in the scripts/ folder of the relevant ticket.

**Follow-up user prompt (verbatim):**

> add the description of each ticket and make an upfront design doc with the ticket goal and content so that we can keep track as we go on, to not lose the overall picture.

**Follow-up user prompt (verbatim):**

> since each document will require a significant amount of planning / investigation, so we want to make sure we don't forget what the do cis. Once that overview doc is written, update the goal to reference it.

**Assistant interpretation:** Expose sealed intent/catalog facts for reads and adapt compile/preview/seal applications to HTTP without mixing them into specialist projections.

**Inferred user intent:** Give React authoritative read and command contracts while preserving historical determinism.

### What I did

- Mapped current GET-only ServeMux, fact loader, comparison projection, and application services.
- Specified candidate/catalog responses, command routes, authorization actions, typed errors, idempotency, and preview envelopes.
- Designed historical fallback, sealed catalog labeling, strict HTTP, security, concurrency, and live smoke tests.
- Added a substantive ticket index description, implementation tasks, and this strict-format diary.
- Kept program-wide ticket scaffolding/generation scripts under the parent OPTKIT-011 `scripts/` directory.

### Why

- Each ticket requires significant independent investigation, but its contracts must remain aligned with the parent roadmap.
- An intern should understand ownership, runtime behavior, invariants, and test evidence before editing code.
- Up-front acceptance gates prevent frontend or persistence work from outrunning foundational contracts.

### What worked

- Baseline focused tests passed before documentation: `optkit/space`, `optkit/examples/numbergame`, and the relevant RAG-TTC optimization/workbench/campaign/search/specialist packages.
- Current CLI help confirmed the existing `config` and `campaign` Glazed command surfaces.
- Concrete source evidence was sufficient to define this ticket without speculative production changes.

### What didn't work

- N/A for this ticket's own design. During the program-wide index update, the first batch used guessed generated timestamps and exact-text edits failed for OPTKIT-013 through OPTKIT-020; the retry used smaller stable frontmatter blocks and succeeded.

### What I learned

- One server may host both route groups, but package/application boundaries must keep historical reads incapable of compiling or sealing proposals.
- Ticket boundaries are most reliable when each names both its upstream contract and the next consumer.
- Documentation should distinguish observed runtime behavior from proposed APIs and future implementation choices.

### What was tricky to build

- One server may host both route groups, but package/application boundaries must keep historical reads incapable of compiling or sealing proposals.
- The guide had to be self-contained without copying the entire parent architecture. It summarizes prerequisites, then links each claim to the file that owns it.

### What warrants a second pair of eyes

- Review actor spoofing, sensitive diagnostics/logging, typed error mapping, restart-safe idempotency, and sealed-versus-current catalog behavior.
- Review all proposed public schema/API names before implementation makes them expensive to change.

### What should be done in the future

- Implement candidate fact projections first, then catalog reads, command applications, strict HTTP adapters, and frontend fixtures.
- Update this diary immediately when implementation reveals a false assumption or accepted contract change.

### Code review instructions

Start with the ticket's `design-doc/01-*.md`, then inspect these decision-shaping files:

- `rag-ttc/pkg/ttc/specialistapi/projector.go`
- `rag-ttc/pkg/ttc/specialistapi/types.go`
- `rag-ttc/pkg/ttc/specialistapi/http.go`
- `rag-ttc/pkg/ttc/experimentworkbench/service.go`
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/serve.go`

Validate the design workspace with:

```bash
docmgr validate frontmatter --doc optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/reference/01-implementation-diary.md
docmgr doctor --ticket OPTKIT-018 --stale-after 30
```

Run the focused code/test commands listed in the guide before and after implementation.

### Technical details

- Parent program map: `OPTKIT-011/design-doc/04-backend-first-optimization-workbench-program-roadmap.md`.
- Ticket state: `index.md`, `tasks.md`, and `changelog.md` in this workspace.
- No production code behavior changed while writing this step.

## Step 2: Validate, commit, and deliver the guide

This step converted the researched guide from a working document into a reviewed ticket deliverable. The ticket's frontmatter, relations, tasks, and changelog were validated; the documentation was committed in a dependency-coherent batch; and the index, guide, and diary were rendered and uploaded as one reMarkable PDF with a table of contents.

The implementation tasks intentionally remain open. This delivery completes the up-front planning package and gives the future implementer an evidence-backed starting point, not a false claim that production behavior has already changed.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Finish the ticket documentation package with validation, coherent Git history, strict diary evidence, and a ticket-specific reMarkable delivery.

**Inferred user intent:** Make the plan durable and reviewable before implementation starts, while preserving a clear distinction between completed design work and pending code tasks.

**Commit (documentation):** `ec784b0d84e47ad45a7429f98482e9678df5dac7` — "OPTKIT-018-020: design workbench delivery layers"

### What I did

- Ran `docmgr doctor --ticket OPTKIT-018 --stale-after 30` and obtained `All checks passed`.
- Scanned the index, guide, and diary for generated placeholder sections; none remained.
- Ran focused baseline Go tests before documentation changes; all selected Optkit and RAG-TTC packages passed.
- Ran `git diff --check`/`git diff --cached --check`, corrected whitespace findings, and committed the ticket package.
- Ran the program upload script in `--dry-run` mode, then rendered and uploaded `OPTKIT-018 Workbench API Guide.pdf`.
- Preserved upload evidence in `various/remarkable-dry-run.log` and `various/remarkable-upload.log`.

### Why

- A detailed guide is only useful when frontmatter, links, task state, and delivery artifacts agree.
- Grouping commits by dependency layer keeps review focused while avoiding one nine-ticket mega-commit.
- The reMarkable bundle lets the architecture be reviewed away from the source tree without losing the index or diary context.

### What worked

- `remarquee` reported: `OK: uploaded OPTKIT-018 Workbench API Guide.pdf -> /ai/2026/08/26/OPTKIT-018`.
- The upload rendered the real Markdown rather than only performing a path dry-run.
- The ticket retains open implementation tasks and checked documentation/delivery tasks separately.

### What didn't work

- The first staged diff check failed with `new blank line at EOF` in generated changelogs and `trailing whitespace` on blank quoted prompt lines (`+> `) in generated diaries. The diary generator was corrected to emit `>` on blank quote lines, changelog EOFs were normalized, and the second `git diff --check` passed.
- OPTKIT-020 initially referenced nonexistent `repo://optkit/artifact/artifact.go`; the actual package contains `ref.go`, `store.go`, `read.go`, and `helpers.go`. The invalid relation was removed, `artifact/ref.go` was related, and doctor then passed.

### What I learned

- A dry-run validates upload selection and destination but the real upload is the evidence that Pandoc/LaTeX can render the complete guide.
- Generated prose containing blockquotes needs whitespace validation just like source code.
- Keeping implementation tasks open while checking documentation/delivery tasks makes ticket status truthful.

### What was tricky to build

- The uploaded diary necessarily describes the work up to its render time. This final local step records the upload receipt after the PDF has been created; re-uploading with `--force` solely to include its own receipt would overwrite a new document and risk future annotations.
- Cross-ticket commit hashes and ticket-specific upload names had to remain consistent across nine independent workspaces.

### What warrants a second pair of eyes

- Review the proposed APIs and compatibility decisions before implementation; successful document delivery is not architecture acceptance.
- Confirm the reMarkable bundle name and ticket folder are the intended long-term review locations before adding annotations.

### What should be done in the future

- Begin only after upstream entry gates in the OPTKIT-011 roadmap are satisfied.
- During implementation, append new diary steps with exact code commit hashes, failures, commands, and fresh validation evidence.

### Code review instructions

- Start with `index.md`, then read the complete `design-doc/01-*.md`, then this diary.
- Inspect commit `ec784b0d84e47ad45a7429f98482e9678df5dac7` for the documentation batch.
- Validate locally with `docmgr doctor --ticket OPTKIT-018 --stale-after 30`.
- Consult `tasks.md` for the still-open implementation sequence.

### Technical details

```text
bundle: OPTKIT-018 Workbench API Guide.pdf
remote: /ai/2026/08/26/OPTKIT-018
commit: ec784b0d84e47ad45a7429f98482e9678df5dac7
doctor: clean
production code changes: none
```

## Step 3: Print the API plan and project sealed candidate meaning

This step began after OPTKIT-017 closed. It printed the seven-phase API plan and P1 start receipt, then extended historical comparison output with structured candidate intent and labeled mutations loaded from frozen campaign facts.

The projector verifies candidate lineage and loads the exact sealed catalog artifact. It never consults the current registry for historical labels, and direct/full-arm history remains candidate-absent rather than fabricated.

### Prompt Context

**User prompt (verbatim):**

> OPTKIT-016 - OPTKIT-018 in fact, commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill). Print out a brutalist work slip with the plan / different phases for the ticket. then before stsarting a phase, plrint a split about the phase, and print one when the phase is done. budget 2M [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Complete OPTKIT-018 with sealed candidate projections, separate catalog/command APIs, authorization, typed errors, strict HTTP, live composition, handoff types, coherent commits, and every required phase receipt.

**Inferred user intent:** Finish the backend-to-client contract without moving historical reads or mutation semantics into transport/frontend code.

**Commit (RAG-TTC):** `3144759e4aed8f6293d1e7a51869e9f1e73c53f7` — "OPTKIT-018: project sealed candidate meaning"

### What I did

- Printed one plan and P1 start/done slips.
- Added `CandidateSummary` and `MutationSummary` to comparison DTOs.
- Verified candidate baseline/treatment snapshot ancestry.
- Loaded/validated sealed catalog artifacts and semantic identity.
- Labeled normalized mutations only from sealed descriptors.
- Tested candidate presence, absence, lineage mismatch, and missing catalog.

### Why

- Historical descriptions must not change when current registry documentation changes.
- Absence of candidate provenance is a valid historical state and must not be guessed from arm differences.

### What worked

- Comparison projected proposer, strategy, hypothesis, metric, risks, motivation, and exact before/after values.
- The fixture mutation label is `Final result limit` from the sealed catalog.
- Missing sealed catalog produced a corruption error.
- Full tests and lint passed.

### What didn't work

- N/A.

### What I learned

- Campaign spec v3 already provides canonical candidate indexing by treatment arm; event payload decoding is unnecessary for normal historical projection.

### What was tricky to build

- The projector had to verify baseline/treatment ancestry before exposing intent, because an arm ID alone is not sufficient proof that a candidate belongs to the requested comparison.

### What warrants a second pair of eyes

- Review strict no-fallback behavior for sealed descriptor lookup and corruption handling.

### What should be done in the future

- Candidate list/detail views may reuse this summary but should remain read projections.

### Code review instructions

- Review `CandidateSummary`, `Projector.candidateSummary`, and comparison tests.
- Confirm no current registry import appears in specialist projection code.

### Technical details

```text
candidate projection: optional
historical label source: sealed catalog only
missing candidate: nil, no fabrication
P1 start/done: printed true
```

## Step 4: Expose the current catalog with full-ID ETags

This phase added a current catalog application service and the initial separate workbench HTTP package. Catalog and variable routes emit complete native descriptors and lossless domains.

The full catalog ID is a strong cache validator. Conditional GET returns 304 without changing the historical specialist API.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Provide authoring discovery from the current registry under a new route namespace.

**Inferred user intent:** Let clients render legal controls without exposing executable Go bindings.

**Commit (RAG-TTC):** `8c263f753a8c1583f61c7c0180bf5ffc4e3a13ab` — "OPTKIT-018: expose current catalog with ETags"

### What I did

- Added `CatalogService` and full variable lookup.
- Added `workbenchapi` response types and catalog routes.
- Added quoted full-ID ETag and `If-None-Match` handling.
- Preserved native defaults, ranges, choices, docs, probes, and IDs.
- Added invalid/missing variable and 304 tests.
- Printed P2 start/done slips.

### Why

- Current authoring legality and historical sealed provenance are different read models.
- ETag should change for documentation/presentation edits, so it uses full rather than semantic catalog identity.

### What worked

- Catalog returned two ordered sections and complete descriptors.
- RRF variable returned range `0.001…1000` and default `60`.
- Conditional request returned empty 304 with equal ETag.
- Full pre-commit validation passed.

### What didn't work

- N/A.

### What I learned

- Catalog cache policy can permit private revalidation while command responses remain `no-store`.

### What was tricky to build

- `space.Catalog` uses custom JSON serialization and immutable detached sections; HTTP output had to preserve that rather than reconstructing a parallel DTO.

### What warrants a second pair of eyes

- Review ETag/full-ID policy and cache headers.

### What should be done in the future

- Clients should key descriptor caches by full ID and proposal identity by semantic ID.

### Code review instructions

- Review `catalog_service.go`, workbench catalog handlers, and conditional tests.

### Technical details

```text
API version: rag-ttc.workbench-api/v1
ETag source: full catalog ID
conditional status: 304
P2 start/done: printed true
```

## Step 5: Define principals, actions, and typed errors

This phase fixed authorization and error contracts before adding command routes. The application layer now owns validated principals, a closed action set, resources, policy interface, and safe typed errors.

Internal causes remain available to `errors.Is`/operators but are excluded from public messages.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Make policy and failure semantics explicit before HTTP can invoke mutations.

**Inferred user intent:** Prevent transport convenience from creating actor spoofing, string-matched statuses, or secret leakage.

**Commit (RAG-TTC):** `9952f13d3433d5637358ab6307263a0e4b48a345` — "OPTKIT-018: define authorization and error contracts"

### What I did

- Defined six action constants and validation.
- Added `actor:` principal and resource contracts.
- Added `Authorizer` and fixed action-set implementation.
- Added stable application error codes and unwrap support.
- Tested allowed/denied/unauthenticated/canceled policy paths.
- Tested safe message versus secret internal cause.
- Printed P3 start/done slips.

### Why

- HTTP status mapping must branch on typed codes, never substrings.
- Authenticated actor identity must be available before parent/case resolution.

### What worked

- Closed action validation rejects unknown policy names.
- Principal/action mismatch returns typed forbidden.
- Internal public message omits secret cause.
- Full CI passed.

### What didn't work

- N/A.

### What I learned

- Catalog access still benefits from an explicit action even when a deployment grants it broadly.

### What was tricky to build

- Application errors need safe public text and an unwrap chain simultaneously; `ApplicationError` separates `Message` from `Cause`.

### What warrants a second pair of eyes

- Review action granularity before restricted asset variables are implemented.

### What should be done in the future

- Multi-principal production policy may replace the fixed local authorizer without changing command applications.

### Code review instructions

- Review `workbench_contracts.go` and its complete action/error tests.

### Technical details

```text
actions: 6
principal namespace: actor:
status mapping input: ErrorCode
P3 start/done: printed true
```

## Step 6: Implement the transport-independent command service

This phase added campaign parent/case resolution, authorized compile/preview/seal operations, actual deterministic RRF preview execution, proposer binding, and typed idempotency conflict mapping.

The service resolves persisted campaign specs/snapshots and calls the existing compiler/sealer. It contains policy and orchestration but no HTTP parsing.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Build one command application that CLI/HTTP and future UI can share.

**Inferred user intent:** Keep transport handlers thin and prevent business logic duplication.

**Commit (RAG-TTC):** `bf734fda0065fb32c3f767070a258386eae24343` — "OPTKIT-018: add authorized workbench command service"

### What I did

- Added campaign parent and case resolver over stored spec/snapshot artifacts.
- Added current catalog, compile, preview, and seal application methods.
- Added deterministic before/after semantic fixture preview for the registered RRF probe.
- Added authenticated proposer identity binding and spoof rejection.
- Added application composition constructor and completeness validation.
- Added typed candidate idempotency sentinel and mapping.
- Tested compile, preview score/policy change, seal, exact retry, conflict, denial, missing parent, unsupported probe, and stale digest.
- Printed P4 start/done slips.

### Why

- HTTP cannot safely own parent verification, preview dispatch, or idempotency.
- Preview capability names are backend probe contracts rather than component IDs.

### What worked

- RRF preview changed retrieval policy and result scores through actual runtime execution.
- Seal bound proposer to `actor:operator` regardless of absent body identity.
- Retry returned the exact stored envelope and left journal head unchanged.
- Same key with changed compiled draft returned typed idempotency conflict.
- Full CI passed.

### What didn't work

- The first preview test referenced a nonexistent field and failed to compile:

```text
pkg/ttc/experimentworkbench/workbench_service_test.go:87:44: previewPayload.Before.Policy undefined (type search.SearchOutput has no field or method Policy)
pkg/ttc/experimentworkbench/workbench_service_test.go:87:76: previewPayload.After.Policy undefined (type search.SearchOutput has no field or method Policy)
```

`SearchOutput` stores policy provenance at `Identity.RetrievalPolicyID`. I inspected the type, corrected both assertions, and reran focused/full tests.

### What I learned

- Runtime identity is nested under `SearchOutput.Identity`; preview output should preserve the complete native search envelope rather than inventing a reduced policy field.

### What was tricky to build

- Terminal-state retries must succeed by command lookup while new seals must fail state preflight. The sealer and candidate recorder preserve this ordering.

### What warrants a second pair of eyes

- Review parent/case resolver error redaction, preview sensitivity, and actor binding.

### What should be done in the future

- Additional preview probes should implement `PreviewService` dispatch with explicit sensitivity and resource budgets.

### Code review instructions

- Start with `WorkbenchService`, `CampaignResolver`, and `DeterministicPreviewService`.
- Run command-service tests under race.

### Technical details

```text
preview probe: fusion.rrf-contributions/v1
preview schema: schema:rag-ttc.preview.rrf-contributions/v1
preview sensitivity: internal
P4 start/done: printed true
```

## Step 7: Add strict authorized command HTTP adapters

This phase completed the separate workbench route tree. Every route authenticates before application access, command bodies use one bounded strict JSON value, and response/error headers enforce no-store and correlation.

Seal accepts exactly one `Idempotency-Key` header; an `idempotency_key` body field is an unknown-field error.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Adapt command applications to HTTP without moving domain behavior into handlers.

**Inferred user intent:** Publish secure stable browser contracts before frontend implementation.

**Commit (RAG-TTC):** `b675fb6d91670136d94490afa22ef62ba7f2fd1b` — "OPTKIT-018: add strict authorized workbench HTTP API"

### What I did

- Added constant-time bearer authentication.
- Added compile, preview, and seal request/response DTOs/routes.
- Added strict content type, 1 MiB body limit, unknown-field, trailing-value, and cancellation checks.
- Added typed status/error response mapping and internal redaction.
- Added request ID and security/no-store headers.
- Added complete auth/body/diagnostic/spoof/retry/conflict/method/error matrix tests.
- Migrated catalog routes through application authorization.
- Printed P5 start/done slips.

### Why

- Authentication must precede restricted parent/case resolution.
- Body/header duplication creates ambiguous idempotency semantics.

### What worked

- Invalid operator mutation returned HTTP 200 with `sealable:false` and structured diagnostics.
- Missing parent returned typed 404; unsupported probe 422; idempotency conflict 409.
- Wrong/missing token returned 401 without echoing token.
- Internal cause/message were replaced by generic `internal server error`.
- Full CI passed.

### What didn't work

- N/A.

### What I learned

- Request correlation must be generated once in middleware and copied into the request header so downstream error writers reuse the same ID.

### What was tricky to build

- `http.MaxBytesReader` errors, strict JSON errors, application errors, and cancellation need distinct safe responses while preserving one response envelope.

### What warrants a second pair of eyes

- Review bearer-token deployment assumptions, body limit, status mapping, and idempotency header normalization.

### What should be done in the future

- Production authentication can replace the local bearer adapter while retaining `Authenticator`/`Principal` contracts.

### Code review instructions

- Review `workbenchapi/server.go` handlers and `command_test.go` security matrix.

### Technical details

```text
command request media type: application/json
maximum body: 1048576 bytes
seal key source: exactly one Idempotency-Key header
P5 start/done: printed true
```

## Step 8: Compose separate route trees and run a live workflow

This phase updated `campaign serve` to require a secret workbench token and actor, construct the authorized application, and mount specialist/workbench handlers on disjoint paths in one standard-library mux.

A built-binary live smoke created a running campaign, started the real server on a free loopback port, then exercised unauthenticated rejection, catalog, compile, preview, seal, idempotent retry, and specialist cockpit.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Prove one process can host both semantics without merging packages or route responsibilities.

**Inferred user intent:** Deliver an operational backend endpoint for the next frontend ticket.

**Commit (RAG-TTC):** `4c38094ab9edf7b664bd51d5ad90129836852db4` — "OPTKIT-018: compose read and workbench API handlers"

### What I did

- Added required secret token and actor serve settings.
- Added fixed local action policy and application composition.
- Mounted `/api/rag/v1/` and `/api/rag/workbench/v1/` separately.
- Added composition tests preserving GET-only specialist behavior.
- Added a ticket-local running-campaign helper and live server script.
- Ensured server shutdown/temporary-store cleanup with traps.
- Printed P6 start/done slips.

### Why

- Separation is semantic/package-level; one local process remains operationally convenient.
- Live network evidence covers mux, Cobra/Glazed settings, listener, JSON, auth headers, persistence, and shutdown together.

### What worked

- Specialist health returned `read_only:true` without workbench auth.
- Unauthenticated catalog returned 401.
- Authenticated compile/preview/seal succeeded.
- Retry response bytes were identical.
- Candidate proposer was bound to `actor:live-smoke`.
- The live script exited and cleaned the server/store.

### What didn't work

- The first readiness loop printed one expected transient `curl: (7) Failed to connect` before the listener was ready, then the workflow passed. I redirected readiness-probe stderr on subsequent runs while retaining a final mandatory successful health request.

### What I learned

- A readiness poll must suppress only transient attempts; the final health call must still fail the script if the server never starts.

### What was tricky to build

- A seal smoke requires a running campaign. The ticket-local Go helper creates a real checkpointed campaign, closes it, and lets the server reopen the store before commands run.

### What warrants a second pair of eyes

- Review token secrecy in process invocation and whether production deployment should load it from a secret file/environment source.

### What should be done in the future

- Production process management should supply token/actor through secured deployment configuration.

### Code review instructions

- Review `serve.go`, composition test, setup helper, and live smoke script.
- Run the script and confirm no server process remains.

### Technical details

```text
specialist prefix: /api/rag/v1/
workbench prefix: /api/rag/workbench/v1/
live retry bytes: equal
P6 start/done: printed true
```

## Step 9: Publish handoff types, validate, and deliver

This phase exported complete TypeScript contracts and retained sanitized live JSON fixtures for the PBUI/propose implementation. It then ran full backend/frontend validation, evidence audits, docmgr closure, completed PDF delivery, and the final phase receipt.

No React authoring UI was added. The accepted program assigns UI composition to OPTKIT-021–023 after this backend contract gate.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Close backend-to-frontend contracts with generated evidence and full validation, not by implementing the downstream UI ticket.

**Inferred user intent:** Give frontend work authoritative stable types/examples and a proven live endpoint.

**Commit (RAG-TTC):** `f3d42719dd7cb0fef8c75db54295d09ea6320135` — "OPTKIT-018: publish frontend workbench contract types"

### What I did

- Added TypeScript candidate, catalog, domain, pipeline, graph, draft, preview, seal, and error types.
- Exported sanitized live health/catalog/compile/preview/seal/cockpit/error JSON.
- Ran full RAG-TTC lint, tests, vet, build, focused race, and dependency scan.
- Ran specialist TypeScript check, 45 tests, and production build.
- Ran fresh live HTTP smoke and inspected response fields.
- Updated guide, diary, tasks, relations, roadmap, doctor, completed PDF, and slips.

### Why

- Type declarations plus actual live payloads reduce frontend guessing and expose drift immediately.
- Backend completion should not silently expand into the accepted downstream React/PBUI scope.

### What worked

- TypeScript accepted all handoff contracts.
- Nine frontend test files/45 tests passed.
- Production Vite build succeeded.
- Full backend and focused race suites passed.
- Live JSON matched candidate actor, catalog ID, draft digest, preview probe, and read-only specialist health.
- Doctor, upload, and slip audit passed.

### What didn't work

- N/A.

### What I learned

- Sanitized contract fixtures are useful even without UI rendering because they preserve exact nested server output for future client tests.

### What was tricky to build

- TypeScript declarations must represent raw mutation JSON as `unknown`, not guess values as numbers/strings, while still giving exact discriminants for catalog value kinds and preview modes.

### What warrants a second pair of eyes

- Review TypeScript sealed proposal intersections, optional artifact schemas, and future code-generation opportunities.

### What should be done in the future

- OPTKIT-023 should consume these types and fixtures through generated/validated API clients and render the RRF vertical slice.

### Code review instructions

- Review RAG-TTC commits `3144759e` through `f3d42719` in order.
- Run full Go/frontend validation and `scripts/01-live-workbench-api-smoke.sh`.
- Inspect `various/api-contracts/` and run `docmgr doctor --ticket OPTKIT-018 --stale-after 30`.

### Technical details

```text
frontend typecheck: pass
frontend tests: 9 files / 45 tests
frontend build: pass
dependency packages: 373
import cycles: 0
live candidate: candidate:3fc3c902b5d459be5f60ffc40c60ad5941c4c9287cb437afd45e478b93b83314
work slips: 15/15 printed
```
