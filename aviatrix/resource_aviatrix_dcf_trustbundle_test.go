package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDCFTrustBundle_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_TRUSTBUNDLE")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Trust Bundle test as SKIP_DCF_TRUSTBUNDLE is set")
	}
	resourceName := "aviatrix_dcf_trustbundle.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFTrustBundleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFTrustBundleBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFTrustBundleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-dcf-trustbundle"),
					resource.TestCheckResourceAttrSet(resourceName, "bundle_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "bundle_content", testCertificateContent()),
				),
			},
		},
	})
}

func TestAccAviatrixDCFTrustBundle_update(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_TRUSTBUNDLE")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Trust Bundle test as SKIP_DCF_TRUSTBUNDLE is set")
	}
	resourceName := "aviatrix_dcf_trustbundle.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFTrustBundleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFTrustBundleBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFTrustBundleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-dcf-trustbundle"),
				),
			},
			{
				Config: testAccDCFTrustBundleUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFTrustBundleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-dcf-trustbundle-updated"),
				),
			},
		},
	})
}

func TestAccAviatrixDCFTrustBundle_invalidCertificate(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_TRUSTBUNDLE")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Trust Bundle test as SKIP_DCF_TRUSTBUNDLE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFTrustBundleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDCFTrustBundleInvalidCertificate(),
				ExpectError: regexp.MustCompile("no certificates found in bundle"),
			},
		},
	})
}

func testAccDCFTrustBundleBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_trustbundle" "test" {
  display_name   = "test-dcf-trustbundle"
  bundle_content = %q
}
`, testCertificateContent())
}

func testAccDCFTrustBundleUpdated() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_trustbundle" "test" {
  display_name   = "test-dcf-trustbundle-updated"
  bundle_content = %q
}
`, testCertificateContent())
}

func testAccDCFTrustBundleInvalidCertificate() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_trustbundle" "test" {
  display_name   = "test-dcf-trustbundle-invalid"
  bundle_content = %q
}
`, testInvalidCertificateContent())
}

func testCertificateContent() string {
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

func testInvalidCertificateContent() string {
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
FQfVdG2q7fM3lGeyx/HFfaOvgYMi`
}

func testAccCheckDCFTrustBundleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF Trust Bundle resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF Trust Bundle ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		trustBundle, err := client.GetDCFTrustBundleByID(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF Trust Bundle status: %w", err)
		}

		if trustBundle.BundleID != rs.Primary.ID {
			return fmt.Errorf("DCF Trust Bundle ID not found")
		}

		return nil
	}
}

func testAccCheckDCFTrustBundleDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_trustbundle" {
			continue
		}

		_, err := client.GetDCFTrustBundleByID(context.Background(), rs.Primary.ID)
		if err == nil || !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("DCF Trust Bundle still exists when it should be destroyed")
		}
	}

	return nil
}
