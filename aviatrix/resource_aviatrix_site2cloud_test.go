package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixS2C_basic(t *testing.T) {
	var s2c goaviatrix.Site2Cloud

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_site2cloud.foo"

	skipAcc := os.Getenv("SKIP_S2C")
	if skipAcc == "yes" {
		t.Skip("Skipping Site2Cloud test as SKIP_S2C is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set SKIP_S2C to yes to skip Site2Cloud tests")

		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS2CDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS2CConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS2CExists("aviatrix_site2cloud.foo", &s2c),
					resource.TestCheckResourceAttr(resourceName, "connection_name", fmt.Sprintf("tfs-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "tunnel_type", "policy"),
					resource.TestCheckResourceAttr(resourceName, "primary_cloud_gateway_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "remote_subnet_cidr", "10.23.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_type", "generic"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "unmapped"),
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

func testAccS2CConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
resource "aviatrix_site2cloud" "foo" {
	vpc_id                     = aviatrix_gateway.test.vpc_id
	connection_name            = "tfs-%[1]s"
	connection_type            = "unmapped"
	remote_gateway_type        = "generic"
	tunnel_type                = "policy"
	primary_cloud_gateway_name = aviatrix_gateway.test.gw_name
	remote_gateway_ip          = "8.8.8.8"
	remote_subnet_cidr         = "10.23.0.0/24"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccCheckS2CExists(n string, s2c *goaviatrix.Site2Cloud) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("site2cloud Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no site2cloud ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundS2C := &goaviatrix.Site2Cloud{
			TunnelName: rs.Primary.Attributes["connection_name"],
			VpcID:      rs.Primary.Attributes["vpc_id"],
		}

		_, err := client.GetSite2Cloud(foundS2C)
		if err != nil {
			return err
		}
		if foundS2C.TunnelName+"~"+foundS2C.VpcID != rs.Primary.ID {
			return fmt.Errorf("site2cloud connection not found")
		}

		*s2c = *foundS2C
		return nil
	}
}

func testAccCheckS2CDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_site2cloud" {
			continue
		}

		foundS2C := &goaviatrix.Site2Cloud{
			TunnelName: rs.Primary.Attributes["connection_name"],
			VpcID:      rs.Primary.Attributes["vpc_id"],
		}

		_, err := client.GetSite2Cloud(foundS2C)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("site2cloud still exists")
		}
	}

	return nil
}
