# Changelog

## 2026-08-25

- Initial workspace created


## 2026-08-25

Scenarios/questions design doc written; UI restyled flat with hyperslop accents, hashes demoted to trays/toggles, verdict-first comparison; 28/28 tests green, live-validated with screenshots (commit 1e77d16cd)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/25/OPTKIT-009--humanize-the-specialist-ui-scenarios-color-accents-hash-demotion/design-doc/01-specialist-work-scenarios-and-questions.md — Product roadmap grounded in usage scenarios


## 2026-08-25

Step 2: scenarios doc rephrased around failure-first IR work; Tufte graphics (slope/chunk-grid/stage-flow/meters); description fields added to manifest/arm/case through campaign spec and API; manifest rewritten with authored prose; ledes on all panels (commits f5e5a7e29, baf6ce93f)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/optkitcampaign/system.go — RetrievalCase.Description joins the persisted case contract


## 2026-08-25

Step 3: per-layer widget registry — retrieval value diff (limit 2 was 1) and ranked-evidence output widget (scores, fusion contributions, chunk text); ArmSummary.Config surfaced; preview cap 16KiB; no data recreation needed (commit 4ddd59f03)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/apps/specialist/web/src/layerwidgets/retrieval.tsx — Retrieval diff and ranked-output widgets


## 2026-08-25

Step 4: per-stage candidate recording (scores, channels, fusion contributions) in RetrievalStage; pipeline shows query, case prose, stage glossary, and candidate tables; store recreated (commit 20e8d4621)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/search/service.go — StageCandidate recording in stage constructors


## 2026-08-25

Step 5: chunk catalog recorded with every run (content of all touched chunks incl. filtered ones); pipeline shows readable chunk panel and What-it-says column; store recreated (commit 4ff578c1c)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/search/search.go — ChunkCatalog collection in RunRoute


## 2026-08-25

Step 6: match lineage — representation_id recorded per hit, representations registered into chunk catalog, UI shows 'matched on' text; query confirmed recorded and verbatim to both channels (commit 340b7b02a)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/search/semantic_fixture.go — Fixture registers shipped representations

