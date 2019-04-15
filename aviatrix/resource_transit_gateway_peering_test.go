package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func preAvxTransitGatewayPeeringCheck(t *testing.T, msgCommon string) (string, string, string, string, string, string,
	string, string) {
	vpcID1, region1, subnet1 := preGatewayCheck(t, msgCommon)
	vpcID2, region2, subnet2 := preGateway2Check(t, msgCommon)
	return vpcID1, region1, subnet1, subnet1, vpcID2, region2, subnet2, subnet2
}

func TestAvxTransitGatewayPeering_basic(t *testing.T) {
	rName := acctest.RandString(5)

	sourceName := "aviatrix_transit_gateway_peering.foo"

	skipAcc := os.Getenv("SKIP_TRANSIT_GATEWAY_PEERING")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix transit gateway peering test as SKIP_TRANSIT_GATEWAY_PEERING is set")
	}
	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_PEERING to yes to skip Aviatrix transit gateway peering tests"

	vpcID1, region1, subnet1, haSubnet1, vpcID2, region2, subnet2, haSubnet2 := preAvxTransitGatewayPeeringCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAvxTransitGatewayPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAvxTransitGatewayPeeringConfigBasic(rName, vpcID1, vpcID2, region1, region2,
					subnet1, subnet2, haSubnet1, haSubnet2),
				Check: resource.ComposeTestCheckFunc(
					tesAvxTransitGatewayPeeringExists("aviatrix_transit_gateway_peering.foo"),
					resource.TestCheckResourceAttr(
						sourceName, "transit_gateway_name1", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(
						sourceName, "transit_gateway_name2", fmt.Sprintf("tfg2-%s", rName)),
				),
			},
		},
	})
}

func testAvxTransitGatewayPeeringConfigBasic(rName string, vpcID1 string, vpcID2 string, region1 string, region2 string,
	subnet1 string, subnet2 string, haSubnet1 string, haSubnet2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name = "tfa-%s"
	cloud_type = 1
	aws_account_number = "%s"
	aws_iam = "false"
	aws_access_key = "%s"
	aws_secret_key = "%s"
}

resource "aviatrix_transit_vpc" "transitGw1" {
    cloud_type   = 1
    account_name = "${aviatrix_account.test.account_name}"
    gw_name      = "tfg-%s"
    vpc_id       = "%s"
    vpc_reg      = "%s"
    vpc_size     = "t2.micro"
    subnet       = "%s"
    ha_subnet    = "%s"
    ha_gw_size   = "t2.micro"
}

resource "aviatrix_transit_vpc" "transitGw2" {
    cloud_type   = 1
    account_name = "${aviatrix_account.test.account_name}"
    gw_name      = "tfg2-%s"
    vpc_id       = "%s"
    vpc_reg      = "%s"
    vpc_size     = "t2.micro"
    subnet       = "%s"
    ha_subnet    = "%s"
    ha_gw_size   = "t2.micro"
}

resource "aviatrix_transit_gateway_peering" "foo" {
	transit_gateway_name1 = "${aviatrix_transit_vpc.transitGw1.gw_name}"
	transit_gateway_name2 = "${aviatrix_transit_vpc.transitGw2.gw_name}"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, vpcID1, region1, subnet1, haSubnet1, rName, vpcID2, region2, subnet2, haSubnet2)
}

func tesAvxTransitGatewayPeeringExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix transit gateway peering Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix transit gateway peering ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

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

func testAvxTransitGatewayPeeringDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_gateway_peering" {
			continue
		}
		foundTransitGatewayPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: rs.Primary.Attributes["transit_gateway_name1"],
			TransitGatewayName2: rs.Primary.Attributes["transit_gateway_name2"],
		}
		err := client.GetTransitGatewayPeering(foundTransitGatewayPeering)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("aviatrix transit gateway peering still exists")
		}
	}
	return nil
}
