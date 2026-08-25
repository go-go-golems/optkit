---
Title: Intern Guide to the RAG Optimization Workbench and Specialist UI
Ticket: OPTKIT-004
Status: complete
Topics:
    - architecture
    - implementation
    - judgekit
    - optkit
    - rag
    - rag-ttc
    - scientific-workflow
    - ui
    - visualization
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://internal/web/static/app.js
      Note: Current dependency-free campaign explorer and UI baseline
    - Path: repo://query/service.go
      Note: Current bounded read-only campaign query plane
    - Path: repo://system/registry.go
      Note: Optkit product preparation and execution boundary
    - Path: ws://judgekit/protocol/protocol.go
      Note: Attributed reproducible LLM measurement protocol
    - Path: ws://rag-ttc/pkg/ttc/optkitcampaign/campaign.go
      Note: Durable product campaign and reconciliation reference
    - Path: ws://rag-ttc/pkg/ttc/search/service.go
      Note: Policy-safe attributable staged retrieval pipeline
    - Path: ws://ragkit/rag/indexbundle/types.go
      Note: Immutable RAG bundle and embedding identities
ExternalSources: []
Summary: Evidence-backed architecture, workflow, UI design, and phased implementation guide for an end-to-end TTC RAG optimization workbench.
LastUpdated: 2026-08-25T19:03:00Z
WhatFor: Onboard an intern to the current system and provide an implementation-ready design for optimizing chunking, representations, embeddings, retrieval, fusion, reranking, context construction, answering, and LLM judging.
WhenToUse: Use before implementing new RAG optimization campaigns, projections, specialist UI screens, or Judgekit-backed measurements.
---


# Intern Guide to the RAG Optimization Workbench and Specialist UI

## 1. Executive Summary

This document describes a complete research and engineering environment for improving The Tree Center's retrieval-augmented generation system. The target is not a generic prompt playground. It is a durable scientific workbench in which a RAG specialist can propose a change, execute controlled experiments, inspect exactly where relevant evidence was gained or lost, compare quality against cost and latency, and promote only configurations supported by attributable evidence.

The optimized system spans the entire path from source documents to judged answers:

```text
source documents
  -> parsing and normalization
  -> chunking
  -> derived retrieval representations
  -> embeddings and immutable indexes
  -> query transformation and prepared route selection
  -> lexical and vector retrieval
  -> source-policy filtering
  -> collapse and fusion
  -> optional route augmentation
  -> bounded reranking
  -> evidence hydration and admission
  -> context construction and compression
  -> answer generation
  -> answer-contract validation
  -> claim and citation extraction
  -> deterministic measurements
  -> LLM-judge measurements
  -> paired estimates and promotion decisions
```

Every stage must retain an immutable identity, bounded artifacts, stage-level diagnostics, resource use, and links to its inputs. That lineage allows the workbench to answer a concrete question such as:

> The required product-height chunk survived lexical retrieval but fell below the rerank pool after vector-heavy fusion; increasing the lexical fusion weight restored it, improved paired answer correctness by 8.2 percentage points, added 24 ms at p95, and caused no source-policy violations.

The existing repositories already provide most of the foundation:

- **RagKit** owns reusable RAG-domain algorithms and immutable bundle contracts.
- **Optkit** owns domain-neutral snapshots, trials, episodes, artifacts, budgets, observations, estimates, and campaign journals.
- **RAG-TTC** owns Tree Center policy, prepared routes, direct retrieval and answering services, Optkit integration, and the product-facing specialist experience.
- **Judgekit** owns reusable, content-addressed LLM measurement protocols, reports, reliability probes, and calibration.

The proposed specialist UI remains read-only. Campaigns, parameters, and candidate configurations are authored through CLI/application workflows, often by LLM agents. The browser provides navigation, comparison, search, replay, provenance, and visualization. This separation keeps scientific writes explicit and scriptable while making evidence unusually easy to inspect.

The implementation should be incremental. It should not attempt to optimize the full Cartesian product of chunker, embedding, route, reranker, answer model, and judge in one campaign. Instead, it should use a progressive funnel, reuse immutable upstream artifacts, and advance candidates only after cheaper deterministic checks pass.

## 2. Audience and Reading Strategy

This guide is written for an intern who knows Go and basic RAG concepts but has not worked in this workspace. Read it in this order:

1. Read Sections 3 through 6 to learn the vocabulary, repository boundaries, and current implementation.
2. Read Sections 7 through 12 to understand the proposed pipeline, identity model, and optimization workflow.
3. Read Sections 13 through 17 before implementing query-plane projections or UI screens.
4. Read Sections 18 through 23 before changing campaign execution, measurement, or testing.
5. Use Sections 24 through 28 as the implementation checklist and review guide.

The words **current** and **proposed** are deliberate:

- **Current** means behavior observed in source files at the time of this ticket.
- **Proposed** means the architecture this document recommends but which may not exist yet.
- **Accepted decision** means a constraint already established by the implementation and prior tickets.
- **Proposed decision** means a design choice that still needs implementation review.

## 3. Core Vocabulary

### 3.1 Corpus

A corpus is a set of immutable source document revisions. Its digest changes when source content changes. A production comparison is meaningful only when every arm either uses the same corpus or explicitly declares that corpus construction is the treatment.

### 3.2 Chunk

A chunk is an exact byte range of one document. It remains source evidence. RagKit's `rag.Chunk` records document identity, ordinal, byte range, text, digest, and chunker identity (`ragkit/rag/types.go:20-30`).

### 3.3 Representation

A representation is searchable text derived from a chunk. Examples include raw text, a summary, a synthetic question, or a product-fact view. It is retrieval material, not final evidence. A representation hit must hydrate back to the source chunk before reaching a reranker, answer model, judge, or human (`ragkit/rag/types.go:32-42`).

### 3.4 Vector

A vector is an embedding of one representation. The representation identity, provider, model, dimensions, normalization, and index implementation are all part of the treatment.

### 3.5 Prepared route

A prepared route is a server-owned retrieval strategy. It binds already-open lexical/vector indexes, representations, source-role policy, pool sizes, fusion behavior, augmentation, and reranking. A caller may request a route name but may not submit providers, corpus paths, searchers, or policy overrides.

### 3.6 Snapshot

An Optkit snapshot is a content-derived identity over a system, schema, and canonical configuration artifact. It contains declarative configuration, never credentials or open provider/index handles (`optkit/space/snapshot.go:12-65`).

### 3.7 Arm

An arm is a named snapshot in an experiment. For example, `limit-8` and `rerank-cohere-30` can be arms.

### 3.8 Complete block

A complete-block trial runs every arm on every case and repeat. This controls query difficulty and enables paired differences. Optkit builds content-derived trial and episode identities from arm snapshots, dataset digests, cases, repeats, seeds, and protocol (`optkit/experiment/trial.go:58-153`).

### 3.9 Episode and trajectory

An episode is one arm/case/repeat execution. Its trajectory is an immutable ordered list of typed events and payload references. The trajectory is the reusable historical record from which later instruments can remeasure behavior.

### 3.10 Measurement epoch

A measurement epoch identifies one construct, instrument, protocol, implementation, calibration, and redaction policy. Changing the judge or deterministic evaluator creates a new epoch rather than mutating old observations (`optkit/measure/epoch.go:9-38`).

### 3.11 Observation and estimate

An observation is one measured value attached to a subject and epoch. An estimate aggregates observations, usually with paired comparisons. `experiment.PairedMean` requires complete valid pairs and exposes missing-pair failures (`optkit/experiment/estimate.go:17-77`).

### 3.12 Candidate artifact

A candidate artifact is a bounded, content-addressed stage output such as a ranked list of chunk IDs, a rerank pool, admitted evidence set, answer contract, or Judgekit report.

## 4. Product Thesis

The workbench should transform RAG tuning from intuition-driven configuration edits into a repeatable loop:

```text
hypothesis
  -> typed candidate configuration
  -> compiled immutable treatment
  -> complete-block execution
  -> stage and answer measurements
  -> paired estimates
  -> diagnosis
  -> next candidate or promotion
```

A useful system must satisfy five properties.

### 4.1 Attribution

A result must identify the corpus, bundle, configuration, query transform, retrieval policy, evidence policy, reranker, answer prompt/model, and judge protocol that produced it. A quality number without this identity is not scientific evidence.

### 4.2 Reuse

Changing a judge should not rebuild indexes or rerun answers. Changing an answer prompt should not rerun retrieval. Changing fusion weights should reuse lexical and vector channel outputs. The workbench must model invalidation as a directed acyclic graph rather than treating every candidate as a complete cold run.

### 4.3 Diagnosis

The system must retain enough stage information to identify the first point where a required evidence group disappears. Final Recall@K alone cannot distinguish a chunker failure from a policy filter, fusion demotion, reranker truncation, or context admission failure.

### 4.4 Safety

Source-role and authorization policy must execute before fusion and before any external reranker. Unauthorized text must not affect rankings, leave the process, or appear in retained previews.

### 4.5 Specialist usability

The browser should make complex evidence navigable without making experimental writes implicit. The specialist should be able to move from a campaign-level Pareto frontier into one query, one stage, one chunk, one claim, and one raw artifact while preserving browser history and shareable routes.

## 5. Repository Responsibilities

### 5.1 RagKit: reusable RAG-domain mechanics

RagKit defines the domain contracts used by products:

