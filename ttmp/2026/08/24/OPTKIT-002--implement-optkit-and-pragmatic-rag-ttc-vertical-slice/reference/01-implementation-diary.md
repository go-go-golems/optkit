---
Title: Implementation Diary
Ticket: OPTKIT-002
Status: active
Topics:
    - optkit
    - rag
    - rag-ttc
    - coinvault
    - judgekit
    - implementation
    - migration
    - local-development
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://.goreleaser.yaml
      Note: Normalized release and Homebrew cask metadata
    - Path: repo://Makefile
      Note: Repository-level CGO, no-CGO, race, vet, build, lint, and demo gates
    - Path: repo://artifact/filesystem/store.go
      Note: Filesystem CAS cleanup and durability normalization
    - Path: repo://campaign/reducer.go
      Note: Explicit lifecycle and evidence-event projection behavior
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/design-doc/01-phased-implementation-plan.md
      Note: Controls phase gates and work-slip cadence
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/scripts/01-import-optkit-baseline.sh
      Note: Reproducible archive import and collision policy
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/scripts/02-validate-p0.sh
      Note: Complete P0 validation and smoke workflow
ExternalSources: []
Summary: Chronological implementation record for the Optkit and pragmatic RAG-TTC vertical slice, including phase slips, commands, commits, failures, decisions, and review instructions.
LastUpdated: 2026-08-24T22:50:00-04:00
WhatFor: Preserve enough operational and technical context to review, reproduce, or continue every OPTKIT-002 phase.
WhenToUse: Read before starting a phase, reviewing a phase commit, debugging validation, or preparing final delivery.
---



# Diary

## Goal

Implement the phased Optkit and RAG-TTC vertical slice while preserving exact evidence for decisions, tests, failures, commits, and thermal work-slip checkpoints. Optkit remains domain neutral; RagKit remains the reusable RAG data plane.

## Step 1: Create the implementation ticket and executable phase plan

I created a fresh implementation ticket rather than reopening the completed architecture assessment. The plan converts OPTKIT-001's recommendations into nine gated implementation phases plus final validation, with explicit entry conditions, tasks, acceptance gates, rollback points, commit policy, and before/after thermal slip requirements.

The overall plan slip was printed successfully after the phase design was written. This means the physical checklist and ticket plan begin from the same ownership decision and phase ordering.

### Prompt Context

**User prompt (verbatim):**

> Create detailed phases + tasks.
>
> Then print a brutalist work slip with the overall plan.
>
> Then before each phase implementatino, print a brutalist wor slip about the phase to be implemented, then at the end of the phase pritn out a slip about the work that was done.
>
> commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Create a tracked implementation program, print one overall thermal plan, then enforce a plan-slip/implementation/status-slip loop for every phase with focused commits and a strict diary.

**Inferred user intent:** Make a large multi-repository implementation visible, reviewable, restartable, and physically trackable instead of letting it become an opaque long-running refactor.

### What I did

- Created `OPTKIT-002 — Implement Optkit and Pragmatic RAG-TTC Vertical Slice`.
- Added a design document and implementation diary.
- Added phase-level docmgr tasks for P0-P8 and final validation/delivery.
- Wrote the detailed phase plan with ownership boundaries, tasks, acceptance gates, commit policy, rollback guidance, and cross-phase invariants.
- Printed the overall brutalist plan slip through the remote Almanach renderer and AtomS3R printer.
- Used a QR target pointing to the ticket path on the implementation branch.

### Why

- OPTKIT-001 is a completed assessment and should remain immutable historical context.
- Phase gates reduce the risk of mixing archive import, semantic refactors, evaluation changes, and deletions.
- The plan must state what is *not* moving, because repository collapse is not itself the objective.
- Thermal checkpoints make the requested operational cadence explicit before code changes begin.

### What worked

- Docmgr created the ticket and stable task IDs without vocabulary errors.
- The plan slip printed successfully with `printed: true`, two printer segments, width 384, and height 816.
- The plan fits the full program on one overall work slip while each later phase gets its own smaller slip.

### What didn't work

- N/A. Ticket creation, task creation, document generation, and the first print all succeeded.

### What I learned

- The implementation naturally separates into foundation, semantic service, attribution/evaluation, direct application, durable campaign, judge measurement, and cleanup phases.
- P0 must preserve existing repository plumbing and tickets while importing archive domain code; a blind archive overlay would be unsafe.
- Cross-repository commits need separate hashes even when one phase spans multiple repositories.

