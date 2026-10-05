---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_telix_profile"
description: |-
  Creates and manages an Aviatrix Telix telemetry export profile
---

# aviatrix_telix_profile

The **aviatrix_telix_profile** resource creates and manages a Telix profile, which configures how the Aviatrix platform exports telemetry (DCF logs, gateway operational syslog, Prometheus node exporter metrics) to an external OTLP destination. Each profile pairs a set of telemetry sources with a destination and a gateway scope with optional filters.

## Example Usage

### Minimal profile exporting DCF logs to an OTLP gRPC collector

```hcl
resource "aviatrix_telix_profile" "dcf_logs" {
  display_name = "dcf-logs-prod"
  sources      = ["TELEMETRY_SOURCE_DCF_LOGS"]

  destination {
    otlp {
      endpoint = "otel-collector.example.com:4317"
      protocol = "TELIX_OTLP_PROTOCOL_GRPC"
    }
  }

  scope {
    all_gateways {}
  }
}
```

### Profile with mTLS, custom headers, selected gateways, and DCF filters

```hcl
resource "aviatrix_telix_profile" "secure_export" {
  display_name = "secure-export"
  sources      = ["TELEMETRY_SOURCE_DCF_LOGS"]
  enabled      = true

  destination {
    otlp {
      endpoint = "https://collector.example.com/v1/logs"
      protocol = "TELIX_OTLP_PROTOCOL_HTTP"

      headers = {
        "X-Tenant-Id" = "acme-prod"
        "X-Auth"      = var.collector_auth_token
      }

      tls {
        ca_certificate_pem     = file("${path.module}/ca.pem")
        client_certificate_pem = file("${path.module}/client.pem")
        client_private_key_pem = file("${path.module}/client.key")
        server_name_override   = "collector.example.com"
        insecure_skip_verify   = false
      }
    }
  }

  scope {
    selected_gateways {
      gateway_names = [
        "transit-aws-us-east-1",
        "transit-aws-us-west-2",
      ]
    }
  }

  filters {
    dcf_log_types = ["FLOW", "INTRUSION"]
    dcf_actions   = ["ALLOW", "DENY"]
  }
}
```

## Argument Reference

### Required

* `display_name` - (Required) Human-readable name for the profile. Must be unique across profiles in the controller.
* `sources` - (Required, ForceNew) Telemetry sources exported by this profile. Must contain at least one of:
  * `TELEMETRY_SOURCE_DCF_LOGS`
  * `TELEMETRY_SOURCE_NODE_EXPORTER`
  * `TELEMETRY_SOURCE_GATEWAY_OPERATIONAL_SYSLOG`

  Changing this list forces resource replacement.
* `destination` - (Required) Destination configuration block. Must contain exactly one nested `otlp` block.
* `scope` - (Required) Gateway scope block. Must contain exactly one of `all_gateways {}` or `selected_gateways { gateway_names = [...] }`.

### Optional

