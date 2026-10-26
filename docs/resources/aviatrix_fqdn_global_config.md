---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_fqdn_global_config"
description: |-
  Manages Aviatrix FQDN Global Config
---

# aviatrix_fqdn_global_config

The **aviatrix_fqdn_global_config** resource manages FQDN global config.

## Example Usage

```hcl
# Create an Aviatrix FQDN Global Config with Private Network Filtering enabled
resource "aviatrix_fqdn_global_config" "test" {
  enable_exception_rule            = false
  enable_private_network_filtering = true
  enable_caching                   = false
  enable_exact_match               = true
}
```
```hcl
# Create an Aviatrix FQDN Global Config with Custom Network Filtering enabled
resource "aviatrix_fqdn_global_config" "test" {
  enable_exception_rule           = false
  enable_custom_network_filtering = true
  configured_ips                  = [
    "172.16.0.0/12~~RFC-1918",
    "10.0.0.0/8~~RFC-1918",
    "168.16.0.0/32"
  ]
  enable_caching                  = false
  enable_exact_match              = true
}
```

## Argument Reference

The following arguments are supported:

* `enable_exception_rule` - (Optional) If enabled, it allows packets passing through the gateway without an SNI field. Only applies to whitelist. Valid values: true, false. Default value: true.
* `enable_private_network_filtering` - (Optional) If enabled, destination FQDN names that translate to private IP address range (RFC 1918) are subject to FQDN whitelist filtering function. Valid values: true, false. Default value: false.
* `enable_custom_network_filtering` - (Optional) If enabled, it customizes packet destination address ranges not to be filtered by FQDN. Valid values: true, false. Default value: false.
* `configured_ips` - (Optional) Customized packet destination address ranges not to be filtered by FQDN. Can be selected from pre-defined RFC 1918 range, or own network range. Required with `enable_custom_network_filtering` enabled.
* `enable_caching` - (Optional) If enabled, it caches the resolved IP address from FQDN filter. Valid values: true, false. Default value: true.
* `enable_exact_match` - (Optional) If enabled, the resolved IP address from FQDN filter is cached so that if subsequent TCP session matches the cached IP address list, FQDN domain name is not checked and the session is allowed to pass. Valid values: true, false. Default value: false.

## Import

Instance fqdn_global_config can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_fqdn_global_config.test 10-11-12-13
```
