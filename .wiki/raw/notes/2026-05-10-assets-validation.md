---
title: "Assets Validation Implementation"
source: "D:/Work/spritey/docs/spec/004-assets-validation; D:/Work/spritey/app; D:/Work/spritey/cmd/spritey; D:/Work/spritey/testdata/fixtures/basic-assets"
type: notes
ingested: 2026-05-10
tags: [spritey, assets, validate, cli, implementation]
summary: "Fourth Spritey product slice implements `spritey assets validate --assets <dir> --json` for structural assets-pack preflight checks."
---

# Assets Validation Implementation

Spritey's fourth implemented product slice is:

```bash
spritey assets validate --assets ./assets --json
```

The feature lives under `docs/spec/004-assets-validation/` with `spec.md`, `plan.md`, and `tasks.md`.

## Behavior

- `assets` requires a subcommand.
- The only supported assets subcommand in this slice is `validate`.
- `assets validate` requires `--assets <dir>`.
- It validates the assets directory, `pack.json`, `sheet_definitions/`, `spritesheets/`, and `palette_definitions/`.
- It recursively parses non-`meta_*.json` sheet definition files by reusing catalog loading behavior.
- It returns top-level `ok`, `assets`, `pack`, `summary`, `warnings`, and `errors`.
- `summary` includes `category_count` and `layer_count`.

## Error Contract

- Missing assets subcommand returns exit code `2` and `MISSING_ASSETS_SUBCOMMAND`.
- Unsupported assets subcommand returns exit code `2` and `UNSUPPORTED_ASSETS_SUBCOMMAND`.
- Missing `--assets` returns exit code `2` and `MISSING_ASSETS`.
- Invalid assets structure returns exit code `3`.
- Reused invalid-assets codes include `ASSETS_DIRECTORY_NOT_FOUND`, `MISSING_PACK_JSON`, `INVALID_PACK_JSON`, `MISSING_SHEET_DEFINITIONS`, and `INVALID_SHEET_DEFINITION_JSON`.
- Assets-validation-specific invalid-assets codes include `MISSING_SPRITESHEETS` and `MISSING_PALETTE_DEFINITIONS`.

## Implementation Shape

- `cmd/spritey/main.go` routes `assets validate`.
- `app/controllers` owns assets subcommand orchestration and exit-code mapping.
- `app/services` validates assets-pack structure and computes summary counts.
- `app/models` contains assets-validation target, summary, and result structures.
- `app/views` writes assets-validation JSON responses.