- `rag.Chunker`, `rag.Embedder`, `rag.Searcher`, and `rag.Reranker` are provider-neutral interfaces (`ragkit/rag/components.go:8-91`).
- `rag.Document`, `rag.Chunk`, `rag.Representation`, `rag.Vector`, `rag.Hit`, `rag.FusedHit`, and `rag.Evidence` define lineage-preserving values (`ragkit/rag/types.go:3-111`).
- `chunking.MarkdownHeading` keeps structural sections whole and falls back to overlapping fixed windows only for oversized sections (`ragkit/rag/chunking/markdown_heading.go:12-103`).
- `representations.Summaries` and `representations.Questions` build non-raw retrieval material whose IDs include source, content, kind, model, and prompt provenance (`ragkit/rag/representations/representations.go:1-194`).
- `retrieval.Collapse`, `WeightedRRF`, and `HydrateFromStore` provide deterministic collapse, fusion, and bounded source hydration (`ragkit/rag/retrieval/retrieval.go:14-194`).
- `indexbundle.Manifest` records corpus, chunks, representations, lexical backend, vector backend, embedding model, dimensions, and content digests (`ragkit/rag/indexbundle/types.go:14-74`).

RagKit must remain independently reusable. It should not import Optkit campaign orchestration or RAG-TTC product code.

### 5.2 Optkit: domain-neutral experiment control plane

Optkit owns the scientific record:

- canonical artifacts and filesystem/memory content-addressed stores;
- semantic IDs and canonical JSON;
- immutable typed snapshots;
- candidate patches and configuration lineage;
- complete-block trials and deterministic episode seeds;
- episode writers and sealed trajectories;
- campaign commands, append-only control events, and foldable state;
- SQLite work leases and resource budgets;
- measurement epochs, observations, estimates, and decisions;
- read-only campaign query projections and the embedded explorer.

The system registry deliberately stores heterogeneous product factories behind artifact/schema boundaries. `system.Factory` prepares an immutable snapshot, and `system.Prepared.Run` consumes a case artifact, seed, and episode sink (`optkit/system/registry.go:15-34`). The registry verifies the prepared system, snapshot, and case schema before execution (`optkit/system/registry.go:65-96`).

Optkit must not learn RAG-specific concepts such as chunkers, embeddings, source roles, or rerankers.

### 5.3 RAG-TTC: product composition and policy

RAG-TTC owns the Tree Center-specific application:

- verified corpus and bundle loading;
- prepared route compilation;
- source-role policy;
- direct retrieval service composition;
- evidence limits and tool description;
- customer question-answer application service;
- deterministic retrieval evaluation suites;
- Optkit campaign adapter and runner;
- specialist-facing query projectors and UI semantics.

`search.Service` is the direct, transport-independent retrieval boundary (`rag-ttc/pkg/ttc/search/service.go:40-59`). Its request contains only a query and checked route name (`service.go:14-21`). This prevents a model or browser from constructing arbitrary runtime dependencies.

### 5.4 Judgekit: reusable LLM measurement

Judgekit defines a content-addressed measurement system:

- `eval.Instance` contains input, candidate answer, admitted evidence, optional reference, required facts, metadata, and a digest (`judgekit/eval/instance.go:3-27`).
- `protocol.Protocol` pins measurement digest, actual model identity, rendered prompt digests, decoding, evidence order, parser, aggregator, and retries (`judgekit/protocol/protocol.go:3-67`).
- `assessment.Report` binds instance and protocol digests to claims, claim assessments, dimensions, raw artifacts, timing, and report digest (`judgekit/assessment/report.go:5-29`).
- `judging.ConfigurableJudge` supports explicit cache bypass for reliability experiments (`judgekit/judging/interfaces.go:38-67`).
- `audit.Reliability` compares perturbation pairs per construct rather than hiding instability in one score (`judgekit/audit/reliability.go:16-147`).
- `calibration.Calibrate` reports extraction recall, sensitivity, specificity, false-support rate, Brier score, ECE, and group slices (`judgekit/calibration/report.go:9-154`).

Judgekit remains reusable. RAG-TTC should integrate it through an Optkit instrument rather than moving its contracts into Optkit.

## 6. Current System: What Already Exists

### 6.1 Direct retrieval with stage traces

The current TTC retrieval service performs:

1. query validation and prepared route selection;
2. lexical search;
3. lexical collapse;
4. lexical source-role filtering;
5. vector search;
6. vector collapse;
7. vector source-role filtering;
8. weighted reciprocal-rank fusion;
9. optional route augmentation;
10. post-augmentation policy recheck;
11. optional bounded reranking with degraded fallback;
12. source hydration.

The implementation emits named `RetrievalStage` values with counts, chunk IDs, candidate-set digest, status, and error class (`rag-ttc/pkg/ttc/search/service.go:23-49`, `200-282`).

The policy boundary is already correctly placed. Hits are filtered before fusion, fused candidates are rechecked after augmentation, and only then may candidates enter reranking (`service.go:218-267`). The reranker validates that every returned item came from the bounded pool, preserves fused order on provider/hydration/shape failures, and emits a degraded stage (`service.go:347-430`).

### 6.2 Runtime research identity

`RuntimeIdentity` records:

- bundle ID;
- corpus digest;
- resolved configuration digest;
- query transform identity;
- retrieval policy identity;
- evidence policy identity;
- reranker identity;
- tool-description identity.

These are local research fingerprints, not signatures (`rag-ttc/pkg/ttc/search/identity.go:13-32`). `retrievalPolicyID` includes route, enabled channels, representation names, pool sizes, RRF constant, roles, augmentation, reranker, and rerank pool (`identity.go:50-76`).

### 6.3 Deterministic retrieval evaluation

The current evaluator checks intended versus observed treatment, preserves stage rankings, calculates rank movement, and diagnoses first loss. Its metrics include recall, precision, MRR, NDCG, group coverage, and source diversity (`rag-ttc/pkg/ttc/retrievaleval/evaluate.go:11-109`). Authorization-negative cases separately detect forbidden chunks after policy (`evaluate.go:157-201`).

This is an important design pattern: deterministic mechanical evidence remains distinct from answer-level and judge-level measurement.

### 6.4 Canonical answer application service

`customerapp.Service.RunTurn` composes a provider engine, prepared search session, bounded tool loop, evidence ledger, answer-contract interpretation, and typed domain events (`rag-ttc/pkg/ttc/customerapp/service.go:45-190`). It records provider calls, retrieval identities, evidence, failures, abstention, and contract validity without depending on HTTP or WebSocket transport.

This direct boundary is what answer-level Optkit episodes should call. A campaign should not drive the web transport to evaluate the domain application.

### 6.5 Durable retrieval campaign

The RAG-TTC Optkit adapter defines strict typed retrieval configuration and case schemas (`rag-ttc/pkg/ttc/optkitcampaign/system.go:18-71`). Product-owned `Executor` implementations retain providers and open indexes in memory while snapshots contain only preparation, route, and limit (`system.go:76-120`). The prepared system writes retrieval inputs, every retrieval stage, and completion output into the episode trajectory (`system.go:132-193`).

The campaign runner initializes complete blocks, schedules leased work, reconciles terminal queue results into the campaign journal, records deterministic observations, and calculates paired estimates (`rag-ttc/pkg/ttc/optkitcampaign/campaign.go:134-591`). Tests simulate interruption after lease, after terminal result, after observation, and during execution (`campaign_test.go:17-112`).

### 6.6 Read-only campaign explorer

Optkit's current query service exposes campaign summaries and bounded event pages. Preview limits are 256 KiB, confidential/restricted artifacts are not rendered, and event pages are capped at 500 (`optkit/query/service.go:15-24`, `173-281`).

The HTTP server registers only GET routes for health, campaigns, campaign detail, events, SSE, static assets, and the root application. It uses Go 1.22 `http.ServeMux`, embeds static assets, and applies a restrictive content-security policy (`optkit/internal/web/server.go:1-66`).

The current JavaScript renders:

- campaign list and global search;
- lifecycle stage rail;
- overview metrics;
- immutable configuration lineage;
- complete-block trial matrix;
- paired analysis;
- budget and event-kind summaries;
- journal timeline and artifact drawer.

Those functions are visible in `optkit/internal/web/static/app.js:113-330`. This is a strong read-only foundation, but it currently loads every campaign event into browser memory (`app.js:113-123`) and does not yet project RAG-specific stage, chunk, evidence, answer, or judge views.

## 7. Gap Analysis

The current implementation proves durability and stage attribution for a small retrieval campaign. A full optimization workbench still needs the following.

### 7.1 Unified but layered configuration

Current snapshots cover a small prepared route and limit. The full system needs typed configuration for:

- document normalization;
- chunking;
- representation generation;
- embedding and vector index;
- lexical index;
- query transformation;
- prepared route and fusion;
- reranking;
- context admission/compression;
- answer generation;
- deterministic instruments;
- Judgekit protocols.

This must remain layered. A judge configuration must not be embedded into an index identity, because doing so would force unnecessary recomputation.

### 7.2 Artifact dependency planning

The system needs an explicit invalidation planner. Today a campaign arm is one snapshot. The workbench should understand that changing fusion weights can reuse channel outputs while changing chunking invalidates all downstream artifacts.

### 7.3 Build-time campaign support

The first durable campaign starts from prepared indexes. Optimizing chunkers, representations, and embeddings requires episodes or jobs that build immutable bundles and measure build-stage behavior.

### 7.4 Answer and judge episodes

The direct answer service exists, but the Optkit adapter currently executes retrieval only. The next product systems must seal answer trajectories and measure historical answers under multiple deterministic and Judgekit epochs.

### 7.5 Rich domain projections

The query plane exposes generic journal events. The specialist UI needs server-side projectors for:

- stage funnels;
- candidate rankings;
- rank movements;
- first-loss diagnosis;
- chunk boundaries;
- evidence-to-claim relations;
- judge dimensions and calibration;
- cross-arm paired comparisons;
- cost/latency Pareto fronts.

