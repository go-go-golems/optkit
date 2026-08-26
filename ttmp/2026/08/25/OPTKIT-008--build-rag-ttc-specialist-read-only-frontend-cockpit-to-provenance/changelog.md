# Changelog

## 2026-08-25

- Initial workspace created


## 2026-08-25

Step 1: scaffolded apps/specialist/web with all five read-only screens, monochrome macos1 design system, and 28 passing fixture-driven MSW tests (commit 09a911ae5)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/apps/specialist/web/src/api/types.ts — v1 API types derived from handoff sketches and archived fixtures


## 2026-08-25

Step 2: validated against live server (campaign serve on 8091 behind SPECIALIST_API proxy override), walked all screens with Playwright, archived screenshots, fixed row-header uppercasing of opaque IDs, disabled-button rendering, hatch density (commits 316800111, 27aa79671)

### Related Files

- /home/manuel/workspaces/2026-08-24/use-optkit/rag-ttc/apps/specialist/web/vite.config.ts — SPECIALIST_API proxy override for occupied port 8090


## 2026-08-25

Specialist read-only frontend delivered and validated against live server (rag-ttc commits 09a911ae5, 316800111, 27aa79671)

