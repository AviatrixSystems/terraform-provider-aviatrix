## 1.5.24 (Apr 15 2019)

CHANGES
  - Supported controller version: 4.2.764
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
  - Updated doc for aws_peer resource
  - Updated fqdn resource to block updating "fqdn_tag" 
  - Created new resource "aws_tgw_vpc_attachment" to simply manage attaching/detaching vpc to/from AWS TGW
  - Updated aws_tgw resource to allow creating aws tgw only, and attaching/detaching vpc to/from tgw using aws_tgw_vpc_attachment


## 1.3.12 (Mar 21 2019)

CHANGES
  - Supported controller version: 4.1.982 and 4.2.634
  - Fixed firewall resource base_allow_deny on refresh
  - Fixed site2cloud resource arguments on refresh, update and import
  - Fixed aws_peer resource arguments on refresh, update and import
  - Deprecated dc_extn resource
  - Added version information 
  

## 1.2.12 (Mar 15 2019)

CHANGES
  - Supported controller version: 4.1.981 
  - Temporarily reverted peering resource refresh changes
  - Temporarily reverted site2cloud resource refresh changes
  - Updated site2cloud resource to ignore local_cidr changes

  
## 1.2.10 (Mar 14 2019)

CHANGES
  - Supported controller version: 4.1.981
  - Updated peering resource to support refresh
  - Updated site2cloud resource to support refresh of some paramters
  - Corrected taq list reordering on gateway resource refresh
  - Corrected VGW resource on refresh

  
## 1.1.66 (Mar 6 2019)

CHANGES
  - Supported controller version: 4.1.981
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
