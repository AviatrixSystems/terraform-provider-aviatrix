---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_vpc_attachment"
sidebar_current: "docs-aviatrix-resource-aws_tgw_vpc_attachment"
description: |-
  Manages attaching or detaching VPCs to aws tgw
---

# aviatrix_aws_tgw

The AviatrixAwsTgwVpcAttachment resource manages attaching or detaching VPCs to AWS TGW

## Example Usage

```hcl
# Manage attaching or detaching VPCs to aws tgw
resource "aviatrix_aws_tgw_vpc_attachment" "test" {
  tgw_name             = "tgwTest"
  region               = "us-east-1"
  security_domain_name = "mySdn"
  vpc_account_name     = "accountTest"
  vpc_id               = "vpc-0e2fac2b91c6697b3"
}
```

## Argument Reference

The following arguments are supported:

* `tgw_name` - (Required) Name of the AWS TGW.
* `region` - (Required) Region of cloud provider(AWS).
* `security_domain_name` - (Required & ForceNew) The name of the security domain, to which the VPC will be attached. If changed, the VPC will be detached from the old domain, and attached to the new domain.
* `vpc_account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller, which is associated with the VPC. 
* `vpc_id` - (Required) This parameter represents the ID of the VPC which is going to be attached to the security domain (name: `security_domain_name`).

## Import

Instance aws_tgw_vpc_attachment can be imported using the tgw_name, security_domain_name and vpc_id, e.g.

```
$ terraform import aviatrix_aws_tgw_vpc_attachment.test tgw_name~security_domain_name~vpc_id
```