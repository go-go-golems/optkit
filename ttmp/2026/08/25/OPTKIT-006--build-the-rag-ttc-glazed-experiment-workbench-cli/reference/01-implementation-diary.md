---
Title: Implementation Diary
Ticket: OPTKIT-006
Status: active
Topics:
    - implementation
    - optkit
    - rag-ttc
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://rag-ttc/cmd/rag-ttc/main.go
      Note: CLI composition root discovered during architecture mapping
    - Path: repo://rag-ttc/pkg/ttc/experimentworkbench/manifest_test.go
      Note: Manifest failures and identity evidence
    - Path: repo://rag-ttc/pkg/ttc/experimentworkbench/service_test.go
      Note: Durable run and resume evidence
ExternalSources: []
Summary: Chronological implementation record for the manifest-driven RAG-TTC Glazed experiment workbench.
LastUpdated: 2026-08-25T23:59:00Z
WhatFor: Preserve design reasoning, commands, failures, commits, validation evidence, and continuation guidance.
WhenToUse: Read before resuming OPTKIT-006 or reviewing its implementation.
---



# Diary

## Goal

Design and implement a manifest-driven Glazed CLI that validates, explains, runs, resumes, inspects, and verifies durable RAG-TTC experiments before the specialist UI is built.

## Step 1: Map the workbench and freeze the implementation plan

The first step established the workbench as the operator-facing proof layer between reviewed scientific intent and durable Optkit facts. I traced the current Glazed command into `optkitcampaign.Run`, then through snapshots, complete-block expansion, budgets, leases, trajectories, observations, estimates, and journal verification.

The design deliberately keeps RAG semantics in RAG-TTC, generic durability in Optkit, and transport concerns in Glazed commands. It separates side-effect-free validation and planning from campaign creation and resume, and it makes the future UI a consumer of already-characterized facts rather than the first implementation of their meaning.

### Prompt Context

**User prompt (verbatim):**

```text
ok, let's create a new ticket for the CLI, and implementthat first, that way we can already validate that what the UI is going to display is valid. Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


Btw, don't forget to commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)
```

**Assistant interpretation:** Create OPTKIT-006, research the existing experiment stack, write an intern-ready architecture and implementation guide, implement the manifest-driven Glazed workbench in coherent phases, validate it, maintain this diary, and publish the deliverables to reMarkable.

**Inferred user intent:** Establish a trustworthy operator workflow before UI implementation so command output can validate the same scientific facts and projections the browser will later display.

### What I did

- Created `OPTKIT-006 — Build the RAG-TTC Glazed Experiment Workbench CLI`.
- Added six implementation tasks covering architecture, manifest, config commands, campaign commands, end-to-end validation, and publication.
- Read the current Glazed command, experiment command root, campaign runner, prepared system, semantic fixture, layered graph, invalidation planner, Optkit query service, and generic Optkit CLI.
- Captured line-anchored evidence with `nl -ba` and repository-wide command discovery with `rg`.
- Wrote the initial intern guide with architecture maps, command APIs, manifest schema, flows, pseudocode, decisions, phases, tests, failure handling, risks, and file references.
- Printed the ticket plan slip.

### Why

- The current command runs a fixed embedded plan; implementation needs a precise contract before adding file-driven experiment setup.
- The guide must orient an engineer across three repositories/layers without encouraging RAG concepts to leak into Optkit.
- Separating plan from execution is necessary for CI safety and operator review.

### What worked

- Existing package boundaries are strong enough to support a thin CLI: `RunOptions` already accepts arms, cases, repeats, and an executor.
- The layered graph and invalidation planner already provide the exact semantics needed by config commands.
- Glazed's current command pattern is already used correctly by the small `optkit-rag` command.
- The deterministic semantic fixture provides an offline end-to-end implementation target.

### What didn't work

- One attempted read targeted a nonexistent file:

  ```text
  ENOENT: no such file or directory, access '/home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/cmd/rag-ttc/cmds/root.go'
  ```

  RAG-TTC's composition root is `cmd/rag-ttc/main.go`; command groups are registered directly from there.

### What I learned

