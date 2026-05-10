# Technical Plan: Assets Validation

## Architecture

Follow the existing MVC-style CLI shape and keep parsing thin:

- `cmd/spritey` parses `assets validate --assets <dir> --json`.
- `app/controllers` owns assets subcommand orchestration, argument validation, and exit-code decisions.
- `app/services` validates the compatible assets directory structure and reuses the existing catalog-loading behavior where possible.
- `app/models` adds an assets-validation result and summary structure if existing catalog models are not sufficient.
- `app/views` writes stable assets-validation JSON responses.

## Data Flow

1. CLI receives `assets validate --assets <dir> --json`.
2. CLI routing validates that `assets` has the `validate` subcommand.
3. The assets validation controller validates the required assets path.
4. The service checks the assets directory, `pack.json`, `sheet_definitions/`, `spritesheets/`, and `palette_definitions/`.
5. The service parses pack JSON and sheet definition JSON, reusing existing catalog loader errors where possible.
6. The service computes `summary.category_count` and `summary.layer_count` from discovered catalog metadata.
7. The view emits success JSON or a structured error response.

## JSON Contract

Success response shape:

```json
{
  "ok": true,
  "assets": {
    "path": "./assets"
  },
  "pack": {
    "schema_version": "1",
    "id": "example",
    "name": "Example Assets"
  },
  "summary": {
    "category_count": 1,
    "layer_count": 3
  },
  "warnings": [],
  "errors": []
}
```

Failure responses keep the same top-level keys, set `ok` to `false`, and include one structured error in `errors`.

## Compatibility

This slice validates the minimum compatible assets directory shape from the constitution:

```text
assets/
  pack.json
  sheet_definitions/
  spritesheets/
  palette_definitions/
```

The validation is metadata and directory-structure focused. It must not require sprite PNG decoding, palette semantic validation, sprite dimensions, recipe files, or render output.

## Validation

Use TDD for all CLI and validation behavior when implementation begins. Add command tests before implementation, then service tests for structure validation and summary counts. Validate with:

```bash
go test ./...
make native-ci
go run ./cmd/spritey assets validate --assets testdata/fixtures/basic-assets --json
```
