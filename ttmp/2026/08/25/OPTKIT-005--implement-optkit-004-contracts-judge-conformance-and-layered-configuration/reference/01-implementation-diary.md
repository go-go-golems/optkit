---
Title: Implementation Diary
Ticket: OPTKIT-005
Status: active
Topics:
    - optkit
    - rag-ttc
    - implementation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/budget/types.go
      Note: Defines enforced resource accounting contrasted with scheduler claims
    - Path: repo://optkit/campaign/event.go
      Note: Defines control-event identity and unused causation/correlation fields
    - Path: repo://optkit/docs/01-optkit-records-artifacts-and-control-model.md
      Note: Public Glazed-style field ownership and lifecycle reference
    - Path: repo://optkit/scheduler/types.go
      Note: Defines work, lease, and removal-candidate fields documented by the audit
    - Path: repo://optkit/ttmp/2026/08/25/OPTKIT-005--implement-optkit-004-contracts-judge-conformance-and-layered-configuration/scripts/01-validate-release.sh
      Note: Reproduces Phase R isolated validation
    - Path: repo://rag-ttc/go.mod
      Note: Pins public Optkit and Judgekit pseudo-versions for isolated builds
    - Path: repo://rag-ttc/go.sum
      Note: Records checksums for published integration commits
ExternalSources: []
Summary: Chronological implementation diary for OPTKIT-004 Phases 0-2.
LastUpdated: 2026-08-25T16:55:00-04:00
WhatFor: Record implementation decisions, commands, failures, commits, review risks, and continuation guidance.
WhenToUse: Read before resuming OPTKIT-005 or reviewing its implementation.
---



# Diary

## Goal

Record the evidence-first implementation of the release baseline and OPTKIT-004 Phases 0-2: contract freeze, sealed-answer Judgekit conformance, layered configuration identity, and deterministic invalidation planning.

## Step 1: Open the implementation ticket and freeze the work plan

This step turned the accepted sequencing discussion into a bounded implementation ticket. The plan separates operational release stabilization from the three OPTKIT-004 phases and gives each phase an explicit acceptance gate before projector or browser work begins.

The ticket deliberately preserves the unrelated Optkit SQLite worktree edit and treats the existing Judgekit adapter as an implementation to audit rather than replace.

### Prompt Context

**User prompt (verbatim):** "Ok, go. commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill) . Print out a brutalist work slip with the plan / different phases for the ticket. then before stsarting a phase, plrint a split about the phase, and print one when the phase is done."

**Assistant interpretation:** Start implementing the agreed OPTKIT-004 prerequisites, create a tracked ticket and strict diary, commit coherent phases, and print plan/start/completion thermal slips throughout the work.

**Inferred user intent:** Execute the work as an auditable multi-phase engineering project with physical progress markers and continuation-quality documentation.

### What I did

- Read the full brutalist work-slip instructions and diary reference.
- Inspected the Optkit, RAG-TTC, Judgekit, and RagKit worktrees.
- Created `OPTKIT-005` with a design document and implementation diary.
- Defined release stabilization, Phase 0, Phase 1, Phase 2, and final handoff gates.
- Explicitly excluded the unrelated `optkit/store/sqlite/rows.go` modification.

### Why

- A new implementation ticket avoids reopening the completed OPTKIT-004 design ticket.
- Explicit phase gates keep UI contracts from depending on unstable or inferred event shapes.
- Recording unrelated worktree state before the first commit prevents accidental inclusion.

### What worked

- `docmgr status --summary-only` located the expected Optkit ticket root and configuration.
- Ticket and document creation completed successfully.
- RAG-TTC, Judgekit, and RagKit began clean; only the known Optkit SQLite edit was present.

### What didn't work

- N/A.

### What I learned

- The existing work naturally separates into one operational phase and three semantic phases.
- Phase 1 should be primarily a conformance and evidence task because OPTKIT-002 P7 already implemented its behavioral requirements.

### What was tricky to build

- The main planning edge was distinguishing the closed OPTKIT-004 design deliverable from new implementation work. The solution was a separate OPTKIT-005 ticket that references, but does not mutate the status of, OPTKIT-004.
- The Optkit worktree is not clean because of a known unrelated edit. Every future staging command must name exact ticket or code paths rather than using broad staging.

### What warrants a second pair of eyes

- Confirm that Phase 0 freezes only contracts needed by the first specialist UI and does not prematurely encode future bundle-build behavior.
- Confirm that release stabilization does not rewrite unpublished dependency history or absorb unrelated local work.

### What should be done in the future

- Every phase must add its exact commands, failures, commit hashes, review instructions, and validation evidence here before being marked complete.

### Code review instructions

- Start with the implementation plan in `design-doc/01-optkit-004-phases-0-2-implementation-plan.md`.
- Compare its phase gates with Section 25 of the OPTKIT-004 intern guide.
- Run `git -C optkit status --short` and verify `store/sqlite/rows.go` remains unstaged and unrelated.

### Technical details

Initial worktree state:

```text
optkit:  M store/sqlite/rows.go
rag-ttc: clean
judgekit: clean
ragkit: clean
```

## Step 2: Replace workspace-only dependencies with published commits

