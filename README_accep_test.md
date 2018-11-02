# Acceptance Tests

#### Pre-requisites

- The controller must be launched before hand and must be up and running. 
- IAM roles (aviatrix-role-ec2 and aviatrix-role-app) also must be created and attached if any IAM role related tests are to be run. Currently all tests are based on Access key, Secret key
- The VPC's with public subnet to launch the gateways must be created before the tests.
- If you are running aviatrix_aws_peer or aviatrix_peer, two VPC's with non overlapping CIDR's must be created before hand
- If you are running the tests on a BYOL controller, the customer ID must be set prior to the tests, otherwise the tests them on a PayG metered controller.

#### Skip parameters and variables

Passing an environment value of "yes" to the skip parameter allows you to skip the particular resource. If it is not skipped, it checks for the existence of other required variables. Generic variables are required for any acceptance test

| Test module name      | Skip parameter    | Required variables                                           |
| --------------------- | ----------------- | ------------------------------------------------------------ |
| Generic               | N/A               | AVIATRIX_USERNAME, AVIATRIX_PASSWORD, AVIATRIX_CONTROLLER_IP |
| aviatrix_account      | SKIP_ACCOUNT      | AWS_ACCOUNT_NUMBER, AWS_ACCESS_KEY, AWS_SECRET_KEY           |
| aviatrix_account_user | SKIP_ACCOUNT_USER |                                                              |
| aviatrix_aws_peer     | SKIP_AWS_PEER     | aviatrix_account+AWS_VPC_ID, AWS_VPC_ID2, AWS_REGION, AWS_REGION2 |
| aviatrix_gateway      | SKIP_GATEWAY      | aviatrix_account+AWS_VPC_ID2, AWS_REGION, AWS_VPC_NET        |


