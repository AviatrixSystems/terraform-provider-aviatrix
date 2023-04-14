package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixQosPolicy_basic(t *testing.T) {
	if os.Getenv("SKIP_QOS_POLICY") == "yes" {
		t.Skip("Skipping QoS policy test as SKIP_QOS_POLICY is set")
	}

	resourceName := "aviatrix_qos_policy.test"
	qosClassName := "qos-class-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckQosPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQosPolicyBasic(qosClassName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQosPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "policies.0.name", "qos_policy1"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.dscp_values.0", "1"),
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

func testAccQosPolicyBasic(qosClassName string) string {
	return fmt.Sprintf(`
resource "aviatrix_qos_class" "test" {
	name     = "%s"
	priority = "1"
}
resource "aviatrix_qos_policy" "test" {
	policies {
		name           = "qos_policy1"
		dscp_values    = ["1", "AF11"]
		qos_class_uuid = aviatrix_qos_class.test.uuid
	}
}
 `, qosClassName)
}

func testAccCheckQosPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("qos policy not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetQosPolicy(context.Background())
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("qos policy not found")
		}

		return nil
	}
}

func testAccCheckQosPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_qos_policy" {
			continue
		}

		_, err := client.GetQosPolicy(context.Background())
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("qos policy still exists")
		}
	}

	return nil
}
