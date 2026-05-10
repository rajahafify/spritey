# 2026-05-11 Make Frame-Readiness Parity Implementation

Date: 2026-05-11  
Scope: spec 015 (`docs/spec/015-make-frame-readiness-parity/`)

Implemented Python-compositor-style make frame-readiness behavior: missing per-animation layer frames are non-fatal in `make`, while `validate` remains strict.

Changes:

- Updated `app/services/recipe_validator.go`:
  - introduced internal validation mode split;
  - `Validate(...)` remains strict and still fails on `MISSING_SPRITE_FRAME`;
  - added `ValidateForMake(...)` that preserves metadata/path validation but skips strict required-frame preflight.
- Updated `app/services/make_service.go`:
  - switched `make` validation call to `ValidateForMake(...)`;
  - row composition continues to skip missing layer frames non-fatally and emits only available rows.
- Updated tests:
  - `app/services/make_service_test.go` now asserts make success for partial missing-frame scenarios;
  - `cmd/spritey/main_test.go` now asserts make CLI success and emitted-row summaries for those scenarios;
  - strict validate missing-frame tests remain unchanged and passing.
- Preserved command/output contracts:
  - make CLI flags unchanged;
  - make JSON/report envelope keys unchanged;
  - make exit-code mapping unchanged.

Validation:

- `go test ./...`
- `make native-ci`
