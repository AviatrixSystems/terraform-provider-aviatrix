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
  name     = "test"
  server   = "1.2.3.4"
  port     = 10
  protocol = "TCP"
}
```

```hcl
# Enable remote syslog with TLS
resource "aviatrix_remote_syslog" "test_remote_syslog" {
  index                   = 0
  name                    = "rsyslog-profilename"
  server                  = "1.2.3.4"
  port                    = 10
  protocol                = "TCP"
  ca_certificate_file     = file("/path/to/ca.pem")
  public_certificate_file = file("/path/to/server.pem")
  private_key_file        = file("/path/to/client.pem")
}
```

## Argument Reference

The following arguments are supported:

### Required
* `index` - (Optional) Profile index. An index from 0 to 9 is supported. 0 by default.
* `server` - (Required) Server IP.
* `port` - (Required) Port number.
* `protocol` - (Optional) TCP or UDP. TCP by default.

### Optional
* `name` - (Optional) Profile name.
* `ca_certificate_file` - (Optional) The Certificate Authority (CA) certificate. Use the `file` function to read from a file.
* `public_certificate_file` - (Optional) The public certificate of the controller signed by the same CA. Use the `file` function to read from a file.
* `private_key_file` - (Optional) The private key of the controller that pairs with the public certificate. Use the `file` function to read from a file.

~> **NOTE:** To enable TLS, either `ca_certificate_file`, or the combination of `ca_certificate_file`, `public_certificate_file` and `private_key_file` should be used.

* `template` - (Optional) Optional custom template.
* `excluded_gateways` - (Optional) List of gateways to be excluded from logging. e.g.: ["gateway01", "gateway02", "gateway01-hagw"].

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - The status of remote syslog.
* `notls` - This attribute is true if the remote syslog is not protected by TLS.

## Import

**remote_syslog** can be imported using "remote_syslog_" + `index`, e.g.

```
$ terraform import aviatrix_remote_syslog.test remote_syslog_0
```
