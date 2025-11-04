package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixDcfLogProfile_basic(t *testing.T) {
	resourceName := "data.aviatrix_dcf_log_profile.test"

	skipAcc := os.Getenv("SKIP_DATA_DCF_LOG_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source DCF Log Profile test as SKIP_DATA_DCF_LOG_PROFILE is set")
	}

	testLogProfileName := "start/end"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, ". Set SKIP_DATA_DCF_LOG_PROFILE to yes to skip Data Source DCF Log Profile tests")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixDcfLogProfileConfigBasic(testLogProfileName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixDcfLogProfile(resourceName),
					resource.TestCheckResourceAttr(resourceName, "profile_name", testLogProfileName),
					resource.TestCheckResourceAttrSet(resourceName, "profile_id"),
					resource.TestCheckResourceAttrSet(resourceName, "session_end"),
					resource.TestCheckResourceAttrSet(resourceName, "session_start"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixDcfLogProfileConfigBasic(profileName string) string {
	return fmt.Sprintf(`
data "aviatrix_dcf_log_profile" "test" {
	profile_name = "%q"
}
	`, profileName)
}

func testAccDataSourceAviatrixDcfLogProfile(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
