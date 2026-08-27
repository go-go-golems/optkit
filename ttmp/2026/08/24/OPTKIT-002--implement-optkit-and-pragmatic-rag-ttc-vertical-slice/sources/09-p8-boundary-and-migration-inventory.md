# P8 Boundary and Migration Inventory

Date: 2026-08-25

## Scope

This inventory classifies the final vertical-slice ownership boundaries and records which old paths are deleted, retained, or deferred. The controlling rule is the accepted OPTKIT-002 decision: replace RagOpt orchestration only after behavioral parity; preserve RAG-domain algorithms, measurement-domain behavior, product policy, historical readers, and product projections.

## Repository classification

| Repository/package | Classification | Owns | Must not own |
|---|---|---|---|
| `optkit/{record,artifact,space,episode,measure,experiment,campaign,budget,scheduler,projection,query,system}` | domain-neutral control/query plane | identity, artifacts, snapshots, trials, trajectories, durable commands/events, leases, budgets, observations, estimates, read-only generic projections | chunks, embeddings, source roles, TTC routes, judge prompts, product providers |
| `ragkit/rag/...` | reusable RAG domain | documents, chunks, representations, embedding/search/rerank interfaces, bundle manifests, retrieval/fusion/hydration, grounded-answer mechanics | campaigns, product routes, source admission, Judgekit instruments |
| `judgekit/{spec,eval,protocol,assessment,judging,audit,calibration,suite}` | reusable measurement domain | constructs, instances, attributed protocols/reports, judge execution, reliability, calibration | TTC source meaning, Optkit campaign control, provider profiles |
| `rag-ttc/pkg/ttc/search` | TTC product data plane | prepared routes, source-role policy, runtime identity, retrieval stages | campaign journal/queue |
| `rag-ttc/pkg/ttc/customerapp` | TTC product application | direct customer turn, answer contract, evidence admission, redacted trajectory | HTTP/WebSocket orchestration |
| `rag-ttc/pkg/ttc/optkitcampaign` | product integration | typed TTC snapshots/cases and durable retrieval campaign composition | generic RAG algorithms or Optkit internals |
| `rag-ttc/pkg/ttc/judgeinstrument` | product integration | sealed-answer projection, TTC prompts/contract, Judgekit-to-Optkit observation mapping | provider credentials or campaign mutation |
| retained `ragopt/pkg/runstore` and `ragopt/pkg/review` callers | legacy artifact/query/report support | historical run readers, blinded review projections, retained reports | new campaign orchestration |
| retained `ragopt/pkg/candidate` and `ragopt/pkg/eval` I5 caller | legacy orchestration pending parity | frozen I5 answer-generation/judge proof cycle | expansion into new workflows |

## Dependency guards added

### Optkit

`optkit/internal/boundary/boundary_test.go` scans every Optkit package's direct, test, and external-test imports and rejects:

- Coinvault;
- Judgekit;
- RagKit;
- RagOpt;
- RAG-TTC.

Product integrations may import Optkit; the reverse direction is forbidden.

### RagKit

`ragkit/boundary_test.go` now applies a second all-package guard that rejects imports of:

- Coinvault;
- Judgekit;
- Optkit;
- RagOpt;
- RAG-TTC.

The existing core guard still allows Geppetto only below `rag/provider/...` while rejecting UI/CLI/application frameworks from core.

### Judgekit

`judgekit/boundary_test.go` now additionally rejects Optkit and RAG-TTC imports from core. It already rejected RagKit, RagOpt, Coinvault, application frameworks, and provider SDKs.

## Deleted superseded RAG-TTC path

The following outer customer tool-loop seam had no production caller after P5 switched provider serving to `customerapp.EngineAdapter`:

- `realruntime.ComposerOptions.ToolRegistryFactory`;
- `realruntime.ResolverOptions.ToolRegistryFactory`;
- `Composer.toolRegistryFactory`;
- the outer `enginebuilder.Builder.Registry` branch;
- `webchatcmd.buildProviderToolRegistry`;
- tests that exercised only that dead factory.

The production path is now singular:

```text
provider engine
  -> realruntime ApplicationEngineFactory
  -> ragsearch.Handle.NewApplicationEngine
  -> customerapp.EngineAdapter
  -> customerapp.Service.RunTurn
  -> per-turn ragsearch.NewSessionRegistry
```

This deletion removed 78 lines and did not delete `widgetintent`, `statusintent`, `ttcwidgets`, historical UI entities, or product projections. Those product packages have separate migration/design work and are not experiment orchestration.

