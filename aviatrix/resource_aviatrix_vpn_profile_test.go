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

func TestAccAviatrixVPNProfile_basic(t *testing.T) {
	var vpnProfile goaviatrix.Profile
	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_vpn_profile.test_vpn_profile"

	skipAcc := os.Getenv("SKIP_VPN_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping VPN Profile test as SKIP_VPN_PROFILE is set")
	}
	msg := ". Set SKIP_VPN_PROFILE to yes to skip VPN Profile tests"

	preGatewayCheck(t, msg)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPNProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNProfileConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNProfileExists("aviatrix_vpn_profile.test_vpn_profile", &vpnProfile),
					resource.TestCheckResourceAttr(
						resourceName, "name", fmt.Sprintf("tfp-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "base_rule", "allow_all"),
					resource.TestCheckResourceAttr(
						resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "users.0", fmt.Sprintf("tfu-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "policy.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "policy.0.action", "deny"),
					resource.TestCheckResourceAttr(
						resourceName, "policy.0.proto", "tcp"),
					resource.TestCheckResourceAttr(
						resourceName, "policy.0.port", "443"),
					resource.TestCheckResourceAttr(
						resourceName, "policy.0.target", "10.0.0.0/32"),
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

func testAccVPNProfileConfigBasic(rName string) string {
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

resource "aviatrix_vpn_profile" "test_vpn_profile" {
    name = "tfp-%s"
    base_rule = "allow_all"
    users = ["${aviatrix_vpn_user.test_vpn_user.user_name}"]
    policy = [
    {
        action = "deny"
        proto = "tcp"
        port = "443"
        target = "10.0.0.0/32"
    }
    ]
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_VPC_NET"), rName, rName,
		rName)
}

func testAccCheckVPNProfileExists(n string, vpnProfile *goaviatrix.Profile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN Profile Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no VPN Profile ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundVPNProfile := &goaviatrix.Profile{
			Name: rs.Primary.Attributes["name"],
		}

		foundVPNProfile2, err := client.GetProfile(foundVPNProfile)
		if err != nil {
			return err
		}
		if foundVPNProfile2.Name != rs.Primary.ID {
			return fmt.Errorf("VPN profile not found")
		}
		*vpnProfile = *foundVPNProfile

		return nil
	}
}

func testAccCheckVPNProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpn_profile" {
			continue
		}
		foundVPNProfile := &goaviatrix.Profile{
			Name: rs.Primary.Attributes["name"],
		}
		_, err := client.GetProfile(foundVPNProfile)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("VPN Profile still exists")
		}
	}
	return nil
}
