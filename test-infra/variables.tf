variable "gcp_credentials_file_path" {
  description = "Location of the GCP credentials to use."
}

variable "gcp_project_id1" {}
variable "gcp_region1" {default="us-west1"}
variable "gcp_vpc_cidr1" {default="10.50.0.0/16"}
variable "gcp_zone1" {default="us-west1-b"}

variable "azure_region1" {default="West US"}
variable "azure_region2" {default="West US"}
variable "azure_version" {default="1.30.1"}
variable "azure_vpc_cidr1" {default="10.30.0.0/16"}
variable "azure_vpc_cidr2" {default="10.40.0.0/16"}
variable "azure_vpc_subnet1" {default="10.30.0.0/24"}
variable "azure_vpc_subnet2" {default="10.40.0.0/24"}
variable "azure_vpc_subnet3" {default="10.30.1.0/24"}
variable "azure_vpc_subnet4" {default="10.40.1.0/24"}
variable "azure_gw_size" {default="Standard_B2ms"}
variable "azure_subscription_id" {}
variable "azure_tenant_id" {}
variable "azure_client_id" {}
variable "azure_client_secret" {}

variable "aws_region1" {default="us-west-1"}
variable "aws_region2" {default="us-west-1"}
variable "aws_region3" {default="us-west-1"}
variable "aws_vpc_cidr1" {default="10.10.0.0/16"}
variable "aws_vpc_cidr2" {default="10.20.0.0/16"}
variable "aws_vpc_cidr3" {default="192.168.0.0.0/16"}
variable "aws_vpc_subnet1" {default="10.10.0.0/24"}
variable "aws_vpc_subnet2" {default="10.20.0.0/24"}
variable "aws_vpc_subnet3" {default="192.168.0.0/24"}
variable "aws_access_key" {}
variable "aws_secret_key" {}
variable "domain_name" {}

variable "awsgov_region1" {default="us-gov-west-1"}
variable "awsgov_vpc_cidr1" {default="10.10.0.0/16"}
variable "awsgov_vpc_subnet1" {default="10.10.0.0/24"}
variable "awsgov_access_key" {}
variable "awsgov_secret_key" {}
variable "awsgov_account_number" {}

variable "oci_tenancy_id" {}
variable "oci_user_id" {}
variable "oci_compartment_id" {}
variable "oci_api_key_filepath" {
  description = "Location of the OCI credentials to use."
}
variable "oci_region1" {default="us-ashburn-1"}
variable "oci_vpc_cidr1" {default="10.10.0.0/16"}
variable "oci_subnet1" {default="10.10.0.0/16"}
variable "oci_fingerprint" {}

variable "IDP_METADATA" {default="ThisIsATest"}
variable "IDP_METADATA_TYPE" {default="Text"}

variable "admin_password" {default="Aviatrix123#"}
variable "admin_email" {default="abc@xyz.com"}
variable "access_account_name" {default="aws_init_acc"}
variable "keypair" {default="aviatrix-key"}
variable "type" {default="metered"}
variable "customer_id" {default = ""}

variable "datadog_api_key" {}
