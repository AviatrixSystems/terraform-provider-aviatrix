package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix/cloud"
)

func TestAccAviatrixGeoVPN_basic(t *testing.T) {
	var geoVPN goaviatrix.GeoVPN

	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipAcc := os.Getenv("SKIP_GEO_VPN")
	if skipAcc == "yes" {
		t.Skip("Skipping Geo VPN test as SKIP_GEO_VPN is set")
	}

	resourceName := "aviatrix_geo_vpn.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set SKIP_GEO_VPN to yes to skip Geo VPN tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGeoVPNDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGeoVPNConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeoVPNExists(resourceName, &geoVPN),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "service_name", "vpn"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", os.Getenv("DOMAIN_NAME")),
					resource.TestCheckResourceAttr(resourceName, "elb_dns_names.#", "1"),
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

func testAccGeoVPNConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "test_vpn_gw" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
	vpn_access   = true
	vpn_cidr     = "192.168.50.0/24"
	max_vpn_conn = "100"
	enable_elb   = true
	elb_name     = "elb-test1"
}
resource "aviatrix_geo_vpn" "foo" {
	cloud_type    = 1
	account_name  = aviatrix_account.test.account_name
	service_name  = "vpn"
	domain_name   = "%[8]s"
	elb_dns_names = [
		aviatrix_gateway.test_vpn_gw.elb_dns_name,
	]
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), os.Getenv("DOMAIN_NAME"))
}

func testAccCheckGeoVPNExists(n string, geoVPN *goaviatrix.GeoVPN) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("GeoVPN Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no GeoVPN ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundGeoVPN := &goaviatrix.GeoVPN{
			CloudType:   cloud.AWS,
			ServiceName: rs.Primary.Attributes["service_name"],
			DomainName:  rs.Primary.Attributes["domain_name"],
		}

		foundGeoVPN2, err := client.GetGeoVPNInfo(foundGeoVPN)
		if err != nil {
			return err
		}
		if foundGeoVPN2.ServiceName+"~"+foundGeoVPN2.DomainName != rs.Primary.ID {
			return fmt.Errorf("GeoVPN not found")
		}

		*geoVPN = *foundGeoVPN
		return nil
	}
}

func testAccCheckGeoVPNDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_geo_vpn" {
			continue
		}
		foundGeoVPN := &goaviatrix.GeoVPN{
			CloudType:   cloud.AWS,
			ServiceName: rs.Primary.Attributes["service_name"],
			DomainName:  rs.Primary.Attributes["domain_name"],
		}

		_, err := client.GetGeoVPNInfo(foundGeoVPN)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("GeoVPN still enabled")
		}
	}

	return nil
}
