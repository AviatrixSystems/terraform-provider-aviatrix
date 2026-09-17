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

func TestAccAviatrixEdgeGatewaySelfmanaged_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_GATEWAY_SELFMANAGED") == "yes" {
		t.Skip("Skipping Edge Gateway Selfmanaged test as SKIP_EDGE_GATEWAY_SELFMANAGED is set")
	}

	resourceName := "aviatrix_edge_gateway_selfmanaged.test"
	gwName := "edge-" + acctest.RandString(5)
	siteID := "site-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeGatewaySelfmanagedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeGatewaySelfmanagedBasic(gwName, siteID, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeGatewaySelfmanagedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteID),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "172.16.15.162/20"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address", "10.230.3.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.ip_address", "10.230.5.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.dns_server_ip", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.secondary_dns_server_ip", "9.9.9.9"),
					resource.TestCheckResourceAttr(resourceName, "bgp_polling_time", "50"),
					resource.TestCheckResourceAttr(resourceName, "bgp_neighbor_status_polling_time", "5"),
					resource.TestCheckResourceAttr(resourceName, "included_advertised_spoke_routes.0", "10.231.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "included_advertised_spoke_routes.1", "10.232.5.0/24"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_encryption_cipher", "strong"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_forward_secrecy", "enable"),
				),
			},
			{
				Config: testAccEdgeGatewaySelfmanagedUpdate(gwName, siteID, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeGatewaySelfmanagedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteID),
					resource.TestCheckResourceAttr(resourceName, "tunnel_encryption_cipher", "default"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_forward_secrecy", "disable"),
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

func testAccEdgeGatewaySelfmanagedBasic(gwName, siteID, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_edge_gateway_selfmanaged" "test" {
	gw_name                            = "%s"
	site_id                            = "%s"
	ztp_file_type                      = "iso"
	ztp_file_download_path             = "%s"
	bgp_polling_time                   = 50
	bgp_neighbor_status_polling_time   = 5
	tunnel_encryption_cipher           = "strong"
	tunnel_forward_secrecy             = "enable"

	interfaces {
		name          = "eth0"
		type          = "WAN"
		ip_address    = "10.230.5.32/24"
		gateway_ip    = "10.230.5.100"
		wan_public_ip = "64.71.24.221"
		dns_server_ip = "8.8.8.8"
		secondary_dns_server_ip = "9.9.9.9"
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

	included_advertised_spoke_routes = [
		"10.231.3.0/24",
		"10.232.5.0/24"
	]
}
  `, gwName, siteID, path)
}

func testAccEdgeGatewaySelfmanagedUpdate(gwName, siteID, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_edge_gateway_selfmanaged" "test" {
	gw_name                            = "%s"
	site_id                            = "%s"
	ztp_file_type                      = "iso"
	ztp_file_download_path             = "%s"
	bgp_polling_time                   = 50
	bgp_neighbor_status_polling_time   = 5
	tunnel_encryption_cipher           = "default"
	tunnel_forward_secrecy             = "disable"

	interfaces {
		name          = "eth0"
		type          = "WAN"
		ip_address    = "10.230.5.32/24"
		gateway_ip    = "10.230.5.100"
		wan_public_ip = "64.71.24.221"
		dns_server_ip = "8.8.8.8"
		secondary_dns_server_ip = "9.9.9.9"
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

	included_advertised_spoke_routes = [
		"10.230.3.0/24",
		"10.230.5.0/24"
	]
}
  `, gwName, siteID, path)
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

		client := mustClient(testAccProvider.Meta())

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
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_gateway_selfmanaged" {
			continue
		}

		_, err := client.GetEdgeSpoke(context.Background(), rs.Primary.Attributes["gw_name"])
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("edge gateway selfmanaged still exists")
		}
	}

	return nil
}

func TestAccAviatrixEdgeGatewaySelfmanaged_tunnelPolicies(t *testing.T) {
	if os.Getenv("SKIP_EDGE_GATEWAY_SELFMANAGED") == "yes" {
		t.Skip("Skipping Edge Gateway Selfmanaged test as SKIP_EDGE_GATEWAY_SELFMANAGED is set")
	}

	resourceName := "aviatrix_edge_gateway_selfmanaged.test_tunnel"
	gwName := "edge-tunnel-" + acctest.RandString(5)
	siteID := "site-tunnel-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeGatewaySelfmanagedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeGatewaySelfmanagedTunnelPoliciesConfig(gwName, siteID, path, "default", "disable"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeGatewaySelfmanagedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteID),
					resource.TestCheckResourceAttr(resourceName, "tunnel_encryption_cipher", "default"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_forward_secrecy", "disable"),
				),
			},
			{
				Config: testAccEdgeGatewaySelfmanagedTunnelPoliciesConfig(gwName, siteID, path, "strong", "enable"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeGatewaySelfmanagedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteID),
					resource.TestCheckResourceAttr(resourceName, "tunnel_encryption_cipher", "strong"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_forward_secrecy", "enable"),
				),
			},
			{
				Config: testAccEdgeGatewaySelfmanagedTunnelPoliciesConfig(gwName, siteID, path, "default", "disable"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeGatewaySelfmanagedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteID),
					resource.TestCheckResourceAttr(resourceName, "tunnel_encryption_cipher", "default"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_forward_secrecy", "disable"),
				),
			},
		},
	})
}

func testAccEdgeGatewaySelfmanagedTunnelPoliciesConfig(gwName, siteID, path, encCipher, forwardSecrecy string) string {
	return fmt.Sprintf(`
resource "aviatrix_edge_gateway_selfmanaged" "test_tunnel" {
	gw_name                            = "%s"
	site_id                            = "%s"
	ztp_file_type                      = "iso"
	ztp_file_download_path             = "%s"
	tunnel_encryption_cipher           = "%s"
	tunnel_forward_secrecy             = "%s"

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
  `, gwName, siteID, path, encCipher, forwardSecrecy)
}
