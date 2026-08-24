---
Title: Investigation Diary
Ticket: OPTKIT-001
Status: active
Topics:
    - optkit
    - rag
    - rag-ttc
    - coinvault
    - judgekit
    - architecture
    - migration
    - local-development
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://judgekit/ttmp/2026/08/17/JUDGEKIT-001--design-and-implement-judgekit/design-doc/03-lightweight-research-guarantees-after-merge.md
      Note: Accepted trusted-research guarantee model
    - Path: ws://judgekit/ttmp/2026/08/17/JUDGEKIT-001--design-and-implement-judgekit/design-doc/04-fully-hardened-judgekit-architecture.md
      Note: Explicitly deferred adversarial reference architecture
    - Path: ws://sources/implementation-diary.md
      Note: Chronological evidence for the supplied implementation
    - Path: ws://sources/optkit-clean-slate-architecture-and-porting-guide.md
      Note: Primary clean-slate architecture under assessment
    - Path: ws://sources/optkit-implementation-source.zip
      Note: Tested implementation source artifact
ExternalSources: []
Summary: Chronological, evidence-backed record of the Optkit assessment and pragmatic Coinvault-to-RAG-TTC port design.
LastUpdated: 2026-08-24T22:20:00-04:00
WhatFor: Preserve investigation evidence, commands, failures, decisions, and review guidance so another engineer can continue the unification work.
WhenToUse: Read before implementing OPTKIT-001 or revisiting its architecture and scope decisions.
---


# Diary

## Goal

This diary records the assessment of the clean-slate Optkit design and implementation archive, the current repositories, Judgekit's lightweight local-research policy, and the practical path for applying Coinvault's RAG work to RAG-TTC. It captures the evidence and tradeoffs behind the intern-oriented implementation guide in this ticket.

## Step 1: Establish the ticket and verify the supplied implementation

The first step created a repository-local ticket in `optkit/ttmp` and turned the broad unification request into an explicit investigation checklist. This gives the assessment a durable home before conclusions are written and separates evidence gathering from implementation recommendations.

I also unpacked the supplied implementation archive into a temporary review directory and ran its complete test suite. This immediately exposed the most important current-state fact: the archive contains a coherent, green Optkit vertical slice, while the checked-out `optkit/` repository still contains only the unnormalized Go template.

### Prompt Context

**User prompt (verbatim):**

> Create a new docmgr in ./optkit to analyze all this, and assess where we are at. Also read the last tickets in judgekit/ and how we have decided to relax the hardening of it all to support local development, not adversarial work. 
>
> Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.
>
>
> You can run any experiments and scripts and compile what you want and all that, store scripts and code you write in the ticket's scripts/ folder.
>
> commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Create and maintain OPTKIT-001 in the `optkit` repository; assess the supplied architecture, implementation diary, archive, and six related repositories; incorporate Judgekit's trusted-local-development policy; produce an intern-ready technical design and implementation plan; commit the documentation in logical stages; and upload the final bundle to reMarkable.

**Inferred user intent:** Convert extensive prior design work into a truthful status assessment and a practical next implementation sequence that a new engineer can execute without re-deriving the architecture or overbuilding adversarial integrity machinery.

### What I did

- Read the repository and workspace `AGENT.md` guidance plus the docmgr, diary, ticket-research, and reMarkable skill references.
- Confirmed that the workspace root is an aggregate workspace rather than a Git repository and that each component repository has its own `.git` directory.
- Added ticket vocabulary for `optkit`, `rag`, `rag-ttc`, `coinvault`, `judgekit`, `architecture`, `migration`, and `local-development`.
- Created `OPTKIT-001`, its primary design document, this diary, and eight evidence-to-delivery tasks.
- Inventoried `sources/`: the architecture guide is 8,841 lines, the implementation diary is 741 lines, and the implementation ZIP contains a complete `github.com/go-go-golems/optkit` module.
- Unpacked `sources/optkit-implementation-source.zip` under `/tmp/optkit-implementation-source` for read-only review.
- Ran `go test ./... -count=1` in the unpacked archive.
- Located Judgekit's latest repository-local ticket, `JUDGEKIT-001`, and read its lightweight and fully hardened design documents.

### Why

