---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_mitm_ca_selection"
description: |-
  Selects the active DCF MITM CA for the Aviatrix Controller
---

# aviatrix_dcf_mitm_ca_selection

The **aviatrix_dcf_mitm_ca_selection** resource manages the selection of the active MITM (Man-in-the-Middle) Certificate Authority for Distributed Cloud Firewall. Only one CA can be active at a time, so selecting a new CA will automatically deactivate the previously active CA.

## Example Usage

### Select a MITM CA as Active

```hcl
# Create a DCF MITM CA
resource "aviatrix_dcf_mitm_ca" "example" {
  name              = "corporate-mitm-ca"
  key               = file("${path.module}/mitm-ca-key.pem")
  certificate_chain = file("${path.module}/mitm-ca-cert.pem")
}

# Select the CA as the active MITM CA
resource "aviatrix_dcf_mitm_ca_selection" "example" {
  mitm_ca_id = aviatrix_dcf_mitm_ca.example.ca_id
}
```

### Switch Between CAs

```hcl
resource "aviatrix_dcf_mitm_ca" "primary" {
  name              = "primary-mitm-ca"
  key               = file("${path.module}/primary-ca-key.pem")
  certificate_chain = file("${path.module}/primary-ca-cert.pem")
}

resource "aviatrix_dcf_mitm_ca" "backup" {
  name              = "backup-mitm-ca"
  key               = file("${path.module}/backup-ca-key.pem")
  certificate_chain = file("${path.module}/backup-ca-cert.pem")
}

# Select primary CA as active (change mitm_ca_id to switch)
resource "aviatrix_dcf_mitm_ca_selection" "active" {
  mitm_ca_id = aviatrix_dcf_mitm_ca.primary.ca_id
}
```

## Argument Reference

The following arguments are supported:

### Required

* `mitm_ca_id` - (Required) The UUID of the DCF MITM CA to select as the active CA. This can be obtained from the `ca_id` attribute of an `aviatrix_dcf_mitm_ca` resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for this resource, derived from the controller IP address.

## Import

**aviatrix_dcf_mitm_ca_selection** can be imported using the resource ID (controller IP with dots replaced by dashes), e.g.

```
$ terraform import aviatrix_dcf_mitm_ca_selection.example 10-0-0-1
```

## Notes

* This is a singleton resource - only one `aviatrix_dcf_mitm_ca_selection` resource should exist per Aviatrix Controller.
* When this resource is created or updated, the specified CA becomes `active` and the previously active CA becomes `inactive`.
* When this resource is destroyed, the system automatically falls back to the built-in system CA.
* The resource ID is derived from the controller IP address to ensure uniqueness per controller.

