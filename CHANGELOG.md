## 2.1.29 (Aug 19 2019)

CHANGES
  - Supported controller version: 4.7.585
  - Supported Terraform version: 0.12.*
  - Added new attribute support of allocate_new_eip, eip, ha_eip in spoke_gateway
  - Added new attribute support of allocate_new_eip, eip, ha_eip in transit_gateway
  - Added new resource saml_endpoint, currently only "Text" is supported
  
  
## 2.0.36 (Jul 25 2019)

CHANGES
  - Supported controller version: 4.7.520
  - Supported Terraform version: 0.12.*
  - Added Terraform Feature Changelist v2
  - Added new resource aviatrix_spoke_gateway
  - Added new resource aviatrix_transit_gateway
  - Renamed "vpc_size" and "vpc_net" to "gw_size" and "subnet" respectively, changed "enable_nat" to "enable_snat" and its type from
  string to boolean in gateway resource
  - Renamed "vpc_name1" and "vpc_name2" to "gw_name1" and "gw_name2" respectively in tunnel resource
  - Renamed "base_allow_deny" to "base_policy", renamed "base_log_enable" to "base_log_enabled" and changed its type from string
  to boolean, renamed "allow_deny" to "action", renamed "log_enable" to "log_enabled" and changed its type from string to boolean
  in firewall resource
  - Renamed "fqdn_status" to "fqdn_enabled" in fqdn resource
  - Changed type from string to boolean for "vpn_access", "enable_elb", "split_tunnel", "saml_enabled", "enable_ldap", "single_az_ha",
  and "allocate_new_eip" in gateway resource
  - Changed type from string to boolean for "ha_enabled" in site2cloud resource
  - Changed type from string to boolean for "enable_ha" in tunnel resource
  
  
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
  - Added max_vpn_conn support in gateway
  - Updated Terraform Feature Changelist
  

## 1.13.14 (Jun 28 2019)

CHANGES
  - Supported controller version: 4.7.419
  - Supported Terraform version: 0.11.*
  - Added defer function in transit_vpc, spoke_vpc, aws_tgw, site2cloud, vgw_conn, and fqdn
  - Added Terraform Feature Changelist
  - Added test-infra for Hashicorp acceptance
  
  
## 1.12.12 (Jun 20 2019)

CHANGES
  - Supported controller version: 4.7.378
  - Supported Terraform version: 0.11.*
  - Added support of Inside IP CIDR and Pre Shared Key for tunnel 1 and 2
  - Added defer function in gateway


## 1.11.16 (Jun 18 2019)

CHANGES
  - Supported controller version: 4.7.378
  - Supported Terraform version: 0.11.*
  - Deprecated version resource, and moved version functionality to controller config
  - Added "bgp_manual_spoke_advertise_cidr" support in vgw_conn
  - Created a new resource vpn_user_accelerator
  - Created a new resource to support attaching/detaching vpn to tgw in aws_tgw
  
  
## 1.10.10 (Jun 7 2019)

CHANGES
  - Supported controller version: 4.6.604
  - Supported Terraform version: 0.11.*
  - Deprecated "vnet_and_resource_group_names" in spoke_vpc and replaced it with "vpc_id"
  - Deprecated "vnet_name_resource_group" in transit_vpc and replaced it with "vpc_id"
  
    
## 1.9.28 (Jun 3 2019)

CHANGES
  - Supported controller version: 4.6.569
  - Supported Terraform version: 0.11.*
  - Added private route encryption support in site2cloud
  - Added custom algorithm support in site2cloud
  - Added ssl server pool support for "tcp" tunnel type in site2cloud
  - Added dead peer detection support in site2cloud
  - Added advertise transit CIDR support in vgw_conn
  - Added aviatrix firenet vpc support in vpc
  - Added enable firenet interfaces support in transit_vpc
  - Deprecated "cluster" in tunnel resource
  - Deprecated admin_email and customer_id resources
  
  
## 1.8.26 (May 30 2019)

CHANGES
  - Supported controller version: 4.3.1275
  - Supported Terraform version: 0.11.*
  - Deprecated “ha_subnet” (gateway HA) completely from gateway
  - Added ability to configure gateway size for peering HA gateway in gateway resource
  - Added acceptance test support for import feature for all resources
  - Added insane mode support in transit_vpc
  - Added new resource arm_peer
  - Added GCP support in gateway
  
  
