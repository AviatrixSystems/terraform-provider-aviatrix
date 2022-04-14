---
layout: "aviatrix"
page_title: "Migrating from Security Domain to Network Domain"
description: |-
  Aviatrix Network Domain Migration Guide
---

# Aviatrix Network Domain Migration Guide

## USAGE
As of Controller version 6.7+ and provider version 2.22.0+, security domain will be renamed to network domain. Resources and attributes whose name includes security domain will be deprecated in future releases. Please follow along for guidance on the migration process to ensure a smooth transition.

---
## Migration Steps

- For resource **aviatrix_aws_tgw**:
  - If `security_domains` is currently being used in the resource:
    - Please set attribute `manage_security_domain` to false, and use the standalone resource **aviatrix_aws_tgw_network_domain** instead.
  - If `manage_security_domain` is already set to false, and the resource **aviatrix_aws_tgw_security_domain** is being used:
    - Please refer to the instructions for **aviatrix_aws_tgw_security_domain** below.
     

- For resources **aviatrix_aws_tgw_security_domain**, **aviatrix_segmentation_security_domain**, **aviatrix_segmentation_security_domain_association** and **aviatrix_segmentation_security_domain_connection_policy**:
  1. Remove the states of the resources.
     - Please refer to this [link](https://www.terraform.io/cli/commands/state/rm) for instructions on removing the state of a resource.
  2. In the configuration file, rename the resources by replacing `security_domain` with `network_domain`.
  3. Import the existing infrastructures to the renamed resources.
     - Please refer to this [link](https://www.terraform.io/cli/import) for instructions on importing the resource.

- For resources **aviatrix_aws_tgw_connect**, **aviatrix_aws_tgw_directconnect**, **aviatrix_aws_tgw_vpc_attachment**:
  1. Remove the states of the resources.
     - Please refer to this [link](https://www.terraform.io/cli/commands/state/rm) for instructions on removing the state of a resource.
  2. In the configuration file, rename the attribute `security_domain_name` to `network_domain_name`.
  3. Import the existing infrastructures to the resources.
     - Please refer to this [link](https://www.terraform.io/cli/import) for instructions on importing the resource.
  