---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_tag"
description: |-
  Creates and manages Aviatrix Firewall Tags
---

# aviatrix_firewall_tag

The aviatrix_firewall_tag resource allows the creation and management of Aviatrix Firewall tags.

## Example Usage

```hcl
# Create an Aviatrix Firewall Tag
resource "aviatrix_firewall_tag" "test_firewall_tag" {
  firewall_tag = "test-firewall-tag"

  cidr_list {
    cidr_tag_name = "a1"
    cidr          = "10.1.0.0/24"
  }

  cidr_list {
    cidr_tag_name = "b1"
    cidr          = "10.2.0.0/24"
  }
}
```

## Argument Reference

The following arguments are supported:

* `firewall_tag` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller.
* `cidr_list` - (Optional) A JSON file with the following:
  * `cidr_tag_name` - (Required) The name attribute of a policy. Example: "policy1".
  * `cidr` - (Required) The CIDR attribute of a policy. Example: "10.88.88.88/32".

## Import

Instance firewall_tag can be imported using the firewall_tag, e.g.

```
$ terraform import aviatrix_firewall_tag.test firewall_tag
```
