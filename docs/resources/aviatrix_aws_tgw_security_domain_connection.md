---
subcategory: "TGW Orchestrator"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_security_domain_connection"
description: |-
  Creates and manages the connections between security domains in an AWS TGW
---

# aviatrix_aws_tgw_security_domain_connection

!> **WARNING:** Resource 'aviatrix_aws_tgw_security_domain_connection' will be deprecated in future releases. Please use resource 'aviatrix_aws_tgw_peering_domain_conn' instead.

The **aviatrix_aws_tgw_security_domain_connection** resource allows the creation and management of the connections between security domains in an AWS TGW.

## Example Usage

```hcl
# Create an Aviatrix Security Domain Connection
resource "aviatrix_aws_tgw_security_domain_connection" "test" {
  tgw_name     = "tgw"
  domain_name1 = "domain1"
  domain_name2 = "domain2"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `tgw_name` - (Required) The AWS TGW name.
* `domain_name1` - (Required) The name of a security domain to make a connection.
* `domain_name2` - (Required) The name of another security domain to make a connection.

-> **NOTE:** In the resource ID , `domain_name1` and `domain_name2` will be sorted lexicographically. 

## Import

**aws_tgw_security_domain_connection** can be imported using the `tgw_name`, `domain_name1` and `domain_name2`, e.g.

```
$ terraform import aviatrix_aws_tgw_security_domain_connection.test tgw_name~domain_name1~domain_name2
```

-> **NOTE:** In the resource ID , `domain_name1` and `domain_name2` are sorted lexicographically. When importing, using `tgw_name~domain_name1~domain_name2` or `tgw_name~domain_name2~domain_name1` has the same effect.