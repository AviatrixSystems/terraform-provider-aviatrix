package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixSpokeGatewayInspectionSubnets_basic(t *testing.T) {
	resourceName := "data.aviatrix_spoke_gateway_inspection_subnets.foo"

	skipAcc := os.Getenv("SKIP_DATA_SPOKE_GATEWAY_INSPECTION_SUBNETS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Spoke Gateway Inspection Subnets test as SKIP_DATA_SPOKE_GATEWAY_INSPECTION_SUBNETS is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixSpokeGatewayInspectionSubnetsConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixSpokeGatewayInspectionSubnets(resourceName),
					testAccCheckSpokeGatewayInspectionSubnetsMatch(resourceName, []string{"18.9.16.0/20~~test-vpc-Public-subnet-1", "18.9.32.0/20~~test-vpc-Private-subnet-1", "18.9.48.0/20~~test-vpc-Public-subnet-2", "18.9.64.0/20~~test-vpc-Private-subnet-2"}),
				),
			},
		},
	})

}

func testAccDataSourceAviatrixSpokeGatewayInspectionSubnetsConfigBasic() string {
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
data "aviatrix_spoke_gateway_inspection_subnets" "foo" {
	gw_name = aviatrix_spoke_gateway.test_spoke.gw_name
}
    `, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"),
		os.Getenv("ARM_APPLICATION_KEY"))
}

func testAccDataSourceAviatrixSpokeGatewayInspectionSubnets(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}

func testAccCheckSpokeGatewayInspectionSubnetsMatch(resourceName string, input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("data source spoke gateway inspection subnets not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		subnetsForInspection, err := client.GetSubnetsForInspection(rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if !goaviatrix.Equivalent(subnetsForInspection, input) {
			return fmt.Errorf("subnets don't match with the input")
		}
		return nil
	}
}
