---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_http_filter"
description: |-
  Creates and manages an Aviatrix Distributed-firewalling HTTP Filter
---

# aviatrix_dcf_http_filter

The **aviatrix_dcf_http_filter** resource handles the creation and management of Aviatrix Distributed-firewalling HTTP Filters. An HTTP filter matches HTTP requests by header and can be attached to DCF rules through the `http_filters` argument of [aviatrix_dcf_ruleset](aviatrix_dcf_ruleset.md).

A filter is expressed in disjunctive normal form: it matches a request when **any** `filters` block matches, and a `filters` block matches when **all** of its `all` header matches hold.

## Example Usage

```hcl
# Match POST or PUT requests that carry "app: web", or any request that carries an "x-debug" header
resource "aviatrix_dcf_http_filter" "example" {
  name        = "web-writes-or-debug"
  description = "Writes to the web app, or debug traffic"

  filters {
    all {
      header = ":method"
      values = ["POST", "PUT"]
    }
    all {
      header = "app"
      values = ["web"]
    }
  }

  filters {
    all {
      header = "x-debug"
    }
  }
}

# Reference the filter from a DCF rule
resource "aviatrix_dcf_ruleset" "example" {
  name = "example-ruleset"

  rules {
    name             = "web-writes"
    action           = "PERMIT"
    protocol         = "TCP"
    src_smart_groups = [aviatrix_smart_group.src.uuid]
    dst_smart_groups = [aviatrix_smart_group.dst.uuid]
    http_filters     = [aviatrix_dcf_http_filter.example.uuid]
  }
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) Name of the HTTP filter. Must be unique across HTTP filters.
* `filters` - (Required) One or more clauses. The filter matches when any clause matches.
  * `all` - (Required) One or more header matches. The clause matches when all of them hold.
    * `header` - (Required) HTTP header name to match. Must be a lowercase HTTP header name (for example `user-agent`), or one of the HTTP/2 pseudo-headers `:method`, `:scheme`, `:authority`, or `:path`.
    * `values` - (Optional) List of header values to match.

### Optional
* `description` - (Optional) Description of the HTTP filter.

### Computed
* `uuid` - (Computed) The unique identifier for the HTTP filter.

## Import

**aviatrix_dcf_http_filter** can be imported using the HTTP filter UUID, e.g.

```
$ terraform import aviatrix_dcf_http_filter.test <http_filter_uuid>
```

## Notes

* Distributed Cloud Firewall must be enabled on the controller to manage HTTP filters.
* An HTTP filter that is referenced by a DCF rule cannot be deleted. Remove it from all rules first. When the rule refers to the filter's `uuid` attribute in the same configuration, Terraform orders the rule update before the filter delete.
* The order of `filters` blocks, `all` blocks, and `values` is preserved. Reordering them in configuration produces an update even though matching behavior is unchanged.
