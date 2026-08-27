---
Title: Diary
Ticket: OPTKIT-022
Status: active
Topics:
    - design
    - implementation
    - ui
    - rag-ttc
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-08-26T23:18:00.562906497-04:00
WhatFor: ""
WhenToUse: ""
---

# Diary

## Goal

<!-- What is the purpose of this reference document? -->

## Context

<!-- Provide background context needed to use this reference -->

## Quick Reference

<!-- Provide copy/paste-ready content, API contracts, or quick-look tables -->

## Usage Examples

<!-- Show how to use this reference in practice -->

## Related

<!-- Link to related documents or resources -->

## Goal

Chronological record of the ragttc workbench implementation: the guide
retargeting to pbui 0.8.0 and the build of the product package and evidence
tiles.

## Step 1: Retarget the 022/023/024 Guides to pbui 0.8.0

The greenfield ruling deleted pbui's legacy surfaces (pbui repo
PBUI-ACTIONS-3, commit 6efeaeb → version 0.8.0), so the three
implementation guides were rewritten before any code: this product must be
born kernel-native, and its guides must never name a deleted mechanism.

### Prompt Context

**User prompt (verbatim):** (see OPTKIT-021 diary step 6 — same directive:
"perfect, update the design docs to have the simplification that not having
to worry about lgacy affords us." followed by the five-step execution plan)

**Assistant interpretation:** Rewrite the 022/023/024 guides against pbui
0.8.0 per the approved simplification plan.

**Inferred user intent:** Implementation (which starts immediately after)
follows the guides literally; stale guides would rebuild deleted APIs.

### What I did
- 022 §3: menu mechanics rewritten to kernel resolution (four-state
  availability, no bound verb on disabled rows, ambiguity as data, fresh
  revalidation, primary invocation for click-to-open).
- 022 §4 skeleton: added `pbui/actions.ts` (graph + rules + snapshotFor +
  `workbenchTileContributions()`) and `translators.ts`; runtime signature
  updated; descriptors marked representation-only.
- 022 §5: the case example split into a three-line descriptor and four
  kernel rules (autopsy as `primary: true`; evidence-attach unavailable with
  reason until 023).
- 022 §10/§11: "descriptor tests assert verbs" replaced by GOLDENS-FIRST
  menu row specs (`registry.resolve()` against literal snapshots — the
  pbui-workbench standalone-row-spec pattern); build order step 2 updated.
- 023: variable primary rule wording; translator edge names
  (`ragttc.delta-to-case`, `ragttc.layer-to-section`); SealBar renders the
  resolved `proposal.seal` action (capability + sealability in one rule) and
  clicks go through `performAction` — the hand-rolled stale guard now
  protects only the preview cache.
- 024: vocabulary export declared GENERATED from the registry/graph/
  translators/descriptors (PBUI-ACTIONS-3 Phase B generator, golden-pinned);
  approval implemented as a one-shot capability grant re-resolved through
  the kernel.

### Why
- The guides are the implementation contract; OPTKIT-022 starts from them.

### What worked
- No ADR decision changed — every edit was vocabulary-level substitution,
  which is itself evidence the kernel matches the design's intent.

### What didn't work
- N/A

### What I learned
- The goldens-first inversion is the natural greenfield form of the
  PBUI-ACTIONS-2 equivalence fence: with no old implementation to record
  from, the golden IS the spec rather than the fence.

### What was tricky to build
- The SealBar/kernel split in 023: deciding that fresh revalidation replaces
  the seal-path stale guard (but NOT the compile-cache guard) required
  tracing which staleness each guard actually protects against.

### What warrants a second pair of eyes
- 022 §5's rule ids (`ragttc.case.*`) and action ids (`case.*`) — they
  become frozen wire names under ADR L once implemented.

### What should be done in the future
- Implementation begins at build-order step 1 (scaffold) in the next step.

### Code review instructions
- Diff the three guides at this commit; grep them for `disabledBecause`,
  `conversions`, `activate` — all gone except historical notes.

### Technical details
- pbui 0.8.0 API: `createPbui({registry, defaultEnvironment, actions,
  snapshotFor, translators})`; `metadata.primary`; four-state availability.
