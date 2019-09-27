## 2.4.1 (Unreleased)
## 2.4.0 (September 27, 2019)

NOTES:
  - Supported controller version: UserConnect-5.0.2761
  - Supported Terraform version: 0.12.*  
  
FEATURES:
  - Add support for OCI (Oracle Cloud Infrastructure) in resource_aviatrix_account
  - Add support for OCI in resource_aviatrix_gateway, resource_aviatrix_spoke_gateway and resource_aviatrix_transit_gateway
  - Add support for GCP in resource_aviatrix_transit_gateway
  - Update test-infra to support acceptance test for OCI 
  
ENHANCEMENTS
  - Add "description" as an attribute under policy in resource_aviatrix_firewall
  
BUG FIXES:
  - Fix the issue that HA gateway could not be deleted before primary gateway for GCP transit gateway

## 2.3.36 (September 16, 2019)

  - Supported controller version: 5.0.2675
  - Supported Terraform version: 0.12.*  
  
BUG FIXES:
  - Acceptance test cases
  
## 2.3.35 (September 10, 2019)

CHANGES
  - Supported controller version: 5.0.2632
  - Supported Terraform version: 0.12.*
  - Added support for Azure RM insane mode for resource_aviatrix_gateway, resource_aviatrix_spoke_gateway and resource_aviatrix_transit_gateway
  - "sing_az_ha" is changed to default true for resource_aviatrix_spoke_gateway
  - Added subnet ID as an attribute of the output in resource_aviatrix_vpc
  - Added support to edit vpn CIDR by gateway instead of LB
  - Added "vgw_account" and "vgw_region" support in resource_vaviatrix_gw_conn
  - Added support of creating "aviatrix_firewall", "native_egress" or "native_aviatrix" domain in resource_aviatrix_aws_tgw
  - Fixed enable/disable advertise_cidrs issue in resource_aviatrix_vgw_conn
  - Added active_mesh mode support for resource_aviatrix_spoke_gateway and resource_aviatrix_transit_gateway, default value: false 

## 2.2.0 (August 30, 2019)

CHANGES
  - Supported controller version: 4.7.591
  - Supported Terraform version: 0.12.*
  - Initial Release for 'terraform init'

## 2.1.29 (Aug 19 2019)

CHANGES
  - Supported controller version: 4.7.585
  - Supported Terraform version: 0.12.*
  - Added new attribute support of "allocate_new_eip", "eip", "ha_eip" in resource_aviatrix_spoke_gateway
  - Added new attribute support of "allocate_new_eip", "eip", "ha_eip" in resource_aviatrix_transit_gateway
  - Added new resource_aviatrix_saml_endpoint, currently only supporting "Text"
  
  
## 2.0.36 (Jul 25 2019)

CHANGES
  - Supported controller version: 4.7.520
  - Supported Terraform version: 0.12.*
  - Added Terraform Feature Changelist v2
  - Added new resource_aviatrix_spoke_gateway
  - Added new resource_aviatrix_transit_gateway
  - Renamed "vpc_size" and "vpc_net" to "gw_size" and "subnet" respectively, changed "enable_nat" to "enable_snat" and its type from
  string to boolean in resource_aviatrix_gateway
  - Renamed "vpc_name1" and "vpc_name2" to "gw_name1" and "gw_name2" respectively in resource_aviatrix_tunnel
  - Renamed "base_allow_deny" to "base_policy", renamed "base_log_enable" to "base_log_enabled" and changed its type from string
  to boolean, renamed "allow_deny" to "action", renamed "log_enable" to "log_enabled" and changed its type from string to boolean
  in resource_aviatrix_firewall
  - Renamed "fqdn_status" to "fqdn_enabled" in resource_aviatrix_fqdn
  - Changed type from string to boolean for "vpn_access", "enable_elb", "split_tunnel", "saml_enabled", "enable_ldap", "single_az_ha",
  and "allocate_new_eip" in resource_aviatrix_gateway
  - Changed type from string to boolean for "ha_enabled" in resource_aviatrix_site2cloud
  - Changed type from string to boolean for "enable_ha" in resource_aviatrix_tunnel
  
  
## 1.16.20 (Jul 25 2019)

CHANGES
  - Supported controller version: 4.7.520
  - Supported Terraform version: 0.12.*
  - Updated Terraform Feature Changelist
  - Use Go Mod
  
  
## 1.15.05 (Jul 15 2019)

CHANGES
  - Supported controller version: 4.7.474
  - Supported Terraform version: 0.11.*
  - Updated Terraform Feature Changelist
  - Updated test infra
  - Added 10s sleep time before updating split tunnel for vpn gateway creation
  
  
