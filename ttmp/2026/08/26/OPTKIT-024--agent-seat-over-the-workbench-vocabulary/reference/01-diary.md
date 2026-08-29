---
Title: Diary
Ticket: OPTKIT-024
Status: active
Topics:
    - design
    - ui
    - rag-ttc
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-08-27T02:19:51.138919015-04:00
WhatFor: ""
WhenToUse: ""
---

# Diary

## Goal

Chronological record of the agent-seat implementation: vocabulary export,
chat integration, verb router with attribution, approval flow, agent
authorization, and the exit-criterion session.

## Step 1: Vocabulary Export Build Step (rag-ttc c2204f985)

The vocabulary the agent reasons over is now a generated, committed
artifact. pbui 0.9.0's registry.vocabulary() supplies the kernel half —
types with the graph, every rule with danger/primary — and the product
composes it with the translator edges (read live), one doc line per
concrete type (OPTKIT-021 §7.2 prose), and a verb table typed against the
ProductVerb union. `pnpm vocab` writes src/agent/vocabulary.json; the
golden test pins it.

### Prompt Context

**User prompt (verbatim):** "4. PBUI-ACTIONS-3 Phase B, then OPTKIT-024.
The perform envelope (onPerform(verb, {action, candidateId, actor…})) and
the registry-generated vocabulary export land as one small pbui release —
both are pulled directly by 024's task list — and then the agent seat is
implementable: the agent principal is just a snapshot without the seal
capability, and the approval flow is a"

**Assistant interpretation:** With Phase B shipped (pbui 0.9.0), start
OPTKIT-024 at its first task: the vocabulary export with a golden test.

**Inferred user intent:** The agent seat should be grounded in a
generated, pinned statement of what exists before any chat wiring.

**Commit (code):** rag-ttc c2204f985 — "workbench-ui: agent vocabulary
export build step with golden (OPTKIT-024 P1)"

### What I did
- `src/agent/vocabulary.ts`: buildAgentVocabulary() — schema_version,
  product, vocabulary_version "optkit-021/v1", the kernel registry
  vocabulary, TYPE_DOCS (one line per concrete type), VERB_SPECS
  (`Record<ProductVerb["kind"], {doc, fields, danger?}>` — exhaustiveness
  typechecked), conversions from workbenchTranslators.
- `pnpm vocab` = the golden test with UPDATE_VOCAB=1; the test otherwise
  compares committed vs generated exactly.
- Cross-check tests: danger set == {proposal.seal, trial.run}; the kernel
  trial rule carries danger:true; every non-abstract registry type is
  documented; conversions mirror the eight translator edges.
- Declared `inspectable`/`watchable` abstract in the type graph — the
  comment always claimed it, only the vocabulary reads the flag.

### Why
- Task spir; ADR L reserved the generated vocabulary, and the agent prompt,
  validator, and describe-types tool all read this artifact.

### What worked
- The typed verb table: TypeScript's exhaustiveness makes "verb kind
  renamed but table stale" a compile error, which is most of what a
  hand-maintained vocabulary gets wrong.

### What didn't work
- First golden run failed usefully: `inspectable` surfaced as a concrete
  undocumented type because the graph never DECLARED it abstract. The
  export audit caught a latent modeling slip within minutes of existing.

### What was tricky to build
- Deciding the seal has no kernel rule to cross-check (the SealBar is
  chrome-owned pbui.perform), so the danger agreement pins the trial rule
  and the verb-table set — documented in the test so nobody "fixes" the
  asymmetry blindly.

### What warrants a second pair of eyes
- VERB_SPECS field prose is hand-written (the shapes are typed, the prose
  is not) — a reviewer should read it against verbs.ts once.

### What should be done in the future
- P2: chat tile + mention resolver; P3 verb router; P4 approval; P5 agent
  principal; P6 exit criterion.

### Code review instructions
- Start at apps/workbench/web/src/agent/vocabulary.ts, then the golden
  test; regenerate with `pnpm vocab` and confirm no diff.

## Step 2: The Codec Ruling and the Chat Layer Mount (pbui e09ab55, 917c04a; rag-ttc b76000a1d)

The wire-convention fork was put to the program owner: migrate the
product's 22 value shapes to value-IS-the-wire-reference, or give pbui-chat
a product reference codec. Ruling: the codec ("ok, do it" on option 2).
pbui-chat 0.3.0 gained `createPbuiChat({referenceAdapter})` with the old
convention as the default; rag-ttc supplies one adapter over its existing
refKey/refFromKey and keeps its proven presentation layer untouched. On
top of that, the whole chat layer mounted: chat.Provider now routes EVERY
verb — human click or agent output — through one validated, attributed
path to the sink.

### Prompt Context

