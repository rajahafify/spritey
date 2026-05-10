# Technical Plan: Download LPC Assets

## Architecture

- Add a top-level CLI route in `cmd/spritey/main.go` for `--download-lpc-assets`.
- Keep parser thin by adding a small options parser for `--json` and `--force`.
- Add MVC slice:
  - Model: download-assets result payload
  - Controller: command orchestration, exit mapping, JSON/text output routing
  - Service: download/extract/install/validate pipeline with injected fetch and cache-dir resolvers
  - View: stable JSON envelope writer
- Validate staged assets via existing `AssetsValidator` before promoting install directory.

## TDD Sequence

1. Add failing CLI tests in `cmd/spritey/main_test.go`:
   - success JSON
   - unknown argument mapping
   - operational download failure mapping
   - invalid downloaded assets mapping
2. Add failing service tests for:
   - default path resolution
   - successful install behavior
   - force reinstall behavior
3. Implement service + controller + view + model + main routing.
4. Update usage/wiki/spec docs for new command contract.
5. Validate:

```bash
go test ./...
make native-ci
```
