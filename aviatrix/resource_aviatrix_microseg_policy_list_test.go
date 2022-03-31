package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixMicroseg_Policy_List_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_MICROSEG_POLICY_LIST")
	if skipAcc == "yes" {
		t.Skip("Skipping Microseg Policy List test as SKIP_MICROSEG_POLICY_LIST is set")
	}
	resourceName := "aviatrix_microseg_policy_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccMicrosegPolicyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMicrosegPolicyListBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicrosegPolicyListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "policies.0.name", "test-microseg-policy"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.src_app_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.dst_app_domains.#", "1"),
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

func testAccMicrosegPolicyListBasic() string {
	return `
resource "aviatrix_app_domain" "ad1" {
	name      = "test-app-domain-1"
	selector {
		match_expressions {
			cidr = "10.0.0.0/16"
		}
	}
}

resource "aviatrix_app_domain" "ad2" {
	name       = "test-app-domain-2"
	selector {
		match_expressions {
			cidr = "11.0.0.0/16"
		}
	}
}

resource "aviatrix_microseg_policy_list" "test" {
	policies {
		name            = "test-microseg-policy"
		action          = "PERMIT"
		logging         = true
		priority        = 0
		protocol        = "TCP"
		src_app_domains = [
		  aviatrix_app_domain.ad1.uuid
		]
		dst_app_domains = [
		  aviatrix_app_domain.ad2.uuid
		]
	
		port_ranges {
		  hi = 0
		  lo = 0
		}
  }
}
`
}

func testAccCheckMicrosegPolicyListExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Micro-segmentation Policy List resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Micro-segmentation Policy List ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetMicrosegPolicyList(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get Micro-segmentation Policy List status: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("micro-segmentation policy list ID not found")
		}

		return nil
	}
}

func testAccMicrosegPolicyListDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_app_domain" {
			continue
		}

		_, err := client.GetMicrosegPolicyList(context.Background())
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("micro-segmentation policy list configured when it should be destroyed")
		}
	}

	return nil
}
