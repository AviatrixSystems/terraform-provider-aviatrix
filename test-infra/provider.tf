provider "aws" {
  region          = "${var.aws_region3}"
  access_key      = "${var.aws_access_key}"
  secret_key      = "${var.aws_secret_key}"
}
provider "google" {
  credentials     = "${file("${var.gcp_credentials_file_path}")}"
  region          = "${var.gcp_region1}"
  project         = "${var.gcp_project_id1}"
}
provider "azurerm" {
  subscription_id = "${var.azure_subscription_id}"
  tenant_id       = "${var.azure_tenant_id}"
  client_id       = "${var.azure_client_id}"
  client_secret   = "${var.azure_client_secret}"
}
#provider "aviatrix" {
#  controller_ip   = "${var.aviatrix_controller_ip}"
#  username        = "${var.aviatrix_username}"
#  password        = "${var.aviatrix_password}"
#}
provider "aviatrix" {
    username      = "admin"
    password      = "${module.aviatrix-controller-build.private_ip}"
    controller_ip = "${module.aviatrix-controller-build.public_ip}"
}

