## 2.16.2 (Unreleased)
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

### Features
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

### Enhancements
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
