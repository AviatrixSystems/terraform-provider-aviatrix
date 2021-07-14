---
subcategory: "Firewall Network"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firenet"
description: |-
  Creates and manages Aviatrix FireNets
---

# aviatrix_firenet

The **aviatrix_firenet** resource allows the creation and management of [Aviatrix Firewall Networks](https://docs.aviatrix.com/HowTos/firewall_network_faq.html).

~> **NOTE:** This resource is used in conjunction with multiple other resources that may include, and are not limited to: **firewall_instance**, **firewall_instance_association**, **aws_tgw**, and **transit_gateway** resources or even **aviatrix_fqdn**, under the Aviatrix FireNet solution. Explicit dependencies may be set using `depends_on`. For more information on proper FireNet configuration, please see the workflow [here](https://docs.aviatrix.com/HowTos/firewall_network_workflow.html).

## Example Usage

```hcl
# Create an Aviatrix FireNet
resource "aviatrix_firenet" "test_firenet" {
  vpc_id                               = "vpc-032005cc371"
  inspection_enabled                   = true
  egress_enabled                       = false
  keep_alive_via_lan_interface_enabled = false
  manage_firewall_instance_association = false

  depends_on = [aviatrix_firewall_instance_association.association_1]
}
```

```hcl
# Create an Aviatrix GCP FireNet
resource "aviatrix_firenet" "gcp_firenet" {
  vpc_id              = format("%s~-~%s", aviatrix_transit_gateway.test_transit_gateway.vpc_id, aviatrix_account.gcp.gcloud_project_id)
  inspection_enabled  = true
  egress_enabled      = true
  keep_alive_via_lan_interface_enabled = false
  manage_firewall_instance_association = false

  depends_on = [aviatrix_firewall_instance_association.association2]
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

* `egress_static_cidrs` - (Optional) List of egress static CIDRs. Egress is required to be enabled. Example: ["1.171.15.184/32", "1.171.15.185/32"]. Available as of provider version R2.19+.
* `east_west_inspection_excluded_cidrs` - (Optional) Network List Excluded From East-West Inspection. CIDRs to be excluded from inspection. Type: Set(String). Available as of provider version R2.19.5+.
* `fail_close_enabled` - (Optional/Computed) Enable Fail Close. When Fail Close is enabled, FireNet gateway drops all traffic when there are no firewalls attached to the FireNet gateways. Type: Boolean. Available as of provider version R2.19.5+.
* `tgw_segmentation_for_egress_enabled` - (Optional) Enable TGW segmentation for egress. Valid values: true or false. Default value: false. Available as of provider version R2.19+.
* `hashing_algorithm` - (Optional) Hashing algorithm to load balance traffic across the firewall. Valid values: "2-Tuple", "5-Tuple". Default value: "5-Tuple".
* `keep_alive_via_lan_interface_enabled` - (Optional) Enable Keep Alive via Firewall LAN Interface. Valid values: true or false. Default value: false. Available as of provider version R2.18.1+.
* `manage_firewall_instance_association` - (Optional) Enable this attribute to manage firewall associations in-line. If set to true, in-line `firewall_instance_association` blocks can be used. If set to false, all firewall associations must be managed via standalone `aviatrix_firewall_instance_association` resources. Default value: true. Valid values: true or false. Available in provider version R2.17.1+.

### Firewall Association

!> **WARNING:** Attribute `firewall_instance_association` has been deprecated as of provider version R2.18+ and will not receive further updates. Please set `manage_firewall_instance_association` to false, and use the standalone `aviatrix_firewall_instance_association` resource instead.

-> **NOTE:** `firewall_instance_association` - Associating a firewall instance with a Native GWLB enabled VPC is not supported in the in-line `firewall_instance_association` attribute. Please use the standalone `aviatrix_firewall_instance_association` resource instead.

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