### 7.6 Scale

The current browser recursively loads all events. Large campaigns need database-side pagination, entity routes, exact-sequence replay, structured search, and virtualization.

## 8. Proposed End-to-End Architecture

### 8.1 Control-plane and data-plane separation

```text
                       AUTHORING PLANE
            CLI commands and LLM campaign agents
                              |
                              v
              typed candidates / snapshots / trials
                              |
                              v
+------------------------------------------------------------------+
|                     OPTKIT CONTROL PLANE                         |
| campaign journal | queue/leases | budgets | artifacts | measures |
+------------------------------------------------------------------+
                              |
                 prepare immutable snapshot
                              |
                              v
+------------------------------------------------------------------+
|                  RAG-TTC PRODUCT DATA PLANE                      |
| RagKit bundles -> prepared routes -> retrieval -> answer service |
|                                -> Judgekit instruments           |
+------------------------------------------------------------------+
                              |
                 sealed events and artifacts
                              |
                              v
+------------------------------------------------------------------+
|                    READ-ONLY QUERY PLANE                         |
| generic projectors | RAG projectors | SSE | bounded artifact API |
+------------------------------------------------------------------+
                              |
                              v
                    SPECIALIST EXPLORER UI
```

The authoring plane decides what to run. The control plane records what should happen and what happened. The data plane executes product behavior. The query plane derives bounded views. The UI never mutates campaign state.

### 8.2 Full pipeline stage graph

```text
S0 source revision
 |
 v
S1 normalized document
 |
 v
S2 chunks -------------------------------+
 |                                        |
 +--> S3 raw representation               |
 +--> S3 summary representation           |
 +--> S3 synthetic-question representation|
 +--> S3 structured-fact representation   |
              |                            |
              v                            |
         S4 vectors                        |
              |                            |
              +---------+                  |
                        v                  v
query -> S5 query transform -> S6 lexical/vector channel results
                                   |
                             S7 policy filter
                                   |
                             S8 collapse/fusion
                                   |
                             S9 augmentation
                                   |
                            S10 policy recheck
                                   |
                            S11 rerank pool/order
                                   |
                            S12 source hydration
                                   |
                            S13 evidence admission
                                   |
                            S14 context construction
                                   |
                            S15 answer generation
                                   |
                            S16 answer contract
                                   |
                            S17 claims/citations
                              /             \
                             v               v
                    S18 deterministic     S19 Judgekit
                         measures             measures
                              \             /
                               v           v
                         S20 paired estimates
                                   |
                         S21 promotion decision
```

Each stage should have a stable schema and an identity derived from semantic inputs. A stage artifact should not include volatile timestamps in its semantic digest.

## 9. Proposed Layered Configuration Model

A configuration should be represented as composable references rather than one monolithic struct.

```go
type RAGSystemConfig struct {
    Corpus        CorpusConfigRef
    Chunking      ChunkingConfigRef
    Representation RepresentationSetRef
    Embedding     EmbeddingConfigRef
    Index         IndexConfigRef
    Retrieval     RetrievalConfigRef
    Context       ContextConfigRef
    Answer        AnswerConfigRef
}

type MeasurementPlan struct {
    Deterministic []InstrumentConfigRef
    Judges        []JudgeProtocolRef
}
```

A human-readable candidate might look like:

```yaml
corpus:
  revision: ttc-catalog-2026-08-25
  source_policy: public-product-guides-v3

chunking:
  algorithm: markdown-heading
  max_section_runes: 2200
  min_section_runes: 160
  overlap_runes: 160
  preserve_tables: true

representations:
  kinds:
    - raw
    - summary
    - synthetic-question
  summary:
    model: gpt-5-mini
    prompt_digest: sha256:...
  question:
    model: gpt-5-mini
    prompt_digest: sha256:...

embedding:
  provider: openai
  model: text-embedding-3-large
  dimensions: 1536
  normalization: cosine

lexical:
  backend: bleve
  analyzer: english-product
  title_boost: 2.0
  body_boost: 1.0
  top_k: 40

vector:
  backend: hnsw
  representation_kinds: [raw, summary, synthetic-question]
  top_k: 40
  ef_search: 128

fusion:
  algorithm: weighted-rrf
  rank_constant: 60
  weights:
    bm25: 1.0
    vector: 1.2

reranker:
  provider: cohere
  model: rerank-v3.5
  pool: 30
  output: 8
  fallback: fused-order

context:
  maximum_tokens: 6000
  maximum_items: 8
  deduplicate: chunk
  compression: source-preserving-extractive-v1
  order: evidence-utility

answer:
  provider: openai
  model: gpt-5
  orchestration_digest: sha256:...
  output_schema_digest: sha256:...
  maximum_provider_calls: 4
  require_citations: true
  abstention_policy: strict-v2
```

### 9.1 Configuration identity rule

Every semantic field that can change behavior must participate in identity. Operational fields that do not change behavior, such as a local cache directory, should not.

Bad identity:

```text
embedding identity = "openai/text-embedding-3-large"
```

Better identity:

```text
embedding identity = semantic digest of:
  provider
  model
  revision
  dimensions
  normalization
  input representation digest
  provider settings affecting vectors
```

### 9.2 Runtime handle rule

The snapshot stores identifiers and parameters. The composition root resolves them into:

- provider clients;
- credentials;
- open bundle handles;
- content stores;
- lexical indexes;
- vector indexes;
- reranker clients;
- answer engines;
- Judgekit generators.

The pattern already exists in `optkitcampaign.Factory` and should be extended, not replaced.

## 10. Artifact Dependency DAG and Reuse

### 10.1 Invalidation table

| Change | Reusable artifacts | Recomputed artifacts |
|---|---|---|
| Judge protocol | Everything through sealed answer | Judge reports, observations, estimates |
| Deterministic instrument | Everything through relevant trajectory | Observations, estimates |
| Answer prompt/model | Admitted evidence and earlier | Answer, claims, judges, estimates |
| Context policy | Hydrated evidence and earlier | Context, answer, claims, judges |
| Reranker | Fused pool and earlier | Reranked order and downstream |
| Fusion weights | Channel results and earlier | Fusion and downstream |
| Query transform | Bundle/indexes | Query results and downstream |
| Embedding model | Documents, chunks, representations | Vectors, vector index, retrieval and downstream |
| Representation prompt | Documents and chunks | Derived representations and affected downstream indexes |
| Chunker | Source documents only | Chunks and every downstream artifact |
| Corpus revision | Nothing except unrelated provider caches | Entire branch |

### 10.2 Proposed planner API

```go
type StageID string

type Materialization struct {
    Stage      StageID
    Identity   record.Digest
    Artifact   artifact.Ref
    Inputs     []record.Digest
}

type PlanRequest struct {
    Baseline SnapshotGraph
    Candidate SnapshotGraph
    Available map[record.Digest]Materialization
}

type Plan struct {
    Reuse   []Materialization
    Execute []StageExecution
}

func CompileReusePlan(req PlanRequest) (Plan, error)
```

### 10.3 Planner pseudocode

```text
function compileReusePlan(baseline, candidate, artifactIndex):
    plan = empty

    for stage in topologicalPipelineOrder:
        identity = semanticIdentity(stage, candidate.config, stage.inputIdentities)

        if artifactIndex contains verified(identity):
            plan.reuse(stage, artifactIndex[identity])
        else:
            plan.execute(stage, identity, stage.inputArtifacts)

    return plan
```

Verification is mandatory. A matching filename or user-supplied label is not sufficient evidence that an artifact may be reused.

## 11. Stage-by-Stage Optimization Guide

### 11.1 Source normalization

Treat normalization as an explicit stage. It may include:

- decoding and Unicode normalization;
- boilerplate removal;
- Markdown or HTML conversion;
- heading repair;
- table normalization;
- canonical URL and source-role assignment;
- document-level metadata extraction.

Measurements:

- source count and byte/rune count;
- conversion failures;
- removed-text ratio;
- malformed structure count;
- stable document-ID rate;
- source-role and access-scope distribution.

Do not silently mutate source text during chunking. A source revision should remain inspectable.

### 11.2 Chunking

Candidate variables:

- fixed, Markdown, Markdown-heading, semantic, or source-specific chunker;
- maximum and minimum size;
- overlap;
- sentence snapping;
- heading ancestry;
- parent/child hierarchy;
- table and list preservation;
- source-role-specific policies.

Mechanical measurements:

- chunk count;
- size histogram;
- overlap token ratio;
- duplicated source coverage;
- orphaned source ranges;
- boundary integrity;
- table/list split count;
- required-fact containment;
- downstream oracle recall.

A chunking campaign should not begin with expensive answer generation. It should first eliminate configurations that lose labeled facts or produce pathological fragmentation.

### 11.3 Derived representations and summarization

Index-time summaries and synthetic questions are retrieval representations. They are not evidence. This distinction is already encoded in RagKit.

Candidate variables:

- representation kind;
- source roles receiving each kind;
- prompt and model;
- extractive versus generative summarization;
- number of generated questions;
- representation token limit;
- deterministic validation rules.

Measurements:

- representation coverage by chunk;
- empty/failure rate;
- source-claim support;
- representation-to-source lexical overlap;
- retrieval lift versus raw-only;
- generation cost;
- index-size increase.

A generated representation must retain:

- source chunk ID;
- representation kind;
- generated text digest;
- model identity;
- prompt digest.

### 11.4 Embeddings

Candidate variables:

- provider and model revision;
- dimensions;
- input representation kinds;
- normalization;
- batching;
- quantization;
- exact versus ANN index;
- HNSW construction/search parameters.

