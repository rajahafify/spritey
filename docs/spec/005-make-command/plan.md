# Technical Plan: Make Command

## Architecture

Follow the existing MVC-style CLI architecture and keep `main.go` thin:

- `cmd/spritey`: parse `make <recipe> --assets <dir> --out <png> [--report <json>] --json`.
- `app/controllers`: validate command arguments, orchestrate validation + render + optional report, and map domain errors to exit codes.
- `app/services`: compose recipe + assets into render inputs, execute render, and produce report v1 data.
- `app/models`: represent make result payload, report v1 payload, and structured error payload.
- `app/views`: emit stable JSON success/error envelopes for stdout.

## Data Flow

1. CLI receives `make` invocation with positional recipe and required flags.
2. Controller validates required arguments and JSON mode requirement for this slice.
3. Controller delegates recipe validation and assets validation through existing service pathways.
4. Render service resolves layers/animations/canvas metadata from validated inputs.
5. Render service writes PNG to `--out`.
6. If `--report` is present, report service writes report v1 JSON to requested path.
7. View returns success envelope with output paths and summary.
8. Any failure is normalized into stable JSON error envelope with mapped exit code.

## Exit Code Mapping

- `2`: CLI misuse (missing required args, missing JSON mode in this spec slice).
- `3`: input validation failures (recipe/assets invalid or not found).
- `6`: render pipeline failure (`RENDER_FAILED`) and PNG write failures inside render flow.

This plan preserves the existing render-failure exit-code convention.

## JSON Contract

Success envelope (stdout):

```json
{
  "ok": true,
  "command": "make",
  "outputs": {
    "png": { "path": "output/sprite.png" },
    "report": { "path": "output/sprite.report.json" }
  },
  "summary": {
    "frame_count": 8,
    "canvas": { "width": 256, "height": 256 },
    "animation_count": 2
  },
  "warnings": [],
  "errors": []
}
```

Error envelope (stdout):

```json
{
  "ok": false,
  "command": "make",
  "outputs": {
    "png": { "path": "output/sprite.png" }
  },
  "summary": {},
  "warnings": [],
  "errors": [
    {
      "code": "MISSING_OUT",
      "message": "Missing required flag: --out",
      "field": "--out",
      "details": {}
    }
  ]
}
```

Notes:
- Envelope keys remain stable across success/failure.
- `outputs.report` is omitted when `--report` is not requested.
- `summary` may be empty on early failures.

## Report v1 Contract

Write report only when `--report` is supplied. Minimal deterministic shape:

```json
{
  "schema_version": "1",
  "command": "make",
  "recipe": { "path": "recipes/hero.json" },
  "assets": { "path": "assets" },
  "output": { "png": { "path": "output/sprite.png" } },
  "render": {
    "canvas": { "width": 256, "height": 256 },
    "frame_count": 8,
    "animation_ids": ["idle", "walk"]
  },
  "layers": {
    "applied": ["body/base", "hair/short_a"]
  },
  "warnings": []
}
```

Determinism expectations:
- Stable key naming and array ordering.
- Deterministic output for same fixture inputs.

## Validation Strategy

Implement with TDD when execution begins:

- Write failing CLI tests for command contract and exit-code mapping.
- Write failing service tests for report v1 shape and deterministic ordering.
- Implement minimum logic to pass tests.
- Refactor only after green tests.

Required validation commands after implementation:

```bash
go test ./...
make native-ci
go run ./cmd/spritey make testdata/fixtures/recipes/basic.json --assets testdata/fixtures/basic-assets --out testdata/tmp/out.png --report testdata/tmp/out.report.json --json
```

