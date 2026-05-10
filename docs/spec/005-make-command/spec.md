# Feature Spec: Make Command

## Goal

Spritey must render a recipe into a spritesheet PNG with an optional machine-readable report:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] --json
```

This is the first production rendering slice for Spritey CLI and must provide stable JSON envelopes for success and error outcomes.

## User Stories

- As an agent, I can run `spritey make` with a recipe file and compatible assets directory to produce a PNG at an explicit output path.
- As an agent, I can optionally request a report JSON file and rely on a stable report v1 schema for downstream automation.
- As an agent, I can distinguish misuse, validation failures, and render failures through stable exit codes and structured JSON error envelopes.

## Functional Requirements

- The command is `spritey make <recipe> --assets <dir> --out <png> [--report <json>] --json`.
- `<recipe>` is a required positional file path.
- `--assets <dir>` is required and points to a compatible assets directory.
- `--out <png>` is required and points to the destination spritesheet file path.
- `--report <json>` is optional. When provided, Spritey writes a report JSON artifact using report schema v1.
- `--json` is optional. When provided, Spritey writes stable machine-readable output to stdout.
- Spritey validates recipe input using existing recipe-validation contract before attempting render.
- Spritey validates assets input using existing assets-validation contract before attempting render.
- Spritey resolves required layer selections from recipe and pack metadata and renders one spritesheet PNG.
- On success, Spritey writes the PNG to `--out` and writes report JSON only when `--report` is present.
- Success JSON response is written to stdout with a stable envelope:
  - `ok`: `true`
  - `command`: `"make"`
  - `outputs`: object containing at least `png.path` and optional `report.path`
  - `summary`: object with at minimum `frame_count`, `canvas.width`, `canvas.height`, and `animation_count`
  - `warnings`: array (may be empty)
  - `errors`: empty array
- Failure JSON response is written to stdout with the same top-level envelope keys and:
  - `ok`: `false`
  - `errors`: non-empty structured error list
- JSON error objects use a stable shape:
  - `code`: stable machine code
  - `message`: human-readable error message
  - `field`: optional CLI argument or logical field
  - `details`: optional object with structured context

## Error Requirements

- CLI misuse returns exit code `2`.
- Input validation failures (recipe/assets invalid) return exit code `3`.
- Render failure returns exit code `6` and preserves existing convention.
- Missing `<recipe>` returns exit code `2` with error code `MISSING_RECIPE`.
- Missing `--assets` returns exit code `2` with error code `MISSING_ASSETS`.
- Missing `--out` returns exit code `2` with error code `MISSING_OUT`.
- Nonexistent recipe file returns exit code `3` with error code `RECIPE_FILE_NOT_FOUND`.
- Nonexistent assets directory returns exit code `3` and reuses assets validation error code conventions.
- Render pipeline failure (composition/write failure) returns exit code `6` with error code `RENDER_FAILED`.

## Report v1 Requirements

- Report v1 is minimal and deterministic for this slice.
- Required report fields:
  - `schema_version`: `"1"`
  - `command`: `"make"`
  - `recipe.path`
  - `assets.path`
  - `output.png.path`
  - `render.canvas.width`
  - `render.canvas.height`
  - `render.frame_count`
  - `render.animation_ids` (ordered list)
  - `layers.applied` (ordered list of applied layer ids)
  - `warnings` (array)
- Report content must be deterministic for identical inputs in the same runtime environment.

## Testing Requirements

- CLI command tests for success path with `--report` and without `--report`.
- CLI command tests for each required-argument misuse case (`<recipe>`, `--assets`, `--out`, `--json`).
- CLI command tests for recipe-not-found and assets-not-found behavior.
- CLI command tests that assert stable JSON envelope keys for success and error.
- CLI command tests that assert exit code `6` for render failure.
- Service tests for report v1 generation and deterministic field ordering.
- Tests must use runtime-generated deterministic fixture assets and recipe fixtures created within test setup or deterministic helpers under `testdata/`; no third-party assets.
- PNG output verification includes deterministic dimension checks and hash/pixel checks where deterministic.

## Out Of Scope

- Alternate output formats (GIF/APNG/webp).
- Advanced report analytics (timings, per-layer bounding boxes, palette diffs).
- Batch make workflows or multi-recipe orchestration.
- Asset pack authoring or third-party asset ingestion workflows.
- GUI/editor or plugin system expansion.
- Parallel rendering optimization and caching strategy.
