---
Title: Optkit Records, Artifacts, Experiments, and Campaign Control
Slug: optkit-records-artifacts-control-model
Short: Field-level reference for Optkit's public data model, ownership boundaries, writers, readers, and removal decisions.
Topics:
- optkit
- records
- artifacts
- experiments
- campaigns
- scheduler
- budgets
Commands:
- optkit demo
- optkit campaign inspect
- optkit campaign verify
- optkit artifact verify
Flags:
- store
- id
- tail
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

Optkit separates immutable scientific facts from mutable execution control. Product adapters write typed configurations, episode evidence, observations, and campaign commands; Optkit stores those values as canonical records and artifacts; projectors and query services rebuild read models without becoming a second authority. This separation makes experiment identity, restart behavior, and measurement attribution explicit.

This reference documents the exported records and interfaces that participate in that model. Each table names the writer, reader, and operational status of every field. Items marked **removed before v0 freeze** were present in the imported implementation but had no complete behavioral consumer. They remain documented so readers understand the deletion and the evidence required to reintroduce them.

## Ownership and status vocabulary

Optkit uses a small set of ownership roles. Distinguishing these roles prevents a stored field from being mistaken for enforced behavior.

| Term | Meaning |
|---|---|
| Product adapter | Domain code that prepares configurations, executes cases, and emits domain events. RAG-TTC is one adapter; Numbergame is the proof adapter. |
| Record constructor | Optkit function that validates a value and derives its semantic identity. |
| Artifact store | Immutable byte store that verifies exact content by digest. |
| Journal | Append-only campaign fact store with optimistic version checks and hash-chain verification. |
| Queue | Mutable execution-control store for ready, leased, completed, and failed work. |
| Budget ledger | Mutable accounting store that reserves finite quantities and commits actual usage. |
| Projector | Pure reader that folds authoritative facts into a rebuildable view. |
| Active | Written and consumed by a current RAG-TTC or Optkit execution path. |
| Implemented foundation | Has executable behavior and tests, but the first RAG-TTC campaign does not exercise it. |
| Declared only | Stored or exposed without a behavioral consumer. |
| Removed before v0 freeze | Former declared-only API deleted because it created a false capability signal or duplicated another authority. |

The removal rule is strict: a field remains only when a current product, executable proof system, integrity check, or read model gives it semantics. A field does not remain merely because a future distributed system might need it.

## End-to-end write and read flow

A campaign moves through distinct data and control boundaries. Product code never writes SQLite tables directly, and query code never mutates campaign state.

```text
product configuration
    -> space.SnapshotRecord + immutable config artifact
    -> experiment.TrialPlan + EpisodeSpec
    -> scheduler.WorkItem + budget.Reservation
    -> scheduler.Lease
    -> episode.Event payload artifacts
    -> episode.Trajectory + episode.Result
    -> measure.Observation
    -> campaign.ControlEvent journal
    -> projection.Overview / query views
```

The queue controls who may execute work. The budget ledger controls how much finite resource a campaign may consume. The journal records durable scientific and lifecycle facts. These responsibilities overlap in time but not in authority.

## Record identities

The `record` package supplies semantic digests and namespaced string IDs. Constructors derive content IDs from canonical semantic payloads; random operational IDs are limited to authorities such as leases and commands.

### `record.Digest`

A digest is a validated `sha256:<hex>` string. `record.SemanticDigest` writes it from a schema ID and canonical payload, while artifact stores write exact-byte digests. Snapshot, candidate, episode, observation, event, and work constructors read digests when deriving higher-level identities.

| Operation | Writer | Reader | Purpose |
|---|---|---|---|
| `SumBytes` | Artifact and fixture code | Artifact verification | Identify exact bytes. |
| `SemanticDigest` | Record constructors | Identity verification | Identify schema-qualified meaning. |
| `DigestParts` | Deterministic seed and key code | Experiment expansion | Hash ordered components without JSON. |
| `ParseDigest` / `Validate` | CLI and stores | Every custody boundary | Reject malformed identities before use. |

### Namespaced IDs

Each exported ID type prevents accidental interchange at compile time while retaining stable JSON strings.

| Type | Written by | Read by | Purpose and status |
|---|---|---|---|
| `SchemaID` | Schema authors and codecs | Stores, trajectories, journal, query API | Names the decoding contract. Active. |
| `SystemID` | Product system registration | Snapshots and `system.Registry` | Selects the executable factory. Active. |
| `SnapshotID` | `space.MaterializeSnapshot` | Arms, candidates, registry preparation | Identifies complete executable configuration. Active. |
| `PatchID` | `PatchBuilder.Build` | `space.Candidate` | Identifies an ordered, canonical intervention. Implemented foundation. |
| `CandidateID` | `space.NewCandidate` | Campaign facts and future optimization projectors | Identifies a proposed child snapshot plus hypothesis. Used by Numbergame. |
| `CampaignID` | Product campaign creation | Journal, queue, budget, query API | Root authority for one campaign. Active. |
| `EventID` | Episode writer or campaign materializer | Trajectories, journal chain, query API | Identifies immutable events. Active. |
| `CommandID` | Command caller | Controller and command lookup | Makes lifecycle commands idempotent. Active. |
| `CorrelationID` **(removed)** | No meaningful writer | Former campaign event storage and query API | Removed before v0 freeze because no command, worker, projector, or UI assigned or queried it. |
| `EpisodeID` | Experiment expansion | Trajectory writer, campaign reducer | Identifies one arm/case/repeat execution. Active. |
| `TrialID` | `NewCompleteBlockTrial` | Episode specs and campaign summary | Identifies an experimental design. Active. |
| `ObservationID` | `measure.NewObservation` | Journal and estimators | Identifies one attributed measured fact. Active. |
| `EpochID` | `measure.NewEpoch` | Observations and measurement queries | Separates incompatible measurement populations. Active. |
| `WorkID` | `scheduler.NewWorkItem` | Queue and budget ledger | Identifies schedulable work deterministically. Active. |
| `LeaseID` | SQLite queue `Lease` | `Complete`, `Fail`, and lease fencing | Identifies one temporary execution authority. Active. |
| `ReservationID` | `budget.ReservationID` | Budget commit and release | Identifies finite resource custody for one work item. Active. |
| `SpanID` | Product executable | Episode writer | Connects parent and child episode events. Active. |
| `ActorRef` | Command caller and fact appender | Journal, audit views, candidate records | Attributes a mutation or proposal. Active. |

