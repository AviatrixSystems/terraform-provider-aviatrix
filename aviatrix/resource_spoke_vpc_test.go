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

func preSpokeGatewayCheck(t *testing.T, msgCommon string) string {
	preAccountCheck(t, msgCommon)

	armGwSize := os.Getenv("ARM_GW_SIZE")
	if armGwSize == "" {
		t.Fatal("Environment variable ARM_GW_SIZE is not set" + msgCommon)
	}
	return armGwSize
}

func TestAccAviatrixSpokeGw_basic(t *testing.T) {
	var gateway goaviatrix.Gateway
	rName := fmt.Sprintf("%s", acctest.RandString(5))
	importStateVerifyIgnore := []string{"gcloud_project_credentials_filepath", "vnet_and_resource_group_names"}

	resourceName := "aviatrix_spoke_vpc.test_spoke_vpc"

	msgCommon := ". Set SKIP_SPOKE to yes to skip Spoke Gateway tests"

	skipGw := os.Getenv("SKIP_SPOKE")
	skipAWS := os.Getenv("SKIP_AWS_SPOKE")
	skipGCP := os.Getenv("SKIP_GCP_SPOKE")
	skipARM := os.Getenv("SKIP_ARM_SPOKE")

	if skipGw == "yes" {
		t.Skip("Skipping Spoke Gateway test as SKIP_SPOKE is set")
	}

	if skipAWS == "yes" && skipGCP == "yes" && skipARM == "yes" {
		t.Skip("Skipping Spoke Gateway test as SKIP_AWS_SPOKE, SKIP_GCP_SPOKE, and SKIP_ARM_SPOKE are all set, " +
			"even though SKIP_SPOKE isn't set")
	}

	preGatewayCheck(t, msgCommon)
	preSpokeGatewayCheck(t, msgCommon)

	//Setting default values for AWS_GW_SIZE and GCP_GW_SIZE
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	gcpGwSize := os.Getenv("GCP_GW_SIZE")

	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	if gcpGwSize == "" {
		gcpGwSize = "f1-micro"
	}
	if skipAWS == "yes" {
		t.Log("Skipping AWS Spoke Gateway test as SKIP_AWS_SPOKE is set")
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGwDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGwConfigAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGwExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName,
							"gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_size", awsGwSize),
						resource.TestCheckResourceAttr(resourceName,
							"account_name", fmt.Sprintf("tfa-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_VPC_NET")),
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
		t.Log("Skipping GCP Spoke Gateway test as SKIP_GCP_SPOKE is set")
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGwDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGwConfigGCP(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGwExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName,
							"gw_name", fmt.Sprintf("tfg-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_size", gcpGwSize),
						resource.TestCheckResourceAttr(resourceName,
							"account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
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
	if skipARM == "yes" {
		t.Log("Skipping ARM Spoke Gateway test as SKIP_ARM_SPOKE is set")
	} else {
		importStateVerifyIgnore = append(importStateVerifyIgnore, "vpc_id")
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGwDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGwConfigARM(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGwExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName,
							"gw_name", fmt.Sprintf("tfg-arm-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_size", os.Getenv("ARM_GW_SIZE")),
						resource.TestCheckResourceAttr(resourceName,
							"account_name", fmt.Sprintf("tfa-arm-%s", rName)),
						resource.TestCheckResourceAttr(resourceName,
							"vnet_and_resource_group_names", os.Getenv("ARM_VNET_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("ARM_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("ARM_REGION")),
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
	account_name 	   = "tfa-aws-%s"
	cloud_type 		   = 1
	aws_account_number = "%s"
	aws_iam 		   = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_vpc" "test_spoke_vpc" {
    cloud_type   = 1
    account_name = "${aviatrix_account.test.account_name}"
    gw_name      = "tfg-aws-%[1]s"
    vpc_id       = "%[5]s"
    vpc_reg      = "%[6]s"
    vpc_size     = "%[7]s"
    subnet       = "%[8]s"
    enable_nat   = "no"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_VPC_NET"))
}

func testAccSpokeGwConfigGCP(rName string) string {
	gcpGwSize := os.Getenv("GCP_GW_SIZE")
	if gcpGwSize == "" {
		gcpGwSize = "f1-micro"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name 						= "tfa-gcp-%s"
	cloud_type 							= 4
	gcloud_project_id 					= "%s"
	gcloud_project_credentials_filepath = "%s"
}

resource "aviatrix_spoke_vpc" "test_spoke_vpc" {
    cloud_type   = 4
    account_name = "${aviatrix_account.test.account_name}"
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

func testAccSpokeGwConfigARM(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name 		= "tfa-arm-%s"
	cloud_type 			= 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}

resource "aviatrix_spoke_vpc" "test_spoke_vpc" {
    cloud_type                    = 8
    account_name                  = "${aviatrix_account.test.account_name}"
    gw_name                       = "tfg-arm-%[1]s"
    vnet_and_resource_group_names = "%[6]s"
    vpc_reg                       = "%[7]s"
    vpc_size                      = "%[8]s"
    subnet                        = "%[9]s"
    enable_nat                    = "no"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("ARM_VNET_ID"), os.Getenv("ARM_REGION"),
		os.Getenv("ARM_GW_SIZE"), os.Getenv("ARM_SUBNET"))
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
