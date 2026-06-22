---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_distributed_firewalling_config"
description: |-
  Creates and manages an Aviatrix Distributed-firewalling Config
---

# aviatrix_distributed_firewalling_config

The **aviatrix_distributed_firewalling_config** resource allows management of an Aviatrix Distributed Firewalling configuration. This resource is available as of provider version R3.0+.

## Example Usage

```hcl
# Create an Aviatrix Distributed Firewalling config
resource "aviatrix_distributed_firewalling_config" "test" {
  enable_distributed_firewalling = true
}
```


## Argument Reference

The following arguments are supported:

### Required
* `enable_distributed_firewalling` - (Optional) Whether to enable Aviatrix Distributed Firewalling on an Aviatrix Controller. Valid values: true, false. Default value: false.

## Import

**aviatrix_distributed_firewalling_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_distributed_firewalling_config.test 10-11-12-13
```
