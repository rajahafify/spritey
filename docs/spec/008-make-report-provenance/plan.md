# Technical Plan: Make Report Provenance

## Architecture

Follow current MVC boundaries:

- `app/services/make_service.go`
  - build additive report provenance fields from catalog metadata and validator render inputs.
- `app/models/catalog.go`
  - carry requested body type internally from validation into make reporting.
- `app/models/*` (dedicated report model)
  - represent report v1 additive provenance payload without changing make stdout envelopes.
- `app/services/make_service_test.go`
  - verify report provenance shape and ordering.
- `cmd/spritey/main_test.go`
  - verify CLI make JSON envelope remains stable while report contains new fields.

## Data Flow

1. `make` validates recipe and render-input readiness as before.
2. `make` loads catalog as before.
3. `make` composes layers as before.
4. When `--report` is provided, service emits report v1 with additive provenance:
   - pack metadata
   - requested/effective body type metadata
   - ordered composed layer provenance entries
   - existing `layers.applied` retained.

## Determinism Rules

- `layers.applied` ordering remains existing behavior for backward compatibility.
- `layers.composed` ordering follows render order (`z_pos` asc, `id` asc).
- `resolved_path` is slash-normalized before writing report.
- `credits` is always an array (empty when absent).

## TDD Sequence

1. Extend failing tests in:
   - `app/services/make_service_test.go`
   - `cmd/spritey/main_test.go`
2. Implement additive report/model changes.
3. Run validation:

```bash
go test ./...
make native-ci
```
