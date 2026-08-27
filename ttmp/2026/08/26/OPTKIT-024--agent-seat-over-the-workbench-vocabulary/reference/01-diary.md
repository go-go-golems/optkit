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
