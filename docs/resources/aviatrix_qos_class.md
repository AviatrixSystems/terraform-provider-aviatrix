---
subcategory: "Copilot"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_qos_class"
description: |-
  Creates Aviatrix QoS Class
---

# aviatrix_qos_class

The **aviatrix_qos_class** resource creates a class for the Quality of Service (QoS) mechanism in Aviatrix Edge. This feature allows classification based on DSCP value in IP Header.

## Example Usage

```hcl
# Create a QoS Class
resource "aviatrix_qos_class" "test" {
  name     = "priority1"
  priority = 1
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Name of QoS class.
* `priority` - (Optional) Priority of QoS class, requires an integer value.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `uuid` - UUID of QoS class.

## Import

**qos_class** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_qos_class.test uuid
```
