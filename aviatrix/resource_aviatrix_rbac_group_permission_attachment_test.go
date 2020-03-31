package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixRbacGroupPermissionAttachment_basic(t *testing.T) {
	var rbacGroupPermissionAttachment goaviatrix.RbacGroupPermissionAttachment

	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipAcc := os.Getenv("SKIP_RBAC_GROUP_PERMISSION_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping rbac group permission attachment tests as 'SKIP_RBAC_GROUP_PERMISSION_ATTACHMENT' is set")
	}

	resourceName := "aviatrix_rbac_group_permission_attachment.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRbacGroupPermissionAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRbacGroupPermissionAttachmentConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRbacGroupPermissionAttachmentExists(resourceName, &rbacGroupPermissionAttachment),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tf-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "permission_name", fmt.Sprintf("all_write")),
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

func testAccRbacGroupPermissionAttachmentConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_rbac_group" "test" {
	group_name = "tf-%s"
}
resource "aviatrix_rbac_group_permission_attachment" "test" {
	group_name      = aviatrix_rbac_group.test.group_name
	permission_name = "all_write"
}
	`, rName)
}

func testAccCheckRbacGroupPermissionAttachmentExists(n string, rAttachment *goaviatrix.RbacGroupPermissionAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("RbacGroupPermissionAttachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no RbacGroupPermissionAttachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAttachment := &goaviatrix.RbacGroupPermissionAttachment{
			GroupName:      rs.Primary.Attributes["group_name"],
			PermissionName: rs.Primary.Attributes["permission_name"],
		}

		foundAttachment2, err := client.GetRbacGroupPermissionAttachment(foundAttachment)
		if err != nil {
			return err
		}
		if foundAttachment2.GroupName != rs.Primary.Attributes["group_name"] {
			return fmt.Errorf("'group_name' Not found in created attributes")
		}
		if foundAttachment2.PermissionName != rs.Primary.Attributes["permission_name"] {
			return fmt.Errorf("'permission_name' Not found in created attributes")
		}

		*rAttachment = *foundAttachment2
		return nil
	}
}

func testAccCheckRbacGroupPermissionAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_rbac_group_permission_attachment" {
			continue
		}
		foundAttachment := &goaviatrix.RbacGroupPermissionAttachment{
			GroupName:      rs.Primary.Attributes["group_name"],
			PermissionName: rs.Primary.Attributes["permission_name"],
		}

		_, err := client.GetRbacGroupPermissionAttachment(foundAttachment)
		if err != goaviatrix.ErrNotFound {
			if strings.Contains(err.Error(), "is invalid") {
				return nil
			}
			return fmt.Errorf("rbac group user attachment still exists")
		}
	}

	return nil
}
