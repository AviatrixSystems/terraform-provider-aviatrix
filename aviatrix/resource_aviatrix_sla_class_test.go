package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixSLAClass_basic(t *testing.T) {
	if os.Getenv("SKIP_SLA_CLASS") == "yes" {
		t.Skip("Skipping sla class test as SKIP_SLA_CLASS is set")
	}

	resourceName := "aviatrix_sla_class.test"
	slaClassName := "sla-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSLAClassDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSLAClassBasic(slaClassName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSLAClassExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", slaClassName),
					resource.TestCheckResourceAttr(resourceName, "latency", "43"),
					resource.TestCheckResourceAttr(resourceName, "jitter", "1"),
					resource.TestCheckResourceAttr(resourceName, "packet_drop_rate", "3"),
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

func testAccSLAClassBasic(slaClassName string) string {
	return fmt.Sprintf(`
resource "aviatrix_sla_class" "test" {
	name             = "%s"
	latency          = 43
	jitter           = 1
	packet_drop_rate = 3
}
 `, slaClassName)
}

func testAccCheckSLAClassExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("sla class not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetSLAClass(context.Background(), rs.Primary.Attributes["uuid"])
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("sla class not found")
		}

		return nil
	}
}

func testAccCheckSLAClassDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_sla_class" {
			continue
		}

		_, err := client.GetSLAClass(context.Background(), rs.Primary.Attributes["uuid"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("sla class still exists")
		}
	}

	return nil
}
