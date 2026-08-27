---
Title: Implementation Diary
Ticket: OPTKIT-016
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
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-016--proposal-compiler-and-glazed-cli/design-doc/01-intern-guide-to-pure-proposal-compilation-and-glazed-cli-authoring.md
      Note: Primary design deliverable whose research and delivery this diary records
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-016--proposal-compiler-and-glazed-cli/various/cli-no-write-proof.log
      Note: Exact 100-run no-write result
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-016--proposal-compiler-and-glazed-cli/various/remarkable-implementation-dry-run.log
      Note: Completed implementation bundle selection and destination evidence
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-016--proposal-compiler-and-glazed-cli/various/remarkable-implementation-upload.log
      Note: Successful completed implementation PDF upload receipt
    - Path: repo://rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/proposal.go
      Note: CLI implementation recorded by diary
    - Path: repo://rag-ttc/pkg/ttc/experimentworkbench/proposal.go
      Note: Central implementation recorded by diary
ExternalSources: []
Summary: Chronological research, design, implementation, validation, and delivery record for OPTKIT-016.
LastUpdated: 2026-08-26T14:20:25.711695329-04:00
WhatFor: ""
WhenToUse: ""
---





# Diary

## Goal

Design side-effect-free proposal compilation and backend-first Glazed authoring commands. This diary records research, design decisions, validation, commits, delivery evidence, failures, and continuation instructions.

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

**Assistant interpretation:** Build one deterministic application service for manifests, CLI, and future browser drafts without writing artifacts or journal events.

**Inferred user intent:** Prove the complete mutation/diff/plan contract before durable sealing or UI work.

### What I did

- Mapped current experimentworkbench service and Glazed config/campaign command boundaries.
- Specified request, normalized mutation, draft, diagnostic, preview-capability, and digest contracts.
- Designed deterministic errors/order, catalog/proposal commands, and a repeated-compilation no-write proof.
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

- Operator-correctable errors should be complete deterministic diagnostics, while service-state failures remain errors; partial drafts must never be sealable.
- Ticket boundaries are most reliable when each names both its upstream contract and the next consumer.
- Documentation should distinguish observed runtime behavior from proposed APIs and future implementation choices.

### What was tricky to build

- Operator-correctable errors should be complete deterministic diagnostics, while service-state failures remain errors; partial drafts must never be sealable.
- The guide had to be self-contained without copying the entire parent architecture. It summarizes prerequisites, then links each claim to the file that owns it.

### What warrants a second pair of eyes

- Verify no artifact/store dependency reaches ProposalCompiler and that CLI --set parsing does not guess JSON types.
- Review all proposed public schema/API names before implementation makes them expensive to change.

### What should be done in the future

- Implement DTOs/compiler tests first, then catalog/proposal Glazed adapters and store-count smoke evidence.
- Update this diary immediately when implementation reveals a false assumption or accepted contract change.

### Code review instructions

Start with the ticket's `design-doc/01-*.md`, then inspect these decision-shaping files:

- `rag-ttc/pkg/ttc/experimentworkbench/service.go`
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/config.go`
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/command.go`
- `optkit/space/patch.go`
- `rag-ttc/pkg/ttc/optimization/invalidation.go`

Validate the design workspace with:

```bash
docmgr validate frontmatter --doc optkit/ttmp/2026/08/26/OPTKIT-016--proposal-compiler-and-glazed-cli/reference/01-implementation-diary.md
docmgr doctor --ticket OPTKIT-016 --stale-after 30
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

- Ran `docmgr doctor --ticket OPTKIT-016 --stale-after 30` and obtained `All checks passed`.
- Scanned the index, guide, and diary for generated placeholder sections; none remained.
- Ran focused baseline Go tests before documentation changes; all selected Optkit and RAG-TTC packages passed.
- Ran `git diff --check`/`git diff --cached --check`, corrected whitespace findings, and committed the ticket package.
- Ran the program upload script in `--dry-run` mode, then rendered and uploaded `OPTKIT-016 Proposal Compiler CLI Guide.pdf`.
- Preserved upload evidence in `various/remarkable-dry-run.log` and `various/remarkable-upload.log`.

### Why

- A detailed guide is only useful when frontmatter, links, task state, and delivery artifacts agree.
- Grouping commits by dependency layer keeps review focused while avoiding one nine-ticket mega-commit.
- The reMarkable bundle lets the architecture be reviewed away from the source tree without losing the index or diary context.

### What worked

- `remarquee` reported: `OK: uploaded OPTKIT-016 Proposal Compiler CLI Guide.pdf -> /ai/2026/08/26/OPTKIT-016`.
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
- Validate locally with `docmgr doctor --ticket OPTKIT-016 --stale-after 30`.
- Consult `tasks.md` for the still-open implementation sequence.

### Technical details

```text
bundle: OPTKIT-016 Proposal Compiler CLI Guide.pdf
remote: /ai/2026/08/26/OPTKIT-016
commit: 83d0f4f201ae59a0d8983d476e7873312fca8142
doctor: clean
production code changes: none
```

## Step 3: Print the implementation plan and lock compiler contracts

This step began production work only after OPTKIT-013 through OPTKIT-015 were complete. It printed the ticket plan and phase-one start receipt, then defined stable application DTOs, schema-qualified draft identity, preview modes, diagnostic codes, and generic binding error classifications.

The binding boundary now exposes machine-readable failure kinds rather than forcing RAG-TTC to parse codec/domain error strings. Draft identity excludes human wording while covering the semantic fields that OPTKIT-017 must verify.

### Prompt Context

**User prompt (verbatim):**

> OPTKIT-016, commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill). Print out a brutalist work slip with the plan / different phases for the ticket. then before stsarting a phase, plrint a split about the phase, and print one when the phase is done. [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Follow-up user prompt (verbatim):**

> OPTKIT-016 - OPTKIT-018 in fact, commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill). Print out a brutalist work slip with the plan / different phases for the ticket. then before stsarting a phase, plrint a split about the phase, and print one when the phase is done. budget 2M [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Implement OPTKIT-016 completely first, preserving its phase evidence, then continue to OPTKIT-017 and OPTKIT-018 in dependency order under one durable two-million-token completion goal.

**Inferred user intent:** Complete the backend authoring chain with reviewable code history, strict evidence, and physical phase boundaries rather than stopping after design.

**Commit (Optkit):** `759ef01fc7b64da8e12229d596f99567169f9815` — "OPTKIT-016: classify executable binding failures"

**Commit (RAG-TTC):** `1e926542e75405cbbe429180f3d5349ed586c141` — "OPTKIT-016: define pure proposal draft contracts"

### What I did

- Printed one seven-phase plan slip and the P1 start slip.
- Added `space.BindingError` with decode/domain/encode/apply classifications and unwrap support.
- Added candidate-draft v1 and draft-identity v1 schemas.
- Defined request, normalized mutation, diagnostics, capabilities, compiler, and complete draft DTOs.
- Defined a semantic digest projection that excludes message/reason prose.
- Added contract tests for deterministic digest behavior and error classification.

### Why

- Stable diagnostics cannot depend on parsing implementation-specific error messages.
- OPTKIT-017 needs one exact, timestamp-free draft contract to recompile and verify.
- Physical phase receipts were an explicit user requirement and must precede each phase.

### What worked

- Binding tests distinguished JSON/type decode failures from legal-value domain failures.
- Diagnostic wording and capability-reason changes did not alter draft identity.
- Semantic mutation changes did alter draft identity.
- Focused tests and vet passed before both commits.

### What didn't work

- N/A in this phase.

### What I learned

- The generic binding had enough runtime behavior but lacked one application-facing contract: a typed failure category.
- Preview reason prose is presentation and must be treated like diagnostic message prose in identity projections.

### What was tricky to build

- The digest must represent enough information for optimistic verification without becoming a durable candidate ID. The solution uses its own schema and includes semantic result fields, codes, modes, and probes but excludes human copy and time.

### What warrants a second pair of eyes

- Review every field in `digestCandidateDraft` against the future sealer's recompile comparison.
- Review whether the exported binding error strings should remain stable in addition to their typed kinds.

### What should be done in the future

- OPTKIT-017 should compare the recomputed draft digest, not trust a client-supplied child graph or config.

### Code review instructions

- Review Optkit `space/binding.go` first, then RAG-TTC `experimentworkbench/proposal.go` contract definitions.
- Run `go test ./space -count=1` in Optkit and the candidate-draft test in RAG-TTC.

### Technical details

```text
candidate schema: schema:rag-ttc.candidate-draft/v1
identity schema: schema:rag-ttc.candidate-draft-identity/v1
binding kinds: decode, domain, encode, apply
plan slip: printed true
P1 start/done: printed true
```

## Step 4: Implement verified parent handling and the pure compiler engine

This step added store-free verification of loaded snapshot values and implemented the compiler algorithm over the real RAG registry. The compiler validates its service state, semantic catalog, parent system/schema/value/identity, then applies requests in sorted variable order and derives complete planning output.

Successfully applied mutations remain visible in a partial preview if another mutation fails, but any error diagnostic makes the draft unsealable. Invalid parent contracts remain returned errors because no meaningful draft graph can be built from them.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Build the application engine without any artifact-store or journal dependency.

**Inferred user intent:** Make CLI, future HTTP, and sealing all consume one deterministic compiler.

**Commit (Optkit):** `f0dafc6f09b4dabd4e5edfa1d4356f7cf8a16c03` — "OPTKIT-016: verify snapshot values without stores"

**Commit (RAG-TTC):** `4c90cd51611b234541856f7abc34d36a6370bd4b` — "OPTKIT-016: implement pure proposal compiler"

### What I did

- Added `space.VerifySnapshotValue` and refactored snapshot identity derivation into one helper.
- Published the Optkit revision and upgraded RAG-TTC's isolated `GOWORK=off` dependency.
- Implemented catalog validation, parent verification, duplicate grouping, canonical sorting, normalization, pure application, diff, plan, capabilities, diagnostics, sealability, and draft identity.
- Added focused tests for a valid RRF mutation and forged parent/value drift.
- Printed P2 start/done receipts around the phase.

### Why

- A loaded parent record and typed value must agree before a proposal can safely inherit identity from it.
- The strongest architecture boundary is a compiler struct that cannot receive a store or journal.

### What worked

- A `60 → 20` RRF request produced canonical `60`/`20` values, changed only the fusion local identity, and advertised deterministic local contribution preview.
- Full RAG-TTC tests and lint passed during the code commit.
- The compiler source contains no artifact store, local profile, or journal reference.

### What didn't work

- The first formatting/test command failed verbatim:

```text
pkg/ttc/experimentworkbench/proposal.go:11:41: missing import path

