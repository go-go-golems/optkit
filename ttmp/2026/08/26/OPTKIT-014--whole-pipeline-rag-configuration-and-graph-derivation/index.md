---
Title: Whole-Pipeline RAG Configuration and Graph Derivation
Ticket: OPTKIT-014
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
Summary: Introduce a patchable whole-pipeline RAG semantic configuration and derive the twelve-layer graph from it.
LastUpdated: 2026-08-26T14:20:22.368506158-04:00
WhatFor: Eliminate drift between executable values, Optkit snapshots, and manually authored graph identities.
WhenToUse: Implement after OPTKIT-012/013 and before adding real fusion variables.
---

# Whole-Pipeline RAG Configuration and Graph Derivation

## Overview

This ticket repairs the current retrieval-only snapshot asymmetry. It defines `PipelineConfig`, relocates semantic config ownership away from the campaign adapter, lifts layer-local lenses into the aggregate, and makes graph derivation a deterministic consequence of semantic values.

**Program position:** RAG configuration foundation; required by OPTKIT-015 and proposal compilation.

- [Intern guide](design-doc/01-intern-guide-to-pipelineconfig-layer-lenses-and-derived-graphs.md)
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
