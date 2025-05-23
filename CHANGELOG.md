## 8.0.0 (Unreleased)
### Notes:
- Supported Controller version: **UserConnect-8.0.0-1000.2432**
- Supported Terraform version **v1.x**

### Enhancements:
1. Added field tls_profile to **aviatrix_distributed_firewalling_policy_list** to be able specify a TLS profile in a DCF policy.
2. Added support for importing the **aviatrix_edge_proxy_profile** resource.
3. Allow special characters in Azure tags on gateway instances, to align with UI behavior.
4. Extended k8s examples for the **aviatrix_smart_group** resource.
5. Implemented new resources **aviatrix_edge_megaport** and **aviatrix_edge_megaport_ha** for deployment of spoke gateway in Megaport.
6. The **aviatrix_transit_gateway** resource now supports deploying in Equinix and Megaport.
7. Added support for `enable_bgp_multihop` to **aviatrix_edge_spoke_external_device_conn**, **aviatrix_spoke_external_device_conn** and **aviatrix_transit_external_device_conn** resources.
8. The **aviatrix_edge_spoke_external_device_conn** resource now supports `enable_edge_underlay` which allows for configuring BGP on the underlay.
9. The **aviatrix_edge_spoke_external_device_conn** resource now supports `enable_bgp_lan_activemesh` enabling activemesh support for BGP peering on the LAN.
10. The **aviatrix_edge_spoke_external_device_conn** resource now supports `enable_jumbo_frame`.

### Deprecations
1. The argument **aviatrix_dns_profile** has been removed from the following resources, as it was never used.
  - **aviatrix_edge_csp**
  - **aviatrix_edge_equinix**
  - **aviatrix_edge_neo**
  - **aviatrix_edge_platform**
  - **aviatrix_edge_zededa**
2. The CloudN related resources **aviatrix_cloudn_registration** and **aviatrix_cloudn_transit_gateway_attachment** have been removed, as they are no longer supported.

### Bug Fixes:
1. Fixed an issue where `ha_eip` was not applied to an Azure transit HA gateway.
2. Fixed an issue where the **aviatrix_vpc_tracker** data source was returning incorrect `cloud_type` information.
3. Fixed an issue where retrieving a large number of objects could result in a 502 Proxy Error during terraform operations.


## 3.2.2
### Notes:
- Supported Controller version: **UserConnect-7.2.5090**
- Supported Terraform version **v1.x**

### Enhancements:
1. Updated documentation for **aviatrix_edge_equinix**.
    - Added examples for static and DHCP management interface.
    - Added examples of required Equinic resources.
    - Clarified the description of **management_egress_ip_prefix_list**.

### Bug Fixes:
1. Fixed an issue where retrieving a large number of objects could result in a 502 Proxy Error during terraform operations.

## 3.2.1
### Notes:
- Supported Controller version: **UserConnect-7.2.4996**
- Supported Terraform version **v1.x**

### Bug Fixes:
1. Fixed a typo in the constants for the Kubernetes related properties for the smart group resource.

### Enhancements:
1. Added new attribute ``bgp_neighbor_status_polling_time`` to support the bgp bfd configuration in the following resources.
    - **aviatrix_edge_csp**
    - **aviatrix_edge_equinix**
    - **aviatrix_edge_gateway_selfmanaged**
    - **avaitrix_edge_platform**
    - **aviatrix_edge_zededa**
    - **aviatrix_spoke_gateway**
    - **aviatrix_edge_spoke_gateway**
    - **aviatrix_transit_gateway**
2. Added new attribute ``bgp_bfd`` and ``enable_bfd`` to support bgp_bfd configuration in the following resources.
    - **aviatrix_transit_external_device_conn**
    - **aviatrix_edge_spoke_external_device_conn**
    - **aviatrix_spoke_external_device_conn**
3. Added/Updated Edge Terraform documentation to include interface mapping.
4. Updated documentation for **aviatrix_firewall_instance** where ``user_data`` was not applicable to Palo Alto firewalls. There is no longer
   a restriction on using ``user_data`` on Palo Alto Firewalls.
5. Updated documentation for **aviatrix_spoke_gateway**. Clarified use of the ``included_advertised_spoke_routes`` attribute.
6. Updated documentation for **aviatrix_smart_group** with examples of external groups usage.
7. Improved performance for **aviatrix_smart_group** and **aviatrix_web_group** resources.
8. Updated the attribute ``enable_jumbo_frame`` default value in **aviatrix_transit_gateway**. For CSP transit gateways, ``enable_jumbo_frame`` will be set to true by default and for Edge EAT it will be set to false by default.

### Deprecations:
Customers can no longer re-bootstrap their PKI with a custom root CA using Terraform. However, this functionality remains available through the Controller UI for added flexibility.

The **aviatrix_dns_profile** resource has been removed.

## 3.2.0 (October 15, 2024)
### Notes:
- Supported Controller version: **UserConnect-7.2.4820**
- Supported Terraform version **v1.x**

### Bug Fixes:
1. Fixed issue in **aviatrix_edge_platform_device_onboarding** where performing subsequent a apply would continue to update the ``network`` configuration even when there were no changes made.
2. Fixed issue in **aviatrix_gateway_dnat** where configuring the interface in a policy-based Site-to-Cloud DNAT rule would trigger an error.
3. Fixed issue in **aviatrix_edge_platform_device_onboarding** where ``dns_server_ips`` configuration order was not preserved.
4. Fixed issue in **aviatrix_edge_platform_device_onboarding** where onboard resource failed to properly import.
5. Fixed issue in **aviatrix_site2cloud** where ``remote_subnet_cidr`` was not properly applied.
6. Fixed issue in **aviatrix_transit_gateway** that prevented the deployment Firenet Gateways in Azure China
7. Optimized **aviatrix_account** to significantly reduce the time required for Terraform operations (e.g., update, add, delete) involving hundreds of accounts. Previously, these operations could take tens of minutes, but with this fix, they now complete in tens of seconds.

### Enhancements:
1. Added proxy profile support to **aviatrix_edge_platform_device_onboarding**, enabling the specification of a proxy for onboarding Aviatrix Edge devices.
2. Updated documentation references by consolidating the legacy terms ``Cloudn`` and ``Multi-Cloud Transit`` under a single ``Edge`` subcategory.
3. Added support for Kubernetes Smart Groups, updating the **aviatrix_smart_group** resource to allow Smart Groups to be created from artifacts within Kubernetes clusters

#### Provider:
1. Added support for the Terraform provider to properly set the user-agent when making requests.

### Deprecations:
1. Removed the ``bandwidth`` attribute from the interface configuration for all Edge related resources.
2. Removed the ``http_access`` in **aviatrix_controller_config** as it longer has any effect.
3. Removed the ``keep_alive_via_lan_interface_enabled`` in **aviatrix_firenet** resource.

## 3.1.6 (October 15, 2024)
### Notes:
- Supported Controller version: **UserConnect-7.1.4183**
- Supported Terraform version **v1.x**

### Bug Fixes:
1. Fixed issue in **aviatrix_edge_platform_device_onboarding** where performing subsequent a apply would continue to update the ``network`` configuration even when there were no changes made.
2. Fixed issue in **aviatrix_gateway_dnat** where configuring the interface in a policy-based Site-to-Cloud DNAT rule would trigger an error.
3. Fixed issue in **aviatrix_edge_platform_device_onboarding** where ``dns_server_ips`` configuration order was not preserved.
4. Fixed issue in **aviatrix_site2cloud** where ``remote_subnet_cidr`` was not properly applied.

## 3.1.5 (July 18, 2024)
### Notes:
- Supported Controller version: **UserConnect-7.1.4105**
- Supported Terraform version: **v1.x**

### Bug Fixes:

1. Fixed issue in **resource_aviatrix_firewall_instance_association** and **resource_aviatrix_gateway** for Azure where we no longer require special handling of ``fqdn_lan_interface``  and ``lan_interface``.
2. Fixed issue in **aviatrix_edge_platform_device_onboarding** where importing was failing.
### Features:
#### Provider:
1. Added support for the Terraform provider to properly set the user-agent when making requests.


### Multi-Cloud Transit:
1. Add new attribute ``dns_server_ip`` and ``secondary_dns_server_ip`` in **aviatrix_edge_gateway_selfmanaged_ha** resource.


### Deprecations:

1. Deprecated ``http_access`` in **aviatrix_controller_config**. This configuration value no longer has any effect. It will be removed from the Aviatrix provider in the 3.2.0 release.

## 3.1.4 (January 11, 2024)
### Notes:
- Supported Controller version: **UserConnect-7.1.3006**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented support for VLAN and VRRP in **aviatrix_edge_gateway_selfmanaged**
2. Implemented a new attribute in **aviatrix_edge_spoke_external_device_conn** to support configuring BGP Advertised CIDRs:
   - ``manual_bgp_advertised_cidrs``
3. Implemented a new resource to replace **aviatrix_edge_vm_selfmanged** and **aviatrix_edge_spoke**:
   - **aviatrix_edge_gateway_selfmanged**
4. Implemented a new attribute in **aviatrix_spoke_transit_attachment** to support configuring tunnel count:
   - ``tunnel_count``

#### Secured Networking:
1. Implemented a new data source to list all **aviatrix_smart_group**:
   - **aviatrix_smart_groups**

### Enhancements:
1. Added support for setting ``ztp_file_type`` into state in **aviatrix_edge_gateway_selfmanaged**
2. Changed ``transit_gateway_name`` to "Optional" and "Computed" in **aviatrix_segmentation_network_domain_association**
3. Removed "ForceNew" property of ``number_of_retries`` and ``retry_interval`` in **aviatrix_edge_spoke_external_device_conn**
4. Enhanced retry mechanism for attaching spoke to transit in **resource_aviatrix_spoke_transit_attachment**

### Bug Fixes:
1. Fixed issue where the format of image file generated for Edge gateways was not correct
2. Fixed issue where the association fails in **aviatrix_segmentation_network_domain_association**
3. Fixed issue where ``edge_wan_interfaces`` being not set caused force replacement of **aviatrix_edge_spoke_transit_attachment**
4. Fixed issue where the format of ``attachment_name`` was not correct in **aviatrix_segmentation_network_domain_association**
5. Fixed issue where enabling Active Standby mode for Edge gateways fails
6. Fixed issue where ``transit_gateway_name`` was not set in TF state in **aviatrix_cloudn_transit_gateway_attachment**
7. Fixed a decoding issue in **aviatrix_cloudn_transit_gateway_attachment**
8. Fixed issue where resizing gateway size fails in **aviatrix_spoke_ha_gateway**
9. Fixed issue where updating ``remote_subnet`` fails in **aviatrix_transit_external_device_conn**
10. Fixed issue where updating ``ha_private_mode_subnet_zone`` fails for Private Mode in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
11. Fixed issue where Single IP SNAT was not set correctly in **aviatrix_spoke_gateway**

### Deprecations:
1. Deprecated ``keep_alive_via_lan_interface_enabled`` in **aviatrix_firenet**. It will be removed from the Aviatrix provider in the next upcoming 3.2.0 release


## 3.1.3 (November 02, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.1.2131**
- Supported Terraform version: **v1.x**

### Bug Fixes:
1. Fixed issue where ``terraform plan`` fails to read CloudN transit gateway attachment due to JSON decode error after controller was upgraded to 7.1.x in **aviatrix_cloudn_transit_gateway_attachment**


## 3.1.2 (August 29, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.1.2131**
- Supported Terraform version: **v1.x**

### Enhancements:
1. Changed ``insane_mode`` to "ForceNew" in **aviatrix_transit_gateway**

### Bug Fixes:
1. Fixed issue where default PSK was not configured correctly for default ``auth_type`` in **aviatrix_site2cloud**


## 3.1.1 (June 14, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.1.1794**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented a new resource to support Edge VM Selfmanaged:
   - **aviatrix_edge_vm_selfmanaged**
2. Implemented new resources to support Edge Platform device onboarding, Edge Platform and HA:
   - **aviatrix_edge_platform_device_onboarding**
   - **aviatrix_edge_platform**
   - **aviatrix_edge_platform_ha**
3. Implemented a new attribute to support ignoring DFW policy for SG orchestration in **aviatrix_distributed_firewalling_policy_list**:
   - ``exclude_sg_orchestration``
4. Implemented a new resource to support configuring how proxy handles bad origin certificates:
   - **aviatrix_distributed_firewalling_origin_cert_enforcement_config**
5. Implemented a new resource to allow customer directed change of MITM CA cert/key:
   - **aviatrix_distributed_firewalling_proxy_ca_config**

#### Copilot:
1. Implemented new resources to support Copilot simple and fault-tolerant deployments:
   - **aviatrix_copilot_simple_deployment**
   - **aviatrix_copilot_fault_tolerant_deployment**

### Enhancements:
1. Added support for updating ``enable_edge_active_standby`` and ``enable_edge_active_standby_preemptive`` in Edge resources
2. Updated valid range of ``insane_mode_tunnel_number`` to "2-50" in **aviatrix_edge_spoke_transit_attachment**
3. Added support for ``manual_bgp_advertised_cidrs`` in **aviatrix_edge_spoke_external_device_conn**
4. Changed ``profile_name`` to be optional in **aviatrix_remote_syslog**

### Bug Fixes:
1. Fixed issue where creating BGP underlay for Edge HA fails for **aviatrix_edge_spoke_external_device_conn**
2. Fixed issue where retrieving transit gateway peering info from a corrupted database caused provider crash
3. Fixed issue where task status check could fail due to proxy error
4. Fixed issue where enabling single IP SNAT fails during spoke gateway creation


## 3.1.0 (May 11, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.1**
- Supported Terraform version: **v1.x**

### Features:
#### Gateway:
1. Implemented support for configuring GRO/GSO in **aviatrix_gateway**

#### Multi-Cloud Transit:
1. Implemented support for configuring Local Identifier in **aviatrix_site2cloud**, **aviatrix_spoke_external_device_conn** and **aviatrix_transit_external_device_conn**
2. Implemented support for configuring GRO/GSO in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
3. Implemented support for BGP over GRE on a BGP Spoke in **aviatrix_spoke_external_device_conn**
4. Implemented new resources to support Edge Equinix and HA:
   - **aviatrix_edge_equinix**
   - **aviatrix_edge_equinix_ha**
5. Implemented new resources to support Edge NEO device onboarding, Edge NEO and HA:
   - **aviatrix_edge_neo_device_onboarding**
   - **aviatrix_edge_neo**
   - **aviatrix_edge_neo_ha**