This step removed the synthetic `v0.0.0` dependency boundary that prevented RAG-TTC from building outside the parent workspace. Optkit and Judgekit integration commits are now reachable from public task branches, and RAG-TTC pins reproducible pseudo-versions with complete module sums.

The final validator tests the committed Optkit tree in a detached worktree so the unrelated malformed SQLite edit remains untouched. RAG-TTC now passes its normal pre-push test, lint, and Glazed-vet gates with `GOWORK=off`.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Complete the release/isolation phase before implementing OPTKIT-004 semantic contracts.

**Inferred user intent:** Ensure subsequent UI-enabling work rests on independently consumable modules rather than a fragile local workspace.

**Commit (code):** `2930b5fb7ca62bda2f51e2f46cf948bd1a9f4616` — "Build: pin published Optkit integration commits"

### What I did

- Inventoried `go.work`, module requirements, tags, remotes, branches, and GitHub repository visibility.
- Pushed `task/use-optkit` to `go-go-golems/optkit` and `go-go-golems/judgekit`.
- Resolved public pseudo-versions with `go list -m -json`.
- Replaced RAG-TTC's synthetic Judgekit and Optkit `v0.0.0` requirements with those pseudo-versions.
- Ran `GOWORK=off go mod tidy` and committed `go.mod` plus `go.sum`.
- Removed the obsolete exact-version replacements from the untracked parent `go.work`.
- Added `scripts/01-validate-release.sh` and archived its successful output in `sources/01-phase-r-validation.txt`.
- Pushed the RAG-TTC task branch after its pre-push hooks passed.

### Why

- The previous module graph worked only because `go.work` replaced nonexistent `v0.0.0` revisions.
- Public pseudo-versions preserve exact attribution without prematurely tagging an unreleased `v1` API.
- A detached Optkit worktree validates committed source without modifying or hiding unrelated user work.

### What worked

- Judgekit isolated unit and race tests passed.
- The clean committed Optkit tree passed all package tests and vet with `GOWORK=off`.
- RAG-TTC passed full tests, build, vet, golangci-lint, focused race tests, and Glazed vet without workspace replacements.
- The reproducible validator ended with `PHASE_R_VALIDATION=PASS`.
- The RAG-TTC pre-push hook passed and published `task/use-optkit`.

### What didn't work

- Two initial test commands were accidentally run from the workspace root instead of module directories:

  ```text
  # ./...
  pattern ./...: directory prefix . does not contain main module or its selected dependencies
  FAIL ./... [setup failed]
  ```

- Running Optkit tests in the active worktree exposed the unrelated malformed local edit exactly as expected:

  ```text
  # github.com/go-go-golems/optkit/store/sqlite
  store/sqlite/rows.go:59:1: syntax error: non-declaration statement outside function body
  ```

  The validator was changed to test a detached committed worktree rather than touching that file.

- The first RAG-TTC push failed because its pre-push hook correctly used `GOWORK=off` while `go.mod` still required unpublished `v0.0.0` modules. Representative errors were:

  ```text
  missing go.sum entry for module providing package github.com/go-go-golems/optkit/record
  missing go.sum entry for module providing package github.com/go-go-golems/judgekit/assessment
  ```

  Pinning public pseudo-versions and running `GOWORK=off go mod tidy` resolved the failure.

### What I learned

- A branch-reachable pseudo-version is sufficient for an isolated, reproducible integration baseline; a semantic release tag can wait for API review.
- Pre-push hooks provided the correct release test: they immediately exposed assumptions hidden by the workspace.
- Clean detached worktrees are the right validation mechanism when unrelated local edits must remain visible and untouched.

### What was tricky to build

- The active Optkit worktree cannot compile due to unrelated malformed content, but reverting or stashing it would violate the preservation constraint. The validator therefore creates a detached worktree at committed `HEAD`, runs isolated checks, and removes it through a trap.
- `go.work` replacements were version-specific (`v0.0.0`), so changing RAG-TTC to pseudo-versions made them ineffective before they were removed. This prevented accidental continued use of local modules during validation.

### What warrants a second pair of eyes

- Confirm that using task-branch pseudo-versions is the desired pre-release policy; formal tags may be preferable after API review.
- Review `rag-ttc/go.mod` and `go.sum` to verify only Judgekit and Optkit dependency identities changed.
- Confirm CI has credentials/network access compatible with public GitHub module resolution.

### What should be done in the future

- Replace pseudo-versions with reviewed semantic release tags when Optkit and Judgekit establish their first public release contract.
- Keep all isolated-module validation independent of the parent `go.work`.

### Code review instructions

- Start at `rag-ttc/go.mod` and verify both pseudo-versions resolve to the documented commits.
- Run `scripts/01-validate-release.sh` from any working directory.
- Inspect `sources/01-phase-r-validation.txt` for `PHASE_R_VALIDATION=PASS`.
- Verify `git -C optkit status --short` still reports only the unrelated `store/sqlite/rows.go` edit plus ticket documentation awaiting commit.

### Technical details

Pinned dependency identities:

```text
Optkit   v0.0.0-20260825205238-f92264c48eaa
Judgekit v0.0.0-20260825195635-870fd3bb8391
RAG-TTC  2930b5fb7ca62bda2f51e2f46cf948bd1a9f4616
```

