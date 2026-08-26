---
Title: Implementation Diary
Ticket: OPTKIT-008
Status: active
Topics:
    - implementation
    - rag-ttc
    - ui
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ws://rag-ttc/apps/specialist/web/src/api/specialistApi.ts
      Note: RTK Query endpoints and error normalization for specialist-api/v1 (commit 09a911ae5)
    - Path: ws://rag-ttc/apps/specialist/web/src/screens/ComparisonScreen.tsx
      Note: Config diff, invalidation plan, cursor-paginated paired cases (commit 09a911ae5)
    - Path: ws://rag-ttc/apps/specialist/web/src/styles/system.css
      Note: Monochrome Mac System 1 design system and material vocabulary (commit 09a911ae5)
    - Path: ws://rag-ttc/apps/specialist/web/src/test/comparison.test.tsx
      Note: Fixture-driven MSW tests for ordering, change-kind, pagination (commit 09a911ae5)
ExternalSources: []
Summary: Chronological diary of building the read-only RAG-TTC specialist frontend (cockpit → comparison → paired cases → pipeline → provenance).
LastUpdated: 2026-08-25T21:25:00-04:00
WhatFor: Record what was built, what failed, and how to review the specialist UI work.
WhenToUse: Read before resuming or reviewing OPTKIT-008 frontend work.
---


# Diary

## Goal

Capture the implementation journey of the first read-only RAG-TTC specialist workflow UI: campaign cockpit → baseline/challenger comparison → paired case → episode pipeline → provenance/configuration diff, built against the OPTKIT-007 API handoff and its archived fixtures, in a monochrome Mac System 1 visual style.

## Step 1: Scaffold the specialist app, all five screens, and fixture-driven tests

Read the OPTKIT-007 frontend handoff and all six archived API fixtures, then built the complete first pass of the app in `rag-ttc/apps/specialist/web`: a Vite + React 19 + TypeScript SPA with RTK Query for data fetching, react-router for deep-linkable screens, a hand-written monochrome "Mac System 1 without window chrome" design system, and a vitest + MSW test suite that replays the archived fixtures. All 28 tests pass and the production build succeeds.

### Prompt Context

**User prompt (verbatim):** "Build the first read-only RAG-TTC specialist workflow: Campaign cockpit → baseline/challenger comparison → paired case → episode pipeline → provenance/configuration diff. Start with the handoff document: optkit/ttmp/2026/08/25/OPTKIT-007--build-rag-ttc-specialist-ui-backend-and-frontend-handoff/reference/02-frontend-engineer-api-handoff.md Use the real API responses in: .../sources/api-fixtures/ Use those fixtures for component tests, Storybook, or MSW before connecting to the live server. Important constraints: - The UI is strictly read-only. - Treat IDs and pagination cursors as opaque. - Never convert missing, failed, unknown, or inapplicable measurements to zero. - Preserve the server's pipeline and configuration-layer ordering. - Visually distinguish direct configuration changes from transitive invalidation. - Display why an artifact preview is unavailable; do not fetch artifacts through filesystem paths. - Allow navigation to both baseline and challenger episode pipelines. - Keep diagnostics and journal-integrity status visible. - Use a same-origin development proxy; do not request wildcard CORS. - Do not implement promotion, approval, chunk-lab, answer, calibration, or Pareto screens yet. Please make every screen deep-linkable and accessible by keyboard. Include loading, empty, diagnostic, and API-error states. Validate the implementation against the acceptance checklists in the handoff document. The backend can be started with: rag-ttc experiment optkit-rag campaign serve --store /tmp/rag-specialist --listen 127.0.0.1:8090 Use docmgr, commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill). Use a simple monochrome macos1 retro look, with modern fonts (no chicago), no window chrome and no menu bar, we want clean retro black and white UI."

**Assistant interpretation:** Act as the frontend engineer for the OPTKIT-007 handoff: build the five read-only specialist screens as a real SPA in the rag-ttc repo, test them against the archived fixtures first, then validate against the live server, while doing docmgr/diary bookkeeping and following the strict data-semantics constraints (opaque IDs, no zero-coercion, server ordering, preview reasons, both-arm navigation).

**Inferred user intent:** Get a working, contract-faithful specialist UI that a scientist can actually navigate, with the semantics of the sealed journal preserved (missing ≠ zero, direct ≠ transitive change), styled as a clean 1-bit retro instrument rather than a generic dashboard.

**Commit (code):** 09a911ae5 — "Specialist UI: scaffold read-only workflow app with tests"

### What I did

