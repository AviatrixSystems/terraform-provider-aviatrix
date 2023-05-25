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

func TestAccAviatrixEdgeZededaHa_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_ZEDEDA_HA") == "yes" {
		t.Skip("Skipping Edge Zededa HA test as SKIP_EDGE_ZEDEDA_HA is set")
	}

	resourceName := "aviatrix_edge_zededa_ha.test"
	accountName := "edge-zededa-acc-" + acctest.RandString(5)
	gwName := "edge-zededa-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeZededaHaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeZededaHaBasic(accountName, gwName, siteId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeZededaHaExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "10.220.11.20/24"),
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

func testAccEdgeZededaHaBasic(accountName, gwName, siteId string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
 	account_name         = "%s"
	cloud_type           = 65536
	edge_zededa_username = "%s"
	edge_zededa_password = "%s"
}
resource "aviatrix_edge_zededa" "test" {
	account_name      = aviatrix_account.test_account.account_name
	gw_name           = "%s"
	site_id           = "%s"
 	project_uuid      = "%s"
 	compute_node_uuid = "%s"
 	template_uuid     = "%s"

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
resource "aviatrix_edge_zededa_ha" "test" {
	primary_gw_name   = aviatrix_edge_zededa.test.gw_name
	compute_node_uuid = "%s"

	interfaces {
		name       = "eth1"
		type       = "LAN"
		ip_address = "10.220.11.20/24"
		gateway_ip = "10.220.11.1"
	}

	interfaces {
		name        = "eth2"
		type        = "MANAGEMENT"
		enable_dhcp = true
	}
}
 `, accountName, os.Getenv("EDGE_ZEDEDA_USERNAME"), os.Getenv("EDGE_ZEDEDA_PASSWORD"), gwName, siteId,
		os.Getenv("EDGE_ZEDEDA_PROJECT_UUID"), os.Getenv("EDGE_ZEDEDA_COMPUTE_NODE_UUID"),
		os.Getenv("EDGE_ZEDEDA_TEMPLATE_UUID"), os.Getenv("EDGE_ZEDEDA_HA_COMPUTE_NODE_UUID"))
}

func testAccCheckEdgeZededaHaExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge zededa ha not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge zededa ha id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeCSPHa, err := client.GetEdgeCSPHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != nil {
			return err
		}
		if edgeCSPHa.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge zededa ha")
		}
		return nil
	}
}

func testAccCheckEdgeZededaHaDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_zededa_ha" {
			continue
		}

		_, err := client.GetEdgeCSPHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge zededa ha still exists")
		}
	}

	return nil
}
