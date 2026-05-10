---
title: "Agent-Friendly CLI"
category: concept
sources:
  - raw/notes/2026-05-10-recipe-validation.md
  - raw/notes/2026-05-10-inspect-layer.md
  - raw/notes/2026-05-10-catalog-foundation.md
  - raw/notes/2026-05-10-spritey-cli-contract.md
  - raw/notes/2026-05-10-agent-rules-and-constitution.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, cli, agents, automation]
aliases: [Spritey CLI, automation-friendly CLI]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "Spritey's CLI is designed for agents through file-based inputs, explicit output paths, validation-before-generation, JSON output, and stable error behavior."
---

# Agent-Friendly CLI

> Spritey's command surface is intended to be predictable for agents and automation scripts, not only comfortable for humans.

The central product contract is that an agent can discover available layers, inspect options, validate a recipe, generate a PNG, and read a machine-readable report without guessing repository internals.

## Core Commands

The intended command family is:

```bash
spritey catalog --assets ./assets --json
spritey inspect layer <layer-id> --assets ./assets --json
spritey validate recipe.json --assets ./assets --json
spritey make recipe.json --assets ./assets --out output/sprite.png --report output/sprite.report.json
```

The first implemented product slice is `spritey catalog --assets <dir> --json`. The second implemented slice is `spritey inspect layer <layer-id> --assets <dir> --json`. The third implemented slice is `spritey validate <recipe-path> --assets <dir> --json`.

## Design Implications

- Recipes should be files, not inline JSON.
- All generated files should use explicit output paths.
- JSON output must remain stable enough for scripts to parse.
- Validation should be available before rendering.
- Error codes and exit codes are part of the public contract.

## See Also

- [[compatible-assets|Compatible Assets]] ([Compatible Assets](compatible-assets.md)) - asset pack contract.
- [[catalog-foundation|Catalog Foundation]] ([Catalog Foundation](catalog-foundation.md)) - implemented catalog command.
- [[inspect-layer|Inspect Layer]] ([Inspect Layer](inspect-layer.md)) - implemented layer inspection command.
- [[recipe-validation|Recipe Validation]] ([Recipe Validation](recipe-validation.md)) - implemented recipe validation command.
- [[spec-driven-development|Spec-Driven Development]] ([Spec-Driven Development](../topics/spec-driven-development.md)) - development workflow that protects the CLI contract.

## Sources

- [Recipe Validation Implementation](../../raw/notes/2026-05-10-recipe-validation.md) - third implemented product slice.
- [Inspect Layer Implementation](../../raw/notes/2026-05-10-inspect-layer.md) - second implemented product slice.
- [Catalog Foundation Implementation](../../raw/notes/2026-05-10-catalog-foundation.md) - first implemented product slice.
- [Spritey CLI Contract](../../raw/notes/2026-05-10-spritey-cli-contract.md) - command and reporting expectations.
- [Agent Rules and Constitution](../../raw/notes/2026-05-10-agent-rules-and-constitution.md) - agent-first requirement.
