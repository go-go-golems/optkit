#!/usr/bin/env python3
"""Write the initial strict diary step for OPTKIT-012 through OPTKIT-020.

The generated entries share the exact umbrella prompts but retain ticket-specific
research, architecture, review, and continuation details. Run from workspace root.
"""
from pathlib import Path

ROOT = Path("optkit/ttmp/2026/08/26")
PROMPT = """Ok, create all these tickets, and for each ticket, Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.

commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)

Youc an run experiments and such, store all scripts you write in the scripts/ folder of the relevant ticket."""
FOLLOWUP_1 = """add the description of each ticket and make an upfront design doc with the ticket goal and content so that we can keep track as we go on, to not lose the overall picture."""
FOLLOWUP_2 = """since each document will require a significant amount of planning / investigation, so we want to make sure we don't forget what the do cis. Once that overview doc is written, update the goal to reference it."""

TICKETS = {
    "OPTKIT-012": {
        "slug": "architecture-closure-and-optimization-workbench-contracts",
        "goal": "Record the architecture closure that every later workbench ticket must implement.",
        "interpretation": "Turn the architect brief into concrete, reviewable whole-configuration, catalog/binding, value-schema, identity, provenance, service, and compatibility contracts.",
        "intent": "Prevent nine substantial implementation efforts from drifting into incompatible local designs.",
        "actions": [
            "Mapped the existing typed variable, patch, graph, manifest, campaign, and projector boundaries.",
            "Specified concrete Go shapes for PipelineConfig, Catalog, Binding, CompileProposal, SealProposal, and durable candidate records.",
            "Recorded six proposed decisions and an explicit compatibility-policy gate.",
        ],
        "files": ["optkit/space/variable.go", "optkit/space/patch.go", "rag-ttc/pkg/ttc/optimization/graph.go", "rag-ttc/pkg/ttc/experimentworkbench/service.go", "rag-ttc/pkg/ttc/optkitcampaign/campaign.go"],
        "tricky": "The proposed package ownership must avoid an optimization/optkitcampaign import cycle while still making semantic configuration—not a campaign adapter—the source of graph truth.",
        "risk": "Review candidate/catalog identity inputs and decide v1 manifest/store compatibility before child implementations begin.",
        "future": "Accept or reject each proposed ADR and compile the interface sketches before opening OPTKIT-013 implementation changes.",
    },
    "OPTKIT-013": {
        "slug": "optkit-semantic-catalog-and-executable-variable-bindings",
        "goal": "Design the domain-neutral semantic catalog, lossless domains, executable bindings, and candidate intent extension.",
        "interpretation": "Explain how serialized values invoke existing typed Optkit variables without adding a second mutation algebra, then prove the design with numbergame.",
        "intent": "Create the generic framework all RAG coordinates and authoring transports can share.",
        "actions": [
            "Inspected Variable, Domain, Lens, Codec, PatchBuilder, Snapshot, Candidate, and numbergame candidate flow.",
            "Designed one-construction-site descriptor/binding registration and semantic/full catalog identities.",
            "Specified domain, binding, candidate identity, and numbergame conformance tests.",
        ],
        "files": ["optkit/space/variable.go", "optkit/space/domain.go", "optkit/space/lens.go", "optkit/space/patch.go", "optkit/examples/numbergame/demo.go"],
        "tricky": "Type erasure must retain the real Codec[V], Domain[V], and Lens[C,V]; decoding through any/map values would corrupt integer and choice semantics.",
        "risk": "Ensure descriptors cannot be registered independently from bindings and that documentation identity is separated from semantic identity.",
        "future": "Implement value specs, catalog construction, typed adapters, candidate intent, and the numbergame proof in the documented order.",
    },
    "OPTKIT-014": {
        "slug": "whole-pipeline-rag-configuration-and-graph-derivation",
        "goal": "Design PipelineConfig, lifted layer lenses, and deterministic graph derivation.",
        "interpretation": "Replace retrieval-only snapshots plus manually authored graph identities with one typed aggregate semantic configuration.",
        "intent": "Remove config/graph drift before any new RAG coordinate is exposed.",
        "actions": [
            "Mapped the twelve-layer graph and current retrieval-only Factory/Executor snapshot contract.",
            "Designed explicit frozen values, typed retrieval/fusion fields, and one dependency topology table.",
            "Specified lens-law, graph identity, runtime parity, package-cycle, and migration tests.",
        ],
        "files": ["rag-ttc/pkg/ttc/optimization/contracts.go", "rag-ttc/pkg/ttc/optimization/graph.go", "rag-ttc/pkg/ttc/optkitcampaign/system.go", "rag-ttc/pkg/ttc/experimentworkbench/manifest.go", "optkit/space/lens.go"],
        "tricky": "The aggregate must model frozen layers explicitly without becoming an untyped map, and lifted lens setters must preserve every unrelated layer value.",
        "risk": "Review package direction, direct dependency topology, and all identity/schema changes before migration.",
        "future": "Land config/schema types first, then graph derivation and lenses, then migrate execution/manifests/campaign projections with parity evidence.",
    },
    "OPTKIT-015": {
        "slug": "real-fusion-configuration-and-first-rag-optimization-catalog",
        "goal": "Design the first runtime-honest RAG variables: final result limit and float64 RRF k.",
        "interpretation": "Plumb typed FusionConfig into actual WeightedRRF execution and register reviewed retrieval/fusion catalog entries.",
        "intent": "Prove a registered variable changes real arithmetic and matching graph invalidation.",
        "actions": [
            "Traced RRFConstant through SearchConfig, SearchRoute, route identity, validation, and WeightedRRF.",
            "Traced RetrievalConfig.Limit to SearchInput and confirmed it caps final returned evidence rather than channel top-K.",
            "Designed runtime, contribution, catalog/binding, graph-plan, and deterministic fixture tests.",
        ],
        "files": ["rag-ttc/pkg/ttc/search/search.go", "rag-ttc/pkg/ttc/search/service.go", "rag-ttc/pkg/ttc/search/semantic_fixture.go", "rag-ttc/pkg/ttc/optkitcampaign/fixture.go", "rag-ttc/pkg/ttc/optimization/invalidation.go"],
        "tricky": "Several production and fixture routes carry RRF values; the implementation must change the intended Optkit path without overriding unrelated route configuration.",
        "risk": "Review float canonicalization, recorded contribution precision, and the final-limit/per-channel-top-K naming boundary.",
        "future": "Implement semantic naming/config first, then runtime injection, graph parity, registry declarations, and an auditable experiment.",
    },
    "OPTKIT-016": {
        "slug": "proposal-compiler-and-glazed-cli",
        "goal": "Design side-effect-free proposal compilation and backend-first Glazed authoring commands.",
        "interpretation": "Build one deterministic application service for manifests, CLI, and future browser drafts without writing artifacts or journal events.",
        "intent": "Prove the complete mutation/diff/plan contract before durable sealing or UI work.",
        "actions": [
            "Mapped current experimentworkbench service and Glazed config/campaign command boundaries.",
            "Specified request, normalized mutation, draft, diagnostic, preview-capability, and digest contracts.",
            "Designed deterministic errors/order, catalog/proposal commands, and a repeated-compilation no-write proof.",
        ],
        "files": ["rag-ttc/pkg/ttc/experimentworkbench/service.go", "rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/config.go", "rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/command.go", "optkit/space/patch.go", "rag-ttc/pkg/ttc/optimization/invalidation.go"],
        "tricky": "Operator-correctable errors should be complete deterministic diagnostics, while service-state failures remain errors; partial drafts must never be sealable.",
        "risk": "Verify no artifact/store dependency reaches ProposalCompiler and that CLI --set parsing does not guess JSON types.",
        "future": "Implement DTOs/compiler tests first, then catalog/proposal Glazed adapters and store-count smoke evidence.",
    },
    "OPTKIT-017": {
        "slug": "proposal-sealing-candidate-manifests-and-campaign-persistence",
        "goal": "Design canonical sealing, strict candidate manifests, idempotent events, and self-contained campaign persistence.",
        "interpretation": "Replay verified drafts through typed bindings and PatchBuilder so concise candidate authoring becomes immutable campaign fact.",
        "intent": "Make campaigns executable and explainable after their source manifests disappear.",
        "actions": [
            "Mapped PatchBuilder, numbergame CandidateProposed precedent, strict manifest loading, and campaign initialization ordering.",
            "Specified seal/materialize/record APIs, candidate v2 YAML, durable records, catalog provenance, and idempotency.",
            "Designed stale-draft, retry, restart, manifest-removal, and artifact verification tests.",
        ],
        "files": ["optkit/space/patch.go", "optkit/space/candidate.go", "optkit/examples/numbergame/demo.go", "rag-ttc/pkg/ttc/experimentworkbench/manifest.go", "rag-ttc/pkg/ttc/optkitcampaign/campaign.go"],
        "tricky": "Manifest-driven materialization may occur before a campaign journal exists, while event recording needs campaign state; the internal boundary must support both without special mutation logic.",
        "risk": "Review canonical relationships between creation-spec candidate records and CandidateProposed payloads, plus retry conflicts and stale parent/catalog handling.",
        "future": "Implement strict v2 resolution, materialization, durable schemas, event idempotency, and the manifest-deletion proof.",
    },
    "OPTKIT-018": {
        "slug": "candidate-projections-and-workbench-command-api",
        "goal": "Design historical candidate projections and a separate authorized workbench command API.",
        "interpretation": "Expose sealed intent/catalog facts for reads and adapt compile/preview/seal applications to HTTP without mixing them into specialist projections.",
        "intent": "Give React authoritative read and command contracts while preserving historical determinism.",
        "actions": [
            "Mapped current GET-only ServeMux, fact loader, comparison projection, and application services.",
            "Specified candidate/catalog responses, command routes, authorization actions, typed errors, idempotency, and preview envelopes.",
            "Designed historical fallback, sealed catalog labeling, strict HTTP, security, concurrency, and live smoke tests.",
        ],
        "files": ["rag-ttc/pkg/ttc/specialistapi/projector.go", "rag-ttc/pkg/ttc/specialistapi/types.go", "rag-ttc/pkg/ttc/specialistapi/http.go", "rag-ttc/pkg/ttc/experimentworkbench/service.go", "rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/serve.go"],
        "tricky": "One server may host both route groups, but package/application boundaries must keep historical reads incapable of compiling or sealing proposals.",
        "risk": "Review actor spoofing, sensitive diagnostics/logging, typed error mapping, restart-safe idempotency, and sealed-versus-current catalog behavior.",
        "future": "Implement candidate fact projections first, then catalog reads, command applications, strict HTTP adapters, and frontend fixtures.",
    },
    "OPTKIT-019": {
        "slug": "react-workbench-framework-and-rrf-vertical-slice",
        "goal": "Design the reusable React workbench and complete RRF scalar vertical slice.",
        "interpretation": "Combine generic catalog editors, specialized plugins, shared workflow shell, backend compilation, and historical comparison without a server-driven UI DSL.",
        "intent": "Deliver the first product-complete authoring experience after backend semantics are proven.",
        "actions": [
            "Mapped RTK Query APIs, layer/output registries, hardcoded retrieval copy, ComparisonScreen, and numbergame LabScreen/parity tests.",
            "Specified generic discriminated editors, WorkbenchRegistry, shell/context, compile race handling, RRF inspector, intent, and sealing flow.",
            "Designed MSW, parity, accessibility, deep-link, typecheck/build, and live Playwright validation.",
        ],
        "files": ["rag-ttc/apps/specialist/web/src/layerwidgets/index.tsx", "rag-ttc/apps/specialist/web/src/screens/LabScreen.tsx", "rag-ttc/apps/specialist/web/src/screens/ComparisonScreen.tsx", "rag-ttc/apps/specialist/web/src/api/specialistApi.ts", "rag-ttc/apps/specialist/web/src/test/lab-sim.test.ts"],
        "tricky": "Debounced compile responses can arrive out of order; the UI must never seal a draft whose digest no longer matches current local edits.",
        "risk": "Review backend semantic duplication, RRF parity/precision, broad context rerenders, missing-value handling, and keyboard behavior.",
        "future": "Add API fixtures/types and generic editors first, then registry/shell, compile integration, RRF plugin, sealing, and live inspection.",
    },
    "OPTKIT-020": {
        "slug": "asset-variable-proof-for-representation-prompts",
        "goal": "Design the artifact-valued summary-prompt mutation and bounded one-chunk preview.",
        "interpretation": "Generalize the scalar workbench to sensitive, content-addressed, expensive configuration without prompt text in patch metadata or client-side fake generation.",
        "intent": "Prove one core workflow supports both scalar and asset coordinates honestly.",
        "actions": [
            "Mapped Optkit assignment refs, representation summaries, artifact rendering precedents, and the required producer discovery boundary.",
            "Specified RepresentationsConfig, artifact-ref variable, provisional draft refs, seal-time materialization, preview provenance, and sensitivity policy.",
            "Designed producer inventory, one-chunk experiment, cancellation/security, scalar-versus-asset, frontend, and campaign tests.",
        ],
        "files": ["optkit/artifact", "optkit/space/patch.go", "rag-ttc/pkg/ttc/search/types.go", "rag-ttc/apps/specialist/web/src/components/Artifact.tsx", "rag-ttc/apps/specialist/web/src/layerwidgets/index.tsx"],
        "tricky": "The semantic fixture may expose only precomputed representation text; implementation must find or build an honest producer instead of editing recorded output and calling it generation.",
        "risk": "Review prompt leakage, provisional/stored digest parity, provider/model identity, preview evidence labeling, and resource bounds.",
        "future": "Inventory/prove the producer first; only then implement config/variable, provisional assets, bounded preview, sealing, and editor conformance.",
    },
}

