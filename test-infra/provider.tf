provider "aws" {
  region     = var.aws_region3
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
  alias      = "reg"
}

provider "aws" {
  region     = var.awsgov_region1
  access_key = var.awsgov_access_key
  secret_key = var.awsgov_secret_key
  alias      = "gov"
}

provider "google" {
  credentials = "${file(var.gcp_credentials_file_path)}"
  region      = var.gcp_region1
  project     = var.gcp_project_id1
}

provider "azurerm" {
  subscription_id = var.azure_subscription_id
  tenant_id       = var.azure_tenant_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
}

#provider "aviatrix" {
#  controller_ip = var.aviatrix_controller_ip
#  username      = var.aviatrix_username
#  password      = var.aviatrix_password
#}

provider "oci" {
  tenancy_ocid     = var.oci_tenancy_id
  user_ocid        = var.oci_user_id
  fingerprint      = var.oci_fingerprint
  private_key_path = var.oci_api_key_filepath
  region           = var.oci_region1
}
data "oci_identity_availability_domains" "ads" {
  compartment_id = var.oci_compartment_id
}

provider "aviatrix" {
  username      = "admin"
  password      = module.aviatrix-controller-build.private_ip
  controller_ip = module.aviatrix-controller-build.public_ip
}