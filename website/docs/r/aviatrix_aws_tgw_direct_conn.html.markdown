---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_direct_conn"
description: |-
  Creates and manages Aviatrix AWS TGW Direct Connections
---
 
# aviatrix_aws_tgw_direct_conn
 
The aviatrix_aws_tgw_direct_conn resource allows the creation and management of Aviatrix AWS TGW Direct Connections.
 
## Example Usage
 
```hcl
# Create an Aviatrix AWS TGW Direct Connection
resource "aviatrix_aws_tgw_direct_conn" "test_aws_tgw_direct_conn" {
  tgw_name                 = "myawstgw1"
  direct_conn_account_name = "username"
  direct_conn_gw_id        = "30321d76-dd01-49bf"
  route_domain_name        = "mysdn1"
  allowed_prefix           = "10.12.0.0/24"
}
```
 
## Argument Reference
 
The following arguments are supported:
 
* `tgw_name` - (Required) This parameter represents the name of an AWS TGW.
* `direct_conn_account_name` - (Required) This parameter represents the name of an Account in Aviatrix controller.
* `direct_conn_gw_id` - (Required) This parameter represents the name of a Direct Connect Gateway ID.
* `route_domain_name` - (Required) The name of a route domain, to which the direct connect gateway will be attached.
* `allowed_prefix` - (Required) A list of comma separated CIDRs for DXGW to advertise to remote(on-prem).
 
## Import
 
Instance aws_tgw_direct_conn can be imported using the tgw_name and direct_conn_gw_id, e.g.
 
```
$ terraform import aviatrix_aws_tgw_direct_conn.test tgw_name~direct_conn_gw_id
```