- `campaign resume` can and should be manifest-free because `CampaignSpec` is already persisted in the creation event.
- The first workbench can support arbitrary deterministic fixture arms and cases without pretending to support every connected runtime.
- Graph plans describe semantic reuse but do not prove materialized artifact availability.
- Journal verification and direct event-payload verification are distinct claims and should have distinct output fields.

### What was tricky to build

- Scope had to be useful without overclaiming. The manifest can model all twelve layers, while the first executor only materializes the retrieval slice. The design explicitly labels the graph as planning metadata until later producers persist all stages.
- Existing `run --campaign` overloads create and resume. The design separates those authority boundaries while leaving compatibility decisions to implementation evidence.

### What warrants a second pair of eyes

- Whether v1 should persist the manifest and graph as campaign-reachable artifacts immediately or defer that to the projector ticket.
- Whether existing automation requires the old top-level `run` and `inspect` command paths.
- Whether YAML should be accepted alongside JSON in v1; the design chooses strict YAML because it is already a direct dependency and remains JSON-compatible.

### What should be done in the future

- Implement Phase 1 strict manifest loading and identity.
- Keep connected/provider-backed preparations out until their runtime and budget identities can be represented without secrets.

### Code review instructions

- Start with the design guide's repository map and decisions.
- Compare proposed service boundaries with `optkitcampaign.RunOptions` and `optimization.NewGraph`.
- Confirm every mutation command maps to an existing durable application path.

### Technical details

Evidence commands included:

```bash
rg -n "optkit|campaign|experiment" rag-ttc/cmd/rag-ttc/cmds rag-ttc/cmd/ttc-admin optkit/cmd -g '*.go'
nl -ba rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/command.go
nl -ba rag-ttc/pkg/ttc/optkitcampaign/campaign.go
nl -ba rag-ttc/pkg/ttc/optimization/graph.go
nl -ba rag-ttc/pkg/ttc/optimization/invalidation.go
nl -ba optkit/query/service.go
```

## Step 2: Implement strict manifests and durable attribution

This step introduced the versioned operator input and application service behind the CLI. A strict YAML manifest now defines preparation, arms, cases, repeats, and configuration layers, derives a content identity, and creates deterministic campaign options without importing Cobra or Glazed into domain code.

The durable campaign specification was extended with manifest and configuration-graph attribution. Resume continues to load the stored `CampaignSpec`; it never rereads a local manifest.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Implement the reviewed experiment contract and durable application boundary before adding terminal commands.

**Inferred user intent:** Ensure the CLI and future UI can identify exactly which reviewed configuration produced each arm and campaign.

**Commit (RAG-TTC code):** `1ebd2c0a6` — "Workbench: add versioned experiment manifests"

### What I did

- Added `pkg/ttc/experimentworkbench/manifest.go` with bounded reads, strict YAML fields, single-document enforcement, semantic identity, and validation.
- Added `service.go` with dry-run, graph comparison, invalidation, run, resume, status, and direct payload custody services.
- Added baseline and challenger semantic-fixture manifests.
- Added unit tests for malformed input, duplicate IDs, unsupported preparations, invalid graphs, cancellation, deterministic identity, no-mutation dry-run, diff, and invalidation.
- Added an integration test that runs, verifies, inspects, and resumes a six-episode campaign.
- Persisted manifest and graph attribution in `CampaignSpec`.
- Added YAML tags to reused typed contracts rather than decoding into untyped maps.

### Why

- Scientific setup needs a reviewable, hashable input rather than hard-coded defaults or an unreviewable flag list.
- Domain services must be reusable by CLI characterization tests and later projectors.
- Resume must be governed by persisted facts, not mutable operator files.

### What worked

- Strict decoding rejected unknown fields and multiple YAML documents.
- The fixture produced stable manifest and graph identities.
- Dry-run calculated six episodes, six retrieval-query units, and 600 maximum retrieval-result units without creating its store path.
- End-to-end execution completed six episodes, recorded paired estimates, verified direct payloads, and resumed without identity or count changes.
- Full repository pre-commit lint, Glazed vet, and tests passed.

### What didn't work

