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

func TestAccAviatrixSpokeTransitAttachment_basic(t *testing.T) {
	var spokeTransitAttachment goaviatrix.SpokeTransitAttachment

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_transit_attachment.test"

	skipAcc := os.Getenv("SKIP_SPOKE_TRANSIT_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping spoke transit attachment tests as 'SKIP_SPOKE_TRANSIT_ATTACHMENT' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_SPOKE_TRANSIT_ATTACHMENT' to 'yes' to skip spoke transit attachment tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeTransitAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeTransitAttachmentConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeTransitAttachmentExists(resourceName, &spokeTransitAttachment),
					resource.TestCheckResourceAttr(resourceName, "spoke_gw_name", fmt.Sprintf("tfs-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "transit_gw_name", fmt.Sprintf("tft-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "tunnel_count", "4"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpokeTransitAttachmentConfigEdge("acc-"+rName, "spoke-"+rName, os.Getenv("PWD"), "transit-"+rName, "site-"+rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeTransitAttachmentExists(resourceName, &spokeTransitAttachment),
					resource.TestCheckResourceAttr(resourceName, "spoke_gw_name", "spoke-"+rName),
					resource.TestCheckResourceAttr(resourceName, "transit_gw_name", "transit-"+rName),
					resource.TestCheckResourceAttr(resourceName, "tunnel_count", "4"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_logical_ifnames.0", "wan1"),
				),
			},
		},
	})
}

func testAccSpokeTransitAttachmentConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	cloud_type         = 1
	account_name       = "tfa-%s"
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test" {
	cloud_type           = 1
	account_name         = aviatrix_account.test.account_name
	region               = "us-west-1"
	name                 = "aws-vpc-test-0"
	cidr                 = "16.0.0.0/20"
	aviatrix_transit_vpc = true
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type     = 1
	account_name   = aviatrix_account.test.account_name
	gw_name        = "tft-%s"
	vpc_id         = aviatrix_vpc.test.vpc_id
	vpc_reg        = aviatrix_vpc.test.region
	gw_size        = "c5.xlarge"
	insane_mode    = true
	subnet         = join(".", [join(".", slice(split(".", aviatrix_vpc.test.public_subnets[0].cidr), 0, 3)), "128/26"]) #"16.0.0.128/26"
	insane_mode_az = "us-west-1b"
}
resource "aviatrix_vpc" "test1" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	cidr         = "173.31.0.0/20"
	name         = "aws-vpc-test-1"
	region       = "us-east-1"
}
resource "aviatrix_spoke_gateway" "test" {
	cloud_type     = 1
	account_name   = aviatrix_account.test.account_name
	gw_name        = "tfs-%s"
	vpc_id         = aviatrix_vpc.test1.vpc_id
	vpc_reg        = aviatrix_vpc.test1.region
	gw_size        = "c5.xlarge"
	insane_mode    = true
	subnet         = join(".", [join(".", slice(split(".", aviatrix_vpc.test1.public_subnets[1].cidr), 0, 2)), "12.0/26"]) #"173.31.12.0/26"
	insane_mode_az = "us-east-1a"
}
resource "aviatrix_spoke_transit_attachment" "test" {
	spoke_gw_name   = aviatrix_spoke_gateway.test.id
	transit_gw_name = aviatrix_transit_gateway.test.id
	tunnel_count    = 4
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, rName)
}

func testAccSpokeTransitAttachmentConfigEdge(accountName, spokeGwName, path, transitGwName, transitSiteID string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_aws" {
	cloud_type         = 1
	account_name       = "aws-%s"
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_account" "test_acc_edge_megaport" {
	account_name       = "megaport-%s"
	cloud_type         = 1048576
}

resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_aws.account_name
	region               = "us-west-1"
	name                 = "aws-vpc-test-1"
	cidr                 = "16.0.0.0/20"
	aviatrix_transit_vpc = true
}

resource "aviatrix_spoke_gateway" "test_spoke" {
	cloud_type     = 1
	account_name   = aviatrix_account.test.account_name
	gw_name        = "tfs-%s"
	vpc_id         = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg        = aviatrix_vpc.test_vpc.region
	gw_size        = "c5.xlarge"
	insane_mode    = true
	subnet         = join(".", [join(".", slice(split(".", aviatrix_vpc.test1.public_subnets[1].cidr), 0, 2)), "12.0/26"]) #"173.31.12.0/26"
	insane_mode_az = "us-east-1a"
}

resource "aviatrix_transit_gateway" "test_edge_transit" {
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

resource "aviatrix_spoke_transit_attachment" "test" {
	spoke_gw_name   = aviatrix_spoke_gateway.test_spoke.gw_name
	transit_gw_name = aviatrix_transit_gateway.test_edge_transit.gw_name
	tunnel_count    = 4
	transit_gateway_logical_ifnames = ["wan1"]

}
	`, accountName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), accountName, spokeGwName, transitGwName, transitSiteID, path)
}

func testAccCheckSpokeTransitAttachmentExists(n string, spokeTransitAttachment *goaviatrix.SpokeTransitAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke transit attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke transit attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundSpokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}
		foundSpokeTransitAttachment2, err := client.GetSpokeTransitAttachment(foundSpokeTransitAttachment)
		if err != nil {
			return err
		}
		if foundSpokeTransitAttachment2.SpokeGwName+"~"+foundSpokeTransitAttachment2.TransitGwName != rs.Primary.ID {
			return fmt.Errorf("spoke transit attachment not found")
		}

		*spokeTransitAttachment = *foundSpokeTransitAttachment2
		return nil
	}
}

func testAccCheckSpokeTransitAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_transit_attachment" {
			continue
		}

		foundSpokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}

		_, err := client.GetSpokeTransitAttachment(foundSpokeTransitAttachment)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("spoke transit attachment still exists %s", err.Error())
		}
	}

	return nil
}
