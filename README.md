# OptKit clean-slate vertical slice

This repository is an executable, local-first vertical slice of the supplied OptKit architecture. The initial implementation was imported from source revision `1786d1da86c9e03316ed71336bbc993fa30531f0`; OPTKIT-002 records the archive checksum, collision policy, and validation evidence. It begins as one Go module with enforced package boundaries and no compatibility layer.

It is intentionally deeper than a collection of interfaces: the included number-game campaign creates typed snapshots and a patch, expands a complete-block trial, leases durable work from SQLite, records artifact-backed trajectories and measurement epochs, restarts the process twice, analyzes paired evidence, records a decision, completes the campaign, and rebuilds its overview solely from the journal.

## Implemented scope

- canonical semantic records, typed identifiers, SHA-256 digests, and schema registration;
- memory and filesystem content-addressed artifact stores with corruption checks;
- typed configuration variables, lawful lenses, snapshots, patches, and candidates;
- episode trajectories with ordered events, causal spans, and artifact-backed payloads;
- immutable measurement epochs and typed observations;
- deterministic complete-block trial expansion and paired estimates with explicit missing-pair failure;
- append-only, hash-chained campaign journal in SQLite;
- durable SQLite work queue with leases, retries, expiry reclaim, and idempotent terminal commits;
- finite integer resource budgets with atomic reservations, releases, committed actual usage, and visible overage custody;
- pure reducers and rebuildable campaign overview projection;
- command-line demo, campaign inspection, campaign verification, and artifact verification;
- production-import architecture tests that enforce the current package DAG;
- a deterministic fake judge in the toy system, used where an external provider would otherwise be required.

The local profile is SQLite in WAL mode plus a filesystem CAS. Control metadata and large immutable evidence remain separate.

## Requirements

- Go 1.23 or newer;
- a C toolchain, SQLite headers, and the SQLite shared library for local persistence.

This implementation uses a narrow CGO binding over the system SQLite library because the build environment could not download a third-party driver. `CGO_ENABLED=0 go test ./...` still validates all driver-independent packages; opening the SQLite profile without CGO returns an explicit error.

## Quick start

```bash
make check
go run ./cmd/optkit demo --store ./tmp/demo --reset
```

The demo prints a JSON summary containing the campaign ID, paired estimate, decision, terminal projection, sample trajectory reference, and local paths. Inspect or verify the persisted campaign with that ID:

```bash
go run ./cmd/optkit campaign inspect \
  --store ./tmp/demo \
  --id campaign:...

go run ./cmd/optkit campaign verify \
  --store ./tmp/demo \
  --id campaign:...
```

The first command folds the full control history and returns the current budget snapshot, an event-kind count, and a bounded event tail. The second verifies the journal hash chain and every directly referenced control-event payload.

## Package map

- `record`, `artifact`: canonical identity and immutable byte custody;
- `space`: domains, lenses, variables, snapshots, patches, and candidates;
- `episode`: executable results, event writer, trajectories, usage, and failures;
- `measure`: epochs and observations;
- `experiment`: datasets, complete-block trials, episode expansion, and paired analysis;
- `campaign`: commands, control events, reducer, and journal interfaces;
- `budget`: finite resource limits, reservations, actual usage, and conservation snapshots;
- `scheduler`: durable-work vocabulary and queue interface;
- `store/sqlite`: local journal and queue implementation;
- `projection`: read models rebuilt from authoritative events;
- `local`: composition root for SQLite plus filesystem CAS;
- `examples/numbergame`: end-to-end proof system and deterministic fake instrument;
- `cmd/optkit`: operator command surface.

## Verification

```bash
make check          # formatting, vet, CGO tests, and no-CGO tests
make race           # race detector over the complete module
make demo           # reset and run ./tmp/demo
```

The field-level ownership and lifecycle reference is [`docs/01-optkit-records-artifacts-and-control-model.md`](docs/01-optkit-records-artifacts-and-control-model.md). Historical implementation diaries, ADRs, and journals live in ticket workspaces under `ttmp/`; they are not product documentation.

## Deliberate boundaries

Optkit is a domain-neutral control plane. No Optkit package or test may import RagKit, RagOpt, Judgekit, Coinvault, or RAG-TTC; `internal/boundary/boundary_test.go` enforces this over the module's direct imports. Products compose domain behavior through `system.Factory` and `system.Prepared`, while snapshots and queues retain only canonical configuration and artifact references.

RAG algorithms remain in RagKit, evaluator semantics remain in Judgekit, and product policy remains in product repositories. The embedded web application is a read-only query plane over Optkit campaigns; it does not make Optkit own product-specific projections or browser mutation workflows.

This is a substantial foundation and durable vertical slice, not the complete product-porting program. PostgreSQL/S3 adapters, production deployment, hierarchical organization/stage budgets, provider pricing catalogs, and adaptive optimization remain named follow-on slices rather than mocked as completed work.
