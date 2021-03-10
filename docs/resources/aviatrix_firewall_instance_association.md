---
subcategory: "Firewall Network"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_instance_association"
description: |-
  Create and manage a single firewall instance association
---

# aviatrix_firewall_instance_association

The **aviatrix_firewall_instance_association** resource allows for the creation and management of a firewall instance association. To use this resource you must also have an `aviatrix_firenet` resource with it's `manage_firewall_instance_association` attribute set to false.

Available in provider version R2.17.1+.

## Example Usage

```hcl
# Associate an Aviatrix FireNet Gateway with a Firewall Instance
resource "aviatrix_firewall_instance_association" "firewall_instance_association_1" {
  vpc_id               = aviatrix_firewall_instance.firewall_instance_1.vpc_id
  firenet_gw_name      = aviatrix_transit_gateway.transit_gateway_1.gw_name
  instance_id          = aviatrix_firewall_instance.firewall_instance_1.instance_id
  firewall_name        = aviatrix_firewall_instance.firewall_instance_1.firewall_name
  lan_interface        = aviatrix_firewall_instance.firewall_instance_1.lan_interface
  management_interface = aviatrix_firewall_instance.firewall_instance_1.management_interface
  egress_interface     = aviatrix_firewall_instance.firewall_instance_1.egress_interface
  attached             = true
}
```
```hcl
# Associate an GCP Aviatrix FireNet Gateway with a Firewall Instance
resource "aviatrix_firewall_instance_association" "firewall_instance_association_1" {
  vpc_id               = aviatrix_firewall_instance.firewall_instance_1.vpc_id
  firenet_gw_name      = aviatrix_transit_gateway.transit_gateway_1.gw_name
  instance_id          = aviatrix_firewall_instance.firewall_instance_1.instance_id
  lan_interface        = aviatrix_firewall_instance.firewall_instance_1.lan_interface
  management_interface = aviatrix_firewall_instance.firewall_instance_1.management_interface
  egress_interface     = aviatrix_firewall_instance.firewall_instance_1.egress_interface
  attached             = true
}
```

## Argument Reference

The following arguments are supported:

### Required
* `vpc_id` - (Required) VPC ID of the Security VPC.

-> **NOTE:** If associating FQDN gateway to FireNet, `single_az_ha` needs to be enabled for the FQDN gateway.

* `firenet_gw_name` - (Optional) Name of the primary FireNet gateway. Required for FireNet without Native GWLB VPC.
* `instance_id` - (Required) ID of Firewall instance.

-> **NOTE:** If associating FQDN gateway to FireNet, it is FQDN gateway's `gw_name`. For Azure FireNet, it is the `firewall_name` concatenated with a ":" and the Resource Group of the `vpc_id` set for that instance.

* `vendor_type` - (Optional) Type of firewall. Valid values: "Generic", "fqdn_gateway". Default value: "Generic". Value "fqdn_gateway" is required for FQDN gateway.  
* `firewall_name` - (Optional) Firewall instance name. **Required for non-GCP firewall instance. For GCP, this field should not be set.**
* `lan_interface`- (Optional) Lan interface ID. **Required if it is a firewall instance.**
* `management_interface` - (Optional) Management interface ID. **Required if it is a firewall instance.**
* `egress_interface`- (Optional) Egress interface ID. **Required if it is a firewall instance.**
* `attached`- (Optional) Switch to attach/detach firewall instance to/from FireNet. Valid values: true, false. Default value: false.


## Import

**firewall_instance_association** can be imported using the `vpc_id`, `firenet_gw_name` and `instance_id` in the form `vpc_id~~firenet_gw_name~~instance_id` e.g.

```
$ terraform import aviatrix_firewall_instance_association.test vpc_id~~firenet_gw_name~~instance_id
```

When using a Native GWLB VPC where there is no `firenet_gw_name` but the ID is in the same form e.g.

```
$ terraform import aviatrix_firewall_instance_association.test vpc_id~~~~instance_id
```
