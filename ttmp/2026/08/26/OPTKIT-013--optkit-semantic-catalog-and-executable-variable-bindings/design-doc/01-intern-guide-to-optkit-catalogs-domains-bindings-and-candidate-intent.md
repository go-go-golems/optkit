---
Title: Intern Guide to Optkit Catalogs Domains Bindings and Candidate Intent
Ticket: OPTKIT-013
Status: active
Topics:
    - architecture
    - design
    - implementation
    - optkit
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/examples/numbergame/demo.go
      Note: |-
        Complete generic proof case for variables patches candidates and campaign events
        Complete generic proof consumer
    - Path: repo://optkit/space/candidate.go
      Note: Candidate intent and identity model to extend
    - Path: repo://optkit/space/domain.go
      Note: |-
        Existing integer and lossy choice domains to extend
        Lossy current domain descriptors to replace
    - Path: repo://optkit/space/patch.go
      Note: |-
        Canonical durable typed assignment path that bindings must invoke
        Durable typed assignment path bindings must preserve
    - Path: repo://optkit/space/variable.go
      Note: |-
        Existing typed variable and descriptor implementation
        Typed variable declaration extended by catalog metadata
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/design-doc/04-backend-first-optimization-workbench-program-roadmap.md
      Note: Parent program goals dependencies exclusions and exit gates
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-012--architecture-closure-and-optimization-workbench-contracts/design-doc/01-intern-guide-to-optimization-workbench-architecture-and-contracts.md
      Note: Accepted workbench contract revision b1fcf17a29f89921e9e1c42049de0486a35511f9
ExternalSources: []
Summary: Implementation design for ordered semantic catalogs, complete value descriptors, serialized-to-typed executable bindings, richer candidate intent, and the numbergame proof.
LastUpdated: 2026-08-26T14:20:21.132009032-04:00
WhatFor: Teach a new Optkit contributor how to make variables discoverable and remotely invocable without creating a parallel mutation mechanism.
WhenToUse: Implement after OPTKIT-012 contracts are accepted and before RAG-TTC registers application variables.
---



# Intern Guide to Optkit Catalogs, Domains, Bindings, and Candidate Intent

## 1. Executive summary

Optkit already knows how to mutate a typed configuration safely. A `Variable[C,V]` combines a lens, a domain, and a codec; `PatchBuilder` records a canonical assignment and materializes a child snapshot. What Optkit cannot yet do is enumerate variables as an ordered, documented catalog or receive a serialized value and route it back into the correct typed variable.

This ticket adds those capabilities while preserving the existing mutation algebra. It introduces serializable sections and value specifications, executable type-erased bindings constructed from typed variables, deterministic catalog identities, and structured candidate intent. Numbergame proves the generic mechanism before RAG-TTC depends on it.

The key invariant is:

> Every registered descriptor has exactly one executable binding, and both come from the same typed variable declaration.

## 2. Concepts an intern must understand

### 2.1 Configuration type `C`

`C` is the complete value being optimized. In numbergame it is `numbergame.Config`. In the later RAG work it is `optimization.PipelineConfig`. Optkit remains generic over `C`.

### 2.2 Value type `V`

`V` is one coordinate's Go type: `int`, `float64`, `bool`, `string`, or `artifact.Ref`. The compiler cannot erase this type until a binding has captured the codec, domain, and lens.

### 2.3 Lens

`space.Lens[C,V]` (`lens.go:5-8`) reads and immutably writes a `V` inside a `C`. `CheckLensLaws` verifies get-put, put-get, and put-put. These laws matter because proposal compilation and sealing apply the same variable through different execution paths and must reach the same child configuration.

### 2.4 Codec

`space.Codec[V]` (`codec.go:11-15`) is the canonical transport boundary. A binding must decode with the variable's codec and re-encode canonically. It must not use loose `map[string]any` conversion because JSON numbers and unknown fields can change meaning.

### 2.5 Domain

A domain states which values are legal. The existing `IntRangeDomain` is adequate for integers. The existing `ChoiceDomain.Descriptor` loses machine values and exposes only sorted labels. This ticket fixes the serialization without weakening typed validation.

### 2.6 Patch

`PatchBuilder.Build` stores assignment values and a child snapshot. Bindings use it only during sealing. Pure draft application uses the same lens/domain/codec but does not call `Build`.

## 3. Current state and observed gaps

