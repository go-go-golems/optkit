---
Title: 'Propose Workspace: Draft Document and RRF Vertical Slice'
Ticket: OPTKIT-023
Status: complete
Topics:
    - design
    - implementation
    - ui
    - rag-ttc
    - optkit
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Build the authoring workspace on the PBUI product — proposal draft documents, the catalog/proposal/invalidation/preview/intent tiles, generic ValueSpec editors with a plugin registry, and the fusion.rrf_k vertical slice from case to sealed candidate to comparison.
LastUpdated: 2026-08-26T16:04:32.065457171-04:00
WhatFor: Prove the complete scalar-variable workflow (exit criterion inherited verbatim from the superseded OPTKIT-019) on the PBUI architecture.
WhenToUse: Read when implementing authoring tiles or workbench plugins; assumes OPTKIT-021 contracts, the OPTKIT-022 scaffold, and OPTKIT-016–018 backend services.
---

# Propose Workspace: Draft Document and RRF Vertical Slice

## Overview

Successor (with OPTKIT-021/022) to the superseded OPTKIT-019, inheriting its
exit criterion verbatim: open a case, change RRF k, understand invalidation,
preview honestly, declare intent, seal, run, and compare against the declared
intent. On PBUI this becomes the Propose workspace: four linked tiles bound to
one `ragttc.proposal-draft/v1` document, accept-driven mutation and evidence
attachment, a debounced pure-compile loop, a plugin registry whose proving
plugin is fusion, and a danger-verb seal crossing into the optkit journal.

**Program position:** Gate 6 first half (scalar variable end to end);
prerequisite for OPTKIT-020's asset plugin and OPTKIT-024's agent seat.
Depends on OPTKIT-021/022 plus the backend chain OPTKIT-015–018 (steps 1–6
can run on compiler-CLI fixtures before OPTKIT-018 lands).

- [Intern guide](design-doc/01-intern-guide-propose-workspace-draft-document-and-rrf-slice.md)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- design
- implementation
- ui
- rag-ttc
- optkit

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
