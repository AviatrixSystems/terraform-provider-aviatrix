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

func TestAccAviatrixAppDomain_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_APP_DOMAIN")
	if skipAcc == "yes" {
		t.Skip("Skipping App Domain test as SKIP_APP_DOMAIN is set")
	}
	resourceName := "aviatrix_app_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccAppDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppDomainBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-app-domain"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.0.cidr", "11.0.0.0/16"),

					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.type", "vm"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.account_name", "mlin-aviatrix"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.region", "us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "selector.0.match_expressions.1.tags.k3", "v3"),
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

func testAccAppDomainBasic() string {
	return `
resource "aviatrix_app_domain" "test" {
	name = "test-app-domain"

	selector {
		match_expressions {
			cidr = "11.0.0.0/16"
		}

		match_expressions {
			type         = "vm"
			account_name = "mlin-aviatrix"
			region       = "us-west-2"
			tags         = {
				k3 = "v3"
			}
		}
	}
}
`
}

func testAccCheckAppDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no App Domain resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no App Domain ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		appDomain, err := client.GetAppDomain(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get App Domain status: %v", err)
		}

		if appDomain.UUID != rs.Primary.ID {
			return fmt.Errorf("app domain ID not found")
		}

		return nil
	}
}

func testAccAppDomainDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_app_domain" {
			continue
		}

		_, err := client.GetAppDomain(context.Background(), rs.Primary.ID)
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("app domain configured when it should be destroyed")
		}
	}

	return nil
}
