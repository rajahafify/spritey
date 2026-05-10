# Feature Spec: Render-Input Readiness

## Goal

Add render-input readiness checks to recipe validation so missing spritesheet frames fail early and consistently before rendering.

## User Stories

- As an agent, I can run `spritey validate` and know whether selected layers are render-ready for required animations.
- As an agent, I can run `spritey make` and receive input-validation failures (not render failures) when required source frames are missing.
- As an operator, I can see deterministic warnings when body-type fallback paths were used.

## Functional Requirements

- Validation resolves effective body type as:
  - recipe `body_type`, else
  - `pack.defaults.body_type`.
- For each selected layer:
  - resolve body-type path for effective body type;
  - if missing for requested body type, use `pack.defaults.missing_body_type_fallback` only when that fallback path exists for the same layer.
- Required animation IDs come from `pack.defaults.animations`.
  - If empty (or effectively empty), fallback to `["idle"]`.
- For every selected layer and every required animation ID, validate frame existence at:
  - `assets/spritesheets/<resolved-path>/<animation>.png`
- When requested body type was missing and fallback path was used, emit deterministic warning text.
- Validate JSON response keeps top-level keys unchanged:
  - `ok`, `recipe`, `warnings`, `errors`.
- Make command propagates validation warnings into:
  - JSON `warnings`
  - report `warnings`.

## Error & Exit Requirements

- Missing required frame returns problem code:
  - `MISSING_SPRITE_FRAME`
- For `validate` command:
  - `MISSING_SPRITE_FRAME` maps to validation-failed exit `5`.
- For `make` command:
  - `MISSING_SPRITE_FRAME` maps to input-validation exit `3`.
  - no output PNG is emitted on this failure path.

## Testing Requirements

- Add service and CLI tests for:
  - validate success with complete frames
  - validate missing frame -> `MISSING_SPRITE_FRAME`, exit `5`
  - validate fallback-path success with warnings
  - make missing frame -> exit `3`, `MISSING_SPRITE_FRAME`, no output PNG
  - make fallback success with warnings in JSON and report
- Keep existing JSON envelope top-level keys unchanged.

## Out Of Scope

- New CLI commands.
- New top-level JSON keys.
- Rendering algorithm changes beyond readiness propagation.
