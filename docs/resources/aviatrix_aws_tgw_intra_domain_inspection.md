---
subcategory: "TGW Orchestrator"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_intra_domain_inspection"
description: |-
  Creates and manages the intra domain inspection of security domains in an AWS TGW
---

# aviatrix_aws_tgw_intra_domain_inspection

The **aviatrix_aws_tgw_intra_domain_inspection** resource allows the creation and management of intra domain inspection of security domains in an AWS TGW.

## Example Usage

```hcl
# Create an Aviatrix Intra Domain Inspection
resource "aviatrix_aws_tgw_intra_domain_inspection" "test" {
  tgw_name             = "test-AWS-TGW"
  route_domain_name    = "mysd"
  firewall_domain_name = "firewall-domain"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `tgw_name` - (Required) The AWS TGW name.
* `route_domain_name` - (Required) The name of a security domain.
* `firewall_domain_name` - (Required) The name of a firewall security domain.

## Import

**aviatrix_aws_tgw_intra_domain_inspection** can be imported using the `tgw_name` and `route_domain_name`, e.g.

```
$ terraform import aviatrix_aws_tgw_intra_domain_inspection.test tgw_name~route_domain_name
```
