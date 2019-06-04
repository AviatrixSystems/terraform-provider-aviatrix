package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixFirewallTag_basic(t *testing.T) {
	var ftag goaviatrix.FirewallTag
	rInt := acctest.RandInt()
	resourceName := "aviatrix_firewall_tag.foo"

	skipAcc := os.Getenv("SKIP_FIREWALL_TAG")
	if skipAcc == "yes" {
		t.Skip("Skipping firewall tag test as SKIP_FIREWALL_TAG is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallTagConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallTagExists("aviatrix_firewall_tag.foo", &ftag),
					resource.TestCheckResourceAttr(
						resourceName, "firewall_tag", fmt.Sprintf("tft-%d", rInt)),
					resource.TestCheckResourceAttr(
						resourceName, "cidr_list.#", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "cidr_list.0.cidr", "10.1.0.0/24"),
					resource.TestCheckResourceAttr(
						resourceName, "cidr_list.0.cidr_tag_name", "a1"),
					resource.TestCheckResourceAttr(
						resourceName, "cidr_list.1.cidr", "10.2.0.0/24"),
					resource.TestCheckResourceAttr(
						resourceName, "cidr_list.1.cidr_tag_name", "b1"),
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

func testAccFirewallTagConfigBasic(rInt int) string {
	return fmt.Sprintf(`
resource "aviatrix_firewall_tag" "foo" {
	firewall_tag = "tft-%d"
	cidr_list = [
	{
		cidr_tag_name = "a1"
		cidr = "10.1.0.0/24"
	},
	{
		cidr_tag_name = "b1"
		cidr = "10.2.0.0/24"
	}
	]
}
	`, rInt)
}

func testAccCheckFirewallTagExists(n string, firewallTag *goaviatrix.FirewallTag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("firewall tag Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no tag ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundTag := &goaviatrix.FirewallTag{
			Name: rs.Primary.Attributes["firewall_tag"],
		}

		_, err := client.GetFirewallTag(foundTag)
		if err != nil {
			return err
		}
		if foundTag.Name != rs.Primary.ID {
			return fmt.Errorf("firewall tag not found")
		}
		*firewallTag = *foundTag

		return nil
	}
}

func testAccCheckFirewallTagDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_firewall_tag" {
			continue
		}
		foundTag := &goaviatrix.FirewallTag{
			Name: rs.Primary.Attributes["firewall_tag"],
		}
		_, err := client.GetFirewallTag(foundTag)

		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("firewall tag still exists after destroy")
		}
	}
	return nil
}
