resource "aviatrix_account" "test_aws" {
  cloud_type         = 1
  account_name       = "aws-%s"
  aws_account_number = "%s"
  aws_iam            = false
  aws_access_key     = "%s"
  aws_secret_key     = "%s"
}

resource "aviatrix_account" "test_acc_edge_megaport" {
  account_name = "megaport-%s"
  cloud_type   = 1048576
}

resource "aviatrix_vpc" "test_vpc" {
  cloud_type           = 1
  account_name         = aviatrix_account.test_aws.account_name
  region               = "us-west-1"
  name                 = "aws-vpc-test-1"
  cidr                 = "16.0.0.0/20"
  aviatrix_transit_vpc = true
}

resource "aviatrix_spoke_gateway" "test_spoke" {
  cloud_type     = 1
  account_name   = aviatrix_account.test.account_name
  gw_name        = "tfs-%s"
  vpc_id         = aviatrix_vpc.test_vpc.vpc_id
  vpc_reg        = aviatrix_vpc.test_vpc.region
  gw_size        = "c5.xlarge"
  insane_mode    = true
  subnet         = join(".", [join(".", slice(split(".", aviatrix_vpc.test1.public_subnets[1].cidr), 0, 2)), "12.0/26"]) #"173.31.12.0/26"
  insane_mode_az = "us-east-1a"
}

resource "aviatrix_transit_gateway" "test_edge_transit" {
  cloud_type             = 1048576
  account_name           = aviatrix_account.test_acc_edge_megaport.account_name
  gw_name                = "%s"
  vpc_id                 = "%s"
  gw_size                = "SMALL"
  ztp_file_download_path = "%s"
  interfaces {
    gateway_ip                  = "192.168.20.1"
    ip_address                  = "192.168.20.11/24"
    public_ip                   = "67.207.104.19"
    logical_ifname              = "wan0"
    secondary_private_cidr_list = ["192.168.20.16/29"]
  }

  interfaces {
    gateway_ip                  = "192.168.21.1"
    ip_address                  = "192.168.21.11/24"
    public_ip                   = "67.71.12.148"
    logical_ifname              = "wan1"
    secondary_private_cidr_list = ["192.168.21.16/29"]
  }

  interfaces {
    dhcp           = true
    logical_ifname = "mgmt0"
  }

  interfaces {
    gateway_ip     = "192.168.22.1"
    ip_address     = "192.168.22.11/24"
    logical_ifname = "wan2"
  }

  interfaces {
    gateway_ip     = "192.168.23.1"
    ip_address     = "192.168.23.11/24"
    logical_ifname = "wan3"
  }
}

resource "aviatrix_spoke_transit_attachment" "test" {
  spoke_gw_name                   = aviatrix_spoke_gateway.test_spoke.gw_name
  transit_gw_name                 = aviatrix_transit_gateway.test_edge_transit.gw_name
  tunnel_count                    = 4
  transit_gateway_logical_ifnames = ["wan1"]

}
