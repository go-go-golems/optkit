# P1 Cross-Product Semantic Fixture Validation

- Fixture schema: `rag.semantic-fixture/v1`
- Fixture corpus: `cross-product-tiny-corpus-v1`
- Canonical SHA-256: `2fa045999a8a89039e00dd60b3fec2bc17b732d557eb00746e207620a5fbdc7f`
- Canonical source: `sources/rag-semantic-fixture-v1.json`
- RAG-TTC copy: `pkg/ttc/search/testdata/rag-semantic-fixture-v1.json`
- Coinvault copy: `internal/knowledge/testdata/rag-semantic-fixture-v1.json`

## Product commits

- RAG-TTC: `b8aaf41d9a2ded9daca400cb16f9083fbcf37225` — `test(rag): freeze cross-product semantic fixture`
- Coinvault: `e3090be05b657b4ce6f0015028f3eba562ef909b` — `test(rag): adopt cross-product semantic fixture`

## Verified laws

- Stable lexical and vector candidate IDs and ranks.
- Deterministic fused and returned chunk order.
- Public access removes the analyst-only vector candidate.
- Product-role filtering retains only the public product candidate.
- Evidence labels are stable and repeated chunks reuse labels.
- Item budgets reject a third distinct chunk.
- Positive, comparison, authorization-negative, and no-answer cases remain explicit.
- Answer expectations retain exact citation labels and abstention case identity.
- Product copies are byte-identical to the canonical ticket fixture.

## Commands

```bash
cd optkit
./ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice/scripts/03-sync-rag-semantic-fixture.sh --check

cd ../rag-ttc
GOWORK=off go test ./pkg/ttc/search -count=1
GOWORK=off go test ./... -count=1

cd ../coinvault
GOWORK=off go test ./internal/knowledge -count=1
GOWORK=off go test ./... -count=1
```

All final commands passed. RAG-TTC's first full run observed the unrelated existing heartbeat timing flake `heartbeat timeout was not observable`; the isolated test passed 10/10 and the full rerun passed.
