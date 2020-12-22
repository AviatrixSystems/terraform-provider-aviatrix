---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_segmentation_security_domain"
description: |-
  Creates and manages an Aviatrix Segmentation Security Domain
---

# aviatrix_segmentation_security_domain

The **aviatrix_segmentation_security_domain** resource handles creation of [Transit Segmentation](https://docs.aviatrix.com/HowTos/transit_segmentation_faq.html) Security Domains.

## Example Usage

```hcl
# Create an Aviatrix Segmentation Security Domain
resource "aviatrix_segmentation_security_domain" "test_segmentation_security_domain" {
  domain_name = "domain-a"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `domain_name` - (Required) Name of the Security Domain.

## Import

**aviatrix_segmentation_security_domain** can be imported using the `domain_name`, e.g.

```
$ terraform import aviatrix_segmentation_security_domain.test domain_name
```
