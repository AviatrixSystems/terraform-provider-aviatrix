package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAwsTgwVpnConn_basic(t *testing.T) {
	var awsTgwVpnConn goaviatrix.AwsTgwVpnConn

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_aws_tgw_vpn_conn.test"

	skipAcc := os.Getenv("SKIP_AWS_TGW_VPN_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping AWS TGW VPN CONN test as SKIP_AWS_TGW_VPN_CONN is set")
	}

	msg := ". Set SKIP_AWS_TGW_VPN_CONN to yes to skip AWS TGW VPN CONN tests"

	preAccountCheck(t, msg)

	awsSideAsNumber := "12"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsTgwVpnConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsTgwVpnConnConfigBasic(rName, awsSideAsNumber),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAwsTgwVpnConnExists(resourceName, &awsTgwVpnConn),
					resource.TestCheckResourceAttr(
						resourceName, "tgw_name", fmt.Sprintf("tft-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "route_domain_name", "Default_Domain"),
					resource.TestCheckResourceAttr(
						resourceName, "connection_name", fmt.Sprintf("tfc-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "public_ip", "40.0.0.0"),
					resource.TestCheckResourceAttr(
						resourceName, "remote_as_number", awsSideAsNumber),
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

func testAccAwsTgwVpnConnConfigBasic(rName string, awsSideAsNumber string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam		         = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_tgw" "test_aws_tgw" {
	account_name          = aviatrix_account.test_account.account_name
	aws_side_as_number    = "64512"
	manage_vpc_attachment = false
	region                = "%s"
	tgw_name              = "tft-%s"
	security_domains {
		connected_domains    = [
			"Default_Domain",
			"Shared_Service_Domain"
		]
		security_domain_name = "Aviatrix_Edge_Domain"
	}
	security_domains {
		connected_domains    = [
			"Aviatrix_Edge_Domain",
			"Shared_Service_Domain"
		]
		security_domain_name = "Default_Domain"
	}
	security_domains {
		connected_domains    = [
			"Aviatrix_Edge_Domain",
			"Default_Domain"
		]
		security_domain_name = "Shared_Service_Domain"
	}
}
resource "aviatrix_aws_tgw_vpn_conn" "test" {
	tgw_name          = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	route_domain_name = "Default_Domain"
	connection_name   = "tfc-%s"
	public_ip         = "40.0.0.0"
	remote_as_number  = "%s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName, rName, awsSideAsNumber)
}

func tesAccCheckAwsTgwVpnConnExists(n string, awsTgwVpnConn *goaviatrix.AwsTgwVpnConn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("AWS TGW VPN CONN Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no AWS TGW VPN CONN ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAwsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
			TgwName: rs.Primary.Attributes["tgw_name"],
			VpnID:   rs.Primary.Attributes["vpn_id"],
		}

		foundAwsTgwVpnConn2, err := client.GetAwsTgwVpnConn(foundAwsTgwVpnConn)
		if err != nil {
			return err
		}

		if foundAwsTgwVpnConn2.TgwName != rs.Primary.Attributes["tgw_name"] {
			return fmt.Errorf("tgw_name Not found in created attributes")
		}

		if foundAwsTgwVpnConn2.ConnName != rs.Primary.Attributes["connection_name"] {
			return fmt.Errorf("connection_name Not found in created attributes")
		}

		*awsTgwVpnConn = *foundAwsTgwVpnConn

		return nil
	}
}

func testAccCheckAwsTgwVpnConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_vpn_conn" {
			continue
		}

		foundAwsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
			TgwName: rs.Primary.Attributes["tgw_name"],
			VpnID:   rs.Primary.Attributes["vpn_id"],
		}

		_, err := client.GetAwsTgwVpnConn(foundAwsTgwVpnConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("aviatrix AWS TGW VPN CONN still exists")
		}
	}

	return nil
}
