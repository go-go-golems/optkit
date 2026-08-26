---
Title: React Workbench Framework and RRF Vertical Slice
Ticket: OPTKIT-019
Status: superseded
Topics:
    - architecture
    - design
    - implementation
    - optkit
    - rag-ttc
    - ui
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Build a reusable catalog-driven React workbench and prove it through a complete RRF candidate experiment.
LastUpdated: 2026-08-26T14:20:29.258849198-04:00
WhatFor: Provide generic authoring controls plus specialized expert experiences without duplicating backend semantics.
WhenToUse: Implement only after OPTKIT-018 exposes stable read and command contracts.
---

# React Workbench Framework and RRF Vertical Slice

## Overview

This ticket creates the frontend framework after backend semantics are proven. Generic editors cover ordinary catalog values; a plugin registry supplies expert inspectors and previews; a reusable shell carries case, mutation, invalidation, intent, evidence, and sealing state. The RRF vertical slice proves the complete experience.

**Program position:** first product-complete scalar workflow; prerequisite for the asset proof.

- [Intern guide](design-doc/01-intern-guide-to-the-react-workbench-framework-and-rrf-vertical-slice.md)
- [Implementation diary](reference/01-implementation-diary.md)

## Superseded — 2026-08-26

This ticket was closed **before implementation** and replaced by the PBUI adoption
track. The program decided to build the frontend on the presentation-based UI
stack (`@hyperslop-systems/pbui`, `pbui-workbench`, `plot`) instead of the bespoke
`WorkbenchShell`/`WorkbenchProvider` architecture this guide specifies. The
successors are:

- **OPTKIT-021** — PBUI adoption ADRs and workbench vocabulary (packaging, Go
  workbench service placement, document formats, presentation types/verbs).
- **OPTKIT-022** — RAG-TTC workbench product package and evidence tiles (the
  pbui product scaffold plus the read-only tiles over `specialistapi`).
- **OPTKIT-023** — Propose workspace: draft document and RRF vertical slice
  (inherits this ticket's exit criterion verbatim).
- **OPTKIT-024** — Agent seat over the workbench vocabulary (deferred).

The intern guide in this ticket remains valuable as a design record and as a
source of salvaged content. The following parts survive into the successors and
must not be re-invented:

- generic editors keyed by `ValueSpec` kind with a visible unsupported fallback
  (never silent coercion) — lifted into OPTKIT-023;
- plugin resolution order: specialized variable editor → generic by value kind →
  visible fallback — becomes the pbui product's plugin registry rule;
- the debounced-compile / stale-response / `sealState` machine
  (`editing → compiling → sealable → sealing → sealed`) — lifted into the
  `ragttc.proposal-draft/v1` document flow in OPTKIT-023;
- the `InvalidationStrip` rule that recomputation status renders in **server
  order** from `InvalidationPlan.Steps`, never inferred client-side;
- the RRF contribution inspector (`contribution = weight / (k + rank)` with
  auditable per-channel operands) and the parity-pinning discipline against
  recorded six-decimal Go results (OPTKIT-010 precedent);
- accessibility, deep-link, loading/empty/error state requirements.

What does **not** survive: the bespoke shell (`WorkbenchShell.tsx`,
`WorkbenchProvider.tsx`), the proposed `apps/specialist/web/src/workbench/`
module tree, and the `/workbench/campaigns/:campaign/cases/:case/proposal`
route pattern — pbui-workbench's split-tree surface, workbench documents, and
accept protocol replace them.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **superseded** (by OPTKIT-021, OPTKIT-022, OPTKIT-023, OPTKIT-024)

## Topics

- architecture
- design
- implementation
- optkit
- rag-ttc
- ui

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
