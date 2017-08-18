provider "aviatrix" {
  controller_ip = "13.126.166.7"
  username = "rakesh"
  password = "av1@Tr1x"
}

resource "aviatrix_gateway" "test_gateway2" {
  cloud_type = 1
  account_name = "devops"
  gw_name = "avtxgw3"
  vpc_id = "vpc-0d7b3664"
  vpc_reg = "ap-south-1"
  vpc_size = "t2.micro"
  vpc_net = "avtxgw3_sub1~~10.3.0.0/24~~ap-south-1a"
}
