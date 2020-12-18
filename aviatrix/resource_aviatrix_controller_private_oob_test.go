package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixControllerPrivateOob_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_CONTROLLER_PRIVATE_OOB")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Private OOB test as SKIP_CONTROLLER_PRIVATE_OOB is set")
	}
	msgCommon := ". Set SKIP_CONTROLLER_PRIVATE_OOB to yes to skip Controller Private OOB tests"
	resourceName := "aviatrix_controller_private_oob.test_private_oob"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControllerPrivateOobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerPrivateOobBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerPrivateOobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_private_oob", "true"),
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

func testAccControllerPrivateOobBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_controller_private_oob" "test_private_oob" {
	enable_private_oob = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccCheckControllerPrivateOobExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller private oob ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller private oob ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller private oob ID not found")
		}

		return nil
	}
}

func testAccCheckControllerPrivateOobDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_private_oob" {
			continue
		}

		privateOobState, err := client.GetPrivateOobState()
		if err != nil {
			return fmt.Errorf("could not retrieve controller private oob Status")
		}
		if privateOobState {
			return fmt.Errorf("controller private oob is still enabled")
		}
	}

	return nil
}