**User prompt (verbatim):** "1. would be tsx side? or what would it mean?"
then "ok, do it."

**Assistant interpretation:** Explain what the full value-shape migration
would entail; on the ruling, implement the codec option and continue P2.

**Inferred user intent:** Keep the proven 022/023 layer stable; pay the
integration cost in the chat package where it belongs.

**Commits (code):** pbui e09ab55 (referenceAdapter, chat 0.3.0), pbui
917c04a ('any' verb field type, TS+Go), rag-ttc b76000a1d (chat layer).

### What I did
- pbui-chat 0.3.0: `ReferenceAdapter{toProduct, fromProduct}` option;
  every wire↔product crossing (RefPresentation, composer insert-object,
  accept tool, form ref fields, label lookups, router accept) goes through
  the context's adapter; identity default keeps all consumers green.
- Wire contract: new 'any' verb field type (TS schema+validate+zod
  derivation, Go vocabulary.go) — a mutation value is typed by the
  CATALOG, and a coarse wire type would refuse legitimate scalars.
- rag-ttc src/chat/: zod verb schema typechecked against ProductVerb;
  ONE wire vocabulary (replaces the P1 bespoke artifact — pnpm vocab now
  writes src/chat/vocabulary.json, Go-embeddable); refs codec;
  router (conversation verbs → chat layer, all else → sink with actor
  attribution); chat instance; conversation pointer documents
  (ragttc.conversation/v1) reconciled so chat tiles satisfy the host.
- `unresolved` presentation type: a mention nothing answers renders as an
  inert chip, never a fabricated value.
- sink.onPerform(verb, {actor}); all trace writes through one bound
  recorder.
- Go host: conversation + pbui.widget formats, chat catalog entries; pbui
  module pin bumped to upstream main (pkg/pbuichat exists there), which
  pulled pinocchio 0.11.16 — the customer appserver's manifest handler now
  forwards the required executor identity, with its test updated.
- Live verification: app boots with the chat layer; Inspect on an arm
  lands as '#1 human inspect the arm …' through router → sink.

### What didn't work
- Module-cycle TDZ crashes, twice: the runtime's facts imported the
  workbench (fixed with workbenchSlot.ts, the reference product's
  dependency-light-slot pattern), and the router's sink import pulled the
  workbench into the chat's evaluation chain (fixed with a lazy import).
- A menu click silently did nothing: the router validated the PRODUCT-
  shaped ref (no wire id) and rejected — rejection is data, invisible.
  Fixed with the OUTBOUND codec on router.perform, plus a console.warn on
  rejected outcomes so the next silent pit announces itself.
- GOWORK=off lefthook gate failed on the old pbui pin (no pkg/pbuichat);
  bumping to upstream main broke the colleague's appserver via pinocchio
  0.11.16's new manifest identity — forwarded, tests green.

### What I learned
- Upstream pbui main moved TODAY (e54210fefb3b) — the local pbui branch
  task/use-optkit (0.9.0 + envelope + codec) will need a merge/rebase
  before anything is pushed; flag for the user.
- HMR + dynamic import("/src/…") probes create second module instances;
  assert through the DOM or window hooks set by the app itself.

### What warrants a second pair of eyes
- The double codec round-trip (product ref → wire → product) on every
  local menu verb — correct, but a reviewer should confirm the wire trace
  is worth the conversion (it is what the verb log will store).
- The appserver identity forward is minimal — whether the customer app's
  real client sends clientInstanceId/connectionId is the colleague's call.

### What should be done in the future
- P4 approval flow; P5 agent principal; P6 chat server mount + scripted
  exit session; verb-log reporting turns on with the server.

### Code review instructions
- Start at rag-ttc apps/workbench/web/src/chat/ (router.ts, refs.ts,
  vocabulary.ts), then the pbui-chat referenceAdapter (packages/pbui-chat
  src/types.ts + createPbuiChat.tsx).
- Validate: pnpm test + typecheck in the web app; GOWORK=off make test in
  rag-ttc; pnpm -r test in pbui.

## Step 3: Approval as a One-Shot Grant (rag-ttc c8bbe78f7)

The safety half of the agent seat. An agent's danger verb never executes:
the router's approval gate parks it as a pending request, the SealBar
renders it beside the hypothesis it argues from, and a human's approve
click spends a one-shot grant that re-performs the agent's OWN verb —
same sink, fresh revalidation, approvalId riding in the verb where the
trace and verb log keep it.

### Prompt Context

**User prompt (verbatim):** "ok, do it. finish the ticket. commit at
appropriate intervals and keep a detailed diary as you work (using the
diary format from the skill) Print out a brutalist work slip with the
plan / different phases for the ticket. then before stsarting a phase,
plrint a split about the phase, and print one when the phase is done."

