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

func TestAccAviatrixTransitInstance_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)

	skipInstance := os.Getenv("SKIP_TRANSIT_INSTANCE")
	if skipInstance == "yes" {
		t.Skip("Skipping Transit instance test as SKIP_TRANSIT_INSTANCE is set")
	}

	skipInstanceAWS := os.Getenv("SKIP_TRANSIT_INSTANCE_AWS")
	skipInstanceAZURE := os.Getenv("SKIP_TRANSIT_INSTANCE_AZURE")
	skipInstanceGCP := os.Getenv("SKIP_TRANSIT_INSTANCE_GCP")
	skipInstanceOCI := os.Getenv("SKIP_TRANSIT_INSTANCE_OCI")

	if skipInstanceAWS == "yes" && skipInstanceAZURE == "yes" && skipInstanceGCP == "yes" && skipInstanceOCI == "yes" {
		t.Skip("Skipping Transit instance test as SKIP_TRANSIT_INSTANCE_AWS, SKIP_TRANSIT_INSTANCE_AZURE, " +
			"SKIP_TRANSIT_INSTANCE_GCP and SKIP_TRANSIT_INSTANCE_OCI are all set")
	}

	if skipInstanceAWS != "yes" {
		resourceNameAws := "aviatrix_transit_instance.test_transit_instance_aws"
		msgCommonAws := ". Set SKIP_TRANSIT_INSTANCE_AWS to yes to skip Transit Instance tests in AWS"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommonAws)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitInstanceDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitInstanceConfigBasicAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitInstanceExists(resourceNameAws, &gateway),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_name", fmt.Sprintf("tfi-aws-%s", rName)),
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
		t.Log("Skipping Transit instance test in AWS as SKIP_TRANSIT_INSTANCE_AWS is set")
	}

	if skipInstanceAZURE != "yes" {
		resourceNameAzure := "aviatrix_transit_instance.test_transit_instance_azure"
		msgCommonAzure := ". Set SKIP_TRANSIT_INSTANCE_AZURE to yes to skip Transit Instance tests in Azure"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckAZURE(t, msgCommonAzure)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitInstanceDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitInstanceConfigBasicAZURE(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitInstanceExists(resourceNameAzure, &gateway),
						resource.TestCheckResourceAttr(resourceNameAzure, "gw_name", fmt.Sprintf("tfi-azure-%s", rName)),
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
		t.Log("Skipping Transit instance test in Azure as SKIP_TRANSIT_INSTANCE_AZURE is set")
	}

	if skipInstanceGCP != "yes" {
		resourceNameGCP := "aviatrix_transit_instance.test_transit_instance_gcp"
		msgCommonGCP := ". Set SKIP_TRANSIT_INSTANCE_GCP to yes to skip Transit Instance tests in GCP"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckGCP(t, msgCommonGCP)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitInstanceDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitInstanceConfigBasicGCP(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitInstanceExists(resourceNameGCP, &gateway),
						resource.TestCheckResourceAttr(resourceNameGCP, "gw_name", fmt.Sprintf("tfi-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameGCP, "gw_size", os.Getenv("GCP_GW_SIZE")),
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
		t.Log("Skipping Transit instance test in GCP as SKIP_TRANSIT_INSTANCE_GCP is set")
	}

	if skipInstanceOCI != "yes" {
		resourceNameOCI := "aviatrix_transit_instance.test_transit_instance_oci"
		msgCommonOCI := ". Set SKIP_TRANSIT_INSTANCE_OCI to yes to skip Transit Instance tests in OCI"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckOCI(t, msgCommonOCI)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitInstanceDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitInstanceConfigBasicOCI(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitInstanceExists(resourceNameOCI, &gateway),
						resource.TestCheckResourceAttr(resourceNameOCI, "gw_name", fmt.Sprintf("tfi-oci-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameOCI, "gw_size", os.Getenv("OCI_GW_SIZE")),
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
		t.Log("Skipping Transit instance test in OCI as SKIP_TRANSIT_INSTANCE_OCI is set")
	}
}

func TestAccAviatrixTransitInstance_withTags(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_instance.test_transit_instance_tags"

	skipInstance := os.Getenv("SKIP_TRANSIT_INSTANCE")
	if skipInstance == "yes" {
		t.Skip("Skipping Transit instance test as SKIP_TRANSIT_INSTANCE is set")
	}

	skipInstanceAWS := os.Getenv("SKIP_TRANSIT_INSTANCE_AWS")
	if skipInstanceAWS == "yes" {
		t.Skip("Skipping Transit instance tags test as SKIP_TRANSIT_INSTANCE_AWS is set")
	}

	msgCommon := ". Set SKIP_TRANSIT_INSTANCE_AWS to yes to skip Transit Instance tests in AWS"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitInstanceConfigWithTags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitInstanceExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfi-tags-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Project", "transit-instance"),
				),
			},
			{
				Config: testAccTransitInstanceConfigWithTagsUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitInstanceExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfi-tags-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "production"),
					resource.TestCheckResourceAttr(resourceName, "tags.Team", "networking"),
				),
			},
		},
	})
}

