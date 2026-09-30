---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_spoke_group_transparent_inspection"
description: |-
  Enables and manages transparent inspection for an Aviatrix spoke group
---

# aviatrix_spoke_group_transparent_inspection

The **aviatrix_spoke_group_transparent_inspection** resource enables and manages transparent inspection for the VPC of an Aviatrix spoke group. Traffic in the selected private route tables is steered through the group's spoke gateways.

~> **NOTE:** Transparent inspection can only be enabled after the group's spoke instances are launched. Add a `depends_on` on the `aviatrix_spoke_instance` resources of the group so that transparent inspection is enabled after they are created and disabled before they are destroyed.

~> **NOTE:** Transparent inspection manages SNAT on the group's gateways while this resource exists. The spoke group's `enable_nat` value is ignored during that time, and changing it from true to false is rejected.

## Example Usage

```hcl
resource "aviatrix_spoke_group" "aws_spoke_group" {
  group_name          = "my-spoke-group"
  cloud_type          = 1
  account_name        = "my-aws-account"
  gw_type             = "spoke"
  group_instance_size = "t3.medium"
  vpc_id              = "vpc-abcd1234"
  vpc_region          = "us-west-2"
}

resource "aviatrix_spoke_instance" "aws_spoke" {
  group_uuid = aviatrix_spoke_group.aws_spoke_group.group_uuid
  gw_name    = "my-aws-spoke-instance"
  subnet     = "10.0.1.0/24"
  gw_size    = "t3.medium"
}

resource "aviatrix_spoke_group_transparent_inspection" "aws_spoke_group" {
  group_uuid   = aviatrix_spoke_group.aws_spoke_group.group_uuid
  route_tables = ["rtb-0a1b2c3d4e5f67890"]
  egress       = true

  depends_on = [aviatrix_spoke_instance.aws_spoke]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `group_uuid` - (Required) UUID of the spoke group. The group must have `gw_type` SPOKE and `cloud_type` AWS (1). Changing this value requires resource replacement.

### Optional

* `route_tables` - (Optional) Set of private route table IDs to use for transparent inspection. If not set, the controller selects the VPC's private route tables.
* `egress` - (Optional) Enable egress through the spoke gateways for transparent inspection. Valid values: true, false. Default value: false.

## Import

**spoke_group_transparent_inspection** can be imported using the `group_uuid`, e.g.

```
$ terraform import aviatrix_spoke_group_transparent_inspection.test group_uuid
```

## Notes

* Transparent inspection is configured per VPC. A VPC has at most one spoke group, so at most one of these resources exists per VPC.
* Deleting this resource disables transparent inspection for the VPC, which also disables SNAT on the group's gateways.
