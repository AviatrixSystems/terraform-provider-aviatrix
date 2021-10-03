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

func TestAccAviatrixSpokeTransitAttachment_basic(t *testing.T) {
	var spokeTransitAttachment goaviatrix.SpokeTransitAttachment

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_transit_attachment.test"

	skipAcc := os.Getenv("SKIP_SPOKE_TRANSIT_ATTACHMENT")
	if skipAcc == "yes" {
		t.Skip("Skipping spoke transit attachment tests as 'SKIP_SPOKE_TRANSIT_ATTACHMENT' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_SPOKE_TRANSIT_ATTACHMENT' to 'yes' to skip spoke transit attachment tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeTransitAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeTransitAttachmentConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeTransitAttachmentExists(resourceName, &spokeTransitAttachment),
					resource.TestCheckResourceAttr(resourceName, "spoke_gw_name", fmt.Sprintf("tfs-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "transit_gw_name", fmt.Sprintf("tft-%s", rName)),
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

func testAccSpokeTransitAttachmentConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_spoke_gateway" "test"{
	cloud_type                        = 1
	account_name                      = aviatrix_account.test.account_name
	gw_name                           = "tfs-%s"
	vpc_id                            = "%s"
	vpc_reg                           = "%s"
	gw_size                           = "t2.micro"
	subnet                            = "%s"
	enable_active_mesh                = true
	manage_transit_gateway_attachment = false
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type         = 1
	account_name       = aviatrix_account.test.account_name
	gw_name            = "tft-%s"
	vpc_id             = "%s"
	vpc_reg            = "%s"
	gw_size            = "t2.micro"
	subnet             = "%s"
	enable_active_mesh = true
}
resource "aviatrix_spoke_transit_attachment" "test" {
	spoke_gw_name   = aviatrix_spoke_gateway.test.gw_name
	transit_gw_name = aviatrix_transit_gateway.test.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"),
		rName, os.Getenv("AWS_VPC_ID2"), os.Getenv("AWS_REGION2"), os.Getenv("AWS_SUBNET2"))
}

func testAccCheckSpokeTransitAttachmentExists(n string, spokeTransitAttachment *goaviatrix.SpokeTransitAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke transit attachment Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke transit attachment ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundSpokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}
		foundSpokeTransitAttachment2, err := client.GetSpokeTransitAttachment(foundSpokeTransitAttachment)
		if err != nil {
			return err
		}
		if foundSpokeTransitAttachment2.SpokeGwName+"~"+foundSpokeTransitAttachment2.TransitGwName != rs.Primary.ID {
			return fmt.Errorf("spoke transit attachment not found")
		}

		*spokeTransitAttachment = *foundSpokeTransitAttachment2
		return nil
	}
}

func testAccCheckSpokeTransitAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_transit_attachment" {
			continue
		}

		foundSpokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
			SpokeGwName:   rs.Primary.Attributes["spoke_gw_name"],
			TransitGwName: rs.Primary.Attributes["transit_gw_name"],
		}

		_, err := client.GetSpokeTransitAttachment(foundSpokeTransitAttachment)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("spoke transit attachment still exists %s", err.Error())
		}
	}

	return nil
}
