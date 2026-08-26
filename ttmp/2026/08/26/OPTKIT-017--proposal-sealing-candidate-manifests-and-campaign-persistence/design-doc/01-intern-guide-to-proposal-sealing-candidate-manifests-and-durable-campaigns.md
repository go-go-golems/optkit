---
Title: Intern Guide to Proposal Sealing Candidate Manifests and Durable Campaigns
Ticket: OPTKIT-017
Status: active
Topics:
    - architecture
    - design
    - implementation
    - optkit
    - rag-ttc
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/examples/numbergame/demo.go
      Note: |-
        Existing CandidateProposed and SnapshotMaterialized event precedent
        Candidate event precedent
    - Path: repo://optkit/space/candidate.go
      Note: Candidate identity and structured intent produced by sealing
    - Path: repo://optkit/space/patch.go
      Note: |-
        Canonical assignment artifact and child snapshot mechanics used by sealing
        Canonical sealing mechanics
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/design-doc/04-backend-first-optimization-workbench-program-roadmap.md
      Note: Parent program goals dependencies exclusions and exit gates
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-012--architecture-closure-and-optimization-workbench-contracts/design-doc/01-intern-guide-to-optimization-workbench-architecture-and-contracts.md
      Note: Accepted workbench contract revision b1fcf17a29f89921e9e1c42049de0486a35511f9
    - Path: repo://rag-ttc/pkg/ttc/experimentworkbench/manifest.go
      Note: |-
        Strict full-arm manifest loader extended with candidate authoring
        Strict manifest input boundary
    - Path: repo://rag-ttc/pkg/ttc/optkitcampaign/campaign.go
      Note: |-
        Durable campaign specification and initialization sequence
        Durable spec and journal ordering
ExternalSources: []
Summary: Design for replaying pure drafts through canonical Optkit sealing, strict patch-style candidate manifests, idempotent journal facts, and self-contained campaign persistence.
LastUpdated: 2026-08-26T14:20:26.678499694-04:00
WhatFor: Explain how a reviewed candidate becomes immutable patch/snapshot/candidate records that remain understandable without the source manifest.
WhenToUse: Implement after OPTKIT-016 provides stable pure draft compilation.
---



# Intern Guide to Proposal Sealing, Candidate Manifests, and Durable Campaigns

## 1. Executive summary

OPTKIT-016 creates pure candidate drafts. This ticket performs the irreversible transition from a draft to durable experiment facts. Sealing replays normalized mutations through the registered typed bindings and Optkit's existing `PatchBuilder`, verifies the resulting child against the draft, constructs structured `space.Candidate` intent, persists catalog provenance, and appends candidate/snapshot facts to the campaign.

The same operation supports strict patch-style manifest candidates. A manifest describes a baseline once and challenger deltas as `parent + mutations + intent`. It does not receive a special mutation implementation: manifest loading invokes the same proposal compiler and sealer used by CLI and future browser authoring.

The defining invariant is:

> After campaign creation, deleting the manifest must not remove any fact required to explain or execute the experiment.

## 2. Existing durable mechanics

### 2.1 Patch and snapshot

`PatchBuilder.Build` sorts assignment IDs, applies each typed variable, stores canonical assignment values, materializes the child snapshot, and derives patch identity (`optkit/space/patch.go:104-151`). This is the only mutation-record implementation sealing should use.

### 2.2 Candidate event precedent

Numbergame constructs a `Candidate`, stores `CandidateProposal{Candidate, Patch}`, appends `CandidateProposed`, stores the child snapshot record, and appends `SnapshotMaterialized` (`examples/numbergame/demo.go:103-207`). RAG campaigns currently skip that candidate path and create full arms directly.

### 2.3 Campaign creation

`optkitcampaign.initialize` (`campaign.go:220-338`) materializes one snapshot for each full arm, creates a trial, stores `CampaignSpec`, defines budget, emits lifecycle events, schedules episodes, and starts execution. Candidate facts must be ordered so the campaign creation artifact and journal can be replayed deterministically.

## 3. Draft/seal trust boundary

