---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_controller_metadata"
description: |-
  Gets the Aviatrix controller metadata.
---

# aviatrix_controller_metadata

The **aviatrix_controller_metadata** data source provides the controller metadata for use in other resources.

## Example Usage

```hcl
# Aviatrix Controller Metadata Data Source
data "aviatrix_controller_metadata" "foo" {
}
```

## Argument Reference

The following arguments are supported:

* None.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `region` - Controller region.
* `vpc_id` - Controller VPC ID.
* `instance_id` - Controller instance ID.
* `cloud_type` - Controller cloud type, only supported for AWS and GCP now.
