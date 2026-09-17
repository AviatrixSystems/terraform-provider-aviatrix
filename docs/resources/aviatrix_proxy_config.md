---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_proxy_config"
description: |-
  Creates and manages an Aviatrix controller proxy config resource
---

# aviatrix_proxy_config

The **aviatrix_proxy_config** resource allows management of an Aviatrix Controller's proxy configurations.

## Example Usage

```hcl
# Create an Aviatrix Controller Proxy Config
resource "aviatrix_proxy_config" "test_proxy_config" {
  http_proxy           = "172.31.52.145:3127"
  https_proxy          = "172.31.52.145:3129"
  proxy_ca_certificate = file("/path/to/ca.pem")
}
```

## Argument Reference

The following arguments are supported:

* `http_proxy` - (Required) Http proxy URL.
* `https_proxy` - (Required) Https proxy URL.
* `proxy_ca_certificate` - (Optional) Server CA Certificate file. Use the `file` function to read from a file.

## Import

**controller_proxy_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_proxy_config.test 10-11-12-13
```
