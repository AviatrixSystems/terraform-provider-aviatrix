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

func preTransitGroupCheck(t *testing.T, msgCommon string) {
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

func TestAccAviatrixTransitGroup_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tfg-transit-group-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "cloud_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "gw_type", "transit"),
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

func TestAccAviatrixTransitGroup_update(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
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
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with initial values
			{
				Config: testAccTransitGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tfg-transit-group-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "group_instance_size", awsGwSize),
					resource.TestCheckResourceAttr(resourceName, "enable_nat", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_vpc_dns_server", "false"),
					resource.TestCheckResourceAttr(resourceName, "bgp_polling_time", "50"),
					resource.TestCheckResourceAttr(resourceName, "bgp_neighbor_status_polling_time", "5"),
					resource.TestCheckResourceAttr(resourceName, "bgp_hold_time", "180"),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_ecmp", "false"),
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
				Config: testAccTransitGroupConfigUpdateSize(rName, awsGwSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_instance_size", awsGwSizeUpdated),
					resource.TestCheckResourceAttr(resourceName, "enable_nat", "true"),
				),
			},
			// Step 3: Update BGP settings
			{
				Config: testAccTransitGroupConfigUpdateBgp(rName, awsGwSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bgp_polling_time", "30"),
					resource.TestCheckResourceAttr(resourceName, "bgp_hold_time", "120"),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_ecmp", "true"),
				),
			},
			// Step 4: Update feature flags
			{
				Config: testAccTransitGroupConfigUpdateFeatures(rName, awsGwSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_vpc_dns_server", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_connected_transit", "true"),
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

func TestAccAviatrixTransitGroup_activeStandby(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with active-standby disabled
			{
				Config: testAccTransitGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby_preemptive", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "explicitly_created"),
				),
			},
			// Step 2: Enable active-standby
			{
				Config: testAccTransitGroupConfigActiveStandby(rName, awsGwSize, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby_preemptive", "false"),
				),
			},
			// Step 3: Enable active-standby with preemptive
			{
				Config: testAccTransitGroupConfigActiveStandby(rName, awsGwSize, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby_preemptive", "true"),
				),
			},
			// Step 4: Disable active-standby
			{
				Config: testAccTransitGroupConfigActiveStandby(rName, awsGwSize, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_active_standby_preemptive", "false"),
				),
			},
		},
	})
}

func TestAccAviatrixTransitGroup_learnedCidrsApproval(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Enable learned CIDRs approval
			{
				Config: testAccTransitGroupConfigLearnedCidrs(rName, awsGwSize, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_learned_cidrs_approval", "true"),
					resource.TestCheckResourceAttr(resourceName, "learned_cidrs_approval_mode", "gateway"),
					resource.TestCheckResourceAttrSet(resourceName, "explicitly_created"),
				),
			},
			// Step 2: Disable learned CIDRs approval
			{
				Config: testAccTransitGroupConfigLearnedCidrs(rName, awsGwSize, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_learned_cidrs_approval", "false"),
					resource.TestCheckResourceAttr(resourceName, "learned_cidrs_approval_mode", "gateway"),
					resource.TestCheckResourceAttrSet(resourceName, "explicitly_created"),
				),
			},
		},
	})
}

func TestAccAviatrixTransitGroup_bgp(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with BGP enabled and local_as_number
			{
				Config: testAccTransitGroupConfigBgp(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp", "true"),
					resource.TestCheckResourceAttr(resourceName, "local_as_number", "65001"),
				),
			},
			// Step 2: Add prepend_as_path
			{
				Config: testAccTransitGroupConfigBgpWithPrepend(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
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

func TestAccAviatrixTransitGroup_featureFlags(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with default feature flags (jumbo_frame and gro_gso default to true)
			{
				Config: testAccTransitGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_jumbo_frame", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_gro_gso", "true"),
				),
			},
			// Step 2: Disable jumbo_frame and gro_gso
			{
				Config: testAccTransitGroupConfigFeatureFlagsDisabled(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_jumbo_frame", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_gro_gso", "false"),
				),
			},
			// Step 3: Re-enable jumbo_frame and gro_gso
			{
				Config: testAccTransitGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_jumbo_frame", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_gro_gso", "true"),
				),
			},
		},
	})
}

