---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_datadog_agent"
description: |-
  Enables and disables datadog agent
---

# aviatrix_datadog_agent

The **aviatrix_datadog_agent** resource allows the enabling and disabling of datadog agent.

## Example Usage

```hcl
# Enable datadog agent
resource "aviatrix_datadog_agent" "test_datadog_agent" {
  api_key           = "your_api_key"
  site              = "datadoghq.com"
  excluded_gateways = ["a", "b"]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `api_key` (Required) API key.
* `site` (Optional) Site preference ("datadoghq.com" or" datadoghq.eu"). "datadoghq.com" by default.

### Optional
* `excluded_gateways` (Optional) List of gateways to be excluded from logging. e.g.: ["gateway01", "gateway02", "gateway01-hagw"].

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - The status of datadog agent.

## Import

**datadog_agent** can be imported using "datadog_agent", e.g.

```
$ terraform import aviatrix_datadog_agent.test datadog_agent
```
