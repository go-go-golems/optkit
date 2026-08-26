---
Title: Intern Guide to the RAG-TTC Experiment Workbench CLI
Ticket: OPTKIT-006
Status: active
Topics:
    - implementation
    - optkit
    - rag-ttc
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/query/service.go
      Note: Read-only bounded projection precedent
    - Path: repo://rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/command.go
      Note: Current Glazed durable campaign surface
    - Path: repo://rag-ttc/pkg/ttc/optimization/graph.go
      Note: Layered identity resolver
    - Path: repo://rag-ttc/pkg/ttc/optimization/invalidation.go
      Note: Configuration diff and invalidation authority
    - Path: repo://rag-ttc/pkg/ttc/optkitcampaign/campaign.go
      Note: Durable campaign application service
ExternalSources: []
Summary: Architecture, contracts, commands, implementation phases, and validation plan for the manifest-driven RAG-TTC experiment workbench CLI.
LastUpdated: 2026-08-25T23:58:00Z
WhatFor: Teach a new engineer how experiment definitions become validated configuration graphs, durable Optkit campaigns, and structured operator output before the specialist UI is built.
WhenToUse: Read before changing the RAG-TTC Optkit campaign CLI, experiment manifest, layered configuration model, or future projector API.
---


# Intern Guide to the RAG-TTC Experiment Workbench CLI

## Executive summary

RAG-TTC already contains the core of a durable experiment system. Optkit provides a content-addressed artifact store, an append-only campaign journal, leases, budget reservations, immutable system snapshots, complete-block trial expansion, observations, and estimates. RAG-TTC adapts its retrieval runtime to those contracts. A small Glazed command can start or resume one fixed deterministic campaign.

What the repository does not yet provide is an operator-facing workbench that turns a reviewed experiment definition into an explainable execution. Today, `rag-ttc experiment optkit-rag run` silently chooses the embedded semantic fixture, two hard-coded retrieval arms, one repeat, and the semantic-fixture executor. An operator cannot save that plan, validate it without running, compare it with another plan, or see which configuration layers would be reused. The future browser would therefore risk becoming the first place where those semantics are assembled and interpreted.

This ticket makes the CLI the first consumer of the same scientific facts that the specialist UI will display. It introduces a versioned experiment manifest, a small application service that loads and validates it, and Glazed commands for four jobs:

1. **Understand configuration** — inspect, validate, compare, and derive an invalidation plan.
2. **Prepare execution** — render a dry-run containing exact graph, arm, case, episode, and budget counts.
3. **Operate campaigns** — start, resume, inspect, and verify a durable campaign.
4. **Export facts** — emit stable rows through Glazed in table, JSON, JSONL, CSV, TSV, or YAML form.

The workbench belongs in RAG-TTC because its manifest names RAG layers, retrieval preparations, cases, target groups, and evaluation semantics. Optkit remains domain-neutral and continues to own durable mechanisms. The CLI remains a thin adapter: parsing and row emission happen in `cmd/rag-ttc`, while manifest validation, planning, and execution happen in `pkg/ttc` packages that HTTP projectors can later reuse.

The first implementation deliberately supports the deterministic semantic-fixture preparation. That is enough to prove arbitrary arm limits, selected cases, repeat counts, graph comparison, restart behavior, budgets, and structured output without network credentials. Future runtime preparations can register behind the same manifest and executor boundary without changing the durable campaign model.

## 1. The problem in operator terms

An experiment operator needs to answer five questions before spending work:

- **What exactly will run?** Which system, data, cases, arms, repeats, and configuration graph?
- **What changed from the baseline?** Which local configuration values differ?
- **What must be recomputed?** Which downstream products are invalidated by those direct changes?
- **Can execution resume safely?** Is the campaign journal valid, and will completed work be reused?
- **Can another consumer trust the result?** Are emitted identities, counts, budgets, observations, and estimates derived from the same authoritative records that the UI will use?

The existing CLI only answers part of the fourth question. It can start or resume a fixed campaign and emit a one-row summary. The fixed plan is selected in code, so the invocation itself is not a complete scientific record.

A good workbench must separate three operations that are often accidentally conflated:

```text
manifest validation     graph planning          durable execution
-------------------     --------------          -----------------
Is the input legal?     What changes?           What happened?
Are refs complete?      What is reusable?       What was measured?
Are cases coherent?     What is invalidated?    Is the journal valid?
No store mutation       No store mutation       Journal + CAS writes
```

That separation is central to this design. `config validate`, `config diff`, `config plan`, and `campaign dry-run` must be side-effect free. Only `campaign run` creates a campaign. `campaign resume` may advance an existing campaign but may not replace its stored specification with a local file.

## 2. Repository map for a new engineer

The system spans three ownership layers.

### 2.1 Optkit: domain-neutral experiment control

Optkit is the durable substrate. The important packages are:

