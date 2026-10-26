package aviatrix

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixDcfMitmCa_basic(t *testing.T) {
	resourceName := "data.aviatrix_dcf_mitm_ca.test"

	skipAcc := os.Getenv("SKIP_DATA_DCF_MITM_CA")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source DCF MITM CA test as SKIP_DATA_DCF_MITM_CA is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFMitmCaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixDcfMitmCaConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixDcfMitmCa(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-mitm-ca-data-source"),
					resource.TestCheckResourceAttrSet(resourceName, "ca_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ca_hash"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_chain"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "origin"),
				),
			},
		},
	})
}

func TestAccDataSourceAviatrixDcfMitmCa_notFound(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DATA_DCF_MITM_CA")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source DCF MITM CA test as SKIP_DATA_DCF_MITM_CA is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAviatrixDcfMitmCaConfigNotFound(),
				ExpectError: regexp.MustCompile(`DCF MITM CA with name .* not found`),
			},
		},
	})
}

func testAccDataSourceAviatrixDcfMitmCaConfigBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_mitm_ca" "test" {
	name              = "test-dcf-mitm-ca-data-source"
	key               = %q
	certificate_chain = %q
}

data "aviatrix_dcf_mitm_ca" "test" {
	depends_on = [aviatrix_dcf_mitm_ca.test]
	name       = "test-dcf-mitm-ca-data-source"
}
`, testMitmCaPrivateKey(), testMitmCaCertificate())
}

func testAccDataSourceAviatrixDcfMitmCaConfigNotFound() string {
	return `
data "aviatrix_dcf_mitm_ca" "test" {
	name = "non-existent-mitm-ca-12345"
}
`
}

func testAccDataSourceAviatrixDcfMitmCa(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
