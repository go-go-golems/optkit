---
Title: Build the RAG-TTC Glazed Experiment Workbench CLI
Ticket: OPTKIT-006
Status: active
Topics:
    - implementation
    - optkit
    - rag-ttc
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Implement a manifest-driven Glazed CLI that validates configuration and operates durable RAG-TTC experiments before UI work begins."
LastUpdated: 2026-08-25T23:59:00Z
WhatFor: "Track design, implementation, validation, and publication of the experiment operator workbench."
WhenToUse: "Use for all OPTKIT-006 planning and status."
---

# Build the RAG-TTC Glazed Experiment Workbench CLI

## Overview

Build the CLI proof layer for the specialist UI. The workbench will strictly validate versioned experiment manifests, resolve and compare layered configuration graphs, explain invalidation, dry-run exact work and budgets, and create, resume, inspect, and verify deterministic durable campaigns through Glazed structured output.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

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