- `optkit/artifact` — immutable content-addressed blobs and typed references;
- `optkit/record` — validated IDs, canonical semantic digests, and content-derived IDs;
- `optkit/space` — immutable system configuration snapshots;
- `optkit/experiment` — arms, datasets, complete-block trials, episode expansion, and estimators;
- `optkit/campaign` — commands, append-only control events, lifecycle transitions, and journal verification;
- `optkit/scheduler` — work items, leases, attempts, completion, failure, and expiry reclaim;
- `optkit/budget` — atomic reservation and actual-usage commitment;
- `optkit/episode` — append-only trajectory construction and sealed episode results;
- `optkit/measure` — measurement epochs and observations;
- `optkit/projection` and `optkit/query` — read-only generic campaign views.

`optkit/query/service.go:19-24` defines current read limits: 200 events by default, 500 maximum, and 256 KiB payload previews. `query.Service` verifies and folds journal data, but it intentionally does not understand RAG stages or retrieval cases.

The standalone `optkit` binary is an integrity and demonstration tool. `optkit/cmd/optkit/main.go:34-70` dispatches demo, serve, campaign inspection/verification, and artifact verification. It is not the place for RAG experiment manifests.

### 2.2 RAG-TTC application services: domain semantics

The durable RAG campaign adapter is in:

```text
rag-ttc/pkg/ttc/optkitcampaign/
```

Key contracts:

- `system.go` defines `RetrievalConfig`, `RetrievalCase`, `Executor`, and the Optkit system factory;
- `fixture.go` adapts the embedded semantic fixture and constructs the current two-arm plan;
- `campaign.go` initializes, executes, reconciles, resumes, verifies, and summarizes a campaign.

`RunOptions` in `campaign.go:48-59` is the current application boundary. It accepts a store root, reset flag, optional campaign ID, arms, cases, repeats, an executor, clock, lease duration, and checkpoint hook.

`CampaignSpec` in `campaign.go:61-68` is persisted as the `CampaignCreated` payload. This is important: once created, the durable campaign owns its exact trial, snapshots, arms, cases, and budgets. A resume operation loads that stored specification rather than trusting the caller's current manifest.

`Run` in `campaign.go:134-205` implements the high-level state machine:

```text
open profile
  -> create ID or use requested campaign
  -> read journal
  -> initialize when journal is empty
     OR load immutable stored specification
  -> register process-local executor
  -> reclaim expired leases
  -> execute ready work
  -> reconcile completed work into usage, observations, and estimates
  -> verify journal
  -> summarize
```

The implementation is already restart-safe. The workbench should expose it, not replace it.

### 2.3 RAG-TTC CLI: parsing and structured output

The Glazed commands live under:

```text
rag-ttc/cmd/rag-ttc/cmds/
```

The root binary initializes logging and Cobra once in `cmd/rag-ttc/main.go`. The experiment group is registered in `cmds/experiments/root.go`. The current durable campaign command is `cmds/experiments/optkitrag/command.go`.

That command already follows the core Glazed pattern:

- `cmds.CommandDescription` declares fields;
- a settings struct uses `glazed` tags;
- `RunIntoGlazeProcessor` decodes values and calls an application service;
- `types.NewRow` emits structured facts;
- `cli.BuildCobraCommandFromCommand` injects the three universal output flags.

The universal output flags are only:

```text
--format table|json|jsonl|csv|tsv|yaml
--output-fields field1,field2,...
--max-output-rows N
```

A domain limit such as case pagination or event count must remain a command field. `--max-output-rows` caps serialization; it does not alter source work.

## 3. Current behavior and evidence

### 3.1 What is already good

The current `run` command is small and delegates correctly. At `optkitrag/command.go:39-59`, it decodes settings, constructs `RunOptions`, selects the semantic fixture only for a new campaign, invokes `optkitcampaign.Run`, and emits the summary.

The application layer protects durable identity:

- new campaigns require at least two arms and one case (`campaign.go:208-210`);
- arm IDs and case IDs are unique;
- each arm becomes an immutable Optkit snapshot;
- each case becomes a CAS artifact;
- the trial is complete-block, so every arm is evaluated on every case and repeat;
- scheduler work uses semantic episode keys;
- budgets are reserved before enqueueing;
- completions are reconciled idempotently into actual usage and observations;
- campaign completion is recorded only after all expanded episodes have valid rows.

The current command also gets Glazed output flags automatically (`optkitrag/command.go:109-118`) and has tests proving the flags and processor-error propagation.

### 3.2 Current gaps

The command hard-codes these scientific choices:

```go
Executor: optkitcampaign.SemanticFixtureExecutor{}
arms, cases := optkitcampaign.SemanticFixturePlan()
Repeats = 1
```

The effects are:

