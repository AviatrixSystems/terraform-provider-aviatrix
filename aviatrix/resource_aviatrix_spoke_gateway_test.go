package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func preSpokeGatewayCheck(t *testing.T, msgCommon string) string {
	preAccountCheck(t, msgCommon)

	azureGwSize := os.Getenv("AZURE_GW_SIZE")
	if azureGwSize == "" {
		t.Fatal("Environment variable AZURE_GW_SIZE is not set" + msgCommon)
	}
	return azureGwSize
}

func preAwsSpokeGatewayCheck(t *testing.T, msgCommon string) string {
	requiredEnvVars := []string{
		"AWS_VPC_ID4",
		"AWS_SUBNET4",
		"AWS_REGION",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
	return ""
}

func preAwsSpokeGatewayInsertionCheck(t *testing.T, msgCommon string) string {
	requiredEnvVars := []string{
		"AWS_VPC_ID4",
		"AWS_REGION",
		"AWS_AVAILABILITY_ZONE",
		"AWS_INSERTION_SUBNET",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
	return ""
}

func TestAccAviatrixSpokeGateway_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	importStateVerifyIgnore := []string{"gcloud_project_credentials_filepath", "vnet_and_resource_group_names"}

	resourceName := "aviatrix_spoke_gateway.test_spoke_gateway"

	msgCommon := ". Set SKIP_SPOKE_GATEWAY to yes to skip Spoke Gateway tests"

	skipGw := os.Getenv("SKIP_SPOKE_GATEWAY")
	skipAWS := os.Getenv("SKIP_SPOKE_GATEWAY_AWS")
	skipGCP := os.Getenv("SKIP_SPOKE_GATEWAY_GCP")
	skipAZURE := os.Getenv("SKIP_SPOKE_GATEWAY_AZURE")
	skipOCI := os.Getenv("SKIP_SPOKE_GATEWAY_OCI")

	if skipGw == "yes" {
		t.Skip("Skipping Spoke Gateway test as SKIP_SPOKE_GATEWAY is set")
	}

	if skipAWS == "yes" && skipGCP == "yes" && skipAZURE == "yes" && skipOCI == "yes" {
		t.Skip("Skipping Spoke Gateway test as SKIP_SPOKE_GATEWAY_AWS, SKIP_SPOKE_GATEWAY_GCP, " +
			"SKIP_SPOKE_GATEWAY_AZURE, and SKIP_SPOKE_GATEWAY_OCI are all set, even though SKIP_SPOKE_GATEWAY isn't set")
	}

	// Setting default values for AWS_GW_SIZE and GCP_GW_SIZE
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
				preAwsSpokeGatewayCheck(t, msgCommon)
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
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID4")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_SUBNET4")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
						resource.TestCheckResourceAttr(resourceName, "single_ip_snat", "false"),
						resource.TestCheckResourceAttr(resourceName, "bgp_polling_time", "50"),
						resource.TestCheckResourceAttr(resourceName, "bgp_neighbor_status_polling_time", "5"),
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

	skipAWSInsertion := os.Getenv("SKIP_SPOKE_GATEWAY_AWS_INSERTION")
	if skipAWSInsertion == "yes" {
		t.Log("Skipping AWS Spoke Gateway Insertion test as SKIP_SPOKE_GATEWAY_AWS_INSERTION is set")
	} else if skipAWS != "yes" {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preAwsSpokeGatewayInsertionCheck(t, msgCommon)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSpokeGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSpokeGatewayConfigAWSInsertion(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGatewayExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-insertion-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSize),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID4")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
						resource.TestCheckResourceAttr(resourceName, "insertion_gateway", "true"),
						resource.TestCheckResourceAttr(resourceName, "insertion_gateway_az", os.Getenv("AWS_AVAILABILITY_ZONE")),
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

	if skipAZURE == "yes" {
		t.Log("Skipping Azure Spoke Gateway test as SKIP_SPOKE_GATEWAY_AZURE is set")
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
					Config: testAccSpokeGatewayConfigAZURE(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSpokeGatewayExists(resourceName, &gateway),
						resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "gw_size", os.Getenv("AZURE_GW_SIZE")),
						resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AZURE_VNET_ID")),
						resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AZURE_SUBNET")),
						resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AZURE_REGION")),
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
		// importStateVerifyIgnore = append(importStateVerifyIgnore, "vpc_id")
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
	cloud_type                       = 1
	account_name                     = aviatrix_account.test_acc_aws.account_name
	gw_name                          = "tfg-aws-%[1]s"
	vpc_id                           = "%[5]s"
	vpc_reg                          = "%[6]s"
	gw_size                          = "%[7]s"
	subnet                           = "%[8]s"
	single_ip_snat                   = false
	bgp_polling_time                 = 50
	bgp_neighbor_status_polling_time = 5
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID4"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET4"))
}

