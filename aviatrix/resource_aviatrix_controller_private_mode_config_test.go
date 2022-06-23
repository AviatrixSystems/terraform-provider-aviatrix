package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixControllerPrivateModeConfig_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_CONTROLLER_PRIVATE_MODE_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Private Mode config tests as SKIP_CONTROLLER_PRIVATE_MODE_CONFIG is set")
	}
	msgCommon := ". Set SKIP_CONTROLLER_PRIVATE_MODE_CONFIG to yes to skip Controller Private Mode config tests"
	resourceName := "aviatrix_controller_private_mode_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccControllerPrivateModeConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerPrivateModeConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccControllerPrivateModeConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_private_mode", "true"),
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

func testAccControllerPrivateModeConfigBasic(rName string) string {
	return `
resource "aviatrix_controller_private_mode_config" "test" {
	enable_private_mode = true
}
	`
}

func testAccControllerPrivateModeConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller private mode config Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller private mode config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller private mode config ID not found")
		}

		return nil
	}
}

func testAccControllerPrivateModeConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_private_mode_config" {
			continue
		}

		info, err := client.GetPrivateModeInfo(context.Background())
		if err != nil {
			return fmt.Errorf("could not retrieve controller private mode config")
		}
		if info.EnablePrivateMode {
			return fmt.Errorf("controller private mode is still enabled")
		}
	}

	return nil
}
