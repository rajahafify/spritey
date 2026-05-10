---
title: "Spritey"
description: "Project-local wiki for the Spritey Go CLI rewrite."
created: 2026-05-10
freshness_threshold: 70
---

# Wiki Configuration

## Scope

This wiki covers the Spritey project: a Go CLI for generating animated 2D character spritesheets from recipe files and compatible assets.

The wiki tracks product decisions, CLI contracts, Spec Kit workflow, Rails-style MVC architecture, TDD expectations, multi-agent workflow rules, asset policy, and implementation context.

## Conventions

- Treat ignored local prototype directories such as `python_source/` as outside the repository.
- Treat root-level Go folders as the future implementation home.
- Use `assets` as the user-facing term for compatible asset packs.
- Do not track third-party sprite assets in this repository; assets are user-provided at runtime.
- Prefer spec-backed claims. If implementation has not started, say so explicitly.