### What was tricky to build

- The plan had to be detailed without pretending all later API details are already settled. Each phase therefore names semantic outcomes and tests while allowing local symbol names to follow existing repository conventions.
- The user requested a slip before and after *every* phase. This requirement is now part of the global working protocol and final acceptance gate so it cannot be treated as optional bookkeeping.

### What warrants a second pair of eyes

- Review phase ordering before P2: P1 fixtures intentionally precede retrieval extraction.
- Review whether the final-validation step should remain a separate task rather than P9; it is operationally equivalent to a final phase and will receive slips.
- Verify the archive collision policy during P0, especially Makefile, README, module files, and template package removal.

### What should be done in the future

- Follow the plan-slip/code/status-slip loop without skipping phases.
- Update later phase tasks only when code evidence reveals a necessary change; record the change here.

### Code review instructions

- Start at `design-doc/01-phased-implementation-plan.md` Sections 2-6 and 16-17.
- Compare the phase summary against `tasks.md`.
- Verify the printer receipt in this step against the command record below.
- Run `docmgr doctor --ticket OPTKIT-002 --stale-after 30` before approving ticket initialization.

### Technical details

- Ticket path: `ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice`
- Overall slip task: `OPTKIT-002`, label `PLAN`, title `Optkit RAG-TTC Vertical Slice`.
- Print result: `printed: true`, rendered at `2026-08-24T22:37:00Z`.
- Immediate next phase: `P0 import Optkit baseline`.

## Step 2: Keep RagKit separate and narrow the consolidation target

Before implementation began, the ownership model was challenged directly: should RagKit's primitives move into Optkit at all? The answer is no. The plan now treats Optkit as a general-purpose optimization and experiment framework, RagKit as a cohesive RAG-domain library, and product integration code as the composition point.

This changes P8 from “move generic kit code” to “stabilize boundaries and delete superseded orchestration.” RagOpt's overlapping control-plane semantics remain a consolidation target; RagKit algorithms do not.

### Prompt Context

**User prompt (verbatim):**

> do we need to move ragkit primitives into optkit? or can we keep them separate and consider optkit to be more general purpose?
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Re-evaluate whether repository unification is technically necessary and define a cleaner long-term dependency boundary.

**Inferred user intent:** Avoid bloating Optkit or coupling a general framework to RAG simply to reduce the number of repositories.

### What I did

- Updated the P8 phase-level task to `Stabilize Optkit-RagKit boundaries and delete superseded orchestration paths`.
- Made the separate RagKit boundary a formal accepted decision in the phase plan.
- Defined the integration direction: Optkit campaign → product-owned executable → RAG-TTC service → RagKit algorithms.
- Kept Judgekit separate behind an Optkit instrument unless a later explicit consolidation decision says otherwise.
- Narrowed deletion targets to RagOpt orchestration and duplicated product runners proven superseded.

### Why

- RagKit's document, chunk, representation, index, retrieval, fusion, reranking, and hydration concepts are cohesive domain semantics.
- Optkit can optimize RAG, prompt systems, numeric systems, and other applications without owning every domain's algorithms.
- Product composition is the right place to import both a general control plane and domain data plane.
- Removing duplicate orchestration produces architectural value; moving files merely for repository count does not.

### What worked

- The revised boundary simplifies dependency rules and P8 acceptance criteria.
- P6 can still register a RAG-TTC executable in Optkit without creating an Optkit → RagKit dependency.
- Coinvault and RAG-TTC can continue sharing RagKit independently of Optkit campaign adoption.

### What didn't work

- The initial OPTKIT-001 recommendation included a possible later move of exercised generic RagKit packages into `optkit/rag`. That recommendation is superseded by this decision and must not drive implementation.

### What I learned

- “Unified framework” should mean one experiment/control-plane model, not one source tree.
- Judgekit has the same useful separation property: it is a measurement provider, while Optkit owns measurement scheduling, epochs, and observation custody.
- Permanent domain boundaries can be thin without being backwards-compatibility adapters.

### What was tricky to build

- The dependency direction must prevent Optkit core from importing RagKit while still allowing an Optkit campaign to execute RAG. The solution is product-owned integration code that implements Optkit system/executable contracts and composes RagKit-backed services.
- Deletion language needed to be precise. “Delete old kits” would wrongly imply deleting useful domain libraries; “delete superseded orchestration paths after last-caller switch” is the actual requirement.

