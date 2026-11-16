package aviatrix

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
					testAccCheckDcfTLSProfileExists(resourceName, t),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-dcf-tls-profile"),
					resource.TestCheckResourceAttr(resourceName, "certificate_validation", "CERTIFICATE_VALIDATION_LOG_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "verify_sni", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				Config: testAccCheckDcfTLSProfileUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfTLSProfileExists(resourceName, t),
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
	caBundleID := uuid.New().String()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfTLSProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfTLSProfileWithCABundle(caBundleID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfTLSProfileExists(resourceName, t),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test-dcf-tls-profile-with-ca"),
					resource.TestCheckResourceAttr(resourceName, "certificate_validation", "CERTIFICATE_VALIDATION_ENFORCE"),
					resource.TestCheckResourceAttr(resourceName, "verify_sni", "true"),
					resource.TestCheckResourceAttr(resourceName, "ca_bundle_id", caBundleID),
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

func TestAccAviatrixDcfTLSProfile_tlsLogFields(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_TLS_PROFILE")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF TLS Profile test as SKIP_DCF_TLS_PROFILE is set")
	}
	resourceName := "aviatrix_dcf_tls_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return testAccProviderVersionValidation, nil },
		},
		CheckDestroy: testAccCheckDcfTLSProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDcfTLSProfileLogFields(`["TLS_LOG_FIELDS_ALL", "TLS_CLIENT_SNI"]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("TLS_LOG_FIELDS_ALL must be the only value"),
			},
			{
				Config: testAccCheckDcfTLSProfileLogFields(`["TLS_CLIENT_SNI", "TLS_CLIENT_JA4_FINGERPRINT"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfTLSProfileExists(resourceName, t),
					resource.TestCheckResourceAttr(resourceName, "tls_log_fields_extension.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tls_log_fields_extension.*", "TLS_CLIENT_SNI"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tls_log_fields_extension.*", "TLS_CLIENT_JA4_FINGERPRINT"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCheckDcfTLSProfileLogFields(`["TLS_LOG_FIELDS_ALL"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tls_log_fields_extension.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tls_log_fields_extension.*", "TLS_LOG_FIELDS_ALL"),
				),
			},
			{
				Config: testAccCheckDcfTLSProfileBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tls_log_fields_extension.#", "0"),
				),
			},
		},
	})
}

func TestMarshalDCFTLSProfileInput_TLSLogFields(t *testing.T) {
	tests := []struct {
		name   string
		fields []any
		want   []string
	}{
		{name: "omitted", fields: nil, want: nil},
		{name: "subset", fields: []any{"TLS_CLIENT_SNI"}, want: []string{"TLS_CLIENT_SNI"}},
		{name: "all shorthand", fields: []any{"TLS_LOG_FIELDS_ALL"}, want: allTLSLogFields},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := map[string]any{
				"display_name":           "p",
				"certificate_validation": "CERTIFICATE_VALIDATION_NONE",
				"verify_sni":             true,
			}
			if tt.fields != nil {
				raw["tls_log_fields_extension"] = tt.fields
			}
			d := schema.TestResourceDataRaw(t, resourceAviatrixDCFTLSProfile().Schema, raw)

			profile, err := marshalDCFTLSProfileInput(d)
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, profile.TLSLogFieldsExtension)
		})
	}
}

func TestMarshalDCFTLSProfileInput_OmitsEmptyTLSLogFieldsKey(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAviatrixDCFTLSProfile().Schema, map[string]any{
		"display_name":           "p",
		"certificate_validation": "CERTIFICATE_VALIDATION_NONE",
		"verify_sni":             true,
	})

	profile, err := marshalDCFTLSProfileInput(d)
	require.NoError(t, err)
	body, err := json.Marshal(profile)
	require.NoError(t, err)
	assert.NotContains(t, string(body), "tls_log_fields_extension")
}

