package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixAwsTgwSecurityDomain_basic(t *testing.T) {
	rName := acctest.RandString(5)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tgwName := acctest.RandStringFromCharSet(5, charset) + acctest.RandString(5)
	awsSideAsNumber := "64512"
	sdName := acctest.RandStringFromCharSet(5, charset) + acctest.RandString(5)
	resourceName := "aviatrix_aws_tgw_security_domain.test"

	skipAcc := os.Getenv("SKIP_AWS_TGW_SECURITY_DOMAIN")
	if skipAcc == "yes" {
		t.Skip("Skipping AWS TGW SECURITY DOMAIN test as SKIP_AWS_TGW_SECURITY_DOMAIN is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsTgwSecurityDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsTgwSecurityDomainBasic(rName, tgwName, awsSideAsNumber, sdName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsTgwSecurityDomainExists(resourceName, tgwName, sdName),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", tgwName),
					resource.TestCheckResourceAttr(resourceName, "name", sdName),
				),
			},
		},
	})
}

func testAccAwsTgwSecurityDomainBasic(rName string, tgwName string, awsSideAsNumber string, sdName string) string {
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
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		awsSideAsNumber, tgwName, sdName)
}

func testAccCheckAwsTgwSecurityDomainExists(resourceName string, tgwName string, sdName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		sd := &goaviatrix.SecurityDomain{
			Name:       sdName,
			AwsTgwName: tgwName,
		}

		_, err := client.GetSecurityDomain(sd)
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("security domain %s not found", sdName)
		}

		return nil
	}
}

func testAccCheckAwsTgwSecurityDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_security_domain" {
			continue
		}

		if rs.Primary.Attributes["name"] == "Default_Domain" || rs.Primary.Attributes["name"] == "Shared_Service_Domain" ||
			rs.Primary.Attributes["name"] == "Aviatrix_Edge_Domain" {
			continue
		}

		awsTgw := &goaviatrix.AWSTgw{
			Name: rs.Primary.Attributes["tgw_name"],
		}

		_, err := client.ListTgwDetails(awsTgw)

		if err != goaviatrix.ErrNotFound {
			sd := &goaviatrix.SecurityDomain{
				Name:       rs.Primary.Attributes["name"],
				AwsTgwName: rs.Primary.Attributes["tgw_name"],
			}

			_, err := client.GetSecurityDomain(sd)
			if err != goaviatrix.ErrNotFound {
				return fmt.Errorf("security domain still exists: %v", err)
			}
		} else {
			break
		}
	}

	return nil
}
