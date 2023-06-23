package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAviatrixSmartGroups_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_smart_groups.test"

	skipAcc := os.Getenv("SKIP_DATA_SMART_GROUPS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Smart Groups tests as SKIP_DATA_SMART_GROUPS is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixSmartGroupsConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "smart_groups.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "smart_groups.0.uuid"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixSmartGroupsConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_smart_group" "test" {
	name = "aaa-smart-group"
	selector {
		match_expressions {
			cidr = "11.0.0.0/16"
		}
	}
}
data "aviatrix_smart_groups" "test"{
	depends_on = [
        aviatrix_smart_group.test
  ]
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}
