---
subcategory: "Secured Networking"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_dcf_http_filter"
description: |-
  Gets details about a specific DCF HTTP filter.
---

# aviatrix_dcf_http_filter

The **aviatrix_dcf_http_filter** data source provides details about a specific Distributed Cloud Firewall (DCF) HTTP filter, looked up by name. Use it to reference an HTTP filter that is not managed in the same Terraform configuration.

## Example Usage

```hcl
# Aviatrix DCF HTTP Filter Data Source
data "aviatrix_dcf_http_filter" "example" {
  name = "my-http-filter"
}

resource "aviatrix_dcf_ruleset" "example" {
  name = "example-ruleset"

  rules {
    name             = "filtered-rule"
    action           = "PERMIT"
    protocol         = "TCP"
    src_smart_groups = [aviatrix_smart_group.src.uuid]
    dst_smart_groups = [aviatrix_smart_group.dst.uuid]
    http_filters     = [data.aviatrix_dcf_http_filter.example.uuid]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) (String) Name of the HTTP filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - (String) The unique identifier for the HTTP filter.
* `description` - (String) Description of the HTTP filter.
* `filters` - (List) Clauses of the HTTP filter. The filter matches when any clause matches.
  * `all` - (List) Header matches that must all hold for the clause to match.
    * `header` - (String) HTTP header name to match.
    * `values` - (List of String) Header values to match.
