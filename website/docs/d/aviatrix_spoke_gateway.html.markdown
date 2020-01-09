---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_gateway"
description: |-
  Gets the Aviatrix spoke gateway.
---

# aviatrix_spoke_gateway

Use this data source to get the Aviatrix spoke gateway for use in other resources.

## Example Usage

```hcl
# Aviatrix Spoke Gateway Data Source
data "aviatrix_spoke_gateway" "foo" {
  gw_name      = "gatewayname"
  account_name = "username"
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Spoke gateway name. This can be used for getting spoke gateway.
* `account_name` - (Optional) Account name. This can be used for logging in to CloudN console or UserConnect controller.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `gw_name` - Aviatrix spoke gateway name.
* `account_name` - Aviatrix account name.
* `cloud_type` - Type of cloud service provider.
* `vpc_id` - VPC ID.
* `vpc_reg` - VPC Region.
* `gw_size` - Instance type.
* `subnet` - Range of the subnet where the spoke gateway is launched.
* `public_ip` - Public IP address of the spoke gateway created.

