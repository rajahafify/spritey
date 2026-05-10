---
title: "Spec-Driven Development"
category: topic
sources:
  - raw/notes/2026-05-10-agent-rules-and-constitution.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, sdd, tdd, spec-kit, workflow]
aliases: [SDD, Spec Kit Workflow]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "Spritey uses Spec Kit-style SDD with practical TDD: constitution, spec, plan, tasks, implementation, and validation."
---

# Spec-Driven Development

> Spritey uses specs to define behavior before implementation, then uses tests where practical to drive implementation.

The project references GitHub Spec Kit:

```text
https://github.com/github/spec-kit
```

## Required Flow

```text
constitution -> feature spec -> technical plan -> tasks -> implementation -> validation
```

## Practical TDD

For CLI behavior, schema validation, recipe parsing, catalog output, reports, and rendering logic, implementation should follow:

```text
spec task -> failing test -> minimal implementation -> passing test -> refactor -> validation
```

Docs-only, exploratory, and repo-maintenance tasks may skip TDD when a test would not add useful signal.

## Multi-Agent Rule

When multi-agent work is requested, the main agent is Orchestrator. Advisor findings happen before Worker implementation, and Worker tasks should be atomic, ideally one spec task per agent task.

## Sources

- [Agent Rules and Constitution](../../raw/notes/2026-05-10-agent-rules-and-constitution.md) - SDD, TDD, and multi-agent rules.