### What warrants a second pair of eyes

- Review the future architecture guard rules so they allow integration packages but prevent core dependency inversion.
- Review whether Judgekit should remain a separate module long-term; this ticket assumes yes unless implementation evidence changes the decision.
- Confirm no P0 archive package already embeds a RAG-specific assumption into Optkit core.

### What should be done in the future

- Add dependency guards in P8 after real integration package locations exist.
- Update OPTKIT-001 with a superseding decision note only if readers are likely to follow its older optional package-move sequence without reading OPTKIT-002.

### Code review instructions

- Review `design-doc/01-phased-implementation-plan.md` Sections 2, 14, and 16.
- Ensure P8 tasks preserve RagKit and target RagOpt/product orchestration only.
- During P6 review, reject any Optkit-core import of RagKit or RAG-TTC.

### Technical details

- Updated task ID: `o795`.
- Stable dependency direction: product integration imports both Optkit and RagKit; Optkit core imports neither product nor RAG packages.
- Superseded concept: mandatory migration of RagKit primitives into `optkit/rag`.

## Step 3: Phase P0 — Import and normalize the Optkit baseline

Phase P0 replaced the checked-out template with the complete supplied Optkit vertical slice while preserving GitHub workflows, release configuration, lint policy, project guidance, and both docmgr tickets. The source import is provenance-isolated from repository normalization so reviewers can compare the domain implementation directly with archive revision `1786d1d`.

The imported module now contains 46 core Go files under artifact, budget, campaign, episode, experiment, measurement, scheduling, configuration-space, and storage packages, plus the CLI and Numbergame proof system. Full CGO, no-CGO, race, vet, build, lint, GoReleaser configuration, demo, inspect, and journal/artifact verification checks pass.

### Prompt Context

**User prompt (verbatim):**

> ok, update the design as necessary, continue
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Apply the separate-RagKit decision to the implementation design, then continue the agreed phase/slip/commit/diary workflow without pausing at planning.

**Inferred user intent:** Begin real implementation from a corrected architecture and preserve the promised operational evidence.

**Commit (source import):** `ea51f8fedc3d74d7ccb81911eacb73d380780aa0` — "feat: import optkit baseline at 1786d1d"

**Commit (repository normalization):** `b49aece6daf446fe717d4aa19e42c0fc9e129a95` — "chore: normalize imported optkit repository"

### What I did

- Printed the P0 plan slip before modifying production files; printer result was `printed: true` at `2026-08-24T22:38:26Z`.
- Extracted the source archive into a temporary directory rather than over the worktree.
- Added `scripts/01-import-optkit-baseline.sh` with explicit source directories and collision behavior.
- Generated `sources/01-optkit-baseline-manifest.txt` containing archive SHA-256 `f211877fa4b845f6d0750601e5dfe65c209fc4c12eb86dc9caa454e4baa9f975`, source revision, and per-file hashes.
- Imported the archive packages, CLI, examples, architecture tests, source docs, module path, and README.
- Removed template `cmd/XXX`, placeholder `pkg`, logcopter generator, and stale dependency graph.
- Preserved `AGENT.md`, `.github`, lint, release, lefthook, license, and `ttmp` history.
- Rebuilt the Makefile around the actual CGO/no-CGO/race/vet/build/lint/demo commands.
- Normalized GoReleaser project, binary, command, Homebrew cask, and package metadata and removed the irrelevant disabled Glazed docs-publishing template.
- Updated CI to install SQLite headers and run `make ci-check` without a deleted template logcopter package.
- Fixed imported code issues exposed by the repository's existing lint policy.
- Added `scripts/02-validate-p0.sh` and generated a complete command transcript in `sources/02-p0-validation.txt`.
- Ran the Numbergame demo, inspected its 81-event completed campaign, and verified the journal and 80 unique event payloads.

### Why

- An explicit import script and checksums make provenance reproducible and reviewable.
- Keeping archive semantics in the first commit distinguishes imported behavior from local cleanup.
- Template dependencies and commands referred to a nonexistent `XXX` product and would have made CI/release behavior misleading.
- The no-CGO path is part of the archive's declared portability contract even though local persistence requires SQLite/CGO.
- Lint fixes improve resource cleanup and make evidence-only campaign events explicit without changing architecture.

### What worked

