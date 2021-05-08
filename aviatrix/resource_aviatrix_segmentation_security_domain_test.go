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

func TestAccAviatrixSegmentationSecurityDomain_basic(t *testing.T) {
	if os.Getenv("SKIP_SEGMENTATION_SECURITY_DOMAIN") == "yes" {
		t.Skip("Skipping segmentation security domain test as SKIP_SEGMENTATION_SECURITY_DOMAIN is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_segmentation_security_domain.test_segmentation_security_domain"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentationSecurityDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentationSecurityDomainBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentationSecurityDomainExists(resourceName),
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

func testAccSegmentationSecurityDomainBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_segmentation_security_domain" "test_segmentation_security_domain" {
	domain_name = "segmentation-sd-%s"
}
`, rName)
}

func testAccCheckSegmentationSecurityDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("segmentation_security_domain Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no segmentation_security_domain ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundSegmentationSecurityDomain := &goaviatrix.SegmentationSecurityDomain{
			DomainName: rs.Primary.Attributes["domain_name"],
		}

		_, err := client.GetSegmentationSecurityDomain(foundSegmentationSecurityDomain)
		if err != nil {
			return err
		}
		if foundSegmentationSecurityDomain.DomainName != rs.Primary.ID {
			return fmt.Errorf("segmentation_security_domain not found")
		}

		return nil
	}
}

func testAccCheckSegmentationSecurityDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_segmentation_security_domain" {
			continue
		}
		foundSegmentationSecurityDomain := &goaviatrix.SegmentationSecurityDomain{
			DomainName: rs.Primary.Attributes["domain_name"],
		}
		_, err := client.GetSegmentationSecurityDomain(foundSegmentationSecurityDomain)
		if err == nil {
			return fmt.Errorf("segmentation_security_domain still exists")
		}
	}

	return nil
}