Validation commands include:

```bash
GOWORK=off go test ./... -count=1
GOWORK=off go build ./...
GOWORK=off go vet ./...
GOWORK=off golangci-lint run
GOWORK=off go test -race ./pkg/ttc/... ./internal/admin/chatserver/... -count=1
```

## Step 3: Publish the field-level Optkit model and relocate bootstrap history

This step converted the scheduler-field audit into permanent product documentation. The new Glazed-style help entry covers every public layer from record identities and artifact custody through typed configuration, episodes, measurements, experiments, campaign control, leases, budgets, projections, and the query API. Every record table explains purpose, writer, reader, enforcement status, and failure behavior.

The old ADR, implementation diary, and placeholder component-journal directory were not product documentation. They now live under this ticket's `sources/legacy-bootstrap-docs/` tree so the historical record remains attributable without presenting implementation chronology as current API guidance.

### Prompt Context

**User prompt (verbatim):** "ok, ths is great, actually write a full document in optkit about all these fields and structs (not just the ones above). DOcument the ones that are slated for removal as well, but clearly mark why they are removed. Store it in optkit/docs, use the glazed writing help entries standard. Remove the adr / journals / diary (move those to ttmp, they are leftovers from our first implementer who didn't have ttmp / docmgr)"

**Assistant interpretation:** Create a comprehensive public Optkit data-model reference using Glazed help-page conventions, clearly register pre-release removals, and relocate legacy implementation-history documents into docmgr-managed ticket storage.

**Inferred user intent:** Replace stale implementation-era documentation with a durable, discoverable, field-level contract that lets maintainers distinguish enforced behavior from unused declarations before freezing UI schemas.

**Commit (docs):** `367d7b06081c7271a820b3def855148a895bbcc2` — "Docs: publish Optkit control model reference"

### What I did

- Read `glaze help how-to-write-good-documentation-pages` and `glaze help writing-help-entries` in full.
- Inventoried exported types across `record`, `artifact`, `space`, `system`, `episode`, `measure`, `experiment`, `campaign`, `scheduler`, `budget`, `projection`, `query`, `local`, and SQLite composition.
- Authored `docs/01-optkit-records-artifacts-and-control-model.md` with Glazed frontmatter, concept-first sections, writer/reader tables, failure modes, troubleshooting, and cross-references.
- Added an explicit removal register for scheduler resource claims, the unwired heartbeat method, scheduler failure evidence, campaign causation/correlation fields, matching query fields, and `CorrelationID`.
- Moved the legacy ADR, 741-line bootstrap diary, and component-journal placeholder into `sources/legacy-bootstrap-docs/` with Git history preserved.
- Updated `README.md` to link the public model reference and direct historical engineering records to `ttmp/`.
- Validated exact Glazed frontmatter keys, `GeneralTopic`, unique intended slug, absence of a Markdown H1, troubleshooting, and See Also sections.

### Why

- Public documentation should describe current contracts and operational ownership, not preserve chronological implementation notes at the product-doc root.
- Writer/reader/enforcement columns make declared-only fields visible before they become accidental compatibility commitments.
- The ticket source tree is the correct custody location for historical ADRs and diaries because it retains provenance, chronology, and review context.

### What worked

- The new document contains 728 lines and 7,532 words.
- Frontmatter validation ended with `GLAZED_FRONTMATTER=PASS`.
- The only remaining references to old `docs/adr`, `docs/journals`, or `docs/implementation-diary` paths are preserved historical evidence under `ttmp/`.
- Git recognized all three legacy files as renames rather than delete-and-recreate operations.

### What didn't work

- N/A. The documentation move and format validation succeeded on the first attempt.

### What I learned

- The strongest cleanup candidates are not merely unused by RAG-TTC; they have writers and storage but no behavioral reader. That distinction identified `ResourceClaims` as more misleading than implemented Numbergame-only paths such as queue failure and budget release.
- Glazed help pages intentionally omit a top-level Markdown heading because the renderer supplies the frontmatter title.
- Optkit's current CLI uses the standard `flag` package rather than Cobra/Glazed, so this step follows Glazed's document contract without adding a large command-framework migration to a field-reference task.

### What was tricky to build

- “All fields and structs” spans several authorities that intentionally share vocabulary. The document had to distinguish episode failure evidence from scheduler failure evidence, resource usage from budget claims, and lease authority from campaign facts rather than grouping fields by similar names.
- Causation and correlation are unused but participate in the current control-event identity shape. The removal register therefore marks a journal-version consequence instead of presenting deletion as a harmless struct edit.
- The historical diary contains stale paths by design. Editing them after relocation would rewrite source evidence, so they remain verbatim under the ticket source tree.

### What warrants a second pair of eyes

- Review the removal register before code deletion, especially the event-identity implications of removing causation and correlation.
- Confirm whether the Optkit CLI should later migrate to Cobra/Glazed and load this page at runtime; that migration is intentionally outside this documentation commit.
- Verify the writer/reader classifications for Numbergame-only behavior remain accurate if the example is reduced later.

### What should be done in the future

- Update the document in the same commit whenever exported record fields or ownership semantics change.
- Add runtime Glazed help integration only as a dedicated CLI migration with command-parity tests.
- Preserve the moved bootstrap documents as historical sources; do not copy them back into `docs/`.

