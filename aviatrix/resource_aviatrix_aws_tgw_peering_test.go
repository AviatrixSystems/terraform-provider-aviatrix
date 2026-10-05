package aviatrix

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAWSTgwPeering_basic(t *testing.T) {
	var awsTgwPeering goaviatrix.AwsTgwPeering
	rName := acctest.RandString(5)
	resourceName := "aviatrix_aws_tgw_peering.test"

	skipAcc := os.Getenv("SKIP_AWS_TGW_PEERING")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix AWS tgw peering tests as 'SKIP_AWS_TGW_PEERING' is set")
	}
	msgCommon := ". Set 'SKIP_AWS_TGW_PEERING' to 'yes' to skip Aviatrix AWS tgw peering tests"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSTgwPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSTgwPeeringConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAWSTgwPeeringExists(resourceName, &awsTgwPeering),
					resource.TestCheckResourceAttr(resourceName, "tgw_name1", "tgw1"),
					resource.TestCheckResourceAttr(resourceName, "tgw_name2", "tgw2"),
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

func testAccAWSTgwPeeringConfigBasic(rName string) string {
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
	account_name       = aviatrix_account.test.account_name
	aws_side_as_number = "64512"
	region             = "us-east-1"
	tgw_name           = "tgw1"
}
resource "aviatrix_aws_tgw_network_domain" "test1_Default_Domain" {
	name     = "Default_Domain"
	tgw_name = aviatrix_aws_tgw.test1.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "test1_Shared_Service_Domain" {
	name     = "Shared_Service_Domain"
	tgw_name = aviatrix_aws_tgw.test1.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "test1_Aviatrix_Edge_Domain" {
	name     = "Aviatrix_Edge_Domain"
	tgw_name = aviatrix_aws_tgw.test1.tgw_name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "test1_default_sd_conn1" {
	tgw_name1    = aviatrix_aws_tgw.test1.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.test1_Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test1.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.test1_Default_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "test1_default_sd_conn2" {
	tgw_name1    = aviatrix_aws_tgw.test1.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.test1_Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test1.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.test1_Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "test1_default_sd_conn3" {
	tgw_name1    = aviatrix_aws_tgw.test1.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.test1_Default_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test1.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.test1_Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw" "test2" {
	account_name       = aviatrix_account.test.account_name
	aws_side_as_number = "64512"
	region             = "us-east-2"
	tgw_name           = "tgw2"
}
resource "aviatrix_aws_tgw_network_domain" "test2_Default_Domain" {
	name     = "Default_Domain"
	tgw_name = aviatrix_aws_tgw.test2.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "test2_Shared_Service_Domain" {
	name     = "Shared_Service_Domain"
	tgw_name = aviatrix_aws_tgw.test2.tgw_name
}
resource "aviatrix_aws_tgw_network_domain" "test2_Aviatrix_Edge_Domain" {
	name     = "Aviatrix_Edge_Domain"
	tgw_name = aviatrix_aws_tgw.test2.tgw_name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "test2_default_sd_conn1" {
	tgw_name1    = aviatrix_aws_tgw.test2.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.test2_Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test2.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.test2_Default_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "test2_default_sd_conn2" {
	tgw_name1    = aviatrix_aws_tgw.test2.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.test2_Aviatrix_Edge_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test2.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.test2_Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_peering_domain_conn" "test2_default_sd_conn3" {
	tgw_name1    = aviatrix_aws_tgw.test2.tgw_name
	domain_name1 = aviatrix_aws_tgw_network_domain.test2_Default_Domain.name
	tgw_name2    = aviatrix_aws_tgw.test2.tgw_name
	domain_name2 = aviatrix_aws_tgw_network_domain.test2_Shared_Service_Domain.name
}
resource "aviatrix_aws_tgw_peering" "test" {
	tgw_name1 = aviatrix_aws_tgw.test1.tgw_name
	tgw_name2 = aviatrix_aws_tgw.test2.tgw_name
	depends_on = [
		aviatrix_aws_tgw_network_domain.test1_Default_Domain,
		aviatrix_aws_tgw_network_domain.test1_Shared_Service_Domain,
		aviatrix_aws_tgw_network_domain.test1_Aviatrix_Edge_Domain,
		aviatrix_aws_tgw_network_domain.test2_Default_Domain,
		aviatrix_aws_tgw_network_domain.test2_Shared_Service_Domain,
		aviatrix_aws_tgw_network_domain.test2_Aviatrix_Edge_Domain
	]
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"))
}

func tesAccCheckAWSTgwPeeringExists(n string, awsTgwPeering *goaviatrix.AwsTgwPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix AWS tgw peering Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no aviatrix AWS tgw peering ID is set")
		}

		client := mustClient(testAccProvider.Meta())
		foundAwsTgwPeering := &goaviatrix.AwsTgwPeering{
			TgwName1: rs.Primary.Attributes["tgw_name1"],
			TgwName2: rs.Primary.Attributes["tgw_name2"],
		}
		err := client.GetAwsTgwPeering(foundAwsTgwPeering)
		if err != nil {
			if errors.Is(err, goaviatrix.ErrNotFound) {
				return fmt.Errorf("no aviatrix AWS tgw peering is found")
			}
			return err
		}

		*awsTgwPeering = *foundAwsTgwPeering
		return nil
	}
}

func testAccCheckAWSTgwPeeringDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_peering" {
			continue
		}

		foundAwsTgwPeering := &goaviatrix.AwsTgwPeering{
			TgwName1: rs.Primary.Attributes["tgw_name1"],
			TgwName2: rs.Primary.Attributes["tgw_name2"],
		}

		err := client.GetAwsTgwPeering(foundAwsTgwPeering)
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("aviatrix AWS tgw peering still exists")
		}
	}

	return nil
}
