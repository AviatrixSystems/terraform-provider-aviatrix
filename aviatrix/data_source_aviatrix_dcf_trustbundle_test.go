package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAviatrixDcfTrustbundle_basic(t *testing.T) {
	resourceName := "data.aviatrix_dcf_trustbundle.test"

	skipAcc := os.Getenv("SKIP_DATA_DCF_TRUSTBUNDLE")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source DCF Trust Bundle test as SKIP_DATA_DCF_TRUSTBUNDLE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, ". Set SKIP_DATA_DCF_TRUSTBUNDLE to yes to skip Data Source DCF Trust Bundle tests")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixDcfTrustbundleConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixDcfTrustbundle(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "bundle_id"),
					resource.TestCheckResourceAttrSet(resourceName, "bundle_content"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
		},
	})
}

func testAccDataSourceAviatrixDcfTrustbundleConfigBasic() string {
	return `
resource "aviatrix_dcf_trustbundle" "test_trustbundle" {
	display_name    = "test-trustbundle-data-source"
	bundle_content = <<EOF
-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAKoK/heBjcOuMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTMwODI3MjM1NzUwWhcNMTQwODI3MjM1NzUwWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAwUdHPiQnY1MxrE6iPCtJ3cIw/XhKsjsKlswiynxOzN4tOW9LsAyCChJH
2mx6kRYE5wD/T8Qk3SzZGw7Bcd8sCnvvODyC5XQwNMJkMu2YyLT6yl8DbV/gLIny
LdM6EKVGM8zo7FKsD3ZmzY9JN+sLO4Pwy9w8AGU8C5aA/7MCwqJhYzUfKLQx9Yl6
gcEw5VF8Ma8o72VexSdENQOSW/Gsc8QdpB3VxGhGrMsOXjOdaMdAStvGyLjK+2w3
C4XeYEtFxXE9ctSP9OFGP8Ee4XltZDiIoIJ5vYBKCpqFrWsb3RwDjkZGdE8fW5yJ
mHSI3f3HgP9vehtRfFmjZy0TItS8CQIDAQABMA0GCSqGSIb3DQEBBQUAA4IBAQC6
SKqr5r8DlQl1K7BV9+j0iKp7vr19LQVoQQzcOgl7sBLlNJvQaZCvOmHf3ES3UhHZ
4HqyFy8hRRi0sRluOSSdrBfOZYDgXdB3Hv+i9JBHrWKUtr9YG8a7z1VqCcqNKvrK
TLLvw6YLJ7c4kVuQb2sFyMPS9j0RRQ/B7VoIWO5VsO0QU+CzICFCpwxcUKapF8YL
OdPKuHUP9DhNXDJHSM5j+QWi8u9K9QVzKWi3g9XdKYmmjghtjlNoHgDeBmF7dMQJ
5D8ZG7DQ2ZvQBGFVzZ9nf/JcaOO8IOBdXGFQg8sCHN5KgZaGt2GU8vxQPgzMklha
0YJdD5hHuWfxG5N4o0S2
-----END CERTIFICATE-----
EOF
}

data "aviatrix_dcf_trustbundle" "test" {
	depends_on   = [aviatrix_dcf_trustbundle.test_trustbundle]
	display_name = "test-trustbundle-data-source"
}
	`
}

func testAccDataSourceAviatrixDcfTrustbundle(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
