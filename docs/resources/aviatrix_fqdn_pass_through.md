---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_fqdn_pass_through"
description: |-
  Manages Aviatrix FQDN filter pass-through
---

# aviatrix_fqdn_pass_through

The **aviatrix_fqdn_pass_through** resource manages FQDN filter pass-through for Aviatrix gateways.

~> **NOTE:** The **aviatrix_fqdn_pass_through** resource must be created after the **aviatrix_fqdn** resource to be created successfully. To ensure that Terraform orders the creation of these resources correctly please use a `depends_on` meta-argument so that the **aviatrix_fqdn** resource is created before the **aviatrix_fqdn_pass_through** resource.

## Example Usage

```hcl
# Create an Aviatrix Gateway FQDN filter pass-through
resource "aviatrix_fqdn_pass_through" "test_fqdn_pass_through" {
  gw_name            = aviatrix_gateway.test_gw_aws.gw_name
  pass_through_cidrs = [
    "10.0.0.0/24",
    "10.0.1.0/24",
    "10.0.2.0/24",
  ]

  depends_on         = [aviatrix_fqdn.test_fqdn]
}
```

## Argument Reference

The following arguments are supported:

* `gw_name` - (Required) Gateway name to apply pass-through rules to.
* `pass_through_cidrs` - (Required) List of origin CIDR's to allow to pass-through FQDN filtering rules. Minimum list length: 1.

## Import

**fqdn_pass_through** can be imported using the `gw_name`, e.g.

```
$ terraform import aviatrix_fqdn_pass_through.test gw_name
```
