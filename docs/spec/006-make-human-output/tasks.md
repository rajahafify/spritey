# Tasks: Make Human Output

## Spec

- [x] Create feature spec, technical plan, and tasks under `docs/spec/006-make-human-output/`.

## Tests

- [x] Add CLI test for successful non-JSON make with `--report`.
- [x] Add CLI test for successful non-JSON make without `--report`.
- [x] Add CLI test for non-JSON failure message and exit code.
- [x] Keep existing JSON make tests passing unchanged.

## Implementation

- [x] Add non-JSON make success writer in `app/views`.
- [x] Update make controller to use the new text writer for non-JSON success.
- [x] Preserve current non-JSON error format and existing exit-code mapping.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Deferred To Next Spec (Out Of Scope Here)

- [ ] Extended report sections beyond report v1.
- [ ] Alternate render output formats (GIF/APNG/webp).
- [ ] Batch make and multi-recipe orchestration.
- [ ] Render caching and performance optimization.