- no reviewable input file;
- no dry-run;
- no way to select or parameterize arms;
- no way to select cases or repeat count;
- no explicit layered graph attached to the operator plan;
- no graph validation or comparison command;
- no invalidation explanation;
- `run --campaign` ambiguously means resume;
- `--reset` and `--campaign` are not separated by command intent;
- `inspect` emits only one summary row and does not explicitly report integrity status;
- no command verifies all journal payload artifacts;
- no output contract designed for later HTTP parity.

The CLI should not solve these gaps by embedding orchestration in command methods. That would make HTTP and tests duplicate it.

## 4. Proposed command tree

The first workbench remains under the existing public path to avoid an unnecessary command migration:

```text
rag-ttc experiment optkit-rag
├── config
│   ├── inspect
│   ├── validate
│   ├── diff
│   └── plan
└── campaign
    ├── dry-run
    ├── run
    ├── resume
    ├── status
    └── verify
```

The old top-level `run` and `inspect` commands may remain temporary aliases only if existing tests or automation require them. The preferred interface is explicit: creating and resuming are different commands.

### 4.1 Configuration commands

#### `config inspect`

Loads one manifest, validates it, resolves the graph, and emits one row per layer. Intended fields:

```text
manifest_schema
manifest_id
graph_id
layer
value_schema
config_identity
resolved_digest
depends_on
```

This is the operator's equivalent of the future config inspector.

#### `config validate`

Loads one manifest and emits a compact result row only after every validation gate passes:

```text
manifest_schema
manifest_id
graph_id
preparation
arms
cases
repeats
episodes
status=valid
```

Invalid input returns an error and no success row.

#### `config diff`

Loads baseline and challenger manifests, resolves both graphs, calls `optimization.Diff`, and emits one row per canonical layer:

```text
before_graph
after_graph
layer
before_identity
after_identity
changed
```

#### `config plan`

Calls `optimization.Plan` and emits one row per canonical layer:

```text
before_graph
after_graph
layer
action
reason
before_resolved
after_resolved
```

The CLI does not independently infer invalidation.

### 4.2 Campaign commands

#### `campaign dry-run`

Loads and validates the manifest without opening an Optkit profile. It emits exact counts and conservative budget ceilings:

```text
manifest_id
graph_id
arms
cases
repeats
episodes
retrieval_query_budget
retrieval_result_budget
store
mutation=false
```

#### `campaign run`

Creates a campaign from a manifest. Inputs:

```text
--manifest PATH
--store PATH
--reset
```

`--reset` is permitted only for a new run and retains Optkit's safe-path guard. The command emits the durable campaign summary after execution.

#### `campaign resume`

Resumes an existing campaign:

```text
--campaign campaign:...
--store PATH
```

It does not accept a manifest or reset flag. The stored `CampaignSpec` is authoritative. The process still selects an executor compatible with the stored preparation; the first version supports only the embedded semantic fixture.

#### `campaign status`

Verifies and rebuilds the campaign summary without advancing work. It reports:

- campaign and trial IDs;
- lifecycle status and version;
- scheduled/running/completed/failed counts;
- budget reservation/commit state;
- arm estimates and first paired delta;
- integrity state;
- store root.

#### `campaign verify`

Verifies the hash-chained journal and every unique artifact directly referenced by a journal event. It emits a success row with event count and unique payload count. Later versions can traverse known nested artifacts through schema-specific custody walkers.

## 5. Experiment manifest contract

### 5.1 Goals

The manifest is a reviewed operator input, not the authoritative run record. Its responsibilities are:

- identify its own schema;
- name an execution preparation;
- define at least two arms;
- define at least one evaluation case;
- define repeat count;
- carry one complete layered configuration graph;
- be strictly decoded;
- derive a deterministic semantic ID.

The journal and CAS remain authoritative for what actually ran.

### 5.2 Proposed v1 shape

```yaml
schema: rag-ttc.experiment-manifest/v1
name: semantic-limit-comparison
preparation: rag.semantic-fixture/v1
repeats: 1

arms:
  - id: limit-1
    retrieval:
      preparation: rag.semantic-fixture/v1
      route: default
      limit: 1
  - id: limit-2
    retrieval:
      preparation: rag.semantic-fixture/v1
      route: default
      limit: 2

cases:
  - id: placement-basics
    mode: positive
    query: Where should the tree be planted?
    required_groups:
      - id: answer
        targets: [chunk:placement]
  - id: private-notes-denied
    mode: authorization_negative
    query: Show me private staff notes
    forbidden_targets: [chunk:private-notes]

layers:
  - schema: rag-ttc.layered-config-ref/v2
    layer: corpus
    value_schema: rag-ttc.config.corpus/v1
    identity: config:...
  # all twelve canonical layers, with direct dependency identities
```

For the first implementation, the checked-in example may derive its `layers` from the frozen optimization fixture. Operators can copy and modify complete manifests. A later authoring command can construct layer refs from typed value documents; it is not required to validate and execute v1.

