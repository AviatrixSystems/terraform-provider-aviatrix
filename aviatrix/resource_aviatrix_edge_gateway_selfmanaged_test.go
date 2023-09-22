package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixEdgeGatewaySelfmanaged_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_GATEWAY_SELFMANAGED") == "yes" {
		t.Skip("Skipping Edge Gateway Selfmanaged test as SKIP_EDGE_GATEWAY_SELFMANAGED is set")
	}

	resourceName := "aviatrix_edge_gateway_selfmanaged.test"
	gwName := "edge-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeGatewaySelfmanagedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeGatewaySelfmanagedBasic(gwName, siteId, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeGatewaySelfmanagedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteId),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "10.230.5.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address", "10.230.3.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.ip_address", "172.16.15.162/20"),
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

func testAccEdgeGatewaySelfmanagedBasic(gwName, siteId, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_edge_gateway_selfmanaged" "test" {
	gw_name                = "%s"
	site_id                = "%s"
	ztp_file_type          = "iso"
	ztp_file_download_path = "%s"

	interfaces {
		name          = "eth0"
		type          = "WAN"
		ip_address    = "10.230.5.32/24"
		gateway_ip    = "10.230.5.100"
		wan_public_ip = "64.71.24.221"
	}
	
	interfaces {
		name       = "eth1"
		type       = "LAN"
		ip_address = "10.230.3.32/24"
	}
	
	interfaces {
		name        = "eth2"
		type        = "MANAGEMENT"
		enable_dhcp = false
		ip_address  = "172.16.15.162/20"
		gateway_ip  = "172.16.0.1"
	}
}
  `, gwName, siteId, path)
}

func testAccCheckEdgeGatewaySelfmanagedExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge gateway selfmanaged not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge gateway selfmanaged id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeSpoke, err := client.GetEdgeSpoke(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge gateway selfmanaged")
		}
		return nil
	}
}

func testAccCheckEdgeGatewaySelfmanagedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_gateway_selfmanaged" {
			continue
		}

		_, err := client.GetEdgeSpoke(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge gateway selfmanaged still exists")
		}
	}

	return nil
}
