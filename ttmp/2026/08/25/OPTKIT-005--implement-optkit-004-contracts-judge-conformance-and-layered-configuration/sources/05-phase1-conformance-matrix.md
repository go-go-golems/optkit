# OPTKIT-004 Phase 1 Judgekit Conformance Matrix

## Result

The existing OPTKIT-002 P7 implementation satisfies OPTKIT-004 Phase 1. Phase 1 requires conformance evidence and one stronger sealed-byte regression assertion; it does not require a second judge adapter.

## Requirement mapping

| OPTKIT-004 requirement | Implementation evidence | Test or validation evidence | Result |
|---|---|---|---|
| Decode a sealed historical answer trajectory | `rag-ttc/pkg/ttc/judgeinstrument/record.go`: `LoadSealedAnswer` loads and verifies the trajectory, requires exactly one `customer.answer.sealed` event, strictly decodes it, validates the record, and verifies nested deterministic artifacts. | `TestHistoricalAnswerCanBeRemeasuredUnderNewEpoch`; `TestSealedAnswerDecoderRejectsUnknownAndTrailingFields` | PASS |
| Recompute Judgekit instance identity from current content | `rag-ttc/pkg/ttc/judgeinstrument/instrument.go`: `BuildInstance` verifies the current contract digest, reconstructs evidence, and calls `eval.NewInstance`; it never trusts a caller-supplied instance digest. | `TestHistoricalAnswerCanBeRemeasuredUnderNewEpoch`; Judgekit `TestClaimJudgeRebindsCurrentInstanceIdentity` | PASS |
| Hide evidence from claim extraction | `judgekit/judging/claimjudge.go`: `ClaimExtractionInput` contains only ID, input, candidate, and metadata; `extractionInput` does not copy evidence or reference material. | Judgekit `TestClaimJudgeEndToEnd` inspects the extraction prompt and rejects evidence leakage. | PASS |
| Persist attributed protocol and report artifacts | `Instrument.Measure` stores `judgekit.instance/v1`, `judgekit.assessment/v1`, and `optkit.observation/v1` artifacts. Report validation requires instance, protocol, contract, and cache-mode agreement. | `TestHistoricalAnswerCanBeRemeasuredUnderNewEpoch`; Judgekit `TestClaimJudgeRecordsPromptModelAndCacheAttribution`; assessment report validation tests | PASS |
| Record prompt, model, usage, duration, and cache provenance | Judgekit `assessment.RunProvenance` records prompt template and rendered prompt digests, expected and observed model identity, cache mode, token usage, duration, and generation attempts. | `TestClaimJudgeRecordsPromptModelAndCacheAttribution`; `TestClaimJudgeRejectsObservedModelMismatch` | PASS |
| Support cache-bypass reliability probes | `Instrument.Measure` accepts only `CacheUse` and `CacheBypass`; `evaluate` requires `ConfigurableJudge` for bypass. Judgekit bypass neither reads nor writes cache. | RAG-TTC `TestCacheBypassProducesFreshAttributedRepeat`; Judgekit `TestClaimJudgeBypassesCacheWhenRequested` | PASS |
| Map measurements to Optkit epochs | `Instrument.observations` calls `measure.NewEpoch` with construct, instrument, protocol digest, implementation, contract calibration digest, and answer redaction schema. | `TestHistoricalAnswerCanBeRemeasuredUnderNewEpoch` requires changed protocol to produce a changed epoch. | PASS |
| Preserve deterministic application evidence separately | `SealAnswerEpisode` emits `customer.answer.contract` separately and adds its artifact to `SealedAnswer.DeterministicArtifacts`; observations retain those refs in diagnostics rather than replacing them with judge scores. | `TestHistoricalAnswerCanBeRemeasuredUnderNewEpoch` rejects deterministic artifacts appearing as substituted judge evidence. | PASS |
| Preserve typed judge failure, missing output, and inapplicability | `failedOutcome` emits `judge_failed`; observation construction distinguishes failed, unknown/missing output, inapplicable, and measured states. | `TestJudgeFailureAndMissingDimensionBecomeTypedObservations` | PASS |
| Remeasure without retrieval or answer execution | `Instrument.Measure` accepts only an artifact store and sealed `episode.Result`; it has no retrieval service, customer application, generator factory, or prepared system dependency. The regression test snapshots exact trajectory and output bytes around two measurements. | `TestHistoricalAnswerCanBeRemeasuredUnderNewEpoch` verifies identical trajectory/output bytes after both epochs. | PASS |
| Preserve old observations | Observations are immutable content-addressed artifacts. The second measurement creates a different epoch and new observation IDs. | `TestHistoricalAnswerCanBeRemeasuredUnderNewEpoch` retains and deep-compares the old observation after the second measurement. | PASS |

## Acceptance evidence

The acceptance case performs these operations:

1. seal one historical answer episode once;
2. capture exact trajectory and output bytes;
3. measure under protocol `ttc-faithfulness-v1`;
4. measure the same `episode.Result` under protocol `ttc-faithfulness-v2`;
5. require identical Judgekit instance identity because historical content is unchanged;
6. require different Optkit epoch identities because the protocol changed;
7. require the first observation to remain byte-for-byte unchanged in memory;
8. require the sealed trajectory and output artifacts to retain exact bytes;
9. verify all referenced artifacts.

No retrieval or answer executor is available to `Instrument.Measure`, so a successful second measurement cannot invoke product execution through a hidden callback.

## Deferred work

- Gold-set calibration views belong to the specialist Judgekit UI slice.
- Reliability summaries across multiple perturbation probes belong to later measurement projectors.
- Cross-trust signatures and remote custody remain explicitly out of scope.
- Promotion policy must consume compatible epochs but remains a Phase 8 concern.
