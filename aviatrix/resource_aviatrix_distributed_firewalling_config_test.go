package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDistributedFirewallingConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed-firewalling config tests as SKIP_DISTRIBUTED_FIREWALLING_CONFIG is set")
	}

	resourceName := "aviatrix_distributed_firewalling_config.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccDistributedFirewallingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDistributedFirewallingConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDistributedFirewallingConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_distributed_firewalling", "true"),
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

func testAccDistributedFirewallingConfigBasic() string {
	return `
resource "aviatrix_distributed_firewalling_config" "test" {
	enable_distributed_firewalling = true
}
	`
}

func testAccDistributedFirewallingConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("distributed-firewalling config Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no distributed-firewalling config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("distributed-firewalling config ID not found")
		}

		return nil
	}
}

func testAccDistributedFirewallingConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_distributed_firewalling_config" {
			continue
		}

		distributedFirewalling, err := client.GetDistributedFirewallingStatus(context.Background())
		if err != nil {
			return fmt.Errorf("could not retrieve distributed firewalling config")
		}
		if distributedFirewalling.EnableDistributedFirewalling {
			return fmt.Errorf("distributed firewalling is still enabled")
		}
	}

	return nil
}
