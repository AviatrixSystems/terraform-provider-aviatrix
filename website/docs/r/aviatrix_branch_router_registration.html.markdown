---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_branch_router_registration"
description: |-
  Creates and manages branch router entries for CloudWAN
---

# aviatrix_branch_router_registration

The **aviatrix_branch_router_registration** resource allows the creation and management of branch router entries for CloudWAN.

~> **NOTE:** Before this branch router can be attached to any Aviatrix Transit Gateway, AWS TGW or Azure Virtual WAN you must configure its WAN interface and IP via the `aviatrix_branch_router_interface_config` resource.

## Example Usage

```hcl
# Create an Aviatrix Branch Router entry with private key authentication
resource "aviatrix_branch_router_registration" "test_branch_router" {
  name      = "test-branch-router"
  public_ip = "58.151.114.231"
  username  = "ec2-user"
  key_file  = "/path/to/key_file.pem"
}
```

```hcl
# Create an Aviatrix Branch Router entry with password authentication
resource "aviatrix_branch_router_registration" "test_branch_router" {
  name      = "test-branch-router"
  public_ip = "58.151.114.231"
  username  = "ec2-user"
  password  = "secret"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Name of the router.
* `public_ip` - (Required) Public IP address of the router.
* `username` - (Required) Username for SSH into the router.
* `key_file` - (Optional) Path to private key file for SSH into the router. Either `key_file` or `password` must be set to create a branch router successfully.
* `password` - (Optional) Password for SSH into the router. Either `key_file` or `password` must be set to create a branch router successfully. This attribute can also be set via environment variable 'AVIATRIX_BRANCH_ROUTER_PASSWORD'. If both are set, the value in the config file will be used.

### Optional
* `host_os` - (Optional) Router host OS.  Default value is 'ios'. Only valid value is 'ios'.
* `ssh_port` - (Optional) SSH port for connecting to the router. Default value is 22.
* `address_1` - (Optional) Address line 1.
* `address_2` - (Optional) Address line 2.
* `city` - (Optional) City.
* `state` - (Optional) State.
* `country` - (Optional) ISO two-letter country code.
* `zip_code` - (Optional) Zip code.
* `description` - (Optional) Description.

## Import

**branch_router_registration** can be imported using the `name`, e.g.

```
$ terraform import aviatrix_branch_router_registration.test name
```
