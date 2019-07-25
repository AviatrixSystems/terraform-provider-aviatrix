output "AVIATRIX_CONTROLLER_IP" {
  value = module.aviatrix-controller-build.public_ip
}

output "AVIATRIX_USERNAME" {
  value = "admin"
}

output "AVIATRIX_PASSWORD" {
  value = var.admin_password
}

output "AWS_ACCOUNT_NUMBER" {
  value = data.aws_caller_identity.current.account_id
}

output "AWS_ACCESS_KEY" {
  value = var.aws_access_key
}

output "AWS_SECRET_KEY" {
  value = var.aws_secret_key
}

output "ARM_SUBSCRIPTION_ID" {
  value = var.azure_subscription_id
}

output "ARM_DIRECTORY_ID" {
  value = var.azure_tenant_id
}

output "ARM_APPLICATION_ID" {
  value = var.azure_client_id
}

output "ARM_APPLICATION_KEY"{
  value = var.azure_client_secret
}

output "GCP_ID" {
  value = var.gcp_project_id1
}

output "GCP_CREDENTIALS_FILEPATH"{
  value = var.gcp_credentials_file_path
}

output "AWS_BGP_VGW_ID" {
  value = aws_vpn_gateway.vgw.id
} 

output "GCP_VPC_ID" {
  value = module.aviatrix_gcp_vpc1.vpc_id
} 

output "GCP_SUBNET" {
  value = module.aviatrix_gcp_vpc1.subnet
} 

output "ARM_REGION" {
  value = var.azure_region1
} 

output "ARM_VNET_ID" {
  value = "${module.aviatrix_arm_vpc1.vnet}:${module.aviatrix_arm_vpc1.group}"
}

output "ARM_SUBNET" {
  value = module.aviatrix_arm_vpc1.subnet
}

output "ARM_VNET_ID2" {
  value = "${module.aviatrix_arm_vpc2.vnet}:${module.aviatrix_arm_vpc2.group}"
} 

output "ARM_REGION2" {
  value = var.azure_region2
} 

output "AWS_VPC_ID" {
  value = module.aviatrix_aws_vpc1.vpc
} 

output "AWS_SUBNET" {
  value = module.aviatrix_aws_vpc1.subnet
} 

output "AWS_VPC_ID2" {
  value = module.aviatrix_aws_vpc2.vpc
} 

output "AWS_SUBNET2" {
  value = module.aviatrix_aws_vpc2.subnet
} 

