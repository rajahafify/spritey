# Technical Plan: Catalog Foundation

## Architecture

Keep CLI parsing thin in `cmd/spritey`. The command entrypoint parses arguments, then delegates catalog behavior to app packages:

- `app/models`: pack, layer, category, catalog, and structured error data.
- `app/services`: asset-pack catalog loading from `pack.json` and `sheet_definitions`.
- `app/controllers`: command orchestration and exit-code decisions.
- `app/views`: JSON response formatting.

## Data Flow

1. `cmd/spritey` receives `catalog --assets <dir> --json`.
2. The catalog controller validates required arguments.
3. The catalog service reads `pack.json`.
4. The catalog service walks `sheet_definitions`, ignores `meta_*.json`, parses layer definitions, and groups them by category.
5. The view writes deterministic JSON to stdout.

## Compatibility

The first slice supports the compatible asset-pack shape from the constitution:

```text
assets/
  pack.json
  sheet_definitions/
  spritesheets/
  palette_definitions/
```

Only JSON metadata is required for the catalog command. Sprite image files are not required until rendering.

## Validation

Use TDD for behavior that affects the CLI contract. Add CLI tests first, then app-level loader tests. Validate with:

```bash
go test ./...
make native-ci
```
