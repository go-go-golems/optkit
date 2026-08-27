# Changelog

## 2026-08-26

- Initial workspace created


## 2026-08-26

Created ticket and wrote the intern guide: pbui-chat vocabulary export, mentions, verb-router families with actor attribution, draft co-editing over the shared workbench document, Proposer/approval provenance, authorization-as-safety, implementation sketch and exclusions.

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/optkit/ttmp/2026/08/26/OPTKIT-024--agent-seat-over-the-workbench-vocabulary/design-doc/01-intern-guide-agent-seat-over-the-workbench-vocabulary.md — Primary design guide


## 2026-08-26

Uploaded the four-guide PBUI track bundle (OPTKIT-021-024 PBUI Workbench Track Guides.pdf) to /ai/2026/08/26/OPTKIT-021-PBUI-workbench on the reMarkable.


## 2026-08-26

Guide retargeted to pbui 0.8.0: kernel rules with goldens-first row spec, translator edges, primary invocation, capability-gated seal, registry-generated vocabulary (pbui commit 6efeaeb)


## 2026-08-27

P1: agent vocabulary export build step + golden (rag-ttc c2204f985) — registry.vocabulary() composed with translator edges, type docs, typed verb table

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/apps/workbench/web/src/agent/vocabulary.ts — The vocabulary generator


## 2026-08-27

P2+P3 core: pbui-chat 0.3.0 reference codec (user ruling) + 'any' field type; rag-ttc chat layer mounted — one routed, attributed verb path live (pbui e09ab55/917c04a, rag-ttc b76000a1d)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/apps/workbench/web/src/chat/router.ts — The single dispatcher — codec, validation, attribution


## 2026-08-27

P4: approval as one-shot capability grant — router gate, SealBar + global surfaces, approvalId in the verb, proposer llm on approved seals (rag-ttc c8bbe78f7); §4 open question resolved as journal-metadata-only, flagged for ADR review

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/apps/workbench/web/src/chat/router.ts — The approval gate — park, verify, consume


## 2026-08-27

P5: agent principal — multi-principal bearer auth on serve, agent granted compile/preview and NOT seal, proven 403 at HTTP level (rag-ttc 456a1d037)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/cmd/rag-ttc/cmds/experiments/optkitrag/principals.go — The grants and the multi-token authenticator