### 5.3 Deterministic identity

After strict decoding and normalization:

```text
manifest_digest = SemanticDigest(
  schema:rag-ttc.experiment-manifest/v1,
  normalized_manifest_without_derived_id
)
manifest_id = ContentID("experiment-manifest", manifest_digest)
```

Path, file modification time, YAML comments, map insertion order, and output format do not affect identity.

### 5.4 Validation invariants

Manifest validation is layered:

1. **Envelope**
   - exact schema;
   - non-empty name;
   - supported preparation;
   - repeats greater than zero.
2. **Arms**
   - at least two;
   - unique non-empty IDs;
   - valid retrieval config;
   - arm preparation equals manifest preparation.
3. **Cases**
   - at least one;
   - unique IDs;
   - each case passes `RetrievalCase.Validate`;
   - positive cases have required groups;
   - authorization-negative cases have forbidden targets.
4. **Graph**
   - exactly twelve canonical layers;
   - exact layer-ref schema;
   - unique layer and identity;
   - valid value schema IDs;
   - dependencies exist and point upstream;
   - graph resolves through `optimization.NewGraph`.
5. **Execution compatibility**
   - preparation has a registered executor factory;
   - first version allows `rag.semantic-fixture/v1` only.

Strict decoding must reject unknown fields. A typo such as `repeat: 3` must fail instead of silently producing one repeat.

## 6. Application architecture

### 6.1 Package structure

```text
rag-ttc/pkg/ttc/experimentworkbench/
├── manifest.go       # schema, strict load, identity, validation
├── plan.go           # dry-run and graph comparison services
├── execution.go      # preparation registry and RunOptions construction
├── verify.go         # journal + direct payload custody check
└── testdata/
    ├── semantic-limit-v1.yaml
    └── semantic-limit-challenger-v1.yaml

rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/
├── root.go
├── config.go
├── campaign.go
├── rows.go
└── command_test.go
```

Domain code must not import Cobra, Glazed, or command packages. Command code may import application packages.

### 6.2 Service boundary

A minimal service API is:

```go
type Manifest struct { /* versioned operator input */ }

type LoadedManifest struct {
    Path     string
    ID       string
    Manifest Manifest
    Graph    optimization.Graph
}

func LoadManifest(ctx context.Context, path string) (LoadedManifest, error)
func ValidateManifest(manifest Manifest) (optimization.Graph, error)
func DryRun(loaded LoadedManifest, store string) DryRun
func Compare(before, after LoadedManifest) (optimization.ConfigDiff, error)
func Invalidation(before, after LoadedManifest) (optimization.InvalidationPlan, error)

func NewRunOptions(loaded LoadedManifest, store string, reset bool) (optkitcampaign.RunOptions, error)
func Resume(ctx context.Context, store string, campaign record.CampaignID) (optkitcampaign.Summary, error)
func Status(ctx context.Context, store string, campaign record.CampaignID) (optkitcampaign.Summary, error)
func Verify(ctx context.Context, store string, campaign record.CampaignID) (Verification, error)
```

The exact names may evolve during implementation, but the ownership boundary should not.

### 6.3 Preparation registry

The manifest names a preparation; it does not serialize process-local clients or open handles.

```go
type Preparation interface {
    Name() string
    Executor(context.Context, LoadedManifest) (optkitcampaign.Executor, error)
}

registry = {
    "rag.semantic-fixture/v1": SemanticFixturePreparation{},
    // future local-index or connected-runtime preparations
}
```

The first version may implement this as a closed switch rather than a general plugin registry. The stable contract is the preparation name and error behavior, not a premature extension framework.

## 7. End-to-end flows

### 7.1 Validate a manifest

```text
operator
  |
  | rag-ttc experiment optkit-rag config validate --manifest experiment.yaml
  v
Glazed command
  |
  | DecodeSectionInto(settings)
  v
experimentworkbench.LoadManifest
  |
  +--> read bounded file
  +--> strict YAML decode
  +--> validate envelope, arms, cases
  +--> optimization.NewGraph(layers)
  +--> derive semantic manifest ID
  v
LoadedManifest
  |
  | emit one structured row
  v
table/json/jsonl/csv/tsv/yaml
```

Pseudocode:

```go
func validateCommand(ctx, settings, processor) error {
    loaded, err := workbench.LoadManifest(ctx, settings.Manifest)
    if err != nil { return err }

    return processor.AddRow(ctx, row(
        "manifest_id", loaded.ID,
        "graph_id", loaded.Graph.ID,
        "arms", len(loaded.Manifest.Arms),
        "cases", len(loaded.Manifest.Cases),
        "episodes", loaded.EpisodeCount(),
        "status", "valid",
    ))
}
```

### 7.2 Compare two configurations

