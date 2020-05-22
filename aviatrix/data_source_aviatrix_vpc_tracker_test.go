package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAviatrixVpcTracker_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_vpc_tracker.test"

	skipAcc := os.Getenv("SKIP_DATA_VPC_TRACKER")
	if skipAcc == "yes" {
		t.Skip("Skipping data source vpc_tracker tests as 'SKIP_DATA_VPC_TRACKER' is set")
	}

	msg := ". Set 'SKIP_DATA_VPC_TRACKER' to 'yes' to skip data source vpc_tracker tests"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msg)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixVpcTrackerConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixVpcTracker(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_list.0.account_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_list.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_list.0.cloud_type"),
					resource.TestCheckResourceAttrSet(resourceName, "cloud_type"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixVpcTrackerConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
data "aviatrix_vpc_tracker" "test" {
	cloud_type = 1
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccDataSourceAviatrixVpcTracker(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
