#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
OPTKIT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)
WORKSPACE=$(dirname "$OPTKIT")
RAG_TTC="$WORKSPACE/rag-ttc"
PHASE2="$OPTKIT/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice"

(cd "$OPTKIT" && "$PHASE2/scripts/03-sync-rag-semantic-fixture.sh" --check)

(
  cd "$RAG_TTC"
  GOWORK=off go test \
    ./pkg/ttc/optimization \
    ./pkg/ttc/search \
    ./pkg/ttc/customerapp \
    ./pkg/ttc/judgeinstrument \
    ./pkg/ttc/optkitcampaign \
    -count=1
  GOWORK=off go test -race \
    ./pkg/ttc/optimization \
    ./pkg/ttc/judgeinstrument \
    ./pkg/ttc/optkitcampaign \
    -count=1
  GOWORK=off go vet \
    ./pkg/ttc/optimization \
    ./pkg/ttc/search \
    ./pkg/ttc/customerapp \
    ./pkg/ttc/judgeinstrument \
    ./pkg/ttc/optkitcampaign
)

fixture="$RAG_TTC/pkg/ttc/optimization/testdata/rag-optimization-semantic-fixture-v1.json"
digest=$(sha256sum "$fixture" | awk '{print $1}')
constant=$(awk -F'"' '/const OptimizationFixtureSHA256/ {print $2}' "$RAG_TTC/pkg/ttc/optimization/fixture.go")
[[ "$digest" == "$constant" ]]

printf 'OPTIMIZATION_FIXTURE_SCHEMA=rag-ttc.optimization-semantic-fixture/v1\n'
printf 'OPTIMIZATION_FIXTURE_SHA256=%s\n' "$digest"
printf 'BASE_RETRIEVAL_FIXTURE_SHA256=2fa045999a8a89039e00dd60b3fec2bc17b732d557eb00746e207620a5fbdc7f\n'
printf 'PHASE_0_VALIDATION=PASS\n'
