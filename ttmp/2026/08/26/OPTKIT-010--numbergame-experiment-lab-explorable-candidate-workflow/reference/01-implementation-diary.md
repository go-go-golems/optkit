---
Title: Implementation Diary
Ticket: OPTKIT-010
Status: active
Topics:
    - implementation
    - optkit
    - ui
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Diary of the Experiment Lab — an explorable in-browser numbergame campaign with exact Go parity, compared against the sealed store.
LastUpdated: 2026-08-26T10:18:11.644405681-04:00
WhatFor: ""
WhenToUse: ""
---


# Diary

## Goal

Record the construction of the Experiment Lab: a browser screen where a person runs a complete optkit campaign by hand — propose, seal, execute, estimate, decide — with exact parity to the Go numbergame example, grounded against a bundle exported from a real sealed store.

## Step 1: Exporter, simulation, and the explorable lab screen

The user asked to generate numbergame data and use it to explore UI concepts — then sharpened the goal: "not just a way to visualize results... a way to play with them... experience the experiments for myself." Numbergame is the right vehicle because it is the only optkit consumer of `space.Candidate` (hypothesis, targets, risks, decision policy) and its mechanics are pure arithmetic, so the browser can execute episodes with the exact recorded semantics.

### Prompt Context

**User prompt (verbatim):** "Can we generate the data for numbergame? and then maybe use that as a thing to visualize as well, with maybe a more effective UI / explore other UI concepts?" — followed mid-turn by: "like I want not just a way to visualize results and such, but a way to play with them to understand how the different cases and arms and campaigns work, how the individual steps work, experience the experiments for myself."

**Assistant interpretation:** Run the numbergame demo into a store, export its data, and build an interactive lab where the user performs each campaign step themselves, rather than a passive results viewer.

**Inferred user intent:** Build intuition for optkit's concepts (arms, cases, episodes, trajectories, paired estimates, decision policies, the journal) by doing them, with the sealed record as ground truth.

**Commit (code):** optkit `99baa92` — "Add numbergame-demo exporter for the Experiment Lab UI"; rag-ttc `53e1251dd` — "Experiment Lab: run a numbergame campaign by hand in the browser"

### What I did

- `optkit/cmd/numbergame-demo`: runs `RunDemo` with a deterministic sequence clock, reopens the store, and reconstructs a `numbergame.lab-bundle/v1` JSON entirely from the journal and artifacts: candidate proposal (hypothesis/targets/risks/proposer/strategy/lineage), baseline and challenger configs (loaded from snapshot records), cases, all eight episode measurements (judge score, absolute error, intervention status, judge rationale), estimate, decision, budget, and the 81-event journal timeline. Output goes to the specialist app's `public/numbergame-demo.json`.
- `src/lab/sim.ts`: TypeScript port of the mechanics — the 64-bit xorshift noise hash (BigInt), episode execution with the five trajectory events, judge score 1/(1+|error|), `pairedMean`, and the lexicographic decision with Go's exact reason strings.
- `src/screens/LabScreen.tsx` at `/lab/numbergame`: seven stations (system under test → propose with live per-case preview → trial plan grid with budget → run episodes one click at a time with trajectory and a growing sandbox journal → paired estimate with slope graph → decision gates → sealed-record comparison). Sealing the proposal freezes the controls: changing your mind means proposing a new candidate — the immutability of a trial plan, experienced. Linked from the campaign entry screen.
- 9 parity tests pin the simulation against the exported bundle (every episode score/error, the estimate, the decision object including its reason string) plus decision-order and noise-range properties. 45/45 total.

### What worked

- Live walkthrough: slider to 3, seal, run all → sandbox Δ 0.710417 = recorded Δ 0.710417, both `eligible`, and the screen says "reproduced" — the deterministic-store thesis demonstrated interactively (screenshot 18).
- The bundle exposes the full candidate metadata the rag-ttc manifests lack: hypothesis, declared regression risks, proposer, strategy — displayed verbatim from the sealed record.

### What didn't work

- First parity run failed on the estimate: exact arithmetic gives 0.71041666…, the store says 0.71041675. Cause: the Go judge records its score via FormatFloat(…, 'f', 6) and the paired mean is computed from those six-decimal recorded values. The simulation now quantizes scores to recorded precision before differencing — instrument precision is part of the measurement, and the failing test taught it.
- Pre-existing breakage: optkit's `store/sqlite/rows.go` had an uncommitted stray "k" line (accidental keystroke) that broke compilation; reverted that one uncommitted hunk.
- The exporter initially referenced a nonexistent `space.LoadSnapshotValue`; the store-faithful path decodes the baseline snapshot record from the CampaignCreated payload and the challenger from SnapshotMaterialized, then `space.LoadSnapshot` with the config codec.

### What was tricky to build

- Faithfulness boundaries: the sandbox intervention gate is always satisfied by construction, so the UI says so rather than pretending to probe a trajectory; noise modes get an honest caption that a preview shows one draw, which is why paired trials exist.

### What warrants a second pair of eyes

- The committed `public/numbergame-demo.json` snapshots one run; regenerating it after numbergame changes requires re-running the exporter (parity tests will fail loudly if the two drift — by design).
- Slope-graph bottom padding now scales with plotted pairs; verify no clipping at larger case counts.

### What should be done in the future

- The lab pattern (do it by hand, compare with the sealed record) generalizes to rag-ttc: a retrieval lab where the user assembles a fused ranking by hand before seeing the pipeline's, once candidate metadata lands in rag-ttc manifests.

### Code review instructions

- optkit: `git show 99baa92` (`cmd/numbergame-demo/main.go`); run `GOWORK=off go run ./cmd/numbergame-demo -out /tmp/bundle.json`.
- rag-ttc: `git show 53e1251dd` — start at `src/lab/sim.ts` (parity core), then `LabScreen.tsx`; validate `pnpm test`; live at http://127.0.0.1:5197/lab/numbergame. Evidence: `various/screenshots/18-experiment-lab.png`.
