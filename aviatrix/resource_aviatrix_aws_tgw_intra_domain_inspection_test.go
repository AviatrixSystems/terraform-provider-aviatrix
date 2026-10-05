package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAwsTgwIntraDomainInspection_basic(t *testing.T) {
	accName := "acc-" + acctest.RandString(5)
	tgwName := "tgw-" + acctest.RandString(5)
	routeDomainName := "sd-" + acctest.RandString(5)
	firewallDomainName := "sd-" + acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_AWS_TGW_INTRA_DOMAIN_INSPECTION")
	if skipAcc == "yes" {
		t.Skip("Skipping Aws Tgw Intra Domain Inspection test as SKIP_AWS_TGW_INTRA_DOMAIN_INSPECTION is set")
	}
	resourceName := "aviatrix_aws_tgw_intra_domain_inspection.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckAwsTgwIntraDomainInspectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsTgwIntraDomainInspectionBasic(accName, tgwName, routeDomainName, firewallDomainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsTgwIntraDomainInspectionExists(resourceName, tgwName, routeDomainName),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", tgwName),
					resource.TestCheckResourceAttr(resourceName, "route_domain_name", routeDomainName),
					resource.TestCheckResourceAttr(resourceName, "firewall_domain_name", firewallDomainName),
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

func testAccAwsTgwIntraDomainInspectionBasic(accName string, tgwName string, routeDomainName string, firewallDomainName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test" {
 	cloud_type           = 1
 	account_name         = aviatrix_account.test.account_name
 	region               = "us-west-1"
	name                 = "firenet-vpc"
 	cidr                 = "10.0.0.0/16"
 	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type               = 1
	account_name             = aviatrix_account.test.account_name
	gw_name                  = "transit"
	vpc_id                   = aviatrix_vpc.test.vpc_id
	vpc_reg                  = aviatrix_vpc.test.region
	gw_size                  = "c5.xlarge"
	subnet                   = "10.0.0.0/28"
	enable_firenet           = true
	enable_hybrid_connection = true
}
resource "aviatrix_aws_tgw" "test" {
	account_name       = aviatrix_account.test.account_name
	aws_side_as_number = "64512"
	region             = aviatrix_vpc.test.region
	tgw_name           = "%s"
}
resource "aviatrix_aws_tgw_network_domain" "Default_Domain" {
	name     = "Default_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Shared_Service_Domain" {
	name     = "Shared_Service_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "Aviatrix_Edge_Domain" {
	name     = "Aviatrix_Edge_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn1" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Default_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn2" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "default_nd_conn3" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.Default_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_network_domain" "route_domain" {
	name       = "%s"
	tgw_name   = aviatrix_aws_tgw.test.tgw_name
	depends_on = [
    	aviatrix_aws_tgw_network_domain.Default_Domain,
    	aviatrix_aws_tgw_network_domain.Shared_Service_Domain,
    	aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain
	]
}
resource "aviatrix_aws_tgw_network_domain" "firewall_domain" {
	name              = "%s"
	tgw_name          = aviatrix_aws_tgw.test.tgw_name
	aviatrix_firewall = true
	depends_on = [
    	aviatrix_aws_tgw_network_domain.Default_Domain,
    	aviatrix_aws_tgw_network_domain.Shared_Service_Domain,
    	aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain
	]
}
resource "aviatrix_aws_tgw_peering_domain_conn" "nd_conn" {
	tgw_name1    = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.firewall_domain.name
	tgw_name2    = aviatrix_aws_tgw.test.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.route_domain.name
}
resource "aviatrix_aws_tgw_vpc_attachment" "test" {
	tgw_name            = aviatrix_aws_tgw.test.tgw_name
	region              = aviatrix_vpc.test.region
	network_domain_name = aviatrix_aws_tgw_network_domain.firewall_domain.name
	vpc_account_name    = aviatrix_vpc.test.account_name
	vpc_id              = aviatrix_vpc.test.vpc_id
   	depends_on = [aviatrix_transit_gateway.test]
}
resource "aviatrix_aws_tgw_intra_domain_inspection" "test" {
	tgw_name             = aviatrix_aws_tgw.test.tgw_name
	route_domain_name    = aviatrix_aws_tgw_network_domain.route_domain.name
	firewall_domain_name = aviatrix_aws_tgw_network_domain.firewall_domain.name
	depends_on = [aviatrix_aws_tgw_vpc_attachment.test]
}
	`, accName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		tgwName, routeDomainName, firewallDomainName)
}

func testAccCheckAwsTgwIntraDomainInspectionExists(resourceName string, tgwName string, domainName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("aws tgw intra domain inspection ID Not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("aws tgw intra domain inspection ID is not set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		intraDomainInspection := &goaviatrix.IntraDomainInspection{
			TgwName:         tgwName,
			RouteDomainName: domainName,
		}

		err := client.GetIntraDomainInspectionStatus(context.Background(), intraDomainInspection)
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("aws tgw intra domain inspection disabled")
		}
		if err != nil {
			return fmt.Errorf("failed to get aws tgw intra domain inspection status: %w", err)
		}

		return nil
	}
}

func testAccCheckAwsTgwIntraDomainInspectionDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_intra_domain_inspection" {
			continue
		}

		intraDomainInspection := &goaviatrix.IntraDomainInspection{
			TgwName:         rs.Primary.Attributes["tgw_name"],
			RouteDomainName: rs.Primary.Attributes["route_domain_name"],
		}

		err := client.GetIntraDomainInspectionStatus(context.Background(), intraDomainInspection)

		if err == nil {
			return fmt.Errorf("aws tgw intra domain inspection still exists")
		}
	}

	return nil
}
