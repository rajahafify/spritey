# Feature Spec: Download LPC Assets

## Goal

Add a top-level Spritey CLI mode to download compatible LPC assets into user cache so users can run Spritey without manual asset bootstrapping steps.

## User Stories

- As a user, I can run `spritey --download-lpc-assets --json` and get a validated local assets directory path.
- As an automation agent, I receive a stable JSON envelope with path/source metadata and structured errors.
- As an operator, I can force reinstall via `--force`.

## Functional Requirements

- Support top-level invocation:
  - `spritey --download-lpc-assets`
  - `spritey --download-lpc-assets --json`
  - `spritey --download-lpc-assets --json --force`
- Default install path is `os.UserCacheDir()/spritey/assets/lpc`.
- Download LPC upstream archive and extract compatible subset:
  - `pack.json` (or generate compatible `pack.json` when missing)
  - `sheet_definitions/`
  - `spritesheets/`
  - `palette_definitions/`
- Validate installed assets with existing assets validator before success.
- Keep all existing command behavior unchanged.

## Contract Requirements

- JSON response envelope keys for this mode are stable and ordered:
  - `ok`, `command`, `assets`, `warnings`, `errors`
- `command` is `download-lpc-assets`.
- Exit codes:
  - `0` success
  - `2` invalid CLI args for this mode
  - `3` invalid downloaded assets (validator failure)
  - `1` operational download/extract/write failures

## Testing Requirements

- CLI tests:
  - success `--download-lpc-assets --json`
  - unknown arg after this mode returns CLI usage failure
  - download failure maps to exit `1`
  - invalid downloaded assets maps to exit `3`
- Service tests:
  - default path resolution uses `UserCacheDir()/spritey/assets/lpc`
  - install behavior succeeds with injected archive payload
  - force reinstall refetches archive
- Tests must not perform real network calls; inject fetch behavior.

## Out Of Scope

- Additional remote providers or mirrors.
- Custom install path flags.
- Asset version pinning beyond default upstream branch snapshot.