### `record.Schema`

A schema record binds a stable ID to a content digest and optional description. Schema registration rejects a second definition that uses the same ID with different content.

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `ID` | Schema author | `SchemaRegistry` | Stable schema name. |
| `Digest` | Schema author or schema loader | `SchemaRegistry` | Detects conflicting definitions. |
| `Description` | Schema author | Documentation and diagnostics | Human explanation; it does not affect runtime decoding by itself. |

`SchemaRegistry` is an in-memory concurrency-safe registry. `Register` writes definitions; `Get` and `List` read them. It is not a durable schema migration service.

## Artifact custody

The `artifact` package stores immutable bytes separately from control metadata. A reference states what bytes are expected; `Store.Verify` proves that the current object still matches that expectation.

### `artifact.Ref`

| Field | Writer | Reader | Purpose and failure mode |
|---|---|---|---|
| `Digest` | Artifact store after streaming bytes | `Open`, `Stat`, `Verify`, record constructors | Exact-byte identity. A mismatch is corruption, not a cache miss. |
| `MediaType` | `PutRequest` caller | Renderers and decoders | Declares representation format. Empty values fail validation. |
| `Schema` | Typed writer | Codecs, projectors, query clients | Optional semantic decoding contract. Unknown schemas must fail closed at strict boundaries. |
| `Size` | Artifact store | Verification and bounded query previews | Exact byte length. Negative values fail validation. |
| `Sensitivity` | Product or framework writer | Query service and future authorization wrappers | Classification metadata. It is not authorization by itself. |

Sensitivity values are `public`, `internal`, `confidential`, and `restricted`. Query code may omit or bound content based on sensitivity, but callers must not infer that a label alone enforces access control.

### `artifact.PutRequest`

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `MediaType` | Product/framework artifact writer | Store implementation | Populates the resulting reference. |
| `Schema` | Typed writer | Store and later decoder | Associates bytes with a schema. |
| `Sensitivity` | Product policy | Store and query boundary | Carries classification into the reference. |
| `ExpectedDigest` | Importer or known-content writer | Store `Put` | Requires input bytes to match a prior identity. |

### `artifact.Info` and `artifact.Store`

`Info.Ref` is the authoritative metadata returned by `Stat`. Store methods divide custody as follows:

| Method | Writer or caller | What it guarantees |
|---|---|---|
| `Put` | Record, episode, and product writers | Streams bytes, derives identity, and publishes immutable content. |
| `Open` | Decoders and inspectors | Opens the exact referenced object or reports absence/corruption. |
| `Stat` | Query and verification code | Returns metadata without exposing content. |
| `Verify` | Trajectory loader, campaign verification, CLI | Recomputes size and digest custody. |

The memory store supports tests. The filesystem store is the local immutable data plane and publishes content without replacing an existing digest path.

## Configuration space

The `space` package keeps product configuration typed while persisting type-erased records and artifact references. Product code owns configuration structs; Optkit owns lawful intervention and identity rules.

### Domains and codecs

| Type | Writer | Reader | Purpose |
|---|---|---|---|
| `DomainDescriptor` | Variable author | Documentation and candidate tooling | Persistable domain kind and parameters. |
| `Domain[V]` | Variable author | Patch application | Validates candidate values and returns a descriptor. |
| `IntRangeDomain` | Product adapter | Variable validation | Accepts integers between `Min` and `Max`. |
| `ChoiceDomain[V]` | Product adapter | Variable validation | Accepts one of a declared finite set. |
| `Codec[V]` | Product adapter | Snapshot and patch code | Canonically encodes and strictly decodes typed values. |
| `JSONCodec[V]` | `space.NewJSONCodec` | Snapshot loader and patch builder | Uses schema-qualified canonical JSON; strict decoding is enabled by default. |

`JSONCodec.SchemaID` names the value contract. `JSONCodec.Strict` controls unknown-field rejection. Product-facing schemas should retain `Strict=true`.

### `Lens[C,V]`

A lens identifies one typed value inside a complete configuration.

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `Get` | Variable author | Patch builder and law tests | Reads the current value. |
| `Put` | Variable author | Patch builder and law tests | Returns a modified complete configuration. |

A valid lens must satisfy get-put, put-get, and put-put laws. Without those laws, patch order and expected-old-value checks become unreliable.

### `VariableDescriptor`

| Field | Writer | Reader | Purpose and status |
|---|---|---|---|
| `ID` | Product variable author | Patch canonicalization and UI | Stable machine-readable variable name. Active. |
| `Name` | Product variable author | Human-facing tools | Display name. Active. |
| `Description` | Product variable author | Documentation and UI | Explains intervention meaning. Active. |
| `ValueSchema` | Product variable author | `Variable.Validate` | Must match the codec schema. Active. |
| `Domain` | Product variable author | Candidate tooling | Persistable legal-value description. Active. |
| `Sensitive` | Product policy | Patch artifact writer | Stores assigned values as restricted instead of internal. Active. |
| `Probes` | Product variable author | Future intervention diagnostics | Names checks that prove the intervention took effect. Declared and product-defined; RAG-TTC will consume it in later campaign work. |

