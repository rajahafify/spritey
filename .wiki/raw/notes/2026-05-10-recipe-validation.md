---
title: "Recipe Validation Implementation"
source: "D:/Work/spritey/docs/spec/003-recipe-validation; D:/Work/spritey/app; D:/Work/spritey/cmd/spritey; D:/Work/spritey/testdata/fixtures/recipes"
type: notes
ingested: 2026-05-10
tags: [spritey, validate, recipes, cli, implementation]
summary: "Third Spritey product slice implements `spritey validate <recipe> --assets <dir> --json` for metadata-only recipe/catalog validation."
---

# Recipe Validation Implementation

Spritey's third implemented product slice is:

```bash
spritey validate <recipe-path> --assets ./assets --json
```

The feature lives under `docs/spec/003-recipe-validation/` with `spec.md`, `plan.md`, and `tasks.md`.

## Behavior

- `validate` requires a file-based recipe path.
- It reuses the catalog loader.
- It applies `pack.json` default `body_type` when the recipe omits `body_type`.
- It requires at least one selection.
- It requires every selection to contain an `id`.
- It verifies selected layer IDs exist.
- It verifies each selected layer supports the effective body type.
- It returns top-level `ok`, `recipe`, `warnings`, and `errors`.

## Error Contract

- Missing recipe path returns exit code `2` and error code `MISSING_RECIPE`.
- Missing `--assets` returns exit code `2` and error code `MISSING_ASSETS`.
- Missing recipe file returns exit code `4` and error code `RECIPE_FILE_NOT_FOUND`.
- Invalid recipe JSON returns exit code `4` and error code `INVALID_RECIPE_JSON`.
- Missing selections returns exit code `4` and error code `MISSING_SELECTIONS`.
- Missing selection ID returns exit code `4` and error code `MISSING_SELECTION_ID`.
- Invalid assets errors reuse catalog loader exit code `3`.
- Unknown layer ID and unsupported body type return exit code `5`.

## Implementation Shape

- `cmd/spritey/main.go` routes `validate`.
- `app/controllers` owns validation orchestration and exit-code mapping.
- `app/services` reads recipe files and validates them against catalog metadata.
- `app/models` contains recipe and validation result structures.
- `app/views` writes validation JSON responses.