6. Added a new attribute in **aviatrix_spoke_gateway** to support GCP spoke global VPC:
   - ``enable_global_vpc``
7. Implemented new resources to support GCP global VPC excluded instance and tagging settings:
   - **aviatrix_global_vpc_excluded_instance**
   - **aviatrix_global_vpc_tagging_settings**
8. Added a new attribute in **aviatrix_edge_spoke_external_device_conn** to support BGP over WAN underlay:
   - ``enable_edge_underlay``
9. Added a new attribute in **aviatrix_edge_spoke_external_device_conn** to support Prepend AS Path:
   - ``prepend_as_path``
10. Added a new attribute in **aviatrix_edge_spoke_transit_attachment** to support multiple WAN interfaces:
    - ``edge_wan_interfaces``
11. Implemented a new data source to get Edge gateway WAN interface IP address:
    - **aviatrix_edge_gateway_wan_interface_discovery**

#### Settings:
1. Implemented a new data source to collect controller metadata:
   - **aviatrix_controller_metadata**

#### Copilot:
1. Implemented new resources to support QoS class and QoS policy list:
   - **aviatrix_qos_class**
   - **aviatrix_qos_policy_list**

### Enhancements:
1. Added support for ``fqdn`` as one of the attributes under ``selector`` in **aviatrix_smart_group**
2. Added support for the "#" character in Azure gateway tags
3. Added support for enabling BGP over LAN for Azure transit in update in **aviatrix_transit_gateway**
4. Changed ``cloud_type``, ``account_name``, ``insane_mode`` and ``insane_mode_az`` to "ForceNew" in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
5. Removed ``bgp_lan_interfaces_count``'s default value in **aviatrix_transit_gateway**
6. Removed the option to config proxies in private mode config in **aviatrix_controller_private_mode_config**
7. Restored the following attributes in **aviatrix_cloudn_transit_gateway_attachment**:
   - ``enable_dead_peer_detection``
   - ``enable_learned_cidrs_approval``
   - ``approved_cidrs``
8. Added support for "ddog-gov.com" in **aviatrix_datadog_agent**
9. Added support for connection with HA in **aviatrix_edge_spoke_external_device_conn**

### Bug Fixes:
1. Fixed issue where ``max_vpn_conn`` was not properly set in TF state and could not be updated in **aviatrix_gateway**
2. Fixed issue where retries inconsistently fails in **aviatrix_edge_spoke_transit_attachment**

### Preview Features:
1. Implemented support of MITM/IDS: Request URL filtering:
   - New resource: **aviatrix_web_group**
   - New attributes in **aviatrix_distributed_firewalling_policy_list**:
      - ``web_groups``
      - ``flow_app_requirement``
      - ``decrypt_policy``

### Deprecations:
1. Removed support of **aviatrix_splunk_logging**, **aviatrix_filebeat_forwarder** and **aviatrix_sumologic_forwarder**


## 3.0.7 (January 10, 2024)
### Notes:
- Supported Controller version: **UserConnect-7.0.2239**
- Supported Terraform version: **v1.x**

### Enhancements:
1. Enhanced retry mechanism for attaching spoke to transit in **resource_aviatrix_spoke_transit_attachment**


## 3.0.6 (June 30, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.0.1768**
- Supported Terraform version: **v1.x**

### Bug Fixes:
1. Fixed issue where in-place update to enable Transit FireNet with GWLB fails in **aviatrix_transit_gateway**
2. Fixed issue where newly launched Controller fails to upgrade to 7.1
3. Fixed issue where retrieving Transit Gateway Peering info might cause provider to crash in **aviatrix_transit_gateway_peering**


## 3.0.5 (April 13, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.0.1724**
- Supported Terraform version: **v1.x**

### Features:
#### Settings:
1. Implemented a new resource to support Controller Security Group Management:
   - **aviatrix_controller_access_allow_list_config**


## 3.0.4 (April 13, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.0.1724**
- Supported Terraform version: **v1.x**

### Enhancements:
1. Added support for whitespace characters in **aviatrix_firewall_policy** ``description``


## 3.0.3 (March 24, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.0.1601**
- Supported Terraform version: **v1.x**

### Features:
#### CloudN:
1. Restored support of CloudN transit attachment CIDR approval features in **aviatrix_cloudn_transit_gateway_attachment**:
    - ``enable_dead_peer_detection``
    - ``enable_learned_cidrs_approval``
    - ``approved_cidrs``

#### Multi-Cloud Transit:
1. Implemented support for adding additional BGP over LAN interfaces to Azure Transit without redeploying **aviatrix_transit_gateway**

### Enhancements:
1. Added support of "datadoghq.com" for ``site`` in **aviatrix_datadog_agent**
2. Updated following attributes to "ForceNew" in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**:
   - ``cloud_type``
   - ``account_name``
   - ``insane_mode_az``
   - ``insane_mode``


## 3.0.2 (March 07, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.0.1577**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented support for BGP over LAN on Spoke in the following resources:
    - **aviatrix_spoke_gateway**
    - **aviatrix_spoke_external_device_conn**

### Enhancements:
1. Added support of ``#`` as a valid character in resource tags for Azure CSP
2. Added support of ``name`` as one of the attributes under ``selector`` in **aviatrix_smart_group**

### Bug Fixes:
1. Fixed issue where white space is allowed for ``phase1_remote_identifier`` in **aviatrix_site2cloud**, **aviatrix_transit_external_device_conn** and **aviatrix_spoke_external_device_conn**
2. Fixed issue where ``cloud_image_id`` is allowed for Azure in **aviatrix_firewall_instance**
3. Fixed issue where S2C creation fails for an exception error
4. Fixed re-ordering issue of ``rtb_list1`` and ``rtb_list2`` in **aviatrix_aws_peer**


## 3.0.1 (January 09, 2023)
### Notes:
- Supported Controller version: **UserConnect-7.0.1373**
- Supported Terraform version: **v1.x**

### Features:
#### Site2Cloud:
1. Implemented support for remote identification using empty string in the following resources:
    - **aviatrix_site2cloud**
    - **aviatrix_spoke_external_device_conn**
    - **aviatrix_transit_external_device_conn**

### Bug Fixes:
1. Fixed issue where FQDN tag's ``source_ip_list`` requires executing terraform apply twice for more than 2 gateways
2. Fixed issue where ``eip`` is not valid for creating Azure **aviatrix_spoke_ha_gateway**

### Deprecations:
1. The following resource is removed:
    - **aviatrix_transit_cloudn_conn**


## 3.0.0 (October 8, 2022)
### Notes:
- Supported Controller version: **UserConnect-7.0**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented a new resource to support Edge CSP:
   - **aviatrix_edge_csp**
2. Implemented a new resource to support Centralized Transit FireNet:
   - **aviatrix_centralized_transit_firenet**
3. Implemented a new data source to get the list of Spoke Gateways:
   - **aviatrix_spoke_gateways**
4. Implemented support for Gateway Group:
   Please check [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/introduction_to_gateway_group) for more information:
   - New resource: **aviatrix_spoke_ha_gateway**
   - New attributes in **aviatrix_spoke_gateway**:
      - ``manage_ha_gateway``
5. Implemented a new resource to support configuring Distributed Firewalling:
   - **aviatrix_distributed_firewalling_config**
6. Implemented a new resource to support Distributed Firewalling Intra VPC:
   - **aviatrix_distributed_firewalling_intra_vpc**

#### Security:
1. Implemented a new resource to support configuring FQDN Global settings:
   - **aviatrix_fqdn_global_config**

### Enhancements:
1. Updated ``transit_gateway_name`` to be an optional attribute in **aviatrix_segmentation_network_domain_association**
2. Removed validation for controller IP address during import in **aviatrix_controller_security_group_management_config**
3. Renamed **aviatrix_app_domain** to **aviatrix_smart_group**
4. Renamed **aviatrix_microseg_policy_list** to **aviatrix_distributed_firewalling_policy_list**
5. Added API key for server and client authentication

### Bug Fixes:
1. Fixed issue where configuring large number of policies errors out in **aviatrix_gateway_dnat**
2. Fixed issue where delta shows on ``region`` when same ``region`` is provided in different formats
3. Fixed issue where ``rx_queue_size`` does not apply on HA in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
4. Fixed issue where route table list is missing in state file while importing **aviatrix_aws_peering**
5. Fixed issue where HA creation does not use the right CMK for volume encryption when CMK is provided in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**

### Deprecations:
1. The following resources are removed:
   - **aviatrix_arm_peer**
   - **aviatrix_aws_tgw_security_domain**
   - **aviatrix_aws_tgw_security_domain_connection**
   - **aviatrix_edge_caag**
   - **aviatrix_segmentation_security_domain**
   - **aviatrix_segmentation_security_domain_association**
   - **aviatrix_segmentation_security_domain_connection_policy**
   - **aviatrix_spoke_vpc**
   - **aviatrix_transit_vpc**
2. Removed support for the attribute ``sync_to_ha`` from the following resources:
   - **aviatrix_gateway_snat**
   - **aviatrix_gateway_dnat**
3. Removed support for the attribute ``tag_list`` from the following resources and their respective data sources:
   - **aviatrix_gateway**
   - **aviatrix_spoke_gateway**
   - **aviatrix_transit_gateway**
4. Removed support for ``manage_transit_gateway_attachment`` and ``transit_gw`` from the resource **aviatrix_spoke_gateway**
5. Removed support for managing in-line firewall instance associations by removing the following attributes from the resource **aviatrix_firenet** and its respective data source:
   - ``manage_firewall_instance_association`` (from the resource)
   - ``firewall_instance_association`` (from both)
6. Removed support for managing in-line TGW network domains, VPC attachments and transit gateway attachments by removing the following attributes from the resource **aviatrix_aws_tgw**:
   - ``manage_security_domain``
   - ``security_domains``
   - ``manage_vpc_attachment``
   - ``attached_vpc``
   - ``manage_transit_gateway_attachment``
   - ``attached_transit_gateway``
7. Removed support for the attribute ``security_domain_name`` from the following resources:
   - **aviatrix_aws_tgw_connect**
   - **aviatrix_aws_tgw_directconnect**
   - **aviatrix_aws_tgw_vpc_attachment**
8. Deprecated **aviatrix_trans_peer** and it will be removed in Aviatrix provider 3.0.1


## 2.24.4 (January 10, 2024)
### Notes:
- Supported Controller version: **UserConnect-6.9.822**
- Supported Terraform version: **v1.x**

### Enhancements:
1. Enhanced retry mechanism for attaching spoke to transit in **resource_aviatrix_spoke_transit_attachment**


## 2.24.3 (March 20, 2023)
### Notes:
- Supported Controller version: **UserConnect-6.9.349**
- Supported Terraform version: **v1.x**

### Features:
#### CloudN:
1. Restored support of CloudN transit attachment CIDR approval features in **aviatrix_cloudn_transit_gateway_attachment**:
    - ``enable_dead_peer_detection``
    - ``enable_learned_cidrs_approval``
    - ``approved_cidrs``

### Enhancements:
1. Added support of "datadoghq.com" for ``site`` in **aviatrix_datadog_agent**


## 2.24.2 (January 06, 2023)
### Notes:
- Supported Controller version: **UserConnect-6.9.282**
- Supported Terraform version: **v1.x**

### Features:
#### Site2Cloud:
1. Implemented support for remote identification using empty string in the following resources:
    - **aviatrix_site2cloud**
    - **aviatrix_spoke_external_device_conn**
    - **aviatrix_transit_external_device_conn**

### Bug Fixes:
1. Fixed issue where FQDN tag's ``source_ip_list`` requires executing terraform apply twice for more than 2 gateways


## 2.24.1 (September 30, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.9.161**
- Supported Terraform version: **v1.x**

### Features:
#### Site2Cloud:
1. Implemented support for Certificate-based HA Gateway Remote Identifier for Site2Cloud VPN:
    - New attributes in **aviatrix_site2cloud**:
        - ``backup_remote_identifier``

### Bug Fixes:
1. Fixed issue where route-based Single IP HA tunnel S2C creation fails for **aviatrix_site2cloud**


## 2.24.0 (September 09, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.9**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented support for exposing BGP LAN interface info on transit in Azure via ``bgp_lan_ip_list`` and ``ha_bgp_lan_ip_list``
2. Implemented support for multiple disjoint port ranges for **aviatrix_microseg_policy_list**

### Enhancements:
1. Added support for updating ``bgp_md5_key`` and ``backup_bgp_md5_key`` for **aviatrix_spoke_external_device_conn** and **aviatrix_transit_external_device_conn**
2. Optimized the read operation for **aviatrix_transit_firenet_policy**

### Bug Fixes:
1. Fixed issue where NAT config has ordering issues that would cause refresh problems for **aviatrix_gateway_dnat** and **aviatrix_gateway_snat**
2. Fixed issue where provider crashes for private mode config
3. Fixed issue where exported SNAT/DNAT interface shows tunnel ID when using Transit connection or route-based S2C
4. Fixed issue where creating FQDN gateway with ``fqdn_lan_interface`` causes replacement


## 2.23.6 (January 09, 2024)
### Notes:
- Supported Controller version: **UserConnect-6.8.1826**
- Supported Terraform version: **v1.x**

### Enhancements:
1. Enhanced retry mechanism for attaching spoke to transit in **resource_aviatrix_spoke_transit_attachment**


## 2.23.5 (March 14, 2023)
### Notes:
- Supported Controller version: **UserConnect-6.8.1509**
- Supported Terraform version: **v1.x**

### Features:
#### CloudN:
1. Restored support of CloudN transit attachment CIDR approval features in **aviatrix_cloudn_transit_gateway_attachment**:
    - ``enable_dead_peer_detection``
    - ``enable_learned_cidrs_approval``
    - ``approved_cidrs``


## 2.23.4 (February 02, 2023)
### Notes:
- Supported Controller version: **UserConnect-6.8.1483**
- Supported Terraform version: **v1.x**

### Bug Fixes:
1. Fixed issue where creating multiple **aviatrix_trans_peer** between the same gateways with different ``reachable_cidrs`` errors out


## 2.23.3 (January 06, 2023)
### Notes:
- Supported Controller version: **UserConnect-6.8.1455**
- Supported Terraform version: **v1.x**

### Features:
#### Site2Cloud:
1. Implemented support for remote identification using empty string in the following resources:
    - **aviatrix_site2cloud**
    - **aviatrix_spoke_external_device_conn**
    - **aviatrix_transit_external_device_conn**

### Bug Fixes:
1. Fixed issue where FQDN tag's ``source_ip_list`` requires executing terraform apply twice for more than 2 gateways


