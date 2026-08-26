---
Title: React Workbench Framework and RRF Vertical Slice
Ticket: OPTKIT-019
Status: active
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