```text
baseline YAML --> strict load --> Graph A --+
                                           +--> optimization.Diff/Plan --> rows
challenger YAML -> strict load --> Graph B -+
```

The workbench does not compare YAML text. Comments, formatting, and field order are irrelevant. `Diff` compares local semantic identities. `Plan` compares resolved dependency digests and distinguishes:

- `direct_change` — the layer's own semantic identity changed;
- `upstream_change` — local identity stayed the same but an input changed;
- `unchanged` — resolved meaning is identical and output is reusable.

### 7.3 Create a campaign

```text
manifest
  -> strict load + graph resolution
  -> preparation lookup
  -> RunOptions{arms, cases, repeats, executor}
  -> optkitcampaign.Run
     -> immutable snapshots
     -> case artifacts
     -> complete-block trial
     -> budgets + reservations
     -> work queue + leases
     -> trajectories + results
     -> usage + observations
     -> paired estimates
     -> journal verification
  -> structured campaign summary
```

### 7.4 Resume a campaign

```text
campaign ID + store
  -> validate ID
  -> select supported process-local executor
  -> optkitcampaign.Run{Campaign: id}
     -> load stored CampaignSpec
     -> reclaim expired leases
     -> execute only ready work
     -> reconcile idempotently
     -> verify
  -> summary
```

Resume does not use a local manifest. If future preparations cannot be inferred from stored snapshots, add an explicit `--preparation` selector that must match stored configuration; never allow a manifest to replace stored arms or cases.

## 8. Structured output contract

A command emits rows, not human-formatted prose. This makes the same operation usable by a terminal operator, shell script, test, and future API characterization.

### 8.1 Row cardinality

| Command | Row cardinality |
|---|---:|
| `config validate` | exactly one |
| `config inspect` | one per layer |
| `config diff` | one per layer |
| `config plan` | one per layer |
| `campaign dry-run` | exactly one |
| `campaign run` | exactly one summary |
| `campaign resume` | exactly one summary |
| `campaign status` | exactly one summary |
| `campaign verify` | exactly one result |

This is intentionally boring. It gives CSV and tables coherent records while JSON remains machine-readable.

### 8.2 Error behavior

Errors propagate through Cobra `RunE` to the root. Reusable commands must not call `os.Exit`, `cobra.CheckErr`, or close the builder-owned processor. A validation failure emits no success row.

Useful errors include context:

```text
load manifest experiments/foo.yaml: unknown field "repeat"
validate arm "limit-zero": TTC retrieval limit must be positive
resolve graph: layer "answer" depends on non-upstream layer "judge"
resume campaign campaign:...: journal verification failed at sequence 19
```

Sensitive payloads and credentials must never be included in errors.

## 9. Scientific and safety invariants

### 9.1 Manifest is intent; journal is fact

The manifest explains requested work. The campaign journal and artifact store prove completed work. Status and verification commands read the latter.

### 9.2 Resume never changes the experiment

A campaign's persisted `CampaignSpec` controls arms, cases, trial, snapshots, and budgets after creation. Resume may only continue missing work and reconcile completed results.

### 9.3 Graph plans do not claim materialized reuse

An invalidation plan says which layer outputs are semantically reusable if available and verified. It does not prove that a corresponding artifact exists. The future execution planner must intersect semantic reuse with artifact custody.

### 9.4 Missing is not zero

Campaign and measurement rows must preserve absent, failed, and inapplicable states. A missing paired estimate is emitted as null/absent, never `0`.

### 9.5 Dry-run is side-effect free

A test should pass a nonexistent store path to `campaign dry-run` and prove the path remains absent.

### 9.6 Reset is explicit and guarded

Only new campaign execution accepts `--reset`. `optkitcampaign.resetRoot` refuses unsafe paths. Resume and read-only commands never delete a store.

### 9.7 No provider secrets in manifests

The manifest may name a profile or preparation in future versions, but provider tokens remain environment or credential-store inputs. They do not participate directly in a committed manifest.

## 10. Design decisions

### Decision: Keep the workbench in RAG-TTC

- **Context:** The command names RAG layers, retrieval configurations, target groups, answer/judge lineage, and product fixtures.
- **Options considered:** Put all commands in Optkit; put generic mechanics in Optkit and semantic commands in RAG-TTC; create a separate binary.
- **Decision:** Keep the workbench under `rag-ttc experiment optkit-rag` and call Optkit libraries.
- **Rationale:** Optkit remains reusable while RAG-TTC composes retrieval and evaluation semantics.
- **Consequences:** A second product may later motivate generic command helpers, but semantic manifests remain product-owned.
- **Status:** accepted.

### Decision: Use a versioned manifest rather than flags alone

