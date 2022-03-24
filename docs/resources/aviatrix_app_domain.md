---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_app_domain"
description: |-
Creates and manages an Aviatrix App Domain
---

# aviatrix_app_domain

The **aviatrix_app_domain** resource handles the creation and management of App Domains. Available as of Provider R2.22.0+

## Example Usage

```hcl
# Create an Aviatrix App Domain with IP Filters
resource "aviatrix_app_domain" "test_app_domain_ip" {
  name      = "app-domain"
  ip_filter = [
    "10.0.0.0/16",
    "11.0.0.0/16"
  ]
}
```

```hcl
# Create an Aviatrix App Domain with Tag Filters
resource "aviatrix_app_domain" "test_app_domain_ip" {
  name       = "app-domain"
  tag_filter = {
    k1 = "v1"
    k2 = "v2"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the App Domain.

### Filters

* `ip_filter` - (Optional) List of IP CIDRs to filter the App Domain.
* `tag_filter` - (Optional) Map of key-value tags to filter the App Domain. 


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - UUID of the App Domain.

## Import

**aviatrix_app_domain** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_app_domain.test 41984f8b-5a37-4272-89b3-57c79e9ff77c
```
