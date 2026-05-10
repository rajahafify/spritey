# Technical Plan: Make Human Output

## Architecture

Preserve MVC layering and keep CLI parsing thin:

- `cmd/spritey/main.go`: no contract change; existing `--json` option remains optional.
- `app/controllers/make_controller.go`: branch to JSON view or text view.
- `app/views`: add dedicated non-JSON writer for deterministic success output lines.

## Data Flow

1. CLI parses `make` options.
2. Controller validates required arguments and executes existing make service.
3. On success:
   - JSON mode: existing `WriteMakeJSON`.
   - Non-JSON mode: new text writer prints deterministic lines.
4. On failure:
   - JSON mode: existing structured error JSON.
   - Non-JSON mode: existing `<CODE>: <message>` stderr output.

## Output Contract (Non-JSON Success)

Exact field order:

1. `ok: make`
2. `png: <path>`
3. `report: <path>` (if present)
4. `frame_count: <n>`
5. `canvas: <width>x<height>`
6. `animation_count: <n>`

## Validation

TDD sequence:

1. Add failing CLI tests for non-JSON success/error behavior.
2. Implement text view + controller wiring.
3. Run:

```bash
go test ./...
make native-ci
```
