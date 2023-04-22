package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixEdgeNEODeviceOnboarding_basic(t *testing.T) {
	if os.Getenv("SKIP_EDGE_NEO_DEVICE_ONBOARDING") == "yes" {
		t.Skip("Skipping Edge NEO Device Onboarding test as SKIP_EDGE_NEO_DEVICE_ONBOARDING is set")
	}

	resourceName := "aviatrix_edge_neo_device_onboarding.test"
	accountName := "acc-" + acctest.RandString(5)
	deviceName := "device-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeNEODeviceOnboardingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeNEODeviceOnboardingBasic(accountName, deviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeNEODeviceOnboardingExists(resourceName, accountName, deviceName),
					resource.TestCheckResourceAttr(resourceName, "device_name", deviceName),
					resource.TestCheckResourceAttr(resourceName, "serial_number", os.Getenv("EDGE_NEO_DEVICE_SERIAL_NUMBER")),
					resource.TestCheckResourceAttr(resourceName, "hardware_model", os.Getenv("EDGE_NEO_DEVICE_HARDWARE_MODEL")),
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

func testAccEdgeNEODeviceOnboardingBasic(accountName, deviceName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name = "%s"
	cloud_type   = 262144
}
resource "aviatrix_edge_neo_device_onboarding" "test" {
	account_name   = aviatrix_account.test.account_name
	device_name    = "%s"
	serial_number  = "%s"
	hardware_model = "%s"
}
 `, accountName, deviceName, os.Getenv("EDGE_NEO_DEVICE_SERIAL_NUMBER"), os.Getenv("EDGE_NEO_DEVICE_HARDWARE_MODEL"))
}

func testAccCheckEdgeNEODeviceOnboardingExists(resourceName, accountName, deviceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge neo device onboarding not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge neo device onboarding id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetEdgeNEODevice(context.Background(), accountName, deviceName)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("could not find edge neo device")
			}
			return err
		}

		return nil
	}
}

func testAccCheckEdgeNEODeviceOnboardingDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_neo_device_onboarding" {
			continue
		}

		_, err := client.GetEdgeNEODevice(context.Background(), rs.Primary.Attributes["account_name"], rs.Primary.Attributes["device_name"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("edge neo device still exists")
		}
	}

	return nil
}
