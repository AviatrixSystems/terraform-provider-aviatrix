package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixFirewallInstanceImages_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_firewall_instance_images.foo"

	skipAcc := os.Getenv("SKIP_DATA_FIREWALL_INSTANCE_IMAGES")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source Firewall Instance Images tests as SKIP_DATA_FIREWALL_INSTANCE_IMAGES is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFirewallInstanceImagesConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixFirewallInstanceImages(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "firewall_images.0.firewall_image"),
				),
			},
		},
	})
}

func testAccDataSourceFirewallInstanceImagesConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%s"
	name                 = "vpc-for-firewall-images"
	cidr                 = "10.10.0.0/24"
	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test_transit_gateway" {
	cloud_type               = aviatrix_vpc.test_vpc.cloud_type
	account_name             = aviatrix_account.test_account.account_name
	gw_name                  = "tftg-%s"
	vpc_id                   = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg                  = aviatrix_vpc.test_vpc.region
	gw_size                  = "c5.xlarge"
	subnet                   = aviatrix_vpc.test_vpc.subnets[0].cidr
	enable_hybrid_connection = true
	enable_firenet           = true
}
data "aviatrix_firewall_instance_images" "foo" {
	vpc_id = aviatrix_transit_gateway.test_transit_gateway.vpc_id
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName)
}

func testAccDataSourceAviatrixFirewallInstanceImages(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
