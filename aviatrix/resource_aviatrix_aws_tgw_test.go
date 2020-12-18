package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAWSTgw_basic(t *testing.T) {
	var awsTgw goaviatrix.AWSTgw

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_aws_tgw.aws_tgw_test"

	skipAcc := os.Getenv("SKIP_AWS_TGW")
	if skipAcc == "yes" {
		t.Skip("Skipping AWS TGW test as SKIP_AWS_TGW is set")
	}
	msg := ". Set SKIP_AWS_TGW to yes to skip AWS TGW  tests"

	awsSideAsNumber := "64512"
	sDm := "zSecurityDomain"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSTgwDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSTgwConfigBasic(rName, awsSideAsNumber, sDm),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSTgwExists(resourceName, &awsTgw),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", fmt.Sprintf("tft-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfaa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "aws_side_as_number", awsSideAsNumber),
					resource.TestCheckResourceAttr(resourceName, "attached_aviatrix_transit_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "attached_aviatrix_transit_gateway.0", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "security_domains.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.0.security_domain_name", "Aviatrix_Edge_Domain"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.0.connected_domains.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.0.connected_domains.0", "Default_Domain"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.0.connected_domains.1", "Shared_Service_Domain"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.0.connected_domains.2", sDm),
					resource.TestCheckResourceAttr(resourceName, "security_domains.1.security_domain_name", "Default_Domain"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.1.connected_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.1.connected_domains.0", "Aviatrix_Edge_Domain"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.2.security_domain_name", "Shared_Service_Domain"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.2.connected_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.2.connected_domains.0", "Aviatrix_Edge_Domain"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.3.security_domain_name", sDm),
					resource.TestCheckResourceAttr(resourceName, "security_domains.3.connected_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.3.connected_domains.0", "Aviatrix_Edge_Domain"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.3.attached_vpc.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_domains.3.attached_vpc.0.vpc_id", os.Getenv("AWS_VPC_ID2")),
				),
			},
		},
	})
}

func testAccAWSTgwConfigBasic(rName string, awsSideAsNumber string, sDm string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account1" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_account" "test_account2" {
	account_name       = "tfaa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "transit_gw_test" {
	cloud_type               = 1
	account_name             = aviatrix_account.test_account1.account_name
	gw_name                  = "tfg-%s"
	vpc_id                   = "%s"
	vpc_reg                  = "%s"
	gw_size                  = "t2.micro"
	subnet                   = "%s"
	enable_hybrid_connection = true
}
resource "aviatrix_aws_tgw" "aws_tgw_test" {
	account_name                      = aviatrix_account.test_account2.account_name
	attached_aviatrix_transit_gateway = [aviatrix_transit_gateway.transit_gw_test.gw_name]
	aws_side_as_number                = "%s"
	manage_vpc_attachment             = true
	region                            = "%s"
	tgw_name                          = "tft-%s"

	security_domains {
		connected_domains    = [
			"Default_Domain",
			"Shared_Service_Domain",
			"%s"
		]
		security_domain_name = "Aviatrix_Edge_Domain"
	}
	security_domains {
		connected_domains    = [
			"Aviatrix_Edge_Domain"
		]
		security_domain_name = "Default_Domain"
	}
	security_domains {
		connected_domains    = [
			"Aviatrix_Edge_Domain"
		]
		security_domain_name = "Shared_Service_Domain"
	}
	security_domains {
		connected_domains    = [
			"Aviatrix_Edge_Domain"
		]
		security_domain_name = "%s"
		attached_vpc {
			vpc_account_name = aviatrix_account.test_account2.account_name
			vpc_id           = "%s"
			vpc_region       = "%s"
		}
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"),
		awsSideAsNumber, os.Getenv("AWS_REGION"), rName, sDm, sDm, os.Getenv("AWS_VPC_ID2"),
		os.Getenv("AWS_REGION2"))
}

func testAccCheckAWSTgwExists(n string, awsTgw *goaviatrix.AWSTgw) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("AWS TGW Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no AWS TGW ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAwsTgw := &goaviatrix.AWSTgw{
			Name: rs.Primary.Attributes["tgw_name"],
		}

		foundAwsTgw2, err := client.GetAWSTgw(foundAwsTgw)
		if err != nil {
			return err
		}
		if foundAwsTgw2.Name != rs.Primary.ID {
			return fmt.Errorf("AWS TGW not found")
		}

		*awsTgw = *foundAwsTgw
		return nil
	}
}

func testAccCheckAWSTgwDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw" {
			continue
		}

		foundAWSTgw := &goaviatrix.AWSTgw{
			Name: rs.Primary.Attributes["tgw_name"],
		}

		_, err := client.GetAWSTgw(foundAWSTgw)
		if err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return nil
			}
			return fmt.Errorf("AWS TGW still exists: %v", err)
		}

		return fmt.Errorf("AWS TGW still exists")
	}

	return nil
}