Command exited with code 2
```

An exact edit accidentally inserted `},{` between two imports. I inspected the import block, replaced the malformed token with a newline, reran gofmt/tests/vet, and the command passed.

### What I learned

- Snapshot identity can be verified from canonical typed bytes and the record alone; no artifact read is required once the value is already loaded.
- Parent invalidity belongs to service errors while mutation invalidity belongs to complete draft diagnostics.

### What was tricky to build

- Invalid mutations must not be applied, duplicate IDs must skip every conflicting occurrence, and valid independent mutations must remain visible for UI feedback while the draft remains globally unsealable.
- Error categorization had to distinguish operator errors from impossible codec encode/service failures.

### What warrants a second pair of eyes

- Review partial-preview behavior and confirm no caller mistakes `ChildConfig` for sealable output without checking `Sealable`.
- Review exact snapshot ref media type, schema, sensitivity, size, and digest checks.

### What should be done in the future

- The sealer must reject `Sealable=false` and recompile instead of replaying partial-preview fields.

### Code review instructions

- Start at `ProposalCompiler.CompileProposal`, then `space.VerifySnapshotValue`.
- Run `GOWORK=off go test ./pkg/ttc/experimentworkbench -run TestCompileProposal -count=1`.

### Technical details

```text
parent system: system:rag-ttc-pipeline/v2
parent schema: schema:rag-ttc.pipeline-config/v2
atomicity: retain partial preview, never seal on any error
P2 start/done: printed true
```

## Step 5: Prove the exhaustive diagnostic, determinism, and no-write matrix

This phase expanded the compiler tests from a successful example to the complete acceptance matrix. It exercises every operator failure, service-state errors, cancellation, request-order independence, partial preview behavior, synthetic lens failures, child invalidity, graph failure, and repeated compilation.

The tests also lock the rule that capability and diagnostic presentation wording cannot change semantic draft identity.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Prove deterministic behavior exhaustively before mounting transport commands.

**Inferred user intent:** Prevent CLI and later API testing from hiding compiler edge cases.

**Commit (RAG-TTC):** `b7c4e4283cae8cd541aefc0d11c56bf81a41c4cc` — "OPTKIT-016: prove compiler determinism and no-write behavior"

### What I did

- Added empty, unknown, duplicate, malformed JSON, wrong type, domain, no-op, apply, child, and graph failure cases.
- Added parent system/schema/value/digest mismatch cases.
- Added reversed multi-variable request order and exact output equality assertions.
- Added deterministic diagnostic ordering and partial-preview assertions.
- Added a counting artifact-store wrapper and compiled the same draft 100 times.
- Ran focused race tests and full pre-commit validation.
- Printed P3 start/done receipts.

### Why

- Deterministic happy-path output is insufficient if diagnostics or invalid partial output vary by request order.
- Repeated compilation must be observably pure, not only described as pure.

### What worked

- Both mutation orders returned equal child config, normalized mutations, graph, diff, plan, capabilities, and digest.
- One parent materialization write remained exactly one after 100 compilations.
- Error diagnostics sorted by severity, variable, code, and path.
- Focused race and repository-wide CI passed.

### What didn't work

- N/A.

### What I learned

- Child validation and graph derivation errors can only be reached with a deliberately unlawful test lens because production lifted lenses validate complete pipelines during each write.
- The test registry is useful evidence that the compiler handles generic future variables rather than only two hard-coded RAG IDs.

### What was tricky to build

- A meaningful child-invalid test required constructing a legal registry whose synthetic lens returns an invalid complete pipeline without an immediate lens error. This isolates compiler post-validation rather than faking internal fields.

### What warrants a second pair of eyes

- Review custom test lenses to ensure they test application defenses without implying production registration should permit unlawful lenses.

### What should be done in the future

- Add property/fuzz testing if the catalog grows to overlapping lenses; current two production variables are independent.

### Code review instructions

- Review `proposal_test.go` as the executable contract.
- Run focused tests under `-race` and inspect the counting-store assertion.

### Technical details

```text
operator diagnostics tested: 10
reversed request output: equal
repeat count: 100
additional compiler store puts: 0
P3 start/done: printed true
```

## Step 6: Expose the complete catalog through Glazed

This step added ordered `catalog list` and full `catalog show` commands under the existing Optkit-RAG command group. Both commands construct the authoritative runtime registry and emit native structured values rather than prose-flattened JSON.

Command tests verify Glazed v1.4's three universal output flags, stable row order, complete domain/default/probe fields, unknown-variable errors, and processor-error propagation conventions.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Make the real registry reviewable from structured CLI output before proposal authoring.

**Inferred user intent:** Give backend operators and future clients one catalog surface generated from executable declarations.

**Commit (RAG-TTC):** `5ce7876791f5325a29a9f3baf8f25235e9757e1d` — "OPTKIT-016: expose RAG optimization catalog CLI"

### What I did

- Added GlazeCommand implementations and interface assertions for list/show.
- Added the `catalog` Cobra group through existing `addGlazeCommands` composition.
- Emitted semantic/full IDs, section metadata, descriptor fields, raw default, full value specification, and probes.
- Added command-level tests and actual JSON smoke parsing.
- Printed P4 start/done receipts.

### Why

- The compiler's legal variable set must be discoverable from the same catalog that owns its bindings.
- Nested domains must remain machine-readable in JSON/YAML output.

### What worked

- List emitted retrieval then fusion, exactly matching registration order.
- Show emitted RRF float range `{minimum:0.001, maximum:1000}`, default `60`, and its contribution probe.
- Unknown variable IDs failed clearly.
- Full pre-commit tests and lint passed.

### What didn't work

- The first plain `--help` inspection displayed only short-help domain flags, so a grep found `--variable` but not universal structured-output flags. This was not a missing-flag bug: the command advertises `--long-help`. Running `--help --long-help` showed exactly `--format`, `--output-fields`, and `--max-output-rows` plus the domain flag.

### What I learned

- Glazed short-help sections intentionally omit universal output details; flag existence tests plus long-help inspection are the correct evidence.

### What was tricky to build

- `json.RawMessage`, `ValueSpec`, and `VariableDescriptor` had to stay as native row values so JSON preserved numbers and nested discriminated-union fields.

### What warrants a second pair of eyes

- Review table readability and whether future catalogs require pagination; do not change JSON shape to solve presentation concerns.

### What should be done in the future

- OPTKIT-018 should expose this same serializable catalog without executable Go bindings.

### Code review instructions

- Review `catalog.go`, command composition, then `TestCatalogCommandsEmitOrderedCompleteDescriptors`.
- Run both catalog commands with `--format json` and inspect long help.

### Technical details

```text
list rows: 2
catalog semantic ID: sha256:d3034d1d61cb5da92649bf9d199015e25a6e5223ed50741f594ceba2093730b6
catalog full ID: sha256:d20f66171bfe0c5ef1a7ba4aade490c1d98b7c4d4d791e52b4b1f2a417c6f6ca
P4 start/done: printed true
```

## Step 7: Add strict proposal compilation commands

This step added the manifest/arm parent adapter and `proposal compile` GlazeCommand. A selected arm becomes a verified snapshot record through pure derivation, then repeated `--set` inputs or one strict YAML/JSON mutation array enter the same application compiler.

The command emits every candidate-draft field as one nested structured row. Invalid mutation values return successful structured output with `sealable:false`; malformed command envelopes and files return command errors.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Add a thin CLI adapter without a second mutation or graph-planning implementation.

**Inferred user intent:** Make pure proposal behavior directly inspectable and scriptable before persistence and HTTP work.

**Commit (Optkit):** `b45db5582e75117d17c0f976eb2c7665be3057e7` — "OPTKIT-016: derive snapshot records without writes"

**Commit (RAG-TTC):** `fad4920849dbb41c6f6ce5baabf20b6debd1def3` — "OPTKIT-016: add pure proposal compile CLI"

### What I did

- Added `space.DeriveSnapshotValue` and made materialization consume/verify its exact derivation.
- Added `LoadedManifest.Pipeline` and `CompileManifestProposal`.
- Added strict repeated `variable=JSON` parsing using first-`=` splitting.
- Added strict array-file parsing with known-field checks, size limit, required values, one-document rule, and JSON canonical transport.
- Added full draft row emission and the proposal command group.
- Added CLI/direct-service parity, mutation file, malformed source, help/flag, and processor failure tests.
- Ran actual valid and invalid JSON CLI smoke checks.
- Printed P5 start/done receipts.

### Why

- Manifest parents should produce the same deterministic snapshot ID later materialization will use without performing a draft write.
- Complex values need a file format that maps to the core mutation DTO rather than a CLI-only semantic schema.

### What worked

- `--set fusion.rrf_k=20` returned sealable draft `sha256:1e9eaa…`.
- A two-variable file returned sorted fusion/retrieval mutations regardless of file order.
- `fusion.rrf_k=oops` returned `invalid_json` in a non-sealable draft.
- CLI and direct application compilation returned equal digests.
- Full pre-commit validation passed.

### What didn't work

- N/A.

### What I learned

- A pure snapshot derivation is stronger than writing to an in-memory store: the boundary performs literally zero store calls while preserving future durable identity.
- Splitting only the first equals sign preserves JSON strings such as `"a=b"`.

### What was tricky to build

- YAML scalar/object values must be converted to JSON bytes without guessing Go scalar types. Decoding the value node, then JSON marshaling it, preserves the core JSON-boundary contract.
- `--set` and `--mutations` are mutually exclusive to avoid unreviewed precedence or merging rules.

### What warrants a second pair of eyes

- Review YAML numeric conversion for very large integers if future artifact schemas accept values beyond Go's current decoded scalar ranges.
- Review one-row draft output for table ergonomics without weakening nested JSON output.

### What should be done in the future

- OPTKIT-018 can reuse `RequestedMutation` directly in HTTP compile envelopes.

### Code review instructions

- Start at `CompileManifestProposal`, then `proposal.go` in the CLI package.
- Run both `--set` and `--mutations` forms with `--format json`.

### Technical details

```text
input sources: exactly one of --set or --mutations
mutation file: strict YAML/JSON array of {variable,value}
output rows: 1 complete CandidateDraft
P5 start/done: printed true
```

## Step 8: Unify graph paths and run the real CLI no-write experiment

This phase centralized graph comparison/planning at the application boundary and proved that config comparison and proposal compilation return equal results for equivalent pipelines. It then built the real CLI and compiled the same proposal 100 times inside a watched directory.

The watched tree included only the manifest and mutation input. Its aggregate SHA-256 was equal before and after, and no SQLite, artifact, or journal path appeared.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Remove comparison drift risk and produce filesystem-level no-write evidence.

**Inferred user intent:** Prove purity through both architecture and observed command behavior.

**Commit (RAG-TTC):** `eaef2024e94c843b2e92bcbbea278922906fe23b` — "OPTKIT-016: unify proposal and config graph planning"

### What I did

- Added shared `CompareGraphs` and `PlanGraphs` application helpers.
- Routed manifest config comparison and proposal compilation through those helpers.
- Added an equality test for a `limit-1 → limit-2` proposal versus existing config diff/plan.
- Added `scripts/01-prove-cli-compilation-no-write.sh` and retained its exact output.
- Added `scripts/02-run-cli-contract-proof.sh` for deterministic nested-output assertions.
- Printed P6 start/done receipts.

### Why

- Equivalent before/after pipeline values must not produce two planning implementations.
- Core unit purity and built-command filesystem purity are independent evidence.

### What worked

- The same pipeline pair returned equal complete diff and plan values.
- All 100 CLI runs returned one stable draft digest.
- Watched tree hash stayed `ad2257c…` before and after.
- No durable database, artifact directory, or journal path existed.
- Focused and full pre-commit validation passed.

### What didn't work

- N/A.

### What I learned

- A command that has no store flag can still accidentally write caches or defaults relative to its working directory; the watched-directory test guards that observable behavior.

### What was tricky to build

- The proof binary and output had to live outside the watched directory so build/output writes did not invalidate the experiment. Only command side effects were measured.

### What warrants a second pair of eyes

- Review the watched-tree boundary and temporary-directory cleanup in the proof script.

### What should be done in the future

- Reuse these scripts after OPTKIT-017 to show the contrast: compile writes zero, seal writes canonical facts exactly once.

### Code review instructions

- Review shared helper call sites and run both ticket scripts from the workspace root.

### Technical details

```text
runs: 100
draft digest: sha256:1e9eaa7482f129cdd3069d0aa8cd016e600bbc6de6dccb2445d9133b871afb2d
before tree: ad2257c2be2188627b84d292ac194444c0fb4922af35c66edb43ef4737287e22
after tree:  ad2257c2be2188627b84d292ac194444c0fb4922af35c66edb43ef4737287e22
P6 start/done: printed true
```

## Step 9: Validate, close, and deliver OPTKIT-016

This phase ran fresh repository-wide and isolated validation after all production commits, reconciled the implementation guide with actual behavior, completed docmgr state, audited all physical work slips, and published the current guide/diary bundle.

The completed ticket leaves no durable write path in compilation. OPTKIT-017 now receives canonical normalized mutations, verified parent identity, deterministic digest, and explicit sealability.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Close every acceptance criterion with independent test, CLI, documentation, print, and delivery evidence before beginning sealing.

**Inferred user intent:** Produce an auditable implementation gate rather than declaring completion from focused tests alone.

### What I did

- Ran full Optkit CI, CGO/non-CGO tests/builds, race, lint, and vet.
- Ran full RAG-TTC lint, tests, vet, build, focused race suites, CLI proofs, and a 372-package dependency scan.
- Checked all six implementation tasks and updated guide, diary, changelog, relations, index, and parent roadmap.
- Dry-ran, rendered, and uploaded the completed implementation bundle.
- Printed the P7 done slip and audited one plan plus seven start/done pairs.
- Preserved the unrelated `optkit/numbergame-demo` directory and concurrent OPTKIT-021–024 history.

### Why

- Ticket completion requires code, observed behavior, documentation, and operational delivery to agree.
- OPTKIT-017 must start from a clean, independently reviewed compiler gate.

### What worked

- All fresh validation commands passed.
- CLI proof reproduced catalog IDs, draft digest, before/after graph IDs, direct fusion invalidation, and invalid diagnostic code.
- `docmgr doctor` passed and no implementation tasks remained.
- The completed PDF uploaded to `/ai/2026/08/26/OPTKIT-016`.
- Every required slip log contained `printed: true`/`printed: yes`.

### What didn't work

- N/A.

### What I learned

- The phase-slip audit is simplest and most reliable when filenames encode global sequence, ticket phase, and boundary type.

### What was tricky to build

- The P7 done receipt can only be printed after all phase evidence exists, but that receipt is itself required evidence. The closure sequence therefore validates code/docs/delivery first, prints P7 done, runs the slip audit, and commits the final receipt/audit afterward.

### What warrants a second pair of eyes

- Review the partial-preview sealability contract, semantic digest projection, strict mutation-file boundary, and pure snapshot derivation before implementing replay.

### What should be done in the future

- Begin OPTKIT-017 by recompiling drafts and replaying only normalized mutations through `Binding.Assign` and `PatchBuilder`.

### Code review instructions

- Review RAG-TTC commits `1e926542`, `4c90cd51`, `b7c4e428`, `5ce78767`, `fad49208`, and `eaef2024` in order.
- Review Optkit commits `759ef01f`, `f0dafc6f`, and `b45db558`.
- Run all three ticket scripts and the validation commands in the implementation guide.
- Run `docmgr doctor --ticket OPTKIT-016 --stale-after 30` from the workspace root.

### Technical details

```text
Optkit full CI/race/lint/CGO and non-CGO: pass
RAG-TTC full lint/test/vet/build: pass
focused race: pass
dependency packages: 372
import cycles: 0
implementation tasks: 6/6 complete
work slips: 15/15 printed
completed bundle: /ai/2026/08/26/OPTKIT-016
```
