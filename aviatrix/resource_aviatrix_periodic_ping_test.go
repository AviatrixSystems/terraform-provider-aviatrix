package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixPeriodicPing_basic(t *testing.T) {
	if os.Getenv("SKIP_PERIODIC_PING") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_PERIODIC_PING is set")
	}

	resourceName := "aviatrix_periodic_ping.test_periodic_ping"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			periodicPingPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPeriodicPingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPeriodicPingBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPeriodicPingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "127.0.0.1"),
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

func testAccPeriodicPingBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_periodic_ping" "test_periodic_ping" {
	gw_name    = "%s"
	interval   = 5
	ip_address = "127.0.0.1"
}
`, os.Getenv("GATEWAY_NAME"))
}

func testAccCheckPeriodicPingExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("periodic_ping Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no periodic_ping ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundPeriodicPing := &goaviatrix.PeriodicPing{
			GwName: rs.Primary.Attributes["gw_name"],
		}

		_, err := client.GetPeriodicPing(foundPeriodicPing)
		if err != nil {
			return err
		}
		if foundPeriodicPing.GwName != rs.Primary.ID {
			return fmt.Errorf("periodic_ping not found")
		}

		return nil
	}
}

func testAccCheckPeriodicPingDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_periodic_ping" {
			continue
		}
		foundPeriodicPing := &goaviatrix.PeriodicPing{
			GwName: rs.Primary.Attributes["gw_name"],
		}
		_, err := client.GetPeriodicPing(foundPeriodicPing)
		if err == nil {
			return fmt.Errorf("periodic_ping still exists")
		}
	}

	return nil
}

func periodicPingPreCheck(t *testing.T) {
	if os.Getenv("GATEWAY_NAME") == "" {
		t.Fatal("GATEWAY_NAME must be set for aviatrix_periodic_ping acceptance test.")
	}
}
