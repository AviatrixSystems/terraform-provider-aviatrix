---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_data_gateway"
sidebar_current: "docs-aviatrix-data_source-gateway"
description: |-
  Gets the Aviatrix gateway.
---

# aviatrix_gateway

Use this data source to get the Aviatrix gateway for use in other resources.

## Example Usage

```hcl
# Create Aviatrix gateway data source
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

* `account_name` - Aviatrix account name.
* `gw_name` - Aviatrix gateway name.


## The following arguments are computed - please do not edit in the resource file:

* `cloud_type` - Type of cloud service provider. (Only AWS is supported currently. Value of 1 for AWS.)
* `vpc_id` - AWS VPC ID.
* `vpc_reg` - AWS VPC Region. 
* `vpc_size` - Instance type.
* `public_ip` - Public IP address of the Gateway created
 