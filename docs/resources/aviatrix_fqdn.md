---
subcategory: "Security"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_fqdn"
description: |-
  Manages Aviatrix FQDN filtering for gateways
---

# aviatrix_fqdn

The **aviatrix_fqdn** resource manages [FQDN filtering](https://docs.aviatrix.com/HowTos/fqdn_faq.html) for Aviatrix gateways.

~> **NOTE on FQDN and FQDN Tag Rule resources:** Terraform currently provides both a standalone FQDN Tag Rule resource and an FQDN resource with domain name rules defined in-line. At this time, you cannot use an FQDN resource with in-line rules in conjunction with any FQDN Tag Rule resources. Doing so will cause a conflict of rule settings and will overwrite rules. In order to use the **aviatrix_fqdn_tag_rule** resource, `manage_domain_names` must be set to false in this resource.

~> **NOTE:** Please see the [Notes](#notes) section below for troubleshooting known issues/deltas that may occur when enabling this feature

~> **NOTE:** Please note that there is no need to attach FQDN resource/enable this feature for the HA gateway. Enabling FQDN for the primary gateway will automatically handle this for the HA

## Example Usage

```hcl
# Create an Aviatrix Gateway FQDN filter
resource "aviatrix_fqdn" "test_fqdn" {
  fqdn_tag     = "my_tag"
  fqdn_enabled = true
  fqdn_mode    = "white"

  gw_filter_tag_list {
    gw_name        = "test-gw1"
    source_ip_list = [
      "172.31.0.0/16",
      "172.31.0.0/20"
    ]
  }

  gw_filter_tag_list {
    gw_name        = "test-gw2"
    source_ip_list = [
      "30.0.0.0/16"
    ]
  }
}
```

## Argument Reference

!> **WARNING:** Attribute `domain_names` has been deprecated as of provider version R2.18.1+ and will not receive further updates. Please set `manage_domain_names` to false, and use the standalone `aviatrix_fqdn_tag_rule` resource instead.

The following arguments are supported:

* `fqdn_tag` - (Required) FQDN Filter tag name.
* `fqdn_enabled` - (Optional) FQDN Filter tag status. Valid values: true, false.
* `fqdn_mode` - (Optional) Specify FQDN mode: whitelist or blacklist. Valid values: "white", "black".
* `manage_domain_names` - (Optional) Enable to manage domain name rules in-line. If false, domain name rules must be managed using `aviatrix_fqdn_tag_rule` resources. Default: true. Valid values: true, false. Available in provider version R2.17+.
* `gw_filter_tag_list` - (Optional) A list of gateways to attach to the specific tag.
  * `gw_name` - (Required) Name of the gateway to attach to the specific tag.
  * `source_ip_list` - (Optional) List of source IPs in the VPC qualified for a specific tag.
* `domain_names` - (Optional) One or more domain names in a list with details as listed below:
  * `fqdn` - (Required) FQDN. Example: "facebook.com".
  * `proto` - (Required) Protocol. Valid values: "all", "tcp", "udp", "icmp".
  * `port` - (Required) Port. Example "25".
  * `action` - (Optional) What action should happen to matching requests. Possible values are: 'Base Policy', 'Allow' or 'Deny'. Defaults to 'Base Policy' if no value provided.
    * For protocol "all", port must be set to "all".
    * For protocol “icmp”, port must be set to “ping”.

-> **NOTE:** If you are using/upgraded to Aviatrix Terraform Provider R1.5+, and an FQDN resource was originally created with a provider version <R1.5, you must modify your configuration file to match current format, and do ‘terraform refresh’ to update the state file to current format.


## Import

**fqdn** can be imported using the `fqdn_tag`, e.g.

```
$ terraform import aviatrix_fqdn.test fqdn_tag
```

## Notes
### FireNet
If FQDN is enabled on a gateway for the purposes of the Aviatrix FireNet Solution, you may run into an error requiring SNAT to be disabled when associating the gateway with the firewall (for reasons as described in the [note](#single_ip_snat) below). Please add an explicit dependency (`depends_on`) on the **aviatrix_firenet** resource to ensure the FireNet attachment completes first, before FQDN is enabled for that gateway.

### enable_vpc_dns_server
In order for the FQDN feature to be enabled, the corresponding gateway's `enable_vpc_dns_server` must be set to `false` at creation. FQDN will automatically enable that feature, which will cause a diff in the state. Please add `lifecycle { ignore_changes = [enable_vpc_dns_server] }` within that gateway's resource block in order to workaround this known issue. See [here](https://www.terraform.io/language/meta-arguments/lifecycle) for more information about the `lifecycle` attribute in Terraform.

### single_ip_snat
In order for the FQDN feature to be enabled, `single_ip_snat` must be set to true in the corresponding gateway. If it is not set at gateway creation, creation of FQDN resource will automatically enable SNAT and users must rectify the diff in the Terraform state by setting `single_ip_snat = true` in their gateway resource. An alternative is to add the utilise the `lifecycle` options in that gateway to ignore any changes, as described in the above bullet point.