## 1.14.15 (Jul 11 2019)

CHANGES
  - Supported controller version: 4.7.474
  - Supported Terraform version: 0.11.*
  - Added max_vpn_conn support in resource_aviatrix_gateway
  - Updated Terraform Feature Changelist
  

## 1.13.14 (Jun 28 2019)

CHANGES
  - Supported controller version: 4.7.419
  - Supported Terraform version: 0.11.*
  - Added defer function in resource_aviatrix_transit_vpc, resource_aviatrix_spoke_vpc, resource_aviatrix_aws_tgw, 
  resource_aviatrix_site2cloud, resource_aviatrix_vgw_conn, and resource_aviatrix_fqdn
  - Added Terraform Feature Changelist
  - Added test-infra for HashiCorp acceptance
  
  
## 1.12.12 (Jun 20 2019)

CHANGES
  - Supported controller version: 4.7.378
  - Supported Terraform version: 0.11.*
  - Added support of Inside IP CIDR and Pre Shared Key for tunnel 1 and 2
  - Added defer function in resource_aviatrix_gateway


## 1.11.16 (Jun 18 2019)

CHANGES
  - Supported controller version: 4.7.378
  - Supported Terraform version: 0.11.*
  - Deprecated resource_aviatrix_version, and moved version functionality to controller config
  - Added "bgp_manual_spoke_advertise_cidr" support in resource_aviatrix_vgw_conn
  - Created new resource_aviatrix_vpn_user_accelerator
  - Created new to support attaching/detaching vpn to tgw in resource_aviatrix_aws_tgw
  
  
## 1.10.10 (Jun 7 2019)

CHANGES
  - Supported controller version: 4.6.604
  - Supported Terraform version: 0.11.*
  - Deprecated "vnet_and_resource_group_names" in resource_aviatrix_spoke_vpc and replaced it with "vpc_id"
  - Deprecated "vnet_name_resource_group" in resource_aviatrix_transit_vpc and replaced it with "vpc_id"
  
    
## 1.9.28 (Jun 3 2019)

CHANGES
  - Supported controller version: 4.6.569
  - Supported Terraform version: 0.11.*
  - Added private route encryption support in resource_aviatrix_site2cloud
  - Added custom algorithm support in resource_aviatrix_site2cloud
  - Added ssl server pool support for "tcp" tunnel type in resource_aviatrix_site2cloud
  - Added dead peer detection support in resource_aviatrix_site2cloud
  - Added advertise transit CIDR support in resource_aviatrix_vgw_conn
  - Added aviatrix firenet vpc support in resource_aviatrix_vpc
  - Added enable firenet interfaces support in resource_aviatrix_transit_vpc
  - Deprecated "cluster" in resource_aviatrix_tunnel
  - Deprecated aviatrix_admin_email and customer_id resources
  
  
## 1.8.26 (May 30 2019)

CHANGES
  - Supported controller version: 4.3.1275
  - Supported Terraform version: 0.11.*
  - Deprecated “ha_subnet” (gateway HA) completely from resource_aviatrix_gateway
  - Added ability to configure gateway size for peering HA gateway in resource_aviatrix_gateway
  - Added acceptance test support for import feature for all resources
  - Added insane mode support in resource_aviatrix_transit_vpc
  - Added new resource_aviatrix_arm_peer
  - Added GCP support in resource_aviatrix_gateway
  
  
## 1.7.18 (May 9 2019)

CHANGES
  - Supported controller version: 4.3.1253
  - Supported Terraform version: 0.11.*
  - Added new resource_aviatrix_vpc
  - Added support for connection type “mapped”
  - Fixed connection_type read/refresh issue
  - Fixed resource_aviatrix_vgwConn refresh/import issue
  - Set supportedVersion as a global variable
  - Updated GevVPNUser to by get_vpn_user_by_name instead of list_vpn_users 
  
  
## 1.6.29 (May 3 2019)

CHANGES
  - Supported controller version: 4.2.764
  - Supported Terraform version: 0.11.*
  - GCP and Azure support added for resource_aviatrix_spoke_gateway
  - Azure support added for resource_aviatrix_transit_vpc
  - Added support for FQDN source_ip filtering
  - Added new resource_aviatrix_controller_config
  - Added migration support for resource_aviatrix_aws_tgw
  - Added exception_rule support
  - Added security_group_management support
  - Added controller version check functionality
  

## 1.5.24 (Apr 15 2019)

