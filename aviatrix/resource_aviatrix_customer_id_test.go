package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func preCustomerIDCheck(t *testing.T, msgCommon string) string {
	preAccountCheck(t, msgCommon)

	customerId := os.Getenv("CUSTOMER_ID")
	if customerId == "" {
		t.Fatal("customer ID is not set" + msgCommon)
	}

	return customerId
}

func TestAccAviatrixCustomerID_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CUSTOMER_ID")
	if skipAcc == "yes" {
		t.Skip("Skipping customer ID test as SKIP_CUSTOMER_ID is set")
	}
	msgCommon := ". Set SKIP_CUSTOMER_ID to yes to skip Customer ID tests"
	customerId := preCustomerIDCheck(t, msgCommon)

	resourceName := "aviatrix_customer_id.test_customer_id"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerIDDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomerIDConfigBasic(customerId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomerIDExists("aviatrix_customer_id.test_customer_id"),
					resource.TestCheckResourceAttr(
						resourceName, "customer_id",
						os.Getenv("CUSTOMER_ID")),
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

func testAccCustomerIDConfigBasic(customerId string) string {
	return fmt.Sprintf(`
resource"aviatrix_customer_id" "test_customer_id" {
    customer_id = "%s"
}
	`, customerId)
}

func testAccCheckCustomerIDExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("customer ID Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no customer ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundCustomerID, err := client.GetCustomerID()

		if err != nil {
			return err
		}

		if foundCustomerID != rs.Primary.Attributes["customer_id"] {
			return fmt.Errorf("customer ID not found")
		}

		return nil
	}
}

func testAccCheckCustomerIDDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_customer_id" {
			continue
		}
		_, err := client.GetCustomerID()
		if err != nil {
			return fmt.Errorf("could not retrieve Customer ID due to err: %v", err)
		}
	}
	return nil
}
