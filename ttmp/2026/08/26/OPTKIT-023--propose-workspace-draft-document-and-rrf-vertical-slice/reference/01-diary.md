---
Title: Diary
Ticket: OPTKIT-023
Status: active
Topics:
    - design
    - implementation
    - ui
    - rag-ttc
    - optkit
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-08-27T00:05:53.163266352-04:00
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

Record the Propose-workspace implementation: the draft document, the compile
loop, the editors, the SealBar, and the live RRF vertical slice.

## Step 1: The Propose Workspace, Live to the Armed Seal (rag-ttc d195d9528, e2fda8c13)

The whole authoring path now runs in the browser against the live command
API: right-click an arm → "Draft a proposal from this arm"; the proposal
tile opens bound to a fresh ragttc.proposal-draft/v1 document; a left click
on "RRF rank constant" in the catalog adds the mutation through the
kernel's primary invocation; the RrfEditor sets k=20; the compile loop
re-queries and renders the server's digest, the normalized 60 → 20, and
"none — sealable"; evidence attaches from a case's menu through the
now-live accept rule; and the SealBar arms its confirmation. The final
"SEAL — durable" click was deliberately left to the user — the permission
classifier refused to let the agent press a durable-write button, which is
the approval symmetry OPTKIT-024 will formalize.

### Prompt Context

**User prompt (verbatim):** (see OPTKIT-022 diary step 1 — the five-step
plan; mid-implementation: "make sure to match the look of new widgets and
such to the overall theme for consistency" and "[Image] like here" showing
the native slider breaking the theme)

**Assistant interpretation:** Implement OPTKIT-023 per the retargeted
guide, keeping every new control in the turboproof family look.

**Inferred user intent:** A working authoring slice whose UI is of one
piece with the evidence tiles.

**Commit (code):** rag-ttc d195d9528 — draft format + Propose catalog
(Go); e2fda8c13 — the frontend workspace.

### What I did
- Go: ragttc.proposal-draft/v1 validator (authoring input only; derived
  keys rejected BY NAME — draft_digest/plan/diff/diagnostics/sealable/…),
  Propose apps in the host catalog, 9 table-driven validator tests.
- Authoring vocabulary: seven types as descriptors + rules (variable
  primary = mutation.add, greyed without a draft; arm → draft.create;
  draft → open intent/invalidation/preview siblings, discard as danger;
  candidate → compare primary + capability-gated trial).
- Verb families completed per ADR K: draft verbs are documentPut through
  the protocol applier and never touch the campaign store; command verbs
  dispatch RTK endpoints with the Bearer token; the token's presence IS the
  `seal` capability in the snapshot.
- The compile loop: the COMPILE REQUEST is the RTK cache key, so a stale
  response cannot attach to newer input — the newer input is a different
  key; 350ms debounce; two tiles bound to one draft share one in-flight
  compile.
- ValueSpec editors: specialized→generic→visible-fallback registry;
  RrfEditor (slider + exact number) as the first plugin.
- IntentApp with the SealBar: renders resolved sealability, names what is
  missing, mints the Idempotency-Key at confirmation, and after success
  writes exactly one field back (sealed: candidate id) — the derived-state
  rule's only sanctioned write.
- Theme pass on new controls: squared range slider (2px ink track, square
  authoring-tone thumb), checkbox accent, select, textarea, token input.

### Why
- OPTKIT-023's task list; the backend chain and pbui 0.8.0 were ready.

### What worked
- The kernel earned its keep again: "Attach as evidence" went from greyed-
  with-reason to live purely by a fact changing — no rule edits.
- The server's own diagnostics (empty_mutations, invalid_json) rendered
  with stable codes at each intermediate authoring state, exactly as the
  guide's honest-states section intended.

### What didn't work
- The doc-bound Propose tiles are invisible to the default launcher (it
  skips doc-bound apps) — first walk-through dead-ended after the proposal
  tile. Fixed with draft-menu navigation rows (open.draft-tile verbs), the
  draft chip's primary being the intent tile.
- The native range slider broke the family look (user screenshot); fixed
  with themed track/thumb CSS.

### What I learned
- "The request is the cache key" collapses the guide's debounce+revision-
  cache+stale-guard trio into one RTK property; the hand-rolled stale guard
  survives only as the debounce timer.

### What was tricky to build
- Seal request assembly: arm_id is authoring input the ADR I body does not
  carry; it is generated at confirmation (cand-<digest6>) like the
  idempotency key, keeping the document format unchanged. Flagged for ADR
  review.

### What warrants a second pair of eyes
- The seal path builds CandidateIntent from the draft body (motivation
  case_ids parsed from evidence keys) — worth checking against
  experimentworkbench's intent validation.
- ProposalApp activates its draft on render; two proposal tiles bound to
  different drafts will ping-pong the active fact. One-draft-at-a-time is
  the v1 assumption.

### What should be done in the future
- The user's seal click completes the exit criterion; then the fusion
  recorded-precision parity fixture (task gs7w) and a fresh-store replay
  (task 5thg). trial.run stays CLI-side this slice.

### Code review instructions
- Start at src/draft.ts and src/sink.ts (performDraft/performCommand),
  then src/compile.ts, then IntentApp. Validate: pnpm typecheck && pnpm
  test && pnpm build; GOWORK=off go test ./pkg/ttc/workbenchhost/.

### Technical details
- Live evidence: catalog semantic id sha256:d3034d1d…30b6 (matches the
  OPTKIT-015 recorded id); compile digest sha256:b56f41ebc27ea… for
  fusion.rrf_k=20 on limit-1.
