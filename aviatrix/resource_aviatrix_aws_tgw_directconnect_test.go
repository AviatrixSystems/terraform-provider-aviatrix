package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAwsTgwDirectConnect_basic(t *testing.T) {
	var awsTgwDirectConnect goaviatrix.AwsTgwDirectConnect

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_aws_tgw_directconnect.test"

	skipAcc := os.Getenv("SKIP_AWS_TGW_DIRECTCONNECT")
	if skipAcc == "yes" {
		t.Skip("Skipping AWS TGW DIRECTCONNECT test as SKIP_AWS_TGW_DIRECTCONNECT is set")
	}

	msg := ". Set SKIP_AWS_TGW_DIRECTCONNECT to yes to skip AWS TGW DIRECTCONNECT tests"

	awsSideAsNumber := "12"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsTgwDirectConnectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsTgwDirectConnectConfigBasic(rName, awsSideAsNumber),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAwsTgwDirectConnectExists(resourceName, &awsTgwDirectConnect),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", fmt.Sprintf("tft-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "directconnect_account_name", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "dx_gateway_id", os.Getenv("AWS_DX_GATEWAY_ID")),
					resource.TestCheckResourceAttr(resourceName, "security_domain_name", "Default_Domain"),
					resource.TestCheckResourceAttr(resourceName, "allowed_prefix", "10.12.0.0/24"),
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

func testAccAwsTgwDirectConnectConfigBasic(rName string, awsSideAsNumber string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam	           = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_tgw" "test_aws_tgw" {
	account_name          = aviatrix_account.test_account.account_name
	aws_side_as_number    = "64513"
	manage_vpc_attachment = false
	region                = "%s"
	tgw_name              = "tft-%s"
	security_domains {
		connected_domains    = [
			"Default_Domain",
			"Shared_Service_Domain"
		]
		security_domain_name = "Aviatrix_Edge_Domain"
	}
	security_domains {
		connected_domains    = [
			"Aviatrix_Edge_Domain",
			"Shared_Service_Domain"
		]
		security_domain_name = "Default_Domain"
	}
	security_domains {
		connected_domains    = [
			"Aviatrix_Edge_Domain",
			"Default_Domain"
		]
		security_domain_name = "Shared_Service_Domain"
	}
}
resource "aviatrix_aws_tgw_directconnect" "test" {
	tgw_name                   = aviatrix_aws_tgw.test_aws_tgw.tgw_name
	directconnect_account_name = aviatrix_account.test_account.account_name
	dx_gateway_id              = "%s"
	security_domain_name       = "Default_Domain"
	allowed_prefix             = "10.12.0.0/24"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName, os.Getenv("AWS_DX_GATEWAY_ID"))
}

func tesAccCheckAwsTgwDirectConnectExists(n string, awsTgwDirectConnect *goaviatrix.AwsTgwDirectConnect) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("AWS TGW DIRECTCONNECT Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no AWS TGW DIRECTCONNECT ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAwsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
			TgwName:     rs.Primary.Attributes["tgw_name"],
			DxGatewayID: rs.Primary.Attributes["dx_gateway_id"],
		}

		foundAwsTgwDirectConnect2, err := client.GetAwsTgwDirectConnect(foundAwsTgwDirectConnect)
		if err != nil {
			return err
		}
		if foundAwsTgwDirectConnect2.TgwName != rs.Primary.Attributes["tgw_name"] {
			return fmt.Errorf("tgw_name Not found in created attributes")
		}
		if foundAwsTgwDirectConnect2.DxGatewayID != rs.Primary.Attributes["dx_gateway_id"] {
			return fmt.Errorf("dx_gateway_id Not found in created attributes")
		}

		*awsTgwDirectConnect = *foundAwsTgwDirectConnect
		return nil
	}
}

func testAccCheckAwsTgwDirectConnectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_directconnect" {
			continue
		}

		foundAwsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
			TgwName:     rs.Primary.Attributes["tgw_name"],
			DxGatewayID: rs.Primary.Attributes["dx_gateway_id"],
		}

		_, err := client.GetAwsTgwDirectConnect(foundAwsTgwDirectConnect)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("aviatrix AWS TGW DIRECTCONNECT still exists")
		}
	}

	return nil
}
