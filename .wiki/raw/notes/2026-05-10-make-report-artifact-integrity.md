# 2026-05-10 Make Report Artifact Integrity Implementation

Date: 2026-05-10
Scope: spec 010 (`docs/spec/010-make-report-artifact-integrity/`)

Expanded make report v1 with additive artifact-integrity metadata derived from the actual written output PNG while preserving make CLI behavior and stdout envelope contracts.

Changes:

- Added additive report section:
  - `artifacts.output_png.sha256` (lowercase hex SHA-256)
  - `artifacts.output_png.bytes` (file size in bytes)
- Artifact metadata is computed from the already-written `--out` PNG file.
- `--report` omitted path remains unchanged (no metadata/report behavior added).
- Metadata computation failures return `RENDER_FAILED` (existing exit-code mapping remains unchanged).
- Metadata failure path avoids writing report output.
- Added/updated tests in:
  - `app/services/make_service_test.go`
  - `cmd/spritey/main_test.go`

Validation:

- `go test ./...`
- `make native-ci`
