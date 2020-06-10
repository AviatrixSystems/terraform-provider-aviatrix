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
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"commit"},
			},
		},
	})
}

func testAccBranchRouterTagBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_branch_router" "test_branch_router_a" {
	name        = "branchrouter-A%[1]s"
	public_ip   = "18.144.102.14"
	username    = "ec2-user"
	password    = "testing"
	zip_code    = "19232"
}

resource "aviatrix_branch_router" "test_branch_router_b" {
	name        = "branchrouter-B%[1]s"
	public_ip   = "18.144.102.14"
	username    = "ec2-user"
	password    = "testing"
	zip_code    = "19232"
}

resource "aviatrix_branch_router_tag" "test_branch_router_tag" {
	name     = "branchroutertag-%[1]s"
	config   = <<EOD
hostname myrouter
EOD
	branches = [aviatrix_branch_router.test_branch_router_a.name, aviatrix_branch_router.test_branch_router_b.name]
	commit   = false
}
`, rName)
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
