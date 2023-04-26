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

func TestAccAviatrixWebGroup_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_WEB_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Web Group test as SKIP_WEB_GROUP is set")
	}
	resourceName := "aviatrix_web_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccWebGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-web-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.urlfilter", "https://aviatrix.com/test"),
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

func testAccWebGroupBasic() string {
	return `
resource "aviatrix_web_group" "test" {
	name = "test-web-group"

	selector {
		match_expressions {
			urlfilter = "https://aviatrix.com/test"
		}
	}
}
`
}

func testAccCheckWebGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Web Group resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Web Group ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		webGroup, err := client.GetWebGroup(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get Web Group status: %v", err)
		}

		if webGroup.UUID != rs.Primary.ID {
			return fmt.Errorf("web Group ID not found")
		}

		return nil
	}
}

func testAccWebGroupDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_web_group" {
			continue
		}

		_, err := client.GetWebGroup(context.Background(), rs.Primary.ID)
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("web group configured when it should be destroyed")
		}
	}

	return nil
}
