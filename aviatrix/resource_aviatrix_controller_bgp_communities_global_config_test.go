package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixControllerBgpCommunitiesGlobalConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_BGP_COMMUNITIES_GLOBAL_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller BGP Communities Global Config test as SKIP_CONTROLLER_BGP_COMMUNITIES_GLOBAL_CONFIG is set")
	}
	resourceName := "aviatrix_controller_bgp_communities_global_config.test_bgp_communities_global"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerBgpCommunitiesGlobalConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerBgpCommunitiesGlobalConfigExists(resourceName),
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

func testAccControllerBgpCommunitiesGlobalConfigBasic() string {
	return `
resource "aviatrix_controller_bgp_communities_global_config" "test_bgp_communities_global" {
	bgp_communities_global = true
}
`
}

func testAccCheckControllerBgpCommunitiesGlobalConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller bgp communities global config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller bgp communities global config ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		_, err := client.GetControllerBgpCommunitiesGlobal(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get controller bgp communities global config status")
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller bgp communities global config ID not found")
		}

		return nil
	}
}
