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
					resource.TestCheckResourceAttr(resourceName, "management_interface_config", "DHCP"),
					resource.TestCheckResourceAttr(resourceName, "lan_interface_ip_prefix", "10.60.0.0/24"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ztp_file_type", "ztp_file_download_path"},
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
	account_name                = aviatrix_account.test_account.account_name
	gw_name                     = "%s"
	site_id                     = "%s"
 	project_uuid                = "%s"
 	compute_node_uuid           = "%s"
 	template_uuid               = "%s"
	management_interface_config = "DHCP"
	lan_interface_ip_prefix     = "10.60.0.0/24"
}
 `, accountName, os.Getenv("EDGE_CSP_USERNAME"), os.Getenv("EDGE_CSP_PASSWORD"), gwName, siteId,
		os.Getenv("EDGE_CSP_PROJECT_UUID"), os.Getenv("EDGE_CSP_COMPUTE_UUID"), os.Getenv("EDGE_CSP_TEMPLATE_UUID"))
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

		client := testAccProvider.Meta().(*goaviatrix.Client)

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
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_csp" {
			continue
		}

		_, err := client.GetEdgeCSP(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge csp still exists")
		}
	}

	return nil
}
