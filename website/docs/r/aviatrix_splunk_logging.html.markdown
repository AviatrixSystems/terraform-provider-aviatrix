---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_splunk_logging"
description: |-
  Enables and disables splunk logging
---

# aviatrix_splunk_logging

The **aviatrix_splunk_logging** resource allows the enabling and disabling of splunk logging.

## Example Usage

```hcl
# Enable splunk logging using server and port combination
resource "aviatrix_splunk_logging" "test_splunk_logging" {
  server = "1.2.3.4"
  port   = 10
}
```

```hcl
# Enable splunk logging using configuration file
resource "aviatrix_splunk_logging" "test_splunk_logging" {
  custom_output_configuration = file("/path/to/configuration.file")
}
```

## Argument Reference

The following arguments are supported:

### Required
* `server` (Optional) Server IP. Either `server` and `port` combination or `cu_output_configuration` is required.
* `port` (Optional) Port number.
* `custom_output_configuration` (Optional) Configuration file. Use the file function to read from a file.

### Optional
* `custom_input_configuration` (Optional) Custom configuration.
* `excluded_gateways` (Optional) List of gateways to be excluded from logging. e.g.: ["gateway01", "gateway02", "gateway01-hagw"].

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - The status of splunk logging.

## Import

**splunk_logging** can be imported using `splunk_logging`, e.g.

```
$ terraform import aviatrix_splunk_logging.test splunk_logging
```
