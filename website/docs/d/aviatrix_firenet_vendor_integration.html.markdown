---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firenet_vendor_integration"
description: |-
  Does 'save' or 'sync' for vendor integration purpose for Aviatrix FireNet.
---

# aviatrix_firenet_vendor_integration

Use this data source to do 'save' or 'sync' for vendor integration purpose for Aviatrix FireNet.

## Example Usage

```hcl
# Aviatrix FireNet Vendor Integration Data Source
data "aviatrix_firenet_vendor_integration" "foo" {
  vpc_id        = "vpc-abcd123"
  instance_id   = "i-09ade2592661316f8"
  vendor_type   = "Palo Alto VM Series"
  public_ip     = "10.11.12.13"
  username      = "admin"
  password      = "Avx123456#"
  firewall_name = "Avx-Firewall-Instance"
  save_enabled  = true
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) VPC ID.
* `instance_id` - (Required) ID of Firewall instance.
* `vendor_type` - (Required) Select PAN. Valid values: "Generic", "Palo Alto VM Series", "Palo Alto VM Panorama", "Aviatrix FQDN Gateway".
* `public_ip` - (Required) The public IP address of the firewall management interface for API calls from the Aviatrix Controller.
* `username` - (Required) Firewall login name for API calls from the Controller. For example, admin-api, as shown in the screen shot.
* `password` - (Required) Firewall login password for API calls.
* `firewall_name` - (Optional) Name of firewall instance.
* `route_table` - (Optional) Specify the firewall virtual Router name you wish the Controller to program. If left unspecified, the Controller programs the firewall’s default router.
* `save_enabled` - (Optional) Switch to save or not.
* `sync_enabled` - (Optional) Switch to sync or not.