func TestAccAviatrixTransitGroup_transitFeatures(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with transit-specific features enabled
			{
				Config: testAccTransitGroupConfigTransitFeatures(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_connected_transit", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_segmentation", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_advertise_transit_cidr", "true"),
				),
			},
			// Step 2: Disable transit features
			{
				Config: testAccTransitGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_connected_transit", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_segmentation", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_advertise_transit_cidr", "false"),
				),
			},
		},
	})
}

func TestAccAviatrixTransitGroup_firenet(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with FireNet disabled
			{
				Config: testAccTransitGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
				),
			},
			// Step 2: Enable FireNet
			{
				Config: testAccTransitGroupConfigFireNet(rName, awsGwSize, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
				),
			},
			// Step 3: Disable FireNet
			{
				Config: testAccTransitGroupConfigFireNet(rName, awsGwSize, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
				),
			},
			// Step 4: Import state verification
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

func TestAccAviatrixTransitGroup_transitFirenet(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with Transit FireNet disabled
			{
				Config: testAccTransitGroupConfigBasic(rName, awsGwSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_gateway_load_balancer", "false"),
				),
			},
			// Step 2: Enable Transit FireNet
			{
				Config: testAccTransitGroupConfigTransitFireNet(rName, awsGwSize, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_gateway_load_balancer", "false"),
				),
			},
			// Step 3: Disable Transit FireNet
			{
				Config: testAccTransitGroupConfigTransitFireNet(rName, awsGwSize, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_gateway_load_balancer", "false"),
				),
			},
			// Step 4: Import state verification
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

func TestAccAviatrixTransitGroup_firenetToggle(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Enable FireNet
			{
				Config: testAccTransitGroupConfigFireNetToggle(rName, awsGwSize, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
				),
			},
			// Step 2: Switch to Transit FireNet
			{
				Config: testAccTransitGroupConfigFireNetToggle(rName, awsGwSize, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "true"),
				),
			},
			// Step 3: Switch back to FireNet
			{
				Config: testAccTransitGroupConfigFireNetToggle(rName, awsGwSize, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
				),
			},
			// Step 4: Disable both
			{
				Config: testAccTransitGroupConfigFireNetToggle(rName, awsGwSize, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
				),
			},
			// Step 5: Import state verification
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

func TestAccAviatrixTransitGroup_transitFirenetGWLB(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_group.test"

	msgCommon := ". Set SKIP_TRANSIT_GROUP to yes to skip Transit Group tests"

	skipTransitGroup := os.Getenv("SKIP_TRANSIT_GROUP")
	if skipTransitGroup == "yes" {
		t.Skip("Skipping Transit Group test as SKIP_TRANSIT_GROUP is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t3.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preTransitGroupCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Enable Transit FireNet without GWLB
			{
				Config: testAccTransitGroupConfigTransitFireNetGWLB(rName, awsGwSize, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_gateway_load_balancer", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
				),
			},
			// Step 2: Disable Transit FireNet (required before enabling with GWLB)
			{
				Config: testAccTransitGroupConfigTransitFireNetGWLB(rName, awsGwSize, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_gateway_load_balancer", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
				),
			},
			// Step 3: Enable Transit FireNet with GWLB
			{
				Config: testAccTransitGroupConfigTransitFireNetGWLB(rName, awsGwSize, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_gateway_load_balancer", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
				),
			},
			// Step 4: Disable both
			{
				Config: testAccTransitGroupConfigTransitFireNetGWLB(rName, awsGwSize, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_transit_firenet", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_gateway_load_balancer", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_firenet", "false"),
				),
			},
			// Step 5: Import state verification
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

// ============================================================================
// Config Generators
// ============================================================================

func testAccTransitGroupConfigBasic(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%[1]s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%[7]s"
	vpc_region          = "%[8]s"
	gw_name             = "tfg-transit-group-gw-%[1]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"))
}

func testAccTransitGroupConfigUpdateSize(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%[1]s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%[7]s"
	vpc_region          = "%[8]s"
	gw_name             = "tfg-transit-group-gw-%[1]s"

	enable_nat          = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"))
}

func testAccTransitGroupConfigUpdateBgp(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%[1]s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%[7]s"
	vpc_region          = "%[8]s"
	gw_name             = "tfg-transit-group-gw-%[1]s"

	enable_nat          = true
	bgp_polling_time    = 30
	bgp_hold_time       = 120
	enable_bgp_ecmp     = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"))
}

func testAccTransitGroupConfigUpdateFeatures(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%[1]s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%[5]s"
	vpc_id              = "%[6]s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%[7]s"
	vpc_region          = "%[8]s"
	gw_name             = "tfg-transit-group-gw-%[1]s"

	enable_nat                = true
	bgp_polling_time          = 30
	bgp_hold_time             = 120
	enable_bgp_ecmp           = true
	enable_vpc_dns_server     = true
	enable_connected_transit  = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"))
}

func testAccTransitGroupConfigActiveStandby(rName, gwSize string, enableActiveStandby, enablePreemptive bool) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"

	enable_active_standby            = %t
	enable_active_standby_preemptive = %t
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName,
		enableActiveStandby, enablePreemptive)
}

