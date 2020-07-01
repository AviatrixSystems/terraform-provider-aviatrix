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

func TestAccAviatrixBranchRouterVirtualWanAttachment_basic(t *testing.T) {
	if os.Getenv("SKIP_BRANCH_ROUTER_VIRTUAL_WAN_ATTACHMENT") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_BRANCH_ROUTER_VIRTUAL_WAN_ATTACHMENT is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_branch_router_virtual_wan_attachment.test_branch_router_virtual_wan_attachment"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			branchRouterVirtualWanAttachmentPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBranchRouterVirtualWanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchRouterVirtualWanAttachmentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBranchRouterVirtualWanAttachmentExists(resourceName),
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

func testAccBranchRouterVirtualWanAttachmentBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_branch_router_virtual_wan_attachment" "test_branch_router_virtual_wan_attachment" {
	connection_name       = "conn-%s"
	branch_name           = "%s"
	account_name          = "%s"
	resource_group        = "%s"
	hub_name              = "%s"
	branch_router_bgp_asn = %s
}
`, rName, os.Getenv("BRANCH_ROUTER_NAME"), os.Getenv("AZURE_ACCOUNT_NAME"),
		os.Getenv("AZURE_RESOURCE_GROUP"), os.Getenv("AZURE_HUB_NAME"), os.Getenv("BRANCH_ROUTER_ASN"))
}

func testAccCheckBranchRouterVirtualWanAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("branch_router_virtual_wan_attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no branch_router_virtual_wan_attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundBranchRouterVirtualWanAttachment := &goaviatrix.BranchRouterVirtualWanAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetBranchRouterVirtualWanAttachment(foundBranchRouterVirtualWanAttachment)
		if err != nil {
			return err
		}
		if foundBranchRouterVirtualWanAttachment.ConnectionName != rs.Primary.ID {
			return fmt.Errorf("branch_router_virtual_wan_attachment not found")
		}

		return nil
	}
}

func testAccCheckBranchRouterVirtualWanAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_branch_router_virtual_wan_attachment" {
			continue
		}
		foundBranchRouterVirtualWanAttachment := &goaviatrix.BranchRouterVirtualWanAttachment{
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}
		_, err := client.GetBranchRouterVirtualWanAttachment(foundBranchRouterVirtualWanAttachment)
		if err == nil {
			return fmt.Errorf("branch_router_virtual_wan_attachment still exists")
		}
	}

	return nil
}

func branchRouterVirtualWanAttachmentPreCheck(t *testing.T) {
	if os.Getenv("BRANCH_ROUTER_NAME") == "" {
		t.Fatal("environment variable BRANCH_ROUTER_NAME must be set for aviatrix_branch_router_virtual_wan_attachment acceptance test")
	}
	if os.Getenv("AZURE_ACCOUNT_NAME") == "" {
		t.Fatal("environment variable AZURE_ACCOUNT_NAME must be set for aviatrix_branch_router_virtual_wan_attachment acceptance test")
	}
	if os.Getenv("AZURE_RESOURCE_GROUP") == "" {
		t.Fatal("environment variable AZURE_RESOURCE_GROUP must be set for aviatrix_branch_router_virtual_wan_attachment acceptance test")
	}
	if os.Getenv("AZURE_HUB_NAME") == "" {
		t.Fatal("environment variable AZURE_HUB_NAME must be set for aviatrix_branch_router_virtual_wan_attachment acceptance test")
	}
	if os.Getenv("BRANCH_ROUTER_ASN") == "" {
		t.Fatal("environment variable BRANCH_ROUTER_ASN must be set for aviatrix_branch_router_virtual_wan_attachment acceptance test")
	}
}
