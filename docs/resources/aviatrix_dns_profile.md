---
subcategory: "Copilot"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dns_profile"
description: |-
  Creates Aviatrix DNS Profile
---

# aviatrix_dns_profile

The **aviatrix_dns_profile** resource creates the Aviatrix DNS profile.

## Example Usage

```hcl
# Create a DNS profile
resource "aviatrix_dns_profile" "test" {
  name               = "profileA"
  global_dns_servers = ["8.8.8.8", "8.8.3.4"]
  local_domain_names = ["avx.internal.com", "avx.media.com"]
  lan_dns_servers    = ["1.2.3.4", "5.6.7.8"]
  wan_dns_servers    = ["2.3.4.5", "6.7.8.9"]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) DNS profile name.

### Optional
* `global_dns_servers` - (Optional) List of global DNS servers. Example: ["8.8.8.8", "8.8.3.4"].
* `local_domain_names` - (Optional) List of local domain names. Example: ["avx.internal.com", "avx.media.com"].
* `lan_dns_servers` - (Optional) List of LAN DNS servers. Example: ["1.2.3.4", "5.6.7.8"].
* `wan_dns_servers` - (Optional) List of WAN DNS servers. Example: ["2.3.4.5", "6.7.8.9"].

## Import

**dns_profile** can be imported using the `name`, e.g.

```
$ terraform import aviatrix_dns_profile.test name
```
