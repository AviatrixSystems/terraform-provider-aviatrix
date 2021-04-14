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

func TestAccAviatrixAwsTgwSecurityDomainConnection_basic(t *testing.T) {
	rName := acctest.RandString(5)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tgwName := acctest.RandStringFromCharSet(5, charset) + acctest.RandString(5)
	awsSideAsNumber := "64512"
	sdName := acctest.RandStringFromCharSet(5, charset) + acctest.RandString(5)
	resourceName := "aviatrix_aws_tgw_security_domain.test"

	skipAcc := os.Getenv("SKIP_AWS_TGW_SECURITY_DOMAIN_CONNECTION")
	if skipAcc == "yes" {
		t.Skip("Skipping AWS TGW SECURITY DOMAIN CONNECTION test as SKIP_AWS_TGW_SECURITY_DOMAIN_CONNECTION is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsTgwSecurityDomainConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsTgwSecurityDomainConnectionBasic(rName, tgwName, awsSideAsNumber, sdName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsTgwSecurityDomainConnectionExists(resourceName, tgwName, sdName),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", tgwName),
					resource.TestCheckResourceAttr(resourceName, "name", sdName),
				),
			},
		},
	})
}

func testAccAwsTgwSecurityDomainConnectionBasic(rName string, tgwName string, awsSideAsNumber string, sdName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_tgw" "test" {
	account_name                      = aviatrix_account.test.account_name
	aws_side_as_number                = "%s"
	region                            = "us-west-1"
	tgw_name                          = "%s"
	manage_security_domain            = false
	manage_vpc_attachment             = false
	manage_transit_gateway_attachment = false
}
resource "aviatrix_aws_tgw_security_domain" "Default_Domain" {
	name     = "Default_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_security_domain" "Shared_Service_Domain" {
	name     = "Shared_Service_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_security_domain" "Aviatrix_Edge_Domain" {
	name     = "Aviatrix_Edge_Domain"
	tgw_name = aviatrix_aws_tgw.test.tgw_name
}
resource "aviatrix_aws_tgw_security_domain" "test" {
	name       = "%s"
	tgw_name   = aviatrix_aws_tgw.test.tgw_name
	depends_on = [
    	aviatrix_aws_tgw_security_domain.Default_Domain,
    	aviatrix_aws_tgw_security_domain.Shared_Service_Domain,
    	aviatrix_aws_tgw_security_domain.Aviatrix_Edge_Domain
  ]
}
resource "aviatrix_aws_tgw_security_domain_connection" "test" {
	tgw_name     = aviatrix_aws_tgw.test.tgw_name
	domain_name1 = aviatrix_aws_tgw_security_domain.Default_Domain.name
	domain_name2 = aviatrix_aws_tgw_security_domain.test.name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		awsSideAsNumber, tgwName, sdName)
}

func testAccCheckAwsTgwSecurityDomainConnectionExists(resourceName string, tgwName string, sdName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		securityDomain := &goaviatrix.SecurityDomain{
			Name:       sdName,
			AwsTgwName: tgwName,
		}

		securityDomainDetails, err := client.GetSecurityDomainDetails(context.TODO(), securityDomain)
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("couldn't find security domain %s", sdName)
		}
		if err != nil {
			return fmt.Errorf("couldn't get the details of the security domain %s due to %v", sdName, err)
		}

		for _, sd := range securityDomainDetails.ConnectedDomain {
			if sd == "Default_Domain" {
				return nil
			}
		}

		return fmt.Errorf("couldn't find security domain connection")
	}
}

func testAccCheckAwsTgwSecurityDomainConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_security_domain_connection" {
			continue
		}

		securityDomain := &goaviatrix.SecurityDomain{
			Name:       rs.Primary.Attributes["domain_name1"],
			AwsTgwName: rs.Primary.Attributes["tgw_name"],
		}

		securityDomainDetails, err := client.GetSecurityDomainDetails(context.TODO(), securityDomain)
		if err == goaviatrix.ErrNotFound || err != nil {
			return nil
		} else {
			for _, sd := range securityDomainDetails.ConnectedDomain {
				if sd == rs.Primary.Attributes["domain_name2"] {
					return fmt.Errorf("security domain connection still exits")
				}
			}
		}
	}

	return nil
}
