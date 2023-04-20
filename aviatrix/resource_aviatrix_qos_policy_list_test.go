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

func TestAccAviatrixQosPolicyList_basic(t *testing.T) {
	if os.Getenv("SKIP_QOS_POLICY_LIST") == "yes" {
		t.Skip("Skipping QoS policy test as SKIP_QOS_POLICY_LIST is set")
	}

	resourceName := "aviatrix_qos_policy_list.test"
	qosClassName := "qos-class-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckQosPolicyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQosPolicyListBasic(qosClassName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQosPolicyListExists(resourceName),
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

func testAccQosPolicyListBasic(qosClassName string) string {
	return fmt.Sprintf(`
resource "aviatrix_qos_class" "test" {
	name     = "%s"
	priority = "1"
}
resource "aviatrix_qos_policy_list" "test" {
	policies {
		name           = "qos_policy1"
		dscp_values    = ["1", "AF11"]
		qos_class_uuid = aviatrix_qos_class.test.uuid
	}
}
 `, qosClassName)
}

func testAccCheckQosPolicyListExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("qos policy list not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetQosPolicyList(context.Background())
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("qos policy list not found")
		}

		return nil
	}
}

func testAccCheckQosPolicyListDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_qos_policy_list" {
			continue
		}

		_, err := client.GetQosPolicyList(context.Background())
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("qos policy list still exists")
		}
	}

	return nil
}
