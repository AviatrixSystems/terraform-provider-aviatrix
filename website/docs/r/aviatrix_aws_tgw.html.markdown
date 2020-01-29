---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw"
description: |-
  Creates and manages Aviatrix AWS TGWs
---

# aviatrix_aws_tgw

The aviatrix_aws_tgw resource allows the creation and management of AWS TGWs.

~> **NOTE:** If you are planning to attach VPCs to the **aviatrix_aws_tgw** resource and anticipate updating it often and/or using advanced options such as customized route advertisement, we highly recommend managing those VPCs outside this resource by setting `manage_vpc_attachment` to false and using the **aviatrix_aws_tgw_vpc_attachment** resource instead of the in-line `attached_vpc {}` block.

## Example Usage

```hcl
# Create an Aviatrix AWS TGW
resource "aviatrix_aws_tgw" "test_aws_tgw" {
  account_name                      = "devops"
  attached_aviatrix_transit_gateway = [
     "avxtransitgw"
  ]
  aws_side_as_number                = "64512"
  manage_vpc_attachment             = true
  region                            = "us-east-1"
  tgw_name                          = "testAWSTgw"

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

    attached_vpc {
      vpc_account_name = "devops1"
      vpc_id           = "vpc-0e2fac2b91"
      vpc_region       = "us-east-1"
    }

    attached_vpc {
      vpc_account_name = "devops1"
      vpc_id           = "vpc-0c63660a16"
      vpc_region       = "us-east-1"
    }

    attached_vpc {
      vpc_account_name = "devops"
      vpc_id           = "vpc-032005cc444"
      vpc_region       = "us-east-1"
    }
  }

  security_domains {
    security_domain_name = "mysdn2"

    attached_vpc {
      vpc_region                      = "us-east-1"
      vpc_account_name                = "devops"
      vpc_id                          = "vpc-03200566666"
      customized_routes               = "10.8.0.0/16,10.9.0.0/16"
      disable_local_route_propagation = true
    }
  }

  security_domains {
    security_domain_name = "firewall-domain"
    aviatrix_firewall    = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `tgw_name` - (Required) Name of the AWS TGW which is going to be created.
* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `region` - (Required) Region of cloud provider(AWS).
* `aws_side_as_number` - (Required) BGP Local ASN (Autonomous System Number). Integer between 1-65535. Example: "65001".
* `security_domains` - (Required) Security Domains to create together with AWS TGW's creation. Three default domains are created automatically together with the AWS TGW's creation, so are the connections between any two of them. These three domains can't be deleted, but the connection between any two of them can be deleted.
  * `security_domain_name` - (Required) Three default domains ("Aviatrix_Edge_Domain", "Default_Domain" and "Shared_Service_Domain") are required with AWS TGW's creation.
  * `aviatrix_firewall` - (Optional) Set to true if the security domain is an aviatrix firewall domain. Valid values: true, false. Default value: false.
  * `native_egress` - (Optional) Set to true if the security domain is a native egress domain. Valid values: true, false. Default value: false.
  * `native_firewall` - (Optional) Set to true if the security domain is a native firewall domain. Valid values: true, false. Default value: false.
  * `connected_domains` - (Optional) A list of domains connected to the domain (name: `security_domain_name`) together with its creation.
  * `attached_vpc` - (Optional) A list of VPCs attached to the domain (name: `security_domain_name`) together with its creation. This list needs to be null for "Aviatrix_Edge_Domain".
    * `vpc_region` - (Required) Region of the vpc, needs to be consistent with AWS TGW's region.
    * `vpc_account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
    * `vpc_id` - (Required) This parameter represents the ID of the VPC which is going to be attached to the security domain (name: `security_domain_name`) which is going to be created.
    * `subnets` - (Optional) Advanced option. VPC subnets separated by ',' to attach to the VPC. If left blank, Aviatrix Controller automatically selects a subnet representing each AZ for the VPC attachment. Example: "subnet-214f5646,subnet-085e8c81a89d70846".
    * `route_tables` - (Optional) Advanced option. Route tables separated by ',' to participate in TGW Orchestrator, i.e., learned routes will be propagated to these route tables. Example: "rtb-212ff547,rtb-045397874c170c745".
    * `customized_routes` - (Optional) Advanced option. Customized Spoke VPC Routes. It allows the admin to enter non-RFC1918 routes in the VPC route table targeting the TGW. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16".
    * `customized_route_advertisement` - (Optional) Advanced option. Customized route(s) to advertise. Example: "10.8.0.0/16,10.9.0.0/16,10.10.0.0/16".    
    * `disable_local_route_propagation` - (Optional) Advanced option. Switch to allow admin not to propagate the VPC CIDR to the security domain/TGW route table that it is being attached to. Valid values: true, false. Default value: false.
* `attached_aviatrix_transit_gateway` - (Optional) A list of Names of Aviatrix Transit Gateway to attach to one of the three default domains: Aviatrix_Edge_Domain.
* `manage_vpc_attachment` - (Optional) This parameter is a switch used to allow attaching VPCs to tgw using the aviatrix_aws_tgw resource. If it is set to false, attachment of VPC must be done using the aviatrix_aws_tgw_vpc_attachment resource. Valid values: true or false. Default value is true.

-> **NOTE:** `manage_vpc_attachment` - If you are using/upgraded to Aviatrix Terraform Provider R1.5+, and an aws_tgw resource was originally created with a provider version <R1.5, you must do ‘terraform refresh’ to update and apply the attribute’s default value (true) into the state file.

## Import

Instance aws_tgw can be imported using the tgw_name, e.g.

```
$ terraform import aviatrix_aws_tgw.test tgw_name
```

-> **NOTE:** If `manage_vpc_attachment` is set to "false", import action will also import the information of the VPCs attached to tgw into the state file. Will need to do `terraform apply` to sync `manage_vpc_attachment` to "false".