func TestAccAviatrixTransitInstance_withRoutes(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_instance.test_transit_instance_routes"

	skipInstance := os.Getenv("SKIP_TRANSIT_INSTANCE")
	if skipInstance == "yes" {
		t.Skip("Skipping Transit instance test as SKIP_TRANSIT_INSTANCE is set")
	}

	skipInstanceAWS := os.Getenv("SKIP_TRANSIT_INSTANCE_AWS")
	if skipInstanceAWS == "yes" {
		t.Skip("Skipping Transit instance routes test as SKIP_TRANSIT_INSTANCE_AWS is set")
	}

	msgCommon := ". Set SKIP_TRANSIT_INSTANCE_AWS to yes to skip Transit Instance tests in AWS"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitInstanceConfigWithRoutes(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitInstanceExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfi-routes-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "customized_spoke_vpc_routes", "10.0.0.0/16,10.1.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "filtered_spoke_vpc_routes", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "excluded_advertised_spoke_routes", "172.16.0.0/16"),
				),
			},
		},
	})
}

func TestAccAviatrixTransitInstance_withBgpManualSpokeAdvertiseCidrs(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_instance.test_transit_instance_bgp"

	skipInstance := os.Getenv("SKIP_TRANSIT_INSTANCE")
	if skipInstance == "yes" {
		t.Skip("Skipping Transit instance test as SKIP_TRANSIT_INSTANCE is set")
	}

	skipInstanceAWS := os.Getenv("SKIP_TRANSIT_INSTANCE_AWS")
	if skipInstanceAWS == "yes" {
		t.Skip("Skipping Transit instance BGP test as SKIP_TRANSIT_INSTANCE_AWS is set")
	}

	msgCommon := ". Set SKIP_TRANSIT_INSTANCE_AWS to yes to skip Transit Instance tests in AWS"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitInstanceConfigWithBgpManualSpokeAdvertiseCidrs(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitInstanceExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfi-bgp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "bgp_manual_spoke_advertise_cidrs", "10.10.0.0/16,10.20.0.0/16"),
				),
			},
		},
	})
}

func TestAccAviatrixTransitInstance_withFireNet(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_instance.test_transit_instance_firenet"

	skipInstance := os.Getenv("SKIP_TRANSIT_INSTANCE")
	if skipInstance == "yes" {
		t.Skip("Skipping Transit instance test as SKIP_TRANSIT_INSTANCE is set")
	}

	skipInstanceAWS := os.Getenv("SKIP_TRANSIT_INSTANCE_AWS")
	if skipInstanceAWS == "yes" {
		t.Skip("Skipping Transit instance FireNet test as SKIP_TRANSIT_INSTANCE_AWS is set")
	}

	msgCommon := ". Set SKIP_TRANSIT_INSTANCE_AWS to yes to skip Transit Instance tests in AWS"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitInstanceConfigWithFireNet(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitInstanceExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfi-firenet-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "true"),
				),
			},
		},
	})
}

func TestAccAviatrixTransitInstance_withTransitFireNet(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_instance.test_transit_instance_transit_firenet"

	skipInstance := os.Getenv("SKIP_TRANSIT_INSTANCE")
	if skipInstance == "yes" {
		t.Skip("Skipping Transit instance test as SKIP_TRANSIT_INSTANCE is set")
	}

	skipInstanceAWS := os.Getenv("SKIP_TRANSIT_INSTANCE_AWS")
	if skipInstanceAWS == "yes" {
		t.Skip("Skipping Transit instance Transit FireNet test as SKIP_TRANSIT_INSTANCE_AWS is set")
	}

	msgCommon := ". Set SKIP_TRANSIT_INSTANCE_AWS to yes to skip Transit Instance tests in AWS"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitInstanceConfigWithTransitFireNet(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitInstanceExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfi-tfirenet-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "true"),
				),
			},
		},
	})
}

