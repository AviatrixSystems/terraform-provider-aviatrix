---
subcategory: "TGW Orchestrator"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_network_domain"
description: |-
  Creates and manages Aviatrix network domains
---

# aviatrix_aws_tgw_network_domain

The **aviatrix_aws_tgw_network_domain** resource allows the creation and management of Aviatrix network domains.

## Example Usage

```hcl
# Create an Aviatrix AWS TGW network domain
resource "aviatrix_aws_tgw" "test_aws_tgw" {
  account_name       = "devops"
  aws_side_as_number = "64512"
  region             = "us-east-1"
  tgw_name           = "test-AWS-TGW"
}

resource "aviatrix_aws_tgw_network_domain" "Default_Domain" {
  name     = "Default_Domain"
  tgw_name = aviatrix_aws_tgw.test_aws_tgw.id
}

resource "aviatrix_aws_tgw_network_domain" "Shared_Service_Domain" {
  name     = "Shared_Service_Domain"
  tgw_name = aviatrix_aws_tgw.test_aws_tgw.id
}

resource "aviatrix_aws_tgw_network_domain" "Aviatrix_Edge_Domain" {
  name     = "Aviatrix_Edge_Domain"
  tgw_name = aviatrix_aws_tgw.test_aws_tgw.id
}

resource "aviatrix_aws_tgw_peering_domain_conn" "default_sd_conn1" {
  tgw_name1     = aviatrix_aws_tgw.test_aws_tgw.id
  tgw_name2     = aviatrix_aws_tgw.test_aws_tgw.id
  domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
  domain_name2 = aviatrix_aws_tgw_network_domain.Default_Domain.name
}

resource "aviatrix_aws_tgw_peering_domain_conn" "default_sd_conn2" {
  tgw_name1     = aviatrix_aws_tgw.test_aws_tgw.id
  tgw_name2     = aviatrix_aws_tgw.test_aws_tgw.id
  domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
  domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}

resource "aviatrix_aws_tgw_peering_domain_conn" "default_sd_conn3" {
  tgw_name1     = aviatrix_aws_tgw.test_aws_tgw.id
  tgw_name2     = aviatrix_aws_tgw.test_aws_tgw.id
  domain_name1 = aviatrix_aws_tgw_network_domain.Default_Domain.name
  domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}

resource "aviatrix_aws_tgw_network_domain" "test" {
  name       = "test_domain"
  tgw_name   = aviatrix_aws_tgw.test_aws_tgw.id
  depends_on = [
    aviatrix_aws_tgw_network_domain.Default_Domain,
    aviatrix_aws_tgw_network_domain.Shared_Service_Domain,
    aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `name` - (Required) The name of the network domain.
* `tgw_name` - (Required) The AWS TGW name of the network domain.
* `aviatrix_firewall` - (Optional) Set to true if the network domain is to be used as an Aviatrix Firewall Domain for the Aviatrix Firewall Network. Valid values: true, false. Default value: false.
* `native_egress` - (Optional) Set to true if the network domain is to be used as a native egress domain (for non-Aviatrix Firewall Network-based central Internet bound traffic). Valid values: true, false. Default value: false.
* `native_firewall` - (Optional) Set to true if the network domain is to be used as a native firewall domain (for non-Aviatrix Firewall Network-based firewall traffic inspection). Valid values: true, false. Default value: false.

-> **NOTE:** Three default domains ("Aviatrix_Edge_Domain", "Default_Domain" and "Shared_Service_Domain") are required before the creation of other domains. Non-default domains should depend on default domains in order to get proper destroy sequence. The connections between three default domains should also be created using the resource `aviatrix_aws_tgw_peering_domain_conn`. 

## Import

**aws_tgw_network_domain** can be imported using the `name` and `tgw_name`, e.g.

```
$ terraform import aviatrix_aws_tgw_network_domain.test tgw_name~name
```
