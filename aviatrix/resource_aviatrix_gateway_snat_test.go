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

func TestAccAviatrixGatewaySNat_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	importStateVerifyIgnore := []string{"vnet_and_resource_group_names"}

	resourceName := "aviatrix_gateway_snat.test"

	msgCommon := ". Set SKIP_GATEWAY_SNAT to yes to skip gateway source NAT tests"

	skipSNat := os.Getenv("SKIP_GATEWAY_SNAT")
	skipSNatAWS := os.Getenv("SKIP_GATEWAY_SNAT_AWS")
	skipSNatAZURE := os.Getenv("SKIP_GATEWAY_SNAT_AZURE")

	if skipSNat == "yes" {
		t.Skip("Skipping gateway source NAT tests as SKIP_GATEWAY_SNAT is set")
	}
	if skipSNatAWS == "yes" && skipSNatAZURE == "yes" {
		t.Skip("Skipping gateway source NAT tests as SKIP_GATEWAY_SNAT_AWS and SKIP_GATEWAY_SNAT_AZURE " +
			"are all set, even though SKIP_GATEWAY_SNAT isn't set")
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
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.protocol", "tcp"),
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

	if skipSNatAZURE == "yes" {
		t.Log("Skipping Azure gateway source NAT tests as SKIP_GATEWAY_SNAT_AZURE is set")
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
					Config: testAccGatewaySNatConfigAZURE(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewaySNatExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "snat_mode", "customized_snat"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "snat_policy.0.protocol", "tcp"),
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
		src_cidr    = ""
		src_port    = ""
		dst_cidr    = ""
		dst_port    = ""
		protocol    = "tcp"
		interface   = "eth0"
		connection  = "None"
		mark        = ""
		snat_ips    = ""
		snat_port   = "12"
		exclude_rtb = ""
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID2"), os.Getenv("AWS_REGION2"), awsGwSize, os.Getenv("AWS_SUBNET2"))
}

func testAccGatewaySNatConfigAZURE(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_azure" {
	account_name        = "tfa-azure-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc_azure.account_name
	gw_name      = "tfg-azure-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
resource "aviatrix_gateway_snat" "test" {
	gw_name   = aviatrix_spoke_gateway.test_spoke_gateway.gw_name
	snat_mode = "customized_snat"
	snat_policy {
		src_cidr    = ""
		src_port    = ""
		dst_cidr    = ""
		dst_port    = ""
		protocol    = "tcp"
		interface   = "eth0"
		connection  = "None"
		mark        = ""
		snat_ips    = ""
		snat_port   = "12"
		exclude_rtb = ""
	}
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"))
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
		if err != nil && err != goaviatrix.ErrNotFound {
			return err
		} else if err == nil && gw.SnatMode == "customized" {
			return fmt.Errorf("resource 'aviatrix_gateway_snat' still exists")
		}
	}

	return nil
}
