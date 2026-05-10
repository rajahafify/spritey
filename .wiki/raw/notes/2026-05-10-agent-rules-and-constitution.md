---
title: "Agent Rules and Constitution"
source: "D:/Work/spritey/AGENTS.md; D:/Work/spritey/docs/spec/constitution.md"
type: notes
ingested: 2026-05-10
tags: [spritey, agents, constitution, sdd, tdd, multi-agent]
summary: "Project rules define Spritey as an agent-first Go CLI built through Spec Kit-style SDD with practical TDD, Rails-style MVC, and enforced multi-agent orchestration when requested."
---

# Agent Rules and Constitution

Spritey is defined as a Go CLI for generating animated 2D character spritesheets from recipe files and compatible assets.

The command name is `spritey`.

The Go implementation lives at the repository root. Legacy prototypes may exist locally in ignored directories, but they are not part of the repository.

## Development Workflow

Spritey uses Spec Kit-style spec-driven development.

Reference: `https://github.com/github/spec-kit`

Feature work should follow:

```text
constitution -> feature spec -> technical plan -> tasks -> implementation -> validation
```

Implementation tasks should use TDD when practical:

```text
spec task -> failing test -> minimal implementation -> passing test -> refactor -> validation
```

## Multi-Agent Workflow

If the user prompts for multi-agent workflow, subagents, parallel agents, background workers, or delegated agent work, the main agent must act as Orchestrator.

Roles:

- Orchestrator: owns scope, sequencing, integration, and final validation.
- Advisor: read-only reviewer for specs, architecture, risks, or implementation plans.
- Worker: implementation agent for one bounded code or documentation task.

Advisor findings must come before Worker implementation. The Orchestrator converts Advisor findings into atomic Worker tasks.

Aim for one spec task per subagent task.
