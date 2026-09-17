package aviatrix

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// configSubnet builds a "subnets" list element as Terraform hands it to the
// resource, mirroring the schema keys resourceAviatrixVpc declares.
func configSubnet(region, cidr, name string) map[string]any {
	return map[string]any{
		"region":           region,
		"cidr":             cidr,
		"name":             name,
		"subnet_id":        "",
		"ipv6_cidr":        "",
		"ipv6_access_type": "",
	}
}

// TestBuildSubnetsForStatePreservesConfigOrder is a regression test for
// AVX-75737: a multi-region GCP VPC showed a spurious force-replacement on
// every plan because the controller returns subnets grouped by region rather
// than in configured order. buildSubnetsForState must echo configured subnets
// back in configured order so the TypeList indices match and Terraform sees no
// diff.
func TestBuildSubnetsForStatePreservesConfigOrder(t *testing.T) {
	// Configuration declares us-east1 subnets first, us-west1 second.
	configured := []any{
		configSubnet("us-east1", "172.16.1.0/24", "subnet-us-east1-1"),
		configSubnet("us-east1", "172.16.2.0/24", "subnet-us-east1-2"),
		configSubnet("us-west1", "172.16.3.0/24", "subnet-us-west1-1"),
		configSubnet("us-west1", "172.16.4.0/24", "subnet-us-west1-2"),
	}

	// Controller returns them grouped with us-west1 first (a different order).
	apiSubnets := []goaviatrix.SubnetInfo{
		{Region: "us-west1", Cidr: "172.16.3.0/24", Name: "subnet-us-west1-1"},
		{Region: "us-west1", Cidr: "172.16.4.0/24", Name: "subnet-us-west1-2"},
		{Region: "us-east1", Cidr: "172.16.1.0/24", Name: "subnet-us-east1-1"},
		{Region: "us-east1", Cidr: "172.16.2.0/24", Name: "subnet-us-east1-2"},
	}

	got := buildSubnetsForState(goaviatrix.GCP, apiSubnets, configured)

	gotNames := make([]string, len(got))
	for i, s := range got {
		gotNames[i] = mustString(s["name"])
	}
	wantNames := []string{
		"subnet-us-east1-1",
		"subnet-us-east1-2",
		"subnet-us-west1-1",
		"subnet-us-west1-2",
	}
	assert.Equal(t, wantNames, gotNames, "subnets should be returned in configured order")

	// The region and cidr must travel with the correct subnet.
	assert.Equal(t, "us-east1", got[0]["region"])
	assert.Equal(t, "172.16.1.0/24", got[0]["cidr"])
	assert.Equal(t, "us-west1", got[2]["region"])
	assert.Equal(t, "172.16.3.0/24", got[2]["cidr"])
}

// TestBuildSubnetsForStateImportUsesAPIOrder covers the import case where there
// is no prior configuration: all subnets should still be emitted, in the
// controller's order.
func TestBuildSubnetsForStateImportUsesAPIOrder(t *testing.T) {
	apiSubnets := []goaviatrix.SubnetInfo{
		{Region: "us-west1", Cidr: "172.16.3.0/24", Name: "subnet-us-west1-1"},
		{Region: "us-east1", Cidr: "172.16.1.0/24", Name: "subnet-us-east1-1"},
	}

	got := buildSubnetsForState(goaviatrix.GCP, apiSubnets, nil)

	gotNames := make([]string, len(got))
	for i, s := range got {
		gotNames[i] = mustString(s["name"])
	}
	assert.Equal(t, []string{"subnet-us-west1-1", "subnet-us-east1-1"}, gotNames)
}

// TestBuildSubnetsForStateGCPExcludesSubnetID verifies GCP state elements omit
// the AWS-only subnet_id key (the mismatch on this key was the original cause
// of the config-order logic silently never matching).
func TestBuildSubnetsForStateGCPExcludesSubnetID(t *testing.T) {
	configured := []any{configSubnet("us-east1", "172.16.1.0/24", "subnet-1")}
	apiSubnets := []goaviatrix.SubnetInfo{
		{Region: "us-east1", Cidr: "172.16.1.0/24", Name: "subnet-1", SubnetID: "ignored-for-gcp"},
	}

	got := buildSubnetsForState(goaviatrix.GCP, apiSubnets, configured)

	assert.Len(t, got, 1)
	_, hasSubnetID := got[0]["subnet_id"]
	assert.False(t, hasSubnetID, "GCP subnets must not carry subnet_id in state")
	assert.Contains(t, got[0], "region")
}

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

		client := mustClient(testAccProvider.Meta())

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
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpc" {
			continue
		}

		foundVpc := &goaviatrix.Vpc{
			Name: rs.Primary.Attributes["name"],
		}

		_, err := client.GetVpc(foundVpc)
		if !errors.Is(err, goaviatrix.ErrNotFound) {
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