`space.VariableDescriptor` currently contains ID, name, one description, value schema, a loose domain descriptor, sensitivity, and probes (`variable.go:21-29`). The executable `Variable[C,V]` adds lens/domain/codec (`variable.go:31-36`). Validation checks ID, name, lens, non-nil domain/codec, and descriptor/codec schema equality.

The current descriptor cannot support the workbench completely:

- no section/order information;
- no separate machine key and human label;
- no Short/Long documentation;
- no float, boolean, string, or artifact domain shape;
- choices expose labels but not canonical values;
- no default semantics;
- no binding lookup by variable ID;
- no catalog content/version identity.

The current `Candidate` already treats hypothesis, targets, and risks as semantic identity. It lacks structured expected improvement and motivating evidence.

## 4. Proposed package layout

Keep the work in `optkit/space` unless implementation review finds a package-cycle or cohesion reason to split:

```text
optkit/space/
  section.go          SectionID, Section, catalog ordering
  valuespec.go        discriminated serializable value specification
  catalog.go          Catalog construction, identity, JSON
  binding.go          Binding[C], Bindings[C], typed adapter
  variable.go         existing typed variable + enriched descriptor
  domain.go           typed domains + lossless descriptors
  candidate.go        structured intent and identity
  *_test.go           contracts and numbergame-independent tests
```

Catalogs are part of the coordinate space, so `space` is a coherent home. Avoid importing Glazed schema types: Optkit needs domains, canonical values, sensitivity, artifact schema, and binding versions that CLI flag definitions do not provide.

## 5. Descriptor and catalog API

```go
type SectionID string

type Section struct {
    ID        SectionID            `json:"id"`
    Label     string               `json:"label"`
    Short     string               `json:"short"`
    Long      string               `json:"long"`
    Variables []VariableDescriptor `json:"variables"`
}

type VariableDescriptor struct {
    ID             VariableID      `json:"id"`
    Key            string          `json:"key"`
    Label          string          `json:"label"`
    Short          string          `json:"short"`
    Long           string          `json:"long"`
    ValueSchema    record.SchemaID `json:"value_schema"`
    Value          ValueSpec       `json:"value"`
    Default        json.RawMessage `json:"default,omitempty"`
    Sensitive      bool            `json:"sensitive"`
    CostHint       string          `json:"cost_hint,omitempty"`
    BindingVersion string          `json:"binding_version"`
    Probes         []string        `json:"probes,omitempty"`
}

type Catalog struct {
    Schema     record.SchemaID `json:"schema"`
    SemanticID record.Digest   `json:"semantic_id"`
    ID         record.Digest   `json:"id"`
    Sections   []Section       `json:"sections"`
}
```

Validation rules:

- section IDs and variable IDs use stable machine syntax;
- section IDs are unique and preserve registration order;
- variable IDs are globally unique;
- every variable belongs to exactly one section;
- labels, Short, and Long documentation are non-empty;
- default values decode canonically and pass the typed domain;
- artifact kinds name an artifact schema;
- non-artifact kinds do not name an artifact schema;
- cost hints are advisory and from a documented finite vocabulary;
- semantic IDs ignore labels/docs but include legality and binding version;
- full catalog ID includes exact docs and ordering.

## 6. Value and domain design

Use a discriminated value shape. Do not infer the kind from whichever optional field happens to be present.

```go
type ValueSpec struct {
    Kind           ValueKind        `json:"kind"`
    IntegerRange   *IntegerRangeSpec `json:"integer_range,omitempty"`
    FloatRange     *FloatRangeSpec   `json:"float_range,omitempty"`
    Choices        []ChoiceSpec      `json:"choices,omitempty"`
    String         *StringSpec       `json:"string,omitempty"`
    ArtifactSchema *record.SchemaID  `json:"artifact_schema,omitempty"`
}

type ChoiceSpec struct {
    Value json.RawMessage `json:"value"`
    Label string          `json:"label"`
}
```

Typed domains still implement:

```go
type Domain[V any] interface {
    Validate(V) error
    Descriptor(Codec[V]) (ValueSpec, error)
}
```

Passing the codec when deriving a descriptor ensures choice machine values use the same canonical encoding as assignments. An alternative is to store canonical values in the domain constructor. Choose one and test it; do not format generic values with `fmt.Sprint`.

### Float range invariants

- reject NaN and positive/negative infinity;
- require finite minimum/maximum and `minimum <= maximum`;
- specify inclusive bounds;
- do not add an arbitrary slider step to semantic metadata unless the domain truly quantizes values.

### Artifact domain invariants