- **Context:** Arms, cases, layer dependencies, and repeats are too structured to review or reproduce as a long flag list.
- **Options considered:** Flags only; arbitrary Go templates; YAML/JSON manifest plus small operational flags.
- **Decision:** Use strict YAML/JSON-compatible manifests for scientific intent and flags for paths, IDs, reset, and output.
- **Rationale:** Manifests are reviewable, hashable, diffable, and suitable for durable attribution.
- **Consequences:** Schema versioning and migration become explicit responsibilities.
- **Status:** accepted.

### Decision: Separate create and resume

- **Context:** The old `run --campaign` changes meaning based on whether an ID is present.
- **Options considered:** Retain overloaded run; use `run --resume`; create separate commands.
- **Decision:** `campaign run` creates from a manifest; `campaign resume` continues a stored campaign.
- **Rationale:** Different mutation and authority rules deserve distinct interfaces.
- **Consequences:** Existing shorthand may need a temporary compatibility alias.
- **Status:** accepted.

### Decision: Make planning read-only

- **Context:** Operators and CI need to inspect intent before creating a store or spending provider budget.
- **Options considered:** Create a draft campaign for every plan; calculate locally without persistence.
- **Decision:** Validation, diff, plan, and dry-run are pure application operations.
- **Rationale:** Review does not need mutation, and dry-run should be safe in CI.
- **Consequences:** A dry-run is not an authoritative campaign fact until `run` stores the plan.
- **Status:** accepted.

### Decision: Emit normalized rows through Glazed

- **Context:** Humans want tables while scripts and future characterization tests need structured records.
- **Options considered:** Handwritten text; JSON only; Glazed rows.
- **Decision:** Every command is a `GlazeCommand` and emits stable rows.
- **Rationale:** One source supports table, JSON, JSONL, CSV, TSV, YAML, projection, and row caps.
- **Consequences:** Nested records should be represented carefully so tabular output remains useful.
- **Status:** accepted.

### Decision: Support one offline preparation first

- **Context:** Real connected runtimes add credentials, provider budgets, caches, and mutable external dependencies.
- **Options considered:** Generalize all runtimes immediately; prove the workbench with the deterministic fixture; ship planning only.
- **Decision:** Fully implement the semantic-fixture preparation first and leave an explicit preparation boundary.
- **Rationale:** It proves contracts, durability, output, and UI parity without network flakiness.
- **Consequences:** The first CLI is useful for scientific validation but not yet a replacement for all legacy answer-quality commands.
- **Status:** accepted.

### Decision: Verify direct journal payload custody in v1

- **Context:** Journal verification proves the event hash chain but not every nested artifact referenced inside payload objects.
- **Options considered:** Journal only; direct event payload verification; full transitive schema-aware traversal.
- **Decision:** Verify the journal and every unique direct event payload in v1; document nested traversal as a later projector/custody feature.
- **Rationale:** It is deterministic and materially stronger without inventing an unsafe generic JSON reference crawler.
- **Consequences:** “verified” output must state its scope.
- **Status:** accepted.

## 11. Implementation phases

### Phase 0: Architecture and contract audit

- Map current command, campaign, fixture, graph, query, and Glazed boundaries.
- Record current hard-coded behavior and durable invariants.
- Write this intern guide and diary.
- Freeze command names, row cardinality, and first preparation scope.

**Exit gate:** every proposed command maps to an existing or explicitly new application service.

### Phase 1: Versioned manifest

- Add `pkg/ttc/experimentworkbench`.
- Define strict schema and loader.
- Validate envelope, arms, cases, graph, and preparation compatibility.
- Derive deterministic manifest identity.
- Add baseline and challenger example manifests.
- Add malformed, duplicate, unknown-field, invalid-graph, and deterministic-ID tests.

**Exit gate:** valid manifests resolve to stable graph and manifest IDs; malformed inputs fail before mutation.

### Phase 2: Configuration commands

- Split the command package into root/config/campaign/rows files.
- Add `config inspect`, `validate`, `diff`, and `plan`.
- Use `cli.BuildCobraCommandFromCommand`.
- Test settings decode, universal flags, emitted rows, processor errors, and deterministic output.
- Characterize table and JSON output from the built binary.

**Exit gate:** CLI plan rows exactly match `optimization.Plan` for judge-, answer-, reranker-, fusion-, and chunking-change fixtures.

### Phase 3: Campaign commands

- Add side-effect-free `campaign dry-run`.
- Add manifest-driven `campaign run`.
- Add manifest-free `campaign resume`.
- Add `campaign status` and scoped `campaign verify`.
- Preserve existing fixed command compatibility only if needed.
- Validate IDs before opening stores.

**Exit gate:** a checked-in manifest creates, completes, inspects, verifies, and resumes the same campaign without duplicate episodes or observations.

### Phase 4: End-to-end characterization

- Build the binary with `GOWORK=off`.
- Validate all help-visible fields and only three universal output flags.
- Run baseline and challenger graph comparisons.
- Run an offline campaign in a temporary store.
- Capture JSONL or JSON output as ticket evidence.
- Run a second resume and prove stable campaign/trial IDs and counts.
- Run full tests, race tests, build, vet, golangci-lint, Glazed lint, and vulnerability checks where configured.

