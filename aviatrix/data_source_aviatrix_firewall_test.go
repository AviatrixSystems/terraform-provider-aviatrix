package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAviatrixFirewall_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DATA_FIREWALL")
	if skipAcc == "yes" {
		t.Skip("Skipping data source firewall tests as 'SKIP_DATA_FIREWALL' is set")
	}

	msg := ". Set 'SKIP_DATA_FIREWALL' to 'yes' to skip data source firewall tests"
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_firewall.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msg)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixFirewallConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixFirewall(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", "test-gw-"+rName),
					resource.TestCheckResourceAttr(resourceName, "base_policy", "allow-all"),
					resource.TestCheckResourceAttr(resourceName, "base_log_enabled", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixFirewallConfigBasic(rName string) string {
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
	gw_name          = aviatrix_gateway.test.gw_name
	base_policy      = "allow-all"
	base_log_enabled = true

	policy {
		protocol    = "tcp"
		src_ip      = "10.15.0.224/32"
		log_enabled = false
		dst_ip      = "10.12.0.172/32"
		action      = "allow"
		port        = "0:65535"
		description = "This is policy no.1"
	}

	policy {
		protocol    = "tcp"
		src_ip      = "10.15.1.224/32"
		log_enabled = true
		dst_ip      = "10.12.1.172/32"
		action      = "deny"
		port        = "0:65535"
		description = "This is policy no.2"
	}
}
data "aviatrix_firewall" "test" {
	gw_name = aviatrix_firewall.test.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_REGION"), rName, os.Getenv("AWS_REGION"))
}

func testAccDataSourceAviatrixFirewall(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