## 2.23.2 (September 30, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.8.1342**
- Supported Terraform version: **v1.x**

### Features:
#### Site2Cloud:
1. Implemented support for Certificate-based HA Gateway Remote Identifier for Site2Cloud VPN:
   - New attributes in **aviatrix_site2cloud**:
     - ``backup_remote_identifier``

### Bug Fixes:
1. Fixed issue where route-based Single IP HA tunnel S2C creation fails for **aviatrix_site2cloud**


## 2.23.1 (September 13, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.8.1311**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented support for exposing BGP LAN interface info on transit in Azure via ``bgp_lan_ip_list`` and ``ha_bgp_lan_ip_list``
2. Implemented support for multiple disjoint port ranges for **aviatrix_microseg_policy_list**

### Enhancements:
1. Added support for updating ``bgp_md5_key`` and ``backup_bgp_md5_key`` for **aviatrix_spoke_external_device_conn** and **aviatrix_transit_external_device_conn**
2. Optimized the read operation for **aviatrix_transit_firenet_policy**

### Bug Fixes:
1. Fixed issue where NAT config has ordering issues that would cause refresh problems for **aviatrix_gateway_dnat** and **aviatrix_gateway_snat**
2. Fixed issue where provider crashes for private mode config
3. Fixed issue where exported SNAT/DNAT interface shows tunnel ID when using Transit connection or route-based S2C
4. Fixed issue where creating FQDN gateway with ``fqdn_lan_interface`` causes replacement


## 2.23.0 (August 09, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.8**
- Supported Terraform version: **v1.x**

### Features:
#### Provider:
1. Implemented support to ignore changes in selected tag keys across all resources on the provider-level:
  - New configuration block ``ignore_tags {}`` with the following options:
    - ``keys``
    - ``key_prefixes``

#### Multi-Cloud Transit:
1. Implemented support for Private Mode:
  - New attributes in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**:
    - ``private_mode_lb_vpc_id``
    - ``private_mode_subnet_zone``
    - ``ha_private_mode_subnet_zone``
  - New attribute in **aviatrix_vpc**:
    - ``private_mode_subnets``
  - New resources:
    - **aviatrix_controller_private_mode_config**
    - **aviatrix_private_mode_lb**
    - **aviatrix_private_mode_multicloud_endpoint**
2. Implemented a new resource to support Edge as a Spoke:
  - **aviatrix_edge_spoke**
3. Implemented a new resource to support attaching Edge as a Spoke to Transit Gateway:
  - **aviatrix_edge_spoke_transit_attachment**
4. Implemented a new resource to support Edge as a Spoke External Device Connection:
  - **aviatrix_edge_spoke_external_device_conn**
5. Implemented support for connection based AS path prepend for BGP Spoke Transit attachment in **aviatrix_spoke_transit_attachment** with the following new attributes:
  - ``spoke_prepend_as_path``
  - ``transit_prepend_as_path``
6. Implemented support for creating multiple BGP over LAN interfaces in **aviatrix_transit_gateway** for Azure with the following new attribute:
  - ``bgp_lan_interfaces_count``

#### Security:
1. Implemented support for order of rules and rule addition to any place in **aviatrix_firewall_policy** with the following new attribute:
  - ``position``

#### Settings:
1. Implemented a new resource to support CoPilot Security Group Management:
  - **aviatrix_copilot_security_group_management_config**

#### Site2Cloud:
1. Implemented support for Certificate based Authentication for Site2Cloud VPN:
  - New attributes in **aviatrix_site2cloud**:
    - ``auth_type``
    - ``ca_cert_tag_name``
    - ``remote_identifier``
  - New resource:
    - **aviatrix_site2cloud_ca_cert_tag**

#### TGW Orchestrator:
1. Implemented support for setting AWS TGW inspection mode in **aviatrix_aws_tgw** with the following new attribute:
  - ``inspection_mode``

### Enhancements:
1. Increased maximum number of policies allowed for **aviatrix_dnat** and **aviatrix_snat**
2. Removed ``fail_close_enabled`` from **aviatrix_firenet**. ``fail_close_enabled`` will automatically be true for all **aviatrix_firenet** for R2.23.0+
3. Updated ``account_name`` to ForceNew in **aviatrix_account**
4. Added support for ``insane_mode`` for **aviatrix_gateway**, **aviatrix_spoke_gateway**, and **aviatrix_transit_gateway** for Azure China

### Bug Fixes:
1. Fixed issue where duplicate **aviatrix_account** resources would be set into state even after giving an error
2. Fixed issue where ``username`` could not be specified with ``private_key_file`` in **aviatrix_firenet_vendor_integration**
3. Fixed issue where setting ``custom_algorithms`` to true would still use default values, causing tunnel replacement in **aviatrix_transit_external_device_conn**


## 2.22.5 (March 14, 2023)
### Notes:
- Supported Controller version: **UserConnect-6.7.1574**
- Supported Terraform version: **v1.x**

### Features:
#### CloudN:
1. Restored support of CloudN transit attachment CIDR approval features in **aviatrix_cloudn_transit_gateway_attachment**:
    - ``enable_dead_peer_detection``
    - ``enable_learned_cidrs_approval``
    - ``approved_cidrs``


## 2.22.4 (September 20, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.7.1480**
- Supported Terraform version: **v1.x**

### Enhancements:
1. Optimized the read operation for **aviatrix_transit_firenet_policy**

### Bug Fixes:
1. Fixed issue where NAT config has ordering issues that would cause refresh problems for **aviatrix_gateway_dnat** and **aviatrix_gateway_snat**


## 2.22.3 (August 02, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.7.1376**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented support for creating one HPE tunnel per instance size in **aviatrix_spoke_transit_attachment** and **aviatrix_transit_gateway_peering**:
  - ``enable_max_performance``


## 2.22.2 (July 20, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.7.1324**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented support for enabling/disabling Jumbo Frames on GRE tunnels under BGP connections in **aviatrix_transit_external_device_conn**:
  - ``enable_jumbo_frame``

### Enhancements:
1. Added duplicate rules check and removed deprecation message for ``domain_names`` in **aviatrix_fqdn** to continue support in-line tag rules and the standalone **aviatrix_fqdn_tag_rule** resource
2. Added duplicate rules check and removed deprecation message for ``policy`` in **aviatrix_firewall** to continue support in-line policy rules and the standalone **aviatrix_firewall_policy** resource

### Bug Fixes:
1. Fixed issue where adding more custom SNAT policy rules to ``snat_policy`` after creation on policy-based S2C fails
2. Fixed issue where editing FQDN default policy from allow-all to deny-all errors out
3. Fixed issue where importing invalid ID crashes plugin for **aviatrix_firewall_policy**


## 2.22.1 (June 10, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.7.1319**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented support for enabling preserve AS path when advertising manual summary CIDRs in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**:
  - ``enable_preserve_as_path``
2. Implemented new resources to support micro-segmentation:
  - **aviatrix_app_domain**
  - **aviatrix_microseg_policy_list**

#### Settings:
1. Implemented a new resource to support setting email configs for critical alerts and security events:
  - **aviatrix_controller_email_config**

### Enhancements:
1. Added support for "ANY" protocol for micro-segmentation policies in **aviatrix_microseg_policy_list**

### Bug Fixes:
1. Fixed issue where Terraform tries to disable the certificates when uploading renewed certificates
2. Fixed issue where destroying app domains created with Terraform errors out


## 2.22.0 (May 09, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.7**
- Supported Terraform version: **v1.x**

### Features:
#### Gateway:
1. Implemented support for setting rx queue size in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** with the following new attribute:
  - ``rx_queue_size``

#### Multi-Cloud Transit:
1. Implemented support for modifying BGP connection's MD5 signature in **aviatrix_spoke_external_device_conn** and **aviatrix_transit_external_device_conn**:
  - ``bgp_md5_key``
  - ``backup_bgp_md5_key``

#### CloudN:
1. Implemented a new resource to support Edge as a CaaG:
  - **aviatrix_edge_caag**
2. Implemented a new data source to get the list of device WAN interfaces:
  - **aviatrix_device_interfaces**

### Enhancements:
1. Implemented new resources to support the renaming from security domain to network domain. Resources and attributes whose name include "security_domain" will be deprecated in future releases.
Please follow the guide [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/migrating_from_security_domain_to_network_domain) for migration:
  - **aviatrix_aws_tgw_network_domain**
  - **aviatrix_segmentation_network_domain**
  - **aviatrix_segmentation_network_domain_association**
  - **aviatrix_segmentation_network_domain_connection_policy**
2. Renamed the attribute ``security_domain_name`` to ``network_domain_name`` in resources **aviatrix_aws_tgw_connect**, **aviatrix_aws_tgw_directconnect** and **aviatrix_aws_tgw_vpc_attachment**
to support the renaming from security domain to network domain. Resources and attributes whose name includes security domain will be deprecated in future releases.
Please follow the guide [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/migrating_from_security_domain_to_network_domain) for migration
3. Updated the ``vpc_id`` attribute for **aviatrix_gateway**, **aviatrix_spoke_gateway**, **aviatrix_transit_gateway** and **aviatrix_vpc** created in GCP to include the project id:
  - New format: ``"<vpc_name>~-~<project_id>"``
4. Added support for ``insane_mode`` for **aviatrix_gateway**, **aviatrix_spoke_gateway**, and **aviatrix_transit_gateway** created in AWS China
5. Sorted the lists of ``firewall_image_version`` and ``firewall_size`` in data source **aviatrix_firewall_instance_images**

### Bug Fixes:
1. Fixed issue where the forced replacement of the resource **aviatrix_cloudn_registration** errors out
2. Fixed issue where the creation of the resource **aviatrix_aws_tgw_vpc_attachment** errors out
3. Fixed issue where ``interface`` attribute in **aviatrix_snat** and **aviatrix_dnat** ``policy`` could not be set when using policy-based connections
4. Fixed issue with **aviatrix_transit_gateway_peering** creation when using gateways that do not exist

### Preview Features:
1. Implemented new resources to support micro-segmentation:
  - **aviatrix_app_domain**
  - **aviatrix_microseg_policy_list**

### Deprecations:
1. Deprecated support for CloudWAN. The following resources are removed:
  - **aviatrix_device_registration**
  - **aviatrix_device_tag**
  - **aviatrix_device_transit_gateway_attachment**
  - **aviatrix_device_aws_tgw_attachment**
  - **aviatrix_device_virtual_wan_attachment**
2. Removed support for the following attributes from the resource **aviatrix_cloudn_transit_gateway_attachment**:
  - ``enable_dead_peer_detection``
  - ``enable_learned_cidrs_approval``
  - ``approved_cidrs``


## 2.21.2 (March 31, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.6.5544**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented support for S2C RX steering toggle with a new attribute for **aviatrix_transit_gateway**:
  - ``enable_s2c_rx_balancing``

### Enhancements:
1. Updated the ``vpc_id`` attribute for **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** created in OCI to use the VCN OCID
2. Implemented support for uploading certificate content with the following new attributes in resource **aviatrix_controller_config**:
  - ``ca_certificate_file``
  - ``server_public_certificate_file``
  - ``server_private_key_file``

### Bug Fixes:
1. Fixed issue where the ``peering_ha_zone`` attribute in **aviatrix_gateway** would not be set to the correct value


## 2.21.1 (February 28, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.6.5404**
- Supported Terraform version: **v1.x**

### Features:
#### Multi-Cloud Transit:
1. Implemented a new resource and a new data source to support the Azure subnet inspection feature:
  - new resource: **aviatrix_spoke_gateway_subnet_group**
  - new data source: **aviatrix_spoke_gateway_inspection_subnets**
2. Implemented support for Active-Standby behavior backward compatibility with a new attribute for **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**:
  - ``enable_active_standby_preemptive``
3. Implemented support for disabling route propagation on BGP Spoke to attached Transit Gateway with a new attribute for **aviatrix_spoke_gateway**:
  - ``disable_route_propagation``
4. Implemented support for BGP MD5 Authentication with the following new attributes in **aviatrix_spoke_external_device_conn** and **aviatrix_transit_external_device_conn**:
  - ``bgp_md5_key``
  - ``backup_bgp_md5_key``
5. Renamed RBAC CloudWAN "all_cloudwan_write" to "all_cloudn_write" for ``permission_name`` in **aviatrix_rbac_group_permission_attachment**

### Enhancements:
1. Implemented a new data source to output Firewall Instance Images information:
  - **aviatrix_firewall_instance_images**
2. Updated attributes in data sources for **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
3. Made asynchronous calls to the API with constant polling for updates in order to prevent timeouts in those requests for some long-running HTTP requests
4. Added support for "NULL-ENCR" for ``phase_2_encryption`` in **aviatrix_transit_external_device_conn**
5. Extended GCM encryption in IPSec for **aviatrix_site2cloud**, **aviatrix_spoke_external_device_conn** and **aviatrix_transit_external_device_conn**:
  - Added support for "AES-128-GCM-64", "AES-128-GCM-96", "AES-128-GCM-128", "AES-256-GCM-64", "AES-256-GCM-96" and "AES-256-GCM-128" in ``phase_1_encryption``
  - Added support for "AES-256-GCM-64", "AES-256-GCM-96" and "AES-256-GCM-128" in ``phase_2_encryption``

### Bug Fixes:
1. Fixed issue where importing a resource with "symmetric" ID causes force replacement in **aviatrix_aws_tgw_peering** and **aviatrix_aws_tgw_peering_domain_conn**
2. Fixed issue where setting ``enable_public_subnet_filtering`` attribute in **aviatrix_gateway** would prevent ``tags`` from being set during creation
3. Fixed issue where ``terraform plan`` shows diff when creating a GCP transit with LAN interface without HA
4. Fixed issue where Aviatrix Terraform provider fails to upgrade controller from a version lower than latest, when target_version is set to "latest"


## 2.21.0 (January 23, 2022)
### Notes:
- Supported Controller version: **UserConnect-6.6**
- Supported Terraform version: **v1.x**

### Features:
#### Provider:
1. Implemented support for SSL certificate verification with the following new attributes in provider:
  - ``verify_ssl_certificate``
  - ``path_to_ca_certificate``

#### Gateway:
1. Implemented support to enable the feature to apply route entries into cloud platform routing table when using source NAT by adding the following attribute for **aviatrix_gateway_snat**:
  - ``apply_route_entry``

#### Multi-Cloud Transit:
1. Implemented a new resource to support registering a managed CloudN device to the controller:
  - **aviatrix_cloudn_registration**
