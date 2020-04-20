package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
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

		client := testAccProvider.Meta().(*goaviatrix.Client)
		foundAwsTgwPeering := &goaviatrix.AwsTgwPeering{
			TgwName1: rs.Primary.Attributes["tgw_name1"],
			TgwName2: rs.Primary.Attributes["tgw_name2"],
		}
		err := client.GetAwsTgwPeering(foundAwsTgwPeering)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("no aviatrix AWS tgw peering is found")
			}
			return err
		}

		*awsTgwPeering = *foundAwsTgwPeering
		return nil
	}
}

func testAccCheckAWSTgwPeeringDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_peering" {
			continue
		}

		foundAwsTgwPeering := &goaviatrix.AwsTgwPeering{
			TgwName1: rs.Primary.Attributes["tgw_name1"],
			TgwName2: rs.Primary.Attributes["tgw_name2"],
		}

		err := client.GetAwsTgwPeering(foundAwsTgwPeering)
		if err != goaviatrix.ErrNotFound {
			if strings.Contains(err.Error(), "does not exist") {
				return nil
			}
			return fmt.Errorf("aviatrix AWS tgw peering still exists")
		}
	}

	return nil
}