**Assistant interpretation:** Complete OPTKIT-024's remaining phases (P4
approval, P5 agent principal, P6 exit session) with full working
discipline.

**Inferred user intent:** The ticket lands whole: agent proposes, human
approves, provenance durable.

**Commit (code):** rag-ttc c8bbe78f7

### What I did
- approvals Redux slice: requested → approved|denied → consumed; the
  four-state ledger, one-shot by construction.
- Router approval gate: agent + {proposal.seal, trial.run} without a
  grant → park + attributed trace refusal; with an approvalId → verify
  entry, status, verb kind, SUBJECT (docId/candidateId), then consume and
  delegate. Mismatch spends nothing.
- approvalId became a field of the two danger verbs (product union, zod
  schema, vocabulary regenerated) — the linkage lives in the verb, not in
  a side channel.
- ApprovalPrompt/PendingApprovals components: inline in the SealBar for
  the draft's own seal request, product-wide strip under the masthead so
  no request depends on the right tile being open.
- Seal proposer: {kind: "llm", identity: "agent:workbench"} when the
  actor is the agent. The §4 open question (approving human in candidate
  identity vs journal metadata) lands as recommended: metadata only — the
  approving human is the authenticated actor of the seal call. FLAGGED
  for ADR review, not silently decided: noted here and in the changelog.
- 5 loop tests: park, one-shot consume + replay refusal, deny, subject
  mismatch, human-passes-through.

### What worked
- The tests bound the REAL router against the REAL vocabulary — the gate,
  the validation, and the sink's honest refusals all exercised together;
  green on first run.

### What was tricky to build
- One-shot semantics: consume ON the attempt, not on success — a failed
  seal spends the grant, and the agent must re-request. The alternative
  (consume on success) would let a failing verb retry unbounded under one
  human decision, which is not what "approved once" means.

### What warrants a second pair of eyes
- The approving human's identity is only the seal call's bearer actor; if
  the program wants a named human IN the intent, that is an OPTKIT-018
  schema addition (ADR note).
- The pending-approvals strip renders from local Redux — approvals are
  per-browser until the verb log (P6) makes them conversation state.

### Code review instructions
- src/chat/router.ts (approvalGate), src/components/approvals.tsx,
  src/store/store.ts approvals slice; validate with
  `pnpm vitest run src/test/approvals.test.ts`.

## Step 4: The Agent Principal (rag-ttc 456a1d037)

Authorization catches up with the UI: the command API now serves multiple
bearer principals, and the agent's grant simply does not contain
proposal.seal. A 403 before any domain logic runs is what "the agent
cannot seal alone" means as a fact rather than a convention.

### Prompt Context

**User prompt (verbatim):** (see Step 3 — "finish the ticket")

**Assistant interpretation:** P5: the agent principal on the OPTKIT-018
authorizer.

**Commit (code):** rag-ttc 456a1d037

### What I did
- principals.go in the serve command package (zero footprint on the
  colleague's packages): grantAuthorizer (per-actor action sets, mirrors
  FixedAuthorizer semantics) and multiBearerAuthenticator (constant-time
  resolution over every candidate token).
