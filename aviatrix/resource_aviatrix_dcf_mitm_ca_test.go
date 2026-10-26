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

func TestAccAviatrixDCFMitmCa_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_MITM_CA")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF MITM CA test as SKIP_DCF_MITM_CA is set")
	}
	resourceName := "aviatrix_dcf_mitm_ca.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFMitmCaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFMitmCaBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFMitmCaExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-mitm-ca"),
					resource.TestCheckResourceAttrSet(resourceName, "ca_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ca_hash"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "state", "inactive"),
					resource.TestCheckResourceAttrSet(resourceName, "origin"),
				),
			},
		},
	})
}

func TestAccAviatrixDCFMitmCa_update(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_MITM_CA")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF MITM CA test as SKIP_DCF_MITM_CA is set")
	}
	resourceName := "aviatrix_dcf_mitm_ca.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFMitmCaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFMitmCaBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFMitmCaExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-mitm-ca"),
				),
			},
			{
				Config: testAccDCFMitmCaUpdatedName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCFMitmCaExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-mitm-ca-updated"),
				),
			},
		},
	})
}

func TestAccAviatrixDCFMitmCa_invalidCertificate(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_MITM_CA")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF MITM CA test as SKIP_DCF_MITM_CA is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFMitmCaDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDCFMitmCaInvalidCertificate(),
				ExpectError: regexp.MustCompile("no certificates found in bundle"),
			},
		},
	})
}

func TestAccAviatrixDCFMitmCa_nonCaCertificate(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_MITM_CA")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF MITM CA test as SKIP_DCF_MITM_CA is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDCFMitmCaDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDCFMitmCaNonCaCertificate(),
				ExpectError: regexp.MustCompile("is not a CA"),
			},
		},
	})
}

func testAccDCFMitmCaBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_mitm_ca" "test" {
  name              = "test-dcf-mitm-ca"
  key               = %q
  certificate_chain = %q
}
`, testMitmCaPrivateKey(), testMitmCaCertificate())
}

func testAccDCFMitmCaUpdatedName() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_mitm_ca" "test" {
  name              = "test-dcf-mitm-ca-updated"
  key               = %q
  certificate_chain = %q
}
`, testMitmCaPrivateKey(), testMitmCaCertificate())
}

func testAccDCFMitmCaInvalidCertificate() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_mitm_ca" "test" {
  name              = "test-dcf-mitm-ca-invalid"
  key               = %q
  certificate_chain = %q
}
`, testMitmCaPrivateKey(), testInvalidCertificateContent())
}

func testAccDCFMitmCaNonCaCertificate() string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_mitm_ca" "test" {
  name              = "test-dcf-mitm-ca-non-ca"
  key               = %q
  certificate_chain = %q
}
`, testMitmCaPrivateKey(), testNonCaCertificate())
}

// testMitmCaPrivateKey returns a test RSA private key in PEM format
// This is a test key and should only be used for testing purposes
func testMitmCaPrivateKey() string {
	return `-----BEGIN PRIVATE KEY-----
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
I0wvJbHgLn0PILUwUylRQ0bp7iNk8Q==
-----END PRIVATE KEY-----`
}

// testMitmCaCertificate returns a test CA certificate in PEM format
// This is a self-signed CA certificate for testing purposes only
func testMitmCaCertificate() string {
	return `-----BEGIN CERTIFICATE-----
MIIFJTCCAw2gAwIBAgIUV0+9A4jifR/BtURx0/JAg3jb1CwwDQYJKoZIhvcNAQEL
BQAwGjEYMBYGA1UEAwwPYXZ4LW5ldy1iZXRhLWNhMB4XDTI2MDIwMjEzNTcyOFoX
DTI3MDIwMjEzNTcyOFowGjEYMBYGA1UEAwwPYXZ4LW5ldy1iZXRhLWNhMIICIjAN
BgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA8f3KXl+hvJ+/InpVoWmBw6NZdVWt
M83rbEH5Wms6oFZOprC77TPpnGp1vv3JWTPd/X1G49xb3vpyzuByJ+sfIfI0XySz
JMbaSnts1xkmzs1jTuBLECA7w0EN4okKgTiqgEBdh8/Pws9IyX9K6oDppBO57h/j
hzn9UIxs8zbFmWUmMw1FE/H1TcPGBW8vx5x8tarHZz7stOYCUb8AyaEhecBZ4r7B
w2IHMCWKbId6xtV9zjSouKRbFzYAFHtq1sMjm3lo04Qo6gLzIWGtSkfdnPvfSzz/
6I+WEekxIZcWfq5WRd4SBQbaaQICZhhPvRXRagG4q96a/QLi6isFy8AH5g1wz8sb
UvgjT4SNDA9mikxID4wPwuEfaxxCDUB+8Kdut5dYFHnHblAMiezlUrJ75/E8p3CF
959QPym7viVmWqCzNtCKkF14smRglqm1l7LKEEbIJLigDf1p3gv4YH/dmCdJUdWd
ZqYiwxkbOtwEoeCj+pNcX6BmKQoOgUX8ZpOHpzfWsp5SDLbhORYjpR0d0/Kn7jI4
YTeU5erXrRwfOa5wZQMKbPQ74F7h9FEwqYd/vzT3bLlor50+BqitpdZEUdIfPXw3
0MMAa0jMXxif/GidNpcqHQ02Ee+TKLJ2WDWVHnDPQm3I1CI7btV1lqks4Ycg1A07
YJVa2CNkjbxjebUCAwEAAaNjMGEwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8E
BAMCAQYwHQYDVR0OBBYEFIrWd9p5L28NJFWnLxiuFbz23LI5MB8GA1UdIwQYMBaA
FIrWd9p5L28NJFWnLxiuFbz23LI5MA0GCSqGSIb3DQEBCwUAA4ICAQDvgMmZQyiZ
C5ZEkhqEKTivTlCXfMt+9qlT8x9x0Rnzzt5scwdT8fQ80DMEc20ODv6OuYW677oA
S9fpT66VHUjljB3ipHZnB1/iIjvjzExWLYDBMzugdIGonEgIn11Kg6k5FCrOmmfv
wBBbSerJryv443zdS84VsPe4y8Bv3Nnhge0hq/3qrK4IgX8AEZGyPs0Uo+cl5/ql
Vv5PX4cscvZ5ySlJ0elLFshLy/HphlDG98fIO1A3R0H3uzTjxsYB00eLAIU+Lx5W
niYaIZUbuwqYKVCKFeHO1jII30vI8EgHrhGYsMiERqc6BnJeSvjROb/t2HIPu4M+
E/371vT+8m58PWvAqvz9TuqWn9VxHp8ynmFxfZS8Fxo+FxlmAsXoZXAwKRJ102sT
lkv2rFaBnVL129fjPG48wsiMgmLQ6hVTcQzK7z5GGkBAkTU6AXpl2nO9JUK+epkh
4EjETvHCoT8C0WYoME+24MQA4P20KBUIuhEhMHIpIGRNXH3rYBD6ptaIcmaocqVB
Yt4Ebyc5FDLSIHe1/7DhMpAeu98n8Neig9Gv/AO9vOVKcMk891wOytKOA+he/h/T
i8sKG5gsMn1YMO8taIKsfeariz956+YRBxyntFzV3sD4GVcAc3zV0ykfxkIeHd9x
9s13Pw3eUfNHKhn9EdmBUNJyZPVDZyfc9w==
-----END CERTIFICATE-----`
}

