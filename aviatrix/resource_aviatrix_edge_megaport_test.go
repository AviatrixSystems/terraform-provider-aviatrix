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

func TestAccAviatrixEdgeMegaport_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_MEGAPORT") == "yes" {
		t.Skip("Skipping Edge Megaport test as SKIP_EDGE_MEGAPORT is set")
	}

	resourceName := "aviatrix_edge_megaport.test_spoke"
	accountName := "acc-" + acctest.RandString(5)
	gwName := "gw-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeMegaportDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeMegaportBasic(accountName, gwName, siteId, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeMegaportExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteId),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "10.220.14.10/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.logical_ifname", "lan0"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.tag", "LAN"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address", "192.168.99.14/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.logical_ifname", "wan0"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.tag", "WAN"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.ip_address", "192.168.88.14/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.logical_ifname", "wan1"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.3.ip_address", "192.168.77.14/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.3.logical_ifname", "wan2"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.4.logical_ifname", "mgmt0"),
					resource.TestCheckResourceAttr(resourceName, "vlan.0.parent_logical_interface_name", "lan0"),
					resource.TestCheckResourceAttr(resourceName, "vlan.0.vlan_id", "21"),
					resource.TestCheckResourceAttr(resourceName, "vlan.0.ip_address", "10.220.21.11/24"),
					resource.TestCheckResourceAttr(resourceName, "included_advertised_spoke_routes.0", "10.230.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "included_advertised_spoke_routes.1", "10.230.5.0/24"),
				),
			},
			{
				Config: testAccEdgeMegaportUpdate(accountName, gwName, siteId, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeMegaportExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "10.220.14.10/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address", "192.168.95.14/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.ip_address", "192.168.85.14/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.3.ip_address", "192.168.75.14/24"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ztp_file_download_path"},
			},
		},
	})
}

func testAccEdgeMegaportBasic(accountName, gwName, siteId, path string) string {
	return fmt.Sprintf(`
	resource "aviatrix_account" "test" {
		account_name          = "%s"
		cloud_type            = 1048576
	}
	resource "aviatrix_edge_megaport" "test_spoke" {
		account_name                       = aviatrix_account.test.account_name
		gw_name                            = "%s"
		site_id                            = "%s"
		ztp_file_download_path             = "%s"

		interfaces {
			gateway_ip     = "10.220.14.1"
			ip_address     = "10.220.14.10/24"
			logical_ifname = "lan0"
			tag            = "LAN"
		}

		interfaces {
			gateway_ip     = "192.168.99.1"
			ip_address     = "192.168.99.14/24"
			logical_ifname = "wan0"
			wan_public_ip  = "67.207.104.19"
			tag            = "WAN"
		}

		interfaces {
			gateway_ip     = "192.168.88.1"
			ip_address     = "192.168.88.14/24"
			logical_ifname = "wan1"
			wan_public_ip  = "67.71.12.148"
		}

		interfaces {
			gateway_ip     = "192.168.77.1"
			ip_address     = "192.168.77.14/24"
			logical_ifname = "wan2"
			wan_public_ip  = "67.72.12.149"
		}

		interfaces {
			enable_dhcp   = true
			logical_ifname = "mgmt0"
		}

		vlan {
			parent_logical_interface_name = "lan0"
			vlan_id                        = 21
			ip_address                     = "10.220.21.11/24"
		}

	    included_advertised_spoke_routes = [
			"10.230.3.0/24",
			"10.230.5.0/24"
		]
	}
 `, accountName, gwName, siteId, path)
}

func testAccEdgeMegaportUpdate(accountName, gwName, siteID, path string) string {
	return fmt.Sprintf(`
	resource "aviatrix_account" "test" {
		account_name          = "%s"
		cloud_type            = 1048576
	}
	resource "aviatrix_edge_megaport" "test_spoke" {
		account_name                       = aviatrix_account.test.account_name
		gw_name                            = "%s"
		site_id                            = "%s"
		ztp_file_download_path             = "%s"

		interfaces {
			gateway_ip     = "10.220.15.1"
			ip_address     = "10.220.15.10/24"
			logical_ifname = "lan0"
		}

		interfaces {
			gateway_ip     = "192.168.95.1"
			ip_address     = "192.168.95.14/24"
			logical_ifname = "wan0"
			wan_public_ip  = "67.207.104.19"
		}

		interfaces {
			gateway_ip     = "192.168.85.1"
			ip_address     = "192.168.85.14/24"
			logical_ifname = "wan1"
			wan_public_ip  = "67.71.12.148"
		}

		interfaces {
			gateway_ip     = "192.168.75.1"
			ip_address     = "192.168.75.14/24"
			logical_ifname = "wan2"
			wan_public_ip  = "67.72.12.149"
		}

		interfaces {
			enable_dhcp   = true
			logical_ifname = "mgmt0"
		}

		vlan {
			parent_logical_interface_name = "lan0"
			vlan_id                        = 21
			ip_address                     = "10.220.21.11/24"
		}
	}
 `, accountName, gwName, siteID, path)
}

func testAccCheckEdgeMegaportExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge megaport not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge megaport id is set")
		}

		client := mustClient(testAccProvider.Meta())

		edgeSpoke, err := client.GetEdgeMegaport(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge megaport")
		}
		return nil
	}
}

func testAccCheckEdgeMegaportDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_megaport" {
			continue
		}

		_, err := client.GetEdgeMegaport(context.Background(), rs.Primary.Attributes["gw_name"])
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("edge megaport still exists")
		}
	}

	return nil
}
