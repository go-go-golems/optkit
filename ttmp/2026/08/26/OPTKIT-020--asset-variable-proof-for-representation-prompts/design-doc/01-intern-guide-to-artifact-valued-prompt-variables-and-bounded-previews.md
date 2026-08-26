---
Title: Intern Guide to Artifact-Valued Prompt Variables and Bounded Previews
Ticket: OPTKIT-020
Status: active
Topics:
    - architecture
    - design
    - implementation
    - optkit
    - rag-ttc
    - ui
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/artifact
      Note: Content-addressed sensitivity-aware artifact contracts
    - Path: repo://optkit/artifact/ref.go
      Note: Content-addressed artifact reference and sensitivity contract
    - Path: repo://optkit/space/patch.go
      Note: |-
        Assignment records store artifact refs rather than embedding values
        Assignment artifact-ref behavior
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/design-doc/04-backend-first-optimization-workbench-program-roadmap.md
      Note: Parent program goals dependencies exclusions and exit gates
    - Path: repo://rag-ttc/apps/specialist/web/src/components/Artifact.tsx
      Note: |-
        Existing schema-keyed artifact renderer fallback
        Current artifact preview fallback
    - Path: repo://rag-ttc/apps/specialist/web/src/layerwidgets/index.tsx
      Note: Existing output renderer registry extended by the workbench plugin
    - Path: repo://rag-ttc/pkg/ttc/search/types.go
      Note: |-
        Recorded representation summaries and chunk catalog used by previews
        Recorded representation summaries
ExternalSources: []
Summary: Design for sensitivity-aware artifact-valued representation prompt variables, content materialization, truthful recomputation, side-by-side editing, and bounded one-chunk previews.
LastUpdated: 2026-08-26T14:20:30.531858526-04:00
WhatFor: Prove the scalar workbench architecture generalizes to expensive content-addressed configuration without embedding prompt text in patch metadata or faking instant generation.
WhenToUse: Implement after OPTKIT-019 completes the scalar RRF vertical slice and OPTKIT-018 provides preview authorization.
---




# Intern Guide to Artifact-Valued Prompt Variables and Bounded Previews

## 1. Executive summary

A representation prompt is not an ordinary scalar. Its value is text that may be sensitive, large enough to deserve independent provenance, and expensive in effect: changing it can regenerate representations, embeddings, indexes, retrieval outputs, and every downstream result. The workbench must still treat it as the same kind of optimization coordinate—registered, compiled, planned, sealed, run, and compared—but it cannot store raw prompt text inside patch metadata or simulate generation in TypeScript.

This ticket registers `representations.summary_prompt` as an artifact-reference variable. The editor creates or references a new prompt artifact, displays an authorized old/new diff, and compiles the ref mutation to obtain the real graph invalidation plan. A bounded server preview regenerates one selected chunk representation and records enough provenance to explain exactly what ran.

This is the final abstraction proof: scalar and asset coordinates use one catalog/compiler/sealer/persistence/workbench flow, differing only in value materialization and preview capability.

## 2. Artifact orientation

Optkit artifacts are immutable bytes addressed by digest, with media type, optional schema, and sensitivity. `space.AssignmentRecord` stores `Value artifact.Ref` (`space/patch.go:13-17`), so the patch already supports artifact-valued encoded values.

For a scalar variable, the assignment value artifact contains canonical JSON such as `20`. For an asset variable, the variable's canonical value is itself an `artifact.Ref`; prompt text is stored in a separate prompt artifact. The assignment value artifact therefore records the ref, not the text.

```text
prompt text bytes
  → prompt artifact Ref(digest, schema, sensitivity)
  → canonical encoded artifact.Ref value
  → assignment value artifact
  → PatchRecord.Assignment.Value
```

This extra level is intentional: patches remain small and generic.

## 3. Representation orientation

The current retrieval record can expose representations associated with touched chunks. `search.RepresentationSummary` and `ChunkSummary.Representations` (`search/types.go:46-72`) contain representation ID, kind, and text for observability. That is output evidence, not necessarily the producer that creates summaries.

