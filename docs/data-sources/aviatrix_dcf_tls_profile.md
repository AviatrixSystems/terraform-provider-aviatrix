---
subcategory: "Distributed Cloud Firewall"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_tls_profile"
description: |-
  Gets details about a specific DCF TLS profile.
---

# aviatrix_dcf_tls_profile

The **aviatrix_dcf_tls_profile** data source provides details about a specific Distributed Cloud Firewall (DCF) TLS profile.

## Example Usage

```hcl
# Aviatrix DCF TLS Profile Data Source
data "aviatrix_dcf_tls_profile" "example" {
  display_name = "my-tls-profile"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) (String) Display name of the TLS profile.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - (String) The unique identifier for the TLS profile.
* `certificate_validation` - (String) Certificate validation mode.
* `verify_sni` - (Boolean) Toggle to enable advanced SNI verification of client connections.
* `ca_bundle_id` - (String) UUID of the CA bundle used for origin certificate validation.
