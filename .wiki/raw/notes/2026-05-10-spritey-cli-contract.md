---
title: "Spritey CLI Contract"
source: "D:/Work/spritey/README.md; D:/Work/spritey/docs/spec/constitution.md"
type: notes
ingested: 2026-05-10
tags: [spritey, cli, contract, assets, recipes, agents]
summary: "Spritey's intended complete CLI exposes catalog, inspect, validate, make, and assets validate commands with file-based recipes, explicit outputs, JSON modes, and stable exit codes."
---

# Spritey CLI Contract

Spritey is designed for agents and automation. Inputs are file-based, outputs are explicit, validation is available before generation, and JSON output is stable enough for scripts to consume.

## Expected Commands

```bash
spritey catalog --assets ./assets --json
spritey inspect layer torso_armour_plate --assets ./assets --json
spritey validate recipes/knight.json --assets ./assets --json
spritey make recipes/knight.json --assets ./assets --out output/knight.png --report output/knight.report.json
spritey assets validate --assets ./assets --json
```

## Assets Directory

```text
assets/
  pack.json
  sheet_definitions/
  spritesheets/
  palette_definitions/
```

`pack.json` owns defaults such as body type, animation list, canvas width, palette fallback rules, output format, and missing body-type fallback behavior.

## Recipe Shape

Recipes are JSON files that describe selected layers and optional palette variants. Inline JSON is not the primary workflow.

## Reporting

`spritey make --report` writes machine-readable JSON with output path, size, pack metadata, selected layers, palette behavior, and warnings.

## Exit Codes

The README defines stable exit codes:

- `0`: success
- `1`: general error
- `2`: invalid CLI usage
- `3`: invalid assets directory
- `4`: invalid recipe
- `5`: validation failed
- `6`: render failed
