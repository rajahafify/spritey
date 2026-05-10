# Feature Spec: Make Report Provenance

## Goal

Expand make report v1 with additive provenance metadata while keeping make command behavior, stdout contracts, and exit-code mappings unchanged.

## User Stories

- As an agent, I can read report v1 and identify which assets pack and body-type resolution were used for a render.
- As an agent, I can read a deterministic composed-layer provenance list with per-layer body-type path resolution metadata.
- As an existing integration, I can keep using `layers.applied` and make JSON envelopes without migration breakage.

## Functional Requirements

- Keep command and stdout contracts unchanged:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]
```

- Report v1 stays schema version `"1"` and is additive-only for this slice.
- Add report metadata:
  - `pack.id`
  - `pack.name`
  - `recipe.body_type_effective`
  - `recipe.body_type_requested`
- Add `layers.composed` array ordered by:
  1. `z_pos` ascending
  2. `id` ascending when `z_pos` is tied
- Every `layers.composed` entry includes:
  - `category`
  - `id`
  - `z_pos`
  - `resolved_body_type`
  - `resolved_path` (slash-normalized)
  - `palette_variant` (empty string when omitted)
  - `credits` (array, never `null`)
- Keep `layers.applied` unchanged for backward compatibility.
- Make JSON envelope written to stdout remains unchanged.

## Error & Exit Requirements

- No exit-code mapping changes from existing make behavior.
- No new error envelope keys.

## Testing Requirements

- Add service tests for:
  - report v1 provenance fields
  - deterministic `layers.composed` ordering
  - fallback-body-type scenario reflected in composed entries
- Add CLI test for:
  - make JSON stdout envelope unchanged
  - report artifact contains provenance fields

## Out Of Scope

- Report schema version bump.
- Non-additive report key renames/removals.
- Rendering algorithm or composition changes.
