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

func preGatewayCheck(t *testing.T, msgCommon string) {
	preAccountCheck(t, msgCommon)

	awsVpcId := os.Getenv("AWS_VPC_ID")
	if awsVpcId == "" {
		t.Fatal("Environment variable AWS_VPC_ID is not set" + msgCommon)
	}
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		t.Fatal("Environment variable AWS_REGION is not set" + msgCommon)
	}
	awsSubnet := os.Getenv("AWS_SUBNET")
	if awsSubnet == "" {
		t.Fatal("Environment variable AWS_SUBNET is not set" + msgCommon)
	}
}

func preGatewayCheckGCP(t *testing.T, msgCommon string) {
	preAccountCheck(t, msgCommon)

	gcpVpcId := os.Getenv("GCP_VPC_ID")
	if gcpVpcId == "" {
		t.Fatal("Environment variable GCP_VPC_ID is not set" + msgCommon)
	}
	gcpZone := os.Getenv("GCP_ZONE")
	if gcpZone == "" {
		t.Fatal("Environment variable GCP_ZONE is not set" + msgCommon)
	}
	gcpSubnet := os.Getenv("GCP_SUBNET")
	if gcpSubnet == "" {
		t.Fatal("Environment variable GCP_SUBNET is not set" + msgCommon)
	}
}

func preGatewayCheckARM(t *testing.T, msgCommon string) {
	preAccountCheck(t, msgCommon)

	armVnetId := os.Getenv("ARM_VNET_ID")
	if armVnetId == "" {
		t.Fatal("Environment variable ARM_VNET_ID is not set" + msgCommon)
	}
	armRegion := os.Getenv("ARM_REGION")
	if armRegion == "" {
		t.Fatal("Environment variable ARM_REGION is not set" + msgCommon)
	}
	armSubnet := os.Getenv("ARM_SUBNET")
	if armSubnet == "" {
		t.Fatal("Environment variable ARM_SUBNET is not set" + msgCommon)
	}
	armGwSize := os.Getenv("ARM_GW_SIZE")
	if armGwSize == "" {
		t.Fatal("Environment variable ARM_GW_SIZE is not set" + msgCommon)
	}
}

func preGatewayCheckOCI(t *testing.T, msgCommon string) {
	preAccountCheck(t, msgCommon)

	ociVpcId := os.Getenv("OCI_VPC_ID")
	if ociVpcId == "" {
		t.Fatal("Environment variable OCI_VPC_ID is not set" + msgCommon)
	}
	ociRegion := os.Getenv("OCI_REGION")
	if ociRegion == "" {
		t.Fatal("Environment variable OCI_REGION is not set" + msgCommon)
	}
	ociSubnet := os.Getenv("OCI_SUBNET")
	if ociSubnet == "" {
		t.Fatal("Environment variable OCI_SUBNET is not set" + msgCommon)
	}
}

func preGatewayCheckAWSGOV(t *testing.T, msgCommon string) {
	preAccountCheck(t, msgCommon)

	awsgovVpcId := os.Getenv("AWSGOV_VPC_ID")
	if awsgovVpcId == "" {
		t.Fatal("Environment variable AWSGOV_VPC_ID is not set" + msgCommon)
	}
	awsgovRegion := os.Getenv("AWSGOV_REGION")
	if awsgovRegion == "" {
		t.Fatal("Environment variable AWSGOV_REGION is not set" + msgCommon)
	}
	awsgovSubnet := os.Getenv("AWSGOV_SUBNET")
	if awsgovSubnet == "" {
		t.Fatal("Environment variable AWSGOV_SUBNET is not set" + msgCommon)
	}
}

