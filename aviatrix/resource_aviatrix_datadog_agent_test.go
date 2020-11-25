package aviatrix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDatadogAgent_basic(t *testing.T) {
	resourceName := "aviatrix_datadog_agent.test_datadog_agent"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatadogAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatadogAgentBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogAgentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "api_key", "your_api_key"),
					resource.TestCheckResourceAttr(resourceName, "site", "datadoghq.com"),
					testAccCheckDatadogAgentExcludedGatewaysMatch([]string{"a", "b"}),
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

func testAccDatadogAgentBasic() string {
	return `
resource "aviatrix_datadog_agent" "test_datadog_agent" {
	api_key           = "your_api_key"
	site              = "datadoghq.com"
	excluded_gateways = ["a", "b"]
}
`
}

func testAccCheckDatadogAgentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("datadog agent not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetDatadogAgentStatus()
		if resp.Status != "enabled" {
			return fmt.Errorf("datadog agent not found")
		}

		return nil
	}
}

func testAccCheckDatadogAgentExcludedGatewaysMatch(input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetDatadogAgentStatus()
		if !goaviatrix.Equivalent(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}
		return nil
	}
}

func testAccCheckDatadogAgentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_datadog_agent" {
			continue
		}

		resp, _ := client.GetDatadogAgentStatus()
		if resp.Status == "enabled" {
			return fmt.Errorf("datadog_agent still exists")
		}
	}

	return nil
}
