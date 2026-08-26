---
Title: Proposal Sealing Candidate Manifests and Campaign Persistence
Ticket: OPTKIT-017
Status: complete
Topics:
    - architecture
    - design
    - implementation
    - optkit
    - rag-ttc
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Seal validated proposal drafts through canonical Optkit mechanics and persist patch-style manifest candidates as complete campaign facts.
LastUpdated: 2026-08-26T17:37:20.657943926-04:00
WhatFor: Make concise candidate authoring durable, reproducible, and explainable without retaining the source manifest.
WhenToUse: Implement after OPTKIT-016 establishes pure compilation.
---


# Proposal Sealing Candidate Manifests and Campaign Persistence

## Overview

This ticket defines the durable half of proposal authoring. It replays normalized draft values through typed bindings and `PatchBuilder`, extends strict manifests with candidate proposals, persists snapshots/patches/intent/catalog provenance, and emits the candidate journal fact used by later projections.

**Program position:** persistence gate; required before candidate-aware APIs and UI.

- [Intern guide](design-doc/01-intern-guide-to-proposal-sealing-candidate-manifests-and-durable-campaigns.md)
- [Implementation diary](reference/01-implementation-diary.md)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **complete** — strict v3 candidate manifests compile through the shared service, seal through PatchBuilder, persist complete campaign facts, record idempotently, and survive source-manifest removal/restart.

## Topics

- architecture
- design
- implementation
- optkit
- rag-ttc

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
