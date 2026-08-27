---
Title: OPTKIT-004 Phases 0-2 Implementation Plan
Ticket: OPTKIT-005
Status: active
Topics:
    - optkit
    - rag-ttc
    - implementation
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Implementation plan for release stabilization and OPTKIT-004 contract, Judgekit, and layered-configuration phases."
LastUpdated: 2026-08-25T16:55:00-04:00
WhatFor: "Guide implementation and review of the prerequisites for the specialist UI."
WhenToUse: "Use while implementing OPTKIT-004 Phases 0-2 and validating readiness for projector and UI work."
---

# OPTKIT-004 Phases 0-2 Implementation Plan

## Executive Summary

This ticket prepares the factual and identity contracts required by the OPTKIT-004 specialist UI. It first makes the current local integration independently buildable, then freezes the cross-layer RAG schemas and fixture, proves the existing Judgekit adapter conforms to the sealed-answer measurement contract, and implements a deterministic layered configuration diff and invalidation planner.

The ticket does not implement bundle-build campaigns, adaptive optimization, browser mutation, or the complete specialist UI. Its exit condition is a stable backend contract on which bounded server-side projectors and read-only screens can rely.

## Problem Statement

The completed retrieval campaign proves durable execution, replay, attribution, deterministic evaluation, and historical Judgekit measurement. The planned specialist UI still lacks three prerequisites: a canonical fixture spanning context, answer, and judge lineage; an audited conformance statement for the already-built Judgekit integration; and a semantic model that explains which configuration layer changed, which artifacts can be reused, and which stages require recomputation.

Opaque digests and generic journal events are insufficient for a scientific comparison interface. The UI must present stable schema names, typed identities, explicit dependency edges, and deterministic invalidation explanations without becoming an alternate authority.

## Scope and Phases

### Phase R: Release and isolated-module stabilization

- Inventory current module versions, replacements, remotes, and release state.
- Publish or pin real Optkit and Judgekit dependencies where repository authority permits.
- Remove synthetic `v0.0.0` and workspace-only assumptions from isolated RAG-TTC builds.
- Run tests, race tests, builds, vet, lint, and release hooks without local workspace replacement.

Exit gate: affected modules pass their isolated validation or any external publication blocker is precisely documented with a non-synthetic pinned alternative.

### Phase 0: Contract freeze and fixtures

- Define strict versioned schemas for layered configuration references, context, answer, and judge lineage.
- Preserve the stable retrieval-stage vocabulary.
- Extend the canonical semantic fixture through context construction, sealed answer, measurement epoch, report, and observation relationships.
- Define hard constraints and primary metrics as machine-readable contract data.

Exit gate: the canonical fixture has a version and digest, every decoder rejects unknown fields, and retrieval parity remains green.

### Phase 1: Judgekit conformance over sealed answers

- Map every OPTKIT-004 Phase 1 requirement to concrete implementation and tests.
- Add only missing conformance tests or fixture links; do not duplicate the adapter.
- Prove historical remeasurement under a second epoch without retrieval or answer execution.
- Preserve all old observations unchanged.

Exit gate: a focused conformance validator passes and archives exact evidence.

### Phase 2: Layered configuration graph

- Add RAG-TTC-owned typed configuration references and dependency edges.
- Define semantic identity rules and tests.
- Implement deterministic configuration diff, invalidation, and reuse planning.
- Add a promotion-manifest skeleton that references candidate, graph root, evidence, and policy identity without making promotion decisions.

Exit gate: judge-only changes schedule no upstream work, answer-only changes reuse admitted evidence, fusion-only changes reuse channel outputs, chunker changes invalidate every downstream layer, and repeated planner runs are byte-stable.

### Phase F: Final validation and UI handoff

- Run focused and repository-wide tests, race tests, builds, vet, lint, boundary guards, fixture verification, and deterministic planner checks.
- Update the diary, task state, changelog, and related-file evidence.
- Produce a projector/UI handoff enumerating stable contracts and intentionally deferred producers.

## Design Decisions

1. **Contracts precede projectors.** Browser and query schemas must not infer domain meaning from generic events.
2. **Phase 1 is audited, not rebuilt.** OPTKIT-002 P7 already supplies the product-owned Judgekit instrument.
3. **The configuration graph starts in RAG-TTC.** RAG layer names do not belong in domain-neutral Optkit core without a second-domain reuse proof.
4. **Invalidation is semantic.** Content-addressed identities identify values; explicit dependency rules determine recomputation.
5. **Promotion remains skeletal.** This ticket defines attributable references but defers gates, review, Pareto policy, and production mutation.
6. **Existing unrelated worktree state is preserved.** The unstaged `optkit/store/sqlite/rows.go` edit is excluded from all ticket commits.

## Validation Strategy

- Strict JSON decoder tests including unknown-field rejection.
- Golden fixture digest and cross-reference verification.
- Existing retrieval semantic fixture comparison.
- Historical Judgekit remeasurement with execution counters proving no product rerun.
- Table-driven invalidation tests for every layer.
- Property checks for deterministic diff ordering and transitive invalidation.
- Isolated module test/build/race/vet/lint gates.

## Deferred Work

- Bundle-build campaigns and partial-publication recovery.
- Rich retrieval/fusion/reranker search spaces.
- Context and answer campaigns.
- RAG-specific query projectors.
- Specialist browser screens.
- Pareto computation, review, decisions, and executable promotion.
- Adaptive search.

## References

- `sources/optkit-clean-slate-architecture-and-porting-guide.md`
- `optkit/ttmp/2026/08/25/OPTKIT-004--end-to-end-rag-optimization-workbench-and-specialist-ui/design-doc/01-intern-guide-to-the-rag-optimization-workbench-and-specialist-ui.md`
- `optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/reference/01-implementation-diary.md`
- `rag-ttc/pkg/ttc/judgeinstrument/`
- `rag-ttc/pkg/ttc/optkitcampaign/`
- `optkit/space/`