Measurements:

- Recall@K against labeled targets;
- group coverage;
- nearest-neighbor stability;
- irrelevant-neighbor rate;
- vector build time;
- query latency;
- memory/disk size;
- embedding tokens and cost;
- missing/misaligned vector count.

Never compare embedding models on different unlabeled corpora without making corpus identity explicit.

### 11.5 Lexical retrieval

Candidate variables:

- analyzer and stemming;
- title/body fields and boosts;
- product-name synonyms;
- source-role fields;
- raw versus derived representation indexing;
- top-K.

Measurements:

- Recall@K and MRR;
- exact product-name success;
- rare-token success;
- target rank;
- query latency;
- channel-only wins and losses.

### 11.6 Query transformation

Candidate variables:

- verbatim;
- spelling normalization;
- entity-aware rewrite;
- multi-query expansion;
- hypothetical answer/HyDE;
- route classification;
- structured-first intent detection.

Measurements:

- original versus transformed query artifact;
- transform model/prompt identity;
- treatment adherence;
- incremental retrieval coverage;
- false expansion rate;
- provider calls, latency, and cost.

The server must still select only prepared routes. A query transform may propose a route name but cannot construct a route.

### 11.7 Source policy

Policy is a hard constraint, not an optimization weight. Candidate sets should record counts before and after filtering, but a policy-violating candidate is invalid regardless of average quality.

Required invariants:

- filtering occurs before fusion;
- filtering occurs before external reranking;
- post-augmentation results are rechecked;
- source text is hydrated only for allowed candidates;
- negative authorization cases remain separate from positive retrieval averages.

### 11.8 Collapse and hybrid fusion

Candidate variables:

- collapse by representation, chunk, document, or parent;
- lexical/vector pool depths;
- RRF constant;
- channel weights;
- representation-specific channels;
- learned fusion, if introduced later.

Measurements:

- channel contribution per final result;
- target rank movement;
- unique chunk/document coverage;
- dominance of one channel;
- first loss after fusion;
- fusion latency.

Weighted RRF is an excellent baseline because it combines ranks rather than incomparable raw channel scores. Learned fusion should be considered only after feature and holdout discipline exists.

### 11.9 Route augmentation

Examples include:

- parent expansion;
- neighboring chunk expansion;
- graph or connected-document expansion;
- product/guide cross-link expansion;
- structured entity lookup.

Augmentation must emit its own identity and trace. The post-augmentation policy recheck remains mandatory.

### 11.10 Reranking

Candidate variables:

- provider/model/revision;
- pool depth;
- output depth;
- document versus chunk text;
- query plus metadata formatting;
- score threshold;
- fallback behavior.

Measurements:

- NDCG and MRR delta from fused order;
- relevant promotions and demotions;
- first relevant rank;
- targets outside the pool;
- invalid provider response rate;
- degraded fallback rate;
- latency, tokens, and cost.

A reranker should not be credited for gains caused by a changed fusion pool. Those are different treatments.

### 11.11 Evidence admission

Admission applies conversational budgets and deduplication after retrieval. Candidate variables include:

- maximum items;
- maximum runes/tokens;
- dedupe key;
- source diversity floor;
- repeat evidence behavior;
- ordering.

Measurements:

- required groups admitted;
- evidence rejected by budget;
- duplicate rate;
- source diversity;
- evidence token efficiency;
- previously seen versus new evidence.

### 11.12 Context construction and compression

Context construction may reorder, group, trim, or summarize admitted evidence. Any summary must preserve source references.

Candidate variables:

- chunk versus parent grouping;
- redundancy removal;
- extractive compression;
- source-preserving generative compression;
- evidence ordering;
- token budget allocation;
- citation labels and formatting.

Measurements:

- required evidence retained;
- unsupported compression claims;
- token compression ratio;
- redundancy ratio;
- citation addressability;
- downstream answer lift.

### 11.13 Answer generation

Candidate variables:

- model/provider/revision;
- orchestration prompt;
- output schema;
- tool-call budget;
- maximum parallel tools;
- reserve-final-call behavior;
- abstention policy;
- temperature and decoding;
- answer length and style constraints.

Mechanical measurements:

- output-schema validity;
- citation existence and range;
- citation-to-evidence resolution;
- required field presence;
- provider/tool call counts;
- input/output/reasoning tokens;
- latency and cost;
- abstention status.

LLM-judge measurements should not replace these deterministic contracts.

### 11.14 Claim extraction and LLM judging

A Judgekit instance should be derived from a sealed answer trajectory, not from mutable caller values.

Recommended inputs:

- original user question;
- final candidate answer;
- admitted evidence only;
- optional reference answer;
- required facts;
- product/query metadata;
- exact retrieval and answer identities.

Recommended constructs:

- claim support/faithfulness;
- answer relevance;
- completeness;
- citation correctness;
- abstention appropriateness;
- policy/style compliance.

Judge attribution must include:

- instance digest;
- measurement contract digest;
- protocol digest;
- prompt template and rendered prompt digests;
- expected and observed provider/model identity;
- decoding and evidence order;
- parser and aggregator versions;
- cache mode;
- raw generation artifacts;
- report digest.

## 12. Progressive Optimization Funnel

The workbench should optimize in gates.

### Gate 0: data validity

Reject candidates with:

- malformed source or chunks;
- missing lineage;
- representation/vector count mismatch;
- invalid bundle manifest;
- policy metadata gaps.

### Gate 1: zero-provider mechanical tests

Run:

- chunk containment;
- source coverage;
- deterministic extractive representations;
- exact index integrity;
- policy-negative tests;
- fixture parity.

### Gate 2: retrieval tests

Run labeled retrieval suites with fixed corpus and case splits. Evaluate stage-aware metrics and first-loss diagnosis.

### Gate 3: reranker/context tests

Advance only competitive retrieval candidates. Measure rerank and context deltas while reusing upstream artifacts.

### Gate 4: answer-contract tests

Run answer generation on a development subset. Reject invalid contracts, citation failures, or runaway resource use.

### Gate 5: judged answer quality

Judge only sealed answers that passed deterministic contracts. Run calibration and reliability checks for the judge protocol.

### Gate 6: hidden holdout

Evaluate a small number of predeclared finalists. Do not use holdout results to repeatedly tune the same candidates.

### Gate 7: promotion/canary

Produce a frozen promotion manifest and, later, a bounded production canary. Online behavior should not silently overwrite offline evidence.

## 13. Experiment Design

### 13.1 Dataset roles

Maintain explicit roles:

- **fixture**: tiny deterministic cross-product contract;
- **development**: used for candidate generation and debugging;
- **validation**: used for periodic selection;
- **holdout**: used rarely for final comparison;
- **policy-negative**: forbidden/source-boundary cases;
- **judge calibration**: human-labeled claims and dimensions;
- **reliability probes**: semantically controlled perturbation pairs.

### 13.2 Complete blocks and pairing

For a baseline and candidate, run every case under both arms and compare within case:

```text
case delta = candidate score - baseline score
estimate   = mean(case deltas)
```

Pairing controls the fact that some queries are intrinsically harder. Missing or invalid pairs should invalidate the estimate rather than silently changing its denominator.

### 13.3 Repeats

Use repeats when the treatment contains stochastic generation or judging. Retrieval over frozen deterministic indexes generally does not need repeated identical executions unless auditing nondeterminism.

A repeat must have a deterministic seed and distinct episode identity.

### 13.4 Multiple comparisons

The optimizer should avoid selecting a winner from hundreds of noisy arms based only on the maximum observed mean. Recommended controls:

- progressive elimination;
- predeclared primary metrics;
- minimum effect size;
- bootstrap confidence intervals over paired deltas;
- validation split confirmation;
- limited holdout access;
- explicit number of candidates considered.

### 13.5 Multi-objective outcomes

Hard constraints:

- zero authorization leakage;
- valid artifact and answer contracts;
- bounded provider calls;
- bounded p95 latency and cost;
- no protected-group regression beyond tolerance.

Optimization objectives:

- answer correctness;
- faithfulness;
- retrieval group coverage;
- abstention correctness;
- latency;
- cost;
- context size.

Do not collapse these into one opaque score. Preserve a Pareto frontier.

### 13.6 Proposed estimate types

```go
type EstimateSet struct {
    Metric       string
    Baseline     string
    Treatment    string
    MeanDelta    float64
    MedianDelta  float64
    Confidence   Interval
    SampleSize   int
    Missing      int
    ByGroup      map[string]GroupEstimate
}

type ParetoPoint struct {
    Snapshot record.SnapshotID
    Quality  float64
    CostUSD  float64
    P95MS    float64
    Valid    bool
    DominatedBy []record.SnapshotID
}
```

## 14. Optimization Algorithms

The first implementation should not begin with Bayesian optimization. Start with understandable candidate generation.

### 14.1 Phase-one search strategy

1. Define a baseline.
2. Change one layer at a time.
3. Generate a small grid or hand-authored set.
4. Run mechanical gates.
5. Keep the Pareto-nondominated candidates.
6. Compose only proven improvements.
7. confirm composed candidates against validation and holdout.

### 14.2 Later strategies

Once the data model and invalidation planner are stable, consider:

- random search for mixed discrete/continuous spaces;
- successive halving;
- Bayesian optimization with conditional parameters;
- evolutionary search for representation/route combinations;
- contextual selection by query category.

Every optimizer proposal must still compile to a typed immutable snapshot and fixed trial before execution.

### 14.3 Candidate proposal pseudocode

