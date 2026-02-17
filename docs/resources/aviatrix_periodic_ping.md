---
subcategory: "Gateway"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_periodic_ping"
description: |-
  Manages periodic pings on Aviatrix gateways
---

# aviatrix_periodic_ping

The **aviatrix_periodic_ping** resource manages the periodic ping feature for Aviatrix gateways.

## Example Usage

```hcl
# Enable Periodic Ping for a Gateway
resource "aviatrix_periodic_ping" "test_ping" {
	gw_name    = "test-gw"
	interval   = 600
	ip_address = "127.0.0.1"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Name of the gateway.
* `interval` - (Required) Interval between pings in seconds.
* `ip_address` - (Required) IP Address to ping.

## Import

**aviatrix_periodic_ping** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_periodic_ping.test gw_name
```
