---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_sumologic_forwarder"
description: |-
  Enables and disables sumologic forwarder
---

# aviatrix_sumologic_forwarder

The **aviatrix_sumologic_forwarder** resource allows the enabling and disabling of sumologic forwarder.

## Example Usage

```hcl
# Enable sumologic forwarder
resource "aviatrix_sumologic_forwarder" "test_sumologic_forwarder" {
  access_id       = 0
  access_key      = "1.2.3.4"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `access_id` (Required) Access ID.
* `access_key` (Required) Access Key.

### Optional
* `source_category` (Optional) Source category.
* `custom_configuration` (Optional) Custom configuration. The format should be key=value pairs.
* `excluded_gateways` (Optional) List of gateways to be excluded from logging. e.g.: ["gateway01", "gateway02", "gateway01-hagw"].

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - The status of sumologic forwarder.

## Import

**sumologic_forwarder** can be imported using "sumologic_forwarder", e.g.

```
$ terraform import aviatrix_sumologic_forwarder.test sumologic_forwarder
```
