---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_caller_identity"
description: |-
  Gets the an Aviatrix caller identity.
---

# aviatrix_caller_identity

Use this data source to get the Aviatrix caller identity for use in other resources.

## Example Usage

```hcl
# Aviatrix caller identity data source
data "aviatrix_caller_identity" "foo" {

}
```

## Argument Reference

The following arguments are supported:

* None.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `cid` - Aviatrix caller identity.
