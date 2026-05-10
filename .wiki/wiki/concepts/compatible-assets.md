---
title: "Compatible Assets"
category: concept
sources:
  - raw/notes/2026-05-10-spritey-cli-contract.md
  - raw/repos/2026-05-10-python-reference-implementation.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, assets, pack-json, lpc-compatible]
aliases: [Assets Pack, Spritey Assets]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "Spritey treats assets as a compatible pack format with pack.json defaults plus sheet definitions, spritesheets, and palette definitions."
---

# Compatible Assets

> Spritey's user-facing term is `assets`; internally the first compatible format follows the LPC-style folder and metadata shape used by the Python reference.

A compatible assets directory should contain:

```text
assets/
  pack.json
  sheet_definitions/
  spritesheets/
  palette_definitions/
```

`pack.json` owns pack-level defaults. This prevents Spritey from hardcoding behavior that belongs to the asset pack.

## Pack Defaults

Expected defaults include:

- body type
- animation list
- canvas width
- output format
- palette fallback rules
- missing body-type fallback behavior

## Reference Behavior

The Python reference currently assumes LPC-compatible `sheet_definitions/`, `spritesheets/`, and `palette_definitions/` folders. The Go implementation should preserve the proven behavior while making `pack.json` the formal owner of defaults.

## See Also

- [[python-reference|Python Reference]] ([Python Reference](../references/python-reference.md)) - behavior to preserve.
- [[agent-friendly-cli|Agent-Friendly CLI]] ([Agent-Friendly CLI](agent-friendly-cli.md)) - commands that consume assets.

## Sources

- [Spritey CLI Contract](../../raw/notes/2026-05-10-spritey-cli-contract.md) - assets directory and `pack.json` defaults.
- [Python Reference Implementation](../../raw/repos/2026-05-10-python-reference-implementation.md) - current asset assumptions.
