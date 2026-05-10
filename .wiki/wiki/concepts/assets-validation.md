---
title: "Assets Validation"
category: concept
sources:
  - raw/notes/2026-05-10-assets-validation.md
  - raw/notes/2026-05-10-spritey-cli-contract.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, assets, validate, cli]
aliases: [Spritey Assets Validate Command, Assets Preflight]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "Spritey's assets validate command checks the minimum compatible assets-pack structure before catalog, recipe, or rendering work."
---

# Assets Validation

> The assets validate command gives agents a structural preflight for compatible assets directories.

The implemented command is:

```bash
spritey assets validate --assets ./assets --json
```

It checks that the assets directory exists, loads `pack.json`, parses non-meta sheet definition JSON, and requires `sheet_definitions/`, `spritesheets/`, and `palette_definitions/`.

## JSON Contract

Success output includes `ok`, `assets`, `pack`, `summary`, `warnings`, and `errors`. The `assets` payload records the input path, and `summary` reports `category_count` and `layer_count`.

Structured CLI errors include `MISSING_ASSETS_SUBCOMMAND`, `UNSUPPORTED_ASSETS_SUBCOMMAND`, and `MISSING_ASSETS`. Structural assets errors include reused catalog codes plus `MISSING_SPRITESHEETS` and `MISSING_PALETTE_DEFINITIONS`.

## See Also

- [[compatible-assets|Compatible Assets]] ([Compatible Assets](compatible-assets.md)) - minimum assets directory contract.
- [[catalog-foundation|Catalog Foundation]] ([Catalog Foundation](catalog-foundation.md)) - metadata loader reused by assets validation.
- [[agent-friendly-cli|Agent-Friendly CLI]] ([Agent-Friendly CLI](agent-friendly-cli.md)) - larger CLI workflow.

## Sources

- [Assets Validation Implementation](../../raw/notes/2026-05-10-assets-validation.md) - implemented behavior and tests.
- [Spritey CLI Contract](../../raw/notes/2026-05-10-spritey-cli-contract.md) - intended command family.