func TestAccAviatrixGateway_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	msgCommon := ". Set SKIP_GATEWAY to yes to skip Gateway tests"

	skipGw := os.Getenv("SKIP_GATEWAY")
	skipAWS := os.Getenv("SKIP_GATEWAY_AWS")
	skipGCP := os.Getenv("SKIP_GATEWAY_GCP")
	skipARM := os.Getenv("SKIP_GATEWAY_ARM")
	skipOCI := os.Getenv("SKIP_GATEWAY_OCI")
	skipAWSGOV := os.Getenv("SKIP_GATEWAY_AWSGOV")

	if skipGw == "yes" {
		t.Skip("Skipping Gateway test as SKIP_GATEWAY is set")
	}
	if skipAWS == "yes" && skipGCP == "yes" && skipARM == "yes" && skipOCI == "yes" && skipAWSGOV == "yes" {
		t.Skip("Skipping Gateway test as SKIP_GATEWAY_AWS, SKIP_GATEWAY_GCP, SKIP_GATEWAY_ARM " +
			",SKIP_GATEWAY_OCI, and SKIP_GATEWAY_AWSGOV are all set, even though SKIP_GATEWAY isn't set")
	}

	//Setting default values for AWS_GW_SIZE and GCP_GW_SIZE
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	gcpGwSize := os.Getenv("GCP_GW_SIZE")
	ociGwSize := os.Getenv("OCI_GW_SIZE")
	awsgovGwSize := os.Getenv("AWSGOV_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	if gcpGwSize == "" {
		gcpGwSize = "n1-standard-1"
	}
	if ociGwSize == "" {
		ociGwSize = "VM.Standard2.2"
	}
	if awsgovGwSize == "" {
		awsgovGwSize = "t2.micro"
	}

	if skipAWS == "yes" {
		t.Log("Skipping AWS Gateway test as SKIP_GATEWAY_AWS is set")
	} else {
		awsVpcId := os.Getenv("AWS_VPC_ID")
		awsRegion := os.Getenv("AWS_REGION")
		awsVpcNet := os.Getenv("AWS_SUBNET")
		resourceNameAws := "aviatrix_gateway.test_gw_aws"
		msgCommonAws := ". Set SKIP_GATEWAY_AWS to yes to skip AWS Gateway tests"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				//Checking resources have needed environment variables set
				preAccountCheck(t, msgCommon)
				preGatewayCheck(t, msgCommonAws)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccGatewayConfigBasicAWS(rName, awsGwSize, awsVpcId, awsRegion, awsVpcNet),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewayExists(resourceNameAws, &gateway),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_size", awsGwSize),
						resource.TestCheckResourceAttr(resourceNameAws, "vpc_id", awsVpcId),
						resource.TestCheckResourceAttr(resourceNameAws, "subnet", awsVpcNet),
						resource.TestCheckResourceAttr(resourceNameAws, "vpc_reg", awsRegion),
					),
				},
				{
					ResourceName:      resourceNameAws,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}

	if skipGCP == "yes" {
		t.Log("Skipping GCP Gateway test as SKIP_GATEWAY_GCP is set")
	} else {
		gcpZone := os.Getenv("GCP_ZONE")
		gcpVpcId := os.Getenv("GCP_VPC_ID")
		gcpSubnet := os.Getenv("GCP_SUBNET")
		resourceNameGcp := "aviatrix_gateway.test_gw_gcp"
		msgCommonGcp := ". Set SKIP_GATEWAY_GCP to yes to skip GCP Gateway tests"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				//Checking resources have needed environment variables set
				preAccountCheck(t, msgCommon)
				preGatewayCheckGCP(t, msgCommonGcp)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccGatewayConfigBasicGCP(rName, gcpGwSize, gcpVpcId, gcpZone, gcpSubnet),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewayExists(resourceNameGcp, &gateway),
						resource.TestCheckResourceAttr(resourceNameGcp, "gw_name", fmt.Sprintf("tfg-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameGcp, "gw_size", gcpGwSize),
						resource.TestCheckResourceAttr(resourceNameGcp, "vpc_id", gcpVpcId),
						resource.TestCheckResourceAttr(resourceNameGcp, "subnet", gcpSubnet),
						resource.TestCheckResourceAttr(resourceNameGcp, "vpc_reg", gcpZone),
					),
				},
				{
					ResourceName:      resourceNameGcp,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}

	if skipARM == "yes" {
		t.Log("Skipping ARM Gateway test as SKIP_GATEWAY_ARM is set")
	} else {
		armVnetId := os.Getenv("ARM_VNET_ID")
		armRegion := os.Getenv("ARM_REGION")
		armSubnet := os.Getenv("ARM_SUBNET")
		armGwSize := os.Getenv("ARM_GW_SIZE")
		resourceNameArm := "aviatrix_gateway.test_gw_arm"
		msgCommonArm := ". Set SKIP_GATEWAY_ARM to yes to skip ARM Gateway tests"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				//Checking resources have needed environment variables set
				preAccountCheck(t, msgCommon)
				preGatewayCheckARM(t, msgCommonArm)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccGatewayConfigBasicARM(rName, armGwSize, armVnetId, armRegion, armSubnet),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewayExists(resourceNameArm, &gateway),
						resource.TestCheckResourceAttr(resourceNameArm, "gw_name", fmt.Sprintf("tfg-arm-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameArm, "gw_size", armGwSize),
						resource.TestCheckResourceAttr(resourceNameArm, "vpc_id", armVnetId),
						resource.TestCheckResourceAttr(resourceNameArm, "subnet", armSubnet),
						resource.TestCheckResourceAttr(resourceNameArm, "vpc_reg", armRegion),
					),
				},
				{
					ResourceName:      resourceNameArm,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}

	if skipOCI == "yes" {
		t.Log("Skipping OCI Gateway test as SKIP_GATEWAY_OCI is set")
	} else {
		ociVpcId := os.Getenv("OCI_VPC_ID")
		ociRegion := os.Getenv("OCI_REGION")
		ociSubnet := os.Getenv("OCI_SUBNET")
		resourceNameOci := "aviatrix_gateway.test_gw_oci"
		msgCommonOci := ". Set SKIP_GATEWAY_OCI to yes to skip OCI Gateway tests"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				//Checking resources have needed environment variables set
				//preAccountCheck(t, msgCommon)
				preGatewayCheckOCI(t, msgCommonOci)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccGatewayConfigBasicOCI(rName, ociGwSize, ociVpcId, ociRegion, ociSubnet),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewayExists(resourceNameOci, &gateway),
						resource.TestCheckResourceAttr(resourceNameOci, "gw_name", fmt.Sprintf("tfg-oci-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameOci, "gw_size", ociGwSize),
						resource.TestCheckResourceAttr(resourceNameOci, "vpc_id", ociVpcId),
						resource.TestCheckResourceAttr(resourceNameOci, "subnet", ociSubnet),
						resource.TestCheckResourceAttr(resourceNameOci, "vpc_reg", ociRegion),
					),
				},
				{
					ResourceName:      resourceNameOci,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}

	if skipAWSGOV == "yes" {
		t.Log("Skipping AWSGOV Gateway test as SKIP_GATEWAY_AWSGOV is set")
	} else {
		awsgovVpcId := os.Getenv("AWSGOV_VPC_ID")
		awsgovRegion := os.Getenv("AWSGOV_REGION")
		awsgovVpcNet := os.Getenv("AWSGOV_SUBNET")
		resourceNameAwsgov := "aviatrix_gateway.test_gw_awsgov"
		msgCommonAwsgov := ". Set SKIP_GATEWAY_AWSGOV to yes to skip AWSGOV Gateway tests"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				//Checking resources have needed environment variables set
				preAccountCheck(t, msgCommon)
				preGatewayCheckAWSGOV(t, msgCommonAwsgov)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccGatewayConfigBasicAWSGOV(rName, awsgovGwSize, awsgovVpcId, awsgovRegion, awsgovVpcNet),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGatewayExists(resourceNameAwsgov, &gateway),
						resource.TestCheckResourceAttr(resourceNameAwsgov, "gw_name", fmt.Sprintf("tfg-awsgov-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAwsgov, "gw_size", awsgovGwSize),
						resource.TestCheckResourceAttr(resourceNameAwsgov, "vpc_id", awsgovVpcId),
						resource.TestCheckResourceAttr(resourceNameAwsgov, "subnet", awsgovVpcNet),
						resource.TestCheckResourceAttr(resourceNameAwsgov, "vpc_reg", awsgovRegion),
					),
				},
				{
					ResourceName:      resourceNameAwsgov,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccGatewayConfigBasicAWS(rName string, awsGwSize string, awsVpcId string, awsRegion string, awsVpcNet string) string {
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
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		awsVpcId, awsRegion, awsGwSize, awsVpcNet)
}

func testAccGatewayConfigBasicGCP(rName string, gcpGwSize string, gcpVpcId string, gcpZone string, gcpSubnet string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}

resource "aviatrix_gateway" "test_gw_gcp" {
	cloud_type   = 4
	account_name = aviatrix_account.test_acc_gcp.account_name
	gw_name      = "tfg-gcp-%[1]s"
	vpc_id       = "%[4]s"
	vpc_reg      = "%[5]s"
	gw_size      = "%[6]s"
	subnet       = "%[7]s"
}
	`, rName, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"),
		gcpVpcId, gcpZone, gcpGwSize, gcpSubnet)
}

func testAccGatewayConfigBasicARM(rName string, armGwSize string, armVnetId string, armRegion string, armSubnet string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_arm" {
	account_name        = "tfa-arm-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}

resource "aviatrix_gateway" "test_gw_arm" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc_arm.account_name
	gw_name      = "tfg-arm-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		armVnetId, armRegion, armGwSize, armSubnet)
}

func testAccGatewayConfigBasicOCI(rName string, ociGwSize string, ociVpcId string, ociRegion string, ociSubnet string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_oci" {
	account_name                 = "tfa-oci-%s"
	cloud_type                   = 16
	oci_tenancy_id               = "%s"
	oci_user_id                  = "%s"
	oci_compartment_id           = "%s"
	oci_api_private_key_filepath = "%s"
}

resource "aviatrix_gateway" "test_gw_oci" {
	cloud_type   = 16
	account_name = aviatrix_account.test_acc_oci.account_name
	gw_name      = "tfg-oci-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
	`,
		rName, os.Getenv("OCI_TENANCY_ID"), os.Getenv("OCI_USER_ID"), os.Getenv("OCI_COMPARTMENT_ID"),
		os.Getenv("OCI_API_KEY_FILEPATH"), ociVpcId, ociRegion, ociGwSize, ociSubnet)
}

func testAccGatewayConfigBasicAWSGOV(rName string, awsgovGwSize string, awsgovVpcId string, awsgovRegion string, awsgovVpcNet string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_awsgov" {
	account_name       = "tf-acc-awsgov-%s"
	cloud_type         = 256
	awsgov_account_number = "%s"
	awsgov_access_key     = "%s"
	awsgov_secret_key     = "%s"
}

resource "aviatrix_gateway" "test_gw_awsgov" {
	cloud_type   = 256
	account_name = aviatrix_account.test_acc_awsgov.account_name
	gw_name      = "tfg-awsgov-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "%[7]s"
	subnet       = "%[8]s"
}
	`, rName, os.Getenv("AWSGOV_ACCOUNT_NUMBER"), os.Getenv("AWSGOV_ACCESS_KEY"), os.Getenv("AWSGOV_SECRET_KEY"),
		awsgovVpcId, awsgovRegion, awsgovGwSize, awsgovVpcNet)
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
