package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAviatrixGateway_basic(t *testing.T) {
	var gateway goaviatrix.Gateway
	rName := fmt.Sprintf("tf-testing-%s", acctest.RandString(5))
	resourceName := "aviatrix_gateway.test"
	accountID := os.Getenv("AWS_ACCOUNT_NUMBER")

	vpcID := os.Getenv("AWS_VPC_ID")
	if vpcID == "" {
		t.Skip("Environment variable AWS_VPC_ID is not set")
	}

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		t.Skip("Environment variable AWS_DEFAULT_REGION is not set")
	}

	vpcNet := os.Getenv("AWS_VPC_NET")
	if vpcNet == "" {
		t.Skip("Environment variable AWS_VPC_NET is not set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfig_basic(rName, accountID, vpcID, region, vpcNet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(
						resourceName, "gw_name", rName),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_size", "t2.micro"),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_net", vpcNet),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_reg", region),
				),
			},
		},
	})
}

func testAccGatewayConfig_basic(rName string, accountID string, vpcID string, region string, vpcNet string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name = "%[1]s"
	account_email = "noone@aviatrix.com"
	cloud_type = 1
	aws_account_number = "%[2]s"
	aws_iam = "true"
	aws_role_app = "arn:aws:iam::%[2]s:role/aviatrix-role-app"
	aws_role_ec2 = "arn:aws:iam::%[2]s:role/aviatrix-role-ec2"

	// account_email currently doesn't behave according to API docs, ignoring changes:
	lifecycle {
		ignore_changes = ["account_email"]
	}
}

resource "aviatrix_gateway" "test" {
	cloud_type = 1
	account_name = "${aviatrix_account.test.account_name}"
	gw_name = "%[1]s"
	vpc_id = "%[3]s"
	vpc_reg = "%[4]s"
	vpc_size = "t2.micro"
	vpc_net = "%[5]s"
}
	`, rName, accountID, vpcID, region, vpcNet)
}

func testAccCheckGatewayExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Gateway Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Account ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetGateway(foundGateway)

		if err != nil {
			return err
		}

		if foundGateway.GwName != rs.Primary.ID {
			return fmt.Errorf("Gateway not found")
		}

		*gateway = *foundGateway

		return nil
	}
}

func testAccCheckGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_gateway" {
			continue
		}
		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}
		_, err := client.GetGateway(foundGateway)

		if err == nil {
			return fmt.Errorf("Gateway still exists")
		}
	}
	return nil
}
