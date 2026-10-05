---
subcategory: "Copilot"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_qos_policy_list"
description: |-
  Creates Aviatrix QoS Policy List
---

# aviatrix_qos_policy_list

The **aviatrix_qos_policy_list** resource creates a list of policies (and rules) under the Quality of Service (QoS) mechanism in Aviatrix Edge. This is to be used in conjunction with **aviatrix_qos_class**es to classify based on DSCP value in IP Header.

!> **WARNING:** Creating the **aviatrix_qos_policy_list** resource will overwrite all the QoS policies. Deleting the **aviatrix_qos_policy_list** resource will remove all the QoS policies.

## Example Usage

```hcl
# Create a QoS Policy List
resource "aviatrix_qos_class" "test" {
  name     = "test-qos-class"
  priority = 3
}

resource "aviatrix_qos_policy_list" "test" {
  policies {
    name           = "qos_policy_1"
    dscp_values    = ["1", "AF11"]
    qos_class_uuid = aviatrix_qos_class.test.uuid
  }

  policies {
    name           = "qos_policy_2"
    dscp_values    = ["AF22"]
    qos_class_uuid = aviatrix_qos_class.test.uuid
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `policies` - (Required) List of QoS policies.
  * `name` - (Required) Name of QoS class.
  * `dscp_values` - (Required) List of DSCP values.
  * `qos_class_uuid` - (Required) QoS class UUID.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `uuid` in `policies` - UUID of QoS policy.

## Import

**qos_policy_list** can be imported using the "qos_policy_list", e.g.

```
$ terraform import aviatrix_qos_policy_list.test qos_policy_list
```
