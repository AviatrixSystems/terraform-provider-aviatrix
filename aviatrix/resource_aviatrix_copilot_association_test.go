package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixCopilotAssociation_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_COPILOT_ASSOCIATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Copilot Association test as SKIP_COPILOT_ASSOCIATION is set")
	}
	resourceName := "aviatrix_copilot_association.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckCopilotAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCopilotAssociationBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCopilotAssociationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "copilot_address", "aviatrix.com"),
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

func testAccCopilotAssociationBasic() string {
	return `
resource "aviatrix_copilot_association" "test" {
    copilot_address = "aviatrix.com"
}
`
}

func testAccCheckCopilotAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("copilot association Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no copilot association ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetCopilotAssociationStatus(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get copilot association status: %v", err)
		}

		return nil
	}
}

func testAccCheckCopilotAssociationDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_copilot_association" {
			continue
		}

		_, err := client.GetCopilotAssociationStatus(context.Background())
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("copilot association exists when it should be destroyed")
		}
	}

	return nil
}
