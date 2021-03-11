package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixTunnelDetectionTime_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_TUNNEL_DETECTION_TIME")
	if skipAcc == "yes" {
		t.Skip("Skipping Tunnel Detection Time test as SKIP_TUNNEL_DETECTION_TIME is set")
	}
	resourceName := "aviatrix_tunnel_detection_time.test_detection_time"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccTunnelDetectionTimeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelDetectionTimeBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccTunnelDetectionTimeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "detection_time", "40"),
					resource.TestCheckResourceAttr(resourceName, "aviatrix_entity", "Controller"),
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

func testAccTunnelDetectionTimeBasic() string {
	return `
resource "aviatrix_tunnel_detection_time" "test_detection_time" {
	detection_time = 40
	aviatrix_entity = "Controller"
}
`
}

func testAccTunnelDetectionTimeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("tunnel detection time ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no tunnel detection time ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		aviatrixEntity := rs.Primary.Attributes["aviatrix_entity"]
		_, err := client.GetTunnelDetectionTime(context.Background(), aviatrixEntity)
		if err != nil {
			return fmt.Errorf("failed to get tunnel detection time status for %s: %v", aviatrixEntity, err)
		}

		if "Controller" != rs.Primary.ID {
			return fmt.Errorf("tunnel detection time ID not found")
		}

		return nil
	}
}

func testAccTunnelDetectionTimeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_tunnel_detection_time" {
			continue
		}

		aviatrixEntity := rs.Primary.Attributes["aviatrix_entity"]
		detectionTime, err := client.GetTunnelDetectionTime(context.Background(), aviatrixEntity)
		if err != nil || detectionTime != 60 {
			return fmt.Errorf("tunnel detection time configured when it should be destroyed")
		}
	}

	return nil
}
