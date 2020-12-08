package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAviatrixAwsGuardDuty_basic(t *testing.T) {
	if os.Getenv("SKIP_AWS_GUARD_DUTY") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_AWS_GUARD_DUTY is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_aws_guard_duty.test_aws_guard_duty"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsGuardDutyBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanning_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "enabled_accounts.#", "1"),
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

func testAccAwsGuardDutyBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc" {
	account_name       = "tf-testing-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_guard_duty" "test_aws_guard_duty" {
	scanning_interval = 5
	enabled_accounts {
		account_name = aviatrix_account.test_acc.account_name
		region = "us-west-1"
	}
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}
