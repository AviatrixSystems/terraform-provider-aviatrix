---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_distributed_firewalling_origin_cert_enforcement_config"
description: |-
  Creates and manages an Aviatrix Distributed-firewalling Origin Cert Enforcement Config
---

# aviatrix_distributed_firewalling_origin_cert_enforcement_config

The **aviatrix_distributed_firewalling_origin_cert_enforcement_config** resource allows management of an Aviatrix Distributed Firewalling Origin Cert Enforcement level configuration. This resource is available as of provider version R3.1.1+.

## Example Usage

```hcl
# Create an Aviatrix Distributed Firewalling Origin Cert Enforcement config
resource "aviatrix_distributed_firewalling_origin_cert_enforcement_config" "test" {
  enforcement_level = "Strict"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `enforcement_level` - (Optional) Which origin cert enforcement level to set to for distributed firewalling on an Aviatrix Controller. Valid values: "Strict", "Permissive" and "Ignore". Default value: "Permissive".

## Import

**aviatrix_distributed_firewalling_origin_cert_enforcement_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_distributed_firewalling_origin_cert_enforcement_config.test 10-11-12-13
```
