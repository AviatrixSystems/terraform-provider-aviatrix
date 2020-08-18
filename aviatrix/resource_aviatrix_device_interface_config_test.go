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

func TestAccAviatrixDeviceInterfaceConfig_basic(t *testing.T) {
	if os.Getenv("SKIP_DEVICE_INTERFACE_CONFIG") == "yes" {
		t.Skip("Skipping device interface config test as SKIP_DEVICE_INTERFACE_CONFIG is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_device_interface_config.test_device_interface_config"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			deviceRegistrationPreCheck(t)
			deviceInterfaceConfigPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceInterfaceConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceInterfaceConfigExists(resourceName),
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

func testAccDeviceInterfaceConfigBasic(rName string) string {
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
	wan_primary_interface           = "%[4]s"
	wan_primary_interface_public_ip = "%[2]s"
}
`, rName, os.Getenv("DEVICE_PUBLIC_IP"), os.Getenv("DEVICE_KEY_FILE_PATH"), os.Getenv("DEVICE_PRIMARY_INTERFACE"))
}

func testAccCheckDeviceInterfaceConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("device_interface_config Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no device_interface_config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		device := &goaviatrix.Device{Name: rs.Primary.Attributes["device_name"]}

		device, err := client.GetDevice(device)
		if err != nil {
			return err
		}

		if device.Name != rs.Primary.ID ||
			device.PrimaryInterface != rs.Primary.Attributes["wan_primary_interface"] ||
			device.PrimaryInterfaceIP != rs.Primary.Attributes["wan_primary_interface_public_ip"] {
			return fmt.Errorf("device_interface_config not found")
		}

		return nil
	}
}

func deviceInterfaceConfigPreCheck(t *testing.T) {
	if os.Getenv("DEVICE_PRIMARY_INTERFACE") == "" {
		t.Fatal("environment variable DEVICE_PRIMARY_INTERFACE must be set for " +
			"device_interface_config acceptance test")
	}
}