for ticket, data in TICKETS.items():
    directory = ROOT / f"{ticket}--{data['slug']}"
    diary = directory / "reference/01-implementation-diary.md"
    text = diary.read_text()
    first = text.find("---")
    second = text.find("---", first + 3)
    frontmatter = text[: second + 3]
    actions = "\n".join(f"- {item}" for item in data["actions"])
    files = "\n".join(f"- `{path}`" for path in data["files"])
    body = f'''

# Diary

## Goal

{data["goal"]} This diary records research, design decisions, validation, commits, delivery evidence, failures, and continuation instructions.

## Step 1: Investigate the current system and write the intern design guide

This step established the ticket before implementation starts. I read the parent OPTKIT-011 roadmap and architect brief, traced the current source paths that constrain this ticket, and wrote a self-contained intern guide that explains both current behavior and the proposed implementation sequence.

The design preserves the backend-first dependency order. It records concrete APIs, pseudocode, diagrams, file/line references, tests, risks, exclusions, and an exit gate so later work can proceed without reconstructing the program context.

### Prompt Context

**User prompt (verbatim):**

{chr(10).join(('> ' + line) if line else '>' for line in PROMPT.split(chr(10)))}

**Follow-up user prompt (verbatim):**

> {FOLLOWUP_1}

**Follow-up user prompt (verbatim):**

> {FOLLOWUP_2}

**Assistant interpretation:** {data["interpretation"]}

**Inferred user intent:** {data["intent"]}

### What I did

{actions}
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

- {data["tricky"]}
- Ticket boundaries are most reliable when each names both its upstream contract and the next consumer.
- Documentation should distinguish observed runtime behavior from proposed APIs and future implementation choices.

### What was tricky to build

- {data["tricky"]}
- The guide had to be self-contained without copying the entire parent architecture. It summarizes prerequisites, then links each claim to the file that owns it.

### What warrants a second pair of eyes

- {data["risk"]}
- Review all proposed public schema/API names before implementation makes them expensive to change.

### What should be done in the future

- {data["future"]}
- Update this diary immediately when implementation reveals a false assumption or accepted contract change.

### Code review instructions

Start with the ticket's `design-doc/01-*.md`, then inspect these decision-shaping files:

{files}

Validate the design workspace with:

```bash
docmgr validate frontmatter --doc {diary}
docmgr doctor --ticket {ticket} --stale-after 30
```

Run the focused code/test commands listed in the guide before and after implementation.

### Technical details

- Parent program map: `OPTKIT-011/design-doc/04-backend-first-optimization-workbench-program-roadmap.md`.
- Ticket state: `index.md`, `tasks.md`, and `changelog.md` in this workspace.
- No production code behavior changed while writing this step.
'''
    diary.write_text(frontmatter + body)
    print(diary)
