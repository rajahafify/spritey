---
title: "Spritey"
description: "Project-local wiki for the Spritey Go CLI rewrite and its Python reference implementation."
created: 2026-05-10
freshness_threshold: 70
---

# Wiki Configuration

## Scope

This wiki covers the Spritey project: a Go CLI for generating animated 2D character spritesheets from recipe files and compatible assets.

The wiki tracks product decisions, CLI contracts, Spec Kit workflow, Rails-style MVC architecture, TDD expectations, multi-agent workflow rules, and the Python prototype used as the behavioral reference.

## Conventions

- Treat `python_source/` as reference implementation material, not the Go application root.
- Treat root-level Go folders as the future implementation home.
- Use `assets` as the user-facing term for compatible asset packs.
- Prefer spec-backed claims. If implementation has not started, say so explicitly.
