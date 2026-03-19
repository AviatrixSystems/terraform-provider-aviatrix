---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_mitm_ca"
description: |-
  Creates and manages an Aviatrix DCF MITM CA
---

# aviatrix_dcf_mitm_ca

The **aviatrix_dcf_mitm_ca** resource handles the creation and management of MITM (Man-in-the-Middle) Certificate Authorities for Distributed Cloud Firewall. These CAs are used to sign certificates during MITM inspection.

## Example Usage

### Basic MITM CA

```hcl
# Create a DCF MITM CA
resource "aviatrix_dcf_mitm_ca" "example" {
  name              = "corporate-mitm-ca"
  key               = file("${path.module}/mitm-ca-key.pem")
  certificate_chain = file("${path.module}/mitm-ca-cert.pem")
}

# To activate the CA, use the aviatrix_dcf_mitm_ca_selection resource
resource "aviatrix_dcf_mitm_ca_selection" "example" {
  ca_id = aviatrix_dcf_mitm_ca.example.ca_id
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The name for the MITM CA. Every CA must have a unique name.
* `key` - (Required, Sensitive, ForceNew) The private key in PEM format. Changing this forces a new resource to be created.
* `certificate_chain` - (Required, ForceNew) The certificate chain in PEM format. The first certificate must be a signing CA certificate and should match the provided private key. Changing this forces a new resource to be created.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `ca_id` - The unique identifier (UUID) assigned to the MITM CA by the Aviatrix Controller.
* `ca_hash` - Hash of the certificate.
* `created_at` - Timestamp in RFC3339 format indicating when the MITM CA was created.
* `state` - The state of the MITM CA (`active` or `inactive`). To activate a CA, use the `aviatrix_dcf_mitm_ca_selection` resource.
* `origin` - The origin of the MITM CA, custom - Customer uploaded, aviatrix - system provided

## Import

**aviatrix_dcf_mitm_ca** can be imported using the `ca_id` (UUID), e.g.

```
$ terraform import aviatrix_dcf_mitm_ca.example 41984f8b-5a37-4272-89b3-57c79e9ff77c
```

## Notes

* The `key` and `certificate_chain` must be valid PEM formatted content.
* The certificate must be a CA certificate (signing certificate) that can be used to sign other certificates.
* Only one MITM CA can be `active` at a time. To activate a CA, use the `aviatrix_dcf_mitm_ca_selection` resource.
* The `key` and `certificate_chain` cannot be updated after creation. To use a different certificate, you must create a new resource.
* The private key is marked as sensitive and will not be displayed in Terraform output or logs.
* Changes to `name` will result in an update operation (PATCH).
