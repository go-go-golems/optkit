# Diary

## Goal

This diary records the construction of a clean-slate OptKit vertical slice from the supplied architecture. It documents code changes, commands, failures, design decisions, review risks, and the exact point from which another contributor can continue.

## Step 1: Establish the clean-slate repository and implementation boundary

**Commit (repository):** `0ba8f9cb7d7b18228c49320084cc5a2126fbe69f` — "Initialize clean-slate OptKit vertical slice"

The workspace contained the architecture PDF and the diary workflow, but no application or library repository. I therefore initialized a new Go module and recorded an explicit vertical-slice boundary rather than pretending to port unavailable Coinvault or RAG-TTC code.

The initial scope follows the guide's local profile and package build order: immutable records and artifacts first, then typed configuration and episodes, followed by SQLite-backed campaign control. The toy number-game system will prove the end-to-end path while external providers and judges remain deterministic fakes.

### What I did

- Created `/mnt/data/optkit` as a Git repository.
- Initialized module `github.com/go-go-golems/optkit`.
- Added the repository README, `.gitignore`, and ADR 0001.
- Read and rendered the architecture sections covering the local SQLite/filesystem profile and Slices A-G.
- Read the attached diary workflow and adopted its mandatory step sections.

### Why

- There was no existing repository to modify.
- The architecture explicitly prefers one clean module and vertical slices over compatibility scaffolding or empty package trees.
- SQLite plus filesystem CAS is the specified local mode and can exercise durable semantics without external services.

### What worked

- Go 1.23.2 and Git 2.47.3 are available.
- The architecture PDF rendered successfully with `pdftoppm` for visual verification.

### What didn't work

- `docmgr` is not installed. `command -v docmgr` returned no path, so diary and changelog bookkeeping are maintained directly in Git-tracked Markdown.
- The browser PDF screenshot API could not resolve the uploaded file reference. The exact error was `Unable to resolve open call ... due to invalid ref_id argument`; local PDF rendering was used instead.

### What I learned

- The absence of product source makes a toy vertical slice the only honest implementation target.
- The local profile is intentionally relational for control data and content-addressed for large immutable evidence; those two stores should remain distinct from the first commit.

### What was tricky to build

- Defining a useful implementation boundary from a broad architecture without reducing the work to empty interfaces.

### What warrants a second pair of eyes

- Confirm that the selected vertical slice is sufficiently deep to validate the ontology while remaining honest about excluded product ports and UI work.

### What should be done in the future

- Product ports must begin only after their repositories and semantic fixtures are available.
- A later UI slice should consume projections rather than reading SQLite tables or artifact paths directly.

### Code review instructions

- Start with `README.md` and `docs/adr/0001-clean-slate-local-vertical-slice.md`.
- Compare the declared scope with the package graph after subsequent steps.

### Technical details

- Module: `github.com/go-go-golems/optkit`
- Local store target: `.optkit/optkit.db` plus `.optkit/artifacts/sha256/...`


## Step 2: Implement canonical records and immutable artifact custody

This step established the bottom of the package graph. Semantic records now have a single digest and canonical-JSON implementation, schema registration rejects conflicting definitions, and all large immutable bytes flow through one artifact interface.

The filesystem store deliberately uses an atomic hard-link publication step instead of ordinary rename semantics. That choice prevents a concurrent writer from replacing an existing object and lets the store fail closed when bytes under a known digest are corrupt.

**Commit (code):** `34e90b0bfb7029ecaff490bf1d2bad384454d078` — "Add canonical records and artifact stores"

### What I did

- Added `record.Digest`, SHA-256 parsing/validation, semantic digest composition, namespaced IDs, canonical JSON, and a concurrency-safe schema registry.
- Added the shared `artifact.Ref` and `artifact.Store` interfaces.
- Implemented `artifact/memory` for focused tests.
- Implemented `artifact/filesystem` with sharded digest paths, temp-file writes, fsync, atomic no-replace publication, and independent verification.
- Added tests for canonical map ordering, schema identity, idempotent registration, common store behavior, concurrent identical writes, and corruption custody.
- Ran:
  - `gofmt -w record artifact`
  - `go test ./record ./artifact/... -count=1`
  - `go vet ./record ./artifact/...`

### Why

- Every later snapshot, trajectory, observation, event payload, and report requires one stable identity vocabulary.
- Artifact storage must be idempotent while refusing to overwrite corrupt content under a claimed digest.
- The memory implementation keeps unit tests cheap; the filesystem implementation is the real local data plane.

### What worked

- All record and artifact tests passed.
- Sixteen concurrent writers of identical bytes converged on one verified object.
- After deliberate on-disk corruption, both verification and a later identical `Put` failed instead of repairing or replacing the object silently.

### What didn't work

- N/A. The first implementation and test run passed without a compile or behavioral failure.

### What I learned

- Unix `os.Rename` can replace an existing destination, which is the wrong default for immutable CAS custody. A same-filesystem hard link gives atomic create-if-absent behavior and exposes the concurrent-writer case cleanly.
- Artifact metadata such as media type and sensitivity belongs in references and owning records; byte identity remains the digest of exact stored content.

### What was tricky to build

