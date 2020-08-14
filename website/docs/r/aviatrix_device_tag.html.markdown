---
subcategory: "CloudWAN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_device_tag"
description: |-
  Creates and manages device config tags
---

# aviatrix_device_tag

The **aviatrix_device_tag** resource allows the creation and management of device config tags.

~> **NOTE:** Creating this resource will automatically commit the config to the specified devices.

## Example Usage

```hcl
# Create an Aviatrix Device Tag and commit it
resource "aviatrix_device_tag" "test_device_tag" {
  name                = "tag_hostname"
  config              = <<EOT
hostname myrouter
EOT
  device_names        = [aviatrix_device_registration.test_device.name]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Name of the tag.
* `config` - (Required) Config to apply to devices that are attached to the tag.
* `device_names` - (Required) List of device names to attach to this tag.

## Import

**device_tag** can be imported using the `name`, e.g.

```
$ terraform import aviatrix_device_tag.test name
```