func testAccSpokeGatewayConfigAWSInsertion(rName string) string {
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
	cloud_type           = 1
	account_name         = aviatrix_account.test_acc_aws.account_name
	gw_name              = "tfg-aws-insertion-%[1]s"
	vpc_id               = "%[5]s"
	vpc_reg              = "%[6]s"
	gw_size              = "%[7]s"
	subnet               = "%[8]s"
	insertion_gateway    = true
	insertion_gateway_az = "%[9]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID4"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_INSERTION_SUBNET"), os.Getenv("AWS_AVAILABILITY_ZONE"))
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

func testAccSpokeGatewayConfigAZURE(rName string) string {
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
	cloud_type     = 8
	account_name   = aviatrix_account.test_acc_azure.account_name
	gw_name        = "tfg-azure-%[1]s"
	vpc_id         = "%[6]s"
	vpc_reg        = "%[7]s"
	gw_size        = "%[8]s"
	subnet         = "%[9]s"
	single_ip_snat = false
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"))
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

		client := mustClient(testAccProvider.Meta())

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
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_gateway" {
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

// TestAccAviatrixSpokeGateway_ipv6AWS tests IPv6 CIDR fields for AWS spoke gateway
func TestAccAviatrixSpokeGateway_ipv6AWS(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_gateway.test_spoke_gateway_ipv6"

	msgCommon := ". Set SKIP_SPOKE_GATEWAY_IPV6 to yes to skip Spoke Gateway IPv6 tests"

	skipGwIPv6 := os.Getenv("SKIP_SPOKE_GATEWAY_IPV6")
	if skipGwIPv6 == "yes" {
		t.Skip("Skipping Spoke Gateway IPv6 test as SKIP_SPOKE_GATEWAY_IPV6 is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAwsSpokeGatewayIPv6Check(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeGatewayConfigAWSIPv6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-ipv6-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSize),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID4")),
					resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_SUBNET4")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ipv6_cidr", os.Getenv("AWS_SUBNET_IPV6_CIDR")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
				},
			},
		},
	})
}

// TestAccAviatrixSpokeGateway_ipv6WithHA tests IPv6 CIDR fields with HA enabled
func TestAccAviatrixSpokeGateway_ipv6WithHA(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_gateway.test_spoke_gateway_ipv6_ha"

	msgCommon := ". Set SKIP_SPOKE_GATEWAY_IPV6 to yes to skip Spoke Gateway IPv6 tests"

	skipGwIPv6 := os.Getenv("SKIP_SPOKE_GATEWAY_IPV6")
	if skipGwIPv6 == "yes" {
		t.Skip("Skipping Spoke Gateway IPv6 HA test as SKIP_SPOKE_GATEWAY_IPV6 is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAwsSpokeGatewayIPv6HACheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeGatewayConfigAWSIPv6WithHA(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-ipv6-ha-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSize),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID4")),
					resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_SUBNET4")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ipv6_cidr", os.Getenv("AWS_SUBNET_IPV6_CIDR")),
					resource.TestCheckResourceAttr(resourceName, "ha_subnet", os.Getenv("AWS_HA_SUBNET")),
					resource.TestCheckResourceAttr(resourceName, "ha_subnet_ipv6_cidr", os.Getenv("AWS_HA_SUBNET_IPV6_CIDR")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
				},
			},
		},
	})
}

