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

func TestAccAviatrixFirewallInstanceAssociation_basic(t *testing.T) {
	if os.Getenv("SKIP_FIREWALL_INSTANCE_ASSOCIATION") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_FIREWALL_INSTANCE_ASSOCIATION is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_firewall_instance_association.test_firewall_instance_association"

	msg := ". Set SKIP_FIREWALL_INSTANCE_ASSOCIATION to 'yes' to skip firewall instance association tests."

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallInstanceAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallInstanceAssociationBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallInstanceAssociationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firenet_gw_name", fmt.Sprintf("tftg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "firewall_name", fmt.Sprintf("tffw-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "attached", "true"),
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

func testAccFirewallInstanceAssociationBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%[1]s"
	cloud_type         = 1
	aws_account_number = "%[2]s"
	aws_iam            = false
	aws_access_key     = "%[3]s"
	aws_secret_key     = "%[4]s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%[5]s"
	name                 = "vpc-%[1]s"
	cidr                 = "10.10.0.0/24"
	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test_transit_gateway" {
	cloud_type               = aviatrix_vpc.test_vpc.cloud_type
	account_name             = aviatrix_account.test_account.account_name
	gw_name                  = "tftg-%[1]s"
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
	firewall_name     = "tffw-%[1]s"
	firewall_image    = "Palo Alto Networks VM-Series Next-Generation Firewall Bundle 1"
	firewall_size     = "m5.xlarge"
	management_subnet = aviatrix_vpc.test_vpc.subnets[0].cidr
	egress_subnet     = aviatrix_vpc.test_vpc.subnets[1].cidr
}
resource "aviatrix_firenet" "test_firenet" {
	vpc_id                                = aviatrix_vpc.test_vpc.vpc_id
	inspection_enabled                    = true
	egress_enabled                        = false
	manage_firewall_instance_association  = false
 	depends_on = [aviatrix_transit_gateway.test_transit_gateway]
}
resource "aviatrix_firewall_instance_association" "test_firewall_instance_association" {
	vpc_id               = aviatrix_firenet.test_firenet.vpc_id
	firenet_gw_name      = aviatrix_transit_gateway.test_transit_gateway.gw_name
	instance_id          = aviatrix_firewall_instance.test_firewall_instance.instance_id
	firewall_name        = aviatrix_firewall_instance.test_firewall_instance.firewall_name
	attached             = true
	lan_interface        = aviatrix_firewall_instance.test_firewall_instance.lan_interface
	management_interface = aviatrix_firewall_instance.test_firewall_instance.management_interface
	egress_interface     = aviatrix_firewall_instance.test_firewall_instance.egress_interface
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"))
}

func testAccCheckFirewallInstanceAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("firewall_instance_association Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no firewall_instance_association ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		fireNet := &goaviatrix.FireNet{
			VpcID: rs.Primary.Attributes["vpc_id"],
		}
		fireNetDetail, err := client.GetFireNet(fireNet)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("could not find firenet: %v", err)
			}
			return fmt.Errorf("couldn't find FireNet: %s", err)
		}

		var instanceInfo *goaviatrix.FirewallInstanceInfo
		for _, v := range fireNetDetail.FirewallInstance {
			if v.GwName == rs.Primary.Attributes["firenet_gw_name"] &&
				v.InstanceID == rs.Primary.Attributes["instance_id"] {
				instanceInfo = &v
			}
		}
		if instanceInfo == nil {
			return fmt.Errorf("could not find firewall association")
		}

		return nil
	}
}

func testAccCheckFirewallInstanceAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_firewall_instance_association" {
			continue
		}
		fireNet := &goaviatrix.FireNet{
			VpcID: rs.Primary.Attributes["vpc_id"],
		}
		_, err := client.GetFireNet(fireNet)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("firenet still exists")
		}
	}
	return nil
}