### Code review instructions

- Start with `docs/01-optkit-records-artifacts-and-control-model.md`, especially “Durable work scheduling,” “Campaign commands, events, and state,” and “Removal register.”
- Compare tables against `scheduler/types.go`, `scheduler/queue.go`, `campaign/event.go`, `budget/types.go`, and `query/service.go`.
- Verify `docs/` contains only current public help entries.
- Run the frontmatter validation command recorded below and `docmgr doctor --ticket OPTKIT-005 --stale-after 30`.

### Technical details

Glazed document validation:

```bash
python3 - <<'PY'
from pathlib import Path
import yaml
text = Path('docs/01-optkit-records-artifacts-and-control-model.md').read_text()
_, frontmatter, body = text.split('---\\n', 2)
data = yaml.safe_load(frontmatter)
assert data['SectionType'] == 'GeneralTopic'
assert not any(line.startswith('# ') for line in body.splitlines())
assert '## Troubleshooting' in body
assert '## See Also' in body
print('GLAZED_FRONTMATTER=PASS')
PY
```

Legacy source destination:

```text
ttmp/2026/08/25/OPTKIT-005--implement-optkit-004-contracts-judge-conformance-and-layered-configuration/
  sources/legacy-bootstrap-docs/
    adr/0001-clean-slate-local-vertical-slice.md
    implementation-diary.md
    journals/README.md
```

## Step 4: Remove declared-only scheduler and event fields

This step converted the documented removal register into a clean pre-v0 API. Scheduler resource claims, the unwired heartbeat method, scheduler-level failure evidence, and unused campaign causation/correlation fields no longer imply capabilities that no worker, policy, projector, or UI actually implements.

Lease fencing, expiry reclamation, attempts, retryable failures, episode evidence, and enforced budget reservations remain intact. RAG-TTC now consumes the lean work-item constructor through a published Optkit pseudo-version, so isolated module builds continue to prove the real dependency boundary.

### Prompt Context

**User prompt (verbatim):** "Ok, do it."

**Follow-up user prompt (verbatim):** "commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)"

**Assistant interpretation:** Delete every approved item in the removal register, migrate active callers, validate behavior at repository boundaries, commit semantic units, and document the exact result.

**Inferred user intent:** Freeze a smaller honest Optkit API before Phase 0 schemas and the UI make unused implementation placeholders expensive to remove.

**Commit (Optkit code):** `fda4ad62c97de2db4d85eab78aea926a7d474f5f` — "Scheduler: remove unenforced control fields"

**Commit (RAG-TTC caller):** `ebb6cbc63f37cd0090ccdd1661a7e512640b21db` — "Campaign: adopt lean Optkit work contract"

### What I did

- Removed `scheduler.ResourceClaim`, `WorkItem.ResourceClaims`, and the claims parameter from `NewWorkItem`.
- Removed `resource_claims_json` from the new SQLite work-item schema, inserts, selects, and decoders.
- Preserved enforced `budget.Quantity` reservations and episode actual-usage commitment.
- Removed `Queue.Heartbeat` and the unused SQLite heartbeat implementation.
- Removed `scheduler.WorkFailure.Evidence` while preserving `episode.Failure.Evidence`.
- Removed campaign causation and correlation from `NewEvent`, `ControlEvent`, event identity, SQLite storage, query views, and `record.CorrelationID`.
- Migrated Numbergame and RAG-TTC work-item constructors.
- Published Optkit commit `fda4ad6` and pinned RAG-TTC to pseudo-version `v0.0.0-20260825222033-fda4ad62c97d`.
- Updated the public model document from “slated” to a completed removal record with reintroduction gates.
- Added `scripts/03-validate-control-cleanup.sh` and archived exact output in `sources/03-control-cleanup-validation.txt`.

### Why

- Resource claims duplicated budget quantities but had no capacity-aware lease reader.
- Heartbeat had no worker loop and overstated multi-worker liveness guarantees.
- Scheduler failure evidence duplicated the episode evidence authority.
- Causation and correlation had storage and query fields but no producer or consumer semantics.
- Removing these contracts before UI schema freeze avoids empty fields, fake controls, and permanent compatibility obligations.

### What worked

- The clean committed Optkit tree passed every package test and `go vet` under `GOWORK=off`.
- The Numbergame demo completed against a fresh post-cleanup SQLite schema.
- RAG-TTC focused tests and race tests passed against the published Optkit pseudo-version.
- RAG-TTC's full pre-commit and pre-push test, golangci-lint, and Glazed-vet hooks passed.
- Source guards found none of the removed identifiers in active Optkit Go or SQL source.
- The validator ended with `CONTROL_CLEANUP_VALIDATION=PASS`.

### What didn't work

- N/A. The source migration and validation passed without a failed implementation attempt.

### What I learned

- Nil `causation` and `correlation` fields used `omitempty` in the event identity, so removing them preserves canonical identity for existing events that never populated those fields. Only hypothetical pre-freeze events with nonempty values require their original verifier.
- Removing a `NOT NULL` SQLite column from new schema code creates an intentional pre-freeze storage break for old databases: an old `resource_claims_json` column without a default prevents new inserts. Existing stores must be reset rather than receiving a compatibility write shim.
- Publishing Optkit before updating RAG-TTC preserved isolated-module validation and prevented the workspace from hiding an API mismatch.

