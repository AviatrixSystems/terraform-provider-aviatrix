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

func preAvxTransitGatewayPeeringCheck(t *testing.T, msgCommon string) {
	preGatewayCheck(t, msgCommon)
	preGateway2Check(t, msgCommon)
}

func TestAccAviatrixTransitGatewayPeering_basic(t *testing.T) {
	rName := acctest.RandString(5)
	vpcID1 := os.Getenv("AWS_VPC_ID")
	region1 := os.Getenv("AWS_REGION")
	subnet1 := os.Getenv("AWS_SUBNET")
	haSubnet1 := os.Getenv("AWS_SUBNET")

	vpcID2 := os.Getenv("AWS_VPC_ID2")
	region2 := os.Getenv("AWS_REGION2")
	subnet2 := os.Getenv("AWS_SUBNET2")
	haSubnet2 := os.Getenv("AWS_SUBNET2")

	resourceName := "aviatrix_transit_gateway_peering.foo"

	accountName := "megaport-" + rName
	transit1GwName := "spoke-" + rName
	transit1SiteID := "site-" + rName
	path, _ := os.Getwd()
	transit2GwName := "transit-" + rName
	transit2SiteID := "site-" + rName
	resourceNameEdge := "aviatrix_transit_gateway_peering.test_transit_gateway_peering"

	skipAcc := os.Getenv("SKIP_TRANSIT_GATEWAY_PEERING")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix transit gateway peering test as SKIP_TRANSIT_GATEWAY_PEERING is set")
	}
	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_PEERING to yes to skip Aviatrix transit gateway peering tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAvxTransitGatewayPeeringCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayPeeringConfigBasic(rName, vpcID1, vpcID2, region1, region2,
					subnet1, subnet2, haSubnet1, haSubnet2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckTransitGatewayPeeringExists("aviatrix_transit_gateway_peering.foo"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name1", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name2", fmt.Sprintf("tfg2-%s", rName)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayPeeringConfigEdge(accountName, transit1GwName, transit1SiteID, path, transit2GwName, transit2SiteID),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckTransitGatewayPeeringExists(resourceNameEdge),
					resource.TestCheckResourceAttr(resourceNameEdge, "transit_gateway_name1", transit1GwName),
					resource.TestCheckResourceAttr(resourceNameEdge, "transit_gateway_name2", transit2GwName),
					resource.TestCheckResourceAttr(resourceNameEdge, "enable_peering_over_private_network", "true"),
					resource.TestCheckResourceAttr(resourceNameEdge, "jumbo_frame", "false"),
					resource.TestCheckResourceAttr(resourceNameEdge, "insane_mode", "true"),
					resource.TestCheckResourceAttr(resourceNameEdge, "gateway1_logical_ifnames.0", "wan1"),
					resource.TestCheckResourceAttr(resourceNameEdge, "gateway2_logical_ifnames.0", "wan1"),
					resource.TestCheckResourceAttr(resourceNameEdge, "disable_activemesh", "true"),
				),
			},
		},
	})
}

func testAccTransitGatewayPeeringConfigBasic(rName string, vpcID1 string, vpcID2 string, region1 string, region2 string,
	subnet1 string, subnet2 string, haSubnet1 string, haSubnet2 string,
) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "transitGw1" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
	ha_subnet    = "%s"
	ha_gw_size   = "t2.micro"
}
resource "aviatrix_transit_gateway" "transitGw2" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg2-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
	ha_subnet    = "%s"
	ha_gw_size   = "t2.micro"
}
resource "aviatrix_transit_gateway_peering" "foo" {
	transit_gateway_name1 = aviatrix_transit_gateway.transitGw1.gw_name
	transit_gateway_name2 = aviatrix_transit_gateway.transitGw2.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, vpcID1, region1, subnet1, haSubnet1, rName, vpcID2, region2, subnet2, haSubnet2)
}

func tesAccCheckTransitGatewayPeeringExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix transit gateway peering Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix transit gateway peering ID is set")
		}

		client := mustClient(testAccProvider.Meta())

		foundTransitGatewayPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: rs.Primary.Attributes["transit_gateway_name1"],
			TransitGatewayName2: rs.Primary.Attributes["transit_gateway_name2"],
		}

		err := client.GetTransitGatewayPeering(foundTransitGatewayPeering)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccTransitGatewayPeeringConfigEdge(accountName, transit1GwName, transit1SiteID, path, transit2GwName, transit2SiteID string) string {
	return fmt.Sprintf(`
	resource "aviatrix_account" "test_acc_edge_megaport" {
		account_name       = "edge-%s"
		cloud_type         = 1048576
	}
	resource "aviatrix_transit_gateway" "test_edge_transit_1" {
		cloud_type   = 1048576
		account_name = aviatrix_account.test_acc_edge_megaport.account_name
		gw_name      = "%s"
		vpc_id       = "%s"
		gw_size      = "SMALL"
		ztp_file_download_path = "%s"
		interfaces {
			gateway_ip     = "192.168.20.1"
			ip_address     = "192.168.20.11/24"
			public_ip      = "67.207.104.19"
			logical_ifname = "wan0"
			secondary_private_cidr_list = ["192.168.20.16/29"]
		}
		interfaces {
			gateway_ip     = "192.168.21.1"
			ip_address     = "192.168.21.11/24"
			public_ip      = "67.71.12.148"
			logical_ifname = "wan1"
			secondary_private_cidr_list = ["192.168.21.16/29"]
		}
		interfaces {
			dhcp           = true
			logical_ifname = "mgmt0"
		}
		interfaces {
			gateway_ip     = "192.168.22.1"
			ip_address     = "192.168.22.11/24"
			logical_ifname = "wan2"
		}
		interfaces {
			gateway_ip     = "192.168.23.1"
			ip_address     = "192.168.23.11/24"
			logical_ifname = "wan3"
		}
	}
	resource "aviatrix_transit_gateway" "test_edge_transit_2" {
		cloud_type   = 1048576
		account_name = aviatrix_account.test_acc_edge_megaport.account_name
		gw_name      = "%s"
		vpc_id       = "%s"
		gw_size      = "SMALL"
		ztp_file_download_path = "%s"
		interfaces {
			gateway_ip     = "192.168.24.1"
			ip_address     = "192.168.24.11/24"
			public_ip      = "67.207.104.24"
			logical_ifname = "wan0"
			secondary_private_cidr_list = ["192.168.24.16/29"]
		}
		interfaces {
			gateway_ip     = "192.168.25.1"
			ip_address     = "192.168.25.11/24"
			public_ip      = "67.71.12.25"
			logical_ifname = "wan1"
			secondary_private_cidr_list = ["192.168.25.16/29"]
		}
		interfaces {
			dhcp           = true
			logical_ifname = "mgmt0"
		}
		interfaces {
			gateway_ip     = "192.168.26.1"
			ip_address     = "192.168.26.11/24"
			logical_ifname = "wan2"
		}
		interfaces {
			gateway_ip     = "192.168.27.1"
			ip_address     = "192.168.27.11/24"
			logical_ifname = "wan3"
		}
	}
	resource "aviatrix_transit_gateway_peering" "test_transit_gateway_peering" {
		transit_gateway_name1 = aviatrix_transit_gateway.test_edge_transit_1.gw_name
		transit_gateway_name2 = aviatrix_transit_gateway.test_edge_transit_2.gw_name
		enable_peering_over_private_network = true
		jumbo_frame = false
		insane_mode = true
		gateway1_logical_ifnames = ["wan1"]
		gateway2_logical_ifnames = ["wan1"]
		disable_activemesh = true
	}
		`, accountName, transit1GwName, transit1SiteID, path, transit2GwName, transit2SiteID, path)
}

func testAccCheckTransitGatewayPeeringDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_gateway_peering" {
			continue
		}

		foundTransitGatewayPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: rs.Primary.Attributes["transit_gateway_name1"],
			TransitGatewayName2: rs.Primary.Attributes["transit_gateway_name2"],
		}

		err := client.GetTransitGatewayPeering(foundTransitGatewayPeering)
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("aviatrix transit gateway peering still exists")
		}
	}

	return nil
}

func TestAccAviatrixTransitGatewayPeering_insaneModeNotSetInState(t *testing.T) {
	rName := acctest.RandString(5)
	vpcID1 := os.Getenv("AWS_VPC_ID")
	region1 := os.Getenv("AWS_REGION")
	subnet1 := os.Getenv("AWS_SUBNET")

	vpcID2 := os.Getenv("AWS_VPC_ID2")
	region2 := os.Getenv("AWS_REGION2")
	subnet2 := os.Getenv("AWS_SUBNET2")

	resourceName := "aviatrix_transit_gateway_peering.test_insane_mode_not_set"

	skipAcc := os.Getenv("SKIP_TRANSIT_GATEWAY_PEERING")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix transit gateway peering test as SKIP_TRANSIT_GATEWAY_PEERING is set")
	}
	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_PEERING to yes to skip Aviatrix transit gateway peering tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAvxTransitGatewayPeeringCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayPeeringConfigInsaneModeNotSet(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckTransitGatewayPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name1", "transit-"+rName+"-1"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name2", "transit-"+rName+"-2"),
					resource.TestCheckResourceAttr(resourceName, "enable_insane_mode_encryption_over_internet", "true"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_count", "4"),
					// Verify that insane_mode is NOT set in state when user doesn't explicitly set it
					resource.TestCheckNoResourceAttr(resourceName, "insane_mode"),
				),
			},
			{
				// Import the resource to verify it can be imported without insane_mode
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Skip insane_mode during import verification since it's not set
				ImportStateVerifyIgnore: []string{"insane_mode"},
			},
		},
	})
}

func TestAccAviatrixTransitGatewayPeering_insaneModeDefaultValueBehavior(t *testing.T) {
	rName := acctest.RandString(5)
	vpcID1 := os.Getenv("AWS_VPC_ID")
	region1 := os.Getenv("AWS_REGION")
	subnet1 := os.Getenv("AWS_SUBNET")

	vpcID2 := os.Getenv("AWS_VPC_ID2")
	region2 := os.Getenv("AWS_REGION2")
	subnet2 := os.Getenv("AWS_SUBNET2")

	resourceName := "aviatrix_transit_gateway_peering.test_insane_mode_default"

	skipAcc := os.Getenv("SKIP_TRANSIT_GATEWAY_PEERING")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix transit gateway peering test as SKIP_TRANSIT_GATEWAY_PEERING is set")
	}
	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_PEERING to yes to skip Aviatrix transit gateway peering tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAvxTransitGatewayPeeringCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayPeeringDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1: Create without insane_mode (should use default value but not set in state)
				Config: testAccTransitGatewayPeeringConfigInsaneModeDefault(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckTransitGatewayPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name1", "transit-"+rName+"-1"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name2", "transit-"+rName+"-2"),
					resource.TestCheckResourceAttr(resourceName, "enable_insane_mode_encryption_over_internet", "true"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_count", "4"),
					// Verify that insane_mode is NOT set in state even though default is false
					resource.TestCheckNoResourceAttr(resourceName, "insane_mode"),
				),
			},
			{
				// Step 2: Update to explicitly set insane_mode = false (should now be in state)
				Config: testAccTransitGatewayPeeringConfigInsaneModeExplicitFalse(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckTransitGatewayPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name1", "transit-"+rName+"-1"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name2", "transit-"+rName+"-2"),
					resource.TestCheckResourceAttr(resourceName, "enable_insane_mode_encryption_over_internet", "true"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_count", "4"),
					// Verify that insane_mode IS now set in state because user explicitly set it
					resource.TestCheckResourceAttr(resourceName, "insane_mode", "false"),
				),
			},
			{
				// Step 3: Update to explicitly set insane_mode = true (should still be in state)
				Config: testAccTransitGatewayPeeringConfigInsaneModeExplicitTrue(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckTransitGatewayPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name1", "transit-"+rName+"-1"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name2", "transit-"+rName+"-2"),
					resource.TestCheckResourceAttr(resourceName, "enable_insane_mode_encryption_over_internet", "true"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_count", "4"),
					// Verify that insane_mode IS set in state because user explicitly set it
					resource.TestCheckResourceAttr(resourceName, "insane_mode", "true"),
				),
			},
			{
				// Step 4: Remove insane_mode again (should remove from state)
				Config: testAccTransitGatewayPeeringConfigInsaneModeDefault(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckTransitGatewayPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name1", "transit-"+rName+"-1"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name2", "transit-"+rName+"-2"),
					resource.TestCheckResourceAttr(resourceName, "enable_insane_mode_encryption_over_internet", "true"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_count", "4"),
					// Verify that insane_mode is NOT set in state again
					resource.TestCheckNoResourceAttr(resourceName, "insane_mode"),
				),
			},
		},
	})
}

