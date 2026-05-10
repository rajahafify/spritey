# Tasks: Batch Make Manifest v1

## Spec

- [x] Add `011-batch-make-manifest-v1` spec, plan, and tasks.

## Tests

- [x] Add service tests for two-job success and manifest-order execution.
- [x] Add service test for manifest-relative path resolution.
- [x] Add service fail-fast test that stops after first failure.
- [x] Add CLI tests for JSON envelope stability and ordering.
- [x] Add CLI tests for misuse and manifest error exits.
- [x] Add CLI tests for job-failure exit mappings (`3` and `6`).
- [x] Add CLI test for deterministic non-JSON summary output.

## Implementation

- [x] Add batch manifest/result models.
- [x] Add batch make service that reuses `MakeService` per job.
- [x] Add batch make controller with required validation and exit mapping.
- [x] Add batch make JSON and text views.
- [x] Add CLI parse/dispatch path for `make batch`.
- [x] Keep single-make behavior unchanged.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Wiki

- [x] Add raw note for spec 011.
- [x] Update wiki indexes/concepts/log for `make batch` behavior.
