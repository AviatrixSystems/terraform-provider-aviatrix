## 1.1.66 Mar 6 2019)

CHANGES
 
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
