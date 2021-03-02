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

func TestAccAviatrixAwsTgwTransitGatewayAttachment_basic(t *testing.T) {
	var awsTgwTransitGwAttachment goaviatrix.AwsTgwTransitGwAttachment

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_aws_tgw_transit_gateway_attachment.test"

	skipAcc := os.Getenv("SKIP_AWS_TGW_TRANSIT_GATEWAY_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping AWS tgw transit gateway test as 'SKIP_AWS_TGW_TRANSIT_GATEWAY_ATTACHMENT' is set")
	}

	msg := ". Set 'SKIP_AWS_TGW_TRANSIT_GATEWAY_ATTACHMENT' to 'yes' to skip AWS tgw transit gateway tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsTgwTransitGatewayAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsTgwTransitGatewayAttachmentConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAwsTgwTransitGatewayAttachmentExists(resourceName, &awsTgwTransitGwAttachment),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", fmt.Sprintf("tft-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "vpc_account_name", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name", fmt.Sprintf("tfg-%s", rName)),
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

func testAccAwsTgwTransitGatewayAttachmentConfigBasic(rName string) string {
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
	aws_side_as_number                = "64512"
	manage_vpc_attachment             = false
	manage_transit_gateway_attachment = false
	region                            = "%s"
	tgw_name                          = "tft-%s"

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
resource "aviatrix_transit_gateway" "test" {
	cloud_type               = 1
	account_name             = aviatrix_account.test.account_name
	gw_name                  = "tfg-%s"
	vpc_id                   = "%s"
	vpc_reg                  = aviatrix_aws_tgw.test.region
	gw_size                  = "c5.xlarge"
	subnet                   = "%s"
	enable_active_mesh       = true
	enable_hybrid_connection = true
	connected_transit        = true
}
resource "aviatrix_aws_tgw_transit_gateway_attachment" "test" {
	tgw_name             = aviatrix_aws_tgw.test.tgw_name
	region               = aviatrix_aws_tgw.test.region
	vpc_account_name     = aviatrix_transit_gateway.test.account_name
	vpc_id               = aviatrix_transit_gateway.test.vpc_id
	transit_gateway_name = aviatrix_transit_gateway.test.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName, rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_SUBNET"))
}

func tesAccCheckAwsTgwTransitGatewayAttachmentExists(n string, awsTgwTransitGwAttachment *goaviatrix.AwsTgwTransitGwAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("AWS tgw transit gatway attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no AWS tgw transit gatway attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundAwsTgwTransitGwAttachment := &goaviatrix.AwsTgwTransitGwAttachment{
			TgwName: rs.Primary.Attributes["tgw_name"],
			VpcID:   rs.Primary.Attributes["vpc_id"],
		}

		foundAwsTgwTransitGwAttachment2, err := client.GetAwsTgwTransitGwAttachment(foundAwsTgwTransitGwAttachment)
		if err != nil {
			return err
		}
		if foundAwsTgwTransitGwAttachment2.TgwName != rs.Primary.Attributes["tgw_name"] {
			return fmt.Errorf("'tgw_name' Not found in created attributes")
		}
		if foundAwsTgwTransitGwAttachment2.VpcID != rs.Primary.Attributes["vpc_id"] {
			return fmt.Errorf("'vpc_id' Not found in created attributes")
		}

		*awsTgwTransitGwAttachment = *foundAwsTgwTransitGwAttachment2
		return nil
	}
}

func testAccCheckAwsTgwTransitGatewayAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw_transit_gateway_attachment" {
			continue
		}

		foundAwsTgwTransitGwAttachment := &goaviatrix.AwsTgwTransitGwAttachment{
			TgwName: rs.Primary.Attributes["tgw_name"],
			VpcID:   rs.Primary.Attributes["vpc_id"],
		}
		_, err := client.GetAwsTgwTransitGwAttachment(foundAwsTgwTransitGwAttachment)
		if err == nil {
			return fmt.Errorf("aviatrix AWS tgw transit gateway attachment still exists: %s", err.Error())
		}

		return nil
	}

	return nil
}
