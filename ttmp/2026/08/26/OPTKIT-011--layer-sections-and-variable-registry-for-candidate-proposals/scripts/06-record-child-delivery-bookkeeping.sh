#!/usr/bin/env bash
set -euo pipefail
ROOT="/home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/26"

record() {
  local ticket="$1" dir="$2" commit="$3" bundle="$4"
  local diary="$ROOT/$dir/reference/01-implementation-diary.md"
  local dry="$ROOT/$dir/various/remarkable-dry-run.log"
  local upload="$ROOT/$dir/various/remarkable-upload.log"
  docmgr doc relate --doc "$diary" \
    --file-note "$dry:Ticket bundle selection and destination dry-run evidence" \
    --file-note "$upload:Successful rendered PDF upload receipt"
  docmgr changelog update --ticket "$ticket" \
    --entry "Validated and committed the intern guide/diary package in $commit; uploaded $bundle.pdf to /ai/2026/08/26/$ticket" \
    --file-note "$diary:Validation commit and delivery record" \
    --file-note "$upload:Explicit successful reMarkable upload evidence"
}

record OPTKIT-012 "OPTKIT-012--architecture-closure-and-optimization-workbench-contracts" 428f6b8f5391dc851d989364b21a9d727c78cfc0 "OPTKIT-012 Architecture Contracts Guide"
record OPTKIT-013 "OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings" 428f6b8f5391dc851d989364b21a9d727c78cfc0 "OPTKIT-013 Catalog and Bindings Guide"
record OPTKIT-014 "OPTKIT-014--whole-pipeline-rag-configuration-and-graph-derivation" 428f6b8f5391dc851d989364b21a9d727c78cfc0 "OPTKIT-014 Pipeline Config Guide"
record OPTKIT-015 "OPTKIT-015--real-fusion-configuration-and-first-rag-optimization-catalog" 83d0f4f201ae59a0d8983d476e7873312fca8142 "OPTKIT-015 Fusion and RAG Catalog Guide"
record OPTKIT-016 "OPTKIT-016--proposal-compiler-and-glazed-cli" 83d0f4f201ae59a0d8983d476e7873312fca8142 "OPTKIT-016 Proposal Compiler CLI Guide"
record OPTKIT-017 "OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence" 83d0f4f201ae59a0d8983d476e7873312fca8142 "OPTKIT-017 Proposal Sealing Guide"
record OPTKIT-018 "OPTKIT-018--candidate-projections-and-workbench-command-api" ec784b0d84e47ad45a7429f98482e9678df5dac7 "OPTKIT-018 Workbench API Guide"
record OPTKIT-019 "OPTKIT-019--react-workbench-framework-and-rrf-vertical-slice" ec784b0d84e47ad45a7429f98482e9678df5dac7 "OPTKIT-019 React Workbench RRF Guide"
record OPTKIT-020 "OPTKIT-020--asset-variable-proof-for-representation-prompts" ec784b0d84e47ad45a7429f98482e9678df5dac7 "OPTKIT-020 Prompt Asset Variable Guide"
