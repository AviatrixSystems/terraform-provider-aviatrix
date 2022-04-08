---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_segmentation_network_domain_association"
description: |-
  Creates and manages an Aviatrix Segmentation Network Domain Association
---

# aviatrix_segmentation_network_domain_association

The **aviatrix_segmentation_network_domain_association** resource handles creation of [Transit Segmentation](https://docs.aviatrix.com/HowTos/transit_segmentation_faq.html) Network Domain and Transit Gateway Attachment Associations.

## Example Usage

```hcl
# Create an Aviatrix Segmentation Network Domain Association
resource "aviatrix_segmentation_network_domain_association" "test_segmentation_network_domain_association" {
  transit_gateway_name = "transit-gw-name"
  network_domain_name  = "network-domain-name"
  attachment_name      = "attachment-name"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `transit_gateway_name` - (Required) Name of the Transit Gateway.
* `network_domain_name` - (Required) Name of the Segmentation Network Domain.
* `attachment_name` - (Required) Name of the transit gateway attachment, Spoke or Edge, to associate with the network domain.

## Import

**aviatrix_segmentation_network_domain_association** can be imported using `transit_gateway_name`, `network_domain_name` and `attachment_name` separated by a `~` e.g.

```
$ terraform import aviatrix_segmentation_network_domain_association.test transit_gateway_name~network_domain_name~attachment_name
```
