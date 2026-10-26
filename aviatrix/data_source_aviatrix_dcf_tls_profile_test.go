package aviatrix

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixDcfTLSProfile_basic(t *testing.T) {
	resourceName := "data.aviatrix_dcf_tls_profile.test"

	skipAcc := os.Getenv("SKIP_DATA_DCF_TLS_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source DCF TLS Profile test as SKIP_DATA_DCF_TLS_PROFILE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfTLSProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixDcfTLSProfileConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixDcfTLSProfile(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-dcf-tls-profile-data-source"),
					resource.TestCheckResourceAttr(resourceName, "certificate_validation", "CERTIFICATE_VALIDATION_NONE"),
					resource.TestCheckResourceAttr(resourceName, "verify_sni", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
		},
	})
}

func TestAccDataSourceAviatrixDcfTLSProfile_notFound(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DATA_DCF_TLS_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source DCF TLS Profile test as SKIP_DATA_DCF_TLS_PROFILE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAviatrixDcfTLSProfileConfigNotFound(),
				ExpectError: regexp.MustCompile(`DCF TLS Profile with display_name .* not found`),
			},
		},
	})
}

func testAccDataSourceAviatrixDcfTLSProfileConfigBasic() string {
	return `
resource "aviatrix_dcf_tls_profile" "test" {
	display_name           = "test-dcf-tls-profile-data-source"
	certificate_validation = "CERTIFICATE_VALIDATION_NONE"
	verify_sni             = false
}

data "aviatrix_dcf_tls_profile" "test" {
	depends_on   = [aviatrix_dcf_tls_profile.test]
	display_name = "test-dcf-tls-profile-data-source"
}
`
}

func testAccDataSourceAviatrixDcfTLSProfileConfigNotFound() string {
	return `
data "aviatrix_dcf_tls_profile" "test" {
	display_name = "non-existent-tls-profile-12345"
}
`
}

func testAccDataSourceAviatrixDcfTLSProfile(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