`Variable[C,V]` combines the descriptor, lens, domain, and codec. The descriptor is persistable; the function-bearing lens and interfaces remain runtime values.

### `SnapshotRecord` and `Snapshot[C]`

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `ID` | `MaterializeSnapshot` | Trial arms, candidates, system registry | Content-derived complete configuration identity. |
| `System` | Product adapter | `system.Registry.Prepare` | Selects the factory allowed to decode and execute the snapshot. |
| `Schema` | Config codec | Snapshot loader and factory | Names the complete configuration schema. |
| `Config` | Artifact store | Snapshot loader and factory | Immutable canonical configuration artifact. |
| `Value` | `MaterializeSnapshot` or `LoadSnapshot` | Typed authoring code | In-memory typed value; excluded from JSON. |

`LoadSnapshot` decodes the config, rematerializes its identity, and rejects an ID or digest mismatch.

### Patch records

| Type and field | Writer | Reader | Purpose |
|---|---|---|---|
| `AssignmentRecord.Variable` | `PatchBuilder` | Candidate/projector code | Variable being changed. |
| `ExpectedOldDigest` | `PatchBuilder` | Verification and review | Digest of the base value seen when constructing the patch. |
| `Value` | `PatchBuilder` | Patch review and materialization | Immutable new-value artifact. |
| `PatchRecord.ID` | `PatchBuilder.Build` | Candidate and campaign facts | Semantic identity of the intervention. |
| `Base` | `PatchBuilder.Build` | `VerifyBase` and review | Snapshot to which the patch applies. |
| `Assignments` | `PatchBuilder.Build` | Review and future diff projectors | Canonically sorted variable assignments. |
| `Result` | `PatchBuilder.Build` | Candidate and trial planning | Child snapshot produced by the patch. |
| `Patch.ResultConfig` | `PatchBuilder.Build` | Typed authoring code | In-memory child config; excluded from JSON. |

`PatchBuilder` holds a base snapshot, artifact store, complete-config codec, and a private map of typed assignments. It rejects duplicate variable assignments and empty patches.

### `Candidate`

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `ID` | `NewCandidate` | Campaign facts and future UI | Content-derived proposal identity. |
| `Parent` | Candidate proposer | Review and lineage projectors | Incumbent snapshot. |
| `Patch` | Candidate proposer | Review and lineage projectors | Proposed intervention. |
| `Child` | Candidate proposer | Trial planner | Resulting executable snapshot. |
| `Proposer` | Agent, user, or strategy | Audit views | Attribution. |
| `Strategy` | Candidate proposer | Analysis and UI | Proposal mechanism. |
| `Hypothesis` | Candidate proposer | Reviewer and optimizer | Required expected causal explanation. |
| `Targets` | Candidate proposer | Reviewer | Intended constructs or metrics. |
| `Risks` | Candidate proposer | Reviewer | Known regressions or policy risks. |
| `CreatedAt` | Candidate proposer | UI ordering | Display time; excluded from candidate semantic identity. |

Numbergame writes candidates today. Product optimization campaigns will use the same record after layered configuration identities are available.

## System registration and execution

The `system` package is the domain-neutral composition seam. Product code registers factories; workers select only registered systems and prepared snapshots.

### `system.Factory`

| Method | Implemented by | Called by | Purpose |
|---|---|---|---|
| `SystemID` | Product adapter | Registry registration | Stable system authority. |
| `ConfigSchema` | Product adapter | Registry preparation | Declares the accepted snapshot schema. |
| `CaseSchema` | Product adapter | Prepared execution and validation | Declares accepted case artifacts. |
| `Prepare` | Product adapter | Worker through registry | Verifies and compiles one immutable snapshot into runtime behavior. |

### `system.Prepared`

| Method | Implemented by | Called by | Purpose |
|---|---|---|---|
| `SystemID` | Prepared product system | Worker | Confirms selected system identity. |
| `SnapshotID` | Prepared product system | Worker and diagnostics | Attributes execution to exact configuration. |
| `CaseSchema` | Prepared product system | Worker | Rejects incompatible case artifacts. |
| `Run` | Prepared product system | Campaign worker | Executes one case and emits episode evidence. |

`Registry` uses a concurrency-safe map of factories. `Register` rejects duplicate system IDs. `Prepare` rejects unknown systems and schema mismatches before product behavior runs.

## Episode evidence

The `episode` package records one execution as ordered immutable events followed by a sealed trajectory and terminal result. Episode facts describe product behavior; they do not mutate campaign lifecycle directly.

### `episode.Event`

| Field | Writer | Reader | Purpose and status |
|---|---|---|---|
| `ID` | `episode.Writer.Emit` | Trajectory verification and projectors | Content-derived event identity. Active. |
| `Episode` | Writer constructor | Trajectory loader | Prevents cross-episode event injection. Active. |
| `Seq` | Writer | Loader and UI | Strict one-based order. Active. |
| `Time` | Writer clock | UI and diagnostics | Observed event time; participates in identity. Active. |
| `Kind` | Product executable | Domain projector | Stable event vocabulary. Active. |
| `Schema` | Product executable | Strict decoder | Payload contract. Active. |
| `Span` | Product executable | Pipeline projector | Operation grouping. Active. |
| `Parent` | Product executable | Writer validation and pipeline projector | Causal span parent. Active. |
| `Payload` | Writer after storing payload | Loader and projector | Immutable event data. Active. |
| `Tags` | Product executable | Bounded filtering and diagnostics | Small non-authoritative labels. Active. |
| `Sensitive` | Writer from sensitivity | Query policy | Signals content-bearing event payload. Active. |

### `episode.Emission`

`Emission` is the pre-persistence form of an event. Product code writes `Kind`, `Schema`, `Span`, optional `Parent`, `Payload`, `Tags`, and `Sensitivity`. The writer validates spans, stores the payload, derives sequence and identity, and returns the immutable `Event`.

