---
Title: Optkit Semantic Catalog and Executable Variable Bindings
Ticket: OPTKIT-013
Status: active
Topics:
    - architecture
    - design
    - implementation
    - optkit
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Add discoverable semantic catalogs and serialized-to-typed executable bindings without replacing Optkit's existing variable and patch algebra.
LastUpdated: 2026-08-26T14:20:21.132009032-04:00
WhatFor: Let manifests, CLIs, and UIs enumerate and safely apply the same typed variables used by Go code.
WhenToUse: Implement after OPTKIT-012 and before registering RAG variables or compiling proposals.
---

# Optkit Semantic Catalog and Executable Variable Bindings

## Overview

This ticket provides the domain-neutral framework beneath every optimization surface. It adds ordered sections, complete value schemas, deterministic catalog identity, executable bindings, and richer candidate intent while preserving `Variable[C,V]`, lenses, codecs, domains, and `PatchBuilder` as the canonical mutation mechanics.

**Program position:** generic Optkit foundation; consumed by OPTKIT-014 through OPTKIT-020.

- [Intern guide](design-doc/01-intern-guide-to-optkit-catalogs-domains-bindings-and-candidate-intent.md)
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
