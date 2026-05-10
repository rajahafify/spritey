---
title: "Recipe Validation"
category: concept
sources:
  - raw/notes/2026-05-10-render-input-readiness.md
  - raw/notes/2026-05-10-recipe-validation.md
  - raw/notes/2026-05-10-spritey-cli-contract.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, validate, recipes, cli]
aliases: [Spritey Validate Command, Recipe Validation CLI]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "Spritey's validate command checks file-based recipes against catalog metadata and render-input readiness before rendering."
---

# Recipe Validation

> The validate command gives agents a deterministic gate before rendering, including render-input readiness.

The implemented command is:

```bash
spritey validate <recipe-path> --assets ./assets --json
```

It loads a recipe file, applies the pack default body type when needed, checks selections, verifies selected layer IDs, validates body-type support with optional per-layer fallback path behavior, and verifies required source frames exist for required animations.

## JSON Contract

Success output includes `ok`, `recipe`, `warnings`, and `errors`. The recipe payload contains the recipe path, effective body type, and sorted validated selections.

Validation-specific structured errors include `MISSING_RECIPE`, `RECIPE_FILE_NOT_FOUND`, `INVALID_RECIPE_JSON`, `MISSING_SELECTIONS`, `MISSING_SELECTION_ID`, `UNKNOWN_LAYER_ID`, `UNSUPPORTED_BODY_TYPE`, and `MISSING_SPRITE_FRAME`.

Warnings are emitted when a selected layer does not provide the requested body-type path and Spritey uses `pack.defaults.missing_body_type_fallback` for that layer.

## See Also

- [[catalog-foundation|Catalog Foundation]] ([Catalog Foundation](catalog-foundation.md)) - metadata loader reused by validation.
- [[inspect-layer|Inspect Layer]] ([Inspect Layer](inspect-layer.md)) - command agents use before recipe authoring.
- [[agent-friendly-cli|Agent-Friendly CLI]] ([Agent-Friendly CLI](agent-friendly-cli.md)) - larger CLI workflow.

## Sources

- [Recipe Validation Implementation](../../raw/notes/2026-05-10-recipe-validation.md) - implemented behavior and tests.
- [Render-Input Readiness Implementation](../../raw/notes/2026-05-10-render-input-readiness.md) - readiness checks, missing-frame failures, and fallback warnings.
- [Spritey CLI Contract](../../raw/notes/2026-05-10-spritey-cli-contract.md) - intended command family.
