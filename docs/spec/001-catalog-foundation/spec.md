# Feature Spec: Catalog Foundation

## Goal

Spritey must expose the first useful agent workflow:

```bash
spritey catalog --assets ./assets --json
```

The command lets agents discover compatible asset categories and layer IDs before writing recipe files.

## User Stories

- As an agent, I can point Spritey at a compatible assets directory and receive stable JSON listing the pack and available layer categories.
- As an agent, I can tell when an assets directory is missing or malformed because Spritey returns structured errors and stable exit codes.

## Functional Requirements

- The command is `spritey catalog --assets <dir> --json`.
- `--assets` is required for `catalog`.
- Spritey loads `<dir>/pack.json`.
- Spritey recursively loads JSON files under `<dir>/sheet_definitions/`.
- Files whose names start with `meta_` are ignored as catalog metadata.
- Each sheet definition contributes one catalog layer using:
  - `id`: sheet definition file stem.
  - `name`: sheet `name` when present, otherwise `id`.
  - `category`: sheet `type_name` when present, otherwise a folder-derived category.
  - `z_pos`: `layer_1.zPos`, defaulting to `0`.
  - `body_types`: body-type keys from `layer_1`, excluding metadata keys.
  - `animations`: sheet `animations`, defaulting to an empty list.
  - `recolor_material`: `recolors.material` when present.
  - `path_prefix`: the first body-type path from `layer_1`, when present.
  - `credits`: sheet `credits`, defaulting to an empty list.
- JSON output contains top-level `ok`, `pack`, `categories`, `warnings`, and `errors`.
- Catalog categories and layers are sorted by ID for deterministic output.

## Error Requirements

- Missing `--assets` returns exit code `2` and a structured error with code `MISSING_ASSETS`.
- A nonexistent assets directory returns exit code `3` and a structured error with code `ASSETS_DIRECTORY_NOT_FOUND`.
- Missing `pack.json` returns exit code `3` and a structured error with code `MISSING_PACK_JSON`.
- Invalid sheet definition JSON returns exit code `3` and a structured error with code `INVALID_SHEET_DEFINITION_JSON`.
- JSON mode writes structured responses to stdout.

## Out Of Scope

- Recipe parsing.
- Rendering.
- `inspect layer`.
- Full JSON Schema validation.
- Loading or committing third-party sprite assets.
