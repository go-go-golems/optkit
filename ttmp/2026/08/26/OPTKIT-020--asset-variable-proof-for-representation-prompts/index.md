---
Title: Asset Variable Proof for Representation Prompts
Ticket: OPTKIT-020
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
Summary: Prove the workbench supports sensitivity-aware artifact mutations using a representation summary prompt and bounded preview.
LastUpdated: 2026-08-26T14:20:30.531858526-04:00
WhatFor: Demonstrate that the scalar architecture generalizes to expensive content-addressed configuration without fake instant behavior.
WhenToUse: Implement after the RRF workbench vertical slice is complete.
---

# Asset Variable Proof for Representation Prompts

## Overview

This ticket is the final architecture proof. It registers an artifact-valued representation prompt, stores content by digest and sensitivity, renders an old/new diff, computes the full recomputation bill, and runs an honest bounded server preview for one chunk through the same compiler/sealer/workbench flow used by scalar variables.

**Program position:** final generalization gate for OPTKIT-011.

- [Intern guide](design-doc/01-intern-guide-to-artifact-valued-prompt-variables-and-bounded-previews.md)
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
