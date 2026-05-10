---
title: "Make Command"
category: concept
sources:
  - raw/notes/2026-05-11-make-recolor-parity.md
  - raw/notes/2026-05-11-make-animation-layout-parity.md
  - raw/notes/2026-05-10-path-resolver-parity.md
  - raw/notes/2026-05-10-batch-make-manifest-v1.md
  - raw/notes/2026-05-10-make-report-artifact-integrity.md
  - raw/notes/2026-05-10-make-animation-strip-output.md
  - raw/notes/2026-05-10-make-report-provenance.md
  - raw/notes/2026-05-10-render-input-readiness.md
  - raw/notes/2026-05-10-make-human-output.md
  - raw/notes/2026-05-10-make-command.md
created: 2026-05-10
updated: 2026-05-11
tags: [spritey, make, render, report]
summary: "Spritey's make commands support Python-parity path resolution, palette recoloring, and LPC row/layout behavior with deterministic strip rendering, stable single/batch JSON and non-JSON output, readiness-gated validation, and additive report provenance plus output-artifact integrity metadata."
---

# Make Command

`spritey make` is Spritey's core rendering slice:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]
```

Contract points in current implemented slices:

- recipe/assets validation is required before rendering;
- frame lookup uses shared resolver parity order:
  1) `{prefix}/{anim}.png`,
  2) `{prefix}/fg/{anim}.png`,
  3) `{prefix}/{anim}/*.png`,
  4) mapped slash/thrust weapon dirs (`attack_slash`/`attack_thrust`) excluding lowercase `behind` filenames;
- slash/thrust structure-C paths win over mapped structure-D paths when both exist;
- per-layer recolor parity is applied before composition only when both `recolor_material` and selection `palette_variant` are present;
- recolor palette resolution scans `assets/palette_definitions/**/{material}_*.json` in lexicographic order with exact key first, then fuzzy case-insensitive lookup;
- recolor source/base palette fallback order is `light`, `base`, `beige`, `brown`;
- recolor pixel swap uses per-channel tolerance `<=2` for non-transparent pixels only (`alpha > 0`), and missing/invalid palettes are non-fatal no-op;
- render row emission follows fixed LPC order (`spellcast`, `thrust`, `walk`, `slash`, `shoot`, `hurt`) instead of pack default animation order;
- per LPC row, missing layer frame is skipped non-fatally and rows emit only when at least one layer contributes;
- first contributing layer sets row height; subsequent layers are padded/clipped to that row height before alpha compositing;
- strip width is fixed at `832` and emitted rows are stacked vertically;
- if no LPC rows emit, make writes a transparent `832x256` PNG;
- render-input readiness failures such as `MISSING_SPRITE_FRAME` are treated as input-validation failures (exit `3`), not render failures;
- output PNG path is explicit and required;
- report output is optional and uses additive report schema v1 provenance fields;
- when report output is requested, report v1 also includes output PNG artifact metadata (`sha256`, `bytes`) from the written file;
- output PNG is a vertical strip with one frame row per required animation ID;
- JSON envelope is stable in `--json` mode;
- readiness warnings are propagated into make JSON `warnings` and report `warnings`;
- non-JSON success output is deterministic and line-oriented for human operators;
- render failures use `RENDER_FAILED` with exit code `6`.

Batch make is available for manifest-driven runs:

```bash
spritey make batch <manifest.json> --assets <dir> [--json]
```

Batch behavior points:

- manifest v1 shape is `{"schema_version":"1","jobs":[...]}`;
- jobs execute sequentially in manifest order;
- relative recipe/output/report paths resolve from manifest file directory;
- existing `MakeService` is reused per job;
- execution fails fast on first failed job;
- JSON envelope remains stable with top-level keys:
  - `ok`, `command`, `summary`, `jobs`, `warnings`, `errors`;
- batch `command` value is `make-batch`;
- manifest file/json/schema/empty-jobs issues return exit `4`;
- CLI misuse returns exit `2`;
- job failures map with existing make exit behavior (`3`/`6`) and include failing-job context in JSON error details.

Report v1 now includes additive provenance and artifact integrity metadata:

- pack metadata;
- recipe requested/effective body-type metadata;
- deterministic composed-layer entries with category, z-order, resolved body-type/path, palette variant, and credits;
- `artifacts.output_png.sha256` and `artifacts.output_png.bytes` from the written output PNG.

Current render pipeline remains intentionally minimal and deterministic for testing: LPC-order rows are composed from available layer frames, then stacked into one deterministic strip with deterministic report field ordering.

Non-JSON success output format:

1. `ok: make`
2. `png: <path>`
3. `report: <path>` (when requested)
4. `frame_count: <n>`
5. `canvas: <width>x<height>`
6. `animation_count: <n>`
