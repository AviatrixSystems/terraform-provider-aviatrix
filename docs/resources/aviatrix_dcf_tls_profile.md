---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_tls_profile"
description: |-
  Creates and manages an Aviatrix Distributed-firewalling TLS Profile
---

# aviatrix_dcf_tls_profile

The **aviatrix_dcf_tls_profile** resource handles the creation and management of Aviatrix Distributed-firewalling TLS Profiles.

## Example Usage

```hcl
# Use a data source to get the bundle UUID by its name
data "aviatrix_dcf_trustbundle" "bundle_sample" {
  display_name = "sample-bundle-1"
}

# Create an Aviatrix Distributed Firewalling TLS Profile with log only cerificate validation
resource "aviatrix_dcf_tls_profile" "basic" {
  display_name           = "basic-tls-profile"
  certificate_validation = "CERTIFICATE_VALIDATION_LOG_ONLY"
  verify_sni            = false
  ca_bundle_id          = data.aviatrix_dcf_trustbundle.bundle_sample.bundle_id
}

# Create an Aviatrix Distributed Firewalling TLS Profile with certificate enforcement
resource "aviatrix_dcf_tls_profile" "enforce" {
  display_name           = "enforce-tls-profile"
  certificate_validation = "CERTIFICATE_VALIDATION_ENFORCE"
  verify_sni            = true
  ca_bundle_id          = data.aviatrix_dcf_trustbundle.bundle_sample.bundle_id
}

# Create an Aviatrix Distributed Firewalling TLS Profile with no certificate validation
resource "aviatrix_dcf_tls_profile" "no_validation" {
  display_name           = "no-validation-tls-profile"
  certificate_validation = "CERTIFICATE_VALIDATION_NONE"
  verify_sni            = true
}
```

## Argument Reference

The following arguments are supported:

### Required
* `display_name` - (Required) Display name for the TLS profile.
* `certificate_validation` - (Required) Certificate validation mode for origin certificate. Must be one of CERTIFICATE_VALIDATION_NONE, CERTIFICATE_VALIDATION_LOG_ONLY, or CERTIFICATE_VALIDATION_ENFORCE.
* `verify_sni` - (Required) Toggle to enable advanced SNI verification of client connections. Type: Boolean.

### Optional
* `ca_bundle_id` - (Optional) UUID of the CA bundle that should be used for origin certificate validation. If not populated, the default bundle would be used. The aviatrix_dcf_trustbundle data source can be used to get the UUID from the bundle name.

### Computed
* `uuid` - (Computed) The unique identifier for the TLS profile.

## Certificate Validation Modes

The `certificate_validation` parameter supports the following modes:

* `CERTIFICATE_VALIDATION_NONE` - No certificate validation is performed
* `CERTIFICATE_VALIDATION_LOG_ONLY` - Certificate validation is performed but only logged, traffic is not blocked
* `CERTIFICATE_VALIDATION_ENFORCE` - Certificate validation is enforced and connections to origins with invalid certificates will
be blocked.

-> **NOTE** If `certificate_validation` mode is anything other than `CERTIFICATE_VALIDATION_NONE` the `ca_bundle_id` must be specified.

## Import

**aviatrix_dcf_tls_profile** can be imported using the TLS profile UUID, e.g.

```
$ terraform import aviatrix_dcf_tls_profile.test <tls_profile_uuid>
```

## Notes

* TLS profiles are used in DCF policies to provide granular TLS validation capability.
