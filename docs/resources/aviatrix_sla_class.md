---
subcategory: "Copilot"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_sla_class"
description: |-
  Creates Aviatrix SLA Class
---

# aviatrix_sla_class

The **aviatrix_sla_class** resource creates the Aviatrix SLA Class.

## Example Usage

```hcl
# Create a SLA Class
resource "aviatrix_sla_class" "test" {
  name             = "gold"
  latency          = 43
  jitter           = 1
  packet_drop_rate = 3
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Name of SLA class.

### Optional
* `latency` - (Optional) Latency of sla class in ms.
* `jitter` - (Optional) Jitter of sla class in ms.
* `packet_drop_rate` - (Optional) Packet drop rate of sla class.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `uuid` - UUID of SLA class.

## Import

**sla_class** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_sla_class.test uuid
```