func testAccTransitGroupConfigLearnedCidrs(rName, gwSize string, enableLearnedCidrs bool) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"
	enable_bgp          = true
	local_as_number     = "65001"

	enable_learned_cidrs_approval = %t
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName,
		enableLearnedCidrs)
}

func testAccTransitGroupConfigBgp(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"

	enable_bgp          = true
	local_as_number     = "65001"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName)
}

func testAccTransitGroupConfigBgpWithPrepend(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"

	enable_bgp          = true
	local_as_number     = "65001"
	prepend_as_path     = ["65001", "65001"]
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName)
}

func testAccTransitGroupConfigFeatureFlagsDisabled(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"

	enable_jumbo_frame  = false
	enable_gro_gso      = false
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName)
}

func testAccTransitGroupConfigTransitFeatures(rName, gwSize string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"

	enable_connected_transit      = true
	enable_segmentation           = true
	enable_advertise_transit_cidr = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName)
}

func testAccTransitGroupConfigFireNet(rName, gwSize string, enableFireNet bool) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"

	enable_firenet         = %t
	enable_transit_firenet = false
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName,
		enableFireNet)
}

func testAccTransitGroupConfigTransitFireNet(rName, gwSize string, enableTransitFireNet bool) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"

	enable_firenet            = false
	enable_transit_firenet    = %t
	enable_connected_transit  = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName,
		enableTransitFireNet)
}

func testAccTransitGroupConfigFireNetToggle(rName, gwSize string, enableFireNet, enableTransitFireNet bool) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"

	enable_firenet         = %t
	enable_transit_firenet = %t
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName,
		enableFireNet, enableTransitFireNet)
}

func testAccTransitGroupConfigTransitFireNetGWLB(rName, gwSize string, enableTransitFireNet, enableGWLB bool) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_group" "test" {
	group_name          = "tfg-transit-group-%s"
	cloud_type          = 1
	gw_type             = "transit"
	group_instance_size = "%s"
	vpc_id              = "%s"
	account_name        = aviatrix_account.test_acc_aws.account_name
	subnet              = "%s"
	vpc_region          = "%s"
	gw_name             = "tfg-transit-group-gw-%s"

	enable_firenet                = false
	enable_transit_firenet        = %t
	enable_gateway_load_balancer  = %t
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, gwSize, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"), os.Getenv("AWS_REGION"), rName,
		enableTransitFireNet, enableGWLB)
}

// ============================================================================
// Helper Functions
// ============================================================================

func testAccCheckTransitGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transit group not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit group ID is set")
		}

		client, ok := testAccProvider.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to get client")
		}

		_, err := client.GetGatewayGroup(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("transit group not found: %w", err)
		}

		return nil
	}
}

func testAccCheckTransitGroupDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to get client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_group" {
			continue
		}

		_, err := client.GetGatewayGroup(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("transit group still exists: %s", rs.Primary.ID)
		}

		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("error checking transit group destroyed: %w", err)
		}
	}

	return nil
}
