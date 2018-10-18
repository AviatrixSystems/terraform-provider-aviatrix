---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_firewall_tag"
sidebar_current: "docs-aviatrix-resource-firewall-tag"
description: |-
  Creates and manages Aviatrix Firewall Tags
---

# aviatrix_firewall_tag

The FirewallTag resource allows the creation and management of Aviatrix Firewall Tags.

## Example Usage

```hcl
# Create Aviatrix Firewall Tag
resource "aviatrix_firewall_tag" "test_firewall_tag" {
  firewall_tag = "test-firewall-tag"
  cidr_list = [
                {
                  cidr_tag_name = "a1"
                  cidr = "10.1.0.0/24"
                },
                {
                  cidr_tag_name = "b1"
                  cidr = "10.2.0.0/24"
                }
              ]
}

### review write up of cidr_list (cidr_tag_name and cidr)
### in API cidr_list looks for an array-like arguments. What are we looking for here?
}
```

## Argument Reference

The following arguments are supported:

* `firewall_tag` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller. Type: String
* `cidr_list` - (Optional) A JSON file with the following:
    * `cidr_tag_name` - The name attribute of a policy. Example: "policy1"
    * `cidr` - The CIDR attribute of a policy. Example: "10.88.88.88/32"