- Preserving cancellation while streaming one input simultaneously to disk and SHA-256 without buffering the entire artifact.
- Handling the race between multiple publishers without allowing replacement of the winning bytes.
- Distinguishing semantic record identity (`schema || canonical payload`) from exact artifact byte identity.

### What warrants a second pair of eyes

- Review `record.CanonicalJSON`. It relies on Go's stable key ordering and schema-defined omitted/zero semantics; before v1, the project should lock numeric, time, and custom-marshaler rules with a larger golden corpus.
- Review `artifact/filesystem.Store.Put`, especially hard-link behavior and directory fsync portability on non-Linux filesystems.
- Confirm the sensitivity labels are sufficient for the first product fixtures; they are classification metadata, not authorization enforcement.

### What should be done in the future

- Add authorization wrappers before sensitive product artifacts are exposed through an API.
- Extend canonical golden tests when decimal, time, and schema-evolution choices are made.
- Add mark-and-sweep reachability only after authoritative roots exist; do not infer liveness from filesystem contents.

### Code review instructions

- Start in `record/canonical.go` and `record/digest.go`.
- Review the CAS invariant in `artifact/filesystem/store.go`, then run `go test ./record ./artifact/... -count=1`.
- Run `go test -race ./artifact/... -count=1` to stress concurrent publication.

### Technical details

- Artifact path: `<root>/sha256/<hex[0:2]>/<hex[2:4]>/<full-hex>`.
- Publication protocol: stream to `<root>/tmp`, hash and fsync, hard-link to final path, chmod read-only, fsync parent directory.
- Existing content is accepted only after size and digest verification.

## Step 3: Add typed configuration, trajectories, epochs, and complete-block trials

This step turned immutable bytes into executable experimental objects. A product can now define a typed configuration and variables, materialize a complete snapshot, apply a canonical patch, record a causal episode trajectory, emit epoch-bound observations, and expand a complete-block trial into deterministic episode identities.

The implementation deliberately keeps typed authoring above type-erased records. Snapshot and patch records persist IDs and artifact references, while the generic `Snapshot[C]` and `Patch[C]` carry typed values only inside the authoring/runtime process.

**Commit (code):** `869bed649e272d6b77823f264d766f1a923eb590` — "Add configuration, episode, measurement, and trial model"

### What I did

- Added domains, JSON codecs, lawful lenses, variable descriptors, typed snapshots, canonical patches, candidates, and patch provenance records under `space/`.
- Added artifact JSON read/decode helpers.
- Added an artifact-backed episode writer with monotonic sequence numbers, span-parent validation, payload references, trajectory sealing, and trajectory verification.
- Added measurement epoch and immutable observation identities.
- Added immutable datasets, complete-block trial plans, deterministic episode expansion, and a paired-mean estimator that refuses missing pairs.
- Added tests for lens laws, disjoint assignment commutation, conflict rejection, domain validation, trajectory sequence/parent rules, epoch isolation, deterministic trial expansion, and missing-pair custody.
- Ran:
  - `gofmt -w artifact space episode measure experiment`
  - `go test ./artifact/... ./space ./episode ./measure ./experiment -count=1`
  - `go vet ./space ./episode ./measure ./experiment`
  - `go test -race ./space ./episode ./measure ./experiment -count=1`

### Why

- The architecture's central claim is that optimization acts on complete typed systems and observable trajectories, not isolated strings or scalar metric maps.
- Patch support and experimental design must remain general; one-coordinate and paired experiments are campaign policies, not hard-coded ontology limits.
- Measurement epochs must alter identity so historical trajectories can be remeasured without rewriting previous observations.

### What worked

- Reversing the order of two disjoint assignments produced the same patch and child snapshot IDs.
- Conflicting assignments and out-of-domain values were rejected before a child snapshot was accepted.
- Trajectory loading verifies every referenced event payload.
- Protocol changes produced different measurement epoch and observation IDs.
- Complete-block expansion produced stable episode IDs and seeds across repeated expansion.
- The race detector passed for all new packages.

### What didn't work

- The first compilation failed because `Snapshot[C]` embedded `SnapshotRecord.Config artifact.Ref` and also declared a typed field named `Config C`. The compiler resolved `materialized.Config` to the typed value, producing:

  ```text
  # github.com/go-go-golems/optkit/space [github.com/go-go-golems/optkit/space.test]
  space/snapshot.go:80:54: materialized.Config.Digest undefined (type C has no field or method Digest)
  space/snapshot.go:81:137: materialized.Config.Digest undefined (type C has no field or method Digest)
  ```

  I first qualified the embedded field, then removed the ambiguity entirely by renaming the typed field to `Snapshot.Value`.

### What I learned

- Persisted record fields and typed runtime values need visibly distinct names. Embedded records are convenient, but field shadowing can make generic code misleading even when the compiler catches it.
- Deterministic expansion requires the trial, snapshot, input artifact, repeat, seed, and protocol to participate in the episode key; worker and scheduling details remain excluded.
- A patch builder can preserve typed values while storing only canonical value artifacts and erased assignment records.

### What was tricky to build

- Go methods cannot introduce their own type parameters, so the typed patch operation is a generic free function: `space.Set(builder, variable, value)`.
- Event payload artifacts must be written before the final event envelope can be identified; failed event validation can therefore leave unreferenced but immutable bytes, which later garbage collection must handle through reachability rather than mutation.
- Observation identity must include the epoch and evidence digests but exclude non-semantic creation time.