- Baseline tests passed before normalization in normal, no-CGO, race, and vet modes.
- `make ci-check`, `make race`, and `make lint` pass after normalization.
- `goreleaser check` reports `1 configuration file(s) validated`.
- The demo completed eight episodes, committed both finite budgets exactly, recorded a positive paired estimate, and reached campaign status `completed`.
- `campaign verify` validated the journal and all event payload artifacts after reopening the local profile.
- `docmgr doctor` is clean after adding the missing `implementation` vocabulary topic.

### What didn't work

- The first ticket doctor run warned:

  `unknown_topics — unknown topics value(s): implementation (3 docs)`

  I added `topics/implementation` to `ttmp/vocabulary.yaml` and reran doctor successfully.
- The first import commit command included `git add -u cmd/XXX pkg ...` after those paths had already been staged by parent directories. Git returned:

  `fatal: pathspec 'cmd/XXX' did not match any files`

  I used a repository-wide `git add -u`, reviewed the staged list, and committed successfully.
- The first lint run found 21 issues: 16 unchecked concrete `Close`/`Remove` calls, one non-exhaustive campaign event switch, one predeclared `min` variable, two redundant embedded selectors, and one unused SQLite row helper. Production issues were fixed; test cleanup is explicitly excluded from `errcheck` because temporary-store cleanup cannot alter completed assertions.
- The initial `goreleaser check` rejected two deprecated properties:

  `snapshot.name_template should not be used anymore`

  `brews is being phased out in favor of homebrew_casks`

  I removed the snapshot override, migrated to `homebrew_casks`, inspected the installed GoReleaser JSON schema for supported fields, and reran successfully.
- The validation script's first smoke attempt inherited the workspace and failed because the parent `go.work` version is older than sibling module requirements. I set `GOWORK=off` on every direct `go run`.
- The next smoke attempt parsed the demo field as `campaign_id`; the actual stable JSON field is `campaign`, producing `KeyError: 'campaign_id'`. I corrected the parser. The next full validation passed.

### What I learned

- The archive is genuinely dependency-free at the Go module level; SQLite is accessed through a narrow system-library CGO binding.
- The current parent `go.work` is not usable for direct commands across all sibling modules. Repository scripts must use `GOWORK=off` until the workspace Go directive is upgraded.
- The campaign reducer intentionally receives evidence events that change journal custody but not lifecycle projection. Listing them explicitly is clearer than a default branch and satisfies exhaustive checking.
- GoReleaser's current binary distribution mechanism is `homebrew_casks`; preserving old `brews` configuration would leave a known migration warning.
- A completed Numbergame run contains 81 control events but only 80 unique payloads because content addressing legitimately deduplicates equal semantic payload bytes.

### What was tricky to build

- Collision handling had three classes: archive-owned semantic source, repository-owned lifecycle plumbing, and template-only placeholders. Blind copying would have dropped CI; blind preservation would have retained `XXX` metadata and a huge irrelevant dependency graph. The import script copies only archive semantic directories and module/readme files, while normalization handles repository plumbing deliberately.
- Resource cleanup had to avoid changing primary error semantics. Error-path closes are best-effort (`_ = Close`) because the operation's original stat/migration error remains authoritative. Successful API callers still receive explicit close errors when they call `Close` themselves.
- The campaign reducer's newly explicit no-op event cases are not ignored events: they still advance version, digest, and event count after transition and digest verification. The comment documents this invariant.
- The validation script needed isolation from the parent workspace and a machine-readable campaign identifier from demo output. Both assumptions were corrected from observed failures rather than hidden with shell fallbacks.

### What warrants a second pair of eyes

- Review `artifact/filesystem/store.go` cleanup behavior around temp-file removal, close-on-stat-error, and directory sync.
- Review the explicit evidence-event cases in `campaign.Apply` to confirm no event should update the current overview beyond version/digest/count.
- Review the GoReleaser CGO matrix before the first actual tagged release; configuration validation does not prove cross-compiled SQLite availability.
- Review the lint exclusion for concrete close errors in `_test.go`; production files remain fully checked.
- Compare the import manifest and `ea51f8f` tree with the source ZIP to verify collision policy did not omit a domain file.

### What should be done in the future

- Upgrade the parent `go.work` Go directive in a workspace-management change if multi-module workspace commands are desired.
- Exercise the actual GoReleaser split build in release-readiness work before tagging.
- Begin P1 with credential-free semantic fixtures; do not change retrieval architecture until characterization tests exist.

### Code review instructions

