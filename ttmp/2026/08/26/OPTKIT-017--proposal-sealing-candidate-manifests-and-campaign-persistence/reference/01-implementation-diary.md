---
Title: Implementation Diary
Ticket: OPTKIT-017
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
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence/design-doc/01-intern-guide-to-proposal-sealing-candidate-manifests-and-durable-campaigns.md
      Note: Primary design deliverable whose research and delivery this diary records
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence/various/manifest-removal-restart-proof.log
      Note: Exact source-removal and reopen evidence
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence/various/remarkable-implementation-dry-run.log
      Note: Completed sealing implementation bundle selection evidence
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence/various/remarkable-implementation-upload.log
      Note: Successful sealing implementation bundle upload receipt
    - Path: repo://rag-ttc/pkg/ttc/experimentworkbench/sealing.go
      Note: Primary sealing implementation recorded by diary
    - Path: repo://rag-ttc/pkg/ttc/optkitcampaign/candidate_facts.go
      Note: Idempotent event implementation recorded by diary
ExternalSources: []
Summary: Chronological research, design, implementation, validation, and delivery record for OPTKIT-017.
LastUpdated: 2026-08-26T14:20:26.968891089-04:00
WhatFor: ""
WhenToUse: ""
---





# Diary

## Goal

Design canonical sealing, strict candidate manifests, idempotent events, and self-contained campaign persistence. This diary records research, design decisions, validation, commits, delivery evidence, failures, and continuation instructions.

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

**Assistant interpretation:** Replay verified drafts through typed bindings and PatchBuilder so concise candidate authoring becomes immutable campaign fact.

**Inferred user intent:** Make campaigns executable and explainable after their source manifests disappear.

### What I did

- Mapped PatchBuilder, numbergame CandidateProposed precedent, strict manifest loading, and campaign initialization ordering.
- Specified seal/materialize/record APIs, candidate v2 YAML, durable records, catalog provenance, and idempotency.
- Designed stale-draft, retry, restart, manifest-removal, and artifact verification tests.
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

- Manifest-driven materialization may occur before a campaign journal exists, while event recording needs campaign state; the internal boundary must support both without special mutation logic.
- Ticket boundaries are most reliable when each names both its upstream contract and the next consumer.
- Documentation should distinguish observed runtime behavior from proposed APIs and future implementation choices.

### What was tricky to build

- Manifest-driven materialization may occur before a campaign journal exists, while event recording needs campaign state; the internal boundary must support both without special mutation logic.
- The guide had to be self-contained without copying the entire parent architecture. It summarizes prerequisites, then links each claim to the file that owns it.

### What warrants a second pair of eyes

- Review canonical relationships between creation-spec candidate records and CandidateProposed payloads, plus retry conflicts and stale parent/catalog handling.
- Review all proposed public schema/API names before implementation makes them expensive to change.

### What should be done in the future

- Implement strict v2 resolution, materialization, durable schemas, event idempotency, and the manifest-deletion proof.
- Update this diary immediately when implementation reveals a false assumption or accepted contract change.

### Code review instructions

Start with the ticket's `design-doc/01-*.md`, then inspect these decision-shaping files:

- `optkit/space/patch.go`
- `optkit/space/candidate.go`
- `optkit/examples/numbergame/demo.go`
- `rag-ttc/pkg/ttc/experimentworkbench/manifest.go`
- `rag-ttc/pkg/ttc/optkitcampaign/campaign.go`

Validate the design workspace with:

