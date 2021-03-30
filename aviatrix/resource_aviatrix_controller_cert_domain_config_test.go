package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixControllerCertDomainConfig_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CONTROLLER_CERT_DOMAIN_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Cert Domain Config test as SKIP_CONTROLLER_CERT_DOMAIN_CONFIG is set")
	}

	rName := acctest.RandString(5) + ".com"
	resourceName := "aviatrix_controller_cert_domain_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckControllerCertDomainConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControllerCertDomainConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControllerCertDomainConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cert_domain", rName),
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

func testAccControllerCertDomainConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_controller_cert_domain_config" "test" {
  cert_domain = "%s"
}
`, rName)
}

func testAccCheckControllerCertDomainConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("controller cert domain config ID Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("controller cert domain config ID is not set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetCertDomain(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get cert domain config status: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("controller cert domain config ID not found")
		}

		return nil
	}
}

func testAccCheckControllerCertDomainConfigDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_controller_cert_domain_config" {
			continue
		}

		certDomainConfig, _ := client.GetCertDomain(context.Background())
		if !certDomainConfig.IsDefault {
			return fmt.Errorf("controller cert domain configured when it should be destroyed")
		}
	}

	return nil
}
