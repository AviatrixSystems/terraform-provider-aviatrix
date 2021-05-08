package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixDatadogAgent_basic(t *testing.T) {
	if os.Getenv("SKIP_DATADOG_AGENT") == "yes" {
		t.Skip("Skipping datadog agent test as SKIP_DATADOG_AGENT is set")
	}

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
					resource.TestCheckResourceAttr(resourceName, "api_key", os.Getenv("DATADOG_API_KEY")),
					resource.TestCheckResourceAttr(resourceName, "site", "datadoghq.com"),
					testAccCheckDatadogAgentExcludedGatewaysMatch([]string{"a", "b"}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

func testAccDatadogAgentBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_datadog_agent" "test_datadog_agent" {
	api_key           = "%s"
	site              = "datadoghq.com"
	excluded_gateways = ["a", "b"]
}
`, os.Getenv("DATADOG_API_KEY"))
}

func testAccCheckDatadogAgentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("datadog agent not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetDatadogAgentStatus()
		if err == goaviatrix.ErrNotFound {
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

		_, err := client.GetDatadogAgentStatus()
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("datadog_agent still exists")
		}
	}

	return nil
}
