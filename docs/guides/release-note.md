---
layout: "aviatrix"
page_title: "Release Note"
description: |-
  The Aviatrix provider Release Note
---

# Aviatrix Provider: Release Note

## 8.1.30
### Notes:
- Supported Controller version: **8.1.30**

### Bug Fixes:
| Issue | Description |
| :--- | :--- |
| AVX-72574 | Fixed **aviatrix_site2cloud** custom mapped NAT CIDR updates (`remote_destination_real_cidrs`, `remote_destination_virtual_cidrs`, and related fields) not taking effect in the controller when modified in-place; changes now propagate correctly without requiring resource recreation |
