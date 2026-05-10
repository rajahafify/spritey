# Technical Plan: Path Resolver Parity with python_source/compositor.py

## Architecture

- Add `app/services/sprite_frame_resolver.go`:
  - shared resolver service implementing A/B/C/D path search order;
  - slash/thrust mapping and `behind` exclusion rule.
- Wire resolver into `RecipeValidator`:
  - replace direct `{prefix}/{anim}.png` readiness check with resolver lookup;
  - preserve `MISSING_SPRITE_FRAME` and `READ_SPRITE_FRAME_FAILED` classes.
- Wire resolver into `MakeService`:
  - use resolved frame path for source PNG opens during rendering.

## Data Flow

1. `validate` resolves each required animation via shared resolver.
2. Missing resolution returns `MISSING_SPRITE_FRAME`.
3. `make` still starts with `validate`.
4. Render loop resolves each layer/animation via same resolver and opens that exact path.

## TDD Sequence

1. Add resolver unit tests for A/B/C/D, filter, precedence, not-found.
2. Add make service + CLI tests for mapped D success and full miss failure.
3. Implement resolver service and wire validator/make services.
4. Validate:

```bash
go test ./...
make native-ci
```