func TestAccAviatrixTransitGatewayPeering_insaneModeSetInState(t *testing.T) {
	rName := acctest.RandString(5)
	vpcID1 := os.Getenv("AWS_VPC_ID")
	region1 := os.Getenv("AWS_REGION")
	subnet1 := os.Getenv("AWS_SUBNET")

	vpcID2 := os.Getenv("AWS_VPC_ID2")
	region2 := os.Getenv("AWS_REGION2")
	subnet2 := os.Getenv("AWS_SUBNET2")

	resourceName := "aviatrix_transit_gateway_peering.test_insane_mode_set"

	skipAcc := os.Getenv("SKIP_TRANSIT_GATEWAY_PEERING")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix transit gateway peering test as SKIP_TRANSIT_GATEWAY_PEERING is set")
	}
	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_PEERING to yes to skip Aviatrix transit gateway peering tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAvxTransitGatewayPeeringCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayPeeringConfigInsaneModeSet(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckTransitGatewayPeeringExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name1", "transit-"+rName+"-1"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name2", "transit-"+rName+"-2"),
					resource.TestCheckResourceAttr(resourceName, "enable_insane_mode_encryption_over_internet", "true"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_count", "4"),
					// Verify that insane_mode IS set in state when user explicitly sets it
					resource.TestCheckResourceAttr(resourceName, "insane_mode", "true"),
				),
			},
			{
				// Import the resource to verify it can be imported with insane_mode
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTransitGatewayPeeringConfigInsaneModeNotSet(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account_1" {
	account_name       = "tfa-%s-1"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_account" "test_account_2" {
	account_name       = "tfa-%s-2"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test_transit_gateway_1" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_1.account_name
	gw_name      = "transit-%s-1"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway" "test_transit_gateway_2" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_2.account_name
	gw_name      = "transit-%s-2"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway_peering" "test_insane_mode_not_set" {
	transit_gateway_name1 = aviatrix_transit_gateway.test_transit_gateway_1.gw_name
	transit_gateway_name2 = aviatrix_transit_gateway.test_transit_gateway_2.gw_name
	enable_insane_mode_encryption_over_internet = true
	tunnel_count = 4
	# Note: insane_mode is NOT explicitly set here
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_ACCOUNT_NUMBER2"), os.Getenv("AWS_ACCESS_KEY2"), os.Getenv("AWS_SECRET_KEY2"),
		rName, vpcID1, region1, subnet1,
		rName, vpcID2, region2, subnet2)
}

