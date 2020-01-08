---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firenet"
description: |-
  Creates and manages Aviatrix FireNets
---

# aviatrix_firenet

The aviatrix_firenet resource allows the creation and management of Aviatrix FireNets.

~> **NOTE:** This resource is used in conjunction with multiple other resources that may include, and are not limited to: "firewall_instance", "aws_tgw", and "transit_gateway" resources, under the Aviatrix FireNet solution. Explicit dependencies may be set using `depends_on`. For more information on proper FireNet configuration, please see the workflow [here](https://docs.aviatrix.com/HowTos/firewall_network_workflow.html).

## Example Usage

```hcl
# Create an Aviatrix FireNet associated to a Firewall Instance
resource "aviatrix_firenet" "test_firenet" {
  vpc_id             = "vpc-032005cc371"
  inspection_enabled = true
  egress_enabled     = false

  firewall_instance_association {
    firenet_gw_name      = "avx_firenet_gw"
    instance_id          = "i-09dc118db6a1eb901"
    firewall_name        = "avx_firewall_instance"
    attached             = true
    lan_interface        = "eni-0a34b1827bf222353"
    management_interface = "eni-030e53176c7f7d34a"
    egress_interface     = "eni-03b8dd53a1a731481"
  }
}
```
```hcl
# Create an Aviatrix FireNet associated to an FQDN Gateway
resource "aviatrix_firenet" "test_firenet" {
  vpc_id             = "vpc-032005cc371"
  inspection_enabled = true
  egress_enabled     = false

  firewall_instance_association {
    firenet_gw_name = "avx_firenet_gw"
    instance_id     = "avx_fqdn_gateway"
    vendor_type     = "fqdn_gateway"
    attached        = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) ID of the Security VPC.
* `inspection_enabled` - (Optional) Enable/Disable traffic inspection. Valid values: true, false. Default value: true.
* `egress_enabled` - (Optional) Enable/Disable egress through firewall. Valid values: true, false. Default value: false.
* `firewall_instance_association` - (Optional) List of firewall instances to be associated with fireNet.
  * `firenet_gw_name` - (Required) Name of the primary FireNet gateway.
  * `instance_id` - (Required) ID of Firewall instance, if associating FQDN gateway to fireNet, it is FQDN gateway's gw_name..
  * `vendor_type` - (Optional) Type of the firewall. Valid values: "Generic", "fqdn_gateway". Default value: "Generic". Value "fqdn_gateway" is required for FQDN gateway.  
  * `firewall_name` - (Optional) Firewall instance name, required if it is a firewall instance.
  * `lan_interface`- (Optional) Lan interface ID, required if it is a firewall instance.
  * `management_interface` - (Optional) Management interface ID, required if it is a firewall instance.
  * `egress_interface`- (Optional) Egress interface ID, required if it is a firewall instance.
  * `attached`- (Optional) Switch to attach/detach firewall instance to/from fireNet. Valid values: true, false. Default value: false.

-> **NOTE:** `inspection_enabled` - Default value is true for associating firewall instance to fireNet. Only false is supported for associating FQDN gateway to fireNet.

-> **NOTE:** `egress_enabled` - Default value is false for associating firewall instance to fireNet. Only true is supported for associating FQDN gateway to fireNet.

-> **NOTE:** `firewall_instance_association` - If associating FQDN gateway to fireNet, "single_az_ha" needs to be enabled for the FQDN gateway.

## Import

Instance firenet can be imported using the vpc_id, e.g.

```
$ terraform import aviatrix_firenet.test vpc_id
```
