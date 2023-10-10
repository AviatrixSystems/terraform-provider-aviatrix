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
  network_domain_name  = "network-domain-name"
  attachment_name      = "attachment-name"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `network_domain_name` - (Required) Name of the Segmentation Network Domain.
* `attachment_name` - (Required) Attachment name to associate with the network domain. For spoke gateways, use spoke gateway name. For VLAN, use <site-id>:<vlan-id>.

### Optional

-> **NOTE:** `transit_gateway_name` is an optional and computed attribute now, and it will only be a computed attribute in the V3.2.0 release. 

* `transit_gateway_name` - (Optional) Name of the Transit Gateway.

## Import

-> **NOTE:** Starting from Aviatrix Terraform Provider R3.0+, the resource ID will not contain `transit_gateway_name` since it is optional. If you are using/upgraded to Aviatrix Terraform Provider R3.0+, and an **aviatrix_segmentation_network_domain_association** resource was originally created with a provider version <R3.0, please perform a 'terraform refresh' to rectify the state file.

**aviatrix_segmentation_network_domain_association** can be imported using `network_domain_name` and `attachment_name` separated by a `~` e.g.

```
$ terraform import aviatrix_segmentation_network_domain_association.test network_domain_name~attachment_name
```
