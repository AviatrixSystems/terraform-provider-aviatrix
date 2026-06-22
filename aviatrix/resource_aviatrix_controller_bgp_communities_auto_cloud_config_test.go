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

func TestAccAviatrixControllerBgpCommunitiesAutoCloudConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_BGP_AUTO_CLOUD_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller BGP Communities Auto Cloud Config test as SKIP_CONTROLLER_BGP_AUTO_CLOUD_CONFIG is set")
	}
	resourceName := "aviatrix_controller_bgp_communities_auto_cloud_config.test_bgp_communities_auto_cloud"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckControllerBgpCommunitiesAutoCloudConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerBgpCommunitiesAutoCloudConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerBgpCommunitiesAutoCloudConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_cloud_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "community_prefix", "12345"),
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

func testAccControllerBgpCommunitiesAutoCloudConfigBasic() string {
	return `
resource "aviatrix_controller_bgp_communities_auto_cloud_config" "test_bgp_communities_auto_cloud" {
	auto_cloud_enabled = true
	community_prefix = 12345
}
`
}

func testAccCheckControllerBgpCommunitiesAutoCloudConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller bgp communities auto cloud config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller bgp communities auto cloud config ID is set")
		}

		client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
		}

		_, err := client.GetControllerBgpCommunitiesAutoCloud(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get controller bgp communities auto cloud config status")
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller bgp communities auto cloud config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerBgpCommunitiesAutoCloudConfigDestroy(s *terraform.State) error {
	client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_bgp_communities_auto_cloud_config" {
			continue
		}

		_, err := client.GetControllerBgpCommunitiesAutoCloud(context.Background())
		if err == nil || !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("controller bgp communities auto cloud configured when it should be destroyed")
		}
	}

	return nil
}
