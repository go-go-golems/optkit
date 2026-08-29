# Changelog

## 2026-08-27

- Initial workspace created


## 2026-08-27

Initial user-story set: 36 stories across 8 themes, written needs-first (corpus knowledge, ground truth, exploration, experiments, judging, decisions, longevity, delegation)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/27/OPTKIT-025--live-ttc-campaign-real-executor-llm-judge-and-operations-tiles/analysis/01-from-fixture-to-live-workflow-gaps-and-the-llm-judge.md — The gap-derived slicing this document supersedes


## 2026-08-27

Design for themes 1-2: 5 presentation types, 5 tiles, 3 workspaces, 2 document formats, verb table, backend projections; question set designed as a proposal-draft isomorph

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/apps/workbench/web/src/pbui/actions.ts — The action registry every new verb contributes to


## 2026-08-27

Design for theme 3: run/hit/chunkPreview types, trail + split tiles, ask lanes, autopsy candidate rows, StageSummary un-coarsening; most of the theme needs no live executor

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/specialistapi/types.go — StageSummary — the lossy projection the design fixes


## 2026-08-27

Theme 9 (Running the work): 12 stories S37-S48 for long-running jobs — build phases, judge batches watched live, multi-turn conversation evaluation, control, dependencies, overnight handoff


## 2026-08-27

Design for theme 9: job/phase/unit/failureGroup/turn types, jobs+job+unit+conversation tiles with kind panels, work projection over BOTH engines, engine-B progress snapshots; multi-turn needs no measurement-chain change

### Related Files

- /home/manuel/go/pkg/mod/github.com/go-go-golems/flowkit@v0.1.1/flow/report.go — StepReport — the per-phase counters the work projection reads


## 2026-08-28

Coverage analysis: audited all 48 stories against the code — 9 covered, 11 partial, 23 designed, 5 not covered; mapped open tickets (025/027/028/029) to stories; named 5 unowned gaps (recommended OPTKIT-030 for Theme 2)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/27/OPTKIT-026--operator-user-stories-for-rag-design-and-optimization/analysis/02-coverage-of-user-stories-by-pbui-optkit-and-rag-ttc.md — The coverage report stored next to the user-stories doc


## 2026-08-28

Intern guide: end-to-end technical orientation of the whole system (substrate, retrieval, optkit control plane, measurement, campaign machinery, optimization/propose loop, five HTTP APIs, pbui workbench, implemented tiles, what is designed-but-not-built, CLI, invariants, glossary) — integrates and extends the 025/027/029 per-ticket intern guides

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/27/OPTKIT-026--operator-user-stories-for-rag-design-and-optimization/design-doc/04-intern-guide-the-rag-ttc-optimization-system-end-to-end.md — The canonical system analysis guide for a new intern


## 2026-08-28

Expanded the needs-first story set from 48 to 50: strengthened S38/S42/S46/S48 for durable progress history, generated-answer supervision, provenance, and exact dependencies; added Theme 10 with S49 end-to-end grounded-answer benchmarking and S50 one durable build-to-benchmark workflow; marked the prior coverage audit as S1-S48 scope.

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/27/OPTKIT-026--operator-user-stories-for-rag-design-and-optimization/analysis/01-user-stories-the-person-who-makes-the-rag-system-good.md — Canonical fifty-story user-needs document
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/27/OPTKIT-026--operator-user-stories-for-rag-design-and-optimization/analysis/02-coverage-of-user-stories-by-pbui-optkit-and-rag-ttc.md — Scope note prevents the original 48-story totals from being misread as current

