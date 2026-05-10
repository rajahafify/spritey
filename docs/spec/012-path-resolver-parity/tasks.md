# Tasks: Path Resolver Parity with python_source/compositor.py

## Spec

- [x] Add `012-path-resolver-parity` spec, plan, and tasks.

## Tests

- [x] Add resolver unit tests for A/B/C/D lookup order.
- [x] Add resolver test for `behind` filename exclusion in mapped slash/thrust directories.
- [x] Add resolver test for C-over-D precedence on slash/thrust.
- [x] Add resolver not-found test.
- [x] Add make service test for mapped D-path success.
- [x] Add make service test for missing resolution returning `MISSING_SPRITE_FRAME`.
- [x] Add CLI test for mapped D-path success.
- [x] Add CLI test for missing resolution returning `MISSING_SPRITE_FRAME`.

## Implementation

- [x] Add `SpriteFrameResolver` shared service.
- [x] Wire resolver into `RecipeValidator` readiness checks.
- [x] Wire resolver into `MakeService` render frame lookup.
- [x] Keep make/validate contracts and exit mappings unchanged.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Wiki

- [x] Add raw note for spec 012.
- [x] Update wiki indexes and make-command concept coverage.
