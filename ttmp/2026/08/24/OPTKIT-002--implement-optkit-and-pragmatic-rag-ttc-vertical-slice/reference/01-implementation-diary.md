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
    - Path: repo://system/registry.go
      Note: P6 domain-neutral system preparation registry (commit b8e233e86)
    - Path: repo://system/registry_test.go
      Note: P6 registry identity cancellation and race laws (commit b8e233e86)
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/design-doc/01-phased-implementation-plan.md
      Note: Controls phase gates and work-slip cadence
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/scripts/01-import-optkit-baseline.sh
      Note: Reproducible archive import and collision policy
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/scripts/02-validate-p0.sh
      Note: Complete P0 validation and smoke workflow
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/scripts/04-validate-p3.sh
      Note: P3 reproducible validation gate
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/scripts/05-validate-p4.sh
      Note: P4 reproducible validation gate
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/scripts/06-validate-p5.sh
      Note: P5 reproducible validation gate
    - Path: ws://rag-ttc/cmd/rag-ttc/cmds/experiments/answerquality/stageaware.go
      Note: P4 existing-runner integration (commit d7701685d)
    - Path: ws://rag-ttc/internal/customer/ragsearch/connected.go
      Note: P3 connected-RAG augmentation adapter (commit d4c5adab4)
    - Path: ws://rag-ttc/internal/customer/ragsearch/routes.go
      Note: P3 verified configured-route compiler (commit d4c5adab4)
    - Path: ws://rag-ttc/internal/customer/realruntime/composer.go
      Note: P5 transport composition switch (commit 8853613a4)
    - Path: ws://rag-ttc/internal/customer/webchatcmd/run.go
      Note: P5 provider serving uses direct application (commit 8853613a4)
    - Path: ws://rag-ttc/pkg/ttc/customerapp/service.go
      Note: P5 canonical direct customer turn (commit 8853613a4)
    - Path: ws://rag-ttc/pkg/ttc/customerapp/service_test.go
      Note: P5 direct-served parity and redaction laws (commit 8853613a4)
    - Path: ws://rag-ttc/pkg/ttc/customerapp/types.go
      Note: P5 request result event failure and trajectory contracts (commit 8853613a4)
    - Path: ws://rag-ttc/pkg/ttc/retrievaleval/diagnose.go
      Note: P4 stage normalization rank movement and target-loss diagnosis (commit d7701685d)
    - Path: ws://rag-ttc/pkg/ttc/retrievaleval/evaluate.go
      Note: P4 treatment verification metrics and denominators (commit d7701685d)
    - Path: ws://rag-ttc/pkg/ttc/retrievaleval/types.go
      Note: P4 versioned suite and report contracts (commit d7701685d)
    - Path: ws://rag-ttc/pkg/ttc/search/identity.go
      Note: P3 runtime identity and semantic fingerprints (commit d4c5adab4)
    - Path: ws://rag-ttc/pkg/ttc/search/service.go
      Note: P3 policy-safe stage-traced retrieval and reranker fallback (commit d4c5adab4)
    - Path: ws://rag-ttc/pkg/ttc/search/service_policy_test.go
      Note: P3 identity policy fusion and fallback laws (commit d4c5adab4)
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

## Step 5: Phase P2 — Extract the canonical RAG-TTC retrieval service

Phase P2 moved channel execution, collapse, fusion, route augmentation, hydration, and source-catalog verification into a direct `search.Service`. The Geppetto-facing `SearchTool` now owns only model input bounds, conversation-scoped evidence admission, result shaping, registration, and the structured-first tool outcome.

The shared P1 fixture now calls the direct service before the tool and proves identical fused evidence order. There is one retrieval implementation and no transport, sessionstream, HTTP, or Geppetto dependency in the service.

### Prompt Context

**User prompt (verbatim):** (same as Step 3)

**Assistant interpretation:** Continue through the next gated implementation phase after P1 completed.

**Inferred user intent:** Establish the direct service seam required for evaluation and Optkit without changing current serving behavior.

**Commit (RAG-TTC):** `ea6c629be8496a8ac0d262594df8161708dd6d5a` — "refactor(rag): extract canonical retrieval service"

### What I did

- Printed the P2 plan slip before code changes at `2026-08-24T22:59:02Z`.
- Added `RetrievalRequest`, `RetrievalResult`, and `Service` in `pkg/ttc/search/service.go`.
- Moved default/named route selection, lexical/vector execution, collapse, weighted RRF, augmentation, hydration, and verified source checks into `Service.Retrieve`.
- Moved route validation and registration into `Service.AddRoute`.
- Changed `SearchTool` to compose one `Service` and delegate retrieval.
- Kept the evidence ledger on each tool instance so conversation scope did not move into the reusable service.
- Added direct-service fixture assertions for fused and hydrated chunk order.
- Ran focused tests, focused race tests, the complete RAG-TTC suite, lint, vet, and pre-commit tests.

### Why

- Serving and experiments need one semantic execution path below model and transport adapters.
- Retrieval is reusable across calls; evidence labels and budgets are session state and must not become service-global.
- Moving source verification into the service prevents direct callers from bypassing a guarantee previously enforced only while shaping tool output.

### What worked

- The fixture preserves fused order `chunk-b, chunk-a, chunk-c` through the direct service.
- Existing default, named, fallback, augmentation, structured, cancellation, bounds, source verification, and Geppetto registration tests pass unchanged in intent.
- `go test -race ./pkg/ttc/search -count=1` passes.
- The full repository test and lint hooks pass.
- The refactor removed 95 lines from the tool while adding an independently callable service.

### What didn't work

- N/A. The extraction and all validation passed on the first implementation attempt.

### What I learned

- The existing package boundary was already suitable; a separate new module/package was unnecessary.
- Source-catalog verification is retrieval semantics, while citation assignment is conversation semantics.
- Returning channel, fused, and hydrated evidence together gives P3/P4 observability without reconstructing stages from tool output.

