---
Title: Implementation Diary
Ticket: OPTKIT-007
Status: active
Topics: [implementation, optkit, rag-ttc]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Chronological implementation record for the specialist UI backend."
LastUpdated: 2026-08-26T00:30:00Z
WhatFor: "Preserve prompts, decisions, failures, commits, validation, and frontend handoff evidence."
WhenToUse: "Read before resuming or reviewing OPTKIT-007."
---

# Diary

## Goal

Build the verified, bounded, GET-only RAG-TTC backend that carries durable experiment facts up to the frontend integration boundary.

## Step 1: Freeze the read model and ticket plan

This step created the ticket, split implementation into six explicit phases, mapped existing HTTP composition roots and durable artifact schemas, and froze the first version of the specialist DTOs before projector implementation.

The API remains product-owned and read-only. Its first vertical path is cockpit → arm comparison → paired case → episode pipeline → episode provenance.

### Prompt Context

**User prompt (verbatim):**

```text
Ok, can we now Let's do the server backend side of things, all the way up to the UI, which we'll give to our frontend engineer.

commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)
Print out a brutalist work slip with the plan / different phases for the ticket. then before stsarting a phase, plrint a split about the phase, and print one when the phase is done.
```

**Assistant interpretation:** Implement the complete specialist read backend, characterize it end to end, and leave a precise frontend integration package, with phase slips, commits, and diary updates.

**Inferred user intent:** Give the frontend engineer stable, trustworthy APIs backed by the same durable facts already validated through the CLI.

### What I did

- Created OPTKIT-007 and six phase tasks.
- Printed the overall plan and Phase 0 start slips.
- Mapped RAG-TTC server roots, `http.ServeMux` conventions, Optkit query limits, campaign events, episode results, trajectories, and artifact sensitivity.
- Added versioned cockpit, comparison, case-page, pipeline, provenance, diagnostic, stage, and artifact DTOs.
- Added opaque cursor and bounded case-limit contracts.
- Wrote the initial API/projector design.

### Why

- DTOs and failure semantics must stabilize before multiple projector and HTTP implementations depend on them.
- The browser must not become a journal reducer or trajectory decoder.
- Explicit preview policy is necessary before returning artifact content.

### What worked

- Existing Optkit and RAG-TTC records contain enough information for a retrieval-focused vertical path.
- Go 1.22 method-aware `ServeMux` patterns match repository conventions.
- OPTKIT-006 already persists per-arm graphs and manifest identity in `CampaignSpec`.

### What didn't work

- N/A during initial contract creation.

### What I learned

- The current durable fixture can fully support retrieval pipeline and provenance, but not context/answer/judge stages.
- Direct event payloads are sufficient to find campaign spec, completions, observations, and estimates.
- Pipeline stage payloads are separately referenced from the sealed trajectory and can follow bounded preview policy.

### What was tricky to build

- Missing observations must have explicit `unknown` status and no numeric field. Zero is a valid measured result and cannot represent absence.
- Case pagination needs an opaque token even though the first fixture is small; this prevents the frontend from depending on array offsets.

### What warrants a second pair of eyes

- Review DTO field names before a frontend ships against v1.
- Review whether 4096 bytes is the right initial preview ceiling.
- Confirm first provenance kind should remain episode-only.

### What should be done in the future

- Add context, answer, and judge DTO fields only when durable producers emit them.
- Replace in-memory case slicing with database pagination for large campaigns.

### Code review instructions

- Start with `specialistapi/types.go` and the route list in the design doc.
- Confirm every top-level DTO carries API version, schema, and journal sequence.
- Run `go test ./pkg/ttc/specialistapi -count=1`.

### Technical details

```text
DefaultCaseLimit=50
MaximumCaseLimit=100
MaximumArtifactPreviewSize=4096
API version=rag-ttc.specialist-api/v1
```

## Step 2: Project verified cockpit and comparison views

This step implemented the product reducer that verifies one campaign journal, decodes known RAG-TTC payloads, expands the stored trial, and projects cockpit, comparison, and paginated paired-case views.

Configuration comparison delegates to the existing `optimization.Diff` and `optimization.Plan` authorities. Measurement rows preserve explicit measured, failed, unknown, and inapplicable states.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Turn durable workbench facts into frontend-ready campaign and comparison read models.

**Inferred user intent:** Make the first two UI screens trustworthy before adding HTTP transport.

**Commit (RAG-TTC code):** `346ecdfc9` — "Specialist API: project cockpit and comparisons"

### What I did

