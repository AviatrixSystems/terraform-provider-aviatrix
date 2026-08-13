package aviatrix

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAwsTgwVpnConn_basic(t *testing.T) {
	var awsTgwVpnConn goaviatrix.AwsTgwVpnConn

	rName := acctest.RandString(5)
	resourceName := "aviatrix_aws_tgw_vpn_conn.test"
	importStateVerifyIgnore := []string{"vpn_tunnel_data"}

	skipAcc := os.Getenv("SKIP_AWS_TGW_VPN_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping AWS TGW VPN CONN test as SKIP_AWS_TGW_VPN_CONN is set")
	}

	msg := ". Set SKIP_AWS_TGW_VPN_CONN to yes to skip AWS TGW VPN CONN tests"

	awsSideAsNumber := "12"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsTgwVpnConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsTgwVpnConnConfigBasic(rName, awsSideAsNumber),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAwsTgwVpnConnExists(resourceName, &awsTgwVpnConn),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", fmt.Sprintf("tft-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "route_domain_name", "Default_Domain"),
					resource.TestCheckResourceAttr(resourceName, "connection_name", fmt.Sprintf("tfc-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "public_ip", "40.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "remote_as_number", awsSideAsNumber),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
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
	aws_iam	           = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_tgw" "test_aws_tgw" {
	account_name       = aviatrix_account.test_account.account_name
	aws_side_as_number = "64512"
	region             = "%s"
	tgw_name           = "tft-%s"
}
resource "aviatrix_aws_tgw_network_domain" "Default_Domain" {
	name     = "Default_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Shared_Service_Domain" {
	name     = "Shared_Service_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Aviatrix_Edge_Domain" {
	name     = "Aviatrix_Edge_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn1" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Default_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn2" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn3" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Default_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_vpn_conn" "test" {
	tgw_name          = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	route_domain_name = aviatrix_aws_tgw_network_domain.Default_Domain.name
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

		client := mustClient(testAccProvider.Meta())

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
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_vpn_conn" {
			continue
		}

		foundAwsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
			TgwName: rs.Primary.Attributes["tgw_name"],
			VpnID:   rs.Primary.Attributes["vpn_id"],
		}

		_, err := client.GetAwsTgwVpnConn(foundAwsTgwVpnConn)
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("aviatrix AWS TGW VPN CONN still exists")
		}
	}

	return nil
}
