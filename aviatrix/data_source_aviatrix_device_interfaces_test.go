package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixDeviceInterfaces_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_device_interfaces.foo"

	skipAcc := os.Getenv("SKIP_DATA_DEVICE_INTERFACES")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Device Interfaces tests as SKIP_DATA_DEVICE_INTERFACES is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeviceInterfacesConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixDeviceInterfaces(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "wan_interfaces.0.wan_primary_interface"),
					resource.TestCheckResourceAttrSet(resourceName, "wan_interfaces.0.wan_primary_interface_public_ip"),
				),
			},
		},
	})
}

func testAccDataSourceDeviceInterfacesConfigBasic(rName string) string {
	return fmt.Sprintf(`
data "aviatrix_device_interfaces" "foo" {
	device_name = "%s"
}
	`, os.Getenv("CLOUDN_DEVICE_NAME"))
}

func testAccDataSourceAviatrixDeviceInterfaces(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
