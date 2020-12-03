data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
data "aws_caller_identity" "current_awsgov" {
  provider = aws.gov
}
data "aws_region" "current_awsgov" {
  provider = aws.gov
}
# This is not part of this role should not be destroyed
# module "aviatrix-iam-roles" {
#   source            = "github.com/AviatrixSystems/terraform-modules.git//aviatrix-controller-iam-roles?ref=terraform_0.12"
#   master-account-id = "${data.aws_caller_identity.current.account_id}"
# }

module "aviatrix_controller_vpc" {
  source         = "./aws"
  aws_vpc_cidr   = var.aws_vpc_cidr2
  aws_vpc_subnet = var.aws_vpc_subnet2
  aws_region     = var.aws_region2
}
module "aviatrix-controller-build" {
  source                 = "github.com/AviatrixSystems/terraform-modules.git//aviatrix-controller-build?ref=terraform_0.12"
  vpc                    = module.aviatrix_controller_vpc.vpc
  subnet                 = module.aviatrix_controller_vpc.subnet_id
  keypair                = var.keypair
  #ec2role                = module.aviatrix-iam-roles.aviatrix-role-ec2-name  # This can be used from the module aviatrix-iam-roles above, but since it cannot be deleted, it is harcoded
  ec2role                = "aviatrix-role-ec2"
  termination_protection = false
  type                   = var.type
}
module "aviatrix-controller-initialize" {
  source              = "github.com/AviatrixSystems/terraform-modules.git//aviatrix-controller-initialize?ref=terraform_0.12"
  admin_password      = var.admin_password
  admin_email         = var.admin_email
  private_ip          = module.aviatrix-controller-build.private_ip
  public_ip           = module.aviatrix-controller-build.public_ip
  access_account_name = var.access_account_name
  aws_account_id      = data.aws_caller_identity.current.account_id
  customer_license_id = var.customer_id
  vpc_id              = module.aviatrix-controller-build.vpc_id
  subnet_id           = module.aviatrix-controller-build.subnet_id
}

output "result" {
  value = module.aviatrix-controller-initialize.result
}
output "controller_private_ip" {
  value = module.aviatrix-controller-build.private_ip
}

module "aviatrix_gcp_vpc1" {
  source         = "./gcp"
  gcp_project_id = var.gcp_project_id1
  gcp_region     = var.gcp_region1
  gcp_vpc_cidr   = var.gcp_vpc_cidr1
}
module "aviatrix_azure_vpc1" {
  source               = "./azure"
  name                 = "vpc1"
  azure_region         = var.azure_region1
  azure_vpc_cidr       = var.azure_vpc_cidr1
  azure_vpc_subnet     = var.azure_vpc_subnet1
  azure_address_prefix = var.azure_vpc_subnet3
}
module "aviatrix_azure_vpc2" {
  source               = "./azure"
  name                 = "vpc2"
  azure_region         = var.azure_region2
  azure_vpc_cidr       = var.azure_vpc_cidr2
  azure_vpc_subnet     = var.azure_vpc_subnet2
  azure_address_prefix = var.azure_vpc_subnet4
}
module "aviatrix_aws_vpc1" {
  source         = "./aws"
  aws_vpc_cidr   = var.aws_vpc_cidr1
  aws_vpc_subnet = var.aws_vpc_subnet1
  aws_region     = var.aws_region1
}
module "aviatrix_aws_vpc2" {
  source         = "./aws"
  aws_vpc_cidr   = var.aws_vpc_cidr2
  aws_vpc_subnet = var.aws_vpc_subnet2
  aws_region     = var.aws_region2
}
module "aviatrix_awsgov_vpc" {
  providers      = {
    aws = aws.gov
  }
  source         = "./aws"
  aws_vpc_cidr   = var.awsgov_vpc_cidr1
  aws_vpc_subnet = var.awsgov_vpc_subnet1
  aws_region     = var.awsgov_region1
}
module "aviatrix_oci_vpc1" {
  source             = "./oci"
  oci_compartment_id = var.oci_compartment_id
  oci_vpc_cidr1      = var.oci_vpc_cidr1
}
resource "aws_vpn_gateway" "vgw" {
  vpc_id = module.aviatrix_aws_vpc2.vpc
  tags   = {
    Name = "aviatrix-vgw"
  }
}
resource "aws_dx_gateway" "dx-gateway" {
  name            = "aws-dx-gateway"
  amazon_side_asn = "64512"
}

resource "aviatrix_transit_gateway" "cwan-transitgw" {
  cloud_type   = 1
  account_name = var.access_account_name
  gw_name      = "cwan-transitgw"
  vpc_id       = module.aviatrix_aws_vpc1.vpc
  vpc_reg      = data.aws_region.current.name
  gw_size      = "t2.micro"
  subnet       = module.aviatrix_aws_vpc1.subnet
}

resource "aviatrix_aws_tgw" "cwan-awstgw" {
  account_name                      = var.access_account_name
  aws_side_as_number                = "65011"
  region                            = data.aws_region.current.name
  tgw_name                          = "cwan-awstgw"

  security_domains {
    connected_domains    = [
      "Default_Domain",
      "Shared_Service_Domain",
    ]
    security_domain_name = "Aviatrix_Edge_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain",
      "Shared_Service_Domain"
    ]
    security_domain_name = "Default_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain",
      "Default_Domain"
    ]
    security_domain_name = "Shared_Service_Domain"
  }
}

module "azure-vwan" {
  source       = "./azure-vwan"
  azure_region = var.azure_region1
}

module "cisco-csr" {
  source     = "./cisco-csr"
  aws_region = var.aws_region1
}

module "azure-vng" {
  source            = "./azure-vng"
  azure_region      = var.azure_region1
  azure_vpc_cidr    = var.azure_vpc_cidr1
  azure_vpc_subnet1 = var.azure_vpc_subnet1
  azure_vpc_subnet2 = var.azure_vpc_subnet3
}