## 1.7.18 (May 9 2019)

CHANGES
  - Supported controller version: 4.3.1253
  - Supported Terraform version: 0.11.*
  - Added new vpc resource
  - Added support for connection type “mapped”
  - Fixed connection_type read/refresh issue
  - Fixed vgwConn resource refresh/import issue
  - Set supportedVersion as a global variable
  - Updated GevVPNUser to by get_vpn_user_by_name instead of list_vpn_users 
  
  
## 1.6.29 (May 3 2019)

CHANGES
  - Supported controller version: 4.2.764
  - Supported Terraform version: 0.11.*
  - GCP and Azure support added for resource_spoke_gateway
  - Azure support added for resource_transit_vpc
  - Added support for FQDN source_ip filtering
  - Added controller configuration resource
  - Added migration support for aws_tgw resource
  - Added exception_rule support
  - Added security_group_management support
  - Added controller version check functionality
  

## 1.5.24 (Apr 15 2019)

CHANGES
  - Supported controller version: 4.2.764
  - Supported Terraform version: 0.11.*
  - Description added for all argument
  - GCP and Azure support added for resource_account
  - Updated gateway resource for Split_tunnel import support
  - Fixed migration/update issue for "manage_vpc_attachment" in aws_tgw resource
  - Fixed failing to destroy vgw_conn deleted through UI issue
  - Fixed refresh issue for fqdn deleted through UI
  - Moved goaviatrix library from vendor to root folder
  - Fixed read/refresh issue for more than 3 site2cloud instances
  - Deprecated dns_server for gateway, transit gw, and spoke gw
  
  
## 1.4.4 (Mar 28 2019)

CHANGES
  - Supported controller version: 4.2.634
  - Supported Terraform version: 0.11.*
  - Updated doc for aws_peer resource
  - Updated fqdn resource to block updating "fqdn_tag" 
  - Created new resource "aws_tgw_vpc_attachment" to simply manage attaching/detaching vpc to/from AWS TGW
  - Updated aws_tgw resource to allow creating aws tgw only, and attaching/detaching vpc to/from tgw using aws_tgw_vpc_attachment


## 1.3.12 (Mar 21 2019)

CHANGES
  - Supported controller version: 4.1.982 and 4.2.634
  - Supported Terraform version: 0.11.*
  - Fixed firewall resource base_allow_deny on refresh
  - Fixed site2cloud resource arguments on refresh, update and import
  - Fixed aws_peer resource arguments on refresh, update and import
  - Deprecated dc_extn resource
  - Added version information 
  

## 1.2.12 (Mar 15 2019)

CHANGES
  - Supported controller version: 4.1.981
  - Supported Terraform version: 0.11.* 
  - Temporarily reverted peering resource refresh changes
  - Temporarily reverted site2cloud resource refresh changes
  - Updated site2cloud resource to ignore local_cidr changes

  
## 1.2.10 (Mar 14 2019)

CHANGES
  - Supported controller version: 4.1.981
  - Supported Terraform version: 0.11.*
  - Updated peering resource to support refresh
  - Updated site2cloud resource to support refresh of some paramters
  - Corrected taq list reordering on gateway resource refresh
  - Corrected VGW resource on refresh

  
## 1.1.66 (Mar 6 2019)

CHANGES
  - Supported controller version: 4.1.981
  - Supported Terraform version: 0.11.*
  - Supports import feature for all resources
  - URL encode error is fixed for all resources
  - Error messages show REST api names for better understanding
  - Added EIP for peering HA gateways
  - Fixed port requirement for ICMP protocol in FQDN resource
  - Deprecated over_aws_peering in aviatrix_tunnel resource
  - Updated refresh for tgw, admin_email resource
  - Policy validation in firewall resource
  - Support empty tag list in transit_vpc resource
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
      - resource_account
      - resource_account_user
      - resource_admin_email
      - resource_aws_peer
      - resource_aws_tgw
      - resource_customer_id
      - resource_dc_extn
      - resource_firewall
      - resource_firewall_tag
      - resource_fqdn
      - resource_gateway
      - resource_site2cloud
      - resource_spoke_vpc
      - resource_transit_gateway_peering
      - resource_transit_vpc
      - resource_transitive_peering
      - resource_tunnel
      - resource_version
      - resource_vgw_conn
      - resource_vpn_profile
      - resource_vpn_user 
