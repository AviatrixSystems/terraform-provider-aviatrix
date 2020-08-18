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

func TestAccAviatrixFirewallPolicy_basic(t *testing.T) {
	if os.Getenv("SKIP_FIREWALL_POLICY") == "yes" {
		t.Skip("Skipping firewall policy test as SKIP_FIREWALL_POLICY is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_firewall_policy.test_firewall_policy"

	msg := ". Set SKIP_FIREWALL_POLICY to yes to skip firewall policy tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallPolicyBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallPolicyExists(resourceName),
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

func testAccFirewallPolicyBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	name         = "tfv-%s"
	region       = "%s"
	cidr         = "10.0.0.0/16"
}
data "aviatrix_vpc" "test" {
	name = aviatrix_vpc.test.name
}
resource "aviatrix_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "test-gw-%s"
	vpc_id       = aviatrix_vpc.test.vpc_id
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = data.aviatrix_vpc.test.public_subnets[0].cidr
}
resource "aviatrix_firewall" "test" {
	gw_name                  = aviatrix_gateway.test.gw_name
	base_policy              = "allow-all"
	base_log_enabled         = true
	manage_firewall_policies = false
}

resource "aviatrix_firewall_policy" "test_firewall_policy" {
	gw_name     = aviatrix_firewall.test.gw_name
	protocol    = "tcp"
	src_ip      = "10.15.0.224/32"
	log_enabled = true
	dst_ip      = "10.12.0.172/32"
	action      = "allow"
	port        = "0:65535"
	description = "This is policy no.1"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_REGION"), rName, os.Getenv("AWS_REGION"))
}

func testAccCheckFirewallPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("firewall_policy Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no firewall_policy ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		logEnabled := "on"
		if rs.Primary.Attributes["log_enabled"] == "false" {
			logEnabled = "off"
		}
		fw := &goaviatrix.Firewall{
			GwName: rs.Primary.Attributes["gw_name"],
			PolicyList: []*goaviatrix.Policy{
				{
					Protocol:    rs.Primary.Attributes["protocol"],
					Port:        rs.Primary.Attributes["port"],
					SrcIP:       rs.Primary.Attributes["src_ip"],
					DstIP:       rs.Primary.Attributes["dst_ip"],
					LogEnabled:  logEnabled,
					Action:      rs.Primary.Attributes["action"],
					Description: rs.Primary.Attributes["description"],
				},
			},
		}

		_, err := client.GetFirewallPolicy(fw)
		if err != nil {
			return err
		}
		if getFirewallPolicyID(fw) != rs.Primary.ID {
			return fmt.Errorf("firewall_policy not found")
		}

		return nil
	}
}

func testAccCheckFirewallPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_firewall_policy" {
			continue
		}
		logEnabled := "on"
		if rs.Primary.Attributes["log_enabled"] == "false" {
			logEnabled = "off"
		}
		fw := &goaviatrix.Firewall{
			GwName: rs.Primary.Attributes["gw_name"],
			PolicyList: []*goaviatrix.Policy{
				{
					Protocol:    rs.Primary.Attributes["protocol"],
					Port:        rs.Primary.Attributes["port"],
					SrcIP:       rs.Primary.Attributes["src_ip"],
					DstIP:       rs.Primary.Attributes["dst_ip"],
					LogEnabled:  logEnabled,
					Action:      rs.Primary.Attributes["action"],
					Description: rs.Primary.Attributes["description"],
				},
			},
		}
		_, err := client.GetFirewallPolicy(fw)
		if err == nil {
			return fmt.Errorf("firewall_policy still exists")
		}
	}

	return nil
}
