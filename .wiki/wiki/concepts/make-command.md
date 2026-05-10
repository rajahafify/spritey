---
title: "Make Command"
category: concept
sources:
  - raw/notes/2026-05-10-render-input-readiness.md
  - raw/notes/2026-05-10-make-human-output.md
  - raw/notes/2026-05-10-make-command.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, make, render, report]
summary: "Spritey's make command supports deterministic JSON and non-JSON success output, plus render-input readiness gating with warning propagation."
---

# Make Command

`spritey make` is Spritey's first rendering slice:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]
```

Contract points in current implemented slices:

- recipe/assets validation is required before rendering;
- render-input readiness failures such as `MISSING_SPRITE_FRAME` are treated as input-validation failures (exit `3`), not render failures;
- output PNG path is explicit and required;
- report output is optional and uses minimal report schema v1;
- JSON envelope is stable in `--json` mode;
- readiness warnings are propagated into make JSON `warnings` and report `warnings`;
- non-JSON success output is deterministic and line-oriented for human operators;
- render failures use `RENDER_FAILED` with exit code `6`.

Current render pipeline is intentionally minimal and deterministic for testing: ordered layer composition using runtime fixture assets plus deterministic report field ordering.

Non-JSON success output format:

1. `ok: make`
2. `png: <path>`
3. `report: <path>` (when requested)
4. `frame_count: <n>`
5. `canvas: <width>x<height>`
6. `animation_count: <n>`