2. Implemented a new resource to support connecting a standalone CloudN device to an **aviatrix_transit_gateway**:
  - **aviatrix_cloudn_transit_conn**
3. Implemented support for AWSChina in **aviatrix_firewall_instance**
4. Implemented support for BGP Prepending AS-PATH with the following new attribute for **aviatrix_transit_gateway_attachment**:
  - ``prepend_as_path``
5. Implemented support for BGP over LAN for GCP:
  - New attributes in **aviatrix_transit_gateway**:
    - ``bgp_lan_interfaces``
    - ``ha_bgp_lan_interfaces``
    - ``bgp_lan_ip_list``
    - ``ha_bgp_lan_ip_list``
  - New attribute in **aviatrix_transit_external_device_conn**
    - ``enable_bgp_lan_activemesh``
6. Implemented support for BGP over LAN on Spoke:
  - New attributes in **aviatrix_spoke_gateway**
    - ``enable_bgp``
    - ``spoke_bgp_manual_advertise_cidrs``
    - ``bgp_ecmp``
    - ``enable_active_standby``
    - ``prepend_as_path``
    - ``bgp_polling_time``
    - ``bgp_hold_time``
    - ``enable_learned_cidrs_approval``
    - ``learned_cidrs_approval_mode``
    - ``approved_learned_cidrs``
    - ``local_as_number``
  - New resource
    - **aviatrix_spoke_external_device_conn**
7. Implemented support for updating approved learned CIDRs with the following new attribute for **aviatrix_transit_gateway** :
  - ``approved_learned_cidrs``
8. Implemented support for BGP over LAN for GCP in **aviatrix_transit_external_device_conn**

### Enhancements:
1. Added support for updating ``remote_subnet`` in **aviatrix_transit_external_device_conn**
2. Updated ``key_name`` to a sensitive attribute in **aviatrix_firewall_instance**
3. Added retry when creating the following resources fails due to HA Transit is not up:
  - **aviatrix_transit_external_device_conn**
  - **aviatrix_vgw_conn**
4. Added support for scaling up to 64 netmap CIDRs in **aviatrix_site2cloud**

### Bug Fixes:
1. Fixed issue where ``bgp_manual_spoke_advertise_cidrs`` attribute in **aviatrix_transit_gateway** would have incorrect values when using **aviatrix_gateway_snat**
2. Removed the default value for ``interface`` attribute in **aviatrix_gateway_snat**
3. Fixed issue where the spaces in ``remote_subnet`` cause force replacement in **aviatrix_transit_external_device_conn**
4. Fixed issue where ``phase1_remote_identifier`` is set to two IP addresses when ``remote_gateway_ip`` and ``backup_remote_gateway_ip`` are with the same value
5. Fixed issue where ``active_active_ha`` causes diff when ActiveActive HA is enabled by default in some cases in **aviatrix_site2cloud**
6. Fixed issue where Terraform scripts with empty content is exported for **aviatrix_controller_cert_domain_config**, **aviatrix_controller_email_exception_notification_config** and **aviatrix_splunk_logging**
7. Fixed issue where an EOF error is returned when deleting transit HA gateway
8. Fixed issue where a service unavailable error may return when upgrading controller
9. Fixed issue where deleting HA with insane mode enabled returns error in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**

### Deprecations:
1. Removed support for ``storage_name`` attribute from **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** in AzureChina
2. Removed support for Non-ActiveMesh features from **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**


## 2.20.3 (November 22, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.5.2721**
- Supported Terraform version: **v1.x**

#### Multi-Cloud Transit
1. Implemented support for Transit FireNet for AWSChina

### Bug Fixes:
1. Fixed issue where Terraform Plan shows diff for a use case of **aviatrix_transit_external_device_conn** when controller is upgraded from 6.5.c- to 6.5.c+


## 2.20.2 (November 16, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.5.2608**
- Supported Terraform version: **v1.x**

### Bug Fixes:
1. Fixed issue where upgrading Controller using **aviatrix_controller_config** fails due to async action


## 2.20.1 (October 28, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.5.2608**
- Supported Terraform version: **v1.x**

### Features:
#### Firewall Network:
1. Implemented support for AzureGov cloud in **aviatrix_firewall_instance**

### Enhancements:
1. Added more validity checks for advanced option attributes in **aviatrix_transit_gateway_peering**
2. Added new standalone resource **aviatrix_controller_security_group_management_config** to configure Controller's Security Group Management settings

### Bug Fixes:
1. Fixed issue where ``phase1_remote_identifier`` would always be unset when two IP addressed are used for ``remote_gateway_ip`` in **aviatrix_transit_external_device_conn**
2. Fixed issue where OCI cloud **aviatrix_firewall_instance**s couldn't be launched with CheckPoint images
3. Fixed issue where refreshing **aviatrix_cloudn_transit_gateway_attachment** state would fail if attachment is deleted from UI
4. Fixed issue where refreshing **aviatrix_vgw_conn** state would fail it connection is deleted from UI

