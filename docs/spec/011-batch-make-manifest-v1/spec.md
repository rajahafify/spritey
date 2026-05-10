# Feature Spec: Batch Make Manifest v1

## Goal

Add `make batch` to execute multiple make jobs from a manifest file while preserving existing single-make behavior and stable machine-readable output envelopes.

## User Stories

- As an agent operator, I can execute many sprite builds from one manifest file.
- As an integration consumer, I can parse a stable batch JSON envelope for success and failure.
- As a CLI user, I can use concise deterministic text output for quick batch status checks.

## Functional Requirements

- Add command:

```bash
spritey make batch <manifest.json> --assets <dir> [--json]
```

- Keep existing single-make command and behavior unchanged:

```bash
spritey make <recipe> --assets <dir> --out <png> [--report <json>] [--json]
```

- Manifest v1 schema:

```json
{"schema_version":"1","jobs":[{"id":"...","recipe":"...","out":"...","report":"..."}]}
```

- Execution rules:
  - run jobs sequentially in manifest order;
  - resolve relative `recipe`, `out`, and `report` paths from manifest file directory;
  - reuse existing `MakeService` for each job.

- Batch JSON success envelope must keep top-level keys:
  - `ok`, `command`, `summary`, `jobs`, `warnings`, `errors`
- `command` must be `"make-batch"`.
- `summary` must include:
  - `job_count`, `success_count`, `failed_count`
- `jobs` entries must include:
  - `id`, `recipe`, `outputs`, `summary`, `warnings`, `errors`

- Non-JSON success output must be deterministic and concise.

## Error & Exit Requirements

- Exit `2` for CLI misuse:
  - missing manifest argument;
  - missing `--assets`;
  - unknown arguments.
- Exit `4` for manifest issues:
  - manifest file read/not found issues;
  - invalid manifest JSON;
  - unsupported manifest schema version;
  - empty jobs list;
  - missing required per-job fields.
- Fail fast on first failed job.
- Job failure exit mapping must follow existing make behavior:
  - invalid-assets/input-validation class failures -> exit `3`
  - render failures -> exit `6`
- JSON failure envelope stays stable and includes failing-job context.

## Testing Requirements

- Service tests:
  - two-job success;
  - manifest-relative path resolution;
  - fail-fast behavior.
- CLI tests:
  - JSON envelope stability and manifest-order jobs;
  - misuse and manifest-error exits;
  - job-failure exit mapping (`3` and `6`);
  - deterministic non-JSON summary.
- Existing make tests must remain green.

## Out Of Scope

- Parallel batch execution.
- Manifest schema version bump.
- Changes to single-make CLI contract.