A client-supplied draft is not trusted merely because it contains graph IDs. Sealing must either recompile from canonical parent + requested mutations or verify every draft contract against current registry and parent.

Recommended approach:

```text
SealProposal(request):
    load and verify parent snapshot
    recompile request.Draft.Mutations with current semantic catalog ID
    compare recomputed draft digest with supplied draft digest
    validate intent
    replay normalized mutations through PatchBuilder
    verify child config/digest equals recomputed preview
    build candidate and durable envelope
    append idempotent events
```

This prevents a caller from altering `ChildConfig`, diff, or plan fields between compile and seal.

## 4. Sealing API

```go
type SealProposalRequest struct {
    Campaign       record.CampaignID `json:"campaign"`
    Parent         record.SnapshotID `json:"parent"`
    DraftDigest    record.Digest     `json:"draft_digest"`
    Mutations      []RequestedMutation `json:"mutations"`
    Intent         CandidateIntent   `json:"intent"`
    IdempotencyKey string            `json:"idempotency_key"`
}

type SealedProposal struct {
    Candidate      space.Candidate      `json:"candidate"`
    Patch          space.PatchRecord    `json:"patch"`
    Child          space.SnapshotRecord `json:"child"`
    Catalog        artifact.Ref         `json:"catalog"`
    ParentGraph    optimization.Graph   `json:"parent_graph"`
    ChildGraph     optimization.Graph   `json:"child_graph"`
    Plan           optimization.InvalidationPlan `json:"invalidation_plan"`
}
```

A manifest-driven campaign may seal before a campaign journal exists. In that flow, split pure materialization from event append internally:

```go
type ProposalMaterializer interface {
    Materialize(ctx, store, verifiedDraft, intent) (SealedProposal, error)
}

type ProposalRecorder interface {
    Record(ctx, campaignID, idempotencyKey, sealed) error
}
```

The public application can orchestrate both. This separation supports campaign initialization without making manifest logic special.

## 5. Strict candidate manifest

Recommended v2 shape:

```yaml
schema: rag-ttc.experiment-manifest/v2
name: rrf-k-comparison
preparation: rag.semantic-fixture/v1
repeats: 1
baseline:
  id: baseline
  pipeline:
    # complete PipelineConfig values
    retrieval:
      preparation: rag.semantic-fixture/v1
      route: default
      final_result_limit: 2
    fusion:
      rrf_k: 60
    # explicit frozen versions for remaining layers
candidates:
- id: rrf-k-20
  parent: baseline
  proposer:
    kind: human
    identity: actor:manuel
  strategy: manual-coordinate/v1
  mutations:
  - variable: fusion.rrf_k
    value: 20
  hypothesis: >-
    Lowering k will increase the influence of early ranks in each channel.
  expected_improvement:
    metric: retrieval.target-coverage
    groups: [multi-source]
  regression_risks:
  - A single noisy high-ranked channel may dominate.
  motivating_evidence:
    case_ids: [q-comparison]
cases:
  # existing cases
```

Use fully qualified variable ID as the canonical mutation key. A redundant section field can drift (`section:fusion`, `variable:fusion.rrf_k`); omit it unless required strictly for readability and validated as derived.

`yaml.Decoder.KnownFields(true)` remains mandatory. Limit size and reject multiple documents as the v1 loader does (`manifest.go:63-94`).

## 6. Manifest resolution

```text
ResolveManifest(v2):
    strict-decode and validate top-level metadata
    validate baseline PipelineConfig
    baselineGraph = DeriveGraph(baseline)
    build current RAG registry
    for candidate in declaration order:
        resolve parent (baseline or previously resolved candidate if chaining accepted)
        compile parent + mutations
        reject diagnostics/non-sealable draft
        validate candidate intent
        retain resolved arm draft
    validate cases/repeats/preparation
    derive manifest semantic ID from reviewed authoring values
    return LoadedManifest with baseline and resolved candidate drafts
```

### Candidate chaining

Decide whether `parent` may name only the baseline or any earlier candidate. Baseline-only is simpler and matches comparisons. If chaining is accepted:

