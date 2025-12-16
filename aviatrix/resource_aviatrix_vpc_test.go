package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixVpc_basic(t *testing.T) {
	var vpc goaviatrix.Vpc

	rName := acctest.RandString(5)
	resourceName := "aviatrix_vpc.test_vpc"

	skipAcc := os.Getenv("SKIP_VPC")
	if skipAcc == "yes" {
		t.Skip("Skipping VPC tests as 'SKIP_VPC' is set")
	}

	skipAccAWS := os.Getenv("SKIP_VPC_AWS")
	skipAccAZURE := os.Getenv("SKIP_VPC_AZURE")
	skipAccGCP := os.Getenv("SKIP_VPC_GCP")
	if skipAccAWS == "yes" && skipAccAZURE == "yes" && skipAccGCP == "yes" {
		t.Skip("Skipping VPC tests as 'SKIP_VPC_AWS', 'SKIP_VPC_AZURE' and 'SKIP_VPC_GCP' are all set")
	}

	if skipAccAWS != "yes" {
		msgCommon := ". Set 'SKIP_VPC_AWS' to 'yes' to skip VPC tests in AWS"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preAccountCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckVpcDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccVpcConfigBasicAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckVpcExists(resourceName, &vpc),
						resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "cloud_type", "1"),
						resource.TestCheckResourceAttr(resourceName, "aviatrix_transit_vpc", "false"),
						resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
						resource.TestCheckResourceAttr(resourceName, "cidr", "10.0.0.0/16"),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping VPC tests in AWS as 'SKIP_VPC_AWS' is set")
	}

	if skipAccGCP != "yes" {
		msgCommon := ". Set 'SKIP_VPC_GCP' to 'yes' to skip VPC tests in GCP"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preAccountCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckVpcDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccVpcConfigBasicGCP(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckVpcExists(resourceName, &vpc),
						resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "cloud_type", "4"),
						resource.TestCheckResourceAttr(resourceName, "subnets.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "subnets.0.region", "us-east1"),
						resource.TestCheckResourceAttr(resourceName, "subnets.0.name", "us-east1-subnet"),
						resource.TestCheckResourceAttr(resourceName, "subnets.0.cidr", "10.0.0.0/16"),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping VPC tests in GCP as 'SKIP_VPC_GCP' is set")
	}

	if skipAccAZURE != "yes" {
		msgCommon := ". Set 'SKIP_VPC_AZURE' to 'yes' to skip VPC tests in Azure"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preAccountCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckVpcDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccVpcConfigBasicAZURE(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckVpcExists(resourceName, &vpc),
						resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "cloud_type", "8"),
						resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AZURE_REGION")),
						resource.TestCheckResourceAttr(resourceName, "cidr", "10.0.0.0/16"),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping VPC tests in Azure as 'SKIP_VPC_AZURE' is set")
	}
}

func TestAccAviatrixVpcGCPIPv6(t *testing.T) {
	var vpc goaviatrix.Vpc

	rName := acctest.RandString(5)
	resourceName := "aviatrix_vpc.test_vpc_ipv6"

	skipAcc := os.Getenv("SKIP_VPC_IPV6")
	if skipAcc == "yes" {
		t.Skip("Skipping VPC IPv6 tests as 'SKIP_VPC_IPV6' is set")
	}

	// Test GCP IPv6 support with ipv6_access_type
	skipGCP := os.Getenv("SKIP_VPC_IPV6_GCP")
	if skipGCP != "yes" {
		msgCommon := ". Set 'SKIP_VPC_IPV6_GCP' to 'yes' to skip VPC IPv6 tests in GCP"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preVpcIPv6GCPCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckVpcDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccVpcConfigGCPIPv6(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckVpcExists(resourceName, &vpc),
						resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfg-ipv6-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "cloud_type", "4"),
						resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
						resource.TestCheckResourceAttr(resourceName, "ipv6_access_type", "INTERNAL"),
						resource.TestCheckResourceAttr(resourceName, "subnets.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "subnets.0.region", "us-east1"),
						resource.TestCheckResourceAttr(resourceName, "subnets.0.name", "us-east1-subnet-ipv6"),
						resource.TestCheckResourceAttr(resourceName, "subnets.0.cidr", "10.0.0.0/16"),
						resource.TestCheckResourceAttr(resourceName, "subnets.0.ipv6_access_type", "EXTERNAL"),
						resource.TestCheckResourceAttrSet(resourceName, "subnets.0.ipv6_cidr"),
						resource.TestCheckResourceAttrSet(resourceName, "vpc_ipv6_cidr"),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"gcloud_project_credentials_filepath",
					},
				},
			},
		})
	} else {
		t.Log("Skipping VPC IPv6 tests in GCP as 'SKIP_VPC_IPV6_GCP' is set")
	}
}

func preVpcIPv6GCPCheck(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"GCP_PROJECT_ID",
		"GOOGLE_CREDENTIALS_FILEPATH",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func testAccVpcConfigBasicAWS(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc.account_name
	name         = "tfg-%s"
	region       = "%s"
	cidr         = "10.0.0.0/16"
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_REGION"))
}

func testAccVpcConfigBasicGCP(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc" {
	account_name                        = "tfa-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type   = 4
	account_name = aviatrix_account.test_acc.account_name
	name         = "tfg-%s"

	subnets {
		region = "us-east1"
		cidr   = "10.0.0.0/16"
		name   = "us-east1-subnet"
	}
}
`, rName, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"), rName)
}

func testAccVpcConfigBasicAZURE(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc" {
	account_name        = "tfa-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc.account_name
	name         = "tfg-%s"
	region       = "%s"
	cidr         = "10.0.0.0/16"
}
`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		rName, os.Getenv("AZURE_REGION"))
}

func testAccCheckVpcExists(n string, vpc *goaviatrix.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPC Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no VPC ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundVpc := &goaviatrix.Vpc{
			Name: rs.Primary.Attributes["name"],
		}

		foundVpc2, err := client.GetVpc(foundVpc)
		if err != nil {
			return err
		}
		if foundVpc2.Name != rs.Primary.ID {
			return fmt.Errorf("VPC not found")
		}

		*vpc = *foundVpc2
		return nil
	}
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpc" {
			continue
		}

		foundVpc := &goaviatrix.Vpc{
			Name: rs.Primary.Attributes["name"],
		}

		_, err := client.GetVpc(foundVpc)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("VPC still exists")
		}
	}

	return nil
}

func testAccVpcConfigGCPIPv6(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_vpc" "test_vpc_ipv6" {
	cloud_type       = 4
	account_name     = aviatrix_account.test_acc_gcp.account_name
	name             = "tfg-ipv6-%[1]s"
	enable_ipv6      = true
	ipv6_access_type = "INTERNAL"
	vpc_ipv6_cidr    = "fd00::/56"
	subnets {
		name             = "us-east1-subnet-ipv6"
		region           = "us-east1"
		cidr             = "10.0.0.0/16"
		ipv6_access_type = "EXTERNAL"
	}
}
	`, rName, os.Getenv("GCP_PROJECT_ID"), os.Getenv("GOOGLE_CREDENTIALS_FILEPATH"))
}
