package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixControllerConfig_basic(t *testing.T) {

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
				Config: testAccControllerConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerConfigExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "http_access", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "fqdn_exception_rule", "false"),
				),
			},
		},
	})
}

func testAccControllerConfigBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_controller_config" "test_controller_config" {
	http_access         = true
	fqdn_exception_rule = false
}
	`)
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
