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

func TestAccAviatrixSmartGroup_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_SMART_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping Smart Group test as SKIP_SMART_GROUP is set")
	}
	resourceName := "aviatrix_smart_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccSmartGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-smart-group"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.cidr", "11.0.0.0/16"),

					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.type", "vm"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.account_name", "devops"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.region", "us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.tags.k3", "v3"),
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

func testAccSmartGroupBasic() string {
	return `
resource "aviatrix_smart_group" "test" {
	name = "test-smart-group"

	selector {
		match_expressions {
			cidr = "11.0.0.0/16"
		}

		match_expressions {
			type         = "vm"
			account_name = "devops"
			region       = "us-west-2"
			tags         = {
				k3 = "v3"
			}
		}
	}
}
`
}

func testAccCheckSmartGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Smart Group resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Smart Group ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		smartGroup, err := client.GetSmartGroup(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get Smart Group status: %v", err)
		}

		if smartGroup.UUID != rs.Primary.ID {
			return fmt.Errorf("smart Group ID not found")
		}

		return nil
	}
}

func testAccSmartGroupDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_smart_group" {
			continue
		}

		_, err := client.GetSmartGroup(context.Background(), rs.Primary.ID)
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("smart group configured when it should be destroyed")
		}
	}

	return nil
}
