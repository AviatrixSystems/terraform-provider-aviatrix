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

func preGateway2Check(t *testing.T, msgCommon string) (string, string, string) {
	preAccountCheck(t, msgCommon)

	vpcID2 := os.Getenv("AWS_VPC_ID2")
	if vpcID2 == "" {
		t.Fatal("Environment variable AWS_VPC_ID2 is not set" + msgCommon)
	}

	region2 := os.Getenv("AWS_REGION2")
	if region2 == "" {
		t.Fatal("Environment variable AWS_REGION2 is not set" + msgCommon)
	}

	vpcNet2 := os.Getenv("AWS_VPC_NET2")
	if vpcNet2 == "" {
		t.Fatal("Environment variable AWS_VPC_NET2 is not set" + msgCommon)
	}
	return vpcID2, region2, vpcNet2
}

func preAvxTunnelCheck(t *testing.T, msgCommon string) (string, string, string, string, string, string) {
	vpcID1, region1, subnet1 := preGatewayCheck(t, msgCommon)
	vpcID2, region2, subnet2 := preGateway2Check(t, msgCommon)
	return vpcID1, region1, subnet1, vpcID2, region2, subnet2
}

func TestAccAviatrixTunnel_basic(t *testing.T) {
	var tun goaviatrix.Tunnel
	rName := acctest.RandString(5)
	resourceName := "aviatrix_tunnel.foo"

	skipAcc := os.Getenv("SKIP_TUNNEL")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix peering tunnel test as SKIP_TUNNEL is set")
	}
	msgCommon := ". Set SKIP_TUNNEL to yes to skip Aviatrix peering tunnel tests"

	vpcID1, region1, subnet1, vpcID2, region2, subnet2 := preAvxTunnelCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelConfigBasic(rName, vpcID1, vpcID2, region1, region2, subnet1, subnet2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckTunnelExists("aviatrix_tunnel.foo", &tun),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_name1", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_name2", fmt.Sprintf("tfg2-%s", rName)),
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

func testAccTunnelConfigBasic(rName string, vpcID1 string, vpcID2 string, region1 string, region2 string,
	subnet1 string, subnet2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name = "tfa-%s"
	cloud_type = 1
	aws_account_number = "%s"
	aws_iam = "false"
	aws_access_key = "%s"
	aws_secret_key = "%s"
}

resource "aviatrix_gateway" "gw1" {
	cloud_type = 1
	account_name = "${aviatrix_account.test.account_name}"
	gw_name = "tfg-%[1]s"
	vpc_id = "%[5]s"
	vpc_reg = "%[7]s"
	vpc_size = "t2.micro"
	vpc_net = "%[9]s"
}


resource "aviatrix_gateway" "gw2" {
	cloud_type = 1
	account_name = "${aviatrix_account.test.account_name}"
	gw_name = "tfg2-%[1]s"
	vpc_id = "%[6]s"
	vpc_reg = "%[8]s"
	vpc_size = "t2.micro"
	vpc_net = "%[10]s"
}

resource "aviatrix_tunnel" "foo" {
	vpc_name1 = "${aviatrix_gateway.gw1.gw_name}"
	vpc_name2 = "${aviatrix_gateway.gw2.gw_name}"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		vpcID1, vpcID2, region1, region2, subnet1, subnet2)
}

func tesAccCheckTunnelExists(n string, tunnel *goaviatrix.Tunnel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix tunnel Not Created: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix tunnel ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundTunnel := &goaviatrix.Tunnel{
			VpcName1: rs.Primary.Attributes["vpc_name1"],
			VpcName2: rs.Primary.Attributes["vpc_name2"],
		}

		_, err := client.GetTunnel(foundTunnel)

		if err != nil {
			return err
		}

		if foundTunnel.VpcName1 != rs.Primary.Attributes["vpc_name1"] {
			return fmt.Errorf("vpc_name1 Not found in created attributes")
		}

		if foundTunnel.VpcName2 != rs.Primary.Attributes["vpc_name2"] {
			return fmt.Errorf("vpc_name2 Not found in created attributes")
		}

		*tunnel = *foundTunnel

		return nil
	}
}

func testAccCheckTunnelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_tunnel" {
			continue
		}
		foundTunnel := &goaviatrix.Tunnel{
			VpcName1: rs.Primary.Attributes["vpc_name1"],
			VpcName2: rs.Primary.Attributes["vpc_name2"],
		}
		_, err := client.GetTunnel(foundTunnel)

		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("aviatrix tunnel still exists")
		}
	}
	return nil
}