### What warrants a second pair of eyes

- Review `space.PatchBuilder.Build` for the exact distinction between expected-old-value digest, canonical new-value artifact, child snapshot identity, and patch identity.
- Review `episode.Writer.Emit`: it is concurrency-safe, but a race between payload attachment and sealing can leave an orphan artifact. This is safe for custody but should be considered when GC is implemented.
- Review the v0 canonical JSON codec's strict decoding and the choice to use decimal strings in observation values while the statistical helper currently exposes `float64`.
- Confirm whether candidate identity should continue to include hypothesis prose or be split into proposal and materialization entities before public APIs stabilize.

### What should be done in the future

- Add a registered, type-erased system descriptor only when the campaign worker needs to execute more than the toy system.
- Add explicit estimator identities and intervals when the first product fixture defines required statistical semantics.
- Add redaction hooks before trajectories can contain sensitive product payloads.
- Add exact-epoch query enforcement in the SQLite analysis path; the current in-memory estimator accepts already-filtered rows.

### Code review instructions

- Start with `space/snapshot.go`, `space/patch.go`, and `space/space_test.go`.
- Follow one episode through `episode/writer.go` and `episode/writer_test.go`.
- Check epoch identity in `measure/epoch.go` and missingness custody in `experiment/estimate.go`.
- Validate with `go test -race ./space ./episode ./measure ./experiment -count=1`.

### Technical details

- Snapshot identity includes system ID, configuration schema ID, and canonical configuration artifact digest.
- Patch identity includes base snapshot, sorted assignments, and verified child snapshot.
- Episode specs are ordered case -> repeat -> arm, with deterministic seeds derived from trial/case/repeat.
- Paired analysis returns an error when any required pair is missing or invalid; it does not silently shrink the denominator.

## Step 4: Make campaign control and work leasing durable in SQLite

This step introduced the first authoritative control plane. Campaign facts are appended to a hash-chained SQLite journal under optimistic expected-version checks, reduced through a pure state machine, and rendered through a rebuildable overview. Durable work is separately enqueued, leased, heartbeated, retried, reclaimed after expiry, and committed exactly once.

SQLite is real rather than mocked. Because the environment could not download a Go driver, I implemented a narrow internal CGO binding over the installed system SQLite library. The binding serializes access through one full-mutex connection and exposes only scripts, parameterized statements, materialized queries, and `BEGIN IMMEDIATE` transactions needed by the store.

**Commit (code):** `c91382fc15747f13c52b4d117561d1fc93272707` — "Add SQLite campaign journal and durable lease queue"

### What I did

- Added campaign control-event envelopes, semantic event IDs, previous-digest chaining, journal interfaces, lifecycle/episode reducers, and an idempotent command controller.
- Added a rebuildable campaign overview projection.
- Added scheduler work items, semantic work IDs, leases, retry/terminal failures, results, and queue interfaces.
- Added `internal/sqlite`, a minimal SQLite C API wrapper using the installed `libsqlite3`.
- Added `store/sqlite` migrations and implementations for:
  - campaign heads, control events, and command results;
  - optimistic append and command idempotency;
  - journal read/head/hash verification;
  - durable work enqueue, lease, heartbeat, completion, failure, backoff, expiry reclaim, and lookup.
- Added crash/reopen, tamper, version-conflict, retry, semantic-work-conflict, and terminal-idempotency tests.
- Added a non-CGO stub so the rest of the module still compiles with `CGO_ENABLED=0`; opening SQLite then returns an explicit error.
- Ran:
  - `gofmt -w .`
  - `go vet ./...`
  - `CGO_ENABLED=1 go test ./... -count=1`
  - `CGO_ENABLED=1 go test -race ./campaign ./internal/sqlite ./scheduler ./store/sqlite -count=1`
  - `CGO_ENABLED=0 go test ./... -count=1`

### Why

- A campaign must outlive controller and worker processes; mutable directories or in-memory queues cannot prove recovery, idempotency, or event ancestry.
- SQLite in WAL mode plus a filesystem CAS is the architecture's specified local profile.
- Control events and high-volume trajectory payloads remain separate: SQLite stores relational custody and references, while the CAS stores immutable evidence bytes.
- Work leases require durable attempt history and terminal-result conflict detection independently of campaign projection state.

### What worked

- Closing the database after a work lease, advancing the clock beyond expiry, reopening, reclaiming, and leasing again produced attempt 2 with a new lease ID.
- Completing the recovered lease twice with the same result was idempotent; a different result for the same terminal lease was rejected.
- Retryable failure respected `earliest_start` backoff; terminal failure remained visible.
- Reusing a campaign command ID returned the original event even after the campaign had advanced to a later state.
- Deliberately altering persisted event tags caused journal verification to fail as corrupt.
- WAL mode was verified directly in the SQLite wrapper test.
- All concurrency-sensitive packages passed the race detector.

### What didn't work

