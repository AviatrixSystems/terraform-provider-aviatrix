data "aws_caller_identity" "current_aws" {provider = aws.reg}
data "aws_region" "current_aws" {provider = aws.reg}

data "aws_caller_identity" "current_awsgov" {provider = aws.gov}
data "aws_region" "current_awsgov" {provider = aws.gov}
# This is not part of this role should not be destroyed
# module "aviatrix-iam-roles" {
#   source            = "github.com/AviatrixSystems/terraform-modules.git//aviatrix-controller-iam-roles?ref=terraform_0.12"
#   master-account-id = "${data.aws_caller_identity.current.account_id}"
# }

module "aviatrix_controller_vpc" {
  providers {
    aws = aws.reg
  }
  source         = "./aws"
  aws_vpc_cidr   = var.aws_vpc_cidr2
  aws_vpc_subnet = var.aws_vpc_subnet2
  aws_region     = var.aws_region2
}
module "aviatrix-controller-build" {
  providers {
    aws = aws.reg
  }
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
  providers {
    aws = aws.reg
  }
  source              = "github.com/AviatrixSystems/terraform-modules.git//aviatrix-controller-initialize?ref=terraform_0.12"
  admin_password      = var.admin_password
  admin_email         = var.admin_email
  private_ip          = module.aviatrix-controller-build.private_ip
  public_ip           = module.aviatrix-controller-build.public_ip
  access_account_name = var.access_account_name
  aws_account_id      = data.aws_caller_identity.current.account_id
  customer_license_id = var.customer_id
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
module "aviatrix_arm_vpc1" {
  source           = "./arm"
  name             = "vpc1"
  azure_region     = var.azure_region1
  azure_vpc_cidr   = var.azure_vpc_cidr1
  azure_vpc_subnet = var.azure_vpc_subnet1
}
module "aviatrix_arm_vpc2" {
  source           = "./arm"
  name             = "vpc2"
  azure_region     = var.azure_region2
  azure_vpc_cidr   = var.azure_vpc_cidr2
  azure_vpc_subnet = var.azure_vpc_subnet2
}
module "aviatrix_aws_vpc1" {
  providers {
    aws = aws.reg
  }
  source         = "./aws"
  aws_vpc_cidr   = var.aws_vpc_cidr1
  aws_vpc_subnet = var.aws_vpc_subnet1
  aws_region     = var.aws_region1
}
module "aviatrix_aws_vpc2" {
  providers {
    aws = aws.reg
  }
  source         = "./aws"
  aws_vpc_cidr   = var.aws_vpc_cidr2
  aws_vpc_subnet = var.aws_vpc_subnet2
  aws_region     = var.aws_region2
}
module "aviatrix_awsgov_vpc" {
  providers {
    aws = aws.gov
  }
  source         = "./aws"
  aws_vpc_cidr   = var.awsgov_vpc_cidr1
  aws_vpc_subnet = var.awsgov_vpc_subnet1
  aws_region     = var.awsgov_region1
}
module "aviatrix_oci_vpc1" {
  source         = "./oci"
  oci_compartment_id = var.oci_compartment_id
  oci_vpc_cidr1 = var.oci_vpc_cidr1
}
resource "aws_vpn_gateway" "vgw" {
  provider = aws.reg
  vpc_id = module.aviatrix_aws_vpc2.vpc
  tags = {
    Name = "aviatrix-vgw"
  }
}
