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

func TestAccAviatrixSegmentationNetworkDomainConnectionPolicy_basic(t *testing.T) {
	if os.Getenv("SKIP_SEGMENTATION_NETWORK_DOMAIN_CONNECTION_POLICY") == "yes" {
		t.Skip("Skipping segmentation network domain conn policy test as SKIP_SEGMENTATION_NETWORK_DOMAIN_CONNECTION_POLICY is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_segmentation_network_domain_connection_policy.test_segmentation_network_domain_connection_policy"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentationNetworkDomainConnectionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentationNetworkDomainConnectionPolicyBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentationNetworkDomainConnectionPolicyExists(resourceName),
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

func testAccSegmentationNetworkDomainConnectionPolicyBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_segmentation_network_domain" "test_segmentation_network_domain_1" {
	domain_name = "segmentation-sd-1-%[1]s"
}

resource "aviatrix_segmentation_network_domain" "test_segmentation_network_domain_2" {
	domain_name = "segmentation-sd-2-%[1]s"
}

resource "aviatrix_segmentation_network_domain_connection_policy" "test_segmentation_network_domain_connection_policy" {
	domain_name_1 = aviatrix_segmentation_network_domain.test_segmentation_network_domain_1.domain_name
	domain_name_2 = aviatrix_segmentation_network_domain.test_segmentation_network_domain_2.domain_name
}
`, rName)
}

func testAccCheckSegmentationNetworkDomainConnectionPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("segmentation_network_domain_connection_policy Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no segmentation_network_domain_connection_policy ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundSegmentationNetworkDomainConnectionPolicy := &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
			Domain1: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_1"],
			},
			Domain2: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_2"],
			},
		}

		_, err := client.GetSegmentationSecurityDomainConnectionPolicy(foundSegmentationNetworkDomainConnectionPolicy)
		if err != nil {
			return err
		}
		id := foundSegmentationNetworkDomainConnectionPolicy.Domain1.DomainName + "~" + foundSegmentationNetworkDomainConnectionPolicy.Domain2.DomainName
		if id != rs.Primary.ID {
			return fmt.Errorf("segmentation_network_domain_connection_policy not found")
		}

		return nil
	}
}

func testAccCheckSegmentationNetworkDomainConnectionPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_segmentation_network_domain_connection_policy" {
			continue
		}
		foundSegmentationNetworkDomainConnectionPolicy := &goaviatrix.SegmentationSecurityDomainConnectionPolicy{
			Domain1: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_1"],
			},
			Domain2: &goaviatrix.SegmentationSecurityDomain{
				DomainName: rs.Primary.Attributes["domain_name_2"],
			},
		}
		_, err := client.GetSegmentationSecurityDomainConnectionPolicy(foundSegmentationNetworkDomainConnectionPolicy)
		if err == nil {
			return fmt.Errorf("segmentation_network_domain_connection_policy still exists")
		}
	}

	return nil
}
