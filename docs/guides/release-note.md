---
layout: "aviatrix"
page_title: "Release Note"
description: |-
  The Aviatrix provider Release Note
---

# Aviatrix Provider: Release Note

## 10.1.10
### Notes:
- Supported Controller version: **10.1.10**

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-82151 | Fixed **`aviatrix_dcf_ruleset`** showing the entire policy list as changed on every `terraform plan`/apply even when no rules were modified. The hash function now normalizes `watch` and `enforcement` fields before hashing, preventing spurious diffs. |
