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
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixControllerBgpMaxAsLimitConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_BGP_MAX_AS_LIMIT_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller BGP Max AS Limit Config test as SKIP_CONTROLLER_BGP_MAX_AS_LIMIT_CONFIG is set")
	}
	msgCommon := ". Set SKIP_CONTROLLER_BGP_MAX_AS_LIMIT_CONFIG to yes to skip Controller BGP Max AS Limit Config tests"
	resourceName := "aviatrix_controller_bgp_max_as_limit_config.test_bgp_max_as_limit"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
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
	return fmt.Sprintf(`
resource "aviatrix_controller_bgp_max_as_limit_config" "test_bgp_max_as_limit" {
	max_as_limit = 1
}
`)
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

		maxAsLimit, err := client.GetControllerBgpMaxAsLimit(context.Background())
		if err != nil {
			return fmt.Errorf("could not retrieve controller bgp max as limit config status: %v", err)
		}
		if maxAsLimit != -1 {
			return errors.New("controller bgp max as limit is still enabled")
		}
	}

	return nil
}