func TestMarshalDCFTLSProfileInput_RejectsAllWithOtherFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAviatrixDCFTLSProfile().Schema, map[string]any{
		"display_name":             "p",
		"certificate_validation":   "CERTIFICATE_VALIDATION_NONE",
		"verify_sni":               true,
		"tls_log_fields_extension": []any{"TLS_LOG_FIELDS_ALL", "TLS_CLIENT_SNI"},
	})

	_, err := marshalDCFTLSProfileInput(d)
	require.ErrorContains(t, err, "TLS_LOG_FIELDS_ALL must be the only value")
}

func TestValidateTLSLogFields(t *testing.T) {
	require.NoError(t, validateTLSLogFields(nil))
	require.NoError(t, validateTLSLogFields([]string{"TLS_LOG_FIELDS_ALL"}))
	require.NoError(t, validateTLSLogFields([]string{"TLS_CLIENT_SNI", "TLS_CLIENT_ALPNS"}))
	require.Error(t, validateTLSLogFields([]string{"TLS_LOG_FIELDS_ALL", "TLS_CLIENT_SNI"}))
}

func TestFlattenTLSLogFields(t *testing.T) {
	reversedAll := slices.Clone(allTLSLogFields)
	slices.Reverse(reversedAll)

	tests := []struct {
		name  string
		api   []string
		prior []string
		want  []string
	}{
		{name: "all with prior shorthand collapses", api: reversedAll, prior: []string{"TLS_LOG_FIELDS_ALL"}, want: []string{"TLS_LOG_FIELDS_ALL"}},
		{name: "all without prior shorthand stays expanded", api: allTLSLogFields, prior: nil, want: allTLSLogFields},
		{name: "all with prior explicit list stays expanded", api: allTLSLogFields, prior: allTLSLogFields, want: allTLSLogFields},
		{name: "subset with prior shorthand shows drift", api: []string{"TLS_CLIENT_SNI"}, prior: []string{"TLS_LOG_FIELDS_ALL"}, want: []string{"TLS_CLIENT_SNI"}},
		{name: "empty", api: nil, prior: []string{"TLS_LOG_FIELDS_ALL"}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, flattenTLSLogFields(tt.api, tt.prior))
		})
	}
}

func TestResourceAviatrixDCFTLSProfileRead_TLSLogFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAviatrixDCFTLSProfile().Schema, map[string]any{
		"tls_log_fields_extension": []any{"TLS_LOG_FIELDS_ALL"},
	})
	mustSet(d, "tls_log_fields_extension", flattenTLSLogFields(allTLSLogFields, getStringSet(d, "tls_log_fields_extension")))
	assert.Equal(t, []string{"TLS_LOG_FIELDS_ALL"}, getStringSet(d, "tls_log_fields_extension"))

	mustSet(d, "tls_log_fields_extension", flattenTLSLogFields(nil, getStringSet(d, "tls_log_fields_extension")))
	assert.Empty(t, getStringSet(d, "tls_log_fields_extension"))
}

func testAccCheckDcfTLSProfileLogFields(fields string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_tls_profile" "test" {
	display_name             = "test-dcf-tls-profile"
	certificate_validation   = "CERTIFICATE_VALIDATION_LOG_ONLY"
	verify_sni               = true
	tls_log_fields_extension = %s
}
`, fields)
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

func testAccCheckDcfTLSProfileWithCABundle(caBundleID string) string {
	return fmt.Sprintf(`
resource "aviatrix_dcf_tls_profile" "test" {
	display_name           = "test-dcf-tls-profile-with-ca"
	certificate_validation = "CERTIFICATE_VALIDATION_ENFORCE"
	verify_sni            = true
	ca_bundle_id          =  "%s"
}
`, caBundleID)
}

func testAccCheckDcfTLSProfileExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF TLS Profile resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF TLS Profile ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		_, err := client.GetTLSProfile(t.Context(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF TLS Profile status: %w", err)
		}

		return nil
	}
}

func testAccCheckDcfTLSProfileDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_tls_profile" {
			continue
		}
		_, err := client.GetTLSProfile(ctx, rs.Primary.ID)
		if err == nil || !strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("dcf tls profile configured when it should be destroyed %w", err)
		}
	}

	return nil
}
