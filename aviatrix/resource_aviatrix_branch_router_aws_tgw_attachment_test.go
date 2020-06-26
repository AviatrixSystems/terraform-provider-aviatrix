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

func TestAccAviatrixBranchRouterAwsTgwAttachment_basic(t *testing.T) {
	if os.Getenv("SKIP_BRANCH_ROUTER_AWS_TGW_ATTACHMENT") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_BRANCH_ROUTER_AWS_TGW_ATTACHMENT is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_branch_router_aws_tgw_attachment.test_branch_router_aws_tgw_attachment"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAviatrixBranchRouterAwsTgwAttachmentPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBranchRouterAwsTgwAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchRouterAwsTgwAttachmentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBranchRouterAwsTgwAttachmentExists(resourceName),
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

func testAccBranchRouterAwsTgwAttachmentBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_branch_router_aws_tgw_attachment" "test_branch_router_aws_tgw_attachment" {
	connection_name           = "conn-%s"
	branch_name               = "%s"
	aws_tgw_name              = "%s"
	branch_router_bgp_asn     = 65001
	security_domain_name      = "Default_Domain"
}
`, rName, os.Getenv("BRANCH_ROUTER_NAME"), os.Getenv("AWS_TGW_NAME"))
}

func testAccCheckBranchRouterAwsTgwAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("branch_router_aws_tgw_attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no branch_router_aws_tgw_attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		attachment := &goaviatrix.BranchRouterAwsTgwAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
			BranchName:     rs.Primary.Attributes["branch_name"],
			AwsTgwName:     rs.Primary.Attributes["aws_tgw_name"],
		}

		_, err := client.GetBranchRouterAwsTgwAttachment(attachment)
		if err != nil {
			return err
		}
		if attachment.ID() != rs.Primary.ID {
			return fmt.Errorf("branch_router_aws_tgw_attachment not found")
		}

		return nil
	}
}

func testAccCheckBranchRouterAwsTgwAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_branch_router_aws_tgw_attachment" {
			continue
		}
		attachment := &goaviatrix.BranchRouterAwsTgwAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
			BranchName:     rs.Primary.Attributes["branch_name"],
			AwsTgwName:     rs.Primary.Attributes["aws_tgw_name"],
		}
		_, err := client.GetBranchRouterAwsTgwAttachment(attachment)
		if err == nil {
			return fmt.Errorf("branch_router_aws_tgw_attachment still exists")
		}
	}

	return nil
}

func testAccAviatrixBranchRouterAwsTgwAttachmentPreCheck(t *testing.T) {
	if os.Getenv("BRANCH_ROUTER_NAME") == "" {
		t.Fatal("BRANCH_ROUTER_NAME must be set for aviatrix_branch_router_aws_tgw_attachment acceptance test.")
	}
	if os.Getenv("AWS_TGW_NAME") == "" {
		t.Fatal("AWS_TGW_NAME must be set for aviatrix_branch_router_aws_tgw_attachment acceptance test.")
	}
}
