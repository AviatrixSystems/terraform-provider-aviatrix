package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixDcfWebgroups_basic(t *testing.T) {
	resourceName := "data.aviatrix_web_group.test"

	skipAcc := os.Getenv("SKIP_DATA_DCF_WEBGROUPS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source DCF Webgroups test as SKIP_DATA_DCF_WEBGROUPS is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, ". Set SKIP_DATA_DCF_WEBGROUPS to yes to skip Data Source DCF Webgroups tests")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixDcfWebgroupsConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixDcfWebgroups(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixDcfWebgroupsConfigBasic() string {
	return `
resource aviatrix_web_group "test_webgroup" {
	name = "test-webgroup"
	selector {
		match_expressions {
			snifilter = "example.com"
		}
	}
}

data "aviatrix_web_group" "test" {
	depends_on = [aviatrix_web_group.test_webgroup]
	name = "test-webgroup"
}
	`
}

func testAccDataSourceAviatrixDcfWebgroups(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
