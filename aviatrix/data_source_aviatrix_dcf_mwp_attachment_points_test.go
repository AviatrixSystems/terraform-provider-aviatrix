package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixDcfMwpAttachmentPoints_basic(t *testing.T) {
	resourceName := "data.aviatrix_dcf_mwp_attachment_points.test"

	skipAcc := os.Getenv("SKIP_DATA_DCF_MWP_ATTACHMENT_POINTS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source DCF MWP Attachment Points test as SKIP_DATA_DCF_MWP_ATTACHMENT_POINTS is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, ". Set SKIP_DATA_DCF_MWP_ATTACHMENT_POINTS to yes to skip Data Source DCF MWP Attachment Points tests")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixDcfMwpAttachmentPointsConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixDcfMwpAttachmentPoints(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment_point_id"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixDcfMwpAttachmentPointsConfigBasic() string {
	return `
data "aviatrix_dcf_mwp_attachment_point" "test" {
	name = "test-attachment-point"
}
	`
}

func testAccDataSourceAviatrixDcfMwpAttachmentPoints(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
