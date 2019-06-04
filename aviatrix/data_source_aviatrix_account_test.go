package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceAviatrixAccount_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_account.foo"

	skipAcc := os.Getenv("SKIP_DATA_ACCOUNT")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Account test as SKIP_DATA_ACCOUNT is set")
	}

	preAccountCheck(t, ". Set SKIP_DATA_ACCOUNT to yes to skip Data Source Account tests")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixAccountConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixAccount(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "account_name", fmt.Sprintf("tf-testing-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "aws_access_key",
						os.Getenv("AWS_ACCESS_KEY")),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixAccountConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name = "tf-testing-%s"
	cloud_type = 1
	aws_account_number = "%s"
	aws_iam = "false"
	aws_access_key = "%s"
	aws_secret_key = "%s"
}

data "aviatrix_account" "foo" {
	account_name = "${aviatrix_account.test.id}"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccDataSourceAviatrixAccount(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		return nil
	}
}
