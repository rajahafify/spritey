# 2026-05-11 Download LPC Assets Implementation

Date: 2026-05-11  
Scope: spec 016 (`docs/spec/016-download-lpc-assets/`)

Implemented top-level LPC asset bootstrap command:

```bash
spritey --download-lpc-assets [--json] [--force]
```

Changes:

- Added CLI routing and option parsing in `cmd/spritey/main.go`:
  - accepts only `--json` and `--force` in this mode;
  - unknown args return `UNKNOWN_ARGUMENT` (exit 2).
- Added MVC slice:
  - `app/models/download_assets.go`
  - `app/views/download_assets_json.go`
  - `app/controllers/download_assets_controller.go`
  - `app/services/lpc_assets_downloader.go`
- Service behavior:
  - resolves install target via `os.UserCacheDir()/spritey/assets/lpc`;
  - downloads upstream LPC archive from GitHub codeload;
  - extracts required subset (`pack.json`, `sheet_definitions`, `spritesheets`, `palette_definitions`);
  - generates compatible `pack.json` when upstream archive lacks one;
  - validates staged assets via existing `AssetsValidator` before install;
  - supports `--force` reinstall and non-force idempotent reuse of existing valid install.
- Added tests:
  - CLI tests in `cmd/spritey/main_test.go` for success and exit-code mappings;
  - service tests in `app/services/lpc_assets_downloader_test.go` with injected archive fetcher (no real network).

Contracts:

- JSON envelope keys for download mode are stable:
  - `ok`, `command`, `assets`, `warnings`, `errors`
- Exit mapping:
  - `0` success
  - `2` invalid CLI args
  - `3` invalid downloaded assets (validator failures)
  - `1` operational failures (download/extract/write)

Validation:

- `go test ./...`
- `make native-ci`
