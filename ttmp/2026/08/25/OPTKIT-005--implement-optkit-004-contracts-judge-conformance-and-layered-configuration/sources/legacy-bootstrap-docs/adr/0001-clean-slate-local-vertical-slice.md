# ADR 0001: Build a local clean-slate vertical slice

- Status: accepted
- Date: 2026-08-24

## Context

The supplied architecture specifies one modular Go module, immutable semantic records and artifacts, typed configurations and episodes, append-only campaign control, SQLite plus a filesystem CAS for local mode, and vertical slices before broad framework construction.

No source repository was supplied. There are therefore no compatibility obligations or product call sites to preserve.

## Decision

Create `github.com/go-go-golems/optkit` as one Go module. Implement the smallest end-to-end local campaign that proves the shared ontology:

1. canonical records and artifact custody;
2. typed snapshots and patches;
3. episode trajectories;
4. measurements and complete-block experiments;
5. hash-chained campaign events;
6. a durable SQLite lease queue;
7. a rebuildable overview projection;
8. one toy campaign and CLI.

Use `modernc.org/sqlite` to avoid a CGO requirement. Use deterministic fakes for judge-like behavior and external providers.

## Consequences

The repository does not contain compatibility adapters, product transport integrations, a web UI, distributed workers, or the Coinvault/RAG-TTC ports. Those require their source repositories and are later vertical slices.
