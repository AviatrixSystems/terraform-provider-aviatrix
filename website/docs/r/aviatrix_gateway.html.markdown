---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway"
sidebar_current: "docs-aviatrix-resource-gateway"
description: |-
  Creates and manages an Aviatrix gateway.
---

# aviatrix_gateway

The Account resource allows the creation and management of an Aviatrix gateway.

## Example Usage

```hcl
# Create Aviatrix AWS gateway
resource "aviatrix_gateway" "test_gateway1" {
  cloud_type = 1
  account_name = "devops"
  gw_name = "avtxgw1"
  vpc_id = "vpc-abcdef"
  vpc_reg = "us-west-1"
  vpc_size = "t2.micro"
  vpc_net = "10.0.0.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Type of cloud service provider. (Only AWS is supported currently. Enter 1 for AWS.)
* `account_name` - (Required) Account name. This account will be used to launch Aviatrix gateway.
* `gw_name` - (Required) Aviatrix gateway unique name.
* `vpc_id` - (Required) ID of legacy VPC/Vnet to be connected.
* `vpc_reg` - (Required) AWS region where this gateway will be launched.
* `vpc_size` - (Required) Size of Gateway Instance.
* `vpc_net` - (Required) A VPC Network address range selected from one of the available network ranges.
* `saml_enabled` - (Optional) Enables Gateway SAML support.
