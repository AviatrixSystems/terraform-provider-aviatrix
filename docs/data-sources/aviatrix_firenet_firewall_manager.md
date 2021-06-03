---
subcategory: "Firewall Network"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firenet_firewall_manager"
description: |-
  Performs 'save' or 'sync' for Aviatrix FireNet firewall manager.
---

# aviatrix_firenet_firewall_manager

Use this data source to do 'save' or 'sync' for Aviatrix FireNet firewall manager.

## Example Usage

```hcl
# Aviatrix FireNet Firewall Manager Data Source
data "aviatrix_firenet_firewall_manager" "foo" {
  vpc_id         = "vpc-abcd123"
  gateway_name   = "transit"
  vendor_type    = "Palo Alto Networks Panorama"
  public_ip      = "1.2.3.4"
  username       = "admin-api"
  password       = "password"
  template       = "template"
  template_stack = "templatestack"
  route_table    = "router"
  save           = true
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) VPC ID.
* `gateway_name` - (Required) The FireNet gateway name.
* `vendor_type` - (Required) Vendor type. Valid values: "Generic" and "Palo Alto Networks Panorama".
* `public_ip` - (Optional) The public IP address of the Panorama instance. Required for vendor type "Palo Alto Networks Panorama". 
* `username` - (Optional) Panorama login name for API calls from the Controller. Required for vendor type "Palo Alto Networks Panorama".
* `password` - (Optional) Panorama login password for API calls. Required for vendor type "Palo Alto Networks Panorama".
* `template` - (Optional) Panorama template for each FireNet gateway. Required for vendor type "Palo Alto Networks Panorama".
* `template_stack` - (Optional) Panorama template stack for each FireNet gateway. Required for vendor type "Palo Alto Networks Panorama".
* `route_table` - (Optional) The name of firewall virtual router to program. If left unspecified, the Controller programs the Panorama template’s first router.
* `number_of_retries` - (Optional) Number of retries for `save` or `synchronize`. Example: 1. Default value: 0.
* `retry_interval` - (Optional) Retry interval in seconds for `save` or `synchronize`. Example: 120. Default value: 300.
* `save` - (Optional) Switch to save or not.
* `synchronize` - (Optional) Switch to sync or not.
