# Tasks: Download LPC Assets

## Spec

- [x] Add `016-download-lpc-assets` spec, plan, and tasks.

## Tests

- [x] Add failing CLI tests for success, unknown-arg, download-failure, and invalid-assets mapping in `cmd/spritey/main_test.go`.
- [x] Add service tests for cache-dir path resolution and install behavior in `app/services/lpc_assets_downloader_test.go`.
- [x] Ensure tests use injected archive fetchers (no real network in tests).

## Implementation

- [x] Add top-level CLI parse/routing in `cmd/spritey/main.go` for `--download-lpc-assets`.
- [x] Add `app/controllers/download_assets_controller.go`.
- [x] Add `app/services/lpc_assets_downloader.go`.
- [x] Add `app/views/download_assets_json.go`.
- [x] Add `app/models/download_assets.go`.
- [x] Keep existing command behavior unchanged.

## Validation

- [x] Run `go test ./...`.
- [x] Run `make native-ci`.

## Wiki/Docs

- [x] Add raw note for spec 016.
- [x] Update wiki indexes and CLI concept docs.
- [x] Update `USAGE.md` with `--download-lpc-assets` contract.
