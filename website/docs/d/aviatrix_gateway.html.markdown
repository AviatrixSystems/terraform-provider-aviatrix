---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway"
description: |-
  Gets the Aviatrix gateway.
---

# aviatrix_gateway

Use this data source to get the Aviatrix gateway for use in other resources.

## Example Usage

```hcl
# Aviatrix gateway data source
data "aviatrix_gateway" "foo" {
  account_name = "username"
  gw_name      = "gatewayname"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Gateway name. This can be used for getting gateway.
* `account_name` - (Optional) Account name. This can be used for logging in to CloudN console or UserConnect controller.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `account_name` - Aviatrix account name.
* `gw_name` - Aviatrix gateway name.
* `cloud_type` - Type of cloud service provider.
* `vpc_id` - VPC ID.
* `vpc_reg` - VPC Region.
* `vpc_size` - Instance type.
* `public_ip` - Public IP address of the Gateway created.
