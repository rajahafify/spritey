# 2026-05-11 Make Animation/Layout Parity Implementation

Date: 2026-05-11  
Scope: spec 013 (`docs/spec/013-make-animation-layout-parity-with-python-compositor/`)

Implemented make animation/layout parity with `python_source/compositor.py` for LPC row emission and strip layout (recolor excluded).

Changes:

- Updated `app/services/make_service.go` render behavior:
  - fixed animation row order: `spellcast`, `thrust`, `walk`, `slash`, `shoot`, `hurt`;
  - row order is no longer driven by `pack.defaults.animations`;
  - per animation row, missing layer frame resolution is skipped non-fatally;
  - row emitted only when at least one layer contributed;
  - first contributing layer defines row height;
  - subsequent layers are padded/clipped to that row height before alpha composite;
  - final strip width fixed to `832`;
  - final canvas stacks emitted rows vertically;
  - no emitted rows returns transparent `832x256` output.
- Preserved make contract shape:
  - CLI flags unchanged;
  - make JSON envelope keys unchanged;
  - report envelope keys unchanged;
  - exit-code mapping unchanged.
- Updated tests:
  - `app/services/make_service_test.go` for fixed LPC order, emitted-row filtering, row-height parity, fixed-width output, fallback canvas, and report emitted-row metadata;
  - `cmd/spritey/main_test.go` make assertions keep envelope keys stable while summary/report values shift to emitted-row behavior.

Validation:

- `go test ./...`
- `make native-ci`
