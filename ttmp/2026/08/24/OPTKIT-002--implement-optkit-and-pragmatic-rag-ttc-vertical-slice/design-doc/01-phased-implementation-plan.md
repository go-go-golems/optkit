---
Title: Phased Implementation Plan
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
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/09-p8-boundary-and-migration-inventory.md
      Note: P8 completion outcome and remaining parity gates
    - Path: ws://coinvault/internal/knowledge/semantic_fixture_test.go
      Note: P1 authorization and evidence reference laws
    - Path: ws://coinvault/internal/knowledge/service.go
      Note: Cross-product RAG behavioral reference
    - Path: ws://judgekit/judging/claimjudge.go
      Note: P7 measurement integration boundary
    - Path: ws://rag-ttc/README.md
      Note: Final product integration and retained legacy boundaries
    - Path: ws://rag-ttc/pkg/ttc/search/search.go
      Note: P2 extraction starting point
    - Path: ws://rag-ttc/pkg/ttc/search/semantic_fixture_test.go
      Note: P1 TTC retrieval and evidence characterization
    - Path: ws://ragkit/rag/retrieval/retrieval.go
      Note: Retained reusable RAG-domain boundary
    - Path: ws://sources/optkit-implementation-source.zip
      Note: P0 source baseline and provenance
ExternalSources: []
Summary: Executable phase plan for landing the Optkit baseline and using a Coinvault-informed, RagKit-backed RAG-TTC vertical slice to prove durable experiments without folding RAG primitives into Optkit.
LastUpdated: 2026-08-24T22:45:00-04:00
WhatFor: Define phase boundaries, concrete tasks, tests, commits, work-slip checkpoints, and acceptance gates for OPTKIT-002.
WhenToUse: Before starting or reviewing any implementation phase in the Optkit and RAG-TTC vertical-slice program.
---




# Phased Implementation Plan

## 1. Objective

Deliver a working, restartable optimization vertical slice in which:

1. Optkit owns domain-neutral configuration snapshots, episodes, measurements, experiments, campaigns, scheduling, budgets, and provenance.
2. RagKit remains a separate reusable RAG-domain library.
3. RAG-TTC owns TTC-specific source policy, retrieval composition, evidence admission, answer behavior, benchmarks, and application services.
4. Coinvault supplies proven behavioral patterns and semantic fixtures, not copied product code.
5. Judgekit remains a separately reusable measurement provider integrated through explicit Optkit instruments.
6. RagOpt orchestration paths are replaced by Optkit where the vertical slice proves equivalent behavior.

The implementation should produce product value before repository consolidation. There is no requirement to move RagKit algorithms into Optkit.

## 2. Stable ownership boundaries

```text
Optkit (domain neutral)
  record/artifact/space
  episode/measure/experiment
  campaign/scheduler/budget/projection
  system registration and preparation contracts

RagKit (RAG domain)
  document/chunk/representation
  index bundles and content stores
  lexical/vector search
  collapse/fusion/reranking/hydration
  optional generic answering primitives

RAG-TTC (product)
  corpus admission and role policy
  TTC retrieval routes and configuration
  customer/admin application services
  evidence limits and answer contracts
  benchmark cases and promotion rules

Judgekit (measurement domain)
  claim extraction and judging protocols
  provider execution and cache policy
  attributable reports

Coinvault (product/reference)
  retained product service
  cross-product semantic fixtures
  reference for runtime identity and diagnostics
```

### Integration direction

```text
Optkit campaign
  -> RAG-TTC system executable
       -> RAG-TTC retrieval service
            -> RagKit algorithms
       -> Optkit trajectory sink
  -> deterministic Optkit instruments
  -> optional Judgekit instrument
```

Optkit does not import product packages. Product composition roots import Optkit and RagKit.

## 3. Global working protocol

Every phase follows the same sequence:

1. Print a **phase plan** brutalist work slip before code changes.
2. Read the prior diary and current task state.
3. Establish a failing or characterizing test before semantic refactoring.
4. Implement the smallest complete vertical change.
5. Format, run focused tests, then run repository-level tests.
6. Commit production/test changes at a coherent boundary.
7. Update the diary using the strict diary format, relate files, check the phase task, and update the changelog.
8. Commit ticket documentation.
9. Print a **phase status** brutalist work slip with commit, work done, tricky points, and next phase.
10. Do not start the next phase until the current acceptance gate is green.

No backwards-compatibility layer is added unless a concrete caller requires a transition. Thin model-tool and transport boundaries are permanent application boundaries, not legacy adapters.

