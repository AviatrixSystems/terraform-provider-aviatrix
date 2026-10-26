---
subcategory: "Copilot"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_link_hierarchy"
description: |-
  Creates Aviatrix Link Hierarchy
---

# aviatrix_link_hierarchy

The **aviatrix_link_hierarchy** resource creates the Aviatrix Link Hierarchy.

## Example Usage

```hcl
# Create a Link Hierarchy
resource "aviatrix_link_hierarchy" "test" {
  name = "test"

  links {
    name = "test1"
    wan_link {
      wan_tag = "wan3.10"
    }
    wan_link {
      wan_tag = "wan3.11"
    }
  }

  links {
    name = "test2"
    wan_link {
      wan_tag = "wan4.10"
    }
  }
}

```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Link hierarchy name.

### Optional
* `links` - (Optional) List of named links.
    * `name` - (Optional) Link name.
    * `wan_link` - (Optional) Set of WAN links.
        * `wan_tag` - (Optional) WAN tag.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `uuid` - UUID of link hierarchy.

## Import

**link_hierarchy** can be imported using the `uuid`, e.g.

```
$ terraform import aviatrix_link_hierarchy.test uuid
```
