package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixControllerBgpMaxAsLimitConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_BGP_MAX_AS_LIMIT_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller BGP Max AS Limit Config test as SKIP_CONTROLLER_BGP_MAX_AS_LIMIT_CONFIG is set")
	}
	resourceName := "aviatrix_controller_bgp_max_as_limit_config.test_bgp_max_as_limit"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckControllerBgpMaxAsLimitConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerBgpMaxAsLimitConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerBgpMaxAsLimitConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "max_as_limit", "1"),
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

func testAccControllerBgpMaxAsLimitConfigBasic() string {
	return `
resource "aviatrix_controller_bgp_max_as_limit_config" "test_bgp_max_as_limit" {
	max_as_limit = 1
}
`
}

func testAccCheckControllerBgpMaxAsLimitConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller bgp max as limit config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller bgp max as limit config ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		maxAsLimit, err := client.GetControllerBgpMaxAsLimit(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get controller bgp max as limit config status: %v", err)
		} else if maxAsLimit != 1 {
			return fmt.Errorf("API returned the wrong value for controller bgp max as limit: expected %d but got %d", 1, maxAsLimit)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller bgp max as limit config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerBgpMaxAsLimitConfigDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_bgp_max_as_limit_config" {
			continue
		}

		_, err := client.GetControllerBgpMaxAsLimit(context.Background())
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("could not retrieve controller bgp max as limit config status: %v", err)
		}
	}

	return nil
}