### Deprecations:
1. Deprecated ``enable_active_mesh`` in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
  - Non-ActiveMesh features will be removed in Aviatrix provider v2.21.0. Please follow the guide [here](https://registry.terraform.io/providers/AviatrixSystems/aviatrix/latest/docs/guides/migrating_to_active_mesh_transit_network) to migrate from Classic Aviatrix Encrypted Transit Network to Aviatrix ActiveMesh Transit Network
2. Deprecated ``sg_management_account_name`` and ``security_group_management`` in **aviatrix_controller_config**
  - Please remove the attributes from this resource, perform a refresh, and use the new **aviatrix_controller_security_group_management_config** resource to configure the Controller's Security Group Management settings


## 2.20.0 (August 17, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.5**
- Supported Terraform version: **v1.x**

### Features:
#### Gateway:
1. Implemented support for Selective Gateway Upgrade in **aviatrix_gateway** with the following new attributes:
  - ``software_version``
  - ``peering_ha_software_version``
  - ``image_version``
  - ``peering_ha_image_version``
2. Implemented new data source **aviatrix_gateway_image**
3. Implemented support for preallocated IP for Azure in **aviatrix_gateway** with the following attributes:
  - ``eip``
  - ``peering_ha_eip``
  - ``azure_eip_name_resource_group``
  - ``peering_ha_azure_eip_name_resource_group``
4. Implemented support for preallocated IP for OCI in **aviatrix_gateway** by updating the following attributes:
  - ``eip``
  - ``peering_ha_eip``

#### Multi-Cloud Transit:
1. Implemented support for Selective Gateway Upgrade in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** with the following new attributes:
  - ``software_version``
  - ``ha_software_version``
  - ``image_version``
  - ``ha_image_version``
2. Implemented support for preallocated IP for Azure in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** with the following attributes:
  - ``eip``
  - ``ha_eip``
  - ``azure_eip_name_resource_group``
  - ``ha_azure_eip_name_resource_group``
3. Implemented support for preallocated IP for OCI in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** by updating the following attributes:
  - ``eip``
  - ``ha_eip``
4. Updated the format for ``remote_vpc_name`` in **aviatrix_transit_external_device_conn** for BGP over LAN connections to "<vnet_name>:<vnet_resource_group>:<subscription_id>"

#### CloudWAN:
1. Implemented support for Selective Gateway Upgrade in **aviatrix_device_registration** when used for CloudN as a Gateway with the following new attributes:
  - ``software_version``
  - ``is_caag``

#### Useful Tools:
1. Implemented cross-subscription support for **aviatrix_vpc** for Azure by updating ``vpc_id`` to the new following 3-tuple format: "<vnet-name>:<resource-group-name>:<GUID>"

#### Settings:
1. Implemented support for Selective Gateway Upgrade in **aviatrix_controller_config** with the following new attributes:
  - ``manage_gateway_upgrades``
  - ``current_version``
  - ``previous_version``

### Enhancements:
1. Improved refresh performance of **aviatrix_firenet_firewall_manager** resource and data source
2. Added ``vpn_tunnel_data`` in **aviatrix_aws_tgw_vpn_conn** resource
3. Added ``private_key_file`` in **aviatrix_firenet_vendor_integration** data source to allow the user to use private key file instead of username/password for Check Point Cloud Guard

### Bug Fixes:
1. Fixed issue in **aviatrix_firenet** where creating with ``keep_alive_via_lan_interface_enabled`` set to false would still set ``keep_alive_via_lan_interface_enabled`` to true
2. Fixed issue where HA related attribute would be left in the state file after disabling HA on an **aviatrix_gatetway**, **aviatrix_spoke_gateway** or **aviatrix_transit_gateway**


## 2.19.5 (July 14, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.4.2776**
- Supported Terraform version: **v1.x**

### Features:
#### Accounts:
1. Implemented support for AWSTS in **aviatrix_account** and data source with the following new attributes:
  - ``awsts_account_number``
  - ``awsts_cap_url``
  - ``awsts_cap_agency``
  - ``awsts_cap_mission``
  - ``awsts_cap_role_name``
  - ``awsts_cap_cert``
  - ``awsts_cap_cert_key``
  - ``awsts_ca_chain_cert``
2. Implemented support for AWSS in **aviatrix_account** and data source with the following new attributes:
  - ``awss_account_number``
  - ``awss_cap_url``
  - ``awss_cap_agency``
  - ``awss_cap_account_name``
  - ``awss_cap_role_name``
  - ``awss_cap_cert``
  - ``awss_cap_cert_key``
  - ``awss_ca_chain_cert``

#### Firewall Network:
1. Implemented support for Fail Close and Network List Excluded From East-West Inspection in **aviatrix_firenet**

#### Gateway:
1. Implemented support for applying route entry in **aviatrix_gateway_dnat**
2. Implemented support for AWS Top Secret cloud in **aviatrix_gateway**
3. Implemented support for AWS Secret cloud in **aviatrix_gateway**
4. Implemented support for configuring gateway keepalive settings
  - **aviatrix_controller_gateway_keepalive_config**

#### Multi-Cloud Transit:
1. Implemented support for AWS Top Secret cloud  in **aviatrix_spoke_gateway**
2. Implemented support for AWS Secret cloud in **aviatrix_transit_gateway**
3. Implemented support for connection based BGP prepending in **aviatrix_transit_external_device_conn** and **aviatrix_vgw_conn**

#### TGW Orchestrator:
1. Implemented support for the following attribute in **aviatrix_aws_tgw_vpn_conn**
  - ``enable_global_acceleration``

### Enhancements:
1. Allowed the value "aviatrix" for the attribute ``host_os`` to support managed cloudN deployment
2. Added support for computed attribute``peering_ha_security_group_id`` in **aviatrix_gateway**
3. Added support for computed attributes ``availability_domains`` and ``fault_domains`` in **aviatrix_vpc** and data source
4. Added support for Panorama setup in **aviatrix_firenet_firewall_manager** data source

### Bug Fixes:
1. Fixed issue where creating, updating or deleting **aviatrix_controller_cert_domain_config** may cause timeout
2. Fixed issue where disabling Egress fails when Egress is enabled without setting Egress Static CIDRs in **aviatrix_firenet**
3. Fixed issue where setting "account_name" will cause panic in **aviatrix_rbac_group_access_account_attachment**
4. Fixed issue where context deadline exceeded error happens in the following resources
  - **aviatrix_account**
  - **aviatrix_aws_tgw_connect**
  - **aviatrix_aws_tgw_connect_peer**
  - **aviatrix_aws_tgw_intra_domain_inspection**
  - **aviatrix_aws_tgw_security_domain**
  - **aviatrix_aws_tgw_security_domain_connection**
  - **aviatrix_cloudn_transit_gateway_attachment**
  - **aviatrix_controller_bgp_max_as_limit_config**
  - **aviatrix_controller_cert_domain_config**
  - **aviatrix_controller_email_exception_notification_config**
  - **aviatrix_copilot_association**
  - **aviatrix_gateway_certificate_config**
5. Fixed issue where ``local_subnet_cidr`` can't be updated for a mapped connection in **aviatrix_site2cloud**
6. Fixed issue where updating access account to swap custom IAM roles for gateways fails
7. Fixed issue where updating ``single_az_ha`` does not apply to HA gateway in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
8. Fixed issue where enabling EBS volume encryption after initial gateway deployment only applies to primary gateway in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**


## 2.19.4 (June 24, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.4.2672**
- Supported Terraform version: **v0.12.x**, **v0.13.x**, **v0.14.x** and **v0.15.x**

### Enhancements:
1. Added retries for failed GET requests
2. Optimized state refresh performance for **aviatrix_transit_gateway_peering**
3. Updated Aviatrix HTTP Client to try to look for proxies in the default env variables HTTP_PROXY/http_proxy and HTTPS_PROXY/https_proxy


## 2.19.3 (June 14, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.4.2672**
- Supported Terraform version: **v0.12.x**, **v0.13.x**, **v0.14.x** and **v0.15.x**

### Features:
#### Firewall Network:
1. Implemented support for the following attributes for OCI in **aviatrix_firewall_instance**:
  - ``availability_domain``
  - ``fault_domain``

#### Gateway:
1. Implemented support for the following attributes for OCI in **aviatrix_gateway** and data source:
  - ``availability_domain``
  - ``fault_domain``
  - ``peering_ha_availability_domain``
  - ``peering_ha_fault_domain``

#### Multi-Cloud Transit:
1. Implemented support for the following attributes for OCI in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** and data source:
  - ``availability_domain``
  - ``fault_domain``
  - ``ha_availability_domain``
  - ``ha_fault_domain``


## 2.19.2 (June 11, 2021)
### Notes:
- Due to technical issues, 2.19.2 was not released correctly. Please use 2.19.3 instead.


## 2.19.1 (May 18, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.4.2561**
- Supported Terraform version: **v0.12.x**, **v0.13.x** and **v0.14.x**

### Features:
#### Accounts:
1. Implemented support for AWSGov IAM role-based in **aviatrix_account** and data source with the following new attributes:
  - ``awsgov_iam``
  - ``awsgov_role_app``
  - ``awsgov_role_ec2``
2. Implemented support for separate IAM role and policy for gateways in AWSChina and AWSGov **aviatrix_account** and data source

### Bug Fixes:
1. Fixed crashing issue when creating an **aviatrix_transit_external_device_conn** without ``phase1_remote_identifier``
2. Fixed an issue where enabling Single IP HA failover for an **aviatrix_site2cloud** with mapped config will read deltas in the ``phase1_remote_identifier`` values


## 2.19.0 (May 09, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.4**
- Supported Terraform version: **v0.12.x**, **v0.13.x** and **v0.14.x**

### Features:
#### Accounts:
1. Implemented support for Alibaba Cloud in **aviatrix_account** and data source
2. Implemented support for AzureChina, AzureGov and AWSChina clouds in **aviatrix_account**
3. Implemented support for separate IAM role and policy for gateways in AWS **aviatrix_account** with new attributes
  - ``aws_gateway_role_app``
  - ``aws_gateway_role_ec2``
4. Implemented support for enabling auditing in **aviatrix_account**:
  - New attribute ``audit_account``

#### CloudWAN:
1. Implemented support for enabling event triggered HA for Site2Cloud type connection resources:
  - New attribute ``enable_event_triggered_ha`` in **aviatrix_device_transit_gateway_attachment**

#### Firewall Network:
1. Implemented support for GCP FireNet with Fortinet and CheckPoint firewall vendors
2. Implemented support for TGW segmentation for Egress in TGW FireNet workflows:
  - New attribute ``tgw_segmentation_for_egress_enabled`` in **aviatrix_firenet** and data source
3. Implemented support for OCI FireNet
4. Implemented support for Egress FireNet route injection:
  - New attribute ``egress_static_cidrs`` in **aviatrix_firenet** and data source
5. Implemented custom AMI support for Firewall instance, allowing customers to launch special images provided by firewall vendors:
  - New attribute ``firewall_image_id`` in **aviatrix_firewall_instance**

#### Gateway:
1. Implemented support for Alibaba Cloud in **aviatrix_gateway** and data source
2. Implemented support for AzureGov, AWSChina and AzureChina clouds in **aviatrix_gateway**
3. Implemented support for IPSec tunnel down detection time in **aviatrix_gateway**:
  - New attribute ``tunnel_detection_time``

#### Multi-Cloud Transit:
1. Implemented support for the following attributes in **aviatrix_spoke_gateway**:
  - ``enable_private_vpc_default_route``
  - ``enable_skip_public_route_table_update``
  - ``enable_auto_advertise_s2c_cidrs``
2. Implemented support for enabling Event Triggered HA for Site2Cloud type connection resources:
  - New attribute ``enable_event_triggered_ha`` in **aviatrix_transit_external_device_conn**, **aviatrix_vgw_conn**
3. Implemented Insane Mode support over Public Network for Transit Peering in **aviatrix_transit_gateway_peering**:
  - ``enable_insane_mode_encryption_over_internet``
  - ``tunnel_count``
4. Implemented support for attaching a managed CloudN device to an **aviatrix_transit_gateway**:
  - New resource **aviatrix_cloudn_transit_gateway_attachment**
5. Implemented support for setting approved CIDRs in **aviatrix_transit_external_device_conn**:
  - New attribute ``approved_cidrs``
6. Implemented support for Multi-Tier Transit feature:
  - New attribute ``enable_multi_tier_transit`` in **aviatrix_transit_gateway** and data source
7. Implemented support for Alibaba Cloud in **aviatrix_transit_gateway** and **aviatrix_spoke_gateway** and data sources
8. Implemented support for AzureGov, AWSChina and AzureChina clouds in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
9. Implemented support for IPSec tunnel down detection time in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**:
  - New attribute ``tunnel_detection_time``
10. Implemented OCI transit Insane Mode support in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
11. Implemented support for ``enable_egress_transit_firenet`` for Azure and OCI in **aviatrix_transit_gateway**
12. Implemented support for phase 1 remote identifier in **aviatrix_transit_external_device_conn**:
  - New attribute ``phase1_remote_identifier``

#### Settings:
1. Implemented support for associating a Controller with a CoPilot instance, allowing user login without a username and password:
  - New resource **aviatrix_copilot_association**
2. Implemented support for adding profile names to Remote Syslog configs:
  - New attribute ``name`` in **aviatrix_remote_syslog**
3. Implemented support for enabling/disabling Controller from sending exception emails to Aviatrix:
  - New resource **aviatrix_controller_email_exception_notification_config**
4. Implemented support for updating Controller's certificate domain, required for Aviatrix China Solution:
  - New resource **aviatrix_controller_cert_domain_config**
5. Implemented support for BGP max AS limit controller configuration:
  - New resource **aviatrix_controller_bgp_max_as_limit_config**

#### Site2Cloud:
1. Implemented support for enabling event triggered HA for Site2Cloud resource:
  - New attribute ``enable_event_triggered_ha`` in **aviatrix_site2cloud**
2. Implemented support for setting optional tunnel IP address with the following attributes in **aviatrix_site2cloud**:
  - ``local_tunnel_ip``
  - ``remote_tunnel_ip``
  - ``backup_local_tunnel_ip``
  - ``backup_remote_tunnel_ip``
3. Implemented single public IP failover support for **aviatrix_site2cloud** connections:
  - New attribute ``single_ip_ha``
4. Implemented support for phase 1 remote identifier for Site2Cloud:
  - New attribute ``phase1_remote_identifier`` in **aviatrix_site2cloud**

#### TGW Orchestrator:
1. Implemented new resources to decouple ``security_domains`` out of **aviatrix_aws_tgw**:
  - **aviatrix_aws_tgw_security_domain**
  - **aviatrix_aws_tgw_security_domain_connection**
2. Implemented support for TGW intra-domain inspection:
  - New resource **aviatrix_aws_tgw_intra_domain_inspection**

#### Useful Tools:
1. Implemented support for Alibaba Cloud in **aviatrix_vpc** and data source
2. Implemented support for AzureGov, AWSChina and AzureChina clouds in **aviatrix_vpc**
3. Implemented support for creating an **aviatrix_vpc** in Azure with an existing ``resource_group``

### Enhancements:
1. Added following attributes in **aviatrix_account** data source:
  - ``gcloud_project_id``
  - ``arm_subscription_id``
  - ``awsgov_account_number``
  - ``awsgov_access_key``
2. Changed ``aws_access_key`` and ``aws_gov_access_key`` in **aviatrix_acount** to be sensitive values
3. Optimized state refresh performance for **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
4. Added new map type attribute ``tags`` to replace ``tag_list`` in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
5. Added support for Fortinet Fortigate in **aviatrix_firenet_vendor_integration** data source
6. Added computed value ``tgw_id`` in **aviatrix_aws_tgw**

### Bug Fixes:
1. Fixed an edge case in **aviatrix_gateway** that could cause the provider to crash when refreshing the resource
2. Fixed **aviatrix_transit_gateway_peering** to allow setting duplicate AS Numbers in the ``prepend_as_path1`` and ``prepend_as_path2`` attributes
3. Fixed **aviatrix_fqdn** to not remove ``domain_names`` after importing the resource with ``manage_domain_names`` set to false
4. Fixed reordering issue for ``security_domains`` in **aviatrix_aws_tgw**
5. Fixed issue where Transit FireNet option and downsizing the gateway can't be completed in one Terraform operation
6. Fixed issue where enabling HA for Insane Mode **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** will cause Terraform to errors out
7. Fixed issue where disabling Transit FireNet and Egress Transit FireNet options can't be completed in one Terraform operation in **aviatrix_transit_gateway**

### Deprecations:
1. Deprecated the in-line attributes ``security_domains``, ``security_domain_name``, ``connected_domains``, ``aviatrix_firewall``, ``native_egress`` and ``native_firewall`` in **aviatrix_aws_tgw**. Please use the standalone resources **aviatrix_aws_tgw_security_domain** and **aviatrix_aws_tgw_security_domain_connection** instead
2. Deprecated ``tag_list`` in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**. Please use map type attribute ``tags`` instead


## 2.18.2 (March 22, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.3.2364**
- Supported Terraform version: **v0.12.x** and **v0.13.x**

### Bug Fixes:
1. Fixed an issue where associating an out-of-band firewall instance, not created by the specified controller, was not supported in **aviatrix_firewall_instance_association**


## 2.18.1 (March 18, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.3.2364**
- Supported Terraform version: **v0.12.x** and **v0.13.x**

### Features:
1. Implemented new resources to support TGW Connect and Connect peers:
  - **aviatrix_aws_tgw_connect**
  - **aviatrix_aws_tgw_connect_peer**
2. Implemented support for GCP FireNet:
  - New attributes ``lan_vpc_id`` and ``lan_private_subnet`` in **aviatrix_transit_gateway**
  - New attribute ``fqdn_lan_vpc_id`` in **aviatrix_gateway**
  - New attributes ``egress_vpc_id`` and ``management_vpc_id`` in **aviatrix_firewall_instance**
3. Implemented support for FireNet Keep Alive via Firewall LAN Interface:
  - New attribute ``keep_alive_via_lan_interface_enabled`` in **aviatrix_firenet** resource and data source
4. Implemented support for Gateway Certificate import:
  - New resource **aviatrix_gateway_certificate_config**
5. Implemented support for configuring AWS TGW CIDRs in **aviatrix_aws_tgw** using attribute ``cidrs``
6. Implemented support for IKEv2 for route-based Site2Cloud connections in **aviatrix_site2cloud**
7. Implemented support for ``metrics_only`` option in **aviatrix_datadog_agent**
8. Implemented support for building OOB Transit/Spoke gateway and HA in different AZs/Subnets
9. Implemented support for controller backup for AWSGov, Azure, GCP and OCI providers
10. Implemented support for attribute ``route_tables`` in **aviatrix_vpc** resource and data source
11. Implemented support for Management Access from on-prem in **aviatrix_site2cloud**
12. Implemented support for Enable Transit Summarize CIDR to TGW in **aviatrix_transit_gateway** using ``enable_transit_summarize_cidr_to_tgw``
13. Implemented support for Jumbo Frames in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** using ``enable_jumbo_frame``
14. Implemented support for Tags in **aviatrix_firewall_instance** using ``tags``

### Enhancements:
1. Added check function to ignore whitespace for following attributes in **aviatrix_transit_external_device_conn**:
  - ``local_tunnel_cidr``
  - ``remote_tunnel_cidr``
  - ``backup_local_tunnel_cidr``
  - ``backup_remote_tunnel_cidr``
2. Added support for DH-group 19, 20 and 21 when IKEv2 enabled in **aviatrix_transit_external_device_conn**
3. Added support for DH-group 20 and 21 when IKEv2 enabled in **aviatrix_site2cloud**
4. Updated following attributes to ForceNew in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**:
  - ``gw_name``
  - ``vpc_id``
  - ``vpc_reg``
  - ``subnet``
  - ``zone``
  - ``oob_management_subnet``
  - ``oob_availability_zone``
5. Updated following attributes to ForceNew in **aviatrix_aws_tgw**:
  - ``tgw_name``
  - ``aws_side_as_number``
6. Updated following attributes to ForceNew in **aviatrix_aws_tgw_vpc_attachment**:
  - ``tgw_name``
  - ``vpc_id``
7. Updated attribute ``local_as_number`` to Optional and Computed in **aviatrix_transit_gateway**:
8. Optimized API list_vpcs_summary to reduce terraform refresh time for **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**

### Bug Fixes:
1. Fixed an issue where **aviatrix_firewall_instance** would not import attribute ``key_name`` correctly
2. Fixed an issue where updating ``ha_subnet`` fails in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
3. Fixed an issue where terraform refresh/destroy does not work if site2cloud connection has been removed from UI for **aviatrix_vgw_conn**
4. Fixed an issue where upgrading controller causes CID to expire, which fails other functions in **aviatrix_controller_config**
5. Fixed an issue where dot is not supported in ``spoke_vpc_id`` in **aviatrix_azure_spoke_native_peering**
6. Fixed an issue where enabling encrypt volume with a customer managed keys fails in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**

### Deprecations:
1. Deprecated the in-line ``attached_vpc`` and ``attached_aviatrix_transit_gateway`` attributes in **aviatrix_aws_tgw**. Please use the standalone **aviatrix_aws_tgw_vpc_attachment** and **aviatrix_aws_tgw_transit_gateway_attachment** resources instead
2. Deprecated the in-line ``transit_gw`` attribute in **aviatrix_spoke_gateway**. Please use the standalone **aviatrix_spoke_transit_attachment** resource instead
3. Deprecated the in-line ``policy`` attribute in **aviatrix_firewall**. Please use the standalone **aviatrix_firewall_policy** resource instead
4. Deprecated the in-line ``domain_names`` attribute in **aviatrix_fqdn**. Please use the standalone **aviatrix_fqdn_tag_rule** resource instead
5. Deprecated the in-line ``firewall_instance_association`` attribute in **aviatrix_firenet**. Please use the standalone **aviatrix_firewall_instance_association** resource instead


## 2.18.0 (January 31, 2021)
### Notes:
- Supported Controller version: **UserConnect-6.3.2092**
- Supported Terraform version: **v0.12.x** and **v0.13.x**

### Features:
1. Implemented support for BGP over GRE and BGP over LAN through ``enable_bgp_over_lan`` in **aviatrix_transit_gateway**, and the following attributes in **aviatrix_transit_external_device_conn**:
  - ``tunnel_protocol``
  - ``remote_lan_ip``
  - ``backup_remote_lan_ip``
  - ``local_lan_ip``
  - ``backup_local_lan_ip``
  - ``remote_vpc_name``
2. Implemented support for the controller HTTPS certificate import with the following attributes in **aviatrix_controller_config**:
  - ``ca_certificate_file_path``
  - ``server_public_certificate_file_path``
  - ``server_private_key_file_path``
3. Implemented support for creating a Public Subnet Filtering gateway with the following attributes in **aviatrix_gateway**:
  - ``enable_public_subnet_filtering``
  - ``public_subnet_filtering_route_tables``
  - ``public_subnet_filtering_ha_route_tables``
  - ``public_subnet_filtering_guard_duty_enforced``
4. Implemented support for configuring AWS Guard Duty:
  - New resource **aviatrix_aws_guard_duty**
  - New attribute ``aws_guard_duty_scanning_interval`` in **aviatrix_controller_config**
5. Implemented support for configuring Learned CIDR Approval per connection:
  - New attribute ``learned_cidrs_approval_mode`` in **aviatrix_transit_gateway**
  - New attribute ``enable_learned_cidrs_approval`` in **aviatrix_device_transit_gateway_attachment**, **aviatrix_transit_external_device_conn** and **aviatrix_vgw_conn**
6. Implemented support for configuring Manual Advertised CIDRs per connection:
  - New attribute ``manual_bgp_advertised_cidrs`` in **aviatrix_device_transit_gateway_attachment**, **aviatrix_transit_external_device_conn** and **aviatrix_vgw_conn**
7. Implemented support for FireNet with AWS Gateway Load Balancer (GWLB):
  - New attribute ``enable_gateway_load_balancer`` in **aviatrix_transit_gateway**
  - New attribute ``enable_native_gwlb`` in **aviatrix_vpc**
  - Make ``firenet_gw_name`` Optional in **aviatrix_firewall_instance** and **aviatrix_firewall_instance_association**
8. Implemented support for Monitor Gateway Subnets feature in **aviatrix_transit_gateway** and **aviatrix_spoke_gateway** using the following attributes:
  - ``enable_monitor_gateway_subnets``
  - ``monitor_exclude_list``
9. Implemented support for private transit gateway peering with single-tunnel mode in **aviatrix_transit_gateway_peering** using attribute ``enable_single_tunnel_mode``
10. Implemented support for IKEv2 protocol in transit to external device connections in **aviatrix_transit_external_device_conn** using attribute ``enable_ikev2``
11. Implemented new resource to support transit in Azure with ExpressRoute:
  - **aviatrix_azure_vng_conn**
12. Implemented support for Private OOB feature:
  - New resource **aviatrix_controller_private_oob** to enable Controller-wide setting
  - New attributes ``enable_private_oob``, ``oob_management_subnet``, and ``oob_availability_zone`` in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
  - New attribute ``enable_private_oob_subnet`` in **aviatrix_vpc**
13. Implemented support for proxy configuration:
  - New resource: **aviatrix_proxy_config**
14. Implemented support for OCI in **aviatrix_vpc**
15. Implemented support for Aviatrix client/ovpn file download from the controller when SAML authentication is used:
  - New resource: **aviatrix_vpn_cert_download**
16. Implemented new resources to support Controller logging configurations:
  - **aviatrix_remote_syslog**
  - **aviatrix_splunk_logging**
  - **aviatrix_filebeat_forwarder**
  - **aviatrix_sumologic_forwarder**
  - **aviatrix_datadog_agent**
  - **aviatrix_netflow_agent**
  - **aviatrix_cloudwatch_agent**

### Enhancements:
1. Added Computed value ``ha_lan_interface_cidr`` in **aviatrix_transit_gateway**
2. Changed **aviatrix_gateway** attribute ``monitor_exclude_list`` type from String to Set of Strings
3. Added support of ``tag_list`` for Azure provider in **aviatrix_gateway**, **aviatrix_transit_gateway**, and **aviatrix_spoke_gateway** resources and data sources
4. Added ``customized_transit_vpc_routes`` in **aviatrix_transit_gateway** resource and data source
5. Added ``azure_vnet_resource_id`` as output for **aviatrix_vpc** resource and data source

### Bug Fixes:
1. Fixed issue where users could not create an **aviatrix_firewall_instance** if the VPC/VNET was not managed by the Aviatrix controller
2. Fixed an argument ordering issue in **aviatrix_site2cloud** Custom Mapped attributes by changing from type Set to List
3. Fixed race condition when deploying spoke gateway (HA disabled) using ``customized_spoke_vpc_routes`` and ``transit_gw``
4. Fixed issue where creating **aviatrix_site2cloud** for ActiveActive-enabled gateway causes deltas in state
5. Fixed issue where attribute ``bgp_manual_spoke_advertise_cidrs`` in **aviatrix_transit_gateway** causes delta in every apply
6. Fixed issue where Egress Transit Gateway can't be created due to blocking on the provider end
7. Fixed issue where an **aviatrix_spoke_gateway** with advertised spoke VPC CIDRs can't connect to an **aviatrix_transit_gateway**

### Deprecations:
1. Deprecated the in-line ``firewall_instance_association`` attribute in **aviatrix_firenet**. Please use the standalone **aviatrix_firewall_instance_association** resource instead


## 2.17.2 (December 08, 2020)
### Notes:
- Supported Controller version: **UserConnect-6.2.1914**
- Supported Terraform version: **v0.12.x** and **v0.13.x**

### Features:
1. Implemented further support for Custom Mapped and overlapping CIDR scenarios for **aviatrix_site2cloud** with attribute ``forward_traffic_to_transit``
2. Implemented Connection-based BGP Prepending AS-PATH support with the following attributes for **aviatrix_transit_gateway_peering**:
  - ``prepend_as_path1``
  - ``prepend_as_path2``

### Bug Fixes:
1. Fixed issue where the following parameters caused reordering issues for **aviatrix_transit_gateway_peering**:
  - ``gateway1_excluded_network_cidrs``
  - ``gateway2_excluded_network_cidrs``
  - ``gateway1_excluded_tgw_connections``
  - ``gateway2_excluded_tgw_connections``


## 2.17.1 (November 22, 2020)
### Notes:
- Supported Controller version: **UserConnect-6.2.1891**
- Supported Terraform version: **v0.12.x** and **v0.13.x**

### Features:
1. Implemented support for monitoring gateway subnets in **aviatrix_gateway** through ``enable_monitor_gateway_subnets`` and ``monitor_exclude_list``
2. Implemented support for managing Aviatrix VPN timeout configurations through ``idle_timeout`` and ``renegotiation_interval`` in **aviatrix_gateway**
3. Implemented support for ``enable_active_standby`` in **aviatrix_transit_gateway**
4. Implemented Active-Standby support for Transit Network workflows:
  - ``enable_active_standby`` in **aviatrix_transit_gateway**
  - ``switch_to_ha_standby_gateway`` in **aviatrix_transit_external_device_conn**
5. Implemented new resource to decouple ``firewall_instance_association`` out of ``aviatrix_firenet``:
  - **aviatrix_firewall_instance_association**
6. Implemented support for transit gateway peering over private networks through the ``enable_peering_over_private_network`` attribute in **aviatrix_transit_gateway_peering**
7. Implemented support for FQDN gateway in Azure FireNet:
  - ``fqdn_lan_cidr`` as an attribute, and ``fqdn_lan_interface`` as a computed output in **aviatrix_gateway**
  - ``lan_interface_cidr`` as an attribute in **aviatrix_transit_gateway**
8. Implemented support for ``local_login`` in **aviatrix_rbac_group**
9. Implemented Support for IDP Metadata URLs for SAML endpoints
10. Implemented support for ``sign_authn_requests`` in **aviatrix_saml_endpoint**
11. Implemented Bootstrap support for AWS and Azure FireNet solutions in aviatrix_firewall_instance:
  - ``bootstrap_storage_name``
  - ``storage_access_key``
  - ``file_share_folder``
  - ``share_directory``
  - ``sic_key``
  - ``user_data``
  - ``container_folder``
  - ``sas_url_config``
  - ``sas_url_license``
12. Implemented support for DH Group 19 in **aviatrix_site2cloud**
13. Implemented support for Custom Mapped in **aviatrix_site2cloud**

### Enhancements:
1. Changed ``management_subnet`` to optional to support Check Point and Fortinet instances in **aviatrix_firewall_instance**
2. Added support for Terraform state migration due to resource-decoupling implementation for the following resources:
  - **aviatrix_aws_tgw**
  - **aviatrix_firenet**
  - **aviatrix_fqdn**
  - **aviatrix_spoke_gateway**
  - **aviatrix_vpn_profile**
  - **aviatrix_vpn_user**
3. Official support for Terraform 0.13

### Bug Fixes:
1. Fixed issue with deltas in the state after creating non-AWS VPN gateways with ELB disabled and ``vpn_protocol`` set as "UDP"


## 2.17.0 (October 15, 2020)
### Notes:
- Supported Controller version: **UserConnect-6.2** (tested on **UserConnect-6.2.1700**)
- Supported Terraform version: **v0.12.x**

### Features:
1. Implemented new resources to support CloudWAN:
  - **aviatrix_device_aws_tgw_attachment**
  - **aviatrix_device_interface_config**
  - **aviatrix_device_registration**
  - **aviatrix_device_tag**
  - **aviatrix_device_transit_gateway_attachment**
  - **aviatrix_device_virtual_wan_attachment**
2. Implemented new resource to decouple ``domain_names`` out of ``aviatrix_fqdn``:
  - **aviatrix_fqdn_tag_rule**
3. Implemented new resource to decouple ``policy`` out of ``aviatrix_firewall``:
  - **aviatrix_firewall_policy**
4. Implemented new resources to support Multi-Cloud Segmentation:
  - **aviatrix_segmentation_security_domain**
  - **aviatrix_segmentation_security_domain_connection_policy**
  - **aviatrix_segmentation_security_domain_association**
5. Implemented support for updating **aviatrix_saml_endpoint**
6. Implemented support for advanced options to specify ``subnet_size`` and ``num_of_subnet_pairs`` for AWS, AWSGov, and Azure VPCs/VNets in **aviatrix_vpc** resource and data source
7. Implemented support for launching AWS TGWs with Multicast capability through the ``enable_multicast`` attribute for **aviatrix_aws_tgw** resource
8. Implemented Insane Mode support for GCP **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
9. Implemented support for updating ``customized_routes`` and ``customized_route_advertisement`` for **aviatrix_aws_tgw**
10. Implemented support for Availability Zone selection for the following resources in Azure:
  - ``zone`` and ``peering_ha_zone`` for **aviatrix_gateway**
  - ``zone`` and ``ha_zone`` for **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
  - ``zone`` for **aviatrix_firewall_instance**
11. Implemented new resource to decouple attaching **aviatrix_spoke_gateway** to **aviatrix_transit_gateway** out of **aviatrix_spoke_gateway**
  - **aviatrix_spoke_transit_attachment**

### Enhancements:
1. Blocked updating ``allocate_new_eip``, ``eip`` and ``ha_eip`` for **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
2. Added ``private_subnets`` and ``public_subnets`` as output for **aviatrix_vpc** resource and data source
3. Added support of ``resource_group`` for Azure provider in **aviatrix_vpc** data source

### Bug Fixes:
1. Fixed issue where there was a delta in state after creating a GCP **aviatrix_vpc**
2. Fixed import issue for **aviatrix_firewall**
3. Fixed issue where long metadata text was unable to be handled in **aviatrix_saml_endpoint** by updating operations from GET to POST method


## 2.16.3 (September 17, 2020)
### Notes:
- Supported Controller version: **UserConnect-6.1.1309**
- Supported Terraform version: **v0.12.x**

### Features:
1. Implemented support for allowing multiple **aviatrix_transit_gateway** attachments to **aviatrix_spoke_gateway**
2. Implemented support for Dual Transit FireNet through new attribute ``enable_egress_transit_firenet`` in **aviatrix_transit_gateway**
3. Implemented support for AWSGov cloud in the following resources:
  - **aviatrix_vpc**
  - **aviatrix_gateway**
  - **aviatrix_spoke_gateway**
  - **aviatrix_transit_gateway**
  - **aviatrix_aws_tgw**

### Enhancements:
1. Added validation function for ``username`` in **aviatrix_account_user** to block using upper letters in ``username`` since it is case insensitive in controller


## 2.16.2 (August 18, 2020)
### Notes:
- Supported Controller version: **UserConnect-6.1.1280**
- Supported Terraform version: **v0.12.x**

### Bug Fixes:
1. Fixed issue where peered TGWs with connected domain policies caused the **aviatrix_aws_tgw** to read deltas due to backend change


## 2.16.1 (August 07, 2020)
### Notes:
- Moved provider to HashiCorp Terraform Registry
- Supported Controller version: **UserConnect-6.1** (tested on **UserConnect-6.1.1162**)
- Supported Terraform version: **v0.12.x**


## 2.16.0 (August 04, 2020)
### Notes:
- Supported Controller version: **UserConnect-6.1** (tested on **UserConnect-6.1.1162**)
- Supported Terraform version: **v0.12.x**

### Features:
1. Implemented new resource to support periodic ping from gateways:
  - **aviatrix_periodic_ping**
2. Implemented new resource to support FQDN pass-through:
  - **aviatrix_fqdn_pass_through**
3. Implemented support for specifying and updating ``gateway1_excluded_network_cidrs``, ``gateway1_excluded_tgw_connections``, ``gateway2_excluded_network_cidrs``, and ``gateway2_excluded_tgw_connections`` for **aviatrix_transit_gateway_peering**
4. Implemented support for configuring ``bgp_polling_time``, ``prepend_as_path``, ``local_as_number``, and ``bgp_ecmp`` for **aviatrix_transit_gateway**
5. Implemented support for ``enable_vpc_dns_server`` in **aviatrix_controller_config**
6. Implemented support for updating name servers individually on ELBs under the **aviatrix_geo_vpn**
7. Implemented support for specifying EIPs to use for launching GCP **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** by setting ``allocate_new_eip`` to false and their respective ``eip`` and/or ``peering_ha_eip/ha_eip`` attributes
8. Implemented support for syncing **aviatrix_gateway_dnat** and **aviatrix_gateway_snat** policies to HA gateways through the ``sync_to_ha`` argument

### Enhancements:
1. Removed condition requiring ``single_az_ha`` to be disabled to in order to set ``enable_encrypt_volume`` for **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
2. Enhanced reading ``allocate_new_eip`` for GCP **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** data sources

### Bug Fixes:
1. Fixed issue where peered TGWs showing in domain connection list causes **aviatrix_aws_tgw_peering** to read deltas due to backend change


## 2.15.1 (July 10, 2020)
### Notes:
- Supported Controller version: **UserConnect-6.0** (tested on **UserConnect-6.0.2383**)
- Supported Terraform version: **v0.12.x**

### Enhancements:
1. Implemented support for 4-Byte ASN (Autonomous System Number) in **aviatrix_aws_tgw**, **aviatrix_aws_tgw_vpn_conn**, **aviatrix_transit_external_device_conn** and **aviatrix_vgw_conn**


## 2.15.0 (June 22, 2020)
### Notes:
- Supported Controller version: **UserConnect-6.0** (tested on **UserConnect-6.0.2269**)
- Supported Terraform version: **v0.12.x**

### Features:
1. New data sources:
  - **aviatrix_firewall**
  - **aviatrix_vpc_tracker**
2. Implemented support for the option to manage attachment on either **aviatrix_vpn_profile** or **aviatrix_vpn_user** using ``manage_user_attachment`` (and ``profiles`` for the user)
3. Implemented support for ``action`` under domain_names filters for **aviatrix_fqdn**
4. Implemented support for adding VPN users under GeoVPN workflow
5. Implemented support for specifying ``ha_peering_subnet`` for GCP **aviatrix_gateway**
6. Implemented support for specifying ``ha_subnet`` for GCP **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
7. Implemented support for ``enable_ikev2`` for **aviatrix_site2cloud**

### Enhancements:
1. Updated **aviatrix_site2cloud**'s``tunnel_type`` to support "policy" and "route"-based options in Controller 6.0
2. Added ``route_tables`` and ``route_tables_filter`` in **aviatrix_vpc** data source
3. Updated **aviatrix_vpc** to return parsed vpc_id for GCP VPC Networks
4. Updated terraform provider to support unencrypted gateway volumes as an option for backward compatibility between existing and new **aviatrix_gateways** created in Controller version 6.0. New gateway volumes are encrypted by default by the Controller in 6.0, but will not be, if created by Terraform unless otherwise specified by ``enable_encrypt_volume``
5. Enhanced GCP access account creation by supporting uploading credential files directly from local
6. Updated **aviatrix_gateway_snat** to support custom SNAT in cases of spoke to transit peering using ``connection``

### Bug Fixes:
1. Fix issue where **aviatrix_aws_tgw** could not be found in terraform state after creation due to backend change
2. Fix issue where HA gateways could not be created in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
3. Fix issue where **aviatrix_saml_endpoint**'s ``custom_saml_request_template`` return output was null after creation


## 2.14.1 (May 19, 2020)
### Notes:
- Supported Controller version: **UserConnect-5.4.1232**
- Supported Terraform version: **v0.12.x**

### Bug Fixes:
1. Fixed issue where **aviatrix_transit_external_device_conn** is forced to recreate due to ``connection_type`` not being set correctly


## 2.14.0 (May 08, 2020)
### Notes:
- Supported Controller version: **UserConnect-5.4.1201**
- Supported Terraform version: **v0.12.x**

### Features:
1. Implemented support for dynamically updating ``customized_route_advertisement`` in **aviatrix_aws_tgw_vpc_attachment**
2. Implemented support for SAML authentication for Controller login in **aviatrix_saml_endpoint**
3. New data source to support referencing specific private/public subnets:
  - **aviatrix_vpc**
4. New resources to support AWS TGW inter-region peering:
  - **aviatrix_aws_tgw_peering**
  - **aviatrix_aws_tgw_peering_domain_conn**
5. Implemented new resource to support connection to External Devices for Transit Network:
  - **aviatrix_transit_external_device_conn**

### Enhancements:
1. Added ``peering_ha_gw_name`` in **aviatrix_gateway**, and ``ha_gw_name`` in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** as computed values
2. Added ``peering_ha_private_ip`` in **aviatrix_gateway** data source, and ``ha_private_ip`` in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway** data sources as computed values

### Bug Fixes:
1. Fixed issue where OpenVPN configurations are unable to be modified when attached to a GeoVPN


## 2.13.0 (April 02, 2020)
### Notes:
- Supported Controller version: **UserConnect-5.4.1074**
- Supported Terraform version: **v0.12.x**

### Features:
1. Implemented support for ``enable_learned_cidrs_approval`` in **aviatrix_transit_gateway**, **aviatrix_aws_tgw_vpn_conn** and **aviatrix_aws_tgw_directconnect**
2. Implemented a new parameter ``manage_transit_gateway_attachment`` to provide the option of attaching/detaching transit gateway to/from AWS TGW outside of **aviatrix_aws_tgw** resource
3. New resources to support Role-Based Access Control (RBAC) feature in Controller 5.4 release:
  - **aviatrix_rbac_group**
  - **aviatrix_rbac_group_access_account_attachment**
  - **aviatrix_rbac_group_permission_attachment**
  - **aviatrix_rbac_group_user_attachment**
4. New resources:
  - **aviatrix_aws_tgw_transit_gateway_attachment**

### Enhancements:
1. Enhanced read-back of ``attached_aviatrix_transit_gateway`` to cover cases where multiple transit gateways are launched on the same VPC as the one already attached to the AWS TGW
2. Removed ``account_name`` from **aviatrix_account_user** for RBAC implementation

### Bug Fixes:
1. Fixed issue where changes in ``vpc_name`` in **aviatrix_aws_tgw** results in ``subnets`` being mismatched in the Terraform state


## 2.12.0 (March 12, 2020)
### Notes:
- Supported Controller version: **UserConnect-5.3.1491**
- Supported Terraform version: **v0.12.x**

### Features:
1. Implemented support for Transit FireNet:
  - ``enable_transit_firenet`` in **aviatrix_transit_gateway**
  - ``username`` and ``password`` in **aviatrix_firewall_instance** for Azure (Transit) FireNet
2. New resources for Transit FireNet:
  - **aviatrix_transit_firenet_policy**
  - **aviatrix_firewall_management_access**
3. New resources:
  - **aviatrix_azure_spoke_native_peering**
4. New resource **aviatrix_azure_peer** to replace **aviatrix_arm_peer**
5. Implemented support for Azure VNet in **aviatrix_vpc** resource

### Enhancements:
1. Enhanced handling enabling/disabling active-mesh and attaching/detaching to/from transit actions during updates in **aviatrix_spoke_gateway**
2. The following computed attributes are now available in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**:
  - ``private_ip``
  - ``instance_id``
  - ``security_group_id``
3. ``ha_cloud_instance_id`` is now a computed attribute available in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
4. Replaced ``cloudn_bkup_gateway_inst_id`` with ``peering_ha_cloud_instance_id`` in **aviatrix_gateway**
5. Deprecated **aviatrix_arm_peer** resource and replaced it with **aviatrix_azure_peer**

### Bug Fixes:
1. Fixed issue where **aviatrix_firewall_instance** forces replacement if ``firewall_image_version`` is not set
2. Fixed issue where **aviatrix_gateway_dnat** resource creation fails


## 2.11.0 (February 18, 2020)
### Notes:
- Supported Controller version: **UserConnect-5.3.1391**
- Supported Terraform version: **v0.12.x**

### Features:
1. Implemented support for ``firewall_image_version`` in **aviatrix_firewall_instance**
2. Implemented support for "UDP" ``vpn_protocol`` for AWS ELB-enabled VPN gateways
3. Implemented support for Active-Active HA (``enable_active_active``) in **aviatrix_site2cloud**

### Enhancements:
1. Implemented coverage for ``tag_list`` formatting change due to Boto3
2. Implemented support for attaching TGW VPN connections to different security domains besides the default domain in **aviatrix_aws_tgw_vpn_conn**
3. Implemented cloud_type check to catch incorrect ha_subnet usage for **aviatrix_gateway** **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
4. Implemented ha_gw_size check to catch incorrect usage when enabling HA for **aviatrix_gateway** **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**

### Bug Fixes:
1. Fixed issue where ``filtered_spoke_vpc_routes`` caused reordering issues for **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**


## 2.10.0 (February 06, 2020)
### Notes:
- Supported Controller version: **UserConnect-5.2.2122**
- Supported Terraform version: **v0.12.x**

### Features:
1. Implemented advanced VPC attachment options for both **aviatrix_aws_tgw** and **aviatrix_aws_tgw_vpc_attachment**
2. Implemented support for updating ``customized_routes`` in **aviatrix_aws_tgw_vpc_attachment**
3. Implemented string length verification for ``aws_account_number`` in **aviatrix_account**
4. Implemented support for ``customized_spoke_vpc_routes``, ``filtered_spoke_vpc_routes`` and ``include/exclude_advertised_spoke_routes`` options in **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**
5. Implemented support for configuring CloudN backup for controller in **aviatrix_controller_config**
6. New resources:
  - **aviatrix_gateway_dnat**
  - **aviatrix_gateway_snat**
7. New data sources:
  - **aviatrix_spoke_gateway**
  - **aviatrix_transit_gateway**
  - **aviatrix_firenet**

### Enhancements:
1. Added coverage for the new resources **aviatrix_gateway_dnat** and **aviatrix_gateway_snat** in test-infra
2. Added coverage for the new data sources **aviatrix_spoke_gateway**, **aviatrix_transit_gateway** and **aviatrix_firenet** in test-infra
3. Deprecated ``dnat_policy`` in **aviatrix_gateway**
4. Deprecated ``dnat_policy``, ``snat_policy`` and ``snat_mode`` in **aviatrix_spoke_gateway**
5. Replaced ``enable_snat`` with ``single_ip_snat`` in **aviatrix_gateway**, **aviatrix_spoke_gateway** and **aviatrix_transit_gateway**

### Bug Fixes:
1. Fixed issue where importing the **aviatrix_aws_tgw** resource results in deltas that could not be rectified through apply


## 2.9.1 (January 28, 2020)
### Notes:
- Supported Controller version: **UserConnect-5.2.2122**
- Supported Terraform version: **v0.12.x**

### Bug Fixes:
1. Fixed issue where JSON Decode ``get_site2cloud_conn_detail`` fails for **aviatrix_site2cloud** and **aviatrix_vgw_conn**


## 2.9.0 (December 20, 2019)
### Notes:
- Supported Controller version: **UserConnect-5.2.2048**
- Supported Terraform version: **v0.12.x**

### Features:
1. Added support for "Designated Gateway" feature in **aviatrix_gateway**
2. Added support for encrypting the AWS EBS volume in **aviatrix_gateway**
3. Added support for "secondary" and "custom" Source NAT in **aviatrix_spoke_gateway**
4. Added support for Destination NAT in **aviatrix_gateway** and **aviatrix_spoke_gateway**
5. New resources:
  - **aviatrix_geo_vpn**

### Enhancements:
1. Migrated from Terraform Core to new Terraform Plugin SDK
2. Added ``elb_dns_name`` as a computed attribute in **aviatrix_gateway**
3. Added coverage for **aviatrix_geo_vpn** in test-infra

### Bug Fixes:
1. Fixed issue where read-back for **aviatrix_gateway**'s ``additional_cidrs_designated_gateway`` incorrectly displayed deltas


## 2.8.0 (December 05, 2019)
### Notes:
- Supported Controller versions: **UserConnect-5.1.1179** and **UserConnect-5.2.1987**
- Supported Terraform version: **v0.12.x**

### Features:
1. Added support for AWS GovCloud access account in **aviatrix_account**
2. Added support for ``customized_routes`` and ``disable_local_route_propagation`` in **aviatrix_aws_tgw_vpc_attachment**
3. Added a link to view the feature compatibility doc online
4. New resources:
  - **aviatrix_aws_tgw_directconnect**

### Enhancements:
1. Added support for updating in **aviatrix_site2cloud** by ForceNew

### Bug Fixes:
1. Fixed an issue that caused an inability to manage a VPN gateway's ``split_tunnel`` attributes after creating the **aviatrix_gateway**


## 2.7.0 (November 07, 2019)
### Notes:
- Supported Controller version: **UserConnect-5.1.973**
- Supported Terraform version: **v0.12.x**

### Features:
1. Added support for attaching/detaching FireNet VPC to/from TGW in **aviatrix_aws_tgw_vpc_attachment**
2. Added support for creating GCP VPC with GCP provider in **aviatrix_vpc**
3. Added support for ``custom_saml_request_template`` in **aviatrix_saml_endpoint**
4. Added support for ``customized_routes`` and ``disable_local_route_propagation`` in **aviatrix_aws_tgw**
5. Added option of retries for ``save`` or ``synchronize`` in **aviatrix_firenet_vendor_integration** data source
6. Added support for VPN NAT for VPN **aviatrix_gateway**
7. Added support for “force-drop” option for policy actions in **aviatrix_firewall**

### Enhancements:
1. Reverted separating ``subnets`` to ``public_subnets`` and ``private_subnets`` in **aviatrix_vpc**
2. Changed calling ``update_access_policy`` from GET to POST in **aviatrix_firewall**

### Bug Fixes:
1. Fixed issue where **aviatrix_gateway** was unable to disable ``split_tunnel``
2. Fixed issue where terraform refresh was not working for firewall policy
3. Fixed issue where **aviatrix_vpc** ``subnets`` were reordering after an import
4. Fixed the issue where creating with special characters causes parsing issue in **aviatrix_firewall_instance**


## 2.6.0 (October 22, 2019)
### Notes:
- Supported Controller version: **UserConnect-5.1.935**
- Supported Terraform version: **v0.12.x**

### Features:
1. New resources:
  - **aviatrix_firewall_instance**
  - **aviatrix_firenet**
2. New data source:
  - **aviatrix_firenet_vendor_integration**
3. Added support to create security domain of ``aviatrix_firewall``, ``native_egress`` or ``native_firewall`` in **aviatrix_aws_tgw**
4. Added support to attach/detach firenet vpc to/from tgw in **aviatrix_aws_tgw**

### Enhancements:
1. Separated ``subnets`` to ``public_subnets`` and ``private_subnets`` in **aviatrix_vpc**
2. Moved ``enable_advertise_transit_cidr`` and ``bgp_manual_spoke_advertise_cidrs`` from **aviatrix_vgw_conn** to **aviatrix_transit_gateway**, and made **aviatrix_vgw_conn** non-updatable
3. Added option to use ``byol`` for test-infra, and updated test-infra to support acceptance test for new resources and data sources
4. Added err body printing for the err that can not decode output of rest api
5. Renamed ``enable_firenet_interfaces`` to ``enable_firenet`` in **aviatrix_transit_gateway**
6. Added option to enable/disable ``single_az_ha`` in **aviatrix_transit_gateway**

### Bug Fixes:
1. Fixed issue where updating aviatrix_account's aws_account_number causes crash


## 2.5.0 (October 02, 2019)
### Notes:
- Supported Controller version: **UserConnect-5.1.738**
- Supported Terraform version: **v0.12.x**

### Features:
1. Added support for enabling/ disabling vpc_dns_server (``enable_vpc_dns_server``) under the AWS (Amazon Web Services) cloud provider for the following resources:
  - **aviatrix_gateway**
  - **aviatrix_spoke_gateway**
  - **aviatrix_transit_gateway**

### Enhancements:
1. Implemented a shell script tool to export test-infra output for acceptance test


## 2.4.0 (September 27, 2019)
### Notes:
- Supported Controller version: **UserConnect-5.0.2761**
- Supported Terraform version: **v0.12.x**

### Features:
1. Added support for OCI (Oracle Cloud Infrastructure) in the following resources:
  - **aviatrix_account**
  - **aviatrix_gateway**
  - **aviatrix_spoke_gateway**
  - **aviatrix_transit_gateway**
2. Added support for GCP (Google Cloud Platform) in **aviatrix_transit_gateway**
3. Updated test-infra to support acceptance test for OCI

### Enhancements:
1. Added ``description`` as an attribute under policy in **aviatrix_firewall**

### Bug Fixes:
1. Fixed issue where HA gateway could not be deleted before the primary gateway for GCP transit gateway


## 2.3.36 (September 16, 2019)
### Notes:
- Supported Controller version: **UserConnect-5.0.2675**
- Supported Controller version: **v0.12.x**

### Bug Fixes:
1. Fixed acceptance test cases


## 2.3.35 (September 10, 2019)
### Notes:
- Supported Controller version: **UserConnect-5.0.2632**
- Supported Terraform version: **v0.12.x**

### Features:
1. Added support for Insane Mode for ARM (Azure Resource Manager) in the following resources:
  - **aviatrix_gateway**
  - **aviatrix_spoke_gateway**
  - **aviatrix_transit_gateway**
2. Added support for ``vgw_account`` and ``vgw_region`` in **aviatrix_vgw_conn**
3. Added support for creating ``aviatrix_firewall``, ``native_egress``, and ``native_aviatrix`` domain in **aviatrix_aws_tgw**
4. Added support for ActiveMesh mode for the following resources:
  - **aviatrix_gateway**
  - **aviatrix_spoke_gateway**
  - **aviatrix_transit_gateway**

### Enhancements:
1. Added ``subnet_id`` as an output attribute for **aviatrix_vpc**
2. Added support to edit ``vpn_cidr`` by gateway instead of just load balancer

### Bug Fixes:
1. Fixed enabling/ disabling advertising CIDRs issue in **aviatrix_vgw_conn**


## 2.2.0 (August 30, 2019)
### Notes:
- Supported Controller version: **UserConnect-4.7.591**
- Supported Terraform version: **v0.12.x**
- Initial Release for Official provider to allow: ``terraform init`` setup


## 2.1.29 (Aug 19 2019)
### Notes:
- Supported Controller version: **UserConnect-4.7.585**
- Supported Terraform version: **v0.12.x**

### Features:
1. Added support for specifying EIP (``allocate_new_eip``, ``eip``, ``ha_eip``) of the primary and HA gateway under the AWS (Amazon Web Services) cloud provider for the following resources:
  - **spoke_gateway**
  - **transit_gateway**
2. Added new resource: **aviatrix_saml_endpoint**. Currently only supports text IDP metadata type


## 2.0.36 (Jul 25 2019)
### Notes:
- Supported Controller version: **UserConnect-4.7.520**
- Supported Terraform version: **v0.12.x**

Major code-base restructuring, featuring renaming of attributes, resources, and attribute values. All these changes are all in the name of standardization of naming conventions and resources

### Changes:
Please see the [R2.0 feature changelist table](https://www.terraform.io/docs/providers/aviatrix/guides/feature-changelist-v2.html#r2-0-userconnect-4-7-patch-terraform-v0-12-) for full details on the changes

---

## 1.16.20 (Jul 25 2019)
### Notes:
- Supported Controller version: **UserConnect-4.7.520**
- Supported Terraform version: **v0.12.x**
- Updated R1.x Feature Changelist

### Enhancements:
1. Now supports Terraform v0.12.x
2. Now uses Go Mod


## 1.15.05 (Jul 15 2019)
### Notes:
- Supported Controller version: **UserConnect-4.7.474**
- Supported Terraform version: **v0.11.x**
- Updated R1.x Feature Changelist

### Enhancements:
1. Added 10s sleep time before updating ``split_tunnel`` for VPN gateway creation
2. Updated test-infra


## 1.14.15 (Jul 11 2019)
### Notes:
- Supported Controller version: **UserConnect-4.7.474**
- Supported Terraform version: **v0.11.x**
- Updated R1.x Feature Changelist

### Features:
1. Added support for ``max_vpn_conn`` in **aviatrix_gateway** resource


## 1.13.14 (Jun 28 2019)
### Notes:
- Supported Controller version: **UserConnect-4.7.419**
- Supported Terraform version: **v0.11.x**
- Added R1.x Feature Changelist

### Enhancements:
1. Added defer function for the following resources:
  - **aviatrix_aws_tgw**
  - **aviatrix_fqdn**
  - **aviatrix_spoke_vpc**
  - **aviatrix_transit_vpc**
  - **aviatrix_site2cloud**
  - **aviatrix_vgw_conn**
2. Added test-infra for Hashicorp acceptance


## 1.12.12 (Jun 20 2019)
### Notes:
- Supported Controller version: **UserConnect-4.7.378**
- Supported Terraform version: **v0.11.x**

### Features:
1. Added support for inside IP CIDR and pre-shared key for tunnel1 and tunnel2 of **aviatrix_aws_tgw_vpn_conn**
  - ``inside_ip_cidr_tun_1``
  - ``inside_ip_cidr_tun_2``
  - ``pre_shared_key_tun_1``
  - ``pre_shared_key_tun_2``

### Enhancements:
1. Added defer function for **aviatrix_gateway**


## 1.11.16 (Jun 18 2019)
### Notes:
- Supported Controller version: **UserConnect-4.7.378**
- Supported Terraform version: **v0.11.x**

### Features:
1. Added support for ``bgp_manual_spoke_advertise_cidrs`` for **aviatrix_vgw_conn** resource
2. Added new resource **aviatrix_vpn_user_accelerator** to support VPN user acceleration through Terraform
3. Added new resource **aviatrix_aws_tgw_vpn_conn** to support attaching/ detaching VPN to TGWs

### Enhancements:
1. Deprecated **version** resource, and changed to an attribute ``target_version`` under **aviatrix_controller_config** to consolidate controller configuration behaviors under one resource


## 1.10.10 (Jun 7 2019)
### Notes:
- Supported Controller version: **UserConnect-4.6.604**
- Supported Terraform version: **v0.11.x**

### Enhancements:
1. Deprecated ``vnet_and_resource_group_names`` and ``vnet_name_resource_group`` in **aviatrix_spoke_vpc** and **aviatrix_transit_vpc**, respectively and replaced with ``vpc_id`` in order to standardize attributes across various cloud providers


## 1.9.28 (Jun 3 2019)
### Notes:
- Supported Controller version: **UserConnect-4.6.569**
- Supported Terraform version: **v0.11.x**

### Features:
1. Added the following support for **aviatrix_site2cloud**:
  - private route encryption (``private_route_encryption``)
  - custom algorithm (``custom_algorithms``)
  - SSL server pool for TCP tunnel types (``ssl_server_pool``)
  - dead peer detection (``enable_dead_peer_detection``)
2. Added support for advertising transit CIDRs (``enable_advertise_transit_cidr``) for **aviatrix_vgw_conn**
3. Added support creating an Aviatrix FireNet VPC (``aviatrix_firenet_vpc``) for **aviatrix_vpc**
4. Added support for enabling a transit gateway for Aviatrix FireNet; (``enable_firenet_interfaces``) in **aviatrix_transit_vpc**

### Enhancements:
1. Deprecated the following resources to consolidate workflow:
  - **aviatrix_admin_email**
  - **aviatrix_customer_id**
2. Deprecated ``cluster`` from **aviatrix_tunnel** resource due to being a deprecated feature in the Controller


## 1.8.26 (May 30 2019)
### Notes:
- Supported Controller version: **UserConnect-4.3.1275**
- Supported Terraform version: **v0.11.x**

### Features:
1. Added support for configuring gateway size for peering HA gateway (``peering_ha_gw_size``) for **aviatrix_gateway**
2. Added Insane Mode support (``insane_mode``, ``insane_mode_az``) for **aviatrix_transit_vpc**
3. Added support for GCP (Google Cloud Platform) in **aviatrix_gateway**
4. Added new resource **aviatrix_arm_peer** to support ARM (Azure Resource Manager) VNet peering
5. Added acceptance test support for import feature for all resources

### Enhancements:
1. Deprecated ``ha_subnet`` from **aviatrix_gateway**


## 1.7.18 (May 9 2019)
### Notes:
- Supported Controller version: **UserConnect-4.3.1253**
- Supported Terraform version: **v0.11.x**

### Features:
1. Added new resource **aviatrix_vpc** to support Controller's Create VPC Tool to create easily create VPCs, subnets
2. Added support for "mapped" connection types (``connection_type``) in **aviatrix_site2cloud**

### Enhancements:
1. Set supportedVersion as a global variable
2. Updated GetVPNUser to call get_vpn_user_by_name instead of list_vpn_user

### Bug Fixes:
1. Fixed **aviatrix_site2cloud**'s ``connection_type`` read/ refresh issue
2. Fixed **aviatrix_vgw_conn** read/ refresh/ import issue


## 1.6.29 (May 3 2019)
### Notes:
- Supported Controller version: **UserConnect-4.2.764**
- Supported Terraform version: **v0.11.x**

### Features:
1. Added ARM (Azure Resource Manager) and GCP (Google Cloud Platform) for **aviatrix_spoke_vpc**
2. Added ARM support for **aviatrix_transit_vpc**
3. Added support for FQDN source IP filtering ``source_ip_list`` in **aviatrix_fqdn** resource
4. Added migration support for **aviatrix_aws_tgw** resource
5. Added **aviatrix_controller_config** resource that supports the following features:
  - system-wide FQDN exception rule (``fqdn_exception_rule``)
  - security group management (``security_group_management``)
  - http access (``http_access``)

### Enhancements:
1. Added controller version checking in the provider to ensure compatibility between Aviatrix Terraform provider and Controller


## 1.5.24 (Apr 15 2019)
### Notes:
- Supported Controller version: **UserConnect-4.2.764**
- Supported Terraform version: **v0.11.x**

### Features:
1. Added ARM (Azure Resource Manager) and GCP (Google Cloud Platform) support for **aviatrix_account**

### Enhancements:
1. Moved goaviatrix library from vendor to root folder
2. Deprecated ``dns_server`` for the following resources:
  - **aviatrix_gateway**
  - **aviatrix_spoke_vpc**
  - **aviatrix_transit_vpc**
3. Added description for all attributes
4. Added import support for **aviatrix_gateway**'s ``split_tunnel``

### Bug Fixes:
1. Fixed migration/ update issue for ``manage_vpc_attachment`` in **aviatrix_aws_tgw** resource
2. Fixed failing to destroy **aviatrix_vgw_conn** despite being destroyed in Controller UI
3. Fixed refresh issue for deleted **aviatrix_fqdn** through Controller UI
4. Fixed read/ refresh issue for **aviatrix_site2cloud** where resource count exceeds 3


## 1.4.4 (Mar 28 2019)
### Notes:
- Supported Controller version: **UserConnect-4.2.634**
- Supported Terraform version: **v0.11.x**

### Features:
1. Added new resource **aviatrix_aws_tgw_vpc_attachment** to simplify/ add an option on how users can choose to manage attaching/ detaching VPCs to and from their **aviatrix_aws_tgw**

### Enhancements:
1. Updated **aviatrix_aws_tgw** to allow creation of only the TGW, as well as allowing management of VPC attachments to be done either within the resource, or though **aviatrix_aws_tgw_vpc_attachment**
2. updated documentation for **aviatrix_aws_peer** resource
3. updated **aviatrix_fqdn** to block updating ``fqdn_tag``


## 1.3.12 (Mar 21 2019)
### Notes:
- Supported Controller version: **UserConnect-4.1.982** and **4.2.634**
- Supported Terraform version: **v0.11.x**

### Enhancements:
1. Deprecated **aviatrix_dc_extn** resource due to removed support from Controller
2. Added version information

### Bug Fixes:
1. Fixed **aviatrix_firewall**'s ``base_allow_deny`` on refresh
2. Fixed **aviatrix_site2cloud**'s refresh, update and import issues
3. Fixed **aviatrix_aws_peer**'s refresh, update and import issues


## 1.2.12 (Mar 15 2019)
### Notes:
- Supported Controller version: **UserConnect-4.1.981**
- Supported Terraform version: **v0.11.x**

### Changes:
1. Temporarily reverted refresh changes for the following resources:
  - **aviatrix_aws_peer**
  - **aviatrix_site2cloud**

### Bug Fixes:
1. Fixed **aviatrix_site2cloud** to ignore ``local_subnet_cidr`` changes


## 1.2.10 (Mar 14 2019)

-> **NOTE:** This release is unsupported and deprecated

### Notes:
- Supported Controller version: **UserConnect-4.1.981**
- Supported Terraform version: **v0.11.x**

### Bug Fixes:
1. Fixed ``tag_list`` reordering issue on **aviatrix_gateway**
2. Fixed refresh issues for the following resources:
  - **aviatrix_aws_peer**
  - **aviatrix_site2cloud**
  - **aviatrix_vgw_conn**


## 1.1.66 (Mar 6 2019)
### Notes:
- Supported Controller version: **UserConnect-4.1.981**
- Supported Terraform version: **v0.11.x**

### Features:
1. Added support for specifying EIP (``peering_ha_eip``) for the HA gateway in **aviatrix_gateway** resource
2. All resources now support ``terraform import``

### Enhancements:
1. Enhanced returned error messages to show REST API names
2. Deprecated ``over_aws_peering`` in **aviatrix_tunnel** resource
3. Enhanced refresh functionality for the following resources:
  - **aviatrix_aws_tgw**
  - **aviatrix_admin_email**
4. **aviatrix_firewall** resource enhanced to have policy validation

### Bug Fixes:
1. Fixed URL encode error for all resources
2. Fixed port requirement for ICMP protocol in **aviatrix_fqdn**
3. Fixed **aviatrix_transit_vpc** resource to support empty ``tag_list``
4. Fixed **aviatrix_vpn_user** re-ordering issue on refresh


## 1.0.242 (Feb 26 2019)
### Notes:
- Supported Controller version: **UserConnect-4.1.981**
- Supported Terraform version: **v0.11.x**

### Features:
1. Support for Terraform's ``create``, ``destroy``, ``refresh``, ``update``, and acceptance tests for most of the following resources:
  - **data_source_aviatrix_account**
  - **data_source_aviatrix_caller_identity**
  - **data_source_aviatrix_gateway**
  - **resource_aviatrix_account**
  - **resource_aviatrix_account_user**
  - **resource_aviatrix_admin_email**
  - **resource_aviatrix_aws_peer**
  - **resource_aviatrix_aws_tgw**
  - **resource_aviatrix_customer_id**
  - **resource_aviatrix_dc_extn**
  - **resource_aviatrix_firewall**
  - **resource_aviatrix_firewall_tag**
  - **resource_aviatrix_fqdn**
  - **resource_aviatrix_gateway**
  - **resource_aviatrix_site2cloud**
  - **resource_aviatrix_spoke_vpc**
  - **resource_aviatrix_transit_gateway_peering**
  - **resource_aviatrix_transit_vpc**
  - **resource_aviatrix_transitive_peering**
  - **resource_aviatrix_tunnel**
  - **resource_aviatrix_version**
  - **resource_aviatrix_vgw_conn**
  - **resource_aviatrix_vpn_profile**
- **resource_aviatrix_vpn_user**
