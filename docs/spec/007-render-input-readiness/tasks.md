# Tasks: Render-Input Readiness

## Spec

- [x] Add spec, plan, and tasks for `007-render-input-readiness`.

## Tests

- [x] Add recipe-validator test for missing required frame (`MISSING_SPRITE_FRAME`).
- [x] Add recipe-validator test for fallback body-type path warning.
- [x] Add make-service test for missing required frame with no output PNG.
- [x] Add make-service test for fallback warning propagation to result/report.
- [x] Add CLI validate test for missing frame exit `5`.
- [x] Add CLI validate test for fallback warning success path.
- [x] Add CLI make test for missing frame exit `3` and no output PNG.
- [x] Add CLI make test for fallback warnings in JSON/report.

## Implementation

- [x] Extend recipe validation result model with internal render-input and warning fields.
- [x] Implement per-layer body-type path resolution with optional missing-body-type fallback.
- [x] Implement required-frame existence checks across required animations.
- [x] Return deterministic `MISSING_SPRITE_FRAME` problems.
- [x] Surface validate warnings through existing JSON envelope.
- [x] Reuse validator readiness metadata in make service and propagate warnings.
- [x] Map `MISSING_SPRITE_FRAME` to validate exit `5` and make exit `3`.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.
