---
subcategory: "OpenVPN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_vpn_cert_download"
description: |-
  Creates and Manages Aviatrix VPN Users
---

# aviatrix_vpn_cert_download

The **aviatrix_vpn_cert_download** resource manages the VPN Certificate Download configuration for SAML Authentication

## Example Usage

```hcl
# Set up the Aviatrix VPN Certificate Download configuration
resource "aviatrix_vpn_cert_download" "test_vpn_cert_download" {
  download_enabled = true
  saml_endpoints = ["saml_endpoint_name"]
}
```
## Argument Reference

The following arguments are supported:

### Optional

* `download_enabled` - (Optional) Whether the VPN Certificate download is enabled `gw_name`. Supported Values: "true", "false".
* `saml_endpoints` - (Optional) List of SAML endpoint names for which the downloading should be enabled . Currently, only a single endpoint is supported. Example: ["saml_endpoint_1"].

## Import

**vpn_cert_download** can be imported using the default id `vpn_cert_download`, e.g.

```
$ terraform import aviatrix_vpn_cert_download.test_vpn_cert_download vpn_cert_download
```
