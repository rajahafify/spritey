# Technical Plan: Render-Input Readiness

## Architecture

Keep current MVC layering:

- `app/services/recipe_validator.go`
  - add render-readiness checks and fallback warning generation.
- `app/views/validate_json.go`
  - expose validator warnings through existing top-level `warnings`.
- `app/services/make_service.go`
  - consume validator-resolved render inputs and propagate warnings into make result/report.
- `app/controllers/{validate,make}_controller.go`
  - map `MISSING_SPRITE_FRAME` to required command exit codes.

## Data Flow

1. Validate loads catalog + recipe.
2. Validate resolves effective body type and required animations.
3. Validate resolves per-layer path by requested body type, optionally fallback body type.
4. Validate checks required source frames exist for each selected layer and animation.
5. Validate returns:
   - existing recipe payload,
   - warning list,
   - internal render-input metadata for make.
6. Make calls validator first, reuses resolved readiness metadata, and only renders when validation passed.

## Determinism Rules

- Selection iteration remains sorted by selection key.
- Required animation IDs are normalized deterministically.
- Missing-frame error message/field uses deterministic relative frame path.
- Fallback warning strings are deterministic with layer id, requested type, fallback type, and resolved path.

## TDD Sequence

1. Add failing tests in:
   - `app/services/recipe_validation_test.go`
   - `app/services/make_service_test.go`
   - `cmd/spritey/main_test.go`
2. Implement model/service/controller/view updates.
3. Run:

```bash
go test ./...
make native-ci
```