### What was tricky to build

- Augmenters require both channel rankings and the content store, so they remain prepared route behavior inside the service rather than tool behavior.
- Tool results copy contributions by hydrated rank. Preserving the complete fused list in `RetrievalResult` keeps this exact behavior while avoiding a second fusion path.
- `AddRoute` mutates setup state and is intentionally restricted by documentation to pre-concurrency composition; P3 will compile routes before exposure.

### What warrants a second pair of eyes

- Review `Service.Retrieve` for exact parity with the deleted block in `SearchTool.RunRoute`.
- Confirm that returning full bounded hydrated evidence is the right direct-service contract before P4 diagnostics.
- Review the permanent tool/service boundary; it is not a compatibility shim.
- Consider whether `Source` should eventually be replaced by source metadata carried on direct evidence results.

### What should be done in the future

- P3 should add immutable runtime identity and prepared route compilation without reintroducing tool-owned retrieval.
- P4 should consume `Channels`, `Fused`, and `Evidence` directly for stage diagnosis.

### Code review instructions

- Start at `pkg/ttc/search/service.go`, then compare the reduced `SearchTool.RunRoute`.
- Review `semantic_fixture_test.go` direct-versus-tool assertions.
- Run `GOWORK=off go test -race ./pkg/ttc/search -count=1` and `GOWORK=off go test ./... -count=1`.

### Technical details

- Direct API: `Service.Retrieve(context.Context, RetrievalRequest) (RetrievalResult, error)`.
- Service output: normalized query, route observation, collapsed channels, fused hits, hydrated evidence.
- Session output: citations, evidence budget effects, new-evidence count, and model-facing result limit.

## Step 6: Refocus the UI as a read-only scientific query plane

The UI scope is now deliberately read-only. Campaign creation, parameter selection, candidate proposal, trial compilation, lifecycle commands, and decisions remain in the CLI/application command path used primarily by LLM agents. The browser is responsible for helping scientists find, navigate, replay, compare, and understand durable evidence.

I inserted P2.5 before the remaining RAG phases so UI work can start against Numbergame immediately. The phase defines query projections, stable deep links, search, paginated timelines, artifact inspection, checked-in fixtures, and replayable SSE, with an explicit acceptance test that no HTTP mutation route exists.

### Prompt Context

**User prompt (verbatim):**

> we care abut navigation and visualization, the creation of campaigns and parameters and all that will be done through the CLI by LLM agents mostly, so the UI should basically allow navigation / search / visualization, no need for editing.
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Remove browser authoring and lifecycle controls from the UI design; optimize the UI/backend for read-only scientific exploration and agent-to-human deep-link handoff.

**Inferred user intent:** Avoid duplicating complex command validation in forms and build the interface that is actually valuable to a scientist reviewing agent-operated campaigns.

### What I did

- Added task `gym7`: `P2.5 Build the read-only campaign query plane and UI navigation contract`.
- Added detailed P2.5 tasks and acceptance gates to the phased implementation plan.
- Defined read models for campaign overview, lineage, episode matrix, events, trajectories, measurements, estimates, decisions, budgets, and artifacts.
- Defined GET-only query/search/SSE responsibilities and stable entity deep links.
- Added an accepted decision that the UI has no campaign mutation API.
- Added the read-only UI rule to cross-phase invariants.

### Why

- CLI and LLM-agent workflows remain the single command path for validated writes.
- Read-only projections can evolve for comprehension without changing authoritative campaign semantics.
- Deep links let agents create a campaign and hand a scientist the exact result to inspect.
- Eliminating browser writes simplifies security, concurrency, validation, and auditability.

### What worked

- Numbergame already supplies enough events, artifacts, snapshots, observations, estimates, and decisions to design every navigation view.
- P2.5 can proceed independently of the real RAG-TTC campaign.
- The query plane aligns with Optkit's existing separation between authoritative journal facts and rebuildable projections.

### What didn't work

- The prior scenario included browser forms and buttons for creating baselines, candidates, trials, and lifecycle commands. That UI model is superseded and should not guide implementation.
- A P3 plan slip had already been printed before this clarification. No P3 code was changed; P3 remains queued and its phase slip can be reprinted when implementation actually resumes.

### What I learned

- The most valuable UI unit is an addressable evidence object, not an editable campaign form.
- Search and provenance navigation are first-class scientific tools: users need to move from decision to estimate to observations to episode trajectory to immutable payload.
- Agent-operated CLI writes should return campaign IDs and URLs as part of the human handoff.

### What was tricky to build

- Read-only does not mean static. Running campaigns still need cursor-based SSE, reconnection, and replay, but these streams carry committed facts and never commands.
- Search must be built over derived indexes/projections rather than allowing UI queries to depend on authoritative SQLite table layouts.
- Artifact previews need sensitivity and size enforcement even in a trusted local environment.

### What warrants a second pair of eyes

- Review the exact boundary between safe artifact metadata/preview and raw artifact download.
- Review whether global search begins with structured filters only or includes SQLite FTS in the first slice.
- Confirm that no convenience lifecycle action is added to the browser later without a new decision.

### What should be done in the future

- Implement and print P2.5 slips before resuming P3.
- Add CLI output fields for canonical campaign and entity URLs once the query server address is configurable.
- Build frontend navigation against checked-in Numbergame API fixtures before requiring a live server.

### Code review instructions

- Review the P2.5 section and `UI is a read-only scientific query plane` decision in the phase plan.
- Verify `tasks.md` contains task `gym7`.
- Reject any P2.5 endpoint that mutates campaigns, queues, budgets, artifacts, or decisions.

### Technical details

- Query transport: GET JSON plus GET SSE.
- Live cursor: campaign journal sequence (`after=<seq>`).
- Primary navigation entities: campaign, candidate, snapshot, patch, trial, episode, observation, estimate, decision, artifact.
- Write transport: existing/future CLI application commands only.

