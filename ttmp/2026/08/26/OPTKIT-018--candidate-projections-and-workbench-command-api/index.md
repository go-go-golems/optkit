---
Title: Candidate Projections and Workbench Command API
Ticket: OPTKIT-018
Status: active
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
Summary: Project sealed candidate meaning for historical reads and expose proposal commands through a separate workbench application/API boundary.
LastUpdated: 2026-08-26T14:20:27.945954284-04:00
WhatFor: Let browsers consume authoritative facts and invoke authoring services without embedding storage or mutation rules.
WhenToUse: Implement after candidate sealing and campaign persistence are complete.
---

# Candidate Projections and Workbench Command API

## Overview

This ticket connects durable backend records to clients while preserving boundary clarity. `specialistapi` continues to project historical facts; a workbench command service owns compile, preview, and seal operations; thin HTTP adapters add authorization, idempotency, sensitivity, and stable errors.

**Program position:** backend-to-frontend contract; blocks the React workbench.

- [Intern guide](design-doc/01-intern-guide-to-candidate-read-projections-and-the-workbench-command-api.md)
- [Implementation diary](reference/01-implementation-diary.md)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

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
