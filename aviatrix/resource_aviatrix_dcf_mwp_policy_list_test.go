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

func TestAccAviatriDcfMwpPolicyList_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_MWP_POLICY_LIST")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF MWP Policy List test as SKIP_DCF_MWP_POLICY_LIST is set")
	}
	resourceName := "aviatrix_dcf_mwp_policy_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfMwpPolicyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfMwpPolicyListBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfMwpPolicyListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-mwp-policy-list"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.name", "test-distributed-firewalling-policy"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.src_smart_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.dst_smart_groups.#", "1"),
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

func testAccCheckDcfMwpPolicyListBasic() string {
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

resource "aviatrix_dcf_mwp_policy_list" "test" {
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
`
}

func testAccCheckDcfMwpPolicyListExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF MWP Policy List resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF MWP Policy List ID is set")
		}

		client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
		}

		_, err := client.GetDCFPolicyList(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF MWP Policy List status: %w", err)
		}

		return nil
	}
}

func testAccCheckDcfMwpPolicyListDestroy(s *terraform.State) error {
	client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_smart_group" {
			continue
		}

		_, err := client.GetDCFPolicyList(context.Background(), rs.Primary.ID)
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("dcf mwp list configured when it should be destroyed %w", err)
		}
	}

	return nil
}
