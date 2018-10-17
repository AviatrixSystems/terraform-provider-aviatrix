---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_peer"
sidebar_current: "docs-aviatrix-resource-aws-peer"
description: |-
  Creates and manages Aviatrix AWS Peering
---

# aviatrix_aws_peer

The AWSPeer resource allows the creation and management of Aviatrix AWS Peering.

## Example Usage

```hcl
# Create Aviatrix AWS Peering
resource "aviatrix_aws_peer" "test_awspeer" {
  account_name1 = "test1-account"
  account_name2 = "test2-account"
  vpc_id1 = "vpc-abcd1234"
  vpc_id2 = "vpc-rdef3333"
  vpc_reg1 = "us-east-1"
  vpc_reg2 = "us-west-1"
  rtb_list1 = "rtb-abcd1234"
  rtb_list2 = "rtb-wxyz5678"
}
```

## Argument Reference

The following arguments are supported:

* `account_name1` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `account_name2` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `vpc_id1` - (Required) VPC-ID/VNet-Name of cloud provider. Example: AWS: "vpc-abcd1234", ARM: "VNet_Name:Resource_Group_Name", etc...
* `vpc_id2` - (Required) VPC-ID/VNet-Name of cloud provider. Example: AWS: "vpc-abcd1234", ARM: "VNet_Name:Resource_Group_Name", etc...
* `vpc_reg1` - (Required) Region of cloud provider. Example: AWS: "us-east-1", ARM: "East US 2", etc...
* `vpc_reg2` - (Required) Region of cloud provider. Example: AWS: "us-east-1", ARM: "East US 2", etc...
* `rtb_list1` - (Optional) Route table ID. Valid Values: "all" OR "rtb-abcd1234" OR "rtb-abcd1234,rtb-wxyz5678" etc...
* `rtb_list2` - (Optional) Route table ID. Valid Values: "all" OR "rtb-abcd1234" OR "rtb-abcd1234,rtb-wxyz5678" etc...
