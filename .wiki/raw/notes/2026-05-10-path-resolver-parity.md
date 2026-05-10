# 2026-05-10 Path Resolver Parity Implementation

Date: 2026-05-10  
Scope: spec 012 (`docs/spec/012-path-resolver-parity/`)

Implemented shared sprite frame resolver parity with `python_source/compositor.py` lookup rules for validate/make frame discovery.

Changes:

- Added shared resolver service:
  - `app/services/sprite_frame_resolver.go`
  - order:
    1) `{prefix}/{anim}.png`
    2) `{prefix}/fg/{anim}.png`
    3) `{prefix}/{anim}/*.png` (first hit)
    4) for `slash`/`thrust`: `{prefix}/attack_slash/*.png` or `{prefix}/attack_thrust/*.png`, excluding filenames containing lowercase `behind`
  - C path precedence over D path is preserved by lookup order.
- Wired resolver into readiness checks in `app/services/recipe_validator.go`.
- Wired resolver into render frame loading in `app/services/make_service.go` so make uses the same resolution path as validate.
- Kept existing contracts:
  - missing frame remains `MISSING_SPRITE_FRAME` in validation class;
  - existing make/validate exit mappings unchanged.
- Added tests:
  - `app/services/sprite_frame_resolver_test.go` for A/B/C/D, behind-filter, precedence, and not-found.
  - `app/services/make_service_test.go` parity case where only mapped D path exists.
  - `cmd/spritey/main_test.go` CLI parity success/failure cases for slash mapping.

Validation:

- `go test ./...`
- `make native-ci`