## 4. Commit policy

Use focused commits:

- provenance import;
- repository normalization;
- one semantic service/refactor unit;
- one evaluation or measurement feature;
- phase diary/bookkeeping;
- deletion of superseded paths.

Do not combine cross-repository changes in one Git commit. Record all related commit hashes in the Optkit diary.

## 5. Phase summary

| Phase | Outcome | Primary repository | Acceptance signal |
|---|---|---|---|
| P0 | supplied Optkit implementation is the actual repository | Optkit | full, race, and static-link tests green |
| P1 | deterministic shared RAG semantic fixtures | RAG-TTC, Coinvault | fixture laws pass without provider credentials |
| P2 | canonical direct TTC retrieval service | RAG-TTC | direct/tool parity and full tests green |
| P2.5 | read-only campaign query plane and UI navigation contract | Optkit | replayable projections, search, deep links, and SSE |
| P3 | complete runtime identity, routes, and source policy | RAG-TTC | configured behavior and trace attribution proven |
| P4 | stage-aware deterministic retrieval evaluation | RAG-TTC | target-loss diagnosis is machine readable |
| P5 | direct answer application service and contracts | RAG-TTC | direct/served domain parity proven |
| P6 | durable multi-arm Optkit RAG-TTC campaign | Optkit + RAG-TTC | restartable campaign and estimates |
| P7 | lightweight Judgekit measurement integration | Judgekit + RAG-TTC | remeasurement without retrieval rerun |
| P8 | stable package boundaries and old orchestration deletion | all relevant | no superseded RagOpt/product runner path |
| Final | validation, docs, and delivery | Optkit | doctor clean, all tests green, bundle uploaded |

## 6. P0 — Import and normalize the supplied Optkit baseline

### Goal

Replace the template module with the supplied green Optkit source while preserving repository-owned CI/release plumbing and ticket history. Keep archive provenance explicit and avoid semantic refactoring in the import commit.

### Entry conditions

- Archive revision is known: `1786d1da86c9e03316ed71336bbc993fa30531f0`.
- Archive tests pass from an isolated extraction.
- Current Optkit working tree is clean.
- OPTKIT-001 and OPTKIT-002 documentation is committed.

### Tasks

1. Print the P0 plan slip.
2. Snapshot archive file list and checksums under the ticket `sources/` directory.
3. Extract the archive to a temporary directory, never over the Git worktree.
4. Classify collisions:
   - keep repository `AGENT.md`, CI, GoReleaser, lint, lefthook, logcopter, license, and `ttmp`;
   - import archive domain packages, command, README content, Makefile behavior, `go.mod`, and `go.sum`;
   - remove template `cmd/XXX` and placeholder `pkg` only when replaced.
5. Copy archive implementation into the worktree with an explicit include/exclude list.
6. Normalize module path to `github.com/go-go-golems/optkit`.
7. Merge Makefile targets rather than silently dropping repository checks.
8. Update README to describe real packages and record source provenance.
9. Run `gofmt`/`go fmt`, `go mod tidy`, and generation checks where applicable.
10. Run:
    - `go test ./... -count=1`;
    - `go test -race ./... -count=1`;
    - `CGO_ENABLED=0 go test ./... -count=1`;
    - `go vet ./...`;
    - CLI demo and verification smoke tests.
11. Commit the source import separately from normalization fixes when possible.
12. Update diary, relations, task, and changelog; print P0 status slip.

### Expected file impact

```text
artifact/ budget/ campaign/ episode/ experiment/
local/ measure/ numbergame/ projection/ record/
scheduler/ space/ store/ cmd/optkit/
go.mod go.sum README.md Makefile
```

### Acceptance gate

- Checked-out module path is correct.
- All archive packages and tests are present.
- Full, race, static-link, and vet checks pass.
- Numbergame demo creates a store; inspect/verify commands read it.
- Repository CI/release files and both tickets remain present.
- Git history identifies archive provenance.

### Rollback

Revert the import commit. No external data migration occurs in P0.

## 7. P1 — Freeze cross-product semantic RAG fixtures

### Goal

Create deterministic fixtures that capture the Coinvault laws to preserve and the current RAG-TTC behavior to refactor. Tests must run without external provider credentials.

### Tasks

