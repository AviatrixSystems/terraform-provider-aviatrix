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

func TestAccAviatrixDistributedFirewallingOriginCertEnforcementConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_ORIGIN_CERT_ENFORCEMENT_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed-firewalling origin cert enforcement config tests as SKIP_DISTRIBUTED_FIREWALLING_ORIGIN_CERT_ENFORCEMENT_CONFIG is set")
	}

	resourceName := "aviatrix_distributed_firewalling_origin_cert_enforcement_config.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccDistributedFirewallingOriginCertEnforcementConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDistributedFirewallingOriginCertEnforcementConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDistributedFirewallingOriginCertEnforcementConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enforcement_level", "Strict"),
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

func testAccDistributedFirewallingOriginCertEnforcementConfigBasic() string {
	return `
resource "aviatrix_distributed_firewalling_origin_cert_enforcement_config" "test" {
	enforcement_level = "Strict"
}
	`
}

func testAccDistributedFirewallingOriginCertEnforcementConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("distributed-firewalling origin cert enforcement config Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no distributed-firewalling origin cert enforcement config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("distributed-firewalling origin cert enforcement config ID not found")
		}

		return nil
	}
}

func testAccDistributedFirewallingOriginCertEnforcementConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_distributed_firewalling_origin_cert_enforcement_config" {
			continue
		}

		enforcementLevel, err := client.GetEnforcementLevel(context.Background())
		if err != nil {
			return fmt.Errorf("could not retrieve distributed firewalling origin cert enforcement config")
		}
		if enforcementLevel.Level == "Strict" || enforcementLevel.Level == "Ignore" {
			return fmt.Errorf("distributed firewalling origin cert enforcement config still exists")
		}
	}

	return nil
}
