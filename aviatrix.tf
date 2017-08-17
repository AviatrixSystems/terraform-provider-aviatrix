provider "aviatrix" {
  controller_ip = "13.126.166.7"
  username = "rakesh"
  password = "av1@Tr1x"
}

resource "aviatrix_tunnel" "test_tunnel1" {
  vpc_name1 = "avtxgw1"
  vpc_name2 = "avtxgw2"
}
