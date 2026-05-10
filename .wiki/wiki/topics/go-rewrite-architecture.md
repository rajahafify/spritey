---
title: "Go Rewrite Architecture"
category: topic
sources:
  - raw/notes/2026-05-10-assets-validation.md
  - raw/notes/2026-05-10-recipe-validation.md
  - raw/notes/2026-05-10-inspect-layer.md
  - raw/notes/2026-05-10-catalog-foundation.md
  - raw/repos/2026-05-10-spritey-repo-baseline.md
  - raw/notes/2026-05-10-agent-rules-and-constitution.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, go, architecture, mvc]
aliases: [Spritey Go Architecture, Rails-style MVC]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "Spritey's Go implementation should live at the repository root and follow a Rails-style MVC split adapted for CLI controllers, models, views, and services."
---

# Go Rewrite Architecture

> Spritey's Go app is intended to live at the repository root. Ignored local prototype directories are not repository content.

The scaffolded root layout is:

```text
cmd/spritey/
app/controllers/
app/models/
app/services/
app/views/
config/
schemas/
testdata/
```

## Responsibility Boundaries

- Models hold domain data and validation rules.
- Controllers handle command orchestration.
- Views format human-facing and JSON outputs.
- Services perform reusable operations such as catalog loading, palette resolution, rendering, and report generation.

## CLI Boundary

CLI parsing should stay thin. `cmd/spritey/main.go` should delegate to controllers and services rather than embedding business logic.

## Current State

As of the assets-validation slice, the root Go app has four metadata commands:

- `go.mod`
- `cmd/spritey/main.go` routes `catalog`, `assets validate`, `inspect layer`, and `validate`
- `app/models` defines catalog, inspect, recipe, validation, and error data
- `app/services` loads asset pack metadata, validates assets-pack structure, performs exact layer lookup, and validates recipes against catalog metadata
- `app/controllers` owns command orchestration and exit codes
- `app/views` writes JSON responses
- `Dockerfile`
- `compose.yaml`
- `Makefile`

The app still has no renderer or report generator.

Normal development should use Docker Compose through Make targets. Native Go is optional.

## Test Baseline

The first test is intentionally minimal and exists to prove that local Go and GitHub Actions can execute the project test path.

GitHub Actions should use Docker Compose through `make docker-ci`. Local development should use `make ci`, which runs Docker Compose when available and falls back to native Go when Docker is unavailable.

## Sources

- [Spritey Repository Baseline](../../raw/repos/2026-05-10-spritey-repo-baseline.md) - current tree and implementation state.
- [Assets Validation Implementation](../../raw/notes/2026-05-10-assets-validation.md) - fourth implemented metadata command.
- [Recipe Validation Implementation](../../raw/notes/2026-05-10-recipe-validation.md) - third implemented metadata command.
- [Inspect Layer Implementation](../../raw/notes/2026-05-10-inspect-layer.md) - second implemented metadata command.
- [Catalog Foundation Implementation](../../raw/notes/2026-05-10-catalog-foundation.md) - first implemented command and MVC package split.
- [Agent Rules and Constitution](../../raw/notes/2026-05-10-agent-rules-and-constitution.md) - architecture rules.
