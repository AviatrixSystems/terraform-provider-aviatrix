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

func TestAccAviatrixSpokeGateway_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	importStateVerifyIgnore := []string{"gcloud_project_credentials_filepath", "vnet_and_resource_group_names"}

	resourceName := "aviatrix_spoke_gateway.test_spoke_gateway"

	msgCommon := ". Set SKIP_SPOKE_GATEWAY to yes to skip Spoke Gateway tests"

	skipGw := os.Getenv("SKIP_SPOKE_GATEWAY")
	skipAWS := os.Getenv("SKIP_SPOKE_GATEWAY_AWS")
	skipGCP := os.Getenv("SKIP_SPOKE_GATEWAY_GCP")
	skipARM := os.Getenv("SKIP_SPOKE_GATEWAY_ARM")
	skipOCI := os.Getenv("SKIP_SPOKE_GATEWAY_OCI")

	if skipGw == "yes" {
		t.Skip("Skipping Spoke Gateway test as SKIP_SPOKE_GATEWAY is set")
	}

	if skipAWS == "yes" && skipGCP == "yes" && skipARM == "yes" && skipOCI == "yes" {
		t.Skip("Skipping Spoke Gateway test as SKIP_SPOKE_GATEWAY_AWS, SKIP_SPOKE_GATEWAY_GCP, " +
			"SKIP_SPOKE_GATEWAY_ARM, and SKIP_SPOKE_GATEWAY_OCI are all set, even though SKIP_SPOKE_GATEWAY isn't set")
	}

	//Setting default values for AWS_GW_SIZE and GCP_GW_SIZE
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	gcpGwSize := os.Getenv("GCP_GW_SIZE")
	ociGwSize := os.Getenv("OCI_GW_SIZE")

	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	if gcpGwSize == "" {
		gcpGwSize = "n1-standard-1"
	}
	if ociGwSize == "" {
		ociGwSize = "VM.Standard2.2"
	}

	if skipAWS == "yes" {
		t.Log("Skipping AWS Spoke Gateway test as SKIP_SPOKE_GATEWAY_AWS is set")
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommon)
				preSpokeGatewayCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGatewayConfigAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGatewayExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSize),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
						resource.TestCheckResourceAttr(resourceName, "single_ip_snat", "false"),
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

	if skipGCP == "yes" {
		t.Log("Skipping GCP Spoke Gateway test as SKIP_SPOKE_GATEWAY_GCP is set")
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommon)
				preSpokeGatewayCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGatewayConfigGCP(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGatewayExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "gw_size", gcpGwSize),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("GCP_VPC_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("GCP_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("GCP_ZONE")),
						resource.TestCheckResourceAttr(resourceName, "single_ip_snat", "false"),
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

	if skipARM == "yes" {
		t.Log("Skipping ARM Spoke Gateway test as SKIP_SPOKE_GATEWAY_ARM is set")
	} else {
		importStateVerifyIgnore = append(importStateVerifyIgnore, "vpc_id")
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommon)
				preSpokeGatewayCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGatewayConfigARM(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGatewayExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-arm-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "gw_size", os.Getenv("ARM_GW_SIZE")),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-arm-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("ARM_VNET_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("ARM_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("ARM_REGION")),
						resource.TestCheckResourceAttr(resourceName, "single_ip_snat", "false"),
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

	if skipOCI == "yes" {
		t.Log("Skipping OCI Spoke Gateway test as SKIP_SPOKE_GATEWAY_OCI is set")
	} else {
		//importStateVerifyIgnore = append(importStateVerifyIgnore, "vpc_id")
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckOCI(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGatewayConfigOCI(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGatewayExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-oci-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "gw_size", ociGwSize),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-oci-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("OCI_VPC_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("OCI_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("OCI_REGION")),
						resource.TestCheckResourceAttr(resourceName, "single_ip_snat", "false"),
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

func testAccSpokeGatewayConfigAWS(rName string) string {
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
	cloud_type     = 1
	account_name   = aviatrix_account.test_acc_aws.account_name
	gw_name        = "tfg-aws-%[1]s"
	vpc_id         = "%[5]s"
	vpc_reg        = "%[6]s"
	gw_size        = "%[7]s"
	subnet         = "%[8]s"
	single_ip_snat = false
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET"))
}

func testAccSpokeGatewayConfigGCP(rName string) string {
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
resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type     = 4
	account_name   = aviatrix_account.test_acc_gcp.account_name
	gw_name        = "tfg-gcp-%[1]s"
	vpc_id         = "%[4]s"
	vpc_reg        = "%[5]s"
	gw_size        = "%[6]s"
	subnet         = "%[7]s"
	single_ip_snat = false
}
	`, rName, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), gcpGwSize, os.Getenv("GCP_SUBNET"))
}

func testAccSpokeGatewayConfigARM(rName string) string {
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
	cloud_type     = 8
	account_name   = aviatrix_account.test_acc_arm.account_name
	gw_name        = "tfg-arm-%[1]s"
	vpc_id         = "%[6]s"
	vpc_reg        = "%[7]s"
	gw_size        = "%[8]s"
	subnet         = "%[9]s"
	single_ip_snat = false
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("ARM_VNET_ID"), os.Getenv("ARM_REGION"),
		os.Getenv("ARM_GW_SIZE"), os.Getenv("ARM_SUBNET"))
}

func testAccSpokeGatewayConfigOCI(rName string) string {
	ociGwSize := os.Getenv("OCI_GW_SIZE")
	if ociGwSize == "" {
		ociGwSize = "VM.Standard2.2"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_oci" {
	account_name                 = "tfa-oci-%s"
	cloud_type                   = 16
	oci_tenancy_id               = "%s"
	oci_user_id                  = "%s"
	oci_compartment_id           = "%s"
	oci_api_private_key_filepath = "%s"
}

resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type     = 16
	account_name   = aviatrix_account.test_acc_oci.account_name
	gw_name        = "tfg-oci-%[1]s"
	vpc_id         = "%[6]s"
	vpc_reg        = "%[7]s"
	gw_size        = "%[8]s"
	subnet         = "%[9]s"
	single_ip_snat = false
}
	`, rName, os.Getenv("OCI_TENANCY_ID"), os.Getenv("OCI_USER_ID"), os.Getenv("OCI_COMPARTMENT_ID"),
		os.Getenv("OCI_API_KEY_FILEPATH"), os.Getenv("OCI_VPC_ID"), os.Getenv("OCI_REGION"),
		ociGwSize, os.Getenv("OCI_SUBNET"))
}

func testAccCheckSpokeGatewayExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke gateway Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke gateway ID is set")
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
			return fmt.Errorf("spoke gateway not found")
		}

		*gateway = *foundGateway
		return nil
	}
}

func testAccCheckSpokeGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_vpc" {
			continue
		}
		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetGateway(foundGateway)
		if err == nil {
			return fmt.Errorf("spoke gateway still exists")
		}
	}

	return nil
}
