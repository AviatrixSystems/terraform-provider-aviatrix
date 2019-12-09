package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixFireNet_basic(t *testing.T) {
	var fireNet goaviatrix.FireNet

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_firenet.test_firenet"

	skipAcc := os.Getenv("SKIP_FIRENET")
	if skipAcc == "yes" {
		t.Skip("Skipping FireNet test as SKIP_FIRENET is set")
	}
	msg := ". Set SKIP_FIRENET to yes to skip FireNet tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFireNetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFireNetConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireNetExists(resourceName, &fireNet),
					resource.TestCheckResourceAttr(resourceName, "inspection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "egress_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "firewall_instance_association.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "firewall_instance_association.0.firenet_gw_name", fmt.Sprintf("tftg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "firewall_instance_association.0.firewall_name", fmt.Sprintf("tffw-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "firewall_instance_association.0.attached", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFireNetConfigBasic(rName string) string {
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
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName, rName)
}

func testAccCheckFireNetExists(n string, fireNet *goaviatrix.FireNet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("fireNet Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no FireNet ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundFireNet := &goaviatrix.FireNet{
			VpcID: rs.Primary.Attributes["vpc_id"],
		}

		foundFireNet2, err := client.GetFireNet(foundFireNet)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("fireNet not found")
			}
			return err
		}
		if foundFireNet2.VpcID != rs.Primary.ID {
			return fmt.Errorf("fireNet not found")
		}

		*fireNet = *foundFireNet
		return nil
	}
}

func testAccCheckFireNetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_firenet" {
			continue
		}

		foundFireNet := &goaviatrix.FireNet{
			VpcID: rs.Primary.Attributes["vpc_id"],
		}

		_, err := client.GetFireNet(foundFireNet)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("fireNet still exists")
		}
	}

	return nil
}
