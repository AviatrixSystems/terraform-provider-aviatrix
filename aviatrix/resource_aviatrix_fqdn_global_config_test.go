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

func TestAccAviatrixFQDNGlobalConfig_basic(t *testing.T) {
	if os.Getenv("SKIP_FQDN_GLOBAL_CONFIG") == "yes" {
		t.Skip("Skipping FQDN global config test as SKIP_FQDN_GLOBAL_CONFIG is set")
	}

	resourceName := "aviatrix_fqdn_global_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFQDNGlobalConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFQDNGlobalConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFQDNGlobalConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_exception_rule", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_private_network_filtering", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_caching", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_exact_match", "true"),
				),
			},
		},
	})
}

func testAccFQDNGlobalConfigBasic() string {
	return `
resource "aviatrix_fqdn_global_config" "test" {
	enable_exception_rule            = false
	enable_private_network_filtering = true
	enable_caching                   = false
	enable_exact_match               = true
}
`
}

func testAccCheckFQDNGlobalConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("fqdn global config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no fqdn global config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("fqdn global config ID not found")
		}

		return nil
	}
}

func testAccCheckFQDNGlobalConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_fqdn_global_config" {
			continue
		}

		privateSubFilter, err := client.GetFQDNPrivateNetworkFilteringStatus(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get FQDN private network filter status: %s", err)
		}
		if privateSubFilter.PrivateSubFilter == "enabled" {
			return fmt.Errorf("failed to destroy fqdn global config")
		}
	}

	return nil
}
