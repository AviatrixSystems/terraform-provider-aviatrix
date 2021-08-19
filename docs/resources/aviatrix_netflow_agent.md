---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_netflow_agent"
description: |-
  Enables and disables netflow_agent
---

# aviatrix_netflow_agent

The **aviatrix_netflow_agent** resource allows the enabling and disabling of netflow agent.

## Example Usage

```hcl
# Enable netflow agent
resource "aviatrix_netflow_agent" "test_netflow_agent" {
  server_ip         = "1.2.3.4"
  port              = 10
  version           = 5
  excluded_gateways = [
    "a", 
    "b"
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `server_ip` (Required) Netflow server IP address.
* `port` (Required) Netflow server port.

### Optional
* `version` (Optional) Netflow version (5 or 9). 5 by default. 
* `excluded_gateways` (Optional) List of gateways to be excluded from logging. e.g.: ["gateway01", "gateway02", "gateway01-hagw"].

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - The status of netflow agent.

## Import

**netflow_agent** can be imported using "netflow_agent", e.g.

```
$ terraform import aviatrix_netflow_agent.test netflow_agent
```