1. Print the P1 plan slip.
2. Inventory existing Coinvault and RAG-TTC test corpora and helper builders.
3. Define a small shared semantic fixture schema containing:
   - stable documents, chunks, representations, roles, and expected bundle identity;
   - lexical-only, vector-only, and hybrid queries;
   - one comparison/multi-source query;
   - one role-filtered negative query;
   - one no-answer query;
   - one malformed output and one provider failure.
4. Add expected stage candidate IDs and rankings.
5. Add expected evidence labels, item/rune limits, and omission counts.
6. Add expected answer citation and abstention outcomes.
7. Characterize current RAG-TTC tool behavior without changing semantics.
8. Express Coinvault reference laws as table tests rather than importing Coinvault internals.
9. Add fixture version and semantic digest.
10. Document intentional differences between products.
11. Run focused and full product tests.
12. Commit each repository separately, update diary, and print P1 status slip.

### Acceptance gate

- Fixtures are deterministic and credential-free.
- Repeated runs produce identical candidate order and evidence labels.
- Both products' relevant tests pass.
- Product-specific scopes, prompts, and SQL behavior are not forced into a common package.

## 8. P2 — Extract the canonical RAG-TTC retrieval service

### Goal

Move semantic retrieval out of the model-facing `SearchTool` into one TTC-owned direct service used by both model tools and evaluation.

### Tasks

1. Print the P2 plan slip.
2. Add parity tests around `SearchTool.RunRoute` before moving behavior.
3. Introduce a TTC retrieval package with explicit request, result, route observation, and failure types.
4. Move channel execution, representation collapse, weighted RRF, route augmentation, source resolution, and hydration into the service.
5. Keep conversation-scoped evidence admission outside the reusable retrieval service.
6. Make `SearchTool` a permanent model boundary that invokes the service; do not retain a second retrieval implementation.
7. Preserve fresh ledger construction per session.
8. Add direct-service tests for lexical, vector, hybrid, augmentation, failure, and cancellation paths.
9. Add compile-time interface assertions where interfaces are used.
10. Remove superseded private search helpers once the new path is authoritative.
11. Run focused race tests and the full RAG-TTC suite.
12. Commit code, update diary, and print P2 status slip.

### Acceptance gate

- Direct service and model tool yield the same ranked chunks for fixtures.
- There is only one implementation of channel/collapse/fusion semantics.
- The direct service has no Geppetto, HTTP, WebSocket, or sessionstream dependency.
- Evidence remains conversation scoped.

## 8A. P2.5 — Build the read-only campaign query plane and UI navigation contract

### Goal

Enable UI development against real Numbergame campaigns without moving campaign creation, configuration, compilation, or lifecycle mutation into the browser. LLM agents and scientists use the CLI for writes; the UI is an addressable explorer for navigation, search, replay, comparison, and visualization.

### Tasks

1. Print the P2.5 plan slip.
2. Define versioned read-model contracts for campaign summaries, lineage, episode matrices, event timelines, trajectories, measurements, estimates, decisions, budgets, and artifact previews.
3. Add a query service that rebuilds projections from the journal and resolves immutable artifact references without exposing SQLite tables.
4. Add read-only `net/http` endpoints:
   - campaign list and campaign overview;
   - campaign lineage graph;
   - case-by-arm episode matrix;
   - paginated control events after a sequence;
   - episode detail and trajectory;
   - measurements, estimates, decisions, and budgets;
   - bounded artifact metadata/JSON preview;
   - search across IDs, kinds, schemas, actors, tags, and indexed text projections.
5. Add replayable SSE for campaign events using `after` sequence cursors; SSE carries facts, not browser commands.
6. Define stable URL routes and deep links for campaigns, candidates, snapshots, trials, episodes, observations, estimates, decisions, and artifacts.
7. Generate checked-in Numbergame API fixtures for frontend work and Storybook/mock-server use.
8. Add navigation-oriented projections:
   - stage rail and campaign status;
   - provenance/lineage graph;
   - episode matrix;
   - event timeline with replay slider;
   - trajectory span tree;
   - measurement table and paired-delta chart;
   - budget utilization;
   - artifact/evidence drawer.
9. Add query pagination, response size limits, artifact sensitivity checks, and cancellation.
10. Add tests proving the server exposes no mutation routes and cannot alter journal head, budgets, queue state, or artifacts.
11. Document CLI-to-UI handoff: agents print or return campaign IDs/URLs after creating or changing campaigns.
12. Run full/race tests, commit code, update diary, and print the P2.5 status slip.

### Acceptance gate

