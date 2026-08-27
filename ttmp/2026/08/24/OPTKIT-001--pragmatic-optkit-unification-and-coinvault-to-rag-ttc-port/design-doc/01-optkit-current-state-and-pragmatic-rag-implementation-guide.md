---
Title: Optkit Current State and Pragmatic RAG Implementation Guide
Ticket: OPTKIT-001
Status: active
Topics:
    - optkit
    - rag
    - rag-ttc
    - coinvault
    - judgekit
    - architecture
    - migration
    - local-development
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://coinvault/internal/knowledge/eval.go
      Note: Reference evaluation modes and failure custody
    - Path: ws://coinvault/internal/knowledge/service.go
      Note: Behavioral reference for verified hybrid retrieval and policy ordering
    - Path: ws://judgekit/ttmp/2026/08/17/JUDGEKIT-001--design-and-implement-judgekit/design-doc/03-lightweight-research-guarantees-after-merge.md
      Note: Accepted local research attribution policy
    - Path: ws://rag-ttc/internal/customer/ragsearch/ragsearch.go
      Note: Current customer composition and verified-bundle lifecycle
    - Path: ws://rag-ttc/pkg/ttc/search/search.go
      Note: Current TTC hybrid tool and extraction target
    - Path: ws://rag-ttc/pkg/ttc/toolconfig/types.go
      Note: Current strict serving and routing configuration contract
    - Path: ws://sources/optkit-clean-slate-architecture-and-porting-guide.md
      Note: Original clean-slate architecture and product-port design
ExternalSources: []
Summary: Evidence-backed assessment of the Optkit design and implementation, plus a pragmatic intern-ready plan to transfer Coinvault's proven RAG architecture into RAG-TTC before completing framework consolidation.
LastUpdated: 2026-08-24T22:30:00-04:00
WhatFor: Orient a new engineer, state what exists versus what is only designed, and provide file-level APIs, pseudocode, tests, and deletion gates for the next implementation sequence.
WhenToUse: Before importing the supplied Optkit archive, changing RAG-TTC retrieval or evaluation, integrating Judgekit, or deleting RagOpt/RagKit paths.
---


# Optkit Current State and Pragmatic RAG Implementation Guide

## 1. Executive summary

The clean-slate architecture is conceptually strong and the supplied implementation is real, tested work. The checked-out `optkit/` repository, however, does **not** contain that implementation. It is still the Go template with module path `github.com/go-go-golems/XXX`. The supplied ZIP contains a separate 7,581-line, 59-Go-file `github.com/go-go-golems/optkit` module whose full test suite passes. The first Optkit implementation task is therefore an import and provenance task, not a rewrite.

The supplied implementation completes a durable local vertical slice: canonical records, content-addressed artifacts, typed snapshots and patches, episode trajectories, measurement epochs, complete-block trials, a hash-chained campaign journal, SQLite leases, budgets, projections, and a restartable Numbergame campaign. It does **not** yet contain RAG, a real judge, application-service registration, web UI, exposure control, adaptive search, Coinvault, or RAG-TTC ports.

The product repositories are farther along than Optkit in their domains:

- RagKit has a mature RAG data plane and all tests pass.
- RagOpt has useful candidate, paired-evaluation, gate, report, review, and run-store semantics, but its orchestration role should eventually disappear.
- Judgekit has a working provider-neutral judge implementation and all tests pass. Its latest ticket explicitly chooses trusted local research attribution over adversarial custody. The proposed lightweight cleanup is not fully implemented yet.
- Coinvault has the strongest production-oriented RAG service: verified immutable bundles, hybrid retrieval, pre-fusion authorization, optional reranking, bounded evidence admission, semantic runtime identities, stage diagnostics, and strict evaluation modes.
- RAG-TTC already has a verified hybrid index, strict tool configuration, turn-scoped evidence, a customer runtime below HTTP, answer-quality experiments, and rich product/UI code. Its RAG path is not a blank slate. The missing piece is convergence into one canonical, observable service with Coinvault's provenance and diagnostic discipline.

The pragmatic recommendation is:

> **Use Coinvault as the behavioral reference, make RAG-TTC's current RagKit-backed customer search a proper canonical RAG service, and register that service as Optkit's first real product slice. Do not block product progress on bulk-consolidating every kit first.**

This produces two coordinated tracks:

1. **Foundation track:** land the supplied Optkit archive unchanged enough to preserve its green baseline, then add only capabilities required by the RAG-TTC vertical slice.
2. **Product track:** refactor RAG-TTC around a direct, product-owned retrieval/answer service, preserving its current bundle and tool contracts while adding Coinvault-style identity, stage traces, policy boundaries, deterministic evaluation, and answer constraints.

After the vertical slice proves the ownership boundaries, move generic RAG primitives from RagKit into Optkit in reviewed batches, switch the last callers, and delete old imports. This is a staged cutover, not a compatibility architecture.

## 2. The decision in one diagram

```text
                    CURRENT SOURCE MATERIAL

 RagKit primitives       Coinvault lessons       RAG-TTC product
 (RAG data plane)        (proper RAG behavior)   (customer/admin apps)
        |                         |                       |
        +-------------------------+-----------------------+
                                  |
                                  v
                    RAG-TTC canonical RAG service
              verified bundle -> retrieve -> admit -> answer
                                  |
                       deterministic semantic fixtures
                                  |
                                  v
                    Optkit first real vertical slice
                 snapshots -> episodes -> observations -> trial
                                  |
                                  v
              move generic primitives, switch callers, delete kits
```

The immediate output is a better RAG-TTC system. The long-term output is the unified framework. The order matters because a real product slice tests abstractions better than an empty package consolidation.

## 3. Scope and non-scope

### 3.1 In scope

This guide assesses:

- `sources/optkit-clean-slate-architecture-and-porting-guide.md`;
- `sources/implementation-diary.md`;
- `sources/optkit-implementation-source.zip`;
- the checked-out `optkit`, `ragkit`, `ragopt`, `judgekit`, `coinvault`, and `rag-ttc` repositories;
- Judgekit's latest trust-model decisions;
- Coinvault RAG behavior worth transferring;
- RAG-TTC's current RAG path and concrete gaps;
- a phased implementation and deletion plan.

### 3.2 Explicit non-scope

This ticket does not itself:

- import the Optkit ZIP into the repository;
- move RagKit source;
- refactor Coinvault or RAG-TTC code;
- define production deployment infrastructure;
- introduce signatures, key management, attestations, or hostile-storage defenses;
- promise API compatibility for any v0 kit.

The guide is an implementation contract for subsequent focused tickets.

## 4. Evidence and verification performed

The investigation ran these checks successfully:

```text
supplied Optkit archive:  go test ./... -count=1                 PASS
checked-out Optkit:       GOWORK=off go test ./... -count=1      PASS (template only)
Judgekit:                 GOWORK=off go test ./... -count=1      PASS
RagKit:                   GOWORK=off go test ./... -count=1      PASS
RagOpt:                   GOWORK=off go test ./... -count=1      PASS
Coinvault focused RAG:    knowledge/knowledgebuild/evalchat      PASS
RAG-TTC focused RAG:      search/config/index/knowledge/answer   PASS
```

