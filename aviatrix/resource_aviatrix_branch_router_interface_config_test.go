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

func TestAccAviatrixBranchRouterInterfaceConfig_basic(t *testing.T) {
	if os.Getenv("SKIP_BRANCH_ROUTER_INTERFACE_CONFIG") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_BRANCH_ROUTER_INTERFACE_CONFIG is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_branch_router_interface_config.test_branch_router_interface_config"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			branchRouterPreCheck(t)
			branchRouterInterfaceConfigPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchRouterInterfaceConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBranchRouterInterfaceConfigExists(resourceName),
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

func testAccBranchRouterInterfaceConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_branch_router" "test_branch_router" {
	name        = "branchrouter-%s"
	public_ip   = "%[2]s"
	username    = "ec2-user"
	key_file    = "%[3]s"
	host_os     = "ios"
	ssh_port    = 22
	address_1   = "2901 Tasman Dr"
	address_2   = "Suite #104"
	city        = "Santa Clara"
	state       = "CA"
	zip_code    = "12323"
	description = "Test branch router."
}

resource "aviatrix_branch_router_interface_config" "test_branch_router_interface_config" {
	branch_router_name              = aviatrix_branch_router.test_branch_router.name
	wan_primary_interface           = "%[4]s"
	wan_primary_interface_public_ip = "%[2]s"
}
`, rName, os.Getenv("BRANCH_ROUTER_PUBLIC_IP"), os.Getenv("BRANCH_ROUTER_KEY_FILE_PATH"), os.Getenv("BRANCH_ROUTER_PRIMARY_INTERFACE"))
}

func testAccCheckBranchRouterInterfaceConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("branch_router_interface_config Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no branch_router_interface_config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		br := &goaviatrix.BranchRouter{Name: rs.Primary.Attributes["branch_router_name"]}

		br, err := client.GetBranchRouter(br)
		if err != nil {
			return err
		}

		if br.Name != rs.Primary.ID ||
			br.PrimaryInterface != rs.Primary.Attributes["wan_primary_interface"] ||
			br.PrimaryInterfaceIP != rs.Primary.Attributes["wan_primary_interface_public_ip"] {
			return fmt.Errorf("branch_router_interface_config not found")
		}

		return nil
	}
}

func branchRouterInterfaceConfigPreCheck(t *testing.T) {
	if os.Getenv("BRANCH_ROUTER_PRIMARY_INTERFACE") == "" {
		t.Fatal("environment variable BRANCH_ROUTER_PRIMARY_INTERFACE must be set for " +
			"branch_router_interface_config acceptance test")
	}
}
