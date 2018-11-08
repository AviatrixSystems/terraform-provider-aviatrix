package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func preAccountCheck(t *testing.T, msgEnd string) {
	if os.Getenv("AWS_ACCOUNT_NUMBER") == "" {
		t.Fatal(" AWS_ACCOUNT_NUMBER must be set for acceptance tests" + msgEnd)
	}

	if os.Getenv("AWS_ACCESS_KEY") == "" {
		t.Fatal("AWS_ACCESS_KEY must be set for acceptance tests." + msgEnd)
	}

	if os.Getenv("AWS_SECRET_KEY") == "" {
		t.Fatal("AWS_SECRET_KEY must be set for acceptance tests." + msgEnd)
	}
}

func TestAccAviatrixAccount_basic(t *testing.T) {
	var account goaviatrix.Account
	rInt := acctest.RandInt()

	skipAcc := os.Getenv("SKIP_ACCOUNT")
	if skipAcc == "yes" {
		t.Skip("Skipping Access Account test as SKIP_ACCOUNT is set")
	}

	preAccountCheck(t, ". Set SKIP_ACCOUNT to yes to skip account tests")

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
						"aviatrix_account.foo", "aws_iam", "false"),
					resource.TestCheckResourceAttr(
						"aviatrix_account.foo", "aws_access_key",
						os.Getenv("AWS_ACCESS_KEY")),
					resource.TestCheckResourceAttr(
						"aviatrix_account.foo", "aws_secret_key",
						os.Getenv("AWS_SECRET_KEY")),
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
  aws_iam = "false"
  aws_access_key = "%s"
  aws_secret_key = "%s"
}
	`, rInt, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
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
