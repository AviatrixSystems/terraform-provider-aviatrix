package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	account_name       = "%[1]s"
	cloud_type         = 1
	aws_account_number = "%[2]s"
	aws_iam            = false
	aws_access_key     = "%[3]s"
	aws_secret_key     = "%[4]s"
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
	enable_active_mesh       = true
	enable_firenet           = true
	enable_hybrid_connection = true
}
resource "aviatrix_aws_tgw" "test" {
	account_name                      = aviatrix_account.test.account_name
	aws_side_as_number                = "64512"
	manage_vpc_attachment             = false
	manage_transit_gateway_attachment = false
	region                            = aviatrix_vpc.test.region
	tgw_name                          = "%[5]s"
	security_domains {
		security_domain_name = "Aviatrix_Edge_Domain"
		connected_domains    = [
		  "Default_Domain",
		  "Shared_Service_Domain",
		]
	}
	security_domains {
		security_domain_name = "Default_Domain"
		connected_domains    = [
		  "Aviatrix_Edge_Domain",
		  "Shared_Service_Domain"
		]
	}
	security_domains {
		security_domain_name = "Shared_Service_Domain"
		connected_domains    = [
		  "Aviatrix_Edge_Domain",
		  "Default_Domain"
		]
	}
	security_domains {
		security_domain_name = "%[7]s"
		aviatrix_firewall    = true
		connected_domains    = [
		  "%[6]s"
		]
	}
	security_domains {
		security_domain_name = "%[6]s"
		connected_domains    = [
		  "%[7]s"
		]
	}
}
resource "aviatrix_aws_tgw_vpc_attachment" "test" {
	tgw_name             = aviatrix_aws_tgw.test.tgw_name
	region               = aviatrix_vpc.test.region
	security_domain_name = "%[7]s"
	vpc_account_name     = aviatrix_vpc.test.account_name
	vpc_id               = aviatrix_vpc.test.vpc_id
   	depends_on = [aviatrix_transit_gateway.test]
}
resource "aviatrix_aws_tgw_intra_domain_inspection" "test" {
	tgw_name             = aviatrix_aws_tgw.test.tgw_name
	route_domain_name    = "%[6]s"
	firewall_domain_name = "%[7]s"
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

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		intraDomainInspection := &goaviatrix.IntraDomainInspection{
			TgwName:         tgwName,
			RouteDomainName: domainName,
		}

		err := client.GetIntraDomainInspectionStatus(context.Background(), intraDomainInspection)
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("aws tgw intra domain inspection disabled")
		}
		if err != nil {
			return fmt.Errorf("failed to get aws tgw intra domain inspection status: %v", err)
		}

		return nil
	}
}

func testAccCheckAwsTgwIntraDomainInspectionDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

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
