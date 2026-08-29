---
Title: Diary
Ticket: OPTKIT-021
Status: active
Topics:
    - architecture
    - design
    - ui
    - rag-ttc
    - optkit
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Working diary for the PBUI adoption session — study of pbui/plot/datalab, the workflows/apps/objects brainstorm, the supersession of OPTKIT-019, and the creation of OPTKIT-021 through OPTKIT-024 with their intern guides.
WhatFor: Record how the PBUI track was designed, what evidence it rests on, and what a reviewer should check.
WhenToUse: Read before reviewing the OPTKIT-021 ADRs or resuming any OPTKIT-021..024 work.
---

# Diary

## Goal

Capture the design session that created the PBUI adoption track: studying the
pbui/plot/datalab stack, brainstorming the workbench as presentation-based UI,
closing OPTKIT-019, and writing the OPTKIT-021–024 intern guides.

## Step 1: Study pbui, plot, datalab, the ticket portfolio, and the architect brief

The session started from a study request: understand the three sibling
repositories (`pbui`, `plot`, `datalab`), the OPTKIT-011 ticket family with
emphasis on OPTKIT-015 onward, and the architect brief, then brainstorm what a
PBUI-based UI for optkit campaigns and the rag-ttc pipeline would look like.
Four parallel read-only research agents were dispatched (one per repository
plus one auditing OPTKIT-012–020), and the architect brief and the OPTKIT-011
roadmap were read directly.

The study established the facts every later decision rests on: pbui is a
two-layer system (presentation kernel with descriptors/verbs/accept;
workbench tiling shell with protobuf-persisted documents and parity-tested
Go/TS mutation appliers); plot exposes every mark as a typed `PlotHit` that
datalab already wraps in `<Presentation>`; datalab is the reference product
(explicit app registry mirrored by a Go `ApplicationCatalog`, strict document
validators, three embedded doc topics). The ticket audit found a sharp line:
OPTKIT-012–014 built, OPTKIT-015–020 designed-only — with OPTKIT-019's guide
being the most UI-committed document and the direct collision with PBUI.

### Prompt Context

**User prompt (verbatim):** "Study ./pbui and ./plot and ./datalab and then
especially @ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/
and all the related tickets especially OPTKIT-015 onwards about the UI, and
@~/Downloads/Architect Brief_ Modular Optimization Workbench + Candidate
Authoring.md \n\nThen, we want to study / examine / brainstorm what a PBUI
based UI would look like for working with optkit campaigns and our rag-ttc /
coinvault optimizing RAG pipeline particularly.\n\nAfter studying, make a list
of the workflows you would support and what kind of PBUI tiles (called apps I
think) and the kind of workspaces and objects and actions you would suggest.\n\nAnswer
here in the chat first, so we can iterate a bit. Study things consciensously
first." (Mid-study addition: "we had a colleague implement OPTKIT-012 -
OPTKIT-014 already.")

**Assistant interpretation:** Perform a serious multi-repository study, then
deliver an iterable chat brainstorm listing workflows, pbui apps, workspaces,
objects, and actions for an optimization workbench UI.

**Inferred user intent:** Evaluate PBUI as the real frontend platform for the
optimization workbench program before committing, with enough concrete detail
to make an adoption decision.

### What I did
- Dispatched four Explore agents (pbui architecture, plot, datalab, ticket
  audit) and read the architect brief and roadmap doc 04 directly.
- Synthesized the brainstorm in chat: 8 workflows, ~20 presentation types with
  tones and conversions, 17 product apps in four groups, 5 workspaces, the
  draft-as-DocumentPayload idea, accept-mode signature moments, the agent
  seat, and buildability gating.

### Why
- The presentation model only pays off if the domain objects, verbs, and
  document formats are grounded in what the stack actually provides; guessing
  would have produced a design pbui cannot host.

