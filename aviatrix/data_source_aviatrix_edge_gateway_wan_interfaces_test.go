package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixEdgeGatewayWanInterfaces_basic(t *testing.T) {
	resourceName := "aviatrix_edge_gateway_wan_interfaces.test"
	accountName := "edge-csp-acc-" + acctest.RandString(5)
	gwName := "edge-csp-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_DATA_EDGE_GATEWAY_WAN_INTERFACES")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Edge Gateway WAN Interfaces tests as SKIP_DATA_EDGE_GATEWAY_WAN_INTERFACES is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixEdgeGatewayWanInterfacesConfigBasic(accountName, gwName, siteId),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixEdgeGatewayWanInterfaces(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "wan_interfaces"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixEdgeGatewayWanInterfacesConfigBasic(accountName, gwName, siteId string) string {
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
data "aviatrix_edge_gateway_wan_interfaces" "test" {
  gw_name = aviatrix_edge_csp.test.gw_name
}
 `, accountName, os.Getenv("EDGE_CSP_USERNAME"), os.Getenv("EDGE_CSP_PASSWORD"), gwName, siteId,
		os.Getenv("EDGE_CSP_PROJECT_UUID"), os.Getenv("EDGE_CSP_COMPUTE_NODE_UUID"), os.Getenv("EDGE_CSP_TEMPLATE_UUID"))
}

func testAccDataSourceAviatrixEdgeGatewayWanInterfaces(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
