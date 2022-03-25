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
					resource.TestCheckResourceAttr(resourceName, "ip_filter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_filter.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "tag_filter.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tag_filter.k1", "v1"),
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
	name       = "test-app-domain"
	ip_filter  = [
		"10.0.0.0/16"
	]
	tag_filter = {
		k1 = "v1"
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
		//client := testAccProvider.Meta().(*goaviatrix.Client)

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
	//client := testAccProvider.Meta().(*goaviatrix.Client)

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
