---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_data_caller_identity"
sidebar_current: "docs-aviatrix-data_source-caller_identity"
description: |-
  Gets the an Aviatrix caller identity.
---

# aviatrix_caller_identity

Use this data source to get the Aviatrix caller identity for use in other resources.

## Example Usage

```hcl
# Create Aviatrix caller identity data source
data "aviatrix_caller_identity" "foo" {

}
```

## Argument Reference

The following arguments are supported:

* None.

## Attribute Reference

* `cid` - (Computed) Aviatrix caller identity.