The generated inventory is in `sources/01-repository-inventory.md`. At the inspected revisions it records:

| Repository | Go files | Go lines | Go packages | Role today |
|---|---:|---:|---:|---|
| checked-out Optkit | 4 | 20 | 3 | unnormalized template plus this ticket |
| RagKit | 187 | 23,228 | 27 | reusable RAG data plane |
| RagOpt | 55 | 9,719 | 14 | paired optimization orchestration |
| Judgekit | 61 | 7,041 | 15 | provider-neutral evaluator library |
| Coinvault | 291 | 54,665 | 39 | production product and strongest RAG reference |
| RAG-TTC | 392 | 67,643 | 68 | TTC product, laboratory, and current target |

These counts are orientation, not quality scores.

## 5. Vocabulary an intern must know

### 5.1 Configuration, snapshot, patch, and candidate

A **configuration** is the complete set of semantic values controlling a system. A **snapshot** is the immutable identity of one complete configuration. A **patch** changes named variables against one exact base snapshot. A **candidate** connects parent, patch, child, hypothesis, and provenance.

```text
configuration C0 --materialize--> snapshot S0
       S0 + patch P1 -----------> snapshot S1
       S0 + P1 + hypothesis ----> candidate C1
```

Do not use “candidate” to mean a mutable directory. Do not use a partial map as if it were a complete configuration.

### 5.2 System, preparation, episode, attempt, and trajectory

A **system** maps an input and configuration to observable behavior. **Preparation** resolves environment-owned dependencies such as providers and index bundles. An **episode** is one semantic execution. An **attempt** is one retry within it. A **trajectory** is the ordered event evidence produced during execution.

For RAG retrieval:

```text
system:      ragttc.retrieval
config:      chunking + channels + fusion + reranker + filters + top-k
input:       one benchmark query
preparation: open verified bundle and provider adapters
trajectory:  query -> channels -> collapse -> authorize -> fuse -> rerank -> hydrate
result:      ranked evidence artifact or explicit failure
```

### 5.3 Observation, instrument, epoch, estimand, and estimate

An **instrument** measures an episode or artifact. It emits immutable **observations** under one **measurement epoch**, which identifies the construct, protocol, implementation, and calibration context. An **estimand** says what population quantity is being asked. An **estimate** is the computed value.

Changing the target resolver changes the measurement epoch or estimand; it does not change the retrieval candidate.

### 5.4 Trial, campaign, search, and promotion

A **trial** declares arms, cases, repeats, protocol, instruments, and estimands. A **campaign** durably coordinates one or more trials and candidate decisions. A **search strategy** proposes candidates. A **promotion policy** decides whether evidence authorizes product adoption.

Search and promotion must remain separate. An optimizer may learn from development evidence; promotion uses explicitly protected evidence and may require a human.

## 6. What the supplied Optkit implementation actually provides

### 6.1 Repository mismatch

The checked-out repository's `go.mod` still says:

```go
module github.com/go-go-golems/XXX
```

The archive says:

```go
module github.com/go-go-golems/optkit
```

This is the single most important status distinction. “Implemented in the supplied source” does not mean “available on the current branch.”

### 6.2 Implemented package map

```text
record      canonical JSON, digests, typed IDs, schema registry
artifact    immutable bytes; memory and filesystem CAS
space       domains, codecs, lenses, snapshots, patches, candidates
episode     events, spans, trajectory writer/reader, result/failure/usage
measure     epochs and immutable observations
experiment  datasets, complete-block expansion, paired mean
campaign    commands, events, reducer, journal interfaces
scheduler   durable work vocabulary and lease interface
budget      limits, reservations, actual use, overage custody
store/sqlite journal, queue, and budget persistence
projection  rebuildable campaign overview
local       SQLite + filesystem-CAS composition root
numbergame  complete restartable proof system
cmd/optkit  demo, inspect, campaign verify, artifact verify
```

Key API anchors in the archive are:

- `space.MaterializeSnapshot` at `space/snapshot.go:24`;
- `space.Set` and `PatchBuilder.Build` at `space/patch.go:89,104`;
- `episode.Writer.Emit` and `Seal` at `episode/writer.go:49,130`;
- `experiment.NewCompleteBlockTrial` and `Expand` at `experiment/trial.go:66,115`;
- `experiment.PairedMean` at `experiment/estimate.go:26`;
- `campaign.Controller.Handle` at `campaign/controller.go:41`;
- `campaign.MaterializeEvent` at `campaign/event.go:74`;
- `scheduler.NewWorkItem` at `scheduler/types.go:39`.

### 6.3 Implemented end-to-end flow

The Numbergame example proves this sequence:

```text
materialize baseline snapshot
        |
build typed patch and challenger snapshot
        |
create candidate and complete-block trial
        |
expand deterministic episode IDs
        |
reserve budgets and enqueue work
        |
lease -> execute -> emit trajectory -> seal result
        |
record observations and terminal events
        |
restart SQLite process boundary
        |
rebuild paired rows -> estimate -> decision
        |
restart again -> verify journal -> rebuild overview
```

This is substantive infrastructure. It should be imported as a baseline rather than reconstructed from the 8,841-line design guide.

### 6.4 Missing relative to the architecture

| Architecture area | Status in archive | Consequence |
|---|---|---|
| `record`, `artifact` | useful v0 implemented | import first |
| typed `space` | useful v0 implemented | enough for retrieval configs |
| system registry/preparation | not generalized | needed for RAG-TTC registration |
| `work`/`flow` consolidation | absent | continue using existing FlowKit temporarily |
| trajectory redaction | absent | add before sensitive product traces |
| measurement instruments/DAG | only epoch and observation records | deterministic RAG instruments are next |
| Judgekit integration | fake judge only | defer until retrieval evidence is correct |
| multi-arm estimators | trial expansion supports arms; estimator is paired mean only | add per-arm and paired-delta summaries |
| policy expression tree | absent | use a small product policy initially |
| reconciliation | known queue/journal gap in diary | solve before provider-scale campaigns |
| server/SSE/web UI | absent | not required for first CLI slice |
| exposure/authorization ledger | absent | do not expose hidden datasets to search |
| RAG packages | absent | use current RagKit for first product slice |
| adaptive search | absent | manual fixed arms first |
| product ports | absent | RAG-TTC retrieval is first real port |

### 6.5 Assessment

The archive is approximately “Slices A-C, part of E-F, and a strong local Slice G,” not the whole clean-slate program. It is ready to host one deterministic RAG retrieval campaign after adding a system seam and RAG observations. It is not ready to replace all four kits at once.

## 7. Judgekit's accepted local-development trust model

### 7.1 The three guarantee levels

The latest `JUDGEKIT-001` documents distinguish:

1. **Structural validity:** required fields, ranges, IDs, and references are valid.
2. **Research attribution:** the instance, contract, protocol, prompt, model, cache mode, and result population agree.
3. **Adversarial integrity:** hostile mutation, storage, plugins, and cross-organization exchange are protected by immutable verified values, signatures, custody services, and independent verification.

The accepted target is Level 2. Level 3 is explicitly deferred.

The controlling rule is:

> Reject mistakes that change the meaning or attribution of a measurement. Do not defend every in-memory value against mutation by trusted research code.

