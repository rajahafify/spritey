# Feature Spec: Inspect Layer

## Goal

Spritey must let agents inspect one catalog layer by ID:

```bash
spritey inspect layer <layer-id> --assets ./assets --json
```

This command returns the selected layer metadata needed before recipe authoring.

## User Stories

- As an agent, I can inspect a layer ID discovered from `catalog` and receive stable JSON with body types, animations, paths, recolor material, and credits.
- As an agent, I can distinguish missing arguments, unsupported inspect targets, invalid assets, and unknown layer IDs through structured errors and stable exit codes.

## Functional Requirements

- The command is `spritey inspect layer <layer-id> --assets <dir> --json`.
- `inspect` requires a target. For this slice, the only supported target is `layer`.
- `inspect layer` requires a layer ID.
- `--assets` is required.
- Spritey loads the existing catalog from the assets directory and finds a layer by exact `id`.
- The success response contains top-level `ok`, `layer`, `warnings`, and `errors`.
- The returned layer includes `category`, `id`, `name`, `z_pos`, `body_types`, `animations`, `recolor_material`, `path_prefix`, and `credits`.

## Error Requirements

- Missing inspect target returns exit code `2` and error code `MISSING_INSPECT_TARGET`.
- Unsupported inspect target returns exit code `2` and error code `UNSUPPORTED_INSPECT_TARGET`.
- Missing layer ID returns exit code `2` and error code `MISSING_LAYER_ID`.
- Missing `--assets` returns exit code `2` and error code `MISSING_ASSETS`.
- Invalid assets directory errors reuse the catalog loader's exit code `3` behavior.
- Unknown layer ID returns exit code `5` and error code `UNKNOWN_LAYER_ID`.

## Out Of Scope

- Recipe validation.
- Rendering.
- Fuzzy layer search.
- Duplicate layer ID detection.
- Inspecting packs, categories, recipes, or palettes.
