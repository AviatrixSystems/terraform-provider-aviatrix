package aviatrix

import (
	"fmt"
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
					resource.TestCheckResourceAttr(resourceName, "index", strconv.Itoa(rIndex)),
					resource.TestCheckResourceAttr(resourceName, "server", "1.2.3.4"),
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
	index    = %d
	server   = "1.2.3.4"
	port     = 10
	protocol = "TCP"
}
`, rName)
}

func testAccCheckRemoteSyslogDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_remote_syslog" {
			continue
		}
		idx, _ := strconv.Atoi(rs.Primary.Attributes["index"])

		_, err := client.GetRemoteSyslogStatus(idx)
		if err == nil {
			return fmt.Errorf("remote_syslog still exists")
		}
	}

	return nil
}
