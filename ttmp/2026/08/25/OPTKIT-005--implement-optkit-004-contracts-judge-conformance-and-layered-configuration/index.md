---
Title: Implement OPTKIT-004 Contracts, Judge Conformance, and Layered Configuration
Ticket: OPTKIT-005
Status: active
Topics:
    - optkit
    - rag-ttc
    - implementation
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Implement release stabilization and OPTKIT-004 Phases 0-2 as prerequisites for the specialist UI."
LastUpdated: 2026-08-25T16:55:00-04:00
WhatFor: "Track stable contracts, Judgekit conformance, and layered configuration invalidation before projector and UI work."
WhenToUse: "Use for implementation status and navigation across OPTKIT-005 deliverables."
---

# Implement OPTKIT-004 Contracts, Judge Conformance, and Layered Configuration

## Overview

This ticket implements the prerequisites for the OPTKIT-004 specialist UI. It stabilizes isolated module dependencies, freezes cross-layer RAG contracts and fixtures, audits the existing sealed-answer Judgekit integration, and implements deterministic layered configuration diff and invalidation planning.

## Key Links

- [Implementation plan](./design-doc/01-optkit-004-phases-0-2-implementation-plan.md)
- [Implementation diary](./reference/01-implementation-diary.md)
- [Tasks](./tasks.md)
- [Changelog](./changelog.md)
- Parent design: `../OPTKIT-004--end-to-end-rag-optimization-workbench-and-specialist-ui/`

## Status

Current status: **active**

Current phase: **Phase 1 — Judgekit conformance over sealed answers**

## Topics

- optkit
- rag-ttc
- implementation

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Phase Order

1. Phase R — release and isolated-module stabilization
2. Phase 0 — contract freeze and fixtures
3. Phase 1 — Judgekit conformance over sealed answers
4. Phase 2 — layered configuration graph
5. Phase F — final validation and UI handoff