**Exit gate:** a script prints explicit markers for manifest, diff, plan, run, resume, verify, structured output, and no-mutation dry-run.

### Phase F: Documentation and publication

- Update this guide with final APIs and exact command transcripts.
- Complete the diary with failures, commits, and review instructions.
- Relate implementation files through docmgr.
- Run `docmgr doctor`.
- Upload the guide and diary as one reMarkable bundle.
- Close the ticket only when all code and publication gates pass.

## 12. Test strategy

### 12.1 Manifest unit tests

Test:

- strict unknown-field rejection;
- schema mismatch;
- empty name and preparation;
- unsupported preparation;
- fewer than two arms;
- duplicate arm and case IDs;
- invalid retrieval limit;
- mismatched arm preparation;
- invalid positive and authorization-negative cases;
- incomplete, duplicate, cyclic, downstream, and unknown graph dependencies;
- stable identity across repeated loads;
- YAML key ordering and comments do not change identity;
- one semantic value change does change identity.

### 12.2 Command unit tests

For each command:

- decode fields through `values.Values`;
- emit expected row count and columns;
- propagate loader and processor errors;
- expose domain flags in help;
- expose `format`, `output-fields`, and `max-output-rows`;
- do not expose removed legacy output flags;
- preserve requested sparse output-field order.

### 12.3 Campaign integration tests

Use `t.TempDir()` and `SemanticFixtureExecutor`:

1. load manifest;
2. dry-run and assert no store exists;
3. run campaign;
4. assert status complete and expected episode count;
5. inspect and verify;
6. resume same campaign;
7. assert campaign/trial IDs, completed count, observations, and estimates are unchanged;
8. tamper with a copied artifact and assert verification fails.

### 12.4 Race and cancellation tests

- Run workbench and campaign package tests under `-race`.
- Cancel a context before manifest execution and assert no success row.
- Cancel during command execution and ensure the error propagates.
- Do not claim that `--max-output-rows` cancels source work.

### 12.5 CLI characterization

Representative commands:

```bash
rag-ttc experiment optkit-rag config validate \
  --manifest pkg/ttc/experimentworkbench/testdata/semantic-limit-v1.yaml \
  --format json

rag-ttc experiment optkit-rag config plan \
  --before baseline.yaml --after challenger.yaml \
  --format table

rag-ttc experiment optkit-rag campaign dry-run \
  --manifest baseline.yaml --store /tmp/does-not-exist \
  --format json

rag-ttc experiment optkit-rag campaign run \
  --manifest baseline.yaml --store /tmp/rag-workbench --reset \
  --format json

rag-ttc experiment optkit-rag campaign verify \
  --store /tmp/rag-workbench --campaign campaign:... \
  --format json
```

## 13. Failure model and troubleshooting

### Manifest parse failure

**Symptom:** unknown field, malformed YAML, multiple documents, or trailing data.

**Action:** fix the input. Do not add permissive fallback decoding.

### Graph resolution failure

**Symptom:** missing layer, duplicate identity, unknown dependency, or dependency points downstream.

**Action:** inspect layer refs and direct dependencies. Do not sort away semantic mistakes or infer missing edges.

### Unsupported preparation

**Symptom:** manifest validates structurally but cannot construct an executor.

**Action:** use the semantic fixture or implement a reviewed preparation adapter. Do not silently select a default connected runtime.

### Existing campaign with a different local manifest

**Symptom:** operator tries to pass a manifest while resuming.

**Action:** reject the combination. Resume trusts the stored campaign specification.

### Journal verification failure

**Symptom:** status or verify reports a sequence/digest mismatch.

**Action:** stop. Preserve the store for investigation. Do not append repair events to a journal whose authority is unknown.

### Lease interruption

**Symptom:** process exits after leasing or completing work.

**Action:** run `campaign resume`. Expired leases are reclaimed and completed work is reconciled idempotently.

### Missing paired estimate

**Symptom:** campaign is still running or a required observation is absent.

**Action:** show null/absent paired delta and inspect episode state. Never report zero.

## 14. Alternatives rejected

### Build the UI first

Rejected because it would force the browser to become the first consumer and interpreter of evolving scientific records. The CLI is cheaper to characterize, easier to automate, and can prove projector inputs before visual design.

### Put orchestration directly in Cobra handlers

Rejected because HTTP projectors and tests would need to duplicate it. Commands should decode settings, call services, and emit rows.

### Make Optkit understand RAG manifests

Rejected because corpus, retrieval, fusion, evidence, answer, and judge layers are product semantics. Optkit owns generic snapshots, campaigns, observations, and artifacts.

### Accept arbitrary YAML maps

