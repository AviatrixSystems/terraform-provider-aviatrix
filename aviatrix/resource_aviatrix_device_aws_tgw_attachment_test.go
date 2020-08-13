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

func TestAccAviatrixDeviceAwsTgwAttachment_basic(t *testing.T) {
	if os.Getenv("SKIP_DEVICE_AWS_TGW_ATTACHMENT") == "yes" {
		t.Skip("Skipping device and aws tgw attachment test as SKIP_DEVICE_AWS_TGW_ATTACHMENT is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_device_aws_tgw_attachment.test_device_aws_tgw_attachment"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAviatrixDeviceAwsTgwAttachmentPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeviceAwsTgwAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceAwsTgwAttachmentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceAwsTgwAttachmentExists(resourceName),
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

func testAccDeviceAwsTgwAttachmentBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_device_aws_tgw_attachment" "test_device_aws_tgw_attachment" {
	connection_name      = "conn-%s"
	device_name          = "%s"
	aws_tgw_name         = "%s"
	device_bgp_asn       = 65001
	security_domain_name = "Default_Domain"
}
`, rName, os.Getenv("DEVICE_NAME"), os.Getenv("AWS_TGW_NAME"))
}

func testAccCheckDeviceAwsTgwAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("device_aws_tgw_attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no device_aws_tgw_attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		attachment := &goaviatrix.DeviceAwsTgwAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
			DeviceName:     rs.Primary.Attributes["device_name"],
			AwsTgwName:     rs.Primary.Attributes["aws_tgw_name"],
		}

		_, err := client.GetDeviceAwsTgwAttachment(attachment)
		if err != nil {
			return err
		}
		if attachment.ID() != rs.Primary.ID {
			return fmt.Errorf("device_aws_tgw_attachment not found")
		}

		return nil
	}
}

func testAccCheckDeviceAwsTgwAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_device_aws_tgw_attachment" {
			continue
		}
		attachment := &goaviatrix.DeviceAwsTgwAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
			DeviceName:     rs.Primary.Attributes["device_name"],
			AwsTgwName:     rs.Primary.Attributes["aws_tgw_name"],
		}
		_, err := client.GetDeviceAwsTgwAttachment(attachment)
		if err == nil {
			return fmt.Errorf("device_aws_tgw_attachment still exists")
		}
	}

	return nil
}

func testAccAviatrixDeviceAwsTgwAttachmentPreCheck(t *testing.T) {
	if os.Getenv("DEVICE_NAME") == "" {
		t.Fatal("DEVICE_NAME must be set for aviatrix_device_aws_tgw_attachment acceptance test.")
	}
	if os.Getenv("AWS_TGW_NAME") == "" {
		t.Fatal("AWS_TGW_NAME must be set for aviatrix_device_aws_tgw_attachment acceptance test.")
	}
}
