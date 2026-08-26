---
Title: Proposal Compiler and Glazed CLI
Ticket: OPTKIT-016
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
Summary: Compile serialized mutations into pure auditable proposal drafts and expose the workflow through structured Glazed CLI commands.
LastUpdated: 2026-08-26T16:54:56.593205555-04:00
WhatFor: Give manifests, CLI users, and future browser authoring one deterministic validation and planning service.
WhenToUse: Implement after the catalog, PipelineConfig, and first real RAG variables exist.
---


# Proposal Compiler and Glazed CLI

## Overview

This ticket defines the draft half of the workbench lifecycle. `CompileProposal` decodes, validates, normalizes, and applies mutations in memory, then returns before/after values, graph differences, invalidation, diagnostics, and preview capabilities without creating artifacts or journal events. Glazed commands make that contract inspectable before any UI work.

**Program position:** backend authoring application; direct prerequisite for sealing and HTTP authoring.

- [Intern guide](design-doc/01-intern-guide-to-pure-proposal-compilation-and-glazed-cli-authoring.md)
- [Implementation diary](reference/01-implementation-diary.md)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **complete** — pure typed proposal compilation, deterministic diagnostics/digest, catalog and proposal Glazed commands, shared graph planning, and repeated no-write proofs all pass.

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
