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

func TestAccAviatrixEdgeMegaport_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_Megaport") == "yes" {
		t.Skip("Skipping Edge Megaport test as SKIP_EDGE_Megaport is set")
	}

	resourceName := "aviatrix_edge_Megaport.test"
	accountName := "acc-" + acctest.RandString(5)
	edgeMegaportUsername := "megaport-user-" + acctest.RandString(5)
	edgeMegaportPassword := "megaport-password-" + acctest.RandString(5)
	gwName := "gw-" + acctest.RandString(5)
	siteId := "site-" + acctest.RandString(5)
	path, _ := os.Getwd()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeMegaportDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeMegaportBasic(accountName, edgeMegaportUsername, edgeMegaportPassword, gwName, siteId, path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeMegaportExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", gwName),
					resource.TestCheckResourceAttr(resourceName, "site_id", siteId),
					resource.TestCheckResourceAttr(resourceName, "interfaces.2.ip_address", "192.168.99.14/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.3.ip_address", "192.168.88.14/24"),
					resource.TestCheckResourceAttr(resourceName, "interfaces.4.ip_address", "192.168.77.14/24"),
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

func testAccEdgeMegaportBasic(accountName, edgeMegaportUsername, edgeMegaportPassword, gwName, siteId, path string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name          = "%s"
	cloud_type            = 1048576
	edge_csp_username     = "%s"
	edge_csp_password     = "%s"
}
resource "aviatrix_edge_megaport" "test" {
	account_name                       = aviatrix_account.test.account_name
	gw_name                            = "%s"
	site_id                            = "%s"
	ztp_file_download_path             = "%s"
	
	interfaces {
        gateway_ip = "10.220.14.1"
        ip_address = "10.220.14.10/24"
        type       = "LAN"
        index      = 0
    }

    interfaces {
        enable_dhcp   = true
        type   = "MANAGEMENT"
        index  = 0
    }

    interfaces {
        gateway_ip    = "192.168.99.1"
        ip_address    = "192.168.99.14/24"
        type          = "WAN"
        index         = 0
        wan_public_ip     = "67.207.104.19"
    }

    interfaces {
        gateway_ip    = "192.168.88.1"
        ip_address    = "192.168.88.14/24"
        type          = "WAN"
        index         = 1
        wan_public_ip     = "67.71.12.148"
    }

    interfaces {
        gateway_ip  = "192.168.77.1"
        ip_address  = "192.168.77.14/24"
        type        = "WAN"
        index       = 2
        wan_public_ip   = "67.72.12.149"
    }
}
 `, accountName, edgeMegaportUsername, edgeMegaportPassword, gwName, siteId, path)
}

func testAccCheckEdgeMegaportExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge megaport not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge megaport id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		edgeSpoke, err := client.GetEdgeMegaport(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != nil {
			return err
		}
		if edgeSpoke.GwName != rs.Primary.ID {
			return fmt.Errorf("could not find edge megaport")
		}
		return nil
	}
}

func testAccCheckEdgeMegaportDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_megaport" {
			continue
		}

		_, err := client.GetEdgeMegaport(context.Background(), rs.Primary.Attributes["gw_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge megaport still exists")
		}
	}

	return nil
}