### What was tricky to build

- The active Optkit worktree still contains an unrelated malformed `store/sqlite/rows.go`. Validation therefore applied the cleanup diff to a detached clean worktree before the code commit and validated committed `HEAD` afterward.
- The SQLite schema had to remove both logical fields and active SQL columns. Leaving the old work-item column in current creation SQL would have retained the misleading contract even after deleting the Go type.
- RAG-TTC already had unrelated Phase 0 files in progress. Staging named only `go.mod`, `go.sum`, and `pkg/ttc/optkitcampaign/campaign.go`, preserving the optimization fixture and stage constants for their own phase commit.

### What warrants a second pair of eyes

- Confirm that resetting pre-freeze local SQLite stores is acceptable and that no production store depends on `resource_claims_json`.
- Review `campaign.eventIdentity` and the historical-journal note for nonempty removed metadata.
- Confirm Phase 3 reintroduces heartbeat only with worker cadence, cancellation, reclaim races, and stale-completion tests.
- Verify RAG-TTC should continue returning worker errors directly until its queue-failure integration is implemented.

### What should be done in the future

- Add capacity claims only with atomic capacity-aware leasing and worker capacity declarations.
- Add heartbeat only as a complete long-running worker protocol.
- Add causation or correlation only when a writer establishes validated edges and a projector or operator workflow consumes them.
- Keep the public removal register synchronized with any reintroduction decision.

### Code review instructions

- Start at `scheduler/types.go` and `scheduler/queue.go`, then inspect SQLite removal in `store/sqlite/queue.go` and `store/sqlite/store.go`.
- Review campaign identity changes in `campaign/event.go`, SQLite journal changes in `store/sqlite/journal.go`, and query removal in `query/service.go`.
- Review RAG-TTC's two `NewWorkItem` calls and pinned Optkit pseudo-version.
- Run `scripts/03-validate-control-cleanup.sh` and check for `CONTROL_CLEANUP_VALIDATION=PASS`.

### Technical details

Removed active identifiers:

```text
ResourceClaim
ResourceClaims
resource_claims_json
Heartbeat
WorkFailure.Evidence
NewEvent.Causation
ControlEvent.Causation
NewEvent.Correlation
ControlEvent.Correlation
EventView.Causation
EventView.Correlation
CorrelationID
```

Published dependency:

```text
github.com/go-go-golems/optkit v0.0.0-20260825222033-fda4ad62c97d
```

Validation evidence:

```text
OPTKIT_COMMIT=fda4ad62c97de2db4d85eab78aea926a7d474f5f
RAG_TTC_COMMIT=ebb6cbc63f37cd0090ccdd1661a7e512640b21db
CONTROL_CLEANUP_VALIDATION=PASS
```

## Step 5: Freeze the cross-layer optimization contracts

This step completed OPTKIT-004 Phase 0 by adding one strict, digest-pinned optimization fixture over the existing cross-product retrieval fixture. The new contract links all configuration layers to stable retrieval, context, answer, and judge stage schemas, then carries one admitted evidence set through context construction, sealed answer lineage, two measurement epochs, observations, hard constraints, and primary metrics.

Production retrieval stages now use exported constants instead of repeating string literals. The sealed-answer loader also rejects unknown fields and trailing JSON, preventing later measurement code from silently accepting schema drift.

### Prompt Context

**User prompt (verbatim):** (see Step 4)

**Assistant interpretation:** Resume and complete Phase 0 after the control-contract cleanup, preserving strict contracts and validation evidence.

**Inferred user intent:** Give projectors and the specialist UI one stable, attributable cross-layer vocabulary before implementing configuration diffs and invalidation.

**Commit (RAG-TTC code):** `9cf3b2e368a5cb3ce2993dc9cfaeac57b67c36ff` — "Optimization: freeze cross-layer RAG contracts"

### What I did

- Added `pkg/ttc/optimization` with versioned schema constants and the canonical layer vocabulary from corpus through judge.
- Added typed config references, stage contracts, content-free context lineage, answer lineage, judge lineage, hard constraints, and primary metrics.
- Added `rag-ttc.optimization-semantic-fixture/v1`, pinned to SHA-256 `4c5f87b27e134e464e0da868b44df4c93cb6f344f3880ac905eb7c022effba3f`.
- Linked the optimization fixture to retrieval fixture digest `2fa045999a8a89039e00dd60b3fec2bc17b732d557eb00746e207620a5fbdc7f`.
- Froze 16 stage names and their schemas/layers, including all canonical retrieval stages plus context, answer, and judge stages.
- Added strict unknown-field and trailing-data rejection for the fixture and sealed historical answers.
- Added validation for unique layer, stage, constraint, and metric identities; direct dependencies; lineage layer types; context/answer/judge links; citation admission; and observation alignment.
- Added `scripts/04-validate-phase0.sh` and archived the successful run in `sources/04-phase0-validation.txt`.

### Why

