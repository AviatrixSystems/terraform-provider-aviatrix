package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAviatrixVpc_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_vpc.test"

	skipAcc := os.Getenv("SKIP_DATA_VPC")
	if skipAcc == "yes" {
		t.Skip("Skipping data source vpc tests as 'SKIP_DATA_VPC' is set")
	}

	msg := ". Set 'SKIP_DATA_VPC' to 'yes' to skip data source vpc tests"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msg)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixVpcConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixVpc(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfv-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "cidr", "10.0.0.0/16"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixVpcConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	name         = "tfv-%s"
	region       = "%s"
	cidr         = "10.0.0.0/16"
}
data "aviatrix_vpc" "test" {
	name = aviatrix_vpc.test.name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_REGION"))
}

func testAccDataSourceAviatrixVpc(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
