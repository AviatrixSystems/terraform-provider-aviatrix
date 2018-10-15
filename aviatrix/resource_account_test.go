package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccAviatrixAccount_basic(t *testing.T) {
	var account goaviatrix.Account
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountExists("aviatrix_account.foo", &account),
					resource.TestCheckResourceAttr(
						"aviatrix_account.foo", "account_name", fmt.Sprintf("tf-testing-%d", rInt)),
					resource.TestCheckResourceAttr(
						"aviatrix_account.foo", "aws_iam", "true"),
				),
			},
		},
	})
}

func testAccAccountConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "foo" {
  account_name = "tf-testing-%d"
  cloud_type = 1
  aws_account_number = "%s"
  aws_iam = "true"
  aws_role_app = "arn:aws:iam::%s:role/aviatrix-role-app"
  aws_role_ec2 = "arn:aws:iam::%s:role/aviatrix-role-ec2"
}
	`, rInt, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCOUNT_NUMBER"))
}

func testAccCheckAccountExists(n string, account *goaviatrix.Account) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("account Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Account ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAccount := &goaviatrix.Account{
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetAccount(foundAccount)

		if err != nil {
			return err
		}

		if foundAccount.AccountName != rs.Primary.ID {
			return fmt.Errorf("account not found")
		}

		*account = *foundAccount

		return nil
	}
}

func testAccCheckAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_account" {
			continue
		}
		foundAccount := &goaviatrix.Account{
			AccountName: rs.Primary.Attributes["account_name"],
		}
		_, err := client.GetAccount(foundAccount)

		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("account still exists")
		}
	}
	return nil
}
