package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDcfPolicyGroup_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_POLICY_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Policy Group test as SKIP_DCF_POLICY_GROUP is set")
	}
	resourceName := "aviatrix_dcf_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfPolicyGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfPolicyGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-policy-group"),
					resource.TestCheckResourceAttr(resourceName, "ruleset_reference.0.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_group_reference.0.priority", "1"),
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

func testAccCheckDcfPolicyGroupBasic() string {
	return `resource "aviatrix_smart_group" "ad1" {
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

resource "aviatrix_dcf_ruleset" "test_list" {
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

resource "aviatrix_dcf_policy_group" "nested_group" {
	name = "test-nested-dcf-policy-group"
	ruleset_reference {
		priority = 0
		target_uuid = aviatrix_dcf_ruleset.test_list.id
	}
}

resource "aviatrix_dcf_policy_group" "test" {
	name = "test-dcf-policy-group"
	ruleset_reference {
		priority = 0
		target_uuid = aviatrix_dcf_ruleset.test_list.id
	}
	policy_group_reference {
		priority = 1
		target_uuid = aviatrix_dcf_policy_group.nested_group.id
	}
}
`
}

func testAccCheckDcfPolicyGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF Policy Group resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF Policy Group ID is set")
		}

		client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
		}

		_, err := client.GetDCFPolicyBlock(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF Policy Group status: %w", err)
		}

		return nil
	}
}

func testAccCheckDcfPolicyGroupDestroy(s *terraform.State) error {
	client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_policy_group" {
			continue
		}

		_, err := client.GetDCFPolicyBlock(context.Background(), rs.Primary.ID)
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("dcf policy group configured when it should be destroyed %w", err)
		}
	}

	return nil
}
