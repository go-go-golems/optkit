---
Title: Pragmatic Optkit Unification and Coinvault-to-RAG-TTC Port
Ticket: OPTKIT-001
Status: active
Topics:
    - optkit
    - rag
    - rag-ttc
    - coinvault
    - judgekit
    - architecture
    - migration
    - local-development
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Current-state assessment and product-first plan for transferring Coinvault's proven RAG discipline into RAG-TTC as Optkit's first real vertical slice.
LastUpdated: 2026-08-24T22:40:00-04:00
WhatFor: Navigate the assessment, implementation guide, evidence inventory, diary, and delivery state for OPTKIT-001.
WhenToUse: Start here before implementing or reviewing the Optkit import or Coinvault-to-RAG-TTC RAG work.
---

# Pragmatic Optkit Unification and Coinvault-to-RAG-TTC Port

## Overview

This ticket assesses the clean-slate Optkit architecture, supplied implementation archive, and current RagKit, RagOpt, Judgekit, Coinvault, and RAG-TTC repositories. It recommends a pragmatic sequence: land the green Optkit archive, transfer Coinvault's verified retrieval/evidence/evaluation laws into a canonical RAG-TTC service, and use that service as Optkit's first real product campaign before bulk kit consolidation.

## Primary conclusion

The checked-out `optkit` repository is still a template, while the supplied archive is a tested local vertical slice. RAG-TTC already has substantial RAG infrastructure; the next work is service convergence, runtime attribution, stage diagnostics, deterministic evaluation, and Optkit registration—not a ground-up RAG rewrite.

## Key links

- [Primary analysis, design, and implementation guide](design-doc/01-optkit-current-state-and-pragmatic-rag-implementation-guide.md)
- [Investigation diary](reference/01-investigation-diary.md)
- [Repository inventory](sources/01-repository-inventory.md)
- [Inventory generator](scripts/01-repository-inventory.sh)
- [Task checklist](tasks.md)
- [Changelog](changelog.md)

## Status

- Architecture and repository assessment: complete.
- Intern implementation guide: complete.
- Coinvault/RAG-TTC/Judgekit/RagKit/RagOpt focused and full relevant tests: green.
- `docmgr doctor`: clean.
- reMarkable delivery: uploaded successfully to `/ai/2026/08/24/OPTKIT-001`.
- Product implementation: intentionally deferred to follow-up tickets.

## Accepted direction

1. Import the supplied Optkit source with provenance and preserve its green baseline.
2. Freeze deterministic RAG-TTC semantic fixtures.
3. Extract a canonical TTC retrieval service below the model tool and transport.
4. Add Coinvault-style bundle/config/pipeline/evidence identities and stage diagnosis.
5. Build strict deterministic retrieval and answer-contract evaluation.
6. Register a fixed multi-arm RAG-TTC retrieval campaign in Optkit.
7. Add Judgekit only under the lightweight local research-attribution model.
8. Move exercised generic code, switch last callers, delete old paths, and archive kits.