- A ticket and task list make the broad assessment reviewable and continuation-friendly.
- Testing the supplied archive distinguishes “designed” from “implemented and mechanically verified.”
- Comparing the archive with the checked-out repository prevents the guide from incorrectly claiming that the implementation has already landed.
- Reading the latest Judgekit decision documents establishes the intended trust model before recommending Optkit persistence, identity, and evaluation boundaries.

### What worked

- `docmgr ticket create-ticket` created `optkit/ttmp/2026/08/24/OPTKIT-001--pragmatic-optkit-unification-and-coinvault-to-rag-ttc-port` with the expected index, tasks, changelog, design, reference, and scripts directories.
- Every package in the supplied implementation archive passed its tests, including artifact storage, budgets, campaign reduction, CLI behavior, episode writing, experiments, architecture rules, SQLite storage, records, spaces, and the Numbergame vertical slice.
- The latest Judgekit documents state a clear three-level guarantee model: structural validity, research attribution, and deferred adversarial integrity.

### What didn't work

- Running `git status` at the aggregate workspace root failed exactly as follows:

  `fatal: not a git repository (or any of the parent directories): .git`

  The resolution was to run Git commands with `git -C optkit`, `git -C judgekit`, and equivalent repository-local paths.

- The first repository statistics command hid `go list` stderr and reported zero packages for every repository. This was an instrumentation mistake, not evidence that the repositories have no Go packages. Subsequent validation commands must retain stderr and run from each module with the appropriate workspace setting.

### What I learned

- The checked-out `optkit` repository is still the initial template at commit `dc48856`; its module path is `github.com/go-go-golems/XXX`, and it has no implementation packages from the supplied archive.
- The supplied archive is substantially beyond a sketch: it implements a local modular monolith with typed search spaces, immutable records, artifact stores, episodes, experiments, campaign journaling/reduction, durable SQLite queues, budget reservations, projections, a CLI, and an executable Numbergame example.
- Judgekit's current policy is not “remove correctness checks.” It is “enforce structural validity and research attribution at meaningful boundaries, but do not add hostile-storage, signatures, immutable typestate, or cross-trust custody without a concrete threat and consumer.”

### What was tricky to build

- The workspace has a parent `.ttmp.yaml`, while the user requested a ticket under the nested `optkit` repository. The important symptom was that `docmgr status` reported the effective root as `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp` even though the config path was the workspace-level file. I verified the effective root before creating the ticket rather than assuming config-file location and ticket-root location were identical.
- The supplied implementation exists only as a ZIP and cannot be cited with repository commit history. I treated it as a tested source artifact and kept it out of the checked-out repository during assessment so the guide can distinguish “available to import” from “landed on this branch.”

### What warrants a second pair of eyes

- Confirm that `OPTKIT-001` is the desired ticket identifier before future Optkit tickets are created.
- Review the eventual recommendation about importing the archive as a baseline rather than reimplementing it directly in the current template.
- Verify that the pragmatic trust model is applied consistently: do not remove measurement-attribution checks merely because adversarial hardening is deferred.

### What should be done in the future

- Preserve the archive's green baseline when importing it into the real `optkit` repository.
- Keep source-artifact provenance visible in the first import commit message and ticket changelog.
- Replace the failed package-count probe with an auditable inventory script stored in this ticket if package counts are used in the final report.

### Code review instructions

- Start with this diary and `tasks.md` to verify ticket scope.
- Compare `optkit/go.mod` with `/tmp/optkit-implementation-source/optkit/go.mod` to confirm the checked-out-template versus supplied-implementation distinction.
- Re-run `cd /tmp/optkit-implementation-source/optkit && go test ./... -count=1` to validate the supplied implementation.
- Read Judgekit's `design-doc/03-lightweight-research-guarantees-after-merge.md` before reviewing any integrity or identity recommendation.

### Technical details

- Ticket root: `/home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-001--pragmatic-optkit-unification-and-coinvault-to-rag-ttc-port`
- Supplied architecture: `/home/manuel/workspaces/2026-08-24/use-optkit/sources/optkit-clean-slate-architecture-and-porting-guide.md`
- Supplied diary: `/home/manuel/workspaces/2026-08-24/use-optkit/sources/implementation-diary.md`
- Supplied source archive: `/home/manuel/workspaces/2026-08-24/use-optkit/sources/optkit-implementation-source.zip`
- Archive revision comment: `1786d1da86c9e03316ed71336bbc993fa30531f0`
- Archive module: `github.com/go-go-golems/optkit`
- Verified command: `cd /tmp/optkit-implementation-source/optkit && go test ./... -count=1`

