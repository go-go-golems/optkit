---
Title: Investigation Diary
Ticket: OPTKIT-004
Status: complete
Topics:
    - architecture
    - implementation
    - judgekit
    - optkit
    - rag
    - rag-ttc
    - scientific-workflow
    - ui
    - visualization
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://ttmp/2026/08/25/OPTKIT-004--end-to-end-rag-optimization-workbench-and-specialist-ui/design-doc/01-intern-guide-to-the-rag-optimization-workbench-and-specialist-ui.md
      Note: Primary architecture and implementation deliverable
    - Path: repo://ttmp/2026/08/25/OPTKIT-004--end-to-end-rag-optimization-workbench-and-specialist-ui/scripts/01-validate-guide.py
      Note: Structural and source-reference validation
    - Path: ws://rag-ttc/pkg/ttc/retrievaleval/evaluate.go
      Note: Stage-aware metric and first-loss evidence used in the design
ExternalSources: []
Summary: Chronological evidence, decisions, validation, and delivery record for the RAG optimization workbench intern guide.
LastUpdated: 2026-08-25T19:03:00Z
WhatFor: Preserve how the architecture and UI guide was researched, written, validated, and delivered.
WhenToUse: Read before revising the OPTKIT-004 design or implementing its phases.
---


# Diary

## Goal

Record the evidence-first creation of an intern-ready architecture, workflow, UI design, and phased implementation guide for an end-to-end TTC RAG optimization workbench.

## Step 1: Create the ticket and map the existing architecture

This step created a dedicated docmgr workspace before drawing conclusions. The investigation then mapped concrete ownership and API boundaries across Optkit, RAG-TTC, RagKit, and Judgekit so the proposed workbench would extend working contracts rather than invent a parallel system.

The evidence showed that the workspace already has a durable experiment control plane, a direct policy-safe retrieval service, stage-aware deterministic evaluation, a canonical answer service, reusable RAG algorithms, attributed judge protocols, and a read-only campaign explorer. The new design therefore focuses on composing and extending these foundations.

### Prompt Context

**User prompt (verbatim):** "Ok, createa  new docmgr ticket and write up the work and workflow + UI Design you just did as a document for the intern. Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Create a new ticket, research the current multi-repository implementation, write one exhaustive intern guide covering architecture, optimization workflow, UI, APIs, pseudocode, diagrams, implementation, and evidence, then validate and upload it to reMarkable.

**Inferred user intent:** Turn the conceptual “killer RAG optimizer” vision into a durable implementation reference that a new engineer can follow without rediscovering repository boundaries or scientific constraints.

### What I did

- Ran `docmgr status --summary-only`, `docmgr ticket list`, and `docmgr vocab list` from the Optkit repository.
- Created `OPTKIT-004 — End-to-End RAG Optimization Workbench and Specialist UI` with architecture, implementation, Judgekit, Optkit, RAG, scientific-workflow, UI, and visualization topics.
- Added the primary design doc and this investigation diary.
- Added five tasks covering architecture evidence, optimization workflow, UI/projectors, technical implementation guidance, validation, and reMarkable delivery.
- Enumerated relevant files and symbols across all four repositories.
- Read the current system registry, snapshots, complete-block trials, query service, HTTP server, browser application, retrieval service, runtime identity, evaluator, customer application service, durable campaign adapter, RagKit contracts, and Judgekit contracts.
- Captured line-level anchors for current behavior and implementation gaps.

### Why

- A new intern needs a trustworthy map of what is already implemented versus what is proposed.
- Repository ownership is a core architecture constraint: RagKit, Optkit, RAG-TTC, and Judgekit solve different reusable problems.
- File-backed evidence prevents the design from drifting into generic RAG advice disconnected from the codebase.

### What worked

- Existing docmgr vocabulary already covered every needed topic; no vocabulary changes were required.
- The repositories expose clean boundaries that support the proposed architecture.
- Current code provides strong concrete examples for identity, policy ordering, trajectory sealing, complete-block experiments, restart reconciliation, bounded query APIs, and judge attribution.

### What didn't work

- N/A. Ticket creation and repository discovery completed without command failures.

### What I learned

- The current system is further along than a typical design baseline: Phase P6 already proves a restartable real retrieval campaign with stage artifacts and paired estimates.
- The highest-value next architecture is not another runner. It is a layered configuration graph, artifact reuse planner, answer/judge instruments, RAG-specific projectors, and specialist diagnostics.
- The current explorer's generic event model is useful but will not scale to rich RAG analysis without server-side projectors.

### What was tricky to build

- The design must separate reusable domain mechanics from product policy while still giving the intern one coherent mental model. The solution is to describe the execution graph end-to-end, then mark repository ownership at each boundary.
- “Summarization” appears in two distinct roles: index-time searchable representations and post-retrieval context compression. Conflating them would break lineage and invalidation rules, so the guide defines them separately.
- Judge scores cannot be treated as ordinary final metrics without explaining calibration, reliability, cache mode, and historical remeasurement.

