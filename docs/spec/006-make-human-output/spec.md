# Feature Spec: Make Human Output

## Goal

Provide deterministic human-readable output for `spritey make` when `--json` is not supplied, while preserving the existing JSON contract.

## User Stories

- As a human operator, I can run `spritey make` without `--json` and immediately see output/report locations and render summary.
- As an agent or script, I can keep using `--json` and receive the same stable JSON envelope and exit-code behavior.

## Functional Requirements

- Command remains:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]
```

- If `--json` is present, behavior is unchanged from spec 005.
- If `--json` is absent and `make` succeeds, stdout must include deterministic line-oriented fields:
  - `ok: make`
  - `png: <path>`
  - `report: <path>` (only when `--report` is set)
  - `frame_count: <n>`
  - `canvas: <width>x<height>`
  - `animation_count: <n>`
- If `--json` is absent and `make` fails, stderr continues to use:
  - `<CODE>: <message>`
- Non-JSON mode must not print JSON payloads.

## Error & Exit Requirements

- Exit code mapping remains unchanged from spec 005:
  - `2` invalid CLI usage
  - `3` input validation failures
  - `6` render failures

## Testing Requirements

- Add CLI test for successful `make` without `--json` and with `--report` verifying deterministic line output.
- Add CLI test for successful `make` without `--json` and without `--report`.
- Add CLI test for non-JSON failure message shape (`<CODE>: <message>`) and expected exit code.
- Keep existing JSON-mode tests passing unchanged.

## Out Of Scope

- New JSON fields or JSON envelope changes.
- Additional report sections beyond report v1.
- Alternate render output formats, batch workflows, or performance optimizations.