## Step 7: Phase P2.5 — Deliver the read-only scientific explorer

P2.5 was implemented as the separate OPTKIT-003 ticket so its architecture, evidence capture, implementation diary, and reMarkable guide could remain reviewable without burying the RAG runtime work. The completed explorer satisfies the P2.5 gate with stable campaign deep links, search and filtering, sequence-cursor replay, provenance views, artifact previews, and a GET-only HTTP boundary.

The browser remains a projection consumer rather than a campaign controller. CLI and application workflows retain every write, while committed journal facts become visible through polling-backed SSE and rebuildable query projections.

### Prompt Context

**User prompt (verbatim):** (same as Step 6)

**Assistant interpretation:** Implement the read-only query plane and scientific navigation contract before resuming RAG runtime phases.

**Inferred user intent:** Give scientists a usable evidence explorer for agent-operated campaigns without duplicating campaign command semantics in a browser.

**Commit (explorer implementation):** `1b09c3c` — "Add read-only scientific campaign explorer"

**Commit (OPTKIT-003 closure):** `2c1f57428fc3595f9e2bf767c071285b3d3d76ee` — ticket documentation and completion

### What I did

- Created and completed `OPTKIT-003 — Read-Only Scientific Campaign Explorer`.
- Implemented SQLite query projections, a query service, Go 1.22 `http.ServeMux`, embedded semantic HTML/CSS/JavaScript, JSON APIs, and sequence-cursor SSE.
- Added campaign list/search, stage rail, lineage, trial matrix, paired deltas, budgets, event inventory, timeline, and evidence drawer views.
- Enforced GET-only APIs, strict CSP, bounded event responses, and sensitivity/size checks for payload previews.
- Added CGO-gated SQLite integration tests and a validation script proving `POST_MUTATION_STATUS=405`.
- Verified a clean Chromium session with no console warnings or errors.
- Uploaded the 1,575-line implementation guide to `/ai/2026/08/25/OPTKIT-003/OPTKIT-003 Scientific Campaign Explorer Guide.pdf`.
- Printed the P2.5 completion work slip before resuming P3.

### Why

- The explorer is a coherent product slice with different implementation and documentation concerns from RAG-TTC retrieval.
- A separate ticket preserves OPTKIT-002 phase ordering while allowing the substantial UI investigation to close independently.
- Read-only HTTP keeps campaign validation and custody in one command path.

### What worked

- The Numbergame reference campaign projects 81 verified events, eight completed episodes, 24 observations, paired delta `+0.71041675`, and decision `eligible`.
- Full CGO/no-CGO, race, vet, build, lint, HTTP, SSE, browser, and mutation-denial checks passed.
- `docmgr doctor` passed and all eight OPTKIT-003 tasks closed.

### What didn't work

- See the OPTKIT-003 investigation diary for the detailed implementation failures, including event fixture field mismatch, SSE heartbeat observability, JavaScript octal escape syntax, and CSP rejection of inline style attributes. Each was resolved before closure.

### What I learned

- A useful scientific UI is organized around addressable evidence and provenance, not editable dashboard cards.
- Cross-process live updates require authoritative journal-head polling; an in-memory event bus is insufficient.
- Sequence IDs are the correct common cursor for pagination, SSE resumption, and historical replay.

### What was tricky to build

- Browser payload access had to compose event reachability, artifact visibility, sensitivity policy, valid JSON, and a 256 KiB limit. The server now denies unsafe previews rather than relying on client discipline.
- SQLite integration needed CGO build tags while preserving the repository's no-CGO contract. Test files follow the same gate as the implementation.
- The requested visual restraint required removing paper tones, window chrome, shadows, and inline progress geometry while retaining dense scientific navigation.

### What warrants a second pair of eyes

- Review payload-preview reachability and sensitivity checks in `internal/web/server.go`.
- Review sequence-cursor SSE behavior under long-running, high-volume campaigns.
- Review projection query bounds before using the current event APIs for very large campaigns.

### What should be done in the future

- Add database-side event pagination, entity-level trajectory routes, exact-sequence replay controls, and large-campaign virtualization in focused follow-up tickets.
- Keep all mutation endpoints out of this explorer unless a new architecture decision explicitly supersedes the read-only boundary.

### Code review instructions

- Start with the OPTKIT-003 guide and diary, then review `query/service.go`, `store/sqlite/query.go`, and `internal/web/server.go`.
- Run `OPTKIT-003/scripts/02-validate-explorer.sh` and confirm `EXPLORER_VALIDATION=PASS` and `POST_MUTATION_STATUS=405`.
- Verify `/api/v1/` has no POST, PUT, PATCH, or DELETE implementation.

### Technical details

- API version: `optkit.query/v1`.
- Event page default/hard limits: 200/500.
- Payload preview limit: 256 KiB.
- Browser routes: `#/campaigns` and `#/campaigns/:id`.
- OPTKIT-003 validation and closure are the implementation evidence for OPTKIT-002 task `gym7`.

## Step 8: Phase P3 — Add attributable, policy-safe retrieval routes

P3 turns the direct retrieval seam from P2 into an attributable prepared runtime. Every result now names its bundle, corpus, resolved configuration, query transform, retrieval policy, evidence policy, reranker, tool description, and effective limit; every major candidate boundary emits a stable stage record with a content-addressed candidate-set reference.

Checked-in route configuration now becomes serving behavior during customer search preparation. Representation and source-role references fail closed, connected-RAG augmentation is resolved with its semantic identity, the customer server applies a global public-role floor, and candidates are filtered before fusion and before any external reranker.

### Prompt Context

**User prompt (verbatim):** "Continue implementation from Phase 3 and keep printing the required workslips."

**Assistant interpretation:** Resume the gated program at P3, implement the runtime identity and policy boundary, validate it completely, and preserve the work-slip/commit/diary cadence.

**Inferred user intent:** Make RAG treatment execution scientifically comparable and prevent model-controlled inputs from bypassing prepared source policy before beginning deterministic evaluation.