- UI projectors need stable domain contracts rather than raw generic event interpretation.
- A separate optimization fixture can extend the frozen cross-product retrieval fixture without mutating its v1 digest or forcing Coinvault to adopt RAG-TTC answer/judge contracts.
- Exported production stage constants ensure fixture tests and runtime emissions share one vocabulary.
- Strict historical-answer decoding prevents accidental schema widening during remeasurement.

### What worked

- Focused package tests, race tests, and vet passed.
- RAG-TTC's complete pre-commit and pre-push tests, lint, and Glazed vet passed.
- The existing RAG-TTC and Coinvault retrieval fixtures remain byte-identical to the canonical source.
- Fixture tests reject unknown top-level and nested fields, broken dependencies, broken lineage, invalid metric direction, and schema drift.
- The Phase 0 validator ended with `PHASE_0_VALIDATION=PASS`.

### What didn't work

- The first `scripts/04-validate-phase0.sh` run failed immediately with:

  ```text
  fatal: not a git repository (or any of the parent directories): .git
  ```

  The reused OPTKIT-002 fixture-sync script calls `git rev-parse` in its current directory. The Phase 0 validator was running from the workspace root, which is not a Git repository. The fix was to invoke the sync check inside `(cd "$OPTKIT" && ...)`.

### What I learned

- Cross-product retrieval semantics and RAG-TTC optimization semantics need separate fixture versions. Referencing the base fixture by schema and digest preserves both reuse and product ownership.
- Stable stage names require runtime code to consume constants; a fixture containing the same literals is not sufficient protection against drift.
- Layered refs can freeze vocabulary and lineage in Phase 0 without yet implementing Phase 2's transitive invalidation behavior.

### What was tricky to build

- Retrieval policy recheck and evidence stages share the retrieval-stage payload schema but belong to different invalidation layers. The fixture therefore validates both stage name and semantic layer rather than deriving layer solely from schema.
- The fixture must remain content-free enough for public deterministic testing while still proving answer and judge lineage. It uses stable artifact identities and chunk IDs rather than embedding sensitive provider content.
- `Layer` ordering is exposed through a copying function rather than an exported mutable slice, so callers cannot mutate the canonical graph vocabulary.

### What warrants a second pair of eyes

- Review whether `retrieval.policy_recheck` belongs to the evidence layer and whether `retrieval.augmented` belongs to fusion for invalidation purposes.
- Confirm the first hard constraints and primary metrics are sufficient for UI fixture work without implying final promotion policy.
- Review whether context token count should remain a fixture field before the authoritative tokenizer is selected.
- Inspect strict sealed-answer decoding for any intentionally extensible metadata that should stay inside maps rather than becoming unknown top-level fields.

### What should be done in the future

- Phase 2 must derive semantic identities from actual typed configuration values rather than the fixture's stable example identities.
- Future context, answer, and judge producers must emit records conforming to the frozen schemas or intentionally version them.
- Add machine-readable schema documents if external consumers need independent generation or validation.

### Code review instructions

- Start in `pkg/ttc/optimization/contracts.go`, then review `fixture.go` and the canonical JSON fixture together.
- Verify stage constants in `pkg/ttc/search/service.go` are used by retrieval and evidence paths.
- Review `judgeinstrument.decodeSealedAnswer` and its unknown/trailing-field tests.
- Run `scripts/04-validate-phase0.sh` and check both fixture digests plus `PHASE_0_VALIDATION=PASS`.

### Technical details

Frozen layers:

```text
corpus -> chunking -> representations -> embeddings -> indexes
       -> retrieval -> fusion -> reranking -> evidence -> context
       -> answer -> judge
```

Validation markers:

```text
OPTIMIZATION_FIXTURE_SCHEMA=rag-ttc.optimization-semantic-fixture/v1
OPTIMIZATION_FIXTURE_SHA256=4c5f87b27e134e464e0da868b44df4c93cb6f344f3880ac905eb7c022effba3f
BASE_RETRIEVAL_FIXTURE_SHA256=2fa045999a8a89039e00dd60b3fec2bc17b732d557eb00746e207620a5fbdc7f
PHASE_0_VALIDATION=PASS
```

## Step 6: Prove Judgekit conformance without rebuilding the adapter

This step completed OPTKIT-004 Phase 1 as an evidence and conformance phase. The P7 Judgekit adapter already implemented the required behavior; the work here mapped every requirement to code and tests, strengthened the historical remeasurement test to protect exact sealed product bytes, and archived a focused cross-repository validation.

The resulting evidence proves that a second judge protocol creates new epochs and observations over one unchanged answer trajectory. It also proves evidence-hidden claim extraction, cache-bypass execution, current-content instance rebinding, model/prompt/cache/usage provenance, and typed failure or missing states.

### Prompt Context

**User prompt (verbatim):** (see Step 4)

**Assistant interpretation:** Audit and close Phase 1 after Phase 0, adding only the missing proof that historical product artifacts remain unchanged.

**Inferred user intent:** Avoid duplicate implementation while establishing a reviewable, executable acceptance record before configuration-graph work begins.

**Commit (RAG-TTC test):** `e92644779022b49e59be371f399f6317a9f1ceb8` — "Judge: prove sealed product bytes survive remeasurement"

### What I did