### 7.2 What to keep

For local RAG development, retain checks that prevent false conclusions:

- strict document/config decoding;
- current-content identity at execution boundaries;
- contract/protocol compatibility;
- expected versus observed model identity;
- explicit prompt versions and rendered-prompt digests;
- cache bypass for repeatability probes;
- durable report/gold-set digest verification when consumed;
- evidence provenance and missing-output visibility;
- exact integer semantics in canonical data;
- explicit measurement epochs.

These are research-correctness mechanisms, not adversarial hardening.

### 7.3 What not to add now

Do not add without a concrete cross-trust consumer:

- signed reports;
- key rotation or revocation;
- Merkle execution manifests;
- private immutable typestate wrappers for every value;
- recursively verified object graphs at every call;
- append-only remote custody services;
- provider attestations;
- hostile plugin sandboxes;
- long-term canonicalization registries for hypothetical archival consumers.

### 7.4 Current Judgekit status versus the proposal

The lightweight guide is a post-merge design, not a complete description of current code. Concrete evidence:

- `eval.Instance` still has a caller-maintained `Digest` field at `judgekit/eval/instance.go:17`.
- `ClaimProtocol.ExtractPrompt` still receives a complete `eval.Instance` at `judgekit/judging/claimjudge.go:20`, not a restricted extraction input.
- `ClaimJudge.validateFor` recomputes and verifies contract/protocol digests before generation.
- cache `use` and `bypass` modes are implemented at `judgekit/judging/interfaces.go`.
- the report still seals a durable digest.

Therefore Judgekit is green and usable, but the proposed simplification and evidence-hidden extraction type remain follow-up work. Do not assume they have landed merely because the design ticket accepts them.

### 7.5 How this applies to Optkit

The Optkit archive already has a hash-chained local journal and CAS verification. Keep them: they are implemented, tested, and useful for accidental corruption and replay. Do not expand them into signatures or a custody platform.

For RAG-TTC, compute runtime identities when opening a bundle or beginning an episode. Store identities on terminal episode/measurement artifacts. Do not require every internal helper to reverify the entire object graph.

## 8. Coinvault's RAG system: the reference behavior

Coinvault is the behavioral source, not a package to copy wholesale. Its product-specific authorization, source roles, tools, prompts, and eval cases stay product-owned.

### 8.1 Build and publication boundary

`coinvault/internal/knowledgebuild/manifest.go` defines a strict reviewable corpus manifest:

- admitted source connectors;
- source database identity without credentials;
- product/category/SQL-doc roles;
- chunking limits and overlap;
- optional embedding provider/model/dimensions/batching;
- optional boilerplate policy.

`Manifest.Validate` rejects impossible chunking, missing source choices, incomplete embedding identities, and invalid boilerplate parameters. The build writes an immutable RagKit index bundle. Serving opens the bundle through `knowledge.Open` at `internal/knowledge/service.go:89`, verifies the vector provider identity, and opens content read-only.

The lesson is not “every product needs Coinvault's manifest fields.” The lesson is:

```text
human-reviewable source admission
        +
complete build identity
        +
immutable publication
        +
verified read-only open
```

### 8.2 Canonical retrieval service

`knowledge.Service.Search` at `internal/knowledge/service.go:456` owns the semantic retrieval pipeline:

```text
query
  -> optional reviewed decomposition/intent
  -> lexical ranking -> representation-to-chunk collapse -> authorize
  -> vector ranking  -> representation-to-chunk collapse -> authorize
  -> weighted reciprocal-rank fusion
  -> optional external rerank of an authorized hydrated pool
  -> blend fused and reranked order
  -> final authorization check
  -> hydrate verified chunks/documents
  -> deterministic top-k
```

Important properties:

- model input does not contain access scopes;
- authorization occurs before cross-channel fusion and before external reranking;
- representation hits collapse to source chunks before evidence use;
- reranker failure degrades to the fused order rather than failing retrieval;
- reranker text includes title/heading context;
- final results are deterministic given channel rankings and configuration;
- source documents are read from verified content storage.

### 8.3 Runtime identity

Coinvault's search tool output records:

- immutable bundle ID;
- query transform ID;
- retrieval policy ID;
- evidence ledger policy ID;
- reranker configured/applied/runtime identity;
- tool-description identity;
- effective limit and its source;
- comparison-plan selection.

The relevant functions are:

- `Service.QueryTransformID` at `service.go:228`;
- `Service.RetrievalPolicyID` at `service.go:250`;
- `NewToolEntry` at `tool.go:100`;
- `runSearch` at `tool.go:136`.

This is the most reusable Coinvault lesson. A search result should explain which semantic pipeline produced it without requiring ambient flags or log archaeology.

### 8.4 Evidence admission

`EvidenceLedger` at `internal/knowledge/evidence.go:25`:

- deduplicates by chunk;
- assigns stable `E1..En` labels per answer run;
- bounds item count and total runes;
- reuses the same label for repeated evidence;
- exposes a policy identity.

This makes citations a runtime contract rather than prompt convention alone.

### 8.5 Evaluation and diagnosis

Coinvault separates three eval modes in `internal/knowledge/eval.go`:

- positive required evidence groups;
- authorization-negative expected absence;
- judge-only answer cases.

`RunEval` at line 237 preserves per-question failures as misses rather than dropping failed rows. `RunCandidatePoolEval` at `candidate_pool.go:189` records stage rankings and diagnoses where required evidence disappeared:

```text
raw channel
collapsed channel
scope-authorized channel
fused
scope-authorized fused
reranked
returned
admitted
```

This allows an engineer to distinguish “not indexed” from “below budget,” “removed by policy,” “lost in fusion,” or “blocked by evidence admission.”

### 8.6 Known Coinvault boundary still worth improving

Coinvault's controlled evaluator composes the real runtime but drives it through HTTP/WebSocket in `internal/webchat/evalchat/configured_runner.go`. That preserves fidelity but pays transport complexity. The clean-slate guide correctly recommends extracting a direct application service below transport. RAG-TTC should not copy this transport coupling.

## 9. RAG-TTC's current RAG path

### 9.1 What already works

RAG-TTC has a meaningful base:

- `cmd/rag-ttc/cmds/indexes/build.go` builds RagKit immutable index bundles with chunker, representation, embedding, provider-work, cache, budget, and cost metadata.
- `pkg/ttc/toolconfig.Load` at `toolconfig/load.go:17` strictly loads file-backed configuration, resolves safe repository-relative assets, and computes a resolved digest.
- `internal/customer/ragsearch.Open` at `ragsearch.go:37` verifies source documents, opens a verified bundle with exact embedding identity, and constructs a source catalog.
- `pkg/ttc/search.NewSearchTool` at `search.go:107` creates deterministic hybrid retrieval with bounded, turn-scoped evidence.
- `SearchTool.RunRoute` at `search.go:172` performs channel search, collapse, weighted RRF, optional augmentation, hydration, and evidence admission.
- `internal/customer/realruntime.Composer.Compose` at `composer.go:61` creates the provider runtime below the HTTP handler and obtains a fresh tool registry per conversation.
- `ttc_search_results_show` allows the frontend to display only citations admitted by the exact search-tool instance.
- answer-quality, tool-eval, connected-RAG, knowledge extraction, admin review, and customer feedback fixtures already exist.

