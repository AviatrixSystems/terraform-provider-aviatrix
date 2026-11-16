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

# Create an Aviatrix Distributed Firewalling TLS Profile that logs selected TLS client fields
resource "aviatrix_dcf_tls_profile" "log_fields" {
  display_name             = "log-fields-tls-profile"
  certificate_validation   = "CERTIFICATE_VALIDATION_NONE"
  verify_sni               = true
  tls_log_fields_extension = ["TLS_CLIENT_SNI", "TLS_CLIENT_JA4_FINGERPRINT"]
}

# Create an Aviatrix Distributed Firewalling TLS Profile that logs every TLS client field
resource "aviatrix_dcf_tls_profile" "log_all_fields" {
  display_name             = "log-all-fields-tls-profile"
  certificate_validation   = "CERTIFICATE_VALIDATION_NONE"
  verify_sni               = true
  tls_log_fields_extension = ["TLS_LOG_FIELDS_ALL"]
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
* `tls_log_fields_extension` - (Optional) Set of TLS log fields to include in the logs. Each value must be one of TLS_CLIENT_VERSIONS, TLS_CLIENT_CIPHERS, TLS_CLIENT_KEY_SHARE_GROUPS, TLS_CLIENT_SIGNATURE_ALGORITHMS, TLS_CLIENT_JA4_FINGERPRINT, TLS_CLIENT_SNI, TLS_CLIENT_SUPPORTED_GROUPS, TLS_CLIENT_ALPNS, TLS_CLIENT_ECH_PRESENT, or TLS_LOG_FIELDS_ALL. TLS_LOG_FIELDS_ALL selects every field and must be the only value when used. If omitted, no fields are included and any fields configured outside Terraform are removed on the next apply.

### Computed
* `uuid` - (Computed) The unique identifier for the TLS profile.

## Certificate Validation Modes

The `certificate_validation` parameter supports the following modes:

* `CERTIFICATE_VALIDATION_NONE` - No certificate validation is performed
* `CERTIFICATE_VALIDATION_LOG_ONLY` - Certificate validation is performed but only logged, traffic is not blocked
* `CERTIFICATE_VALIDATION_ENFORCE` - Certificate validation is enforced and connections to origins with invalid certificates will
be blocked.

-> **NOTE** If `certificate_validation` mode is anything other than `CERTIFICATE_VALIDATION_NONE` the `ca_bundle_id` must be specified.

## TLS Log Fields

`TLS_LOG_FIELDS_ALL` is a provider shorthand. The provider sends every TLS log field it knows about to the controller, and the controller records that as "all fields". State keeps `TLS_LOG_FIELDS_ALL` while the controller continues to report every field.

-> **NOTE** An imported TLS profile that logs every field is stored in state as the individual field names. If your configuration uses `TLS_LOG_FIELDS_ALL`, the first plan after import shows an in-place update that does not change the controller.

-> **NOTE** `tls_log_fields_extension` requires a controller version that supports TLS log fields. On older controllers, leave this argument unset.

## Import

**aviatrix_dcf_tls_profile** can be imported using the TLS profile UUID, e.g.

```
$ terraform import aviatrix_dcf_tls_profile.test <tls_profile_uuid>
```

## Notes

* TLS profiles are used in DCF policies to provide granular TLS validation capability.
