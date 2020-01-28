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

| Test module name                     | Skip parameter               | Required variables                                                    |
| ------------------------------------ | ---------------------------- | --------------------------------------------------------------------- |
| Generic                              | N/A                          | AVIATRIX_USERNAME, AVIATRIX_PASSWORD, AVIATRIX_CONTROLLER_IP          |
| aviatrix_account                     | SKIP_ACCOUNT                 |                                                                       |
|		                               | SKIP_ACCOUNT_AWS	          | AWS_ACCOUNT_NUMBER, AWS_ACCESS_KEY, AWS_SECRET_KEY                    |
|                     		           | SKIP_ACCOUNT_GCP	          | GCP_ID, GCP_CREDENTIALS_FILEPATH	                                  |
|		                               | SKIP_ACCOUNT_ARM	          | ARM_SUBSCRIPTION_ID, ARM_DIRECTORY_ID, ARM_APPLICATION_ID, ARM_APPLICATION_KEY |
|                     		           | SKIP_ACCOUNT_OCI	          | OCI_TENANCY_ID, OCI_USER_ID, OCI_COMPARTMENT_ID, OCI_API_KEY_FILEPATH |
|		                               | SKIP_ACCOUNT_AWSGOV          | AWSGOV_ACCOUNT_NUMBER, AWSGOV_ACCESS_KEY, AWSGOV_SECRET_KEY           |
| aviatrix_account_user                | SKIP_ACCOUNT_USER            |                                                                       |
| aviatrix_arm_peer                    | SKIP_ARM_PEER                | aviatrix_account + ARM_VNET_ID, ARM_VNET_ID2, ARM_REGION, ARM_REGION2 |
| aviatrix_aws_peer                    | SKIP_AWS_PEER                | aviatrix_account + AWS_VPC_ID, AWS_VPC_ID2, AWS_REGION, AWS_REGION2   |
| aviatrix_aws_tgw                     | SKIP_AWS_TGW                 | aviatrix_account + AWS_VPC_ID, AWS_REGION, AWS_VPC_TGW_ID             |
| aviatrix_aws_tgw_directconnect       | SKIP_AWS_TGW_DIRECTCONNECT   | aviatrix_aws_tgw + AWS_DX_GATEWAY                                     |
| aviatrix_aws_tgw_vpc_attachment      | SKIP_AWS_TGW_VPC_ATTACHMENT  | aviatrix_aws_tgw                                                      |
| aviatrix_aws_tgw_vpn_conn            | SKIP_AWS_TGW_VPN_CONN        | aviatrix_aws_tgw                                                      |
| aviatrix_controller_config           | SKIP_CONTROLLER_CONFIG       | aviatrix_account                                                      |
| aviatrix_firenet                     | SKIP_FIRENET                 | aviatrix_account + AWS_REGION, Palo Alto VM series                    |
| aviatrix_firewall                    | SKIP_FIREWALL                | aviatrix_gateway                                                      |
| aviatrix_firewall_instance           | SKIP_FIREWALL_INSTANCE       | aviatrix_account + AWS_REGION, Palo Alto VM series                    |
| aviatrix_firewall_tag                | SKIP_FIREWALL_TAG            |                                                                       |
| aviatrix_fqdn                        | SKIP_FQDN                    | aviatrix_gateway                                                      |
| aviatrix_gateway                     | SKIP_GATEWAY                 | aviatrix_account                                                      |
|				                       | SKIP_GATEWAY_AWS             |		    + AWS_VPC_ID, AWS_REGION, AWS_SUBNET, AWS_GW_SIZE (optional)  |
|                                      | SKIP_GATEWAY_GCP             |         + GCP_VPC_ID, GCP_ZONE, GCP_SUBNET, GCP_GW_SIZE (optional)    |
|                                      | SKIP_GATEWAY_ARM             |         + ARM_VNET_ID, ARM_REGION, ARM_SUBNET, ARM_GW_SIZE            |
|                                      | SKIP_GATEWAY_OCI             |         + OCI_VPC_ID, OCI_REGION, OCI_SUBNET, OCI_GW_SIZE(optional)   |
|                                      | SKIP_GATEWAY_AWSGOV          |         + AWSGOV_VPC_ID, AWSGOV_REGION, AWSGOV_SUBNET, AWSGOV_GW_SIZE(optional)   |
| aviatrix_gateway_dnat                | SKIP_GATEWAY_DNAT            | aviatrix_account                                                      |
|				                       | SKIP_GATEWAY_DNAT_AWS        |		    + AWS_VPC_ID, AWS_REGION, AWS_SUBNET, AWS_GW_SIZE (optional)  |
|                                      | SKIP_GATEWAY_DNAT_ARM        |         + ARM_VNET_ID, ARM_REGION, ARM_SUBNET, ARM_GW_SIZE            |
| aviatrix_geo_vpn                     | SKIP_GEO_VPN                 | aviatrix_account + DOMAIN_NAME + AWS_VPC_ID, AWS_REGION, AWS_SUBNET   |
|                                      |                              |                                + AWS_VPC_ID2, AWS_REGION2, AWS_SUBNET2|
| aviatrix_saml_endpoint               | SKIP_SAML_ENDPOINT           | IDP_METADATA, IDP_METADATA_TYPE                                       |
| aviatrix_site2cloud                  | SKIP_S2C                     | aviatrix_gateway                                                      |
| aviatrix_spoke_gateway               | SKIP_SPOKE_GATEWAY           | aviatrix_gateway                                                      |
|                                      | SKIP_SPOKE_GATEWAY_AWS       |         + AWS_VPC_ID, AWS_REGION, AWS_SUBNET, AWS_GW_SIZE (optional)  |
|                                      | SKIP_SPOKE_GATEWAY_GCP       |         + GCP_VPC_ID, GCP_ZONE, GCP_SUBNET, GCP_GW_SIZE (optional)    |
|                                      | SKIP_SPOKE_GATEWAY_ARM       |         + ARM_VNET_ID, ARM_REGION, ARM_SUBNET, ARM_GW_SIZE            |
|                                      | SKIP_SPOKE_GATEWAY_OCI       |         + OCI_VPC_ID, OCI_REGION, OCI_SUBNET, OCI_GW_SIZE(optional)   |
| aviatrix_spoke_vpc                   | SKIP_SPOKE                   | aviatrix_gateway                                                      |
|                                      | SKIP_SPOKE_AWS               |         + AWS_VPC_ID, AWS_REGION, AWS_SUBNET, AWS_GW_SIZE (optional)  |
|                                      | SKIP_SPOKE_GCP               |         + GCP_VPC_ID, GCP_ZONE, GCP_SUBNET, GCP_GW_SIZE (optional)    |
|                                      | SKIP_SPOKE_ARM               |         + ARM_VNET_ID, ARM_REGION, ARM_SUBNET, ARM_GW_SIZE            |
| aviatrix_trans_peer                  | SKIP_TRANS_PEER              | aviatrix_tunnel                                                       |
| aviatrix_transit_gateway             | SKIP_TRANSIT_GATEWAY         | aviatrix_gateway                                                      |
|                                      | SKIP_TRANSIT_GATEWAY_AWS     | aviatrix_gateway in AWS                                               |
|                                      | SKIP_TRANSIT_GATEWAY_ARM     | aviatrix_gateway in ARM                                               |
|                                      | SKIP_GATEWAY_GCP             | aviatrix_gateway in GCP                                               |
|                                      | SKIP_GATEWAY_OCI             | aviatrix_gateway in OCI                                               |
| aviatrix_transit_vpc                 | SKIP_TRANSIT                 | aviatrix_gateway                                                      |
|                                      | SKIP_TRANSIT_AWS             | aviatrix_gateway in AWS                                               |
|                                      | SKIP_TRANSIT_ARM             | aviatrix_gateway in ARM                                               |
| aviatrix_transit_gateway_peering     | SKIP_TRANSIT_GATEWAY_PEERING | aviatrix_gateway + AWS_VPC_ID2, AWS_REGION2, AWS_SUBNET2              |
| aviatrix_tunnel                      | SKIP_TUNNEL                  | aviatrix_gateway + AWS_VPC_ID2, AWS_REGION2, AWS_SUBNET2              |
| aviatrix_version                     | SKIP_VERSION                 |                                                                       |
| aviatrix_vgw_conn                    | SKIP_VGW_CONN                | aviatrix_gateway + AWS_BGP_VGW_ID                                     |
| aviatrix_vpc                         | SKIP_VPC                     | aviatrix_account                                                      |
| aviatrix_vpn_profile                 | SKIP_VPN_PROFILE             | aviatrix_vpn_user                                                     |
| aviatrix_vpn_user                    | SKIP_VPN_USER                | aviatrix_gateway                                                      |
| aviatrix_vpn_user_accelerator	       | SKIP_VPN_USER_ACCELERATOR    | aviatrix_gateway						                              |
| aviatrix_data_source_account         | SKIP_DATA_ACCOUNT            | aviatrix_account                                                      |
| aviatrix_data_source_caller_identity | SKIP_DATA_CALLER_IDENTITY    |                                                                       |
| aviatrix_data_source_firenet_vendor_integration | SKIP_DATA_FIRENET_VENDOR_INTEGRATION    | aviatrix_account + AWS_REGION, Palo Alto VM series |
| aviatrix_data_source_gateway         | SKIP_DATA_GATEWAY            | aviatrix_gateway                                                      |
