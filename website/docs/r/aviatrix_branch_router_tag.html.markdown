---
subcategory: "CloudWAN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_branch_router_tag"
description: |-
  Creates and manages branch router config tags
---

# aviatrix_branch_router_tag

The **aviatrix_branch_router_tag** resource allows the creation and management of branch router config tags.

~> **NOTE:** Creating this resource will automatically commit the config to the specified branch routers.

## Example Usage

```hcl
# Create an Aviatrix Branch Router Tag and commit it
resource "aviatrix_branch_router_tag" "test_branch_router_tag" {
  name                = "tag_hostname"
  config              = <<EOT
hostname myrouter
EOT
  branch_router_names = [aviatrix_branch_router.test_branch_router.name]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Name of the tag.
* `config` - (Required) Config to apply to branches that are attached to the tag.
* `branch_router_names` - (Required) List of branch names to attach to this tag.

## Import

**branch_router_tag** can be imported using the `name`, e.g.

```
$ terraform import aviatrix_branch_router_tag.test name
```