```text
function proposeNext(campaignHistory, searchSpace, constraints):
    valid = filter history where all hard constraints pass
    frontier = pareto(valid, quality up, cost down, latency down)

    losses = aggregate first-loss diagnoses on frontier
    targetLayer = layer with highest recoverable loss

    proposals = mutate only targetLayer around frontier configurations
    proposals = remove duplicate semantic identities
    proposals = remove configurations violating static constraints

    return proposals with hypotheses explaining expected effect
```

The proposal should state a hypothesis, expected metrics, risks, and changed variables. The agent should not submit an unexplained parameter vector.

## 15. Durable Campaign Workflow

### 15.1 Authoring

An LLM agent or specialist creates:

- dataset manifest;
- baseline snapshot;
- candidate snapshots;
- trial protocol;
- resource budget;
- measurement plan;
- hypothesis and risks.

### 15.2 Compilation

The controller:

1. validates schemas and static constraints;
2. resolves configuration graph identities;
3. compiles reuse versus execution stages;
4. expands complete-block episode specs;
5. reserves budgets;
6. enqueues semantic work items;
7. records plan artifacts.

### 15.3 Execution

Workers:

1. lease compact work items;
2. prepare snapshots through the product registry;
3. emit stage events to an episode writer;
4. seal trajectory and result artifacts;
5. commit terminal queue state;
6. allow a reconciler to append journal facts and observations.

### 15.4 Reconciliation

The queue and journal are separate durable systems. The reconciler treats terminal queue results as recoverable truth and idempotently records missing usage, episode, and observation facts.

```text
function reconcile(work):
    record = queue.get(work.id)
    if record is not terminal:
        return

    completion = verifyAndDecode(record.result)

    if journal lacks UsageCommitted(work.id):
        budget.commit(reservation(work.id), completion.usage)
        append UsageCommitted

    if journal lacks EpisodeCompleted(completion.episode):
        append EpisodeCompleted

    for instrument in measurementPlan:
        observation = instrument.measure(completion.trajectory)
        if journal lacks ObservationRecorded(observation.id):
            append ObservationRecorded
```

### 15.5 Analysis and promotion

After all required observations exist:

- compute per-arm summaries;
- compute paired estimates;
- compute group slices;
- update Pareto frontier;
- evaluate hard constraints;
- record a decision or request another campaign.

## 16. Specialist UI Product Design

### 16.1 Design principles

The UI should be:

- **read-only**: no hidden mutations or browser-owned campaign state;
- **evidence first**: every metric links to episodes and artifacts;
- **progressively disclosed**: campaign summary first, raw JSON last;
- **keyboard navigable**: `/` for search, Escape for drawers, stable focus;
- **shareable**: every entity view has a route;
- **bounded**: server-side pagination and previews;
- **honest about missingness**: failed and unavailable are not rendered as zero;
- **retro monochrome**: semantic HTML, pure-white restrained Macintosh styling, no frontend build step;
- **accessible**: color is never the only carrier of status.

### 16.2 Information architecture

```text
Workbench
├── Campaigns
│   ├── Campaign overview
│   ├── Arms and configuration diff
│   ├── Cases and paired matrix
│   ├── Pareto frontier
│   └── Journal/integrity
├── Queries
│   ├── Query replay
│   ├── Pipeline microscope
│   ├── Rank waterfall
│   ├── Evidence admission
│   └── Answer and claims
├── Corpus
│   ├── Documents
│   ├── Chunk laboratory
│   ├── Representations
│   └── Embedding neighbors
├── Measurements
│   ├── Deterministic contracts
│   ├── Judge reports
│   ├── Calibration
│   └── Reliability probes
└── Artifacts
    ├── Provenance graph
    ├── Safe preview
    └── Digest verification
```

### 16.3 Screen A: Optimization cockpit

Purpose: identify what needs attention.

Header:

- campaign name/ID;
- corpus and bundle identity;
- status and integrity;
- baseline and leading candidate;
- active filters.

Summary cards:

- completed/failed episodes;
- primary paired delta and confidence;
- cost delta;
- p95 latency delta;
- authorization violations;
- degraded reranker count;
- promotion status.

Main regions:

1. Pareto plot.
2. arm leaderboard with hard-constraint state;
3. case groups with regressions;
4. first-loss-stage histogram;
5. recent journal events.

### 16.4 Screen B: Campaign comparison

A matrix with cases as rows and arms as columns.

Each cell shows:

- completed, failed, missing, or invalid;
- primary score;
- first loss stage;
- latency and cost badge;
- policy state.

Interactions:

- sort by candidate regression;
- filter by query group, source role, status, or loss stage;
- select two arms for a configuration and outcome diff;
- open one cell as an episode.

### 16.5 Screen C: Pipeline microscope

This is the central diagnostic screen.

```text
QUERY: "How large does Blue Ice grow?"

LEXICAL RAW       40  [target rank 2]
LEXICAL POLICY    35  [target rank 2]
VECTOR RAW        40  [target absent]
VECTOR POLICY     31
FUSED             24  [target rank 7]
AUGMENTED         27  [target rank 7]
POLICY RECHECK    24
RERANKED          24  [target rank 1]
HYDRATED           8
ADMITTED           6  [target E1]
CITED               1 claim
```

For each stage, show:

- status and degradation;
- input/output counts;
- ranked candidates;
- required and forbidden markers;
- rank changes from previous stage;
- source role and document;
- channel contributions;
- candidate artifact digest.

Clicking a chunk opens the rank waterfall and safe source preview.

### 16.6 Screen D: Rank waterfall

For one chunk across stages:

```text
chunk-b
  lexical.raw             #1
  lexical.collapsed       #1
  lexical.policy          #1
  vector.raw              #2
  vector.policy           #1
  fused                   #1
  reranked                #2
  evidence.admitted       E2
  answer                  cited by claims C1, C3
```

The screen should explain absence:

- never retrieved;
- collapsed behind another representation;
- forbidden by policy;
- outside rerank pool;
- demoted below evidence limit;
- rejected by context budget;
- admitted but not cited.

### 16.7 Screen E: Chunk laboratory

Display source text with chunk boundaries overlaid. Compare two chunker snapshots side-by-side.

For each chunk show:

- ID, range, ordinal, chunker;
- rune/token size;
- heading ancestry;
- overlap with neighbors;
- representation kinds;
- retrieval frequency;
- labeled facts contained;
- answers/claims that used it.

Controls are read-only filters and selectors. A `COPY EXPERIMENT COMMAND` action may generate a CLI command, but clicking it does not start work.

### 16.8 Screen F: Representation and embedding explorer

For a selected source chunk:

- raw representation;
- summary representation;
- synthetic questions;
- model and prompt digests;
- nearest neighbors per embedding arm;
- target ranks for benchmark queries.

A two-dimensional projection may be offered as a secondary diagnostic, but never as proof of retrieval quality. Labeled nearest-neighbor tables are primary.

### 16.9 Screen G: Answer and evidence studio

Three synchronized panes:

1. admitted evidence with citation labels;
2. final answer with citation highlights;
3. extracted claims and deterministic/Judgekit verdicts.

Selecting a claim highlights:

- supporting and contradicting evidence;
- citations used by the answer;
- required facts;
- judge label, score, confidence, and rationale;
- protocol and model identity.

### 16.10 Screen H: Judge calibration lab

Show:

- extraction recall;
- sensitivity and specificity;
- false-support rate;
- Brier score and ECE when confidence exists;
- human-versus-judge confusion matrix;
- performance by topic/difficulty/evidence kind;
- reliability agreement by perturbation kind;
- disagreements with raw report links.

The UI must not present one global “judge quality” score.

### 16.11 Screen I: Pareto explorer

Axes can include:

- quality versus cost;
- quality versus p95 latency;
- recall versus context tokens;
- judge score versus deterministic contract score.

Display:

- dominated versus nondominated candidates;
- invalid candidates separately;
- baseline marker;
- confidence/missingness;
- snapshot diff on selection.

### 16.12 Screen J: Provenance inspector

Every entity should traverse:

```text
campaign
  -> trial
  -> arm snapshot
  -> episode
  -> trajectory
  -> stage artifact
  -> evidence artifact
  -> answer artifact
  -> Judgekit instance/report
  -> measurement epoch
  -> observation
  -> estimate
  -> decision
```

The inspector should show identity, schema, digest, size, sensitivity, producing event, input references, and verification state.

## 17. Proposed Query Plane

### 17.1 Why server-side projectors

The browser should not decode every historical schema or load a whole large campaign. Projectors provide versioned, bounded views while the journal and artifacts remain authoritative.

### 17.2 Proposed routes

```text
GET /api/v1/campaigns
GET /api/v1/campaigns/{campaignID}
GET /api/v1/campaigns/{campaignID}/events?after=&limit=
GET /api/v1/campaigns/{campaignID}/stream

GET /api/v1/campaigns/{campaignID}/arms
GET /api/v1/campaigns/{campaignID}/cases?cursor=&limit=&group=&status=
GET /api/v1/campaigns/{campaignID}/comparisons?baseline=&treatment=
GET /api/v1/campaigns/{campaignID}/pareto

GET /api/v1/episodes/{episodeID}
GET /api/v1/episodes/{episodeID}/trajectory?after=&limit=
GET /api/v1/episodes/{episodeID}/pipeline
GET /api/v1/episodes/{episodeID}/answer

GET /api/v1/chunks/{chunkID}
GET /api/v1/chunks/{chunkID}/lineage
GET /api/v1/chunks/{chunkID}/rank-waterfall?episode=

GET /api/v1/observations/{observationID}
GET /api/v1/judge-reports/{digest}
GET /api/v1/artifacts/{digest}/metadata
GET /api/v1/artifacts/{digest}/preview
```

