---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_tag"
description: |-
  Creates and manages Aviatrix Stateful Firewall Tags
---

# aviatrix_firewall_tag

The **aviatrix_firewall_tag** resource allows the creation and management of [Aviatrix Stateful Firewall tags](https://docs.aviatrix.com/HowTos/tag_firewall.html) for tag-based security for gateways.

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

### Required
* `firewall_tag` - (Required) Name of the stateful firewall tag to be created.

### Tag Rules
* `cidr_list` - (Optional) Dynamic block representing a CIDR to filter, and a name to identify it:
  * `cidr_tag_name` - (Required) A name to identify the CIDR. Example: "policy1".
  * `cidr` - (Required) CIDR address to filter. Example: "10.88.88.88/32".

## Import

**firewall_tag** can be imported using the `firewall_tag`, e.g.

```
$ terraform import aviatrix_firewall_tag.test firewall_tag
```
