---
Title: Implementation Diary
Ticket: OPTKIT-013
Status: complete
Topics:
    - architecture
    - design
    - implementation
    - optkit
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/space/binding.go
      Note: Central production implementation recorded in Steps 3 through 8
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/design-doc/04-backend-first-optimization-workbench-program-roadmap.md
      Note: Overall program context that keeps the ticket aligned
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings/design-doc/01-intern-guide-to-optkit-catalogs-domains-bindings-and-candidate-intent.md
      Note: Primary design deliverable whose research and delivery this diary records
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings/various/numbergame-catalog-proof.json
      Note: Fresh end-to-end v2 bundle proof
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings/various/remarkable-implementation-dry-run.log
      Note: Final implementation bundle selection and destination evidence
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings/various/remarkable-implementation-upload.log
      Note: Successful current guide and diary upload receipt
    - Path: repo://optkit/ttmp/2026/08/26/OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings/various/work-slips/00-ticket-plan.log
      Note: Successful plan work-slip print receipt
ExternalSources: []
Summary: Chronological research, design, implementation, validation, and delivery record for OPTKIT-013.
LastUpdated: 2026-08-26T14:20:21.419612535-04:00
WhatFor: ""
WhenToUse: ""
---





# Diary

## Goal

Design the domain-neutral semantic catalog, lossless domains, executable bindings, and candidate intent extension. This diary records research, design decisions, validation, commits, delivery evidence, failures, and continuation instructions.

## Step 1: Investigate the current system and write the intern design guide

This step established the ticket before implementation starts. I read the parent OPTKIT-011 roadmap and architect brief, traced the current source paths that constrain this ticket, and wrote a self-contained intern guide that explains both current behavior and the proposed implementation sequence.

The design preserves the backend-first dependency order. It records concrete APIs, pseudocode, diagrams, file/line references, tests, risks, exclusions, and an exit gate so later work can proceed without reconstructing the program context.

### Prompt Context

**User prompt (verbatim):**

> Ok, create all these tickets, and for each ticket, Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.
>
> commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)
>
> Youc an run experiments and such, store all scripts you write in the scripts/ folder of the relevant ticket.

**Follow-up user prompt (verbatim):**

> add the description of each ticket and make an upfront design doc with the ticket goal and content so that we can keep track as we go on, to not lose the overall picture.

**Follow-up user prompt (verbatim):**

> since each document will require a significant amount of planning / investigation, so we want to make sure we don't forget what the do cis. Once that overview doc is written, update the goal to reference it.

**Assistant interpretation:** Explain how serialized values invoke existing typed Optkit variables without adding a second mutation algebra, then prove the design with numbergame.

**Inferred user intent:** Create the generic framework all RAG coordinates and authoring transports can share.

### What I did

- Inspected Variable, Domain, Lens, Codec, PatchBuilder, Snapshot, Candidate, and numbergame candidate flow.
- Designed one-construction-site descriptor/binding registration and semantic/full catalog identities.
- Specified domain, binding, candidate identity, and numbergame conformance tests.
- Added a substantive ticket index description, implementation tasks, and this strict-format diary.
- Kept program-wide ticket scaffolding/generation scripts under the parent OPTKIT-011 `scripts/` directory.

### Why

- Each ticket requires significant independent investigation, but its contracts must remain aligned with the parent roadmap.
- An intern should understand ownership, runtime behavior, invariants, and test evidence before editing code.
- Up-front acceptance gates prevent frontend or persistence work from outrunning foundational contracts.

### What worked

- Baseline focused tests passed before documentation: `optkit/space`, `optkit/examples/numbergame`, and the relevant RAG-TTC optimization/workbench/campaign/search/specialist packages.
- Current CLI help confirmed the existing `config` and `campaign` Glazed command surfaces.
- Concrete source evidence was sufficient to define this ticket without speculative production changes.

### What didn't work

- N/A for this ticket's own design. During the program-wide index update, the first batch used guessed generated timestamps and exact-text edits failed for OPTKIT-013 through OPTKIT-020; the retry used smaller stable frontmatter blocks and succeeded.

### What I learned

- Type erasure must retain the real Codec[V], Domain[V], and Lens[C,V]; decoding through any/map values would corrupt integer and choice semantics.
- Ticket boundaries are most reliable when each names both its upstream contract and the next consumer.
- Documentation should distinguish observed runtime behavior from proposed APIs and future implementation choices.