- Review commit `ea51f8f` first as the provenance import; compare it with `sources/01-optkit-baseline-manifest.txt`.
- Review commit `b49aece` second for local normalization and lint changes.
- Start code review at `go.mod`, `internal/archtest/architecture_test.go`, `campaign/reducer.go`, `artifact/filesystem/store.go`, `Makefile`, and `.github/workflows/push.yml`.
- Reproduce all checks with `scripts/02-validate-p0.sh`.
- Inspect `sources/02-p0-validation.txt` and require the terminal marker `P0_VALIDATION=PASS`.

### Technical details

- Archive SHA-256: `f211877fa4b845f6d0750601e5dfe65c209fc4c12eb86dc9caa454e4baa9f975`.
- Archive revision: `1786d1da86c9e03316ed71336bbc993fa30531f0`.
- Module: `github.com/go-go-golems/optkit`.
- Validation revision: `b49aece6daf446fe717d4aa19e42c0fc9e129a95`.
- Validation commands: `make ci-check`, `make race`, `make lint`, demo, campaign inspect, campaign verify.
- Validation result: `P0_VALIDATION=PASS`.

## Step 4: Phase P1 — Freeze cross-product semantic RAG fixtures

Phase P1 created one versioned, byte-identical semantic fixture consumed by Coinvault and RAG-TTC tests without adding a new runtime package or requiring an unpublished RagKit dependency. The fixture freezes the behavior that P2 must preserve: raw channel identities, deterministic fusion and return order, policy-negative filtering, evidence labels and budgets, and explicit answer/no-answer expectations.

The fixture is canonical in the OPTKIT-002 ticket and synchronized into each product's testdata by a checked script. Each product decodes strictly, pins the SHA-256, and tests its own current responsibility: RAG-TTC proves retrieval/evidence behavior, while Coinvault proves authorization-before-use and run-scoped evidence admission.

### Prompt Context

**User prompt (verbatim):** (same as Step 3)

**Assistant interpretation:** Continue from the completed Optkit import into the next planned phase, preserving the thermal slip, test, commit, and diary cadence.

**Inferred user intent:** Establish a safe semantic refactoring baseline before extracting RAG-TTC's retrieval service.

**Commit (RAG-TTC):** `b8aaf41d9a2ded9daca400cb16f9083fbcf37225` — "test(rag): freeze cross-product semantic fixture"

**Commit (Coinvault):** `e3090be05b657b4ce6f0015028f3eba562ef909b` — "test(rag): adopt cross-product semantic fixture"

### What I did

- Printed the P1 plan slip before fixture changes; printer result was `printed: true` at `2026-08-24T22:47:47Z`.
- Added canonical `sources/rag-semantic-fixture-v1.json` with three documents, three chunks, three representations, two channel rankings, four query modes, stage expectations, evidence policy, and answer expectations.
- Added `scripts/03-sync-rag-semantic-fixture.sh` with `sync` and `--check` modes.
- Synchronized byte-identical fixtures into RAG-TTC and Coinvault testdata.
- Pinned canonical SHA-256 `2fa045999a8a89039e00dd60b3fec2bc17b732d557eb00746e207620a5fbdc7f` in both product test suites.
- Refactored RAG-TTC's existing search fixture constructor to load documents, chunks, and channel hits from the canonical semantic fixture.
- Added RAG-TTC assertions for returned/admitted IDs, exact citation labels, repeat-label reuse, and explicit abstention-case custody.
- Added Coinvault assertions for public-scope vector filtering, product-role filtering, evidence label reuse, item-budget rejection, and citation-label expectations.
- Ran focused and complete Go tests in both products.
- Added `sources/03-p1-validation.md` with identities, commands, commits, laws, and the observed unrelated flake.

### Why

- P2 will move retrieval semantics. A small deterministic corpus prevents the refactor from silently changing ranks, evidence, or route behavior.
- Product-local test decoders avoid a new runtime dependency and avoid requiring a RagKit release merely to share test data.
- A canonical ticket file plus byte comparison makes intentional duplication auditable.
- Coinvault and RAG-TTC should share laws without pretending their product policies and service APIs are identical.

### What worked

- Both strict decoders accept the final schema and reject unknown top-level/nested struct fields.
- Fixture synchronization reports both product copies verified and the expected SHA-256.
- Existing RAG-TTC search tests now use realistic public guide/product and analyst schema documents while preserving expected returned order.
- Coinvault proves that the analyst-only vector hit is absent from public results and that role filtering returns only the product hit.
- Final focused and complete product tests pass.
- Both repository pre-commit hooks passed after workspace isolation was explicit.

