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

func TestAccAviatrixDeviceVirtualWanAttachment_basic(t *testing.T) {
	if os.Getenv("SKIP_DEVICE_VIRTUAL_WAN_ATTACHMENT") == "yes" {
		t.Skip("Skipping virtual wan and device attachment test as SKIP_DEVICE_VIRTUAL_WAN_ATTACHMENT is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_device_virtual_wan_attachment.test_device_virtual_wan_attachment"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
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
resource "aviatrix_device_virtual_wan_attachment" "test_device_virtual_wan_attachment" {
	connection_name = "conn-%s"
	device_name     = "%s"
	account_name    = "%s"
	resource_group  = "%s"
	hub_name        = "%s"
	device_bgp_asn  = %s
}
`, rName, os.Getenv("DEVICE_NAME"), os.Getenv("AZURE_ACCOUNT_NAME"),
		os.Getenv("AZURE_RESOURCE_GROUP"), os.Getenv("AZURE_HUB_NAME"), os.Getenv("DEVICE_ASN"))
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
	if os.Getenv("DEVICE_NAME") == "" {
		t.Fatal("environment variable DEVICE_NAME must be set for aviatrix_device_virtual_wan_attachment acceptance test")
	}
	if os.Getenv("AZURE_ACCOUNT_NAME") == "" {
		t.Fatal("environment variable AZURE_ACCOUNT_NAME must be set for aviatrix_device_virtual_wan_attachment acceptance test")
	}
	if os.Getenv("AZURE_RESOURCE_GROUP") == "" {
		t.Fatal("environment variable AZURE_RESOURCE_GROUP must be set for aviatrix_device_virtual_wan_attachment acceptance test")
	}
	if os.Getenv("AZURE_HUB_NAME") == "" {
		t.Fatal("environment variable AZURE_HUB_NAME must be set for aviatrix_device_virtual_wan_attachment acceptance test")
	}
	if os.Getenv("DEVICE_ASN") == "" {
		t.Fatal("environment variable DEVICE_ASN must be set for aviatrix_device_virtual_wan_attachment acceptance test")
	}
}
