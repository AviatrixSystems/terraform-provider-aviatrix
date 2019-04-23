---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_fqdn"
sidebar_current: "docs-aviatrix-resource-fqdn"
description: |-
  Manages Aviatrix FQDN filtering for Gateway
---

# aviatrix_fqdn

The FQDN resource manages FQDN filtering for Aviatrix Gateway

## Example Usage

```hcl
# Set Aviatrix Gateway FQDN filter
resource "aviatrix_fqdn" "test_fqdn" {
  fqdn_tag     = "my_tag"
  fqdn_status  = "enabled"
  fqdn_mode    = "white"
  gw_filter_tag_list = [
  {
    gw_name        = "gwTest1"
    source_ip_list = ["172.31.0.0/16", "172.31.0.0/20"]
  },
  {
    gw_name        = "gwTest2"
    source_ip_list = ["30.0.0.0/16"]
  },
  ]
  domain_names = [
  {
    fqdn  = "facebook.com"
    proto = "tcp"
    port  = "443"
  },
  {
    fqdn  = "reddit.com"
    proto = "tcp"
    port  = "443"
  }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `fqdn_tag` - (Required) FQDN Filter Tag Name.
* `fqdn_status` - (Optional) FQDN Filter Tag Status. Valid values: "enabled", "disabled".
* `fqdn_mode` - (Optional) Specify the tag color to be a white-list tag or black-list tag. Valid Values: "white", "black".
* `gw_filter_tag_list` - (Optional) A list of gateways to attach to the specific tag.
  * `gw_name` - (Optional) Name of the gateway to attach to the specific tag. 
  * `source_ip_list` - (Optional) List of source IPs in the VPC qualified for a specific tag.
* `domain_names` - (Optional) One or more domain names in a list with details as listed below:
  * `fqdn` - (Optional) FQDN. Example: "facebook.com" 
  * `proto` - (Optional) Protocol. Valid values: "all", "tcp", "udp", "icmp" 
  * `port` - (Optional) Port. Example "25" 
    * for protocol "all", port must be set to "all"
    * for protocol “icmp”, port must be set to “ping”

-> **NOTE:** 

* If you are using/ upgraded to Aviatrix Terraform Provider v4.2+ , and an fqdn resource was originally created with a provider version <4.2, you must modify your configuration file to match current format, and do ‘terraform refresh’ to update the state file to current format. 

## Import

Instance fqdn can be imported using the fqdn_tag, e.g.

```
$ terraform import aviatrix_fqdn.test fqdn_tag
```