No POST, PUT, PATCH, or DELETE route belongs in this server.

### 17.3 Proposed pipeline projection

```go
type PipelineView struct {
    APIVersion string
    Campaign   record.CampaignID
    Episode    record.EpisodeID
    Arm        string
    Case       string
    Query      string
    Identity   RuntimeIdentityView
    Stages     []PipelineStageView
    FirstLoss  *LossView
    Generated  time.Time
}

type PipelineStageView struct {
    Name       string
    Status     string
    InputCount int
    OutputCount int
    Candidates []RankedCandidateView
    Artifact   artifact.Ref
    ErrorClass string
}
```

### 17.4 Proposed comparison projection

```go
type CaseComparison struct {
    CaseID      string
    Groups      []string
    Baseline    Cell
    Treatment   Cell
    Deltas      map[string]float64
    FirstChange string
}

type Cell struct {
    Episode record.EpisodeID
    Status  string
    Metrics map[string]*float64
    CostUSD *float64
    LatencyMS *float64
    FirstLoss string
}
```

Use pointers or explicit missing states. Missing is not zero.

### 17.5 Pagination and search

Projectors must support:

- opaque cursors or stable sequence cursors;
- bounded page size;
- database-side filtering;
- deterministic order;
- exact-entity lookup;
- cancellation;
- sensitivity-aware previews.

The current event sequence is already appropriate for event pagination and SSE. Entity lists should use their own stable compound cursors.

## 18. Frontend Architecture

### 18.1 Technology

Keep the existing dependency-free approach:

- semantic HTML;
- plain CSS custom properties;
- plain JavaScript modules or one small application file;
- `go:embed` for production assets;
- hash or History API routes;
- no frontend build step.

This is sufficient because the UI is a read-only document explorer, not a collaborative editor.

### 18.2 Proposed client modules

If `app.js` becomes too large, split it without adding a bundler:

```text
internal/web/static/
  index.html
  styles.css
  app.js
  api.js
  router.js
  state.js
  views/
    campaigns.js
    campaign.js
    comparison.js
    pipeline.js
    chunks.js
    answer.js
    judges.js
    artifact.js
  components/
    table.js
    sparkline.js
    drawer.js
    status.js
```

Use native ES modules and embed the files.

### 18.3 State model

```js
const state = {
  route: parseRoute(location),
  page: null,
  filters: {},
  selected: null,
  request: null,
  stream: null,
};
```

Rules:

- route is the source of navigation state;
- server responses are replaceable caches, not authority;
- cancel old fetches on route changes;
- do not retain entire campaign event histories when a projector exists;
- drawers must be deep-linkable when they represent durable entities.

### 18.4 Visual language

Use pure white surfaces, black rules, restrained gray, monospaced IDs, compact labels, and Macintosh-like controls. Status should use symbols and text in addition to color:

```text
● COMPLETE
○ PLANNED
! DEGRADED
× FAILED
— MISSING
```

### 18.5 Accessibility

Required:

- meaningful heading hierarchy;
- table captions and header scopes;
- keyboard-operable rows and drawers;
- visible focus;
- dialog semantics and focus return;
- `aria-live` for connection state;
- charts accompanied by data tables;
- minimum target size;
- no horizontal information available only on hover.

## 19. Agent and CLI Workflow

The specialist UI should make authoring easy without performing it.

### 19.1 Agent proposal contract

An agent proposal should include:

```yaml
hypothesis: >
  Increasing lexical weight should restore exact product-name targets lost
  during fusion without increasing provider calls.

baseline_snapshot: snapshot:...
changes:
  retrieval.fusion.weights.bm25: 1.4
  retrieval.fusion.weights.vector: 1.0

expected_effects:
  - higher MRR on exact-name queries
  - unchanged authorization-negative outcomes

risks:
  - weaker semantic-query recall

trial:
  dataset: validation-product-qa-v3
  repeats: 1
  primary_metric: retrieval.mrr

budget:
  retrieval_queries: 500
  reranker_calls: 0
```

### 19.2 Suggested CLI groups

```text
rag-ttc optimize config validate
rag-ttc optimize candidate propose
rag-ttc optimize campaign plan
rag-ttc optimize campaign run
rag-ttc optimize campaign resume
rag-ttc optimize campaign inspect
rag-ttc optimize campaign compare
rag-ttc optimize artifact verify
rag-ttc optimize command render
```

`command render` can produce copyable commands or agent prompts from a UI-selected configuration diff.

### 19.3 Promotion manifest

A winning configuration should produce a signed-off but not necessarily cryptographically signed manifest containing:

- corpus/bundle identities;
- every layered configuration digest;
- route and policy identities;
- answer and judge protocols;
- campaign/trial/estimate references;
- hard-constraint results;
- approver and timestamp;
- deployment target.

These are research and operational fingerprints, not a custody system.

## 20. Security, Privacy, and Policy

### 20.1 Prepared-route authority

Never expose searcher construction, corpus paths, provider configuration, or allowed roles as model/browser input. Only server-prepared names are selectable.

### 20.2 External provider boundary

Before sending candidates to a reranker or judge:

- enforce source policy;
- bound count and text size;
- record provider/model identity;
- record redaction policy;
- preserve usage even on error;
- retain safe failure evidence.

### 20.3 Artifact previews

Preview endpoints must enforce:

- sensitivity policy;
- media-type allowlist;
- byte ceiling;
- JSON/text escaping;
- no arbitrary filesystem paths;
- no browser-supplied artifact metadata.

### 20.4 Prompt injection

Retrieved content is untrusted data. Context formatting must delimit sources and keep instructions in system/application prompts. Judge prompts should similarly distinguish candidate answer and evidence from instructions.

### 20.5 Evaluation leakage

Do not include hidden target labels in generation prompts. Deterministic target groups belong to evaluation, not product runtime input.

## 21. Performance and Scale

### 21.1 Bounded stages

Every stage requires a maximum:

- documents per build batch;
- embedding batch and concurrency;
- lexical/vector top-K;
- fused pool;
- rerank pool;
- evidence items and runes/tokens;
- provider calls;
- artifact preview bytes;
- query page size.

### 21.2 Storage

Use:

- filesystem CAS for immutable large artifacts;
- SQLite for local journal, queue, budgets, and query indexes;
- content-addressed bundles for indexes;
- compact journal payload references rather than large inline bodies.

### 21.3 Derived query tables

For large campaigns, add rebuildable indexes such as:

```text
campaign_episode_index
campaign_arm_index
campaign_case_index
stage_candidate_index
chunk_usage_index
observation_index
estimate_index
artifact_edge_index
```

These tables are projections. They can be rebuilt from verified journal and artifact data.

### 21.4 UI virtualization

Use windowed rendering only for genuinely large tables. Prefer server pagination first. A virtualized table must preserve keyboard semantics and expose total counts.

## 22. Failure Semantics

Failures are data, not just logs.

### 22.1 Stage statuses

Use a stable vocabulary:

- `completed`;
- `skipped`;
- `degraded`;
- `failed`;
- `canceled`;
- `missing`.

### 22.2 Reranker degradation

If pool hydration, provider execution, or output validation fails, preserve fused order and record a degraded stage. Do not claim reranking succeeded.

### 22.3 Answer failure

Keep failure classes distinct:

- validation;
- retrieval;
- generation;
- contract;
- policy;
- presentation;
- cancellation.

### 22.4 Judge failure

A failed judge observation is not a zero quality score. Record missing/failed status, diagnostics, and evidence. Paired estimates must declare denominator policy.

### 22.5 Process interruption

Test interruption after:

- lease grant;
- attempt start;
- stage artifact persistence;
- sealed trajectory;
- terminal queue result;
- budget commit;
- episode journal event;
- observation event;
- estimate event.

## 23. Testing Strategy

### 23.1 Unit tests

RagKit:

- chunk ranges and Unicode boundaries;
- representation lineage and identity;
- embedding alignment;
- collapse and fusion determinism;
- hydration fail-closed behavior.

RAG-TTC:

- route validation and policy intersection;
- pre-fusion filtering;
- post-augmentation recheck;
- reranker pool validation and fallback;
- runtime identity stability/change laws;
- answer-contract validation.

Optkit:

- snapshot identity;
- trial expansion and deterministic seeds;
- trajectory ordering and sealing;
- queue idempotency;
- budget reservation/commit laws;
- projector replay.

Judgekit:

- strict instance/protocol/report validation;
- cache mode;
- calibration and reliability math;
- observed model mismatch.

### 23.2 Cross-product fixture tests

The canonical semantic fixture should prove equivalent behavior across:

- RagKit retrieval mechanics;
- RAG-TTC direct retrieval;
- Optkit campaign adapter;
- Coinvault reference if retained;
- query-plane pipeline projection;
- UI rendering smoke test.

A fixture change requires an explicit version/digest update.

### 23.3 Campaign integration tests

Test:

- complete-block expansion;
- exact expected episode count;
- no duplicate semantic execution after resume;
- restart at durability boundaries;
- journal verification;
- projection equality before/after reopen;
- paired estimate correctness;
- hard-constraint invalidation.

### 23.4 UI tests

Use Go handler tests and browser smoke tests for:

- GET-only route surface;
- mutation method denial;
- pagination and cursors;
- sensitivity previews;
- route/back-button behavior;
- keyboard search and drawers;
- pipeline and rank rendering;
- empty/missing/degraded states;
- zero console errors;
- CSP compliance.

### 23.5 Statistical tests

Use synthetic known distributions to verify:

- paired means;
- confidence intervals;
- missing-pair handling;
- group slicing;
- Pareto dominance;
- calibration metrics;
- reliability agreement.

