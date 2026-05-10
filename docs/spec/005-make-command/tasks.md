# Tasks: Make Command

## Spec

- [x] Create feature spec, technical plan, and tasks under `docs/spec/005-make-command/`.

## Tests

- [x] Add CLI test for successful make with `--report`.
- [x] Add CLI test for successful make without `--report`.
- [x] Add CLI test for missing `<recipe>` positional argument.
- [x] Add CLI test for missing `--assets`.
- [x] Add CLI test for missing `--out`.
- [x] Validate `--json` optional behavior via stable JSON-mode tests while preserving non-JSON command compatibility.
- [x] Add CLI test for nonexistent recipe file.
- [x] Add CLI test for nonexistent assets directory.
- [x] Add CLI test asserting stable success JSON envelope keys.
- [x] Add CLI test asserting stable error JSON envelope keys and error object fields.
- [x] Add CLI test asserting render failure maps to exit code `6` and `RENDER_FAILED`.
- [x] Add service test for report v1 required fields.
- [x] Add service test for deterministic `animation_ids` ordering in report v1.
- [x] Add service test for deterministic `layers.applied` ordering in report v1.
- [x] Add deterministic output test for PNG dimensions/hash using runtime-generated fixture assets.
- [x] Ensure all fixture assets used by tests are runtime-generated deterministic assets (no third-party assets).

## Implementation

- [x] Add make-command controller orchestration for argument validation, validation delegation, render execution, and report writing.
- [x] Add make result model with stable success/error envelope fields.
- [x] Add structured make error model with `code`, `message`, optional `field`, optional `details`.
- [x] Add render/report services required for minimal make + report v1.
- [x] Add make JSON view for stable stdout contract.
- [x] Add exit-code mapping preserving render failure as exit code `6`.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.
- [x] Run direct make-command smoke test with `--report`.
- [x] Run direct make-command smoke test without `--report`.

## Deferred To Next Spec (Out Of Scope Here)

- [ ] Human-readable non-JSON make output mode.
- [ ] Additional report sections beyond minimal v1 contract.
- [ ] Alternate render output formats (GIF/APNG/webp).
- [ ] Batch make and multi-recipe orchestration.
- [ ] Render caching and performance optimization.
