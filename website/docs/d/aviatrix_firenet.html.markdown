---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firenet"
description: |-
  Gets the Aviatrix FireNet.
---

# aviatrix_firenet

Use this data source to get the Aviatrix firenet for use in other resources.

## Example Usage

```hcl
# Aviatrix FireNet Data Source
data "aviatrix_firenet" "foo" {
	vpc_id = "vpc-abcdef"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) ID of the Security VPC.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_id` - ID of the Security VPC.
* `inspection_enabled` - Enable/Disable traffic inspection.
* `egress_enabled` - Enable/Disable egress through firewall.
* `firewall_instance_association` - List of firewall instances associated with fireNet.
  * `firenet_gw_name` - Name of the primary FireNet gateway.
  * `instance_id` - ID of Firewall instance.
  * `vendor_type` - Type of the firewall.
  * `firewall_name` - Firewall instance name.
  * `lan_interface`- Lan interface ID.
  * `management_interface` - Management interface ID.
  * `egress_interface`- Egress interface ID.
  * `attached`- Switch to attach/detach firewall instance to/from fireNet.

