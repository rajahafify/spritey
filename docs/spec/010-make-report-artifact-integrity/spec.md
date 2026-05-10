# Feature Spec: Make Report Artifact Integrity

## Goal

Expand make report v1 with additive output artifact-integrity metadata derived from the actual written PNG while preserving make CLI behavior, stdout envelopes, schema version, and exit-code mappings.

## User Stories

- As an agent, I can verify report integrity against the generated PNG using deterministic `sha256` and byte-size metadata.
- As an integration consumer, I can keep parsing the same make JSON stdout envelope keys.
- As an existing report consumer, I can keep using report schema version `"1"` with additive-only changes.

## Functional Requirements

- Keep make command contract unchanged:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]
```

- Keep make stdout JSON envelope unchanged.
- Keep report schema version `"1"`.
- Add report metadata section:
  - `artifacts.output_png.sha256` (lowercase hex SHA-256 of written `--out` PNG file bytes)
  - `artifacts.output_png.bytes` (integer byte size of written `--out` PNG file)
- Metadata values must reflect the actual written `--out` file.
- No behavior changes when `--report` is omitted.

## Error & Exit Requirements

- Any artifact-metadata computation failure returns `RENDER_FAILED`.
- Keep existing make exit-code mapping unchanged.
- Avoid writing partial/invalid report output when artifact metadata fails.

## Testing Requirements

- Service tests:
  - report includes `artifacts.output_png` section;
  - `sha256` and `bytes` match the written output file;
  - artifact-metadata failure returns `RENDER_FAILED` and no report file is written.
- CLI tests:
  - `make --json --report` keeps stdout envelope keys unchanged;
  - report contains artifact metadata matching output PNG.
- Existing regressions continue to pass.

## Out Of Scope

- Report schema version bump.
- Non-additive report key renames/removals.
- CLI envelope key changes.
