# Tasks: Make Recolor Parity with python_source/recolorer.py

## Spec

- [x] Add `014-make-recolor-parity` spec, plan, and tasks.

## Tests

- [x] Add failing `app/services/palette_recolorer_test.go` for exact variant recolor.
- [x] Add failing `app/services/palette_recolorer_test.go` for unknown variant unchanged behavior.
- [x] Add failing `app/services/palette_recolorer_test.go` for fuzzy variant matching.
- [x] Add failing `app/services/palette_recolorer_test.go` for tolerance/transparency parity.
- [x] Add failing `app/services/make_service_test.go` for recolor-applied make output.
- [x] Add failing `app/services/make_service_test.go` for unknown-variant non-fatal unchanged output.
- [x] Add failing `app/services/make_service_test.go` for fuzzy-variant make output.
- [x] Add failing `app/services/make_service_test.go` for tolerance/transparency make output parity.
- [x] Keep existing make regressions green.

## Implementation

- [x] Add `app/services/palette_recolorer.go`.
- [x] Implement recursive palette file discovery under `assets/palette_definitions`.
- [x] Implement deterministic lexicographic file ordering for palette scans.
- [x] Implement exact then fuzzy variant lookup parity.
- [x] Implement base palette fallback order: `light`, `base`, `beige`, `brown`.
- [x] Implement recolor replacement with per-channel tolerance `<=2`, non-transparent pixels only.
- [x] Treat missing/invalid source or target palettes as non-fatal no-op.
- [x] Wire recolorer into `app/services/make_service.go` before row compositing.
- [x] Keep make CLI/report envelope keys and exit mapping unchanged.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Wiki

- [x] Add raw note for spec 014.
- [x] Update wiki indexes and make-command concept coverage.