**Commit (RAG-TTC):** `d4c5adab410f351067ffe1b70ade4c6e98239032` — "Add attributable policy-safe retrieval routes"

### What I did

- Reprinted the P3 plan slip before implementation because the earlier slip had been superseded by P2.5.
- Added `RuntimeIdentity` and deterministic semantic fingerprints for bundle, corpus, resolved configuration, verbatim query transform, retrieval policy, evidence policy, reranker, and tool description.
- Added explicit effective result limit and source values: `default`, `request`, and `request_clamped`.
- Added stable stage traces for raw, collapsed, policy-filtered, fused, augmented, policy-rechecked, reranked, hydrated, returned, and admitted candidate sets.
- Added `ttc-retrieval-candidate-set/v1` content-addressed references to every stage.
- Added a server-owned customer role floor for `faq`, `page`, `post`, `product`, and `ttc_guide`; route policy can narrow but cannot widen that floor.
- Added prepared route compilation over verified bundle representations and documents using RagKit representation-kind and source-role searcher wrappers.
- Added strict representation count, digest, identity, kind, and duplicate checks and explicit unknown source-role rejection.
- Resolved connected-RAG intent routes during preparation, checked the minimum-subject gate against the connected runtime, and included its semantic digest in route identity.
- Added bounded optional reranking with provider-result validation, policy-safe candidate pools, fused-order fallback, and visible degraded stages.
- Added policy, identity, candidate-trace, route compilation, unresolved-role, effective-limit, and reranker-failure tests.
- Added `scripts/04-validate-p3.sh` and captured the complete passing transcript in `sources/04-p3-validation.txt`.

### Why

- Optkit must reject a measurement when intended and observed treatment identities differ; ambient flags and route names are not enough.
- Filtering after fusion allows forbidden candidates to influence public ranks, and filtering after reranking leaks forbidden text across an external provider boundary.
- Checked configuration that never reaches runtime is misleading experiment metadata rather than executable policy.
- Reranker outages should preserve deterministic fused order while remaining visible to diagnostics.

### What worked

- The policy fixture proves a forbidden rank-1 vector candidate is removed before fusion: the allowed product receives score `2/61`, not the late-filter score `1/61 + 1/62`.
- The reranker spy receives only the allowed product candidate.
- Route fingerprints change when RRF semantics change and remain stable when source-role ordering changes.
- Focused tests, focused race tests, the full RAG-TTC suite, build, vet, golangci-lint, Glazed lint, and both pre-commit hooks pass.
- Validation ends with `P3_COMMIT=d4c5adab410f351067ffe1b70ade4c6e98239032` and `P3_VALIDATION=PASS`.

### What didn't work

- The first draft of `identity.go` was accidentally written only through the beginning of `semanticID`, leaving a malformed partial return expression. I overwrote it immediately with the complete file before formatting or testing.
- The first full repository test run failed in the unrelated existing WebSocket timing test with the exact error:

  `--- FAIL: TestWebSocketHeartbeatTimeoutAndServerCloseAreDeterministic (0.04s)`

  `websocket_test.go:86: heartbeat timeout was not observable`

  Running `GOWORK=off go test ./internal/admin/chatserver -run TestWebSocketHeartbeatTimeoutAndServerCloseAreDeterministic -count=5` reproduced the timing failure twice. A subsequent uncached full suite passed, and the pre-commit full suite also passed. No chatserver code was changed.

### What I learned

- The default bundle search path currently considers all indexed representation kinds; its route identity must say `all`, while prepared routes explicitly identify `raw` or `summary`.
- A route-level role list is not a sufficient security boundary. The immutable server role floor must be intersected into every prepared route.
- Connected-RAG's configuration is only resolved when an intent actually requests connected facts, avoiding unnecessary database handles for ordinary configurations.
- A bounded reranker must validate that every returned candidate belonged to the exact pool, not merely somewhere in the larger fused list.

### What was tricky to build

- Representation identity is stored in the immutable bundle's `representations.json`, while `indexbundle.Bundle` intentionally exposes serving indexes rather than the complete representation list. The composition root now reloads that schema-v2 payload, performs strict JSON decoding, checks count/kinds/duplicates, and verifies the vector manifest's representation digest before compiling route wrappers.
- Connected augmentation predates the direct `search.Service` contract and consumes RagKit's `answering.RetrievalResult`. A narrow application adapter translates the already policy-filtered baseline into that contract, retains its complete trace, then subjects augmented fused candidates to the service's final role-policy recheck.
- Reranker fallback must preserve the original fused ordering and scores while still attributing the configured reranker identity and a stable error class. Provider failures therefore degrade the rerank stage without failing retrieval.
- Candidate artifact references are canonical candidate-set digests at this phase; P6 must persist the corresponding trace payloads into Optkit CAS artifacts when episodes are recorded.

### What warrants a second pair of eyes

- Review the hard-coded customer role floor in `internal/customer/ragsearch/ragsearch.go` against all production corpus roles.
- Review `loadVerifiedRepresentations` for bundle schema evolution and local filesystem time-of-check/time-of-use assumptions.
- Review connected runtime lifecycle and the adapter's mutation of channel maps.
- Review reranker result validation, especially score attribution and fused-tail preservation.
- Confirm that model selection among server-prepared route names is an authorized experiment surface; model input still cannot supply role lists, searchers, stores, or provider settings.
- Track the unrelated chatserver heartbeat test as a separate flaky-test maintenance issue.

### What should be done in the future

- P4 should persist or project these stage candidate sets and use their digests to classify the first target-loss stage.
- P6 should convert candidate-set references and traces into durable Optkit episode artifacts.
- Add a production reranker composer only when a checked configuration names a concrete provider; the optional service contract is ready but no ambient provider is silently enabled.

### Code review instructions