An artifact-valued variable's `V` should be `artifact.Ref`, not prompt text. Validate schema and allowed sensitivity separately from content creation. The binding mutates refs; an application service materializes new content before sealing.

## 7. Executable binding design

A binding performs four related operations:

1. describe the variable;
2. normalize serialized input with the real codec;
3. apply it purely to a `C` for drafting;
4. add it to `PatchBuilder[C]` for sealing.

```go
type Binding[C any] interface {
    ID() VariableID
    Descriptor() VariableDescriptor
    ReadCanonical(C) (json.RawMessage, error)
    Normalize(json.RawMessage) (json.RawMessage, error)
    ApplyPure(C, json.RawMessage) (C, error)
    Assign(*PatchBuilder[C], json.RawMessage) error
}

type Registry[C any] struct {
    catalog Catalog
    byID    map[VariableID]Binding[C]
}
```

Construction captures the type parameter:

```go
func Register[C, V any](
    builder *RegistryBuilder[C],
    section SectionID,
    variable Variable[C,V],
) error {
    if err := variable.Validate(); err != nil { return err }
    binding := typedBinding[C,V]{variable: variable}
    return builder.add(section, binding)
}
```

### Pure application pseudocode

```text
ApplyPure(config, raw):
    value = variable.Codec.Decode(raw)
    canonical = variable.Codec.EncodeCanonical(value)
    reject trailing or unknown data through codec
    variable.Domain.Validate(value)
    updated = variable.Lens.Put(config, value)
    return updated, canonical
```

### Durable assignment pseudocode

```text
Assign(builder, raw):
    value = decode with the same codec
    validate with the same domain
    return space.Set(builder, variable, value)
```

The compiler and sealer therefore share decoding and legality. Only the final side effect differs.

## 8. Catalog construction and identity

A builder preserves section and variable insertion order for rendering, but semantic hashing must be deterministic. Require callers to register in desired display order and reject duplicate IDs. Canonical JSON already stabilizes object keys; slice order remains intentional.

```text
BuildCatalog():
    validate every section and variable
    semanticProjection = strip labels, Short, Long
    semanticID = SemanticDigest(catalog-semantic/v1, semanticProjection)
    fullID = SemanticDigest(catalog/v1, complete catalog)
    return immutable catalog + bindings map
```

Do not expose mutable internal slices/maps. Return detached copies from catalog accessors.

## 9. Candidate intent extension

Recommended fields:

```go
type Proposer struct {
    Kind     ProposerKind    `json:"kind"`
    Identity record.ActorRef `json:"identity"`
}

type ExpectedImprovement struct {
    Metric string   `json:"metric"`
    Groups []string `json:"groups,omitempty"`
}

type Motivation struct {
    CaseIDs          []string      `json:"case_ids,omitempty"`
    DiagnosticDigest record.Digest `json:"diagnostic_digest,omitempty"`
}

type Candidate struct {
    // parent, patch, child
    Proposer            Proposer            `json:"proposer"`
    Strategy            string              `json:"strategy"`
    Hypothesis          string              `json:"hypothesis"`
    ExpectedImprovement ExpectedImprovement `json:"expected_improvement"`
    Risks               []string            `json:"risks,omitempty"`
    Motivation          Motivation          `json:"motivation,omitempty"`
    SemanticCatalogID   record.Digest       `json:"semantic_catalog_id"`
    CreatedAt           time.Time           `json:"created_at"`
}
```

Identity includes all semantic intent and excludes `CreatedAt`. Normalize whitespace only where the API explicitly defines it; silently trimming prose before hashing can make the stored text differ from reviewed text. Reject empty required fields instead.

`Targets` migration follows the compatibility policy accepted in OPTKIT-012. Do not retain both old and new fields as an undocumented dual source of truth.

## 10. Numbergame proof

Numbergame is the correct first consumer because `examples/numbergame/demo.go:103-193` already demonstrates:

- baseline snapshot materialization;
- typed multiplier patch;
- child snapshot;
- candidate with hypothesis and risks;
- `CandidateProposed` event.

Add a catalog:

```go
func Registry() (*space.Registry[Config], error) {
    b := space.NewRegistryBuilder[Config]()
    b.AddSection(space.Section{
        ID: "numbergame", Label: "Number game",
        Short: "Controls the deterministic arithmetic system.",
        Long: "The numbergame system multiplies each input and may add seeded noise...",
    })
    space.Register(b, "numbergame", MultiplierVariable())
    space.Register(b, "numbergame", NoiseVariable())
    return b.Build()
}
```

