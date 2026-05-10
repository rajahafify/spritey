# Tasks: Make Report Artifact Integrity

## Spec

- [x] Add spec, plan, and tasks for `010-make-report-artifact-integrity`.

## Tests

- [x] Extend make-service report test to require additive `artifacts.output_png` fields.
- [x] Assert report artifact `sha256` and `bytes` match the written output PNG.
- [x] Add make-service metadata-failure regression (`RENDER_FAILED`, no report written).
- [x] Extend CLI make JSON/report test to assert unchanged stdout envelope and additive artifact metadata.

## Implementation

- [x] Add additive report model fields for `artifacts.output_png.sha256` and `artifacts.output_png.bytes`.
- [x] Compute artifact metadata from the actual written `--out` PNG when `--report` is requested.
- [x] Return `RENDER_FAILED` when artifact metadata computation fails.
- [x] Keep make CLI/stdout JSON envelope and report schema version unchanged.
- [x] Preserve behavior when `--report` is omitted.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Wiki

- [x] Add raw wiki note for spec 010 implementation slice.
- [x] Update wiki indexes/concepts/log for additive artifact-integrity behavior.
