---
layout: "aviatrix"
page_title: "Introduction to Gateway Group"
description: |-
  Aviatrix Gateway Group Feature Introduction
---

# Introduction to Aviatrix Gateway Group feature

## What is Gateway Group
As of Controller version 7.0+ and provider version 3.0.0+, gateway group feature will be introduced. Gateway group feature allows creation of multiple ha gateways under the same primary gateway. Only Spoke Gateway is supported in 7.0. Ha gateway are only supported to be created one by one.

## How Spoke Gateway Group is supported in 7.0
As in Controller version 7.0, a new resource aviatrix_spoke_ha_gateway is added to be used for Spoke HA Gateway creation. In order to be able to use the new resource, attribute `network_domain` needs to be set as false, which makes aviatrix_spoke_gateway only be used to create primary spoke gateway.

---
## Migration Steps

- For resource **aviatrix_spoke_gateway**:
  - If HA is currently not enabled in the resource:
    - Please set attribute `manage_ha_gateway` to false, and do a "terraform refresh" to set its value to the state file.
  - If HA is currently being enabled in the resource:
    - Please set attribute `manage_ha_gateway` to false, and do a "terraform refresh" to set its value to the state file and remove Spoke HA gateway status from the state file
    - Please create a new resource **aviatrix_spoke_ha_gateway** with local name and map the settings of the above Spoke HA Gateway into the new resource
    - Please run "terraform import aviatrix_spoke_ha_gateway.local_name primary_spoke_gateway + '-hagw''" to import Spoke HA Gateway status into state file
  