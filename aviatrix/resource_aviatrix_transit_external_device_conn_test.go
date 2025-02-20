package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixTransitExternalDeviceConn_basic(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_external_device_conn.test"

	skipAcc := os.Getenv("SKIP_TRANSIT_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping transit external device connection tests as 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN' to 'yes' to skip Site2Cloud transit external device connection tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitExternalDeviceConnConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "123"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "345"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.14"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	skipBgpBfd := os.Getenv("SKIP_BGP_BFD_DEVICE_CONN")
	if skipBgpBfd != "yes" {
		resourceNameBfd := "aviatrix_transit_external_device_conn.transit-1-transit-2"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, ". Set 'SKIP_BGP_BFD_DEVICE_CONN' to 'yes' to skip Site2Cloud transit external device connection tests")
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitExternalDeviceConnConfigBgpBfd(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitExternalDeviceConnExists(resourceNameBfd, &externalDeviceConn),
						resource.TestCheckResourceAttr(resourceNameBfd, "vpc_id", os.Getenv("AWS_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameBfd, "connection_name", "transit-1-transit-2"),
						resource.TestCheckResourceAttr(resourceNameBfd, "gw_name", fmt.Sprintf("%s-aws-transit-1", rName)),
						resource.TestCheckResourceAttr(resourceNameBfd, "connection_type", "bgp"),
						resource.TestCheckResourceAttr(resourceNameBfd, "tunnel_protocol", "IPsec"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_local_as_num", "65075"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_remote_as_num", "65076"),
						resource.TestCheckResourceAttr(resourceNameBfd, "enable_bfd", "true"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.transmit_interval", "400"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.receive_interval", "400"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.multiplier", "5"),
					),
				},
				{
					Config: testAccTransitExternalDeviceConnConfigBgpBfdUpdated(),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitExternalDeviceConnExists(resourceNameBfd, &externalDeviceConn),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.transmit_interval", "500"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.receive_interval", "500"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.multiplier", "7"),
					),
				},
				{
					Config: testAccTransitExternalDeviceConnConfigBfdDisabled(),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitExternalDeviceConnExists(resourceNameBfd, &externalDeviceConn),
						resource.TestCheckResourceAttr(resourceNameBfd, "enable_bfd", "false"),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Skip("Skipping transit external device connection tests for bgp bfd config as 'SKIP_BGP_BFD_DEVICE_CONN' is set")
	}
}

func testAccTransitExternalDeviceConnConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_transit_external_device_conn" "test" {
	vpc_id            = aviatrix_transit_gateway.test.vpc_id
	connection_name   = "%s"
	gw_name           = aviatrix_transit_gateway.test.gw_name
	connection_type   = "bgp"
	bgp_local_as_num  = "123"
	bgp_remote_as_num = "345"
	remote_gateway_ip = "172.12.13.14"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func testAccTransitExternalDeviceConnConfigBgpBfd(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_transit_gateway" "transit_gateway_1" {
	single_az_ha = true
	gw_name = "%s-transit_gateway_1"
	vpc_id = "%s"
	cloud_type   = 1
	vpc_reg = "%s"
    gw_size = "t3.medium"
    account_name = "%s"
    subnet = "%s"
    connected_transit = true
    bgp_manual_spoke_advertise_cidrs = "192.0.2.0/24"
    bgp_polling_time = "10"
    local_as_number = "65075"
}
resource "aviatrix_transit_gateway" "transit_gateway_2" {
    single_az_ha = true
    gw_name = "%s-aws-transit-2"
    vpc_id = "%s"
    cloud_type = 1
    vpc_reg = "%s"
    gw_size = "t3.medium"
    account_name = "%s"
    subnet = "%s"
    connected_transit = true
    bgp_manual_spoke_advertise_cidrs = "198.51.100.0/24"
    bgp_polling_time = "10"
    local_as_number = "65076"
}
resource "aviatrix_transit_external_device_conn" "transit-1-transit-2" {
	vpc_id            = aviatrix_transit_gateway.transit_gateway_1.vpc_id
	connection_name   = "transit-1-transit-2"
	gw_name           = aviatrix_transit_gateway.transit_gateway_1.gw_name
	connection_type   = "bgp"
	tunnel_protocol    = "IPsec"
	bgp_local_as_num  = aviatrix_transit_gateway.transit_gateway_1.local_as_number
	bgp_remote_as_num = aviatrix_transit_gateway.transit_gateway_2.local_as_number
	remote_gateway_ip = aviatrix_transit_gateway.transit_gateway_2.eip
	local_tunnel_cidr  = "169.254.10.1/30"
	remote_tunnel_cidr = "169.254.10.2/30"
	pre_shared_key = "psk12"
	enable_bfd = true
	bgp_bfd {
	  transmit_interval = 400
	  receive_interval = 400
	  multiplier = 5
	}
  }
  resource "aviatrix_transit_external_device_conn" "transit-2-transit-1" {
	vpc_id            = avaitrix_transit_gateway.transit_gateway_2.vpc_id
	connection_name   = "transit-2-transit-1"
	gw_name           = aviatrix_transit_gateway.transit_gateway_2.gw_name
	connection_type   = "bgp"
	tunnel_protocol    = "IPsec"
	bgp_local_as_num  = aviatrix_transit_gateway.transit_gateway_2.local_as_number
	bgp_remote_as_num = aviatrix_transit_gateway.transit_gateway_1.local_as_number
	remote_gateway_ip = "3.140.172.87"
	local_tunnel_cidr  = "169.254.10.2/30"
	remote_tunnel_cidr = "169.254.10.1/30"
	pre_shared_key = "psk12"
  }
	`, rName, os.Getenv("AWS_VPC_ID_1"), os.Getenv("AWS_REGION"), os.Getenv("AWS_ACCOUNT_NAME"), os.Getenv("AWS_SUBNET_1"),
		rName, os.Getenv("AWS_VPC_ID_2"), os.Getenv("AWS_REGION"), os.Getenv("AWS_ACCOUNT_NAME"), os.Getenv("AWS_SUBNET_2"))
}

func testAccTransitExternalDeviceConnConfigBgpBfdUpdated() string {
	return `
	resource "aviatrix_transit_external_device_conn" "transit-1-transit-2" {
		vpc_id            = aviatrix_transit_gateway.transit_gateway_1.vpc_id
		connection_name   = "transit-1-transit-2"
		gw_name           = aviatrix_transit_gateway.transit_gateway_1.gw_name
		connection_type   = "bgp"
		tunnel_protocol    = "IPsec"
		bgp_local_as_num  = aviatrix_transit_gateway.transit_gateway_1.local_as_number
		bgp_remote_as_num = aviatrix_transit_gateway.transit_gateway_2.local_as_number
		remote_gateway_ip = aviatrix_transit_gateway.transit_gateway_2.eip
		local_tunnel_cidr  = "169.254.10.1/30"
		remote_tunnel_cidr = "169.254.10.2/30"
		pre_shared_key = "psk12"
		enable_bfd = true
		bgp_bfd {
		  transmit_interval = 500
		  receive_interval = 500
		  multiplier = 7
		}
	}`
}

func testAccTransitExternalDeviceConnConfigBfdDisabled() string {
	return `
	resource "aviatrix_transit_external_device_conn" "transit-1-transit-2" {
		vpc_id            = aviatrix_transit_gateway.transit_gateway_1.vpc_id
		connection_name   = "transit-1-transit-2"
		gw_name           = aviatrix_transit_gateway.transit_gateway_1.gw_name
		connection_type   = "bgp"
		tunnel_protocol    = "IPsec"
		bgp_local_as_num  = aviatrix_transit_gateway.transit_gateway_1.local_as_number
		bgp_remote_as_num = aviatrix_transit_gateway.transit_gateway_2.local_as_number
		remote_gateway_ip = aviatrix_transit_gateway.transit_gateway_2.eip
		local_tunnel_cidr  = "169.254.10.1/30"
		remote_tunnel_cidr = "169.254.10.2/30"
		pre_shared_key = "psk12"
		enable_bfd = false
	}`
}

func testAccCheckTransitExternalDeviceConnExists(n string, externalDeviceConn *goaviatrix.ExternalDeviceConn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transit external device connection Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit external device connection ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundExternalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}
		foundExternalDeviceConn2, err := client.GetExternalDeviceConnDetail(foundExternalDeviceConn)
		if err != nil {
			return err
		}
		if foundExternalDeviceConn2.ConnectionName+"~"+foundExternalDeviceConn2.VpcID != rs.Primary.ID {
			return fmt.Errorf("transit external device connection not found")
		}

		*externalDeviceConn = *foundExternalDeviceConn2
		return nil
	}
}

func testAccCheckTransitExternalDeviceConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_external_device_conn" {
			continue
		}

		foundExternalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetExternalDeviceConnDetail(foundExternalDeviceConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("site2cloud still exists %s", err.Error())
		}
	}

	return nil
}

func TestAccAviatrixEdgeTransitExternalDeviceConn(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	rName := "aviatrix_transit_external_device_conn.eat-2-bgpoipsec-1"

	skipAcc := os.Getenv("SKIP_EDGE_TRANSIT_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping edge transit external device connection tests as 'SKIP_EDGE_TRANSIT_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_EDGE_TRANSIT_EXTERNAL_DEVICE_CONN' to 'yes' to skip Edge transit external device connection tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeTransitExternalDeviceConnConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(rName, &externalDeviceConn),
					resource.TestCheckResourceAttr(rName, "connection_name", "eat-2-bgpoipsec-1"),
					resource.TestCheckResourceAttr(rName, "gw_name", "e2e-edge-transit-2"),
					resource.TestCheckResourceAttr(rName, "enable_jumbo_frame", "true"),
					resource.TestCheckResourceAttr(rName, "tunnel_src_ip", "192.168.20.117,192.168.23.16"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEdgeTransitExternalDeviceConnConfig() string {
	return `
	resource "aviatrix_transit_external_device_conn" "eat-2-bgpoipsec-1" {
	vpc_id                    = "tsite-2"
	connection_name           = "eat-2-bgpoipsec-1"
	gw_name                   = "e2e-edge-transit-2"
	connection_type           = "bgp"
	tunnel_protocol           = "IPsec"
	bgp_local_as_num          = "65402"
	bgp_remote_as_num         = "65151"
	remote_gateway_ip         = "3.140.40.45"
	pre_shared_key         = "aviatrix,aviatrix"
	direct_connect            = true
	enable_jumbo_frame            = true
	disable_activemesh = false
	ha_enabled                = true
	local_tunnel_cidr         = "169.254.22.54/30, 169.254.238.22/30"
	remote_tunnel_cidr        = "169.254.22.53/30, 169.254.238.21/30"
	backup_local_tunnel_cidr  = "169.254.33.94/30, 169.254.165.254/30"
	backup_remote_tunnel_cidr = "169.254.33.93/30, 169.254.165.253/30"
	backup_bgp_remote_as_num  = "65151"
	backup_remote_gateway_ip  = "18.223.219.22"
	backup_direct_connect     = true
	backup_pre_shared_key         = "aviatrix,aviatrix"
	tunnel_src_ip  = "192.168.20.117,192.168.23.16"
	}`
}