CHANGES
  - Supported controller version: 4.2.764
  - Supported Terraform version: 0.11.*
  - Description added for all argument
  - GCP and Azure support added for resource_aviatrix_account
  - Updated resource_aviatrix_gateway for Split_tunnel import support
  - Fixed migration/update issue for "manage_vpc_attachment" in resource_aviatrix_aws_tgw
  - Fixed failing to destroy vgw_conn deleted through UI issue
  - Fixed refresh issue for fqdn deleted through UI
  - Moved goaviatrix library from vendor to root folder
  - Fixed read/refresh issue for more than 3 site2cloud instances
  - Deprecated dns_server for resource_aviatrix_gateway, resource_aviatrix_transit_vpc, and resource_aviatrix_spoke_vpc
  
  
## 1.4.4 (Mar 28 2019)

CHANGES
  - Supported controller version: 4.2.634
  - Supported Terraform version: 0.11.*
  - Updated doc for resource_aviatrix_aws_peer
  - Updated resource_aviatrix_fqdn to block updating "fqdn_tag" 
  - Created new resource_aviatrix_aws_tgw_vpc_attachment to simply manage attaching/detaching vpc to/from AWS TGW
  - Updated resource_aviatrix_aws_tgw to allow creating aws tgw only, and attaching/detaching vpc to/from tgw using resource_aviatrix_aws_tgw_vpc_attachment


## 1.3.12 (Mar 21 2019)

CHANGES
  - Supported controller version: 4.1.982 and 4.2.634
  - Supported Terraform version: 0.11.*
  - Fixed resource_aviatrix_firewall "base_allow_deny" on refresh
  - Fixed resource_aviatrix_site2cloud arguments on refresh, update and import
  - Fixed resource_aviatrix_aws_peer arguments on refresh, update and import
  - Deprecated dc_extn resource
  - Added version information 
  

## 1.2.12 (Mar 15 2019)

CHANGES
  - Supported controller version: 4.1.981
  - Supported Terraform version: 0.11.* 
  - Temporarily reverted resource_aviatrix_peering refresh changes
  - Temporarily reverted resource_aviatrix_site2cloud refresh changes
  - Updated resource_aviatrix_site2cloud to ignore local_cidr changes

  
## 1.2.10 (Mar 14 2019)

CHANGES
  - Supported controller version: 4.1.981
  - Supported Terraform version: 0.11.*
  - Updated peering resource to support refresh
  - Updated resource_aviatrix_site2cloud to support refresh of some parameters
  - Corrected taq list reordering on resource_aviatrix_gateway refresh
  - Corrected resource_aviatrix_vgw_conn on refresh

  
## 1.1.66 (Mar 6 2019)

CHANGES
  - Supported controller version: 4.1.981
  - Supported Terraform version: 0.11.*
  - Supports import feature for all resources
  - URL encode error is fixed for all resources
  - Error messages show REST api names for better understanding
  - Added EIP for peering HA gateways
  - Fixed port requirement for ICMP protocol in resource_aviatrix_fqdn
  - Deprecated over_aws_peering in resource_aviatrix_tunnel
  - Updated refresh for tgw, admin_email resource
  - Policy validation in resource_aviatrix_firewall
  - Support empty tag list in resource_aviatrix_transit_vpc
  - Fixed VPN profile user re-ordering on refresh

  
## 1.0.242 (Tue Feb 26 2019)

CHANGES
 
  - First versioned release
  - Supported controller version: 4.1.981
  - Supported Terraform version: 0.11.*
  - Supports create, destroy, refresh, update, acceptance tests for most of the following resources
      - data_source_aviatrix_account
      - data_source_aviatrix_caller_identity
      - data_source_aviatrix_gateway
      - resource_aviatrix_account
      - resource_aviatrix_account_user
      - resource_aviatrix_admin_email
      - resource_aviatrix_aws_peer
      - resource_aviatrix_aws_tgw
      - resource_aviatrix_customer_id
      - resource_aviatrix_dc_extn
      - resource_aviatrix_firewall
      - resource_firewall_tag
      - resource_aviatrix_fqdn
      - resource_aviatrix_gateway
      - resource_aviatrix_site2cloud
      - resource_aviatrix_spoke_vpc
      - resource_aviatrix_transit_gateway_peering
      - resource_aviatrix_transit_vpc
      - resource_aviatrix_transitive_peering
      - resource_aviatrix_tunnel
      - resource_aviatrix_version
      - resource_aviatrix_vgw_conn
      - resource_aviatrix_vpn_profile
      - resource_aviatrix_vpn_user 
