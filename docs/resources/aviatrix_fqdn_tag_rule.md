---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_fqdn_tag_rule"
description: |-
  Manages Aviatrix FQDN filtering domain name rule
---

# aviatrix_fqdn_tag_rule

The **aviatrix_fqdn_tag_rule** resource manages a single FQDN filtering domain name rule.

~> **NOTE on FQDN and FQDN Tag Rule resources:** Terraform currently provides both a standalone FQDN Tag Rule resource and an FQDN resource with domain name rules defined in-line. At this time, you cannot use an FQDN resource with in-line rules in conjunction with any FQDN Tag Rule resources. Doing so will cause a conflict of rule settings and will overwrite rules. In order to use this resource, please set `manage_domain_names` in the **aviatrix_fqdn** resource to false.

## Example Usage

```hcl
# Create an Aviatrix Gateway FQDN Tag Rule filter rule
resource "aviatrix_fqdn_tag_rule" "test_fqdn" {
  fqdn_tag_name = "my_tag"
  fqdn          = "reddit.com"
  protocol      = "tcp"
  port          = "443"
}
```

## Argument Reference

The following arguments are supported:

* `fqdn_tag_name` - (Required) FQDN Filter tag name.
* `fqdn` - (Required) FQDN. Example: "facebook.com".
* `protocol` - (Required) Protocol. Valid values: "all", "tcp", "udp", "icmp".
* `port` - (Required) Port. Example "25".
* `action` - (Optional) What action should happen to matching requests. Possible values are: 'Base Policy', 'Allow' or 'Deny'. Defaults to 'Base Policy' if no value provided.
    * For protocol "all", port must be set to "all".
    * For protocol “icmp”, port must be set to “ping”.

## Import

**fqdn_tag_rule** can be imported using the `fqdn_tag_name`, `fqdn`, `protocol`, `port` and `action` separated by `~`, e.g.

```
$ terraform import aviatrix_fqdn_tag_rule.test "fqdn_tag_name~fqdn~protocol~port~action"
```
