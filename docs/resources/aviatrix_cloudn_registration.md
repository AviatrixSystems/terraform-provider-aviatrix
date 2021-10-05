---
subcategory: "CloudWAN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_cloudn_registration"
description: |-
Register and Deregister CloudN as a Gateway
---

# aviatrix_cloudn_registration

The **aviatrix_cloudn_registration** resource allows the registration and deregistration of Aviatrix CloudN as a Gateway. This resource is available as of provider version R2.21+.

## Example Usage

```hcl
# Create a CloudN Registration
resource "aviatrix_cloudn_registration" "test_cloudn_registration" {
  name            = "cloudn-test"
  username        = "admin"
  password        = "password"
  address         = "10.210.38.100"
  local_as_number = "65000"
  prepend_as_path = [
    "65000",
    "65000",
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Gateway name to assign to the CloudN device
* `address` - (Required) Aviatrix CloudN's public or private IP.
* `username` - (Required) Aviatrix account username which will be used to log in to Aviatrix CloudN.
* `password` - (Required) Aviatrix account password corresponding to above username.

### Optional
* `local_as_number` - (Optional) BGP AS Number to assign to the Transit Gateway. Type: String.
* `prepend_as_path` - (Optional) Connection AS Path Prepend customized by specifying AS PATH for a BGP connection. Requires local_as_number to be set. Type: List.

## Import

**aviatrix_cloudn_registration** can be imported using the `name`, e.g.

```
$ terraform import aviatrix_cloudn_registration.test_cloudn_registration name
```