### What warrants a second pair of eyes

- Confirm that proposed generic materialization-DAG concepts should remain in RAG-TTC until a second product domain proves they belong in Optkit.
- Confirm the first promotion policy's primary metrics and hard constraints.
- Confirm that the read-only UI boundary remains desirable if a later approval workflow is introduced.

### What should be done in the future

- Keep file references current when APIs move.
- Update this design as P7 and later implementation phases settle the answer/Judgekit integration and query projector contracts.

### Code review instructions

- Start with the file-reference map in the primary design document.
- Compare every “current” claim against the cited symbol and line range.
- Verify proposed packages do not move RAG concepts into Optkit or product policy into RagKit.

### Technical details

- Ticket root: `optkit/ttmp/2026/08/25/OPTKIT-004--end-to-end-rag-optimization-workbench-and-specialist-ui`.
- Primary document: `design-doc/01-intern-guide-to-the-rag-optimization-workbench-and-specialist-ui.md`.
- Evidence repositories: `optkit`, `rag-ttc`, `ragkit`, and `judgekit`.

## Step 2: Write the intern architecture, workflow, and UI guide

This step converted the architecture evidence and the earlier product vision into a 32-section implementation guide. The document begins with vocabulary and current-state orientation, then specifies the full optimization graph, layered configuration identity, invalidation and reuse, stage-level metrics, scientific experiment design, durable execution, specialist UI, query APIs, implementation phases, test strategy, decisions, risks, and onboarding tasks.

The guide deliberately uses prose before details, compact text diagrams, typed API sketches, pseudocode, tables, concrete commands, and file references. It is intended to be read as both an onboarding chapter and an implementation plan.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Produce the complete intern-facing technical deliverable after mapping the implementation evidence.

**Inferred user intent:** Give a new engineer enough context and specificity to implement the workbench safely and incrementally.

### What I did

- Wrote a 2,515-line, approximately 10,500-word design and implementation guide.
- Defined the end-to-end pipeline from source revision through promotion decision.
- Explained all core entities: corpus, chunk, representation, vector, prepared route, snapshot, arm, episode, trajectory, measurement epoch, observation, estimate, and candidate artifact.
- Mapped repository responsibilities and cited current APIs.
- Added a layered configuration model and example YAML.
- Added a stage invalidation table and reuse-planner API/pseudocode.
- Specified optimization variables and measurements for normalization, chunking, representations, embeddings, lexical/vector retrieval, query transforms, policy, fusion, augmentation, reranking, admission, context construction, answering, and judging.
- Defined a progressive experimental funnel, dataset roles, pairing, repeats, multi-objective constraints, candidate proposal workflow, and promotion manifest.
- Designed ten specialist screens, including the cockpit, comparison matrix, pipeline microscope, rank waterfall, chunk laboratory, answer/evidence studio, judge calibration lab, Pareto explorer, and provenance inspector.
- Proposed bounded GET-only query routes and projection contracts.
- Added frontend architecture, accessibility, security, scaling, failure semantics, testing, package layout, nine implementation phases, nine decision records, alternatives, risks, open questions, and intern onboarding.

### Why

- A complete workbench requires coordinated changes across data preparation, runtime retrieval, answering, measurement, persistence, query projection, and UI.
- Explicit invalidation rules are necessary to avoid rerunning expensive upstream stages when only downstream configuration changes.
- The specialist UI needs domain projectors and diagnostic concepts before implementation begins; otherwise it will become a raw event dashboard.

### What worked

- Current code provided concrete examples for nearly every proposed boundary.
- Text diagrams render well in Markdown and remain suitable for PDF/reMarkable output.
- Layering the document from vocabulary to current state to proposed architecture makes it usable by an unfamiliar engineer.
- Decision records make the read-only UI, repository boundaries, policy ordering, layered identity, and deterministic-versus-judge measurement choices explicit.

### What didn't work

- N/A. The initial full draft was written successfully.

### What I learned

- The central technical advantage of the proposed workbench is not a specific optimizer algorithm. It is the semantic artifact DAG that lets each experiment reuse exactly the valid upstream evidence.
- The pipeline microscope and rank waterfall are the most important UI features because they turn aggregate regressions into actionable stage diagnoses.
- Judge calibration and reliability must be first-class UI concepts or optimization will over-trust a potentially gameable score.

### What was tricky to build

- The guide had to remain exhaustive without becoming a flat feature list. The organizing spine is the immutable pipeline graph and its stage identities.
- UI design had to remain ambitious while honoring the existing no-build, semantic HTML/CSS/JavaScript, read-only constraint.
- Statistical guidance had to be practical for an intern: complete blocks, paired estimates, missingness, progressive gates, Pareto fronts, and holdouts are specified before more advanced Bayesian or evolutionary methods.

### What warrants a second pair of eyes