- Fetching the preferred pure-Go SQLite driver failed because outbound module access is blocked. Exact command and error:

  ```text
  $ go get modernc.org/sqlite@v1.36.3
  go: modernc.org/sqlite@v1.36.3: Get "https://proxy.golang.org/modernc.org/sqlite/@v/v1.36.3.info": dial tcp: lookup proxy.golang.org on 168.63.129.16:53: read udp 172.26.36.65:49048->168.63.129.16:53: read: connection refused
  ```

- `sqlite3` CLI is not installed. The environment does provide `/usr/lib/x86_64-linux-gnu/libsqlite3.so.0.8.6`, `/usr/include/sqlite3.h`, and Python reports SQLite `3.46.1`, so the implementation linked directly against the system library.

### What I learned

- Command idempotency must be checked before current-state validation. Otherwise retrying an old `CreateCampaign` command after `PlanCompiled` would be rejected even though the journal already accepted it.
- A pure reducer cannot depend on loading event payload artifacts. Minimal routing facts such as episode subject and retryability therefore live in the control-event envelope/tags while rich detail remains artifact-backed.
- A single SQLite connection with `BEGIN IMMEDIATE` gives a simple local single-writer authority while WAL still permits external readers.
- Terminal idempotency is a comparison of semantic result identity, not merely acceptance of any repeated completion call.

### What was tricky to build

- Binding text/blob parameters safely with `SQLITE_TRANSIENT`, finalizing every prepared statement, and preserving rollback behavior without importing a third-party driver.
- Returning a prior command batch inside the same transaction before applying an expected-version check.
- Reconstructing typed event and artifact references from nullable relational columns without letting corrupt rows appear valid.
- Keeping queue reclaim, selection, attempt increment, and lease assignment in one transaction.
- Ensuring non-CGO builds compile while failing only when SQLite is actually opened.

### What warrants a second pair of eyes

- The CGO wrapper in `internal/sqlite/sqlite_cgo.go` is the highest-risk code in this step. Review C allocation/free pairs, `SQLITE_TRANSIENT` binding, statement finalization, transaction rollback, and context interruption behavior.
- Review the choice to serialize all operations through one connection. It is correct for the local profile but intentionally not a production-throughput design.
- Review `campaign.ValidateTransition`, especially which post-campaign facts should eventually be accepted through declared remeasurement workflows.
- Review the control-event envelope additions (`Subject`, `Tags`) and decide before v1 whether they remain generic envelope fields or become small typed summary unions.
- Review lease completion after expiry: the current queue accepts a completion only while the row still carries that lease; once reclaimed and re-leased, the old lease cannot win.

### What should be done in the future

- Replace or supplement the native binding with a pinned, audited driver when dependency access and product CGO constraints are known. Preserve store semantics and tests during that change.
- Add signed checkpoints only after profiling shows replay cost matters; checkpoints must remain caches.
- Add budget reservation/reconciliation tables before executing billable provider work.
- Add projection offsets/tables when the read API exists; the current overview proves rebuild semantics without becoming a second authority.
- Add PostgreSQL only after the local controller behavior is exercised by a real product slice.

### Code review instructions

- Start with `campaign/reducer.go` and `campaign/controller.go` for accepted-state semantics.
- Review journal transactions in `store/sqlite/journal.go` and lease transactions in `store/sqlite/queue.go`.
- Review the native boundary in `internal/sqlite/sqlite_cgo.go` last, with the tests open beside it.
- Validate with:

  ```bash
  CGO_ENABLED=1 go test -race ./campaign ./internal/sqlite ./scheduler ./store/sqlite -count=1
  CGO_ENABLED=0 go test ./... -count=1
  ```

### Technical details

- SQLite pragmas: WAL journal, foreign keys on, 5-second busy timeout, synchronous NORMAL.
- Journal append transaction: create/lock head -> resolve duplicate command -> compare expected version -> insert chained events -> CAS-update head -> record command range.
- Queue lease transaction: reclaim expired -> select bounded ready work -> assign unique lease, expiry, and incremented attempt.
- Work semantic uniqueness: `(campaign, kind, semantic_key)`; payload mismatch under the same key is a visible conflict.

## Step 5: Prove the complete local campaign path with a restartable toy system

This step connects the previously isolated foundations into one executable campaign. The number-game proof system is deliberately small, but it uses the same semantic path expected of a product port: typed configuration, a lawful variable, immutable baseline/challenger snapshots, a patch and candidate hypothesis, complete-block trial expansion, durable episode work, full trajectories, intervention checks, deterministic and judge-like measurements, paired analysis, an immutable decision, campaign completion, and replayed projections.

**Commit (code):** `01f43cd0f736743581cee60d4029088f2098c583` — "Add restartable end-to-end number-game campaign"

### What I did

- Added `local.Profile`, the composition root for `<root>/optkit.db` plus `<root>/artifacts`.
- Refined episode execution so a product executable returns `episode.RunResult`; the worker seals the trajectory and combines both into the persisted `episode.Result`. This avoids making product code responsible for artifact-custody finalization.
- Added `examples/numbergame` with:
  - a typed `Config` containing `Multiplier` and seeded `Noise`;
  - lawful variables/domains and a deterministic execution clock;
  - an executable that emits input, effective-configuration, noise, output, and terminal trajectory events;
  - an intervention probe that reads effective configuration from the sealed trajectory;
  - deterministic absolute-error measurement;
  - a deterministic fake judge with an explicit epoch and sealed assessment artifact;
  - a complete local campaign runner.
