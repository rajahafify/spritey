# 2026-05-10 Render-Input Readiness Implementation

Date: 2026-05-10
Scope: spec 007 (`docs/spec/007-render-input-readiness/`)

Implemented render-input readiness checks so validate/make fail early when required source frames are missing, with deterministic warnings for body-type fallback path usage.

Changes:

- Added readiness checks in `app/services/recipe_validator.go`:
  - effective body type: recipe value or `pack.defaults.body_type`
  - layer path resolution for requested body type, then optional `missing_body_type_fallback` when available per layer
  - required animation IDs from `pack.defaults.animations`, fallback `["idle"]`
  - per-layer/per-animation frame existence checks under `assets/spritesheets/<resolved-path>/<animation>.png`
- Added `MISSING_SPRITE_FRAME` validation problem for missing required frames.
- Exposed validate warnings through existing JSON envelope top-level `warnings` in `app/views/validate_json.go`.
- Propagated validation warnings into make JSON and report warnings in `app/services/make_service.go`.
- Updated exit mappings:
  - validate: `MISSING_SPRITE_FRAME` -> exit `5`
  - make: `MISSING_SPRITE_FRAME` -> exit `3`
- Added readiness-focused tests in:
  - `app/services/recipe_validation_test.go`
  - `app/services/make_service_test.go`
  - `cmd/spritey/main_test.go`
- Added readiness baseline fixture PNGs for existing validate-success fixture:
  - `testdata/fixtures/basic-assets/spritesheets/body/human/male/walk.png`
  - `testdata/fixtures/basic-assets/spritesheets/weapon/sword/training/walk.png`

Validation:

- `go test ./...`
- `make native-ci`
