package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func preSpokeGatewayCheck(t *testing.T, msgCommon string) string {
	preAccountCheck(t, msgCommon)

	azureGwSize := os.Getenv("AZURE_GW_SIZE")
	if azureGwSize == "" {
		t.Fatal("Environment variable AZURE_GW_SIZE is not set" + msgCommon)
	}
	return azureGwSize
}

func TestAccAviatrixSpokeGw_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	importStateVerifyIgnore := []string{"gcloud_project_credentials_filepath", "vnet_and_resource_group_names"}

	resourceName := "aviatrix_spoke_vpc.test_spoke_vpc"

	msgCommon := ". Set SKIP_SPOKE to yes to skip Spoke Gateway tests"

	skipGw := os.Getenv("SKIP_SPOKE")
	skipAWS := os.Getenv("SKIP_SPOKE_AWS")
	skipGCP := os.Getenv("SKIP_SPOKE_GCP")
	skipAZURE := os.Getenv("SKIP_SPOKE_AZURE")

	if skipGw == "yes" {
		t.Skip("Skipping Spoke Gateway test as SKIP_SPOKE is set")
	}

	if skipAWS == "yes" && skipGCP == "yes" && skipAZURE == "yes" {
		t.Skip("Skipping Spoke Gateway test as SKIP_SPOKE_AWS, SKIP_SPOKE_GCP, and SKIP_SPOKE_AZURE are all set, " +
			"even though SKIP_SPOKE isn't set")
	}

	//Setting default values for AWS_GW_SIZE and GCP_GW_SIZE
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	gcpGwSize := os.Getenv("GCP_GW_SIZE")

	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	if gcpGwSize == "" {
		gcpGwSize = "n1-standard-1"
	}
	if skipAWS == "yes" {
		t.Log("Skipping AWS Spoke Gateway test as SKIP_SPOKE_AWS is set")
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommon)
				preSpokeGatewayCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGwDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGwConfigAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGwExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_size", awsGwSize),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
						resource.TestCheckResourceAttr(resourceName, "enable_nat", "no"),
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
		t.Log("Skipping GCP Spoke Gateway test as SKIP_SPOKE_GCP is set")
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommon)
				preSpokeGatewayCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGwDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGwConfigGCP(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGwExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_size", gcpGwSize),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("GCP_VPC_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("GCP_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("GCP_ZONE")),
						resource.TestCheckResourceAttr(resourceName, "enable_nat", "no"),
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
	if skipAZURE == "yes" {
		t.Log("Skipping AZURE Spoke Gateway test as SKIP_SPOKE_AZURE is set")
	} else {
		importStateVerifyIgnore = append(importStateVerifyIgnore, "vpc_id")
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommon)
				preSpokeGatewayCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGwDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGwConfigAZURE(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGwExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_size", os.Getenv("AZURE_GW_SIZE")),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AZURE_VNET_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AZURE_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AZURE_REGION")),
						resource.TestCheckResourceAttr(resourceName, "enable_nat", "no"),
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

func testAccSpokeGwConfigAWS(rName string) string {
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_spoke_vpc" "test_spoke_vpc" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	vpc_size     = "%[7]s"
	subnet       = "%[8]s"
	enable_nat   = "no"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET"))
}

func testAccSpokeGwConfigGCP(rName string) string {
	gcpGwSize := os.Getenv("GCP_GW_SIZE")
	if gcpGwSize == "" {
		gcpGwSize = "n1-standard-1"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_spoke_vpc" "test_spoke_vpc" {
	cloud_type   = 4
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-gcp-%[1]s"
	vpc_id       = "%[4]s"
	vpc_reg      = "%[5]s"
	vpc_size     = "%[6]s"
	subnet       = "%[7]s"
	enable_nat   = "no"
}
	`, rName, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), gcpGwSize, os.Getenv("GCP_SUBNET"))
}

func testAccSpokeGwConfigAZURE(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name        = "tfa-azure-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_spoke_vpc" "test_spoke_vpc" {
	cloud_type   = 8
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-azure-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	vpc_size     = "%[8]s"
	subnet       = "%[9]s"
	single_az_ha = "enabled"
	enable_nat   = "no"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"))
}

func testAccCheckSpokeGwExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
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

func testAccCheckSpokeGwDestroy(s *terraform.State) error {
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
