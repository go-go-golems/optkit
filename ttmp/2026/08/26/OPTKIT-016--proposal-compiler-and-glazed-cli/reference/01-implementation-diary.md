---
Title: Implementation Diary
Ticket: OPTKIT-016
Status: active
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