Proof test:

```text
raw mutation {variable:"math.multiplier", value:3}
  → lookup binding
  → ApplyPure(base config)
  → child preview Config{Multiplier:3}
  → Assign(PatchBuilder)
  → Build
  → exact same child config and snapshot identity as direct space.Set
```

The demo should record the semantic/full catalog provenance required by the accepted architecture, without introducing RAG concepts.

## 11. Implementation phases

### Phase 1 — Value specifications and domains

Implement and test discriminated descriptors first. This settles the public JSON shape before catalog and UI consumers exist.

Files:

- `space/valuespec.go` new;
- `space/domain.go` extend;
- `space/domain_test.go` new or focused additions.

### Phase 2 — Sections and catalog

Implement validation, order, detached accessors, semantic/full identity, and deterministic JSON.

### Phase 3 — Bindings

Implement private generic adapters and registry construction. Test pure and durable paths against the same typed variable.

### Phase 4 — Candidate intent

Update the candidate constructor/schema and identity tests according to OPTKIT-012. Avoid compatibility aliases unless explicitly accepted.

### Phase 5 — Numbergame migration

Declare one section and its variables, apply a serialized mutation in a test or demo path, and preserve the existing complete campaign behavior.

### Phase 6 — Full validation

Run focused tests, `go test ./...`, race tests for registry access if shared concurrently, and the numbergame demo/export parity path.

## 12. Testing strategy

### Domain tests

- min/max accepted; just-outside rejected;
- malformed domain rejected;
- float NaN/infinity rejected;
- choices round-trip machine values and labels;
- artifact schema required only for artifact refs;
- default values pass codec and domain.

### Catalog tests

- duplicate section/variable IDs rejected;
- registration order preserved;
- returned slices cannot mutate catalog internals;
- serialization deterministic across runs;
- doc-only edit changes full ID but not semantic ID;
- domain or binding-version edit changes both IDs.

### Binding tests

- unknown variable;
- wrong JSON type;
- trailing JSON;
- out-of-domain value;
- canonical numeric/string representation;
- pure apply equals durable child config;
- duplicate patch assignment remains rejected by `PatchBuilder`.

### Candidate tests

- hypothesis/risk/expected metric/evidence/proposer changes candidate ID;
- timestamp changes do not;
- collection normalization follows the accepted contract;
- required intent fields fail clearly.

### Commands

```bash
cd optkit
GOWORK=off go test ./space ./examples/numbergame -count=1
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./space ./examples/numbergame -count=1
```

## 13. Review risks

- A descriptor and binding can still drift if callers can construct one without the other. Keep constructors narrow.
- Generic JSON number handling can turn integers into floats if `any` is used. Decode through `Codec[V]` only.
- Documentation identity and semantic identity must remain separate and named clearly.
- Registry accessors must not expose mutable slices/maps if a singleton is shared by handlers.
- Candidate schema changes alter new identities. Record the schema/version change, not just the struct diff.
- Asset sensitivity is a policy input; a boolean may eventually be insufficient, but do not generalize beyond the accepted artifact contract in this ticket.

## 14. Out of scope

- RAG `PipelineConfig` and variables (OPTKIT-014/015);
- proposal compilation orchestration (OPTKIT-016);
- manifest syntax and campaign persistence (OPTKIT-017);
- HTTP and React component choices;
- gate-policy unification;
- moving optimization abstractions to Ragkit.

## 15. Exit criteria

- Catalog and value-spec JSON are stable and tested.
- Every registered descriptor has one executable binding.
- Pure application and durable patching agree.
- Candidate identity behavior is explicit.
- Numbergame proves the complete generic path.
- Full Optkit tests pass.
- Diary, changelog, relations, doctor, and reMarkable delivery are complete.

## 16. File reference map

- `optkit/space/variable.go:10-57` — current ID, descriptor, variable, validation.
- `optkit/space/domain.go:8-73` — current descriptor, int range, choice loss.
- `optkit/space/lens.go:5-45` — lenses and laws.
- `optkit/space/codec.go:11-55` — strict canonical codecs.
- `optkit/space/patch.go:31-151` — typed assignment and canonical patch order.
- `optkit/space/snapshot.go:24-83` — snapshot identity and verification.
- `optkit/space/candidate.go:10-60` — current candidate identity.
- `optkit/space/space_test.go:15-96` — current variable and patch tests.
- `optkit/examples/numbergame/demo.go:103-193` — reference candidate loop.