- serve flags --agent-token (empty = no agent seat) and --agent-actor
  (default actor:agent-workbench — record.ValidateID requires the actor:
  namespace, which also corrected the frontend's proposer identity).
- HTTP test on a real fixture store: agent compile 200, agent seal 403,
  human seal passes authorization, unknown token 401.

### What didn't work
- "agent:workbench" as the actor id — the record layer requires the
  actor: namespace. Renamed to actor:agent-workbench everywhere (flags,
  test, frontend proposer identity), which is the better spelling anyway:
  one id space, one grep.

### What warrants a second pair of eyes
- The multi-authenticator compares the provided token against EVERY
  grant in constant time per grant; with two grants that is fine, with
  hundreds it would want a keyed lookup with a decoy compare.

### Code review instructions
- cmd/rag-ttc/cmds/experiments/optkitrag/principals.go + serve.go diff;
  validate with `GOWORK=off go test ./cmd/rag-ttc/cmds/experiments/optkitrag/`.

## Step 5: The Exit-Criterion Session — and the Contract It Found (rag-ttc 2e28e2cc1)

The session ran twice, and the first run's failure was the most valuable
artifact of the phase. The scripted agent (the router driven with
actor:"agent" — the model-backed agent minus the model) authored a draft
on a fresh checkpointed store, the human adjusted its mutation in the
open proposal tile, the agent's seal parked, and the human's FIRST
Approve click hit: "forbidden: candidate proposer does not match
authenticated principal" — the pre-agent-era seal rule forcing proposer
identity to the authenticated principal, which makes an approved agent
seal impossible by construction (the agent may not seal; the approver may
not carry proposer=llm). The human then refreshed, the parked approval
(per-browser, unpersisted) vanished, and the ordinary human seal button
sealed candidate ae4b46… with HUMAN provenance — a durable, honest, WRONG
answer to the exit criterion, caught by decoding the journal payload.

Second run, after the fixes: journal seq 21 CandidateProposed carries
proposer {kind: "llm", identity: "actor:agent-workbench"}, mutations
fusion.rrf_k 60 → 30 (the HUMAN'S adjusted value), the agent's
hypothesis; the verb trace holds "#31 agent proposal.seal
approval:appr-mtbqxgvm-1 performed"; the journal verifies 22/22 with all
nested payloads. The exit criterion is met.

### Prompt Context

**User prompt (verbatim):** "forbidden candidate proposer does not match
authenticated principal [interrupted] k i tried again and it worked this
time, maybe a refresh issue or so?" then "I do get a lot of
failures_unavailable campaign campaign:3f84c055… not found in the
different tiles" then "done"

**Assistant interpretation:** The first approve failed on a server rule;
the retry that "worked" needed forensic verification (it was the human
path); the stale tiles needed cleanup; the final "done" is the corrected
approve.

**Inferred user intent:** A truthful exit criterion — not a seal that
merely happened, but one whose provenance is right.

**Commit (code):** rag-ttc 2e28e2cc1 (seal rule + approvals persistence +
hook peeks + scripted-session test); pbui 3-file 'object'-accepts-arrays
fix committed separately in pbui.

### What I did
- experimentworkbench.Seal: human proposers must be the principal
  (unchanged anti-impersonation); non-human proposers seal through a
  seal-granted approving principal, declared identity kept as provenance,
  anonymous non-human refused. Test covers all three branches.
- Approvals slice persists per browser — a refresh can no longer orphan a
  parked request into the silent human-path trap.
- src/test/agent-session.test.ts: the scripted session as an integration
  test (mention resolution through the codec, co-edited draft with
  interleaved actor attribution, park → grant → one-shot attempt).
- The 'stop-after' seeding affordance (campaign run --stop-after
  after_lease + experimentworkbench.RunUntil) — 023 needed this and had
  no recorded path.
- Live session driven end to end with the user clicking both Approve
  attempts; stale dead-campaign tiles cleaned (their pointer documents
  discarded), which answered the user's failures_unavailable complaint.
- Wire-contract fix found live: 'object' verb fields now accept arrays
  (string lists) on both TS and Go sides — intent.setRisks was
  unperformable.

### What worked
- The layered honesty: the wrong seal was DETECTABLE because the journal
  records proposer provenance, the trace records actors, and nothing
  fabricates success. Every failure in the session announced itself.

### What didn't work
- open.draft-tile silently "performed" while placing nothing: openView
  needs a splittable target placement, the deep-corner first leaf hit the
  split-depth cap, and the router's outcome ("delegated") is not the
  sink's outcome ("refused"). Worked around with tile.activate on a
  shallow pane; TWO follow-ups flagged below.
- Dynamic import("/src/…") probes under vite again created a second
  module universe (its own workbench + its OWN document sync competing
  with the app's) — the reason first-run tiles vanished. The dev hook now
  carries same-instance read peeks; the diary's standing rule: after HMR,
  assert through the DOM or app-installed hooks only.

### What was tricky to build
- Deciding where the approving human lives durably: the amended rule
  keeps proposer=agent in the candidate; the approving human is the
  authenticated actor of the seal call and the trace's approvalId row.
  A NAMED approver inside the candidate intent would be an OPTKIT-018
  schema addition — flagged as ADR note, not decided silently.

### What warrants a second pair of eyes
- The seal-rule amendment is in the colleague's package
  (experimentworkbench) — small and tested, but the provenance semantics
  deserve their review.
- Router outcome vs sink outcome ("performed" = delegated): honest but
  confusing; consider propagating the sink's refusal into the router
  outcome.
- openView's no-splittable-target refusal is invisible to agents; a
  fallback placement search (any splittable leaf) would make "open"
  reliable regardless of activation state.

### What should be done in the future
- Task 75ds (release-gated): chat-server mount with the embedded
  vocabulary + scripted scenario or model profile; flip verb-log
  reporting on. Needs a published pbui carrying 'any' and the
  array-tolerant 'object'.
- The two openView/outcome follow-ups above.

### Code review instructions
- The seal rule: pkg/ttc/experimentworkbench/workbench_service.go (Seal)
  + workbench_service_test.go; the session test:
  apps/workbench/web/src/test/agent-session.test.ts.
- Validate: GOWORK=off go test ./pkg/ttc/experimentworkbench/ and pnpm
  test (47) in apps/workbench/web; the journal evidence lives in the
  scratchpad store-agent (campaign:51717a2f…, seq 21).
