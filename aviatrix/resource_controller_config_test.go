package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixControllerConfig_basic(t *testing.T) {
	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipAcc := os.Getenv("SKIP_CONTROLLER_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Config test as SKIP_CONTROLLER_CONFIG is set")
	}
	msgCommon := ". Set SKIP_CONTROLLER_CONFIG to yes to skip Controller Config tests"
	preAccountCheck(t, msgCommon)
	resourceName := "aviatrix_controller_config.test_controller_config"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControllerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sg_management_account_name",
						fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "fqdn_exception_rule", "false"),
					resource.TestCheckResourceAttr(resourceName, "http_access", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_group_management", "true"),
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

func testAccControllerConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
    account_name       = "tfa-%s"
    cloud_type         = 1
    aws_account_number = "%s"
    aws_iam            = "false"
    aws_access_key     = "%s"
    aws_secret_key     = "%s"
}

resource "aviatrix_controller_config" "test_controller_config" {
	sg_management_account_name = "${aviatrix_account.test_account.account_name}"
	fqdn_exception_rule 	   = false
	http_access         	   = true
	security_group_management  = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccCheckControllerConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_config" {
			continue
		}
		_, err := client.GetHttpAccessEnabled()
		if err != nil {
			return fmt.Errorf("could not retrieve Http Access Status due to err: %v", err)
		}
	}
	return nil
}
