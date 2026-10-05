package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixLinkHierarchy_basic(t *testing.T) {
	if os.Getenv("SKIP_LINK_HIERARCHY") == "yes" {
		t.Skip("Skipping link hierarchy test as SKIP_LINK_HIERARCHY is set")
	}

	resourceName := "aviatrix_link_hierarchy.test"
	linkHierarchyName := "lh-" + acctest.RandString(5)
	linkName := "l-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinkHierarchyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLinkHierarchyBasic(linkHierarchyName, linkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinkHierarchyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", linkHierarchyName),
					resource.TestCheckResourceAttr(resourceName, "links.0.name", linkName),
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

func testAccLinkHierarchyBasic(linkHierarchyName, linkName string) string {
	return fmt.Sprintf(`
resource "aviatrix_link_hierarchy" "test" {
	name = "%s"

	links {
		name = "%s"
		wan_link {
			wan_tag = "wan3.10"
		}
	}
}
 `, linkHierarchyName, linkName)
}

func testAccCheckLinkHierarchyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("link hierarchy not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetLinkHierarchy(context.Background(), rs.Primary.Attributes["uuid"])
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("link hierarchy not found")
		}

		return nil
	}
}

func testAccCheckLinkHierarchyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_link_hierarchy" {
			continue
		}

		_, err := client.GetLinkHierarchy(context.Background(), rs.Primary.Attributes["uuid"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("link hierarchy still exists")
		}
	}

	return nil
}
