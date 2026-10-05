package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func preSpokeGroupCheck(t *testing.T, msgCommon string) {
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

func TestAccAviatrixSpokeGroup_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_group.test"

	msgCommon := ". Set SKIP_SPOKE_GROUP to yes to skip Spoke Group tests"

	skipSpokeGroup := os.Getenv("SKIP_SPOKE_GROUP")
	if skipSpokeGroup == "yes" {
		t.Skip("Skipping Spoke Group test as SKIP_SPOKE_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tfg-spoke-group-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "cloud_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "gw_type", "spoke"),
					resource.TestCheckResourceAttr(resourceName, "group_instance_size", awsGwSize),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "enable_nat", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_vpc_dns_server", "false"),
					resource.TestCheckResourceAttr(resourceName, "bgp_polling_time", "50"),
					resource.TestCheckResourceAttr(resourceName, "bgp_neighbor_status_polling_time", "5"),
					resource.TestCheckResourceAttr(resourceName, "bgp_hold_time", "180"),
					resource.TestCheckResourceAttr(resourceName, "learned_cidrs_approval_mode", "gateway"),
					// Computed fields
					resource.TestCheckResourceAttrSet(resourceName, "group_uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "explicitly_created"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "vendor_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gw_name",
				},
			},
		},
	})
}

func TestAccAviatrixSpokeGroup_update(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_group.test"

	msgCommon := ". Set SKIP_SPOKE_GROUP to yes to skip Spoke Group tests"

	skipSpokeGroup := os.Getenv("SKIP_SPOKE_GROUP")
	if skipSpokeGroup == "yes" {
		t.Skip("Skipping Spoke Group test as SKIP_SPOKE_GROUP is set")
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
			preSpokeGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with initial values
			{
				Config: testAccSpokeGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tfg-spoke-group-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "group_instance_size", awsGwSize),
					resource.TestCheckResourceAttr(resourceName, "enable_nat", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_vpc_dns_server", "false"),
					resource.TestCheckResourceAttr(resourceName, "bgp_polling_time", "50"),
					resource.TestCheckResourceAttr(resourceName, "bgp_neighbor_status_polling_time", "5"),
					resource.TestCheckResourceAttr(resourceName, "bgp_hold_time", "180"),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_ecmp", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_private_vpc_default_route", "false"),
					resource.TestCheckResourceAttr(resourceName, "learned_cidrs_approval_mode", "gateway"),
					// Computed fields
					resource.TestCheckResourceAttrSet(resourceName, "group_uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "explicitly_created"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "vendor_name"),
				),
			},
			// Step 2: Update gateway size and enable NAT
			{
				Config: testAccSpokeGroupConfigUpdateSize(rName, awsGwSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_instance_size", awsGwSizeUpdated),
					resource.TestCheckResourceAttr(resourceName, "enable_nat", "true"),
				),
			},
			// Step 3: Update BGP settings
			{
				Config: testAccSpokeGroupConfigUpdateBgp(rName, awsGwSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bgp_polling_time", "30"),
					resource.TestCheckResourceAttr(resourceName, "bgp_hold_time", "120"),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_ecmp", "true"),
				),
			},
			// Step 4: Update feature flags
			{
				Config: testAccSpokeGroupConfigUpdateFeatures(rName, awsGwSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_vpc_dns_server", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_private_vpc_default_route", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_skip_public_route_table_update", "true"),
				),
			},
			// Step 5: Import verification
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gw_name",
				},
			},
		},
	})
}

