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

func TestAccAviatrixEdgeEquinix_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_EQUINIX") == "yes" {
		t.Skip("Skipping Edge Equinix test as SKIP_EDGE_EQUINIX is set")
	}

	resourceName := "aviatrix_edge_equinix.test"
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
		CheckDestroy: testAccCheckEdgeEquinixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeEquinixBasic(accountName, edgeEquinixUsername, gwName, siteId, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeEquinixExists(resourceName),
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
				ImportStateVerifyIgnore: []string{"ztp_file_download_path"},
			},
		},
	})
}

func testAccEdgeEquinixBasic(edgeEquinixUsername, accountName, gwName, siteId, path string) string {
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
 `, accountName, edgeEquinixUsername, gwName, siteId, path)
}

func testAccCheckEdgeEquinixExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge equinix not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge equinix id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeSpoke, err := client.GetEdgeEquinix(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge equinix")
		}
		return nil
	}
}

func testAccCheckEdgeEquinixDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_equinix" {
			continue
		}

		_, err := client.GetEdgeEquinix(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge equinix still exists")
		}
	}

	return nil
}
