---
title: "Inspect Layer Implementation"
source: "D:/Work/spritey/docs/spec/002-inspect-layer; D:/Work/spritey/app; D:/Work/spritey/cmd/spritey"
type: notes
ingested: 2026-05-10
tags: [spritey, inspect, cli, assets, implementation]
summary: "Second Spritey product slice implements `spritey inspect layer <layer-id> --assets <dir> --json` using the catalog loader and structured JSON errors."
---

# Inspect Layer Implementation

Spritey's second implemented product slice is:

```bash
spritey inspect layer <layer-id> --assets ./assets --json
```

The feature lives under `docs/spec/002-inspect-layer/` with `spec.md`, `plan.md`, and `tasks.md`.

## Behavior

- `inspect layer` requires an exact layer ID.
- It reuses the catalog loader from the catalog foundation slice.
- It scans sorted catalog categories and layers for the first exact ID match.
- It returns `ok`, `layer`, `warnings`, and `errors`.
- The layer payload includes category, ID, display name, z position, body types, animations, recolor material, path prefix, and credits.

## Error Contract

- Missing inspect target returns exit code `2` and error code `MISSING_INSPECT_TARGET`.
- Unsupported inspect target returns exit code `2` and error code `UNSUPPORTED_INSPECT_TARGET`.
- Missing layer ID returns exit code `2` and error code `MISSING_LAYER_ID`.
- Missing `--assets` returns exit code `2` and error code `MISSING_ASSETS`.
- Invalid assets errors reuse catalog loader exit code `3`.
- Unknown layer ID returns exit code `5` and error code `UNKNOWN_LAYER_ID`.

## Implementation Shape

- `cmd/spritey/main.go` routes `inspect`.
- `app/controllers` owns inspect validation, orchestration, and exit codes.
- `app/services` contains exact layer lookup.
- `app/models` contains the inspect layer result shape.
- `app/views` writes inspect-layer JSON responses.
