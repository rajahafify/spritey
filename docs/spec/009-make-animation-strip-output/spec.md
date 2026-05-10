# Feature Spec: Make Animation Strip Output

## Goal

Update make rendering so output PNG includes one composed frame per required animation ID in a vertical strip while preserving all make CLI contracts, envelopes, schema keys, and exit behavior.

## User Stories

- As an agent, I can render one PNG that contains every required animation frame in deterministic row order.
- As an integration consumer, I can keep parsing the same make JSON envelope keys and report v1 keys without migration.
- As a validator consumer, I can keep relying on existing warnings and exit code behavior.

## Functional Requirements

- Keep make command contract unchanged:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]
```

- Render one frame per required animation ID in existing required-animation order.
- For each frame, compose selected layers using existing z-order rules (`z_pos` ascending, `id` ascending tie-break).
- Output PNG dimensions must be:
  - width = frame width
  - height = frame height * frame count
- Keep make JSON envelope keys unchanged; update summary values to reflect strip output dimensions.
- Keep report schema version and key structure unchanged, including spec 008 provenance fields.
- Keep report `render.frame_count` and `render.animation_ids` aligned to strip row order.

## Error & Exit Requirements

- Preserve warning propagation behavior.
- Preserve make exit code mapping and failure behavior.
- No new make JSON envelope keys.

## Testing Requirements

- Service tests:
  - two-animation output dimensions are `8x16` with deterministic hash baseline;
  - two strip rows differ by pixel assertion;
  - single-animation output remains `8x8`.
- CLI tests:
  - make `--json` summary `canvas.height` matches strip height and envelope keys remain unchanged;
  - report `animation_ids` order aligns with strip row order.
- Existing make regressions continue to pass.

## Out Of Scope

- CLI argument changes.
- Report schema version bump.
- Animation ordering policy changes in validation.