- Added a projector over metadata and artifact-store interfaces.
- Verified journals before decoding any product view.
- Decoded campaign spec, completed estimates, observations, and completions.
- Expanded the stored complete-block trial to recover episode-to-arm/case/repeat identity.
- Added cockpit, comparison, and opaque-cursor case-page projectors.
- Added integration tests against an actual manifest-driven six-episode campaign.

### Why

- The frontend must consume one authoritative interpretation rather than independently joining raw events.
- Trial expansion is deterministic and avoids storing redundant browser indexes.
- Missing observations require explicit status rather than a numeric default.

### What worked

- Cockpit reported two arms, three cases, six completions, one paired estimate, budgets, manifest ID, and verified integrity.
- Comparison returned twelve layer diffs, twelve plan steps, three paired cases, and one metric delta.
- Two-page cursor traversal returned 2 + 1 paired rows without duplication.

### What didn't work

- The first two commit attempts failed `golangci-lint` exhaustive-switch checking even after adding a `default` branch:

  ```text
  missing cases in switch of type campaign.EventKind ... (exhaustive)
  ```

  This linter requires an explicit `//nolint:exhaustive` annotation for intentionally partial enum switches. The switch still has a documented default because other event kinds contribute through the generic overview fold rather than product payload decoding.

### What I learned

- A default branch does not satisfy this repository's exhaustive linter configuration.
- The stored complete-block trial is enough to recover stable pair order and episode IDs.
- `CampaignSpec.ConfigGraphs` makes UI configuration provenance possible without reading CLI manifests.

### What was tricky to build

- Maps make fact lookup efficient but not deterministic for display. Paired rows are ordered by stored dataset cases and repeat number rather than map iteration.
- Availability is true only when both selected episodes have completions; each episode ID remains independently exposed.

### What warrants a second pair of eyes

- Review whether multiple observations per episode require construct-keyed storage when Judgekit epochs enter this projector.
- Review the current one-primary-metric comparison shape before answer metrics are added.

### What should be done in the future

- Key observations by episode plus construct/epoch when multi-instrument campaigns arrive.
- Move large case scans to a database projector.

### Code review instructions

- Start with `projector.go:loadFacts`, then `Cockpit`, `Comparison`, and `pairedCases`.
- Run `go test ./pkg/ttc/specialistapi -count=1`.

### Technical details

```text
cockpit arms=2 cases=3 completed=6
comparison pairs=3 metric_pairs=3
case pagination=2 + 1
```

## Step 3: Project pipeline and provenance

This step connected paired cases to the actual sealed episode trajectories. Pipeline views expose only observed retrieval stages, while provenance views connect manifest, graph, snapshot, episode, completion, result, trajectory, output, and event payload artifacts.

Artifact previews fail closed by sensitivity and size. A reference remains visible even when its bytes are not.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Implement the detailed screens behind a paired case without leaking raw filesystem access or sensitive payloads.

**Inferred user intent:** Give the frontend enough evidence to explain where results came from and where arms diverged.

**Commit (RAG-TTC code):** `dbce24f3c` — "Specialist API: project pipeline provenance"

### What I did

- Loaded and verified sealed trajectories with `episode.LoadTrajectory`.
- Decoded `retrieval.stage` payloads into ordered stage summaries.
- Added bounded output and stage previews.
- Added episode provenance and explicit artifact edges.
- Added pipeline, provenance, sensitivity, and size-policy tests.

### Why

- The browser should render a stage rail, not decode trajectory event schemas.
- Artifact refs and unavailable reasons preserve provenance without violating policy.

### What worked

- Treatment episodes produced observed stage rails and output refs.
- Every stage payload preview was available for the deterministic fixture.
- The complete retrieval output exceeded the 4096-byte preview limit and was correctly represented with `size_limit`.
- Provenance linked the persisted arm graph and immutable snapshot to episode artifacts.

### What didn't work

- The first compile failed because `completion` was loaded but unused in the pipeline projector:

  ```text
  pipeline.go:21:8: declared and not used: completion
  ```

  Pipeline now discards that return value while provenance uses it.

- The first test assumed the complete output would always fit the preview ceiling. It failed when the real fixture output exceeded 4096 bytes. The test now accepts either an available bounded preview or the explicit `size_limit` result.

- The first commit failed the repository's `nonamedreturns` linter because `episodeArtifacts` used five named return values. It was rewritten with unnamed returns and explicit zero-value error returns.

### What I learned

- Tests must enforce policy outcomes rather than assume fixture payload size.
- Trajectory verification already checks every event payload ref, making it the correct pipeline loading primitive.

### What was tricky to build

- The trajectory artifact can be too large to preview while each stage payload remains small enough. Preview policy therefore applies independently to every ref.
- Provenance edges need stable semantic roles instead of frontend inference from schemas.