- Added the full demo path:
  1. materialize multiplier-2 baseline;
  2. apply a multiplier-3 patch;
  3. register candidate, patch, child snapshot, dataset, and trial facts;
  4. expand four cases across baseline and challenger into eight episodes;
  5. enqueue and lease work in batches of two;
  6. execute, seal, measure, complete queue work, and append control facts;
  7. close and reopen SQLite before analysis;
  8. rebuild paired rows from `EpisodeCompleted` payloads and observation artifacts;
  9. compute the paired target-accuracy delta;
  10. apply a small lexicographic decision (intervention validity before target gain);
  11. complete the campaign;
  12. close and reopen again, verify the chain, and rebuild the overview.
- Added `cmd/optkit` commands:
  - `demo`;
  - `campaign inspect`;
  - `campaign verify`;
  - `artifact verify`.
- Added an import-graph architecture test that enforces bottom-level package ownership and prevents core packages from importing concrete stores, commands, examples, or projections.
- Added integration tests for campaign replay, SQLite file identity, trajectory reopening, CLI output, and verification.
- Updated the README and Makefile with CGO and no-CGO validation paths.
- Ran:
  - `make check`;
  - `CGO_ENABLED=1 go test -race ./... -count=1`;
  - `CGO_ENABLED=0 go test ./... -count=1`;
  - a manual demo, inspect, and verify sequence against `/mnt/data/optkit-demo-store`.

### Why

- The architecture explicitly requires a toy campaign before product ports so ontology, storage, recovery, and projections can be tested cheaply.
- Interfaces without a complete run would not prove that identities, event transitions, artifacts, queue state, and analysis actually compose.
- Restarting before analysis and after completion is stronger evidence than retaining in-memory structs and merely asserting that persistence exists.
- The fake judge preserves the instrument/epoch/assessment boundary without pretending an unavailable provider integration was implemented.

### What worked

- The integration test produced eight terminal episodes, twenty-four observations, one estimate, one decision, and a terminal campaign projection.
- A representative manual run produced:

  ```text
  campaign events: 65
  episodes: 8 total / 8 completed / 0 active / 0 failed
  paired samples: 4
  paired accuracy delta: 0.71041675
  decision: eligible
  interventions exercised: true
  ```

- The event-kind distribution in that run was:

  ```text
  EpisodeScheduled       8
  EpisodeLeaseGranted    8
  EpisodeAttemptStarted  8
  EpisodeCompleted       8
  ObservationRecorded   24
  EstimateRecorded       1
  DecisionRecorded       1
  CampaignCompleted      1
  ```

  plus campaign, plan, candidate, snapshot, and trial setup events.
- The SQLite file began with the expected `SQLite format 3` header.
- The campaign hash chain and every directly referenced event payload verified after reopening.
- The sample trajectory reopened independently from the filesystem CAS and contained the expected five ordered events.
- `campaign inspect` rebuilt status from all journal facts and returned a bounded tail without reading mutable campaign-status columns.
- The architecture test passed and now fails CI if foundation packages import upward.
- Full race tests passed.
- Driver-independent packages compile and test with `CGO_ENABLED=0`.

### What didn't work

- The first provisional worker loop used `lease.Item.SemanticKey` as the subject of `EpisodeLeaseGranted`. That was incorrect: the reducer tracks episode lifecycle by semantic episode ID, while the work digest is only queue identity. The provisional code was deliberately stopped with:

  ```text
  internal demo invariant: lease event subject must be episode ID
  ```

  The worker now decodes `EpisodeWork` before appending lease/attempt facts and uses `work.Spec.ID` throughout the episode state machine.

- The first complete no-CGO test run failed because the number-game SQLite integration test was not build-gated. Exact failure:

  ```text
  --- FAIL: TestRunDemoPersistsAndReplaysCompleteCampaign (0.00s)
      demo_test.go:20: SQLite support requires CGO
  FAIL
  github.com/go-go-golems/optkit/examples/numbergame
  ```

  The SQLite integration tests now use the `cgo` build constraint. No-CGO builds still compile the demo command and packages, but do not claim to run the unavailable SQLite backend.

### What I learned

- Durable-work identity and domain-subject identity are related but not interchangeable. The queue needs a semantic key for duplicate suppression; the campaign reducer needs the episode ID for lifecycle reconstruction.
- The executable should not construct a full `episode.Result` before the trajectory is sealed. Splitting `RunResult` from the custody-complete `Result` makes the ownership boundary explicit.
- A campaign specification must retain a reachable baseline snapshot record, not just its ID. Otherwise the baseline configuration artifact is not rooted by journal evidence.
- Candidate evidence must retain the patch record with the proposal. A candidate containing only a patch ID is not independently inspectable without another index.
- Completing queue work before appending `EpisodeCompleted` creates a real reconciliation boundary. The demo handles the happy path; a production controller needs a reconciler that ingests completed queue rows after a crash between those two commits.
- A test matrix must distinguish “module compiles without CGO” from “SQLite integration executes without CGO.” Treating those as the same promise caused the caught regression.

### What was tricky to build