// testNonCaCertificate returns a certificate that is not a CA certificate
// Used to test validation that rejects non-CA certificates
func testNonCaCertificate() string {
	// This is an end-entity certificate (not a CA) for testing validation
	return `-----BEGIN CERTIFICATE-----
MIIDETCCAfmgAwIBAgIUV5P+W0rYm5kF6a5w3UjPq4YdW7QwDQYJKoZIhvcNAQEL
BQAwGDEWMBQGA1UEAwwNdGVzdC1lbmQtY2VydDAeFw0yNDAxMDEwMDAwMDBaFw0y
NTAxMDEwMDAwMDBaMBgxFjAUBgNVBAMMDXRlc3QtZW5kLWNlcnQwggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQDRndVLkklx2zfF+f/KBbIXw9ucbLQAcHsz
GzNHRiM4Z8vdtCKTdjYsoZLRJnEbBtiEEQ1UEv4Upv7x6Pxv+Zevw0xEmIzNOc4r
yrxxmuQSAHUHN9rAuHZLqU8/qQeODz+uQI4E/TAuoNBicMsYHMs5ousDl0iUIqQf
Nm8E3Ryg7VeY8IQioBHwGL6gQkGZTysl9+lNskyZgp1Jg5y4murQVDtaOnAwQSzE
anZVZxL6QZsUoaM+YC3uWygu8lohcvd4X4MX1tD2E2qzbvIgkIlMisvAJ7mnkoSw
3NfoVAoxKLTIRpKwUasu+MonIKBMtIlOH1Su6FcvaX8HxDGHAokBAgMBAAGjUzBR
MB0GA1UdDgQWBBQfq/N26dbZFCwrvUu3uryUr7w75jAfBgNVHSMEGDAWgBQfq/N2
6dbZFCwrvUu3uryUr7w75jAPBgNVHRMBAf8EBTADAQEAMA0GCSqGSIb3DQEBCwUA
A4IBAQBWFkHvNcN1qCLWqGH8E8mJnN1nPzWqLBgZYXWxMHFyI6hF8ldLdlNdl0Jx
q7aAKSm1m3UL25lX0yA7U5IJSc+kQGk7q0Q+RN0Q3R5MXl0S8kcw6zEr7L4mKkL3
VjdGGF9XDl0sGBvP1vP7iJQ7U8x0bCrwqKnILqvLJdWWN5GmhPdADq0PwLt2y0Yx
qSeJx8k7hFCjB3qH4wpwqJnN3LLnGvnMuNjbFdRcLzNcNQ8EGzHba2NfwGO1nG2p
k1kYpRJwz7hO3cZ80MvNgnH+QY8OY8N1QKLTPUD2j3l1cM0prSz9tcZKS3MFwX7j
Z8MoMN0K8zU0l8S0a0v5lV7pcl3E
-----END CERTIFICATE-----`
}

func testAccCheckDCFMitmCaExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF MITM CA resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF MITM CA ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		mitmCa, err := client.GetDCFMitmCa(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF MITM CA status: %w", err)
		}

		if mitmCa.CaID != rs.Primary.ID {
			return fmt.Errorf("DCF MITM CA ID not found")
		}

		return nil
	}
}

func testAccCheckDCFMitmCaDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_mitm_ca" {
			continue
		}

		_, err := client.GetDCFMitmCa(context.Background(), rs.Primary.ID)
		if err == nil || !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("DCF MITM CA still exists when it should be destroyed")
		}
	}

	return nil
}