- Mapped every OPTKIT-004 Phase 1 deliverable to concrete RAG-TTC and Judgekit APIs and tests in `sources/05-phase1-conformance-matrix.md`.
- Strengthened `TestHistoricalAnswerCanBeRemeasuredUnderNewEpoch` to snapshot exact trajectory and answer-output bytes before measurement and compare them after two judge epochs.
- Verified current-content instance identity rebinding, evidence-hidden extraction, strict model attribution, prompt execution provenance, cache bypass, typed missing/failure states, deterministic evidence separation, and epoch isolation.
- Added `scripts/05-validate-phase1.sh` and archived output in `sources/06-phase1-validation.txt`.
- Confirmed the `Instrument.Measure` path contains no search, retrieval, customer turn, prepared system, or product execution call.

### Why

- Phase 1 was already behaviorally implemented by OPTKIT-002 P7, so a second adapter would create divergence rather than confidence.
- Exact artifact-byte assertions make “no retrieval or answer rerun” concrete: the instrument receives only sealed artifacts and the test proves those product artifacts remain unchanged.
- A conformance matrix makes inherited Judgekit guarantees visible without copying Judgekit logic into RAG-TTC tests.

### What worked

- Judgekit judging, assessment, and protocol packages passed unit, race, and vet checks.
- RAG-TTC judgeinstrument focused and race tests passed.
- All five Phase 1 validation markers passed, ending with `PHASE_1_VALIDATION=PASS`.
- The RAG-TTC full pre-commit test and lint hooks passed.

### What didn't work

- N/A. Phase 1 required one additive regression assertion and no implementation repair.

### What I learned

- API shape contributes meaningful conformance evidence: `Instrument.Measure` cannot rerun product behavior because it receives an artifact store and sealed `episode.Result`, not a product executor.
- Artifact immutability alone is not enough for an acceptance claim; capturing and comparing exact bytes makes the no-rewrite property explicit in the regression test.
- Judgekit owns detailed prompt/model provenance tests while the product adapter owns the mapping from sealed product evidence into Optkit epochs and observations.

### What was tricky to build

- The conformance matrix had to distinguish guarantees supplied by Judgekit from adapter guarantees supplied by RAG-TTC. Evidence-hidden extraction belongs to Judgekit; sealed trajectory decoding and epoch mapping belong to RAG-TTC.
- Deterministic contract artifacts appear in observation diagnostics rather than as substituted judge evidence. The test must preserve that distinction while still proving all evidence refs verify.

### What warrants a second pair of eyes

- Confirm that exact trajectory/output byte preservation plus the execution-free `Measure` API is sufficient proof of no product rerun.
- Review whether the conformance validator's static execution-call guard should evolve into a repository boundary test.
- Confirm the current two constructs and built-in protocol are fixture scope, not a frozen final judge calibration policy.

### What should be done in the future

- Add calibration and perturbation projectors when specialist Judgekit screens are implemented.
- Keep product execution out of measurement instruments; future measurement adapters must accept sealed records only.
- Promotion policies must compare only compatible epoch populations.

### Code review instructions

- Read `sources/05-phase1-conformance-matrix.md` first.
- Review the strengthened historical test in `pkg/ttc/judgeinstrument/instrument_test.go`.
- Follow `Instrument.Measure` through `LoadSealedAnswer`, `BuildInstance`, report validation, and observation creation.
- Run `scripts/05-validate-phase1.sh` and check for `PHASE_1_VALIDATION=PASS`.

### Technical details

```text
SEALED_ANSWER_REMEASUREMENT=PASS
EVIDENCE_HIDDEN_EXTRACTION=PASS
CACHE_BYPASS_PROBE=PASS
TYPED_MISSING_FAILURE_STATES=PASS
PHASE_1_VALIDATION=PASS
```

## Step 7: Implement deterministic layered invalidation

This step completed OPTKIT-004 Phase 2 with a product-owned configuration graph. Each RAG layer now has its own typed value schema and local semantic identity; graph resolution combines local identity with resolved upstream dependencies; diffing distinguishes direct changes from transitive invalidation; and plans state exactly which layers can be reused or must be recomputed.

A promotion manifest skeleton links candidate, campaign, baseline graph, candidate graph, policy identity, and sorted evidence digests. It is deliberately fixed to `proposed` and contains no approval or deployment authority.

### Prompt Context

**User prompt (verbatim):** (see Step 4)

**Assistant interpretation:** Implement Phase 2's typed identities, DAG validation, semantic diff, deterministic invalidation planner, and attributable promotion envelope.

**Inferred user intent:** Let specialist UI comparisons explain what changed, what was reused, and why downstream work must rerun before bundle and answer campaigns are added.

**Commit (RAG-TTC code):** `a6ad69b7535d313efdc3dfaf6a32492f8a22fe47` — "Optimization: add layered invalidation planner"

### What I did

- Added content-derived local config refs with explicit typed value schemas.
- Added complete graph construction, canonical layer ordering, unique layer/identity checks, upstream-only dependency validation, duplicate-edge rejection, and resolved dependency digests.
- Added safe `ReplaceLayer` rewiring for direct dependants.
- Added deterministic `ConfigDiff` and `InvalidationPlan` records with `reuse`, `direct_change`, and `upstream_change` explanations.
- Proved judge-only, answer-only, reranker-only, fusion-only, and chunker-change invalidation laws.
- Added an attributable, evidence-bearing `PromotionManifest` skeleton with deterministic evidence ordering and no approval state.
- Upgraded the optimization fixture and layered-ref schema to v2 so each config ref exposes its value schema instead of hiding type information inside its digest.
- Added `scripts/06-validate-phase2.sh` and archived `sources/07-phase2-validation.txt`.

