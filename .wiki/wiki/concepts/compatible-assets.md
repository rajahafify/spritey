---
title: "Compatible Assets"
category: concept
sources:
  - raw/notes/2026-05-10-spritey-cli-contract.md
  - raw/repos/2026-05-10-python-reference-implementation.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, assets, pack-json, compatible-assets]
aliases: [Assets Pack, Spritey Assets]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "Spritey treats assets as a compatible pack format with pack.json defaults plus sheet definitions, spritesheets, and palette definitions; assets are user-provided at runtime."
---

# Compatible Assets

> Spritey's user-facing term is `assets`. Asset packs are not bundled in this repository; users download, install, or provide them at runtime.

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

## Asset Policy

Spritey must not commit third-party sprite assets. The repository contains the CLI, specs, schemas, tests, and fixture-sized test data only.

## See Also

- [[agent-friendly-cli|Agent-Friendly CLI]] ([Agent-Friendly CLI](agent-friendly-cli.md)) - commands that consume assets.

## Sources

- [Spritey CLI Contract](../../raw/notes/2026-05-10-spritey-cli-contract.md) - assets directory and `pack.json` defaults.