- Start at `pkg/ttc/search/identity.go` and `pkg/ttc/search/service.go`.
- Review policy ordering in `Service.Retrieve`, then inspect `service_policy_test.go` for the pre-fusion and pre-reranker laws.
- Review `internal/customer/ragsearch/routes.go`, `connected.go`, and the `Open` lifecycle in `ragsearch.go`.
- Run `OPTKIT-002/scripts/04-validate-p3.sh` from the Optkit repository and confirm `P3_VALIDATION=PASS`.
- Compare the RAG-TTC working tree to commit `d4c5adab4`; it should be clean.

### Technical details

- Runtime identity schema is the `RuntimeIdentity` JSON object embedded in `RetrievalResult` and `SearchOutput`.
- Query transform: `ttc-query-transform/v1;query=verbatim`.
- Candidate reference schema: `ttc-retrieval-candidate-set/v1`.
- Server role floor: `faq`, `page`, `post`, `product`, `ttc_guide`.
- Reranker statuses: `completed`, `skipped`, or `degraded`; current fallback classes include `hydrate_pool`, `provider_error`, and `invalid_result`.
- Validation transcript: `sources/04-p3-validation.txt`.
- P3 completion slip: `printed: true`, 384 × 853, two segments, rendered `2026-08-25T01:10:54Z`.

## Step 9: Phase P4 — Build deterministic stage-aware retrieval evaluation

P4 adds a strict TTC-owned evaluator over P3 search outputs. Versioned cases distinguish positive retrieval, authorization-negative, and answer/judge-only modes; target resolution is preflighted before execution; query failures remain rows under an explicit denominator policy; and every successful row retains treatment verification, per-stage rankings, rank movement, target-loss diagnosis, retrieval metrics, latency, and provider-call counts.

The existing answer-quality runner now emits canonical stage-aware JSON and concise Markdown for every arm. This integration is deliberately an application-side adapter over the legacy RagKit answering result, while the evaluator's primary contract consumes the complete P3 `SearchOutput` identity and stages directly.

### Prompt Context

**User prompt (verbatim):**

> continue
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Continue immediately from the completed P3 gate into P4, preserving the required plan-slip, implementation, validation, commit, diary, and completion-slip sequence.

**Inferred user intent:** Keep advancing the executable vertical slice rather than stopping after explanation or one phase.

**Commit (RAG-TTC):** `d7701685ddbf4b7a273c924e157fbb08914be889` — "Add deterministic stage-aware retrieval evaluation"

### What I did

- Printed the P4 plan slip before implementation; it rendered at `2026-08-25T01:11:31Z`.
- Added `pkg/ttc/retrievaleval` with strict suite loading, validation, target resolver identity, execution, treatment verification, stage normalization, diagnostics, metrics, and report rendering.
- Added schemas `ttc-retrieval-suite/v1`, `ttc-retrieval-report/v1`, and evaluator identity `ttc-retrieval-evaluator/v1`.
- Added case modes `positive`, `authorization_negative`, and `answer_or_judge_only`.
- Added denominator policies `count_as_zero` and `exclude`; executor failures remain explicit case rows with latency/provider-call evidence.
- Added detailed lexical/vector raw, collapsed, and policy-filtered rankings plus fused, augmented, policy-rechecked, reranked, returned, and admitted rankings.
- Added first-loss diagnosis, missing-at-raw groups, relevant-candidate rank movement, and post-policy forbidden-candidate detection.
- Added macro recall, precision, MRR, nDCG, required-group coverage, source-document diversity, mean latency, provider calls, treatment mismatch counts, and authorization violation counts.
- Added strict intended/observed comparison for route, bundle, corpus, resolved config, query transform, retrieval policy, evidence policy, reranker, tool description when supplied, effective limit, and limit source.
- Added canonical indented JSON and concise Markdown report projections.
- Added the checked-in four-mode deterministic suite at `pkg/ttc/retrievaleval/testdata/ttc-retrieval-suite-v1.json`.
- Integrated the evaluator into the existing answer-quality runner; each arm now writes `results/stage-aware-<arm>.json` and `.md`.
- Added an answer-quality target resolver for chunk, representation, document, and evaluation-unit judgments and content-addressed identities for the in-memory bundle, selected suite, settings, query transform, retrieval policy, and evidence policy.
- Added `scripts/05-validate-p4.sh` and captured `sources/05-p4-validation.txt`.

### Why

- Aggregate recall alone cannot distinguish absent targets from policy removal, fusion loss, reranking loss, return-limit truncation, or evidence admission.
- Authorization-negative cases must not share positive-case denominators or reward retrieval of forbidden evidence.
- A treatment score is invalid when intended and observed identities differ, even if the numerical ranking looks good.
- Query errors must not silently shrink denominators.
- Existing answer-quality experiments need the new report now, but replacing their entire retrieval pipeline would violate P4's narrow integration task and preempt P5/P8 boundaries.

### What worked

- The deterministic suite identifies the comparison guide's first loss at `returned` while retaining product evidence.
- The policy case detects a forbidden chunk reintroduced after an initially clean policy-filtered stage.
- Treatment mismatch and execution failure tests both count as zero under the configured denominator instead of disappearing.
- The answer-quality sample run writes stage-aware JSON/Markdown for BM25, vector, and RRF arms with five positive rows, zero query failures, and zero treatment mismatches.
- Focused tests, focused race tests, the full repository suite, build, vet, golangci-lint, Glazed lint, and pre-commit hooks pass.
- Validation ends with `P4_COMMIT=d7701685ddbf4b7a273c924e157fbb08914be889` and `P4_VALIDATION=PASS`.

### What didn't work

- The first stage-normalization sketch merged lexical and vector raw rankings into one `raw` list. That lost the channel-specific evidence required by the phase plan. I changed retained report stages to `lexical_raw`, `vector_raw`, and corresponding collapsed/policy stages, while diagnostics use a separate deterministic union funnel.
- The first answer-quality resolver handled document, chunk, and representation targets but corpus inspection showed the full TTC dataset uses `target: unit`. I added evaluation-unit-to-chunk resolution before running the integration suite.
- The first answer-quality suite ID was a constant. That could label different query selections as the same suite, so it now includes a digest of corpus identity and exact selected cases.
- No test or validation command failed after these design corrections.