### What was tricky to build

- Type erasure must retain the real Codec[V], Domain[V], and Lens[C,V]; decoding through any/map values would corrupt integer and choice semantics.
- The guide had to be self-contained without copying the entire parent architecture. It summarizes prerequisites, then links each claim to the file that owns it.

### What warrants a second pair of eyes

- Ensure descriptors cannot be registered independently from bindings and that documentation identity is separated from semantic identity.
- Review all proposed public schema/API names before implementation makes them expensive to change.

### What should be done in the future

- Implement value specs, catalog construction, typed adapters, candidate intent, and the numbergame proof in the documented order.
- Update this diary immediately when implementation reveals a false assumption or accepted contract change.

### Code review instructions

Start with the ticket's `design-doc/01-*.md`, then inspect these decision-shaping files:

- `optkit/space/variable.go`
- `optkit/space/domain.go`
- `optkit/space/lens.go`
- `optkit/space/patch.go`
- `optkit/examples/numbergame/demo.go`

Validate the design workspace with:

```bash
docmgr validate frontmatter --doc optkit/ttmp/2026/08/26/OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings/reference/01-implementation-diary.md
docmgr doctor --ticket OPTKIT-013 --stale-after 30
```

Run the focused code/test commands listed in the guide before and after implementation.

### Technical details

- Parent program map: `OPTKIT-011/design-doc/04-backend-first-optimization-workbench-program-roadmap.md`.
- Ticket state: `index.md`, `tasks.md`, and `changelog.md` in this workspace.
- No production code behavior changed while writing this step.

## Step 2: Validate, commit, and deliver the guide

This step converted the researched guide from a working document into a reviewed ticket deliverable. The ticket's frontmatter, relations, tasks, and changelog were validated; the documentation was committed in a dependency-coherent batch; and the index, guide, and diary were rendered and uploaded as one reMarkable PDF with a table of contents.

The implementation tasks intentionally remain open. This delivery completes the up-front planning package and gives the future implementer an evidence-backed starting point, not a false claim that production behavior has already changed.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Finish the ticket documentation package with validation, coherent Git history, strict diary evidence, and a ticket-specific reMarkable delivery.

**Inferred user intent:** Make the plan durable and reviewable before implementation starts, while preserving a clear distinction between completed design work and pending code tasks.

**Commit (documentation):** `428f6b8f5391dc851d989364b21a9d727c78cfc0` — "OPTKIT-012-014: design workbench foundations"

### What I did

- Ran `docmgr doctor --ticket OPTKIT-013 --stale-after 30` and obtained `All checks passed`.
- Scanned the index, guide, and diary for generated placeholder sections; none remained.
- Ran focused baseline Go tests before documentation changes; all selected Optkit and RAG-TTC packages passed.
- Ran `git diff --check`/`git diff --cached --check`, corrected whitespace findings, and committed the ticket package.
- Ran the program upload script in `--dry-run` mode, then rendered and uploaded `OPTKIT-013 Catalog and Bindings Guide.pdf`.
- Preserved upload evidence in `various/remarkable-dry-run.log` and `various/remarkable-upload.log`.

### Why

- A detailed guide is only useful when frontmatter, links, task state, and delivery artifacts agree.
- Grouping commits by dependency layer keeps review focused while avoiding one nine-ticket mega-commit.
- The reMarkable bundle lets the architecture be reviewed away from the source tree without losing the index or diary context.

### What worked

- `remarquee` reported: `OK: uploaded OPTKIT-013 Catalog and Bindings Guide.pdf -> /ai/2026/08/26/OPTKIT-013`.
- The upload rendered the real Markdown rather than only performing a path dry-run.
- The ticket retains open implementation tasks and checked documentation/delivery tasks separately.

### What didn't work

- The first staged diff check failed with `new blank line at EOF` in generated changelogs and `trailing whitespace` on blank quoted prompt lines (`+> `) in generated diaries. The diary generator was corrected to emit `>` on blank quote lines, changelog EOFs were normalized, and the second `git diff --check` passed.
- OPTKIT-020 initially referenced nonexistent `repo://optkit/artifact/artifact.go`; the actual package contains `ref.go`, `store.go`, `read.go`, and `helpers.go`. The invalid relation was removed, `artifact/ref.go` was related, and doctor then passed.

