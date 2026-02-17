---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_segmentation_network_domain"
description: |-
  Creates and manages an Aviatrix Segmentation Network Domain
---

# aviatrix_segmentation_network_domain

The **aviatrix_segmentation_network_domain** resource handles creation of [Transit Segmentation](https://docs.aviatrix.com/HowTos/transit_segmentation_faq.html) Network Domains.

## Example Usage

```hcl
# Create an Aviatrix Segmentation Network Domain
resource "aviatrix_segmentation_network_domain" "test_segmentation_network_domain" {
  domain_name = "domain-a"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `domain_name` - (Required) Name of the Network Domain.

## Import

**aviatrix_segmentation_network_domain** can be imported using the `domain_name`, e.g.

```
$ terraform import aviatrix_segmentation_network_domain.test domain_name
```
