provider "aviatrix" {
  controller_ip = "13.126.166.7"
  username = "rakesh"
  password = "av1@Tr1x"
}

resource "aviatrix_gateway" "test_gateway1" {
  cloud_type = 1
  account_name = "devops"
  gw_name = "avtxgw1"
  vpc_id = "vpc-abcdef"
  vpc_reg = "us-west-1"
  vpc_size = "t2.micro"
  vpc_net = "10.0.0.0/24"
}
