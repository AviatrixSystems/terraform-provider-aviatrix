package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAdminEmail_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_ADMIN_EMAIL")
	if skipAcc == "yes" {
		t.Skip("Skipping Admin Email test as SKIP_ADMIN_EMAIL is set")
	}

	preAccountCheck(t, ". Set SKIP_ADMIN_EMAIL to yes to skip admin email tests")

	adminEmail := "abc@xyz.com"

	resourceName := "aviatrix_admin_email.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminEmailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAdminEmailConfigBasic(adminEmail),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminEmailExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "admin_email", "abc@xyz.com"),
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

func testAccAdminEmailConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_admin_email" "foo" {
    admin_email = "%s"
}
	`, rName)
}

func testAccCheckAdminEmailExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("admin email Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Admin Email is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAdminEmail, err := client.GetAdminEmail(client.Username, client.Password)

		if err != nil {
			return err
		}
		if foundAdminEmail != rs.Primary.Attributes["admin_email"] {
			return fmt.Errorf("admin email not found")
		}

		return nil
	}
}

func testAccCheckAdminEmailDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_admin_email" {
			continue
		}

		_, err := client.GetAdminEmail(client.Username, client.Password)

		if err != nil {
			return fmt.Errorf("could not retrieve Admin Email due to err: %v", err)
		}
	}
	return nil
}
