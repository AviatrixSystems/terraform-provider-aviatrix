package aviatrix

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixRemoteSyslog_basic(t *testing.T) {
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
					testAccCheckExcludedGatewaysMatch(rIndex, []string{"a", "b"}),
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

		resp, _ := client.GetRemoteSyslogStatus(index)
		if resp.Status != "enabled" {
			return fmt.Errorf("remote syslog %d not found", index)
		}

		return nil
	}
}

func testAccCheckExcludedGatewaysMatch(index int, input []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, _ := client.GetRemoteSyslogStatus(index)
		if !sliceMatch(resp.ExcludedGateways, input) {
			return fmt.Errorf("excluded gateways don't match with the input")
		}
		return nil
	}
}

func sliceMatch(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func testAccCheckRemoteSyslogDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_remote_syslog" {
			continue
		}
		idx, _ := strconv.Atoi(rs.Primary.Attributes["index"])

		resp, _ := client.GetRemoteSyslogStatus(idx)
		if resp.Status == "enabled" {
			return fmt.Errorf("remote_syslog still exists")
		}
	}

	return nil
}
