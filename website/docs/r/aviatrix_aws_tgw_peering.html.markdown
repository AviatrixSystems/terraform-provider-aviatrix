---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_peering"
description: |-
  Creates and manages Aviatrix inter region AWS tgw peerings
---

# aviatrix_aws_tgw_peering

The **aviatrix_aws_tgw_peering** resource allows the creation and management of Aviatrix inter region AWS tgw peerings.

## Example Usage

```hcl
# Create an Aviatrix AWS Tgw Peering
resource "aviatrix_aws_tgw_peering" "test" {
  tgw_name1 = "tgw1"
  tgw_name2 = "tgw2"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `tgw_name1` - (Required) This parameter represents name of the first AWS tgw to make a peer pair.
* `tgw_name2` - (Required) This parameter represents name of the second AWS tgw to make a peer pair.

## Import

**aws_tgw_peering** can be imported using the `tgw_name1` and `tgw_name2`, e.g.

```
$ terraform import aviatrix_aws_tgw_peering.test tgw_name1~tgw_name2
```
