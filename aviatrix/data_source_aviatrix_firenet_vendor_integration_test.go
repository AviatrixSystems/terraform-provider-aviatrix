package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceAviatrixFireNetVendorIntegration_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_firenet_vendor_integration.test"

	skipAcc := os.Getenv("SKIP_DATA_FIRENET_VENDOR_INTEGRATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source FireNet Vendor Integration test as SKIP_DATA_FIRENET_VENDOR_INTEGRATION is set")
	}
	msg := ". Set SKIP_DATA_FIRENET_VENDOR_INTEGRATION to yes to skip Data Source FireNet Vendor Integration tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixFireNetVendorIntegrationConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixFireNetVendorIntegration(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firewall_name", fmt.Sprintf("tffw-%s", rName)),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixFireNetVendorIntegrationConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%s"
	name                 = "vpc-for-firenet"
	cidr                 = "10.10.0.0/24"
	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test_transit_gateway" {
	cloud_type               = aviatrix_vpc.test_vpc.cloud_type
	account_name             = aviatrix_account.test_account.account_name
	gw_name                  = "tftg-%s"
	vpc_id                   = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg                  = aviatrix_vpc.test_vpc.region
	gw_size                  = "c5.xlarge"
	subnet                   = aviatrix_vpc.test_vpc.subnets[0].cidr
	enable_hybrid_connection = true
	enable_firenet           = true
	enable_active_mesh       = true
}
resource "aviatrix_firewall_instance" "test_firewall_instance" {
	vpc_id            = aviatrix_vpc.test_vpc.vpc_id
	firenet_gw_name   = aviatrix_transit_gateway.test_transit_gateway.gw_name
	firewall_name     = "tffw-%s"
	firewall_image    = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
	firewall_size     = "m5.xlarge"
	management_subnet = aviatrix_vpc.test_vpc.subnets[0].cidr
	egress_subnet     = aviatrix_vpc.test_vpc.subnets[1].cidr
}
resource "aviatrix_firenet" "test_firenet" {
	vpc_id             = aviatrix_vpc.test_vpc.vpc_id
	inspection_enabled = true
	egress_enabled     = false

	firewall_instance_association {
		firenet_gw_name      = aviatrix_transit_gateway.test_transit_gateway.gw_name
		instance_id          = aviatrix_firewall_instance.test_firewall_instance.instance_id
		firewall_name        = aviatrix_firewall_instance.test_firewall_instance.firewall_name
		attached             = true
		lan_interface        = aviatrix_firewall_instance.test_firewall_instance.lan_interface
		management_interface = aviatrix_firewall_instance.test_firewall_instance.management_interface
		egress_interface     = aviatrix_firewall_instance.test_firewall_instance.egress_interface
	}
}
data "aviatrix_firenet_vendor_integration" "test" {
	vpc_id        = aviatrix_vpc.test_vpc.vpc_id
	instance_id   = aviatrix_firewall_instance.test_firewall_instance.instance_id
	vendor_type   = "Generic"
	public_ip     = aviatrix_firewall_instance.test_firewall_instance.public_ip
	username      = "admin"
	password      = "Avx123456#"
	firewall_name = aviatrix_firewall_instance.test_firewall_instance.firewall_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName, rName)
}

func testAccDataSourceAviatrixFireNetVendorIntegration(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
