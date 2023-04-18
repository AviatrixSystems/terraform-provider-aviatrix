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

func TestAccAviatrixEdgeEquinixHa_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_EQUINIX_HA") == "yes" {
		t.Skip("Skipping Edge Equinix HA test as SKIP_EDGE_EQUINIX_HA is set")
	}

	resourceName := "aviatrix_edge_equinix_ha.test"
	accountName := "acc-" + acctest.RandString(5)
	edgeEquinixUsername := "equinix-user-" + acctest.RandString(5)
	gwName := "gw-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeEquinixHaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeEquinixHaBasic(accountName, edgeEquinixUsername, gwName, siteId, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeEquinixHaExists(resourceName),
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

func testAccEdgeEquinixHaBasic(edgeEquinixUsername, accountName, gwName, siteId, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name          = "%s"
	cloud_type            = 524288
	edge_equinix_username = "%s"
}
resource "aviatrix_edge_equinix" "test" {
	account_name           = aviatrix_account.test_account.account_name
	gw_name                = "%s"
	site_id                = "%s"
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
resource "aviatrix_edge_equinix_ha" "test" {
	primary_gw_name = aviatrix_edge_equinix.test.gw_name
	ztp_file_download_path = "%[5]s"

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
 `, accountName, edgeEquinixUsername, gwName, siteId, path)
}

func testAccCheckEdgeEquinixHaExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge equinix ha not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge equinix ha id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeEquinixHa, err := client.GetEdgeEquinixHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != nil {
			return err
		}
		if edgeEquinixHa.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge equinix ha")
		}
		return nil
	}
}

func testAccCheckEdgeEquinixHaDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_equinix_ha" {
			continue
		}

		_, err := client.GetEdgeEquinixHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge equinix ha still exists")
		}
	}

	return nil
}
