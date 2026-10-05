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

func TestAccAviatrixDistributedFirewallingPolicyList_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_POLICY_LIST")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed Firewalling Policy List test as SKIP_DISTRIBUTED_FIREWALLING_POLICY_LIST is set")
	}
	resourceName := "aviatrix_distributed_firewalling_policy_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccDistributedFirewallingPolicyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDistributedFirewallingPolicyListBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedFirewallingPolicyListExists(resourceName),
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

func testAccDistributedFirewallingPolicyListBasic() string {
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

resource "aviatrix_distributed_firewalling_policy_list" "test" {
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

func testAccCheckDistributedFirewallingPolicyListExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Distributed-firewalling Policy List resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Distributed-firewalling Policy List ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetDistributedFirewallingPolicyList(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get Distributed-firewalling Policy List status: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("distributed-firewalling policy list ID not found")
		}

		return nil
	}
}

func testAccDistributedFirewallingPolicyListDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_smart_group" {
			continue
		}

		_, err := client.GetDistributedFirewallingPolicyList(context.Background())
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("distributed-firewalling policy list configured when it should be destroyed")
		}
	}

	return nil
}
