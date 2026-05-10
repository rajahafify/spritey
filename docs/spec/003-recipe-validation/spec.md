# Feature Spec: Recipe Validation

## Goal

Spritey must validate a file-based recipe against a compatible assets directory:

```bash
spritey validate recipe.json --assets ./assets --json
```

This gives agents a validation gate before rendering.

## User Stories

- As an agent, I can validate a recipe file and receive stable JSON before calling `make`.
- As an agent, I can tell whether validation failed because of malformed JSON, missing selections, unknown layer IDs, or unsupported body types.

## Functional Requirements

- The command is `spritey validate <recipe-path> --assets <dir> --json`.
- `validate` requires a recipe path.
- `--assets` is required.
- Spritey loads the recipe JSON from disk.
- If `body_type` is omitted, Spritey uses `assets/pack.json` default `body_type`.
- The recipe must contain at least one selection.
- Each selection must contain a non-empty `id`.
- Each selected layer ID must exist in the loaded catalog.
- Each selected layer must support the effective body type.
- Success JSON contains top-level `ok`, `recipe`, `warnings`, and `errors`.
- Failure JSON uses the same top-level shape with `ok: false`.

## Error Requirements

- Missing recipe path returns exit code `2` and error code `MISSING_RECIPE`.
- Missing `--assets` returns exit code `2` and error code `MISSING_ASSETS`.
- Missing recipe file returns exit code `4` and error code `RECIPE_FILE_NOT_FOUND`.
- Invalid recipe JSON returns exit code `4` and error code `INVALID_RECIPE_JSON`.
- Missing selections returns exit code `4` and error code `MISSING_SELECTIONS`.
- Missing selection ID returns exit code `4` and error code `MISSING_SELECTION_ID`.
- Invalid assets errors reuse exit code `3`.
- Validation failures such as unknown layer ID or unsupported body type return exit code `5`.

## Out Of Scope

- Rendering.
- Sprite PNG existence checks.
- Palette-file resolution.
- JSON Schema files.
- Inline recipe JSON.