- reject cycles;
- require parents appear earlier or perform explicit topological sort;
- preserve deterministic order;
- persist the full ancestry.

Do not support chaining accidentally through map lookup.

## 7. Materialization algorithm

```text
MaterializeProposal(parent, draft, intent):
    registry semantic ID must equal draft catalog ID
    full catalog bytes → artifact store
    builder = NewPatchBuilder(parent, store, PipelineConfigCodec)
    for mutation in canonical order:
        registry.Binding(mutation.Variable).Assign(builder, mutation.After)
    patch, child = builder.Build(ctx)
    assert child.Value canonical digest == draft.ChildConfig digest
    childGraph = DeriveGraph(child.Value)
    assert childGraph.ID == draft.AfterGraph.ID
    candidate = NewCandidate(
        parent.ID, patch.ID, child.ID,
        intent + semantic catalog ID,
    )
    store candidate proposal envelope
    return all records/refs
```

If verification fails, do not append events. Artifacts are immutable and may be left unreferenced by a failed operation depending on store transaction capabilities; document garbage-collection implications.

## 8. Event ordering and idempotency

Recommended event order after `CampaignCreated`, or as accepted by campaign transition rules:

```text
CandidateProposed(candidate envelope)
SnapshotMaterialized(child snapshot record)
TrialPlanned / episodes scheduled
```

Review `campaign.ValidateTransition`; candidate facts are informational and currently allowed in active states. Campaign initialization may need to create the campaign before appending candidate facts, while the creation spec already contains candidate records. That duplication is acceptable only if one is canonical payload and the other is a journal index/fact pointing at it.

Idempotency:

- derive/stash a command ID from caller idempotency key;
- retry returns the same candidate when payload is identical;
- same key with different payload is a conflict;
- existing candidate ID/event prevents duplicate append;
- journal expected version handles races.

## 9. Durable campaign specification

Target concepts:

```go
type Arm struct {
    ID          string               `json:"id"`
    Description string               `json:"description,omitempty"`
    Snapshot    space.SnapshotRecord `json:"snapshot"`
    Graph       optimization.Graph   `json:"graph"`
}

type CandidateRecord struct {
    ArmID       string               `json:"arm_id"`
    Candidate   space.Candidate      `json:"candidate"`
    Patch       space.PatchRecord    `json:"patch"`
    Catalog     artifact.Ref         `json:"catalog"`
    Mutations   []NormalizedMutation `json:"mutations"`
}

type CampaignSpec struct {
    System     record.SystemID              `json:"system"`
    Description string                      `json:"description,omitempty"`
    Trial      experiment.TrialPlan         `json:"trial"`
    Arms       []Arm                        `json:"arms"`
    Candidates map[string]CandidateRecord   `json:"candidates,omitempty"`
    Cases      []RetrievalCase              `json:"cases"`
    Budget     []budget.Limit               `json:"budget"`
    ManifestID string                       `json:"manifest_id,omitempty"`
}
```

Snapshots reference content-addressed config artifacts. Candidate records retain proposal meaning. Graphs are sealed derivations used by read projections; read code does not derive them from the current library.

## 10. Catalog provenance

Persist the exact full catalog artifact used at seal. Candidate semantic identity uses semantic catalog ID under OPTKIT-012. Historical UI loads the sealed descriptor/copy rather than current registry documentation.

If the catalog artifact is restricted because descriptors reveal sensitive coordinate metadata, assign sensitivity deliberately. Artifact variable content itself remains in its own artifact; the catalog contains schema/policy metadata, not prompt text.

## 11. Store-canonical test

A decisive integration test:

```text
copy fixture manifest to temp path
campaign run --manifest temp --store fresh
capture campaign ID
remove temp manifest
campaign status/verify/resume (if incomplete)
project cockpit/comparison/provenance
assert hypothesis, risks, expected metric, mutations,
       parent/child snapshots, graph plan, and catalog docs available
```

Store this shell test in `scripts/` if not implemented as Go integration tests.

## 12. Compatibility policy

Follow OPTKIT-012. Without explicit v1 support:

