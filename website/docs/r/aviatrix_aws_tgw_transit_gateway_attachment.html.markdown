---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_transit_gateway_attachment"
description: |-
  Manages attachment of transit gateways to the AWS TGW
---

# aviatrix_aws_tgw_transit_gateway_attachment

The **aviatrix_aws_tgw_transit_gateway_attachment** resource manages the attachment of the transit gateway to the AWS TGW.

## Example Usage

```hcl
# Create an Aviatrix AWS TGW Transit Gateway Attachment
resource "aviatrix_aws_tgw_transit_gateway_attachment" "test_transit_gateway_attachment" {
  tgw_name             = "test-tgw"
  region               = "us-east-1"
  vpc_account_name     = "test-account"
  vpc_id               = "vpc-0e2fac2b91c6697b3"
  transit_gateway_name = "transit-gw-1"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `tgw_name` - (Required) Name of the AWS TGW.
* `region` - (Required) AWS Region of the TGW.
* `vpc_account_name` - (Required) The name of the cloud account in the Aviatrix controller, which is associated with the VPC.
* `vpc_id` - (Required) VPC ID of the VPC, where transit gateway is launched.
* `transit_gateway_name` - (Required) Name of the transit gateway to be attached to the AWS TGW.

## Import

**aws_tgw_transit_gateway_attachment** can be imported using the `tgw_name` and `vpc_id`, e.g.

```
$ terraform import aviatrix_aws_tgw_transit_gateway_attachment.test tgw_name~vpc_id
```
