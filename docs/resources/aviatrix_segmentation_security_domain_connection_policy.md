---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_segmentation_security_domain_connection_policy"
description: |-
  Creates and manages an Aviatrix Segmentation Security Domain Connection Policy
---

# aviatrix_segmentation_security_domain_connection_policy

The **aviatrix_segmentation_security_domain_connection_policy** resource handles creation of [Transit Segmentation](https://docs.aviatrix.com/HowTos/transit_segmentation_faq.html) Security Domain Connection Policies.

## Example Usage

```hcl
# Create an Aviatrix Segmentation Security Domain
resource "aviatrix_segmentation_security_domain_connection_policy" "test_segmentation_security_domain_connection_policy" {
  domain_name_1 = "domain-a"
  domain_name_2 = "domain-b"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `domain_name_1` - (Required) Name of the Security Domain to connect to Domain 2.
* `domain_name_2` - (Required) Name of the Security Domain to connect to Domain 1.

## Import

**aviatrix_segmentation_security_domain_connection_policy** can be imported using `domain_name_1` and `domain_name_2` separated by a `~`, e.g.

```
$ terraform import aviatrix_segmentation_security_domain_connection_policy.test domain_name_1~domain_name_2
```
