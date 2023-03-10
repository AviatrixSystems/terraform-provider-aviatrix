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

func TestAccAviatrixControllerAccessAllowList_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_ACCESS_ALLOW_LIST")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Access Allow List test as SKIP_CONTROLLER_ACCESS_ALLOW_LIST is set")
	}

	resourceName := "aviatrix_controller_access_allow_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControllerAccessAllowListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerAccessAllowListBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerAccessAllowListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "allow_list.0.ip_address", "0.0.0.0"),
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

func testAccControllerAccessAllowListBasic() string {
	return `
resource "aviatrix_controller_access_allow_list" "test" {
	allow_list {
		ip_address = "0.0.0.0"
	}
}
`
}

func testAccCheckControllerAccessAllowListExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller access allow list Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no controller access allow list ID is set")
		}

		if rs.Primary.ID != "allow_list" {
			return fmt.Errorf("controller access allow list ID not found")
		}

		return nil
	}
}

func testAccCheckControllerAccessAllowListDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_access_allow_list" {
			continue
		}

		_, err := client.GetControllerAccessAllowList(context.Background())
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("controller access allow list still exists")
		}
	}

	return nil
}
