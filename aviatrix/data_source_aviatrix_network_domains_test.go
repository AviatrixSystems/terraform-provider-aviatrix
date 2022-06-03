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

func TestAccDataSourceAviatrixNetworkDomains_basic(t *testing.T) {
	rName := acctest.RandString(5)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tgwName := acctest.RandStringFromCharSet(5, charset) + acctest.RandString(5)
	awsSideAsNumber := "64512"
	ndName := acctest.RandStringFromCharSet(5, charset) + acctest.RandString(5)
	resourceName := "data.aviatrix_network_domains.test"

	skipAcc := os.Getenv("SKIP_DATA_NETWORK_DOMAINS")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source All Network Domains tests as SKIP_DATA_NETWORK_DOMAINS is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixNetworkDomainsConfigBasic(rName, tgwName, awsSideAsNumber, ndName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTgwNetworkDomainExists("aviatrix_aws_tgw_network_domain.test", tgwName, ndName),
				),
				Destroy: false,
			},
			{
				Config: testAccDataSourceAviatrixNetworkDomainsConfigBasic(rName, tgwName, awsSideAsNumber, ndName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "network_domain_list.0.tgw_name", tgwName),
					resource.TestCheckResourceAttr(resourceName, "network_domain_list.3.account", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "network_domain_list.3.name", ndName),
					resource.TestCheckResourceAttr(resourceName, "network_domain_list.2.cloud_type", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "network_domain_list.1.region", "us-west-1"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixNetworkDomainsConfigBasic(rName string, tgwName string, awsSideAsNumber string, ndName string) string {
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
resource "aviatrix_aws_tgw_network_domain" "test" {
	name       = "%s"
	tgw_name   = aviatrix_aws_tgw.test.tgw_name
	depends_on = [
    	aviatrix_aws_tgw_network_domain.Default_Domain,
    	aviatrix_aws_tgw_network_domain.Shared_Service_Domain,
    	aviatrix_aws_tgw_network_domain.Aviatrix_Edge_Domain
  ]
}
data "aviatrix_network_domains" "test"{

}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		awsSideAsNumber, tgwName, ndName)
}

func testAccCheckTgwNetworkDomainExists(resourceName string, tgwName string, ndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		nd := &goaviatrix.SecurityDomain{
			Name:       ndName,
			AwsTgwName: tgwName,
		}

		_, err := client.GetSecurityDomain(nd)
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("network domain %s not found", ndName)
		}

		return nil
	}
}
