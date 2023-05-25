---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_edge_platform_device_onboarding"
description: |-
  Onboards an Edge Platform device
---

# aviatrix_edge_platform_device_onboarding

The **aviatrix_edge_platform_device_onboarding** resource onboards the Edge Platform device.

## Example Usage

```hcl
# Onboard an Edge Platform device
resource "aviatrix_edge_platform_device_onboarding" "test" {
  account_name   = "edge-platform-acc"
  device_name    = "device0"
  serial_number  = "serial-123"
  hardware_model = "model-456"

  network {
    interface_name  = "eth5"
    enable_dhcp     = false
    ipv4_cidr       = "172.16.15.162/20"
    gateway_ip      = "172.16.0.1"
    dns_server_ips  = ["172.16.0.1"]
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `account_name` - (Required) Edge Platform account name.
* `device_name` - (Required) Edge Platform device name.
* `serial_number` - (Required) Device serial number.
* `hardware_model` - (Required) Device hardware model.


### Optional
* `management_egress_ip_prefix_list` - (Optional) Set of management egress gateway IP and subnet prefix. Example: ["67.207.104.16/29", "64.71.12.144/29"].
* `network` - (Optional) Network configurations.
  * `interface_name` - (Required) Interface name.
  * `enable_dhcp` - (Optional) Set to true to enable DHCP. Valid values: true, false. Default value: false.
  * `ipv4_cidr` - (Optional) Interface IPV4 CIDR.
  * `gateway_ip` - (Optional) Gateway IP.
  * `dns_server_ips` - (Optional) Set of DNS server IPs.
  * `proxy_server_ip` - (Optional) Proxy server IP.
* `download_config_file` - (Optional) Set to true to download the Edge Platform static config file. Valid values: true, false. Default value: false.
* `config_file_download_path` - (Optional) The location where the config file will be stored.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `device_id` - Edge Platform device ID.

## Import

**edge_platform_device_onboarding** resource can be imported with the `account_name` and `device_name` in the form "account_name~device_name", e.g.

```
$ terraform import aviatrix_edge_platform_device_onboarding.test account_naem~device_name
```
