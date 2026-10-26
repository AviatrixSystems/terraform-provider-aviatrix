---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_gateway_certificate_config"
description: |-
  Manages Aviatrix gateway certificate configuration
---

# aviatrix_gateway_certificate_config

The **aviatrix_gateway_certificate_config** resource allows the management of Aviatrix [gateway certificate](https://docs.aviatrix.com/HowTos/controller_certificate.html#gateway-certificate-management) configuration. Available as of provider version R2.18.1+.

## Example Usage

```hcl
# Aviatrix Gateway Certificate Management
resource "aviatrix_gateway_certificate_config" "test_gateway_cert" {
  ca_certificate = file("path/to/CA_cert.pem")
  ca_private_key = file("path/to/CA_private.key")
}
```

## Argument Reference

The following arguments are supported:

### Required
* `ca_certificate` - (Required) CA Certificate in PEM format. To read certificate from a file please use the built-in `file` function.
* `ca_private_key` - (Required/Sensitive) CA Private Key. To read the private key from a file please use the built-in `file` function.


## Import

!> **WARNING:** When importing, the provider cannot read your private key or certificate into the state file. After importing, if you do not want to change the values of the CA private key or certificate you must set the attributes `ca_certificate` and `ca_private_key` to the empty string (""). Otherwise, Terraform will see a diff and force replacement.

`aviatrix_gateway_certificate_config` can be imported using controller IP with the dots(.) replaces with dashes(-), e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_gateway_certificate_config.test 10-11-12-13
```
