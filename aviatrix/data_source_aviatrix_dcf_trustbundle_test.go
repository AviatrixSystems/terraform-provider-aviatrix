package aviatrix

import (
	"fmt"
	"os"
	"strings"
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
	return fmt.Sprintf(`
resource "aviatrix_dcf_trustbundle" "test_trustbundle" {
	display_name    = "test-trustbundle-data-source"
	bundle_content = "%s"
}

data "aviatrix_dcf_trustbundle" "test" {
	depends_on   = [aviatrix_dcf_trustbundle.test_trustbundle]
	display_name = "test-trustbundle-data-source"
}
`, strings.ReplaceAll(testAccDataSourceCertificateContent(), "\n", "\\n"))
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

func testAccDataSourceCertificateContent() string {
	return `-----BEGIN CERTIFICATE-----
MIIDQTCCAimgAwIBAgITBmyfz5m/jAo54vB4ikPmljZbyjANBgkqhkiG9w0BAQsF
ADA5MQswCQYDVQQGEwJVUzEPMA0GA1UEChMGQW1hem9uMRkwFwYDVQQDExBBbWF6
b24gUm9vdCBDQSAxMB4XDTE1MDUyNjAwMDAwMFoXDTM4MDExNzAwMDAwMFowOTEL
MAkGA1UEBhMCVVMxDzANBgNVBAoTBkFtYXpvbjEZMBcGA1UEAxMQQW1hem9uIFJv
b3QgQ0EgMTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALJ4gHHKeNXj
ca9HgFB0fW7Y14h29Jlo91ghYPl0hAEvrAIthtOgQ3pOsqTQNroBvo3bSMgHFzZM
9O6II8c+6zf1tRn4SWiw3te5djgdYZ6k/oI2peVKVuRF4fn9tBb6dNqcmzU5L/qw
IFAGbHrQgLKm+a/sRxmPUDgH3KKHOVj4utWp+UhnMJbulHheb4mjUcAwhmahRWa6
VOujw5H5SNz/0egwLX0tdHA114gk957EWW67c4cX8jJGKLhD+rcdqsq08p8kDi1L
93FcXmn/6pUCyziKrlA4b9v7LWIbxcceVOF34GfID5yHI9Y/QCB/IIDEgEw+OyQm
jgSubJrIqg0CAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMC
AYYwHQYDVR0OBBYEFIQYzIU07LwMlJQuCFmcx7IQTgoIMA0GCSqGSIb3DQEBCwUA
A4IBAQCY8jdaQZChGsV2USggNiMOruYou6r4lK5IpDB/G/wkjUu0yKGX9rbxenDI
U5PMCCjjmCXPI6T53iHTfIuJruydjsw2hUwsqdnHnFx9k6Tpdp4xvN0dWQVmIUgX
tc9RiOQTUM8IzG2wDz1oydw+RVF/TmRQ6EQfoQJynfKzKkCzR4LGLd4IySQsv0GB
CbFy9K3VRIs57/m3NY+8R4Z8qFJutMSlV+gYBbXUz/+ibZb5l6j9jCFZ5CNczKx8
iTiYXZ68GCDImLLJqTgCp8SysbyMVGLWwUNzbyBqEjxHqGB/Kryl9SEgvQS0hrLN
FQfVdG2q7fM3lGeyx/HFfaOvgYMi
-----END CERTIFICATE-----`
}