### Why

- A digest alone cannot tell the UI or planner what type of configuration it identifies. Typed value schemas must be explicit on each layer reference.
- Local identities distinguish direct changes; resolved identities include upstream meaning and expose transitive invalidation.
- Product-owned RAG layers keep chunking, fusion, evidence, answer, and judge vocabulary out of domain-neutral Optkit.
- Promotion references are needed for provenance now, while gates, review, and mutation authority remain correctly deferred.

### What worked

- Judge-only changes recompute only judge.
- Answer-only changes reuse admitted evidence and recompute answer plus judge.
- Reranker-only changes reuse fusion output and recompute reranking onward.
- Fusion-only changes reuse retrieval channel outputs and recompute fusion onward.
- Chunker changes reuse corpus and recompute every downstream layer.
- Ten repeated planner test runs, race tests, vet, full RAG-TTC tests, lint, and Glazed vet passed.
- The validator ended with `PHASE_2_VALIDATION=PASS`.

### What didn't work

- The first transitive invalidation test failed after sequential replacement scenarios:

  ```text
  layer "representations" depends on unknown identity "config:f7088c6bd67d9e7ba89c6957a6f50c6265901f37a1f2501dcb8d5b4883c46a8b"
  ```

  `ReplaceLayer` copied `ConfigRef` values but shared their `DependsOn` slice backing arrays. Rewiring a candidate graph mutated the baseline graph. The fix deep-copies every dependency slice before rewriting edges.

- The first commit attempt failed the repository lint hook:

  ```text
  pkg/ttc/optimization/invalidation.go:61:53: QF1008: could remove embedded field "ConfigRef" from selector
  ```

  Changing `left.ConfigRef.Schema` and `right.ConfigRef.Schema` to promoted selectors `left.Schema` and `right.Schema` resolved it; the second commit passed.

- Phase 2 exposed a Phase 0 contract omission: v1 `ConfigRef` recorded only the envelope schema, not the typed value schema. Rather than hide type information in a digest or add a compatibility shim, the pre-release fixture and ref contracts were replaced with v2 and Phase 0 validation was rerun.

### What I learned

- Graph replacement must deep-copy slice-bearing records even when outer structs are value types.
- Direct semantic identity and resolved dependency identity are different and both are necessary: one explains the changed knob, the other explains invalidated materialization.
- A promotion manifest can be useful before promotion logic if its status remains non-authoritative and its evidence ordering is canonical.
- Contract freeze is only credible when the next consumer can use the contract; Phase 2 correctly forced the missing value-schema field into the Phase 0 fixture before UI work.

### What was tricky to build

- Replacing one layer changes the identity referenced by its direct dependants, but should not mark those dependants as direct semantic changes. `ReplaceLayer` rewires direct edges while preserving each dependant's local identity; resolved digests then propagate invalidation transitively.
- Dependencies are identity references rather than layer names. This permits explicit multi-input nodes while requiring uniqueness and canonical upstream ordering.
- Evidence order must not change promotion-manifest identity, so constructor-owned sorting occurs before both identity derivation and storage.

### What warrants a second pair of eyes

- Review the canonical layer order and direct edge model before bundle-build systems rely on it.
- Confirm `retrieval -> fusion -> reranking -> evidence` captures reuse boundaries needed by planned campaigns.
- Review whether graph IDs should later become a shared Optkit concept only after a second non-RAG domain proves the abstraction.
- Confirm the promotion skeleton contains enough references for UI provenance without implying acceptance.

### What should be done in the future

- Phase 3 systems should key materialized outputs by resolved layer digest.
- Projectors should expose both direct layer diffs and full invalidation steps.
- Add graph persistence as artifacts when the first campaign executes a layered plan.
- Keep approval, holdout gates, and deployment commands out of the manifest until Phase 8.

### Code review instructions

- Start with `optimization/graph.go`, especially `NewConfigRef`, `NewGraph`, and `ReplaceLayer`.
- Review `invalidation.go` for direct versus transitive reasons and deterministic ordering.
- Review `promotion.go` for evidence sorting and non-authoritative status.
- Run `scripts/06-validate-phase2.sh` and inspect all seven acceptance markers.

### Technical details

```text
JUDGE_ONLY_REUSES_UPSTREAM=PASS
ANSWER_ONLY_REUSES_EVIDENCE=PASS
RERANK_ONLY_REUSES_FUSION=PASS
FUSION_ONLY_REUSES_CHANNELS=PASS
CHUNKER_INVALIDATES_DOWNSTREAM=PASS
DETERMINISTIC_PLANNER=PASS
PROMOTION_MANIFEST_SKELETON=PASS
PHASE_2_VALIDATION=PASS
```

Current fixture identity:

```text
schema: rag-ttc.optimization-semantic-fixture/v2
sha256: 1bb7d696b7ecc8971de03a8f38e8cacff75eb0b4a86fa5a43063831ae50e6a2d
```