### What I learned

- A dry-run validates upload selection and destination but the real upload is the evidence that Pandoc/LaTeX can render the complete guide.
- Generated prose containing blockquotes needs whitespace validation just like source code.
- Keeping implementation tasks open while checking documentation/delivery tasks makes ticket status truthful.

### What was tricky to build

- The uploaded diary necessarily describes the work up to its render time. This final local step records the upload receipt after the PDF has been created; re-uploading with `--force` solely to include its own receipt would overwrite a new document and risk future annotations.
- Cross-ticket commit hashes and ticket-specific upload names had to remain consistent across nine independent workspaces.

### What warrants a second pair of eyes

- Review the proposed APIs and compatibility decisions before implementation; successful document delivery is not architecture acceptance.
- Confirm the reMarkable bundle name and ticket folder are the intended long-term review locations before adding annotations.

### What should be done in the future

- Begin only after upstream entry gates in the OPTKIT-011 roadmap are satisfied.
- During implementation, append new diary steps with exact code commit hashes, failures, commands, and fresh validation evidence.

### Code review instructions

- Start with `index.md`, then read the complete `design-doc/01-*.md`, then this diary.
- Inspect commit `428f6b8f5391dc851d989364b21a9d727c78cfc0` for the documentation batch.
- Validate locally with `docmgr doctor --ticket OPTKIT-013 --stale-after 30`.
- Consult `tasks.md` for the still-open implementation sequence.

### Technical details

```text
bundle: OPTKIT-013 Catalog and Bindings Guide.pdf
remote: /ai/2026/08/26/OPTKIT-013
commit: 428f6b8f5391dc851d989364b21a9d727c78cfc0
doctor: clean
production code changes: none
```

## Step 3: Implement lossless value specifications and typed domains

This phase replaced the lossy domain descriptor with a strict discriminated contract. Variable construction now derives catalog legality and canonical defaults from the same codec and domain used by runtime assignments, so browser metadata cannot silently disagree with typed mutation behavior.

The phase was bracketed by successful P1 start/done prints. Existing numbergame variables and the OPTKIT-012 proof were migrated directly, with no compatibility aliases for `Name`, `Description`, or the old `DomainDescriptor`.

### Prompt Context

**User prompt (verbatim):**

