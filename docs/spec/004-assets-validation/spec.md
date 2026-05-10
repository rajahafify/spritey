# Feature Spec: Assets Validation

## Goal

Spritey must validate the structural readiness of a compatible assets directory:

```bash
spritey assets validate --assets ./assets --json
```

This gives agents a fast preflight check before catalog inspection, recipe validation, or rendering.

## User Stories

- As an agent, I can validate an assets directory and receive stable JSON confirming the pack metadata and minimum asset-pack structure.
- As an agent, I can distinguish CLI misuse from invalid assets using structured errors and stable exit codes.
- As an agent, I can rely on this command to check required directories before deeper commands attempt catalog, recipe, or render work.

## Functional Requirements

- The command is `spritey assets validate --assets <dir> --json`.
- `assets` requires a subcommand. For this slice, the only supported assets subcommand is `validate`.
- `assets validate` requires `--assets <dir>`.
- Spritey validates that `<dir>` exists and is a directory.
- Spritey validates that `<dir>/pack.json` exists and parses as JSON.
- Spritey validates that `<dir>/sheet_definitions/` exists.
- Spritey recursively parses non-`meta_*.json` sheet definition files under `<dir>/sheet_definitions/`.
- Spritey validates that `<dir>/spritesheets/` exists.
- Spritey validates that `<dir>/palette_definitions/` exists.
- Success JSON contains top-level `ok`, `assets`, `pack`, `summary`, `warnings`, and `errors`.
- `assets.path` contains the input assets path as Spritey resolved it for validation output.
- `pack` uses the same pack metadata shape as the catalog command.
- `summary` contains at minimum:
  - `category_count`: number of discovered catalog categories.
  - `layer_count`: number of discovered catalog layers.
- `warnings` is an array and may be empty.
- `errors` is an array and is empty on success.
- Failure JSON uses the same top-level shape with `ok: false`.

## Error Requirements

- CLI misuse returns exit code `2`.
- Invalid assets returns exit code `3`.
- Missing assets subcommand returns exit code `2` and error code `MISSING_ASSETS_SUBCOMMAND`.
- Unsupported assets subcommand returns exit code `2` and error code `UNSUPPORTED_ASSETS_SUBCOMMAND`.
- Missing `--assets` returns exit code `2` and error code `MISSING_ASSETS`.
- A nonexistent assets directory returns exit code `3` and reuses error code `ASSETS_DIRECTORY_NOT_FOUND`.
- Missing `pack.json` returns exit code `3` and reuses error code `MISSING_PACK_JSON`.
- Invalid `pack.json` returns exit code `3` and reuses error code `INVALID_PACK_JSON`.
- Missing `sheet_definitions/` returns exit code `3` and reuses error code `MISSING_SHEET_DEFINITIONS`.
- Invalid sheet definition JSON returns exit code `3` and reuses error code `INVALID_SHEET_DEFINITION_JSON`.
- Missing `spritesheets/` returns exit code `3` and error code `MISSING_SPRITESHEETS`.
- Missing `palette_definitions/` returns exit code `3` and error code `MISSING_PALETTE_DEFINITIONS`.
- JSON mode writes structured responses to stdout.

## Out Of Scope

- Rendering.
- PNG or image-content checks.
- Palette content validation.
- Sprite dimensions.
- Recipe validation.
- `make`.
- Third-party sprite assets or fixture asset expansion beyond what tests require.
