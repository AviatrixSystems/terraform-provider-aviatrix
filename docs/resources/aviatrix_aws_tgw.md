---
subcategory: "TGW Orchestrator"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw"
description: |-
  Creates and manages Aviatrix AWS TGWs
---

# aviatrix_aws_tgw

The **aviatrix_aws_tgw** resource allows the creation and management of Aviatrix-created AWS TGWs.

~> **NOTE:** If you are planning to attach VPCs to the **aviatrix_aws_tgw** resource and anticipate updating it often and/or using advanced options such as customized route advertisement, we highly recommend managing those VPCs outside this resource by setting `manage_vpc_attachment` to false and using the **aviatrix_aws_tgw_vpc_attachment** resource instead of the in-line `attached_vpc {}` block.

## Example Usage

```hcl
# Create an Aviatrix AWS TGW
resource "aviatrix_aws_tgw" "test_aws_tgw" {
  account_name                      = "devops"
  aws_side_as_number                = "64512"
  manage_vpc_attachment             = false
  manage_transit_gateway_attachment = false
  region                            = "us-east-1"
  tgw_name                          = "test-AWS-TGW"

  security_domains {
    connected_domains    = [
      "Default_Domain",
      "Shared_Service_Domain",
      "mysdn1"
    ]
    security_domain_name = "Aviatrix_Edge_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain",
      "Shared_Service_Domain"
    ]    
    security_domain_name = "Default_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain",
      "Default_Domain"
    ]
    security_domain_name = "Shared_Service_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain"
    ]
    security_domain_name = "SDN1"
  }

  security_domains {
    security_domain_name = "mysdn2"
  }

  security_domains {
    security_domain_name = "firewall-domain"
    aviatrix_firewall    = true
  }
}
```
```hcl
# Create an Aviatrix AWSGov TGW
resource "aviatrix_aws_tgw" "test_aws_gov_tgw" {
  account_name                      = "devops"
  cloud_type                        = 256
  aws_side_as_number                = "64512"
  manage_vpc_attachment             = false
  manage_transit_gateway_attachment = false
  region                            = "us-gov-east-1"
  tgw_name                          = "test-AWSGov-TGW"

  security_domains {
    connected_domains    = [
      "Default_Domain",
      "Shared_Service_Domain",
      "mysdn1"
    ]
    security_domain_name = "Aviatrix_Edge_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain",
      "Shared_Service_Domain"
    ]    
    security_domain_name = "Default_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain",
      "Default_Domain"
    ]
    security_domain_name = "Shared_Service_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain"
    ]
    security_domain_name = "SDN1"
  }

  security_domains {
    security_domain_name = "mysdn2"
  }

  security_domains {
    security_domain_name = "firewall-domain"
    aviatrix_firewall    = true
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `tgw_name` - (Required) Name of the AWS TGW to be created
* `account_name` - (Required) Name of the cloud account in the Aviatrix controller.
* `region` - (Required) AWS region of AWS TGW to be created in
* `aws_side_as_number` - (Required) BGP Local ASN (Autonomous System Number). Integer between 1-4294967294. Example: "65001".

!> **WARNING:** Attribute `security_domain` has been deprecated as of provider version R2.19+ and will not receive further updates. Please set `manage_security_domain` to false, and use the standalone `aviatrix_aws_tgw_security_domain` resource instead.

* `security_domains` - (Required if `manage_security_domain` is true) Security Domains to create together with AWS TGW's creation. Three default domains, along with the connections between them, are created automatically. These three domains can't be deleted, but the connection between any two of them can be.
  * `security_domain_name` - (Required) Three default domains ("Aviatrix_Edge_Domain", "Default_Domain" and "Shared_Service_Domain") are required with AWS TGW's creation.
  * `aviatrix_firewall` - (Optional) Set to true if the security domain is to be used as an Aviatrix Firewall Domain for the Aviatrix Firewall Network. Valid values: true, false. Default value: false.
  * `native_egress` - (Optional) Set to true if the security domain is to be used as a native egress domain (for non-Aviatrix Firewall Network-based central Internet bound traffic). Valid values: true, false. Default value: false.
  * `native_firewall` - (Optional) Set to true if the security domain is to be used as a native firewall domain (for non-Aviatrix Firewall Network-based firewall traffic inspection). Valid values: true, false. Default value: false.
  * `connected_domains` - (Optional) A list of domains connected to the domain (name: `security_domain_name`) together with its creation.

### VPC Attachments

!> **WARNING:** Attribute `attached_vpc` has been deprecated as of provider version R2.18.1+ and will not receive further updates. Please set `manage_vpc_attachment` to false, and use the standalone `aviatrix_aws_tgw_vpc_attachment` resource instead. 

-> **NOTE:** The `attached_vpc` code block is to be nested under the `security_domains` block. Please see the code examples above for more information.

* `attached_vpc` - (Optional) A list of VPCs attached to the domain (name: `security_domain_name`) together with its creation. This list needs to be null for "Aviatrix_Edge_Domain".
  * `vpc_region` - (Required) Region of the VPC, needs to be consistent with AWS TGW's region.
  * `vpc_account_name` - (Required) Cloud account name of the VPC in the Aviatrix controller.
  * `vpc_id` - (Required) VPC ID of the VPC to be attached to the security domain
  * `subnets` - (Optional) Advanced option. VPC subnets separated by ',' to attach to the VPC. If left blank, the Aviatrix Controller automatically selects a subnet representing each AZ for the VPC attachment. Example: "subnet-214f5646,subnet-085e8c81a89d70846".
  * `route_tables` - (Optional) Advanced option. Route tables separated by ',' to participate in TGW Orchestrator, i.e., learned routes will be propagated to these route tables. Example: "rtb-212ff547,rtb-045397874c170c745".
  * `customized_routes` - (Optional) Advanced option. Customized Spoke VPC Routes. It allows the admin to enter non-RFC1918 routes in the VPC route table targeting the TGW. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16".
  * `customized_route_advertisement` - (Optional) Advanced option. Customized route(s) to be advertised to other VPCs that are connected to the same TGW. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16".    
  * `disable_local_route_propagation` - (Optional) Advanced option. If set to true, it disables automatic route propagation of this VPC to other VPCs within the same security domain. Valid values: true, false. Default value: false.

### Misc.

!> **WARNING:** Attribute `attached_aviatrix_transit_gateway` has been deprecated as of provider version R2.18.1+ and will not receive further updates. Please set `manage_transit_gateway_attachment` to false, and use the standalone `aviatrix_aws_tgw_transit_gateway_attachment` resource instead.

* `attached_aviatrix_transit_gateway` - (Optional) A list of names of Aviatrix Transit Gateway(s) (transit VPCs) to attach to the Aviatrix_Edge_Domain.
* `cloud_type` - (Optional) Type of cloud service provider, requires an integer value. Supported for AWS (1) and AWSGov (256). Default value: 1.
* `manage_security_domain` - (Optional) This parameter is a switch used to determine whether or not to manage security domains using the **aviatrix_aws_tgw** resource. If this is set to false, creation and management of security domains must be done using the **aviatrix_aws_tgw_security_domain** resource. Valid values: true, false. Default value: true.

-> **NOTE:** `manage_security_domain` - If you are using/upgraded to Aviatrix Terraform Provider R2.19+, and an **aviatrix_aws_tgw** resource was originally created with a provider version <R2.19, you must do 'terraform refresh' to update and apply the attribute's default value (true) into the state file.

* `manage_transit_gateway_attachment` - (Optional) This parameter is a switch used to determine whether or not to manage transit gateway attachments to the TGW using the **aviatrix_aws_tgw** resource. If this is set to false, attachment of transit gateways must be done using the **aviatrix_aws_tgw_transit_gateway_attachment** resource. Valid values: true, false. Default value: true.

-> **NOTE:** `manage_transit_gateway_attachment` - If you are using/upgraded to Aviatrix Terraform Provider R2.13+, and an **aviatrix_aws_tgw** resource was originally created with a provider version <R2.13, you must do 'terraform refresh' to update and apply the attribute's default value (true) into the state file.

* `manage_vpc_attachment` - (Optional) This parameter is a switch used to determine whether or not to manage VPC attachments to the TGW using the **aviatrix_aws_tgw** resource. If this is set to false, attachment of VPCs must be done using the **aviatrix_aws_tgw_vpc_attachment** resource. Valid values: true, false. Default value: true.

-> **NOTE:** `manage_vpc_attachment` - If you are using/upgraded to Aviatrix Terraform Provider R1.5+, and an **aviatrix_aws_tgw** resource was originally created with a provider version <R1.5, you must do 'terraform refresh' to update and apply the attribute's default value (true) into the state file.

* `enable_multicast` - (Optional) Enable multicast. Default value: false. Valid values: true, false. Available in provider version R2.17+.
* `cidrs` - (Optional) Set of TGW CIDRs. For example, `cidrs = ["10.0.10.0/24", "10.1.10.0/24"]`. Available as of provider version R2.18.1+.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `tgw_id`- TGW ID. Available as of provider version R2.19+.

## Import

**aws_tgw** can be imported using the `tgw_name`, e.g.

```
$ terraform import aviatrix_aws_tgw.test tgw_name
```

-> **NOTE:** If `manage_security_domain` is set to "false", import action will also import the information of the security domains into the state file. Will need to do *terraform apply* to sync `manage_security_domain` to "false".

-> **NOTE:** If `manage_vpc_attachment` is set to "false", import action will also import the information of the VPCs attached to TGW into the state file. Will need to do *terraform apply* to sync `manage_vpc_attachment` to "false".

-> **NOTE:** If `manage_transit_gateway_attachment` is set to "false", import action will also import the information of the transit gateway attached to TGW into the state file. Will need to do *terraform apply* to sync `manage_transit_gateway_attachment` to "false".
