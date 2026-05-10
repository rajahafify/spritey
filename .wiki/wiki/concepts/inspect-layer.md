---
title: "Inspect Layer"
category: concept
sources:
  - raw/notes/2026-05-10-inspect-layer.md
  - raw/notes/2026-05-10-catalog-foundation.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, inspect, cli, assets]
aliases: [Spritey Inspect Layer, Inspect Layer CLI]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "Spritey's inspect-layer command returns one catalog layer's metadata for agent recipe authoring."
---

# Inspect Layer

> The inspect-layer command lets agents inspect a known catalog layer before writing recipes.

The implemented command is:

```bash
spritey inspect layer <layer-id> --assets ./assets --json
```

It reuses catalog loading, finds one exact layer ID, and emits deterministic JSON with `ok`, `layer`, `warnings`, and `errors`.

## JSON Contract

The layer payload includes:

- `category`
- `id`
- `name`
- `z_pos`
- `body_types`
- `animations`
- `recolor_material`
- `path_prefix`
- `credits`

The first inspect-specific structured errors are `MISSING_INSPECT_TARGET`, `UNSUPPORTED_INSPECT_TARGET`, `MISSING_LAYER_ID`, and `UNKNOWN_LAYER_ID`.

## See Also

- [[catalog-foundation|Catalog Foundation]] ([Catalog Foundation](catalog-foundation.md)) - metadata loader reused by inspect.
- [[agent-friendly-cli|Agent-Friendly CLI]] ([Agent-Friendly CLI](agent-friendly-cli.md)) - larger CLI workflow.

## Sources

- [Inspect Layer Implementation](../../raw/notes/2026-05-10-inspect-layer.md) - implemented behavior and tests.
- [Catalog Foundation Implementation](../../raw/notes/2026-05-10-catalog-foundation.md) - loader and first CLI slice.
