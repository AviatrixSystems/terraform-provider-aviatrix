package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixSegmentationSecurityDomainConnectionPolicy_basic(t *testing.T) {
	if os.Getenv("SKIP_SEGMENTATION_SECURITY_DOMAIN_CONNECTION_POLICY") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_SEGMENTATION_SECURITY_DOMAIN_CONNECTION_POLICY is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_segmentation_security_domain_connection_policy.test_segmentation_security_domain_connection_policy"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentationSecurityDomainConnectionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentationSecurityDomainConnectionPolicyBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentationSecurityDomainConnectionPolicyExists(resourceName),
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

func testAccSegmentationSecurityDomainConnectionPolicyBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_segmentation_security_domain" "test_segmentation_security_domain_1" {
	domain_name = "segmentation-security-domain-1-%[1]s"
}

resource "aviatrix_segmentation_security_domain" "test_segmentation_security_domain_2" {
	domain_name = "segmentation-security-domain-2-%[1]s"
}

resource "aviatrix_segmentation_security_domain_connection_policy" "test_segmentation_security_domain_connection_policy" {
	domain_name_1 = aviatrix_segmentation_security_domain.test_segmentation_security_domain_1.domain_name
	domain_name_2 = aviatrix_segmentation_security_domain.test_segmentation_security_domain_2.domain_name
}
`, rName)
}

func testAccCheckSegmentationSecurityDomainConnectionPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("segmentation_security_domain_connection_policy Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no segmentation_security_domain_connection_policy ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundSegmentationSecurityDomainConnectionPolicy := &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
			Domain1: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_1"],
			},
			Domain2: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_2"],
			},
		}

		_, err := client.GetSegmentationSecurityDomainConnectionPolicy(foundSegmentationSecurityDomainConnectionPolicy)
		if err != nil {
			return err
		}
		id := foundSegmentationSecurityDomainConnectionPolicy.Domain1.DomainName + "~" + foundSegmentationSecurityDomainConnectionPolicy.Domain2.DomainName
		if id != rs.Primary.ID {
			return fmt.Errorf("segmentation_security_domain_connection_policy not found")
		}

		return nil
	}
}

func testAccCheckSegmentationSecurityDomainConnectionPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_segmentation_security_domain_connection_policy" {
			continue
		}
		foundSegmentationSecurityDomainConnectionPolicy := &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
			Domain1: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_1"],
			},
			Domain2: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_2"],
			},
		}
		_, err := client.GetSegmentationSecurityDomainConnectionPolicy(foundSegmentationSecurityDomainConnectionPolicy)
		if err == nil {
			return fmt.Errorf("segmentation_security_domain_connection_policy still exists")
		}
	}

	return nil
}
