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

**User prompt (verbatim):** "ok, let's create a new ticket for the CLI, and implementthat first, that way we can already validate that what the UI is going to display is valid. Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.\n\n\nBtw, don't forget to commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)"

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
