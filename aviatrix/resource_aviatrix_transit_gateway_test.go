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

func TestAccAviatrixTransitGateway_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipGw := os.Getenv("SKIP_TRANSIT_GATEWAY")
	if skipGw == "yes" {
		t.Skip("Skipping Transit gateway test as SKIP_TRANSIT_GATEWAY is set")
	}

	skipGwAws := os.Getenv("SKIP_TRANSIT_GATEWAY_AWS")
	skipGwArm := os.Getenv("SKIP_TRANSIT_GATEWAY_ARM")
	if skipGwAws == "yes" && skipGwArm == "yes" {
		t.Skip("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_AWS and SKIP_TRANSIT_GATEWAY_ARM are both set")
	}

	if skipGwAws != "yes" {
		resourceNameAws := "aviatrix_transit_gateway.test_transit_gateway_aws"
		msgCommonAws := ". Set SKIP_TRANSIT_GATEWAY_AWS to yes to skip Transit Gateway tests in aws"

		preGatewayCheck(t, msgCommonAws)

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicAws(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameAws, &gateway),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_size", "t2.micro"),
						resource.TestCheckResourceAttr(resourceNameAws, "account_name", fmt.Sprintf("tfa-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAws, "vpc_id", os.Getenv("AWS_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameAws, "subnet", os.Getenv("AWS_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameAws, "vpc_reg", os.Getenv("AWS_REGION")),
					),
				},
				{
					ResourceName:      resourceNameAws,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_AWS is set")
	}

	if skipGwArm != "yes" {
		resourceNameArm := "aviatrix_transit_gateway.test_transit_gateway_arm"

		msgCommonArm := ". Set SKIP_TRANSIT_GATEWAY_ARM to yes to skip Transit Gateway tests in ARM"
		preGatewayCheckArm(t, msgCommonArm)

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicArm(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameArm, &gateway),
						resource.TestCheckResourceAttr(resourceNameArm, "gw_name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameArm, "gw_size", os.Getenv("ARM_GW_SIZE")),
						resource.TestCheckResourceAttr(resourceNameArm, "account_name", fmt.Sprintf("tfaz-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameArm, "vpc_id", os.Getenv("ARM_VNET_ID")),
						resource.TestCheckResourceAttr(resourceNameArm, "subnet", os.Getenv("ARM_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameArm, "vpc_reg", os.Getenv("ARM_REGION")),
					),
				},
				{
					ResourceName:      resourceNameArm,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_ARM is set")
	}
}

func testAccTransitGatewayConfigBasicAws(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_aws" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test_transit_gateway_aws" {
	cloud_type   = 1
	account_name = aviatrix_account.test_aws.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_VPC_NET"))
}

func testAccTransitGatewayConfigBasicArm(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_arm" {
	account_name        = "tfaz-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}

resource "aviatrix_transit_gateway" "test_transit_gateway_arm" {
	cloud_type   = 8
	account_name = aviatrix_account.test_arm.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"),
		os.Getenv("ARM_APPLICATION_KEY"), os.Getenv("ARM_VNET_ID"), os.Getenv("ARM_REGION"),
		os.Getenv("ARM_GW_SIZE"), os.Getenv("ARM_SUBNET"))
}

func testAccCheckTransitGatewayExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transit gateway Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit gateway ID is set")
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
			return fmt.Errorf("transit gateway not found")
		}

		*gateway = *foundGateway
		return nil
	}
}

func testAccCheckTransitGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_vpc" {
			continue
		}

		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetGateway(foundGateway)
		if err == nil {
			return fmt.Errorf("transit gateway still exists")
		}
	}

	return nil
}