### What I learned

- Required evidence groups are alternatives: retrieving any chunk in a document or evaluation-unit group satisfies that group. Recall and nDCG must therefore credit groups rather than penalize a system for not returning every alternative chunk.
- A report needs both channel-specific rankings for evidence and an ordered union funnel for meaningful first-loss diagnosis.
- Target-resolution errors are suite preparation errors, not query outcomes; all targets are preflighted before executor calls.
- The legacy answer-quality runner can provide useful stage-aware evidence, but its adapter identity must remain distinct from the canonical P3 service identity.

### What was tricky to build

- nDCG over alternative evidence groups needs one gain for the first retrieved member of each group. Resolved groups are required to be disjoint so one chunk cannot receive multiple gains and produce nDCG above one.
- Rank movement is only meaningful between adjacent logical funnel stages, not between lexical and vector channel lists. The evaluator retains detailed channel stages but computes movement over normalized raw/collapsed/policy unions.
- Authorization checks begin at policy filtering and inspect every later stage, catching forbidden evidence added by augmentation or reintroduced before return.
- The answer-quality adapter supports old BM25/vector/RRF/query-strategy results without pretending they are native P3 service executions. It generates a separate in-memory bundle and adapter policy identity and leaves provider-call counts at observed values rather than inventing calls.

### What warrants a second pair of eyes

- Review group-based nDCG and precision semantics in `pkg/ttc/retrievaleval/evaluate.go`.
- Review stage ordering and union behavior in `diagnose.go`, especially augmentation and policy recheck.
- Review whether treatment mismatches should always count as zero or become a hard report error in later Optkit instruments.
- Review the legacy answer-quality adapter identities; actual reranker provider/model identity remains more precise on the native P3 path than on this compatibility integration.
- Review artifact filename sanitization in `writeRetrievalArtifacts` even though current arm names are checked constants.

### What should be done in the future

- P5 should expose direct answer outcomes while retaining the same retrieval report inputs.
- P6 should write native P3 stage traces and candidate sets into Optkit artifacts rather than relying on the answer-quality adapter.
- Add provider-call instrumentation to executors that actually cross embedding, generation, reranking, or judge boundaries; the evaluator already aggregates it.

### Code review instructions

- Start with `pkg/ttc/retrievaleval/types.go`, then `evaluate.go` and `diagnose.go`.
- Review `evaluate_test.go` for target loss, policy violation, treatment mismatch, denominator, latency, provider-call, and deterministic-rendering laws.
- Review `cmd/rag-ttc/cmds/experiments/answerquality/stageaware.go` as an application adapter, not a second evaluator.
- Run `OPTKIT-002/scripts/05-validate-p4.sh` and confirm `P4_VALIDATION=PASS`.

### Technical details

- Deterministic CI suite: `ttc-cross-product-deterministic-v1`.
- Native target resolver: `ttc-target-resolver/chunk-id/v1`.
- Answer-quality resolver: `ttc-target-resolver/answer-quality/v1`.
- Positive final-stage preference: admitted, returned, reranked, policy-rechecked, augmented, fused, policy-filtered, collapsed, raw.
- Rank movement delta is `from_rank - to_rank`; positive values mean improvement.
- Validation transcript: `sources/05-p4-validation.txt`.
- P4 completion slip: `printed: true`, 384 × 831, two segments, rendered `2026-08-25T01:30:56Z`.

## Step 10: Phase P5 — Route customer turns through a direct application service

P5 introduces `pkg/ttc/customerapp` as the canonical customer-domain turn below HTTP, WebSocket, and sessionstream. It owns the bounded provider/tool loop, fresh canonical search session, answer schema, citation mapping, safe abstention, failure taxonomy, content-free trajectory events, and redacted durable projection. Admin tool-answer execution remains separate.

The provider webchat composition now wraps this direct service in a one-method engine adapter. The existing transport still supplies profile system instructions, streaming context, snapshots, and turn persistence, but it no longer owns a separate customer tool loop. Direct and served tests run the same service and produce equivalent answer, citation, contract, evidence, failure, and provider-call outcomes.

### Prompt Context

**User prompt (verbatim):** (same as Step 9)

**Assistant interpretation:** Continue from P4 into the next planned phase and land the canonical direct customer application boundary without pausing.

**Inferred user intent:** Make the real customer product executable directly by future Optkit campaigns while preserving serving fidelity and the required phase evidence cadence.

**Commit (RAG-TTC):** `8853613a4ffac96d6e801580fd300b2b300a036d` — "Route customer turns through direct application service"

### What I did

- Printed the P5 plan slip before implementation; it rendered at `2026-08-25T01:31:23Z`.
- Added `pkg/ttc/customerapp` with direct `Service.RunTurn`, typed request/result, application events, event sink, statuses, failures, trajectory projection, and engine adapter.
- Kept `pkg/ttc/toolanswer` and admin assistant composition unchanged so customer and admin applications remain separate.
- Reused P3 `search.SearchTool` sessions; every turn receives a fresh evidence ledger and retains detached chronological search-call/stage evidence.
- Added structured-output setup, bounded provider calls, reserved final call, tool error continuation, canonical citation-label validation, and immutable chunk-ID mapping.
- Added explicit statuses `completed`, `abstained`, and `failed`.
- Added failure codes for validation, cancellation, retrieval miss, policy removal, admission truncation, generation failure, malformed output, unsupported claim, unresolved citation, presentation failure, provider-call limit, and tool failure.
- Added content-free application events for turn start, retrieval completion, turn completion, and turn failure.
- Added `Result.Trajectory()` as the durable redacted projection: IDs, status, abstention, immutable citation IDs, failure codes, provider calls, and duration only.
- Added retrieval diagnostics over P3 stages and made `SearchTool.Calls()` return deep detached copies.
- Added direct-versus-engine-adapter domain parity tests, fresh-ledger tests, safe-abstention and contract-failure tests, provider/presentation failure tests, retrieval-diagnostic tests, and redaction tests.
- Added `ApplicationEngineFactory` to the real runtime composer/resolver and switched provider webchat composition from an outer registry/tool loop to the customer application engine.
- Preserved outer transport middleware, snapshot context, session metadata, and turn-store persistence; the direct application inherits event sinks and owns the inner tool loop.
- Added the provider's tool-result reorder middleware inside the direct application so every inner provider call retains prior serving behavior.
- Added `scripts/06-validate-p5.sh` and captured `sources/06-p5-validation.txt`.

