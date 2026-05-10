# 2026-05-10 Make Report Provenance Implementation

Date: 2026-05-10
Scope: spec 008 (`docs/spec/008-make-report-provenance/`)

Expanded make report v1 with additive provenance fields while preserving make CLI behavior, stdout envelopes, and exit code mapping.

Changes:

- Added additive report pack metadata:
  - `pack.id`
  - `pack.name`
- Added additive report recipe provenance metadata:
  - `recipe.body_type_requested`
  - `recipe.body_type_effective`
- Added additive `layers.composed` provenance array ordered by render order (`z_pos` asc, `id` asc) with:
  - `category`
  - `id`
  - `z_pos`
  - `resolved_body_type`
  - `resolved_path` (slash-normalized)
  - `palette_variant` (empty string when omitted)
  - `credits` (always array, never null)
- Kept existing `layers.applied` for backward compatibility.
- Added internal validation provenance carry-through for requested body type.
- Added/updated tests in:
  - `app/services/make_service_test.go`
  - `cmd/spritey/main_test.go`

Validation:

- `go test ./...`
- `make native-ci`
