---
subcategory: "TGW Orchestrator"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_vpc_attachment"
description: |-
  Manages attaching/detaching VPC to/from an AWS TGW, and FireNet Gateway to TGW Firewall Domain
---

# aviatrix_aws_tgw_vpc_attachment

The **aviatrix_aws_tgw_vpc_attachment** resource manages the attaching & detaching of the VPC to & from an AWS TGW, and FireNet Gateway to TGW Firewall Domain.

## Example Usage

```hcl
# Create an Aviatrix AWS TGW VPC Attachment
resource "aviatrix_aws_tgw_vpc_attachment" "test_aws_tgw_vpc_attachment" {
  tgw_name             = "test-tgw"
  region               = "us-east-1"
  security_domain_name = "my-sdn"
  vpc_account_name     = "test-account"
  vpc_id               = "vpc-0e2fac2b91c6697b3"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `tgw_name` - (Required) Name of the AWS TGW.
* `region` - (Required) AWS Region of the TGW.
* `security_domain_name` - (Required & ForceNew) The name of the security domain, to which the VPC will be attached to. If changed, the VPC will be detached from the old domain, and attached to the new domain.
* `vpc_account_name` - (Required) The name of the cloud account in the Aviatrix controller, which is associated with the VPC.
* `vpc_id` - (Required) VPC ID of the VPC to be attached to the specified `security_domain_name`.

-> **NOTE:** If used to attach/detach FireNet Transit Gateway to/from TGW Firewall Domain, `vpc_id` is the ID of the Security VPC, and `security_domain_name` is the domain name of the Aviatrix Firewall Domain in TGW.

### Advanced Options
* `subnets` - (Optional and ForceNew) Advanced option. VPC subnets separated by ',' to attach to the VPC. If omitted, the Aviatrix Controller automatically computes a subnet representing each AZ for the VPC attachment and Terraform will not manage this attribute. Example: "subnet-214f5646,subnet-085e8c81a89d70846".
* `route_tables` - (Optional and ForceNew) Advanced option. Route tables separated by ',' to participate in TGW Orchestrator, i.e., learned routes will be propagated to these route tables. Example: "rtb-212ff547,rtb-045397874c170c745".
* `customized_routes` - (Optional) Advanced option. Customized Spoke VPC Routes. It allows the admin to enter non-RFC1918 routes in the VPC route table targeting the TGW. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16".
* `customized_route_advertisement` - (Optional and ForceNew) Advanced option. Customized route(s) to be advertised to other VPCs that are connected to the same TGW. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16".
* `disable_local_route_propagation` - (Optional and ForceNew) Advanced option. If set to true, it disables automatic route propagation of this VPC to other VPCs within the same security domain. Valid values: true, false. Default value: false.
* `edge_attachment` - (Optional) Advanced option. To allow access to the private IP of the MGMT interface of the Firewalls, set this attribute to enable Management Access From Onprem. This feature advertises the Firewalls private MGMT subnet to your Edge domain. Example: "vpn-0068bb31917ff2289".

## Import

**aws_tgw_vpc_attachment** can be imported using the `tgw_name`, `security_domain_name` and `vpc_id`, e.g.

```
$ terraform import aviatrix_aws_tgw_vpc_attachment.test tgw_name~security_domain_name~vpc_id
```
