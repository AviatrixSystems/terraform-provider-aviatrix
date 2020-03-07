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

func TestAccAviatrixTransitGw_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipGw := os.Getenv("SKIP_TRANSIT")
	if skipGw == "yes" {
		t.Skip("Skipping Transit gateway test as SKIP_TRANSIT is set")
	}

	skipGwAws := os.Getenv("SKIP_TRANSIT_AWS")
	skipGwAzure := os.Getenv("SKIP_TRANSIT_AZURE")
	if skipGwAws == "yes" && skipGwAzure == "yes" {
		t.Skip("Skipping Transit gateway test in aws as SKIP_TRANSIT_AWS and SKIP_TRANSIT_AZURE are both set")
	}

	if skipGwAws != "yes" {
		resourceNameAws := "aviatrix_transit_vpc.test_transit_vpc_aws"
		msgCommonAws := ". Set SKIP_TRANSIT_AWS to yes to skip Transit Gateway tests in aws"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommonAws)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGwDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGwConfigBasicAws(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGwExists(resourceNameAws, &gateway),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAws, "vpc_size", "t2.micro"),
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
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_AWS is set")
	}

	if skipGwAzure != "yes" {
		resourceNameAzure := "aviatrix_transit_vpc.test_transit_vpc_azure"

		msgCommonAzure := ". Set SKIP_TRANSIT_AZURE to yes to skip Transit Gateway tests in AZURE"
		preGatewayCheckAZURE(t, msgCommonAzure)

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGwDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGwConfigBasicAzure(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGwExists(resourceNameAzure, &gateway),
						resource.TestCheckResourceAttr(resourceNameAzure, "gw_name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAzure, "vpc_size", os.Getenv("AZURE_GW_SIZE")),
						resource.TestCheckResourceAttr(resourceNameAzure, "account_name", fmt.Sprintf("tfaz-%s", rName)),
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
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_AZURE is set")
	}
}

func testAccTransitGwConfigBasicAws(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_aws" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_vpc" "test_transit_vpc_aws" {
	cloud_type   = 1
	account_name = aviatrix_account.test_aws.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	vpc_size     = "t2.micro"
	subnet       = "%[7]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitGwConfigBasicAzure(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_azure" {
	account_name        = "tfaz-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_vpc" "test_transit_vpc_azure" {
	cloud_type   = 8
	account_name = aviatrix_account.test_azure.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	vpc_size     = "%[8]s"
	subnet       = "%[9]s"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"),
		os.Getenv("ARM_APPLICATION_KEY"), os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"))
}

func testAccCheckTransitGwExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
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

func testAccCheckTransitGwDestroy(s *terraform.State) error {
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
