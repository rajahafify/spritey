# Spritey Constitution

## Product Identity

Spritey is a Go CLI for generating animated 2D character spritesheets from recipe files and compatible assets.

The command name is `spritey`.

The Go implementation lives at the repository root. Legacy prototypes may exist locally in ignored directories, but they are not part of the repository.

## Primary User

Agents are first-class users.

Spritey must be easy for coding agents and automation scripts to operate reliably without guessing hidden state, folder structure, valid layer IDs, palette options, or output behavior.

## Core Product Flow

Spritey reads:

- a compatible assets directory
- a recipe file
- optional CLI overrides

Spritey produces:

- final PNG spritesheets
- machine-readable JSON reports where requested
- clear validation errors and warnings

The expected command family is:

```bash
spritey catalog --assets ./assets --json
spritey inspect layer <layer-id> --assets ./assets --json
spritey validate recipe.json --assets ./assets --json
spritey make recipe.json --assets ./assets --out output/sprite.png --report output/sprite.report.json
```

## Asset Packs

The user-facing term is `assets`.

A compatible assets directory should support:

```text
assets/
  pack.json
  sheet_definitions/
  spritesheets/
  palette_definitions/
```

`pack.json` owns pack-level defaults such as body type, animation list, canvas width, palette fallback rules, and missing body-type fallback behavior.

Spritey should not hardcode behavior that belongs to the asset pack unless a safe internal fallback is explicitly documented in a spec.

Spritey must not bundle third-party sprite assets in this repository. Users download, install, or provide compatible assets at runtime.

## Recipe Files

Recipes are file-based JSON documents.

Inline JSON must not be the primary workflow.

A recipe describes the desired character through selections such as body, clothing, armor, hair, weapons, and palette variants.

## Behavioral Reference

The agreed recipe-to-spritesheet behavior is the behavioral reference for the first Go implementation.

The Go rewrite should preserve proven sprite-generation behavior before adding new product features.

## CLI Contract

Spritey commands must be predictable and automation-friendly.

Commands should use explicit inputs and outputs, stable exit codes, concise human logs, and JSON output modes for agent workflows.

Validation should be available before generation.

## Architecture

Spritey follows a Rails-style MVC organization adapted for a Go CLI.

Use clear responsibility boundaries:

- Models: domain data and validation rules.
- Controllers: command handlers and request orchestration.
- Views: user-facing and machine-facing output formatting.
- Services: reusable business operations such as rendering, catalog loading, palette resolution, and report generation.

CLI parsing should stay thin. Command handlers should delegate to controllers and services instead of embedding business logic in `main.go`.

## Specification-Driven Development

Spritey uses Spec Kit-style spec-driven development.

Reference:

```text
https://github.com/github/spec-kit
```

Specs define behavior before implementation. Feature work should proceed through:

1. specification
2. technical plan
3. task breakdown
4. implementation
5. validation

The constitution constrains all feature specs.

## TDD Standard

Implementation tasks should use TDD when practical:

```text
spec task -> failing test -> minimal implementation -> passing test -> refactor -> validation
```

For CLI behavior, schema validation, recipe parsing, catalog output, reports, and rendering logic, write or update the relevant test before implementing behavior.

Docs-only, exploratory, and repo-maintenance tasks may skip TDD when a test would not add useful signal.

## Multi-Agent Workflow

If a user prompts for multi-agent workflow, subagents, parallel agents, background workers, or delegated agent work, the main agent must act as Orchestrator.

Roles:

- Orchestrator: owns scope, sequencing, integration, and final validation.
- Advisor: read-only reviewer for specs, architecture, risks, or implementation plans.
- Worker: implementation agent for one bounded code or documentation task.

Advisor findings must come before Worker implementation. The Orchestrator converts Advisor findings into atomic Worker tasks.

Aim for one spec task per subagent task.

## Testing Standard

Behavior must be verified with automated tests.

Preferred validation includes:

- JSON Schema for pack, recipe, report, and catalog contracts.
- Fixture asset packs.
- Golden JSON outputs.
- Image dimension checks.
- Pixel or hash checks where deterministic.
- CLI command tests.

## Non-Goals For Early Versions

Early versions should avoid:

- GUI/editor work.
- Asset marketplace behavior.
- Automatic publishing.
- Broad plugin systems.
- Large feature expansion before the Go CLI matches the agreed core generation behavior.
