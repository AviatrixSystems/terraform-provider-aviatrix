package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixAwsTgwConnect_basic(t *testing.T) {
	if os.Getenv("SKIP_AWS_TGW_CONNECT") == "yes" {
		t.Skip("Skipping AWS TGW Connect test as SKIP_AWS_TGW_CONNECT is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_aws_tgw_connect.test_aws_tgw_connect"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, " Set SKIP_AWS_TGW_CONNECT to skip this test.")
			if os.Getenv("AWS_REGION") == "" {
				t.Fatalf("")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsTgwConnectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsTgwConnectBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsTgwConnectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", "aws-tgw-"+rName),
					resource.TestCheckResourceAttr(resourceName, "connection_name", "aws-tgw-connect-"+rName),
					resource.TestCheckResourceAttr(resourceName, "network_domain_name", "Shared_Service_Domain"),
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

func testAccAwsTgwConnectBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "aviatrix_aws_tgw" "test_aws_tgw" {
	account_name       = aviatrix_account.aws.account_name
	aws_side_as_number = "64512"
	region             = "%[3]s"
	tgw_name           = "aws-tgw-%[2]s"

	cidrs = ["10.0.0.0/24", "10.1.0.0/24", "8.0.0.0/24", "5.0.0.0/24"]
}
resource "aviatrix_aws_tgw_network_domain" "Default_Domain" {
	name     = "Default_Domain"
	tgw_name = aviatrix_aws_tgw.test_aws_tgw.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Shared_Service_Domain" {
	name     = "Shared_Service_Domain"
	tgw_name = aviatrix_aws_tgw.test_aws_tgw.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Aviatrix_Edge_Domain" {
	name     = "Aviatrix_Edge_Domain"
	tgw_name = aviatrix_aws_tgw.test_aws_tgw.tgw_name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn1" {
	tgw_name1    = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Default_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn2" {
	tgw_name1    = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn3" {
	tgw_name1    = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Default_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
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
	tgw_name            = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	region              = "%[3]s"
	network_domain_name = "Shared_Service_Domain"
	vpc_account_name    = aviatrix_account.aws.account_name
	vpc_id              = aviatrix_vpc.tgw_attach_vpc.vpc_id
}
resource "aviatrix_aws_tgw_connect" "test_aws_tgw_connect" {
	tgw_name            = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	connection_name     = "aws-tgw-connect-%[2]s"
	transport_vpc_id    = aviatrix_aws_tgw_vpc_attachment.aws_tgw_vpc_attachment.vpc_id
	network_domain_name = aviatrix_aws_tgw_vpc_attachment.aws_tgw_vpc_attachment.network_domain_name
}
`, testAccAccountConfigAWS(acctest.RandInt()), rName, os.Getenv("AWS_REGION"))
}

func testAccCheckAwsTgwConnectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aws_tgw_connect Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aws_tgw_connect ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		c := &goaviatrix.AwsTgwConnect{
			ConnectionName: rs.Primary.Attributes["connection_name"],
			TgwName:        rs.Primary.Attributes["tgw_name"],
		}

		foundConn, err := client.GetTGWConnect(context.Background(), c)
		if err != nil {
			return err
		}
		if foundConn.ID() != rs.Primary.ID {
			return fmt.Errorf("aws_tgw_connect not found")
		}

		return nil
	}
}

func testAccCheckAwsTgwConnectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_connect" {
			continue
		}
		foundAwsTgwConnect := &goaviatrix.AwsTgwConnect{
			ConnectionName: rs.Primary.Attributes["connection_name"],
			TgwName:        rs.Primary.Attributes["tgw_name"],
		}
		_, err := client.GetTGWConnect(context.Background(), foundAwsTgwConnect)
		if err == nil {
			return fmt.Errorf("aws_tgw_connect still exists")
		}
	}

	return nil
}
