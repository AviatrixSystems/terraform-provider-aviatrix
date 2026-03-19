---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_mitm_ca"
description: |-
  Gets details about a specific DCF MITM CA.
---

# aviatrix_dcf_mitm_ca

The **aviatrix_dcf_mitm_ca** data source provides details about a specific Distributed Cloud Firewall (DCF) MITM (Man-in-the-Middle) Certificate Authority.

## Example Usage

```hcl
# Aviatrix DCF MITM CA Data Source
data "aviatrix_dcf_mitm_ca" "example" {
  name = "my-mitm-ca"
}
# To activate the CA, use the aviatrix_dcf_mitm_ca_selection resource
resource "aviatrix_dcf_mitm_ca_selection" "example" {
  ca_id = aviatrix_dcf_mitm_ca.example.ca_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) (String) Name of the MITM CA.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `ca_id` - (String) The unique identifier (UUID) for the MITM CA.
* `ca_hash` - (String) Hash of the certificate.
* `certificate_chain` - (String) The certificate chain in PEM format.
* `state` - (String) The state of the MITM CA (`active` or `inactive`).
* `created_at` - (String) Time when the CA was created in RFC3339 format.
* `origin` - The origin of the MITM CA, custom - Customer uploaded, aviatrix - system provided