### `episode.ResourceUsage`

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `Resource` | Prepared system | Budget reconciliation | Names a consumed quantity such as queries or results. |
| `Units` | Prepared system | Budget reconciliation | Actual integer consumption. |

Episode usage is the source of actual budget commitment. It is different from the removed scheduler resource claims because the budget ledger consumes and enforces it.

### `episode.Failure`

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `Class` | Product executable | Failure projectors | Broad failure category. |
| `Scope` | Product executable | Failure projectors | Stage or subsystem affected. |
| `Retryable` | Product executable | Campaign worker policy | Whether a new attempt is meaningful. |
| `Code` | Product executable | Projectors and gates | Stable machine-readable reason. |
| `Message` | Product executable | Operator diagnostics | Human detail; must not carry secrets into public views. |
| `Evidence` | Product executable | Review and diagnostics | Artifact references supporting classification. Active at episode level. |
| `Diagnostics` | Product executable | Domain-specific projector | Schema-owned structured details. |

Episode failure evidence remains supported. The former `scheduler.WorkFailure.Evidence` field was removed because it had no writer or behavioral reader.

### Results and trajectories

| Type and field | Writer | Reader | Purpose |
|---|---|---|---|
| `RunResult.Status` | Prepared system | Worker and reconciliation | Completed, failed, or cancelled execution result. |
| `Output` | Prepared system | Instruments and projectors | Optional primary output artifact. |
| `Usage` | Prepared system | Budget reconciliation | Actual resource usage. |
| `Failure` | Prepared system | Failure policy and projectors | Typed product failure. |
| `StartedAt` / `FinishedAt` | Prepared system or worker | Latency instruments | Execution interval. |
| `Result.Trajectory` | Worker after sealing writer | Instruments, reconciliation, query views | Immutable trajectory artifact. |
| `Trajectory.Episode` | Writer | Loader | Episode ownership. |
| `Trajectory.Events` | Writer | Loader and projectors | Ordered immutable evidence. |

`Sink` exposes `Emit` and `AttachJSON`; `Executable[I]` runs typed input against that sink. `Writer.Seal` refuses empty or already-sealed trajectories.

## Measurement records

The `measure` package stores observations as attributed facts. A changed protocol, implementation, calibration, or redaction policy creates a new epoch rather than rewriting historical values.

### `EpochDefinition` and `Epoch`

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `Construct` | Instrument adapter | Observation constructor and estimators | Quantity or property being measured. |
| `Instrument` | Instrument adapter | Observation constructor | Measurement implementation identity. |
| `Protocol` | Instrument adapter | Observation constructor | Procedure or prompt protocol. |
| `Implementation` | Instrument adapter | Audit and comparison | Code/version identity. |
| `Calibration` | Instrument adapter | Compatibility filtering | Optional calibration population or contract identity. |
| `Redaction` | Instrument adapter | Compatibility and privacy review | Optional input-view policy identity. |
| `Epoch.ID` | `NewEpoch` | Observation | Content-derived compatibility key. |
| `Epoch.Definition` | `NewEpoch` | Audit and UI | Human- and machine-readable epoch meaning. |

### `SubjectRef` and `Value`

`SubjectRef.Kind` names the entity class and `SubjectRef.ID` names the concrete episode, candidate, or artifact. `Value.Kind` selects exactly one representation: decimal string in `Number`, boolean in `Bool`, or text in `Text`. Decimal strings avoid accidental JSON integer/float coercion at custody boundaries.

### Observation draft and record

| Field | Draft writer | Record reader | Purpose |
|---|---|---|---|
| `Role` | Instrument | Gates and projectors | Fact, intervention, constraint, or measurement. |
| `Subject` | Instrument | Grouping and provenance | Entity measured. |
| `Construct` | Instrument | Estimator | Measurement name; must match epoch. |
| `Instrument` | Instrument | Compatibility check | Must match epoch. |
| `Protocol` | Instrument | Compatibility check | Must match epoch. |
| `Epoch` | Instrument | Observation constructor | Full epoch used to validate matching fields. |
| `Status` | Instrument | Missingness and gate policy | Measured, satisfied, violated, inapplicable, unknown, or failed. |
| `Value` | Instrument | Estimator or projector | Typed result; status determines whether it is usable. |
| `Evidence` | Instrument | Verification and review | Artifact references supporting the observation. |
| `Diagnostics` | Instrument | Domain projector | Canonical JSON diagnostics. |
| `Repeat` | Instrument | Paired analysis | Repeated measurement index. |
| `Observation.ID` | Constructor | Journal and estimators | Derived from semantic fields and evidence digests. |
| `CreatedAt` | Constructor caller | UI ordering | Recording time; excluded from semantic identity. |

## Experimental design

The `experiment` package currently implements deterministic complete-block trials and paired means. It does not yet provide adaptive search, confidence intervals, Pareto dominance, or promotion policy.

### Dataset records

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `Case.ID` | Dataset author | Trial expansion and estimators | Stable case identity. |
| `Case.Input` | Dataset author | Prepared system | Immutable input artifact. |
| `Case.Groups` | Dataset author | Group-sliced analysis | Optional declared strata. |
| `Case.Metadata` | Dataset author | Product projector | Small descriptive metadata. |
| `DatasetManifest.ID` | `NewDatasetManifest` | Trial and campaign identity | Content-derived dataset ID. |
| `Role` | Dataset author | Exposure policy and UI | Development, validation, holdout, or other product role. |
| `Cases` | Dataset author | Trial expansion | Frozen complete block. |
| `Digest` | Constructor | Trial identity | Semantic dataset digest. |

