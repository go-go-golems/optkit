#!/usr/bin/env bash
set -euo pipefail

WORKSPACE="/home/manuel/workspaces/2026-08-24/use-optkit"
cd "$WORKSPACE/rag-ttc"

GOWORK=off go test ./pkg/ttc/optimization \
  -run 'TestRAGRegistry(CatalogIsCompleteAndDeterministic|BindingsApplyPureAndDurableMutations|BindingsRejectWrongTypesAndDomains)$' \
  -count=1 -v

GOWORK=off go test ./pkg/ttc/optkitcampaign \
  -run 'Test(SemanticFixtureExecutesConfiguredFusionCoordinate|FinalResultLimitDoesNotChangeChannelTopKOrFusion)$' \
  -count=1 -v