### Why

- Optkit must invoke the exact product-domain turn directly rather than reconstructing customer behavior through WebSocket choreography.
- HTTP/sessionstream and direct evaluation must not own independent tool-loop, answer-contract, or evidence semantics.
- Customer and admin agents have different tools, policies, and product obligations and should not be collapsed into one application service.
- Durable scientific trajectories must not retain customer questions, model prose, source text, or raw failure messages by accident.

### What worked

- Direct and engine-adapted execution produce equivalent statuses, grounded answers, contracts, citation chunk IDs, evidence IDs, failures, and provider-call counts.
- Two turns in the same session each begin at citation `E1`, proving evidence ledger scope is per turn/session factory rather than service-global.
- Safe abstention is a successful explicit outcome; malformed output, unsupported claims, and unresolved citations are failed contract outcomes.
- Provider and event-sink failures surface as distinct generation and presentation codes.
- Content-free events and `ttc-customer-trajectory/v1` omit the test customer secret and source prose while preserving immutable evidence IDs.
- The composer test proves application-engine composition bypasses the old outer tool registry and remains callable through the transport engine contract.
- Focused tests, focused race tests, the full repository suite, build, vet, golangci-lint, Glazed lint, and pre-commit hooks pass.
- Validation ends with `P5_COMMIT=8853613a4ffac96d6e801580fd300b2b300a036d` and `P5_VALIDATION=PASS`.

### What didn't work

- The first customerapp compile failed because imports retained from the extracted source skeleton were unused:

  `pkg/ttc/customerapp/service.go:7:2: "sync/atomic" imported and not used`

  `pkg/ttc/customerapp/service.go:15:2: "github.com/go-go-golems/ragkit/rag" imported and not used`

  `pkg/ttc/customerapp/service.go:16:2: "github.com/go-go-golems/ragkit/rag/answering" imported and not used`

  I moved provider-budget and answer-contract logic into focused files and removed the stale imports.
- The first webchat compile failed with:

  `internal/customer/webchatcmd/run.go:220:3: unknown field ApplicationEngineFactory in struct literal of type realruntime.ResolverOptions`

  I had added the factory to `ComposerOptions` but not propagated it through `ResolverOptions`; adding the field and forwarding it fixed the composition path.
- The first service draft briefly contained placeholder type anchors for contract and diagnostics functions. I removed them and implemented typed `search.Citation`/`search.SearchOutput` helpers before running tests.

### What I learned

- The repository already contained a direct admin/tool-answer service, but reusing it for customers would violate the separate customer/admin boundary and retain the old duplicate search implementation.
- The existing `realruntime.Composer` is the correct switch point: its outer enginebuilder can remain single-pass for persistence and streaming while the application adapter owns the inner tool loop.
- Tool-result reordering must wrap the provider used by the inner application loop; retaining it only outside the application adapter would not affect later provider calls.
- Domain contract failures should return both a populated failed `Result` and a Go error so transport snapshots cannot mistake malformed model output for a successful turn.

### What was tricky to build

- The served path previously built a tool-loop engine directly. To avoid nesting two tool loops, application composition now sets the outer registry to nil; the outer runner calls the application adapter once, while the direct service performs all provider/tool iterations.
- Profile/Garden system instructions still enter through outer middleware. The application prepends its checked orchestration prompt to the resulting turn, preserving both product contract and profile specialization.
- Search results contain source text for immediate answer validation, but durable events cannot. The implementation uses two representations: rich in-memory `Result` and deliberately content-free `Trajectory`/`Event` values.
- Retrieval miss, policy removal, and admission truncation are nonterminal diagnostic outcomes: a valid safe abstention can coexist with them. Generation and contract failures remain terminal.

### What warrants a second pair of eyes

- Review system-prompt ordering when a long-lived turn already contains historical system blocks.
- Review whether the provider/model metadata stamped from checked tool configuration should be overridden by dynamically selected runtime profiles.
- Review the source-results widget tool in direct non-browser execution; it remains transport-aware and may warrant separation during P8.
- Review whether all domain contract failures should remain Go errors for every direct caller or whether an Optkit adapter should treat populated failed results as observations.
- Review failure messages before any caller persists the rich `Result`; only `Trajectory()` is guaranteed content-free.

### What should be done in the future

- P6 should invoke `customerapp.Service.RunTurn` directly and persist only the redacted trajectory plus explicitly reviewed stage artifacts.
- P8 should remove the now-unused provider `buildProviderToolRegistry` serving path after confirming no caller remains.
- Add a live provider smoke test when credentials and a sealed production bundle are available; this phase deliberately used deterministic scripted engines.

### Code review instructions

- Start with `pkg/ttc/customerapp/types.go`, `service.go`, `contract.go`, and `diagnostics.go`.
- Review `EngineAdapter` and `realruntime.Composer` to confirm there is exactly one customer tool loop.
- Review `webchatcmd/run.go` to verify provider serving constructs `ApplicationEngineFactory`.
- Run `OPTKIT-002/scripts/06-validate-p5.sh` and confirm `P5_VALIDATION=PASS`.

