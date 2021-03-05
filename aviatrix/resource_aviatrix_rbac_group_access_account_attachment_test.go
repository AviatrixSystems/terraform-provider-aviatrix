package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixRbacGroupAccessAccountAttachment_basic(t *testing.T) {
	var rbacGroupAccessAccountAttachment goaviatrix.RbacGroupAccessAccountAttachment

	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_RBAC_GROUP_ACCESS_ACCOUNT_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping rbac group access account attachment tests as SKIP_RBAC_GROUP_ACCESS_ACCOUNT_ATTACHMENT is set")
	}

	resourceName := "aviatrix_rbac_group_access_account_attachment.test"
	msgCommon := ". Set SKIP_RBAC_GROUP_ACCESS_ACCOUNT_ATTACHMENT to yes to skip rbac group access account attachment tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRbacGroupAccessAccountAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRbacGroupAccessAccountAttachmentConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRbacGroupAccessAccountAttachmentExists(resourceName, &rbacGroupAccessAccountAttachment),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tf-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "access_account_name", fmt.Sprintf("tf-acc-%s", rName)),
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

func testAccRbacGroupAccessAccountAttachmentConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_rbac_group" "test" {
	group_name = "tf-%s"
}
resource "aviatrix_account" "test_account" {
	account_name       = "tf-acc-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_rbac_group_access_account_attachment" "test" {
	group_name 			= aviatrix_rbac_group.test.group_name
	access_account_name = aviatrix_account.test_account.account_name
}
	`, rName, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func testAccCheckRbacGroupAccessAccountAttachmentExists(n string, rAttachment *goaviatrix.RbacGroupAccessAccountAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("RbacGroupAccessAccountAttachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no RbacGroupAccessAccountAttachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAttachment := &goaviatrix.RbacGroupAccessAccountAttachment{
			GroupName:         rs.Primary.Attributes["group_name"],
			AccessAccountName: rs.Primary.Attributes["access_account_name"],
		}

		foundAttachment2, err := client.GetRbacGroupAccessAccountAttachment(foundAttachment)
		if err != nil {
			return err
		}
		if foundAttachment2.GroupName != rs.Primary.Attributes["group_name"] {
			return fmt.Errorf("'group_name' Not found in created attributes")
		}
		if foundAttachment2.AccessAccountName != rs.Primary.Attributes["access_account_name"] {
			return fmt.Errorf("'access_account_name' Not found in created attributes")
		}

		*rAttachment = *foundAttachment2
		return nil
	}
}

func testAccCheckRbacGroupAccessAccountAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_rbac_group_access_account_attachment" {
			continue
		}
		foundAttachment := &goaviatrix.RbacGroupAccessAccountAttachment{
			GroupName:         rs.Primary.Attributes["group_name"],
			AccessAccountName: rs.Primary.Attributes["access_account_name"],
		}

		_, err := client.GetRbacGroupAccessAccountAttachment(foundAttachment)
		if err != goaviatrix.ErrNotFound {
			if strings.Contains(err.Error(), "is invalid") {
				return nil
			}
			return fmt.Errorf("rbac group access account attachment still exists")
		}
	}

	return nil
}
