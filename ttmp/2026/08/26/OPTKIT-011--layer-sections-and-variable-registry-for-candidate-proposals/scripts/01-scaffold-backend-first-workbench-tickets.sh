#!/usr/bin/env bash
set -euo pipefail

# Scaffold the backend-first implementation tickets agreed in OPTKIT-011.
# Run from the use-optkit workspace root, where .ttmp.yaml points at optkit/ttmp.

create_ticket() {
  local id="$1" title="$2" topics="$3" guide="$4"
  shift 4

  docmgr ticket create-ticket --ticket "$id" --title "$title" --topics "$topics"
  docmgr doc add --ticket "$id" --doc-type design-doc --title "$guide" \
    --summary "Intern-oriented architecture, design, and implementation guide for $title."
  docmgr doc add --ticket "$id" --doc-type reference --title "Implementation Diary" \
    --summary "Chronological research, design, implementation, validation, and delivery record for $id."

  for task in "$@"; do
    docmgr task add --ticket "$id" --text "$task"
  done
}

create_ticket OPTKIT-012 \
  "Architecture Closure and Optimization Workbench Contracts" \
  "architecture,design,implementation,optkit,rag-ttc" \
  "Intern Guide to Optimization Workbench Architecture and Contracts" \
  "Review and accept ADRs for aggregate configuration, catalog/bindings, value schemas, candidate identity, provenance, and API boundaries" \
  "Approve concrete Go contracts for Catalog, Bindings, PipelineConfig, CompileProposal, and SealProposal" \
  "Choose the v1 manifest and store compatibility policy before implementation" \
  "Validate package dependencies and prove the proposed package graph is acyclic" \
  "Use the accepted contracts as entry criteria for OPTKIT-013 through OPTKIT-018"

create_ticket OPTKIT-013 \
  "Optkit Semantic Catalog and Executable Variable Bindings" \
  "architecture,design,implementation,optkit" \
  "Intern Guide to Optkit Catalogs Domains Bindings and Candidate Intent" \
  "Implement ordered sections and deterministic semantic catalog serialization" \
  "Implement complete value/domain descriptors including float boolean choice string and asset values" \
  "Implement executable Bindings[C] constructed from typed space.Variable declarations" \
  "Extend structured Candidate intent and lock identity behavior with tests" \
  "Migrate numbergame as the catalog and serialized-mutation proof case" \
  "Run focused and full Optkit validation"

create_ticket OPTKIT-014 \
  "Whole-Pipeline RAG Configuration and Graph Derivation" \
  "architecture,design,implementation,optkit,rag-ttc" \
  "Intern Guide to PipelineConfig Layer Lenses and Derived Graphs" \
  "Introduce the accepted whole-pipeline semantic configuration model" \
  "Move retrieval configuration ownership out of the campaign adapter without creating package cycles" \
  "Implement layer-local lenses lifted to PipelineConfig" \
  "Derive the complete optimization graph from semantic configuration values" \
  "Remove independently authored graph/config drift from the new authoring path" \
  "Test identity derivation direct changes transitive invalidation and upstream reuse"

create_ticket OPTKIT-015 \
  "Real Fusion Configuration and First RAG Optimization Catalog" \
  "architecture,design,implementation,optkit,rag-ttc" \
  "Intern Guide to Fusion Configuration RRF and the First RAG Catalog" \
  "Introduce typed FusionConfig with a positive float64 RRF rank constant" \
  "Plumb FusionConfig through fixture preparation into actual WeightedRRF execution" \
  "Rename and document the current retrieval final-result limit accurately" \
  "Register retrieval final limit and fusion.rrf_k with reviewed documentation and bindings" \
  "Prove registered mutations change runtime behavior and graph invalidation" \
  "Update semantic fixtures identities and focused parity tests"

create_ticket OPTKIT-016 \
  "Proposal Compiler and Glazed CLI" \
  "architecture,design,implementation,optkit,rag-ttc" \
  "Intern Guide to Pure Proposal Compilation and Glazed CLI Authoring" \
  "Implement pure CompileProposal with canonical decoding validation and normalization" \
  "Return before/after values graph diff invalidation plan diagnostics and preview capabilities" \
  "Reject unknown duplicate ill-typed out-of-domain and no-op mutations deterministically" \
  "Add Glazed catalog list/show and proposal compile commands with structured output" \
  "Route existing config inspection paths through shared application services where appropriate" \
  "Prove repeated compilation performs no durable writes"

create_ticket OPTKIT-017 \
  "Proposal Sealing Candidate Manifests and Campaign Persistence" \
  "architecture,design,implementation,optkit,rag-ttc" \
  "Intern Guide to Proposal Sealing Candidate Manifests and Durable Campaigns" \
  "Implement SealProposal using normal Optkit PatchBuilder and snapshot mechanics" \
  "Add a strict candidate authoring block to the experiment manifest" \
  "Compile manifest candidates through the shared proposal compiler" \
  "Persist patches snapshots structured candidates and catalog provenance in campaign state" \
  "Emit CandidateProposed and preserve deterministic sealing semantics" \
  "Prove sealed campaigns remain explainable after the source manifest is removed"

create_ticket OPTKIT-018 \
  "Candidate Projections and Workbench Command API" \
  "architecture,design,implementation,optkit,rag-ttc" \
  "Intern Guide to Candidate Read Projections and the Workbench Command API" \
  "Project structured candidate intent and mutation summaries from sealed campaign facts" \
  "Expose the current optimization catalog for authoring" \
  "Define an application command service for compile preview and seal operations" \
  "Add HTTP adapters without moving business logic into handlers" \
  "Preserve specialistapi as the durable historical read boundary" \
  "Define authorization sensitivity idempotency and error contracts before asset editing"

create_ticket OPTKIT-019 \
  "React Workbench Framework and RRF Vertical Slice" \
  "architecture,design,implementation,optkit,rag-ttc,ui" \
  "Intern Guide to the React Workbench Framework and RRF Vertical Slice" \
  "Implement generic catalog-driven variable editors with safe fallbacks" \
  "Implement the frontend WorkbenchRegistry for specialized editors inspectors outputs and previews" \
  "Extract a reusable WorkbenchShell and shared candidate-authoring components" \
  "Build the fusion.rrf_k specialized inspector and deterministic preview" \
  "Complete the case proposal trial verdict RRF vertical slice" \
  "Pin browser calculations to recorded Go results and validate accessibility and deep links"

create_ticket OPTKIT-020 \
  "Asset Variable Proof for Representation Prompts" \
  "architecture,design,implementation,optkit,rag-ttc,ui" \
  "Intern Guide to Artifact-Valued Prompt Variables and Bounded Previews" \
  "Register representations.summary_prompt as a sensitivity-aware artifact variable" \
  "Implement old/new artifact materialization and side-by-side prompt diff authoring" \
  "Compute and display the truthful upstream recomputation bill" \
  "Implement a bounded server preview for one selected chunk" \
  "Persist prompt assets and provenance without placing content in patch metadata" \
  "Prove the workbench supports scalar and asset mutations through one core workflow"