- A completed or running Numbergame campaign can be navigated entirely through stable deep links.
- The UI can reconnect and replay events from an exact journal sequence.
- Search reaches campaigns and evidence without coupling to authoritative storage tables.
- Checked-in API fixtures allow frontend work without a running backend.
- The HTTP surface has no create, update, delete, start, pause, resume, stop, or decision mutation endpoint.
- All writes remain CLI/application-command responsibilities and appear in the UI through journal/projection updates.

## 9. P3 — Add runtime identity, route compilation, and policy boundaries

### Goal

Make every retrieval result attributable to the exact prepared pipeline and ensure model input cannot widen corpus or role policy.

### Tasks

1. Print the P3 plan slip.
2. Define runtime identity fields for bundle, corpus, resolved config, query transform, retrieval policy, evidence policy, and reranker.
3. Compute identities during verified service preparation.
4. Compile checked-in intent-routing configuration into prepared service routes.
5. Reject unresolved representation, source-role, augmentation, and route references at preparation.
6. Define server-selected allowed roles/source policy.
7. Filter before fusion and before external reranking when policy is active.
8. Recheck before hydration as an internal invariant.
9. Add optional reranker identity, bounded candidate pool, and fused-order fallback.
10. Record effective result limit and its source.
11. Emit stable stage trace records with candidate artifact references.
12. Add tests proving:
    - model input cannot widen roles;
    - disallowed candidates do not influence fusion;
    - configuration routes are active behavior;
    - identity changes exactly when semantics change;
    - reranker failure degrades visibly.
13. Run focused and full tests, commit, update diary, and print P3 status slip.

### Acceptance gate

A trace consumer can reject intended/observed mismatches for bundle, route, transform, policy, reranker, evidence policy, and limit without relying on ambient flags.

## 10. P4 — Build stage-aware deterministic retrieval evaluation

### Goal

Produce strict, machine-readable evaluation that can explain where required evidence disappeared.

### Tasks

1. Print the P4 plan slip.
2. Define a versioned TTC retrieval suite with explicit case IDs and roles.
3. Separate positive, authorization-negative, and answer/judge-only case modes.
4. Define required evidence groups and target-resolver identity.
5. Record stage rankings:
   - raw lexical/vector;
   - collapsed;
   - policy filtered;
   - fused;
   - reranked;
   - returned;
   - admitted.
6. Diagnose first target-loss stage and rank movement.
7. Preserve query failures as rows with explicit denominator policy.
8. Compute recall, precision, MRR, nDCG, group coverage, diversity, latency, and provider-call metrics.
9. Emit canonical JSON and concise Markdown projections.
10. Add treatment-verification checks comparing intended and observed identities.
11. Add deterministic CI fixture and larger local suite.
12. Integrate one existing answer-quality runner path without rewriting unrelated stages.
13. Run tests, commit, update diary, and print P4 status slip.

### Acceptance gate

Given any failed positive case, the report identifies whether the target was absent, collapsed, filtered, lost in fusion, lost in reranking, below return limit, or rejected by evidence admission.

## 11. P5 — Extract the direct answer application service and contracts

### Goal

Run the same customer-domain turn below transport for serving and evaluation, with explicit grounding, citation, abstention, and failure outcomes.

### Tasks

1. Print the P5 plan slip.
2. Characterize current customer composition and transport events.
3. Define the direct customer turn request/result and event sink.
4. Reuse the canonical retrieval service and fresh per-turn/session evidence ledger.
5. Validate structured answer output and citation labels.
6. Represent safe abstention explicitly.
7. Separate retrieval miss, policy removal, admission truncation, generation failure, malformed output, unsupported claim, unresolved citation, and presentation failure.
8. Make HTTP/sessionstream invoke the direct application service.
9. Add direct-versus-served domain parity tests.
10. Keep customer and admin application services separate.
11. Add redaction tests before durable trajectory use.
12. Run focused race and full tests, commit, update diary, and print P5 status slip.

### Acceptance gate

Direct execution and transport execution produce equivalent domain events and answer-contract outcomes for deterministic fixtures.

## 12. P6 — Register and run a durable Optkit RAG-TTC campaign

### Goal

Use the direct TTC retrieval service as the first real product system in Optkit and run a restartable fixed multi-arm experiment.

### Tasks

