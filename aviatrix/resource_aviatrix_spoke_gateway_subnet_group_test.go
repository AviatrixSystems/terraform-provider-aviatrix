package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixSpokeGatewaySubnetGroup_basic(t *testing.T) {
	if os.Getenv("SKIP_SPOKE_GATEWAY_SUBNET_GROUP") == "yes" {
		t.Skip("Skipping spoke gateway subnet group test as SKIP_SPOKE_GATEWAY_SUBNET_GROUP is set")
	}

	resourceName := "aviatrix_spoke_gateway_subnet_group.test_group"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGatewaySubnetGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeGatewaySubnetGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGatewaySubnetGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-group"),
					resource.TestCheckResourceAttr(resourceName, "spoke_gateway_name", "test-spoke"),
					testAccCheckSpokeGatewaySubnetGroupSubnetsMatch(resourceName, []string{"18.9.16.0/20~~test-vpc-Public-subnet-1", "18.9.32.0/20~~test-vpc-Private-subnet-1"}),
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

func testAccSpokeGatewaySubnetGroupBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc" {
	account_name        = "tfa-azure"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 8
	account_name         = aviatrix_account.test_acc.account_name
	region               = "West US"
	name                 = "test-vpc"
	cidr                 = "18.9.0.0/16"
	aviatrix_firenet_vpc = false
}
resource "aviatrix_spoke_gateway" "test_spoke" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc.account_name
	gw_name      = "test-spoke"
	vpc_id       = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg      = aviatrix_vpc.test_vpc.region
	gw_size      = "Standard_B1ms"
	subnet       = aviatrix_vpc.test_vpc.subnets[0].cidr
}
data "aviatrix_spoke_gateway" "test" {
	gw_name = aviatrix_spoke_gateway.test_spoke.gw_name
}
resource "aviatrix_spoke_gateway_subnet_group" "test_group" {
	name               = "test-group"
	spoke_gateway_name = aviatrix_spoke_gateway.test_spoke.gw_name 
	subnets            = [
		data.aviatrix_spoke_gateway.test.all_subnets_for_inspection[0],
		data.aviatrix_spoke_gateway.test.all_subnets_for_inspection[1]
	]
}
    `, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"),
		os.Getenv("ARM_APPLICATION_KEY"))
}

func testAccCheckSpokeGatewaySubnetGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("spoke gateway subnet group not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke gateway subnet group ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		spokeGatewaySubnetGroup := &goaviatrix.SpokeGatewaySubnetGroup{
			SubnetGroupName: rs.Primary.Attributes["name"],
			GatewayName:     rs.Primary.Attributes["spoke_gateway_name"],
		}
		err := client.GetSpokeGatewaySubnetGroup(context.Background(), spokeGatewaySubnetGroup)
		if err != nil {
			return err
		}
		if spokeGatewaySubnetGroup.GatewayName+"~"+spokeGatewaySubnetGroup.SubnetGroupName != rs.Primary.ID {
			return fmt.Errorf("spoke gateway subnet group not found")
		}
		return nil
	}
}

func testAccCheckSpokeGatewaySubnetGroupSubnetsMatch(resourceName string, input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("spoke gateway subnet group not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		spokeGatewaySubnetGroup := &goaviatrix.SpokeGatewaySubnetGroup{
			SubnetGroupName: rs.Primary.Attributes["name"],
			GatewayName:     rs.Primary.Attributes["spoke_gateway_name"],
		}
		err := client.GetSpokeGatewaySubnetGroup(context.Background(), spokeGatewaySubnetGroup)
		if err != nil {
			return err
		}

		if !goaviatrix.Equivalent(spokeGatewaySubnetGroup.SubnetList, input) {
			return fmt.Errorf("subnets don't match with the input")
		}
		return nil
	}
}

func testAccCheckSpokeGatewaySubnetGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_gateway_subnet_group" {
			continue
		}

		spokeGatewaySubnetGroup := &goaviatrix.SpokeGatewaySubnetGroup{
			SubnetGroupName: rs.Primary.Attributes["name"],
			GatewayName:     rs.Primary.Attributes["spoke_gateway_name"],
		}

		err := client.GetSpokeGatewaySubnetGroup(context.Background(), spokeGatewaySubnetGroup)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("spoke gateway subnet group still exists")
		}
	}
	return nil
}
