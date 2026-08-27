---
Title: Specialist UI Backend API and Projector Design
Ticket: OPTKIT-007
Status: active
Topics:
    - implementation
    - optkit
    - rag-ttc
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://optkit/query/service.go
      Note: Bounded read precedent
    - Path: repo://rag-ttc/pkg/ttc/specialistapi/types.go
      Note: Versioned DTO contract
ExternalSources: []
Summary: Versioned read models, projector boundaries, GET-only routes, and frontend handoff for the RAG-TTC specialist UI.
LastUpdated: 2026-08-26T00:30:00Z
WhatFor: Define how verified campaign and artifact facts become bounded UI-facing RAG views.
WhenToUse: Read before changing specialist projectors, HTTP routes, pagination, previews, or frontend API integration.
---


# Specialist UI Backend API and Projector Design

## Executive summary

The specialist backend is a read-only translation layer over durable Optkit campaign facts and RAG-TTC artifact schemas. It does not run experiments, mutate campaigns, infer missing measurements as zero, or ask the browser to decode trajectories. The OPTKIT-006 workbench remains the operator interface and backend oracle; this ticket exposes the same facts as stable HTTP DTOs for a frontend engineer.

The implementation lives in RAG-TTC because campaign arms, retrieval stages, evidence, per-arm configuration graphs, and paired cases are product semantics. It reads Optkit journals and CAS artifacts through interfaces, verifies the journal before projecting, strictly decodes known payloads, and reports typed diagnostics when facts are missing or unsupported.

## Architecture

```text
Glazed workbench / durable campaign producer
                 |
                 v
      Optkit journal + artifact CAS
                 |
          verify + decode
                 |
                 v
      RAG-TTC specialist projectors
       /        |         |       \
  cockpit   comparison pipeline provenance
       \        |         |       /
          GET-only HTTP transport
                 |
                 v
       frontend engineer / browser
```

Authority is one-way. The journal and artifacts are facts; projectors interpret known schemas; HTTP transports bounded DTOs; the browser presents them.

## Frozen API namespace

```text
GET /api/rag/v1/health
GET /api/rag/v1/campaigns/{campaign}/cockpit
GET /api/rag/v1/campaigns/{campaign}/comparisons/{baseline}/{treatment}
GET /api/rag/v1/campaigns/{campaign}/cases?baseline=...&treatment=...&after=...&limit=...
GET /api/rag/v1/campaigns/{campaign}/episodes/{episode}/pipeline
GET /api/rag/v1/campaigns/{campaign}/provenance/episode/{episode}
```

There are no POST, PUT, PATCH, or DELETE routes.

## DTO contracts

All responses include `api_version`. Top-level projector responses also include a schema string and `through_seq`, the journal sequence through which the view was constructed.

### Campaign cockpit

The cockpit reports lifecycle, integrity, manifest identity, arms, graph IDs, cases, estimates, budget, episode counts, and diagnostics. It answers whether the run completed cleanly and which arm comparison should be opened.

### Comparison

A comparison selects two arms from one campaign. It includes:

- baseline and treatment summaries;
- direct `optimization.ConfigDiff`;
- transitive `optimization.InvalidationPlan`;
- paired target-coverage metric summaries;
- paired case rows with explicit measurement status;
- diagnostics and journal sequence.

Configuration differences are never recomputed in HTTP code; projectors call the Phase 2 optimization APIs.

### Case page

Case pages use an opaque base64url cursor encoding a server-side offset. Default limit is 50 and maximum is 100. The first fixture is small and folded in memory; production-scale database pagination remains a later storage optimization.

### Pipeline

The pipeline decodes one completed episode result and sealed trajectory. It emits ordered retrieval stages, stage counts, chunk IDs, bounded payload previews, final output, and diagnostics. Only observed stages are shown.

### Provenance

The first provenance kind is `episode`. It connects:

```text
manifest -> arm graph -> immutable snapshot -> episode spec
         -> completion -> result -> trajectory -> output/stage artifacts
```

Other provenance kinds are rejected until their contracts exist.

## Preview policy

- Public and internal JSON artifacts may be previewed.
- Confidential and restricted artifacts are never previewed.
- Preview maximum is 4096 bytes.
- Oversized, unknown, malformed, or unreadable artifacts return refs plus a typed reason.
- Filesystem paths are never returned.

## Error envelope

```json
{
  "error": {
    "code": "invalid_campaign_id",
    "message": "campaign ID: ..."
  }
}
```

Bad IDs, cursors, limits, and missing query parameters return 400. Missing campaigns, arms, and episodes return 404. Verified facts that cannot be decoded due to server defects return 500 with a stable code.

## Implementation phases

1. Freeze DTO, cursor, diagnostics, route, and preview contracts.
2. Implement verified cockpit and comparison projectors.
3. Implement pipeline and episode provenance projectors.
4. Add strict GET-only `http.ServeMux` handlers and serve command.
5. Characterize campaign → comparison → case → pipeline → provenance.
6. Publish frontend examples, TypeScript sketches, and validation evidence.

## Security decisions

- Projectors verify the campaign journal before presenting integrity as true.
- Artifact refs are returned, but bytes follow sensitivity and size policy.
- The server binds to an explicit operator address; the default is loopback.
- CSP and UI static hosting remain frontend integration concerns; JSON handlers set nosniff, frame, referrer, and cache headers.
- No CORS wildcard is added. Same-origin deployment is the default.

## Primary file references

- `rag-ttc/pkg/ttc/specialistapi/types.go` — DTO contract.
- `rag-ttc/pkg/ttc/specialistapi/projector.go` — verified campaign fold.
- `rag-ttc/pkg/ttc/specialistapi/pipeline.go` — trajectory projection.
- `rag-ttc/pkg/ttc/specialistapi/http.go` — GET-only transport.
- `rag-ttc/pkg/ttc/experimentworkbench/` — CLI oracle and deterministic fixture.
- `rag-ttc/pkg/ttc/optkitcampaign/campaign.go` — stored campaign schemas.
- `optkit/query/service.go` — bounded generic read precedent.
