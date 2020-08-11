---
subcategory: "CloudWAN"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_branch_router_transit_gateway_attachment"
description: |-
  Creates and manages a branch router and Aviatrix Transit Gateway attachment
---

# aviatrix_branch_router_transit_gateway_attachment

The **aviatrix_branch_router_transit_gateway_attachment** resource allows the creation and management of a branch router and Aviatrix Transit Gateway attachment

~> **NOTE:** Before creating this attachment the branch router must have its WAN interface and IP configured via the `aviatrix_branch_router_interface_config` resource. To avoid attempting to create the attachment before the interface and IP are configured use a `depends_on` meta-argument so that the `aviatrix_branch_router_interface_config` resource is created before the attachment.  

## Example Usage

```hcl
# Create an Aviatrix Branch Router and Transit Gateway attachment
resource "aviatrix_branch_router_transit_gateway_attachment" "test_branch_router_transit_gateway_attachment" {
  branch_name             = "branch-router"
  transit_gateway_name    = "transit-gw"
  connection_name         = "test-conn"
  transit_gateway_bgp_asn = 65000
  branch_router_bgp_asn   = 65001

  depends_on = [aviatrix_branch_router_interface_config.test_branch_router_interface_config]
}
```

## Argument Reference

The following arguments are supported:

### Required
* `branch_name` - Branch router name.
* `transit_gateway_name` - Aviatrix Transit Gateway name.
* `connection_name` - Connection name.
* `transit_gateway_bgp_asn` - BGP AS Number for transit gateway.
* `branch_router_bgp_asn` - BGP AS Number for branch router.

### Optional
* `phase1_authentication` - Phase 1 authentication algorithm. Default "SHA-256".
* `phase1_dh_groups` - Number of phase 1 Diffie-Hellman groups. Default "14".
* `phase1_encryption` - Phase 1 encryption algorithm. Default "AES-256-CBC".
* `phase2_authentication` - Phase 2 authentication algorithm. Default "HMAC-SHA-256".
* `phase2_dh_groups` - Number of phase 2 Diffie-Hellman groups. Default "14".
* `phase2_encryption` - Phase 2 encryption algorithm. Default "AES-256-CBC".
* `enable_global_accelerator` - Boolean enable AWS Global Accelerator. Default "false".
* `pre_shared_key` - Pre-shared Key.
* `local_tunnel_ip` - Local tunnel IP.
* `remote_tunnel_ip` - Remote tunnel IP.

## Import

**branch_router_transit_gateway_attachment** can be imported using the `connection_name`, e.g.

```
$ terraform import aviatrix_branch_router_transit_gateway_attachment.test connection-name
```
