---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_caller_identity"
description: |-
  Gets the Aviatrix caller identity.
---

# aviatrix_caller_identity

The **aviatrix_caller_identity** data source provides the Aviatrix CID for use in other resources.

## Example Usage

```hcl
# Aviatrix Caller Identity Data Source
data "aviatrix_caller_identity" "foo" {

}
```

## Argument Reference

The following arguments are supported:

* None.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `cid` - Aviatrix caller identity.
