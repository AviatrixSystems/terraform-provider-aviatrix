---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_customer_id"
sidebar_current: "docs-aviatrix-resource-customer_id"
description: |-
  Sets Aviatrix CustomerID and License.
---

# aviatrix_customer_id

The CustomerID resource allows to set Aviatrix Customer ID and License

## Example Usage

```hcl
# Set Aviatrix Customer ID and License
resource "aviatrix_customer_id" "test_customer_id" {
  customer_id = "paloaltodev-1234567898.64"
}
```

## Argument Reference

The following arguments are supported:

* `customer_id` - (Required) The license ID provided by Aviatrix Systems
