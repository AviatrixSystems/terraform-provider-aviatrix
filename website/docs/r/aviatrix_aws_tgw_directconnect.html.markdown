---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_directconnect"
description: |-
  Creates and manages Aviatrix AWS TGW Directconnects
---
 
# aviatrix_aws_tgw_directconnect
 
The aviatrix_aws_tgw_directconnect resource allows the creation and management of Aviatrix AWS TGW Directconnects.
 
## Example Usage
 
```hcl
# Create an Aviatrix AWS TGW Directconnect
resource "aviatrix_aws_tgw_directconnect" "test_aws_tgw_directconnect" {
  tgw_name                   = "myawstgw1"
  directconnect_account_name = "username"
  dx_gateway_id              = "30321d76-dd01-49bf"
  security_domain_name       = "mysdn1"
  allowed_prefix             = "10.12.0.0/24"
}
```
 
## Argument Reference
 
The following arguments are supported:
 
* `tgw_name` - (Required) This parameter represents the name of an AWS TGW.
* `directconnect_account_name` - (Required) This parameter represents the name of an Account in Aviatrix controller.
* `dx_gateway_id` - (Required) This parameter represents the name of a Direct Connect Gateway ID.
* `security_domain_name` - (Required) The name of a security domain, to which the direct connect gateway will be attached.
* `allowed_prefix` - (Required) A list of comma separated CIDRs for DXGW to advertise to remote(on-prem).
 
## Import
 
Instance aws_tgw_directconnect can be imported using the tgw_name and dx_gateway_id, e.g.
 
```
$ terraform import aviatrix_aws_tgw_directconnect.test tgw_name~dx_gateway_id
```