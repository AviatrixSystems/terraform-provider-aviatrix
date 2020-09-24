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

func TestAccAviatrixDeviceTag_basic(t *testing.T) {
	if os.Getenv("SKIP_DEVICE_TAG") == "yes" {
		t.Skip("Skipping Device tag test as SKIP_DEVICE_TAG is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_device_tag.test_device_tag"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			deviceRegistrationPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeviceTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTagBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceTagExists(resourceName),
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

func testAccDeviceTagBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_device_registration" "test_device_registration" {
	name        = "device-registration-%s"
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

resource "aviatrix_device_tag" "test_device_tag" {
	name                = "device-tag-%[1]s"
	config              = <<EOD
hostname myrouter
EOD
	device_names        = [aviatrix_device_registration.test_device_registration.name]
	depends_on          = [aviatrix_device_interface_config.test_device_interface_config]
}
`, rName, os.Getenv("DEVICE_PUBLIC_IP"), os.Getenv("DEVICE_KEY_FILE_PATH"))
}

func testAccCheckDeviceTagExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("device_tag Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no device_tag ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		deviceTag := &goaviatrix.DeviceTag{
			Name: rs.Primary.Attributes["name"],
		}

		_, err := client.GetDeviceTag(deviceTag)
		if err != nil {
			return err
		}
		if deviceTag.Name != rs.Primary.ID {
			return fmt.Errorf("device_tag not found")
		}

		return nil
	}
}

func testAccCheckDeviceTagDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_device_tag" {
			continue
		}
		deviceTag := &goaviatrix.DeviceTag{
			Name: rs.Primary.Attributes["name"],
		}
		_, err := client.GetDeviceTag(deviceTag)
		if err == nil {
			return fmt.Errorf("device_tag still exists")
		}
	}

	return nil
}
