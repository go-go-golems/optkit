---
Title: Implement Optkit and Pragmatic RAG-TTC Vertical Slice
Ticket: OPTKIT-002
Status: active
Topics:
    - optkit
    - rag
    - rag-ttc
    - coinvault
    - judgekit
    - implementation
    - migration
    - local-development
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://rag-ttc/pkg/ttc/search/service.go
      Note: P2 canonical direct retrieval service (commit ea6c629be)
ExternalSources: []
Summary: Phase-gated implementation of the Optkit foundation and a Coinvault-informed, RagKit-backed RAG-TTC product campaign.
LastUpdated: 2026-08-24T22:55:00-04:00
WhatFor: Navigate the implementation plan, phase status, diary, scripts, and validation evidence for OPTKIT-002.
WhenToUse: Start here before implementing, reviewing, or resuming any OPTKIT-002 phase.
---


# Implement Optkit and Pragmatic RAG-TTC Vertical Slice

## Overview

This ticket implements the architecture assessed in OPTKIT-001. Optkit remains a domain-neutral optimization and experiment framework; RagKit remains a separate RAG-domain library; RAG-TTC composes both as the first real product vertical slice.

## Key links

- [Detailed phased implementation plan](design-doc/01-phased-implementation-plan.md)
- [Strict implementation diary](reference/01-implementation-diary.md)
- [Tasks](tasks.md)
- [Changelog](changelog.md)
- [P0 archive import script](scripts/01-import-optkit-baseline.sh)
- [P0 validation script](scripts/02-validate-p0.sh)
- [P0 archive manifest](sources/01-optkit-baseline-manifest.txt)
- [P0 validation transcript](sources/02-p0-validation.txt)

## Current status

- Overall plan slip: printed.
- P0 plan slip: printed.
- **P0 Optkit baseline:** complete and green.
- **P1 semantic RAG fixtures:** complete and green.
- **P2 canonical retrieval service:** complete and green.
- P2.5 read-only campaign query plane: next.
- P3 runtime identity, routes, and policy: queued after UI bootstrap.
- P4-P8 and final delivery: pending.

## P0 result

The checked-out repository now contains the supplied Optkit implementation at archive revision `1786d1d`, normalized into the existing CI, lint, release, and ticket plumbing. Full CGO, no-CGO, race, vet, build, lint, demo, inspect, journal verification, artifact verification, and GoReleaser configuration checks pass.

## Stable boundary decision

RagKit primitives do not move into Optkit. Product-owned integration code composes Optkit orchestration and RagKit algorithms. Consolidation targets duplicated RagOpt and product-runner orchestration after behavioral parity is proven.

The UI is a read-only scientific query plane. LLM agents use CLI/application commands for campaign creation and mutation; the UI provides search, deep links, provenance navigation, replay, comparison, and visualization.
