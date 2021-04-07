---
subcategory: "Firewall Network"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firenet_vendor_integration"
description: |-
  Performs 'save' or 'sync' for vendor integration purposes for Aviatrix FireNet.
---

# aviatrix_firenet_vendor_integration

Use this data source to do 'save' or 'sync' for vendor integration purposes for Aviatrix FireNet.

-> **NOTE:** FireNet with Panorama should be set up using the **aviatrix_firenet_firewall_manager** data source. Do not use `save` or `sync` options listed below.

~> **NOTE:** **aviatrix_firenet_firewall_manager** is currently under development.

## Example Usage

```hcl
# Aviatrix FireNet Vendor Integration Data Source
data "aviatrix_firenet_vendor_integration" "foo" {
  vpc_id        = "vpc-abcd123"
  instance_id   = "i-09ade2592661316f8"
  vendor_type   = "Palo Alto Networks VM-Series"
  public_ip     = "10.11.12.13"
  username      = "admin"
  password      = "Avx123456#"
  firewall_name = "Avx-Firewall-Instance"
  save          = true
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) VPC ID.
* `instance_id` - (Required) ID of Firewall instance.
* `vendor_type` - (Required) Select PAN. Valid values: "Generic", "Palo Alto Networks VM-Series", "Aviatrix FQDN Gateway" and "Fortinet FortiGate".
* `public_ip` - (Required) The public IP address of the firewall management interface for API calls from the Aviatrix Controller.
* `username` - (Optional) Firewall login name for API calls from the Controller. Required for vendor type "Generic", "Palo Alto Networks VM-Series" and "Aviatrix FQDN Gateway".
* `password` - (Optional) Firewall login password for API calls. Required for vendor type "Generic", "Palo Alto Networks VM-Series" and "Aviatrix FQDN Gateway".
* `api_token` - (Optional) API token for API calls. Required for vendor type "Fortinet FortiGate".  
* `firewall_name` - (Optional) Name of firewall instance.
* `route_table` - (Optional) Specify the firewall virtual Router name you wish the Controller to program. If left unspecified, the Controller programs the firewall’s default router.
* `number_of_retries` - (Optional) Number of retries for `save` or `synchronize`. Example: 1. Default value: 0.
* `retry_interval` - (Optional) Retry interval in seconds for `save` or `synchronize`. Example: 120. Default value: 300.
* `save` - (Optional) Switch to save or not.
* `synchronize` - (Optional) Switch to sync or not.