## RagOpt active-call inventory

### RAG-TTC: orchestration retained pending parity

`cmd/rag-ttc/cmds/tooleval/ragopt.go` imports:

- `ragopt/pkg/candidate`;
- `ragopt/pkg/eval`.

The `tool-eval optimize` command remains the only active RAG-TTC candidate/evaluation orchestration caller. It is not behaviorally equivalent to P6/P7 yet. It additionally owns:

- frozen candidate asset loading and locked source/profile/index digests;
- feedback versus validation split selection;
- full answer generation over the product runtime;
- per-cell embedding, generation, and judge budgets;
- judge execution in the same legacy run;
- gate-policy loading;
- native session/outcome artifacts;
- resumable legacy run-directory layout;
- retained promotion-run reproducibility.

P6 proves durable fixed-arm retrieval campaigns. P7 proves historical answer remeasurement. Neither yet proves the complete I5 candidate/gate/report behavior. Deleting this command now would violate the accepted parity-first decision and the retained-promotion evidence recorded by RAG-TTC-CLEANUP-001.

### RAG-TTC: report/query dependencies retained

These remain active readers or product projections rather than duplicate new orchestration:

- `ragopt/pkg/runstore` in knowledge build, answer-quality, connected reports, chunk comparison, TUI run browser, and example 06;
- `ragopt/pkg/review` in answer-quality aggregation and TTC blinded-review projection.

They stay until Optkit has equivalent historical readers/projectors and callers have switched.

### Coinvault: active RagOpt product workflow

Coinvault still imports candidate, evaluation, gate, policy, compare, report, and runstore packages in its knowledge RagOpt commands and tests. The RagOpt repository therefore cannot be archived or have those packages deleted during this TTC-only vertical slice.

## Migration backlog and required parity gates

Before removing the remaining RAG-TTC I5 RagOpt command:

1. Represent locked candidate assets as typed Optkit snapshots/patches without losing digest checks.
2. Execute the canonical direct customer answer service as an Optkit system over the same feedback and validation cases.
3. Preserve exact repeat/cell identity, provider budgets, usage, cancellation, and resume behavior.
4. Apply Judgekit through historical measurement epochs and reproduce the current answer-level constructs.
5. Implement equivalent hard constraints, paired estimates, promotion/gate decision evidence, and report output.
6. Provide a read-only importer or retained reader for historical RagOpt run directories.
7. Run the frozen I5 fixtures through both paths and compare candidate identity, cell outcomes, metrics, failures, decisions, and artifacts.
8. Switch the CLI and all readers.
9. Delete the old command in a separate commit only after the parity report passes.

Before removing `ragopt/pkg/runstore` or `pkg/review`:

1. enumerate every historical reader and artifact schema;
2. provide Optkit projectors/importers with fixture parity;
3. switch RAG-TTC and Coinvault callers;
4. preserve retained promotion artifacts;
5. remove the dependency only after both product repositories pass.

## Workspace/module state

RAG-TTC currently requires unpublished local modules:

```text
github.com/go-go-golems/judgekit v0.0.0
github.com/go-go-golems/optkit v0.0.0
```

The parent workspace contains version-specific replacements to sibling checkouts. Isolated `GOWORK=off` RAG-TTC tidy/release hooks remain deferred until both modules are published or pinned. RagKit and RagOpt remain normal published requirements.

## Evidence commands

```bash
rg -n 'github.com/(go-go-golems/(ragkit|judgekit|ragopt)|the-tree-center/rag-ttc)' optkit --glob '*.go' --glob go.mod
rg -n 'github.com/(go-go-golems/(optkit|judgekit|ragopt)|the-tree-center/rag-ttc)' ragkit --glob '*.go' --glob go.mod
rg -n 'github.com/go-go-golems/ragopt' rag-ttc --glob '*.go' --glob go.mod
rg -n 'github.com/go-go-golems/ragopt' coinvault --glob '*.go' --glob go.mod
rg -n 'ToolRegistryFactory|buildProviderToolRegistry' rag-ttc --glob '*.go'
```

Expected final properties:

- the first, second, and final dead-path searches return no forbidden/deleted active imports;
- RagOpt searches return only the explicitly retained callers described above;
- boundary tests pass in Optkit, RagKit, and Judgekit;
- product integrations in RAG-TTC may import both the domain libraries and Optkit.
