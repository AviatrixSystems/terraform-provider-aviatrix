---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_distributed_firewalling_proxy_ca_config"
description: |-
  Creates and manages an Aviatrix Distributed-firewalling Proxy CA Config
---

# aviatrix_distributed_firewalling_proxy_ca_config

The **aviatrix_distributed_firewalling_proxy_ca_config** resource allows management of an Aviatrix Distributed Firewalling Proxy CA configuration. This resource is available as of provider version R3.1.1+.

## Example Usage

```hcl
# Create an Aviatrix Distributed Firewalling Proxy CA config
resource "aviatrix_distributed_firewalling_proxy_ca_config" "test" {
  ca_cert = file("ca_cert_file")
  ca_key  = file("ca_key_file")
}
```

## Argument Reference

The following arguments are supported:

### Required
* `ca_cert` - (Required) Content of proxy ca certificate to create only one cert.
* `ca_key` - (Required) Content of proxy ca cert key to create only one cert.

In addition to all arguments above, the following attributes are exported:

* `common_name` - (Computed) Common name of created cert.
* `expiration_time` - (Computed) Expiration time of created cert.
* `issuer_name` - (Computed) Issuer name of created cert.
* `unique_serial` - (Computed) Unique serial of created cert.
* `upload_info` - (Computed) Upload info of created cert.

## Import

**aviatrix_distributed_firewalling_origin_cert_enforcement_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_distributed_firewalling_origin_cert_enforcement_config.test 10-11-12-13
```
