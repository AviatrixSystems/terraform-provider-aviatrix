package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixEdgeSpokeTransitAttachment_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_edge_spoke_transit_attachment.test"

	accountName := "megaport-" + rName
	spokeGwName := "spoke-" + rName
	spokeSiteID := "site-" + rName
	path, _ := os.Getwd()
	transitGwName := "transit-" + rName
	transitSiteID := "site-" + rName
	resourceNameEdge := "aviatrix_edge_spoke_transit_attachment.test_edge_spoke_transit_attachment"

	skipAcc := os.Getenv("SKIP_EDGE_SPOKE_TRANSIT_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping Edge as a Spoke transit attachment tests as 'SKIP_EDGE_SPOKE_TRANSIT_ATTACHMENT' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preEdgeSpokeTransitAttachmentCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeSpokeTransitAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeSpokeTransitAttachmentConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeSpokeTransitAttachmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "spoke_gw_name", fmt.Sprintf("tfs-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "transit_gw_name", fmt.Sprintf("tft-%s", rName)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEdgeSpokeTransitAttachmentConfigEdge(accountName, spokeGwName, spokeSiteID, path, transitGwName, transitSiteID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeSpokeTransitAttachmentExists(resourceNameEdge),
					resource.TestCheckResourceAttr(resourceNameEdge, "spoke_gw_name", spokeGwName),
					resource.TestCheckResourceAttr(resourceNameEdge, "transit_gw_name", transitGwName),
					resource.TestCheckResourceAttr(resourceNameEdge, "enable_over_private_network", "true"),
					resource.TestCheckResourceAttr(resourceNameEdge, "enable_jumbo_frame", "false"),
					resource.TestCheckResourceAttr(resourceNameEdge, "enable_insane_mode", "true"),
					resource.TestCheckResourceAttr(resourceNameEdge, "spoke_gateway_logical_ifnames.0", "wan1"),
					resource.TestCheckResourceAttr(resourceNameEdge, "transit_gateway_logical_ifnames.0", "wan1"),
				),
			},
		},
	})
}

func preEdgeSpokeTransitAttachmentCheck(t *testing.T) {
	if os.Getenv("EDGE_SPOKE_NAME") == "" {
		t.Fatal("Environment variable EDGE_SPOKE_NAME is not set")
	}
}

func testAccEdgeSpokeTransitAttachmentConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tft-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_edge_spoke_transit_attachment" "test" {
	spoke_gw_name   = "%s"
	transit_gw_name = aviatrix_transit_gateway.test.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"),
		os.Getenv("EDGE_SPOKE_NAME"))
}

func testAccEdgeSpokeTransitAttachmentConfigEdge(accountName, spokeGwName, spokeSiteID, path, transitGwName, transitSiteID string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_edge_megaport" {
	account_name       = "edge-%s"
	cloud_type         = 1048576
}
resource "aviatrix_transit_gateway" "test_edge_transit" {
	cloud_type   = 1048576
	account_name = aviatrix_account.test_acc_edge_megaport.account_name
	gw_name      = "%s"
	vpc_id       = "%s"
	gw_size      = "SMALL"
	ztp_file_download_path = "%s"
	interfaces {
        gateway_ip     = "192.168.20.1"
        ip_address     = "192.168.20.11/24"
        public_ip      = "67.207.104.19"
        logical_ifname = "wan0"
        secondary_private_cidr_list = ["192.168.20.16/29"]
    }

    interfaces {
        gateway_ip     = "192.168.21.1"
        ip_address     = "192.168.21.11/24"
        public_ip      = "67.71.12.148"
        logical_ifname = "wan1"
        secondary_private_cidr_list = ["192.168.21.16/29"]
    }

    interfaces {
        dhcp           = true
        logical_ifname = "mgmt0"
    }

    interfaces {
        gateway_ip     = "192.168.22.1"
        ip_address     = "192.168.22.11/24"
        logical_ifname = "wan2"
    }

    interfaces {
        gateway_ip     = "192.168.23.1"
        ip_address     = "192.168.23.11/24"
        logical_ifname = "wan3"
    }
}

resource "aviatrix_edge_megaport" "test_edge_spoke" {
	account_name                       = aviatrix_account.test_acc_edge_megaport.account_name
	gw_name                            = "%s"
	site_id                            = "%s"
	ztp_file_download_path             = "%s"

	interfaces {
		gateway_ip     = "10.220.14.1"
		ip_address     = "10.220.14.10/24"
		logical_ifname = "lan0"
	}

	interfaces {
		gateway_ip     = "192.168.99.1"
		ip_address     = "192.168.99.14/24"
		logical_ifname = "wan0"
		wan_public_ip  = "67.207.104.19"
	}

	interfaces {
		gateway_ip     = "192.168.88.1"
		ip_address     = "192.168.88.14/24"
		logical_ifname = "wan1"
		wan_public_ip  = "67.71.12.148"
	}

	interfaces {
		gateway_ip     = "192.168.77.1"
		ip_address     = "192.168.77.14/24"
		logical_ifname = "wan2"
		wan_public_ip  = "67.72.12.149"
	}

	interfaces {
		enable_dhcp   = true
		logical_ifname = "mgmt0"
	}

	vlan {
		parent_logical_interface_name = "lan0"
		vlan_id                        = 21
		ip_address                     = "10.220.21.11/24"
	}
}

resource "aviatrix_edge_spoke_transit_attachment" "test_edge_spoke_transit_attachment" {
	spoke_gw_name   = "aviatrix_edge_megaport.test_edge_spoke.gw_name"
	transit_gw_name = "aviatrix_transit_gateway.test_edge_transit.gw_name"
	enable_over_private_network = true
	enable_jumbo_frame = false
	enable_insane_mode = true
	spoke_gateway_logical_ifnames = ["wan1"]
	transit_gateway_logical_ifnames = ["wan1"]
  }
	`, accountName, spokeGwName, spokeSiteID, path, transitGwName, transitSiteID, path)
}

func testAccCheckEdgeSpokeTransitAttachmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge as a spoke transit attachment not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge as a spoke transit attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		attachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}
		attachment, err := client.GetEdgeSpokeTransitAttachment(context.Background(), attachment)
		if err != nil {
			return err
		}
		if attachment.SpokeGwName+"~"+attachment.TransitGwName != rs.Primary.ID {
			return fmt.Errorf("edge as a spoke transit attachment not found")
		}

		return nil
	}
}

func testAccCheckEdgeSpokeTransitAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_transit_attachment" {
			continue
		}

		attachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}

		_, err := client.GetEdgeSpokeTransitAttachment(context.Background(), attachment)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge as a spoke transit attachment still exists %s", err.Error())
		}
	}

	return nil
}
