---
title: "Python Reference"
category: reference
sources:
  - raw/repos/2026-05-10-python-reference-implementation.md
created: 2026-05-10
updated: 2026-05-10
tags: [spritey, python-source, rendering, reference]
aliases: [Python Prototype, Reference Implementation]
confidence: high
volatility: warm
verified: 2026-05-10
compiled-from: sources
summary: "The Python source is the behavioral reference for the Go rewrite, covering catalog loading, palette recoloring, layer resolution, compositing, and PNG export."
---

# Python Reference

> The first Go implementation should match the Python prototype's proven sprite-generation behavior before expanding product scope.

The reference implementation may exist locally in ignored `python_source/`. It is not repo-owned content.

## Pipeline

1. `setup.py` prepares compatible source assets.
2. `catalog.py` scans sheet definition JSON files.
3. `generate.py` loads a recipe-like config and calls the compositor.
4. `compositor.py` resolves selected layers, sorts by `zPos`, locates animation PNGs, recolors when needed, and alpha-composites strips.
5. `exporter.py` writes the final PNG.

## Preserve First

The Go rewrite should preserve:

- layer definition scanning
- body-type path resolution
- `zPos` ordering
- palette recoloring
- animation PNG discovery
- transparent PNG export

## Improve During Rewrite

The Python command surface uses inline JSON. The Go version should use recipe files as the primary workflow because agents and PowerShell handle file-based input more reliably.

## Sources

- [Python Reference Implementation](../../raw/repos/2026-05-10-python-reference-implementation.md) - source file roles and behavior.
