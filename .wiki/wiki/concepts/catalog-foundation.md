---
title: "Catalog Foundation"
category: concept
sources:
  - raw/notes/2026-05-10-catalog-foundation.md
  - raw/notes/2026-05-10-spritey-cli-contract.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, catalog, cli, assets]
aliases: [Spritey Catalog Command, Catalog CLI]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "Spritey's first product slice implements `spritey catalog --assets <dir> --json` for deterministic asset discovery by agents."
---

# Catalog Foundation

> The catalog command is the first implemented Spritey workflow that agents can use against a compatible assets directory.

The implemented command is:

```bash
spritey catalog --assets ./assets --json
```

It reads `pack.json`, recursively loads `sheet_definitions/**/*.json`, ignores `meta_*.json`, and emits deterministic JSON with `ok`, `pack`, `categories`, `warnings`, and `errors`.

## JSON Contract

Catalog output groups layers by category. Each layer exposes the sheet-definition ID, display name, `z_pos`, supported body types, animations, optional recolor material, optional path prefix, and credits.

Errors use the same top-level response shape. The first implemented structured error codes are:

- `MISSING_ASSETS` with exit code `2`.
- `ASSETS_DIRECTORY_NOT_FOUND` with exit code `3`.
- `MISSING_PACK_JSON` with exit code `3`.
- `INVALID_SHEET_DEFINITION_JSON` with exit code `3`.

## Architecture Notes

The implementation follows the Rails-style CLI MVC split:

- `cmd/spritey` keeps parsing thin.
- `app/controllers` owns orchestration and exit codes.
- `app/services` loads pack and sheet metadata.
- `app/models` owns catalog data structures.
- `app/views` writes JSON.

## See Also

- [[agent-friendly-cli|Agent-Friendly CLI]] ([Agent-Friendly CLI](agent-friendly-cli.md)) - larger CLI contract.
- [[compatible-assets|Compatible Assets]] ([Compatible Assets](compatible-assets.md)) - assets directory shape.

## Sources

- [Catalog Foundation Implementation](../../raw/notes/2026-05-10-catalog-foundation.md) - implemented behavior and tests.
- [Spritey CLI Contract](../../raw/notes/2026-05-10-spritey-cli-contract.md) - intended command family.
