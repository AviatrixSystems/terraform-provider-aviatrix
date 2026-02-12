---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_trustbundle"
description: |-
    Gets details about a DCF Trust Bundle.
---

# aviatrix_dcf_trustbundle

The **aviatrix_dcf_trustbundle** data source provides details about a specific DCF Trust Bundle created by the Aviatrix Controller.

This data source can be useful when you need to reference an existing trust bundle in other resources, such as TLS profiles.

## Example Usage

```hcl
# Aviatrix DCF Trust Bundle Data Source
data "aviatrix_dcf_trustbundle" "example" {
    display_name = "my-trust-bundle"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the DCF Trust Bundle.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `bundle_id` - ID of the DCF Trust Bundle. This is typically used as `target_uuid` in DCF policies.
* `bundle_content` - Content of the DCF Trust Bundle as a string. Contains the certificate data.
* `created_at` - Timestamp when the DCF Trust Bundle was created.