### What worked
- The parallel-agent study: each report came back with exact type names, file
  paths, and the load-bearing rules (explicit app registry, verbs-as-data,
  derived-vs-persisted state split, plot's `renderInteractive` seam).
- The audit caught the in-flight OPTKIT-014 closure and the still-hardcoded
  `RRFConstant: 60`, which anchored the "015 is the true next gate" claim.

### What didn't work
- N/A — study and brainstorm proceeded without failures.

### What I learned
- Datalab's registration style (module-level side-effect `registerApp`) is
  deprecated; the pbui-chat demo's explicit `createWorkbench({apps})` list is
  the current contract. The guides must copy the demo, not datalab.
- The draft/seal distinction maps exactly onto workbench-document vs journal,
  which turns the program's most load-bearing backend rule into a spatial UI
  property.

### What was tricky to build
- Nothing built yet; the tricky part was scoping — separating what PBUI
  replaces (OPTKIT-019's bespoke shell) from what it must not touch (the
  backend chain and its gates).

### What warrants a second pair of eyes
- The claim that evidence tiles need only `specialistapi` plus one additive
  projection; if any tile turns out to need command-API data, the
  parallel-work argument weakens.

### What should be done in the future
- N/A (superseded by Steps 2–4, which executed the decisions).

### Code review instructions
- No code. Review the chat brainstorm against the pbui playbooks
  (`pbui/docs/playbooks/`) and the OPTKIT-019 guide for salvage completeness.

### Technical details
- Study agents' key sources: `pbui/src/presentation/types.ts`,
  `pbui/packages/pbui-workbench/src/apps.ts`,
  `pbui/proto/hyperslop/pbui/workbench/v1/workbench.proto`,
  `plot/src/scene.ts` (PlotHit), `datalab/pkg/workbenchapp/catalog.go`,
  `pbui/packages/datalab-ui/src/components/organisms/ChartPanel/ChartPanel.tsx`.

## Step 2: Decide the ticket surgery — keep 015–018, supersede 019, amend 020

Asked which tickets could be adapted versus closed, the analysis concluded the
backend-first sequencing had contained the damage: OPTKIT-015–018 are
UI-agnostic and survive untouched (018 with two additive notes), OPTKIT-020
needs only its React-editor section repointed at the plugin registry, and
OPTKIT-019 is the sole casualty because its architectural spine is exactly
what pbui-workbench already provides. Proposed successors: OPTKIT-021 (ADRs),
OPTKIT-022 (product + evidence tiles, parallel with the backend chain),
OPTKIT-023 (propose workspace + RRF slice, inheriting 019's exit criterion),
OPTKIT-024 (agent seat, last).

### Prompt Context

**User prompt (verbatim):** "Which tickets can be adapted or should we close
and make new ones?"

**Assistant interpretation:** Produce a keep/adapt/close verdict per ticket
with rationale and a concrete successor plan.

**Inferred user intent:** Minimize churn to the colleague's in-flight backend
work while restructuring the frontend tickets around PBUI.

### What I did
- Delivered the per-ticket verdict table and the four-ticket successor plan in
  chat, including the recommendation to close 019 rather than rewrite it in
  place.

### Why
- Rewriting 019's validated guide in place would blur the program's
  "guide validated → implement" discipline and erase the design record of the
  pre-PBUI architecture.

### What worked
- The salvage list (generic editors, plugin resolution order, seal state
  machine, server-order invalidation, RRF inspector, parity pinning) gave the
  successors concrete inherited content instead of a blank page.

### What didn't work
- N/A.

### What I learned
- The dependency picture improves under the split: 022 runs parallel to the
  backend chain, which the monolithic 019 could not do.

### What was tricky to build
- N/A (analysis step).

### What warrants a second pair of eyes
- Whether OPTKIT-024 should exist now as a designed ticket or wait; it was
  created (per the user's go-ahead) with a deliberately proportionate guide.

### What should be done in the future
- N/A (executed in Steps 3–4).

### Code review instructions
- Compare the OPTKIT-019 index "Superseded" section against the successors'
  guides to confirm every salvaged item actually landed somewhere.

### Technical details
- N/A.

## Step 3: Close OPTKIT-019 and create OPTKIT-021–024

Executed the surgery. OPTKIT-019: `Status: superseded`, a "Superseded" section
in the index enumerating successors and the salvage list, tasks.md rewritten to
mark the six implementation tasks as migrated (to 022/023), changelog entry —
commit `7b81bb2`. Created the four tickets with docmgr, seeded their task
lists (implementation tasks unchecked, guide-writing tasks checked on
delivery), and wrote the OPTKIT-021 intern guide: adoption rationale, ADRs
G–L, vocabulary tables, verb families, package diagram, acceptance gates,
pitfalls — commit `502310a`.

### Prompt Context

**User prompt (verbatim):** "ok, close 019, create the new tickets, write a
detailed design doc per ticket. Create  a detailed analysis / design /
implementation guide that is for a new intern, explaining all the parts of the
system needed to understand what it is, with prose paragraphs and bullet point
sand pseudocode and diagrams and api references and file references. It should
be very clear and technical. Store in the ticket and the nupload to
remarkable.\n\n\ncommit at appropriate intervals and keep a detailed diary as
you work (using the diary format from the skill)\n\nDon't implement though,
our colleague is still working on the backend."

**Assistant interpretation:** Execute the agreed plan: supersede 019, create
021–024 each with an intern-grade design guide (prose, bullets, pseudocode,
diagrams, API and file references), maintain docmgr bookkeeping and this
diary, commit incrementally, bundle to reMarkable, and write no
implementation code.

**Inferred user intent:** Have the PBUI track fully specified and reviewable
(on reMarkable) so implementation can start cleanly once the backend colleague
lands OPTKIT-016–018, without the two workstreams colliding.

**Commit (code):** `7b81bb2` — "OPTKIT-019: supersede in favor of PBUI
adoption track (OPTKIT-021..024)"; `502310a` — "OPTKIT-021: create PBUI
adoption ticket with ADR/vocabulary intern guide; scaffold OPTKIT-022..024"

### What I did
- Edited OPTKIT-019 index/tasks, added the changelog entry, committed.
- `docmgr ticket create-ticket` × 4; `docmgr doc add` for guides and this
  diary; renamed placeholders to descriptive filenames.
- Wrote the OPTKIT-021 guide (ADRs G–L); filled index summary/overview;
  seeded 7 tasks; relates + changelog.

### Why
- ADRs before tiles: the vocabulary and document formats are the contracts
  both 022 and 023 build against, and freezing them first is the whole point
  of the 021 ticket.

### What worked
- `docmgr doc relate` normalized the guide frontmatter (added `ws://` related
  files) without disturbing content.

### What didn't work
- First closure commit failed with exit 128: the shell `cd`'d into the ticket
  directory earlier, so the pathspec resolved wrongly
  (`fatal: pathspec 'ttmp/...' did not match any files`). Re-ran from the repo
  root; lesson restated: always run git from the repo root in this worktree.

### What I learned
- The ticket scaffolds' index frontmatter has empty Summary/WhatFor/WhenToUse
  that docmgr does not fill; each index needed a manual frontmatter pass.

### What was tricky to build
- ADR I's derived-state rule needed a mechanical enforcement story (the
  validator rejects `draft_digest`/`plan`/`diff`/`diagnostics` keys by name)
  to keep "never store compiled output" from being aspirational prose.

### What warrants a second pair of eyes
- ADR H's boundary rule that the document validator does **not** check
  referential validity (opaque ids, tiles render honest errors) — deliberate,
  but a reviewer may expect validation there.
- The `/api/rag/workbench-docs/v1` vs `/api/rag/workbench/v1` path split —
  naming is close; confirm it is distinct enough in practice.

### What should be done in the future
- ADR review by the architect/user; the six review tasks in tasks.md.

### Code review instructions
- Start with `design-doc/01-intern-guide-pbui-adoption-adrs-and-workbench-vocabulary.md`
  sections 4–9 (the ADRs); check the vocabulary tables against
  `pbui/src/presentation/types.ts` constraints and the verb list against
  `pbui/packages/pbui-workbench/src/verbs.ts` conventions.

### Technical details
- Ticket paths: `ttmp/2026/08/26/OPTKIT-021--pbui-adoption-adrs-and-workbench-vocabulary/`
  and siblings `OPTKIT-022--rag-ttc-workbench-product-package-and-evidence-tiles/`,
  `OPTKIT-023--propose-workspace-draft-document-and-rrf-vertical-slice/`,
  `OPTKIT-024--agent-seat-over-the-workbench-vocabulary/`.

## Step 4: Write the OPTKIT-022, 023, and 024 intern guides; update the roadmap

Wrote the three remaining guides. OPTKIT-022 (commit `4e478e3`): the pbui
product model for interns, the package scaffold, a complete descriptor and
sink example, nine evidence tiles with their specialist-SPA organism sources,
the plot adapter pattern, the Go `workbenchhost` design, the worst-first
failures projection contract, testing strategy, and a nine-step
implementation order. OPTKIT-023 (commit `21521df`): draft-document lifecycle
with the derived seal state machine (mermaid flow + compile-loop sequence
diagram), the Propose workspace with four linked tiles, accept flows, generic
`ValueSpec` editors plus the plugin registry, the fusion plugin with the
recorded-precision parity rule, verb-family completion, and a seven-step
order whose steps 1–6 run on compiler-CLI fixtures. OPTKIT-024 (commit
`a9d0aa6`): vocabulary export, mentions, router families with actor
attribution, draft co-editing over the shared document, Proposer/approval
provenance, authorization-as-safety. Then updated roadmap doc 04 (commit
`e4b1530`): 019 section marked superseded, four new ticket sections, new
dependency graph with the parallel PBUI branch, Gate 6 repointed at 023,
status table extended.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

**Assistant interpretation / intent:** (see Step 3)

**Commit (code):** `4e478e3`, `21521df`, `a9d0aa6`, `e4b1530` (docs-only).

### What I did
- Wrote three guides, filled three ticket indexes, seeded task lists (022: 9,
  023: 8, 024: 7), changelog entries, related files.
- Updated roadmap doc 04's ticket portfolio, dependency mermaid, Gate 6,
  OPTKIT-020 dependency line, and status table.

### Why
- The roadmap is declared "the authoritative source for child-ticket goals,
  dependencies, exclusions, gates, and status"; leaving it pointing at a
  superseded ticket would have made the supersession invisible at the program
  level.

### What worked
- Inheriting 019's exit criterion verbatim into 023 kept the program gates
  stable: Gate 6 changed its ticket number, not its meaning.
- OPTKIT-023's fixture strategy (record `proposal compile` CLI output as MSW
  fixtures) decouples frontend steps 1–6 from OPTKIT-018's HTTP layer.

### What didn't work
- The roadmap file already carried the colleague's uncommitted OPTKIT-015
  completion edits, so my edits interleaved with theirs in one file. Resolved
  by committing only the shared roadmap file (including their one status row,
  which reflects already-committed rag-ttc work) and leaving all their
  OPTKIT-015 ticket files uncommitted; noted explicitly in the commit message.

### What I learned
- OPTKIT-015 is now complete per the colleague's in-flight edits ("runtime RRF
  injection, executable retrieval/fusion registry, fixture v3 parity") — the
  backend chain is at OPTKIT-016, further along than the morning audit.

### What was tricky to build
- The 022/023 boundary: the sink skeleton, trace store, and even the
  authoring actions on descriptors (greyed with `disabledBecause`) belong to
  022 so that 023 only fills in families and tiles. Getting that split wrong
  would have made 023 an architecture ticket again. Symptom to watch in
  review: any 023 section that modifies a 022 file's structure rather than
  adding to a registry.

### What warrants a second pair of eyes
- The worst-first failures projection contract (022 §9) — it is a backend
  shape settled unilaterally in a frontend guide; the backend colleague should
  confirm it fits the OPTKIT-018 projector conventions before implementing.
- The 023 parity rule's fixture provenance: it asks OPTKIT-015 (now complete)
  to export recorded RRF contributions; confirm the existing OPTKIT-015 test
  evidence includes or can regenerate such a fixture.
- The decision to commit the shared roadmap file while the colleague's
  OPTKIT-015 doc edits were uncommitted.

### What should be done in the future
- reMarkable upload of the four guides plus the 019 closure (Step 5).
- After ADR review: begin OPTKIT-022 step 1 (scaffold) — implementation is
  explicitly out of scope for this session.

### Code review instructions
- Read guides in order 021 → 022 → 023 → 024; each declares its assumed
  contracts, so forward references are intentional.
- Validate: `docmgr doctor --ticket OPTKIT-021` (and 022/023/024);
  `git log --oneline -8` for the commit trail.

### Technical details
- Guide files:
  `OPTKIT-021.../design-doc/01-intern-guide-pbui-adoption-adrs-and-workbench-vocabulary.md`,
  `OPTKIT-022.../design-doc/01-intern-guide-ragttc-workbench-product-package-and-evidence-tiles.md`,
  `OPTKIT-023.../design-doc/01-intern-guide-propose-workspace-draft-document-and-rrf-slice.md`,
  `OPTKIT-024.../design-doc/01-intern-guide-agent-seat-over-the-workbench-vocabulary.md`.

## Step 5: Publish the bundle to reMarkable

Bundled the four intern guides (021 → 024, reading order) into one PDF and
uploaded to the reMarkable at `/ai/2026/08/26/OPTKIT-021-PBUI-workbench`.
Recorded upload results in each ticket's changelog. Final commit of diary and
bookkeeping.

### Prompt Context

**User prompt (verbatim):** (see Step 3)

### What I did
- `remarquee upload bundle` of the four guide files with `--toc-depth 2
  --non-interactive`; changelog entries; final docs commit.

### Why
- The user reviews long design material on the reMarkable; one bundle in
  reading order beats four separate PDFs for a track review.

### What worked
- See changelog entries per ticket for the upload record.

### What didn't work
- (recorded below if the upload required retries)

### What I learned / tricky / second pair of eyes / future
- N/A beyond the above; the session ends with the PBUI track fully specified,
  zero implementation, and the backend colleague's working tree untouched.

### Code review instructions
- Verify the bundle contents match the four committed guide files.

### Technical details
- Remote dir: `/ai/2026/08/26/OPTKIT-021-PBUI-workbench`.

## Step 6: Retarget the Guide to pbui 0.8.0 (Greenfield, Kernel-Only)

The user ruled that this program has no backwards-compatibility obligations
and that pbui's legacy surfaces should be deleted rather than deprecated.
pbui 0.8.0 shipped that deletion (pbui repo, PBUI-ACTIONS-3, commit
6efeaeb): descriptor `actions()`, `conversions`, and the automatic legacy
engine are gone; `actions`+`snapshotFor` are required; acceptance always has
graph-subtype satisfaction; and a `primary` invocation performs the unique
available primary action on left click. This step rewrote the guide so no
section names a deleted mechanism.

### Prompt Context

**User prompt (verbatim):** "perfect, update the design docs to have the simplification that not having to worry about lgacy affords us."

**Assistant interpretation:** Rewrite the 021–024 guides against the
post-cleanup pbui API (the simplification plan presented and approved in
chat).

**Inferred user intent:** The ragttc workbench guides should describe
exactly what will be built — elegant, kernel-native, no legacy vocabulary —
so implementation never touches a mechanism scheduled to die.

### What I did
- §2.2: descriptors are representation-only; actions are kernel
  declarations (four-state availability, bind-only-available, ambiguity as
  data, fresh revalidation); `createPbui` signature updated; primary-click
  paragraph added; key-files list now names the kernel and pins pbui 0.8.0.
- ADR J §7.4: conversions table became eight named translator EDGES
  (`ragttc.verdict-to-case`, …) with chooser-on-tie semantics and a note on
  subtype satisfaction.
- ADR J §7.5: split environment (representation: name lookups) from
  SelectionSnapshot facts (rules), with `seal` as a CAPABILITY — the SealBar
  and menu render one resolved action and cannot disagree.
- ADR K: click-to-open gestures are `metadata.primary` rules; danger flags
  live in rule metadata; unavailable seal rows carry no bound verb.
- ADR L: the agent vocabulary export is reserved as GENERATED from the
  registry (PBUI-ACTIONS-3 Phase B), never hand-maintained.
- Pitfalls: "descriptors must not fetch" became the snapshot cost-boundary
  rule; the hide-unavailable pitfall now teaches
  unavailable/inapplicable/hidden selection.

### Why
- The guides were written against pbui 0.6.x; OPTKIT-022 starts from them.

### What worked
- Every change was a substitution, not a redesign — the ADR decisions (G–L)
  survive intact; only their pbui-facing vocabulary moved.

### What didn't work
- N/A

### What I learned
- The capability mechanism subsumes what §7.5 had modeled as an environment
  boolean; the ADR text is simpler after the split than before it.

### What was tricky to build
- Keeping §7.4's translator ids stable-worthy: they are wire names the ADR L
  reservation now covers, so they were chosen with the `ragttc.` prefix and
  recorded as part of the frozen v1 contract.

### What warrants a second pair of eyes
- The frozen vocabulary now includes translator ids and the two capability
  names (`seal`, and OPTKIT-024's approval grant); confirm during ADR review
  that these belong to the v1 freeze.

### What should be done in the future
- OPTKIT-022/023/024 guide updates follow in this same pass (their own
  diaries record them).

### Code review instructions
- Diff this guide against git history; check no section references
  descriptor actions, conversions, or `canSeal`-in-environment.

### Technical details
- pbui 0.8.0 = pbui repo `task/use-optkit` commit 6efeaeb; test sweep 1221.
