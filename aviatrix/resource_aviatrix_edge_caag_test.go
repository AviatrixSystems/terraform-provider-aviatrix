package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixEdgeCaag_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_CAAG") == "yes" {
		t.Skip("Skipping Edge as a CaaG test as SKIP_EDGE_CAAG is set")
	}

	resourceName := "aviatrix_edge_caag.test"
	name := "edge-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeCaagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeCaagBasic(name, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeCaagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "management_interface_config", "DHCP"),
					resource.TestCheckResourceAttr(resourceName, "management_egress_ip_prefix", "10.60.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "wan_interface_ip_prefix", "10.60.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "wan_default_gateway_ip", "10.60.0.0"),
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

func testAccEdgeCaagBasic(name string, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_edge_caag" "test" {
	name                        = "%s"
	management_interface_config = "DHCP"
	management_egress_ip_prefix = "10.60.0.0/24"
	wan_interface_ip_prefix     = "10.60.0.0/24"
	wan_default_gateway_ip      = "10.60.0.0"
	lan_interface_ip_prefix     = "10.60.0.0/24"
	ztp_file_type               = "iso"
	ztp_file_download_path      = "%s"
}
   `, name, path)
}

func testAccCheckEdgeCaagExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge as a aaag not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge as a caag id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeCaag, err := client.GetEdgeCaag(context.Background(), rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}
		if edgeCaag.Name != rs.Primary.ID {
			return fmt.Errorf("could not find edge as a caaG")
		}
		return nil
	}
}

func testAccCheckEdgeCaagDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_caag" {
			continue
		}

		_, err := client.GetEdgeCaag(context.Background(), rs.Primary.Attributes["name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge as a caag still exists")
		}
	}
	return nil
}