- Keeping every persisted object reachable after two restarts: campaign spec → baseline snapshot/config, candidate proposal → patch/child snapshot, trial → dataset/cases, completion → result/trajectory/observations/assessment.
- Ordering the queue commit and journal facts without pretending they share one transaction across metadata responsibilities.
- Reconstructing paired rows only from sealed event payloads while preserving case, arm, repeat, epoch, status, and intervention validity.
- Making test time deterministic while leaving semantic IDs independent of worker timing and queue order.
- Keeping CLI verification useful without turning it into a hidden second projection or artifact crawler.

### What warrants a second pair of eyes

- Review `examples/numbergame/demo.go` as controller/worker pseudocode made executable. In particular, review the queue-complete/journal-append reconciliation gap and whether a dedicated completion-ingestion command should be the next slice.
- Review the `RunResult`/`Result` API split in `episode/types.go`; it clarifies custody but is a public shape that should be settled before product adapters use it.
- Review reachability of every artifact rooted by the campaign journal. The current verifier checks direct event payloads, while nested-reference traversal is still a future evidence-export feature.
- Review whether intervention observations should be appended after or atomically summarized with `EpisodeCompleted`; the current model preserves detailed facts but uses multiple journal transactions.
- Review the architecture-test rule set. It intentionally enforces only durable current boundaries and should be expanded as `work`, `flow`, `judge`, `rag`, and server packages appear.
- Review the fake judge naming and output. It is intentionally deterministic and must never be presented as evidence of real model-judge integration.

### What should be done in the future

- Add controller reconciliation that converts already-completed durable queue rows into idempotent campaign completion commands after process loss.
- Add budget reservation, actual usage commitment, and release records before any billable provider work.
- Add artifact-graph verification and export manifests that traverse nested references from retained roots.
- Add projection offsets and disposable SQLite projection tables when a read API/UI is implemented.
- Add authorized history/exposure records before reflective search or hidden datasets.
- Implement manual/coordinate search as a first-class strategy rather than embedding the one patch directly in the example.
- Port a deterministic RAG retrieval slice when the RAG-TTC repository and fixture corpus are supplied.
- Port Coinvault only after its direct application service and retained treatment/constraint fixtures are available.

### Code review instructions

- Run the vertical slice first:

  ```bash
  make check
  make race
  make demo
  ```

- Start review at `examples/numbergame/demo.go` and follow every stored reference into `space`, `episode`, `measure`, `experiment`, `campaign`, `scheduler`, and `store/sqlite`.
- Inspect `examples/numbergame/demo_test.go` to see the restart and independent-trajectory assertions.
- Inspect `cmd/optkit/main.go` for the current operator surface and `internal/archtest/architecture_test.go` for enforced package direction.
- For a manual run, capture the `campaign` field from demo JSON, then run both `campaign inspect` and `campaign verify` against the same store.

### Technical details

- Demo trial: two arms × four cases × one repeat = eight semantic episodes.
- Episode trajectory: five events (`input.received`, effective config, noise sample, output, terminal).
- Per-episode observations: intervention status, absolute error, and fake-judge target accuracy.
- Paired score: `1 / (1 + absolute_error)`; the multiplier-3 challenger exactly matches all targets in the fixture.
- Campaign completion is accepted only after the reducer sees every episode in a terminal state.
- Manual demonstration artifacts:
  - summary: `/mnt/data/demo-summary.json`;
  - inspection: `/mnt/data/demo-inspect.json`;
  - local profile: `/mnt/data/optkit-demo-store`.

## Step 6: Enforce finite resource budgets with durable reservations and actual-use custody

The first complete campaign declared a budget but did not enforce it. This step turns resource limits into an authoritative SQLite ledger. Work now receives capacity before it can enter the runnable set, successful execution commits actual usage, failed/admission-aborted work releases capacity, and unavoidable overage remains recorded even when it violates the cap.

**Commit (code):** `dc505c822ee49f97f9ddae1da77ba0ae8423208b` — "Add durable campaign budget reservations and usage"

### What I did

- Added the `budget` package with:
  - normalized resource names and integer quantities;
  - campaign limits;
  - content-derived reservation IDs;
  - reserved, committed, and released states;
  - actual usage and overage flags;
  - campaign resource snapshots;
  - a storage-neutral ledger interface.
- Added `record.ReservationID`.
- Extended the SQLite schema with `budget_limits` and `budget_reservations`.
- Implemented SQLite budget operations:
  - idempotent limit definition with conflict detection;
  - atomic reservation under `BEGIN IMMEDIATE`;
  - per-work idempotency and conflicting-claim rejection;
  - finite-cap admission checks;
  - idempotent actual-use commit;
  - release of unused reservations;
  - lookup by reservation or work ID;
  - transactionally consistent resource snapshots;
  - visible overage rather than rejected/lost actual usage.
- Added campaign event kinds:
  - `BudgetReserved`;
  - `UsageCommitted`;
  - `BudgetReleased`.
- Allowed actual usage/release custody after a campaign reaches a terminal state, while rejecting new reservations after completion.
- Integrated budgets into the number-game campaign:
  - limits of eight episode units and eight operation units;
  - two one-unit reservations for each of eight work items;
  - committed actual usage from `episode.Result.Usage`;
  - released reservations on failed work or enqueue rollback;
  - budget validity as the first decision condition;
  - budget snapshot in demo and campaign-inspection JSON.