### Trial and episode specifications

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `Arm.ID` | Trial author | Scheduler and estimator | Human-stable treatment label. |
| `Arm.Snapshot` | Trial author | System preparation | Exact treatment configuration. |
| `TrialPlan.ID` | Constructor | Campaign and episode identity | Semantic trial identity. |
| `Arms` | Trial author | `Expand` | Treatments applied to every case. |
| `Dataset` | Trial author | `Expand` | Frozen case block. |
| `Repeats` | Trial author | `Expand` | Number of repeated blocks. |
| `Protocol` | Trial author | Episode identity and worker | Execution procedure identity. |
| `EpisodeSpec.ID` | `Expand` | Queue, trajectory, campaign reducer | Deterministic episode ID. |
| `SemanticKey` | `Expand` | Work identity and conflict detection | Stable equivalence key. |
| `Trial` | `Expand` | Provenance | Parent trial. |
| `Arm` | `Expand` | Worker and estimator | Treatment. |
| `Case` | `Expand` | Worker and estimator | Input. |
| `Repeat` | `Expand` | Pairing | Repeat index. |
| `Seed` | `Expand` | Prepared system | Deterministic randomness seed. |
| `Protocol` | `Expand` | Worker and provenance | Copied execution protocol. |

### Paired estimates

`NumericObservation` is an in-memory estimator input rather than a durable record. `CaseID`, `ArmID`, and `Repeat` form a unique pair key; `Value` carries the number; `Valid` controls missingness.

| `Estimate` field | Written by | Read by | Purpose |
|---|---|---|---|
| `Baseline` | `PairedMean` | Campaign summary | Baseline arm ID. |
| `Treatment` | `PairedMean` | Campaign summary | Challenger arm ID. |
| `Value` | `PairedMean` | Decision code and UI | Mean treatment-minus-baseline delta. |
| `SampleSize` | `PairedMean` | Review and UI | Number of complete pairs. |
| `Missing` | `PairedMean` | Failure handling | Count of absent or invalid pairs. Any positive count currently invalidates the estimate. |
| `Deltas` | `PairedMean` | Review and future intervals | Individual paired differences in deterministic order. |

## Campaign commands, events, and state

The `campaign` package is the control-plane authority. Commands request validated lifecycle transitions; control events record accepted facts; reducers rebuild state solely from those events.

### `campaign.Command`

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `ID` | CLI/product controller | Command lookup | Idempotency key. |
| `Campaign` | CLI/product controller | Journal and reducer | Target campaign. |
| `ExpectedVersion` | Caller after reading journal head | Journal append | Optimistic concurrency fence. |
| `Kind` | Caller | Controller | Requested lifecycle transition. |
| `Actor` | Caller | Event attribution | Mutation authority. |
| `Schema` | Caller | Artifact writer and event | Payload contract. |
| `Payload` | Caller | Controller artifact writer | Command-specific immutable data. |
| `Tags` | Caller | Transition validation and resulting event | Small transition metadata. |

`Controller.Journal` reads state and appends accepted events. `Controller.Artifacts` stores command payloads. The controller looks up prior command IDs before writing, making retries idempotent where the journal implements `CommandLookup`.

### Event kinds

| Group | Writers | Readers | Status |
|---|---|---|---|
| Campaign lifecycle | Controller | Reducer, overview, query API | Active; pause/stop paths are implemented but not exposed by the current RAG-TTC CLI. |
| Candidate and snapshot facts | Numbergame today; future product optimizer | Reducer and future lineage projector | Implemented foundation. |
| Trial and episode facts | RAG-TTC and Numbergame | Reducer, reconciliation, query API | Active. |
| Observation and estimate facts | RAG-TTC and Numbergame | Analysis and query API | Active. |
| Decision facts | Numbergame | Future promotion projector | Implemented foundation. |
| Budget facts | RAG-TTC and Numbergame | Reducer and query API | Active. |

### `NewEvent` and `ControlEvent`

`NewEvent` is the append request. `ControlEvent` is the immutable, sequenced, hash-chained result.

| Field | Writer | Reader | Purpose and status |
|---|---|---|---|
| `Kind` | Controller or fact appender | Reducer/projector | Event semantics. Active. |
| `Schema` | Writer | Payload decoder | Event payload contract. Active. |
| `Subject` | Writer | Reducer/projector | Entity affected. Active. |
| `OccurredAt` | Writer or journal default | UI and event identity | Domain occurrence time. Active. |
| `Actor` | Writer | Audit/query | Attribution. Active. |
| `Command` | Controller | Command lookup | Idempotent command provenance. Active. |
| `Causation` **(removed)** | No meaningful product writer | Former SQLite and query serialization only | Removed because no workflow assigned or queried immediate-cause relationships. |
| `Correlation` **(removed)** | No meaningful product writer | Former SQLite and query serialization only | Removed because no operation established or queried correlation groups. |
| `Payload` | Writer | Reducer/projector/verifier | Immutable event data. Active. |
| `Tags` | Writer | Query API and diagnostics | Small labels. Active. |
| `Campaign` | Journal materializer | Chain verifier/projector | Owning campaign. Active. |
| `Seq` | Journal materializer | Cursor, replay, SSE, chain verifier | Monotonic one-based sequence. Active. |
| `RecordedAt` | Journal | UI and identity | Durable append time. Active. |
| `PreviousDigest` | Journal | `VerifyChain` | Hash-chain link. Active. |
| `Digest` | Materializer | `VerifyChain` and query API | Event semantic digest. Active. |
| `ID` | Materializer | Query API and references | Content-derived event ID. Active. |

Removing causation and correlation requires an event-schema decision because omitted fields currently participate in the materializer's identity shape. Existing journals must either remain on their historical event version or be treated as pre-freeze local fixtures.

### Journal records and interfaces

