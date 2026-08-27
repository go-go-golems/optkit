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

## Step 2: The Seal Landed — and What the Two Attempts Taught (rag-ttc f7d67fcde)

The exit criterion is met: on a fresh checkpointed campaign store, the full
RRF slice ran through the UI and the user's SEAL click produced
candidate:b1393360bf8dcd…, verified in the journal (19 events, every direct
and nested payload). The FIRST seal attempt failed honestly with the
server's own campaign_conflict — the fixture campaign had been resumed to
completion for cockpit evidence, and OPTKIT-017's state machine refuses new
candidates on a completed campaign. That failure, and the re-run, surfaced
two real defects now fixed.

### Prompt Context

**User prompt (verbatim):** "i clicked." (the campaign_conflict attempt),
then "done" (the successful seal on the fresh store); earlier in the same
arc: "shouldnt' case:campaign:…/q-comparison remove be a semantic obejct?"
and "can the \"remove\" be added as a contextual action with our new action
framework?"

**Assistant interpretation:** Complete the exit criterion; fix what the
live run exposes; make evidence entries live presentations with a
contextual kernel remove rule.

**Inferred user intent:** The slice should not merely work once — the
failure modes and interaction seams should be honest and kernel-native.

### What I did
- Restarted the fixture server on a fresh checkpointed store (running
  campaigns accept candidates; campaign ids are journal-assigned, so the
  draft had to be re-authored against the new id).
- Fixed draft-activation ping-pong: ProposalApp activates its draft on
  MOUNT only. The body-keyed effect re-fired on every workbench document
  clone, so an older tile stomped a newer draft's activation — observed as
  a mutation landing on the dead draft.
- SealBar renders the seal attempt's outcome via a fixed RTK cache key
  (`fixedCacheKey: "seal"`); the campaign_conflict had shown only in the
  trace.
- Evidence entries re-materialize as live presentations (refFromKey, the
  read-time inverse of refKey) and "Remove from evidence" became an
  inherited contextual rule on `inspectable` gated by the new
  draftEvidenceKeys fact — appearing on an attached object's menu wherever
  it is presented, with attach turning INAPPLICABLE once attached (an
  attach that would touch nothing must not exist).
- Recorded the fusion parity fixture from the live probe and pinned it:
  every contribution equals weight/(k+rank) at recorded precision; the
  OPTKIT-015 evidence numbers appear verbatim (before 1/61 =
  0.01639344262295082, after 1/21 = 0.047619047619047616).

### Why
- Tasks gs7w and 5thg; the user's two review points; the honest failure.

### What worked
- Every failure surfaced with a stable machine code end to end: server →
  sink → trace (and now SealBar). Nothing was invented client-side.
- The compile digest was IDENTICAL across stores for identical authoring
  input — content-addressing making the re-author cheap to trust.

### What didn't work
- Attempt one: campaign_conflict (completed campaigns refuse candidates) —
  environmental, not a defect, but it exposed the two real defects above.
- Two intent tiles with the same aria-labels broke strict-mode locators
  during the re-drive; targeting by emptiness worked, and discarding the
  dead draft cleaned the board.

### What I learned
- "The workbench store clones the whole document on every mutation" is a
  sharp edge for any effect keyed on derived document state; key on ids.

### What was tricky to build
- The parity fixture needed the REAL draft digest (my reconstruction of it
  was wrong once — the server's 409 on a stale digest is itself the
  contract working); the fixture records the server's response verbatim
  and the test's closed-form check makes hand-patching impossible.

### What warrants a second pair of eyes
- Backend owner: candidate cand-b56f41 exists as a candidate record; the
  comparison projection rightly refuses it as an arm until a trial
  materializes it — confirm that is the intended surface for sealed-but-
  untrialed candidates.

### What should be done in the future
- trial.run over HTTP (or a documented CLI handoff) — refused honestly in
  the UI this slice; frontend sync to the workbench document host; the
  022 tails (plot adapter, a11y sweep, worst-first projection).

### Code review instructions
- rag-ttc f7d67fcde and 0ba28d7fd: src/apps/ProposalApp.tsx (mount-only
  activation), src/apps/IntentApp.tsx (SealBar outcome, evidence
  presentations), src/pbui/actions.ts (evidence-remove rule + attach
  partition), src/test/rrf-parity.test.ts.
- Validate: pnpm test (17); campaign verify on the store shows 19/19/5.

### Technical details
- candidate:b1393360bf8dcd… on campaign:3f84c0554dcfe54df51919e994352c1b,
  arm_id cand-b56f41, parent limit-1, draft digest sha256:b56f41ebc27ea….
