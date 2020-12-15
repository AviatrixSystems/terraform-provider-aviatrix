---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_filebeat_forwarder"
description: |-
Enables and disables filebeat forwarder
---

# aviatrix_filebeat_forwarder

The **aviatrix_filebeat_forwarder** resource allows the enabling and disabling of filebeat forwarder.

## Example Usage

```hcl
# Enable filebeat forwarder
resource "aviatrix_filebeat_forwarder" "test_filebeat_forwarder" {
  server   = "1.2.3.4"
  port     = 10
  excluded_gateways   = ["a", "b"]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `server` (Required) Server IP.
* `port` (Required) Port number.

### Optional
* `trusted_ca_file_path` (Optional) The file path of the trusted CA. 
* `config_file_path` (Optional) The config file path.
* `excluded_gateways` (Optional) List of gateways to be excluded from logging. e.g.: ["gateway01", "gateway02", "gateway01-hagw"].

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - The status of filebeat forwarder.

## Import

**filebeat_forwarder** can be imported using "filebeat_forwarder", e.g.

```
$ terraform import aviatrix_filebeat_forwarder.test filebeat_forwarder
```
