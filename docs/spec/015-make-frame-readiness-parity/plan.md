# Technical Plan: Make Frame-Readiness Parity with Python Compositor

## Architecture

- Split recipe validation frame-readiness behavior by caller intent:
  - `Validate(...)` remains strict for `validate` command and existing strict paths;
  - add `ValidateForMake(...)` non-strict mode that keeps metadata/path validation but skips required-frame preflight failures.
- Update `MakeService` to use non-strict validation path and continue composing available rows.
- Preserve all make/report envelope shapes and exit mapping.

## TDD Sequence

1. Update failing make service tests for non-fatal missing-frame behavior.
2. Update failing CLI make tests for success path when frames are partially missing.
3. Implement validator split (`Validate` strict, `ValidateForMake` non-strict).
4. Wire `MakeService` to `ValidateForMake`.
5. Validate:

```bash
go test ./...
make native-ci
```
