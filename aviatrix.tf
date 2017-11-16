provider "aviatrix" {
  controller_ip = "1.2.3.4"
  username = "admin"
  password = "password"
}

resource "aviatrix_transpeer" "test_transpeer" {
  source = "avtxuseastgw1"
  nexthop = "avtxuseastgw2"
  reachable_cidr = "10.152.0.0/16"
}
