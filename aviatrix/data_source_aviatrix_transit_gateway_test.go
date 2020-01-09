package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceAviatrixTransitGateway_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_transit_gateway.foo"

	skipAcc := os.Getenv("SKIP_DATA_TRANSIT_GATEWAY")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Transit Gateway tests as SKIP_DATA_TRANSIT_GATEWAY is set")
	}

	skipAccAWS := os.Getenv("SKIP_DATA_TRANSIT_GATEWAY_AWS")
	skipAccARM := os.Getenv("SKIP_DATA_TRANSIT_GATEWAY_ARM")
	skipAccGCP := os.Getenv("SKIP_DATA_TRANSIT_GATEWAY_GCP")
	if skipAccAWS == "yes" && skipAccARM == "yes" && skipAccGCP == "yes" {
		t.Skip("Skipping Data Source Transit gateway tests as SKIP_DATA_TRANSIT_GATEWAY_AWS, SKIP_DATA_TRANSIT_GATEWAY_ARM and " +
			"SKIP_DATA_TRANSIT_GATEWAY_GCP are all set")
	}

	if skipAccAWS != "yes" {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, ". Set SKIP_DATA_TRANSIT_GATEWAY_AWS to yes to skip Data Source Transit Gateway tests in AWS")
			},
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceAviatrixTransitGatewayConfigBasic(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccDataSourceAviatrixTransitGateway(resourceName),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
						resource.TestCheckResourceAttr(resourceName, "gw_size", "t2.micro"),
					),
				},
			},
		})
	} else {
		t.Log("Skipping Data Source Transit gateway tests in AWS as SKIP_DATA_TRANSIT_GATEWAY_AWS is set")
	}

	if skipAccARM != "yes" {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckARM(t, ". Set SKIP_DATA_TRANSIT_GATEWAY_ARM to yes to skip Data Source Transit Gateway tests in ARM")
			},
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceAviatrixTransitGatewayConfigBasicARM(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccDataSourceAviatrixTransitGateway(resourceName),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-arm-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "gw_size", os.Getenv("ARM_GW_SIZE")),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-arm-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("ARM_VNET_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("ARM_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("ARM_REGION")),
					),
				},
			},
		})
	} else {
		t.Log("Skipping Data Source Transit gateway tests in ARM as SKIP_DATA_TRANSIT_GATEWAY_ARM is set")
	}

	if skipAccGCP != "yes" {
		gcpGwSize := os.Getenv("GCP_GW_SIZE")
		if gcpGwSize == "" {
			gcpGwSize = "n1-standard-1"
		}
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckGCP(t, ". Set SKIP_DATA_TRANSIT_GATEWAY_GCP to yes to skip Data Source Transit Gateway tests in GCP")
			},
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceAviatrixTransitGatewayConfigBasicGCP(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccDataSourceAviatrixTransitGateway(resourceName),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "gw_size", gcpGwSize),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("GCP_VPC_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("GCP_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("GCP_ZONE")),
					),
				},
			},
		})
	} else {
		t.Log("Skipping Data Source Transit gateway tests in GCP as SKIP_DATA_TRANSIT_GATEWAY_GCP is set")
	}
}

func testAccDataSourceAviatrixTransitGatewayConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name 	   = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_aws" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account.account_name
	gw_name      = "tfg-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
data "aviatrix_transit_gateway" "foo" {
	account_name = aviatrix_transit_gateway.test_transit_gateway_aws.account_name
	gw_name      = aviatrix_transit_gateway.test_transit_gateway_aws.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccDataSourceAviatrixTransitGatewayConfigBasicARM(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_arm" {
	account_name        = "tfa-arm-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_arm" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc_arm.account_name
	gw_name      = "tfg-arm-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
data "aviatrix_transit_gateway" "foo" {
	account_name = aviatrix_transit_gateway.test_transit_gateway_arm.account_name
	gw_name      = aviatrix_transit_gateway.test_transit_gateway_arm.gw_name
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"),
		os.Getenv("ARM_APPLICATION_KEY"), os.Getenv("ARM_VNET_ID"), os.Getenv("ARM_REGION"),
		os.Getenv("ARM_GW_SIZE"), os.Getenv("ARM_SUBNET"))
}

func testAccDataSourceAviatrixTransitGatewayConfigBasicGCP(rName string) string {
	gcpGwSize := os.Getenv("GCP_GW_SIZE")
	if gcpGwSize == "" {
		gcpGwSize = "n1-standard-1"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_gcp" {				
	cloud_type   = 4
	account_name = aviatrix_account.test_acc_gcp.account_name
	gw_name      = "tfg-gcp-%[1]s"
	vpc_id       = "%[4]s"
	vpc_reg      = "%[5]s"
	gw_size      = "%[6]s"
	subnet       = "%[7]s"
}
data "aviatrix_transit_gateway" "foo" {
	account_name = aviatrix_transit_gateway.test_transit_gateway_gcp.account_name
	gw_name      = aviatrix_transit_gateway.test_transit_gateway_gcp.gw_name
}
	`, rName, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), gcpGwSize, os.Getenv("GCP_SUBNET"))
}

func testAccDataSourceAviatrixTransitGateway(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
