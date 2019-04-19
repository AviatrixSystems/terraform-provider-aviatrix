# Acceptance Tests

#### Pre-requisites

- The controller must be launched before hand and must be up and running the latest controller version
- IAM roles (aviatrix-role-ec2 and aviatrix-role-app) also must be created and attached if any IAM role related tests are to be run. Currently all tests are based on Access key, Secret key
- The VPC's with public subnet to launch the gateways must be created before the tests
- If you are running aviatrix_aws_peer or aviatrix_peer, two VPC's with non overlapping CIDR's must be created before hand
- If you are running the tests on a BYOL controller, the customer ID must be set prior to the tests, otherwise run the tests on a PayG metered controller
- aviatrix_aws_tgw test only allows Transit GWs and VPCs to be attached to the TGW in the same region 

#### Skip parameters and variables

Passing an environment value of "yes" to the skip parameter allows you to skip the particular resource. If it is not skipped, it checks for the existence of other required variables. Generic variables are required for any acceptance test

| Test module name                     | Skip parameter               | Required variables                                                  |
| ------------------------------------ | ---------------------------- | ------------------------------------------------------------------- |
| Generic                              | N/A                          | AVIATRIX_USERNAME, AVIATRIX_PASSWORD, AVIATRIX_CONTROLLER_IP        |
| aviatrix_account                     | SKIP_ACCOUNT                 | AWS_ACCOUNT_NUMBER, AWS_ACCESS_KEY, AWS_SECRET_KEY                  |
|		                       | SKIP_AWS_ACCOUNT	      | GCP_ID, GCP_CREDENTIALS_FILEPATH,                                   |
|                     		       | SKIP_GCP_ACCOUNT	      | ARM_SUBSCRIPTION_ID, ARM_DIRECTORY_ID, ARM_APPLICATION_ID,	    |
|		                       | SKIP_ARM_ACCOUNT	      | ARM_APPLICATION_KEY						    |	
| aviatrix_account_user                | SKIP_ACCOUNT_USER            |                                                                     |
| aviatrix_admin_email                 | SKIP_ADMIN_EMAIL             |                                                                     |
| aviatrix_aws_peer                    | SKIP_AWS_PEER                | aviatrix_account + AWS_VPC_ID, AWS_VPC_ID2, AWS_REGION, AWS_REGION2 |
| aviatrix_aws_tgw                     | SKIP_AWS_TGW                 | aviatrix_account + AWS_VPC_ID, AWS_REGION, AWS_VPC_TGW_ID           |
| aviatrix_aws_tgw_vpc_attachment      | SKIP_AWS_TGW_VPC_ATTACHMENT  | aviatrix_account + AWS_VPC_ID, AWS_REGION, AWS_VPC_TGW_ID           |
| aviatrix_customer_id                 | SKIP_CUSTOMER_ID             | CUSTOMER_ID                                                         |
| aviatrix_firewall                    | SKIP_FIREWALL                | aviatrix_gateway                                                    |
| aviatrix_firewall_tag                | SKIP_FIREWALL_TAG            |                                                                     |
| aviatrix_fqdn                        | SKIP_FQDN                    | aviatrix_gateway                                                    |
| aviatrix_gateway                     | SKIP_GATEWAY                 | aviatrix_account + AWS_VPC_ID, AWS_REGION, AWS_VPC_NET              |
| aviatrix_site2cloud                  | SKIP_S2C                     | aviatrix_gateway                                                    |
| aviatrix_spoke_vpc                   | SKIP_SPOKE                   | aviatrix_gateway + GCP_VPC_ID, GCP_ZONE, GCP_SUBNET,		    |
|				       | SKIP_AWS_SPOKE		      |                   ARM_VNET_ID, ARM_REGION, ARM_SUBNET		    |
|				       | SKIP_GCP_SPOKE		      |									    |
|				       | SKIP_ARM_SPOKE		      |									    |
| aviatrix_trans_peer                  | SKIP_TRANS_PEER              | aviatrix_tunnel                                                     |
| aviatrix_transit_vpc                 | SKIP_TRANSIT                 | aviatrix_gateway                                                    |
|                                      | SKIP_TRANSIT_AWS             | aviatrix_gateway in AWS                                             |
|                                      | SKIP_TRANSIT_AZURE           | aviatrix_gateway in AZURE                                           |
| aviatrix_transit_gateway_peering     | SKIP_TRANSIT_GATEWAY_PEERING | aviatrix_gateway + AWS_VPC_ID2, AWS_REGION2, AWS_VPC_NET2           |
| aviatrix_tunnel                      | SKIP_TUNNEL                  | aviatrix_gateway + AWS_VPC_ID2, AWS_REGION2, AWS_VPC_NET2           |
| aviatrix_version                     | SKIP_VERSION                 |                                                                     |
| aviatrix_vgw_conn                    | SKIP_VGW_CONN                | aviatrix_gateway + AWS_BGP_VGW_ID                                   |
| aviatrix_vpn_profile                 | SKIP_VPN_PROFILE             | aviatrix_vpn_user                                                   |
| aviatrix_vpn_user                    | SKIP_VPN_USER                | aviatrix_gateway                                                    |
| aviatrix_data_source_account         | SKIP_DATA_ACCOUNT            | aviatrix_account                                                    |
| aviatrix_data_source_caller_identity | SKIP_DATA_CALLER_IDENTITY    |                                                                     |
| aviatrix_data_source_gateway         | SKIP_DATA_GATEWAY            | aviatrix_account + AWS_VPC_ID, AWS_REGION, AWS_VPC_NET              |