## 24. Proposed API and Package Layout

### 24.1 RAG-TTC product packages

```text
pkg/ttc/optimization/
  config.go             layered configuration graph
  identity.go           semantic identity rules
  invalidation.go       reuse planner
  candidate.go          proposal and hypothesis contract
  promotion.go          promotion manifest

pkg/ttc/optkitcampaign/
  retrieval_system.go   current retrieval adapter
  answer_system.go      direct answer adapter
  judge_instrument.go   Judgekit measurement adapter
  campaign.go           durable orchestration
  reconcile.go          queue/journal reconciliation
  measures.go           deterministic instruments

pkg/ttc/projector/
  campaign.go
  comparison.go
  pipeline.go
  rank_waterfall.go
  chunk.go
  answer.go
  judge.go
  pareto.go
```

### 24.2 Optkit packages

Keep Optkit additions domain neutral:

```text
optkit/
  graph/                optional generic materialization DAG
  projection/           generic projection helpers
  system/               existing factory registry
  experiment/           estimate extensions
```

Do not add RAG-specific stage names to Optkit core.

### 24.3 RagKit additions

Likely additions:

- explicit tokenizer/token count hooks;
- parent/child chunk lineage if required;
- representation-set manifest helpers;
- normalized embedding/index identities;
- optional fusion strategy interface after weighted RRF baseline;
- reusable retrieval diagnostics that remain product-neutral.

## 25. Phased Implementation Plan

### Phase 0: Contract freeze and fixtures

Deliverables:

- freeze layered config schemas;
- define stage schemas and stable names;
- extend semantic fixture for context, answer, and judge lineage;
- document hard constraints and primary metrics.

Acceptance:

- canonical fixture has a version and digest;
- every schema rejects unknown fields;
- current retrieval parity remains green.

### Phase 1: Judgekit instrument over sealed answers

Deliverables:

- answer trajectory decoder;
- recomputed Judgekit instance identity;
- evidence-hidden claim extraction path;
- attributed protocol/report artifacts;
- cache-bypass reliability probes;
- Optkit measurement epoch mapping.

Acceptance:

- one historical answer receives a second judge epoch without retrieval or answer rerun;
- old observations remain unchanged.

### Phase 2: Layered configuration graph

Deliverables:

- typed config references;
- semantic identity tests;
- config diff and invalidation planner;
- promotion manifest skeleton.

Acceptance:

- judge-only change plans no upstream execution;
- chunker change invalidates every downstream stage;
- planner output is deterministic.

### Phase 3: Bundle-build campaigns

Deliverables:

- chunk/representation/embedding build systems;
- build-stage events and usage;
- immutable bundle reuse;
- chunking and embedding metrics.

Acceptance:

- two chunkers build and compare over one frozen corpus;
- interruption resumes without publishing partial bundles.

### Phase 4: Retrieval and reranking search spaces

Deliverables:

- typed route/fusion/reranker candidates;
- reuse of channel and fused artifacts;
- richer estimates and group slices;
- Pareto computation.

Acceptance:

- fusion-only campaign reuses channel outputs;
- reranker-only campaign reuses fused pools;
- first-loss diagnosis remains attributable.

### Phase 5: Answer and context campaigns

Deliverables:

- context-construction artifact;
- direct answer Optkit system;
- deterministic answer instruments;
- answer-level cost/latency metrics.

Acceptance:

- answer prompt change reuses admitted evidence;
- invalid contracts cannot become promotable.

### Phase 6: RAG query projectors

Deliverables:

- episode, pipeline, comparison, chunk, answer, judge, and Pareto projectors;
- database-side pagination and filters;
- entity routes and artifact edges.

Acceptance:

- projectors rebuild from a verified journal;
- no projector mutates source state;
- large synthetic campaign stays bounded.

### Phase 7: Specialist UI

Deliverables:

- cockpit;
- comparison matrix;
- pipeline microscope;
- rank waterfall;
- chunk laboratory;
- answer/evidence studio;
- Judgekit calibration views;
- provenance inspector.

Acceptance:

- all screens work without a frontend build step;
- navigation is deep-linkable and keyboard accessible;
- browser emits no mutation request;
- console and CSP checks pass.

### Phase 8: Optimizer and promotion workflow

Deliverables:

- candidate proposal contract;
- progressive elimination;
- Pareto frontier;
- validation/holdout gates;
- promotion manifest and CLI rendering.

Acceptance:

- optimizer cannot bypass hard constraints;
- every proposed candidate includes a hypothesis and lineage;
- winning manifest links to campaign evidence.

### Phase 9: Scale and production readiness

Deliverables:

- projection indexes;
- event/entity pagination;
- artifact retention policy;
- concurrency/version-conflict retry policy;
- optional bounded canary integration.

Acceptance:

- representative large campaigns remain responsive;
- replay and journal verification remain authoritative.

## 26. Decision Records

### Decision: Keep repository boundaries explicit

- **Context:** RAG algorithms, experiment control, product policy, and LLM measurement have different reuse boundaries.
- **Options considered:** merge all code into Optkit; merge into RAG-TTC; retain separate reusable libraries with product adapters.
- **Decision:** RagKit owns RAG mechanics, Optkit owns domain-neutral science, Judgekit owns LLM measurement, and RAG-TTC composes them.
- **Rationale:** This matches current APIs and prevents domain leakage into Optkit.
- **Consequences:** Cross-module integration and versioning require discipline; adapters are explicit.
- **Status:** accepted.

### Decision: Use layered identities, not one monolithic snapshot

- **Context:** A monolithic snapshot would make a judge change appear to invalidate an index.
- **Options considered:** one full-system digest; layered references; ad hoc cache keys.
- **Decision:** Use content-addressed layered config references and a graph-level root snapshot.
- **Rationale:** This preserves exact identity while enabling principled reuse.
- **Consequences:** The planner and UI must understand graph edges.
- **Status:** proposed.

### Decision: Optimize progressively, not by Cartesian product

- **Context:** Full combinations are expensive and confound diagnosis.
- **Options considered:** exhaustive grid; black-box optimizer; progressive gated search.
- **Decision:** Use cheap mechanical gates, layer-local experiments, paired comparisons, and selective composition.
- **Rationale:** Faster feedback and clearer causal attribution.
- **Consequences:** Some cross-layer interactions may be discovered later and require composition campaigns.
- **Status:** proposed.

### Decision: Keep the specialist UI read-only

- **Context:** Scientific mutations need explicit, reproducible commands and identities.
- **Options considered:** browser editor; mixed read/write UI; read-only explorer with copyable commands.
- **Decision:** CLI/agents author; browser explores and generates copyable command text only.
- **Rationale:** Preserves auditability and keeps the web security surface small.
- **Consequences:** Specialists need a terminal/agent workflow for execution.
- **Status:** accepted.

### Decision: Use server-side domain projectors

- **Context:** Loading all raw events in the browser does not scale and leaks schema complexity into UI code.
- **Options considered:** browser-only folds; generic event API only; versioned domain projectors.
- **Decision:** Add rebuildable, bounded RAG-specific projectors backed by journal/artifacts.
- **Rationale:** Stable UI contracts and bounded performance without creating a second authority.
- **Consequences:** Projector correctness and rebuild tests become critical.
- **Status:** proposed.

### Decision: Treat policy as a hard boundary

- **Context:** Unauthorized candidates must not influence fusion or leave the process.
- **Options considered:** filter after rerank; score penalty; pre-fusion filtering and post-augmentation recheck.
- **Decision:** Enforce before fusion, before reranking, and after augmentation.
- **Rationale:** Current implementation proves this pattern and it fails closed.
- **Consequences:** Policy-negative outcomes are constraints, not averaged quality metrics.
- **Status:** accepted.

### Decision: Keep deterministic and LLM measurements separate

- **Context:** LLM scores can obscure mechanical contract failures and judge instability.
- **Options considered:** one blended score; judge-only evaluation; separate epochs and constructs.
- **Decision:** Retain deterministic contracts separately and add Judgekit observations under attributed epochs.
- **Rationale:** Enables diagnosis, calibration, and remeasurement.
- **Consequences:** Promotion policy must reason over multiple constructs.
- **Status:** accepted.

### Decision: Reuse sealed historical trajectories

- **Context:** Judge and instrument changes should not rerun expensive or stochastic product behavior.
- **Options considered:** rerun full episodes; store only final scores; seal rich trajectories and remeasure.
- **Decision:** Persist enough input, stage, evidence, answer, and usage artifacts for later instruments.
- **Rationale:** Lower cost and exact historical attribution.
- **Consequences:** trajectory schemas and redaction policy are long-lived contracts.
- **Status:** proposed, with retrieval trajectory proof already implemented.

### Decision: Keep UI implementation dependency-free

- **Context:** The current explorer is small, embedded, and read-only.
- **Options considered:** React application; server-rendered templates; semantic HTML/CSS/JS modules.
- **Decision:** Continue with embedded semantic HTML, CSS, and JavaScript without a build step.
- **Rationale:** Minimal deployment and sufficient interaction model.
- **Consequences:** Projectors should keep client logic simple; manual module boundaries are required.
- **Status:** accepted.

## 27. Alternatives Considered

### 27.1 One giant “RAG score”

Rejected because it hides policy failures, missing values, judge weakness, cost, and latency. Keep a vector of outcomes and hard constraints.

### 27.2 Optimize only answer prompts

Rejected because many answer failures originate in chunking, retrieval, fusion, or admission. Prompt tuning cannot recover evidence that never reached the model.

### 27.3 Let an LLM choose arbitrary retrieval parameters per query