This is much closer to “proper RAG” than a new implementation.

### 9.2 Current serving flow

```text
webchat command
   |
   +-> load inference profile and tool config
   +-> ragsearch.Open
          -> verify documents
          -> open immutable index bundle
          -> retain search config + source catalog
   |
   +-> realruntime.Resolver
          -> Composer.Compose per conversation
          -> Handle.NewSessionRegistry
                 -> new SearchTool + fresh evidence ledger
                 -> register ttc_search + source-results widget
   |
   +-> Pinocchio/Geppetto tool loop
          -> ttc_search
          -> answer
          -> timeline/widget projections
```

The service boundary exists in pieces, but retrieval semantics are still attached to the model-facing tool object.

### 9.3 Concrete gaps relative to Coinvault

| Concern | Current RAG-TTC | Desired state |
|---|---|---|
| semantic service | `SearchTool` owns retrieval and ledger together | retrieval service separate from tool/session ledger |
| bundle identity in output | not present in `SearchOutput` | include bundle/index identity |
| query transform identity | not present | record exact lexical/vector transform contract |
| retrieval policy identity | only route observation | record stages, top-k, RRF, reranker, filters |
| ledger policy identity | ledger behavior exists but identity absent | expose scope/dedupe/item/rune policy |
| server-owned source policy | source catalog verifies metadata; tool has no allow policy | preparation selects allowed corpus/roles; model cannot widen |
| authorization timing | no pre-fusion policy filter | filter before fusion/external services when policy exists |
| reranking | production tool has route augmentation but no explicit reranker | optional bounded reranker with identity and fallback |
| route configuration wiring | rich intent-routing config exists; customer `ragsearch.Open` does not register those routes | compile checked config into actual service routes |
| stage trace | contribution list and route augmentation only | canonical stage event/trace with candidate identities |
| eval failure custody | large custom runners and runstore artifacts | deterministic service-level suite plus Optkit observations |
| answer contract | output schemas/prompts exist across tool loops | one explicit grounding/citation/abstention result contract |
| framework integration | custom answer-quality and RagOpt paths | one Optkit retrieval trial first |
| Judgekit | no product imports | integrate only after deterministic retrieval evidence is stable |

### 9.4 A critical wiring observation

`toolconfig.IntentRoutingConfig` defines intents, named routes, source roles, and connected-RAG settings. `SearchTool.AddRoute` and `RunRoute` implement route mechanics. At the inspected revision, production `internal/customer/ragsearch` creates only the default tool and does not compile loaded routing config into `AddRoute` calls. The intent-routing configuration is strongly validated and tested as configuration, but it is not the customer serving path's active routing compiler.

Do not describe a checked-in configuration as production behavior until the composition root wires it.

## 10. What “proper RAG” means for this project

A proper RAG system here is not merely “vector search plus an LLM.” It satisfies these contracts:

### 10.1 Corpus contract

- sources are explicitly admitted;
- documents and chunks have stable identity;
- source metadata is verified;
- the index manifest names chunking, representations, embeddings, and corpus;
- publication is immutable and startup verifies it.

### 10.2 Retrieval contract

- lexical and vector stages are independently observable;
- representation hits collapse to source chunks;
- filtering happens before disallowed text crosses an external boundary;
- fusion is deterministic;
- reranking is bounded and attributable;
- final hydration resolves every result to verified source content;
- failures remain visible.

### 10.3 Evidence contract

- evidence is admitted under explicit item/rune limits;
- repeated chunks reuse labels;
- labels are scoped to one conversation/answer episode;
- answer citations resolve to admitted evidence;
- retrieved text is treated as data, not instructions.

### 10.4 Answer contract

- generation receives the exact admitted evidence;
- the output schema identifies answer, citations, and abstention state;
- malformed output becomes an explicit contract failure or safe abstention;
- unsupported claims, missing citations, and unresolved citations are measured separately;
- provider failure is not rewritten as low quality.

### 10.5 Evaluation contract

- cases have explicit roles and modes;
- target-resolver identity is recorded;
- every arm sees the same complete block of cases;
- missing rows are not silently discarded;
- deterministic retrieval metrics precede model judging;
- judge model/protocol/prompt/cache identities are attributable;
- promotion constraints remain separate from measurements.

## 11. Pragmatic target architecture

### 11.1 Near-term ownership

```text
RAG-TTC repository
  pkg/ttc/retrieval/       TTC semantic service and policy
  pkg/ttc/search/          thin model-facing tool adapter + evidence ledger
  internal/customer/ragsearch/ composition/lifecycle/provider wiring
  internal/customer/application/ direct turn service below HTTP
  benchmarks/...           product cases and target rules

RagKit repository (temporary)
  document/chunk/representation primitives
  indexbundle, searchers, collapse, RRF, hydration, reranking

Optkit repository
  imported durable foundation
  RAG-TTC system registration
  trajectory/observation/trial/campaign execution

Judgekit repository (temporary)
  provider-neutral answer measurements under lightweight attribution
```

This avoids bulk movement before behavior is proved. Product semantics are placed in product-owned packages that remain correct after RagKit primitives move.

### 11.2 Final ownership

```text
Optkit
  record/artifact/space/episode/measure/experiment/campaign/...
  rag/... generic algorithms and event schemas
  judge/... after successful consolidation

RAG-TTC
  TTC source admission, configs, application services, tools, constraints
  customer/admin permissions and projections
  benchmark manifests and promotion policy

Coinvault
  Coinvault source admission, configs, application services, constraints
```

RagOpt disappears. RagKit and Judgekit are archived only after both products switch and retained fixtures pass.

### 11.3 Control, data, and query planes

```text
CONTROL PLANE                 DATA PLANE                    QUERY PLANE
Optkit campaign journal       RAG-TTC application service   projections/reports/UI
trial/arm scheduling          RagKit/Optkit RAG algorithms  retrieval funnel
budgets/leases                provider calls                rank movement
measurements/decisions        trajectory artifacts          episode timeline
```

Do not make the RAG service write campaign state. It emits trajectory evidence; Optkit owns orchestration.

## 12. Proposed RAG-TTC service APIs

These are design sketches. Keep names aligned with existing product vocabulary during implementation.

### 12.1 Prepared runtime identity

```go
type RuntimeIdentity struct {
    BundleID          string
    CorpusDigest      string
    ConfigDigest      string
    QueryTransformID  string
    RetrievalPolicyID string
    EvidencePolicyID  string
    RerankerIdentity  string
}
```

Identity is computed during preparation from the opened bundle and resolved configuration. It is not a signed attestation. It tells an engineer what this process intended and observed.

### 12.2 Retrieval service

```go
type Service struct {
    lexical   rag.Searcher
    vector    rag.Searcher
    content   content.Store
    sources   SourceCatalog
    routes    map[RouteID]PreparedRoute
    identity  RuntimeIdentity
}

type Request struct {
    Query        rag.Query
    Route        RouteID
    Limit        int
    AllowedRoles []string // server-owned; not model input
}

type Result struct {
    Identity RuntimeIdentity
    Query    rag.Query
    Route    RouteObservation
    Stages   []Stage
    Evidence []rag.Evidence
    Failure  *Failure
}

func (s *Service) Retrieve(ctx context.Context, req Request, sink TraceSink) (Result, error)
```

