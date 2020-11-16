---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_remote_syslog"
description: |-
  Enables and disables remote syslog
---

# aviatrix_remote_syslog

The **aviatrix_remote_syslog** resource allows the enabling and disabling of remote syslog.

## Example Usage

```hcl
# Enable remote syslog without TLS
resource "aviatrix_remote_syslog" "test_remote_syslog" {
  index    = 0
  server   = "1.2.3.4"
  port     = 10
  protocol = "TCP"
}
```

```hcl
# Enable remote syslog with TLS
resource "aviatrix_remote_syslog" "test_remote_syslog" {
  index              = 0
  server             = "1.2.3.4"
  port               = 10
  protocol           = "TCP"
  ca_certificate     = "/path/to/ca.pem"
  public_certificate = "/path/to/server.pem"
  private_key        = "/path/to/client.pem"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `index` (Optional) Profile index. An index from 0 to 9 is supported. 0 by default.
* `server` (Required) Server IP.
* `port` (Required) Port number.
* `protocol` (Optional) TCP or UDP. TCP by default.

### Optional
* `ca_certificate` (Optional) Certificate Authority (CA) certificate. `ca_certificate`, `public_certificate` and `private_key` must be used together.
* `public_certificate` (Optional) Public certificate of the controller signed by the same CA.
* `private_key` (Optional) Private key of the controller that pairs with the public certificate.
* `template` (Optional) Optional custom template.
* `exclude_gateway_list` (Optional) List of gateways to be excluded from logging. e.g.: "gateway01", "gateway01, gateway01-hagw". Use a comma delimited string.

## Import

**remote_syslog** can be imported using "remote_syslog_" + `index`, e.g.

```
$ terraform import aviatrix_remote_syslog.test remote_syslog_0
```