### What warrants a second pair of eyes

- Review whether chunk IDs should be sensitivity-filtered for non-public corpora.
- Review future nested artifact traversal before claiming full transitive custody.

### What should be done in the future

- Add context/answer/judge stage decoders only with their durable campaigns.
- Add authorized preview policy if confidential local workflows require it.

### Code review instructions

- Review `pipeline.go`, especially `episodeArtifacts` and `preview`.
- Validate an oversized output still returns metadata and `size_limit`.

### Technical details

```text
preview ceiling=4096 bytes
allowed sensitivity=public|internal
blocked sensitivity=confidential|restricted
```

## Step 4: Expose the GET-only HTTP API

This step mounted the projectors behind Go 1.22 method-aware `http.ServeMux` patterns and added a long-running Glazed `BareCommand` for operators. The server validates all IDs, arm selectors, limits, and cursors before projection.

Successful responses carry JSON content type, ETag, no-store, and security headers. Mutation methods are absent and receive 405 from the mux.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Carry the stable projectors to the exact transport boundary the frontend engineer will use.

**Inferred user intent:** Deliver a runnable backend rather than only library APIs.

**Commit (RAG-TTC code):** `59898ae49` — "Specialist API: expose read-only HTTP workflow"

### What I did

- Added health, cockpit, comparison, cases, pipeline, and provenance GET routes.
- Added JSON error envelopes and status mapping.
- Added ETag, no-store, nosniff, frame-deny, and no-referrer headers.
- Added strict pagination and provenance-kind checks.
- Added `campaign serve --store --listen`.
- Added HTTP workflow, invalid-input, header, and no-mutation tests.

### Why

- HTTP should transport projector results, not duplicate joins or decoding.
- A same-origin, loopback-first server avoids an unnecessary wildcard CORS policy.
- The operator CLI is the natural place to start the read backend against an existing store.

### What worked

- The complete in-memory HTTP path passed cockpit through provenance.
- POST to a GET route returned 405.
- Invalid IDs, limits, cursors, missing arm queries, and provenance kinds returned expected 400/404 envelopes.

### What didn't work

- The first commit used raw Cobra string flags for the server command. Glazed lint rejected both:

  ```text
  define CLI flags with cmds.WithFlags(fields.New(...)) instead of raw Cobra/pflag/flag APIs
  ```

  The server was converted to a `cmds.BareCommand` with a `CommandDescription`, typed settings, and `cli.BuildCobraCommandFromCommand`. It correctly receives no structured-output flags because it is a long-running server rather than a row producer.

### What I learned

- Glazed BareCommand is the correct abstraction for lifecycle commands that do not emit rows.
- Method-aware ServeMux provides the desired 405 behavior without custom mutation handlers.

### What was tricky to build

- The command must keep the local profile open for the server lifetime and close it after cancellation.
- Error classification distinguishes malformed requests from missing scientific entities without leaking internal paths.

### What warrants a second pair of eyes

- Review deployment authentication before binding beyond loopback.
- Review whether ETag/no-store is the desired first frontend caching combination.

### What should be done in the future

- Add authentication and authorization before remote deployment.
- Add optional specialist SSE after frontend polling behavior is measured.

### Code review instructions

- Review `http.go` route registration first.
- Run HTTP tests and inspect POST behavior.
- Build the CLI and inspect `campaign serve --help`.

### Technical details

```text
rag-ttc experiment optkit-rag campaign serve \
  --store /tmp/rag-specialist \
  --listen 127.0.0.1:8090
```

## Step 5: Characterize the live frontend workflow

This step exercised the actual built binary and live TCP server rather than only `httptest`. It created a fresh durable campaign, started the API, traversed every route, followed a treatment episode into pipeline and provenance, tested pagination, and archived all six frontend response fixtures.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Prove the backend integration exactly as the frontend will consume it.

**Inferred user intent:** Hand the frontend engineer real, validated fixtures and startup commands.

### What I did

- Added `scripts/01-validate-specialist-api.sh`.
- Built the CLI and ran focused unit and race tests.
- Created and served a six-episode campaign from the canonical manifest.
- Validated health, cockpit, comparison, two case pages, pipeline, provenance, 405 mutation behavior, bad limits, ETag, and no-store.
- Archived six real JSON responses under `sources/api-fixtures/`.
- Wrote `reference/02-frontend-engineer-api-handoff.md` with TypeScript sketches and screen acceptance guidance.

### Why

- Live serving catches command wiring, profile lifetime, sockets, and URL behavior absent from projector tests.
- Real JSON fixtures remove guesswork from frontend component development.

