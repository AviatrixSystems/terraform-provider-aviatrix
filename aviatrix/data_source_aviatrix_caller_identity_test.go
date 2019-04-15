package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceAviatrixCallerIdentity_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_caller_identity.foo"

	skipAcc := os.Getenv("SKIP_DATA_CALLER_IDENTITY")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Caller Identity test as SKIP_DATA_CALLER_IDENTITY is set")
	}

	preAccountCheck(t, ". Set SKIP_DATA_CALLER_IDENTITY to yes to skip Data Source Caller Identity tests")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixCallerIdentityConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixCallerIdentity(resourceName),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixCallerIdentityConfigBasic(rName string) string {
	return fmt.Sprintf(`

data "aviatrix_caller_identity" "foo" {
}

`)
}

func testAccDataSourceAviatrixCallerIdentity(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)
		client.CID = rs.Primary.Attributes["cid"]

		version, _, err := client.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("valid CID was not returned. Get version API gave the following Error: %v", err)
		}

		if !strings.Contains(version, "UserConnect") {
			return fmt.Errorf("valid CID was not returned. Get version API gave the wrong version")
		}

		return nil
	}
}
