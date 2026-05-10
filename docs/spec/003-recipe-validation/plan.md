# Technical Plan: Recipe Validation

## Architecture

Follow the existing MVC-style CLI shape:

- `cmd/spritey` parses `validate <recipe-path> --assets <dir> --json`.
- `app/controllers` handles argument validation and exit-code mapping.
- `app/services` loads the catalog, reads the recipe file, and validates selected layers.
- `app/models` adds recipe and validation result structures.
- `app/views` writes stable validation JSON.

## Data Flow

1. CLI receives `validate recipes/knight.json --assets <dir> --json`.
2. Controller validates required recipe path and assets path.
3. Service loads the catalog through `CatalogLoader`.
4. Service reads and decodes the recipe.
5. Service applies pack default body type when missing.
6. Service validates selection IDs and body-type support.
7. View emits success or structured errors.

## Compatibility

This slice is metadata-only. It validates recipe/catalog compatibility, but does not require spritesheet PNGs or palette files.

## Validation

Use TDD for command and service behavior. Validate with:

```bash
go test ./...
make native-ci
go run ./cmd/spritey validate testdata/fixtures/recipes/valid-basic.json --assets testdata/fixtures/basic-assets --json
```
