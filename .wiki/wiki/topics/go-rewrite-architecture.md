---
title: "Go Rewrite Architecture"
category: topic
sources:
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

As of the Docker development setup, the root Go app has a minimal module and command entrypoint:

- `go.mod`
- `cmd/spritey/main.go`
- `Dockerfile`
- `compose.yaml`
- `Makefile`

The app still has no real command parser, schema code, catalog loader, recipe validation, renderer, or report generator.

Normal development should use Docker Compose through Make targets. Native Go is optional.

## Sources

- [Spritey Repository Baseline](../../raw/repos/2026-05-10-spritey-repo-baseline.md) - current tree and implementation state.
- [Agent Rules and Constitution](../../raw/notes/2026-05-10-agent-rules-and-constitution.md) - architecture rules.
