# Tasks: Make Report Provenance

## Spec

- [x] Add spec, plan, and tasks for `008-make-report-provenance`.

## Tests

- [x] Add make-service tests for report provenance fields and deterministic `layers.composed` ordering.
- [x] Add make-service test for fallback body-type provenance in composed entries.
- [x] Add CLI make test verifying report provenance while make JSON stdout envelope remains unchanged.

## Implementation

- [x] Add internal model support for requested body type provenance in make flow.
- [x] Add report v1 additive pack metadata fields (`pack.id`, `pack.name`).
- [x] Add report v1 additive recipe metadata fields (`body_type_requested`, `body_type_effective`).
- [x] Add report v1 `layers.composed` entries with deterministic ordering and required fields.
- [x] Keep `layers.applied` unchanged for backward compatibility.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Wiki

- [x] Update wiki raw note + indices + activity log for spec 008 slice.
- [x] Update compiled concept coverage for make report provenance.
