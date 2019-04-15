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

func preGatewayCheck(t *testing.T, msgCommon string) (string, string, string) {
	preAccountCheck(t, msgCommon)

	vpcID := os.Getenv("AWS_VPC_ID")
	if vpcID == "" {
		t.Fatal("Environment variable AWS_VPC_ID is not set" + msgCommon)
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		t.Fatal("Environment variable AWS_REGION is not set" + msgCommon)
	}

	vpcNet := os.Getenv("AWS_VPC_NET")
	if vpcNet == "" {
		t.Fatal("Environment variable AWS_VPC_NET is not set" + msgCommon)
	}
	return vpcID, region, vpcNet
}

func TestAccAviatrixGateway_basic(t *testing.T) {
	var gateway goaviatrix.Gateway
	rName := fmt.Sprintf("tf-testing-%s", acctest.RandString(5))
	resourceName := "aviatrix_gateway.test"

	msgCommon := ". Set SKIP_GATEWAY to yes to skip Gateway tests"

	skipGw := os.Getenv("SKIP_GATEWAY")
	if skipGw == "yes" {
		t.Skip("Skipping Gateway test as SKIP_GATEWAY is set")
	}

	preAccountCheck(t, msgCommon)
	accountID := os.Getenv("AWS_ACCOUNT_NUMBER")

	vpcID, region, vpcNet := preGatewayCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayConfigBasic(rName, accountID, vpcID, region, vpcNet),
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

func testAccGatewayConfigBasic(rName string, accountID string, vpcID string, region string, vpcNet string) string {
	return fmt.Sprintf(`

resource "aviatrix_account" "test" {
	account_name = "%s"
	cloud_type = 1
	aws_account_number = "%s"
	aws_iam = "false"
	aws_access_key = "%s"
	aws_secret_key = "%s"
}

resource "aviatrix_gateway" "test" {
	cloud_type = 1
	account_name = "${aviatrix_account.test.account_name}"
	gw_name = "%[1]s"
	vpc_id = "%[5]s"
	vpc_reg = "%[6]s"
	vpc_size = "t2.micro"
	vpc_net = "%[7]s"
}
	`, rName, accountID, os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), vpcID, region, vpcNet)
}

func testAccCheckGatewayExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("gateway Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Account ID is set")
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
			return fmt.Errorf("gateway not found")
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
			return fmt.Errorf("gateway still exists")
		}
	}
	return nil
}
