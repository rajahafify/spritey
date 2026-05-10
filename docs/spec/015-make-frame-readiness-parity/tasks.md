# Tasks: Make Frame-Readiness Parity with Python Compositor

## Spec

- [x] Add `015-make-frame-readiness-parity` spec, plan, and tasks.

## Tests

- [x] Update failing `app/services/make_service_test.go` for non-fatal missing-frame behavior during make.
- [x] Update failing `app/services/make_service_test.go` for missing mapped slash frame with available body frame.
- [x] Update failing `cmd/spritey/main_test.go` make CLI tests to assert success with available-row composition.
- [x] Keep validate strict missing-frame tests unchanged and passing.

## Implementation

- [x] Add non-strict make-oriented recipe validation path in `app/services/recipe_validator.go`.
- [x] Keep strict `Validate(...)` behavior unchanged for validate command paths.
- [x] Switch `app/services/make_service.go` to non-strict validation path.
- [x] Keep make CLI/report envelope keys and exit mapping unchanged.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Wiki

- [x] Add raw note for spec 015.
- [x] Update wiki indexes and make-command concept coverage.
