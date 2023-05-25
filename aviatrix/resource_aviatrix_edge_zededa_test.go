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

func TestAccAviatrixEdgeZededa_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_ZEDEDA") == "yes" {
		t.Skip("Skipping Edge Zededa test as SKIP_EDGE_ZEDEDA is set")
	}

	resourceName := "aviatrix_edge_zededa.test"
	accountName := "edge-zededa-acc-" + acctest.RandString(5)
	gwName := "edge-zededa-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeZededaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeZededaBasic(accountName, gwName, siteId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeZededaExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteId),
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

func testAccEdgeZededaBasic(accountName, gwName, siteId string) string {
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
 `, accountName, os.Getenv("EDGE_ZEDEDA_USERNAME"), os.Getenv("EDGE_ZEDEDA_PASSWORD"), gwName, siteId,
		os.Getenv("EDGE_ZEDEDA_PROJECT_UUID"), os.Getenv("EDGE_ZEDEDA_COMPUTE_NODE_UUID"), os.Getenv("EDGE_ZEDEDA_TEMPLATE_UUID"))
}

func testAccCheckEdgeZededaExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge zededa not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge zededa id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeSpoke, err := client.GetEdgeCSP(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge zededa")
		}
		return nil
	}
}

func testAccCheckEdgeZededaDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_zededa" {
			continue
		}

		_, err := client.GetEdgeCSP(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge zededa still exists")
		}
	}

	return nil
}
