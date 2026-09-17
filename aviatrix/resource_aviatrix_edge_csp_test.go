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

func TestAccAviatrixEdgeCSP_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_CSP") == "yes" {
		t.Skip("Skipping Edge CSP test as SKIP_EDGE_CSP is set")
	}

	resourceName := "aviatrix_edge_csp.test"
	accountName := "edge-csp-acc-" + acctest.RandString(5)
	gwName := "edge-csp-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeCSPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeCSPBasic(accountName, gwName, siteId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeCSPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteId),
					resource.TestCheckResourceAttr(resourceName, "interfaces.0.ip_address", "10.230.5.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.1.ip_address", "10.230.3.32/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.ip_address", "172.16.15.162/20"),
					resource.TestCheckResourceAttr(resourceName, "bgp_polling_time", "50"),
					resource.TestCheckResourceAttr(resourceName, "bgp_neighbor_status_polling_time", "5"),
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

func testAccEdgeCSPBasic(accountName, gwName, siteId string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
 	account_name      = "%s"
	cloud_type        = 65536
	edge_csp_username = "%s"
	edge_csp_password = "%s"
}
resource "aviatrix_edge_csp" "test" {
	account_name                     = aviatrix_account.test_account.account_name
	gw_name                          = "%s"
	site_id                          = "%s"
 	project_uuid                     = "%s"
 	compute_node_uuid                = "%s"
 	template_uuid                    = "%s"
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
	}
}
 `, accountName, os.Getenv("EDGE_CSP_USERNAME"), os.Getenv("EDGE_CSP_PASSWORD"), gwName, siteId,
		os.Getenv("EDGE_CSP_PROJECT_UUID"), os.Getenv("EDGE_CSP_COMPUTE_NODE_UUID"), os.Getenv("EDGE_CSP_TEMPLATE_UUID"))
}

func testAccCheckEdgeCSPExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge csp not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge csp id is set")
		}

		client := mustClient(testAccProvider.Meta())

		edgeSpoke, err := client.GetEdgeCSP(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge csp")
		}
		return nil
	}
}

func testAccCheckEdgeCSPDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_csp" {
			continue
		}

		_, err := client.GetEdgeCSP(context.Background(), rs.Primary.Attributes["gw_name"])
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("edge csp still exists")
		}
	}

	return nil
}
