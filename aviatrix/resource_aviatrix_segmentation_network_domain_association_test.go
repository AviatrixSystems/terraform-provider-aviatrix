package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixSegmentationNetworkDomainAssociation_basic(t *testing.T) {
	if os.Getenv("SKIP_SEGMENTATION_NETWORK_DOMAIN_ASSOCIATION") == "yes" {
		t.Skip("Skipping segmentation network domain association test as SKIP_SEGMENTATION_NETWORK_DOMAIN_ASSOCIATION is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_segmentation_network_domain_association.test_segmentation_network_domain_association"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckSegmentationNetworkDomainAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentationNetworkDomainAssociationBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentationNetworkDomainAssociationExists(resourceName),
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

func testAccSegmentationNetworkDomainAssociationBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "spoke-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}

resource "aviatrix_segmentation_network_domain" "test_segmentation_network_domain" {
	domain_name = "domain-name-%[1]s"
}

resource "aviatrix_segmentation_network_domain_association" "test_segmentation_network_domain_association" {
	network_domain_name = aviatrix_segmentation_network_domain.test_segmentation_network_domain.domain_name
	attachment_name     = aviatrix_spoke_gateway.test_spoke_gateway.gw_name
}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccCheckSegmentationNetworkDomainAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("segmentation_network_domain_association Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no segmentation_network_domain_association ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		association := &goaviatrix.SegmentationSecurityDomainAssociation{
			SecurityDomainName: rs.Primary.Attributes["network_domain_name"],
			AttachmentName:     rs.Primary.Attributes["attachment_name"],
		}

		_, err := client.GetSegmentationSecurityDomainAssociation(association)
		if err != nil {
			return err
		}

		id := association.SecurityDomainName + "~" + association.AttachmentName
		if id != rs.Primary.ID {
			return fmt.Errorf("segmentation_network_domain_association not found")
		}

		return nil
	}
}

func testAccCheckSegmentationNetworkDomainAssociationDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_segmentation_network_domain_association" {
			continue
		}
		association := &goaviatrix.SegmentationSecurityDomainAssociation{
			SecurityDomainName: rs.Primary.Attributes["network_domain_name"],
			AttachmentName:     rs.Primary.Attributes["attachment_name"],
		}
		_, err := client.GetSegmentationSecurityDomainAssociation(association)
		if err == nil {
			return fmt.Errorf("segmentation_network_domain_association still exists")
		}
	}

	return nil
}