- Added tests for normalized identity, duplicate/negative quantity rejection, define/reopen behavior, reserve/commit/release idempotency, conflicting claims, insufficient capacity, concurrent near-cap reservations, unavoidable overage, terminal usage custody, and the integrated campaign ledger.
- Ran:
  - `make check`;
  - `make race`;
  - a manual demo/inspect/verify cycle.

### Why

- Provider calls cannot safely run against a budget that exists only in a campaign spec.
- Reservation and actual usage are different facts. Scheduling must conservatively reserve capacity, while final custody must preserve what the provider actually consumed.
- A provider can report more usage than estimated after the work is already billable. Rejecting that report would make accounting look safe by deleting evidence; the correct state is committed usage plus a visible violation.
- Concurrent schedulers near a cap must not both see the same available capacity.

### What worked

- Two concurrent reservations of six units against a ten-unit cap produced exactly one success and one `budget.ErrInsufficient`; the resulting snapshot retained six reserved and four available.
- Defining the same limits twice was idempotent. Defining a different limit set after initialization returned `budget.ErrLimitConflict`.
- Repeating the same reservation request returned the same content-derived reservation. Reusing a work ID with different claims returned `budget.ErrReservationConflict`.
- Repeating an identical actual-use commit was idempotent. A different second actual-use value was rejected as a conflict.
- A five-unit reservation committed with eleven actual units was accepted, marked as overage, and projected as eleven committed, negative one available, and `violated: true`.
- A released reservation stopped consuming capacity and remained released after database reopen.
- The full number-game campaign now records eight `BudgetReserved` and eight `UsageCommitted` events. The manual run produced:

  ```text
  campaign events: 81
  episodes: 8 total / 8 completed
  budget episodes: limit 8 / committed 8 / reserved 0 / available 0
  budget operations: limit 8 / committed 8 / reserved 0 / available 0
  budget violated: false
  paired accuracy delta: 0.71041675
  decision: eligible
  journal/event payload verification: passed
  ```

- Full module tests, race tests, and no-CGO tests passed.

### What didn't work

- During the first integration edit, the demo referenced the release helper before it existed and imported `errors` before using it. Exact compiler output:

  ```text
  examples/numbergame/demo.go:5:2: "errors" imported and not used
  examples/numbergame/demo.go:249:11: undefined: releaseBudgetReservation
  ```

  I completed the release/usage helper functions, used `errors.Is` to distinguish a missing reservation during failure cleanup, and reran the integration test successfully.

### What I learned

- Budget limits, reservations, and usage should not be represented as one mutable counter. Separate immutable/terminal reservation records make retries, releases, actuals, and overage reviewable.
- Content-derived reservation identity over `(campaign, work, normalized claims)` gives useful idempotency without another command nonce.
- “Actual exceeds reservation” and “campaign total exceeds limit” are both overage conditions, but actual usage still belongs in custody.
- Resource snapshots need one database transaction. Reading limits and reservations in separate unlocked queries could produce an internally inconsistent view during a concurrent commit.
- Usage facts may legitimately arrive after cancellation or terminal campaign status. The reducer must distinguish “no new work” from “do not record late billable facts.”
- Reserving capacity and enqueueing work are still two metadata operations in the demo. A production scheduler should either combine them atomically or reconcile orphan reservations and unadmitted scheduled facts after a crash.

### What was tricky to build

- Computing availability from mixed reservation states while releasing unused reserved quantity on commit.
- Making overage commits idempotent and visible without allowing a second conflicting terminal actual.
- Preserving deterministic ordering for resource quantities so identities and conflict comparisons do not depend on map iteration.
- Testing near-cap concurrency through the serialized local SQLite connection while still proving only one transaction can admit capacity.
- Integrating resource events without making the campaign reducer own mutable budget arithmetic.

### What warrants a second pair of eyes

- Review `store/sqlite/budget.go`, especially the reservation transaction, actual-use overage calculation, and snapshot aggregation.
- Review whether zero-unit limits and zero-unit actuals should remain legal in v0.
- Review the current policy that actual resources must be a subset of reserved resources. A real provider adapter may need an explicit “unestimated resource” path that still commits custody and marks violation.
- Review whether budget limit definition should itself have a dedicated journal event or remain part of the immutable campaign-spec payload plus relational execution state.
- Review the event rule allowing `UsageCommitted` and `BudgetReleased` after terminal campaign states; this is intentional for late provider accounting.
- Review the remaining reserve/enqueue and queue-complete/journal-append reconciliation boundaries before connecting paid providers.

### What should be done in the future

- Add a single SQLite admission transaction that combines reservation and work enqueue, or a reconciler that releases orphan reservations safely.
- Add hierarchical scopes (campaign, stage, provider/model) after one concrete need exists.
- Use fixed-point smallest units for money and add a versioned pricing catalog; do not infer cost from floating point.
- Add budget reservation/usage events to a dedicated rebuildable ledger projection and UI panel.
- Add cancellation tests that commit late actual usage and release only genuinely unused reservations.
- Add a terminal `budget_violation` campaign/policy check with explicit evidence when an unavoidable overage occurs.

