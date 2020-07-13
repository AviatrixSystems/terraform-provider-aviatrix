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

func TestAccAviatrixBranchRouterRegistration_basic(t *testing.T) {
	if os.Getenv("SKIP_BRANCH_ROUTER_REGISTRATION") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_BRANCH_ROUTER_REGISTRATION is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_branch_router_registration.test_branch_router"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			branchRouterRegistrationPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBranchRouterRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchRouterRegistrationBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBranchRouterRegistrationExists(resourceName),
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

func testAccBranchRouterRegistrationBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_branch_router_registration" "test_branch_router" {
	name        = "branchrouter-%s"
	public_ip   = "%s"
	username    = "ec2-user"
	key_file    = "%s"
	host_os     = "ios"
	ssh_port    = 22
	address_1   = "2901 Tasman Dr"
	address_2   = "Suite #104"
	city        = "Santa Clara"
	state       = "CA"
	zip_code    = "12323"
	description = "Test branch router."
}
`, rName, os.Getenv("BRANCH_ROUTER_PUBLIC_IP"), os.Getenv("BRANCH_ROUTER_KEY_FILE_PATH"))
}

func testAccCheckBranchRouterRegistrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("branch_router_registration Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no branch_router_registration ID is set")
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
			return fmt.Errorf("branch_router_registration not found")
		}

		return nil
	}
}

func testAccCheckBranchRouterRegistrationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_branch_router_registration" {
			continue
		}
		foundBranchRouter := &goaviatrix.BranchRouter{
			Name: rs.Primary.Attributes["name"],
		}
		_, err := client.GetBranchRouter(foundBranchRouter)
		if err == nil {
			return fmt.Errorf("branch_router_registration still exists")
		}
	}

	return nil
}

func branchRouterRegistrationPreCheck(t *testing.T) {
	if os.Getenv("BRANCH_ROUTER_PUBLIC_IP") == "" {
		t.Fatal("environment variable BRANCH_ROUTER_PUBLIC_IP must be set for branch_router_registration acceptance test")
	}

	if os.Getenv("BRANCH_ROUTER_KEY_FILE_PATH") == "" {
		t.Fatal("environment variable BRANCH_ROUTER_KEY_FILE_PATH must be set for " +
			"branch_router_registration acceptance test")
	}
}
