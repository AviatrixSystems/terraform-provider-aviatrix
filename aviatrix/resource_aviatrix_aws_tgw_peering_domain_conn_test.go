package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAWSTgwPeeringDomainConn_basic(t *testing.T) {
	var domainConn goaviatrix.DomainConn
	rName := acctest.RandString(5)
	resourceName := "aviatrix_aws_tgw_peering_domain_conn.test"

	skipAcc := os.Getenv("SKIP_AWS_TGW_PEERING_DOMAIN_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix AWS tgw peering domain connection tests as 'SKIP_AWS_TGW_PEERING_DOMAIN_CONN' is set")
	}
	msgCommon := ". Set 'SKIP_AWS_TGW_PEERING_DOMAIN_CONN' to 'yes' to skip Aviatrix AWS tgw peering domain connection tests"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSTgwPeeringDomainConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSTgwPeeringDomainConnConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAWSTgwPeeringDomainConnExists(resourceName, &domainConn),
					resource.TestCheckResourceAttr(resourceName, "tgw_name1", "tgw1"),
					resource.TestCheckResourceAttr(resourceName, "domain_name1", "Default_Domain"),
					resource.TestCheckResourceAttr(resourceName, "tgw_name2", "tgw2"),
					resource.TestCheckResourceAttr(resourceName, "domain_name2", "Default_Domain"),
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

func testAccAWSTgwPeeringDomainConnConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_tgw" "test1" {
    account_name          = aviatrix_account.test.account_name
    aws_side_as_number    = "64512"
    manage_vpc_attachment = false
    region                = "us-east-1"
    tgw_name              = "tgw1"

    security_domains {
        connected_domains    = [
            "Default_Domain",
            "Shared_Service_Domain",
        ]
        security_domain_name = "Aviatrix_Edge_Domain"
    }
    security_domains {
        connected_domains    = [
            "Aviatrix_Edge_Domain",
            "Shared_Service_Domain",
        ]
        security_domain_name = "Default_Domain"
    }
    security_domains {
        connected_domains    = [
            "Aviatrix_Edge_Domain",
            "Default_Domain",
        ]
        security_domain_name = "Shared_Service_Domain"
    }
}
resource "aviatrix_aws_tgw" "test2" {
    account_name          = aviatrix_account.test.account_name
    aws_side_as_number    = "64512"
    manage_vpc_attachment = false
    region                = "us-east-2"
    tgw_name              = "tgw2"

    security_domains {
        connected_domains    = [
            "Default_Domain",
            "Shared_Service_Domain",
        ]
        security_domain_name = "Aviatrix_Edge_Domain"
    }
    security_domains {
        connected_domains    = [
            "Aviatrix_Edge_Domain",
            "Shared_Service_Domain",
        ]
        security_domain_name = "Default_Domain"
    }
    security_domains {
        connected_domains    = [
            "Aviatrix_Edge_Domain",
            "Default_Domain",
        ]
        security_domain_name = "Shared_Service_Domain"
    }
}
resource "aviatrix_aws_tgw_peering" "test" {
	tgw_name1 = aviatrix_aws_tgw.test1.tgw_name
	tgw_name2 = aviatrix_aws_tgw.test2.tgw_name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "test" {
	tgw_name1    = aviatrix_aws_tgw_peering.test.tgw_name1
	domain_name1 = "Default_Domain"
	tgw_name2    = aviatrix_aws_tgw_peering.test.tgw_name2
	domain_name2 = "Default_Domain"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func tesAccCheckAWSTgwPeeringDomainConnExists(n string, domainConn *goaviatrix.DomainConn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix AWS tgw peering domain connection Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix AWS tgw peering domain connection ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)
		foundDomainConn := &goaviatrix.DomainConn{
			TgwName1:    rs.Primary.Attributes["tgw_name1"],
			DomainName1: rs.Primary.Attributes["domain_name1"],
			TgwName2:    rs.Primary.Attributes["tgw_name2"],
			DomainName2: rs.Primary.Attributes["domain_name2"],
		}
		err := client.GetDomainConn(foundDomainConn)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("no aviatrix AWS tgw peering domain connection is found")
			}
			return err
		}

		*domainConn = *foundDomainConn
		return nil
	}
}

func testAccCheckAWSTgwPeeringDomainConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_peering_domain_conn" {
			continue
		}

		foundDomainConn := &goaviatrix.DomainConn{
			TgwName1:    rs.Primary.Attributes["tgw_name1"],
			DomainName1: rs.Primary.Attributes["domain_name1"],
			TgwName2:    rs.Primary.Attributes["tgw_name2"],
			DomainName2: rs.Primary.Attributes["domain_name2"],
		}

		err := client.GetDomainConn(foundDomainConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("aviatrix AWS tgw peering domain connection still exists")
		}
	}

	return nil
}