| Type or method | Writer | Reader | Purpose |
|---|---|---|---|
| `Head.Version` | Journal | Controller and appenders | Current optimistic version. |
| `Head.LastDigest` | Journal | Appender and integrity views | Current chain tip. |
| `AppendResult.Version` | Journal append | Caller | New head version. |
| `AppendResult.Events` | Journal append | Caller | Materialized accepted facts. |
| `AppendResult.Duplicate` | Controller on command retry | Caller | Signals idempotent replay. |
| `CommandLookup.LookupCommand` | SQLite journal implementation | Controller | Finds prior result for a command ID. |
| `Journal.Append` | Controller/fact appender | SQLite implementation | Atomically append at expected version. |
| `Journal.Read` | Reducer/projectors/query | SQLite implementation | Read events after a sequence cursor. |
| `Journal.Head` | Controller/query | SQLite implementation | Read current chain head. |
| `Journal.Verify` | CLI/query | SQLite implementation | Verify complete chain integrity. |

### `campaign.State`

| Field | Written by | Read by | Purpose |
|---|---|---|---|
| `Campaign` | `Initial` | Transition validation | Owning campaign. |
| `Version` | Reducer | Controller | Folded journal version. |
| `LastDigest` | Reducer | Integrity and append code | Folded chain tip. |
| `Status` | Reducer | Transition validation and projections | Draft, ready, running, paused, stopping, stopped, failed, or completed. |
| `Episodes` | Reducer | Reconciliation and overview | Per-subject queued, leased, running, completed, or terminal-failed state. |
| `Events` | Reducer | Overview | Number of folded events. |

The state is derived and disposable. The journal remains authoritative.

## Durable work scheduling

The `scheduler` package defines mutable execution authority. Work identity is deterministic; lease identity is intentionally unique per attempt.

### `WorkItem`

| Field | Writer | Reader | Purpose and status |
|---|---|---|---|
| `ID` | `NewWorkItem` | Queue, budget ledger, reconciliation | Deterministic work identity. Active. |
| `Campaign` | Product campaign | Queue and budget ledger | Owning campaign. Active. |
| `Kind` | Product campaign | Kind-filtered worker lease | Worker routing. Active. |
| `SemanticKey` | Experiment expansion | `NewWorkItem` and enqueue conflict check | Equivalent-work identity. Active. |
| `Payload` | Product campaign | Worker | Immutable execution request. Active. |
| `Priority` | Product campaign | SQLite lease ordering | Higher values lease first. Active mechanism; current products use zero. |
| `EarliestStart` | Product campaign or retry path | SQLite lease filter | Initial delay and retry backoff. Active. |
| `LeaseDuration` | Product campaign | SQLite lease | Expiry window. Active. |
| `ResourceClaims` **(removed)** | Formerly RAG-TTC and Numbergame | Former SQLite serialization only | Removed because the queue never compared claims with worker capacity and enforced budget claims duplicated the values. |

The removed `ResourceClaim.Resource` named a capacity and `Units` declared an amount, but neither field had a behavioral reader. Commit `fda4ad62c97de2db4d85eab78aea926a7d474f5f` removed the type, the `WorkItem` field, the `NewWorkItem` argument, and active SQL serialization. Pre-freeze local stores created with the old `resource_claims_json NOT NULL` schema must be reset before enqueueing new work.

### Lease request and lease

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `LeaseRequest.Kinds` | Worker | SQLite queue | Select supported work kinds. |
| `Limit` | Worker | SQLite queue | Bound transaction batch size. |
| `Now` | Worker/test or store clock | SQLite queue | Deterministic eligibility and expiry time. |
| `Lease.ID` | SQLite queue | `Complete` and `Fail` | Fences one execution attempt. |
| `Item` | SQLite queue | Worker | Work to execute. |
| `Worker` | Lease caller | Queue and diagnostics | Worker authority name. |
| `Attempt` | SQLite queue | Worker and attempt facts | One-based lease attempt. |
| `ExpiresAt` | SQLite queue | Worker diagnostics and reclamation | Current authority deadline. |

A stale worker cannot commit after reassignment because terminal updates require the exact active lease ID and leased status.

### Work result and failure

| Field | Writer | Reader | Purpose and status |
|---|---|---|---|
| `WorkResult.Artifact` | Worker | Queue and reconciliation | Canonical terminal result artifact. Active. |
| `WorkFailure.Code` | Worker | Queue and diagnostics | Stable failure reason. Implemented in Numbergame. |
| `Message` | Worker | Operator diagnostics | Human detail. Implemented foundation. |
| `Retryable` | Worker | Queue `Fail` | Return to ready state or terminate. Implemented foundation. |
| `Backoff` | Worker | Queue `Fail` | New earliest-start delay. Implemented foundation. |
| `Evidence` **(removed)** | No meaningful writer | Former SQLite serialization only | Removed because episode failures already carry evidence and no queue policy or projector consumed the duplicate field. |

### `WorkRecord`

`WorkRecord` is the queue's current mutable view.

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `Item` | Enqueue | Worker and reconciliation | Immutable work definition. |
| `Status` | Queue transitions | Reconciliation | Ready, leased, completed, or failed. |
| `Attempt` | Lease | Worker and diagnostics | Current attempt count. |
| `LeaseID` | Lease/reclaim | Completion fencing and diagnostics | Active authority, empty when unleased. |
| `LeasedBy` | Lease/reclaim | Diagnostics and heartbeat implementation | Active worker name. |
| `LeaseExpiry` | Lease/reclaim | Reclamation | Active authority deadline. |
| `Result` | Complete | Reconciliation | Terminal output. |
| `Failure` | Fail | Diagnostics and retry state | Latest failure. |

### Queue methods

