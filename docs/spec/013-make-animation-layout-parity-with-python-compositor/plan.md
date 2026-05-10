# Technical Plan: Make Animation/Layout Parity with python_source/compositor.py

## Architecture

- Update `app/services/make_service.go` render loop to:
  - iterate fixed LPC animation order;
  - resolve per-layer frame per animation and skip missing non-fatally during compose;
  - compose only contributed rows;
  - use first contributed layer height as row height;
  - pad/clip subsequent layers to row height before compositing;
  - emit final canvas at fixed width `832`;
  - fallback to transparent `832x256` when no rows emit.
- Preserve CLI/report envelope shapes and make exit-code mapping.
- Update summary/report value population to emitted-row-driven frame count, animation IDs, and canvas dimensions.

## TDD Sequence

1. Add/adjust make service tests for LPC row ordering, row filtering, row-height parity, fixed width, fallback, and report emitted-row metadata.
2. Adjust CLI tests to assert unchanged envelope keys plus shifted summary/report values.
3. Implement make-service row emission/composition changes.
4. Run validation:

```bash
go test ./...
make native-ci
```
