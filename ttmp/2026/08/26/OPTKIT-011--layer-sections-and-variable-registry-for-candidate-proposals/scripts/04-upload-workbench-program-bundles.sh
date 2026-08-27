#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-}"
case "$MODE" in
  --dry-run) EXTRA=(--dry-run); label="dry-run" ;;
  --upload) EXTRA=(); label="upload" ;;
  *) echo "usage: $0 --dry-run|--upload" >&2; exit 2 ;;
esac

ROOT="optkit/ttmp/2026/08/26"

bundle() {
  local ticket="$1" dir="$2" name="$3"
  shift 3
  local log="$ROOT/$dir/various/remarkable-$label.log"
  mkdir -p "$(dirname "$log")"
  echo "=== $ticket: $name ($label) ===" | tee "$log"
  remarquee upload bundle "$@" \
    --name "$name" \
    --remote-dir "/ai/2026/08/26/$ticket" \
    --toc-depth 2 \
    --non-interactive \
    "${EXTRA[@]}" 2>&1 | tee -a "$log"
}

bundle OPTKIT-011 \
  "OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals" \
  "OPTKIT-011 Backend First Workbench Roadmap" \
  "$ROOT/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/index.md" \
  "$ROOT/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/design-doc/03-architect-brief-modular-optimization-workbench-and-candidate-authoring.md" \
  "$ROOT/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/design-doc/04-backend-first-optimization-workbench-program-roadmap.md"

bundle OPTKIT-012 \
  "OPTKIT-012--architecture-closure-and-optimization-workbench-contracts" \
  "OPTKIT-012 Architecture Contracts Guide" \
  "$ROOT/OPTKIT-012--architecture-closure-and-optimization-workbench-contracts/index.md" \
  "$ROOT/OPTKIT-012--architecture-closure-and-optimization-workbench-contracts/design-doc/01-intern-guide-to-optimization-workbench-architecture-and-contracts.md" \
  "$ROOT/OPTKIT-012--architecture-closure-and-optimization-workbench-contracts/reference/01-implementation-diary.md"

bundle OPTKIT-013 \
  "OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings" \
  "OPTKIT-013 Catalog and Bindings Guide" \
  "$ROOT/OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings/index.md" \
  "$ROOT/OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings/design-doc/01-intern-guide-to-optkit-catalogs-domains-bindings-and-candidate-intent.md" \
  "$ROOT/OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings/reference/01-implementation-diary.md"

bundle OPTKIT-014 \
  "OPTKIT-014--whole-pipeline-rag-configuration-and-graph-derivation" \
  "OPTKIT-014 Pipeline Config Guide" \
  "$ROOT/OPTKIT-014--whole-pipeline-rag-configuration-and-graph-derivation/index.md" \
  "$ROOT/OPTKIT-014--whole-pipeline-rag-configuration-and-graph-derivation/design-doc/01-intern-guide-to-pipelineconfig-layer-lenses-and-derived-graphs.md" \
  "$ROOT/OPTKIT-014--whole-pipeline-rag-configuration-and-graph-derivation/reference/01-implementation-diary.md"

bundle OPTKIT-015 \
  "OPTKIT-015--real-fusion-configuration-and-first-rag-optimization-catalog" \
  "OPTKIT-015 Fusion and RAG Catalog Guide" \
  "$ROOT/OPTKIT-015--real-fusion-configuration-and-first-rag-optimization-catalog/index.md" \
  "$ROOT/OPTKIT-015--real-fusion-configuration-and-first-rag-optimization-catalog/design-doc/01-intern-guide-to-fusion-configuration-rrf-and-the-first-rag-catalog.md" \
  "$ROOT/OPTKIT-015--real-fusion-configuration-and-first-rag-optimization-catalog/reference/01-implementation-diary.md"

bundle OPTKIT-016 \
  "OPTKIT-016--proposal-compiler-and-glazed-cli" \
  "OPTKIT-016 Proposal Compiler CLI Guide" \
  "$ROOT/OPTKIT-016--proposal-compiler-and-glazed-cli/index.md" \
  "$ROOT/OPTKIT-016--proposal-compiler-and-glazed-cli/design-doc/01-intern-guide-to-pure-proposal-compilation-and-glazed-cli-authoring.md" \
  "$ROOT/OPTKIT-016--proposal-compiler-and-glazed-cli/reference/01-implementation-diary.md"

bundle OPTKIT-017 \
  "OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence" \
  "OPTKIT-017 Proposal Sealing Guide" \
  "$ROOT/OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence/index.md" \
  "$ROOT/OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence/design-doc/01-intern-guide-to-proposal-sealing-candidate-manifests-and-durable-campaigns.md" \
  "$ROOT/OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence/reference/01-implementation-diary.md"

bundle OPTKIT-018 \
  "OPTKIT-018--candidate-projections-and-workbench-command-api" \
  "OPTKIT-018 Workbench API Guide" \
  "$ROOT/OPTKIT-018--candidate-projections-and-workbench-command-api/index.md" \
  "$ROOT/OPTKIT-018--candidate-projections-and-workbench-command-api/design-doc/01-intern-guide-to-candidate-read-projections-and-the-workbench-command-api.md" \
  "$ROOT/OPTKIT-018--candidate-projections-and-workbench-command-api/reference/01-implementation-diary.md"

bundle OPTKIT-019 \
  "OPTKIT-019--react-workbench-framework-and-rrf-vertical-slice" \
  "OPTKIT-019 React Workbench RRF Guide" \
  "$ROOT/OPTKIT-019--react-workbench-framework-and-rrf-vertical-slice/index.md" \
  "$ROOT/OPTKIT-019--react-workbench-framework-and-rrf-vertical-slice/design-doc/01-intern-guide-to-the-react-workbench-framework-and-rrf-vertical-slice.md" \
  "$ROOT/OPTKIT-019--react-workbench-framework-and-rrf-vertical-slice/reference/01-implementation-diary.md"

bundle OPTKIT-020 \
  "OPTKIT-020--asset-variable-proof-for-representation-prompts" \
  "OPTKIT-020 Prompt Asset Variable Guide" \
  "$ROOT/OPTKIT-020--asset-variable-proof-for-representation-prompts/index.md" \
  "$ROOT/OPTKIT-020--asset-variable-proof-for-representation-prompts/design-doc/01-intern-guide-to-artifact-valued-prompt-variables-and-bounded-previews.md" \
  "$ROOT/OPTKIT-020--asset-variable-proof-for-representation-prompts/reference/01-implementation-diary.md"
