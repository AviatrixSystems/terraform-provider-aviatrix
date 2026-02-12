package aviatrix

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func preSpokeInstanceCheck(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"AWS_VPC_ID",
		"AWS_SUBNET",
		"AWS_REGION",
		"AWS_ACCOUNT_NUMBER",
		"AWS_ACCESS_KEY",
		"AWS_SECRET_KEY",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func TestAccAviatrixSpokeInstance_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_instance.test"

	msgCommon := ". Set SKIP_SPOKE_INSTANCE to yes to skip Spoke Instance tests"

	skipSpokeInstance := os.Getenv("SKIP_SPOKE_INSTANCE")
	if skipSpokeInstance == "yes" {
		t.Skip("Skipping Spoke Instance test as SKIP_SPOKE_INSTANCE is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeInstanceConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "group_uuid"),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-spoke-gw-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "allocate_new_eip", "true"),
					resource.TestCheckResourceAttr(resourceName, "single_az_ha", "true"),
					resource.TestCheckResourceAttr(resourceName, "insane_mode", "false"),
					// Computed fields
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAviatrixSpokeInstance_update(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_instance.test"

	msgCommon := ". Set SKIP_SPOKE_INSTANCE to yes to skip Spoke Instance tests"

	skipSpokeInstance := os.Getenv("SKIP_SPOKE_INSTANCE")
	if skipSpokeInstance == "yes" {
		t.Skip("Skipping Spoke Instance test as SKIP_SPOKE_INSTANCE is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	awsGwSizeUpdated := os.Getenv("AWS_GW_SIZE_UPDATED")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}
	if awsGwSizeUpdated == "" {
		awsGwSizeUpdated = "t3.small"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with initial size
			{
				Config: testAccSpokeInstanceConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSize),
				),
			},
			// Step 2: Update size
			{
				Config: testAccSpokeInstanceConfigUpdateSize(rName, awsGwSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSizeUpdated),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeInstance_routeFeatures(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_instance.test"

	msgCommon := ". Set SKIP_SPOKE_INSTANCE to yes to skip Spoke Instance tests"

	skipSpokeInstance := os.Getenv("SKIP_SPOKE_INSTANCE")
	if skipSpokeInstance == "yes" {
		t.Skip("Skipping Spoke Instance test as SKIP_SPOKE_INSTANCE is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeInstanceConfigRouteFeatures(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_private_vpc_default_route", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_skip_public_route_table_update", "true"),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeInstance_filteredRoutes(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_instance.test"

	msgCommon := ". Set SKIP_SPOKE_INSTANCE to yes to skip Spoke Instance tests"

	skipSpokeInstance := os.Getenv("SKIP_SPOKE_INSTANCE")
	if skipSpokeInstance == "yes" {
		t.Skip("Skipping Spoke Instance test as SKIP_SPOKE_INSTANCE is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with filtered routes
			{
				Config: testAccSpokeInstanceConfigFilteredRoutes(rName, awsGwSize, "10.0.0.0/16,192.168.0.0/16"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "filtered_spoke_vpc_routes", "10.0.0.0/16,192.168.0.0/16"),
				),
			},
			// Step 2: Update filtered routes
			{
				Config: testAccSpokeInstanceConfigFilteredRoutes(rName, awsGwSize, "10.0.0.0/8"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "filtered_spoke_vpc_routes", "10.0.0.0/8"),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeInstance_tags(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_instance.test"

	msgCommon := ". Set SKIP_SPOKE_INSTANCE to yes to skip Spoke Instance tests"

	skipSpokeInstance := os.Getenv("SKIP_SPOKE_INSTANCE")
	if skipSpokeInstance == "yes" {
		t.Skip("Skipping Spoke Instance test as SKIP_SPOKE_INSTANCE is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with tags
			{
				Config: testAccSpokeInstanceConfigWithTags(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Team", "terraform"),
				),
			},
			// Step 2: Update tags
			{
				Config: testAccSpokeInstanceConfigWithTagsUpdated(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "production"),
					resource.TestCheckResourceAttr(resourceName, "tags.Team", "aviatrix"),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeInstance_tunnelDetection(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_instance.test"

	msgCommon := ". Set SKIP_SPOKE_INSTANCE to yes to skip Spoke Instance tests"

	skipSpokeInstance := os.Getenv("SKIP_SPOKE_INSTANCE")
	if skipSpokeInstance == "yes" {
		t.Skip("Skipping Spoke Instance test as SKIP_SPOKE_INSTANCE is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with tunnel detection time
			{
				Config: testAccSpokeInstanceConfigTunnelDetection(rName, awsGwSize, 60),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tunnel_detection_time", "60"),
				),
			},
			// Step 2: Update tunnel detection time
			{
				Config: testAccSpokeInstanceConfigTunnelDetection(rName, awsGwSize, 120),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tunnel_detection_time", "120"),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeInstance_encryption(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_instance.test"

	msgCommon := ". Set SKIP_SPOKE_INSTANCE to yes to skip Spoke Instance tests"

	skipSpokeInstance := os.Getenv("SKIP_SPOKE_INSTANCE")
	if skipSpokeInstance == "yes" {
		t.Skip("Skipping Spoke Instance test as SKIP_SPOKE_INSTANCE is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with encryption disabled (default)
			{
				Config: testAccSpokeInstanceConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_encrypt_volume", "false"),
				),
			},
			// Step 2: Create with encryption enabled (forces replacement)
			{
				Config: testAccSpokeInstanceConfigWithEncryption(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_encrypt_volume", "true"),
				),
			},
		},
	})
}

// ============================================================================
// Config Generators
// ============================================================================

func testAccSpokeInstanceConfigBasic(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_group" "test" {
	group_name          = "tfg-spoke-group-%[1]s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	vpc_region          = "%[7]s"
}

resource "aviatrix_spoke_instance" "test" {
	group_uuid       = aviatrix_spoke_group.test.group_uuid
	gw_name          = "tfg-spoke-gw-%[1]s"
	subnet           = "%[8]s"
	gw_size          = "%[5]s"
	allocate_new_eip = true
	single_az_ha     = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccSpokeInstanceConfigUpdateSize(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_group" "test" {
	group_name          = "tfg-spoke-group-%[1]s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	vpc_region          = "%[7]s"
}

resource "aviatrix_spoke_instance" "test" {
	group_uuid       = aviatrix_spoke_group.test.group_uuid
	gw_name          = "tfg-spoke-gw-%[1]s"
	subnet           = "%[8]s"
	gw_size          = "%[5]s"
	allocate_new_eip = true
	single_az_ha     = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccSpokeInstanceConfigRouteFeatures(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_group" "test" {
	group_name          = "tfg-spoke-group-%[1]s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	vpc_region          = "%[7]s"
}

resource "aviatrix_spoke_instance" "test" {
	group_uuid                            = aviatrix_spoke_group.test.group_uuid
	gw_name                               = "tfg-spoke-gw-%[1]s"
	subnet                                = "%[8]s"
	gw_size                               = "%[5]s"
	allocate_new_eip                      = true
	single_az_ha                          = true
	enable_private_vpc_default_route      = true
	enable_skip_public_route_table_update = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccSpokeInstanceConfigFilteredRoutes(rName, gwSize, filteredRoutes string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_group" "test" {
	group_name          = "tfg-spoke-group-%[1]s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	vpc_region          = "%[7]s"
}

resource "aviatrix_spoke_instance" "test" {
	group_uuid                = aviatrix_spoke_group.test.group_uuid
	gw_name                   = "tfg-spoke-gw-%[1]s"
	subnet                    = "%[8]s"
	gw_size                   = "%[5]s"
	allocate_new_eip          = true
	single_az_ha              = true
	filtered_spoke_vpc_routes = "%[9]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), filteredRoutes)
}

func testAccSpokeInstanceConfigWithTags(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_group" "test" {
	group_name          = "tfg-spoke-group-%[1]s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	vpc_region          = "%[7]s"
}

resource "aviatrix_spoke_instance" "test" {
	group_uuid       = aviatrix_spoke_group.test.group_uuid
	gw_name          = "tfg-spoke-gw-%[1]s"
	subnet           = "%[8]s"
	gw_size          = "%[5]s"
	allocate_new_eip = true
	single_az_ha     = true

	tags = {
		Environment = "test"
		Team        = "terraform"
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccSpokeInstanceConfigWithTagsUpdated(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_group" "test" {
	group_name          = "tfg-spoke-group-%[1]s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	vpc_region          = "%[7]s"
}

resource "aviatrix_spoke_instance" "test" {
	group_uuid       = aviatrix_spoke_group.test.group_uuid
	gw_name          = "tfg-spoke-gw-%[1]s"
	subnet           = "%[8]s"
	gw_size          = "%[5]s"
	allocate_new_eip = true
	single_az_ha     = true

	tags = {
		Environment = "production"
		Team        = "aviatrix"
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccSpokeInstanceConfigTunnelDetection(rName, gwSize string, tunnelDetectionTime int) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_group" "test" {
	group_name          = "tfg-spoke-group-%[1]s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	vpc_region          = "%[7]s"
}

resource "aviatrix_spoke_instance" "test" {
	group_uuid            = aviatrix_spoke_group.test.group_uuid
	gw_name               = "tfg-spoke-gw-%[1]s"
	subnet                = "%[8]s"
	gw_size               = "%[5]s"
	allocate_new_eip      = true
	single_az_ha          = true
	tunnel_detection_time = %[9]d
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), tunnelDetectionTime)
}

func testAccSpokeInstanceConfigWithEncryption(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_group" "test" {
	group_name          = "tfg-spoke-group-%[1]s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	vpc_region          = "%[7]s"
}

resource "aviatrix_spoke_instance" "test" {
	group_uuid            = aviatrix_spoke_group.test.group_uuid
	gw_name               = "tfg-spoke-gw-%[1]s"
	subnet                = "%[8]s"
	gw_size               = "%[5]s"
	allocate_new_eip      = true
	single_az_ha          = true
	enable_encrypt_volume = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

// ============================================================================
// Edge Spoke Instance Tests
// ============================================================================

func preEdgeSpokeInstanceCheck(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"EDGE_SITE_ID",
		"EDGE_ZTP_FILE_DOWNLOAD_PATH",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func TestAccAviatrixSpokeInstance_edgeBasic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_instance.test_edge"

	msgCommon := ". Set SKIP_EDGE_SPOKE_INSTANCE to yes to skip Edge Spoke Instance tests"

	skipEdgeSpokeInstance := os.Getenv("SKIP_EDGE_SPOKE_INSTANCE")
	if skipEdgeSpokeInstance == "yes" {
		t.Skip("Skipping Edge Spoke Instance test as SKIP_EDGE_SPOKE_INSTANCE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preEdgeSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeInstanceConfigEdgeBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "group_uuid"),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-edge-spoke-%s", rName)),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeInstance_edgeUpdate(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_instance.test_edge"

	msgCommon := ". Set SKIP_EDGE_SPOKE_INSTANCE to yes to skip Edge Spoke Instance tests"

	skipEdgeSpokeInstance := os.Getenv("SKIP_EDGE_SPOKE_INSTANCE")
	if skipEdgeSpokeInstance == "yes" {
		t.Skip("Skipping Edge Spoke Instance test as SKIP_EDGE_SPOKE_INSTANCE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preEdgeSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create edge spoke instance
			{
				Config: testAccSpokeInstanceConfigEdgeBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
				),
			},
			// Step 2: Update interfaces
			{
				Config: testAccSpokeInstanceConfigEdgeUpdateInterfaces(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(resourceName),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeInstance_haGateway(t *testing.T) {
	rName := acctest.RandString(5)
	primaryResourceName := "aviatrix_spoke_instance.test_primary"
	haResourceName := "aviatrix_spoke_instance.test_ha"

	msgCommon := ". Set SKIP_SPOKE_INSTANCE to yes to skip Spoke Instance tests"

	skipSpokeInstance := os.Getenv("SKIP_SPOKE_INSTANCE")
	if skipSpokeInstance == "yes" {
		t.Skip("Skipping Spoke Instance test as SKIP_SPOKE_INSTANCE is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	haSubnet := os.Getenv("AWS_HA_SUBNET")
	if haSubnet == "" {
		haSubnet = os.Getenv("AWS_SUBNET") // Use same subnet if HA subnet not specified
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create primary and HA spoke instances
			{
				Config: testAccSpokeInstanceConfigWithHA(rName, awsGwSize, haSubnet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(primaryResourceName),
					testAccCheckSpokeInstanceExists(haResourceName),
					resource.TestCheckResourceAttr(primaryResourceName, "gw_name", fmt.Sprintf("tfg-spoke-primary-%s", rName)),
					resource.TestCheckResourceAttr(haResourceName, "gw_name", fmt.Sprintf("tfg-spoke-ha-%s", rName)),
				),
			},
		},
	})
}

// ============================================================================
// Edge Config Generators
// ============================================================================

func testAccSpokeInstanceConfigEdgeBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_spoke_group" "test_edge" {
	group_name = "tfg-edge-spoke-group-%s"
	cloud_type = 65536  # Edge Self-managed
	gw_type    = "spoke"
	vpc_id     = "%s"
}

resource "aviatrix_spoke_instance" "test_edge" {
	group_uuid             = aviatrix_spoke_group.test_edge.group_uuid
	gw_name                = "tfg-edge-spoke-%[1]s"
	ztp_file_download_path = "%[3]s"
	ztp_file_type          = "cloud-init"

	interfaces {
		logical_ifname = "wan0"
		type           = "WAN"
		dhcp           = true
	}

	interfaces {
		logical_ifname = "lan0"
		type           = "LAN"
		ip_address     = "10.220.5.1/24"
	}

	interfaces {
		logical_ifname = "mgmt0"
		type           = "MANAGEMENT"
		dhcp           = true
	}
}
	`, rName, os.Getenv("EDGE_SITE_ID"), os.Getenv("EDGE_ZTP_FILE_DOWNLOAD_PATH"))
}

func testAccSpokeInstanceConfigEdgeUpdateInterfaces(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_spoke_group" "test_edge" {
	group_name = "tfg-edge-spoke-group-%s"
	cloud_type = 65536  # Edge Self-managed
	gw_type    = "spoke"
	vpc_id     = "%s"
}

resource "aviatrix_spoke_instance" "test_edge" {
	group_uuid             = aviatrix_spoke_group.test_edge.group_uuid
	gw_name                = "tfg-edge-spoke-%[1]s"
	ztp_file_download_path = "%[3]s"
	ztp_file_type          = "cloud-init"

	interfaces {
		logical_ifname = "wan0"
		type           = "WAN"
		dhcp           = false
		ip_address     = "192.168.1.10/24"
		gateway_ip     = "192.168.1.1"
	}

	interfaces {
		logical_ifname = "lan0"
		type           = "LAN"
		ip_address     = "10.220.10.1/24"
	}

	interfaces {
		logical_ifname = "mgmt0"
		type           = "MANAGEMENT"
		dhcp           = true
	}

	management_egress_ip_prefix_list = ["10.0.0.0/8"]
}
	`, rName, os.Getenv("EDGE_SITE_ID"), os.Getenv("EDGE_ZTP_FILE_DOWNLOAD_PATH"))
}

func testAccSpokeInstanceConfigWithHA(rName, gwSize, haSubnet string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_group" "test" {
	group_name          = "tfg-spoke-group-%[1]s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	vpc_region          = "%[7]s"
}

resource "aviatrix_spoke_instance" "test_primary" {
	group_uuid       = aviatrix_spoke_group.test.group_uuid
	gw_name          = "tfg-spoke-primary-%[1]s"
	subnet           = "%[8]s"
	gw_size          = "%[5]s"
	allocate_new_eip = true
	single_az_ha     = true
}

resource "aviatrix_spoke_instance" "test_ha" {
	depends_on       = [aviatrix_spoke_instance.test_primary]
	group_uuid       = aviatrix_spoke_group.test.group_uuid
	gw_name          = "tfg-spoke-ha-%[1]s"
	subnet           = "%[9]s"
	gw_size          = "%[5]s"
	allocate_new_eip = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), haSubnet)
}

func TestAccAviatrixSpokeInstance_edgeHA(t *testing.T) {
	rName := acctest.RandString(5)
	primaryResourceName := "aviatrix_spoke_instance.test_edge_primary"
	haResourceName := "aviatrix_spoke_instance.test_edge_ha"

	msgCommon := ". Set SKIP_EDGE_SPOKE_INSTANCE to yes to skip Edge Spoke Instance tests"

	skipEdgeSpokeInstance := os.Getenv("SKIP_EDGE_SPOKE_INSTANCE")
	if skipEdgeSpokeInstance == "yes" {
		t.Skip("Skipping Edge Spoke Instance test as SKIP_EDGE_SPOKE_INSTANCE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preEdgeSpokeInstanceCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeInstanceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create primary and HA edge spoke instances
			{
				Config: testAccSpokeInstanceConfigEdgeWithHA(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeInstanceExists(primaryResourceName),
					testAccCheckSpokeInstanceExists(haResourceName),
					resource.TestCheckResourceAttr(primaryResourceName, "gw_name", fmt.Sprintf("tfg-edge-spoke-primary-%s", rName)),
					resource.TestCheckResourceAttr(haResourceName, "gw_name", fmt.Sprintf("tfg-edge-spoke-ha-%s", rName)),
				),
			},
		},
	})
}

func testAccSpokeInstanceConfigEdgeWithHA(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_spoke_group" "test_edge" {
	group_name = "tfg-edge-spoke-group-%s"
	cloud_type = 65536  # Edge Self-managed
	gw_type    = "spoke"
	vpc_id     = "%s"
}

resource "aviatrix_spoke_instance" "test_edge_primary" {
	group_uuid             = aviatrix_spoke_group.test_edge.group_uuid
	gw_name                = "tfg-edge-spoke-primary-%[1]s"
	ztp_file_download_path = "%[3]s"
	ztp_file_type          = "cloud-init"

	interfaces {
		logical_ifname = "wan0"
		type           = "WAN"
		dhcp           = true
	}

	interfaces {
		logical_ifname = "lan0"
		type           = "LAN"
		ip_address     = "10.220.5.1/24"
	}

	interfaces {
		logical_ifname = "mgmt0"
		type           = "MANAGEMENT"
		dhcp           = true
	}
}

resource "aviatrix_spoke_instance" "test_edge_ha" {
	depends_on             = [aviatrix_spoke_instance.test_edge_primary]
	group_uuid             = aviatrix_spoke_group.test_edge.group_uuid
	gw_name                = "tfg-edge-spoke-ha-%[1]s"
	ztp_file_download_path = "%[3]s"
	ztp_file_type          = "cloud-init"

	interfaces {
		logical_ifname = "wan0"
		type           = "WAN"
		dhcp           = true
	}

	interfaces {
		logical_ifname = "lan0"
		type           = "LAN"
		ip_address     = "10.220.6.1/24"
	}

	interfaces {
		logical_ifname = "mgmt0"
		type           = "MANAGEMENT"
		dhcp           = true
	}
}
	`, rName, os.Getenv("EDGE_SITE_ID"), os.Getenv("EDGE_ZTP_FILE_DOWNLOAD_PATH"))
}

// ============================================================================
// Helper Functions
// ============================================================================

func testAccCheckSpokeInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke instance not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke instance ID is set")
		}

		client, ok := testAccProvider.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to get client")
		}

		_, err := client.GetGateway(&goaviatrix.Gateway{GwName: rs.Primary.ID})
		if err != nil {
			return fmt.Errorf("spoke instance not found: %w", err)
		}

		return nil
	}
}

func testAccCheckSpokeInstanceDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to get client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_instance" {
			continue
		}

		_, err := client.GetGateway(&goaviatrix.Gateway{GwName: rs.Primary.ID})
		if err == nil {
			return fmt.Errorf("spoke instance still exists: %s", rs.Primary.ID)
		}

		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("error checking spoke instance destroyed: %w", err)
		}
	}

	return nil
}
