package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "10.230.3.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address", "10.230.5.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.ip_address", "172.16.15.162/20"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.dns_server_ip", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.secondary_dns_server_ip", "9.9.9.9"),
					resource.TestCheckResourceAttr(resourceName, "bgp_polling_time", "50"),
					resource.TestCheckResourceAttr(resourceName, "bgp_neighbor_status_polling_time", "5"),
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
	gw_name                          = "%s"
	site_id                          = "%s"
	ztp_file_type                    = "iso"
	ztp_file_download_path           = "%s"
	bgp_polling_time                 = 50
	bgp_neighbor_status_polling_time = 5

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
		dns_server_ip = "8.8.8.8"
		secondary_dns_server_ip = "9.9.9.9"
	}
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

		client := mustClient(testAccProvider.Meta())

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
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_spoke" {
			continue
		}

		_, err := client.GetEdgeSpoke(context.Background(), rs.Primary.Attributes["gw_name"])
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("edge as a spoke still exists")
		}
	}

	return nil
}
