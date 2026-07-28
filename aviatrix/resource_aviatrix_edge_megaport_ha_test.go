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

func TestAccAviatrixEdgeMegaportHa_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_MEGAPORT_HA") == "yes" {
		t.Skip("Skipping Edge Megaport HA test as SKIP_EDGE_MEGAPORT_HA is set")
	}

	resourceName := "aviatrix_edge_megaport_ha.test"
	accountName := "acc-" + acctest.RandString(5)
	edgeMegaportUsername := "megaport-user-" + acctest.RandString(5)
	gwName := "gw-" + acctest.RandString(5)
	siteID := "site-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeMegaportHaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeMegaportHaBasic(accountName, edgeMegaportUsername, gwName, siteID, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeMegaportHaExists(resourceName),
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

func testAccEdgeMegaportHaBasic(edgeMegaportUsername, accountName, gwName, siteID, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name           = "%s"
	cloud_type             = 1048576
	edge_megaport_username = "%s"
}
resource "aviatrix_edge_megaport" "test" {
	account_name           = aviatrix_account.test_account.account_name
	gw_name                = "%s"
	site_id                = "%s"
	ztp_file_download_path = "%s"

	interfaces {
		logical_ifname = "wan0"
		ip_address     = "10.230.5.32/24"
		gateway_ip     = "10.230.5.100"
		wan_public_ip  = "64.71.24.221"
	}

	interfaces {
		logical_ifname = "lan0"
		ip_address     = "10.230.3.32/24"
	}

	interfaces {
		logical_ifname = "mgmt0"
		enable_dhcp    = false
		ip_address     = "172.16.15.162/20"
		gateway_ip     = "172.16.0.1"
	}
}
resource "aviatrix_edge_megaport_ha" "test" {
	primary_gw_name        = aviatrix_edge_megaport.test.gw_name
	ztp_file_download_path = "%[5]s"

	interfaces {
		logical_ifname = "wan0"
		ip_address     = "10.220.11.20/24"
		gateway_ip     = "10.220.11.0"
	}

	interfaces {
		logical_ifname = "lan0"
		ip_address     = "10.220.12.20/24"
		gateway_ip     = "10.220.12.2"
	}

	interfaces {
		logical_ifname = "mgmt0"
		enable_dhcp    = true
	}
}
 `, accountName, edgeMegaportUsername, gwName, siteID, path)
}

func testAccCheckEdgeMegaportHaExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge megaport ha not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge megaport ha id is set")
		}

		client := mustClient(testAccProvider.Meta())

		edgeMegaportHa, err := client.GetEdgeMegaportHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != nil {
			return err
		}
		if edgeMegaportHa.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge megaport ha")
		}
		return nil
	}
}

func testAccCheckEdgeMegaportHaDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_megaport_ha" {
			continue
		}

		_, err := client.GetEdgeMegaportHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("edge megaport ha still exists")
		}
	}

	return nil
}
