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

func TestAccAviatrixEdgeNEOHa_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_NEO_HA") == "yes" {
		t.Skip("Skipping Edge NEO HA test as SKIP_EDGE_NEO_HA is set")
	}

	resourceName := "aviatrix_edge_neo_ha.test"
	accountName := "acc-" + acctest.RandString(5)
	deviceName := "device-" + acctest.RandString(5)
	haDeviceName := "ha-device-" + acctest.RandString(5)
	gwName := "gw-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeNEOHaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeNEOHaBasic(accountName, deviceName, haDeviceName, gwName, siteId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeNEOHaExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "10.220.11.20/24"),
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

func testAccEdgeNEOHaBasic(accountName, deviceName, haDeviceName, gwName, siteId string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name = "%s"
	cloud_type   = 262144
}
resource "aviatrix_edge_neo_device_onboarding" "test" {
	account_name   = aviatrix_account.test.account_name
	device_name    = "%s"
	serial_number  = "%s"
	hardware_model = "%s"
}
resource "aviatrix_edge_neo_device_onboarding" "test_ha" {
	account_name   = aviatrix_account.test.account_name
	device_name    = "%s"
	serial_number  = "%s"
	hardware_model = "%s"
}
resource "aviatrix_edge_neo" "test" {
	account_name               = aviatrix_account.test.account_name
	gw_name                    = "%s"
	site_id                    = "%s"
	device_id                  = aviatrix_edge_neo_device_onboarding.test.device_id
	gw_size                    = "small"
	management_interface_names = ["eth2"]
	lan_interface_names        = ["eth1"]
	wan_interface_names        = ["eth0"]

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
resource "aviatrix_edge_neo_ha" "test" {
	primary_gw_name = aviatrix_edge_neo.test.gw_name
	device_id       = aviatrix_edge_neo_device_onboarding.test_ha.device_id

	interfaces {
		name       = "eth0"
		type       = "WAN"
		ip_address = "10.220.11.20/24"
		gateway_ip = "10.220.11.0"
	}

	interfaces {
		name       = "eth1"
		type       = "LAN"
		ip_address = "10.220.12.20/24"
		gateway_ip = "10.220.12.2"
	}

	interfaces {
		name        = "eth2"
		type        = "MANAGEMENT"
		enable_dhcp = true
	}
}
 `, accountName, deviceName, os.Getenv("EDGE_NEO_DEVICE_SERIAL_NUMBER"), os.Getenv("EDGE_NEO_DEVICE_HARDWARE_MODEL"),
		haDeviceName, os.Getenv("EDGE_NEO_HA_DEVICE_SERIAL_NUMBER"), os.Getenv("EDGE_NEO_HA_DEVICE_HARDWARE_MODEL"),
		gwName, siteId)
}

func testAccCheckEdgeNEOHaExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge neo ha not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge neo ha id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeNEOHa, err := client.GetEdgeNEOHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != nil {
			return err
		}
		if edgeNEOHa.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge neo ha")
		}
		return nil
	}
}

func testAccCheckEdgeNEOHaDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_neo_ha" {
			continue
		}

		_, err := client.GetEdgeNEOHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge neo ha still exists")
		}
	}

	return nil
}
