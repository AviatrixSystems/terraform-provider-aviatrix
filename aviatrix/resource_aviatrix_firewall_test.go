package aviatrix

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixFirewall_basic(t *testing.T) {
	var firewall goaviatrix.Firewall

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_firewall.test_firewall"

	skipAcc := os.Getenv("SKIP_FIREWALL")
	if skipAcc == "yes" {
		t.Skip("Skipping Firewall test as SKIP_FIREWALL is set")
	}
	msg := ". Set SKIP_FIREWALL to yes to skip firewall tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists("aviatrix_firewall.test_firewall", &firewall),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "base_policy", "allow-all"),
					resource.TestCheckResourceAttr(resourceName, "base_log_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "policy.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.src_ip", "10.15.0.224/32"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.log_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.dst_ip", "10.12.0.172/32"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.port", "0:65535"),
					resource.TestCheckResourceAttr(resourceName, "policy.1.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "policy.1.src_ip", fmt.Sprintf("tft-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "policy.1.log_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "policy.1.dst_ip", "10.12.1.172/32"),
					resource.TestCheckResourceAttr(resourceName, "policy.1.action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "policy.1.port", "0:65535"),
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

func testAccFirewallConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "test_gw" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_firewall_tag" "foo" {
	firewall_tag = "tft-%s"
	cidr_list {
		cidr_tag_name = "a1"
		cidr          = "10.1.0.0/24"
	}
	cidr_list {
		cidr_tag_name = "b1"
		cidr          = "10.2.0.0/24"
	}
}
resource "aviatrix_firewall" "test_firewall" {
	gw_name          = aviatrix_gateway.test_gw.gw_name
	base_policy      = "allow-all"
	base_log_enabled = false
	policy {
		protocol    = "tcp"
		src_ip      = "10.15.0.224/32"
		log_enabled = false
		dst_ip      = "10.12.0.172/32"
		action      = "deny"
		port        = "0:65535"
	}
	policy {
		protocol    = "tcp"
		src_ip      = aviatrix_firewall_tag.foo.firewall_tag
		log_enabled = false
		dst_ip      = "10.12.1.172/32"
		action      = "deny"
		port        = "0:65535"
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func testAccCheckFirewallExists(n string, firewall *goaviatrix.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("firewall Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no firewall ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundFirewall := &goaviatrix.Firewall{
			GwName: rs.Primary.Attributes["gw_name"],
		}

		_, err := client.GetPolicy(foundFirewall)
		if err != nil {
			return err
		}
		if foundFirewall.GwName != rs.Primary.ID {
			return fmt.Errorf("firewall not found")
		}

		*firewall = *foundFirewall
		return nil
	}
}

func testAccCheckFirewallDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_firewall" {
			continue
		}

		foundFirewall := &goaviatrix.Firewall{
			GwName: rs.Primary.Attributes["gw_name"],
		}

		_, err := client.GetPolicy(foundFirewall)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("firewall still exists")
		}
	}

	return nil
}

func testResourceFirewallStateDataV0() map[string]interface{} {
	return map[string]interface{}{}
}

func testResourceFirewallStateDataV1() map[string]interface{} {
	return map[string]interface{}{
		"manage_firewall_policies": true,
	}
}

func TestResourceFirewallStateUpgradeV0(t *testing.T) {
	expected := testResourceFirewallStateDataV1()
	actual, err := resourceAviatrixFirewallStateUpgradeV0(testResourceFirewallStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\nexpected:%#v\ngot:%#v\n", expected, actual)
	}
}

func testResourceFirewallStateDataV0ManageAlreadySet() map[string]interface{} {
	return map[string]interface{}{
		"manage_firewall_policies": false,
	}
}

func testResourceFirewallStateDataV1ManageAlreadySet() map[string]interface{} {
	return map[string]interface{}{
		"manage_firewall_policies": false,
	}
}

func TestResourceFirewallStateUpgradeV0ManageAlreadySet(t *testing.T) {
	expected := testResourceFirewallStateDataV1ManageAlreadySet()
	actual, err := resourceAviatrixFirewallStateUpgradeV0(testResourceFirewallStateDataV0ManageAlreadySet(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\nexpected:%#v\ngot:%#v\n", expected, actual)
	}
}
