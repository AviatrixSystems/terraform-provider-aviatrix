package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAwsTgwConnectPeer_basic(t *testing.T) {
	if os.Getenv("SKIP_AWS_TGW_CONNECT_PEER") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_AWS_TGW_CONNECT_PEER is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_aws_tgw_connect_peer.test_aws_tgw_connect_peer"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsTgwConnectPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsTgwConnectPeerBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsTgwConnectPeerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", "aws-tgw-"+rName),
					resource.TestCheckResourceAttr(resourceName, "connection_name", "aws-tgw-connect-"+rName),
					resource.TestCheckResourceAttr(resourceName, "connect_peer_name", "connect-peer-"+rName),
					resource.TestCheckResourceAttr(resourceName, "peer_as_number", "65001"),
					resource.TestCheckResourceAttr(resourceName, "peer_gre_address", "172.31.1.11"),
					resource.TestCheckResourceAttr(resourceName, "bgp_inside_cidrs.0", "169.254.6.0/29"),
					resource.TestCheckResourceAttr(resourceName, "tgw_gre_address", "10.0.0.32"),
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

func testAccAwsTgwConnectPeerBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "aviatrix_aws_tgw" "test_aws_tgw" {
  account_name                      = aviatrix_account.aws.account_name
  aws_side_as_number                = "64512"
  region                            = "%[3]s"
  tgw_name                          = "aws-tgw-%[2]s"
  manage_vpc_attachment             = false
  manage_transit_gateway_attachment = false

  cidrs = ["10.0.0.0/24", "10.1.0.0/24", "8.0.0.0/24", "5.0.0.0/24"]

  security_domains {
    connected_domains    = [
      "Default_Domain",
      "Shared_Service_Domain"
    ]
    security_domain_name = "Aviatrix_Edge_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain",
      "Shared_Service_Domain"
    ]
    security_domain_name = "Default_Domain"
  }

  security_domains {
    connected_domains    = [
      "Aviatrix_Edge_Domain",
      "Default_Domain"
    ]
    security_domain_name = "Shared_Service_Domain"
  }
}

resource aviatrix_vpc tgw_attach_vpc {
  cloud_type           = aviatrix_account.aws.cloud_type
  account_name         = aviatrix_account.aws.account_name
  region               = "%[3]s"
  name                 = "tgw-attach-vpc-%[2]s"
  cidr                 = "10.10.0.0/16"
  aviatrix_firenet_vpc = false
  aviatrix_transit_vpc = false
}

resource "aviatrix_aws_tgw_vpc_attachment" "aws_tgw_vpc_attachment" {
  tgw_name             = aviatrix_aws_tgw.test_aws_tgw.tgw_name
  region               = "%[3]s"
  security_domain_name = "Shared_Service_Domain"
  vpc_account_name     = aviatrix_account.aws.account_name
  vpc_id               = aviatrix_vpc.tgw_attach_vpc.vpc_id
}
resource "aviatrix_aws_tgw_connect" "test_aws_tgw_connect" {
  tgw_name             = aviatrix_aws_tgw.test_aws_tgw.tgw_name
  connection_name      = "aws-tgw-connect-%[2]s"
  attachment_name      = aviatrix_aws_tgw_vpc_attachment.aws_tgw_vpc_attachment.vpc_id
  security_domain_name = aviatrix_aws_tgw_vpc_attachment.aws_tgw_vpc_attachment.security_domain_name
}

resource "aviatrix_aws_tgw_connect_peer" "test_aws_tgw_connect_peer" {
  tgw_name              = aviatrix_aws_tgw.test_aws_tgw.tgw_name
  connection_name       = aviatrix_aws_tgw_connect.test_aws_tgw_connect.connection_name
  connect_peer_name     = "connect-peer-%[2]s"
  connect_attachment_id = aviatrix_aws_tgw_connect.test_aws_tgw_connect.connect_attachment_id
  peer_as_number        = "65001"
  peer_gre_address      = "172.31.1.11"
  bgp_inside_cidrs      = ["169.254.6.0/29"]
  tgw_gre_address       = "10.0.0.32"
}
`, testAccAccountConfigAWS(acctest.RandInt()), rName, os.Getenv("AWS_REGION"))
}

func testAccCheckAwsTgwConnectPeerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aws_tgw_connect_peer Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aws_tgw_connect_peer ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAwsTgwConnectPeer := &goaviatrix.AwsTgwConnectPeer{
			ConnectionName:  rs.Primary.Attributes["connection_name"],
			TgwName:         rs.Primary.Attributes["tgw_name"],
			ConnectPeerName: rs.Primary.Attributes["connect_peer_name"],
		}

		_, err := client.GetTGWConnectPeer(context.Background(), foundAwsTgwConnectPeer)
		if err != nil {
			return err
		}
		if foundAwsTgwConnectPeer.ID() != rs.Primary.ID {
			return fmt.Errorf("aws_tgw_connect_peer not found")
		}

		return nil
	}
}

func testAccCheckAwsTgwConnectPeerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_connect_peer" {
			continue
		}
		foundAwsTgwConnectPeer := &goaviatrix.AwsTgwConnectPeer{
			ConnectionName:  rs.Primary.Attributes["connection_name"],
			TgwName:         rs.Primary.Attributes["tgw_name"],
			ConnectPeerName: rs.Primary.Attributes["connect_peer_name"],
		}
		_, err := client.GetTGWConnectPeer(context.Background(), foundAwsTgwConnectPeer)
		if err == nil {
			return fmt.Errorf("aws_tgw_connect_peer still exists")
		}
	}

	return nil
}
