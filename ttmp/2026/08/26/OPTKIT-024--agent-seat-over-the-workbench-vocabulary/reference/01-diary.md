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
