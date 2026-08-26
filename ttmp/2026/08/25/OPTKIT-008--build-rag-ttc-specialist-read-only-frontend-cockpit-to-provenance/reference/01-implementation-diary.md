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
    - Path: ws://rag-ttc/apps/specialist/README.md
      Note: Run instructions and semantics contract for the app (commit 27aa79671)
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

## Step 2: Live-server validation, screenshot walkthrough, and fidelity fixes

Built the rag-ttc binary, ran the deterministic campaign into `/tmp/rag-specialist`, served the specialist API, and drove the real UI with Playwright through the full workflow: entry → cockpit → arm selection → comparison → baseline episode pipeline → provenance → bad-deep-link error state → keyboard focus check. Fixed three visual/fidelity defects the screenshots exposed, re-ran the suite (28/28), and archived the walkthrough screenshots in the ticket.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Complete the "validate the implementation against the acceptance checklists" part of the brief against the live server, not just fixtures.

**Inferred user intent:** Confidence that the app works end-to-end against the real read-only API, with evidence.

**Commit (code):** 316800111 — "Specialist UI: live-validation fixes"; 27aa79671 — "Specialist UI: add app README with run and semantics notes"

### What I did

- Built rag-ttc (`GOWORK=off go build`), ran `campaign run` (campaign `campaign:d0e2d0dade11386980811d9013bbf1ea`, journal seq 47), started `campaign serve` and `pnpm dev` in tmux session `specialist`.
- Port 8090 was already occupied by an unrelated local app (`pbui-chat`), which I did not kill; instead the API serves on 8091 and `vite.config.ts` now honours a `SPECIALIST_API` env override for the proxy target (default stays 8090 per the handoff).
- Walked every screen with Playwright, took full-page screenshots, and archived them under `various/screenshots/` in this ticket (01-cockpit … 06-focus).
- Fixed defects found in the screenshots (see below), added an inline-SVG monochrome favicon, wrote `apps/specialist/README.md`.
- Checked all acceptance checklists from the handoff (results in Technical details).

### Why

- The fixtures freeze one journal point; only the live server proves the proxy, URL encoding of `:`-bearing IDs, and the real error bodies.

### What worked

- The full workflow works end-to-end on first run: deep links with encoded IDs, both-arm pipeline links, provenance custody table (17 edges), and the 404 screen showing the server's real body (`cockpit_unavailable: campaign campaign:doesnotexist not found`).
- Keyboard-only operation: Tab reaches every control, dotted System-1 focus outlines are visible, radios + submit work without a pointer.

### What didn't work

- `tmux new-session` for the API server died silently (window vanished): `Error: serve specialist API: listen tcp 127.0.0.1:8090: bind: address already in use` — an unrelated `pbui-chat` process owns 8090. Resolved with the `SPECIALIST_API` proxy override on 8091 rather than killing someone else's process.
- The disabled "Compare arms" button rendered garbled: the dither-as-text-color trick (`background-clip: text` over transparent text) breaks with the button's own background. Replaced with a dashed border + #878787 ink, which reads as dithered at UI sizes.
- `table.grid th { text-transform: uppercase }` applied to row headers too, so opaque IDs displayed as `Q-HYBRID` — a data-fidelity bug (IDs' casing is data). Scoped the label treatment to `thead th`.

### What I learned

- Campaign IDs are per-invocation (fixture: `campaign:53e7…`, live run: `campaign:d0e2…`) while the trial ID is deterministic (`trial:505a28c9…` matches the fixture) — reinforcing that campaign IDs must come from the operator, never be assumed.
- The live `/health` route responds at `/api/rag/v1/health` with `read_only: true`; probing any other path 404s, which initially masqueraded as "server down" when the wrong process owned the port.

### What was tricky to build

- Diagnosing the 404s: the health endpoint returned `404 page not found`, which looked like a routing bug in the new serve command. The actual cause was a different process on 8090 answering all paths with its own 404. Symptom that cracked it: the tmux window for the serve command had disappeared, and re-running serve in the foreground printed the bind error. Lesson: when a fresh server 404s everything, check who owns the port before reading route code.

### What warrants a second pair of eyes

- The `SPECIALIST_API` override reads `process.env` in `vite.config.ts` — dev-only, but confirm it should not instead become a documented `.env` convention for the team.
- Screenshots show the hatch at 28% black; verify on a real display that transitive-change badges stay clearly distinct from unchanged (dashed) at 100% zoom.

### What should be done in the future

- N/A — deferred screens (chunk lab, answer studio, calibration, Pareto, promotion) are explicitly out of scope per the handoff.

### Code review instructions

- Diff of this step: `git show 316800111` in rag-ttc (CSS scoping, disabled button, hatch density, favicon, proxy override).
- Reproduce validation: build rag-ttc, `campaign run` + `campaign serve --listen 127.0.0.1:8091`, then `SPECIALIST_API=http://127.0.0.1:8091 pnpm dev` and open `/`, paste the campaign ID from the run output.
- Evidence: `various/screenshots/01-cockpit.png` … `06-focus.png` in this ticket.

### Technical details

Acceptance checklist results (handoff §Screen acceptance checklist):

- Cockpit: integrity badge visible (verified/not-verified inverted on failure); queued/active/failed_terminal counts always rendered, failed count inverted when non-zero; budget violation renders an inverted "budget violated" badge; arm means and estimates use `OptionalNumber` (missing → dithered "unknown", zero → 0.000 only when present); baseline/challenger chosen by explicit per-row radios. ✓
- Comparison: direct = solid ink badge, transitive = hatched badge, unchanged = dashed; metric pair count in table; each case row links to both baseline and challenger pipeline + provenance; non-measured statuses render as words on dither. ✓
- Pipeline: stages in server `seq` order (test-asserted against fixture order); in→out counts with delta and chunk-ID chips; `preview_reason` mapped to sentences with metadata retained; diagnostics panel always present. ✓
- Provenance: manifest/graph/snapshot/episode/trial/result/trajectory/output all copyable `IdChip`s; sensitivity + size on every artifact and custody edge; no filesystem paths anywhere; explicit footnote that digests are not cryptographic signatures. ✓