func TestAccAviatrixSpokeGroup_activeStandby(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_group.test"

	msgCommon := ". Set SKIP_SPOKE_GROUP to yes to skip Spoke Group tests"

	skipSpokeGroup := os.Getenv("SKIP_SPOKE_GROUP")
	if skipSpokeGroup == "yes" {
		t.Skip("Skipping Spoke Group test as SKIP_SPOKE_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with active-standby disabled
			{
				Config: testAccSpokeGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby_preemptive", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "explicitly_created"),
				),
			},
			// Step 2: Enable active-standby
			{
				Config: testAccSpokeGroupConfigActiveStandby(rName, awsGwSize, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby_preemptive", "false"),
				),
			},
			// Step 3: Enable active-standby with preemptive
			{
				Config: testAccSpokeGroupConfigActiveStandby(rName, awsGwSize, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby_preemptive", "true"),
				),
			},
			// Step 4: Disable active-standby
			{
				Config: testAccSpokeGroupConfigActiveStandby(rName, awsGwSize, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby_preemptive", "false"),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeGroup_learnedCidrsApproval(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_group.test"

	msgCommon := ". Set SKIP_SPOKE_GROUP to yes to skip Spoke Group tests"

	skipSpokeGroup := os.Getenv("SKIP_SPOKE_GROUP")
	if skipSpokeGroup == "yes" {
		t.Skip("Skipping Spoke Group test as SKIP_SPOKE_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Enable learned CIDRs approval
			{
				Config: testAccSpokeGroupConfigLearnedCidrs(rName, awsGwSize, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_learned_cidrs_approval", "true"),
					resource.TestCheckResourceAttr(resourceName, "learned_cidrs_approval_mode", "gateway"),
					resource.TestCheckResourceAttrSet(resourceName, "explicitly_created"),
				),
			},
			// Step 2: Disable learned CIDRs approval
			{
				Config: testAccSpokeGroupConfigLearnedCidrs(rName, awsGwSize, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_learned_cidrs_approval", "false"),
					resource.TestCheckResourceAttr(resourceName, "learned_cidrs_approval_mode", "gateway"),
					resource.TestCheckResourceAttrSet(resourceName, "explicitly_created"),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeGroup_bgp(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_group.test"

	msgCommon := ". Set SKIP_SPOKE_GROUP to yes to skip Spoke Group tests"

	skipSpokeGroup := os.Getenv("SKIP_SPOKE_GROUP")
	if skipSpokeGroup == "yes" {
		t.Skip("Skipping Spoke Group test as SKIP_SPOKE_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with BGP enabled and local_as_number
			{
				Config: testAccSpokeGroupConfigBgp(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp", "true"),
					resource.TestCheckResourceAttr(resourceName, "local_as_number", "65001"),
				),
			},
			// Step 2: Add prepend_as_path
			{
				Config: testAccSpokeGroupConfigBgpWithPrepend(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp", "true"),
					resource.TestCheckResourceAttr(resourceName, "local_as_number", "65001"),
					resource.TestCheckResourceAttr(resourceName, "prepend_as_path.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "prepend_as_path.0", "65001"),
					resource.TestCheckResourceAttr(resourceName, "prepend_as_path.1", "65001"),
				),
			},
		},
	})
}

func TestAccAviatrixSpokeGroup_featureFlags(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_group.test"

	msgCommon := ". Set SKIP_SPOKE_GROUP to yes to skip Spoke Group tests"

	skipSpokeGroup := os.Getenv("SKIP_SPOKE_GROUP")
	if skipSpokeGroup == "yes" {
		t.Skip("Skipping Spoke Group test as SKIP_SPOKE_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preSpokeGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with default feature flags (jumbo_frame and gro_gso default to true)
			{
				Config: testAccSpokeGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_jumbo_frame", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_gro_gso", "true"),
				),
			},
			// Step 2: Disable jumbo_frame and gro_gso
			{
				Config: testAccSpokeGroupConfigFeatureFlagsDisabled(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_jumbo_frame", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_gro_gso", "false"),
				),
			},
			// Step 3: Re-enable jumbo_frame and gro_gso
			{
				Config: testAccSpokeGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_jumbo_frame", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_gro_gso", "true"),
				),
			},
		},
	})
}

// ============================================================================
// Config Generators
// ============================================================================

func testAccSpokeGroupConfigBasic(rName, gwSize string) string {
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
	subnet              = "%[7]s"
	vpc_region          = "%[8]s"
	gw_name             = "tfg-spoke-group-gw-%[1]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"))
}

func testAccSpokeGroupConfigUpdateSize(rName, gwSize string) string {
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
	subnet              = "%[7]s"
	vpc_region          = "%[8]s"
	gw_name             = "tfg-spoke-group-gw-%[1]s"

	enable_nat          = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"))
}

func testAccSpokeGroupConfigUpdateBgp(rName, gwSize string) string {
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
	subnet              = "%[7]s"
	vpc_region          = "%[8]s"
	gw_name             = "tfg-spoke-group-gw-%[1]s"

	enable_nat          = true
	bgp_polling_time    = 30
	bgp_hold_time       = 120
	enable_bgp_ecmp     = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"))
}

func testAccSpokeGroupConfigUpdateFeatures(rName, gwSize string) string {
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
	subnet              = "%[7]s"
	vpc_region          = "%[8]s"
	gw_name             = "tfg-spoke-group-gw-%[1]s"

	enable_nat                          = true
	bgp_polling_time                    = 30
	bgp_hold_time                       = 120
	enable_bgp_ecmp                     = true
	enable_vpc_dns_server               = true
	enable_private_vpc_default_route    = true
	enable_skip_public_route_table_update = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"))
}

func testAccSpokeGroupConfigActiveStandby(rName, gwSize string, enableActiveStandby, enablePreemptive bool) string {
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
	group_name          = "tfg-spoke-group-%s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-spoke-group-gw-%s"

	enable_active_standby            = %t
	enable_active_standby_preemptive = %t
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName,
		enableActiveStandby, enablePreemptive)
}

func testAccSpokeGroupConfigLearnedCidrs(rName, gwSize string, enableLearnedCidrs bool) string {
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
	group_name          = "tfg-spoke-group-%s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-spoke-group-gw-%s"
	enable_bgp          = true
	local_as_number     = "65001"

	enable_learned_cidrs_approval = %t
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName,
		enableLearnedCidrs)
}

func testAccSpokeGroupConfigBgp(rName, gwSize string) string {
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
	group_name          = "tfg-spoke-group-%s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-spoke-group-gw-%s"

	enable_bgp          = true
	local_as_number     = "65001"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName)
}

func testAccSpokeGroupConfigBgpWithPrepend(rName, gwSize string) string {
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
	group_name          = "tfg-spoke-group-%s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-spoke-group-gw-%s"

	enable_bgp          = true
	local_as_number     = "65001"
	prepend_as_path     = ["65001", "65001"]
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName)
}

func testAccSpokeGroupConfigFeatureFlagsDisabled(rName, gwSize string) string {
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
	group_name          = "tfg-spoke-group-%s"
	cloud_type          = 1
	gw_type             = "spoke"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-spoke-group-gw-%s"

	enable_jumbo_frame  = false
	enable_gro_gso      = false
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName)
}

// ============================================================================
// Helper Functions
// ============================================================================

func testAccCheckSpokeGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke group not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke group ID is set")
		}

		client, ok := testAccProvider.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to get client")
		}

		_, err := client.GetGatewayGroup(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("spoke group not found: %w", err)
		}

		return nil
	}
}

func testAccCheckSpokeGroupDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to get client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_group" {
			continue
		}

		_, err := client.GetGatewayGroup(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("spoke group still exists: %s", rs.Primary.ID)
		}

		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("error checking spoke group destroyed: %w", err)
		}
	}

	return nil
}
