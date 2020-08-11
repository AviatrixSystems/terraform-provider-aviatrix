---
subcategory: "Accounts"
layout: "aviatrix"
page_title: "Aviatrix: aviatrix_rbac_group_permission_attachment"
description: |-
  Creates and manages Aviatrix RBAC group permission attachments
---

# aviatrix_rbac_group_permission_attachment

The **aviatrix_rbac_group_permission_attachment** resource allows the creation and management of permission attachments to Aviatrix (Role-Based Access Control) RBAC groups.

## Example Usage

```hcl
# Create an Aviatrix Rbac Group Permission Attachment
resource "aviatrix_rbac_group_permission_attachment" "test_attachment" {
  group_name      = "write_only"
  permission_name = "all_write"
}
```

## Argument Reference

The following arguments are supported:

### Required
* `group_name` - (Required) This parameter represents the name of a RBAC group.
* `permission_name` - (Required) This parameter represents the permission to attach to the RBAC group.

Valid `permission_name` values:

* "all_dashboard_write"
* "all_accounts_write"
* "all_gateway_write"
* "all_tgw_orchestrator_write"
* "all_transit_network_write":
* "all_firewall_network_write"
* "all_cloud_wan_write"
* "all_peering_write"
* "all_site2cloud_write"
* "all_openvpn_write"
* "all_security_write"
* "all_useful_tools_write"
* "all_troubleshoot_write"
* "all_write"

-> **NOTE:** If "all_write" is specified as the value for `permission_name`, all permissions will be attached to the specified RBAC group; there is then no need to specify any more permission attachments for that RBAC group.

## Import

**rbac_group_permission_attachment** can be imported using the `group_name` and `permission_name`, e.g.

```
$ terraform import aviatrix_rbac_group_permission_attachment.test group_name~permission_name
```
