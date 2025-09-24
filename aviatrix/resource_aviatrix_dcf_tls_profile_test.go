package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixDcfTLSProfile_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_TLS_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF TLS Profile test as SKIP_DCF_TLS_PROFILE is set")
	}
	resourceName := "aviatrix_dcf_tls_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfTLSProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfTLSProfileBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfTLSProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-dcf-tls-profile"),
					resource.TestCheckResourceAttr(resourceName, "certificate_validation", "CERTIFICATE_VALIDATION_LOG_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "verify_sni", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				Config: testAccCheckDcfTLSProfileUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfTLSProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-dcf-tls-profile-updated"),
					resource.TestCheckResourceAttr(resourceName, "certificate_validation", "CERTIFICATE_VALIDATION_ENFORCE"),
					resource.TestCheckResourceAttr(resourceName, "verify_sni", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
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

func TestAccAviatrixDcfTLSProfile_withCABundle(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_TLS_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF TLS Profile test as SKIP_DCF_TLS_PROFILE is set")
	}
	resourceName := "aviatrix_dcf_tls_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfTLSProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfTLSProfileWithCABundle(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfTLSProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-dcf-tls-profile-with-ca"),
					resource.TestCheckResourceAttr(resourceName, "certificate_validation", "CERTIFICATE_VALIDATION_ENFORCE"),
					resource.TestCheckResourceAttr(resourceName, "verify_sni", "true"),
					resource.TestCheckResourceAttr(resourceName, "ca_bundle_id", "test-ca-bundle-uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
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

func testAccCheckDcfTLSProfileBasic() string {
	return `
resource "aviatrix_dcf_tls_profile" "test" {
	display_name           = "test-dcf-tls-profile"
	certificate_validation = "CERTIFICATE_VALIDATION_LOG_ONLY"
	verify_sni            = true
}
`
}

func testAccCheckDcfTLSProfileUpdate() string {
	return `
resource "aviatrix_dcf_tls_profile" "test" {
	display_name           = "test-dcf-tls-profile-updated"
	certificate_validation = "CERTIFICATE_VALIDATION_ENFORCE"
	verify_sni            = false
}
`
}

func testAccCheckDcfTLSProfileWithCABundle() string {
	return `
resource "aviatrix_dcf_tls_profile" "test" {
	display_name           = "test-dcf-tls-profile-with-ca"
	certificate_validation = "CERTIFICATE_VALIDATION_ENFORCE"
	verify_sni            = true
	ca_bundle_id          = "test-ca-bundle-uuid"
}
`
}

func testAccCheckDcfTLSProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF TLS Profile resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF TLS Profile ID is set")
		}

		client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
		}

		_, err := client.GetTLSProfile(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF TLS Profile status: %w", err)
		}

		return nil
	}
}

func testAccCheckDcfTLSProfileDestroy(s *terraform.State) error {
	client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_tls_profile" {
			continue
		}

		_, err := client.GetTLSProfile(context.Background(), rs.Primary.ID)
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("dcf tls profile configured when it should be destroyed %w", err)
		}
	}

	return nil
}
