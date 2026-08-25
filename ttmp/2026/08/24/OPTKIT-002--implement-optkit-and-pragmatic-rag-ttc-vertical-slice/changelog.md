# Changelog

## 2026-08-24

- Initial workspace created


## 2026-08-24

Created detailed P0-P8 implementation plan; retained RagKit as a separate RAG-domain library; printed overall plan slip

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/design-doc/01-phased-implementation-plan.md — Implementation tasks and acceptance gates


## 2026-08-24

P0 complete: imported archive revision 1786d1d, normalized repository plumbing, and passed full/race/no-CGO/lint/demo/verification gates (commits ea51f8f, b49aece)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/02-p0-validation.txt — Full successful P0 validation transcript


## 2026-08-24

P1 complete: froze byte-identical Coinvault/RAG-TTC semantic fixture and passed focused/full product tests (commits b8aaf41d9, e3090be05)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/rag-semantic-fixture-v1.json — Canonical cross-product RAG fixture v1


## 2026-08-24

P2 complete: extracted one direct RAG-TTC retrieval service below Geppetto and preserved semantic fixture parity (commit ea6c629be)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/pkg/ttc/search/service.go — Canonical channel, collapse, fusion, augmentation, hydration, and source-verification path


## 2026-08-24

Added P2.5 read-only scientific query plane: CLI/agents own writes; UI owns search, navigation, replay, provenance, and visualization

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/design-doc/01-phased-implementation-plan.md — Read-only UI phase and accepted boundary decision


## 2026-08-24

Completed P2.5 through OPTKIT-003 and P3 attributable policy-safe retrieval routes (RAG-TTC commit d4c5adab4); full validation passes.

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/../rag-ttc/pkg/ttc/search/service.go — Canonical P3 runtime behavior
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/04-p3-validation.txt — Passing P3 command transcript


## 2026-08-24

Completed P4 deterministic stage-aware retrieval evaluation with treatment verification and answer-quality runner artifacts (RAG-TTC commit d7701685d).

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/../rag-ttc/pkg/ttc/retrievaleval/evaluate.go — Canonical P4 evaluator
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/05-p4-validation.txt — Passing P4 command transcript


## 2026-08-24

Completed P5 canonical direct customer application service and switched provider transport composition (RAG-TTC commit 8853613a4).

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/../rag-ttc/pkg/ttc/customerapp/service.go — Canonical P5 domain turn
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/06-p5-validation.txt — Passing P5 command transcript


## 2026-08-24

P6 checkpoint: added domain-neutral executable system registry (Optkit commit b8e233e86); real RAG-TTC campaign remains in progress pending reproducible cross-module integration.

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/system/registry.go — P6 foundation checkpoint


## 2026-08-25

P6 complete: durable six-episode TTC retrieval campaign, stage artifacts, deterministic observations/estimates, restart reconciliation, and Glazed run/inspect commands (Optkit b8e233e86; RAG-TTC daacadaad and 2bb1c21af).

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/../rag-ttc/pkg/ttc/optkitcampaign/campaign.go — Durable campaign implementation
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/07-p6-validation.txt — P6 acceptance evidence


## 2026-08-25

P7 complete: attributed Judgekit reports now remeasure sealed TTC answer trajectories under distinct Optkit epochs without rerunning retrieval or generation; cache-bypass, failure, missing-output, and old-observation isolation tests pass.

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/../judgekit/assessment/provenance.go — Required report attribution
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/../rag-ttc/pkg/ttc/judgeinstrument/instrument.go — Historical measurement adapter
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/08-p7-validation.txt — P7_VALIDATION=PASS


## 2026-08-25

P8 complete: guarded Optkit/RagKit/Judgekit ownership, deleted the superseded outer customer tool loop, retained only non-parity RagOpt orchestration and active historical readers, and passed cross-repository validation.

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/../rag-ttc/internal/customer/realruntime/composer.go — Single canonical customer application path
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/internal/boundary/boundary_test.go — Optkit ownership guard
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/09-p8-boundary-and-migration-inventory.md — Retention and migration evidence
- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/sources/10-p8-validation.txt — P8_VALIDATION=PASS


## 2026-08-25

Printed and recorded the P7 and P8 completion slips after both acceptance gates passed.

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/reference/01-implementation-diary.md — P7/P8 slip render evidence

