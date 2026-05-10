---
title: "Make Command"
category: concept
sources:
  - raw/notes/2026-05-10-make-human-output.md
  - raw/notes/2026-05-10-make-command.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, make, render, report]
summary: "Spritey's make command now supports deterministic JSON and non-JSON success output with optional report v1."
---

# Make Command

`spritey make` is Spritey's first rendering slice:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]
```

Contract points in current implemented slices:

- recipe/assets validation is required before rendering;
- output PNG path is explicit and required;
- report output is optional and uses minimal report schema v1;
- JSON envelope is stable in `--json` mode;
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
