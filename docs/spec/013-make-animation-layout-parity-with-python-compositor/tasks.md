# Tasks: Make Animation/Layout Parity with python_source/compositor.py

## Spec

- [x] Add `013-make-animation-layout-parity-with-python-compositor` spec, plan, and tasks.

## Tests

- [x] Update service tests for fixed LPC row order and emitted-row filtering.
- [x] Add service test for row-height parity based on first contributing layer.
- [x] Add service test for fixed width `832`.
- [x] Add service test for transparent fallback `832x256` when no rows emit.
- [x] Update service/report assertions so `render.animation_ids` and `render.frame_count` reflect emitted rows.
- [x] Update CLI make tests to assert unchanged envelope keys with shifted summary values.
- [x] Keep regressions green.

## Implementation

- [x] Update `app/services/make_service.go` to use fixed LPC animation order.
- [x] Stop driving make row order from `pack.defaults.animations`.
- [x] Skip missing layer frames at compose time (non-fatal per layer).
- [x] Emit rows only when at least one layer contributes.
- [x] Implement first-layer row-height and padded/clip compositing parity.
- [x] Fix strip width to `832`.
- [x] Stack emitted rows vertically using summed row heights.
- [x] Emit transparent `832x256` fallback when no rows emit.
- [x] Keep CLI surface, envelope keys, and exit mapping unchanged.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Wiki

- [x] Add raw note for spec 013.
- [x] Update wiki indexes and make-command concept coverage.
