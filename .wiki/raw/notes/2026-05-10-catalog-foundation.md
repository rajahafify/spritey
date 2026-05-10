---
title: "Catalog Foundation Implementation"
source: "D:/Work/spritey/docs/spec/001-catalog-foundation; D:/Work/spritey/app; D:/Work/spritey/cmd/spritey; D:/Work/spritey/testdata/fixtures/basic-assets"
type: notes
ingested: 2026-05-10
tags: [spritey, catalog, cli, assets, implementation]
summary: "First Spritey product slice implements `spritey catalog --assets <dir> --json` with MVC-style app packages, structured JSON output, stable exit codes, and a fixture asset pack."
---

# Catalog Foundation Implementation

Spritey's first implemented product slice is the catalog command:

```bash
spritey catalog --assets ./assets --json
```

The feature lives under `docs/spec/001-catalog-foundation/` with `spec.md`, `plan.md`, and `tasks.md`.

## Behavior

- `catalog --assets <dir> --json` reads `<dir>/pack.json`.
- It recursively reads `<dir>/sheet_definitions/**/*.json`.
- It ignores `meta_*.json`.
- It groups layers by `type_name`, falling back to the first sheet-definition directory segment.
- It sorts categories and layers by ID for deterministic output.
- It returns top-level `ok`, `pack`, `categories`, `warnings`, and `errors`.

## Error Contract

- Missing `--assets` returns exit code `2` and error code `MISSING_ASSETS`.
- A nonexistent assets directory returns exit code `3` and error code `ASSETS_DIRECTORY_NOT_FOUND`.
- Missing `pack.json` returns exit code `3` and error code `MISSING_PACK_JSON`.
- Invalid sheet definition JSON returns exit code `3` and error code `INVALID_SHEET_DEFINITION_JSON`.

## Implementation Shape

- `cmd/spritey/main.go` parses the command and flags, then delegates.
- `app/controllers` owns command orchestration and exit-code decisions.
- `app/services` loads asset metadata and builds the catalog.
- `app/models` owns pack, catalog, layer, credit, and problem structs.
- `app/views` writes stable JSON responses.

## Test Fixture

`testdata/fixtures/basic-assets/` is a tiny compatible asset pack containing only authored metadata and empty directories. It does not contain third-party sprite assets.
