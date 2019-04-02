---
layout: "aviatrix"
page_title: "Provider: Aviatrix"
sidebar_current: "docs-aviatrix-index"
description: |-
  The Aviatrix provider is used to interact with Aviatrix organization resources.
---

# Aviatrix Provider

The Aviatrix provider is used to interact with Aviatrix organization resources.

The provider allows you to manage your Aviatrix organization's gateway,tunnels, and other resources easily.
It needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure Aviatrix provider
provider "aviatrix" {
  controller_ip = "1.2.3.4"
  username      = "admin"
  password      = "password"
}

# Create a record
resource "aviatrix_account" "myacc" {
  # ...
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `controller_ip` - (Required) This is Aviatrix controller's public IP. It must be provided.

* `username` - (Required) This is  Aviatrix account username which will be used to 
  login to Aviatrix controller. It must be provided.

* `password` - (Required) This is Aviatrix account's password corresponding to above username.

## Import

Instances can be imported using the id, e.g.

```
$ terraform import aviatrix_instance.test myAviatrixInstanceID
```