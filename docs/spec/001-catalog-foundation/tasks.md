# Tasks: Catalog Foundation

## Spec And Fixtures

- [x] Create feature spec, technical plan, and tasks under `docs/spec/001-catalog-foundation/`.
- [x] Add a minimal compatible assets fixture under `testdata/fixtures/basic-assets/`.

## Tests

- [x] Add CLI test for successful `catalog --assets testdata/fixtures/basic-assets --json`.
- [x] Add CLI test for missing `--assets`.
- [x] Add CLI test for nonexistent assets directory.
- [x] Add CLI test for missing `pack.json`.
- [x] Add CLI test for invalid sheet definition JSON.
- [x] Add service/model tests for pack loading and catalog grouping.

## Implementation

- [x] Add model types for pack metadata, catalog categories, layers, and errors.
- [x] Add catalog loader service for `pack.json` and recursive sheet definitions.
- [x] Add catalog controller for argument validation, orchestration, and exit codes.
- [x] Add JSON view for stable output.
- [x] Replace scaffold output with real CLI routing.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.
- [x] Update `.wiki/` with the new catalog foundation decision and CLI contract.
