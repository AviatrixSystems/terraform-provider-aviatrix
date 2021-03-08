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

func TestAccAviatrixGatewayDNat_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	msgCommon := ". Set SKIP_GATEWAY_DNAT to yes to skip gateway DNAT tests"

	skipDNat := os.Getenv("SKIP_GATEWAY_DNAT")
	skipDNatAWS := os.Getenv("SKIP_GATEWAY_DNAT_AWS")
	skipDNatAZURE := os.Getenv("SKIP_GATEWAY_DNAT_AZURE")

	if skipDNat == "yes" {
		t.Skip("Skipping gateway destination NAT tests as SKIP_GATEWAY_DNAT is set")
	}
	if skipDNatAWS == "yes" && skipDNatAZURE == "yes" {
		t.Skip("Skipping gateway destination NAT tests as SKIP_GATEWAY_DNAT_AWS and SKIP_GATEWAY_DNAT_AZURE " +
			"are all set, even though SKIP_GATEWAY_DNAT isn't set")
	}

	if skipDNatAWS == "yes" {
		t.Log("Skipping AWS gateway destination NAT tests as SKIP_GATEWAY_DNAT_AWS is set")
	} else {
		resourceName := "aviatrix_gateway_dnat.test"
		msgCommonAws := ". Set SKIP_GATEWAY_DNAT_AWS to yes to skip AWS gateway destination NAT tests"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				//Checking resources have needed environment variables set
				preAccountCheck(t, msgCommon)
				preGatewayCheck(t, msgCommonAws)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGatewayDNatDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccGatewayDNatConfigBasicAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewayDNatExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "dnat_policy.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "dnat_policy.0.protocol", "tcp"),
						resource.TestCheckResourceAttr(resourceName, "dnat_policy.0.dnat_port", "12"),
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

	if skipDNatAZURE == "yes" {
		t.Log("Skipping AWS gateway destination NAT tests as SKIP_GATEWAY_DNAT_AZURE is set")
	} else {
		resourceName := "aviatrix_gateway_dnat.test"
		msgCommonAZURE := ". Set SKIP_GATEWAY_DNAT_AZURE to yes to skip gateway destination NAT tests"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				//Checking resources have needed environment variables set
				preAccountCheck(t, msgCommon)
				preGatewayCheckAZURE(t, msgCommonAZURE)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGatewayDNatDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccGatewayDNatConfigBasicAZURE(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewayExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "dnat_policy.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "dnat_policy.0.protocol", "tcp"),
						resource.TestCheckResourceAttr(resourceName, "dnat_policy.0.dnat_port", "12"),
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
}

func testAccGatewayDNatConfigBasicAWS(rName string) string {
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tf-acc-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "test_gw_aws" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfg-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "%[7]s"
	subnet       = "%[8]s"
}
resource "aviatrix_gateway_dnat" "test" {
	gw_name = aviatrix_gateway.test_gw_aws.gw_name
	dnat_policy {
		src_cidr    = ""
		src_port    = ""
		dst_cidr    = ""
		dst_port    = ""
		protocol    = "tcp"
		interface   = "eth0"
		connection  = "None"
		mark        = ""
		dnat_ips    = ""
		dnat_port   = "12"
		exclude_rtb = ""
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET"))
}

func testAccGatewayDNatConfigBasicAZURE(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_azure" {
	account_name        = "tfa-azure-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_gateway" "test_gw_azure" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc_azure.account_name
	gw_name      = "tfg-azure-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
resource "aviatrix_gateway_dnat" "test" {
	gw_name = aviatrix_gateway.test_gw_azure.gw_name
	dnat_policy {
		src_cidr    = ""
		src_port    = ""
		dst_cidr    = ""
		dst_port    = ""
		protocol    = "tcp"
		interface   = "eth0"
		connection  = "None"
		mark        = ""
		dnat_ips    = ""
		dnat_port   = "12"
		exclude_rtb = ""
	}
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"))
}

func testAccCheckGatewayDNatExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("'aviatrix_gateway_dnat' Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no 'aviatrix_gateway_dnat' ID is set")
		}
		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundGateway := &goaviatrix.Gateway{
			GwName: rs.Primary.Attributes["gw_name"],
		}

		foundGateway2, err := client.GetGateway(foundGateway)
		if err != nil {
			return err
		}
		gwDetail, err := client.GetGatewayDetail(foundGateway)
		if err != nil {
			return fmt.Errorf("couldn't get detail information of Aviatrix gateway(name: %s) due to: %s", foundGateway.GwName, err)
		}
		if len(gwDetail.DnatPolicy) == 0 {
			return fmt.Errorf("resource 'aviatrix_gateway_dnat' not found")
		}

		*gateway = *foundGateway2
		return nil
	}
}

func testAccCheckGatewayDNatDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_gateway_dnat" {
			continue
		}
		foundGateway := &goaviatrix.Gateway{
			GwName: rs.Primary.Attributes["gw_name"],
		}

		_, err := client.GetGateway(foundGateway)
		if err != nil {
			if err != goaviatrix.ErrNotFound {
				return err
			}
		} else {
			gwDetail, err := client.GetGatewayDetail(foundGateway)
			if err == nil && len(gwDetail.DnatPolicy) != 0 {
				return fmt.Errorf("resource 'aviatrix_gateway_snat' still exists")
			}
		}
	}

	return nil
}
