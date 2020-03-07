---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_fqdn"
description: |-
  Manages Aviatrix FQDN filtering for gateways
---

# aviatrix_fqdn

The **aviatrix_fqdn** resource manages FQDN filtering for Aviatrix gateways.

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

  domain_names {
    fqdn  = "facebook.com"
    proto = "tcp"
    port  = "443"
  }

  domain_names {
    fqdn  = "reddit.com"
    proto = "tcp"
    port  = "443"
  }
}
```

## Argument Reference

The following arguments are supported:

* `fqdn_tag` - (Required) FQDN Filter tag name.
* `fqdn_enabled` - (Optional) FQDN Filter tag status. Valid values: true, false.
* `fqdn_mode` - (Optional) Specify FQDN mode: whitelist or blacklist. Valid values: "white", "black".
* `gw_filter_tag_list` - (Optional) A list of gateways to attach to the specific tag.
  * `gw_name` - (Required) Name of the gateway to attach to the specific tag.
  * `source_ip_list` - (Optional) List of source IPs in the VPC qualified for a specific tag.
* `domain_names` - (Optional) One or more domain names in a list with details as listed below:
  * `fqdn` - (Required) FQDN. Example: "facebook.com".
  * `proto` - (Required) Protocol. Valid values: "all", "tcp", "udp", "icmp".
  * `port` - (Required) Port. Example "25".
    * For protocol "all", port must be set to "all".
    * For protocol “icmp”, port must be set to “ping”.

-> **NOTE:** If you are using/upgraded to Aviatrix Terraform Provider R1.5+, and an fqdn resource was originally created with a provider version <R1.5, you must modify your configuration file to match current format, and do ‘terraform refresh’ to update the state file to current format.

-> **NOTE:** In order for the FQDN feature to be enabled, `single_ip_snat` must be set to true in the specified gateway. If it is not set at gateway creation, creation of FQDN resource will automatically enable SNAT and users must rectify the diff in the Terraform state by setting `single_ip_snat = true` in their gateway resource.

-> **NOTE:** In order for the FQDN feature to be enabled, the corresponding gateway's `enable_vpc_dns_server` must be set to `false` at creation. FQDN will automatically enable that feature, which will cause a diff in the state. Please add `lifecycle { ignore_changes = [enable_vpc_dns_server] }` within that gateway's resource block in order to workaround this known issue. See [here](https://www.terraform.io/docs/configuration/resources.html#lifecycle-lifecycle-customizations) for more information about the `lifecycle` attribute in Terraform.

## Import

**fqdn** can be imported using the `fqdn_tag`, e.g.

```
$ terraform import aviatrix_fqdn.test fqdn_tag
```
