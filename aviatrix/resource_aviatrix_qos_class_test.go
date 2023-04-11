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

func TestAccAviatrixQosClass_basic(t *testing.T) {
	if os.Getenv("SKIP_QOS_CLASS") == "yes" {
		t.Skip("Skipping QoS class test as SKIP_QOS_CLASS is set")
	}

	resourceName := "aviatrix_qos_class.test"
	qosClassName := "qos-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckQosClassDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQosClassBasic(qosClassName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQosClassExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", qosClassName),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
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

func testAccQosClassBasic(qosClassName string) string {
	return fmt.Sprintf(`
resource "aviatrix_qos_class" "test" {
	name     = "%s"
	priority = "1"
}
 `, qosClassName)
}

func testAccCheckQosClassExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("qos class not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetQosClass(context.Background(), rs.Primary.Attributes["uuid"])
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("qos class not found")
		}

		return nil
	}
}

func testAccCheckQosClassDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_qos_class" {
			continue
		}

		_, err := client.GetQosClass(context.Background(), rs.Primary.Attributes["uuid"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("qos class still exists")
		}
	}

	return nil
}
