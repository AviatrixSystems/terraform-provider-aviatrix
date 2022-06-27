package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAviatrixFQDNGlobalConfigs_basic(t *testing.T) {
	if os.Getenv("SKIP_FQDN_GLOBAL_CONFIGS") == "yes" {
		t.Skip("Skipping FQDN global configs test as SKIP_FQDN_GLOBAL_CONFIGS is set")
	}

	resourceName := "aviatrix_fqdn_global_configs.test_fqdn_global_configs"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFQDNGlobalConfigsBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "exception_rule", "false"),
					resource.TestCheckResourceAttr(resourceName, "caching", "false"),
					resource.TestCheckResourceAttr(resourceName, "exact_match", "true"),
					resource.TestCheckResourceAttr(resourceName, "network_filtering", "Customize Network Filtering"),
				),
			},
		},
	})
}

func testAccFQDNGlobalConfigsBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_fqdn_global_configs" "test_fqdn_global_configs" {
  exception_rule    = false
  network_filtering = "Customize Network Filtering"
  configured_ips    = ["172.16.0.0/12~~RFC-1918", "10.0.0.0/8~~RFC-1918"]
  caching           = false
  exact_match       = true
}
`)
}