The service knows no Geppetto registry and no session ledger. It can be called by serving, tests, eval, or Optkit.

### 12.3 Tool adapter

```go
type Tool struct {
    service *retrieval.Service
    ledger  *EvidenceLedger
    policy  SessionPolicy
}

func (t *Tool) Run(ctx context.Context, input SearchInput) (SearchOutput, error) {
    req := retrieval.Request{
        Query: rag.Query{Text: input.Query},
        Route: t.policy.RouteFor(input),
        Limit: clamp(input.Limit, t.policy.DefaultResults, t.policy.MaxResults),
        AllowedRoles: t.policy.AllowedRoles,
    }
    result, err := t.service.Retrieve(ctx, req, traceSinkFromContext(ctx))
    if err != nil { return SearchOutput{}, err }
    return t.ledger.Admit(result)
}
```

The model can request a query and bounded limit. It cannot select provider credentials, corpus path, arbitrary source policy, or undeclared route configuration.

### 12.4 Stage trace

```go
type StageKind string

const (
    StageLexicalRaw       StageKind = "retrieval.lexical.raw"
    StageVectorRaw        StageKind = "retrieval.vector.raw"
    StageCollapsed        StageKind = "retrieval.collapsed"
    StagePolicyFiltered   StageKind = "retrieval.policy_filtered"
    StageFused            StageKind = "retrieval.fused"
    StageReranked         StageKind = "retrieval.reranked"
    StageHydrated         StageKind = "evidence.hydrated"
    StageAdmitted         StageKind = "evidence.admitted"
)

type Stage struct {
    Kind       StageKind
    Route      RouteID
    InputCount int
    OutputCount int
    RankedRefs artifact.Ref
    Duration   time.Duration
    Status     string
    ErrorClass string
}
```

Large candidate lists belong in artifacts, not inline campaign events.

### 12.5 Evidence ledger

```go
type EvidencePolicy struct {
    Scope    string // "conversation" or "episode"
    DedupeBy string // "chunk"
    MaxItems int
    MaxRunes int
}

type EvidenceLedger interface {
    Admit([]rag.Evidence) (items []Citation, omitted int)
    Snapshot() []Citation
    PolicyID() string
}
```

RAG-TTC's current ledger is already close. Add explicit policy identity and keep construction per conversation.

### 12.6 Direct customer application service

```go
type CustomerApplication interface {
    RunTurn(
        context.Context,
        CustomerTurn,
        episode.Sink,
    ) (CustomerResult, error)
}
```

The existing `realruntime.Composer` is a good composition boundary, but an intern should expose one direct turn method used by both HTTP/sessionstream and Optkit. The service owns domain execution, not WebSocket behavior.

## 13. Retrieval pseudocode

```text
function Retrieve(request, sink):
    validate request query, route, limit, server policy
    route = preparedRoutes.lookupOrDefault(request.route)

    emit retrieval.query(identity, query, route)

    lexical = []
    if route.lexical enabled:
        raw = lexicalSearcher.search(transformLexical(query), route.lexicalTopK)
        emit lexical.raw(raw refs)
        lexical = collapseRepresentationsToChunks(raw)
        lexical = filterByPreparedSourcePolicy(lexical, request.allowedRoles)
        renumber ranks
        emit lexical.filtered(lexical refs)

    vector = []
    if route.vector enabled:
        raw = vectorSearcher.search(transformVector(query), route.vectorTopK)
        emit vector.raw(raw refs)
        vector = collapseRepresentationsToChunks(raw)
        vector = filterByPreparedSourcePolicy(vector, request.allowedRoles)
        renumber ranks
        emit vector.filtered(vector refs)

    fused = weightedRRF(lexical, vector, route.weights, route.rrfConstant)
    emit retrieval.fused(fused refs)

    if route.reranker enabled:
        authorizedPool = filterAgainBeforeExternalBoundary(fused)
        hydratedPool = hydrate(authorizedPool[0:route.rerankPool])
        reranked, err = reranker.rerank(query, titleAndHeadingPrefixed(hydratedPool))
        if err:
            emit retrieval.reranked(status="degraded", errorClass=classify(err))
        else:
            fused = blendRRF(fused, reranked)
            emit retrieval.reranked(fused refs)

    final = filterAgain(fused)
    evidence = hydrate(final[0:request.limit])
    verify every chunk/document relation
    emit evidence.hydrated(evidence refs)

    return Result(identity, query, route, stages, evidence)
```

The repeated policy checks are defense against future internal call paths, not a claim of hostile-process isolation.

## 14. Optkit registration pseudocode

### 14.1 Configuration

```go
type RetrievalConfig struct {
    Route          RouteID
    BM25TopK       int
    VectorTopK     int
    RRFConstant    float64
    VectorWeight   float64
    Reranker       *RerankerConfig
    ResultLimit    int
    AllowedRoles   []string
}
```

Provider credentials and open handles are not configuration. They belong to the environment/prepared executable.

### 14.2 Variables

Start with a deliberately small search space:

```go
variables := []space.Variable[RetrievalConfig, _]{
    ResultLimit,   // e.g. 5, 8, 10
    VectorWeight,  // e.g. 0.5, 1.0, 1.5
    RerankPool,    // only when reranker enabled
}
```

Do not expose source authorization, customer/admin permissions, provider credentials, or hidden-dataset roles as search variables.

### 14.3 Executable

```go
type RetrievalExecutable struct {
    service *retrieval.Service
    cfg     RetrievalConfig
}

func (e *RetrievalExecutable) Run(
    ctx context.Context,
    c RetrievalCase,
    sink episode.Sink,
) (episode.RunResult, error) {
    result, err := e.service.Retrieve(ctx, retrieval.Request{
        Query: c.Query,
        Route: e.cfg.Route,
        Limit: e.cfg.ResultLimit,
        AllowedRoles: e.cfg.AllowedRoles,
    }, NewEpisodeTraceSink(sink))

    if err != nil {
        return terminalFailureWithPartialResult(result, err), nil
    }
    ref := attachResult(result)
    return completed(ref, result.Usage), nil
}
```

Product failures should normally return a terminal `RunResult`; Go errors mean the caller could not obtain valid terminal custody.

### 14.4 Measurements

First instruments should be deterministic:

```text
retrieval.recall_at_k
retrieval.precision_at_k
retrieval.mrr
retrieval.ndcg_at_k
retrieval.required_group_coverage
retrieval.source_role_coverage
retrieval.unique_documents
retrieval.max_chunks_per_document
retrieval.latency_ms
retrieval.embedding_calls
retrieval.reranker_calls
retrieval.failure
```

Add Judgekit only for answer-level constructs that deterministic contracts cannot settle.

## 15. Implementation phases

### Phase 0: Land the supplied Optkit baseline

**Goal:** make the tested archive the actual repository state.

Steps:

1. Create a clean import branch in `optkit`.
2. Preserve ticket files and repository CI plumbing intentionally.
3. Import archive source with module path `github.com/go-go-golems/optkit`.
4. Record archive revision `1786d1da86c9e03316ed71336bbc993fa30531f0` in the commit and changelog.
5. Run:

   ```bash
   gofmt -w .
   go test ./... -count=1
   go test -race ./... -count=1
   CGO_ENABLED=0 go test ./... -count=1
   ```