func testAccTransitGatewayPeeringConfigInsaneModeSet(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account_1" {
	account_name       = "tfa-%s-1"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_account" "test_account_2" {
	account_name       = "tfa-%s-2"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test_transit_gateway_1" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_1.account_name
	gw_name      = "transit-%s-1"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway" "test_transit_gateway_2" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_2.account_name
	gw_name      = "transit-%s-2"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway_peering" "test_insane_mode_set" {
	transit_gateway_name1 = aviatrix_transit_gateway.test_transit_gateway_1.gw_name
	transit_gateway_name2 = aviatrix_transit_gateway.test_transit_gateway_2.gw_name
	enable_insane_mode_encryption_over_internet = true
	tunnel_count = 4
	insane_mode = true  # Explicitly set insane_mode
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_ACCOUNT_NUMBER2"), os.Getenv("AWS_ACCESS_KEY2"), os.Getenv("AWS_SECRET_KEY2"),
		rName, vpcID1, region1, subnet1,
		rName, vpcID2, region2, subnet2)
}

func testAccTransitGatewayPeeringConfigInsaneModeDefault(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account_1" {
	account_name       = "tfa-%s-1"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_account" "test_account_2" {
	account_name       = "tfa-%s-2"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test_transit_gateway_1" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_1.account_name
	gw_name      = "transit-%s-1"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway" "test_transit_gateway_2" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_2.account_name
	gw_name      = "transit-%s-2"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway_peering" "test_insane_mode_default" {
	transit_gateway_name1 = aviatrix_transit_gateway.test_transit_gateway_1.gw_name
	transit_gateway_name2 = aviatrix_transit_gateway.test_transit_gateway_2.gw_name
	enable_insane_mode_encryption_over_internet = true
	tunnel_count = 4
	# insane_mode not set - should use default value but not appear in state
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_ACCOUNT_NUMBER2"), os.Getenv("AWS_ACCESS_KEY2"), os.Getenv("AWS_SECRET_KEY2"),
		rName, vpcID1, region1, subnet1,
		rName, vpcID2, region2, subnet2)
}

func testAccTransitGatewayPeeringConfigInsaneModeExplicitFalse(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account_1" {
	account_name       = "tfa-%s-1"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_account" "test_account_2" {
	account_name       = "tfa-%s-2"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test_transit_gateway_1" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_1.account_name
	gw_name      = "transit-%s-1"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway" "test_transit_gateway_2" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_2.account_name
	gw_name      = "transit-%s-2"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway_peering" "test_insane_mode_default" {
	transit_gateway_name1 = aviatrix_transit_gateway.test_transit_gateway_1.gw_name
	transit_gateway_name2 = aviatrix_transit_gateway.test_transit_gateway_2.gw_name
	enable_insane_mode_encryption_over_internet = true
	tunnel_count = 4
	insane_mode = false  # Explicitly set to false - should appear in state
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_ACCOUNT_NUMBER2"), os.Getenv("AWS_ACCESS_KEY2"), os.Getenv("AWS_SECRET_KEY2"),
		rName, vpcID1, region1, subnet1,
		rName, vpcID2, region2, subnet2)
}

func testAccTransitGatewayPeeringConfigInsaneModeExplicitTrue(rName, vpcID1, region1, subnet1, vpcID2, region2, subnet2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account_1" {
	account_name       = "tfa-%s-1"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_account" "test_account_2" {
	account_name       = "tfa-%s-2"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test_transit_gateway_1" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_1.account_name
	gw_name      = "transit-%s-1"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway" "test_transit_gateway_2" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account_2.account_name
	gw_name      = "transit-%s-2"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
	enable_insane_mode = true
}

resource "aviatrix_transit_gateway_peering" "test_insane_mode_default" {
	transit_gateway_name1 = aviatrix_transit_gateway.test_transit_gateway_1.gw_name
	transit_gateway_name2 = aviatrix_transit_gateway.test_transit_gateway_2.gw_name
	enable_insane_mode_encryption_over_internet = true
	tunnel_count = 4
	insane_mode = true  # Explicitly set to true - should appear in state
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_ACCOUNT_NUMBER2"), os.Getenv("AWS_ACCESS_KEY2"), os.Getenv("AWS_SECRET_KEY2"),
		rName, vpcID1, region1, subnet1,
		rName, vpcID2, region2, subnet2)
}