Before implementing the variable, locate the actual representation-generation pipeline/configuration used by the selected fixture or production preparation. If no producer can accept a prompt artifact and generate one representation deterministically enough for a bounded preview, the ticket is blocked; do not fake it by editing recorded representation text.

## 4. Semantic configuration

Replace the representations frozen value with a typed config when the producer is available:

```go
const SummaryPromptSchema record.SchemaID =
    "schema:rag-ttc.prompt.representation-summary/v1"

type RepresentationsConfig struct {
    SummaryPrompt artifact.Ref `json:"summary_prompt" yaml:"summary_prompt"`
    Producer      string       `json:"producer" yaml:"producer"`
    Model         string       `json:"model,omitempty" yaml:"model,omitempty"`
}
```

`Producer`/`Model` participate in layer identity if they affect output. Do not claim a prompt-only mutation fully identifies representation semantics while model/template code remains hidden.

Validation:

- prompt ref digest and schema valid;
- schema equals accepted prompt schema;
- sensitivity allowed by environment/policy;
- producer non-empty and registered;
- model required if producer semantics need it.

## 5. Variable contract

```go
func SummaryPromptVariable() space.Variable[
    optimization.PipelineConfig,
    artifact.Ref,
] {
    codec := space.NewJSONCodec[artifact.Ref](SummaryPromptRefValueSchema)
    return space.Variable[optimization.PipelineConfig,artifact.Ref]{
        Descriptor: space.VariableDescriptor{
            ID: "representations.summary_prompt",
            Key: "summary_prompt",
            Label: "Summary representation prompt",
            Short: "Prompt used to derive searchable summary text from each chunk.",
            Long: "Changing this artifact regenerates summary representations and invalidates embeddings, indexes, retrieval, fusion, and downstream results...",
            ValueSchema: codec.Schema(),
            Value: space.ArtifactRefSpec(SummaryPromptSchema),
            Sensitive: true,
            CostHint: "regenerates-corpus",
            BindingVersion: "rag-ttc.representations.summary-prompt/v1",
            Probes: []string{"representations.one-chunk/v1"},
        },
        Lens: liftedSummaryPromptLens,
        Domain: summaryPromptRefDomain,
        Codec: codec,
    }
}
```

The domain validates references. Creating a new artifact from text belongs to authoring/sealing applications, not `space.Domain`.

## 6. Draft asset lifecycle

Interactive editing needs a way to refer to unsealed text. Avoid writing a durable artifact on every keystroke.

Recommended command model:

```go
type DraftAsset struct {
    MediaType   string               `json:"media_type"`
    Schema      record.SchemaID      `json:"schema"`
    Sensitivity artifact.Sensitivity `json:"sensitivity"`
    Content     string               `json:"content"`
}

type RequestedMutation struct {
    Variable   space.VariableID `json:"variable"`
    Value      json.RawMessage  `json:"value,omitempty"`
    DraftAsset *DraftAsset      `json:"draft_asset,omitempty"`
}
```

Compile validates metadata/content limits and computes a provisional content digest/ref without storing it. The child preview can use that provisional ref because content addressing is deterministic. Seal materializes bytes, verifies the actual ref equals provisional ref, then assigns the ref through the binding.

Alternative: create expiring draft artifacts. This adds cleanup/auth complexity and durable writes during editing, so prefer deterministic provisional refs if the artifact store's digest algorithm is available without writing.

### Decision: No per-keystroke durable artifacts

- **Context:** Draft/seal separation forbids uncontrolled history while editing.
- **Options considered:** write every edit; browser-only opaque temp IDs; deterministic provisional refs; expiring draft store.
- **Decision:** Compute provisional content-addressed refs in memory and materialize only on preview/seal, subject to request limits.
- **Rationale:** It preserves pure compilation and lets sealing verify exact bytes.
- **Consequences:** Preview requests carry content or a securely stored draft token; enforce size/sensitivity/logging controls.
- **Status:** proposed.

## 7. Recompute plan

Changing `RepresentationsConfig.SummaryPrompt` changes the representations local identity. With the current graph topology:

```text
reuse: corpus, chunking
recompute direct: representations
recompute upstream-change: embeddings, indexes, retrieval, fusion,
                           reranking, evidence, context, answer, judge
```

Do not hardcode this list in React. `optimization.Plan` derives it from before/after graphs and the UI renders the returned steps.

## 8. Bounded preview

The honest preview is:

> Generate the summary representation for one selected chunk using the proposed prompt and configured producer.

It is not:

- regenerate the corpus;
- predict retrieval metrics;
- update indexes;
- run a fake client-side summarizer.

### Request

```go
type RepresentationPreviewRequest struct {
    DraftDigest record.Digest `json:"draft_digest"`
    CaseID      string        `json:"case_id,omitempty"`
    ChunkID     string        `json:"chunk_id"`
    Prompt      DraftAssetRef `json:"prompt"`
}
```

### Result

```go
type RepresentationPreview struct {
    Schema          record.SchemaID `json:"schema"`
    ChunkID         string          `json:"chunk_id"`
    SourceDigest    record.Digest   `json:"source_digest"`
    PromptDigest    record.Digest   `json:"prompt_digest"`
    Producer        string          `json:"producer"`
    Model           string          `json:"model,omitempty"`
    OutputText      string          `json:"output_text"`
    OutputDigest    record.Digest   `json:"output_digest"`
    StartedAt       time.Time       `json:"started_at"`
    FinishedAt      time.Time       `json:"finished_at"`
    Sensitivity     artifact.Sensitivity `json:"sensitivity"`
}
```

Preview output is not campaign evidence unless explicitly sealed/run. Mark it as a preview and bind it to draft digest.

### Execution pseudocode

```text
PreviewOneChunk(request, principal):
    authorize preview.run and restricted artifact reads
    recompile/verify draft
    verify representations.summary_prompt changed
    load exact chunk by ID from parent corpus revision
    materialize or securely resolve draft prompt bytes
    prepare registered representation producer
    run one generation with context deadline and token/resource limits
    compute output digest and provenance
    return sensitivity-protected preview
```

## 9. Security and sensitivity

- Prompt content must not appear in URLs, logs, diagnostics, or analytics.
- Request body limits protect server/rendering paths.
- Authentication determines proposer; authorization gates prompt/chunk reads and preview execution.
- Result sensitivity is at least the maximum of prompt and chunk sensitivity.
- Browser caching and service workers must not persist restricted prompt/preview payloads casually.
- Clipboard/download actions are explicit and authorized if provided.
- Error messages name digest/operation, not raw content.

## 10. Sealing asset mutations

```text
SealProposal with DraftAsset:
    validate/recompile draft and intent
    promptRef = artifactStore.Put(exact bytes, schema, sensitivity)
    assert promptRef digest == provisional ref in compiled mutation
    binding.Assign(PatchBuilder, canonical promptRef)
    Build patch/child snapshot
    persist candidate/catalog/provenance
    append events idempotently
```

If artifact materialization succeeds but journal append fails, retry uses content addressing and returns the same ref. The idempotency contract from OPTKIT-017 handles event recording.

## 11. React editor

The specialized asset editor renders:

- parent prompt digest and authorized text;
- new text editor;
- side-by-side or unified diff;
- schema/sensitivity badges;
- character/byte limit;
- recomputation strip from compiler plan;
- one-chunk selector and preview action;
- preview status/output/provenance;
- explicit statement that preview is one chunk, not full trial evidence.

Use a generic artifact-ref fallback when content preview is unavailable. Missing authorization should explain why text/diff cannot be shown while still rendering safe metadata if policy allows.

## 12. Producer discovery work

Before code:

1. inventory representation creation in RAG-TTC/Ragkit;
2. identify input chunk type and output `rag.Representation` shape;
3. identify current summary prompt source;
4. identify model/provider and determinism controls;
5. identify content/artifact stores and sensitivity;
6. write a one-chunk command-line experiment under this ticket's `scripts/`;
7. record cost/latency/error output in diary.

