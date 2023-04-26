package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixGlobalVpcTaggingSettings_basic(t *testing.T) {
	if os.Getenv("SKIP_GLOBAL_VPC_TAGGING_SETTINGS") == "yes" {
		t.Skip("Skipping global vpc tagging settings test as SKIP_GLOBAL_VPC_TAGGING_SETTINGS is set")
	}

	resourceName := "aviatrix_global_vpc_tagging_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGlobalVpcTaggingSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalVpcTaggingSettingsBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalVpcTaggingSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service_state", "automatic"),
					resource.TestCheckResourceAttr(resourceName, "enable_alert", "false"),
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

func testAccGlobalVpcTaggingSettingsBasic() string {
	return `
resource "aviatrix_global_vpc_tagging_settings" "test" {
	service_state = "automatic"
	enable_alert  = false
}
`
}

func testAccCheckGlobalVpcTaggingSettingsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("global vpc tagging settings not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		globalVpcTaggingSettings, _ := client.GetGlobalVpcTaggingSettings(context.Background())
		if globalVpcTaggingSettings.ServiceState != "automatic" {
			return fmt.Errorf("global vpc tagging settings not found")
		}

		return nil
	}
}

func testAccCheckGlobalVpcTaggingSettingsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_global_vpc_tagging_settings" {
			continue
		}

		globalVpcTaggingSettings, _ := client.GetGlobalVpcTaggingSettings(context.Background())
		if globalVpcTaggingSettings.ServiceState == "automatic" {
			return fmt.Errorf("global vpc tagging settings still exists")
		}
	}

	return nil
}
