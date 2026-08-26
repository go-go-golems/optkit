---
Title: Specialist Work Scenarios and Questions
Ticket: OPTKIT-009
Status: active
Topics:
    - implementation
    - rag-ttc
    - ui
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources:
    - URL: https://hyperslop.systems/
      Note: Accent palette source (green 2db878, purple 805bd7, red ef4038, yellow f2ad00)
Summary: Persona-driven scenarios and the questions a RAG specialist actually asks, mapped to what the UI can answer today, soon, and only with backend support; plus the design consequences.
WhatFor: Drive the evolution of the specialist UI from journal-shaped tables toward a tool people do work in.
WhenToUse: Read when prioritizing specialist UI features or projector/backend asks.
---

# Specialist Work Scenarios and Questions

The v1 UI is contract-faithful but journal-shaped: it shows what the store
recorded, dominated by identifiers. This document walks the tool as its actual
user — an engineer tuning a RAG system — and derives the questions the UI must
answer, ranked by how often they come up in a working session.

## Scenarios (me, doing the work)

### A. Morning check — "did my sweep finish, can I trust it?"

I kicked off a campaign last night. I open the cockpit and want, in five
seconds: finished or not, anything failed, journal verified, budget headroom.
Then one promising challenger to click into. I do not want to cross-reference
four tables to learn "yes, it's fine."

**Wants:** a verdict-first cockpit header; failed/queued counts as words, not
a table I have to scan; a "most interesting comparison" shortcut.

### B. Did my change help? — the core loop

I bumped `retrieval.limit` from 1 to 2. The comparison screen must lead with
the sentence I came for: *"limit-2 beats limit-1: target-coverage
0.833 → 1.000 (+0.167 over 3 pairs)"*, then show me **what I changed in
actual values** (`limit: 1 → 2`), not two 64-char identities. Cases sorted by
|Δ| so the one case that moved (+0.500) is row one.

**Wants:** verdict strip; config *value* diff (needs backend: inline small
config payloads like stage previews — the retrieval config artifact is 69 B);
|Δ|-sorted cases; hashes out of the reading line.

### C. Why did this case improve/regress? — the debugging loop

`q-comparison` went from 0.5 to 1.0. I open the case and want the two
pipelines **side by side**, stages aligned by name, chunk sets diffed:
baseline retrieved `{chunk-b}`, challenger `{chunk-a, chunk-b}`, and the
divergence started at `vector.raw`. Then I want to read chunk-a's text and
the case's actual question to judge relevance myself.

**Wants:** side-by-side pipeline diff (buildable today — both routes exist);
case input text preview (needs backend, 135 B artifact); chunk text (needs
chunk-lab producer, deferred).

### D. Getting a feel — "let me poke at it"

Before trusting metrics I want to touch the system: browse what questions are
in the eval set, see what a case's ideal evidence is, and type my own query
against an arm's sealed snapshot to watch what comes back, stage by stage.
Reading tables never builds intuition; running five of my own queries does.

**Wants:** case browser with real question text; a read-only **query
playground** against a snapshot (non-mutating retrieval — a natural future
`POST /playground` that executes but never records); chunk browser.

### E. Judge skepticism — "says who?"

A 0.5 score is a judge's opinion. I want to see what the judge was shown and
what it said, and eventually calibration ("how often does this judge agree
with me?"). Deferred until answer/judge stages exist in trajectories.

### F. Writing it up — when hashes finally matter

When I share results or pin an appendix, *then* I need exact identifiers:
campaign, snapshots, digests. They belong in a detail tray on each screen —
copyable, complete, out of the way. Deep links carry the rest.

### G. Planning the next run — "what will it cost me?"

If I change chunking, which layers recompute? The invalidation plan already
answers this for two existing arms; a what-if against a hypothetical config
needs backend support. Budget panel answers "can I afford another sweep."

## The questions, ranked

Working questions (every session):

1. Did the challenger win, by how much, over how many pairs?
2. What did I actually change (values, not identities)?
3. Which cases moved, and which single case should I look at first?
4. Where in the pipeline did baseline and challenger diverge?
5. What was the question asked in this case?
6. What text is in the chunks that made the difference?
7. Anything failed / unverified / over budget that makes this run untrustworthy?

Exploration questions (getting a feel):

8. What's in the eval set — what kinds of questions, which are negatives?
9. What comes back if *I* ask my own question against this arm?
10. Why was this chunk dropped at this stage (filter, dedupe, rerank)?
11. What did the judge see and say; do I agree with it?

Bookkeeping questions (occasionally):

12. What exactly is this artifact (digest, size, sensitivity) — for citation?
13. Which episodes does this number rest on (provenance chain)?
14. How much budget is left; what did this campaign consume?

## Answerability map

| Question | Status |
| --- | --- |
| 1, 3, 7, 14 | Answerable now — presentation changes only |
| 4 | Answerable now — new side-by-side diff view over existing routes |
| 12, 13 | Answered by provenance screen (keep, restyle) |
| 2, 5 | Small backend ask: value previews on config-diff layers and case inputs, same sensitivity/size policy as stage previews |
| 8 | Partially now (case IDs + modes); real text needs case-input previews |
| 6, 10 | Needs chunk-lab producer (explicitly deferred in OPTKIT-007) |
| 9 | Needs a read-only playground endpoint (future contract) |
| 11 | Needs answer/judge trajectory stages + calibration projector (deferred) |

## Design consequences (this ticket)

- Verdict-first comparison header; |Δ|-sorted cases.
- Hashes leave the reading line everywhere: identifier trays/disclosures per
  panel; provenance remains the full-detail screen.
- Friendlier surface: white background stays; drop the brutalist offset
  shadows; hyperslop accent palette for meaning — green = healthy/measured/
  reuse, purple = your direct change / links, yellow = derived recompute /
  warnings, red = failure/violation. Missing data stays uncolored dither:
  absence of color is absence of data.