- The first draft of `manifest_test.go` accidentally contained a placeholder slice type while tests were being assembled. It was removed before the first test run and replaced with a real deep-copy helper for arms, cases, groups, and dependency slices.
- The first service integration test attempted to call unexported `optkitcampaign.loadCampaignSpec`. The test was corrected to find the `CampaignCreated` event and decode its typed payload through the public artifact API.

### What I learned

- YAML v3 does not use JSON tags as substitutes for YAML field names. Reused records need explicit YAML tags for `required_groups`, `forbidden_targets`, `value_schema`, and `depends_on` under strict decoding.
- A campaign creation event is the correct test seam for proving persisted operator attribution.
- Direct event payload custody can be verified generically; nested artifact traversal still requires schema-aware logic.

### What was tricky to build

- The loader must reject unknown fields and multiple documents while also enforcing a size ceiling. It uses a limited reader, `KnownFields(true)`, a second decode requiring `io.EOF`, and a file-size check.
- Tests must deep-copy dependency slices before mutation; otherwise one validation scenario can alter another, repeating the Phase 2 graph aliasing class of bug.

### What warrants a second pair of eyes

- Review the manifest size ceiling and YAML parsing policy.
- Confirm direct payload verification is labeled narrowly enough.
- Review whether persisted configuration graphs should remain inline in `CampaignSpec` or move to separately referenced CAS artifacts when graph size grows.

### What should be done in the future

- Add schema-aware transitive custody traversal.
- Add connected preparations only after runtime, profile, cache, and budget identity is explicit.

### Code review instructions

- Start with `experimentworkbench/manifest.go`, then `service.go`.
- Review the `CampaignSpec` additions in `optkitcampaign/campaign.go`.
- Run `go test ./pkg/ttc/experimentworkbench ./pkg/ttc/optkitcampaign -count=1`.

### Technical details

The initial baseline contained two arms, three cases, one repeat, and six complete-block episodes. The first application test passed full run, status, custody, and idempotent resume.

## Step 3: Correct configuration ownership to per-arm graphs

While implementing config commands, I found an important modeling defect in the initial guide and manifest: one experiment-level graph cannot truthfully describe two arms with different retrieval limits. Since the arm is the unit compared by the complete-block trial, each arm must bind its executable config to its own complete graph.

The implementation was corrected before publication. Every arm now carries a graph, validation recomputes the retrieval-layer identity from the executable `RetrievalConfig`, and campaign facts persist a graph map keyed by arm ID.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Ensure configuration diff and UI provenance describe the actual values executed by each arm.

**Inferred user intent:** Prevent a superficially polished CLI from emitting scientifically inconsistent configuration claims.

### What I did

- Replaced the single manifest `layers` collection with `ManifestArm.Layers`.
- Replaced `LoadedManifest.Graph` with `Graphs map[string]optimization.Graph`.
- Added explicit arm selection to comparison and invalidation APIs.
- Recomputed each retrieval `ConfigRef` from `arm.Config` and rejected identity mismatches.
- Persisted `ConfigGraphs map[string]optimization.Graph` in `CampaignSpec`.
- Moved canonical operator examples to `assets/configs/experiments/optkit-rag/`.
- Generated distinct content-derived retrieval identities for limits 1, 2, and 3.
- Added a graph/config mismatch test.

### Why

- An experiment-level graph would claim one retrieval identity while executing two different retrieval configs.
- The future comparison UI needs graph identity per arm, not merely per manifest.
- Validation must connect declared provenance to executable values rather than trusting both independently.

### What worked

- Baseline arm `limit-1`, baseline arm `limit-2`, and challenger arm `limit-3` resolve to distinct graph IDs.
- Comparing `limit-2` with `limit-3` reports exactly one direct retrieval change, five reusable upstream layers, and seven recomputed layers.
- Campaign creation stores both arm graphs and the manifest ID.

### What didn't work

- The first CLI comparison test requested challenger arm `limit-3`, but the challenger YAML still named that arm `limit-2` after changing only its limit. The error was precise:

  ```text
  manifest experiment-manifest:f97b77a8f378fb1336f1f5d93533a6b8a9481c98acf5e2a801bab686c0aa7883 has no arm "limit-3"
  ```

  The challenger arm was renamed to `limit-3` and its content-derived retrieval identity was regenerated.

