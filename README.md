# Spritey

Spritey is a CLI tool for generating animated 2D character spritesheets from recipe files and compatible assets.

It is designed for agents and automation: inputs are file-based, outputs are explicit, validation is available before generation, and JSON output is stable enough for scripts to consume.

Repository:

```text
https://github.com/rajahafify/spritey
```

## What Spritey Does

Spritey reads:

- an assets directory
- a recipe JSON file
- optional CLI overrides

Spritey produces:

- a final transparent PNG spritesheet
- optional JSON reports
- validation errors and warnings that agents can parse

Typical workflow:

```bash
spritey catalog --assets ./assets --json
spritey inspect layer torso_armour_plate --assets ./assets --json
spritey validate recipes/knight.json --assets ./assets --json
spritey make recipes/knight.json --assets ./assets --out output/knight.png --report output/knight.report.json
```

## Installation

Download the Spritey binary for your platform and place it somewhere on your `PATH`.

Verify the install:

```bash
spritey --version
spritey --help
```

## Development Prerequisites

Normal development is Docker-based and Unix-centric.

Required:

- Git
- Docker with Docker Compose
- `make`

Optional:

- Go installed locally for faster native feedback

## Docker Development

Build the development image:

```bash
make docker-build
```

Run tests:

```bash
make docker-test
```

Run the CLI scaffold:

```bash
make docker-run
```

Open a shell inside the Go development container:

```bash
make docker-shell
```

## Native Development

Native Go is optional. If Go is installed locally, the equivalent commands are:

```bash
git clone git@github.com:rajahafify/spritey.git
cd spritey
go test ./...
go build -o bin/spritey ./cmd/spritey
```

## Assets

Spritey works with compatible assets directories.

Expected shape:

```text
assets/
  pack.json
  sheet_definitions/
  spritesheets/
  palette_definitions/
```

`pack.json` defines pack-level defaults, such as:

- default body type
- animation list
- canvas width
- palette fallback rules
- missing body-type fallback behavior

Example:

```json
{
  "schema_version": "1",
  "id": "universal-lpc-compatible",
  "name": "Universal LPC Compatible Assets",
  "defaults": {
    "body_type": "male",
    "animations": ["spellcast", "thrust", "walk", "slash", "shoot", "hurt"],
    "canvas_width": 832,
    "output_format": "png",
    "missing_body_type_fallback": "male",
    "palette_source_fallbacks": ["light", "base", "beige", "brown"]
  }
}
```

Validate an assets directory:

```bash
spritey assets validate --assets ./assets --json
```

## Recipes

Recipes are JSON files that describe the character to generate.

Example `recipes/knight.json`:

```json
{
  "body_type": "male",
  "selections": {
    "body": {
      "id": "body",
      "palette_variant": "light"
    },
    "chainmail": {
      "id": "torso_chainmail",
      "palette_variant": "steel"
    },
    "armour": {
      "id": "torso_armour_plate",
      "palette_variant": "steel"
    },
    "belt": {
      "id": "belt_leather",
      "palette_variant": "brown"
    },
    "weapon": {
      "id": "weapon_longsword",
      "palette_variant": "steel"
    }
  }
}
```

If `body_type` is omitted, Spritey uses the default from `assets/pack.json`.

Validate a recipe before generating:

```bash
spritey validate recipes/knight.json --assets ./assets --json
```

## Commands

### `spritey catalog`

Lists available asset categories and layer IDs.

```bash
spritey catalog --assets ./assets
spritey catalog --assets ./assets --json
```

Use this when an agent needs to discover valid options before writing a recipe.

### `spritey inspect layer`

Shows details for a single layer.

```bash
spritey inspect layer torso_armour_plate --assets ./assets --json
```

The JSON output includes supported body types, available animations, recolor material, palette options, source paths, and credits when available.

### `spritey validate`

Validates a recipe against an assets directory without generating an image.

```bash
spritey validate recipes/knight.json --assets ./assets --json
```

Validation checks include:

- recipe JSON shape
- missing selections
- unknown layer IDs
- unsupported body type paths
- missing sprite files
- missing or unresolved palette variants
- pack default compatibility

### `spritey make`

Generates a final spritesheet PNG.

```bash
spritey make recipes/knight.json --assets ./assets --out output/knight.png
```

Generate with a machine-readable report:

```bash
spritey make recipes/knight.json --assets ./assets --out output/knight.png --report output/knight.report.json
```

### `spritey assets validate`

Validates a compatible assets directory.

```bash
spritey assets validate --assets ./assets --json
```

## Reports

`spritey make --report` writes a JSON report describing what happened.

Example:

```json
{
  "ok": true,
  "output": "output/knight.png",
  "size": {
    "width": 832,
    "height": 1344
  },
  "pack": {
    "id": "universal-lpc-compatible",
    "name": "Universal LPC Compatible Assets"
  },
  "layers": [
    {
      "category": "body",
      "id": "body",
      "z_pos": 10,
      "path": "body/bodies/male/",
      "recolored": true,
      "palette_variant": "light"
    }
  ],
  "warnings": []
}
```

Reports are intended for agents, CI jobs, build scripts, and asset pipelines.

## Exit Codes

Spritey uses stable exit codes:

```text
0  success
1  general error
2  invalid CLI usage
3  invalid assets directory
4  invalid recipe
5  validation failed
6  render failed
```

When `--json` is used, errors are emitted in a structured format:

```json
{
  "ok": false,
  "errors": [
    {
      "code": "UNKNOWN_LAYER_ID",
      "message": "Layer id not found: fake_helmet",
      "field": "selections.helmet.id"
    }
  ],
  "warnings": []
}
```

## Agent Workflow

When asked to generate a sprite, agents should:

1. Run `spritey catalog --assets ./assets --json`.
2. Inspect likely layers with `spritey inspect layer`.
3. Write a recipe file.
4. Run `spritey validate`.
5. Run `spritey make`.
6. Read the report and summarize the result.

Example:

```bash
spritey catalog --assets ./assets --json > tmp/catalog.json
spritey inspect layer torso_armour_plate --assets ./assets --json > tmp/plate.json
spritey validate recipes/knight.json --assets ./assets --json
spritey make recipes/knight.json --assets ./assets --out output/knight.png --report output/knight.report.json
```

Agents should prefer recipe files over inline JSON.

## Development

Spritey uses Spec Kit-style spec-driven development.

Reference:

```text
https://github.com/github/spec-kit
```

Development flow:

```text
constitution -> feature spec -> technical plan -> tasks -> implementation -> validation
```

Implementation tasks should use TDD when practical:

```text
spec task -> failing test -> minimal implementation -> passing test -> refactor -> validation
```

Project layout:

```text
spritey/
  cmd/
    spritey/
      main.go
  app/
    controllers/
    models/
    services/
    views/
  config/
  docs/
    spec/
  schemas/
  testdata/
```

The Go implementation lives at the repository root. A local ignored Python reference snapshot may exist in `python_source/`, but it is not part of this repository.

## License And Credits

Spritey generates derivative spritesheets from the selected source layers. Generated output may require attribution depending on the selected assets.

Use report and credits output to track which source layers were used.
