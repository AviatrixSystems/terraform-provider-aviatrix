---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_branch_router_aws_tgw_attachment"
description: |-
  Creates and manages a branch router and AWS TGW attachment
---

# aviatrix_branch_router_aws_tgw_attachment

The **aviatrix_branch_router_aws_tgw_attachment** resource allows the creation and management of a branch router and AWS TGW attachment

## Example Usage

```hcl
# Create an Aviatrix Branch Router and AWS TGW attachment
resource "aviatrix_branch_router_aws_tgw_attachment" "test_branch_router_aws_tgw_attachment" {
	connection_name           = "test-conn"
	branch_name               = "branch-router"
	aws_tgw_name              = "tgw-test"
	branch_router_bgp_asn     = 65001
	security_domain_name      = "Default_Domain"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `connection_name` - Connection name.
* `branch_name` - Branch router name.
* `aws_tgw_name` - AWS TGW name.
* `branch_router_bgp_asn` - BGP AS Number for branch router.
* `security_domain_name` - Security Domain Name for the attachment.

## Import

**branch_router_aws_tgw_attachment** can be imported using the `connection_name`, `branch_name` and `aws_tgw_name`, e.g.

```
$ terraform import aviatrix_branch_router_aws_tgw_attachment.test connection-name~branch-name~aws-tgw-name
```
