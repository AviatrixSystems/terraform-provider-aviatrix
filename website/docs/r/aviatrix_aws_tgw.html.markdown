---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw"
sidebar_current: "docs-aviatrix-resource-aws_tgw"
description: |-
  Manages the AWS TGWs
---

# aviatrix_aws_tgw

The AviatrixAWSTgw resource manages the AWS TGWs

## Example Usage

```hcl
# Manage AWS TGWs
resource "aviatrix_aws_tgw" "test_aws_tgw" {
  tgw_name                          = "testAWSTgw"
  account_name                      = "devops"
  region                            = "us-east-1"
  aws_side_as_number                = "64512"
  attached_aviatrix_transit_gateway = ["avxtransitgw", "avxtransitgw2"]
  security_domains = [
  {
    security_domain_name = "Aviatrix_Edge_Domain"
    connected_domains    = ["Default_Domain", "Shared_Service_Domain", "SDN1"]          
  },
  {
    security_domain_name = "Default_Domain"
    connected_domains    = ["Aviatrix_Edge_Domain", "Shared_Service_Domain"]    
    attached_vpc         = []      
  },
  {
    security_domain_name = "Shared_Service_Domain"
    connected_domains    = ["Aviatrix_Edge_Domain", "Default_Domain"]
    attached_vpc         = []          
  },
  {
    security_domain_name = "SDN1"
    connected_domains    = ["Aviatrix_Edge_Domain"]
    attached_vpc         = [
    {
      vpc_region       = "us-east-1"
      vpc_account_name = "devops1"
      vpc_id           = "vpc-0e2fac2b91"  
    },
    {
      vpc_region       = "us-east-1"
      vpc_account_name = "devops1"
      vpc_id           = "vpc-0c63660a16"  
    },
    {
      vpc_region       = "us-east-1"
      vpc_account_name = "devops2"
      vpc_id           = "vpc-032005cc37"  
    },
    ]          
  },
  {
    security_domain_name = "SDN2"
    connected_domains    = []
    attached_vpc         = [
    {
      vpc_region       = "us-east-1"
      vpc_account_name = "devops"
      vpc_id           = "vpc-032005cc371"  
    },
    ]          
  },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `tgw_name` - (Required) Name of the AWS TGW which is going to be created.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `region` - (Required) Region of cloud provider(AWS).
* `aws_side_as_number` - (Required) BGP Local ASN (Autonomous System Number). Integer between 1-65535. Example: "65001"
* `attached_aviatrix_transit_gateway` - (Optional) A list of Names of Aviatrix Transit Gateway to attach to one of the three default domains: Aviatrix_Edge_Domain.
* `security_domains` - (Required) Security Domains to create together with AWS TGW's creation. Three default domains are created automatically together with the AWS TGW's creation, so are the connections between any two of them. These three domains can't be deleted, but the connection between any two of them can be deleted.
  * `security_domain_name` - (Required) Three default domains ("Aviatrix_Edge_Domain", "Default_Domain" and "Shared_Service_Domain") are required with AWS TGW's creation.
  * `connected_domains` - (Optional) A list of domains connected to the domain (name: `security_domain_name`) together with its creation.
  * `attached_vpc` - (Optional) A list of VPCs attached to the domain (name: `security_domain_name`) together with its creation. This list needs to be null for "Aviatrix_Edge_Domain".
    * `vpc_region` - (Required) Region of the vpc, needs to be consistent with AWS TGW's region.
    * `vpc_account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller. 
    * `vpc_id` - (Required) This parameter represents the ID of the VPC which is going to be attached to the security domain (name: `security_domain_name`) which is going to be created.
* `manage_vpc_attachment` - (Optional) This parameter is a switch used to allow attaching VPCs to tgw using the aviatrix_aws_tgw resource. If it is set to false, attachment of vpc must be done using the aviatrix_aws_tgw_vpc_attachment resource. Valid values: true or false. Default value is true. 

-> **NOTE:** 

* `manage_vpc_attachment` - If you are using/ upgraded to Aviatrix Terraform Provider v4.2+ , and an aws_tgw resource was originally created with a provider version <4.2, you must do ‘terraform apply’ to update and apply the attribute’s default value (“true”) into the state file. 

## Import

Instance aws_tgw can be imported using the tgw_name, e.g.

```
$ terraform import aviatrix_aws_tgw.test tgw_name
```

If "manage_vpc_attachment" is set to "no", import action will also import the information of the VPCs attached to tgw into the state file. Will need to do "Terraform Apply" to sync "manage_vpc_attachment" to "no".