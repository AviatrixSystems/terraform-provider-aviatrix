package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixBranchRouterTag_basic(t *testing.T) {
	if os.Getenv("SKIP_BRANCH_ROUTER_TAG") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_BRANCH_ROUTER_TAG is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_branch_router_tag.test_branch_router_tag"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAviatrixBranchRouterTagtPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBranchRouterTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchRouterTagBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBranchRouterTagExists(resourceName),
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

func testAccBranchRouterTagBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_branch_router_tag" "test_branch_router_tag" {
	name                = "branchroutertag-%s"
	config              = <<EOD
hostname myrouter
EOD
	branch_router_names = ["%s"]
}
`, rName, os.Getenv("BRANCH_ROUTER_NAME"))
}

func testAccCheckBranchRouterTagExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("branch_router_tag Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no branch_router_tag ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		brt := &goaviatrix.BranchRouterTag{
			Name: rs.Primary.Attributes["name"],
		}

		_, err := client.GetBranchRouterTag(brt)
		if err != nil {
			return err
		}
		if brt.Name != rs.Primary.ID {
			return fmt.Errorf("branch_router_tag not found")
		}

		return nil
	}
}

func testAccCheckBranchRouterTagDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_branch_router_tag" {
			continue
		}
		brt := &goaviatrix.BranchRouterTag{
			Name: rs.Primary.Attributes["name"],
		}
		_, err := client.GetBranchRouterTag(brt)
		if err == nil {
			return fmt.Errorf("branch_router_tag still exists")
		}
	}

	return nil
}

func testAccAviatrixBranchRouterTagtPreCheck(t *testing.T) {
	if os.Getenv("BRANCH_ROUTER_NAME") == "" {
		t.Fatal("BRANCH_ROUTER_NAME must be set for aviatrix_branch_router_tag acceptance test.")
	}
}
