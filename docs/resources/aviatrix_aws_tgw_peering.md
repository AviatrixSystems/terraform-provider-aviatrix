---
subcategory: "TGW Orchestrator"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_aws_tgw_peering"
description: |-
  Creates and manages Aviatrix AWS TGW peerings, within a region or across regions
---

# aviatrix_aws_tgw_peering

The **aviatrix_aws_tgw_peering** resource allows the creation and management of Aviatrix-created peerings between AWS TGWs. The two TGWs can be in the same region or in different regions.

-> **NOTE:** Peering two TGWs in the same region requires Controller version 10.2.0 or later. Earlier Controllers only allow peering across regions.

## Example Usage

```hcl
# Create an Aviatrix AWS Tgw Peering
resource "aviatrix_aws_tgw_peering" "test" {
  tgw_name1 = "tgw1"
  tgw_name2 = "tgw2"
}
```
```hcl
# Peer an Aviatrix AWS Tgw with a TGW that the Controller does not manage
resource "aviatrix_aws_tgw_peering" "external" {
  tgw_name1       = "tgw1"
  tgw_name2       = "partner-tgw"
  peer_tgw_id     = "tgw-0a1b2c3d4e5f60718"
  peer_account_id = "123456789012"
  peer_region     = "us-west-2"
  remote_cidrs    = ["10.50.0.0/16", "10.60.0.0/16"]
}

# Route the remote CIDRs from the "prod" domain on tgw1 to the peering
resource "aviatrix_aws_tgw_peering_domain_conn" "prod_to_partner" {
  tgw_name1    = "tgw1"
  domain_name1 = "prod"
  tgw_name2    = "tgw1"
  domain_name2 = "peering_partner-tgw"

  depends_on = [aviatrix_aws_tgw_peering.external]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `tgw_name1` - (Required) This parameter represents name of the first AWS TGW to make a peer pair.
* `tgw_name2` - (Required) This parameter represents name of the second AWS TGW to make a peer pair. When peering with a TGW that the Controller does not manage, this is a label for that TGW.

### Optional
* `peer_tgw_id` - (Optional) AWS ID of a TGW that the Controller does not manage. Set it together with `peer_account_id` and `peer_region` to peer `tgw_name1` with that TGW. Changing it replaces the peering. Available as of provider version R10.2.0+.
* `peer_account_id` - (Optional) 12-digit AWS account ID that owns the TGW in `peer_tgw_id`. The account must be onboarded on the Controller, because the Controller accepts the peering on that account's side. Changing it replaces the peering. Available as of provider version R10.2.0+.
* `peer_region` - (Optional) AWS region of the TGW in `peer_tgw_id`. Changing it replaces the peering. Available as of provider version R10.2.0+.
* `remote_cidrs` - (Optional) Set of CIDRs behind the TGW in `peer_tgw_id`. Requires `peer_tgw_id`. Every network domain connected to `peering_<tgw_name2>` gets a route to each CIDR through the peering, and so do the VPCs attached to those domains. Changing it updates the peering in place. Removing all of them removes those routes. Terraform manages the whole list: if `remote_cidrs` is left out, CIDRs set through the Controller API are removed on the next apply. If the peering is created but setting `remote_cidrs` fails, Terraform marks the peering tainted and the next apply replaces it; run `terraform untaint` first to keep the peering and retry only `remote_cidrs`. Available as of provider version R10.2.0+.

-> **NOTE:** Peering with a TGW that the Controller does not manage requires Controller version 10.2.0 or later. `tgw_name1` must be a TGW that the Controller manages. The Controller creates the peering attachment and a `peering_<tgw_name2>` network domain on `tgw_name1` only, and keeps no record of the other TGW.

-> **NOTE:** Such a peering carries traffic once three things are in place. First, `remote_cidrs` lists the CIDRs behind the other TGW. Second, a network domain on `tgw_name1` is connected to `peering_<tgw_name2>`, for example with `aviatrix_aws_tgw_peering_domain_conn` using `tgw_name1` for both TGWs, as in the example above. Third, the owner of the other TGW adds routes back to `tgw_name1`; the Controller does not change route tables on that TGW.

## Import

**aws_tgw_peering** can be imported using the `tgw_name1` and `tgw_name2`, e.g.

```
$ terraform import aviatrix_aws_tgw_peering.test tgw_name1~tgw_name2
```

For a peering with a TGW that the Controller does not manage, use the label as `tgw_name2`. `peer_tgw_id`, `peer_account_id`, `peer_region` and `remote_cidrs` are read from the Controller.
