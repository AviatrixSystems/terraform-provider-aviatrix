package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
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
				Config: testAccDistributedFirewallingDefaultActionRuleBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedFirewallingDefaultActionRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceName, "logging", "true"),
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

func testAccDistributedFirewallingDefaultActionRuleBasic() string {
	return `resource "aviatrix_distributed_firewalling_default_action_rule" "test" {
    			action  = "PERMIT"
    			logging = true
			}`
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
		client, ok := meta.(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to assert meta as *goaviatrix.Client")
		}

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
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to assert meta as *goaviatrix.Client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_distributed_firewalling_default_action_rule" {
			continue
		}

		_, err := client.GetDistributedFirewallingDefaultActionRule(context.Background())
		if err == nil || !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("distributed firewalling default action rule configured when it should be destroyed")
		}
	}

	return nil
}