- One manifest test still referenced the removed singular `first.Graph.ID` and failed to compile:

  ```text
  first.Graph.ID undefined (type func(arm string) (optimization.Graph, error) has no field or method ID)
  ```

  It was updated to assert `first.Graphs["limit-1"].ID`.

### What I learned

- Configuration identity belongs at the treatment-arm boundary.
- A graph can be structurally valid yet scientifically false unless leaf identities are checked against executable typed values.
- Slight duplication in a reviewed manifest is preferable to an ambiguous shared graph in v1; authoring helpers can reduce repetition later.

### What was tricky to build

- Changing a retrieval identity requires rewiring fusion's direct dependency while preserving fusion's local identity. Resolved digests then propagate the change downstream.
- Manifest semantic identity must include all per-arm graphs in stable arm order, while runtime lookup benefits from a map. The manifest slice remains identity input; the map is derived.

### What warrants a second pair of eyes

- Review whether future typed layers should all receive executable-value identity checks as their producers are implemented.
- Review map serialization in persisted `CampaignSpec`; Go JSON orders string keys deterministically, while semantic manifest identity remains slice-based.

### What should be done in the future

- Add typed value documents and identity checks layer by layer as execution expands beyond retrieval.
- Consider a `config derive` authoring command instead of hand-editing content IDs.

### Code review instructions

- Compare `ManifestArm` and `ValidateManifest` with `optimization.NewConfigRef`.
- Inspect both checked-in manifests and verify retrieval IDs differ by limit.
- Run the mismatch validation test and the retrieval invalidation test.

### Technical details

```text
limit=1 -> config:2b5e54c46f90146b11119d3022a1c71f66c23285a3d951ee916f54b4d047f4d4
limit=2 -> config:1c067430c194d35f2a847456a0231d98fca20dd89c23bd7c3a183f25564e4907
limit=3 -> config:b37158680ca7a84360fe0199769fc7ffea5648b09654674f0f3f06c0d733a14a
```

## Step 4: Add the Glazed operator command tree

This step made the application services available as a coherent operator workbench. The public command tree now separates configuration understanding from campaign mutation, and it separates new campaign creation from stored-specification resume.

The original fixed-fixture `run` and `inspect` commands remain hidden deprecated aliases so existing callers do not break while new documentation presents only the explicit workbench.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Implement a polished structured-output CLI that proves configuration and campaign facts before UI work.

**Inferred user intent:** Give operators and automated tests a trustworthy way to exercise the backend and inspect what the future UI will render.

**Commit (RAG-TTC code):** `099dd919b` — "Workbench: add Glazed experiment operator commands"

### What I did

- Added `config inspect`, `validate`, `diff`, and `plan`.
- Added `campaign dry-run`, `run`, `resume`, `status`, and `verify`.
- Added `--before-arm` and `--after-arm` to make comparison subjects explicit.
- Emitted normalized Glazed rows with stable cardinality.
- Expanded campaign summaries with lifecycle version, event and episode counts, terminal failures, budgets, and nullable paired delta.
- Added command-construction, help-flag, row-cardinality, no-mutation, processor-error, and service integration tests.
- Hid and deprecated the overloaded legacy paths.
- Ran a real CLI campaign through create, status, custody verification, and resume.

### Why

- Commands should prove the same facts projectors will later serve.
- Structured rows support humans, shell automation, and contract fixtures without separate serializers.
- Explicit command intent prevents resume from accidentally accepting a mutable manifest or reset flag.

### What worked

- `config inspect` emits 24 rows: 2 arms × 12 layers.
- `config validate` emits one row with two graph IDs and six episodes.
- `config diff` emits one direct retrieval change.
- `config plan` emits five reuse and seven recompute rows.
- `campaign dry-run` emits exact budgets and creates no path.
- A real CLI run completed six episodes with paired delta `0.16666666666666666`.
- Status, direct custody verification, and resume returned the same campaign, trial, version, and completion count.
- Full repository pre-commit lint, Glazed vet, and tests passed.

### What didn't work

- The first import edit for `command_test.go` accidentally left literal `},{` after the `types` import, causing:

  ```text
  command_test.go:13:44: missing import path
  ```

  The malformed token was removed and the file was formatted.

