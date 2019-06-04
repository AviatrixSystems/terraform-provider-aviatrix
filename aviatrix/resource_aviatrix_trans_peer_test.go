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

func preTransPeerCheck(t *testing.T, msgCommon string) (string, string, string, string, string, string, string) {
	source, region1, subnet1, nextHop, region2, subnet2 := preAvxTunnelCheck(t, msgCommon)

	reachableCIDR := "192.168.0.0/16"

	return source, region1, subnet1, nextHop, region2, subnet2, reachableCIDR
}

func TestAccAviatrixTransPeer_basic(t *testing.T) {
	var transpeer goaviatrix.TransPeer
	rName := acctest.RandString(5)
	resourceName := "aviatrix_trans_peer.test_trans_peer"

	skipAcc := os.Getenv("SKIP_TRANS_PEER")
	if skipAcc == "yes" {
		t.Skip("Skipping aviatrix transitive peering test as SKIP_TRANS_PEER is set")
	}
	msgCommon := ". Set SKIP_TRANS_PEER to yes to skip transitive peer tests"

	preAccountCheck(t, msgCommon)

	sourceVPC, region1, subnet1, nextHopVPC, region2, subnet2, reachableCIDR := preTransPeerCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTransPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransPeerConfigBasic(rName, sourceVPC, nextHopVPC, region1, region2, subnet1, subnet2,
					reachableCIDR),
				Check: resource.ComposeTestCheckFunc(
					testAccTransPeerExists("aviatrix_trans_peer.test_trans_peer", &transpeer),
					resource.TestCheckResourceAttr(
						resourceName, "source", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "nexthop", fmt.Sprintf("tfg2-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "reachable_cidr", reachableCIDR),
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

func testAccTransPeerConfigBasic(rName string, source string, nextHop string, region1 string, region2 string,
	subnet1 string, subnet2 string, reachableCIDR string) string {
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

resource "aviatrix_trans_peer" "test_trans_peer" {
	source = "${aviatrix_tunnel.foo.vpc_name1}"
	nexthop = "${aviatrix_tunnel.foo.vpc_name2}"
	reachable_cidr = "%s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		source, nextHop, region1, region2, subnet1, subnet2, reachableCIDR)
}

func testAccTransPeerExists(n string, transpeer *goaviatrix.TransPeer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transpeer Not Created: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no transpeer ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundTransPeer := &goaviatrix.TransPeer{
			Source:        rs.Primary.Attributes["source"],
			Nexthop:       rs.Primary.Attributes["nexthop"],
			ReachableCidr: rs.Primary.Attributes["reachable_cidr"],
		}

		_, err := client.GetTransPeer(foundTransPeer)

		if err != nil {
			return err
		}
		if foundTransPeer.Source != rs.Primary.Attributes["source"] {
			return fmt.Errorf("source Not found in created attributes")
		}
		if foundTransPeer.Nexthop != rs.Primary.Attributes["nexthop"] {
			return fmt.Errorf("nexthop Not found in created attributes")
		}
		if foundTransPeer.ReachableCidr != rs.Primary.Attributes["reachable_cidr"] {
			return fmt.Errorf("reachable_cidr Not found in created attributes")
		}
		*transpeer = *foundTransPeer

		return nil
	}
}

func testAccTransPeerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_trans_peer" {
			continue
		}
		foundTransPeer := &goaviatrix.TransPeer{
			Source:        rs.Primary.Attributes["source"],
			Nexthop:       rs.Primary.Attributes["nexthop"],
			ReachableCidr: rs.Primary.Attributes["reachable_cidr"],
		}
		_, err := client.GetTransPeer(foundTransPeer)

		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("transpeer still exists")
		}
	}
	return nil
}