Rejected because it breaks treatment identity and policy authority. The model may select only prepared route names.

### 27.4 Send every retrieved candidate to a reranker

Rejected because it leaks unauthorized data, increases cost, and destroys bounded behavior. Filter and bound first.

### 27.5 Use LLM judges without calibration

Rejected because an optimizer can exploit judge biases. Judge protocols need gold calibration and perturbation reliability evidence.

### 27.6 Browser-owned experiment editing

Rejected for the first system because it creates hidden mutation state and duplicates CLI/agent validation. Copyable commands preserve usability without weakening reproducibility.

### 27.7 Move RagKit algorithms into Optkit

Rejected because chunking, representations, embedding, retrieval, and reranking are cohesive RAG-domain concerns. Optkit should remain general purpose.

## 28. Risks and Open Questions

### 28.1 Cross-layer interaction

Layer-local optimization can miss interactions, such as a representation that helps only one embedding model. Mitigation: composition campaigns among Pareto finalists.

### 28.2 Overfitting

Repeated development-set tuning can overfit queries and judges. Mitigation: explicit splits, limited holdout use, minimum effect sizes, and proposal history.

### 28.3 Judge gaming

Candidate answers may exploit judge prompts. Mitigation: deterministic contracts, calibration, reliability probes, multiple constructs, and human audits.

### 28.4 Projection drift

Derived query tables may diverge from journal truth. Mitigation: rebuild tests, versioned projectors, digest checks, and no projector-owned authority.

### 28.5 Artifact volume

Stage artifacts across many arms can become large. Mitigation: bounded pools, deduplicated CAS, retention policy, and materialize only necessary text.

### 28.6 Provider nondeterminism

Embedding, generation, reranking, and judging providers may change behind a model name. Mitigation: capture observed model/revision/settings, use repeats and reliability probes, and avoid pretending fingerprints are signatures.

### 28.7 Lease retry semantics

The current campaign model has limited explicit lease-expiry/retry events. A generalized controller may need a retry transition and conflict policy.

### 28.8 Statistical selection bias

Selecting the best of many candidates inflates observed gains. Mitigation: successive halving, validation confirmation, and reporting candidate count.

### 28.9 UI scope growth

The workbench can become a generic analytics platform. Mitigation: design around specialist decisions—diagnose loss, compare treatments, inspect provenance, and decide next experiment.

### 28.10 Open questions

1. Should materialization DAG support live in Optkit core or remain product-owned until a second domain proves reuse?
2. Which tokenizer is authoritative for context budgets across answer providers?
3. Which chunk hierarchy representation best preserves parent/child lineage without complicating source evidence?
4. What confidence interval method should be the first supported estimate?
5. Which human-labeled set is sufficient to calibrate Judgekit for Tree Center product questions?
6. Which Pareto objectives are primary for the first production promotion policy?
7. How long should large stage candidate artifacts be retained?
8. Should a future UI permit signed approval while keeping campaign mutation elsewhere?

## 29. Intern Onboarding Checklist

### Day 1: architecture

- Run tests in RagKit, Optkit, RAG-TTC, and Judgekit.
- Read the files in Section 30.
- Run the deterministic semantic fixture.
- Run and inspect the durable Optkit RAG campaign.
- Open the current campaign explorer.

### Day 2: trace one query

- Choose one fixture query.
- Write down each retrieval stage and chunk ordering.
- Verify source policy removes the forbidden chunk before fusion.
- Follow admitted evidence into the output artifact.
- Locate the episode trajectory and campaign event.

### Day 3: implement one projector

Recommended first task: `PipelineView` over one sealed retrieval episode.

Steps:

1. load and verify the episode result;
2. load the sealed trajectory;
3. decode only known retrieval-stage schemas;
4. create a bounded stage list;
5. expose it through a GET route;
6. test missing, degraded, and sensitive payloads;
7. render a simple stage rail and table.

### Review questions

Before opening a pull request, answer:

- What is the treatment identity?
- Which artifacts are reused?
- Which stages are recomputed?
- Where is source policy enforced?
- How are missing and degraded outcomes represented?
- Can the journal/projector rebuild after restart?
- Does the browser remain read-only?
- Which tests prove the behavior?

## 30. File Reference Map

### Optkit

- `optkit/system/registry.go:15-96` — product factory/prepared boundary and identity checks.
- `optkit/space/codec.go:1-57` — strict canonical JSON codecs.
- `optkit/space/snapshot.go:12-92` — content-derived snapshot materialization and verification.
- `optkit/experiment/trial.go:9-153` — datasets, arms, complete blocks, episodes, deterministic seeds.
- `optkit/experiment/estimate.go:17-77` — paired estimate and missing-pair behavior.
- `optkit/episode/writer.go:22-164` — ordered trajectory events and sealing.
- `optkit/campaign/event.go:12-83` — append-only campaign event vocabulary.
- `optkit/campaign/reducer.go:108-166` — lifecycle and episode transition rules.
- `optkit/scheduler/queue.go:12-29` — lease/complete/fail/get queue contract.
- `optkit/measure/epoch.go:9-38` — measurement epoch identity.
- `optkit/measure/observation.go:37-137` — typed observations and semantic IDs.
- `optkit/query/service.go:15-281` — current bounded read-only projectors and previews.
- `optkit/internal/web/server.go:1-200` — GET-only standard-library server and SSE.
- `optkit/internal/web/static/app.js:113-330` — current campaign explorer projections in the browser.

### RAG-TTC

- `rag-ttc/pkg/ttc/search/service.go:14-49` — retrieval request, stages, and result.
- `rag-ttc/pkg/ttc/search/service.go:53-197` — service/route preparation and policy intersection.
- `rag-ttc/pkg/ttc/search/service.go:200-282` — retrieval pipeline ordering.
- `rag-ttc/pkg/ttc/search/service.go:347-430` — bounded reranker and degraded fallback.
- `rag-ttc/pkg/ttc/search/identity.go:13-76` — runtime fingerprints.
- `rag-ttc/pkg/ttc/retrievaleval/evaluate.go:11-201` — stage rankings, metrics, treatment checks, and diagnosis.
- `rag-ttc/pkg/ttc/customerapp/service.go:18-190` — canonical direct answer application service.
- `rag-ttc/pkg/ttc/optkitcampaign/system.go:18-193` — typed Optkit retrieval adapter and trajectory emission.
- `rag-ttc/pkg/ttc/optkitcampaign/campaign.go:134-591` — durable campaign execution, reconciliation, observations, and estimates.
- `rag-ttc/pkg/ttc/optkitcampaign/campaign_test.go:17-112` — restart and no-duplicate semantic execution laws.
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/command.go:1-131` — Glazed run/resume/inspect commands.

### RagKit

- `ragkit/rag/components.go:8-91` — RAG interfaces.
- `ragkit/rag/types.go:3-121` — document/chunk/representation/vector/hit/evidence values.
- `ragkit/rag/chunking/markdown_heading.go:12-103` — structure-aware chunking.
- `ragkit/rag/representations/representations.go:1-194` — summaries and synthetic questions with lineage.
- `ragkit/rag/retrieval/retrieval.go:14-194` — collapse, weighted RRF, and bounded hydration.
- `ragkit/rag/indexbundle/types.go:14-177` — immutable bundle identities and build/open stages.

### Judgekit

- `judgekit/eval/instance.go:3-27` — evaluation instance.
- `judgekit/judging/interfaces.go:10-67` — generator, judge, and cache-bypass interfaces.
- `judgekit/protocol/protocol.go:3-67` — reproducible judge protocol.
- `judgekit/assessment/report.go:5-29` — sealed report attribution.
- `judgekit/audit/reliability.go:16-147` — perturbation reliability and cache bypass.
- `judgekit/calibration/report.go:9-154` — calibration metrics and report sealing.

## 31. Validation Commands

From the workspace root, while the unpublished local Optkit dependency remains workspace-bound:

```bash
cd ragkit && GOWORK=off go test ./... -count=1
cd ../judgekit && GOWORK=off go test ./... -count=1
cd ../optkit && go test ./... -count=1
cd ../rag-ttc && go test ./... -count=1
```

Run the durable demonstration:

```bash
cd rag-ttc
go run ./cmd/rag-ttc experiment optkit-rag run \
  --store ./tmp/optkit-rag \
  --reset \
  --format json
```

Resume with the returned campaign ID:

```bash
go run ./cmd/rag-ttc experiment optkit-rag run \
  --store ./tmp/optkit-rag \
  --campaign campaign:<id> \
  --format json
```

Inspect without mutation:

```bash
go run ./cmd/rag-ttc experiment optkit-rag inspect \
  --store ./tmp/optkit-rag \
  --campaign campaign:<id> \
  --format json
```

Run the existing explorer against an Optkit store:

```bash
cd ../optkit
go run ./cmd/optkit serve \
  --store ../rag-ttc/tmp/optkit-rag \
  --listen 127.0.0.1:8080
```

## 32. Final Implementation Guidance

Begin with evidence and projectors, not charts. The first useful UI increment is a server-side pipeline projection over one sealed episode plus a plain table showing stage counts and ranked chunk IDs. Once that projection is correct and bounded, add rank movement, source metadata, target markers, and the waterfall view.

Keep every implementation review grounded in three invariants:

1. **Identity:** can we prove exactly which treatment produced this result?
2. **Boundary:** can unauthorized or unbounded data cross a provider/UI boundary?
3. **Replay:** can the result and its projections be reconstructed after process restart without duplicate semantic execution?

If those invariants remain intact, the workbench can grow from the current durable retrieval slice into a complete optimization environment without becoming an untraceable collection of benchmark scripts and dashboards.
