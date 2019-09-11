# Sample Aviatrix terraform configuration to create a full mesh network on AWS
# This configuration creates a cloud account on Aviatrix controller, launches 3 gateways with the created account
# and establishes tunnels between each gateway.


# Edit to enter your controller's IP, username and password to login with
provider "aviatrix" {
  controller_ip = "54.54.54.54"
  username      = "admin"
  password      = "Aviatrix#123"
}

# Increase count default value to add more VPCs and subnets to launch more gateways together.

variable "number" {
  default = 3
}

# Enter VPCs where you want to launch gateways.
variable "vpcs" {
  description = "Launch gateways in different VPCs."
  type        = "list"
  default     = [
    "vpc-7a6b2513",
    "vpc-2ee4a147",
    "vpc-0d7b3664",
  ]
}

# Enter Subnets within VPCs added above.
variable "vpc_nets" {
  description = "Launch gateways in different VPC Subnets."
  type        = "list"
  default     = [
    "10.1.0.0/24",
    "10.2.0.0/24",
    "10.3.0.0/24",
  ]
}

resource "aviatrix_account" "test_acc" {
  account_name       = "devops"
  cloud_type         = 1
  aws_account_number = "123456789012"
  aws_iam            = true
}

# Create count number of gateways
resource "aviatrix_gateway" "test_gw" {
  count        = var.number
  cloud_type   = 1
  account_name = "devops"
  gw_name      = "avtxgw-${count.index}"
  vpc_id       = "element(var.vpcs, ${count.index})"
  vpc_reg      = "ap-south-1"
  gw_size      = "t2.micro"
  subnet       = "element(var.vpc_nets, ${count.index})"
  depends_on   = [
    "aviatrix_account.test_acc"
  ]
}

# Create tunnels between above created gateways
resource "aviatrix_tunnel" "test_tunnel" {
  count      = var.number * (var.number - 1)/2
  gw_name1   = "avtxgw-${count.index}"
  gw_name2   = "avtxgw-${(count.index+1)%3}"
  depends_on = [
    "aviatrix_gateway.test_gw"
  ]
}

# For the complete documentation of all resources visit
https://www.terraform.io/docs/providers/aviatrix/