6. Do not “improve” public APIs in the provenance import commit.
7. Follow with focused repository-normalization/CI fixes.

**Exit gate:** the checked-out repository reproduces the archive's green test suite and Numbergame demo.

### Phase 1: Freeze cross-product semantic fixtures

**Goal:** preserve behavior before refactoring.

Create a small fixture pack in RAG-TTC containing:

- two simple positive queries;
- one multi-source/comparison query;
- one no-answer/abstention query;
- one source-role constrained query;
- one malformed provider output;
- one provider failure;
- expected stage candidate IDs for a tiny deterministic corpus;
- expected citation labels and answer-contract outcomes.

Use Coinvault tests as patterns, especially:

- `internal/knowledge/service_test.go`;
- `internal/knowledge/tool_test.go`;
- `internal/knowledge/eval_test.go`;
- `internal/knowledge/candidate_pool_test.go`;
- `cmd/coinvault/cmds/knowledge_ragopt_trace_test.go`.

**Exit gate:** current RAG-TTC behavior is characterized without live provider credentials.

### Phase 2: Split RAG-TTC retrieval semantics from the model tool

**Goal:** one canonical retrieval service callable from serving and evaluation.

Steps:

1. Introduce `pkg/ttc/retrieval` or an equivalently owned package.
2. Move channel search, collapse, RRF, route augmentation, source resolution, and hydration from `SearchTool.RunRoute` into `Service.Retrieve`.
3. Keep the current `SearchInput` and `SearchOutput` as the model adapter contract, but make them call the service.
4. Keep evidence ledger construction in `Handle.NewSessionRegistry` so evidence remains conversation-scoped.
5. Compile loaded route configuration into prepared routes. Do not leave `IntentRoutingConfig` as validation-only data.
6. Add direct service tests and retain tool tests.

**Exit gate:** the same deterministic request through direct service and tool adapter yields equivalent ranked chunks and route observation.

### Phase 3: Add Coinvault-style attribution and policy boundaries

**Goal:** every result explains its pipeline.

Add:

- `RuntimeIdentity` to service preparation;
- bundle/config/query-transform/retrieval/evidence/reranker IDs to tool output and traces;
- server-selected allowed source roles;
- filtering before fusion and external reranking when role policy is active;
- explicit lexical-only, vector-only, hybrid, and configured named routes;
- reranker identity and graceful fallback if a reranker is selected;
- effective-limit value and source.

Keep identities readable and deterministic. No signatures.

**Exit gate:** a trace consumer can reject an experiment where intended and observed route, bundle, transform, reranker, or evidence policy differ.

### Phase 4: Build the proper deterministic retrieval evaluation

**Goal:** replace opaque runner conclusions with stage-aware evidence.

Steps:

1. Define a strict versioned TTC retrieval suite schema.
2. Separate positive, policy-negative, and answer/judge-only cases.
3. Define required evidence groups rather than only single relevant IDs.
4. Add a stage diagnostic analogous to Coinvault's candidate-pool report.
5. Preserve failed query rows.
6. Compute per-query and aggregate metrics with explicit denominators.
7. Create one deterministic CI corpus and one larger local-development suite.
8. Emit machine-readable JSON plus a concise report projection.

**Exit gate:** an intern can answer “where did the target disappear?” from one artifact.

### Phase 5: Make the answer path explicit

**Goal:** proper end-to-end grounded answering, not retrieval alone.

Steps:

1. Define `CustomerApplication.RunTurn` below HTTP/sessionstream.
2. Make exact admitted evidence available to the answer contract.
3. Validate structured output, citation labels, abstention state, and widget intent.
4. Record route/tool calls, admitted evidence, model requests, answer output, citations, and terminal status.
5. Distinguish:
   - retrieval miss;
   - policy removal;
   - context/admission truncation;
   - generation provider failure;
   - malformed output;
   - unsupported claim;
   - unresolved citation;
   - presentation/widget failure.
6. Add transport-conformance tests over deterministic fixtures.

**Exit gate:** direct application execution and served execution produce equivalent domain events and answer-contract outcomes.

### Phase 6: Register the RAG-TTC retrieval system in Optkit

**Goal:** first real product campaign.

Steps:

1. Add the minimal system/preparation registry seam missing from the archive.
2. Define typed `RetrievalConfig`, variables, and `RetrievalCase`.
3. Materialize three to five fixed snapshots/arms.
4. Run a complete-block trial over the deterministic suite.
5. Emit stage trajectories through `episode.Writer`.
6. Emit deterministic measurements under an explicit epoch.
7. Add per-arm summaries and paired differences against baseline.
8. Run locally through SQLite + filesystem CAS.
9. Kill/restart at tested boundaries and verify terminal custody.

Start with manual fixed arms. Do not implement adaptive search yet.

**Exit gate:** one RAG-TTC multi-arm retrieval campaign survives restart and produces inspectable trajectories, estimates, and a decision artifact.

### Phase 7: Add Judgekit under the lightweight model

**Goal:** answer-quality measurements with correct attribution.

Steps:

1. Create `eval.Instance` from the sealed answer episode and admitted evidence.
2. Implement evidence-hidden claim extraction using a restricted input type.
3. Record contract, protocol, prompt version/rendered digest, expected and observed model, and cache mode.
4. Use `CacheBypass` for reliability repeats.
5. Store the sealed report as an artifact and emit typed Optkit observations.
6. Keep deterministic answer-contract failures as separate observations; do not ask the judge to replace them.
7. Apply calibration only to compatible report populations.

**Exit gate:** the same historical answer trajectory can be remeasured under a new judge epoch without rerunning retrieval or rewriting old observations.

### Phase 8: Consolidate and delete

After the real slice works:

1. Move only exercised generic RagKit packages into `optkit/rag`.
2. Rewrite RAG-TTC imports and rerun semantic fixtures.
3. Port Coinvault to the same generic primitives while retaining product policy.
4. Move Judgekit into `optkit/judge` only after both product measurement adapters are clear.
5. Replace RagOpt candidate/run/gate paths with Optkit snapshots/trials/decisions.
6. Switch last callers.
7. Delete old imports and orchestration paths.
8. Archive repositories with migration pointers.

No compatibility shims are required. The cutover gate is “new path passes, last caller switches, old path is deleted.”

## 16. Recommended first pull requests

### PR A: Import Optkit source baseline

Files: whole archive plus repository normalization. No semantic refactor.

### PR B: RAG-TTC retrieval fixture pack

Files likely under:

```text
rag-ttc/pkg/ttc/search/testdata/
rag-ttc/benchmarks/customer/document-qa/
rag-ttc/internal/testfixture/
```

### PR C: Direct retrieval service

Start review at current `pkg/ttc/search/search.go:172`. Extract service semantics while preserving `SearchTool` externally.

### PR D: Runtime identity and route compiler

Start at `pkg/ttc/toolconfig/types.go`, `toolconfig/load.go`, and `internal/customer/ragsearch/ragsearch.go`. Prove checked-in route configuration becomes prepared behavior.

