#!/usr/bin/env bash
set -euo pipefail

ROOT="/home/manuel/workspaces/2026-08-24/use-optkit"
ROADMAP="$ROOT/optkit/ttmp/2026/08/26/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals/design-doc/04-backend-first-optimization-workbench-program-roadmap.md"

relate_ticket() {
  local ticket="$1" dir="$2" guide="$3"
  shift 3
  local guide_path="$ROOT/optkit/ttmp/2026/08/26/$dir/design-doc/$guide"
  local diary_path="$ROOT/optkit/ttmp/2026/08/26/$dir/reference/01-implementation-diary.md"

  docmgr task add --ticket "$ticket" --text "Write validate and review the intern architecture/design/implementation guide"
  docmgr task add --ticket "$ticket" --text "Dry-run and publish the guide/diary bundle to the ticket reMarkable folder"

  docmgr doc relate --doc "$guide_path" \
    --file-note "$ROADMAP:Parent program goals dependencies exclusions and exit gates" \
    "$@"
  docmgr doc relate --doc "$diary_path" \
    --file-note "$guide_path:Primary design deliverable whose research and delivery this diary records" \
    --file-note "$ROADMAP:Overall program context that keeps the ticket aligned"
  docmgr changelog update --ticket "$ticket" \
    --entry "Created a substantive ticket description, actionable implementation tasks, a detailed intern architecture/design/implementation guide, and a strict research diary" \
    --file-note "$guide_path:Primary evidence-backed implementation guide" \
    --file-note "$diary_path:Chronological research and continuation record"
}

relate_ticket OPTKIT-012 \
  "OPTKIT-012--architecture-closure-and-optimization-workbench-contracts" \
  "01-intern-guide-to-optimization-workbench-architecture-and-contracts.md" \
  --file-note "$ROOT/optkit/space/variable.go:Current typed variable and descriptor contract" \
  --file-note "$ROOT/optkit/space/patch.go:Canonical durable mutation boundary" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/optimization/graph.go:Current graph identity and dependency resolution" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/experimentworkbench/service.go:Current application service boundary"

relate_ticket OPTKIT-013 \
  "OPTKIT-013--optkit-semantic-catalog-and-executable-variable-bindings" \
  "01-intern-guide-to-optkit-catalogs-domains-bindings-and-candidate-intent.md" \
  --file-note "$ROOT/optkit/space/variable.go:Typed variable declaration extended by catalog metadata" \
  --file-note "$ROOT/optkit/space/domain.go:Lossy current domain descriptors to replace" \
  --file-note "$ROOT/optkit/space/patch.go:Durable typed assignment path bindings must preserve" \
  --file-note "$ROOT/optkit/examples/numbergame/demo.go:Complete generic proof consumer"

relate_ticket OPTKIT-014 \
  "OPTKIT-014--whole-pipeline-rag-configuration-and-graph-derivation" \
  "01-intern-guide-to-pipelineconfig-layer-lenses-and-derived-graphs.md" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/optimization/contracts.go:Canonical layer vocabulary" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/optimization/graph.go:Graph derivation target" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/optkitcampaign/system.go:Retrieval-only snapshot asymmetry" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/experimentworkbench/manifest.go:Current duplicated config/graph authoring"

relate_ticket OPTKIT-015 \
  "OPTKIT-015--real-fusion-configuration-and-first-rag-optimization-catalog" \
  "01-intern-guide-to-fusion-configuration-rrf-and-the-first-rag-catalog.md" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/search/search.go:Runtime float64 RRF and final-limit semantics" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/search/service.go:Actual WeightedRRF execution" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/search/semantic_fixture.go:Hardcoded RRF 60 to parameterize" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/optimization/invalidation.go:Required graph plan evidence"

relate_ticket OPTKIT-016 \
  "OPTKIT-016--proposal-compiler-and-glazed-cli" \
  "01-intern-guide-to-pure-proposal-compilation-and-glazed-cli-authoring.md" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/experimentworkbench/service.go:Application layer for pure compilation" \
  --file-note "$ROOT/rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/config.go:Existing Glazed config commands" \
  --file-note "$ROOT/rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/command.go:Command composition" \
  --file-note "$ROOT/optkit/space/patch.go:Write path explicitly excluded from drafts"

relate_ticket OPTKIT-017 \
  "OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence" \
  "01-intern-guide-to-proposal-sealing-candidate-manifests-and-durable-campaigns.md" \
  --file-note "$ROOT/optkit/space/patch.go:Canonical sealing mechanics" \
  --file-note "$ROOT/optkit/examples/numbergame/demo.go:Candidate event precedent" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/experimentworkbench/manifest.go:Strict manifest input boundary" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/optkitcampaign/campaign.go:Durable spec and journal ordering"

relate_ticket OPTKIT-018 \
  "OPTKIT-018--candidate-projections-and-workbench-command-api" \
  "01-intern-guide-to-candidate-read-projections-and-the-workbench-command-api.md" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/specialistapi/projector.go:Historical projection boundary" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/specialistapi/types.go:Read response contracts" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/specialistapi/http.go:GET-only HTTP routes" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/experimentworkbench/service.go:Transport-independent applications"

relate_ticket OPTKIT-019 \
  "OPTKIT-019--react-workbench-framework-and-rrf-vertical-slice" \
  "01-intern-guide-to-the-react-workbench-framework-and-rrf-vertical-slice.md" \
  --file-note "$ROOT/rag-ttc/apps/specialist/web/src/layerwidgets/index.tsx:Existing renderer registry precedent" \
  --file-note "$ROOT/rag-ttc/apps/specialist/web/src/screens/LabScreen.tsx:Numbergame reference experience" \
  --file-note "$ROOT/rag-ttc/apps/specialist/web/src/screens/ComparisonScreen.tsx:Historical verdict/evidence screen" \
  --file-note "$ROOT/rag-ttc/apps/specialist/web/src/api/specialistApi.ts:RTK Query transport boundary"

relate_ticket OPTKIT-020 \
  "OPTKIT-020--asset-variable-proof-for-representation-prompts" \
  "01-intern-guide-to-artifact-valued-prompt-variables-and-bounded-previews.md" \
  --file-note "$ROOT/optkit/artifact:Content-addressed sensitivity-aware artifact contracts" \
  --file-note "$ROOT/optkit/space/patch.go:Assignment artifact-ref behavior" \
  --file-note "$ROOT/rag-ttc/pkg/ttc/search/types.go:Recorded representation summaries" \
  --file-note "$ROOT/rag-ttc/apps/specialist/web/src/components/Artifact.tsx:Current artifact preview fallback"
