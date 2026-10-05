package aviatrix

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixConfigFeature_InvalidFeatureNameValidation(t *testing.T) {
	if os.Getenv("SKIP_CONFIG_FEATURE") == "yes" {
		t.Skip("Skipping config feature acceptance tests as SKIP_CONFIG_FEATURE is set")
	}
	featureName := "invalid_feature"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConfigFeatureBasic(featureName, true),
				ExpectError: regexp.MustCompile("invalid feature name: invalid_feature"),
			},
		},
	})
}

func TestAccAviatrixConfigFeature_toggle(t *testing.T) {
	if os.Getenv("SKIP_CONFIG_FEATURE") == "yes" {
		t.Skip("Skipping config feature acceptance tests as SKIP_CONFIG_FEATURE is set")
	}
	featureName := "microseg"
	resourceName := "aviatrix_config_feature.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccConfigFeatureDestroy(featureName),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigFeatureBasic(featureName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccConfigFeatureCheckStatus(resourceName, featureName, true),
					resource.TestCheckResourceAttr(resourceName, "feature_name", featureName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
				),
			},
			{
				Config: testAccConfigFeatureBasic(featureName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccConfigFeatureCheckStatus(resourceName, featureName, false),
					resource.TestCheckResourceAttr(resourceName, "feature_name", featureName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
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

func testAccConfigFeatureBasic(featureName string, enabled bool) string {
	return fmt.Sprintf(`
resource "aviatrix_config_feature" "test" {
	feature_name = %q
	is_enabled   = %t
}
`, featureName, enabled)
}

func testAccConfigFeatureCheckStatus(n, featureName string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("config feature not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no config feature ID is set")
		}

		if rs.Primary.ID != featureName {
			return fmt.Errorf("unexpected config feature ID: got %s, expected %s", rs.Primary.ID, featureName)
		}

		client := mustClient(testAccProvider.Meta())
		status, err := client.GetFeatureStatus(context.Background(), featureName)
		if err != nil {
			return fmt.Errorf("failed to get feature status for %q: %w", featureName, err)
		}
		if status.Enabled != expected {
			return fmt.Errorf("unexpected feature status for %q: got enabled=%t, expected enabled=%t", featureName, status.Enabled, expected)
		}

		return nil
	}
}

func testAccConfigFeatureDestroy(featureName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := mustClient(testAccProvider.Meta())
		status, err := client.GetFeatureStatus(context.Background(), featureName)
		if err != nil {
			return fmt.Errorf("failed to get feature status for %q during destroy check: %w", featureName, err)
		}
		if status.Enabled {
			return fmt.Errorf("feature %q is still enabled after destroy", featureName)
		}
		return nil
	}
}

func TestAccAviatrixConfigFeature_DocFeatureListMatchesAPI(t *testing.T) {
	if os.Getenv("SKIP_CONFIG_FEATURE") == "yes" {
		t.Skip("Skipping config feature acceptance tests as SKIP_CONFIG_FEATURE is set")
	}
	// Make sure to always update the list in this test when adding a new feature name to the docs.
	currentListInDocs := []string{"microseg", "cost_iq", "ipv6", "dcf_on_s2c", "dcf_on_psf", "dcf_stats_obs_sink", "dcf_logs_obs_sink", "k8s", "sre_metrics_export", "k8s_dcf_policies", "dcf_on_firenet", "interface_mtu_based_clamping", "primary_gateway_deletion", "vrf"}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: `data "aviatrix_caller_identity" "test" {}`,
				Check: func(s *terraform.State) error {
					client := mustClient(testAccProviderVersionValidation.Meta())
					apiFeatures, err := client.GetAllFeatureNames(context.Background())
					if err != nil {
						return err
					}
					filtered := slices.DeleteFunc(slices.Clone(apiFeatures), func(f string) bool {
						return slices.Contains(goaviatrix.FeatureNameExceptions, f)
					})
					assert.ElementsMatch(t, currentListInDocs, filtered, "feature list in docs does not match feature list in API, please update the docs to match the API and update the list in this test to match the docs")
					return nil
				},
			},
		},
	})
}
