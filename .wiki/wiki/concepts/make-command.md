---
title: "Make Command"
category: concept
sources:
  - raw/notes/2026-05-10-make-command.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, make, render, report]
summary: "Spritey's spec-005 slice renders a recipe to PNG with optional report v1 and stable JSON envelopes."
---

# Make Command

`spritey make` is Spritey's first rendering slice:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]
```

Contract points in this slice:

- recipe/assets validation is required before rendering;
- output PNG path is explicit and required;
- report output is optional and uses minimal report schema v1;
- JSON envelope is stable in `--json` mode;
- render failures use `RENDER_FAILED` with exit code `6`.

Current render pipeline is intentionally minimal and deterministic for testing: ordered layer composition using runtime fixture assets plus deterministic report field ordering.
