---
layout: "aviatrix"
page_title: "Release Note"
description: |-
  The Aviatrix provider Release Note
---

# Aviatrix Provider: Release Note

## 9.0.20
### Notes:
- Supported Controller version: **9.0.20**

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-80735 | Fixed **`aviatrix_dcf_ruleset`** in-place rule updates failing with `AVXERR-DFW-0004` (`src_ads not a list`) when a rule change caused a phantom entry with empty `src_smart_groups` to be included in the PUT payload |
