package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAccountUser_basic(t *testing.T) {
	var account goaviatrix.AccountUser

	rInt := acctest.RandInt()
	resourceName := "aviatrix_account_user.foo"
	importStateVerifyIgnore := []string{"password"}

	skipAcc := os.Getenv("SKIP_ACCOUNT_USER")
	if skipAcc == "yes" {
		t.Skip("Skipping Account User test as SKIP_ACCOUNT_USER is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccountUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountUserConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountUserExists("aviatrix_account_user.foo", &account),
					resource.TestCheckResourceAttr(resourceName, "username", fmt.Sprintf("tf-testing-%d", rInt)),
					resource.TestCheckResourceAttr(resourceName, "email", "abc@xyz.com"),
					resource.TestCheckResourceAttr(resourceName, "password", "Password-1234^"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},
		},
	})
}

func testAccAccountUserConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_account_user" "foo" {
	username = "tf-testing-%d"
	email    = "abc@xyz.com"
	password = "Password-1234^"
}
	`, rInt)
}

func testAccCheckAccountUserExists(n string, account *goaviatrix.AccountUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("account Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Account ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAccount := &goaviatrix.AccountUser{
			UserName: rs.Primary.Attributes["username"],
		}

		_, err := client.GetAccountUser(foundAccount)
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("account not found in REST response")
		}
		if foundAccount.UserName != rs.Primary.ID {
			return fmt.Errorf("account not found")
		}

		*account = *foundAccount
		return nil
	}
}

func testAccCheckAccountUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_account" {
			continue
		}

		foundAccount := &goaviatrix.AccountUser{
			UserName: rs.Primary.Attributes["username"],
		}

		_, err := client.GetAccountUser(foundAccount)
		if err != nil {
			return fmt.Errorf("account still exists")
		}
	}
	return nil
}
