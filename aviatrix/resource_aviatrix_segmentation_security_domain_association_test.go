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

func TestAccAviatrixSegmentationSecurityDomainAssociation_basic(t *testing.T) {
	if os.Getenv("SKIP_SEGMENTATION_SECURITY_DOMAIN_ASSOCIATION") == "yes" {
		t.Skip("Skipping segmentation security domain association test as SKIP_SEGMENTATION_SECURITY_DOMAIN_ASSOCIATION is set")
	}

	rName := acctest.RandString(5)
	msg := ". Set SKIP_SEGMENTATION_SECURITY_DOMAIN_ASSOCIATION to yes to skip segmentation security domain association tests."
	resourceName := "aviatrix_segmentation_security_domain_association.test_segmentation_security_domain_association"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, msg)
			preGateway2Check(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentationSecurityDomainAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentationSecurityDomainAssociationBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentationSecurityDomainAssociationExists(resourceName),
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

func testAccSegmentationSecurityDomainAssociationBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test_transit_gateway" {
	cloud_type          = 1
	account_name        = aviatrix_account.test_acc_aws.account_name
	gw_name             = "transit-aws-%[1]s"
	vpc_id              = "%[5]s"
	vpc_reg             = "%[6]s"
	gw_size             = "t2.micro"
	subnet              = "%[7]s"
	enable_active_mesh  = true
	enable_segmentation = true
	connected_transit   = true
}

resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type         = 1
	account_name       = aviatrix_account.test_acc_aws.account_name
	gw_name            = "spoke-aws-%[1]s"
	vpc_id             = "%[8]s"
	vpc_reg            = "%[9]s"
	gw_size            = "t2.micro"
	subnet             = "%[10]s"
	transit_gw         = aviatrix_transit_gateway.test_transit_gateway.gw_name
	enable_active_mesh = true
}

resource "aviatrix_segmentation_security_domain" "test_segmentation_security_domain" {
	domain_name = "domain-name-%[1]s"
}

resource "aviatrix_segmentation_security_domain_association" "test_segmentation_security_domain_association" {
	transit_gateway_name = aviatrix_transit_gateway.test_transit_gateway.gw_name
	security_domain_name = aviatrix_segmentation_security_domain.test_segmentation_security_domain.domain_name
	attachment_name      = aviatrix_spoke_gateway.test_spoke_gateway.gw_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"),
		os.Getenv("AWS_VPC_ID2"), os.Getenv("AWS_REGION2"), os.Getenv("AWS_SUBNET2"))
}

func testAccCheckSegmentationSecurityDomainAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("segmentation_security_domain_association Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no segmentation_security_domain_association ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		association := &goaviatrix.SegmentationSecurityDomainAssociation{
			TransitGatewayName: rs.Primary.Attributes["transit_gateway_name"],
			SecurityDomainName: rs.Primary.Attributes["security_domain_name"],
			AttachmentName:     rs.Primary.Attributes["attachment_name"],
		}

		_, err := client.GetSegmentationSecurityDomainAssociation(association)
		if err != nil {
			return err
		}

		id := association.TransitGatewayName + "~" + association.SecurityDomainName + "~" + association.AttachmentName
		if id != rs.Primary.ID {
			return fmt.Errorf("segmentation_security_domain_association not found")
		}

		return nil
	}
}

func testAccCheckSegmentationSecurityDomainAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_segmentation_security_domain_association" {
			continue
		}
		association := &goaviatrix.SegmentationSecurityDomainAssociation{
			TransitGatewayName: rs.Primary.Attributes["transit_gateway_name"],
			SecurityDomainName: rs.Primary.Attributes["security_domain_name"],
			AttachmentName:     rs.Primary.Attributes["attachment_name"],
		}
		_, err := client.GetSegmentationSecurityDomainAssociation(association)
		if err == nil {
			return fmt.Errorf("segmentation_security_domain_association still exists")
		}
	}

	return nil
}
