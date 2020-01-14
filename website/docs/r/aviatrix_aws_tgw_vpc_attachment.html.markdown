---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_vpc_attachment"
description: |-
  Manages attaching/detaching VPC to/from an AWS TGW, and FireNet Gateway to TGW Firewall Domain
---

# aviatrix_aws_tgw_vpc_attachment

The aviatrix_aws_tgw_vpc_attachment resource manages attaching/detaching VPC to/from an AWS TGW, and FireNet Gateway to TGW Firewall Domain.

## Example Usage

```hcl
# Create an Aviatrix AWS TGW VPC Attachment 
resource "aviatrix_aws_tgw_vpc_attachment" "test_aws_tgw_vpc_attachment" {
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
* `subnets` - (Optional) Advanced option. VPC subnets separated by ',' to attach to the VPC. If left blank, Aviatrix Controller automatically selects a subnet representing each AZ for the VPC attachment. Example: "subnet-214f5646,subnet-085e8c81a89d70846".
* `route_tables` - (Optional) Advanced option. Route tables separated by ',' to participate in TGW Orchestrator, i.e., learned routes will be propagated to these route tables. Example: "rtb-212ff547,rtb-045397874c170c745".
* `customized_routes` - (Optional) Advanced option. Customized Spoke VPC Routes. It allows the admin to enter non-RFC1918 routes in the VPC route table targeting the TGW. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16".
* `customized_route_advertisement` - (Optional) Advanced option. Customized route(s) to advertise. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16". 
* `disable_local_route_propagation` - (Optional and ForceNew) Switch to allow admin not to propagate the VPC CIDR to the security domain/TGW route table that it is being attached to. Valid values: true, false. Default value: false.

-> **NOTE:** If used to attach/detach FireNet Gateway to/from TGW Firewall Domain, "vpc_id" is the ID of the Security VPC, and "security_domain_name" is the domain name of the Aviatrix Firewall Domain in TGW.

## Import

Instance aws_tgw_vpc_attachment can be imported using the tgw_name, security_domain_name and vpc_id, e.g.

```
$ terraform import aviatrix_aws_tgw_vpc_attachment.test tgw_name~security_domain_name~vpc_id
```