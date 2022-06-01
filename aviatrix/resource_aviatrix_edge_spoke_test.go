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

func TestAccAviatrixEdgeSpoke_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_SPOKE") == "yes" {
		t.Skip("Skipping Edge as a Spoke test as SKIP_EDGE_SPOKE is set")
	}

	resourceName := "aviatrix_edge_spoke.test"
	gwName := "edge-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeSpokeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeSpokeBasic(gwName, siteId, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeSpokeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteId),
					resource.TestCheckResourceAttr(resourceName, "management_interface_config", "DHCP"),
					resource.TestCheckResourceAttr(resourceName, "wan_interface_ip_prefix", "10.60.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "wan_default_gateway_ip", "10.60.0.0"),
					resource.TestCheckResourceAttr(resourceName, "lan_interface_ip_prefix", "10.60.0.0/24"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ztp_file_type", "ztp_file_download_path"},
			},
		},
	})
}

func testAccEdgeSpokeBasic(gwName, siteId, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_edge_spoke" "test" {
	gw_name                     = "%s"
	site_id                     = "%s"
	management_interface_config = "DHCP"
	wan_interface_ip_prefix     = "10.60.0.0/24"
	wan_default_gateway_ip      = "10.60.0.0"
	lan_interface_ip_prefix     = "10.60.0.0/24"
	ztp_file_type               = "iso"
	ztp_file_download_path      = "%s"
}
  `, gwName, siteId, path)
}

func testAccCheckEdgeSpokeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge as a spoke not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge as a spoke id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeSpoke, err := client.GetEdgeSpoke(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge as a spoke")
		}
		return nil
	}
}

func testAccCheckEdgeSpokeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_spoke" {
			continue
		}

		_, err := client.GetEdgeSpoke(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge as a spoke still exists")
		}
	}

	return nil
}
