# 2026-05-10 Make Human Output Implementation

Date: 2026-05-10
Scope: spec 006 (`docs/spec/006-make-human-output/`)

Implemented deterministic non-JSON output for `spritey make` while preserving existing JSON contract and exit-code behavior.

Changes:

- Added `app/views/make_text.go` to render line-oriented success output in stable field order.
- Updated `app/controllers/make_controller.go` to delegate non-JSON success output to the view.
- Added CLI tests in `cmd/spritey/main_test.go`:
  - success without `--json` and with `--report`
  - success without `--json` and without `--report`
  - failure without `--json` retains `<CODE>: <message>` stderr format

Non-JSON success output contract:

1. `ok: make`
2. `png: <path>`
3. `report: <path>` (when provided)
4. `frame_count: <n>`
5. `canvas: <width>x<height>`
6. `animation_count: <n>`

Validation:

- `go test ./...`
- `make native-ci`
