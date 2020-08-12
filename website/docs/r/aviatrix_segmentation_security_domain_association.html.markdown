---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_segmentation_security_domain_association"
description: |-
  Creates and manages an Aviatrix Segmentation Security Domain Association
---

# aviatrix_segmentation_security_domain_association

The **aviatrix_segmentation_security_domain_association** resource handles creation of [Transit Segmentation](https://docs.aviatrix.com/HowTos/transit_segmentation_faq.html) Security Domain and Transit Gateway Attachment Associations.

## Example Usage

```hcl
# Create an Aviatrix Segmentation Security Domain Association
resource "aviatrix_segmentation_security_domain_association" "test_segmentation_security_domain_association" {
  transit_gateway_name = "transit-gw-name"
  security_domain_name = "security-domain-name"
  attachment_name      = "attachment-name"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `transit_gateway_name` - (Required) Name of the Transit Gateway.
* `security_domain_name` - (Required) Name of the Segmentation Security Domain.
* `attachment_name` - (Required) Name of the transit gateway attachment, Spoke or Edge, to associate with the security domain.

## Import

**aviatrix_segmentation_security_domain_association** can be imported using `transit_gateway_name`, `security_domain_name` and `attachment_name` separated by a `~` e.g.

```
$ terraform import aviatrix_segmentation_security_domain_association.test transit_gateway_name~security_domain_name~attachment_name
```
