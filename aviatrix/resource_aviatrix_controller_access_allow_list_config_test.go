package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixControllerAccessAllowListConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_ACCESS_ALLOW_LIST_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Access Allow List Config test as SKIP_CONTROLLER_ACCESS_ALLOW_LIST_CONFIG is set")
	}

	resourceName := "aviatrix_controller_access_allow_list_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControllerAccessAllowListConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerAccessAllowListConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerAccessAllowListConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "allow_list.0.ip_address", "0.0.0.0"),
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

func testAccControllerAccessAllowListConfigBasic() string {
	return `
resource "aviatrix_controller_access_allow_list_config" "test" {
	allow_list {
		ip_address = "0.0.0.0"
	}
}
`
}

func testAccCheckControllerAccessAllowListConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller access allow list config not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller access allow list config ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller access allow list config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerAccessAllowListConfigDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_access_allow_list_config" {
			continue
		}

		_, err := client.GetControllerAccessAllowList(context.Background())
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("controller access allow list config still exists")
		}
	}

	return nil
}
