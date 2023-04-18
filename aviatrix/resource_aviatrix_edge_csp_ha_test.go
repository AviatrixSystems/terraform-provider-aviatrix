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

func TestAccAviatrixEdgeCSPHa_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_CSP_HA") == "yes" {
		t.Skip("Skipping Edge CSP HA test as SKIP_EDGE_CSP_HA is set")
	}

	resourceName := "aviatrix_edge_csp_ha.test"
	accountName := "edge-csp-acc-" + acctest.RandString(5)
	gwName := "edge-csp-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeCSPHaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeCSPHaBasic(accountName, gwName, siteId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeCSPHaExists(resourceName),
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

func testAccEdgeCSPHaBasic(accountName, gwName, siteId string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
 	account_name      = "%s"
	cloud_type        = 65536
	edge_csp_username = "%s"
	edge_csp_password = "%s"
}
resource "aviatrix_edge_csp" "test" {
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
resource "aviatrix_edge_csp_ha" "test" {
	primary_gw_name   = aviatrix_edge_csp.test.gw_name
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
 `, accountName, os.Getenv("EDGE_CSP_USERNAME"), os.Getenv("EDGE_CSP_PASSWORD"), gwName, siteId,
		os.Getenv("EDGE_CSP_PROJECT_UUID"), os.Getenv("EDGE_CSP_COMPUTE_NODE_UUID"),
		os.Getenv("EDGE_CSP_TEMPLATE_UUID"), os.Getenv("EDGE_CSP_HA_COMPUTE_NODE_UUID"))
}

func testAccCheckEdgeCSPHaExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge csp ha not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge csp ha id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeCSPHa, err := client.GetEdgeCSPHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != nil {
			return err
		}
		if edgeCSPHa.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge csp ha")
		}
		return nil
	}
}

func testAccCheckEdgeCSPHaDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_csp_ha" {
			continue
		}

		_, err := client.GetEdgeCSPHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge csp ha still exists")
		}
	}

	return nil
}
