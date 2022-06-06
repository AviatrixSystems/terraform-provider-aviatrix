package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixTransitGateways_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_transit_gateways.foo"

	skipAcc := os.Getenv("SKIP_DATA_TRANSIT_GATEWAYS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source All Transit Gateway tests as SKIP_DATA_TRANSIT_GATEWAYS is set")
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set SKIP_DATA_TRANSIT_GATEWAYS to yes to skip Data Source All Transit Gateway tests")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixTransitGatewaysConfigBasic(rName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aviatrix_transit_gateway.test", "gw_name", fmt.Sprintf("aa-tfg-aws-%s", rName)),
					resource.TestCheckResourceAttr("aviatrix_transit_gateway.test2", "gw_name", fmt.Sprintf("aa-tfg-gcp-%s", rName)),
				),

				Destroy: false,
			},
			{
				Config: testAccDataSourceAviatrixTransitGatewaysConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixTransitGateways(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.0.gw_name", fmt.Sprintf("aa-tfg-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.0.vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.0.vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.0.gw_size", "t2.micro"),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.1.gw_name", fmt.Sprintf("aa-tfg-gcp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.1.gw_size", "n1-standard-1"),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.1.account_name", fmt.Sprintf("aa-tfa-gcp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.1.subnet", os.Getenv("GCP_SUBNET")),
					resource.TestCheckResourceAttr(resourceName, "gateway_list.1.vpc_reg", os.Getenv("GCP_ZONE")),
				),
				Destroy: true,
			},
		},
	})
}

func testAccDataSourceAviatrixTransitGatewaysConfigBasic(rName string) string {
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
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "aa-tfa-gcp-%[1]s"
	cloud_type                          = 4
	gcloud_project_id                   = "%[8]s"
	gcloud_project_credentials_filepath = "%[9]s"
}
resource "aviatrix_transit_gateway" "test2" {				
	cloud_type   = 4
	account_name = aviatrix_account.test_acc_gcp.account_name
	gw_name      = "aa-tfg-gcp-%[1]s"
	vpc_id       = "%[10]s"
	vpc_reg      = "%[11]s"
	gw_size      = "n1-standard-1"
	subnet       = "%[12]s"
}
data "aviatrix_transit_gateways" "foo" {

}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), os.Getenv("GCP_SUBNET"))
}

func testAccDataSourceAviatrixTransitGateways(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