- The test helper initially used `command.Schema`, but `cmds.GlazeCommand` exposes it through `command.Description().Schema`:

  ```text
  command.Schema undefined (type cmds.GlazeCommand has no field or method Schema)
  ```

  The helper was corrected to use the public command description.

### What I learned

- Short help intentionally omits the structured-output section; `--help --long-help` exposes the three universal flags.
- Glazed row cardinality should be explicit in tests because formatters differ in buffering but share emitted records.
- Keeping legacy commands hidden allows compatibility without advertising the ambiguous interface.

### What was tricky to build

- Nested command groups still need each `GlazeCommand` built through `cli.BuildCobraCommandFromCommand`; raw Cobra should only own grouping.
- Table-friendly rows need scalar columns, while budgets and graph ID collections remain nested fields for JSON. The summary keeps both stable identifiers and bounded aggregate data.

### What warrants a second pair of eyes

- Review whether nested budget resources should become a separate `campaign budget` row stream in a later ticket.
- Confirm hidden legacy aliases can be removed after caller search and migration.
- Review help discoverability before adding embedded Glazed help entries.

### What should be done in the future

- Add projector/view commands only after their server-side DTOs exist.
- Add measurement commands when answer campaigns become durable workbench preparations.

### Code review instructions

- Run `rag-ttc experiment optkit-rag --help`.
- Exercise config validate/inspect/diff/plan against both checked-in manifests.
- Run the manifest-driven campaign in a temporary store, then status, verify, and resume it.

### Technical details

```text
rag-ttc experiment optkit-rag config {inspect,validate,diff,plan}
rag-ttc experiment optkit-rag campaign {dry-run,run,resume,status,verify}
```

## Step 5: Characterize the full workbench

The final implementation step added one reproducible script that builds the CLI and checks every operator promise: strict manifest validation, per-arm graph inspection, direct diff, invalidation, no-mutation dry-run, durable run, status, custody, idempotent resume, and Glazed help.

The archived transcript is intended to be the first backend characterization fixture for the later projector and UI ticket.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Validate the implemented workbench end to end and preserve machine-checkable evidence.

**Inferred user intent:** Make UI work depend on proven backend behavior rather than assumptions.

### What I did

- Added `scripts/01-validate-workbench.sh`.
- Built a clean CLI binary with `GOWORK=off`.
- Ran focused unit and race tests.
- Parsed every command's JSON output with `jq`.
- Asserted no store exists after dry-run.
- Asserted run and resume preserve campaign, trial, version, and six completions.
- Asserted short and long help expose domain and structured-output flags appropriately.
- Archived `sources/01-workbench-validation.txt`.

### Why

- Direct command execution catches wiring and formatter behavior that package tests cannot.
- Explicit PASS markers make regressions and handoff review cheap.
- The UI ticket can use the same baseline/challenger facts as known-good inputs.

### What worked

- All ten acceptance markers passed.
- Focused race tests passed.
- The final campaign had six completed episodes and no terminal failures.
- Journal and all 47 unique directly referenced event payloads verified.

### What didn't work

- The first validation run called `campaign run --long-help` without `--help`. In this Glazed/Cobra integration, `--long-help` only modifies help rendering and does not itself request help, so the command executed and failed validation:

  ```text
  Error: manifest path is required
  Run 'rag-ttc [command] --help' for usage.
  ```

  The validator now calls `campaign run --help --long-help`, which renders the structured-output section without executing.

### What I learned

- Help-modifier flags must be tested through the exact invocation users need.
- Backend characterization is already rich enough to validate campaign, comparison, and config-diff projector inputs.
- It is not yet rich enough to validate full context, answer, or judge UI stages; those remain correctly gated.

### What was tricky to build

- Random campaign IDs prevent golden whole-file output. The validator extracts IDs and checks relational invariants instead of comparing unstable bytes.
- Resume may legitimately append reconciliation facts after an interrupted campaign, but a completed campaign must preserve version. The fixture asserts the stronger completed-campaign invariant.

### What warrants a second pair of eyes

- Review whether the validation script should enter CI now or remain ticket-local until connected preparations are added.
- Confirm the 100-results-per-episode conservative budget ceiling remains aligned with campaign initialization.

