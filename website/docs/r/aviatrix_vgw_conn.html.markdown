---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vgw_conn"
sidebar_current: "docs-aviatrix-resource-vgw_conn"
description: |-
  Manages the Aviatrix Transit Gateway to VGW Connection
---

# aviatrix_vgw_conn

The AviatrixVGWConn resource manager the Aviatrix Transit Gateway to VGW Connection

## Example Usage

```hcl
# Manage Aviatrix Controller Upgrade process
resource "aviatrix_vgw_conn" "test_vgw_conn" {
  conn_name = " "
  gw_name = " "
  vpc_id = " "
  bgp_vgw_id = " "
  bgp_local_as_num = " "
}
```

## Argument Reference

The following arguments are supported:

* `conn_name` - (Required) The name of for Transit GW to VGW connection connection which is going to be created. Example: "my-connection-vgw-to-tgw"
* `gw_name` - (Required) Name of the Transit Gateway. Example: "my-transit-gw"
* `vpc_id` - (Required) VPC-ID where the Transit Gateway is located. Example: AWS: "vpc-abcd1234"
* `bgp_vgw_id` - (Required)Id of AWS's VGW that is used for this connection. Example: "vgw-abcd1234"
* `bgp_local_as_num` - (Required) BGP Local ASN (Autonomous System Number). Ingeter between 1-65535. Example: "65001"
