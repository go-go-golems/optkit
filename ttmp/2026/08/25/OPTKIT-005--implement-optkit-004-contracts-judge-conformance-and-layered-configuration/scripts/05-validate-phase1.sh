#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
OPTKIT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)
WORKSPACE=$(dirname "$OPTKIT")
RAG_TTC="$WORKSPACE/rag-ttc"
JUDGEKIT="$WORKSPACE/judgekit"

(
  cd "$JUDGEKIT"
  GOWORK=off go test ./judging ./assessment ./protocol -count=1
  GOWORK=off go test -race ./judging ./assessment ./protocol -count=1
  GOWORK=off go vet ./judging ./assessment ./protocol
)

(
  cd "$RAG_TTC"
  GOWORK=off go test ./pkg/ttc/judgeinstrument -run 'TestHistoricalAnswerCanBeRemeasuredUnderNewEpoch|TestCacheBypassProducesFreshAttributedRepeat|TestJudgeFailureAndMissingDimensionBecomeTypedObservations|TestFromCustomerResultCapturesOnlyAdmittedEvidence|TestSealedAnswerDecoderRejectsUnknownAndTrailingFields' -count=1
  GOWORK=off go test -race ./pkg/ttc/judgeinstrument -count=1
  GOWORK=off go vet ./pkg/ttc/judgeinstrument
)

if rg -n 'Search\(|Retrieve\(|RunTurn\(|prepared\.Run|system\.Prepared' "$RAG_TTC/pkg/ttc/judgeinstrument/instrument.go"; then
  echo 'judge instrument measurement path contains a product execution call' >&2
  exit 1
fi

printf 'SEALED_ANSWER_REMEASUREMENT=PASS\n'
printf 'EVIDENCE_HIDDEN_EXTRACTION=PASS\n'
printf 'CACHE_BYPASS_PROBE=PASS\n'
printf 'TYPED_MISSING_FAILURE_STATES=PASS\n'
printf 'PHASE_1_VALIDATION=PASS\n'
