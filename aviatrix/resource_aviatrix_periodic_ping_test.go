package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

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
			preGatewayCheck(t, ". Set SKIP_PERIODIC_PING to yes to skip Periodic Ping testing.")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPeriodicPingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPeriodicPingBasic(acctest.RandString(5)),
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

func testAccPeriodicPingBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%[1]s"
	cloud_type         = 1
	aws_account_number = "%[2]s"
	aws_iam            = false
	aws_access_key     = "%[3]s"
	aws_secret_key     = "%[4]s"
}

resource "aviatrix_gateway" "test_gw" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}

resource "aviatrix_periodic_ping" "test_periodic_ping" {
	gw_name    = aviatrix_gateway.test_gw.gw_name
	interval   = 5
	ip_address = "127.0.0.1"
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
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
