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

func TestAccAviatrixVPNUser_basic(t *testing.T) {
	var vpnUser goaviatrix.VPNUser
	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_vpn_user.test_vpn_user"

	skipAcc := os.Getenv("SKIP_VPN_USER")
	if skipAcc == "yes" {
		t.Skip("Skipping VPN User test as SKIP_VPN_USER is set")
	}
	msg := ". Set SKIP_VPN_USER to yes to skip VPN User tests"

	preGatewayCheck(t, msg)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPNUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNUserConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNUserExists("aviatrix_vpn_user.test_vpn_user", &vpnUser),
					resource.TestCheckResourceAttr(
						resourceName, "gw_name", fmt.Sprintf("tfl-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(
						resourceName, "user_email", "user@xyz.com"),
					resource.TestCheckResourceAttr(
						resourceName, "user_name", fmt.Sprintf("tfu-%s", rName)),
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

func testAccVPNUserConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
    account_name = "tfa-%s"
    cloud_type = 1
    aws_account_number = "%s"
    aws_iam = "false"
    aws_access_key = "%s"
    aws_secret_key = "%s"
}

resource "aviatrix_gateway" "test_gw" {
	cloud_type = 1
	account_name = "${aviatrix_account.test_account.account_name}"
	gw_name = "tfg-%s"
	vpc_id = "%s"
	vpc_reg = "%s"
	vpc_size = "t2.micro"
	vpc_net = "%s"
    vpn_access = "yes"
    vpn_cidr = "192.168.43.0/24" 
	enable_elb = "yes"
	elb_name = "tfl-%s"
}

resource "aviatrix_vpn_user" "test_vpn_user" {
    vpc_id = "${aviatrix_gateway.test_gw.vpc_id}"
    gw_name = "${aviatrix_gateway.test_gw.elb_name}"
    user_name = "tfu-%s"
    user_email = "user@xyz.com"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_VPC_NET"), rName, rName)
}

func testAccCheckVPNUserExists(n string, vpnUser *goaviatrix.VPNUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN User Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no VPN User ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundVPNUser := &goaviatrix.VPNUser{
			UserEmail: rs.Primary.Attributes["user_email"],
			VpcID:     rs.Primary.Attributes["vpc_id"],
			UserName:  rs.Primary.Attributes["user_name"],
			GwName:    rs.Primary.Attributes["gw_name"],
		}

		foundVPNUser2, err := client.GetVPNUser(foundVPNUser)

		if err != nil {
			return err
		}

		if foundVPNUser2.UserName != rs.Primary.ID {
			return fmt.Errorf("VPN user not found")
		}

		*vpnUser = *foundVPNUser

		return nil
	}
}

func testAccCheckVPNUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpn_user" {
			continue
		}
		foundVPNUser := &goaviatrix.VPNUser{
			UserEmail: rs.Primary.Attributes["user_email"],
			VpcID:     rs.Primary.Attributes["vpc_id"],
			UserName:  rs.Primary.Attributes["user_name"],
			GwName:    rs.Primary.Attributes["gw_name"],
		}
		_, err := client.GetVPNUser(foundVPNUser)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("VPN User still exists")
		}
	}
	return nil
}