| Method | Current callers | Behavioral consumer | Status |
|---|---|---|---|
| `Enqueue` | RAG-TTC and Numbergame | SQLite queue | Active. |
| `Lease` | RAG-TTC and Numbergame | SQLite queue | Active. |
| `Heartbeat` **(removed)** | None | Former SQLite expiry update | Removed because it advertised multi-worker liveness without a heartbeat loop, cancellation policy, or reclamation race test. Reintroduce with the first real long-running worker. |
| `Complete` | RAG-TTC and Numbergame | SQLite fenced update | Active. |
| `Fail` | Numbergame | SQLite retry/terminal transition | Implemented foundation; RAG-TTC still needs failure integration. |
| `ReclaimExpired` | RAG-TTC and tests | SQLite queue | Active crash recovery. |
| `Get` | RAG-TTC reconciliation and tests | SQLite queue | Active. |

Lease IDs, lease expiry, and reclamation remain after `Heartbeat` left the v0 interface. The local worker model tolerates expiry and replay; Phase 3 bundle-build work must introduce heartbeat as a complete worker protocol rather than a lone method.

## Budget accounting

The `budget` package enforces finite integer quantities. Unlike scheduler resource claims, budget claims have writers, readers, atomic availability checks, terminal states, and campaign facts.

### Quantity and limit records

| Type and field | Writer | Reader | Purpose |
|---|---|---|---|
| `Resource` | Product campaign | Normalizers and ledger | Validated resource name such as `retrieval.queries`. |
| `Quantity.Resource` | Product campaign or episode reconciliation | Reserve/commit | Resource being consumed. |
| `Quantity.Units` | Product campaign or episode reconciliation | Reserve/commit | Nonnegative amount. |
| `Limit.Resource` | Campaign initialization | Ledger | Finite resource name. |
| `Limit.Units` | Campaign initialization | Ledger | Maximum campaign amount. |

### Reservation records

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `ID` | `ReservationID`/ledger | Commit, release, journal payload | Deterministic reservation identity. |
| `Campaign` | Reserve request | Ledger and projectors | Owning campaign. |
| `Work` | Reserve request | Reconciliation lookup | Work item holding custody. |
| `Requested` | Reserve request | Availability accounting | Pre-execution reservation. |
| `Actual` | Commit | Usage accounting | Observed terminal usage. |
| `Status` | Ledger | Reconciliation | Reserved, committed, or released. |
| `Overage` | Commit | Gates and UI | Actual usage exceeded reservation or limit. |
| `CreatedAt` / `UpdatedAt` | Ledger clock | Audit/UI | Lifecycle timestamps. |

`ReserveRequest` contains `Campaign`, `Work`, and enforced `Claims`. The ledger normalizes quantities, rejects duplicates, and atomically refuses insufficient availability.

### Budget snapshots

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `ResourceSnapshot.Resource` | Ledger | Query/UI | Resource name. |
| `Limit` | Ledger | Query/UI | Configured maximum. |
| `Reserved` | Ledger | Query/UI | Outstanding custody. |
| `Committed` | Ledger | Query/UI | Terminal actual usage. |
| `Available` | Ledger | Scheduler/operator | Limit minus reserved and committed. |
| `Snapshot.Campaign` | Ledger | Query/UI | Owning campaign. |
| `Snapshot.Resources` | Ledger | Query/UI | Per-resource conservation state. |
| `Snapshot.Violated` | Ledger | Gates/UI | At least one committed amount exceeds policy. |

`Ledger.Define`, `Reserve`, `Commit`, `Release`, reservation lookup, and `Snapshot` all have concrete SQLite behavior and tests. These APIs remain.

## Projections and query records

Projection and query types are read models. They may be rebuilt or versioned without changing journal authority.

### `projection.Overview`

| Field | Source | Consumer | Purpose |
|---|---|---|---|
| `Campaign` | Requested campaign ID | CLI/UI | View identity. |
| `Version` | Folded state | Cursor and UI | Journal version represented. |
| `Status` | Folded state | UI | Campaign lifecycle. |
| `Events` | Folded state | UI | Event count. |
| `Episodes` | Folded state | UI | Total episode subjects. |
| `Queued` | Episode states | UI | Awaiting lease. |
| `Active` | Leased/running states | UI | In progress. |
| `Completed` | Episode states | UI | Successful terminal count. |
| `FailedTerminal` | Episode states | UI | Failed terminal count. |
| `LastDigest` | Folded state | Integrity display | Journal chain tip. |

### Query interfaces and service

`CampaignCatalog` lists campaign IDs. `CampaignReader` reads, heads, and verifies journals. `BudgetReader` returns budget snapshots. `MetadataStore` composes those read-only contracts. `query.Service` receives `Metadata`, `Artifacts`, and an optional `Now` clock; it exposes no mutation method.

### Query response records

| Type and important fields | Source | Consumer | Purpose |
|---|---|---|---|
| `Integrity.Verified` / `Error` | Journal verification | UI | Visible integrity state without hiding failures. |
| `CampaignSummary` | Journal fold, creation payload, budget ledger | Campaign explorer | One bounded campaign overview with API version and generation time. |
| `CampaignPage.Campaigns` / `Count` | Catalog and summaries | Campaign list | Read-only campaign collection. |
| `EventView` | `ControlEvent` plus optional payload preview | Event timeline and SSE | Public query representation of one event. |
| `EventPage.After` / `Through` / `Head` | Sequence cursor and journal | Pagination client | Stable cursor bounds. |
| `EventPage.Events` / `HasMore` | Bounded journal read | Pagination client | Event page and continuation signal. |

The former `EventView.Causation` and `EventView.Correlation` fields were removed with their control-event fields. `PayloadJSON`, `PayloadAvailable`, and `PayloadReason` remain because the query service must explain why a preview is present, absent, sensitive, too large, or undecodable.

The query API defaults to 200 events, permits at most 500, and refuses payload previews above 256 KiB.

## Local profile and SQLite composition

`local.Profile` is the local composition root.

