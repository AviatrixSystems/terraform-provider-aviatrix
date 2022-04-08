---
layout: "aviatrix"
page_title: "Migrating from Security Domain to Network Domain"
description: |-
  Aviatrix Network Domain Migration Guide
---

# Aviatrix Network Domain Migration Guide

## USAGE
As of Controller 6.7+, security domain will be renamed to network domain. Resources and attributes whose name includes security domain will be deprecated in future releases. Please follow along for guidance on the migration process to ensure a smooth transition.

---
## Migration Steps

- For resource **aviatrix_aws_tgw**:
  - Please set attribute `manage_security_domain` to false, and use the standalone resource **aviatrix_aws_tgw_network_domain** instead.

- For resources **aviatrix_aws_tgw_security_domain**, **aviatrix_segmentation_security_domain**, **aviatrix_segmentation_security_domain_association** and **aviatrix_segmentation_security_domain_connection_policy**:
  1. Remove the states of the resources.
  2. In the configuration file, rename the resources by replacing "security_domain" with "network_domain".
  3. Import the existing infrastructures to the renamed resources.

- For resources **aviatrix_aws_tgw_connect**, **aviatrix_aws_tgw_directconnect**, **aviatrix_aws_tgw_vpc_attachment**:
  1. Remove the states of the resources.
  2. In the configuration file, rename the attribute `security_domain_name` to `network_domain_name`.
  3. Import the existing infrastructures to the resources.
  