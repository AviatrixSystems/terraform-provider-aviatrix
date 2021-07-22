---
subcategory: "Gateway"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway_image"
description: |- Gets an Aviatrix gateway image version details.
---

# aviatrix_gateway_image

The **aviatrix_gateway_image** data source provides the current image version that pairs with the given software version
and cloud type.

This data source is useful for getting the correct image_version for a gateway when upgrading the software_version of
the gateway.

## Example Usage

```hcl
# Aviatrix Gateway Image Data Source using interpolation from a spoke gateway
data "aviatrix_gateway" "foo" {
  cloud_type       = aviatrix_spoke_gateway.spoke.cloud_type
  software_version = aviatrix_spoke_gateway.spoke.software_version
}
```

```hcl
# Aviatrix Gateway Image Data Source
data "aviatrix_gateway" "foo" {
  cloud_type       = 4
  software_version = "6.4.2487"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_type` - (Required) Cloud type. Type: Integer. Example: 1 (AWS)
* `software_version` - (Required) Software version. Type: String. Example: "6.4.2487"

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `image_version` - Image version that is compatible with the given cloud_type and software_version.
