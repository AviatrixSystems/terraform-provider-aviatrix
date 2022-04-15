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

func TestAccAviatrixSegmentationNetworkDomain_basic(t *testing.T) {
	if os.Getenv("SKIP_SEGMENTATION_NETWORK_DOMAIN") == "yes" {
		t.Skip("Skipping segmentation network domain test as SKIP_SEGMENTATION_NETWORK_DOMAIN is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_segmentation_network_domain.test_segmentation_network_domain"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentationNetworkDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentationNetworkDomainBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentationNetworkDomainExists(resourceName),
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

func testAccSegmentationNetworkDomainBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_segmentation_network_domain" "test_segmentation_network_domain" {
	domain_name = "segmentation-nd-%s"
}
`, rName)
}

func testAccCheckSegmentationNetworkDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("segmentation_network_domain Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no segmentation_network_domain ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundSegmentationNetworkDomain := &goaviatrix.SegmentationSecurityDomain{
			DomainName: rs.Primary.Attributes["domain_name"],
		}

		_, err := client.GetSegmentationSecurityDomain(foundSegmentationNetworkDomain)
		if err != nil {
			return err
		}
		if foundSegmentationNetworkDomain.DomainName != rs.Primary.ID {
			return fmt.Errorf("segmentation_network_domain not found")
		}

		return nil
	}
}

func testAccCheckSegmentationNetworkDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_segmentation_network_domain" {
			continue
		}
		foundSegmentationNetworkDomain := &goaviatrix.SegmentationSecurityDomain{
			DomainName: rs.Primary.Attributes["domain_name"],
		}
		_, err := client.GetSegmentationSecurityDomain(foundSegmentationNetworkDomain)
		if err == nil {
			return fmt.Errorf("segmentation_network_domain still exists")
		}
	}

	return nil
}
