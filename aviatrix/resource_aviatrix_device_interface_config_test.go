package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixDeviceInterfaceConfig_basic(t *testing.T) {
	if os.Getenv("SKIP_DEVICE_INTERFACE_CONFIG") == "yes" {
		t.Skip("Skipping device interface config test as SKIP_DEVICE_INTERFACE_CONFIG is set")
	}

	resourceName := "aviatrix_device_interface_config.test_device_interface_config"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceInterfaceConfigBasic(),
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

func testAccDeviceInterfaceConfigBasic() string {
	return fmt.Sprintf(`
data "aviatrix_device_interfaces" "test" {
	device_name = "%s"
}

resource "aviatrix_device_interface_config" "test_device_interface_config" {
	device_name                     = data.aviatrix_device_interfaces.test.device_name
	wan_primary_interface           = "eth0"
	wan_primary_interface_public_ip = data.aviatrix_device_interfaces.test.wan_interfaces[0].wan_primary_interface_public_ip
}
`, os.Getenv("CLOUDN_DEVICE_NAME"))
}

func testAccCheckDeviceInterfaceConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ONE device_interface_config Not found: %s", n)
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
			device.PrimaryInterface != rs.Primary.Attributes["wan_primary_interface"] {
			return fmt.Errorf("device_interface_config not found")
		}

		return nil
	}
}
