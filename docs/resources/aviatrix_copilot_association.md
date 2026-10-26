---
subcategory: "Settings"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_copilot_association"
description: |-
  Creates and manages a CoPilot Association
---

# aviatrix_copilot_association

The **aviatrix_copilot_association** resource allows management of controller CoPilot Association. This resource is available as of provider version R2.19+.

## Example Usage

```hcl
# Create a CoPilot Association
resource "aviatrix_copilot_association" "test_copilot_association" {
    copilot_address = "copilot.aviatrix.com"
}
```

```hcl
# Create a CoPilot Association with an optional FQDN
resource "aviatrix_copilot_association" "test_copilot_association_with_fqdn" {
    copilot_address = "10.11.12.13"
    public_ip       = "35.184.203.217"
    copilot_fqdn    = "copilot.aviatrix.com"
}
```


## Argument Reference

The following arguments are supported:

* `copilot_address` - (Required) CoPilot instance IP Address or Hostname.
* `copilot_fqdn` - (Optional) CoPilot FQDN association metadata. When provided, it must be ASCII and a valid IP address or DNS hostname. Updating this value updates the association in place without recreating the resource.
* `public_ip` - (Optional) CoPilot public IP address or hostname. Defaults to `copilot_address` when omitted.

## Import

**aviatrix_copilot_association** can be imported using controller IP, e.g. controller IP is : 10.11.12.13

```
$ terraform import aviatrix_copilot_association.test_copilot_association 10-11-12-13
```
