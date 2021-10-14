package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixControllerSecurityGroupManagementConfig_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_CONTROLLER_SECURITY_GROUP_MANAGEMENT_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Config test as SKIP_CONTROLLER_SECURITY_GROUP_MANAGEMENT_CONFIG is set")
	}
	msgCommon := ". Set SKIP_CONTROLLER_SECURITY_GROUP_MANAGEMENT_CONFIG to yes to skip Controller Security Group Management Config tests"
	resourceName := "aviatrix_controller_security_group_management_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControllerSecurityGroupManagementConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerSecurityGroupManagementConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerSecurityGroupManagementConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_security_group_management", "false"),
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

func testAccControllerSecurityGroupManagementConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_controller_security_group_management_config" "test" {
	enable_security_group_management = false
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccCheckControllerSecurityGroupManagementConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller security group management config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller security group management config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller security group management config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerSecurityGroupManagementConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_security_group_management_config" {
			continue
		}

		_, err := client.GetSecurityGroupManagementStatus()
		if err != nil {
			return fmt.Errorf("could not retrieve Controller Security Group Management Status due to err: %v", err)
		}
	}

	return nil
}
