# Feature Spec: Make Recolor Parity with python_source/recolorer.py

## Goal

Apply Python-equivalent palette recoloring during `spritey make` composition while preserving existing CLI contracts, output envelopes, and exit-code mapping.

## User Stories

- As an operator, layer selections with `palette_variant` recolor matching materials before composition.
- As an integration consumer, I keep parsing unchanged make JSON/report envelopes while rendered pixels match recolor parity behavior.

## Functional Requirements

- Recolor is attempted per contributing layer frame only when:
  - layer has non-empty `recolor_material`;
  - selection has non-empty `palette_variant`.
- Palette file discovery:
  - scan `assets/palette_definitions` recursively;
  - include files matching `{material}_*.json`;
  - process file paths in lexicographic order.
- Variant lookup parity per file:
  - exact key match first;
  - fallback fuzzy key match (case-insensitive): `requested in key` OR `key in requested`.
- Source/base palette fallback variants, in order:
  - `light`, `base`, `beige`, `brown`.
- Pixel replacement rules:
  - compare source RGB to source-palette RGB with per-channel tolerance `<= 2`;
  - recolor only pixels with `alpha > 0`;
  - preserve alpha channel.
- If source/target palette is missing or invalid, recolor must be skipped non-fatally and image remains unchanged.

## Contract Requirements

- Keep make CLI arguments unchanged.
- Keep make JSON/report envelope keys unchanged.
- Keep make exit-code mapping unchanged.
- Recolor changes output pixels only; no new top-level output keys.

## Testing Requirements

- Service tests:
  - recolor applied when `recolor_material + palette_variant` are present;
  - unknown variant is non-fatal and leaves image unchanged;
  - fuzzy variant matching works;
  - tolerance and transparency behavior parity (`<=2`, `alpha > 0` only).
- Recolor helper tests:
  - exact, fuzzy, unknown-variant, tolerance/transparency parity behaviors.
- Existing make regressions remain green.

## Out Of Scope

- Palette policy beyond Python parity fallback behavior.
- CLI/report schema changes.
- Exit mapping changes.
