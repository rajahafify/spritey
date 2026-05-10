# Feature Spec: Make Animation/Layout Parity with python_source/compositor.py

## Goal

Align `spritey make` animation row emission and canvas layout behavior with the Python compositor for this slice, without changing CLI surface, envelope keys, or exit-code mapping.

## User Stories

- As an operator using LPC-compatible assets, I get animation rows emitted in Python LPC order regardless of pack animation defaults.
- As an integration consumer, I keep parsing the same make JSON/report envelopes while summary/render values reflect emitted rows.

## Functional Requirements

- Render animations in this fixed order:
  - `spellcast`, `thrust`, `walk`, `slash`, `shoot`, `hurt`
- Do not derive row order from `pack.defaults.animations`.
- Per animation row:
  - include only layers that resolve a frame for that animation;
  - missing frame for a layer is skipped at compose time (non-fatal for compose).
- Emit a row only when at least one layer contributed.
- Row compositing parity:
  - first contributing layer defines row height;
  - subsequent contributing layers are padded/clipped to that row height before alpha composite.
- Strip width is fixed to `832`.
- Final canvas stacks emitted rows vertically; height is the sum of emitted row heights.
- If no rows are emitted, output a transparent `832x256` PNG.

## Contract Requirements

- Keep make CLI arguments unchanged.
- Keep make JSON/report envelope keys unchanged.
- Update only summary/report values (`frame_count`, `animation_count`, canvas dimensions, `render.animation_ids`) to match emitted rows.
- Keep make exit-code mapping unchanged.

## Testing Requirements

- Service tests:
  - fixed LPC ordering behavior with emitted-row filtering;
  - row-height parity from first contributing layer;
  - fixed width `832`;
  - blank `832x256` fallback when no LPC rows emit;
  - report `render.animation_ids` and `render.frame_count` reflect emitted rows.
- CLI tests:
  - envelope keys unchanged while summary values shift to emitted-row layout behavior.
- Existing regressions continue to pass.

## Out Of Scope

- Recolor parity.
- CLI shape changes.
- Exit mapping changes.
