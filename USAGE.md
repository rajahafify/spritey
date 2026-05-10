# Spritey Usage (Agent Guide)

Spritey is a file-based Go CLI for sprite generation workflows. For automation, use `--json` and branch on `ok` + `errors[*].code` instead of parsing text output.

## Quickstart (JSON-first)

```powershell
$bootstrap = spritey --download-lpc-assets --json | ConvertFrom-Json
$assets = $bootstrap.assets.path
spritey assets validate --assets $assets --json
spritey catalog --assets $assets --json
spritey inspect layer body_human --assets $assets --json
spritey validate .\recipes\hero.json --assets $assets --json
spritey make .\recipes\hero.json --assets $assets --out .\output\hero.png --report .\output\hero.report.json --json
```

Recommended flow:
1. `--download-lpc-assets` (bootstrap LPC-compatible assets to user cache)
2. `assets validate` (preflight assets directory)
3. `catalog` / `inspect layer` (discover and inspect IDs)
4. `validate` (strict recipe + semantic checks)
5. `make` or `make batch` (artifact generation)

## Command Reference

- `spritey catalog --assets <dir> --json`  
  List catalog categories/layers from a compatible assets pack.

- `spritey --download-lpc-assets [--json] [--force]`  
  Download/install LPC-compatible assets into the default user cache path.
  - text mode: shows live progress/status on `stderr`
  - json mode: keeps `stdout` JSON-only (no progress payload)

- `spritey inspect layer <layer-id> --assets <dir> --json`  
  Return full metadata for one layer ID (exact match).

- `spritey validate <recipe.json> --assets <dir> --json`  
  Strictly validate recipe structure + catalog compatibility.

- `spritey make <recipe.json> --assets <dir> --out <sprite.png> [--report <report.json>] --json`  
  Generate one spritesheet PNG (and optional report JSON).

- `spritey make batch <manifest.json> --assets <dir> [--json]`  
  Run multiple `make` jobs from one manifest.

- `spritey assets validate --assets <dir> --json`  
  Preflight compatible assets structure and parse readiness.

## Stable JSON Envelope Keys

Depend on these top-level keys:

- `download-lpc-assets`: `ok`, `command`, `assets`, `warnings`, `errors`
- `catalog`: `ok`, `pack`, `categories`, `warnings`, `errors`
- `inspect layer`: `ok`, `layer`, `warnings`, `errors`
- `validate`: `ok`, `recipe`, `warnings`, `errors`
- `make`: `ok`, `command`, `outputs`, `summary`, `warnings`, `errors`
- `make batch`: `ok`, `command`, `summary`, `jobs`, `warnings`, `errors`
- `assets validate`: `ok`, `assets`, `pack`, `summary`, `warnings`, `errors`

Notes:
- `make` returns `command: "make"`.
- `make batch` returns `command: "make-batch"` in JSON output.
- `--download-lpc-assets` returns `command: "download-lpc-assets"` in JSON output.

## Batch Manifest Example (`schema_version: "1"`)

```json
{
  "schema_version": "1",
  "jobs": [
    {
      "id": "hero-female",
      "recipe": "recipes/hero-female.json",
      "out": "output/hero-female.png",
      "report": "output/hero-female.report.json"
    },
    {
      "id": "hero-male",
      "recipe": "recipes/hero-male.json",
      "out": "output/hero-male.png"
    }
  ]
}
```

```powershell
spritey make batch .\manifests\batch.json --assets .\assets --json
```

Batch path rule: `recipe`, `out`, and `report` are resolved relative to the manifest file directory.

## Exit-Code and Error-Code Handling

Primary exit codes:

- `0`: success
- `1`: general/internal error
- `2`: CLI usage/argument errors
- `3`: invalid assets or make input readiness issues
- `4`: invalid recipe or invalid batch manifest
- `5`: semantic validation failures (for example `UNKNOWN_LAYER_ID`)
- `6`: render failure (`RENDER_FAILED`)

Agent handling pattern:
1. Check process exit code.
2. Parse JSON.
3. Branch on `ok`.
4. If `ok == false`, handle `errors[0].code` first.

Common error-code families:
- Download/install: `DOWNLOAD_FAILED`, `EXTRACT_FAILED`, `WRITE_ASSETS_FAILED`, `RESOLVE_CACHE_DIR_FAILED`
- CLI: `MISSING_COMMAND`, `UNKNOWN_COMMAND`, `UNKNOWN_ARGUMENT`, `MISSING_*_VALUE`
- Assets: `ASSETS_DIRECTORY_NOT_FOUND`, `MISSING_PACK_JSON`, `INVALID_SHEET_DEFINITION_JSON`, `MISSING_SPRITESHEETS`, `MISSING_PALETTE_DEFINITIONS`
- Recipe/manifest: `INVALID_RECIPE_JSON`, `MISSING_SELECTIONS`, `INVALID_BATCH_MANIFEST_JSON`, `EMPTY_BATCH_JOBS`
- Semantic/render: `UNKNOWN_LAYER_ID`, `UNSUPPORTED_BODY_TYPE`, `MISSING_SPRITE_FRAME`, `RENDER_FAILED`

Exit-code note:
- Mapping is command-specific. Example: `UNKNOWN_LAYER_ID` is `exit 5` in `validate`, but `exit 3` in `make`.
- For `--download-lpc-assets`, invalid extracted assets map to `exit 3`; operational failures map to `exit 1`.

## Troubleshooting Checklist

- `--assets` points to an existing directory with:
  - `pack.json`
  - `sheet_definitions/`
  - `spritesheets/`
  - `palette_definitions/`
- Recipe file exists and is valid JSON.
- Layer IDs in recipe come from `catalog`/`inspect layer`.
- Spritey creates parent directories for `--out` and `--report` paths automatically.
- For batch manifests:
  - `schema_version` is `"1"`
  - `jobs` is non-empty
  - each job includes `id`, `recipe`, `out`
- Use `--json` for all automated flows; do not rely on human text output format.