## Step 2: Reframe unification around a pragmatic RAG-TTC vertical slice

This step converted the broad clean-slate program into a product-first implementation sequence. Instead of recommending that all four extracted kits move before useful product work can begin, the analysis now treats Coinvault as the behavioral reference and RAG-TTC as the first real Optkit consumer.

The repository audit showed that RAG-TTC already has a verified hybrid index, strict file-backed configuration, turn-scoped evidence, and a runtime composition below HTTP. The proposed work is therefore a convergence and observability refactor: split canonical retrieval semantics from the model-facing tool, add Coinvault-style runtime identity and stage diagnosis, create a deterministic retrieval suite, and then register that direct service as an Optkit system.

### Prompt Context

**User prompt (verbatim):**

> we want to be pragmatic and basically apply the work we did in coinvault to rag-ttc to get a proper RAG system going
>
> [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Prioritize a practical transfer of Coinvault's proven RAG construction, retrieval, evidence, attribution, and evaluation discipline into RAG-TTC; use that product slice to drive Optkit rather than making full framework consolidation a prerequisite.

**Inferred user intent:** Deliver a trustworthy, diagnosable RAG system in RAG-TTC soon while preserving a credible route toward the unified optimization framework.

**Commit (docs and investigation script):** `f5d7366a4eb1d1e0c370e022952d0bfe88fefe69` — "docs(optkit-001): add pragmatic RAG implementation guide"

### What I did

- Read the clean-slate architecture sections covering package consolidation, Coinvault, RAG-TTC, pull-request ordering, worked RAG examples, and final recommendation.
- Read all seven supplied implementation-diary steps and inspected the archive's core APIs.
- Read Coinvault's verified-bundle open path, hybrid service, runtime configuration, evidence ledger, tool output, strict eval schema, candidate-pool diagnosis, and canonical evaluator composition.
- Read RAG-TTC's current hybrid search tool, customer bundle lifecycle, runtime composer/resolver, strict tool configuration, index build, knowledge build, answer-quality runner, and relevant tests.
- Read RagKit's retrieval and broad answering APIs, RagOpt's package responsibilities, and current Judgekit value types/execution boundary.
- Ran focused and full-kit tests with `GOWORK=off` where appropriate.
- Added `scripts/01-repository-inventory.sh` and generated `sources/01-repository-inventory.md`.
- Wrote the primary design document with current-state evidence, diagrams, API sketches, retrieval/controller pseudocode, implementation phases, test strategy, decisions, risks, deletion gates, and intern review instructions.
- Checked the first seven OPTKIT-001 tasks and related seven decision-shaping files to the design document.

### Why

- The clarification changes sequencing: “proper RAG in RAG-TTC” is the first outcome, not the reward after a multi-repository consolidation.
- Existing RAG-TTC behavior is substantial enough that a rewrite would discard useful work; a service extraction and semantic hardening are more pragmatic.
- Coinvault's strongest reusable contribution is not its product code but its laws: verified immutable input, deterministic stage ordering, pre-external-boundary policy filtering, runtime identity, bounded evidence, and stage-aware evaluation.
- An inventory script makes repository size, module identity, package count, and kit dependencies reproducible instead of hand-maintained prose.

### What worked

- All six inspected modules or focused package sets passed tests.
- The inventory found that neither Coinvault nor RAG-TTC currently imports Judgekit, despite both having judge/evaluation paths.
- RAG-TTC already has the right seams to extend: `SearchTool.RunRoute`, `ragsearch.Open`, `Handle.NewSessionRegistry`, and `realruntime.Composer.Compose`.
- The audit found one concrete config/runtime gap: intent-routing configuration is strongly typed and validated, while the customer `ragsearch` composition does not compile it into `SearchTool.AddRoute` registrations.
- The 61,607-byte design document fits the requested intern-oriented scope without reproducing the original 8,841-line architecture wholesale.

### What didn't work

- The first run of `scripts/01-repository-inventory.sh` exited with status 1 and no normal output. Under `set -euo pipefail`, `rg -l` correctly returned status 1 when a product had zero Judgekit imports, causing the command substitution to abort even though zero was valid data.
- Running with `bash -x` isolated the failing pipeline at the Judgekit import count. I changed the counter to wrap `rg` with `|| true` before `wc -l`, then reran successfully.
- An early broad `rg` type inventory exceeded the tool output limit and was truncated. I replaced it with focused symbol and file reads rather than treating truncated output as complete evidence.

### What I learned

- RAG-TTC is not missing all RAG infrastructure. It is missing a single canonical semantic service and complete attribution across infrastructure that already exists.
- Coinvault's pre-fusion authorization is product-critical because it mixes public/analyst scopes. RAG-TTC should adopt a server-selected source policy, but it should not invent multi-tenant machinery if its customer corpus is entirely public.
- Judgekit's lightweight architecture is partly prospective: the current `eval.Instance` still stores a digest, and claim extraction still receives a complete instance. The design direction is accepted, but simplification work remains.
- The Optkit archive's hash-chained journal is compatible with a pragmatic local trust model. Keeping an already tested corruption/replay check does not require expanding into signatures or cross-trust attestations.
- The lowest-risk consolidation sequence is product-owned service first, exercised generic primitive move second, last-caller switch third, and deletion fourth.

### What was tricky to build

- The central design tension was avoiding two bad extremes: blocking RAG-TTC on complete kit consolidation, or creating another product-only runner that makes consolidation harder. The resolution is a TTC-owned canonical service whose API remains product-semantic while its temporary implementation composes RagKit primitives.
- “Apply Coinvault” required separating general laws from Coinvault-specific code. Access scopes, SQL tools, product prompts, and Coinvault cases remain local; verified bundles, stage ordering, policy-before-external-boundary, evidence identity, and diagnostic failure custody transfer.
- The status assessment had to distinguish three layers: the checked-out Optkit template, the green implementation archive, and the much larger proposed architecture. Conflating them would either understate completed work or overstate what is on the branch.

### What warrants a second pair of eyes

- Review the accepted decision to allow temporary RagKit imports while building the first Optkit product slice; confirm the named deletion gates are strong enough to prevent permanence.
- Confirm whether RAG-TTC's customer corpus requires actual source-role policy filtering or only attribution and route selection.
- Validate the finding that checked-in intent-routing configuration is not wired into the customer serving composition.
- Review the proposed direct `CustomerApplication.RunTurn` seam against current Pinocchio/sessionstream behavior before implementation.
- Review the initial Optkit import strategy so `ttmp` and repository CI are preserved without obscuring archive provenance.

### What should be done in the future

- Implement Phase 0 as a provenance-preserving import, not a hand port.
- Start the RAG-TTC work with deterministic semantic fixtures and direct-service/tool parity tests.
- Add Judgekit only after deterministic retrieval and answer-contract evidence are reliable.
- Track each temporary old-kit import with a deletion task and architecture test.

### Code review instructions

- Start with the design document's Sections 1, 6, 8, 9, 11, 15, and 20.
- Reproduce the inventory by running `scripts/01-repository-inventory.sh`.
- Validate focused repositories with the commands listed in the design document's evidence section.
- Compare Coinvault `internal/knowledge/service.go:456` with RAG-TTC `pkg/ttc/search/search.go:172` to review the proposed transfer.
- Compare RAG-TTC `pkg/ttc/toolconfig/types.go` with `internal/customer/ragsearch/ragsearch.go` to review the route-wiring finding.

### Technical details

- Inventory script: `ttmp/2026/08/24/OPTKIT-001--pragmatic-optkit-unification-and-coinvault-to-rag-ttc-port/scripts/01-repository-inventory.sh`
- Inventory output: `ttmp/2026/08/24/OPTKIT-001--pragmatic-optkit-unification-and-coinvault-to-rag-ttc-port/sources/01-repository-inventory.md`
- Primary guide: `ttmp/2026/08/24/OPTKIT-001--pragmatic-optkit-unification-and-coinvault-to-rag-ttc-port/design-doc/01-optkit-current-state-and-pragmatic-rag-implementation-guide.md`
- Focused RAG-TTC test command: `GOWORK=off go test ./pkg/ttc/search ./pkg/ttc/toolconfig ./internal/customer/ragsearch ./internal/customer/realruntime ./cmd/rag-ttc/cmds/indexes ./cmd/rag-ttc/cmds/knowledge ./cmd/rag-ttc/cmds/experiments/answerquality -count=1`
- Focused Coinvault test command: `GOWORK=off go test ./internal/knowledge ./internal/knowledgebuild ./internal/webchat/evalchat -count=1`
