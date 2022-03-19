package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixCloudnEdgeGateway_basic(t *testing.T) {
	if os.Getenv("SKIP_CLOUDN_EDGE_GATEWAY") == "yes" {
		t.Skip("Skipping cloudn edge gateway test as SKIP_CLOUDN_EDGE_GATEWAY is set")
	}

	resourceName := "aviatrix_cloudn_edge_gateway.test"
	gwName := "edge-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudnEdgeGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudnEdgeGatewayBasic(gwName, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudnEdgeGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "management_connection_type", "DHCP"),
					resource.TestCheckResourceAttr(resourceName, "wan_interface_ip", "10.60.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "wan_default_gateway", "10.60.0.0"),
					resource.TestCheckResourceAttr(resourceName, "lan_interface_ip", "10.60.0.0/24"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_download_path"},
			},
		},
	})
}

func testAccCloudnEdgeGatewayBasic(gwName string, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_cloudn_edge_gateway" "test" {
	gw_name = "%s"
	management_connection_type = "DHCP"
	wan_interface_ip = "10.60.0.0/24"
	wan_default_gateway = "10.60.0.0"
	lan_interface_ip = "10.60.0.0/24"
	image_download_path = "%s"
}
    `, gwName, path)
}

func testAccCheckCloudnEdgeGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("cloudn edge gateway not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloudn edge gateway ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		cloudnEdgeGateway, err := client.GetCloudnEdgeGateway(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if cloudnEdgeGateway.GatewayName != rs.Primary.ID {
			return fmt.Errorf("cloudn edge gateway not found")
		}
		return nil
	}
}

func testAccCheckCloudnEdgeGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_cloudn_edge_gateway" {
			continue
		}

		_, err := client.GetCloudnEdgeGateway(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("cloudn edge gateway still exists")
		}
	}
	return nil
}
