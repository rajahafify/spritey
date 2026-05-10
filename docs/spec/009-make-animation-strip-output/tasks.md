# Tasks: Make Animation Strip Output

## Spec

- [x] Add spec, plan, and tasks for `009-make-animation-strip-output`.

## Tests

- [x] Update make-service deterministic output test to assert `8x16` strip dimensions and re-baseline hash.
- [x] Add make-service strip row pixel assertion to verify per-frame variation.
- [x] Add make-service single-animation regression test (`8x8` remains unchanged).
- [x] Update CLI make tests to assert strip-aware summary height and unchanged envelope keys.
- [x] Add CLI/report test verifying `animation_ids` order aligns with strip row order.

## Implementation

- [x] Change make rendering to compose one frame per required animation ID.
- [x] Keep per-frame layer composition order unchanged (`z_pos` then `id`).
- [x] Emit vertical strip PNG dimensions (`width=frameWidth`, `height=frameHeight*frameCount`).
- [x] Keep make JSON/report key structures unchanged while updating summary/render dimension values.
- [x] Preserve warnings and existing exit behavior.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Wiki

- [x] Add raw wiki note for spec 009 implementation slice.
- [x] Update wiki raw/index and compiled concept coverage for strip rendering behavior.
- [x] Append wiki activity log entry for spec 009 ingest/compile updates.
