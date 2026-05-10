# Tasks: Assets Validation

## Spec

- [x] Create feature spec, technical plan, and tasks under `docs/spec/004-assets-validation/`.

## Tests

- [x] Add CLI test for successful `assets validate --assets testdata/fixtures/basic-assets --json`.
- [x] Add CLI test for missing assets subcommand after `assets`.
- [x] Add CLI test for unsupported assets subcommand.
- [x] Add CLI test for missing `--assets`.
- [x] Add CLI test for nonexistent assets directory.
- [x] Add CLI test for missing `pack.json`.
- [x] Add CLI test for invalid `pack.json`.
- [x] Add CLI test for missing `sheet_definitions/`.
- [x] Add CLI test for invalid sheet definition JSON.
- [x] Add CLI test for missing `spritesheets/`.
- [x] Add CLI test for missing `palette_definitions/`.
- [x] Add service tests for assets validation summary counts.

## Implementation

- [x] Add assets validation model/result types for `assets.path`, `pack`, `summary`, `warnings`, and `errors`.
- [x] Add assets validation service that reuses existing pack and sheet-definition validation where possible.
- [x] Add controller handling for assets validation argument checks and exit-code mapping.
- [x] Add JSON view for assets validation success and error responses.
- [x] Extend CLI routing for `assets validate`.
- [x] Add stable structured errors for `MISSING_ASSETS_SUBCOMMAND`, `UNSUPPORTED_ASSETS_SUBCOMMAND`, `MISSING_SPRITESHEETS`, and `MISSING_PALETTE_DEFINITIONS`.
- [x] Reuse `MISSING_ASSETS`, `ASSETS_DIRECTORY_NOT_FOUND`, `MISSING_PACK_JSON`, `INVALID_PACK_JSON`, `MISSING_SHEET_DEFINITIONS`, and `INVALID_SHEET_DEFINITION_JSON`.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.
- [x] Run direct assets-validation smoke test.