1. Print the P6 plan slip.
2. Add the minimal domain-neutral system/preparation registry to Optkit.
3. Define typed TTC retrieval configuration and case codecs in product-owned integration code.
4. Keep provider handles/credentials out of snapshots.
5. Define a deliberately small fixed arm set.
6. Expand a complete-block trial over deterministic fixtures.
7. Convert retrieval stages into Optkit episode events/artifacts.
8. Add deterministic instruments and explicit measurement epochs.
9. Add per-arm and paired-difference estimates.
10. Use local SQLite, filesystem CAS, leases, and budgets.
11. Test cancellation and duplicate completion.
12. Kill/restart after lease, terminal result, and observation boundaries.
13. Add CLI commands or Glazed command group to run and inspect the campaign.
14. Run full tests in both repositories, commit separately, update diary, and print P6 status slip.

### Acceptance gate

One command starts or resumes a local multi-arm TTC retrieval campaign; after process restart, the journal verifies and projections/estimates rebuild without duplicate semantic execution.

## 13. P7 — Integrate Judgekit with lightweight research attribution

### Goal

Measure answer-level constructs without adding adversarial custody or rerunning retrieval when only the judge changes.

### Tasks

1. Print the P7 plan slip.
2. Implement or confirm restricted evidence-hidden claim extraction input.
3. Recompute instance identity at the execution boundary rather than trusting stale caller state.
4. Define an Optkit instrument backed by Judgekit.
5. Build instances from sealed answer trajectories and admitted evidence.
6. Record contract, protocol, prompt version/rendered digest, expected/observed model, and cache mode.
7. Use cache bypass for repeatability probes.
8. Store Judgekit reports as artifacts and map scores/failures to typed observations.
9. Preserve deterministic contract measurements separately.
10. Prove historical trajectories can be remeasured under a new measurement epoch without rerunning retrieval.
11. Add failure and missing-output tests.
12. Do not add signatures, keys, custody services, or immutable typestate.
13. Run full tests, commit separately, update diary, and print P7 status slip.

### Acceptance gate

A sealed historical answer episode can produce a second judge report and new observations under a distinct epoch while old observations remain unchanged and attributable.

## 14. P8 — Stabilize Optkit–RagKit boundaries and delete superseded orchestration paths

### Goal

Make repository ownership intentional without moving RagKit primitives into Optkit. Delete only duplicated orchestration and product runner paths proven unnecessary by the vertical slice.

### Tasks

1. Print the P8 plan slip.
2. Inventory final imports and classify each package as domain-neutral, RAG-domain, measurement-domain, or product-specific.
3. Add architecture tests or dependency guards:
   - Optkit cannot import RagKit, RAG-TTC, Coinvault, or Judgekit core packages;
   - RagKit cannot import product packages or Optkit campaign orchestration;
   - product integration may import both.
4. Keep RagKit algorithms and bundle types in RagKit.
5. Keep Judgekit reusable unless a separate consolidation decision is made.
6. Move only genuinely domain-neutral shared contracts if duplication is proven.
7. Replace exercised RagOpt candidate/trial/run/gate paths with Optkit equivalents.
8. Switch all callers before deleting each old path.
9. Remove superseded RAG-TTC custom orchestration paths covered by the new campaign.
10. Preserve product-specific reports/projections that still have users.
11. Update module dependencies and workspace configuration.
12. Run repository and cross-repository fixture suites.
13. Document retained boundaries and remaining migration backlog.
14. Commit deletions separately, update diary, and print P8 status slip.

### Acceptance gate

- RagKit remains independently testable and reusable.
- Optkit remains domain neutral.
- No new compatibility shim exists.
- Superseded RagOpt/product orchestration has no callers and is deleted.
- Products retain their policies and application services.

### P8 completion outcome

P8 classified and guarded the final repository boundaries. Optkit rejects imports of RagKit, RagOpt, Judgekit, Coinvault, and RAG-TTC; RagKit rejects product and experiment-orchestration imports; Judgekit rejects Optkit and RAG-TTC core imports. Product integration remains in RAG-TTC.

The now-unused outer customer `ToolRegistryFactory` and `buildProviderToolRegistry` path was deleted after P5 made `customerapp.EngineAdapter` authoritative. No compatibility path was added.

The frozen I5 RagOpt candidate/evaluation command and historical `runstore`/`review` readers are explicitly retained, not classified as superseded. P6/P7 prove durable retrieval and historical judge measurement, but not yet the old command's candidate asset locks, full answer execution, gate/report behavior, and historical run-directory parity. Deleting it would violate the accepted parity-first decision. The exact migration gates and all active callers are recorded in `sources/09-p8-boundary-and-migration-inventory.md`.

