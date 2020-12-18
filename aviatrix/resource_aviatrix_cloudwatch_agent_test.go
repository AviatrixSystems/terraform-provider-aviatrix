package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixCloudwatchAgent_basic(t *testing.T) {
	if os.Getenv("SKIP_CLOUDWATCH_AGENT") == "yes" {
		t.Skip("Skipping cloudwatch agent test as SKIP_CLOUDWATCH_AGENT is set")
	}

	resourceName := "aviatrix_cloudwatch_agent.test_cloudwatch_agent"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudwatchAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudwatchAgentBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudwatchAgentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cloudwatch_role_arn", "arn:aws:iam::469550033836:role/aviatrix-role-cloudwatch"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east-1"),
					testAccCheckCloudwatchAgentExcludedGatewaysMatch([]string{"a", "b"}),
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

func testAccCloudwatchAgentBasic() string {
	return `
resource "aviatrix_cloudwatch_agent" "test_cloudwatch_agent" {
	cloudwatch_role_arn = "arn:aws:iam::469550033836:role/aviatrix-role-cloudwatch"
	region              = "us-east-1"
	excluded_gateways   = ["a", "b"]
}
`
}

func testAccCheckCloudwatchAgentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("cloudwatch agent not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetCloudwatchAgentStatus()
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("cloudwatch agent not found")
		}

		return nil
	}
}

func testAccCheckCloudwatchAgentExcludedGatewaysMatch(input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetCloudwatchAgentStatus()
		if !goaviatrix.Equivalent(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}
		return nil
	}
}

func testAccCheckCloudwatchAgentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_cloudwatch_agent" {
			continue
		}

		_, err := client.GetCloudwatchAgentStatus()
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("cloudwatch_agent still exists")
		}
	}

	return nil
}
