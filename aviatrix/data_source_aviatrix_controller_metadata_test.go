package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixControllerMetadata_basic(t *testing.T) {
	resourceName := "data.aviatrix_controller_metadata.foo"

	skipAcc := os.Getenv("SKIP_DATA_CONTROLLER_METADATA")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Controller Metadata test as SKIP_DATA_CONTROLLER_METADATA is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, ". Set SKIP_DATA_CONTROLLER_METADATA to yes to skip Data Source Controller Metadata tests")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixControllerMetadataConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixControllerMetadata(resourceName),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixControllerMetadataConfigBasic() string {
	return `
data "aviatrix_controller_metadata" "foo" {
}
	`
}

func testAccDataSourceAviatrixControllerMetadata(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