### PR E: Stage-aware retrieval evaluation

Replace one path through `cmd/rag-ttc/cmds/experiments/answerquality/runner.go`, not the entire runner at once.

### PR F: Optkit RAG-TTC multi-arm campaign

Use the direct service and deterministic fixtures. Avoid live judges/providers in the first campaign.

## 17. Testing strategy

### 17.1 Unit laws

- chunk/document identity is stable;
- representation collapse produces unique deterministic chunks;
- RRF ties break deterministically;
- source policy cannot be widened by model input;
- disallowed candidates do not affect fusion rank;
- evidence labels are stable and scoped;
- item/rune limits are exact;
- query/retrieval/evidence policy IDs change only when semantics change;
- reranker fallback preserves fused order;
- route compiler rejects unresolved representation/source-role references.

### 17.2 Service integration tests

- open verified bundle and exact embedding identity;
- direct service versus tool-adapter parity;
- lexical-only path requires no query embedder;
- hybrid path makes expected embedding call;
- configured route is actually selected;
- provider failure remains terminal evidence;
- candidate-stage trace resolves all artifact references.

### 17.3 Answer tests

- exact citations resolve to admitted chunks;
- hallucinated citation labels fail;
- safe abstention on malformed output;
- unsupported claims stay separate from citation syntax;
- customer/admin tools and source policy remain separate;
- direct application versus transport event parity.

### 17.4 Experiment tests

- deterministic complete-block expansion;
- every case appears in every arm or is explicitly missing;
- target-resolver epoch is explicit;
- failed cells remain in denominator according to declared policy;
- no cross-epoch aggregation;
- restart after lease, completion, and observation boundaries;
- duplicate terminal completion is idempotent;
- budget actuals survive failure.

### 17.5 Judge tests

- extraction input cannot access evidence/reference/required facts;
- observed model must match expected identity;
- prompt version and rendered digest are recorded;
- repeat probes bypass cache;
- missing output counts as disagreement;
- durable report digest verifies on load;
- no signature/custody claims appear in docs or APIs.

## 18. Failure model

Failures are data, not rows to filter away.

```text
Failure class                   Owner
-----------------------------  -----------------------------
invalid config/route           preparation; no provider call
bundle verification failure    preparation; no episode start
search backend failure         retrieval episode
policy removal                 successful stage + observation
reranker outage                degraded stage; fused fallback
content hydration mismatch     retrieval failure
provider generation failure    answer episode failure
malformed answer               answer contract failure/abstention
judge provider failure         measurement failure, not product failure
missing pair                   experiment analysis failure/inconclusive
budget overage                 committed usage + policy violation
projection failure             query-plane issue; facts remain authoritative
```

A Go error should not erase a partial trajectory. Attach whatever evidence was obtained and return a valid terminal failure record whenever possible.

## 19. Security and privacy under the pragmatic model

“Not adversarial” does not mean “ignore product boundaries.”

Keep:

- customer/admin separation;
- server-selected provider profiles;
- server-selected corpus and source policy;
- no credentials in snapshots or artifacts;
- no model control over scopes;
- path confinement for config/assets;
- evidence as untrusted prompt data;
- bounded payload/reasoning retention;
- no hidden chain-of-thought requirement;
- redaction before persistent trajectories;
- explicit dataset exposure roles.

Defer:

- cryptographic producer attestations;
- cross-organization report verification;
- hostile local-user storage protection;
- key-management infrastructure;
- untrusted plugin execution.

If the trust model changes, open a new threat-model ticket rather than sprinkling cryptographic fields into current structs.

## 20. Decisions

### Decision: Product value before bulk kit consolidation

- **Context:** RAG-TTC already has working RAG components, while checked-out Optkit has no implementation and the archive has no RAG packages.
- **Options considered:** Consolidate every kit first; continue product-only custom runners; build one product vertical slice and consolidate exercised code afterward.
- **Decision:** Build the RAG-TTC proper-RAG vertical slice first, using current RagKit primitives and Coinvault semantics, while landing the Optkit baseline in parallel.
- **Rationale:** This validates ownership with a real consumer and delivers the requested RAG system sooner.
- **Consequences:** Temporary RagKit imports remain, but no new compatibility API is introduced; explicit deletion gates are required.
- **Status:** accepted

### Decision: Coinvault is a behavioral reference, not a source-copy target

- **Context:** Coinvault has strong RAG semantics mixed with Coinvault-specific authorization, tools, SQL, prompts, and transport evaluation.
- **Options considered:** Copy the package; abstract both products immediately; transfer laws and APIs into a TTC-owned service.
- **Decision:** Transfer verified bundle, pipeline ordering, runtime identity, evidence, diagnostics, and eval laws; keep TTC source/policy/application code local.
- **Rationale:** It avoids both duplication of product assumptions and premature universal abstractions.
- **Consequences:** Some code will be newly expressed rather than moved verbatim.
- **Status:** accepted

### Decision: Direct service below the model tool and transport

- **Context:** RAG-TTC retrieval currently lives in `SearchTool`; Coinvault evaluation currently pays transport overhead.
- **Options considered:** Keep tool as service; evaluate over HTTP; expose a canonical service used by both.
- **Decision:** Separate retrieval/answer services from Geppetto and HTTP adapters.
- **Rationale:** Serving and evaluation then execute the same domain composition without transport coupling.
- **Consequences:** Requires parity tests and a clear per-session ledger factory.
- **Status:** accepted

### Decision: Lightweight research attribution, not adversarial custody

- **Context:** Development is local and operated by trusted engineers.
- **Options considered:** Minimal validation; full hardened Judgekit architecture; boundary-focused attribution.
- **Decision:** Validate shape and semantic agreement at load, prepare, execute, and durable-consumption boundaries. Do not add signatures or immutable typestate.
- **Rationale:** Accidental experiment mismatch is the real current threat.
- **Consequences:** Artifacts are not cross-trust attestations.
- **Status:** accepted

### Decision: Deterministic retrieval before LLM judging

- **Context:** Many RAG failures are visible in candidate ranks, source roles, evidence groups, citations, and contracts.
- **Options considered:** Add Judgekit immediately; score only final answers; establish deterministic funnel evidence first.
- **Decision:** Build stage-aware retrieval and answer-contract measurements before judge integration.
- **Rationale:** It localizes failures and prevents a judge score from hiding mechanical defects.
- **Consequences:** Initial Optkit slice measures retrieval, not the complete customer experience.
- **Status:** accepted

### Decision: Manual multi-arm campaign before adaptive optimization

- **Context:** The framework's long-term ambition includes reflective and adaptive search.
- **Options considered:** Implement GEPA first; retain only RagOpt pairs; run fixed complete-block arms.
- **Decision:** Run fixed retrieval arms through the durable substrate first.
- **Rationale:** It tests snapshots, episodes, measurements, missingness, budgets, restart, and UI without optimizer confounding.
- **Consequences:** Search automation remains deferred.
- **Status:** accepted

## 21. Risks and mitigations

### Risk: The archive import and ticket history collide

Mitigation: import with provenance in a dedicated commit; preserve `ttmp`; avoid semantic edits in the import commit.

### Risk: “Temporary” RagKit dependency becomes permanent