### What didn't work

- RAG-TTC's first strict decode failed with:

  `json: unknown field "public_vector_chunk_ids"`

  Its local fixture struct omitted two Coinvault-oriented expectations. I added both fields so strict decoding covers the complete shared schema.
- After adopting realistic evidence text, the named-route fallback test returned `chunk-b, chunk-c` instead of `chunk-b, chunk-a`. The route had already admitted `chunk-c`; the initial 120-rune fixture budget admitted only one more distinct chunk. I increased the canonical rune budget to 240, resynchronized both copies, and updated both pinned hashes. The item budget still proves third-distinct rejection in the dedicated Coinvault ledger test.
- The first full RAG-TTC suite run failed in the unrelated admin WebSocket heartbeat test with:

  `heartbeat timeout was not observable`

  The exact isolated test passed 10 consecutive runs, and the complete suite passed on rerun. No heartbeat code was changed.
- The first Coinvault commit hook inherited the parent workspace and failed every Go command because `go.work` declares Go 1.25 while sibling modules require Go 1.26.x. Re-running as `GOWORK=off git commit ...` allowed the existing hook to generate, lint, vet, test, and commit successfully.

### What I learned

- Cross-product fixtures should share data and semantic expectations, not force test helper APIs into a common runtime module.
- Evidence budgets interact with session history: a named-route call can consume budget before a fallback/default call. Characterization must account for ledger scope, not only one isolated retrieval.
- Coinvault's authorization law can be expressed over RagKit hits and content without opening a full production bundle, which keeps the fixture fast and credential-free.
- The current RAG-TTC full suite contains a low-frequency timing flake unrelated to RAG changes; isolated repetition is useful evidence, but the flake remains a product maintenance concern.

### What was tricky to build

- The fixture needed metadata meaningful to both products. The shared keys `access_scopes` and `source_role` already match Coinvault conventions and are safely ignored or displayed by RAG-TTC, so no translation layer was required.
- Byte-identical copies are deliberate because the products consume released RagKit v0.1.9 under `GOWORK=off`; importing a new local fixture package would either break isolated CI or require an unrelated library release. The sync script makes duplication controlled rather than accidental.
- Rune and item budgets test different laws. Increasing the rune ceiling preserved existing multi-route characterization, while the two-item ceiling in Coinvault still isolates deterministic rejection of the third distinct chunk.

### What warrants a second pair of eyes

- Review whether the fixture's duplicate rank-1 lexical hits accurately preserve the tie behavior P2 must maintain.
- Review the decision to keep fixture data duplicated and synchronized rather than releasing a RagKit test-support package.
- Confirm the public/analyst and guide/product/schema roles represent the minimum policy boundary needed by both products.
- Review whether P2 direct-service parity should include session-history scenarios or only retrieval results before evidence admission.
- Track the admin heartbeat flake separately if it recurs in CI.

### What should be done in the future

- P2 must preserve `stage_chunk_ids.fused`, `returned`, and `admitted` for this fixture.
- P3 should use the policy-negative case and public expectations to prove source-policy ordering in RAG-TTC.
- P4 should consume required groups and classify the first stage where a target disappears.
- Add a second fixture version instead of mutating v1 after implementation phases begin to depend on its digest.

### Code review instructions

- Review the canonical JSON and sync script in OPTKIT-002 first.
- In RAG-TTC, start at `pkg/ttc/search/newFixture`, `semantic_fixture_test.go`, and the local testdata copy.
- In Coinvault, start at `internal/knowledge/semantic_fixture_test.go`, especially `authorizeHits` inputs and ledger assertions.
- Run the exact commands in `sources/03-p1-validation.md` and verify the sync SHA.
- Compare commits `b8aaf41d9` and `e3090be05` independently; they intentionally share only fixture semantics.

### Technical details

- Fixture schema: `rag.semantic-fixture/v1`.
- Corpus ID: `cross-product-tiny-corpus-v1`.
- Canonical SHA-256: `2fa045999a8a89039e00dd60b3fec2bc17b732d557eb00746e207620a5fbdc7f`.
- Modes: `positive`, `authorization_negative`, `answer_only`.
- Expected fused order: `chunk-b`, `chunk-a`, `chunk-c`.
- Expected returned/admitted order: `chunk-b`, `chunk-a`.
- Evidence labels: `chunk-b → E1`, `chunk-a → E2`; repeated `chunk-b → E1`.
