# Feature Spec: Make Frame-Readiness Parity with Python Compositor

## Goal

Align `spritey make` with Python compositor behavior for missing per-animation frames: compose available rows non-fatally while keeping `validate` strict.

## User Stories

- As an operator, I can run `spritey make` even when some layer animation frames are missing, and still get a composed PNG from available rows.
- As an integration consumer, I keep the same make JSON/report envelope keys and exit mapping while rendered output follows Python-style row emission.
- As a quality gate owner, `spritey validate` still fails on missing required frames.

## Functional Requirements

- `make` must not fail when a selected layer is missing frame files for one or more animations.
- `make` composes rows only for animations where at least one selected layer frame resolves.
- If zero rows emit, `make` succeeds and writes a transparent `832x256` PNG.
- `validate` remains strict:
  - missing required frame still returns `MISSING_SPRITE_FRAME`;
  - validate command behavior and exit mapping are unchanged.

## Contract Requirements

- Keep make CLI argument surface unchanged.
- Keep make JSON/report envelope keys unchanged.
- Keep make exit mapping unchanged.
- Missing-frame behavior changes for `make` only (from fatal to non-fatal row skip).

## Testing Requirements

- Service tests:
  - missing frame on one layer still succeeds when another layer contributes row pixels;
  - missing mapped slash frame still succeeds when body slash frame exists;
  - zero emitted rows still returns transparent `832x256`.
- CLI tests:
  - `make --json` succeeds for missing layer frames and produces output summary for emitted rows;
  - slash mapped-path missing scenario succeeds when another layer contributes;
  - validate missing-frame tests remain failing/strict.

## Out Of Scope

- New CLI flags, schema keys, or report envelope changes.
- Any rendering algorithm changes outside frame-readiness parity behavior.