Validation across Optkit, RagKit, Judgekit, RagOpt, RAG-TTC, and Coinvault ends in `P8_VALIDATION=PASS`.

## 15. Final validation and delivery

### Tasks

1. Print final-validation phase plan slip.
2. Run format, generation, vet, lint, unit, race, and static-link checks as supported by each touched repository.
3. Re-run all deterministic semantic fixtures.
4. Run an end-to-end restart campaign and archive command output under ticket `sources/`.
5. Verify all task IDs are complete.
6. Update architecture/current-state docs to match implementation.
7. Complete the diary with commit hashes, failures, sharp edges, review instructions, and follow-ups.
8. Relate modified and decision-shaping files with absolute paths.
9. Run `docmgr doctor --ticket OPTKIT-002 --stale-after 30` until clean.
10. Print final project status slip.
11. Upload the index, plan, diary, verification report, and key API guide to reMarkable after a dry-run.
12. Close the ticket only when no required work remains.

### Acceptance gate

All repository checks pass, the worktrees contain no unexplained changes, ticket doctor is clean, all phase slips were printed, the bundle upload succeeds, and the ticket closes with no open tasks.

## 16. Key decisions

### Decision: RagKit remains separate

- **Context:** RagKit contains cohesive RAG-domain algorithms; Optkit aims to optimize systems beyond RAG.
- **Decision:** Keep RagKit as a separate reusable library and integrate it through product-owned system executables.
- **Rationale:** Repository collapse is not architectural unification. A general orchestration framework should not own every domain algorithm it can evaluate.
- **Consequence:** Optkit and RagKit evolve independently; RAG-TTC composes both.
- **Status:** accepted

### Decision: Absorb RagOpt orchestration, not RagKit semantics

- **Context:** RagOpt overlaps Optkit in candidate, paired trial, run custody, gating, and report orchestration.
- **Decision:** Replace exercised RagOpt orchestration with Optkit after parity is proven.
- **Rationale:** This removes duplicated control-plane semantics while preserving the RAG data plane.
- **Status:** accepted

### Decision: Judgekit remains a provider behind an instrument

- **Context:** Judgekit is a reusable measurement domain and has an independent trust model.
- **Decision:** Integrate through an Optkit instrument first; do not move packages merely for repository count.
- **Rationale:** The instrument boundary is explicit and supports independent evolution.
- **Status:** accepted

### Decision: UI is a read-only scientific query plane

- **Context:** Campaigns and parameter changes will primarily be authored through CLI workflows operated by LLM agents; the UI is needed for comprehension rather than form-based administration.
- **Decision:** Build navigation, search, replay, provenance, and visualization APIs only. Do not add browser mutation endpoints.
- **Rationale:** This preserves one command path, avoids duplicating validation logic in forms, and lets agents hand scientists deep links to durable evidence.
- **Consequence:** The backend needs rich projections, stable URLs, SSE replay, and bounded artifact inspection; CLI ergonomics must return campaign URLs/IDs.
- **Status:** accepted

### Decision: Fixed arms before adaptive search

- **Context:** Durable orchestration and measurement must be proven before optimizer complexity.
- **Decision:** P6 uses fixed complete-block arms.
- **Rationale:** Failures remain diagnosable and the control plane is validated in isolation.
- **Status:** accepted

## 17. Cross-phase invariants

- No hidden provider/network dependency in deterministic CI fixtures.
- No credentials or live handles in snapshots/artifacts.
- No model-controlled authorization or corpus selection.
- No silent failure-row filtering.
- No cross-epoch aggregation.
- No transport-only evaluation path.
- No duplicated retrieval implementation after P2.
- No automatic deployment or production promotion.
- No browser mutation API; UI state is derived from authoritative journal and artifact facts.
- No backwards-compatibility shim without explicit approval.
- Every phase has before/after thermal slips and diary evidence.

## 18. References

- `OPTKIT-001/design-doc/01-optkit-current-state-and-pragmatic-rag-implementation-guide.md`
- `/home/manuel/workspaces/2026-08-24/use-optkit/sources/optkit-implementation-source.zip`
- `/home/manuel/workspaces/2026-08-24/use-optkit/sources/implementation-diary.md`
- `/home/manuel/workspaces/2026-08-24/use-optkit/coinvault/internal/knowledge/service.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/search/search.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/ragkit/rag/retrieval/retrieval.go`
- `/home/manuel/workspaces/2026-08-24/use-optkit/judgekit/judging/claimjudge.go`
