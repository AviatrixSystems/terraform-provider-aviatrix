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

func TestAccAviatrixDistributedFirewallingDefaultActionRule_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_DEFAULT_ACTION_RULE")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed Firewalling Default Action Rule test as SKIP_DISTRIBUTED_FIREWALLING_DEFAULT_ACTION_RULE is set")
	}
	resourceName := "aviatrix_distributed_firewalling_default_action_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccDistributedFirewallingDefaultActionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				data "aviatrix_dcf_log_profile" "start" {
				  profile_name = "start"
				}

				resource "aviatrix_distributed_firewalling_default_action_rule" "test" {
				  action      = "PERMIT"
				  logging     = true
				  log_profile = data.aviatrix_dcf_log_profile.start.profile_id
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedFirewallingDefaultActionRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceName, "logging", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "log_profile"),
					resource.TestCheckResourceAttrPair(resourceName, "log_profile", "data.aviatrix_dcf_log_profile.start", "profile_id"),
				),
			},
		},
	})
}

func testAccCheckDistributedFirewallingDefaultActionRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Distributed Firewalling Default Action Rule resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Distributed Firewalling Default Action Rule ID is set")
		}

		meta := testAccProviderVersionValidation.Meta()
		client := mustClient(meta)

		_, err := client.GetDistributedFirewallingDefaultActionRule(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get Distributed Firewalling Default Action Rule status: %w", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("distributed firewalling default action rule ID not found")
		}

		return nil
	}
}

func testAccDistributedFirewallingDefaultActionRuleDestroy(s *terraform.State) error {
	meta := testAccProviderVersionValidation.Meta()
	client := mustClient(meta)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_distributed_firewalling_default_action_rule" {
			continue
		}

		rule, err := client.GetDistributedFirewallingDefaultActionRule(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get distributed firewalling default action rule: %w", err)
		}

		// Check that the rule has been reset to defaults
		if rule.Action != "PERMIT" || rule.Logging != false || rule.LogProfile != "" {
			return fmt.Errorf("distributed firewalling default action rule not reset to defaults: action=%s, logging=%t", rule.Action, rule.Logging)
		}
	}

	return nil
}
