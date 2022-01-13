package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixFirewallInstance_basic(t *testing.T) {
	var firewallInstance goaviatrix.FirewallInstance

	rName := acctest.RandString(5)
	resourceName := "aviatrix_firewall_instance.test_firewall_instance"

	skipAcc := os.Getenv("SKIP_FIREWALL_INSTANCE")
	if skipAcc == "yes" {
		t.Skip("Skipping Firewall Instance test as SKIP_FIREWALL_INSTANCE is set")
	}
	msg := ". Set SKIP_FIREWALL_INSTANCE to yes to skip Firewall Instance tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallInstanceConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallInstanceExists(resourceName, &firewallInstance),
					resource.TestCheckResourceAttr(resourceName, "firenet_gw_name", fmt.Sprintf("tftg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "firewall_name", fmt.Sprintf("tffw-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "firewall_image", "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"),
					resource.TestCheckResourceAttr(resourceName, "firewall_size", "m5.xlarge"),
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

func testAccFirewallInstanceConfigBasic(rName string) string {
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
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName, rName)
}

func testAccCheckFirewallInstanceExists(n string, firewallInstance *goaviatrix.FirewallInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("firewall instance Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no firewall instance ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundFirewallInstance := &goaviatrix.FirewallInstance{
			InstanceID: rs.Primary.Attributes["instance_id"],
		}

		foundFirewallInstance2, err := client.GetFirewallInstance(foundFirewallInstance)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("firewall instance not found")
			}
			return err
		}
		if foundFirewallInstance2.InstanceID != rs.Primary.ID {
			return fmt.Errorf("firewall instance not found")
		}

		*firewallInstance = *foundFirewallInstance
		return nil
	}
}

func testAccCheckFirewallInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_firewall_instance" {
			continue
		}

		foundFirewallInstance := &goaviatrix.FirewallInstance{
			InstanceID: rs.Primary.Attributes["instance_id"],
		}

		_, err := client.GetFirewallInstance(foundFirewallInstance)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("firewall instance still exists")
		}
	}

	return nil
}
