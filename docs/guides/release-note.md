---
layout: "aviatrix"
page_title: "Release Note"
description: |-
  The Aviatrix provider Release Note
---

# Aviatrix Provider: Release Note

## 8.2.20
### Notes:
- Supported Controller version: **8.2.20**

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-72968 | Fixed an issue where **`aviatrix_dcf_ruleset`** recalculated all rules in the diff when a single rule was added or removed |
| AVX-76197 | Fixed an issue where updating a rule attribute in **`aviatrix_dcf_ruleset`** caused the rule UUID to change unexpectedly |
| AVX-76494 | Fixed an issue where **`aviatrix_dcf_ruleset`** allowed creation of multiple rules with the same priority |
