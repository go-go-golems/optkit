---
Title: Real Fusion Configuration and First RAG Optimization Catalog
Ticket: OPTKIT-015
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
Summary: Make RRF fusion a real typed configuration and register the first runtime-honest RAG optimization variables.
LastUpdated: 2026-08-26T14:20:23.810428915-04:00
WhatFor: Prove catalog variables change actual RAG execution and produce correct graph diff and invalidation results.
WhenToUse: Implement after OPTKIT-013/014 and before the proposal compiler.
---

# Real Fusion Configuration and First RAG Optimization Catalog

## Overview

This ticket is the first RAG semantic proof. It preserves RRF's runtime `float64` contract, removes the fixture's hardcoded `60`, accurately names the current final-result limit, and registers both coordinates with reviewed documentation, domains, and executable bindings.

**Program position:** first real application variables; provides the proving coordinates for OPTKIT-016 and OPTKIT-019.

- [Intern guide](design-doc/01-intern-guide-to-fusion-configuration-rrf-and-the-first-rag-catalog.md)
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
