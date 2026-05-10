# Tasks: Recipe Validation

## Spec And Fixtures

- [x] Create feature spec, technical plan, and tasks under `docs/spec/003-recipe-validation/`.
- [x] Add recipe fixtures under `testdata/fixtures/recipes/`.

## Tests

- [x] Add CLI test for successful `validate <recipe> --assets ... --json`.
- [x] Add CLI test for default body type when recipe omits `body_type`.
- [x] Add CLI test for missing recipe path.
- [x] Add CLI test for missing `--assets`.
- [x] Add CLI test for missing recipe file.
- [x] Add CLI test for invalid recipe JSON.
- [x] Add CLI test for missing selections.
- [x] Add CLI test for unknown layer ID.
- [x] Add CLI test for unsupported body type.
- [x] Add service tests for recipe validation behavior.

## Implementation

- [x] Add recipe and validation result models.
- [x] Add recipe validation service that reuses the catalog loader.
- [x] Add validate controller for orchestration and exit codes.
- [x] Add JSON view for validation success and error responses.
- [x] Extend CLI routing for `validate`.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.
- [x] Run direct recipe-validation smoke test.
- [x] Update `.wiki/` with the validation CLI contract.