### What should be done in the future

- Reuse the JSON command outputs as projector characterization inputs.
- Extend validation when graph artifacts receive transitive custody checks.

### Code review instructions

- Run `scripts/01-validate-workbench.sh` from any directory.
- Inspect `sources/01-workbench-validation.txt` for every PASS marker.

### Technical details

```text
MANIFEST_VALIDATION=PASS
PER_ARM_GRAPH_INSPECTION=PASS
DIRECT_CONFIG_DIFF=PASS
INVALIDATION_PLAN=PASS
DRY_RUN_NO_MUTATION=PASS
MANIFEST_CAMPAIGN_RUN=PASS
CAMPAIGN_STATUS=PASS
CAMPAIGN_CUSTODY=PASS
IDEMPOTENT_RESUME=PASS
GLAZED_HELP_SURFACE=PASS
OPTKIT_006_WORKBENCH_VALIDATION=PASS
```

## Step 6: Publish the intern guide and close the ticket

The final step audited the documentation against the implemented per-arm model, preserved validation evidence, and uploaded the guide plus diary as one reMarkable bundle with a table of contents.

The ticket now contains a 1,000-line intern guide, a chronological implementation diary, an executable validator, and the exact output from a successful end-to-end campaign.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Finish publication and ticket bookkeeping only after implementation and validation are complete.

**Inferred user intent:** Leave a durable teaching and review package that can be read away from the repository and used to start the UI work safely.

### What I did

- Updated the guide from the initial single-graph proposal to the implemented per-arm graph contract.
- Added final command paths, API signatures, file locations, decision records, failure modes, and validation results.
- Ran `docmgr doctor` and frontmatter/placeholder checks.
- Ran a mandatory reMarkable dry-run.
- Uploaded the guide and diary bundle to `/ai/2026/08/25/OPTKIT-006`.
- Prepared all tasks and changelog entries for closure.

### Why

- Published architecture must describe the implementation, including corrections discovered during coding.
- The reMarkable bundle gives a new engineer one navigable document containing both the intended system and its actual construction history.

### What worked

- Dry-run selected the two intended documents in the correct order.
- The corrected bundle rendered and uploaded successfully:

  ```text
  OK: uploaded OPTKIT-006 RAG-TTC Experiment Workbench CLI.pdf -> /ai/2026/08/25/OPTKIT-006
  ```

- `docmgr doctor` passed cleanly before publication.

### What didn't work

- The first real upload failed in Pandoc/LaTeX:

  ```text
  Error: pandoc failed: Error producing PDF.
  ! Undefined control sequence.
  l.2052 ... ticket and the nupload to remarkable.\n
  ```

  Step 1 had encoded the multiline verbatim user prompt with literal `\n` sequences inside a quoted line. LaTeX interpreted `\n` as a control sequence. I replaced it with a fenced text block containing the actual blank lines, which is both more faithful to the prompt and safe for rendering.

### What I learned

- Diary prompt fidelity and document rendering agree when multiline prompts use fenced text rather than escaped newline text.
- Dry-run validates selection and destination but does not invoke Pandoc, so a real render can still reveal Markdown-to-LaTeX problems.

### What was tricky to build

- The guide had to be reconciled after implementation corrected the graph ownership model. Search-based auditing for singular `LoadedManifest.Graph`, old testdata paths, and missing arm selectors prevented stale examples from surviving publication.

### What warrants a second pair of eyes

- Review the reMarkable PDF's long tables and code blocks for comfortable reading.
- Review the final guide's distinction between persisted planning graphs and actually materialized stages.

### What should be done in the future

- Start the projector ticket with the CLI characterization as its backend oracle.
- Add an embedded Glazed help page if operators need the full guide discoverable inside the binary.

### Code review instructions

- Read the design guide first and this diary second.
- Run `scripts/01-validate-workbench.sh`.
- Compare CLI JSON with future projector DTO fixtures.

### Technical details

```text
reMarkable destination: /ai/2026/08/25/OPTKIT-006
bundle: OPTKIT-006 RAG-TTC Experiment Workbench CLI.pdf
contents: intern guide + implementation diary
```
