---
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_geo_vpn"
description: |-
  Enables and Manages the Aviatrix Geo VPN
---

# aviatrix_geo_vpn

The aviatrix_geo_vpn resource enables and manages the Aviatrix Geo VPN.

## Example Usage

```hcl
# Create an Aviatrix Geo VPN
resource "aviatrix_geo_vpn" "test_geo_vpn" {
  account_name  = "devops-aws"
  domain_name   = "aviatrix.live"
  service_name  = "vpn"
  elb_dns_names = [
    "elb-test1-497f5e89.elb.us-west-1.amazonaws.com",
    "elb-test2-974f895e.elb.us-east-2.amazonaws.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) This parameter represents the name of a Cloud-Account in Aviatrix controller. Currently only supports AWS provider (cloud_type: 1).
* `domain_name` - (Required) The hosted domain name. It must be hosted by AWS Route53 or Azure DNS in the selected account.
* `service_name` - (Required) The hostname that users will connect to. A DNS record will be created for this name in the specified domain name.
* `elb_dns_names` - (Required) List of ELB names to attach to this Geo VPN name.

## Import