* `enabled` - (Optional) Whether this export profile is active. Defaults to `true`.
* `filters` - (Optional) Optional pre-export filtering for telemetry data. See [filters block](#filters) below.

### `destination.otlp` block

* `endpoint` - (Required) OTLP collector endpoint URL.
* `protocol` - (Required, ForceNew) OTLP transport protocol. One of `TELIX_OTLP_PROTOCOL_GRPC` or `TELIX_OTLP_PROTOCOL_HTTP`. Changing this value forces resource replacement.
* `headers` - (Optional, Sensitive) Map of static headers sent with each OTLP request. Typically used for authentication tokens. Write-only; see [Handling Sensitive Fields](#handling-sensitive-fields) below.
* `tls` - (Optional) TLS configuration block; see [tls block](#tls) below. Driven by configuration—omitting the block after configuring TLS triggers a planned change so Terraform can clear OTLP TLS on the controller (`has_ca_certificate`, etc. nested under each block remain **computed-only** outputs from read).

### `destination.otlp.tls` block

* `ca_certificate_pem` - (Optional, Sensitive) PEM-encoded CA certificate used to validate the remote endpoint. Write-only; see [Handling Sensitive Fields](#handling-sensitive-fields) below.
* `client_certificate_pem` - (Optional, Sensitive) PEM-encoded client certificate for mutual TLS. Write-only; see [Handling Sensitive Fields](#handling-sensitive-fields) below.
* `client_private_key_pem` - (Optional, Sensitive) PEM-encoded client private key for mutual TLS. Write-only; see [Handling Sensitive Fields](#handling-sensitive-fields) below.
* `server_name_override` - (Optional) Server name override used during TLS certificate validation.
* `insecure_skip_verify` - (Optional) Disables TLS certificate verification. Intended for testing only.

### `scope` block

Exactly one of the following must be specified:

* `all_gateways {}` - Empty marker block. Applies the profile to every gateway compatible with the selected sources.
* `selected_gateways { gateway_names = [...] }` - Restricts the profile to a specific list of gateways.

### `filters` block

* `dcf_log_types` - (Optional) DCF log types to include. Requires `TELEMETRY_SOURCE_DCF_LOGS` in `sources`.
* `dcf_actions` - (Optional) DCF actions to include. Requires `TELEMETRY_SOURCE_DCF_LOGS` in `sources`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `profile_id` - Server-generated unique identifier for the profile.
* `created_at` - RFC3339 timestamp of when the profile was created.
* `last_modified_at` - RFC3339 timestamp of when the profile was last modified.
* `destination.0.otlp.0.has_headers` - Whether the controller currently holds a non-empty `headers` value for this profile.
* `destination.0.otlp.0.tls.0.has_ca_certificate` - Whether the controller currently holds a CA certificate.
* `destination.0.otlp.0.tls.0.has_client_certificate` - Whether the controller currently holds a client certificate.
* `destination.0.otlp.0.tls.0.has_client_private_key` - Whether the controller currently holds a client private key.

## Handling Sensitive Fields

The `headers` map and the three TLS PEM fields (`ca_certificate_pem`, `client_certificate_pem`, `client_private_key_pem`) are write-only on the controller. They are accepted by the API on create and update, but the controller never returns them on read for security reasons. Instead, the GET response returns boolean presence flags (`has_headers`, `has_ca_certificate`, `has_client_certificate`, `has_client_private_key`) which this resource exposes as computed attributes.

This has two consequences for drift detection:

~> **Note:** External deletion of a sensitive value is detected. If one of these fields is cleared outside Terraform (for example through the controller UI or a direct API call), the corresponding presence flag is set to `false` on the next read, the stale value is removed from state, and the next `terraform apply` re-sends the value from the configuration.

~> **Note:** External replacement of a sensitive value with a different value is not detected. Terraform compares configuration to state; if the value is changed only on the controller, configuration and state are unchanged from Terraform’s perspective, and the controller retains the out-of-band value until the configuration is updated, or the resource is replaced (`terraform apply -replace`) or tainted. To keep Terraform authoritative for these fields, perform certificate and secret rotations by changing configuration and applying.

### Clearing headers and TLS (best practice)

* **Removing `headers`:** Omit the entire `headers` argument from configuration when you no longer want OTLP HTTP/gRPC headers on the controller. The provider clears stored headers on `apply` instead of silently leaving prior values behind.
* **Removing TLS customization:** Omit the entire `destination.otlp.tls` block—do **not** keep an empty `tls { }` as your long-term steady state. An optional block shell with every attribute unset is easy to introduce while editing or templating modules, but it is not idiomatic Terraform and can behave differently than removing TLS entirely depending on Terraform’s diff semantics. Prefer **no `tls` block** when OTLP TLS options should mirror “no TLS block” behavior.
* **`server_name_override` and `insecure_skip_verify`** are persisted with TLS material when you configure `tls`; removing the `tls` block clears them along with the PEM fields as part of updates.

Operators calling **PATCH `/telix/profile/{id}`** directly should remember that omission of nested fields preserves existing controller state; Terraform users normally do **not** need to model API merge semantics—the provider emits explicit clears when you drop `headers` or the entire `tls` block.

## Refresh-only plans and `has_*` attributes

The `has_*` presence attributes (`has_headers`, `has_ca_certificate`, `has_client_certificate`, `has_client_private_key`) are returned by the controller on **read** and are **computed-only** (they are not set in `.tf` configuration).

If the profile is modified outside Terraform—for example via the controller UI or API when adding or removing TLS material—the next refresh still reads the updated `has_*` values from the API.

A **default** `terraform plan` reports actions that **`terraform apply` would perform** on infrastructure; it is not a full diff of stored state versus the latest read. **Result:** If configuration is unchanged, Terraform may plan **no in-place update** for the resource even though computed attributes (including `has_*`) would differ in state after refresh. The plan may then show **no changes**, or updates to `has_*` may be difficult to see. That behavior comes from Terraform’s plan output, not from a skipped Read: the provider still refreshes and receives the correct flags.

To **compare stored state to post-refresh values** (including `has_*`) **without** applying configuration-driven infrastructure changes:

```shell
terraform plan -refresh-only
```

To **write** refreshed state to the configured backend **without** changing the controller:

```shell
terraform apply -refresh-only
```

If a profile was changed outside Terraform and visibility into presence flags or other computed-only drift is required, use one of the commands above.

The `headers` map and TLS PEM arguments are stored in Terraform state in plaintext. The `Sensitive: true` flag redacts them in CLI output (`plan`, `apply`, `show`, `state show`, `output`), but the raw state file and any `-json` output contain the actual values. Secure your state backend (encryption at rest, restricted access) accordingly. For background, see [Sensitive state best practices](https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state).

## Notes

* `sources` and `destination.otlp.protocol` are immutable after creation. Changing either forces destruction and recreation of the profile, which interrupts telemetry export and assigns a new `profile_id`.
* Most fields can be updated in place via PATCH without disrupting export, including OTLP endpoints, scopes, filters, TLS, and headers (`headers` clears when removed; omit the entire `tls` block to clear OTLP TLS—see [Clearing headers and TLS (best practice)](#clearing-headers-and-tls-best-practice)).
* `filters.dcf_log_types` and `filters.dcf_actions` require `TELEMETRY_SOURCE_DCF_LOGS` to be present in `sources`. The provider validates this at plan time.
* The maximum number of Telix profiles per controller is bounded.
