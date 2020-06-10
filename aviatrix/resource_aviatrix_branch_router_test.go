package aviatrix

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
	"os"
	"testing"
)

func TestAccAviatrixBranchRouter_basic(t *testing.T) {
	if os.Getenv("SKIP_BRANCH_ROUTER") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_BRANCH_ROUTER is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_branch_router.test_branch_router"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBranchRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchRouterBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBranchRouterExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "key_file"},
			},
		},
	})
}

func testAccBranchRouterBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_branch_router" "test_branch_router" {
	name        = "branchrouter-%s"
	public_ip   = "18.144.102.14"
	username    = "ec2-user"
	password    = "testing"
	host_os     = "ios"
	ssh_port    = 22
	address_1   = "2901 Tasman Dr"
	address_2   = "Suite #104"
	city        = "Santa Clara"
	state       = "CA"
	zip_code    = "12323"
	description = "Test branch router."
}
`, rName)
}

func testAccCheckBranchRouterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("branch_router Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no branch_router ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundBranchRouter := &goaviatrix.BranchRouter{
			Name: rs.Primary.Attributes["name"],
		}

		_, err := client.GetBranchRouter(foundBranchRouter)
		if err != nil {
			return err
		}
		if foundBranchRouter.Name != rs.Primary.ID {
			return fmt.Errorf("branch_router not found")
		}

		return nil
	}
}

func testAccCheckBranchRouterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_branch_router" {
			continue
		}
		foundBranchRouter := &goaviatrix.BranchRouter{
			Name: rs.Primary.Attributes["name"],
		}
		_, err := client.GetBranchRouter(foundBranchRouter)
		if err == nil {
			return fmt.Errorf("branch_router still exists")
		}
	}

	return nil
}
