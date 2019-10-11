---
layout: "aviatrix"
page_title: "Provider: Aviatrix"
description: |-
  The Aviatrix provider is used to interact with Aviatrix organization resources
---

# Aviatrix Provider

The Aviatrix provider is used to interact with Aviatrix organization resources.

This provider allows you to manage your Aviatrix organization's gateways, tunnels, and other resources easily. It needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

-> **NOTE:** Although *version* is an optional argument, we highly recommend users to specify the proper Aviatrix provider release version corresponding to their Controller version in order to avoid potential compatibility issues. Please see the [compatibility chart](https://www.terraform.io/docs/providers/aviatrix/guides/release-compatibility.html) for full details. For more information on versioning, a native Terraform provider argument, see [here](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).


## Example Usage

```hcl
# Configure Aviatrix provider
provider "aviatrix" {
  controller_ip           = "1.2.3.4"
  username                = "admin"
  password                = "password"
  skip_version_validation = false
  version                 = "2.5.0"
}

# Create an access account
resource "aviatrix_account" "myacc" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `controller_ip` - (Required) Aviatrix controller's public IP.
* `username` - (Required) Aviatrix account username which will be used to login to Aviatrix controller.
* `password` - (Required) Aviatrix account password corresponding to above username.
* `skip_version_validation` - (Optional) Default: false. If set to true, it skips checking whether current Terraform provider supports current Controller version.
* `version` - (Optional) - Specify Aviatrix provider release version number. If not specified, Terraform will automatically pull and source the latest release.

## Import

Instances can be imported using the id, e.g.

```
$ terraform import aviatrix_instance.test myAviatrixInstanceID
```
