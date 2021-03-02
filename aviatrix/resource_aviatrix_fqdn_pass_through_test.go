package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixFQDNPassThrough_basic(t *testing.T) {
	if os.Getenv("SKIP_FQDN_PASS_THROUGH") == "yes" {
		t.Skip("Skipping FQDN pass through test as SKIP_FQDN_PASS_THROUGH is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_fqdn_pass_through.test_fqdn_pass_through"
	msg := ". Set SKIP_FQDN_PASS_THROUGH to yes to skip FQDN pass through tests."

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFQDNPassThroughDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFQDNPassThroughBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFQDNPassThroughExists(resourceName),
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

func testAccFQDNPassThroughBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tf-acc-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_vpc" "test_vpc" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	region       = "%[5]s"
	name         = "aws-vpc-%[1]s"
	cidr         = "10.0.10.0/24"
}

data "aviatrix_vpc" "test_vpc" {
	name = aviatrix_vpc.test_vpc.name
}

resource "aviatrix_gateway" "test_gw_aws" {
	cloud_type     = 1
	account_name   = aviatrix_account.test_acc_aws.account_name
	gw_name        = "tfg-aws-%[1]s"
	vpc_id         = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg        = "%[5]s"
	gw_size        = "t2.micro"
	subnet         = data.aviatrix_vpc.test_vpc.public_subnets[0].cidr
	single_ip_snat = true
}

resource "aviatrix_fqdn" "test_fqdn" {
	fqdn_tag     = "tag-%[1]s"
	fqdn_enabled = true
	fqdn_mode    = "white"

	gw_filter_tag_list {
		gw_name        = aviatrix_gateway.test_gw_aws.gw_name
		source_ip_list = [
			"172.31.0.0/16",
			"172.31.0.0/20",
		]
	}

	domain_names {
		fqdn   = "facebook.com"
		proto  = "tcp"
		port   = "443"
		action = "Allow"
	}
}

resource "aviatrix_fqdn_pass_through" "test_fqdn_pass_through" {
	gw_name            = aviatrix_gateway.test_gw_aws.gw_name
	pass_through_cidrs = [
		"10.0.0.0/24",
		"10.0.1.0/24",
		"10.0.2.0/24",
	]

	depends_on         = [aviatrix_fqdn.test_fqdn]
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"))
}

func testAccCheckFQDNPassThroughExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("fqdn_pass_through Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no fqdn_pass_through ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		gw := &goaviatrix.Gateway{GwName: rs.Primary.Attributes["gw_name"]}

		_, err := client.GetFQDNPassThroughCIDRs(gw)
		if err != nil {
			return err
		}
		if gw.GwName != rs.Primary.ID {
			return fmt.Errorf("fqdn_pass_through not found")
		}

		return nil
	}
}

func testAccCheckFQDNPassThroughDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_fqdn_pass_through" {
			continue
		}
		gw := &goaviatrix.Gateway{GwName: rs.Primary.Attributes["gw_name"]}
		_, err := client.GetFQDNPassThroughCIDRs(gw)
		if err == nil {
			return fmt.Errorf("fqdn_pass_through still exists")
		}
	}

	return nil
}
