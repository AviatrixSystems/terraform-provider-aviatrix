---
layout: "aviatrix"
page_title: "Introduction to Gateway Group"
description: |-
  Aviatrix Gateway Group Feature Introduction
---

# Introduction to Aviatrix Gateway Group feature

## What is Gateway Group
The Gateway Group feature allows for better horizontal scaling by moving from the current "primary + HA" gateway-pair
model, to support a gateway grouping, represented by a "primary + N number of HA gateways". The Gateway Group feature
allows for users to create N number of HA gateways under a primary gateway.

For 7.0 release, only spoke gateways will be supported and HA gateways can only be created sequentially (one by one).
In future releases, transit gateways will support this feature as well, and restrictions will be removed.

## How Spoke Gateway Group is supported in 7.0
In Controller 7.0 and provider version 3.0, a new resource **aviatrix_spoke_ha_gateway** has been added to support spoke
HA gateway creation.

In order to be able to use this new resource, the attribute `manage_ha_gateway` must be set to false in
**aviatrix_spoke_gateway**; this will mark the **aviatrix_spoke_gateway** to be used only as the primary.

---
## Migration Steps

- For resource **aviatrix_spoke_gateway**:
  - If HA is currently not enabled in the resource:
    - Please set attribute `manage_ha_gateway` to false, and do a "terraform refresh" to set its value to the state file
  - If HA is currently being enabled in the resource:
    - Please set attribute `manage_ha_gateway` to false and remove all the attributes with prefix "ha_", then do a "terraform apply" to set `manage_ha_gateway = false` to the state file and remove Spoke HA gateway status from the state file as well
    - Please create a new resource **aviatrix_spoke_ha_gateway** and map the settings of the above Spoke HA Gateway into the new resource
    - Please run "terraform import aviatrix_spoke_ha_gateway.local_name primary_spoke_gateway + '-hagw''" to import Spoke HA Gateway status into state file
