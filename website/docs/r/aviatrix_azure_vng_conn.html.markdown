---
subcategory: "Multi-Cloud Transit"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_azure_vng_conn"
description: |-
Creates and manages the connection between Aviatrix Transit Gateway and Azure VNG
---

# aviatrix_azure_vng_conn

The **aviatrix_azure_vng_conn** resource allows the creation and management of the connection between Aviatrix Transit Gateway and Azure VNG.

## Example Usage

```hcl
# Attach an Azure VNG to an Aviatrix Transit Gateway
resource "aviatrix_azure_vng_conn" "test" {
  primary_gateway_name = "primary-gateway"
  connection_name      = "connection"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `primary_gateway_name` - (Required) Primary Aviatrix transit gateway name.
* `connection_name` - (Required) Connection name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_id` - VPC ID.
* `vng_name` - Name of Azure VNG.
* `attached` - The status of the connection.

## Import

**transit_gateway** can be imported using the `connection_name`, e.g.

```
$ terraform import aviatrix_azure_vng_conn.test connection
```
