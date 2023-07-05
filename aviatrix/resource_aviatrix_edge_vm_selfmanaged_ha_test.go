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

func TestAccAviatrixEdgeVmSelfmanagedHa_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_VM_SELFMANAGED_HA") == "yes" {
		t.Skip("Skipping Edge VM Selfmanaged HA test as SKIP_EDGE_VM_SELFMANAGED_HA is set")
	}

	resourceName := "aviatrix_edge_vm_selfmanaged_ha.test"
	gwName := "gw-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeVmSelfmanagedHaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeVmSelfmanagedHaBasic(gwName, siteId, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeVmSelfmanagedHaExists(resourceName),
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

func testAccEdgeVmSelfmanagedHaBasic(gwName, siteId, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_edge_vm_selfmanaged" "test" {
	gw_name                = "%s"
	site_id                = "%s"
	ztp_file_type          = "iso"
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
resource "aviatrix_edge_vm_selfmanaged_ha" "test" {
	primary_gw_name        = aviatrix_edge_vm_selfmanaged.test.gw_name
	site_id                = aviatrix_edge_vm_selfmanaged.test.site_id
	ztp_file_type          = "iso"
	ztp_file_download_path = "%[3]s"

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
 `, gwName, siteId, path)
}

func testAccCheckEdgeVmSelfmanagedHaExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge vm selfmanaged ha not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge vm selfmanaged ha id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeVmSelfmanagedHa, err := client.GetEdgeVmSelfmanagedHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != nil {
			return err
		}
		if edgeVmSelfmanagedHa.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge vm selfmanaged ha")
		}
		return nil
	}
}

func testAccCheckEdgeVmSelfmanagedHaDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_vm_selfmanaged_ha" {
			continue
		}

		_, err := client.GetEdgeVmSelfmanagedHa(context.Background(), rs.Primary.Attributes["primary_gw_name"]+"-hagw")
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge vm selfmanaged ha still exists")
		}
	}

	return nil
}