// TestAccAviatrixSpokeGateway_ipv6Azure tests IPv6 CIDR fields for Azure spoke gateway
func TestAccAviatrixSpokeGateway_ipv6Azure(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_gateway.test_spoke_gateway_ipv6_azure"

	msgCommon := ". Set SKIP_SPOKE_GATEWAY_IPV6_AZURE to yes to skip Azure Spoke Gateway IPv6 tests"

	skipGwIPv6Azure := os.Getenv("SKIP_SPOKE_GATEWAY_IPV6_AZURE")
	if skipGwIPv6Azure == "yes" {
		t.Skip("Skipping Azure Spoke Gateway IPv6 test as SKIP_SPOKE_GATEWAY_IPV6_AZURE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAzureSpokeGatewayIPv6Check(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeGatewayConfigAzureIPv6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-azure-ipv6-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", os.Getenv("AZURE_GW_SIZE")),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-azure-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AZURE_VNET_ID")),
					resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AZURE_SUBNET")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AZURE_REGION")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ipv6_cidr", os.Getenv("AZURE_SUBNET_IPV6_CIDR")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
					"vpc_id",
				},
			},
		},
	})
}

// TestAccAviatrixSpokeGateway_ipv6GCP tests IPv6 for GCP spoke gateway (no subnet_ipv6_cidr required)
func TestAccAviatrixSpokeGateway_ipv6GCP(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_gateway.test_spoke_gateway_ipv6_gcp"

	msgCommon := ". Set SKIP_SPOKE_GATEWAY_IPV6_GCP to yes to skip GCP Spoke Gateway IPv6 tests"

	skipGwIPv6GCP := os.Getenv("SKIP_SPOKE_GATEWAY_IPV6_GCP")
	if skipGwIPv6GCP == "yes" {
		t.Skip("Skipping GCP Spoke Gateway IPv6 test as SKIP_SPOKE_GATEWAY_IPV6_GCP is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGCPSpokeGatewayIPv6Check(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeGatewayConfigGCPIPv6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-gcp-ipv6-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", "n1-standard-1"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("GCP_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("GCP_ZONE")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					// For GCP, subnet_ipv6_cidr should be computed/optional, not required
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ipv6_cidr"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
				},
			},
		},
	})
}

func preAwsSpokeGatewayIPv6Check(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"AWS_VPC_ID4",
		"AWS_SUBNET4",
		"AWS_REGION",
		"AWS_SUBNET_IPV6_CIDR",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func preAwsSpokeGatewayIPv6HACheck(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"AWS_VPC_ID4",
		"AWS_SUBNET4",
		"AWS_REGION",
		"AWS_SUBNET_IPV6_CIDR",
		"AWS_HA_SUBNET",
		"AWS_HA_SUBNET_IPV6_CIDR",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func preAzureSpokeGatewayIPv6Check(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"AZURE_VNET_ID",
		"AZURE_SUBNET",
		"AZURE_REGION",
		"AZURE_GW_SIZE",
		"AZURE_SUBNET_IPV6_CIDR",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func preGCPSpokeGatewayIPv6Check(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"GCP_VPC_ID",
		"GCP_ZONE",
		"GCP_SUBNET",
		"GCP_PROJECT_ID",
		"GOOGLE_CREDENTIALS_FILEPATH",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func testAccSpokeGatewayConfigAWSIPv6(rName string) string {
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
resource "aviatrix_spoke_gateway" "test_spoke_gateway_ipv6" {
	cloud_type        = 1
	account_name      = aviatrix_account.test_acc_aws.account_name
	gw_name           = "tfg-aws-ipv6-%[1]s"
	vpc_id            = "%[5]s"
	vpc_reg           = "%[6]s"
	gw_size           = "%[7]s"
	subnet            = "%[8]s"
	enable_ipv6       = true
	subnet_ipv6_cidr  = "%[9]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID4"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET4"),
		os.Getenv("AWS_SUBNET_IPV6_CIDR"))
}

func testAccSpokeGatewayConfigAWSIPv6WithHA(rName string) string {
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
resource "aviatrix_spoke_gateway" "test_spoke_gateway_ipv6_ha" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_acc_aws.account_name
	gw_name              = "tfg-aws-ipv6-ha-%[1]s"
	vpc_id               = "%[5]s"
	vpc_reg              = "%[6]s"
	gw_size              = "%[7]s"
	subnet               = "%[8]s"
	enable_ipv6          = true
	subnet_ipv6_cidr     = "%[9]s"
	ha_subnet            = "%[10]s"
	ha_subnet_ipv6_cidr  = "%[11]s"
	ha_gw_size           = "%[7]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID4"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET4"),
		os.Getenv("AWS_SUBNET_IPV6_CIDR"), os.Getenv("AWS_HA_SUBNET"), os.Getenv("AWS_HA_SUBNET_IPV6_CIDR"))
}

func testAccSpokeGatewayConfigAzureIPv6(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_azure" {
	account_name        = "tfa-azure-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_spoke_gateway" "test_spoke_gateway_ipv6_azure" {
	cloud_type       = 8
	account_name     = aviatrix_account.test_acc_azure.account_name
	gw_name          = "tfg-azure-ipv6-%[1]s"
	vpc_id           = "%[6]s"
	vpc_reg          = "%[7]s"
	gw_size          = "%[8]s"
	subnet           = "%[9]s"
	enable_ipv6      = true
	subnet_ipv6_cidr = "%[10]s"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"), os.Getenv("AZURE_SUBNET_IPV6_CIDR"))
}

// TestAccAviatrixSpokeGateway_ipv6WithInsaneMode tests IPv6 with Insane Mode enabled
func TestAccAviatrixSpokeGateway_ipv6WithInsaneMode(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_gateway.test_spoke_gateway_ipv6_insane"

	msgCommon := ". Set SKIP_SPOKE_GATEWAY_IPV6_INSANE_MODE to yes to skip Spoke Gateway IPv6 Insane Mode tests"

	skipGwIPv6InsaneMode := os.Getenv("SKIP_SPOKE_GATEWAY_IPV6_INSANE_MODE")
	if skipGwIPv6InsaneMode == "yes" {
		t.Skip("Skipping Spoke Gateway IPv6 Insane Mode test as SKIP_SPOKE_GATEWAY_IPV6_INSANE_MODE is set")
	}

	awsGwSize := os.Getenv("AWS_INSANE_MODE_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "c5.xlarge"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAwsSpokeGatewayIPv6InsaneModeCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeGatewayConfigAWSIPv6WithInsaneMode(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-ipv6-insane-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSize),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID4")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ipv6_cidr", os.Getenv("AWS_SUBNET_IPV6_CIDR")),
					resource.TestCheckResourceAttr(resourceName, "insane_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "insane_mode_az", os.Getenv("AWS_AVAILABILITY_ZONE")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
				},
			},
		},
	})
}

