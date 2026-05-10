# Technical Plan: Make Animation Strip Output

## Architecture

Use existing MVC boundaries and keep the change localized:

- `app/services/make_service.go`
  - change composition from single-frame render to vertical strip render
  - maintain existing make result/report payload keys and warning propagation
- `app/services/make_service_test.go`
  - add strip dimension/hash assertions, frame-row pixel assertion, and single-frame regression
- `cmd/spritey/main_test.go`
  - verify strip-aware summary height and report animation order alignment with strip rows

## Data Flow

1. Validate recipe and render readiness as today.
2. Load catalog and resolve applied layers as today.
3. Build strip canvas:
   - frame width from pack `canvas_width` fallback logic;
   - frame height remains current frame dimension basis;
   - total height = frame height * required animation count.
4. For each animation ID in existing required-animation order:
   - compose all selected layers in existing z-order into that frame row.
5. Emit make result/report with unchanged key shapes and updated summary/render dimensions.

## Determinism Rules

- Frame row order must match `validation.RequiredAnimationIDs`.
- Layer order per frame remains `z_pos` ascending with `id` tie-break.
- Report `render.animation_ids` remains the same ordered list used for strip rows.
- No envelope or schema key churn.

## TDD Sequence

1. Add/adjust failing tests in:
   - `app/services/make_service_test.go`
   - `cmd/spritey/main_test.go`
2. Implement strip rendering in `app/services/make_service.go`.
3. Re-baseline deterministic hash expectation for two-animation output.
4. Run:

```bash
go test ./...
make native-ci
```
