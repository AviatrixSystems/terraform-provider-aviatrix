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

func TestAccAviatrixDcfPolicyBlock_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_POLICY_BLOCK")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Policy Block test as SKIP_DCF_POLICY_BLOCK is set")
	}
	resourceName := "aviatrix_dcf_mwp_policy_block.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfPolicyBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfPolicyBlockBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfPolicyBlockExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-policy-block"),
					resource.TestCheckResourceAttr(resourceName, "policy_list_reference.0.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_block_reference.0.priority", "1"),
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

func testAccCheckDcfPolicyBlockBasic() string {
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

resource "aviatrix_dcf_mwp_policy_list" "test_list" {
	name = "test-dcf-mwp-policy-list"
	policies {
		name             = "test-distributed-firewalling-policy"
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

resource "aviatrix_dcf_mwp_policy_block" "nested_block" {
	name = "test-nested-dcf-policy-block"
	policy_list_reference {
		priority = 0
		target_uuid = aviatrix_dcf_mwp_policy_list.test_list.id
	}
}

resource "aviatrix_dcf_mwp_policy_block" "test" {
	name = "test-dcf-policy-block"
	policy_list_reference {
		priority = 0
		target_uuid = aviatrix_dcf_mwp_policy_list.test_list.id
	}
	policy_block_reference {
		priority = 1
		target_uuid = aviatrix_dcf_mwp_policy_block.nested_block.id
	}
}
`
}

func testAccCheckDcfPolicyBlockExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF Policy Block resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF Policy Block ID is set")
		}

		client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
		}

		_, err := client.GetDCFPolicyBlock(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF Policy Block status: %w", err)
		}

		return nil
	}
}

func testAccCheckDcfPolicyBlockDestroy(s *terraform.State) error {
	client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_mwp_policy_block" {
			continue
		}

		_, err := client.GetDCFPolicyBlock(context.Background(), rs.Primary.ID)
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("dcf policy block configured when it should be destroyed %w", err)
		}
	}

	return nil
}
