# Technical Plan: Make Report Artifact Integrity

## Architecture

Keep scope local to make reporting:

- `app/models/make_report_provenance.go`
  - add additive report artifact metadata types/fields.
- `app/services/make_service.go`
  - compute output PNG artifact metadata from the actual written `--out` file;
  - include additive metadata in report v1;
  - fail with `RENDER_FAILED` when metadata computation fails before report writing.
- `app/services/make_service_test.go`
  - verify additive artifact metadata values and metadata-failure behavior.
- `cmd/spritey/main_test.go`
  - verify make JSON envelope stability and report artifact metadata presence/values.

## Data Flow

1. Validate recipe and render inputs as before.
2. Render and atomically write output PNG as before.
3. If `--report` is provided:
   - compute artifact metadata (`sha256`, `bytes`) from the written output PNG;
   - assemble report v1 with additive `artifacts.output_png` metadata;
   - write report JSON.
4. If `--report` is omitted, skip metadata/report logic.

## Determinism Rules

- `artifacts.output_png.sha256` must be lowercase hex.
- `artifacts.output_png.bytes` must match the output file size in bytes.
- No make stdout envelope key changes.
- Keep report schema version `"1"`.

## TDD Sequence

1. Add/extend failing tests in:
   - `app/services/make_service_test.go`
   - `cmd/spritey/main_test.go`
2. Implement additive model/service changes.
3. Run validation:

```bash
go test ./...
make native-ci
```
