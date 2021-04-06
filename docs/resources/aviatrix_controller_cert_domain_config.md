---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_cert_domain_config"
description: |-
  Creates and manages an Aviatrix controller cert domain config
---

# aviatrix_controller_cert_domain_config

The **aviatrix_controller_cert_domain_config** resource allows management of an Aviatrix Controller's cert domain config. This resource is available as of provider version R2.19+.

## Example Usage

```hcl
# Create an Aviatrix controller cert domain config
resource "aviatrix_controller_cert_domain_config" "test" {
  cert_domain = "abc.com"
}
```


## Argument Reference

The following argument is supported:

* `cert_domain` - (Optional) Domain name that is used in FQDN for generating cert. Default value: "aviatrixnetwork.com".

## Import

**aviatrix_controller_cert_domain_config** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_controller_cert_domain_config.test 10-11-12-13
```
