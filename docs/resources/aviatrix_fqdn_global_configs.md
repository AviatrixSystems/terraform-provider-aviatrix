---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_fqdn_global_configs"
description: |-
  Manages Aviatrix FQDN Global Configs
---

# aviatrix_fqdn_global_configs

The **aviatrix_fqdn_global_configs** resource manages FQDN global configs.

## Example Usage

```hcl
# Create an Aviatrix FQDN Global Configs
resource "aviatrix_fqdn_global_configs" "test" {
  exception_rule    = false
  network_filtering = "Customize Network Filtering"
  configured_ips    = ["172.16.0.0/12~~RFC-1918", "10.0.0.0/8~~RFC-1918"]
  caching           = false
  exact_match       = true
}
```

## Argument Reference

The following arguments are supported:

* `exception_rule` - (Optional) Allow packets passing through the gateway without an SNI field. Valid values: true, false.
* `network_filtering` - (Optional) Packet destination address ranges be filtered by FQDN. Valid values: "Enable Private Network Filtering", "Disable Private Network Filtering", "Customize Network Filtering".
* `configured_ips` - (Optional) Config IP address. Only support when network_filtering set as "Customize Network Filtering".
* `caching` - (Optional) Cached the resolved IP address from FQDN filter. Valid values: true, false.
* `exact_match` - (Optional) Exact match in FQDN filter.