If the current semantic fixture contains precomputed representation text only, add a dedicated deterministic test producer or integrate a real bounded producer according to accepted scope. Do not mutate fixture outputs and call that generation.

## 13. Implementation phases

1. Producer/config inventory and bounded experiment.
2. Typed `RepresentationsConfig` and graph derivation.
3. Artifact-ref domain/variable/registry entry.
4. Provisional draft asset compile contract.
5. Preview service with auth, limits, provenance, and cancellation.
6. Seal-time materialization and idempotency.
7. React diff editor/preview plugin.
8. Scalar-versus-asset conformance and full campaign proof.

## 14. Testing strategy

### Artifact/config

- correct schema/digest/sensitivity accepted;
- wrong schema or malformed ref rejected;
- provisional digest equals stored artifact ref;
- prompt content absent from patch/candidate JSON except artifact refs;
- graph plan changes representations/downstream only.

### Preview

- authorized success;
- forbidden prompt/chunk;
- unknown chunk;
- stale draft;
- unsupported producer;
- timeout/cancellation;
- output/request size limits;
- sensitivity propagation;
- no durable candidate/journal event.

### Seal

- exact previewed bytes can seal;
- content changed after compile rejected by digest mismatch;
- retry stable;
- candidate/provenance survive restart;
- full trial uses new prompt producer/config.

### Frontend

- authorized diff and metadata-only fallback;
- no content in URL/error snapshots;
- preview scope/cost clear;
- stale preview invalidated after edits;
- keyboard-accessible editor/diff/selector;
- scalar and asset variables share shell/seal path.

### Commands

```bash
cd optkit && GOWORK=off go test ./space ./artifact -count=1
cd rag-ttc && GOWORK=off go test ./pkg/ttc/optimization \
  ./pkg/ttc/experimentworkbench ./pkg/ttc/optkitcampaign -count=1
cd rag-ttc/apps/specialist/web && pnpm typecheck && pnpm test && pnpm build
```

Run the one-chunk preview smoke through the relevant ticket script and record exact producer/model/environment.

## 15. Risks and review focus

- No actual representation producer may exist behind the fixture; this is a real blocker, not permission to fake one.
- Draft prompt content can leak through logs, browser caches, or diagnostics.
- Provisional digest computation must exactly match artifact store canonical bytes/metadata.
- Model/provider changes affect representation identity beyond the prompt.
- Preview output is not trial evidence.
- Full recomputation cost may be large; budget and cancellation must be explicit.
- Artifact ref validation cannot prove caller authorization by itself.

## 16. Out of scope

- arbitrary file/asset variables beyond the proof coordinate;
- full-corpus instant previews;
- client-side LLM generation;
- prompt optimization/search algorithms;
- promotion/approval workflow;
- moving representation algorithms to another repository.

## 17. Exit criteria

- real producer accepts prompt artifact through semantic config;
- prompt variable is an artifact ref with complete catalog metadata;
- compile remains non-durable;
- plan truthfully shows representations/downstream recomputation;
- one-chunk preview is bounded, authorized, cancellable, and provenance-rich;
- seal stores exact content/ref and uses normal patch mechanics;
- React provides diff/preview/fallback safely;
- scalar and asset conformance, full tests, live smoke, diary, doctor, and upload pass.

## 18. File reference map

- `optkit/artifact/` — immutable store/ref/sensitivity implementation to inspect before code.
- `optkit/space/patch.go:13-24,43-76` — assignment values are artifact refs and new values are stored canonically.
- `optkit/space/snapshot.go:24-60` — child config artifact and snapshot identity.
- `rag-ttc/pkg/ttc/search/types.go:46-72` — recorded representation/chunk summaries.
- `rag-ttc/pkg/ttc/search/semantic_fixture.go:85-119` — fixture preparation/precomputed inputs.
- `rag-ttc/apps/specialist/web/src/components/Artifact.tsx:1-50` — current preview availability/fallback.
- `rag-ttc/apps/specialist/web/src/layerwidgets/index.tsx:19-29` — existing schema renderer registry.
