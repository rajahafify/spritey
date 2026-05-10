# 2026-05-10 Make Animation Strip Output

Date: 2026-05-10
Scope: spec 009 (`docs/spec/009-make-animation-strip-output/`)

Updated make rendering from single-frame composition to vertical strip composition across required animation IDs, while preserving CLI/report contracts and exit behavior.

Changes:

- Render now composes one frame per required animation ID using existing order.
- Each frame layers selected layers using existing z-order (`z_pos` asc, `id` asc tie-break).
- Output PNG now uses vertical strip dimensions:
  - width = frame width (current pack canvas width)
  - height = frame height * frame count
- Make JSON envelope keys remain unchanged, with summary canvas height updated to strip height.
- Report schema/version/keys remain unchanged (including spec 008 provenance fields), and `render.frame_count` plus `render.animation_ids` align with strip row order.
- Warning propagation and existing exit behavior remain unchanged.

Tests:

- Service:
  - deterministic multi-animation strip `8x16` + re-baselined hash;
  - strip rows differ via pixel assertion;
  - single-animation regression remains `8x8`.
- CLI:
  - `--json` summary canvas height reflects strip height with unchanged envelope keys;
  - report `animation_ids` order aligns with strip row order.
