package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func preVGWConnCheck(t *testing.T, msgCommon string) {
	preAccountCheck(t, msgCommon)
	preGatewayCheck(t, msgCommon)

	bgpVGWId := os.Getenv("AWS_BGP_VGW_ID")
	if bgpVGWId == "" {
		t.Fatal("Environment variable AWS_BGP_VGW_ID is not set" + msgCommon)
	}
}

func TestAccAviatrixVGWConn_basic(t *testing.T) {
	var vgwConn goaviatrix.VGWConn
	vpcID := os.Getenv("AWS_VPC_ID")
	bgpVGWId := os.Getenv("AWS_BGP_VGW_ID")

	rName := acctest.RandString(5)

	resourceName := "aviatrix_vgw_conn.test_vgw_conn"

	skipAcc := os.Getenv("SKIP_VGW_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping aviatrix VGW connection test as SKIP_VGW_CONN is set")
	}
	msgCommon := ". Set SKIP_VGW_CONN to yes to skip VGW connection tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preVGWConnCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVGWConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVGWConnConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVGWConnExists(resourceName, &vgwConn),
					resource.TestCheckResourceAttr(resourceName, "conn_name", fmt.Sprintf("tfc-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "bgp_vgw_id", bgpVGWId),
					resource.TestCheckResourceAttr(resourceName, "bgp_vgw_account", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "bgp_vgw_region", os.Getenv("AWS_REGION2")),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "6451"),
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

func testAccVGWConnConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_vpc" {
	account_name = aviatrix_account.test_account.account_name
	cloud_type   = 1
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_vgw_conn" "test_vgw_conn" {
	conn_name        = "tfc-%s"
	gw_name          = aviatrix_transit_gateway.test_transit_vpc.gw_name
	vpc_id           = aviatrix_transit_gateway.test_transit_vpc.vpc_id
	bgp_vgw_id       = "%s"
	bgp_vgw_account  = aviatrix_account.test_account.account_name
	bgp_vgw_region   = "%s"
	bgp_local_as_num = "6451"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"),
		rName, os.Getenv("AWS_BGP_VGW_ID"), os.Getenv("AWS_REGION2"))
}

func testAccCheckVGWConnExists(n string, vgwConn *goaviatrix.VGWConn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vgw connection Not created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no VGW connection ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundVGWConn := &goaviatrix.VGWConn{
			ConnName:      rs.Primary.Attributes["conn_name"],
			GwName:        rs.Primary.Attributes["gw_name"],
			VPCId:         rs.Primary.Attributes["vpc_id"],
			BgpVGWId:      rs.Primary.Attributes["bgp_vgw_id"],
			BgpVGWAccount: rs.Primary.Attributes["bgp_vgw_account"],
			BgpVGWRegion:  rs.Primary.Attributes["bgp_vgw_region"],
			BgpLocalAsNum: rs.Primary.Attributes["bgp_local_as_num"],
		}

		foundVGWConn2, err := client.GetVGWConn(foundVGWConn)
		if err != nil {
			return err
		}
		if foundVGWConn2.ConnName != rs.Primary.Attributes["conn_name"] {
			return fmt.Errorf("conn_name Not found in created attributes")
		}

		*vgwConn = *foundVGWConn
		return nil
	}
}

func testAccCheckVGWConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vgw_conn" {
			continue
		}

		foundVGWConn := &goaviatrix.VGWConn{
			ConnName:      rs.Primary.Attributes["conn_name"],
			GwName:        rs.Primary.Attributes["gw_name"],
			VPCId:         rs.Primary.Attributes["vpc_id"],
			BgpVGWId:      rs.Primary.Attributes["bgp_vgw_id"],
			BgpVGWAccount: rs.Primary.Attributes["bgp_vgw_account"],
			BgpVGWRegion:  rs.Primary.Attributes["bgp_vgw_region"],
			BgpLocalAsNum: rs.Primary.Attributes["bgp_local_as_num"],
		}

		_, err := client.GetVGWConn(foundVGWConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("vgw connection still exists")
		}
	}

	return nil
}
