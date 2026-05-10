# Feature Spec: Path Resolver Parity with python_source/compositor.py

## Goal

Match frame path resolution behavior used by the legacy Python compositor for animation frame lookup, without changing existing make/validate contracts.

## User Stories

- As an operator using existing LPC-style assets, I can run `validate` and `make` even when frame files are laid out as flat files, `fg` files, animation subfolders, or mapped weapon folders.
- As an integration consumer, I keep existing error class and exit mappings while frame lookup becomes more compatible.

## Functional Requirements

- Introduce shared resolver behavior for required frame lookup in this exact order:

1. `{prefix}/{anim}.png`
2. `{prefix}/fg/{anim}.png`
3. `{prefix}/{anim}/*.png` (first hit)
4. only for `slash` and `thrust`, map to `attack_slash` / `attack_thrust` then scan `{prefix}/{mapped}/*.png` excluding filenames containing lowercase `behind`

- Precedence rule: structure C must win over structure D for `slash`/`thrust` when both are present.
- Use the same resolver logic in:
  - validation readiness checks (`MISSING_SPRITE_FRAME` class);
  - make render frame loading path.

## Error & Contract Requirements

- Keep existing make/validate response contracts.
- Keep existing CLI exit mappings unchanged.
- Missing frame resolution in readiness must still report as `MISSING_SPRITE_FRAME`.
- This slice must not include recolor parity or layout/render-dimension parity changes.

## Testing Requirements

- Add resolver unit tests covering:
  - A/B/C/D structures;
  - D `behind` filter behavior;
  - C-over-D precedence;
  - not-found behavior.
- Add service and CLI tests proving:
  - make succeeds when only mapped D path exists for slash;
  - make fails with `MISSING_SPRITE_FRAME` when no resolver path exists.
- Keep full suite green.

## Out Of Scope

- Recolor behavior parity with Python prototype.
- Layout, padding, strip-height, or compositing behavior parity beyond path resolution.