| Field | Writer | Reader | Purpose |
|---|---|---|---|
| `Root` | `local.Open` | CLI and diagnostics | Selected local profile directory. |
| `Metadata` | `local.Open` | Campaign, queue, budget, query code | SQLite implementation of control interfaces. |
| `Artifacts` | `local.Open` | Product/framework writers and readers | Filesystem content-addressed store. |

`store/sqlite.Store` combines journal, command lookup, campaign catalog, queue, and budget ledger interfaces over one SQLite database. This is a local deployment choice, not permission for packages to bypass their interfaces or join artifact bytes into control tables.

## Completed removal register

The removal register preserves why pre-release fields were deleted in commit `fda4ad62c97de2db4d85eab78aea926a7d474f5f`. These were not deprecated compatibility APIs; they failed the writer-reader-enforcement test before the v0 contract freeze.

| Removed item | Former writer | Former reader | Why removal was correct | Reintroduction gate |
|---|---|---|---|---|
| `scheduler.ResourceClaim` | RAG-TTC and Numbergame constructors | SQLite serializer/decoder only | Duplicates enforced budget quantities and does not affect lease admission. | A worker-capacity model compares claims with declared worker capacities during an atomic lease. |
| `WorkItem.ResourceClaims` | `NewWorkItem` | SQLite round trip only | Makes work appear capacity-aware when it is not. | Same as above. |
| `Queue.Heartbeat` | None | SQLite update implementation | No heartbeat loop, cancellation policy, or concurrent reclamation protocol calls it. | A long-running worker extends leases and has race tests against reclaim/complete. |
| `scheduler.WorkFailure.Evidence` | No meaningful writer | SQLite round trip only | Duplicates `episode.Failure.Evidence`; scheduling policy never reads it. | A queue-level failure policy consumes evidence independently of episode records. |
| `NewEvent.Causation` / `ControlEvent.Causation` | No meaningful writer | SQLite and query serialization only | No operation establishes or queries immediate-cause relationships. | A projector or operator workflow requires and validates causation edges. |
| `NewEvent.Correlation` / `ControlEvent.Correlation` | No meaningful writer | SQLite and query serialization only | No operation establishes or queries correlation groups. | A command or distributed operation assigns one stable correlation identity and a reader uses it. |
| `EventView.Causation` / `Correlation` | Query projection of unused fields | Browser serialization | They expose empty concepts and freeze unnecessary API fields. | Control-event fields return with real semantics. |
| `record.CorrelationID` | No meaningful writer | Only unused event fields | Type has no independent behavior after correlation fields leave. | A real correlation contract is introduced. |

Removal does not include lease IDs, expiry, attempts, reclaim, queue failure, budget release, candidate records, decisions, or lifecycle events. Those all have current execution behavior, proof-system behavior, or near-term Phase 2 semantics.

## Failure and integrity rules

Optkit fails closed at boundaries where accepting ambiguous data could change scientific meaning.

- Invalid IDs, digests, schemas, sizes, and media types fail before storage or execution.
- Snapshot loading rematerializes identity rather than trusting caller-supplied IDs.
- Patch construction rejects duplicate assignments and out-of-domain values.
- Trajectory loading verifies sequence, episode ownership, and every event payload.
- Observation construction requires construct, instrument, and protocol to match the epoch.
- Paired estimation reports missing pairs instead of silently dropping them.
- Journal append requires the expected version; verification checks sequence, previous digest, event digest, and event ID.
- Queue completion requires the active lease ID.
- Budget reserve fails atomically when availability is insufficient.
- Query payload previews remain bounded and explain omission.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| `work lease not found` during completion | The lease expired, was reclaimed, or belongs to another worker | Treat the worker as stale; read the current `WorkRecord` and do not retry completion with the old lease ID. |
| Budget reservation fails despite ready work | Other work holds reserved quantities or committed usage exhausted the limit | Inspect `budget.Snapshot`; release abandoned reservations or define a new campaign with an explicit larger limit. |
| Snapshot identity verification fails | Config bytes, schema, system ID, or artifact reference differ from the recorded snapshot | Verify the config artifact and decode it with the exact registered codec; do not rewrite the snapshot ID. |
| Trajectory load fails on an event payload | Referenced bytes are missing or corrupt | Stop analysis and repair custody from an authoritative copy; do not skip the event. |
| Paired estimate reports missing pairs | One arm/case/repeat observation is absent, invalid, or NaN | Reconcile or rerun the missing episode; do not average only the surviving rows. |
| Query event has no payload preview | Artifact is sensitive, unavailable, larger than 256 KiB, or not eligible for JSON preview | Use the reported `payload_reason`; inspect through an authorized artifact workflow if needed. |
| Scheduler resource claims are unavailable | The unenforced fields were removed before v0 freeze | Use the budget ledger for finite consumption. Wait for an enforced worker-capacity design before declaring capacity claims. |
| Heartbeat cannot be found after cleanup | The unused v0 method was removed | Use bounded local episodes; add the complete heartbeat worker protocol when long-running distributed work is implemented. |
| A pre-freeze journal contains nonempty causation or correlation | Those removed values participated in the historical event identity | Validate that journal with its original code/schema version or regenerate disposable local fixtures; nil omitted values retain the same canonical identity shape. |

## See Also

- `../README.md` — Current package map, commands, and deliberate repository boundaries.
- `../examples/numbergame/demo.go` — Executable proof of candidates, leases, retries, budgets, observations, estimates, and decisions.
- `../campaign/event.go` — Authoritative campaign event envelope and identity materialization.
- `../scheduler/types.go` — Current queue vocabulary, including fields recorded in the removal register.
- `../budget/types.go` — Enforced finite-resource accounting contracts.
- `../query/service.go` — Read-only bounded campaign query model.
- `../ttmp/2026/08/25/OPTKIT-005--implement-optkit-004-contracts-judge-conformance-and-layered-configuration/` — Ticket plan, implementation diary, validation evidence, and preserved legacy bootstrap documents.
