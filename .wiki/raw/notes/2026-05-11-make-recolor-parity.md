# 2026-05-11 Make Recolor Parity Implementation

Date: 2026-05-11  
Scope: spec 014 (`docs/spec/014-make-recolor-parity/`)

Implemented Python-equivalent recolor behavior from `python_source/recolorer.py` in the `make` render pipeline.

Changes:

- Added `app/services/palette_recolorer.go`:
  - recursive palette discovery at `assets/palette_definitions/**/{material}_*.json`;
  - deterministic lexicographic file-order scanning;
  - per-file variant lookup with exact key first, then fuzzy match (`requested in key` OR `key in requested`, case-insensitive);
  - source palette fallback order: `light`, `base`, `beige`, `brown`;
  - per-channel tolerance `<=2` for RGB replacement;
  - recolor only non-transparent pixels (`alpha > 0`);
  - missing/invalid source or target palette returns original image non-fatally.
- Integrated recolor into `app/services/make_service.go`:
  - recolor is attempted only when layer has `recolor_material` and selection has `palette_variant`;
  - recolor occurs before per-row compositing.
- Preserved make contract shape:
  - CLI flags unchanged;
  - make JSON/report envelope keys unchanged;
  - exit-code mapping unchanged.
- Added tests:
  - `app/services/palette_recolorer_test.go` for exact, unknown, fuzzy, and tolerance/transparency parity behaviors;
  - `app/services/make_service_test.go` for recolor-applied make output, unknown non-fatal no-op, fuzzy matching, and tolerance/transparency parity.

Validation:

- `go test ./...`
- `make native-ci`
