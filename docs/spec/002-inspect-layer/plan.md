# Technical Plan: Inspect Layer

## Architecture

Reuse the catalog foundation MVC shape:

- `cmd/spritey` parses `inspect layer <layer-id> --assets <dir> --json`.
- `app/controllers` owns inspect command orchestration and exit-code decisions.
- `app/services` loads the catalog and performs exact layer lookup.
- `app/models` adds an inspect-layer result with the category attached.
- `app/views` writes stable inspect JSON responses.

## Data Flow

1. CLI receives `inspect layer body_human --assets <dir> --json`.
2. The inspect controller validates target, layer ID, and assets path.
3. The inspect service calls the existing catalog loader.
4. The inspect service scans sorted categories and sorted layers for an exact layer ID match.
5. The view emits one layer response or a structured error.

## Compatibility

This is a metadata-only read path. It must not require sprite PNG files and must not change existing `catalog` output.

## Validation

Use TDD for CLI and lookup behavior. Validate with:

```bash
go test ./...
make native-ci
go run ./cmd/spritey inspect layer body_human --assets testdata/fixtures/basic-assets --json
```
