package aviatrix

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixRemoteSyslog_basic(t *testing.T) {
	if os.Getenv("SKIP_REMOTE_SYSLOG") == "yes" {
		t.Skip("Skipping remote syslog test as SKIP_REMOTE_SYSLOG is set")
	}

	rIndex := acctest.RandIntRange(0, 9)
	resourceName := "aviatrix_remote_syslog.test_remote_syslog"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRemoteSyslogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRemoteSyslogBasic(rIndex),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemoteSyslogExists(resourceName, rIndex),
					resource.TestCheckResourceAttr(resourceName, "index", strconv.Itoa(rIndex)),
					resource.TestCheckResourceAttr(resourceName, "server", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "port", "10"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					testAccCheckRemoteSyslogExcludedGatewaysMatch(rIndex, []string{"a", "b"}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ca_certificate_file", "public_certificate_file", "private_key_file"},
			},
		},
	})
}

func testAccRemoteSyslogBasic(rName int) string {
	return fmt.Sprintf(`
resource "aviatrix_remote_syslog" "test_remote_syslog" {
	index             = %d
	server            = "1.2.3.4"
	port              = 10
	protocol          = "TCP"
	excluded_gateways = ["a", "b"]
}
`, rName)
}

func testAccCheckRemoteSyslogExists(resourceName string, index int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("remote syslog not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetRemoteSyslogStatus(index)
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("remote syslog %d not found", index)
		}

		return nil
	}
}

func testAccCheckRemoteSyslogExcludedGatewaysMatch(index int, input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetRemoteSyslogStatus(index)
		if !goaviatrix.Equivalent(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}
		return nil
	}
}

func testAccCheckRemoteSyslogDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_remote_syslog" {
			continue
		}
		idx, _ := strconv.Atoi(rs.Primary.Attributes["index"])

		_, err := client.GetRemoteSyslogStatus(idx)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("remote_syslog still exists")
		}
	}

	return nil
}
