# Tasks: Inspect Layer

## Spec

- [x] Create feature spec, technical plan, and tasks under `docs/spec/002-inspect-layer/`.

## Tests

- [x] Add CLI test for successful `inspect layer body_human --assets ... --json`.
- [x] Add CLI test for missing target after `inspect`.
- [x] Add CLI test for unsupported target such as `inspect pack`.
- [x] Add CLI test for missing layer ID after `inspect layer`.
- [x] Add CLI test for unknown layer ID.
- [x] Add CLI test for missing `--assets`.
- [x] Add CLI test for invalid assets directory.
- [x] Add service tests for exact layer lookup across categories.

## Implementation

- [x] Add inspect layer model/result type.
- [x] Add inspect service that reuses the catalog loader.
- [x] Add inspect controller for validation, orchestration, and exit codes.
- [x] Add JSON view for inspect layer success and error responses.
- [x] Extend CLI routing for `inspect layer`.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.
- [x] Run direct inspect-layer smoke test.
- [x] Update `.wiki/` with the inspect-layer CLI contract.