- Created OPTKIT-008 (`docmgr ticket create-ticket`) with an 11-item task list, since OPTKIT-007 is complete.
- Read the handoff document and all six fixtures; derived `src/api/types.ts` from the sketches plus the actual fixture shapes (e.g. `MetricComparison`, `ConfigLayerRef`, `ArtifactEdge` are fixture-derived, not in the sketches).
- Scaffolded `rag-ttc/apps/specialist/web` (pnpm, matching the sibling `apps/admin/web`): Vite 8, React 19, TypeScript strict + `noUncheckedIndexedAccess`, RTK Query, react-router-dom 7, vitest 4, MSW 2.
- Wired the same-origin dev proxy in `vite.config.ts` (`/api` → `http://127.0.0.1:8090`), no CORS anywhere.
- Built the design system in `src/styles/system.css`: two inks (#000/#fff), System-1 pinstripe panel headers with the title on a white plate, hard 3px offset shadows, checkerboard-dither material for missing values, diagonal hatch for transitive invalidation, dotted focus-visible outlines, reduced-motion-safe loading marquee.
- Implemented the five screens plus a campaign entry screen (the API has no campaign-list route, so `/` is a paste-the-ID form with a localStorage recent list): `CockpitScreen`, `ComparisonScreen` (embeds config diff, invalidation plan, and the paginated paired-case table), `PipelineScreen`, `ProvenanceScreen`.
- Encoded the handoff's semantics in small components: `Observation`/`OptionalNumber` (status word on dithered chip when not measured; zero renders as 0.000 only when measured), `ArtifactBlock` (preview or reason + metadata, never a filesystem fetch), `IdChip` (copyable opaque IDs), `DiagnosticsPanel` (visible even when empty), `ApiErrorState` (400/404/405/500 advice per the handoff).
- Pagination puts the opaque `after` cursor in the URL query (deep-linkable), echoes `next_cursor` verbatim, and only offers "First page"/"Next page" since cursors are forward-only.
- Wrote 28 tests across 6 files: observation semantics, cockpit checklist items, comparison layer-order preservation and direct/upstream distinction, both-pipeline links, cursor pagination against MSW, pipeline stage order and skipped-stage rendering, preview reasons, provenance copyability, and 400/404 error states.

### Why

- Fixtures-first (per the handoff) lets every semantic rule be locked in by tests before touching the live server.
- RTK Query keeps loading/error/caching uniform across five read-only GET endpoints and matches house webGuidelines; bootstrap styling was explicitly overridden by the user's monochrome macos1 direction.
- A separate `apps/specialist` app avoids entangling with the admin app's chat/protobuf stack.

### What worked

- Typecheck clean on first run; production build 338 kB JS / 6.9 kB CSS.
- The dither-as-missing-material idea unified the aesthetic with the core data rule: a missing measurement literally renders as 1-bit "undefined" texture.
- MSW's cursor handler doubles as a contract test: page 2 is only served when the archived cursor comes back byte-identical, and any other cursor returns the server's 400 shape.

### What didn't work

- First full test run: 20/28 failed, every screen showing `unknown_error: Request failed for an unknown reason`. Cause: `fetchBaseQuery({ baseUrl: "/api/rag/v1" })` — Node's undici fetch (used by vitest/jsdom) rejects relative URLs, unlike the browser. Fix: resolve the base against the origin, `baseUrl: new URL("/api/rag/v1", location.origin).toString()`, which is identical in the browser and lets MSW intercept in tests. All 28 tests passed afterwards.

### What I learned

- The comparison fixture's `metrics` entries carry `baseline`/`treatment` numbers not shown in the handoff's TS sketches — deriving types from fixtures (as the handoff instructs) caught this.
- `useParams` already URI-decodes segments, so opaque IDs with `:` round-trip cleanly as long as links are built with `encodeURIComponent` (centralized in `src/routes.ts`).

### What was tricky to build

- **Monochrome semantic encoding.** With only black and white, severity and change-kind can't lean on color. Solution: a material vocabulary — solid ink = direct change/error, diagonal hatch = transitive/derived (upstream invalidation, skipped stages), checkerboard dither = missing/unknown, dashed outline = unchanged/reuse. Each badge also carries its word, so nothing is pattern-only (a11y).
- **Forward-only cursor pagination that is still deep-linkable.** The cursor lives in `?after=`; back/forward and link-sharing work, but a stale cursor deep link yields the server's 400 — which `ApiErrorState` maps to "navigate again from a live screen" with a recovery link, per the handoff.

### What warrants a second pair of eyes

- `normalizeApiError` unwraps RTK Query's error union; if the server ever returns a non-JSON error body, fetchBaseQuery yields `PARSING_ERROR` and we fall back to a generic message — acceptable, but worth confirming against the real server's 404 behavior.
- The cockpit's arm-selection radios are per-row pairs (baseline/challenger columns); verify with a screen reader that the aria-labels ("Use limit-1 as baseline") read sensibly.

### What should be done in the future

- Validate against the live server and take screenshots (next step).
- Optional Storybook if the team wants a component catalog; tests + fixtures cover the contract for now.

### Code review instructions

- Start at `rag-ttc/apps/specialist/web/src/api/types.ts` and `specialistApi.ts` (contract), then `screens/ComparisonScreen.tsx` (densest semantics: diff order, direct-vs-upstream, cursor pagination), then `styles/system.css` (material vocabulary).
- Validate: `cd rag-ttc/apps/specialist/web && pnpm install && pnpm typecheck && pnpm test && pnpm build`.

### Technical details

- Routes: `/` (campaign entry), `/campaigns/:campaign`, `/campaigns/:campaign/compare/:baseline/:treatment` (+ `?after=` cursor), `/campaigns/:campaign/episodes/:episode/pipeline`, `/campaigns/:campaign/episodes/:episode/provenance`.
- Dev proxy: `vite.config.ts` proxies `/api` to `http://127.0.0.1:8090`; dev server on port 5197.
