---
subcategory: "Copilot"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_traffic_classifier"
description: |-
  Creates Aviatrix Traffic Classifier
---

# aviatrix_traffic_classifier

The **aviatrix_traffic_classifier** resource creates the Aviatrix Traffic Classifier.

!> **WARNING:** Creating the **aviatrix_traffic_classifier** resource will overwrite all the traffic classifier policies. Deleting the **aviatrix_traffic_classifier** resource will remove all the traffic classifier policies.

## Example Usage

```hcl
# Create a Traffic Classifier
resource "aviatrix_traffic_classifier" "test" {
  policies {
    name                          = "traffic_classifier_policy"
    source_smart_group_uuids      = ["<<source smart group uuid>>"]
    destination_smart_group_uuids = ["<<destination smart group uuid>>"]
    link_hierarchy_uuid           = "<<link hierarchy uuid>>"
    sla_class_uuid                = "<<sla class uuid>>"
    port_ranges {
      low  = 10
      high = 50
    }
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `policies` - (Required) List of traffic classifier policies.
  * `name` - (Required) Name of traffic classifier.
  * `source_smart_group_uuids` - (Required) List of source smart group UUIDs.
  * `destination_smart_group_uuids` - (Required) List of destination smart group UUIDs.
  * `port_ranges` - (Optional) Port ranges.
    * `low` - (Optional) Low port range.
    * `high` - (Optional) High port range.
  * `protocol` - (Optional) Protocol.
  * `link_hierarchy_uuid` - (Optional) Link hierarchy UUID.
  * `sla_class_uuid` - (Optional) SLA class UUID.
  * `enable_logging` - (Optional) Enable logging.
  * `route_type` - (Optional) Route type.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `uuid` in `policies` - Traffic classifier policy UUID.

## Import

**traffic_classifier** can be imported using the ID "traffic_classifier_policies", e.g.

```
$ terraform import aviatrix_traffic_classifier.test traffic_classifier_policies
```
