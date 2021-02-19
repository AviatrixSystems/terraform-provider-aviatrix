---
subcategory: "TGW Orchestrator"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_connect_peer"
description: |- Creates and manages Aviatrix AWS TGW Connect peers
---

# aviatrix_aws_tgw_connect_peer

The **aviatrix_aws_tgw_connect_peer** resource allows the creation and management of AWS TGW Connect peers. This
resource is available as of provider version R2.18.1+.

## Example Usage

```hcl
# Create an Aviatrix AWS TGW Connect Peer
resource "aviatrix_aws_tgw_connect_peer" "test" {
  tgw_name              = aviatrix_aws_tgw.test_aws_tgw.tgw_name
  connection_name       = aviatrix_aws_tgw_connect.test_aws_tgw_connect.connection_name
  connect_peer_name     = "connect-peer-test"
  connect_attachment_id = aviatrix_aws_tgw_connect.test_aws_tgw_connect.connect_attachment_id
  peer_as_number        = "65001"
  peer_gre_address      = "172.31.1.11"
  bgp_inside_cidrs      = ["169.254.6.0/29"]
  tgw_gre_address       = "10.0.0.32"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `tgw_name` - (Required) AWS TGW name.
* `connection_name` - (Required) TGW Connect connection name.
* `connect_peer_name` - (Required) TGW Connect peer name.
* `connect_attachment_id` - (Required) Connect Attachment ID.
* `peer_as_number` - (Required) Peer AS Number.
* `peer_gre_address` - (Required) Peer GRE IP Address.
* `bgp_inside_cidrs` - (Required) Set of BGP Inside CIDR Block(s).

* `tgw_gre_address` - (Optional) AWS TGW GRE IP Address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `connect_peer_id` - Connect Peer ID.

## Import

**aws_tgw_connect_peer** can be imported using the `tgw_name`, `connection_name` and `connect_peer_name`, e.g.

```
$ terraform import aviatrix_aws_tgw_connect_peer.test tgw_name~~connection_name~~connect_peer_name
```
