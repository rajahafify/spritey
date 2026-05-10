# 2026-05-10 Batch Make Manifest v1 Implementation

Date: 2026-05-10  
Scope: spec 011 (`docs/spec/011-batch-make-manifest-v1/`)

Implemented `spritey make batch` manifest-driven sequential job execution while preserving existing single-make behavior.

Changes:

- Added CLI command path:
  - `spritey make batch <manifest.json> --assets <dir> [--json]`
- Added manifest v1 model and validation:
  - `schema_version` must be `"1"`;
  - jobs list must be non-empty;
  - each job requires `id`, `recipe`, `out`.
- Added relative path resolution:
  - relative `recipe`, `out`, and `report` resolve from manifest file directory.
- Added batch make service:
  - executes jobs sequentially in manifest order;
  - reuses existing `MakeService` per job;
  - fails fast on first job failure.
- Added stable batch JSON envelopes:
  - top-level keys: `ok`, `command`, `summary`, `jobs`, `warnings`, `errors`;
  - `command` = `"make-batch"`;
  - failure payload includes failing job context in top-level error details.
- Added deterministic non-JSON success summary output.
- Exit behavior:
  - `2` for CLI misuse;
  - `4` for manifest file/json/schema/empty-jobs/required-field issues;
  - job failures map using existing make rules (`3`/`6`).

Primary touched areas:

- `cmd/spritey/main.go`, `cmd/spritey/main_test.go`
- `app/controllers/make_batch_controller.go`, `app/controllers/make_controller.go`, `app/controllers/make_exit_codes.go`
- `app/services/make_batch_service.go`, `app/services/make_batch_service_test.go`
- `app/models/make_batch.go`
- `app/views/make_batch_json.go`, `app/views/make_batch_text.go`

Validation:

- `go test ./...`
- `make native-ci`
