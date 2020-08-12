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

func TestAccAviatrixDeviceRegistration_basic(t *testing.T) {
	if os.Getenv("SKIP_DEVICE_REGISTRATION") == "yes" {
		t.Skip("Skipping Device registration test as SKIP_DEVICE_REGISTRATION is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_device_registration.test_device"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			deviceRegistrationPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeviceRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRegistrationBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceRegistrationExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "key_file"},
			},
		},
	})
}

func testAccDeviceRegistrationBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_device_registration" "test_device" {
	name        = "device-registration-%s"
	public_ip   = "%s"
	username    = "ec2-user"
	key_file    = "%s"
	host_os     = "ios"
	ssh_port    = 22
	address_1   = "2901 Tasman Dr"
	address_2   = "Suite #104"
	city        = "Santa Clara"
	state       = "CA"
	zip_code    = "12323"
	description = "Test device."
}
`, rName, os.Getenv("DEVICE_PUBLIC_IP"), os.Getenv("DEVICE_KEY_FILE_PATH"))
}

func testAccCheckDeviceRegistrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("device_registration Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no device_registration ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		device := &goaviatrix.Device{
			Name: rs.Primary.Attributes["name"],
		}

		_, err := client.GetDevice(device)
		if err != nil {
			return err
		}
		if device.Name != rs.Primary.ID {
			return fmt.Errorf("device_registration not found")
		}

		return nil
	}
}

func testAccCheckDeviceRegistrationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_device_registration" {
			continue
		}
		device := &goaviatrix.Device{
			Name: rs.Primary.Attributes["name"],
		}
		_, err := client.GetDevice(device)
		if err == nil {
			return fmt.Errorf("device_registration still exists")
		}
	}

	return nil
}

func deviceRegistrationPreCheck(t *testing.T) {
	if os.Getenv("DEVICE_PUBLIC_IP") == "" {
		t.Fatal("environment variable DEVICE_PUBLIC_IP must be set for device_registration acceptance test")
	}

	if os.Getenv("DEVICE_KEY_FILE_PATH") == "" {
		t.Fatal("environment variable DEVICE_KEY_FILE_PATH must be set for " +
			"device_registration acceptance test")
	}
}
