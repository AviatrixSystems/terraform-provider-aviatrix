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

func TestAccAviatrixCentralizedTransitFireNet_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_centralized_transit_firenet.test"

	if os.Getenv("SKIP_CENTRALIZED_TRANSIT_FIRENET") == "yes" {
		t.Skip("Skipping Centralized Transit FireNet test as SKIP_CENTRALIZED_TRANSIT_FIRENET is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCentralizedTransitFireNetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCentralizedTransitFireNetConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCentralizedTransitFireNetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_firenet_gw_name", fmt.Sprintf("primary-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "secondary_firenet_gw_name", fmt.Sprintf("secondary-%s", rName)),
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

func testAccCentralizedTransitFireNetConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test_vpc_1" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%s"
	name                 = "firenet-vpc-1"
	cidr                 = "10.10.0.0/24"
	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test_transit_gateway_1" {
	cloud_type                = aviatrix_vpc.test_vpc_1.cloud_type
	account_name              = aviatrix_account.test_account.account_name
	gw_name                   = "primary-%[1]s"
	vpc_id                    = aviatrix_vpc.test_vpc_1.vpc_id
	vpc_reg                   = aviatrix_vpc.test_vpc_1.region
	gw_size                   = "c5.xlarge"
	subnet                    = aviatrix_vpc.test_vpc_1.subnets[0].cidr
	enable_hybrid_connection  = true
	enable_transit_firenet    = true
	enable_segmentation       = true
	local_as_number           = "1"
	enable_multi_tier_transit = true
	connected_transit = true
}
resource "aviatrix_firenet" "test_firenet_1" {
	vpc_id = aviatrix_transit_gateway.test_transit_gateway_1.vpc_id
}
resource "aviatrix_vpc" "test_vpc_2" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%[5]s"
	name                 = "firenet-vpc-2"
	cidr                 = "10.11.0.0/24"
	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test_transit_gateway_2" {
	cloud_type               = aviatrix_vpc.test_vpc_2.cloud_type
	account_name             = aviatrix_account.test_account.account_name
	gw_name                  = "secondary-%[1]s"
	vpc_id                   = aviatrix_vpc.test_vpc_2.vpc_id
	vpc_reg                  = aviatrix_vpc.test_vpc_2.region
	gw_size                  = "c5.xlarge"
	subnet                   = aviatrix_vpc.test_vpc_2.subnets[0].cidr
	enable_hybrid_connection = true
	enable_transit_firenet   = true
	enable_segmentation      = true
	local_as_number          = "2"
	connected_transit        = true
}
resource "aviatrix_firenet" "test_firenet_2" {
	vpc_id             = aviatrix_transit_gateway.test_transit_gateway_2.vpc_id
	inspection_enabled = false
}
resource "aviatrix_centralized_transit_firenet" "test" {
	primary_firenet_gw_name   = aviatrix_transit_gateway.test_transit_gateway_1.gw_name
	secondary_firenet_gw_name = aviatrix_transit_gateway.test_transit_gateway_2.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"))
}

func testAccCheckCentralizedTransitFireNetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("centralized transit firenet not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no centralized transit firenet id is set")
		}

		client := mustClient(testAccProvider.Meta())

		centralizedTransitFirenet := &goaviatrix.CentralizedTransitFirenet{
			PrimaryGwName:   rs.Primary.Attributes["primary_firenet_gw_name"],
			SecondaryGwName: rs.Primary.Attributes["secondary_firenet_gw_name"],
		}

		err := client.GetCentralizedTransitFireNet(context.Background(), centralizedTransitFirenet)
		if err != nil {
			if errors.Is(err, goaviatrix.ErrNotFound) {
				return fmt.Errorf("centralized transit firenet does not exist")
			}
			return fmt.Errorf("could not get centralized transit firenet due to: %w", err)
		}

		if rs.Primary.ID != centralizedTransitFirenet.PrimaryGwName+"~"+centralizedTransitFirenet.SecondaryGwName {
			return fmt.Errorf("centralized transit firenet does not exist")
		}

		return nil
	}
}

func testAccCheckCentralizedTransitFireNetDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_centralized_transit_firenet" {
			continue
		}

		centralizedTransitFirenet := &goaviatrix.CentralizedTransitFirenet{
			PrimaryGwName:   rs.Primary.Attributes["primary_firenet_gw_name"],
			SecondaryGwName: rs.Primary.Attributes["secondary_firenet_gw_name"],
		}

		err := client.GetCentralizedTransitFireNet(context.Background(), centralizedTransitFirenet)
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("centralized transit firenet still exists")
		}
	}

	return nil
}
