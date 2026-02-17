---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_copilot_association"
description: |-
  Creates and manages a CoPilot Association
---

# aviatrix_copilot_association

The **aviatrix_copilot_association** resource allows management of controller CoPilot Association. This resource is available as of provider version R2.19+.

## Example Usage

```hcl
# Create a CoPilot Association
resource "aviatrix_copilot_association" "test_copilot_association" {
    copilot_address = "copilot.aviatrix.com"
}
```


## Argument Reference

The following arguments are supported:

* `copilot_address` - (Required) CoPilot instance IP Address or Hostname.

## Import

**aviatrix_copilot_association** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_copilot_association.test_copilot_association 10-11-12-13
```