### Code review instructions

- Start with `budget/types.go`, then review `store/sqlite/budget.go` alongside the new tests at the end of `store/sqlite/store_test.go`.
- Follow one demo work item through reservation in `RunDemo`, actual commit in `executeEpisodeLease`, and the terminal decision budget check.
- Run:

  ```bash
  CGO_ENABLED=1 go test -race ./budget ./campaign ./store/sqlite ./examples/numbergame ./cmd/optkit -count=1
  CGO_ENABLED=0 go test ./... -count=1
  ```

- Inspect demo JSON and confirm both resources have zero reserved, eight committed, and zero available after completion.

### Technical details

- Budget quantities are signed 64-bit integers validated as non-negative at API boundaries.
- Reservation IDs are SHA-256 content IDs over campaign, work, and sorted claims.
- Reservation uniqueness is also enforced by `work_id` in SQLite.
- Admission calculation includes all outstanding requested quantities plus all committed actual quantities.
- Commit calculation removes the current reservation from outstanding usage, adds actuals, and marks overage when actual exceeds its reservation or crosses a campaign limit.
- Released reservations contribute neither reserved nor committed usage.
- Manual demonstration files were refreshed at:
  - `/mnt/data/demo-summary.json`;
  - `/mnt/data/demo-inspect.json`;
  - `/mnt/data/demo-verify.json`;
  - `/mnt/data/optkit-demo-store`.

## Step 7: Verify and package the implementation for independent review

This delivery step verifies that the implementation is reproducible from tracked source rather than relying on the working tree or generated demo state. It also preserves the Git history required to audit the commit hashes referenced throughout this diary.

### What I did

- Confirmed the repository was clean at commit `b458a37bdbbf855a9d5a7b38f48d36750ef67486` before packaging.
- Recorded the final tracked scope at that point:
  - 66 tracked files;
  - 7,581 lines of Go;
  - eleven implementation/journal commits before the delivery entry.
- Created a source-only ZIP with `git archive`.
- Created a Git bundle containing the complete branch/history.
- Created a demonstration archive containing:
  - the SQLite database;
  - filesystem CAS;
  - demo summary JSON;
  - campaign inspection JSON;
  - campaign verification JSON.
- Extracted the source ZIP into a fresh directory and ran both CGO and no-CGO test suites from that extracted copy.
- Generated SHA-256 checksums for every delivery artifact.

### Why

- A working-tree test can accidentally depend on untracked files. Testing the `git archive` proves the tracked repository is sufficient.
- The source ZIP is convenient for review, while the Git bundle preserves the diary's commit ancestry and exact hashes.
- The persisted demo profile lets reviewers inspect real SQLite/CAS output without rerunning the campaign first.
- Delivery checksums make accidental transfer corruption visible.

### What worked

- The extracted source archive passed:

  ```text
  CGO_ENABLED=1 go test ./... -count=1
  CGO_ENABLED=0 go test ./... -count=1
  ```

- The source ZIP, Git bundle, and demo archive were all small enough for direct delivery.
- The demonstration database and JSON files remained internally consistent with the 81-event campaign described in Step 6.

### What didn't work

- The first packaging invocation set the container working directory to the extraction path before that path had been created. The container rejected the call with:

  ```text
  ENOENT: No such file or directory
  ```

  I reran the command from the repository, created and extracted the archive inside the script, changed directory only afterward, and completed both test suites successfully.

### What I learned

- Tool-level working directories are resolved before shell commands run; a script cannot create its own requested initial working directory.
- A source archive and a Git-history bundle serve different review needs and should both be delivered when a diary cites commit hashes.
- Generated SQLite/CAS state is useful evidence but must remain a derived demonstration, not part of the authoritative source repository.

### What was tricky to build

- Avoiding a false reproducibility result by ensuring tests ran from the extracted archive rather than the original repository.
- Packaging the demo database and artifact tree together with the exact JSON outputs that identify its campaign.
- Keeping generated delivery files outside Git so source identity remains stable.

### What warrants a second pair of eyes

- Verify the Git bundle with `git bundle verify` and clone/fetch it into an empty repository if commit-level review is required.
- Compare the source ZIP HEAD with the bundle HEAD before applying additional work.
- Treat `/mnt/data/optkit-demo-store` as non-sensitive synthetic evidence only; future product demonstrations require retention and redaction review.

### What should be done in the future

- Add a release script that performs archive creation, extraction tests, Git-bundle verification, and checksum generation in one reviewed command.
- Sign release checksums when the project has a signing policy.
- Add independent evidence-export verification rather than distributing the raw local profile for product campaigns.

### Code review instructions

- For source-only review, unpack `optkit-implementation-source.zip` and run `make check` plus `make race`.
- For commit-history review, initialize an empty Git repository and fetch from `optkit-implementation-history.bundle`.
- For persisted-state review, unpack `optkit-demo-store.tar.gz`, read `demo-summary.json`, and use the included campaign ID with the CLI.

### Technical details

- Delivery artifacts are generated under `/mnt/data`, not committed into the repository.
- Checksum manifest: `/mnt/data/optkit-delivery-SHA256SUMS.txt`.
- Archive extraction test root: `/mnt/data/_optkit_archive_check/optkit`.
