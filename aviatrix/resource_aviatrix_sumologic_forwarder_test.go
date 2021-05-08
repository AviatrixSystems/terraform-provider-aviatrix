package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixSumologicForwarder_basic(t *testing.T) {
	if os.Getenv("SKIP_SUMOLOGIC_FORWARDER") == "yes" {
		t.Skip("Skipping sumologic forwarder test as SKIP_SUMOLOGIC_FORWARDER is set")
	}

	resourceName := "aviatrix_sumologic_forwarder.test_sumologic_forwarder"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSumologicForwarderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicForwarderBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSumologicForwarderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "access_id", "test_id"),
					resource.TestCheckResourceAttr(resourceName, "access_key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "source_category", "test_category"),
					resource.TestCheckResourceAttr(resourceName, "custom_configuration", "key=value"),
					testAccCheckSumologicForwarderExcludedGatewaysMatch([]string{"a", "b"}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_key"},
			},
		},
	})
}

func testAccSumologicForwarderBasic() string {
	return `
resource "aviatrix_sumologic_forwarder" "test_sumologic_forwarder" {
	access_id            = "test_id"
	access_key           = "test_key"
	source_category      = "test_category"
	custom_configuration = "key=value"
	excluded_gateways    = ["a", "b"]
}
`
}

func testAccCheckSumologicForwarderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("sumologic forwarder not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetSumologicForwarderStatus()
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("sumologic forwarder not found")
		}

		return nil
	}
}

func testAccCheckSumologicForwarderExcludedGatewaysMatch(input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetSumologicForwarderStatus()
		if !goaviatrix.Equivalent(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}
		return nil
	}
}

func testAccCheckSumologicForwarderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_sumologic_forwarder" {
			continue
		}

		_, err := client.GetSumologicForwarderStatus()
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("sumologic_forwarder still exists")
		}
	}

	return nil
}
