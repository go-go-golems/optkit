---
Title: Implementation Diary
Ticket: OPTKIT-018
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
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/design-doc/01-intern-guide-to-candidate-read-projections-and-the-workbench-command-api.md
      Note: Primary design deliverable whose research and delivery this diary records
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/various/remarkable-dry-run.log
      Note: Ticket bundle selection and destination dry-run evidence
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/various/remarkable-upload.log
      Note: Successful rendered PDF upload receipt
ExternalSources: []
Summary: Chronological research, design, implementation, validation, and delivery record for OPTKIT-018.
LastUpdated: 2026-08-26T14:20:28.245641732-04:00
WhatFor: ""
WhenToUse: ""
---


# Diary

## Goal

Design historical candidate projections and a separate authorized workbench command API. This diary records research, design decisions, validation, commits, delivery evidence, failures, and continuation instructions.

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

**Assistant interpretation:** Expose sealed intent/catalog facts for reads and adapt compile/preview/seal applications to HTTP without mixing them into specialist projections.

**Inferred user intent:** Give React authoritative read and command contracts while preserving historical determinism.

### What I did

- Mapped current GET-only ServeMux, fact loader, comparison projection, and application services.
- Specified candidate/catalog responses, command routes, authorization actions, typed errors, idempotency, and preview envelopes.
- Designed historical fallback, sealed catalog labeling, strict HTTP, security, concurrency, and live smoke tests.
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

- One server may host both route groups, but package/application boundaries must keep historical reads incapable of compiling or sealing proposals.
- Ticket boundaries are most reliable when each names both its upstream contract and the next consumer.
- Documentation should distinguish observed runtime behavior from proposed APIs and future implementation choices.

### What was tricky to build

- One server may host both route groups, but package/application boundaries must keep historical reads incapable of compiling or sealing proposals.
- The guide had to be self-contained without copying the entire parent architecture. It summarizes prerequisites, then links each claim to the file that owns it.

### What warrants a second pair of eyes

- Review actor spoofing, sensitive diagnostics/logging, typed error mapping, restart-safe idempotency, and sealed-versus-current catalog behavior.
- Review all proposed public schema/API names before implementation makes them expensive to change.

### What should be done in the future

- Implement candidate fact projections first, then catalog reads, command applications, strict HTTP adapters, and frontend fixtures.
- Update this diary immediately when implementation reveals a false assumption or accepted contract change.

### Code review instructions

Start with the ticket's `design-doc/01-*.md`, then inspect these decision-shaping files:

- `rag-ttc/pkg/ttc/specialistapi/projector.go`
- `rag-ttc/pkg/ttc/specialistapi/types.go`
- `rag-ttc/pkg/ttc/specialistapi/http.go`
- `rag-ttc/pkg/ttc/experimentworkbench/service.go`
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/serve.go`

Validate the design workspace with:

```bash
docmgr validate frontmatter --doc optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/reference/01-implementation-diary.md
docmgr doctor --ticket OPTKIT-018 --stale-after 30
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

**Commit (documentation):** `ec784b0d84e47ad45a7429f98482e9678df5dac7` — "OPTKIT-018-020: design workbench delivery layers"

### What I did

- Ran `docmgr doctor --ticket OPTKIT-018 --stale-after 30` and obtained `All checks passed`.
- Scanned the index, guide, and diary for generated placeholder sections; none remained.
- Ran focused baseline Go tests before documentation changes; all selected Optkit and RAG-TTC packages passed.
- Ran `git diff --check`/`git diff --cached --check`, corrected whitespace findings, and committed the ticket package.
- Ran the program upload script in `--dry-run` mode, then rendered and uploaded `OPTKIT-018 Workbench API Guide.pdf`.
- Preserved upload evidence in `various/remarkable-dry-run.log` and `various/remarkable-upload.log`.

### Why

- A detailed guide is only useful when frontmatter, links, task state, and delivery artifacts agree.
- Grouping commits by dependency layer keeps review focused while avoiding one nine-ticket mega-commit.
- The reMarkable bundle lets the architecture be reviewed away from the source tree without losing the index or diary context.

### What worked

- `remarquee` reported: `OK: uploaded OPTKIT-018 Workbench API Guide.pdf -> /ai/2026/08/26/OPTKIT-018`.
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
- Inspect commit `ec784b0d84e47ad45a7429f98482e9678df5dac7` for the documentation batch.
- Validate locally with `docmgr doctor --ticket OPTKIT-018 --stale-after 30`.
- Consult `tasks.md` for the still-open implementation sequence.

### Technical details

```text
bundle: OPTKIT-018 Workbench API Guide.pdf
remote: /ai/2026/08/26/OPTKIT-018
commit: ec784b0d84e47ad45a7429f98482e9678df5dac7
doctor: clean
production code changes: none
```
