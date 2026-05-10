# Technical Plan: Make Recolor Parity with python_source/recolorer.py

## Architecture

- Add palette recolor service in `app/services/palette_recolorer.go`:
  - recursive palette discovery under `assets/palette_definitions`;
  - deterministic lexicographic file ordering;
  - exact then fuzzy variant lookup;
  - base palette fallback order (`light`, `base`, `beige`, `brown`);
  - per-channel tolerance `<=2` and `alpha > 0` pixel replacement.
- Integrate recolor into `app/services/make_service.go` layer-frame pipeline before row compositing.
- Keep make CLI/report envelopes and exit mapping unchanged.

## TDD Sequence

1. Add failing recolor service tests in `app/services/palette_recolorer_test.go` (exact, unknown, fuzzy, tolerance/transparency).
2. Add failing make-service parity tests in `app/services/make_service_test.go` (applied, unknown non-fatal, fuzzy, tolerance/transparency).
3. Implement `app/services/palette_recolorer.go`.
4. Wire recolorer into `app/services/make_service.go`.
5. Run validation:

```bash
go test ./...
make native-ci
```
