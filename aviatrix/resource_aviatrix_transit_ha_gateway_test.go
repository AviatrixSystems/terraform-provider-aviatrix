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

func TestAccAviatrixTransitHaGateway_basic(t *testing.T) {
	if os.Getenv("SKIP_TRANSIT_HA_GATEWAY") == "yes" {
		t.Skip("Skipping Transit HA Gateway test as SKIP_TRANSIT_HA_GATEWAY is set")
	}

	resourceName := "aviatrix_transit_ha_gateway.test"
	primaryGwName := "tfg-" + acctest.RandString(5)
	haGwName := primaryGwName + "-hagw"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set SKIP_TRANSIT_HA_GATEWAY to yes to skip Transit HA Gateway tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitHaGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitHaGatewayConfigBasic(primaryGwName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitHaGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_gw_name", primaryGwName),
					resource.TestCheckResourceAttr(resourceName, "ha_gw_name", haGwName),
					resource.TestCheckResourceAttr(resourceName, "cloud_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "ha_gw_size", "t3.micro"),
					resource.TestCheckResourceAttrSet(resourceName, "ha_private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "ha_cloud_instance_id"),
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

func TestAccAviatrixTransitHaGateway_insaneMode(t *testing.T) {
	if os.Getenv("SKIP_TRANSIT_HA_GATEWAY") == "yes" {
		t.Skip("Skipping Transit HA Gateway Insane Mode test as SKIP_TRANSIT_HA_GATEWAY is set")
	}

	resourceName := "aviatrix_transit_ha_gateway.test"
	primaryGwName := "tfg-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set SKIP_TRANSIT_HA_GATEWAY to yes to skip Transit HA Gateway tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitHaGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitHaGatewayConfigInsaneMode(primaryGwName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitHaGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "primary_gw_name", primaryGwName),
					resource.TestCheckResourceAttr(resourceName, "insane_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "ha_insane_mode_az", "us-west-1b"),
					resource.TestCheckResourceAttr(resourceName, "ha_gw_size", "c5n.large"),
				),
			},
		},
	})
}

func testAccCheckTransitHaGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("transit HA gateway not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit HA gateway ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		gateway := &goaviatrix.Gateway{
			GwName: rs.Primary.ID,
		}

		_, err := client.GetGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to get transit HA gateway: %s", err)
		}

		if gateway.GwName != rs.Primary.ID {
			return fmt.Errorf("transit HA gateway not found")
		}

		return nil
	}
}

func testAccCheckTransitHaGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_ha_gateway" {
			continue
		}

		gateway := &goaviatrix.Gateway{
			GwName: rs.Primary.ID,
		}

		_, err := client.GetGateway(gateway)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("transit HA gateway still exists")
		}
	}

	return nil
}

func testAccTransitHaGatewayConfigBasic(primaryGwName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "primary" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.micro"
	subnet       = "%s"
}

resource "aviatrix_transit_ha_gateway" "test" {
	primary_gw_name = aviatrix_transit_gateway.primary.gw_name
	cloud_type      = 1
	account_name    = aviatrix_account.test.account_name
	vpc_id          = "%s"
	ha_subnet       = "%s"
	ha_gw_size      = "t3.micro"
}
	`, primaryGwName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		primaryGwName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_HA_SUBNET"))
}

func testAccTransitHaGatewayConfigInsaneMode(primaryGwName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "primary" {
	cloud_type      = 1
	account_name    = aviatrix_account.test.account_name
	gw_name         = "%s"
	vpc_id          = "%s"
	vpc_reg         = "%s"
	gw_size         = "c5n.large"
	subnet          = "%s"
	insane_mode     = true
	insane_mode_az  = "us-west-1a"
}

resource "aviatrix_transit_ha_gateway" "test" {
	primary_gw_name   = aviatrix_transit_gateway.primary.gw_name
	cloud_type        = 1
	account_name      = aviatrix_account.test.account_name
	vpc_id            = "%s"
	ha_subnet         = "%s"
	ha_gw_size        = "c5n.large"
	insane_mode       = true
	ha_insane_mode_az = "us-west-1b"
}
	`, primaryGwName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		primaryGwName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_HA_SUBNET"))
}
