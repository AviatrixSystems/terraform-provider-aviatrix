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

func TestAccAviatrixTransitGateway_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipGw := os.Getenv("SKIP_TRANSIT_GATEWAY")
	if skipGw == "yes" {
		t.Skip("Skipping Transit gateway test as SKIP_TRANSIT_GATEWAY is set")
	}

	skipGwAWS := os.Getenv("SKIP_TRANSIT_GATEWAY_AWS")
	skipGwAZURE := os.Getenv("SKIP_TRANSIT_GATEWAY_AZURE")
	skipGwGCP := os.Getenv("SKIP_TRANSIT_GATEWAY_GCP")
	skipGwOCI := os.Getenv("SKIP_TRANSIT_GATEWAY_OCI")

	if skipGwAWS == "yes" && skipGwAZURE == "yes" && skipGwGCP == "yes" && skipGwOCI == "yes" {
		t.Skip("Skipping Transit gateway test as SKIP_TRANSIT_GATEWAY_AWS, SKIP_TRANSIT_GATEWAY_AZURE, " +
			"SKIP_TRANSIT_GATEWAY_GCP and SKIP_TRANSIT_GATEWAY_OCI are all set")
	}

	if skipGwAWS != "yes" {
		resourceNameAws := "aviatrix_transit_gateway.test_transit_gateway_aws"
		msgCommonAws := ". Set SKIP_TRANSIT_GATEWAY_AWS to yes to skip Transit Gateway tests in aws"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommonAws)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameAws, &gateway),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_size", "t2.micro"),
						resource.TestCheckResourceAttr(resourceNameAws, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
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

	if skipGwAZURE != "yes" {
		resourceNameAzure := "aviatrix_transit_gateway.test_transit_gateway_azure"

		msgCommonAzure := ". Set SKIP_TRANSIT_GATEWAY_AZURE to yes to skip Transit Gateway tests in AZURE"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckAZURE(t, msgCommonAzure)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicAZURE(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameAzure, &gateway),
						resource.TestCheckResourceAttr(resourceNameAzure, "gw_name", fmt.Sprintf("tfg-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAzure, "gw_size", os.Getenv("AZURE_GW_SIZE")),
						resource.TestCheckResourceAttr(resourceNameAzure, "account_name", fmt.Sprintf("tfa-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAzure, "vpc_id", os.Getenv("AZURE_VNET_ID")),
						resource.TestCheckResourceAttr(resourceNameAzure, "subnet", os.Getenv("AZURE_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameAzure, "vpc_reg", os.Getenv("AZURE_REGION")),
					),
				},
				{
					ResourceName:      resourceNameAzure,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_Azure is set")
	}

	if skipGwGCP != "yes" {
		resourceNameGCP := "aviatrix_transit_gateway.test_transit_gateway_gcp"
		gcpGwSize := os.Getenv("GCP_GW_SIZE")
		if gcpGwSize == "" {
			gcpGwSize = "n1-standard-1"
		}

		msgCommonGCP := ". Set SKIP_TRANSIT_GATEWAY_GCP to yes to skip Transit Gateway tests in GCP"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckGCP(t, msgCommonGCP)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicGCP(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameGCP, &gateway),
						resource.TestCheckResourceAttr(resourceNameGCP, "gw_name", fmt.Sprintf("tfg-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameGCP, "gw_size", gcpGwSize),
						resource.TestCheckResourceAttr(resourceNameGCP, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameGCP, "vpc_id", os.Getenv("GCP_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameGCP, "subnet", os.Getenv("GCP_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameGCP, "vpc_reg", os.Getenv("GCP_ZONE")),
					),
				},
				{
					ResourceName:      resourceNameGCP,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_GCP is set")
	}

	if skipGwOCI != "yes" {
		resourceNameOCI := "aviatrix_transit_gateway.test_transit_gateway_oci"
		ociGwSize := os.Getenv("OCI_GW_SIZE")
		if ociGwSize == "" {
			ociGwSize = "VM.Standard2.2"
		}

		msgCommonOCI := ". Set SKIP_TRANSIT_GATEWAY_OCI to yes to skip Transit Gateway tests in OCI"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckGCP(t, msgCommonOCI)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicOCI(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameOCI, &gateway),
						resource.TestCheckResourceAttr(resourceNameOCI, "gw_name", fmt.Sprintf("tfg-oci-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameOCI, "gw_size", ociGwSize),
						resource.TestCheckResourceAttr(resourceNameOCI, "account_name", fmt.Sprintf("tfa-oci-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameOCI, "vpc_id", os.Getenv("OCI_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameOCI, "subnet", os.Getenv("OCI_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameOCI, "vpc_reg", os.Getenv("OCI_REGION")),
					),
				},
				{
					ResourceName:      resourceNameOCI,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_OCI is set")
	}
}

func testAccTransitGatewayConfigBasicAWS(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_aws" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfg-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitGatewayConfigBasicAZURE(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_azure" {
	account_name        = "tfa-azure-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_azure" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc_azure.account_name
	gw_name      = "tfg-azure-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"),
		os.Getenv("ARM_APPLICATION_KEY"), os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"))
}

func testAccTransitGatewayConfigBasicGCP(rName string) string {
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
	`, rName, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), gcpGwSize, os.Getenv("GCP_SUBNET"))
}

func testAccTransitGatewayConfigBasicOCI(rName string) string {
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
resource "aviatrix_transit_gateway" "test_transit_gateway_oci" {
	cloud_type   = 16
	account_name = aviatrix_account.test_acc_oci.account_name
	gw_name      = "tfg-oci-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
	`, rName, os.Getenv("OCI_TENANCY_ID"), os.Getenv("OCI_USER_ID"), os.Getenv("OCI_COMPARTMENT_ID"),
		os.Getenv("OCI_API_KEY_FILEPATH"), os.Getenv("OCI_VPC_ID"), os.Getenv("OCI_REGION"),
		ociGwSize, os.Getenv("OCI_SUBNET"))
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