### What worked

- All workflow markers passed.
- The stable episode ID was followed from comparison to pipeline and provenance.
- Cursor pagination returned 2 + 1 cases.
- Security and read-only behavior matched the contract.

### What didn't work

- The first readiness poll printed one expected connection-refused line while the server was starting:

  ```text
  curl: (7) Failed to connect to 127.0.0.1 ... Couldn't connect to server
  ```

  The retry loop was correct; stderr is now suppressed during readiness only, and real API calls retain normal failure output.

### What I learned

- A random campaign ID is acceptable in fixtures because the frontend must treat IDs as opaque.
- Episode identities remain deterministic for the same trial and snapshots, which makes cross-fixture links stable.

### What was tricky to build

- The validator allocates an ephemeral loopback port, tracks the server PID, and guarantees cleanup on every exit.
- Response fixtures must be generated by the exact live transport and not reconstructed from Go structs.

### What warrants a second pair of eyes

- Review TypeScript optional fields against the archived JSON.
- Review whether Storybook should normalize fixture IDs for readability or preserve real opaque values.

### What should be done in the future

- Frontend should use these fixtures for MSW and component tests.
- Add contract generation only if hand-maintained TypeScript starts drifting.

### Code review instructions

- Run the ticket validator.
- Open all six fixtures and follow IDs across them.
- Read the frontend handoff before starting components.

### Technical details

```text
SPECIALIST_HEALTH=PASS
COCKPIT_API=PASS
COMPARISON_API=PASS
CASE_PAGINATION=PASS
PIPELINE_API=PASS
PROVENANCE_API=PASS
READ_ONLY_AND_ERRORS=PASS
OPTKIT_007_SPECIALIST_API_VALIDATION=PASS
```

## Step 6: Audit and hand off the backend

The final step audited generated fixtures for secrets and local paths, reconciled the frontend guide with the live API, reran the complete validator, and prepared the closed ticket as the frontend engineer's source of truth.

The backend ends at a deliberate boundary: it serves verified scientific read models and real fixtures, but it does not include frontend component or styling decisions.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Complete the server-side vertical slice and leave a safe, testable frontend contract.

**Inferred user intent:** Allow frontend work to begin independently without reopening backend semantics.

### What I did

- Audited all six fixtures for credentials, tokens, passwords, and local filesystem paths.
- Confirmed fixture sizes remain bounded and practical for frontend tests.
- Verified every implementation commit passed full repository tests, lint, Glazed vet, and build hooks.
- Ran focused specialist race tests through the live validator.
- Completed docmgr relations, tasks, changelog, doctor, and closure preparation.

### Why

- Frontend fixtures become copied test assets and must not contain machine-local or sensitive data.
- The handoff is only complete when code, docs, fixtures, and runtime commands agree.

### What worked

- No secrets or `/home/...` paths were present in API fixtures.
- Fixture sizes ranged from 772 to 12,321 bytes.
- The final validator remained green.
- The unrelated Optkit SQLite worktree edit remained untouched.

### What didn't work

- The first final `docmgr doctor` reported `frontend` as an unknown topic; the established vocabulary uses `ui`. The handoff topic was changed to `ui` rather than expanding vocabulary with a synonym.
- `git diff --check` found one extra blank line at the end of this diary. The file was normalized to exactly one terminating newline before the final commit.

### What I learned

- Ticket topic vocabulary should reuse `ui` for frontend-facing work.
- The API is sufficient for the retrieval-focused cockpit-to-provenance UI.
- Rich answer and judge screens remain correctly blocked by missing campaign producers rather than backend transport.

### What was tricky to build

- The handoff needed to state optional numeric semantics precisely enough that TypeScript code does not collapse unknown into zero.
- Real fixtures contain opaque random campaign IDs; preserving them reinforces correct client behavior.

### What warrants a second pair of eyes

- Security review is required before exposing the server beyond loopback.
- Frontend should verify accessibility, CSP, and same-origin proxy behavior in its own ticket.

### What should be done in the future

- Implement the frontend against the archived fixtures.
- Add authenticated deployment configuration before any remote use.
- Add database-side case pagination before large production campaigns.

### Code review instructions

- Run the full validator and follow its campaign/episode IDs.
- Compare every frontend type sketch with its corresponding JSON fixture.
- Start backend review at `specialistapi/http.go`, then follow calls into projectors.

### Technical details

```text
RAG-TTC commits:
202d63a29 Specialist API: freeze read model contracts
346ecdfc9 Specialist API: project cockpit and comparisons
dbce24f3c Specialist API: project pipeline provenance
59898ae49 Specialist API: expose read-only HTTP workflow
```
