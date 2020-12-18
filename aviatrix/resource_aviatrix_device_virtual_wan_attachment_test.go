package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDeviceVirtualWanAttachment_basic(t *testing.T) {
	if os.Getenv("SKIP_DEVICE_VIRTUAL_WAN_ATTACHMENT") == "yes" {
		t.Skip("Skipping virtual wan and device attachment test as SKIP_DEVICE_VIRTUAL_WAN_ATTACHMENT is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_device_virtual_wan_attachment.test_device_virtual_wan_attachment"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			deviceRegistrationPreCheck(t)
			deviceVirtualWanAttachmentPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeviceVirtualWanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceVirtualWanAttachmentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceVirtualWanAttachmentExists(resourceName),
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

func testAccDeviceVirtualWanAttachmentBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "temp_acc_azure" {
  account_name        = "azureacc-%[1]s"
  cloud_type          = 8
  arm_subscription_id = "%[4]s" 
  arm_directory_id    = "%[5]s"
  arm_application_id  = "%[6]s"
  arm_application_key = "%[7]s"
}
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
resource "aviatrix_device_virtual_wan_attachment" "test_device_virtual_wan_attachment" {
	connection_name = "conn-%[1]s"
	device_name     = aviatrix_device_interface_config.test_device_interface_config.device_name
	account_name    = aviatrix_account.temp_acc_azure.account_name
	resource_group  = "%[8]s"
	hub_name        = "%[9]s"
	device_bgp_asn  = 65001
}
`, rName, os.Getenv("DEVICE_PUBLIC_IP"), os.Getenv("DEVICE_KEY_FILE_PATH"),
		os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("ARM_RESOURCE_GROUP"), os.Getenv("ARM_HUB_NAME"))
}

func testAccCheckDeviceVirtualWanAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("device_virtual_wan_attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no device_virtual_wan_attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		attachment := &goaviatrix.DeviceVirtualWanAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetDeviceVirtualWanAttachment(attachment)
		if err != nil {
			return err
		}
		if attachment.ConnectionName != rs.Primary.ID {
			return fmt.Errorf("device_virtual_wan_attachment not found")
		}

		return nil
	}
}

func testAccCheckDeviceVirtualWanAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_device_virtual_wan_attachment" {
			continue
		}
		attachment := &goaviatrix.DeviceVirtualWanAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}
		_, err := client.GetDeviceVirtualWanAttachment(attachment)
		if err == nil {
			return fmt.Errorf("device_virtual_wan_attachment still exists")
		}
	}

	return nil
}

func deviceVirtualWanAttachmentPreCheck(t *testing.T) {
	errStr := "environment variable %s must be set for aviatrix_device_virtual_wan_attachment acceptance tests"
	requiredEnvVars := []string{
		"ARM_SUBSCRIPTION_ID",
		"ARM_DIRECTORY_ID",
		"ARM_APPLICATION_ID",
		"ARM_APPLICATION_KEY",
		"ARM_RESOURCE_GROUP",
		"ARM_HUB_NAME",
	}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf(errStr, envVar)
		}
	}
}
