package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixRbacGroup_basic(t *testing.T) {
	var rbacGroup goaviatrix.RbacGroup

	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipAcc := os.Getenv("SKIP_RBAC_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping rbac group tests as SKIP_RBAC_GROUP is set")
	}

	resourceName := "aviatrix_rbac_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRbacGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRbacGroupConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRbacGroupExists(resourceName, &rbacGroup),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tf-%s", rName)),
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

func testAccRbacGroupConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_rbac_group" "test" {
	group_name = "tf-%s"
}
	`, rName)
}

func testAccCheckRbacGroupExists(n string, rbacGroup *goaviatrix.RbacGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("RbacGroup Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no RbacGroup ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundRbacGroup := &goaviatrix.RbacGroup{
			GroupName: rs.Primary.Attributes["group_name"],
		}

		foundRbacGroup2, err := client.GetPermissionGroup(foundRbacGroup)
		if err != nil {
			return err
		}
		if foundRbacGroup2.GroupName != rs.Primary.ID {
			return fmt.Errorf("RbacGroup not found")
		}

		*rbacGroup = *foundRbacGroup2
		return nil
	}
}

func testAccCheckRbacGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_rbac_group" {
			continue
		}
		foundRbacGroup := &goaviatrix.RbacGroup{
			GroupName: rs.Primary.Attributes["group_name"],
		}

		_, err := client.GetPermissionGroup(foundRbacGroup)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("RbacGroup still enabled")
		}
	}

	return nil
}
