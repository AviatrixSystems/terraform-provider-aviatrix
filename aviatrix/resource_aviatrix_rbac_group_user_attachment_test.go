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

func TestAccAviatrixRbacGroupUserAttachment_basic(t *testing.T) {
	var rbacGroupUserAttachment goaviatrix.RbacGroupUserAttachment

	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipAcc := os.Getenv("SKIP_RBAC_GROUP_USER_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping rbac group user attachment tests as SKIP_RBAC_GROUP_USER_ATTACHMENT is set")
	}

	resourceName := "aviatrix_rbac_group_user_attachment.test"
	msgCommon := ". Set SKIP_RBAC_GROUP_USER_ATTACHMENT to 'yes' to skip rbac group user attachment tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRbacGroupUserAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRbacGroupUserAttachmentConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRbacGroupUserAttachmentExists(resourceName, &rbacGroupUserAttachment),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tf-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "user_name", fmt.Sprintf("tf-user-%s", rName)),
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

func testAccRbacGroupUserAttachmentConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_rbac_group" "test" {
	group_name = "tf-%s"
}
resource "aviatrix_account_user" "test" {
	username = "tf-user-%s"
	email    = "abc@xyz.com"
	password = "Password-1234"
}
resource "aviatrix_rbac_group_user_attachment" "test" {
	group_name = aviatrix_rbac_group.test.group_name
	user_name  = aviatrix_account_user.test.username
}
	`, rName, rName)
}

func testAccCheckRbacGroupUserAttachmentExists(n string, rAttachment *goaviatrix.RbacGroupUserAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("RbacGroupUserAttachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no RbacGroupUserAttachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAttachment := &goaviatrix.RbacGroupUserAttachment{
			GroupName: rs.Primary.Attributes["group_name"],
			UserName:  rs.Primary.Attributes["user_name"],
		}

		foundAttachment2, err := client.GetRbacGroupUserAttachment(foundAttachment)
		if err != nil {
			return err
		}
		if foundAttachment2.GroupName != rs.Primary.Attributes["group_name"] {
			return fmt.Errorf("'group_name' Not found in created attributes")
		}
		if foundAttachment2.UserName != rs.Primary.Attributes["user_name"] {
			return fmt.Errorf("'user_name' Not found in created attributes")
		}

		*rAttachment = *foundAttachment2
		return nil
	}
}

func testAccCheckRbacGroupUserAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_rbac_group_user_attachment" {
			continue
		}
		foundAttachment := &goaviatrix.RbacGroupUserAttachment{
			GroupName: rs.Primary.Attributes["group_name"],
			UserName:  rs.Primary.Attributes["user_name"],
		}

		_, err := client.GetRbacGroupUserAttachment(foundAttachment)
		if err != goaviatrix.ErrNotFound {
			if strings.Contains(err.Error(), "is invalid") {
				return nil
			}
			return fmt.Errorf("rbac group user attachment still exists")
		}
	}

	return nil
}
