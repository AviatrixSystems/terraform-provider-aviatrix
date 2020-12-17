package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixFilebeatForwarder_basic(t *testing.T) {
	if os.Getenv("SKIP_FILEBEAT_FORWARDER") == "yes" {
		t.Skip("Skipping filebeat forwarder test as SKIP_FILEBEAT_FORWARDER is set")
	}

	resourceName := "aviatrix_filebeat_forwarder.test_filebeat_forwarder"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFilebeatForwarderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFilebeatForwarderBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFilebeatForwarderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "server", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "port", "10"),
					testAccCheckFilebeatForwarderExcludedGatewaysMatch([]string{"a", "b"}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"trusted_ca_file", "config_file"},
			},
		},
	})
}

func testAccFilebeatForwarderBasic() string {
	return `
resource "aviatrix_filebeat_forwarder" "test_filebeat_forwarder" {
	server            = "1.2.3.4"
	port              = 10
	excluded_gateways = ["a", "b"]
}
`
}

func testAccCheckFilebeatForwarderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("filebeat forwarder not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetFilebeatForwarderStatus()
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("filebeat forwarder not found")
		}

		return nil
	}
}

func testAccCheckFilebeatForwarderExcludedGatewaysMatch(input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetFilebeatForwarderStatus()
		if !goaviatrix.Equivalent(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}
		return nil
	}
}

func testAccCheckFilebeatForwarderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_filebeat_forwarder" {
			continue
		}

		_, err := client.GetFilebeatForwarderStatus()
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("filebeat_forwarder still exists")
		}
	}

	return nil
}