func testAccTransitInstanceConfigBasicAWS(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_aws" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfi-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitInstanceConfigBasicAZURE(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_azure" {
	account_name        = "tfa-azure-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_azure" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc_azure.account_name
	gw_name      = "tfi-azure-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
	zone         = "az-1"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"))
}

func testAccTransitInstanceConfigBasicGCP(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_gcp" {
	cloud_type   = 4
	account_name = aviatrix_account.test_acc_gcp.account_name
	gw_name      = "tfi-gcp-%[1]s"
	vpc_id       = "%[4]s"
	vpc_reg      = "%[5]s"
	gw_size      = "%[6]s"
	subnet       = "%[7]s"
}
	`, rName, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"),
		os.Getenv("GCP_GW_SIZE"), os.Getenv("GCP_SUBNET"))
}

func testAccTransitInstanceConfigBasicOCI(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_oci" {
	account_name                 = "tfa-oci-%s"
	cloud_type                   = 16
	oci_tenancy_id               = "%s"
	oci_user_id                  = "%s"
	oci_compartment_id           = "%s"
	oci_api_private_key_filepath = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_oci" {
	cloud_type          = 16
	account_name        = aviatrix_account.test_acc_oci.account_name
	gw_name             = "tfi-oci-%[1]s"
	vpc_id              = "%[6]s"
	vpc_reg             = "%[7]s"
	gw_size             = "%[8]s"
	subnet              = "%[9]s"
	availability_domain = "%[10]s"
	fault_domain        = "%[11]s"
}
	`, rName, os.Getenv("OCI_TENANCY_ID"), os.Getenv("OCI_USER_ID"),
		os.Getenv("OCI_COMPARTMENT_ID"), os.Getenv("OCI_API_PRIVATE_KEY_FILEPATH"),
		os.Getenv("OCI_VPC_ID"), os.Getenv("OCI_REGION"),
		os.Getenv("OCI_GW_SIZE"), os.Getenv("OCI_SUBNET"),
		os.Getenv("OCI_AVAILABILITY_DOMAIN"), os.Getenv("OCI_FAULT_DOMAIN"))
}

func testAccTransitInstanceConfigWithTags(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_tags" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfi-tags-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"

	tags = {
		Environment = "test"
		Project     = "transit-instance"
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitInstanceConfigWithTagsUpdated(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_tags" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfi-tags-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"

	tags = {
		Environment = "production"
		Team        = "networking"
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitInstanceConfigWithRoutes(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_routes" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfi-routes-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"

	customized_spoke_vpc_routes        = "10.0.0.0/16,10.1.0.0/16"
	filtered_spoke_vpc_routes          = "192.168.0.0/16"
	excluded_advertised_spoke_routes   = "172.16.0.0/16"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitInstanceConfigWithBgpManualSpokeAdvertiseCidrs(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_bgp" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfi-bgp-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"

	bgp_manual_spoke_advertise_cidrs = "10.10.0.0/16,10.20.0.0/16"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitInstanceConfigWithFireNet(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_firenet" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfi-firenet-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "c5.xlarge"
	subnet       = "%[7]s"

	enable_firenet = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitInstanceConfigWithTransitFireNet(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_instance" "test_transit_instance_transit_firenet" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfi-tfirenet-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "c5.xlarge"
	subnet       = "%[7]s"

	enable_transit_firenet = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccCheckTransitInstanceExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transit instance not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit instance ID is set")
		}

		client := mustClient(testAccProvider.Meta())

		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}
		gw, err := client.GetGateway(foundGateway)
		if err != nil {
			return err
		}
		if gw.GwName != rs.Primary.ID {
			return fmt.Errorf("transit instance not found")
		}

		*gateway = *gw
		return nil
	}
}

func testAccCheckTransitInstanceDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_instance" {
			continue
		}

		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetGateway(foundGateway)
		if err == nil {
			return fmt.Errorf("transit instance still exists")
		}
	}

	return nil
}
