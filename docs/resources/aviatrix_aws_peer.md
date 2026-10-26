---
subcategory: "Peering"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_peer"
description: |-
  Creates and manages native AWS VPC peerings
---

# aviatrix_aws_peer

The **aviatrix_aws_peer** resource allows the creation and management of Aviatrix-created native AWS intra and inter-region VPC peerings.

## Example Usage

```hcl
# Create an Aviatrix AWS Peering
resource "aviatrix_aws_peer" "test_awspeer" {
  account_name1 = "test1-account"
  account_name2 = "test2-account"
  vpc_id1       = "vpc-abcd1234"
  vpc_id2       = "vpc-rdef3333"
  vpc_reg1      = "us-east-1"
  vpc_reg2      = "us-west-1"
  rtb_list1     = [
    "rtb-abcd1234",
  ]
  rtb_list2     = [
    "rtb-wxyz5678",
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `account_name1` - (Required) Account name of AWS VPC1.
* `account_name2` - (Required) Account name of AWS VPC2.
* `vpc_id1` - (Required) VPC ID of AWS VPC1. Example: AWS: "vpc-abcd1234".
* `vpc_id2` - (Required) VPC ID of AWS VPC2. Example: AWS: "vpc-abcd1234".
* `vpc_reg1` - (Required) Region of AWS VPC1. Example: AWS: "us-east-1".
* `vpc_reg2` - (Required) Region of AWS VPC2. Example: AWS: "us-east-1".
* `rtb_list1` - (Optional) List of Route table IDs of VPC1. Example: ["rtb-abcd1234", "rtb-wxyz5678"].
* `rtb_list2` - (Optional) List of Route table IDs of VPC2. Example: ["rtb-abcd1234", "rtb-wxyz5678"].

~> **NOTE:** For attributes `rtb_list1` and `rtb_list2`, only valid route table IDs with prefix "rtb-" are supported.
Therefore, "all" will no longer be supported as a valid input as of 3.0.2 onward. If an **aviatrix_aws_peer** resource
was created with provider 3.0.1- and any `rtb_list1` or `rtb_list2` was set as ["all"], it will need to be updated to
a list of valid route table IDs.

## Import

**aws_peer** can be imported using the `vpc_id1` and `vpc_id2`, e.g.

```
$ terraform import aviatrix_aws_peer.test vpc_id1~vpc_id2
```
