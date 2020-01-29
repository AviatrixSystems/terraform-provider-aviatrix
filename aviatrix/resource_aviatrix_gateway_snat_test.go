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

func TestAccAviatrixGatewaySNat_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	importStateVerifyIgnore := []string{"vnet_and_resource_group_names"}

	resourceName := "aviatrix_gateway_snat.test"

	msgCommon := ". Set SKIP_GATEWAY_SNAT to yes to skip gateway source NAT tests"

	skipSNat := os.Getenv("SKIP_GATEWAY_SNAT")
	skipSNatAWS := os.Getenv("SKIP_GATEWAY_SNAT_AWS")
	skipSNatARM := os.Getenv("SKIP_GATEWAY_SNAT_ARM")

	if skipSNat == "yes" {
		t.Skip("Skipping gateway source NAT tests as SKIP_GATEWAY_SNAT is set")
	}
	if skipSNatAWS == "yes" && skipSNatARM == "yes" {
		t.Skip("Skipping gateway source NAT tests as SKIP_GATEWAY_SNAT_AWS and SKIP_GATEWAY_SNAT_ARM " +
			"are all set, even though SKIP_GATEWAY_SNAT isn't set")
	}

	//Setting default values for AWS_GW_SIZE
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	if skipSNatAWS == "yes" {
		t.Log("Skipping AWS gateway source NAT tests as SKIP_GATEWAY_SNAT_AWS is set")
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommon)
				preSpokeGatewayCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGatewaySNatDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccGatewaySNatConfigAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewaySNatExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "snat_mode", "customized_snat"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.src_cidr", "13.0.0.0/16"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.src_port", "22"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.dst_cidr", "14.0.0.0/16"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.dst_port", "222"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.protocol", "tcp"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.snat_ips", "175.32.12.12"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.snat_port", "12"),
					),
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: importStateVerifyIgnore,
				},
			},
		})
	}

	if skipSNatARM == "yes" {
		t.Log("Skipping ARM gateway source NAT tests as SKIP_GATEWAY_SNAT_ARM is set")
	} else {
		importStateVerifyIgnore = append(importStateVerifyIgnore, "vpc_id")
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommon)
				preSpokeGatewayCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGatewaySNatDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccGatewaySNatConfigARM(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewaySNatExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-arm-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "snat_mode", "customized_snat"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.src_cidr", "13.0.0.0/16"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.src_port", "22"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.dst_cidr", "14.0.0.0/16"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.dst_port", "222"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.protocol", "tcp"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.snat_ips", "175.32.12.12"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.snat_port", "12"),
					),
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: importStateVerifyIgnore,
				},
			},
		})
	}
}

func testAccGatewaySNatConfigAWS(rName string) string {
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfg-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "%[7]s"
	subnet       = "%[8]s"
}
resource "aviatrix_gateway_snat" "test" {
	gw_name   = aviatrix_spoke_gateway.test_spoke_gateway.gw_name
	snat_mode = "customized_snat"
	snat_policy {
		src_cidr    = "13.0.0.0/16"
		src_port    = "22"
		dst_cidr    = "14.0.0.0/16"
		dst_port    = "222"
		protocol    = "tcp"
		interface   = "eth0"
		connection  = "None"
		mark        = "22"
		snat_ips    = "175.32.12.12"
		snat_port   = "12"
		exclude_rtb = ""
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET"))
}

func testAccGatewaySNatConfigARM(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_arm" {
	account_name        = "tfa-arm-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc_arm.account_name
	gw_name      = "tfg-arm-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
resource "aviatrix_gateway_snat" "test" {
	gw_name   = aviatrix_spoke_gateway.test_spoke_gateway.gw_name
	snat_mode = "customized_snat"
	snat_policy {
		src_cidr    = "13.0.0.0/16"
		src_port    = "22"
		dst_cidr    = "14.0.0.0/16"
		dst_port    = "222"
		protocol    = "tcp"
		interface   = "eth0"
		connection  = "None"
		mark        = "22"
		snat_ips    = "175.32.12.12"
		snat_port   = "12"
		exclude_rtb = ""
	}
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("ARM_VNET_ID"), os.Getenv("ARM_REGION"),
		os.Getenv("ARM_GW_SIZE"), os.Getenv("ARM_SUBNET"))
}

func testAccCheckGatewaySNatExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix_gateway_snat Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix_gateway_snat ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundGateway := &goaviatrix.Gateway{
			GwName: rs.Primary.Attributes["gw_name"],
		}

		gw, err := client.GetGateway(foundGateway)
		if err != nil {
			return err
		}
		if foundGateway.GwName != rs.Primary.ID || gw.SnatMode != "customized" {
			return fmt.Errorf("resource 'aviatrix_gateway_snat' not found")
		}

		*gateway = *foundGateway
		return nil
	}
}

func testAccCheckGatewaySNatDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_gateway_snat" {
			continue
		}
		foundGateway := &goaviatrix.Gateway{
			GwName: rs.Primary.Attributes["gw_name"],
		}

		gw, err := client.GetGateway(foundGateway)
		if err == nil && gw.SnatMode == "customized" {
			return fmt.Errorf("resource 'aviatrix_gateway_snat' still exists")
		}
	}

	return nil
}