- Review the proposed route and projector naming before schemas are frozen.
- Review whether confidence intervals should use bootstrap, t-based, or another first implementation.
- Review artifact retention and projection-index scale assumptions against representative production campaign sizes.
- Review how context token counts are normalized across providers.

### What should be done in the future

- Implement Phase 1, Judgekit over sealed historical answers, before expanding the optimization search space.
- Keep the guide synchronized with accepted schema and package decisions.

### Code review instructions

- Read Sections 1–7 for correctness of current-state claims.
- Review Sections 8–15 for architecture and scientific workflow.
- Review Sections 16–18 for UI/query/frontend feasibility.
- Review Sections 23–26 for tests, phases, and decisions.
- Use Section 30 as the source-code review map.

### Technical details

- Draft size: 2,515 lines; approximately 10,508 words; approximately 84 KiB.
- Major proposed API families: layered configuration graph, materialization reuse planner, pipeline projection, case comparison, estimate set, Pareto point, agent proposal, and promotion manifest.
- UI transport remains standard-library `http.ServeMux` with GET-only routes and embedded static files.

## Step 3: Validate, publish, and close the documentation package

This step turned the draft into a checked deliverable. A ticket-local validation script verified document size, section ordering, code fences, decision records, required concepts, and every line-anchored repository file reference. Docmgr frontmatter and doctor checks passed before delivery.

The guide and diary were then rendered as one PDF bundle with a table of contents and uploaded to the ticket-specific reMarkable folder. The upload command returned an explicit success path, so no redundant cloud listing was needed.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Complete validation, ticket bookkeeping, and reMarkable delivery for the intern guide.

**Inferred user intent:** Receive a finished, portable technical package rather than an unvalidated local draft.

### What I did

- Added `scripts/01-validate-guide.py` under the ticket workspace.
- Validated 2,515 lines, 10,508 words, nine decision records, balanced code fences, numbered Sections 1–32, required concepts, and 65 line-anchored Go file references.
- Ran `docmgr validate frontmatter --suggest-fixes` on the guide and diary; both returned `Frontmatter OK`.
- Related the primary architecture files and validation evidence with absolute file notes.
- Checked the first four ticket tasks and updated the changelog.
- Ran `docmgr doctor --ticket OPTKIT-004 --stale-after 30`; all checks passed.
- Ran a dry-run bundle render for the guide and diary.
- Uploaded `OPTKIT-004 RAG Optimization Workbench Intern Guide.pdf` to `/ai/2026/08/25/OPTKIT-004`.

### Why

- Structural checks catch truncated drafts, broken fences, missing sections, and stale source references before PDF delivery.
- Bundling the guide and diary gives the intern both the final design and the reasoning/evidence trail.
- A ticket-specific remote path keeps the document discoverable and avoids name collisions.

### What worked

- The structural script ended with `GUIDE_VALIDATION=PASS`.
- Both frontmatter validations passed.
- Docmgr doctor passed cleanly.
- The dry run resolved both Markdown inputs and the expected remote directory.
- The real upload returned:

  `OK: uploaded OPTKIT-004 RAG Optimization Workbench Intern Guide.pdf -> /ai/2026/08/25/OPTKIT-004`

### What didn't work

- N/A. Validation, rendering, authentication, and upload completed without failures.

### What I learned

- Automated source-path validation is useful for long evidence-backed documents because a single repository rename could otherwise leave dozens of plausible but invalid references.
- The design's plain-text diagrams and non-nested fenced blocks rendered without requiring a special diagram toolchain.

### What was tricky to build

- The validation script needed to resolve ticket location back to the multi-repository workspace root before checking `optkit/`, `rag-ttc/`, `ragkit/`, and `judgekit/` references. It derives the workspace from the script's ticket-local path rather than assuming the caller's current directory.
- Delivery had to follow both the ticket workflow's dry-run requirement and the reMarkable skill's minimal-call guidance. Two calls—dry run and actual upload—were sufficient.

### What warrants a second pair of eyes

- Review PDF readability of the widest configuration and route tables on the physical device.
- Review source line anchors after major refactors; the files are validated but line ranges naturally drift when code changes.

### What should be done in the future

- Re-upload with `--force` only when an intentional revised edition is ready, because overwriting may remove annotations.
- Add an edition/date suffix instead if annotated historical versions should remain available.

### Code review instructions

- Run `python3 scripts/01-validate-guide.py` from any directory.
- Run `docmgr doctor --ticket OPTKIT-004 --stale-after 30` from the Optkit repository.
- Open the reMarkable bundle and inspect the table of contents, text diagrams, tables, and code blocks.

### Technical details

- Structural result: `GUIDE_VALIDATION=PASS`.
- Docmgr result: all checks passed.
- Bundle name: `OPTKIT-004 RAG Optimization Workbench Intern Guide`.
- Remote directory: `/ai/2026/08/25/OPTKIT-004`.
- Uploaded file: `OPTKIT-004 RAG Optimization Workbench Intern Guide.pdf`.
