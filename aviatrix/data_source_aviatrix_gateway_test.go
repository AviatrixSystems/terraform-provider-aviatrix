package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceAviatrixGateway_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_gateway.foo"

	skipAcc := os.Getenv("SKIP_DATA_GATEWAY")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Gateway test as SKIP_DATA_GATEWAY is set")
	}

	preGatewayCheck(t, ". Set SKIP_DATA_GATEWAY to yes to skip Data Source Gateway tests")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixGatewayConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixGateway(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_size", "t2.micro"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixGatewayConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name 	   = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_gateway" "test_gw" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account.id
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	vpc_size     = "t2.micro"
	vpc_net      = "%s"
}

data "aviatrix_gateway" "foo" {
	account_name = aviatrix_gateway.test_gw.account_name
	gw_name      = aviatrix_gateway.test_gw.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_VPC_NET"))
}

func testAccDataSourceAviatrixGateway(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
