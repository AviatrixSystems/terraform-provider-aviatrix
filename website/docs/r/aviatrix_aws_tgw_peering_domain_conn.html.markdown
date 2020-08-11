---
subcategory: "TGW Orchestrator"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_peering_domain_conn"
description: |-
  Creates and manages Aviatrix domain connections between peered AWS TGWs
---

# aviatrix_aws_tgw_peering_domain_conn

The **aviatrix_aws_tgw_peering_domain_conn** resource allows the creation and management of Aviatrix domain connections between peered AWS TGWs.

## Example Usage

```hcl
# Create an Aviatrix Domian Connection between Peered AWS Tgws
resource "aviatrix_aws_tgw_peering_domain_conn" "test" {
  tgw_name1    = "tgw1"
  domain_name1 = "Default_Domain"
  tgw_name2    = "tgw2"
  domain_name2 = "Default_Domain"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `tgw_name1` - (Required) The AWS TGW name of the source domain to make a connection.
* `domain_name1` - (Required) The name of the source domain to make a connection.
* `tgw_name2` - (Required) The AWS TGW name of the destination domain to make a connection.
* `domain_name2` - (Required) The name of the destination domain to make a connection.

## Import

**aws_tgw_peering_domain_conn** can be imported using the `tgw_name1`, `domain_name1`, `tgw_name2` and `domain_name2`, e.g.

```
$ terraform import aviatrix_aws_tgw_peering_domain_conn.test tgw_name1:domain_name1~tgw_name2:domain_name2
```
