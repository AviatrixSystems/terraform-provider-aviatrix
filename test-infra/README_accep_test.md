# Acceptance Tests

#### Pre-requisites

- The controller must be launched before hand and must be up and running the latest controller version
- IAM roles (aviatrix-role-ec2 and aviatrix-role-app) also must be created and attached if any IAM role related tests are to be run. Currently all tests are based on Access key, Secret key
- The VPC's with public subnet to launch the gateways must be created before the tests
- If you are running aviatrix_aws_peer or aviatrix_peer, two VPC's with non overlapping CIDR's must be created before hand
- If you are running the tests on a BYOL controller, the customer ID must be set prior to the tests, otherwise run the tests on a PayG metered controller
- aviatrix_aws_tgw test only allows Transit GWs and VPCs to be attached to the TGW in the same region
- AWS_ACCOUNT_NUMBER should be the same one used for controller launch

#### Skip parameters and variables

Passing an environment value of "yes" to the skip parameter allows you to skip the particular resource. If it is not skipped, it checks for the existence of other required variables. Generic variables are required for any acceptance test

| Test module name                     | Skip parameter                     | Required variables                                                             |
| ------------------------------------ | ---------------------------------- | ------------------------------------------------------------------------------ |
| Generic                              | N/A                                | AVIATRIX_USERNAME, AVIATRIX_PASSWORD, AVIATRIX_CONTROLLER_IP                   |
| aviatrix_account                     | SKIP_ACCOUNT                       |                                                                                |
|		                               | SKIP_ACCOUNT_AWS	                | AWS_ACCOUNT_NUMBER, AWS_ACCESS_KEY, AWS_SECRET_KEY                             |
|                     		           | SKIP_ACCOUNT_GCP	                | GCP_ID, GCP_CREDENTIALS_FILEPATH	                                             |
|		                               | SKIP_ACCOUNT_AZURE	                | ARM_SUBSCRIPTION_ID, ARM_DIRECTORY_ID, ARM_APPLICATION_ID, ARM_APPLICATION_KEY |
|                     		           | SKIP_ACCOUNT_OCI	                | OCI_TENANCY_ID, OCI_USER_ID, OCI_COMPARTMENT_ID, OCI_API_KEY_FILEPATH          |
|		                               | SKIP_ACCOUNT_AWSGOV                | AWSGOV_ACCOUNT_NUMBER, AWSGOV_ACCESS_KEY, AWSGOV_SECRET_KEY                    |
| aviatrix_account_user                | SKIP_ACCOUNT_USER                  |                                                                                |
| aviatrix_arm_peer                    | SKIP_ARM_PEER                      | aviatrix_account + ARM_VNET_ID, ARM_VNET_ID2, ARM_REGION, ARM_REGION2          |
| aviatrix_aws_peer                    | SKIP_AWS_PEER                      | aviatrix_account + AWS_VPC_ID, AWS_VPC_ID2, AWS_REGION, AWS_REGION2            |
| aviatrix_aws_tgw                     | SKIP_AWS_TGW                       | aviatrix_account + AWS_VPC_ID, AWS_REGION, AWS_VPC_TGW_ID                      |
| aviatrix_aws_tgw_directconnect       | SKIP_AWS_TGW_DIRECTCONNECT         | aviatrix_aws_tgw + AWS_DX_GATEWAY                                              |
| aviatrix_aws_tgw_peering             | SKIP_AWS_TGW_PEERING               | aviatrix_account                                                               |
| aviatrix_aws_tgw_peering_domain_conn | SKIP_AWS_TGW_PEERING_DOMAIN_CONN   | aviatrix_account                                                               |
| aviatrix_aws_tgw_vpc_attachment      | SKIP_AWS_TGW_VPC_ATTACHMENT        | aviatrix_aws_tgw                                                               |
| aviatrix_aws_tgw_vpc_attachment      | SKIP_AWS_TGW_TRANSIT_GATEWAY_ATTACHMENT | aviatrix_aws_tgw + aviatrix_transit_gateway                               |
| aviatrix_azure_peer                  | SKIP_AZURE_PEER                    | aviatrix_account + AZURE_VNET_ID, AZURE_VNET_ID2, AZURE_REGION, AZURE_REGION2  |
| aviatrix_azure_spoke_native_peering  | SKIP_AZURE_SPOKE_NATIVE_PEERING    | aviatrix_account + AZURE_VNET_ID, AZURE_VNET_ID2, AZURE_REGION, AZURE_REGION2  |
| aviatrix_azure_vng_conn              | SKIP_AZURE_VNG_CONN                | aviatrix_account + AZURE_VNG_VNET_ID, AZURE_REGION, AZURE_VNG_SUBNET, AZURE_VNG|
| aviatrix_aws_tgw_vpn_conn            | SKIP_AWS_TGW_VPN_CONN              | aviatrix_aws_tgw                                                               |
| aviatrix_controller_config           | SKIP_CONTROLLER_CONFIG             | aviatrix_account                                                               |
| aviatrix_datadog_agent               | SKIP_DATADOG_AGENT                 | datadog_api_key
| aviatrix_device_aws_tgw_attachment   | SKIP_DEVICE_AWS_TGW_ATTACHMENT     | DEVICE_NAME, AWS_TGW_NAME                                                      |
| aviatrix_device_interface_config     | SKIP_DEVICE_INTERFACE_CONFIG       | aviatrix_device_registration                                                   |
| aviatrix_device_registration         | SKIP_DEVICE_REGISTRATION           | DEVICE_PUBLIC_IP, DEVICE_KEY_FILE_PATH                                         |
| aviatrix_device_tag                  | SKIP_DEVICE_TAG                    | aviatrix_device_registration                                                   |
| aviatrix_device_transit_gateway_attachment | SKIP_DEVICE_TRANSIT_GATEWAY_ATTACHMENT | aviatrix_device_registration, TRANSIT_GATEWAY_NAME                   |
| aviatrix_device_virtual_wan_attachment | SKIP_DEVICE_VIRTUAL_WAN_ATTACHMENT | aviatrix_device_registration, aviatrix_account on AZURE, ARM_RESOURCE_GROUP, ARM_HUB_NAME |
| aviatrix_firenet                     | SKIP_FIRENET                       | aviatrix_account + AWS_REGION, Palo Alto VM series                             |
| aviatrix_firewall                    | SKIP_FIREWALL                      | aviatrix_gateway                                                               |
| aviatrix_firewall_instance           | SKIP_FIREWALL_INSTANCE             | aviatrix_account + AWS_REGION, Palo Alto VM series                             |
| aviatrix_firewall_instance_association | SKIP_FIREWALL_INSTANCE_ASSOCIATION | aviatrix_firenet, transit_gateway                                            |
| aviatrix_firewall_management_access  | SKIP_FIREWALL_MANAGEMENT_ACCESS    | aviatrix_account                                                               |
|                                      | SKIP_FIREWALL_MANAGEMENT_ACCESS_AWS|       + aviatrix_transit_gateway + aviatrix_spoke_gateway in AWS               |
|                                      | SKIP_FIREWALL_MANAGEMENT_ACCESS_AZURE|     + aviatrix_transit_gateway + aviatrix_spoke_gateway in AZURE             |   
| aviatrix_firewall_policy             | SKIP_FIREWALL_POLICY               | aviatrix_gateway                                                               |
| aviatrix_firewall_tag                | SKIP_FIREWALL_TAG                  |                                                                                |
| aviatrix_fqdn                        | SKIP_FQDN                          | aviatrix_gateway                                                               |
| aviatrix_fqdn_pass_through           | SKIP_FQDN_PASS_THROUGH             | aviatrix_gateway                                                               |
| aviatrix_fqdn_tag_rule               | SKIP_FQDN_TAG_RULE                 | aviatrix_gateway                                                               |
| aviatrix_gateway                     | SKIP_GATEWAY                       | aviatrix_account                                                               |  
|				                       | SKIP_GATEWAY_AWS                   |		    + AWS_VPC_ID, AWS_REGION, AWS_SUBNET, AWS_GW_SIZE (optional)         |
|                                      | SKIP_GATEWAY_GCP                   |         + GCP_VPC_ID, GCP_ZONE, GCP_SUBNET, GCP_GW_SIZE (optional)             |
|                                      | SKIP_GATEWAY_AZURE                 |         + AZURE_VNET_ID, AZURE_REGION, AZURE_SUBNET, AZURE_GW_SIZE             |
|                                      | SKIP_GATEWAY_OCI                   |         + OCI_VPC_ID, OCI_REGION, OCI_SUBNET, OCI_GW_SIZE(optional)            |
|                                      | SKIP_GATEWAY_AWSGOV                |         + AWSGOV_VPC_ID, AWSGOV_REGION, AWSGOV_SUBNET, AWSGOV_GW_SIZE(optional)|
| aviatrix_gateway_dnat                | SKIP_GATEWAY_DNAT                  | aviatrix_account                                                               |
|				                       | SKIP_GATEWAY_DNAT_AWS              |		    + AWS_VPC_ID, AWS_REGION, AWS_SUBNET, AWS_GW_SIZE (optional)         |
|                                      | SKIP_GATEWAY_DNAT_AZURE            |         + AZURE_VNET_ID, AZURE_REGION, AZURE_SUBNET, AZURE_GW_SIZE             |
| aviatrix_gateway_snat                | SKIP_GATEWAY_SNAT                  | aviatrix_account                                                               |
|				                       | SKIP_GATEWAY_SNAT_AWS              |		    + AWS_VPC_ID, AWS_REGION, AWS_SUBNET, AWS_GW_SIZE (optional)         |
|                                      | SKIP_GATEWAY_SNAT_AZURE            |         + AZURE_VNET_ID, AZURE_REGION, AZURE_SUBNET, AZURE_GW_SIZE             |
| aviatrix_geo_vpn                     | SKIP_GEO_VPN                       | aviatrix_account + DOMAIN_NAME + AWS_VPC_ID, AWS_REGION, AWS_SUBNET            |
|                                      |                                    |                                + AWS_VPC_ID2, AWS_REGION2, AWS_SUBNET2         |
| aviatrix_netflow_agent               | SKIP_NETFLOW_AGENT                 | N/A
| aviatrix_periodic_ping               | SKIP_PERIODIC_PING                 | aviatrix_gateway                                                               |
| aviatrix_rbac_group                  | SKIP_RBAC_GROUP                    | N/A                                                                            |
| aviatrix_rbac_group_access_account_attachment | SKIP_RBAC_GROUP_ACCESS_ACCOUNT_ATTACHMENT | aviatrix_account                                               |
| aviatrix_rbac_group_permission_attachment | SKIP_RBAC_GROUP_PERMISSION_ATTACHMENT | N/A                                                                    |
| aviatrix_rbac_group_user_attachment  | SKIP_RBAC_GROUP_USER_ATTACHMENT    | aviatrix_account_user                                                          |
| aviatrix_remote_syslog               | SKIP_REMOTE_SYSLOG                 | N/A                                                                            |
| aviatrix_saml_endpoint               | SKIP_SAML_ENDPOINT                 | IDP_METADATA, IDP_METADATA_TYPE                                                |
| aviatrix_segmentation_security_domain | SKIP_SEGMENTATION_SECURITY_DOMAIN | N/A                                                                            |
| aviatrix_segmentation_security_domain_association | SKIP_SEGMENTATION_SECURITY_DOMAIN_ASSOCIATION | aviatrix_gateway + AWS_VPC_ID2, AWS_REGION2, AWS_SUBNET2 |
| aviatrix_segmentation_security_domain_connection_policy | SKIP_SEGMENTATION_SECURITY_DOMAIN_CONNECTION_POLICY | N/A                                        |
| aviatrix_site2cloud                  | SKIP_S2C                           | aviatrix_gateway                                                               |
| aviatrix_splunk_logging              | SKIP_SPLUNK_LOGGING                | N/A                                                                            |
| aviatrix_spoke_gateway               | SKIP_SPOKE_GATEWAY                 | aviatrix_gateway                                                               |
|                                      | SKIP_SPOKE_GATEWAY_AWS             |         + AWS_VPC_ID, AWS_REGION, AWS_SUBNET, AWS_GW_SIZE (optional)           |
|                                      | SKIP_SPOKE_GATEWAY_GCP             |         + GCP_VPC_ID, GCP_ZONE, GCP_SUBNET, GCP_GW_SIZE (optional)             |
|                                      | SKIP_SPOKE_GATEWAY_AZURE           |         + AZURE_VNET_ID, AZURE_REGION, AZURE_SUBNET, AZURE_GW_SIZE             |
|                                      | SKIP_SPOKE_GATEWAY_OCI             |         + OCI_VPC_ID, OCI_REGION, OCI_SUBNET, OCI_GW_SIZE(optional)            |
| aviatrix_spoke_vpc                   | SKIP_SPOKE                         | aviatrix_gateway                                                               |
|                                      | SKIP_SPOKE_AWS                     |         + AWS_VPC_ID, AWS_REGION, AWS_SUBNET, AWS_GW_SIZE (optional)           | 
|                                      | SKIP_SPOKE_GCP                     |         + GCP_VPC_ID, GCP_ZONE, GCP_SUBNET, GCP_GW_SIZE (optional)             |
|                                      | SKIP_SPOKE_AZURE                   |         + AZURE_VNET_ID, AZURE_REGION, AZURE_SUBNET, AZURE_GW_SIZE             |
| aviatrix_sumologic_forwarder         | SKIP_SUMOLOGIC_FORWARDER           | N/A
| aviatrix_trans_peer                  | SKIP_TRANS_PEER                    | aviatrix_tunnel                                                                |
| aviatrix_transit_external_device_conn| SKIP_TRANSIT_EXTERNAL_DEVICE_CONN  | aviatrix_account + aviatrix_transit_gateway                                    |
| aviatrix_transit_firenet_policy      | SKIP_TRANSIT_FIRENET_POLICY        | aviatrix_account                                                               |
|                                      | SKIP_TRANSIT_FIRENET_POLICY_AWS    |       + aviatrix_transit_gateway + aviatrix_spoke_gateway in AWS               |
|                                      | SKIP_TRANSIT_FIRENET_POLICY_AZURE  |       + aviatrix_transit_gateway + aviatrix_spoke_gateway in AZURE             |
| aviatrix_transit_gateway             | SKIP_TRANSIT_GATEWAY               | aviatrix_gateway                                                               |
|                                      | SKIP_TRANSIT_GATEWAY_AWS           | aviatrix_gateway in AWS                                                        |
|                                      | SKIP_TRANSIT_GATEWAY_AZURE         | aviatrix_gateway in AZURE                                                      |
|                                      | SKIP_GATEWAY_GCP                   | aviatrix_gateway in GCP                                                        |
|                                      | SKIP_GATEWAY_OCI                   | aviatrix_gateway in OCI                                                        |
| aviatrix_transit_vpc                 | SKIP_TRANSIT                       | aviatrix_gateway                                                               |
|                                      | SKIP_TRANSIT_AWS                   | aviatrix_gateway in AWS                                                        |
|                                      | SKIP_TRANSIT_AZURE                 | aviatrix_gateway in AZURE                                                      |  
| aviatrix_transit_gateway_peering     | SKIP_TRANSIT_GATEWAY_PEERING       | aviatrix_gateway + AWS_VPC_ID2, AWS_REGION2, AWS_SUBNET2                       |
| aviatrix_tunnel                      | SKIP_TUNNEL                        | aviatrix_gateway + AWS_VPC_ID2, AWS_REGION2, AWS_SUBNET2                       |
| aviatrix_version                     | SKIP_VERSION                       |                                                                                |
| aviatrix_vgw_conn                    | SKIP_VGW_CONN                      | aviatrix_gateway + AWS_BGP_VGW_ID                                              |
| aviatrix_vpc                         | SKIP_VPC                           | aviatrix_account                                                               |
| aviatrix_vpn_profile                 | SKIP_VPN_PROFILE                   | aviatrix_vpn_user                                                              |
| aviatrix_vpn_user                    | SKIP_VPN_USER                      | aviatrix_gateway                                                               |
| aviatrix_vpn_user_accelerator	       | SKIP_VPN_USER_ACCELERATOR          | aviatrix_gateway						                                         |
| aviatrix_data_source_account         | SKIP_DATA_ACCOUNT                  | aviatrix_account                                                               |
| aviatrix_data_source_caller_identity | SKIP_DATA_CALLER_IDENTITY          |                                                                                |
| aviatrix_data_source_firenet         | SKIP_DATA_FIRENET                  | aviatrix_firenet                                                               |
| aviatrix_data_source_firenet_vendor_integration | SKIP_DATA_FIRENET_VENDOR_INTEGRATION    | aviatrix_account + AWS_REGION, Palo Alto VM series             |
| aviatrix_data_source_firewall        | SKIP_DATA_FIREWALL                 | aviatrix_gateway                                                               |
| aviatrix_data_source_gateway         | SKIP_DATA_GATEWAY                  | aviatrix_gateway                                                               |
| aviatrix_data_source_spoke_gateway   | SKIP_DATA_SPOKE_GATEWAY            | aviatrix_spoke_gateway                                                         |
| aviatrix_data_source_transit_gateway | SKIP_DATA_TRANSIT_GATEWAY          | aviatrix_transit_gateway                                                       |
| aviatrix_data_source_vpc             | SKIP_DATA_VPC                      | aviatrix_vpc                                                                   |
| aviatrix_data_source_vpc_tracker     | SKIP_DATA_VPC_TRACKER              | aviatrix_vpc                                                           |
