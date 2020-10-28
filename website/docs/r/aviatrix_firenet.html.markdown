---
subcategory: "Firewall Network"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firenet"
description: |-
  Creates and manages Aviatrix FireNets
---

# aviatrix_firenet

The **aviatrix_firenet** resource allows the creation and management of [Aviatrix Firewall Networks](https://docs.aviatrix.com/HowTos/firewall_network_faq.html).

~> **NOTE:** This resource is used in conjunction with multiple other resources that may include, and are not limited to: **firewall_instance**, **aws_tgw**, and **transit_gateway** resources or even **aviatrix_fqdn**, under the Aviatrix FireNet solution. Explicit dependencies may be set using `depends_on`. For more information on proper FireNet configuration, please see the workflow [here](https://docs.aviatrix.com/HowTos/firewall_network_workflow.html).

## Example Usage

```hcl
# Create an Aviatrix FireNet associated to a Firewall Instance
resource "aviatrix_firenet" "test_firenet" {
  vpc_id             = "vpc-032005cc371"
  inspection_enabled = true
  egress_enabled     = false

  firewall_instance_association {
    firenet_gw_name      = "avx-firenet-gw"
    instance_id          = "i-09dc118db6a1eb901"
    firewall_name        = "avx-firewall-instance"
    attached             = true
    lan_interface        = "eni-0a34b1827bf222353"
    management_interface = "eni-030e53176c7f7d34a"
    egress_interface     = "eni-03b8dd53a1a731481"
  }
}
```
```hcl
# Create an Aviatrix FireNet associated to an FQDN Gateway (AWS)
resource "aviatrix_firenet" "test_firenet" {
  vpc_id             = "vpc-032005cc371"
  inspection_enabled = true
  egress_enabled     = false

  firewall_instance_association {
    firenet_gw_name = "avx-firenet-gw"
    instance_id     = "avx-fqdn-gateway"
    vendor_type     = "fqdn_gateway"
    attached        = true
  }
}
```

```hcl
# Create an Aviatrix FireNet associated to an FQDN Gateway (Azure)
resource "aviatrix_firenet" "test_firenet" {
  vpc_id             = "vpc-032005cc371"
  inspection_enabled = true
  egress_enabled     = false

  firewall_instance_association {
    firenet_gw_name = "avx-firenet-gw"
    instance_id     = "avx-fqdn-gateway"
    vendor_type     = "fqdn_gateway"
    attached        = true
    lan_interface   = "<< LAN interface id of the FQDN gateway created with additional LAN interface >>"
  }
}
```
## Argument Reference

The following arguments are supported:

### Required
* `vpc_id` - (Required) VPC ID of the Security VPC.
* `inspection_enabled` - (Optional) Enable/disable traffic inspection. Valid values: true, false. Default value: true.

-> **NOTE:** `inspection_enabled` - Default value is true for associating firewall instance to FireNet. Only false is supported for associating FQDN gateway to FireNet.

* `egress_enabled` - (Optional) Enable/disable egress through firewall. Valid values: true, false. Default value: false.

-> **NOTE:** `egress_enabled` - Default value is false for associating firewall instance to FireNet. Only true is supported for associating FQDN gateway to FireNet.

* `hashing_algorithm` - (Optional) Hashing algorithm to load balance traffic across the firewall. Valid values: "2-Tuple", "5-Tuple". Default value: "5-Tuple".
* `manage_firewall_instance_association` - (Optional) Enable this attribute to manage firewall associations in-line. If set to true, in-line `firewall_instance_association` blocks can be used. If set to false, all firewall associations must be managed via standalone `aviatrix_firewall_instance_association` resources. Default value: true. Valid values: true or false. Available as of provider version R2.18+.
* `manage_firewall_instance_association` - (Optional) Enable this attribute to manage firewall associations in-line. If set to true, in-line `firewall_instance_association` blocks can be used. If set to false, all firewall associations must be managed via standalone `aviatrix_firewall_instance_association` resources. Default value: true. Valid values: true or false. Available as of provider version R2.17.1+.

### Firewall Association

-> **NOTE:** `firewall_instance_association` - If associating FQDN gateway to FireNet, `single_az_ha` needs to be enabled for the FQDN gateway.

* `firewall_instance_association` - (Optional) Dynamic block of firewall instance(s) to be associated with the FireNet.
  * `firenet_gw_name` - (Required) Name of the primary FireNet gateway.
  * `instance_id` - (Required) ID of Firewall instance.

  -> **NOTE:** If associating FQDN gateway to FireNet, it is FQDN gateway's `gw_name`. For Azure FireNet, it is the `firewall_name` concatenated with a ":" and the Resource Group of the `vpc_id` set for that instance.
  -> **NOTE:** If associating FQDN gateway to FireNet in Azure, the `lan_interface` is required. The `lan_interface` can be obtained from the exported attribute `fqdn_lan_interface` in aviatrix_gateway resource block (Available in provider version R2.17.1+).
  
  * `vendor_type` - (Optional) Type of firewall. Valid values: "Generic", "fqdn_gateway". Default value: "Generic". Value "fqdn_gateway" is required for FQDN gateway.  
  * `firewall_name` - (Optional) Firewall instance name. **Required if it is a firewall instance.**
  * `lan_interface`- (Optional) Lan interface ID. **Required if it is a firewall instance or FQDN gateway in Azure.**
  * `management_interface` - (Optional) Management interface ID. **Required if it is a firewall instance.**
  * `egress_interface`- (Optional) Egress interface ID. **Required if it is a firewall instance.**
  * `attached`- (Optional) Switch to attach/detach firewall instance to/from FireNet. Valid values: true, false. Default value: false.


## Import

**firenet** can be imported using the `vpc_id`, e.g.

```
$ terraform import aviatrix_firenet.test vpc_id
```
