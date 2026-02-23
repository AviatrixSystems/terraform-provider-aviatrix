package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatriDcfRuleset_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_RULESET")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Ruleset test as SKIP_DCF_RULESET is set")
	}
	resourceName := "aviatrix_dcf_ruleset.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfRulesetBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfRulesetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-ruleset"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.name", "test-distributed-firewalling-rule"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.src_smart_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.dst_smart_groups.#", "1"),
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

func testAccCheckDcfRulesetBasic() string {
	return `
resource "aviatrix_smart_group" "ad1" {
	name = "test-smart_group-1"
	selector {
		match_expressions {
			cidr = "10.0.0.0/16"
		}
	}
}

resource "aviatrix_smart_group" "ad2" {
	name = "test-smart-group-2"
	selector {
		match_expressions {
			cidr = "11.0.0.0/16"
		}
	}
}

resource "aviatrix_dcf_ruleset" "test" {
	name = "test-dcf-ruleset"
	rules {
		name             = "test-distributed-firewalling-rule"
		action           = "PERMIT"
		logging          = true
		priority         = 0
		protocol         = "TCP"
		src_smart_groups = [
		  aviatrix_smart_group.ad1.uuid
		]
		dst_smart_groups = [
		  aviatrix_smart_group.ad2.uuid
		]

		port_ranges {
		  hi = 10
		  lo = 1
		}
  }
}
`
}

func testAccCheckDcfRulesetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF Ruleset resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF Ruleset ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		_, err := client.GetDCFPolicyList(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF Ruleset status: %w", err)
		}

		return nil
	}
}

func testAccCheckDcfRulesetDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_smart_group" {
			continue
		}

		_, err := client.GetDCFPolicyList(context.Background(), rs.Primary.ID)
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("dcf ruleset configured when it should be destroyed %w", err)
		}
	}

	return nil
}