Mitigation: add named deletion tasks and an architecture test that prevents new product orchestration from entering RagKit.

### Risk: Product service duplicates generic RagKit behavior

Mitigation: product service composes generic search/collapse/fusion primitives; it owns TTC policy and runtime identity, not alternate search algorithms.

### Risk: Route config remains declarative fiction

Mitigation: route compiler integration tests assert each checked-in route selects concrete searchers, representations, source roles, and augmentation.

### Risk: Digests create false confidence

Mitigation: document them as snapshot/pipeline fingerprints, not proof of truth or provider behavior.

### Risk: Evaluation overfits the development suite

Mitigation: explicit data roles, separate promotion cases, exposure records before adaptive search, and human review.

### Risk: Huge answer-quality runner refactor stalls delivery

Mitigation: replace one deterministic retrieval path first; leave unrelated stages until the service/campaign path is proven.

### Risk: Admin and customer systems are accidentally reunited

Mitigation: share generic components only; keep application services, tool registries, policies, configs, and benchmarks separate.

## 22. Open questions

1. Which checked-in TTC corpus and index bundle should be the first deterministic campaign fixture?
2. Are all customer corpus documents public, or must source-role policy enforce customer/admin separation inside one bundle?
3. Which current intent-routing configuration is intended for serving rather than experiment-only use?
4. Which reranker provider is acceptable for a local optional arm, and what fallback latency is acceptable?
5. Should the first Optkit RAG campaign compare retrieval representations, fusion weights, result limits, or all as fixed arms?
6. Which answer output schema is canonical for customer serving today?
7. Which existing RAG-TTC human-review artifacts must be preserved during run-store replacement?
8. Should the imported Optkit SQLite binding remain CGO-based, or can the real repository use its existing SQLite dependency conventions?
9. What exact artifact redaction is required before customer turns become local campaign fixtures?
10. Which Judgekit cleanup items should land before its first RAG-TTC adapter: restricted extraction input, on-demand instance digest, or both?

None of these blocks Phase 0 or the deterministic fixture pack.

## 23. Intern reading and review order

1. Read this document Sections 1-10.
2. Read the source architecture Sections 2-7, 55-60, 62-67, 70-71, and 77.
3. Read the supplied implementation diary Steps 1-7.
4. Run the supplied archive tests.
5. Read RAG-TTC:
   - `pkg/ttc/search/search.go`;
   - `internal/customer/ragsearch/ragsearch.go`;
   - `internal/customer/realruntime/composer.go`;
   - `pkg/ttc/toolconfig/types.go` and `load.go`;
   - `cmd/rag-ttc/cmds/indexes/build.go`;
   - focused tests beside each file.
6. Read Coinvault:
   - `internal/knowledge/service.go`;
   - `internal/knowledge/tool.go`;
   - `internal/knowledge/evidence.go`;
   - `internal/knowledge/eval.go`;
   - `internal/knowledge/candidate_pool.go`;
   - `internal/knowledgebuild/manifest.go`.
7. Read Judgekit's lightweight guide before changing identity or judge APIs.
8. Write a failing semantic fixture before refactoring.
9. Keep each PR vertically executable.
10. Delete superseded paths only after the parity gate passes.

## 24. Code review checklist

For every implementation PR ask:

- What semantic fact is new?
- Is the fact product-specific or generic?
- Is its identity complete but not broader than necessary?
- Can a failure disappear from the result?
- Does model input control anything it should not?
- Does checked-in config actually reach runtime behavior?
- Are exact bundle, route, transform, evidence, and model identities observable?
- Can the behavior run without HTTP/WebSocket?
- Can the deterministic fixture run without provider credentials?
- Is the measurement epoch explicit?
- What happens on cancellation and provider failure?
- Can the campaign recover after the next side effect?
- Which old code becomes deletable?
- Did the PR add a compatibility shim without an explicit requirement?

## 25. Definition of done for the pragmatic milestone

The milestone is complete when:

- the supplied Optkit implementation is landed and green in `optkit`;
- RAG-TTC has one canonical retrieval service below Geppetto/HTTP;
- current customer serving uses that service;
- bundle, config, transform, retrieval, reranker, and evidence-policy identities are recorded;
- checked-in route configuration is compiled into runtime behavior;
- a strict deterministic retrieval suite emits stage-aware diagnostics;
- one direct customer turn path validates grounding/citation/abstention outcomes;
- one multi-arm RAG-TTC retrieval campaign runs through Optkit and survives restart;
- deterministic measurements and failures are inspectable by case and stage;
- Judgekit is integrated only under the lightweight attribution model or remains an explicitly deferred next phase;
- no automatic deployment occurs;
- deletion gates for custom runner and old kit imports are named and tracked.

The full unification program is complete later, when both Coinvault and RAG-TTC use Optkit, old product runner paths are deleted, `pkg/mixedttc` is classified away, and standalone RagOpt/RagKit/Judgekit modules are archived.

## 26. References

### Primary source artifacts

- `/home/manuel/workspaces/2026-08-24/use-optkit/sources/optkit-clean-slate-architecture-and-porting-guide.md`
- `/home/manuel/workspaces/2026-08-24/use-optkit/sources/implementation-diary.md`
- `/home/manuel/workspaces/2026-08-24/use-optkit/sources/optkit-implementation-source.zip`

### Ticket evidence

- `sources/01-repository-inventory.md`
- `reference/01-investigation-diary.md`
- `scripts/01-repository-inventory.sh`

### Judgekit policy

- `judgekit/ttmp/2026/08/17/JUDGEKIT-001--design-and-implement-judgekit/design-doc/03-lightweight-research-guarantees-after-merge.md`
- `judgekit/ttmp/2026/08/17/JUDGEKIT-001--design-and-implement-judgekit/design-doc/04-fully-hardened-judgekit-architecture.md`

### Coinvault anchors

- `coinvault/internal/knowledge/service.go`
- `coinvault/internal/knowledge/runtime_config.go`
- `coinvault/internal/knowledge/tool.go`
- `coinvault/internal/knowledge/evidence.go`
- `coinvault/internal/knowledge/eval.go`
- `coinvault/internal/knowledge/candidate_pool.go`
- `coinvault/internal/knowledgebuild/manifest.go`
- `coinvault/internal/webchat/evalchat/configured_runner.go`

### RAG-TTC anchors

- `rag-ttc/pkg/ttc/search/search.go`
- `rag-ttc/pkg/ttc/toolconfig/types.go`
- `rag-ttc/pkg/ttc/toolconfig/load.go`
- `rag-ttc/internal/customer/ragsearch/ragsearch.go`
- `rag-ttc/internal/customer/realruntime/composer.go`
- `rag-ttc/internal/customer/realruntime/resolver.go`
- `rag-ttc/cmd/rag-ttc/cmds/indexes/build.go`
- `rag-ttc/cmd/rag-ttc/cmds/experiments/answerquality/runner.go`

### Reusable kit anchors

- `ragkit/rag/types.go`
- `ragkit/rag/retrieval/retrieval.go`
- `ragkit/rag/answering/service.go`
- `ragopt/pkg/eval/runner.go`
- `ragopt/pkg/runstore/run.go`
- `judgekit/judging/claimjudge.go`
- `judgekit/judging/interfaces.go`
