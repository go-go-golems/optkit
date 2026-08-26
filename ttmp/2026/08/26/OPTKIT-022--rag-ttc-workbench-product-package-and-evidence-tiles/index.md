---
Title: RAG-TTC Workbench Product Package and Evidence Tiles
Ticket: OPTKIT-022
Status: active
Topics:
    - design
    - implementation
    - ui
    - rag-ttc
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Build the PBUI product package at rag-ttc/apps/workbench/web and prove it with read-only evidence tiles (campaigns, failures, judge, autopsy, chunk, compare, inspector, trace, watch) over the specialist API, plus the Go workbench document host.
LastUpdated: 2026-08-26T16:04:31.997349219-04:00
WhatFor: Establish the product scaffold, presentation registry, verb sink, and workbench host so OPTKIT-023 only adds tiles and plugins, never architecture.
WhenToUse: Read when implementing the product package or any evidence tile; assumes accepted OPTKIT-021 contracts.
---

# RAG-TTC Workbench Product Package and Evidence Tiles

## Overview

Successor (with OPTKIT-021/023) to the superseded OPTKIT-019. Builds the pbui
product package — presentation types, descriptors, verbs, verb sink, apps,
workspaces — and the Go `workbenchhost` package, proven by the read-only
evidence tiles over `specialistapi`. Requires only OPTKIT-021 contracts plus
one additive read projection (worst-first verdicts, coordinated with
OPTKIT-018); therefore runs in parallel with the backend chain OPTKIT-015–018.

**Program position:** the PBUI product scaffold; prerequisite for the
authoring workspace in OPTKIT-023.

- [Intern guide](design-doc/01-intern-guide-ragttc-workbench-product-package-and-evidence-tiles.md)

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