```bash
docmgr validate frontmatter --doc optkit/ttmp/2026/08/26/OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence/reference/01-implementation-diary.md
docmgr doctor --ticket OPTKIT-017 --stale-after 30
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

**Commit (documentation):** `83d0f4f201ae59a0d8983d476e7873312fca8142` — "OPTKIT-015-017: design RAG proposal backend"

### What I did

- Ran `docmgr doctor --ticket OPTKIT-017 --stale-after 30` and obtained `All checks passed`.
- Scanned the index, guide, and diary for generated placeholder sections; none remained.
- Ran focused baseline Go tests before documentation changes; all selected Optkit and RAG-TTC packages passed.
- Ran `git diff --check`/`git diff --cached --check`, corrected whitespace findings, and committed the ticket package.
- Ran the program upload script in `--dry-run` mode, then rendered and uploaded `OPTKIT-017 Proposal Sealing Guide.pdf`.
- Preserved upload evidence in `various/remarkable-dry-run.log` and `various/remarkable-upload.log`.

### Why

- A detailed guide is only useful when frontmatter, links, task state, and delivery artifacts agree.
- Grouping commits by dependency layer keeps review focused while avoiding one nine-ticket mega-commit.
- The reMarkable bundle lets the architecture be reviewed away from the source tree without losing the index or diary context.

### What worked

- `remarquee` reported: `OK: uploaded OPTKIT-017 Proposal Sealing Guide.pdf -> /ai/2026/08/26/OPTKIT-017`.
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
- Inspect commit `83d0f4f201ae59a0d8983d476e7873312fca8142` for the documentation batch.
- Validate locally with `docmgr doctor --ticket OPTKIT-017 --stale-after 30`.
- Consult `tasks.md` for the still-open implementation sequence.

### Technical details

```text
bundle: OPTKIT-017 Proposal Sealing Guide.pdf
remote: /ai/2026/08/26/OPTKIT-017
commit: 83d0f4f201ae59a0d8983d476e7873312fca8142
doctor: clean
production code changes: none
```

## Step 3: Print the sealing plan and define strict candidate authoring

This step began after OPTKIT-016 closed. It printed the seven-phase plan and P1 start receipt, added standalone normalized intent validation to Optkit, and defined a strict candidate authoring schema that has one full baseline and baseline-relative deltas.

The implementation chose an explicit v3 migration because v2 already names the whole-pipeline full-arm schema. There is no dual `arms + candidates` shape and no v2 fallback decoder.

### Prompt Context

**User prompt (verbatim):**

> OPTKIT-016 - OPTKIT-018 in fact, commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill). Print out a brutalist work slip with the plan / different phases for the ticket. then before stsarting a phase, plrint a split about the phase, and print one when the phase is done. budget 2M [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Complete OPTKIT-017 after OPTKIT-016, with strict candidate authoring, canonical sealing, campaign persistence, idempotency, restart evidence, coherent commits, and every required physical phase receipt.

**Inferred user intent:** Make proposal history durable and independently explainable before exposing authoring APIs.

**Commit (Optkit):** `bde7ed031672894bd7d92ea46007e27125c5fef9` — "OPTKIT-017: validate candidate intent before sealing"

**Commit (Optkit):** `88d12f45916ca8f501ddb18190080a35918a9ea3` — "OPTKIT-017: expose candidate intent YAML fields"

**Commit (RAG-TTC):** `1e57f539e850e4941dadbef98a2a539acaf19e6c` — "OPTKIT-017: define strict candidate authoring manifest"

### What I did

- Printed one ticket plan and P1 start/done slips.
- Added public standalone `CandidateIntent.Validate` and normalized set-like intent fields.
- Added exact YAML field tags for structured intent.
- Defined v3 baseline, candidate, mutation, and intent authoring structs.
- Chose baseline-only candidate ancestry and rejected chaining.
- Added strict shape, duplicate, parent, case-reference, and YAML mutation tests.

### Why

- Intent must fail before artifact writes, not only inside `NewCandidate` after patch materialization.
- A schema version cannot silently change meaning from full arms to patch candidates.
- Baseline-only ancestry keeps ordering and cycle semantics explicit.

### What worked

- Valid intent can be checked without synthetic snapshot/patch IDs.
- YAML objects/scalars become JSON bytes through the same `RequestedMutation` DTO.
- Unknown and duplicate mutation fields fail before compilation.
- Full pre-commit tests and lint passed.

### What didn't work

- N/A.

### What I learned

- Existing JSON tags do not guarantee the desired snake-case YAML names; authoring contracts require explicit tags.
- Candidate risk order remains meaningful while expected groups and motivating case IDs normalize as sets.

### What was tricky to build

- YAML `value` must preserve typed structure without creating a second legal-value implementation. A custom mutation decoder accepts only `variable`/`value`, decodes the value node, and JSON-marshals it for the binding codec.

### What warrants a second pair of eyes

- Review the baseline-only decision and v3 schema shape before additional product manifests are authored.

### What should be done in the future

- Candidate chaining requires a separate schema/decision with explicit topological semantics; do not add it through map lookup.

### Code review instructions

- Review `authoring_manifest.go`, the mutation YAML method in `proposal.go`, and Optkit candidate intent normalization.
- Run the authoring shape and YAML tests.

### Technical details

```text
authoring schema: rag-ttc.experiment-manifest/v3
full baseline count: 1
candidate parent rule: baseline only
P1 start/done: printed true
```

## Step 4: Migrate fixtures and resolve candidates through the compiler

This phase replaced runtime v2 assets with strict v3 candidate manifests and rewrote manifest loading around an authoring view plus a resolved execution view. Every candidate resolves by compiling its baseline-relative mutations through `ProposalCompiler`.

Resolved arms retain complete child pipelines for execution, but those values are compiler output rather than a second authored source.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Integrate the strict contract and reject any candidate that the shared compiler marks unsealable.

**Inferred user intent:** Ensure manifests, CLI drafts, and future HTTP compilation have identical mutation semantics.

**Commit (RAG-TTC):** `9e8929659b6fb2a14b85d5569e1ba31ce62e2181` — "OPTKIT-017: resolve candidate manifests through compiler"

### What I did

- Changed `LoadManifest` to decode only v3 `AuthoringManifest` with `KnownFields(true)`.
- Added pure baseline snapshot derivation and candidate compilation.
- Stored resolved graphs and candidate drafts by arm ID.
- Derived manifest identity from normalized mutations.
- Replaced both v2 YAML assets with baseline/candidate v3 assets.
- Added a ticket-local one-time migration script.
- Updated CLI and specialist fixtures and exhaustive loader tests.
- Printed P2 start/done slips.

### Why

- Candidate values must not be independently authored as complete challenger pipelines.
- Removing v2 runtime assets makes stale paths fail visibly.

### What worked

- Config validation reported two resolved arms and six episodes.
- Unknown, domain-invalid, no-op, duplicate, malformed, and old-schema candidate inputs failed deterministically.
- Existing config diff/plan and specialist tests continued to pass on resolved v3 arms.
- Full pre-commit validation passed.

### What didn't work

- N/A.

### What I learned

- Keeping `AuthoringManifest` and resolved `Manifest` distinct preserves existing execution consumers without allowing resolved arms back into the strict decoder.

### What was tricky to build

- Manifest semantic identity should ignore raw JSON whitespace. The identity projection therefore uses compiler-normalized before/after mutations while retaining reviewed prose/intent/cases.

### What warrants a second pair of eyes

- Review the semantic identity projection, especially authored prose participation and normalized mutation inclusion.

### What should be done in the future

- Any additional candidate value kind should enter only through its registered binding and core mutation DTO.

### Code review instructions

- Review `LoadManifest`, `ResolveManifest`, and the two v3 fixture diffs.
- Run strict old-schema/old-arms and unsealable-candidate tests.

### Technical details

```text
old assets removed: semantic-limit-v2.yaml, semantic-limit-challenger-v2.yaml
new assets: semantic-limit-v3.yaml, semantic-limit-challenger-v3.yaml
resolved baseline/candidate: limit-1/limit-2
P2 start/done: printed true
```

## Step 5: Materialize verified drafts through PatchBuilder

This phase implemented the irreversible half of proposal authoring. The materializer validates and recompiles all semantic input before writes, stores the exact catalog, replays normalized assignments through typed bindings and `PatchBuilder`, verifies the durable child, builds the structured candidate, and stores one complete proposal envelope.

Tests decode every stored output and prove candidate, patch, child, graph, plan, catalog, and normalized mutation parity.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Use existing Optkit durable mechanics exactly once and distrust all client-computed child fields.

**Inferred user intent:** Make sealing auditable and eliminate parallel mutation persistence.

**Commit (RAG-TTC):** `800791766a9eb72c0cf4800354231f40dee1316c` — "OPTKIT-017: materialize proposals through PatchBuilder"

### What I did

- Added materialize request/result and candidate proposal envelope contracts.
- Recompiled original mutations against the current registry and parent.
- Rejected stale digest, forged parent, invalid intent, and unsealable drafts before writes.
- Persisted full catalog provenance.
- Replayed sorted normalized values through `Binding.Assign`/`PatchBuilder`.
- Verified child config, snapshot identity, and graph exactly match the draft.
- Added deterministic/retry and display-time identity tests.
- Printed P3 start/done slips.

### Why

- A caller-supplied graph or child config is not trusted evidence.
- Patch assignment artifacts are the durable explanation of each variable change.

### What worked

- Two input mutations produced sorted fusion/retrieval assignment records.
- Loaded child bytes equal compiled `PipelineConfig`.
- Catalog JSON validates with exact semantic/full IDs.
- Fixed-clock materialization is fully repeatable.
- Candidate ID, patch, child, and catalog stay stable across display time.

### What didn't work

- N/A.

### What I learned

- Candidate `CreatedAt` is correctly excluded from semantic identity but included in the proposal envelope artifact; command lookup must happen before rematerializing retries.

### What was tricky to build

- Validation must finish before catalog/assignment writes. Once materialization begins, content-addressed writes may occur before a later store failure; no invalid journal fact is appended.

### What warrants a second pair of eyes

- Review all preview-versus-durable equality checks and candidate envelope contents.

### What should be done in the future

- Store garbage collection may remove immutable artifacts not reachable from a committed journal/spec after a storage or concurrency failure.

### Code review instructions

- Start with `ProposalMaterializer.MaterializeProposal`, then its integration tests.
- Verify `PatchBuilder` remains the only durable assignment implementation.

### Technical details

```text
proposal envelope schema: schema:rag-ttc.candidate-proposal/v1
catalog sensitivity: internal
assignment ordering: fully qualified variable ID
P3 start/done: printed true
```

## Step 6: Persist self-validating campaign candidate records

This phase added an arm-materialization boundary to campaign initialization. Manifest-backed runs materialize one baseline snapshot and seal every candidate through the application materializer before constructing the trial.

Campaign spec v3 stores each arm's pipeline, snapshot, and frozen graph plus complete candidate ancestry and provenance. The old parallel graph map was removed.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Make stored campaign state sufficient for execution and explanation without source authoring values.

**Inferred user intent:** Ensure restart and read projections never recompute historical candidate meaning from current code.

**Commit (RAG-TTC):** `c9342f66ed657d00b70a9b10f81ffa7a49f3d78f` — "OPTKIT-017: persist sealed candidates in campaign spec"

### What I did

- Advanced campaign spec to `schema:rag-ttc.optkit-campaign-spec/v3`.
- Added arm snapshot/graph and candidate record contracts.
- Added `ArmMaterializer` and manifest implementation.
- Added `CampaignSpec.Validate` and validation on creation/reload.
- Removed `ConfigGraphs` and migrated specialist readers to per-arm frozen graphs.
- Verified nested catalog/envelope/assignment/snapshot refs in tests.
- Printed P4 start/done slips.

### Why

- A graph stored both per-arm and in a parallel map can drift.
- Stored specs should reject inconsistent snapshot, graph, trial, or candidate ancestry on every reopen.

### What worked

- Dataset/trial identities reconstruct exactly from stored cases and arms.
- Baseline and candidate snapshots load and validate after campaign creation.
- Specialist cockpit/comparison behavior remained unchanged.
- Full pre-commit validation passed.

### What didn't work

- N/A.

### What I learned

- An interface supplied after the store opens/reset avoids package cycles and lets campaign initialization call application sealing without a second mutation implementation.

### What was tricky to build

- `experimentworkbench` already depends on `optkitcampaign`, so campaign code cannot import the materializer. The `ArmMaterializer` interface and domain-neutral materialization result invert that dependency cleanly.

### What warrants a second pair of eyes

- Review `validateArmMaterialization` and ensure every candidate ref/ancestry relation is checked exactly once.

### What should be done in the future

- Projection APIs may expose candidate records but must not mutate this frozen spec.

### Code review instructions

- Review campaign v3 structs/validation, then `manifestArmMaterializer`.
- Run campaign and specialist focused suites.

### Technical details

```text
campaign spec: schema:rag-ttc.optkit-campaign-spec/v3
parallel graph map: removed
candidate records per treatment arm: 1 in fixture
P4 start/done: printed true
```

## Step 7: Append candidate facts with conflict-safe idempotency

This phase records each sealed candidate as one command-indexed two-event batch after the campaign enters running state. The command key is scoped to campaign identity, and each event carries the normalized semantic seal-request digest.

A retry returns exactly the prior candidate/snapshot events. Reusing the same key for different parent/draft/mutations/intent is a conflict and leaves journal version unchanged.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Make candidate recording atomic and semantically idempotent rather than relying only on content-derived candidate ID.

**Inferred user intent:** Ensure retries and process restarts cannot duplicate or silently replace proposal facts.

**Commit (RAG-TTC):** `2676b6f2b80d632202106428ef98135dbc51e369` — "OPTKIT-017: record candidate facts idempotently"

### What I did

- Derived command IDs from campaign plus caller idempotency key.
- Derived seal-request digest from parent, draft, normalized mutations, and normalized intent.
- Added one SQLite transaction containing `CandidateProposed` then `SnapshotMaterialized`.
- Added command lookup, exact prior-event verification, and digest conflict checks.
- Added retry, conflict, invalid-state, subject, payload, and head-version tests.
- Printed P5 start/done slips.

### Why

- Same candidate ID is insufficient to distinguish a retry from an incorrectly reused client key.
- Candidate and child snapshot custody should enter the journal atomically.

### What worked

- Retry returned two prior events with `Duplicate=true` and did not advance the head.
- Changed semantic request was rejected before append.
- A new key in completed state was rejected without head change.
- Full tests passed.

### What didn't work

- The first commit attempt failed the exhaustive linter verbatim:

```text
pkg/ttc/experimentworkbench/service_test.go:85:3: missing cases in switch of type campaign.EventKind: campaign.CampaignCreated, campaign.PlanCompiled, campaign.CampaignStarted, campaign.CampaignPaused, campaign.CampaignResumed, campaign.CampaignStopping, campaign.CampaignStopped, campaign.CampaignFailed, campaign.CampaignCompleted, campaign.TrialPlanned, campaign.EpisodeScheduled, campaign.EpisodeLeaseGranted, campaign.EpisodeAttemptStarted, campaign.EpisodeCompleted, campaign.EpisodeFailed, campaign.ObservationRecorded, campaign.EstimateRecorded, campaign.DecisionRecorded, campaign.BudgetReserved, campaign.UsageCommitted, campaign.BudgetReleased (exhaustive)
```

The test only filtered two event kinds, so I replaced the intentionally partial switch with two explicit `if` checks. Focused lint passed, then the commit hook passed.

### What I learned

- Optkit's SQLite journal already indexes shared command IDs across multi-event batches; the missing application policy was semantic payload conflict checking.

### What was tricky to build

- Lookup must happen before state validation so a valid retry can return prior facts even after the campaign reaches a terminal state.

### What warrants a second pair of eyes

- Review command ID and seal-request digest schema inputs for future API compatibility.

### What should be done in the future

- HTTP transport should provide exactly one normalized idempotency key source.

### Code review instructions

- Review `candidate_facts.go` and the service-level retry/conflict assertions.
- Inspect both events' shared command ID and tags in SQLite.

### Technical details

```text
batch order: CandidateProposed, SnapshotMaterialized
events per candidate command: 2
same request/key: duplicate result
changed request/same key: conflict
P5 start/done: printed true
```

## Step 8: Integrate campaign CLI sealing and nested verification

This phase confirmed that existing campaign commands enforce draft/seal separation. Dry-run loads and compiles v3 candidates but creates no store. Run opens/reset the store, seals arm materialization, records candidate custody once, executes episodes, and exposes status/verify.

Verification now follows nested references from the campaign spec in addition to direct event payloads.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Prove the real command path, not only package-level materialization.

**Inferred user intent:** Give operators observable evidence that run creates durable candidate facts while dry-run remains pure.

**Commit (RAG-TTC):** `71336ef318e1bf14f65d6d57529f518a2652818e` — "OPTKIT-017: verify nested candidate artifacts"

### What I did

- Extended verification with nested spec artifacts and counts.
- Added structured CLI fields for nested verification.
- Added a built-binary campaign proof script and retained output.
- Asserted dry-run store absence, run completion, candidate event schemas, shared command range, and artifact counts.
- Printed P6 start/done slips.

### Why

- Direct event payload verification does not prove that assignment values, catalog, snapshot configs, and dataset inputs remain present.

### What worked

- Dry-run reported `mutation:false` and did not create its selected store.
- Run completed six episodes with 49 events.
- Exactly one candidate and snapshot fact shared one two-sequence command.
- All 49 direct payloads and nine unique nested payloads verified.
- Full pre-commit validation passed.

### What didn't work

- N/A.

### What I learned

- Nested ref verification is the operational test that a content-addressed campaign specification remains self-contained.

### What was tricky to build

- The proof distinguishes direct journal payloads from refs reachable only through the creation spec; duplicate digests are counted once.

### What warrants a second pair of eyes

- Review whether future candidate envelopes add nested artifact types that must be added to the verifier.

### What should be done in the future

- Keep nested verification explicit when asset-valued coordinates arrive.

### Code review instructions

- Run `scripts/02-run-candidate-campaign-cli-proof.sh` and inspect SQLite command/event assertions.

### Technical details

```text
episodes: 6
events: 49
direct payloads: 49
nested payloads: 9
P6 start/done: printed true
```

## Step 9: Prove restart/removal, expose SealProposal, validate, and deliver

This phase proved that source authoring files are unnecessary after initialization, then completed the public application sealer and final validation. An interrupted campaign resumes from stored spec/work facts after its only manifest copy is deleted; completed campaign status, verify, resume, cockpit, and comparison also work after deletion.

The final service checks idempotency before materialization and restores the original persisted result. This closes the trust boundary needed by OPTKIT-018.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Audit the complete persistence contract and provide the transport-independent sealer before closing the ticket.

**Inferred user intent:** Hand OPTKIT-018 stable compile/seal application operations and self-contained historical facts.

**Commit (Optkit):** `4d8d93f9f60cd0b86895e32df1773c9651d7b199` — "OPTKIT-017: expose normalized candidate intent"

**Commit (RAG-TTC):** `c3ec8b46bb69c932049a5f13a5ca0aa579493e95` — "OPTKIT-017: prove restart without source manifest"

**Commit (RAG-TTC):** `71343d87a01f722d80b49aea7c8f2efad7926e5a` — "OPTKIT-017: expose idempotent proposal sealer"

### What I did

- Added interrupted resume and completed CLI restart proofs after manifest deletion.
- Added specialist projection checks after source removal.
- Exposed normalized candidate intent for semantic seal request identity.
- Added `ProposalSealer.SealProposal`, preflight, lookup-before-materialization, prior-envelope restore, and conflict tests.
- Reran OPTKIT-016 proofs against v3 assets.
- Ran full Optkit and RAG-TTC validation, focused race, and dependency scans.
- Reconciled tasks, guide, diary, relations, roadmap, doctor, completed PDF, and all phase receipts.

### Why

- A materializer is the store half; external authoring needs an application operation that also enforces campaign state and idempotent journal custody.
- Source-removal evidence must include restart and read projection, not only completed execution.

### What worked

- Interrupted campaign completed after manifest deletion.
- Completed CLI status/verify/resume retained exactly `49/49` events.
- Candidate hypothesis, expected metric, motivating case, mutation, snapshots, and graphs were available from stored facts.
- Sealer retry returned the exact original envelope and did not advance journal head.
- Full CI, race, lint, vet, build, doctor, slip, and delivery checks passed.

### What didn't work

- The first interrupted-removal test used `CheckpointAfterLease`; immediate resume saw the unexpired lease and correctly remained `running`. The assertion failed with:

```text
Error: Not equal:
expected: "completed"
actual  : "running"
```

I changed this test to `CheckpointAfterResult`, which leaves terminal scheduler work ready for immediate reconciliation without waiting for lease expiry. Existing campaign tests continue to cover expired-lease recovery with an advanced clock.

- A targeted edit omitted closing quotes in two import paths. gofmt reported:

```text
pkg/ttc/experimentworkbench/sealing.go:12:2: string literal not terminated
pkg/ttc/experimentworkbench/sealing_test.go:13:2: string literal not terminated
```

I inspected both import blocks, restored the quotes, and reran focused/full validation successfully.

### What I learned

- Lease-boundary restart tests must either advance the clock beyond expiry or stop after a terminal result; immediate restart is not evidence of lease reclamation.
- Idempotent service lookup must precede materialization because display timestamps change envelope bytes without changing candidate identity.

### What was tricky to build

- Public sealing must compute normalized request identity before any write, allow terminal-state retries, reject terminal-state new commands, then protect against a journal race during append. Content-addressed orphan bytes after a concurrent race remain unreachable and safe for garbage collection.

### What warrants a second pair of eyes

- Review `ProposalSealer`, request digest normalization, prior-envelope restore, campaign spec validation, and manifest-removal tests.

### What should be done in the future

- OPTKIT-018 should adapt `ProposalCompiler` and `ProposalSealer` through policy-aware command services without placing writes in `specialistapi`.

### Code review instructions

- Review RAG-TTC commits from `1e57f539` through `71343d87` in order and Optkit intent commits.
- Run all ticket scripts, full validation, and `docmgr doctor --ticket OPTKIT-017 --stale-after 30`.
- Inspect the campaign creation spec and candidate/snapshot event batch after deleting the source manifest.

### Technical details

```text
candidate ID: candidate:1e6de48018f0850fe1bb838c4c7d9c920e284a2a0ce867f10df251061610358a
child snapshot: snapshot:6658b47758de8eade33ebfe42e759c2a7cb39c21cd54925f6a6dacb3e0c1f84c
manifest present after run: false
status/resume: completed/completed
events before/after resume: 49/49
dependency packages: 372
import cycles: 0
work slips: 15/15 printed
```