> implement OPTKIT-012 to OPTKIT-014 (included). commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill). for each ticket, Print out a brutalist work slip with the plan / different phases for the ticket. then before stsarting a phase, plrint a split about the phase, and print one when the phase is done. budget 2M [REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.

**Assistant interpretation:** Implement OPTKIT-013 phase by phase after the accepted OPTKIT-012 gate, printing physical boundaries and retaining exact test/commit evidence.

**Inferred user intent:** Build a generic variable foundation whose semantics are inspectable, executable, and reviewable without metadata drift.

**Commit (code):** `1949a1d7253427f5fd72c1c32b62735dc01f83de` — "OPTKIT-013: add lossless variable value specs"

### What I did

- Added `ValueSpec` kinds for integer, finite float, boolean, pattern string, canonical choice, and artifact reference.
- Added float, bool, string, and artifact domains; changed every domain to derive metadata through its real codec.
- Added `NewVariable`, complete metadata, canonical defaults, cost-hint validation, and descriptor/domain agreement checks.
- Added boundary, malformed-discriminator, choice-value, artifact-schema, and drift tests.

### Why

- Choice labels alone cannot encode machine values.
- JSON clients need a closed legality contract, while typed code must remain authoritative.

### What worked

- Choice descriptors retained canonical machine JSON and deterministic ordering.
- NaN/infinity, wrong artifact schemas, malformed discriminators, and noncanonical defaults were rejected.
- Focused space and numbergame tests plus vet passed.

### What didn't work

- N/A.

### What I learned

- Deriving defaults and legality inside `NewVariable` is a stronger drift guard than asking each registry to reconstruct metadata later.

### What was tricky to build

- Canonical choice values must preserve JSON numeric/string types. Encoding through `Codec[V]` avoids the `map[string]any` integer-to-float trap.

### What warrants a second pair of eyes

- Review whether the finite cost-hint vocabulary (`low`, `medium`, `high`) is sufficient before external consumers depend on it.

### What should be done in the future

- Extend value kinds only through new explicit discriminators and validation tests.

### Code review instructions

- Start at `space/valuespec.go`, then `space/domain.go` and `space/variable.go`.
- Run `GOWORK=off go test ./space ./examples/numbergame -count=1`.

### Technical details

```text
phase slips: 01-p1-start, 02-p1-done
value kinds: 6
compatibility aliases: none
```

## Step 4: Implement ordered catalogs and two-level identity

This phase made semantic variable discovery durable. Catalogs preserve intentional display order while separating machine-semantic identity from exact documentation identity, and all read access returns deep copies.

Strict JSON decoding recomputes identities instead of trusting digest fields supplied by a caller. The P2 start/done slips were printed around implementation and fresh tests.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Build the serializable half of the accepted catalog/binding split.

**Inferred user intent:** Let future CLI, HTTP, and React clients render exact historical variables without making copy edits alter semantic proposal identity.

**Commit (code):** `72c0caeea1236146e2feef7506d4fe5cb8ca5cc5` — "OPTKIT-013: add deterministic semantic catalogs"

### What I did

- Added validated `SectionID`, `Section`, deep-clone helpers, and ordered catalog construction.
- Added `schema:optkit.catalog-semantic/v1` and `schema:optkit.catalog/v1` digests.
- Implemented detached `Sections`/`Lookup`, strict JSON round trip, and identity-tamper rejection.
- Tested prose-only edits, legality edits, intentional order, duplicate IDs, empty sections, and attempted caller mutation.

### Why

- Historical presentation requires exact provenance; candidate semantics must not change for punctuation edits.

### What worked

- Documentation edits changed only full ID; domain edits changed both IDs.
- Caller mutation of nested defaults, choices, labels, and slices did not affect the catalog.

### What didn't work

- N/A.

### What I learned

- Keeping catalog sections unexported is simpler and safer than exposing slices while promising immutability.

### What was tricky to build

- The semantic projection must include intentional section/variable order while excluding every display-only field and both computed IDs.

### What warrants a second pair of eyes

- Review whether variable order should remain semantic; the accepted contract treats it as stable catalog organization.

### What should be done in the future

- Persist the exact full catalog artifact when proposals are sealed.

### Code review instructions

- Review `space/catalog.go`, `space/section.go`, and `space/catalog_test.go`.
- Tamper with an encoded `id` in the JSON test and confirm decode fails.

### Technical details

```text
phase slips: 03-p2-start, 04-p2-done
semantic schema: schema:optkit.catalog-semantic/v1
full schema: schema:optkit.catalog/v1
```

## Step 5: Implement executable typed bindings and registries

This phase completed the other half of the catalog split. `Register` captures each typed variable in a private generic adapter and adds its descriptor and executable binding in one operation.

Draft callers can normalize, read, and apply serialized values without a store. Seal callers replay the same value through the canonical `PatchBuilder`; the test suite proves both paths reach the same child configuration.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Implement serialized-to-typed dispatch without reflection or a parallel type switch.

**Inferred user intent:** Make normal future variables pluggable without editing compiler, manifest, HTTP, or frontend switch statements.

**Commit (code):** `9f2d53486b699224cb2b47af42300235d3bb7444` — "OPTKIT-013: add executable variable registries"

### What I did

- Added `Binding[C]`, private `typedBinding[C,V]`, `RegistryBuilder[C]`, and immutable `Registry[C]`.
- Added canonical read/normalize, pure apply, and durable assign operations.
- Rejected unknown sections/variables, duplicate section IDs, duplicate variable IDs/local keys, malformed JSON, illegal values, and metadata drift.
- Tested descriptor/binding agreement and detached descriptor access.

### Why

- Go closures cannot serialize, but descriptor-only data cannot invoke typed lenses. The adapter must retain both without reflection.

### What worked

- Generic erasure compiled naturally as a package function because Go does not allow generic methods.
- Pure and durable application produced equal child configurations.

### What didn't work

- N/A.

### What I learned

- An immutable map is safe for concurrent readers once the builder copies it during `Build`.

### What was tricky to build

- `Assign` must decode/domain-check immediately even though `PatchBuilder` checks the value again during materialization; early diagnostics and durable defense are both intentional.

### What warrants a second pair of eyes

- Review registry singleton construction and ownership when RAG-TTC begins caching registries.

### What should be done in the future

- OPTKIT-015 should register RAG variables with this API rather than introducing domain-specific switches.

### Code review instructions

- Start with the `Binding` interface and `Register` in `space/binding.go`.
- Run `go test -race ./space -count=1`.

### Technical details

```text
phase slips: 05-p3-start, 06-p3-done
registry mutation after build: impossible through public API
pure persistence writes: zero
```

## Step 6: Advance candidate intent and identity to v2

This phase removed the flat v1 proposer/targets shape and implemented the accepted semantic intent. Candidates now carry structured proposer attribution, expected metric/groups, motivating case/digest evidence, ordered risks, and semantic catalog provenance.

Set-like groups and case IDs normalize by sorting and deduplication. Risks retain authored order, stored prose is not silently trimmed, and creation time remains excluded from identity.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Implement candidate identity exactly as accepted in OPTKIT-012, without dual fields.

**Inferred user intent:** Ensure sealed proposals explain who proposed what, why, expected effect, risks, and evidence as semantic review data.

**Commit (code):** `c482572abf2157f5f69a975358af4de2637e65ff` — "OPTKIT-013: structure candidate semantic intent"

### What I did

- Added proposer kinds, `ExpectedImprovement`, `Motivation`, and `CandidateIntent`.
- Advanced identity to `schema:optkit.candidate-identity/v2`.
- Validated parent/patch/child IDs, actor identity, required prose/metric, diagnostic digest, duplicate/blank risks, and semantic catalog digest.
- Added exhaustive identity inclusion/exclusion and normalization tests.
- Migrated numbergame candidate creation and bundle projection directly to v2.

### Why

- Flat targets could not distinguish expected metric/group semantics or motivating evidence.

### What worked

- Timestamp changes retained ID; every accepted semantic field changed it.
- Group/case order and duplicates normalized; risk order remained semantic.

### What didn't work

- N/A.

### What I learned

- Exact prose hashing avoids a stored/reviewed value that differs from the identity input.

### What was tricky to build

- The API must normalize only fields explicitly declared set-like; blanket slice sorting would erase risk priority.

### What warrants a second pair of eyes

- Review strict proposer-kind enumeration before integrating authenticated service principals.

### What should be done in the future

- Add new proposer kinds only with identity and authorization review.

### Code review instructions

- Review `space/candidate.go` beside the accepted OPTKIT-012 identity matrix.
- Run `go test ./space -run Candidate -count=1`.

### Technical details

```text
phase slips: 07-p4-start, 08-p4-done
candidate schema: schema:optkit.candidate-identity/v2
v1 Targets compatibility field: absent
```

## Step 7: Prove the complete generic path in numbergame

This phase routed a real serialized mutation through the numbergame registry, compared pure preview with durable sealing, completed the existing campaign, and exported a v2 lab bundle reconstructed from sealed facts.

The candidate-proposal artifact now contains the exact catalog as well as candidate and patch. This proves semantic and full catalog provenance before RAG-TTC consumes the APIs.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Exercise catalog, binding, patch, candidate, campaign, and readback as one generic proof.

**Inferred user intent:** Avoid declaring the generic API complete based only on isolated unit tests.

**Commit (code):** `4b90f214d3d4ea577b95f7019fa2d5020d4ff6a5` — "OPTKIT-013: prove catalogs through numbergame"

### What I did

- Added a two-section numbergame registry and serialized `math.multiplier = 3` draft/seal path.
- Compared preview value, serialized child/patch IDs, and direct typed child/patch IDs.
- Sealed catalog provenance in `schema:numbergame.candidate-proposal/v2` and verified it on campaign readback.
- Exported and asserted a fresh eight-episode `numbergame.lab-bundle/v2` under this ticket's `various/` directory.
- Updated the core records/control-model documentation.

### Why

- The numbergame campaign is the smallest complete domain-neutral proof with real persistence and replay.

### What worked

- Export proof reported eight episodes, an eligible decision, and semantic catalog digest `sha256:f9150e2ac397b7e701d091e10b42a5f79ba6d81ceb0ab28c8327148ccaf40359`.
- Persisted candidate/catalog semantic IDs matched.

### What didn't work

- N/A.

### What I learned

- Storing `Catalog` in the proposal artifact exercises its strict custom JSON decoder during historical readback.

### What was tricky to build

- The demo creates the candidate before episode execution, so motivating case IDs must be stable manifest IDs rather than runtime observations.

### What warrants a second pair of eyes

- Review whether future large catalogs should remain embedded in the proposal artifact or be stored as a separate referenced artifact (OPTKIT-017 owns that choice).

### What should be done in the future

- Use separate catalog artifact references for durable RAG proposals as accepted in OPTKIT-012.

### Code review instructions

- Review `examples/numbergame/registry_test.go` and the candidate section of `RunDemo`.
- Run `go run ./cmd/numbergame-demo` with a temporary store and inspect candidate provenance.

### Technical details

```text
phase slips: 09-p5-start, 10-p5-done
proposal schema: schema:numbergame.candidate-proposal/v2
bundle schema: numbergame.lab-bundle/v2
episodes: 8
decision: eligible
```

## Step 8: Run full validation and close the ticket

This phase ran the repository's complete local CI, non-CGO suite, race suite, lint, build, format, demo smoke, documentation, task, and docmgr checks. It also fixed one previously latent exhaustive-switch lint finding in the modified bundle exporter.

The P6 done slip is printed only after these checks and ticket bookkeeping pass. Production implementation is complete; later RAG-specific variables remain intentionally owned by OPTKIT-014/015.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation:** Prove the complete change with repository-native validation and truthful closure evidence.

**Inferred user intent:** Do not hand downstream tickets an API that merely passes focused happy-path tests.

**Commit (validation fix):** `6567319561790e09e6ae2f4ce170bcd65c9089d1` — "OPTKIT-013: make bundle event projection exhaustive"

### What I did

- Ran `make ci-check`: format, vet, CGO tests, non-CGO tests, and full build passed.
- Ran `make race`: every package passed under the race detector.
- Ran `make lint`: golangci-lint 2.4.0 passed with zero issues after triage.
- Checked all six implementation tasks and updated the implementation outcome, core docs, roadmap, changelog, relations, and ticket status.
- Rendered and uploaded the final current guide/diary bundle as `OPTKIT-013 Catalog Bindings Implementation Complete.pdf`.
- Kept unrelated untracked `optkit/numbergame-demo` untouched.

### Why

- Registry maps are shared read-only state and deserve race validation.
- Public data/schema changes require full repository compile and both CGO modes.

### What worked

- Final `make ci-check`, `make race`, and `make lint` all passed.
- The complete suite includes the OPTKIT-012 contract proof under `ttmp`, so the accepted type-erasure sketch still compiles against the production API.

### What didn't work

- The first lint run failed at `cmd/numbergame-demo/main.go:154:3` with `missing cases in switch of type campaign.EventKind` from the exhaustive linter. Adding `default` did not satisfy this configuration. I replaced it with one explicit grouped case containing all journal-only event kinds; the next lint run reported `0 issues`.

### What I learned

- This repository's exhaustive configuration requires explicit enum coverage; a default branch is not treated as exhaustive.

### What was tricky to build

- The exporter intentionally projects only four payload kinds but retains every event in its journal. Explicitly grouping all other enum values documents that asymmetry and ensures new event kinds trigger review.

### What warrants a second pair of eyes

- Review public API naming and candidate v2 migration before downstream packages publish it externally.
- Review the committed proof bundle for accidental environment-sensitive fields; it is evidence, not a golden fixture.

### What should be done in the future

- OPTKIT-014 should consume `NewVariable`, `Registry`, and lifted lenses without altering their contracts.

### Code review instructions

- Review commits `1949a1d`, `72c0cae`, `9f2d534`, `c482572`, `4b90f21`, and `6567319` in order.
- Run `make ci-check && make race && make lint` from `optkit/`.
- Run `docmgr doctor --ticket OPTKIT-013 --stale-after 30` from the workspace root.
- Inspect all `various/work-slips/*.log` for `printed: true`.

### Technical details

```text
full CGO tests: pass
full non-CGO tests: pass
race: pass
vet/build/fmt: pass
golangci-lint 2.4.0: 0 issues
implementation tasks: 6/6 complete
phase slips: plan + 6 start + 6 done
```
