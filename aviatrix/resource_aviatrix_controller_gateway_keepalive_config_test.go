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

func TestAccAviatrixControllerGatewayKeepaliveConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_GATEWAY_KEEPALIVE_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Gateway Keepalive Config test as SKIP_CONTROLLER_GATEWAY_KEEPALIVE_CONFIG is set")
	}
	resourceName := "aviatrix_controller_gateway_keepalive_config.test_gateway_keepalive"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControllerGatewayKeepaliveConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerGatewayKeepaliveConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerGatewayKeepaliveConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "keepalive_speed", "slow"),
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

func testAccControllerGatewayKeepaliveConfigBasic() string {
	return `
resource "aviatrix_controller_gateway_keepalive_config" "test_gateway_keepalive" {
	keepalive_speed = "slow"
}
`
}

func testAccCheckControllerGatewayKeepaliveConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("gateway keepalive config resource ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no gateway keepalive config resource ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetGatewayKeepaliveConfig(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get gateway keepalive config status: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("gateway keepalive config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerGatewayKeepaliveConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_gateway_keepalive_config" {
			continue
		}

		speed, err := client.GetGatewayKeepaliveConfig(context.Background())
		if err != nil || speed != "medium" {
			return fmt.Errorf("gateway keepalive config configured when it should be destroyed")
		}
	}

	return nil
}
