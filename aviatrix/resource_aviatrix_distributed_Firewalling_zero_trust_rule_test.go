package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixDistributedFirewallingZeroTrustRule_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_ZERO_TRUST_RULE")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed Firewalling Zero Trust Rule test as SKIP_DISTRIBUTED_FIREWALLING_ZERO_TRUST_RULE is set")
	}
	resourceName := "aviatrix_distributed_firewalling_zero_trust_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccDistributedFirewallingZeroTrustRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDistributedFirewallingZeroTrustRuleBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedFirewallingZeroTrustRuleExists(resourceName),
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

func testAccDistributedFirewallingZeroTrustRuleBasic() string {
	return `
resource "aviatrix_distributed_firewalling_zero_trust_rule" "test" {
    action  = "PERMIT"
    logging = true
}
`
}

func testAccCheckDistributedFirewallingZeroTrustRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Distributed Firewalling Zero Trust Rule resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Distributed Firewalling Zero Trust Rule ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetDistributedFirewallingZeroTrustRule(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get Distributed Firewalling Zero Trust Rule status: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("distributed firewalling zero trust rule ID not found")
		}

		return nil
	}
}

func testAccDistributedFirewallingZeroTrustRuleDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_distributed_firewalling_zero_trust_rule" {
			continue
		}

		_, err := client.GetDistributedFirewallingZeroTrustRule(context.Background())
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("distributed firewalling zero trust rule configured when it should be destroyed")
		}
	}

	return nil
}