- v2 loader handles only v2 shape;
- v1 fixtures remain tests for old code/history, not silently converted;
- fresh v2 stores are required;
- no deprecated duplicate `arms + candidates` ambiguity.

If v1 manifests remain accepted, route by schema and compile v1 full arms into internal resolved arms explicitly. If old stores remain readable, dispatch campaign/snapshot schemas without recomputing identities. Resume support is a separate decision from read-only projection support.

## 13. Implementation phases

1. Finalize v2 manifest structs and strict validation.
2. Resolve baseline and candidate drafts through `CompileProposal`.
3. Implement materialization through bindings and `PatchBuilder`.
4. Extend candidate/campaign durable schemas.
5. Add journal recording and idempotency.
6. Integrate `campaign dry-run/run`; dry-run compiles but does not seal.
7. Add manifest-removal, retry, restart, and verification tests.

## 14. Testing strategy

### Manifest

- unknown fields, multiple YAML documents, oversize input;
- duplicate IDs and unknown parents;
- malformed mutation values;
- invalid/empty intent;
- deterministic manifest ID;
- optional candidate chaining cycle/order rules.

### Sealing

- replay equals draft child config/graph;
- stale parent/draft/catalog rejected;
- candidate identity stable across retry;
- same idempotency key/different payload conflict;
- no event on validation/materialization failure;
- canonical assignment order.

### Persistence

- all facts survive manifest deletion;
- restart/resume uses stored spec;
- `Metadata.Verify` and artifact verification pass;
- comparison can locate candidate for treatment arm;
- old-store behavior matches accepted policy.

### Commands

```bash
cd rag-ttc
GOWORK=off go test ./pkg/ttc/experimentworkbench \
  ./pkg/ttc/optkitcampaign ./pkg/ttc/specialistapi -count=1
GOWORK=off go test ./... -count=1
GOWORK=off go run ./cmd/rag-ttc experiment optkit-rag campaign dry-run \
  --manifest <v2-fixture> --store /tmp/optkit-017
GOWORK=off go run ./cmd/rag-ttc experiment optkit-rag campaign run \
  --manifest <v2-fixture> --store /tmp/optkit-017 --reset
```

## 15. Review risks

- Trusting client draft outputs instead of recompiling canonical inputs.
- Appending candidate facts before campaign state permits them.
- Candidate in spec and event payload drifting; define canonical relationship.
- Retry behavior relying only on candidate ID while journal commands differ.
- Recomputing historical graphs/catalog descriptions with current code.
- Leaving unreferenced artifacts after failed sealing without documenting cleanup.
- Candidate chaining increasing complexity beyond current product needs.

## 16. Out of scope

- Candidate comparison UI and APIs (OPTKIT-018/019);
- interactive draft HTTP;
- asset prompt materialization details (OPTKIT-020);
- gate policies or promotion;
- broad manifest language redesign unrelated to candidates.

## 17. Exit criteria

- v2 candidate manifest is strict and concise.
- manifest and CLI sealing share `CompileProposal`/`SealProposal`.
- sealing uses typed bindings and `PatchBuilder`.
- candidate, patch, child, catalog, graph, and intent are durable.
- retries are deterministic/idempotent.
- source manifest can be removed without losing execution or explanation.
- focused/full tests, fresh-store smoke, diary, doctor, and upload pass.

## 18. File reference map

- `optkit/space/patch.go:43-151` — canonical durable mutation.
- `optkit/space/candidate.go:10-60` — candidate creation/identity.
- `optkit/examples/numbergame/demo.go:103-207` — proposal/snapshot events.
- `optkit/campaign/event.go:20-30` — event kinds including `CandidateProposed`.
- `rag-ttc/pkg/ttc/experimentworkbench/manifest.go:31-174` — current v1 strict manifest.
- `rag-ttc/pkg/ttc/experimentworkbench/service.go:37-88` — dry-run/run split.
- `rag-ttc/pkg/ttc/optkitcampaign/campaign.go:68-80,220-338` — current spec and initialization.
- `rag-ttc/pkg/ttc/optkitcampaign/campaign.go:624-689` — stored spec loading and journal append helpers.
