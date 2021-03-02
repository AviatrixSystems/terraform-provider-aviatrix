package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixFQDNTagRule_basic(t *testing.T) {
	if os.Getenv("SKIP_FQDN_TAG_RULE") == "yes" {
		t.Skip("Skipping fqdn tag rule test as SKIP_FQDN_TAG_RULE is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_fqdn_tag_rule.test_fqdn_tag_rule"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFQDNDomainNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFQDNDomainNameBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFQDNDomainNameExists(resourceName),
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

func testAccFQDNDomainNameBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_fqdn" "foo" {
	fqdn_tag            = "fqdn-%s"
	fqdn_enabled        = true
	fqdn_mode           = "white"
	manage_domain_names = false
}

resource "aviatrix_fqdn_tag_rule" "test_fqdn_tag_rule" {
	fqdn_tag_name = aviatrix_fqdn.foo.fqdn_tag
	fqdn          = "*.aviatrix.com"
	protocol      = "tcp"
	port          = "443"
	action        = "Allow"
}
`, rName)
}

func testAccCheckFQDNDomainNameExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("fqdn_tag_rule Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no fqdn_tag_rule ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		fqdn := &goaviatrix.FQDN{
			FQDNTag: rs.Primary.Attributes["fqdn_tag_name"],
			DomainList: []*goaviatrix.Filters{
				{
					FQDN:     rs.Primary.Attributes["fqdn"],
					Protocol: rs.Primary.Attributes["protocol"],
					Port:     rs.Primary.Attributes["port"],
					Verdict:  rs.Primary.Attributes["action"],
				},
			},
		}

		fqdn, err := client.GetFQDNTagRule(fqdn)
		if err != nil {
			return err
		}
		if getFQDNTagRuleID(fqdn) != rs.Primary.ID {
			return fmt.Errorf("fqdn_tag_rule not found")
		}

		return nil
	}
}

func testAccCheckFQDNDomainNameDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_fqdn_tag_rule" {
			continue
		}
		fqdn := &goaviatrix.FQDN{
			FQDNTag: rs.Primary.Attributes["fqdn_tag_name"],
			DomainList: []*goaviatrix.Filters{
				{
					FQDN:     rs.Primary.Attributes["fqdn"],
					Protocol: rs.Primary.Attributes["protocol"],
					Port:     rs.Primary.Attributes["port"],
					Verdict:  rs.Primary.Attributes["action"],
				},
			},
		}
		_, err := client.GetFQDNTagRule(fqdn)
		if err == nil {
			return fmt.Errorf("fqdn_tag_rule still exists")
		}
	}

	return nil
}
