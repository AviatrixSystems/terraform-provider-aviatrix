package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixSpokeGateways_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_spoke_gateways.foo"

	skipAcc := os.Getenv("SKIP_DATA_SPOKE_GATEWAYS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source All Spoke Gateways tests as SKIP_DATA_SPOKE_GATEWAYS is set")
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixSpokeGatewaysConfigBasic(rName),

				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixSpokeGateways(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.0.gw_name", fmt.Sprintf("aa-tfg-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.0.vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.0.vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.0.gw_size", "t2.micro"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixSpokeGatewaysConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name 	   = "aa-tfa-%[1]s"
	cloud_type         = 1
	aws_account_number = "%[2]s"
	aws_iam            = "false"
	aws_access_key     = "%[3]s"
	aws_secret_key     = "%[4]s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account.account_name
	gw_name      = "aa-tfg-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
data "aviatrix_transit_gateways" "foo" {
    depends_on = [
		aviatrix_transit_gateway.test,
    ]
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccDataSourceAviatrixSpokeGateways(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
