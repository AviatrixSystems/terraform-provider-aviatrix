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

func TestAccAviatrixEdgePlatform_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_PLATFORM") == "yes" {
		t.Skip("Skipping Edge Platform test as SKIP_EDGE_PLATFORM is set")
	}

	resourceName := "aviatrix_edge_platform.test"
	accountName := "acc-" + acctest.RandString(5)
	deviceName := "device-" + acctest.RandString(5)
	gwName := "gw-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgePlatformDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgePlatformBasic(accountName, deviceName, gwName, siteId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgePlatformExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", os.Getenv("EDGE_NEO_SITE_ID")),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "10.230.5.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address", "10.230.3.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.ip_address", "172.16.15.162/20"),
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

func testAccEdgePlatformBasic(accountName, deviceName, gwName, siteId string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name = "%s"
	cloud_type   = 262144
}
resource "aviatrix_edge_platform_device_onboarding" "test" {
	account_name   = aviatrix_account.test.account_name
	device_name    = "%s"
	serial_number  = "%s"
	hardware_model = "%s"
}
resource "aviatrix_edge_platform" "test" {
	account_name               = aviatrix_account.test.account_name
	gw_name                    = "%s"
	site_id                    = "%s"
	device_id                  = aviatrix_edge_platform_device_onboarding.test.device_id
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
 `, accountName, deviceName, os.Getenv("EDGE_PLATFORM_DEVICE_SERIAL_NUMBER"), os.Getenv("EDGE_PLATFORM_DEVICE_HARDWARE_MODEL"),
		gwName, siteId)
}

func testAccCheckEdgePlatformExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge platform not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge platform id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeSpoke, err := client.GetEdgeNEO(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge platform")
		}
		return nil
	}
}

func testAccCheckEdgePlatformDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_platform" {
			continue
		}

		_, err := client.GetEdgeNEO(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge platform still exists")
		}
	}

	return nil
}