func preAwsSpokeGatewayIPv6InsaneModeCheck(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"AWS_VPC_ID4",
		"AWS_REGION",
		"AWS_SUBNET_IPV6_CIDR",
		"AWS_AVAILABILITY_ZONE",
		"AWS_INSANE_MODE_SUBNET",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func testAccSpokeGatewayConfigAWSIPv6WithInsaneMode(rName string) string {
	awsGwSize := os.Getenv("AWS_INSANE_MODE_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "c5.xlarge"
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
resource "aviatrix_spoke_gateway" "test_spoke_gateway_ipv6_insane" {
	cloud_type        = 1
	account_name      = aviatrix_account.test_acc_aws.account_name
	gw_name           = "tfg-aws-ipv6-insane-%[1]s"
	vpc_id            = "%[5]s"
	vpc_reg           = "%[6]s"
	gw_size           = "%[7]s"
	subnet            = "%[8]s"
	enable_ipv6       = true
	subnet_ipv6_cidr  = "%[9]s"
	insane_mode       = true
	insane_mode_az    = "%[10]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID4"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_INSANE_MODE_SUBNET"),
		os.Getenv("AWS_SUBNET_IPV6_CIDR"), os.Getenv("AWS_AVAILABILITY_ZONE"))
}

func testAccSpokeGatewayConfigGCPIPv6(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_spoke_gateway" "test_spoke_gateway_ipv6_gcp" {
	cloud_type   = 4
	account_name = aviatrix_account.test_acc_gcp.account_name
	gw_name      = "tfg-gcp-ipv6-%[1]s"
	vpc_id       = "%[4]s"
	vpc_reg      = "%[5]s"
	gw_size      = "n1-standard-1"
	subnet       = "%[6]s"
	enable_ipv6  = true
}
	`, rName, os.Getenv("GCP_PROJECT_ID"), os.Getenv("GOOGLE_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), os.Getenv("GCP_SUBNET"))
}
