---
title: "Spritey Repository Baseline"
source: "D:/Work/spritey"
type: repos
ingested: 2026-05-10
tags: [spritey, repository, baseline, go-cli]
summary: "Snapshot of the Spritey repository after initialization: root Go app skeleton, Spec Kit constitution, README, AGENTS guidance, local wiki, and Docker Compose development setup."
---

# Spritey Repository Baseline

Repository path: `D:/Work/spritey`

Remote repository: `git@github.com:rajahafify/spritey.git`

Current branch: `main`

Latest committed baseline: `7425b40 Initial Spritey scaffold`

Current uncommitted state at ingestion: `README.md` has local edits adding the GitHub repository URL and build-from-source instructions.

## Root Layout

```text
spritey/
  app/
    controllers/
    models/
    services/
    views/
  cmd/
    spritey/
  config/
  docs/
    spec/
      constitution.md
  schemas/
  testdata/
  AGENTS.md
  README.md
  .gitignore
```

The Go implementation is intended to live at the repository root using the Rails-style MVC folder structure documented in `AGENTS.md` and `docs/spec/constitution.md`.

## Current Implementation State

The root Go implementation is scaffold-only. There is no `go.mod`, `main.go`, command parser, schema code, or renderer yet.

Ignored local prototype directories may exist on a developer machine, but they are not repository content.

## Asset Policy

Spritey must not commit third-party sprite assets. Compatible assets are downloaded, installed, or provided by users at runtime.

## Docker Development Update

After the initial wiki ingest, the repository gained a Docker Compose development setup:

- `go.mod` declares module `github.com/rajahafify/spritey`.
- `cmd/spritey/main.go` provides a minimal CLI scaffold so the module builds.
- `Dockerfile` uses `golang:1.23-bookworm`.
- `compose.yaml` defines a `dev` service with Go build and module caches.
- `Makefile` provides `docker-build`, `docker-test`, `docker-run`, `docker-shell`, `docker-fmt`, and `docker-clean`.

Normal development is Docker-based. Native Go remains optional for faster local feedback.