### Technical details

- Direct API: `Service.RunTurn(context.Context, customerapp.Request) (customerapp.Result, error)`.
- Durable projection schema: `ttc-customer-trajectory/v1`.
- Transport adapter: `customerapp.EngineAdapter` implements Geppetto `engine.Engine`.
- Search call history is detached and ordered; stage candidate slices and contribution slices are copied.
- Validation transcript: `sources/06-p5-validation.txt`.
- P5 completion slip: `printed: true`, 384 × 853, two segments, rendered `2026-08-25T01:47:51Z`.

## Step 11: Phase P6 kickoff — Add the domain-neutral system registry

P6 began with the smallest Optkit-core prerequisite: a registry that binds immutable snapshots to product-owned executable factories. The registry stores only factories in process memory; snapshots retain canonical configuration artifacts, while provider handles and credentials remain in the composition root.

This is an intentionally partial P6 checkpoint, not phase completion. The real RAG-TTC adapter, fixed-arm campaign, durable stage artifacts, restart tests, estimates, and CLI remain open, so task `qlyz` is not checked and no P6 completion slip has been printed.

### Prompt Context

**User prompt (verbatim):** (same as Step 9)

**Assistant interpretation:** Continue into the durable campaign phase after completing and documenting P5.

**Inferred user intent:** Reach the first real Optkit product campaign rather than stopping at service extraction.

**Commit (Optkit):** `b8e233e86c5f01144d4d12c408200bd74ff48253` — "Add domain-neutral executable system registry"

### What I did

- Printed the P6 plan slip before implementation; it rendered `2026-08-25T01:48:18Z`.
- Mapped Numbergame scheduling, leases, budgets, episode execution, measurements, estimates, decisions, and campaign event seams.
- Added `system.Factory`, `system.Prepared`, and a concurrency-safe `system.Registry`.
- Made factories declare exact system, configuration-schema, and case-schema identities.
- Made preparation accept an immutable `space.SnapshotRecord` and artifact store, keeping process-local providers outside snapshots.
- Made prepared execution consume a case artifact reference, deterministic seed, and `episode.Sink`.
- Added duplicate registration, unknown system, schema mismatch, prepared-identity mismatch, cancellation, preparation failure, deterministic listing, focused race, full race, full tests, and lint validation.
- Confirmed the Optkit module currently has no published version and the remote task branch is not present, which affects where the product-owned RAG-TTC integration can compile under `GOWORK=off`.

### Why

- Optkit core needs a domain-neutral way to resolve durable system identity into process-local execution behavior.
- Provider clients, credentials, and open bundle handles must never enter snapshots or queue payloads.
- The product adapter should implement this contract from RAG-TTC rather than teaching Optkit core about RAG.

### What worked

- The registry validates every declared identity and verifies that prepared output matches the requested system, snapshot, and case schema.
- Full Optkit tests, full race tests, vet, and lint pass after the new package.
- The registry adds no dependency from Optkit to RagKit or RAG-TTC.

### What didn't work

- The first vet run failed because Optkit's module targets Go 1.23 while the tests used `testing.T.Context`, which vet correctly reports as requiring Go 1.24:

  `system/registry_test.go:53:33: testing.Context requires go1.24 or later (file is go1.23)`

  The same diagnostic appeared at five additional lines. Replacing those calls with `context.Background()` restored the Go 1.23 contract; focused tests, race, and vet then passed.
- `GOWORK=off go list -m -versions github.com/go-go-golems/optkit` returned the module path with no versions, and `git ls-remote ... refs/heads/task/use-optkit` returned no branch. A product-owned RAG-TTC package therefore cannot yet depend on these local Optkit commits through a normal released module version.

### What I learned

- Numbergame contains the complete durability mechanics P6 needs, but its orchestration is example-specific and should not be copied into a general registry.
- The registry boundary can stay very small: factory preparation plus prepared case execution is enough to bridge snapshots, queue work, and episode sinks.
- Cross-repository product integration needs an explicit module publication or integration-module strategy; silently adding a sibling `replace` would make isolated CI non-reproducible.

### What was tricky to build

- Heterogeneous systems cannot be stored behind Go generic interfaces in one registry. The durable boundary therefore uses artifact references and declared schemas, while each product factory performs typed decoding internally.
- Preparation identity must be checked after the factory returns. Validating only the input snapshot would allow a buggy factory to execute the wrong snapshot or case schema.
- The implementation must be concurrency-safe because workers may prepare different systems while a composition root registers factories during startup; registration is documented and tested as setup behavior.

### What warrants a second pair of eyes

- Review whether `Prepared.Run` should receive the complete `experiment.EpisodeSpec` rather than only case artifact plus seed; the current contract keeps trial/arm metadata in queue and episode orchestration.
- Review the module-publication choice before adding a sibling `replace` to RAG-TTC.
- Review whether registry registration should become immutable after first preparation; current locking permits later nonduplicate registration.

### What should be done in the future

- Decide the reproducible product-integration packaging strategy: publish the Optkit branch/version, create a dedicated integration module with explicit CI, or adopt another reviewed boundary.
- Implement the RAG-TTC factory, typed config/case codecs, and prepared executable after that decision.
- Continue all remaining P6 tasks before checking the phase or printing its completion slip.

### Code review instructions

- Review `system/registry.go` and its identity checks first.
- Run `GOWORK=off go test ./system -count=1`, `GOWORK=off go test -race ./system -count=1`, and the full Optkit test/race/lint gates.
- Confirm no Optkit package imports RagKit or RAG-TTC.

### Technical details

- Registry package: `github.com/go-go-golems/optkit/system`.
- Process-local contract: `Factory.Prepare(...) (Prepared, error)`.
- Durable execution inputs: snapshot record, case artifact reference, deterministic seed, episode sink.
- P6 plan slip: `printed: true`, 384 × 809, two segments, rendered `2026-08-25T01:48:18Z`.
