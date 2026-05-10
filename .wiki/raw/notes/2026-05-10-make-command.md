# 2026-05-10 Make Command Implementation

Date: 2026-05-10
Scope: spec 005 (`docs/spec/005-make-command/`)

Implemented `spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]` with MVC wiring:

- `cmd/spritey/main.go` adds `make` parsing (`recipe`, `--assets`, `--out`, optional `--report`, optional `--json`).
- `app/controllers/make_controller.go` validates required args and orchestrates render/report flow.
- `app/services/make_service.go` performs minimal deterministic render composition and optional report v1 write.
- `app/views/make_json.go` emits stable JSON success/error envelopes.
- `app/models/make.go` defines make result, report v1, and structured error models.

Behavior highlights:

- Validates recipe and assets with existing validator/loader contracts before render.
- Writes PNG atomically via temp file + rename to avoid corrupted output on failure.
- Optional `--report` writes minimal report v1 with deterministic ordering:
  - ordered `render.animation_ids`
  - ordered `layers.applied`
- `RENDER_FAILED` maps to exit code `6`.

Tests added:

- CLI success with and without report.
- CLI argument misuse (`MISSING_RECIPE`, `MISSING_ASSETS`, `MISSING_OUT`).
- nonexistent inputs.
- JSON envelope stability for success/error.
- render failure exit-code mapping (`6`, `RENDER_FAILED`).
- service-level report v1 required fields and ordering checks.
- deterministic PNG dimension + hash test using runtime-generated fixture assets.
