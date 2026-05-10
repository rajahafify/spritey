---
title: "Python Reference Implementation"
source: "D:/Work/spritey/python_source"
type: repos
ingested: 2026-05-10
tags: [spritey, python-source, reference-implementation, rendering]
summary: "The Python prototype provides the behavioral reference: setup sparse-clones LPC-compatible assets, catalog parses sheet definitions, compositor resolves and alpha-composites layers, recolorer applies palette swaps, and generate.py exposes catalog/make commands."
---

# Python Reference Implementation

The Python prototype may live locally in ignored `python_source/` and is the behavioral reference for the first Go implementation when present. It must not be committed to the Spritey repository.

## Files

- `setup.py`: sparse-clones compatible asset folders from the Universal LPC Spritesheet Character Generator repository into `lpc-assets/`.
- `catalog.py`: scans `sheet_definitions/**/*.json`, skips `meta_*.json`, groups layer entries by `type_name`, and returns catalog summaries.
- `compositor.py`: resolves selected layer IDs to definitions, chooses a body-type path with male fallback, sorts by `zPos`, finds animation PNGs, optionally recolors, and alpha-composites strips.
- `recolorer.py`: loads palette definitions and swaps colors with a small per-channel tolerance using NumPy and Pillow.
- `exporter.py`: writes PNG output.
- `generate.py`: exposes `catalog` and `make` CLI commands around the catalog/compositor/exporter pipeline.

## Behavior To Preserve

- Asset root contains `sheet_definitions/`, `spritesheets/`, and `palette_definitions/`.
- Default animation list is `spellcast`, `thrust`, `walk`, `slash`, `shoot`, and `hurt`.
- Canvas width is `832`.
- Layers are sorted by `zPos`.
- Missing body type currently falls back to `male`.
- Weapon animation aliases include `slash -> attack_slash` and `thrust -> attack_thrust`.
- Output is a transparent PNG spritesheet with animation strips stacked vertically.

## Known Rewrite Improvement

The Python CLI accepts inline JSON, which is brittle on PowerShell. The Go CLI should prefer recipe files.
