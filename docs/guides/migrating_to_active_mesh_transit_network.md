---
layout: "aviatrix"
page_title: "Migrating from Classic Aviatrix Encrypted Transit Network to Aviatrix ActiveMesh Transit Network"
description: |-
The Migrating Aviatrix Transit Network Guide
---

# Aviatrix Transit Network Migration Guide

## USAGE
For customers who are currently already using Terraform to manage their infrastructure including Classic Aviatrix Encrypted Transit Network prior to Controller 6.6 and are looking to migrate their Classic Aviatrix Encrypted Transit Network to Aviatrix ActiveMesh Transit Network and upgrade their Controller to 6.6+, please follow along for guidance on the migration process to ensure a smooth transition.

---
## Migration Steps

If your Transit Network is built prior to Aviatrix software release 5.1, it’s very likely that the Transit Network is non-ActiveMesh deployment where the IPSec tunnels between Spoke GWs and Transit GWs are in Active/Standby mode (i.e. only one IPSec tunnel is carrying the data traffic).
1. Detach Aviatrix Spoke Gateway from Transit Network:
   - If Spoke was attached to Transit by aviatrix_spoke_transit_attachment, just destroy aviatrix_spoke_transit_attachment resource;
   - If Spoke was attached to Transit by setting “transit_gw” inside aviatrix_spoke_gateway, just remove “transit_gw” do terraform apply
2. Detach Aviatrix Transit Gateway from all other peerings or connections
3. Set “enable_active_mesh” to true in both aviatrix_spoke_gateway and aviatrix_transit_gateway to enable active mesh for both Spoke and Transit Gateway
4. Reattach Aviatrix Spoke Gateway to Transit Network
   - If Spoke was attached to Transit by aviatrix_spoke_transit_attachment, just recreate aviatrix_spoke_transit_attachment resource;
   - If Spoke was attached to Transit by setting “transit_gw” inside aviatrix_spoke_gateway, just reset “transit_gw” and do terraform apply
