---
subcategory: "Distributed Cloud Firewall"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_trustbundle"
description: |-
  Creates and manages an Aviatrix DCF Trust Bundle
---

# aviatrix_dcf_trustbundle

The **aviatrix_dcf_trustbundle** resource handles the creation and management of DCF Trust Bundles for verifying origin certificates in Distributed Cloud Firewall Man-in-the-Middle (MITM) inspection.

## Example Usage

### Basic Trust Bundle

```hcl
# Create a DCF Trust Bundle
resource "aviatrix_dcf_trustbundle" "example" {
  display_name   = "corporate-root-ca"
  bundle_content = file("${path.module}/corporate-root-ca.pem")
}
```

### Trust Bundle with inline certificate

```hcl
resource "aviatrix_dcf_trustbundle" "example" {
  display_name = "example-trustbundle"
  bundle_content = file("bundle_file")
}
```

## Argument Reference

The following arguments are supported:

### Required

* `display_name` - (Required) The display name for the DCF trust bundle. This name is used to identify the trust bundle in the Aviatrix Controller.
* `bundle_content` - (Required, Sensitive) The CA bundle content in PEM format. This should contain one or more X.509 certificates separated by new lines, that will be used to verify origin certificates during DCF MITM inspection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `bundle_id` - The unique identifier (UUID) assigned to the trust bundle by the Aviatrix Controller.
* `created_at` - ISO 8601 timestamp indicating when the trust bundle was created.

## Import

**aviatrix_dcf_trustbundle** can be imported using the `bundle_id` (UUID), e.g.

```
$ terraform import aviatrix_dcf_trustbundle.example 41984f8b-5a37-4272-89b3-57c79e9ff77c
```

## Notes

* The `bundle_content` must be valid X.509 certificates in PEM format.
* Multiple certificates can be included in a single bundle by concatenating them.
* The trust bundle is used in DCF MITM scenarios to verify the authenticity of origin certificates.
* Changes to `bundle_content` or `display_name` will result in an update operation.
* The certificate content should include proper PEM headers (`-----BEGIN CERTIFICATE-----`) and footers (`-----END CERTIFICATE-----`).
* Invalid certificate formats will result in an error during creation or update.
