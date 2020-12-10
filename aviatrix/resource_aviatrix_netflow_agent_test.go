package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixNetflowAgent_basic(t *testing.T) {
	if os.Getenv("SKIP_NETFLOW_AGENT") == "yes" {
		t.Skip("Skipping netflow agent test as SKIP_NETFLOW_AGENT is set")
	}

	resourceName := "aviatrix_netflow_agent.test_netflow_agent"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetflowAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetflowAgentBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetflowAgentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "server_ip", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "port", "10"),
					resource.TestCheckResourceAttr(resourceName, "version", "5"),
					testAccCheckNetflowAgentExcludedGatewaysMatch([]string{"a", "b"}),
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

func testAccNetflowAgentBasic() string {
	return `
resource "aviatrix_netflow_agent" "test_netflow_agent" {
	server_ip         = "1.2.3.4"
	port              = 10
	excluded_gateways = ["a", "b"]
}
`
}

func testAccCheckNetflowAgentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("netflow agent not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetNetflowAgentStatus()
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("netflow agent not found")
		}

		return nil
	}
}

func testAccCheckNetflowAgentExcludedGatewaysMatch(input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetNetflowAgentStatus()
		if !goaviatrix.Equivalent(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}
		return nil
	}
}

func testAccCheckNetflowAgentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_netflow_agent" {
			continue
		}

		_, err := client.GetNetflowAgentStatus()
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("netflow_agent still exists")
		}
	}

	return nil
}
