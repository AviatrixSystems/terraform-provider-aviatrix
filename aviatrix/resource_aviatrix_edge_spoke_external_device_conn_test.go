package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixEdgeSpokeExternalDeviceConn_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_edge_spoke_external_device_conn.test"

	skipAcc := os.Getenv("SKIP_EDGE_SPOKE_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping Edge as a Spoke external device connection tests as 'SKIP_EDGE_SPOKE_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preEdgeSpokeExternalDeviceConnCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeSpokeExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeSpokeExternalDeviceConnConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeSpokeExternalDeviceConnExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_id", os.Getenv("EDGE_SPOKE_SITE_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", os.Getenv("EDGE_SPOKE_NAME")),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "65001"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "65002"),
					resource.TestCheckResourceAttr(resourceName, "local_lan_ip", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "remote_lan_ip", "5.6.7.8"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_lan_activemesh", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_jumbo_frame", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	skipBgpBfd := os.Getenv("SKIP_BGP_BFD_EDGE_SPOKE_EXTERNAL_DEVICE_CONN")
	if skipBgpBfd != "yes" {
		resourceNameBfd := "aviatrix_edge_spoke_external_device_conn.test-bfd"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preEdgeSpokeExternalDeviceConnCheck(t)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckEdgeSpokeExternalDeviceConnDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccEdgeSpokeExternalDeviceConnConfigBgpBfd(),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckEdgeSpokeExternalDeviceConnExists(resourceNameBfd),
						resource.TestCheckResourceAttr(resourceNameBfd, "site_id", os.Getenv("EDGE_SPOKE_SITE_ID")),
						resource.TestCheckResourceAttr(resourceNameBfd, "connection_name", "cloudn-bgp-lan"),
						resource.TestCheckResourceAttr(resourceNameBfd, "gw_name", os.Getenv("EDGE_SPOKE_NAME")),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_local_as_num", "65182"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_remote_as_num", "65220"),
						resource.TestCheckResourceAttr(resourceNameBfd, "local_lan_ip", "10.220.86.182"),
						resource.TestCheckResourceAttr(resourceNameBfd, "remote_lan_ip", "10.220.86.100"),
						resource.TestCheckResourceAttr(resourceNameBfd, "enable_bfd", "true"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgo_bfd.0.transmit_interval", "400"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgo_bfd.0.receive_interval", "400"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgo_bfd.0.multiplier", "5"),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Skip("Skipping BGP BFD Edge as a Spoke external device connection tests as 'SKIP_BGP_BFD_EDGE_SPOKE_EXTERNAL_DEVICE_CONN' is set")
	}
}

func preEdgeSpokeExternalDeviceConnCheck(t *testing.T) {
	if os.Getenv("EDGE_SPOKE_NAME") == "" {
		t.Fatal("Environment variable EDGE_SPOKE_NAME is not set")
	}

	if os.Getenv("EDGE_SPOKE_SITE_ID") == "" {
		t.Fatal("Environment variable EDGE_SPOKE_SITE_ID is not set")
	}
}

func testAccEdgeSpokeExternalDeviceConnConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_edge_spoke_external_device_conn" "test" {
	site_id           = "%s"
	connection_name   = "%s"
	gw_name           = "%s"
	bgp_local_as_num  = "65001"
	bgp_remote_as_num = "65002"
	local_lan_ip      = "1.2.3.4"
	remote_lan_ip     = "5.6.7.8"
	connection_type   = "bgp"
	enable_bgp_lan_activemesh = true
	enable_jumbo_frame = true
}
	`, os.Getenv("EDGE_SPOKE_SITE_ID"), rName, os.Getenv("EDGE_SPOKE_NAME"))
}

func testAccEdgeSpokeExternalDeviceConnConfigBgpBfd() string {
	return fmt.Sprintf(`
resource "aviatrix_edge_spoke_external_device_conn" "test-bfd" {
	site_id           = "%s"
	connection_name   = "cloudn-bgp-lan"
	gw_name           = "%s"
	bgp_local_as_num  = "65182"
	bgp_remote_as_num = "65220"
	local_lan_ip      = "10.220.86.182"
	remote_lan_ip     = "10.220.86.100"
	enable_bfd        = true
	bgo_bfd {
		transmit_interval = 400
		receive_interval = 400
		multiplier = 5
	}
}
	`, os.Getenv("EDGE_SPOKE_SITE_ID"), os.Getenv("EDGE_SPOKE_NAME"))
}

func testAccCheckEdgeSpokeExternalDeviceConnExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge as a spoke external device conn not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge as a spoke external device conn ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["site_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		conn, err := client.GetExternalDeviceConnDetail(externalDeviceConn)
		if err != nil {
			return err
		}

		if conn.ConnectionName+"~"+conn.VpcID != rs.Primary.ID {
			return fmt.Errorf("edge as a spoke external device conn not found")
		}

		return nil
	}
}

func testAccCheckEdgeSpokeExternalDeviceConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_spoke_external_device_conn" {
			continue
		}

		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["site_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetExternalDeviceConnDetail(externalDeviceConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge as a spoke external device conn still exists %s", err.Error())
		}
	}

	return nil
}