Rejected because silent key errors and untyped values undermine reproducibility. The manifest must be versioned and strictly decoded.

### Store only the manifest path in the campaign

Rejected because paths are mutable and host-local. The campaign stores canonical snapshots, case artifacts, and a complete trial specification.

### Implement cancellation before backend support

Rejected for this ticket. A CLI command must not imply safe campaign cancellation until durable cancellation state, worker behavior, and lease races are designed and tested.

## 15. Risks and follow-up work

### Risk: manifest and campaign spec drift

The manifest contains a layered graph while the current `CampaignSpec` stores retrieval arms and cases but not the graph. Phase 3 should either persist the graph as a campaign-reachable artifact or clearly report that the v1 campaign executes the retrieval subset only. The preferred follow-up is to attach manifest and graph artifacts during campaign creation.

### Risk: semantic-fixture scope looks broader than it is

A manifest can name all twelve layers, but the first executor varies retrieval limits only. The CLI must not claim that context, answer, or judge layers were executed. The graph is planning metadata until corresponding producers materialize those stages.

### Risk: verification wording

Direct event payload verification is not full transitive custody. Output columns and documentation must state `journal_verified` and `direct_payloads_verified`, not a vague global `verified` boolean.

### Risk: legacy command overlap

`experiment answers`, `experiment representations`, and `tool-eval ragopt` still perform work outside the durable Optkit campaign. The new workbench should not wrap them superficially. Migrate only after durable parity exists.

### Follow-up: projector parity

After this ticket, server-side projectors should call the same application services or consume the same stable records. Characterization tests should compare CLI JSON with HTTP DTO source fields.

### Follow-up: connected preparations

Add preparations incrementally with explicit runtime identity, provider profile digest, cache policy, budget ceilings, and secret-free manifests.

### Follow-up: manifest authoring

A future `config init` or `config derive` command can construct typed layer refs from value documents. It should not be required for strict v1 validation.

## 16. Review checklist

A reviewer should verify:

- scientific intent is represented in the manifest;
- the manifest is strictly decoded and content-identified;
- graph semantics come from `optimization`, not command code;
- dry-run creates no store;
- create and resume are distinct;
- resume ignores no local scientific inputs because it accepts none;
- every mutation flows through `optkitcampaign.Run`;
- every read verifies the journal or labels integrity accurately;
- errors propagate rather than exit inside packages;
- Glazed owns serialization but not domain limits;
- command output is stable and bounded;
- the unrelated Optkit `store/sqlite/rows.go` worktree edit is never staged.

## 17. File references

### Current implementation

- `rag-ttc/cmd/rag-ttc/main.go` — CLI composition root, logging, error ownership.
- `rag-ttc/cmd/rag-ttc/cmds/experiments/root.go` — experiment command registration.
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/command.go` — current fixed durable campaign commands.
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/command_test.go` — current Glazed flag and row tests.
- `rag-ttc/pkg/ttc/optkitcampaign/campaign.go` — durable create/resume/execute/reconcile service.
- `rag-ttc/pkg/ttc/optkitcampaign/system.go` — retrieval config, case, executor, and prepared system.
- `rag-ttc/pkg/ttc/optkitcampaign/fixture.go` — deterministic fixture preparation and current hard-coded plan.
- `rag-ttc/pkg/ttc/optimization/contracts.go` — canonical layers and lineage contracts.
- `rag-ttc/pkg/ttc/optimization/graph.go` — local and resolved configuration identity.
- `rag-ttc/pkg/ttc/optimization/invalidation.go` — direct diff and transitive reuse/recompute plan.
- `rag-ttc/pkg/ttc/optimization/fixture.go` — frozen cross-layer fixture loader.
- `optkit/query/service.go` — generic bounded read-only query precedent.
- `optkit/cmd/optkit/main.go` — generic Optkit integrity CLI boundary.

### New implementation target

- `rag-ttc/pkg/ttc/experimentworkbench/` — manifest and operator application services.
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/config.go` — config commands.
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/campaign.go` — campaign commands.
- `rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/rows.go` — stable row mappings.
- `rag-ttc/pkg/ttc/experimentworkbench/testdata/` — baseline/challenger manifests.

## 18. Definition of done

This ticket is complete when a new engineer can:

1. read this guide and explain the Optkit/RAG-TTC/CLI ownership boundary;
2. validate a strict manifest without touching a store;
3. inspect all twelve resolved configuration layers;
4. compare two manifests and explain direct versus upstream invalidation;
5. dry-run exact arm, case, repeat, episode, and budget counts;
6. create a deterministic campaign from the manifest;
7. resume it without changing identity or duplicating scientific records;
8. inspect lifecycle, estimates, budgets, and integrity;
9. verify the journal and direct event payload custody;
10. consume every result through Glazed structured output;
11. run the complete validation script successfully;
12. find the exact implementation files, tests, failures, and commits in the diary.
