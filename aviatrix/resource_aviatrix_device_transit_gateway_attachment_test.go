package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDeviceTransitGatewayAttachment_basic(t *testing.T) {
	if os.Getenv("SKIP_DEVICE_TRANSIT_GATEWAY_ATTACHMENT") == "yes" {
		t.Skip("Skipping transit gateway and device attachment test as SKIP_DEVICE_TRANSIT_GATEWAY_ATTACHMENT is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_device_transit_gateway_attachment.test_device_transit_gateway_attachment"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			deviceRegistrationPreCheck(t)
			testAccAviatrixDeviceTransitGatewayAttachmentPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeviceTransitGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTransitGatewayAttachmentNoOptions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceTransitGatewayAttachmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"pre_shared_key",
					"local_tunnel_ip",
					"remote_tunnel_ip",
				},
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			deviceRegistrationPreCheck(t)
			testAccAviatrixDeviceTransitGatewayAttachmentPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeviceTransitGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTransitGatewayAttachmentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceTransitGatewayAttachmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"pre_shared_key",
				},
			},
		},
	})
}

func testAccDeviceTransitGatewayAttachmentBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_device_registration" "test_device_registration" {
	name        = "device-registration-%[1]s"
	public_ip   = "%[2]s"
	username    = "ec2-user"
	key_file    = "%[3]s"
	host_os     = "ios"
	ssh_port    = 22
	address_1   = "2901 Tasman Dr"
	address_2   = "Suite #104"
	city        = "Santa Clara"
	state       = "CA"
	zip_code    = "12323"
	description = "Test device."
}
resource "aviatrix_device_interface_config" "test_device_interface_config" {
	device_name                     = aviatrix_device_registration.test_device_registration.name
	wan_primary_interface           = "GigabitEthernet1"
	wan_primary_interface_public_ip = "%[2]s"
}
resource "aviatrix_device_transit_gateway_attachment" "test_device_transit_gateway_attachment" {
	device_name               = aviatrix_device_interface_config.test_device_interface_config.device_name 
	transit_gateway_name      = "%[4]s"
	connection_name           = "connection-%[1]s"
	transit_gateway_bgp_asn   = 65000
	device_bgp_asn            = 65001
	phase1_authentication     = "SHA-256"
	phase1_dh_groups          = 14
	phase1_encryption         = "AES-256-CBC"
	phase2_authentication     = "HMAC-SHA-256"
	phase2_dh_groups          = 14
	phase2_encryption         = "AES-256-CBC"
	enable_global_accelerator = true
	pre_shared_key            = "key"
	local_tunnel_ip           = "10.0.0.1/30"
	remote_tunnel_ip          = "10.0.0.2/30"
}
`, rName, os.Getenv("DEVICE_PUBLIC_IP"), os.Getenv("DEVICE_KEY_FILE_PATH"),
		os.Getenv("TRANSIT_GATEWAY_NAME"))
}

func testAccDeviceTransitGatewayAttachmentNoOptions(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_device_registration" "test_device_registration" {
	name        = "device-registration-noopts-%[1]s"
	public_ip   = "%[2]s"
	username    = "ec2-user"
	key_file    = "%[3]s"
	host_os     = "ios"
	ssh_port    = 22
	address_1   = "2901 Tasman Dr"
	address_2   = "Suite #104"
	city        = "Santa Clara"
	state       = "CA"
	zip_code    = "12323"
	description = "Test device."
}
resource "aviatrix_device_interface_config" "test_device_interface_config" {
	device_name                     = aviatrix_device_registration.test_device_registration.name
	wan_primary_interface           = "GigabitEthernet1"
	wan_primary_interface_public_ip = "%[2]s"
}
resource "aviatrix_device_transit_gateway_attachment" "test_device_transit_gateway_attachment" {
	device_name               = aviatrix_device_interface_config.test_device_interface_config.device_name 
	transit_gateway_name      = "%[4]s"
	connection_name           = "connection-noopts-%[1]s"
	transit_gateway_bgp_asn   = 65000
	device_bgp_asn            = 65001

}
`, rName, os.Getenv("DEVICE_PUBLIC_IP"), os.Getenv("DEVICE_KEY_FILE_PATH"),
		os.Getenv("TRANSIT_GATEWAY_NAME"))
}

func testAccCheckDeviceTransitGatewayAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("device_transit_gateway_attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no device_transit_gateway_attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		attachment := &goaviatrix.DeviceTransitGatewayAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetDeviceTransitGatewayAttachment(attachment)
		if err != nil {
			return err
		}
		if attachment.ConnectionName != rs.Primary.ID {
			return fmt.Errorf("device_transit_gateway_attachment not found")
		}

		return nil
	}
}

func testAccCheckDeviceTransitGatewayAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_device_transit_gateway_attachment" {
			continue
		}
		attachment := &goaviatrix.DeviceTransitGatewayAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}
		_, err := client.GetDeviceTransitGatewayAttachment(attachment)
		if err == nil {
			return fmt.Errorf("device_transit_gateway_attachment still exists")
		}
	}

	return nil
}

func testAccAviatrixDeviceTransitGatewayAttachmentPreCheck(t *testing.T) {
	if os.Getenv("TRANSIT_GATEWAY_NAME") == "" {
		t.Fatal("TRANSIT_GATEWAY_NAME must be set for aviatrix_device_transit_gateway_attachment acceptance test.")
	}
}
