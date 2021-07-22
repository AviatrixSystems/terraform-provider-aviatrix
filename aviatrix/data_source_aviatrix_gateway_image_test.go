package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixGatewayImage_basic(t *testing.T) {
	resourceName := "data.aviatrix_gateway_image.foo"

	skipAcc := os.Getenv("SKIP_DATA_GATEWAY_IMAGE")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Gateway Image test as SKIP_DATA_GATEWAY_IMAGE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixGatewayImageConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixGatewayImage(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_version", "hvm-cloudx-aws-022021"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixGatewayImageConfigBasic() string {
	return `
data "aviatrix_gateway_image" "foo" {
	cloud_type       = 1
	software_version = "6.5" 
}
	`
}

func testAccDataSourceAviatrixGatewayImage(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
