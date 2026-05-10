# Technical Plan: Batch Make Manifest v1

## Architecture

Implement a dedicated batch command path that orchestrates existing single-make behavior:

- `cmd/spritey/main.go`
  - parse and dispatch `make batch`.
- `app/controllers/make_batch_controller.go`
  - CLI contract checks and exit mapping for batch.
- `app/services/make_batch_service.go`
  - load/validate manifest;
  - resolve relative paths from manifest directory;
  - execute jobs sequentially through existing `MakeService`;
  - fail fast with partial job results.
- `app/models/make_batch.go`
  - manifest, batch result, summary, and job result models.
- `app/views/make_batch_json.go`
  - stable batch JSON success/error envelopes.
- `app/views/make_batch_text.go`
  - deterministic concise non-JSON summary output.

## Data Flow

1. Parse `spritey make batch <manifest> --assets <dir> [--json]`.
2. Controller validates required CLI options.
3. Service reads and validates manifest (`schema_version = "1"`, non-empty jobs, required job fields).
4. Service resolves relative job paths from manifest directory.
5. Service executes each job in order by calling `MakeService.Make`.
6. On first job failure:
   - append failing job result;
   - return immediately with failing context.
7. Render JSON or text output in stable deterministic order.

## Exit Mapping

- `2`: CLI misuse
- `4`: manifest read/json/schema/jobs issues
- job failures map via existing make codes:
  - invalid-assets/input-validation family -> `3`
  - render failures -> `6`

## TDD Sequence

1. Add failing service tests for success/relative/fail-fast.
2. Add failing CLI tests for envelope, ordering, exits, and text summary.
3. Implement models/service/controller/view/CLI dispatch.
4. Validate:

```bash
go test ./...
make native-ci
```
