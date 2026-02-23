package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixTrafficClassifier_basic(t *testing.T) {
	if os.Getenv("SKIP_TRAFFIC_CLASSIFIER") == "yes" {
		t.Skip("Skipping traffic classifier test as SKIP_TRAFFIC_CLASSIFIER is set")
	}

	resourceName := "aviatrix_traffic_classifier.test"
	policyName := "policy-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficClassifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTrafficClassifierBasic(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficClassifierExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "policies.0.name", policyName),
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

func testAccTrafficClassifierBasic(policyName string) string {
	return fmt.Sprintf(`
resource "aviatrix_smart_group" "src" {
	name = "src"
	selector {
		match_expressions {
      		cidr = "10.0.0.0/16"
    	}
  	}
}

resource "aviatrix_smart_group" "dst" {
	name = "dst"
	selector {
    	match_expressions {
      		cidr = "11.0.0.0/16"
    	}
  	}
}

resource "aviatrix_traffic_classifier" "test" {
	policies {
    	name                          = "%s"
    	source_smart_group_uuids      = [aviatrix_smart_group.src.uuid]
    	destination_smart_group_uuids = [aviatrix_smart_group.dst.uuid]
  	}
}
 `, policyName)
}

func testAccCheckTrafficClassifierExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("traffic classifier not found: %s", resourceName)
		}

		client := mustClient(testAccProvider.Meta())

		_, err := client.GetTrafficClassifier(context.Background())
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("traffic classifier not found")
		}

		return nil
	}
}

func testAccCheckTrafficClassifierDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_traffic_classifier" {
			continue
		}

		_, err := client.GetTrafficClassifier(context.Background())
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("traffic classifier still exists")
		}
	}

	return nil
}
