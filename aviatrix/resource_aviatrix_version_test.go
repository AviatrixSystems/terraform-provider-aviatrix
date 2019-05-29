package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixVersion_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_VERSION")
	if skipAcc == "yes" {
		t.Skip("Skipping Version test as SKIP_VERSION is set")
	}

	preAccountCheck(t, ". Set SKIP_VERSION to yes to skip admin email tests")

	resourceName := "aviatrix_version.foo"
	importStateVerifyIgnore := []string{"target_version"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVersionConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVersionExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},
		},
	})
}

func testAccVersionConfigBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_version" "foo" {
	target_version = "latest"
}
	`)
}

func testAccCheckVersionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("version Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Version is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		foundVersion, _, err := client.GetCurrentVersion()
		if err != nil {
			return err
		}
		if foundVersion != rs.Primary.Attributes["version"] {
			return fmt.Errorf("version information not found")
		}

		return nil
	}
}

func testAccCheckVersionDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_version" {
			continue
		}
		_, _, err := client.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("could not retrieve Version due to err: %v", err)
		}
	}
	return nil
}
